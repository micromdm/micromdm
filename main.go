package main

import (
	"fmt"
	stdlog "log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/micromdm/micromdm/app"
	"github.com/micromdm/micromdm/cli"
)

func main() {
	if isCLI() { // use CLI if the user typed a subcommand.
		cli.Main()
		return
	}

	logger := log.NewLogfmtLogger(os.Stderr)
	stdlog.SetOutput(log.NewStdlibAdapter(logger)) // force structured logs
	mainLogger := log.NewContext(logger).With("component", "main")

	// the default action for no subcommand is to launch the server.
	errs := make(chan error, 2)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()
	go func() {
		status, err := app.Main(logger)
		if err != nil {
			mainLogger.Log("terminated", err)
		}
		os.Exit(status)
	}()
	mainLogger.Log("terminated", <-errs)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// isCLI determines if this is the default app or if there are subcommands.
// ./micromdm or ./micromdm -foo would return false, while ./micromdm version
// would return true
func isCLI() bool {
	switch len(os.Args) {
	case 1:
		return false
	default:
		if strings.HasPrefix(os.Args[1], "-") {
			return false
		} else {
			return true
		}
	}
}
