package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
)

type removeCommand struct {
	config *ServerConfig
	*remoteServices
}

func (cmd *removeCommand) setup() error {
	cfg, err := LoadServerConfig()
	if err != nil {
		return err
	}
	cmd.config = cfg
	logger := log.NewLogfmtLogger(os.Stderr)
	remote, err := setupClient(logger)
	if err != nil {
		return err
	}
	cmd.remoteServices = remote
	return nil
}

func (cmd *removeCommand) Run(args []string) error {
	if len(args) < 1 {
		err := cmd.Usage()
		if err != nil {
			fmt.Println(errors.Wrap(err, "calling usage"))
		}
		os.Exit(1)
	}

	if err := cmd.setup(); err != nil {
		return err
	}

	var run func([]string) error
	switch strings.ToLower(args[0]) {
	case "blueprints":
		run = cmd.removeBlueprints
	case "devices":
		run = cmd.removeDevices
	case "profiles":
		run = cmd.removeProfiles
	case "block":
		run = cmd.removeBlock
	case "dep-autoassigner":
		run = cmd.removeDEPAutoAssigner
	default:
		err := cmd.Usage()
		if err != nil {
			fmt.Println(errors.Wrap(err, "calling usage"))
		}
		os.Exit(1)
	}

	return run(args[1:])
}

func (cmd *removeCommand) Usage() error {
	const getUsage = `
Display one or many resources.

Valid resource types:

  * blueprints
  * devices
  * profiles
  * block
  * dep-autoassigner
`

	fmt.Println(getUsage)
	return nil
}
