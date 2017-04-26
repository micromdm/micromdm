package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/micromdm/micromdm/blueprint"
	"github.com/micromdm/micromdm/core/apply"
	uuid "github.com/satori/go.uuid"
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
	applysvc, err := apply.NewClient(cfg.ServerURL, logger, cfg.APIToken, httptransport.SetClient(skipVerifyHTTPClient(cmd.config.SkipVerify)))
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
	case "dep-tokens":
		run = cmd.applyDEPTokens
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
  * dep-tokens

Examples:
  # Get a list of devices
  mdmctl apply blueprints -f /path/to/blueprint.json

`
	fmt.Println(applyUsage)
	return nil
}

func (cmd *applyCommand) applyBlueprint(args []string) error {
	flagset := flag.NewFlagSet("blueprints", flag.ExitOnError)
	var (
		flBlueprintPath    = flagset.String("f", "", "filename of blueprint JSON to apply")
		flNewBlueprintPath = flagset.String("generate-blueprint", "", "filename of new template blueprint JSON to create")
	)
	flagset.Usage = usageFor(flagset, "mdmctl apply blueprints [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}
	if *flBlueprintPath == "" && *flNewBlueprintPath == "" {
		return errors.New("must provide -f or -generate-blueprint parameter")
	}
	if *flBlueprintPath != "" {
		if _, err := os.Stat(*flBlueprintPath); os.IsNotExist(err) {
			return err
		}
		jsonBytes, err := ioutil.ReadFile(*flBlueprintPath)
		if err != nil {
			return err
		}
		var blpt blueprint.Blueprint
		err = json.Unmarshal(jsonBytes, &blpt)
		if err != nil {
			return err
		}
		ctx := context.Background()
		err = cmd.applysvc.ApplyBlueprint(ctx, &blpt)
		if err != nil {
			return err
		}
		fmt.Println("applied blueprint", *flBlueprintPath)
		return nil
	}
	if *flNewBlueprintPath != "" {
		newBlueprintFile, err := os.Create(*flNewBlueprintPath)
		if err != nil {
			return err
		}
		defer newBlueprintFile.Close()

		newBlueprint := new(blueprint.Blueprint)
		newBlueprint.Name = "exampleName"
		newBlueprint.UUID = uuid.NewV4().String()
		newBlueprint.ApplicationURLs = []string{cmd.config.ServerURL + "repo/exampleAppManifest.plist"}
		newBlueprint.Profiles = []blueprint.Mobileconfig{blueprint.Mobileconfig([]byte("this should be a configuration profile"))}

		enc := json.NewEncoder(newBlueprintFile)
		enc.SetIndent("", "  ")
		err = enc.Encode(newBlueprint)
		if err != nil {
			return err
		}

		fmt.Println("wrote", *flNewBlueprintPath)
	}
	return nil
}

func (cmd *applyCommand) applyDEPTokens(args []string) error {
	flagset := flag.NewFlagSet("dep-tokens", flag.ExitOnError)
	var (
		flPublicKeyPath = flagset.String("import-token", "", "filename of p7m encrypted token file (downloaded from DEP portal)")
	)
	flagset.Usage = usageFor(flagset, "mdmctl apply dep-tokens [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}
	if *flPublicKeyPath == "" {
		return errors.New("must provide -import-token parameter")
	}
	if _, err := os.Stat(*flPublicKeyPath); os.IsNotExist(err) {
		return err
	}
	p7mBytes, err := ioutil.ReadFile(*flPublicKeyPath)
	if err != nil {
		return err
	}
	ctx := context.Background()
	err = cmd.applysvc.ApplyDEPToken(ctx, p7mBytes)
	if err != nil {
		return err
	}
	fmt.Println("imported DEP token")
	return nil
}
