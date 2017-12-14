package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-kit/kit/log"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/micromdm/micromdm/core/remove"
)

type removeCommand struct {
	config *ServerConfig
	remove remove.Service
}

func (cmd *removeCommand) setup() error {
	cfg, err := LoadServerConfig()
	if err != nil {
		return err
	}
	cmd.config = cfg
	logger := log.NewLogfmtLogger(os.Stdout)
	rmsvc, err := remove.NewClient(cfg.ServerURL, logger, cfg.APIToken, httptransport.SetClient(skipVerifyHTTPClient(cmd.config.SkipVerify)))
	if err != nil {
		return err
	}
	cmd.remove = rmsvc
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
  * profiles`

	fmt.Println(getUsage)
	return nil
}
