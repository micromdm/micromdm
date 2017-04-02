package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/micromdm/micromdm/core/apply"
)

type applyCommand struct {
	config   *ClientConfig
	applysvc apply.Service
}

func (cmd *applyCommand) setup() error {
	cfg, err := LoadClientConfig()
	if err != nil {
		return err
	}
	cmd.config = cfg
	logger := log.NewLogfmtLogger(os.Stderr)
	applysvc, err := apply.NewClient(cfg.ServerURL, logger, cfg.APIToken)
	if err != nil {
		return err
	}
	cmd.applysvc = applysvc
	return nil
}

func (cmd *applyCommand) Run(args []string) error {
	if len(args) < 1 {
		cmd.Usage()
		os.Exit(1)
	}
	if err := cmd.setup(); err != nil {
		return err
	}
	var run func([]string) error
	switch strings.ToLower(args[0]) {
	case "blueprints":
		run = cmd.applyBlueprint
	default:
		cmd.Usage()
		os.Exit(1)
	}
	return run(args[1:])
}

func (cmd *applyCommand) Usage() error {
	const applyUsage = `
Apply a resource.

Valid resource types:

  * blueprints

Examples:
  # Get a list of devices
  mdmctl apply blueprints -f /path/to/blueprint.json

`
	fmt.Println(applyUsage)
	return nil
}

func (cmd *applyCommand) applyBlueprint(args []string) error {
	flagset := flag.NewFlagSet("blueprints", flag.ExitOnError)
	flagset.Usage = usageFor(flagset, "mdmctl apply blueprints [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}
	return nil
}
