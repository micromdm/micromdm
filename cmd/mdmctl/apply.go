package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/micromdm/micromdm/platform/apiserver/apply"
	"github.com/micromdm/micromdm/platform/blueprint"
	"github.com/micromdm/micromdm/platform/profile"
)

type applyCommand struct {
	config   *ServerConfig
	applysvc apply.Service
}

func (cmd *applyCommand) setup() error {
	cfg, err := LoadServerConfig()
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
	case "dep-profiles":
		run = cmd.applyDEPProfile
	case "profiles":
		run = cmd.applyProfile
	case "app":
		run = cmd.applyApp
	case "users":
		run = cmd.applyUser
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
  * profiles
  * users
  * dep-tokens
  * dep-profiles
  * app

Examples:
  # Apply a Blueprint.
  mdmctl apply blueprints -f /path/to/blueprint.json

  # Apply a DEP Profile.
  mdmctl apply dep-profiles -f /path/to/dep-profile.json

`
	fmt.Println(applyUsage)
	return nil
}

func (cmd *applyCommand) applyBlueprint(args []string) error {
	flagset := flag.NewFlagSet("blueprints", flag.ExitOnError)
	var (
		flBlueprintPath = flagset.String("f", "", "filename of blueprint JSON to apply")
		flTemplate      = flagset.Bool("template", false, "print a new blueprint template")
	)
	flagset.Usage = usageFor(flagset, "mdmctl apply blueprints [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}

	if *flTemplate {
		newBlueprint := &blueprint.Blueprint{
			Name:               "exampleName",
			UUID:               uuid.NewV4().String(),
			ApplicationURLs:    []string{cmd.config.ServerURL + "repo/exampleAppManifest.plist"},
			ProfileIdentifiers: []string{"com.example.my.profile"},
			UserUUID:           []string{"your-admin-account-uuid"},
			ApplyAt:            []string{"Enroll"},
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(newBlueprint); err != nil {
			return errors.Wrap(err, "encode blueprint template")
		}
		return nil
	}

	if *flBlueprintPath == "" {
		flagset.Usage()
		return errors.New("bad input: must provide -f or -template flag")
	}

	if *flBlueprintPath != "" {
		jsonBytes, err := readBytesFromPath(*flBlueprintPath)
		if err != nil {
			return err
		}
		var blpt blueprint.Blueprint
		err = json.Unmarshal(jsonBytes, &blpt)
		if err != nil {
			return err
		}

		// validate the blueprint account rules
		if len(blpt.UserUUID) == 0 && (blpt.SkipPrimarySetupAccountCreation || blpt.SetPrimarySetupAccountAsRegularUser) {
			return errors.New("SkipPrimarySetupAccountCreation and SetPrimarySetupAccountAsRegularUser can only be true if there is an account in the UserUUID array.")
		}

		ctx := context.Background()
		err = cmd.applysvc.ApplyBlueprint(ctx, &blpt)
		if err != nil {
			return err
		}
		fmt.Println("applied blueprint", *flBlueprintPath)
		return nil
	}

	return nil
}

func (cmd *applyCommand) applyDEPTokens(args []string) error {
	flagset := flag.NewFlagSet("dep-tokens", flag.ExitOnError)
	var (
		flTokenPath = flagset.String(
			"import",
			filepath.Join(defaultmdmctlFilesPath, "DEPOAuthToken.json"),
			"Filename of p7m encrypted token file (downloaded from DEP portal)")
	)

	flagset.Usage = usageFor(flagset, "mdmctl apply dep-tokens [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}

	if *flTokenPath == "" {
		return errors.New("must provide -import-token parameter")
	}
	if _, err := os.Stat(*flTokenPath); os.IsNotExist(err) {
		return err
	}
	p7mBytes, err := ioutil.ReadFile(*flTokenPath)
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

func (cmd *applyCommand) applyProfile(args []string) error {
	flagset := flag.NewFlagSet("profiles", flag.ExitOnError)
	var (
		flProfilePath = flagset.String("f", "", "filename of profile to apply")
	)
	flagset.Usage = usageFor(flagset, "mdmctl apply profiles [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}
	if *flProfilePath == "" {
		flagset.Usage()
		return errors.New("bad input: must provide -f parameter. use - for stdin")
	}
	profileBytes, err := readBytesFromPath(*flProfilePath)
	if err != nil {
		return err
	}

	// TODO: to consider just uploading the Mobileconfig data (without a
	// Profile struct and doing init server side)
	var p profile.Profile
	p.Mobileconfig = profileBytes
	p.Identifier, err = p.Mobileconfig.GetPayloadIdentifier()
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = cmd.applysvc.ApplyProfile(ctx, &p)
	if err != nil {
		return err
	}

	fmt.Printf("applied profile id %s from %s\n", p.Identifier, *flProfilePath)
	return nil
}

func readBytesFromPath(path string) ([]byte, error) {
	if path == "-" {
		return ioutil.ReadAll(os.Stdin)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}
	return ioutil.ReadFile(path)
}
