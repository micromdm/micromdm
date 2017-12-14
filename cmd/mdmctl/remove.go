package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-kit/kit/log"
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
	logger := log.NewLogfmtLogger(os.Stdout)
	remote, err := setupClient(logger)
	if err != nil {
		return err
	}
	cmd.remoteServices = remote
	return nil
}

func (cmd *removeCommand) Run(args []string) error {
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
		run = cmd.removeBlueprints
	case "profiles":
		run = cmd.removeProfiles
	case "block":
		run = cmd.removeBlock
	default:
		cmd.Usage()
		os.Exit(1)
	}

	return run(args[1:])
}

func (cmd *removeCommand) Usage() error {
	const getUsage = `
Display one or many resources.

Valid resource types:

  * blueprints
  * profiles
  * block
`

	fmt.Println(getUsage)
	return nil
}
