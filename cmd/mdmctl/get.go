package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/go-kit/kit/log"
	"github.com/micromdm/micromdm/core/list"
)

type getCommand struct {
	config *ClientConfig
	list   list.Service
}

func (cmd *getCommand) setup() error {
	cfg, err := NewClientConfig()
	if err != nil {
		return err
	}
	cmd.config = cfg
	logger := log.NewLogfmtLogger(os.Stderr)
	listsvc, err := list.NewClient(cfg.ServerURL, logger, cfg.APIToken)
	if err != nil {
		return err
	}
	cmd.list = listsvc
	return nil
}

func (cmd *getCommand) Run(args []string) error {
	if len(args) < 1 {
		cmd.Usage()
		os.Exit(1)
	}

	if err := cmd.setup(); err != nil {
		return err
	}

	var run func([]string) error
	switch strings.ToLower(args[0]) {
	case "devices":
		run = cmd.getDevices
	default:
		cmd.Usage()
		os.Exit(1)
	}

	return run(args[1:])
}

func (cmd *getCommand) Usage() error {
	const getUsage = `
Display one or many resources.

Valid resource types:

  * devices
  * blueprints

Examples:
  # Get a list of devices
  mdmctl get devices

  # Get a device by serial (TODO implement filtering)
  mdmctl get devices -serial=C02ABCDEF
`
	fmt.Println(getUsage)
	return nil
}

type devicesTableOutput struct{ w *tabwriter.Writer }

func (out *devicesTableOutput) BasicHeader() {
	fmt.Fprintf(out.w, "UDID\tSerialNumber\tEnrollmentStatus\tLastSeen\n")
}

func (out *devicesTableOutput) BasicFooter() {
	out.w.Flush()
}

func (cmd *getCommand) getDevices(args []string) error {
	flagset := flag.NewFlagSet("devices", flag.ExitOnError)
	flagset.Usage = usageFor(flagset, "micromdm get devices [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}
	w := tabwriter.NewWriter(os.Stderr, 0, 4, 2, ' ', 0)
	out := &devicesTableOutput{w}
	out.BasicHeader()
	defer out.BasicFooter()
	ctx := context.Background()
	devices, err := cmd.list.ListDevices(ctx, list.ListDevicesOption{})
	if err != nil {
		return err
	}
	for _, d := range devices {
		fmt.Fprintf(out.w, "%s\t%s\t%v\t%s\n", d.UDID, d.SerialNumber, d.EnrollmentStatus, d.LastSeen)
	}
	return nil
}
