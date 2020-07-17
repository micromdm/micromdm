package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/oklog/run"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"

	"micromdm.io/v2/pkg/log"
	"micromdm.io/v2/pkg/version"
)

// write a pid so that the server can be restarted with SIGHUP
func writePID(path string) error {
	if err := ioutil.WriteFile(path, []byte(strconv.Itoa(os.Getpid())), 0600); err != nil {
		return fmt.Errorf("writing pidfile: %w", err)
	}
	return nil
}

type cliFlags struct {
	siteName string
	http     string
}

func micromdm(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	var (
		ctx     = context.Background()
		logger  = log.New(log.Output(stderr))
		cli     = &cliFlags{}
		rootfs  = flag.NewFlagSet("micromdm", flag.ContinueOnError)
		pidfile = rootfs.String("pidfile", "/tmp/micromdm.pid", "Path to server pidfile")
		_       = rootfs.String("config", "", "Path to config file (optional)")
	)

	rootfs.StringVar(&cli.siteName, "site_name", "Acme", "Name of the site as it would appear in the top left of the HTML UI")
	rootfs.StringVar(&cli.http, "http", "localhost:9000", "HTTP service address")

	// default output is os.Stderr.
	// setting the output and flag.ContinueOnError overrides allows testing usage.
	rootfs.SetOutput(stderr)

	version := &ffcli.Command{
		Name:       "version",
		ShortUsage: "version [<arg> ...]",
		ShortHelp:  "Print version information.",
		Exec: func(_ context.Context, args []string) error {
			version.PrintFull()
			return nil
		},
	}

	// add a help subcommand to make usage more discoverable.
	helpCmd := &ffcli.Command{
		Name:      "help",
		ShortHelp: "Print this help text.",
		UsageFunc: func(c *ffcli.Command) string { return "" },
		Exec: func(_ context.Context, args []string) error {
			rootfs.Usage()
			return flag.ErrHelp
		},
	}

	root := &ffcli.Command{
		ShortUsage:  "micromdm [flags] <subcommand>",
		FlagSet:     rootfs,
		Options:     []ff.Option{ff.WithEnvVarPrefix("MICROMDM"), ff.WithConfigFileParser(ff.PlainParser), ff.WithConfigFileFlag("config")},
		Subcommands: []*ffcli.Command{helpCmd, version},
		Exec: func(context.Context, []string) error {
			if err := writePID(*pidfile); err != nil {
				return err
			}

			srv, err := setup(cli, logger)
			if err != nil {
				return err
			}

			// run.Group manages lifecycles of various long running goroutines:
			// - signal handlers for SIGTERM/SIGHUP etc.
			// - http.Server listeners.
			var g run.Group
			{
				server := &http.Server{
					Handler: srv.ui.Handler(),
					Addr:    cli.http,
				}

				g.Add(func() error {
					log.Info(logger).Log("component", "frontend", "msg", "started")
					return server.ListenAndServe()
				}, func(error) {
					ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
					defer cancel()
					server.Shutdown(ctx)
				})
			}
			{
				// when the binary receives SIGINT or SIGTERM, execution is cancelled
				ctx, cancel := context.WithCancel(ctx)
				g.Add(func() error {
					c := make(chan os.Signal, 1)
					signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
					select {
					case <-ctx.Done():
						return ctx.Err()
					case sig := <-c:
						return fmt.Errorf("received signal %s", sig)
					}
				}, func(error) {
					os.Remove(*pidfile)
					cancel()
				})
			}

			{
				// restart the process after SIGHUP. Mainly used for development,
				// restarting for config changes/html template reloading.
				ctx, cancel := context.WithCancel(ctx)
				g.Add(func() error {
					c := make(chan os.Signal, 1)
					signal.Notify(c, syscall.SIGHUP)
					for {
						select {
						case <-ctx.Done():
							return ctx.Err()
						case sig := <-c:
							log.Info(logger).Log("msg", "restarting process", "signal", sig.String())
							syscall.Exec(args[0], args, os.Environ())
						}
					}
				}, func(error) {
					cancel()
				})
			}

			return g.Run()
		},
	}

	switch err := root.ParseAndRun(ctx, args[1:]); {
	case err == nil:
		return 0
	case errors.Is(err, flag.ErrHelp):
		return 2
	default:
		log.Info(logger).Log("exit", err)
		return 1
	}
}

func main() { os.Exit(micromdm(os.Args, os.Stdin, os.Stdout, os.Stderr)) }
