package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/micromdm/go4/version"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	var run func([]string) error
	switch strings.ToLower(os.Args[1]) {
	case "version", "-version":
		version.Print()
		return
	case "serve":
		run = serve
	default:
		usage()
		os.Exit(1)
	}

	if err := run(os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func usage() error {
	helpText := `USAGE: micromdm <COMMAND>

Available Commands:
	serve
	version

Use micromdm <command> -h for additional usage of each command.
Example: micromdm serve -h
`
	fmt.Println(helpText)
	return nil
}

func usageFor(fs *flag.FlagSet, short string) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "USAGE\n")
		fmt.Fprintf(os.Stderr, "  %s\n", short)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "FLAGS\n")
		w := tabwriter.NewWriter(os.Stderr, 0, 2, 2, ' ', 0)
		fs.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(w, "\t-%s %s\t%s\n", f.Name, f.DefValue, f.Usage)
		})
		w.Flush()
		fmt.Fprintf(os.Stderr, "\n")
	}
}
