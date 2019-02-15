package main

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/go-kit/kit/log"
	"github.com/micromdm/micromdm/pkg/crypto"
	"github.com/micromdm/micromdm/platform/blueprint"
	"github.com/micromdm/micromdm/platform/device"
	"github.com/micromdm/micromdm/platform/profile"
	"github.com/pkg/errors"
)

type getCommand struct {
	config *ServerConfig
	*remoteServices
}

func (cmd *getCommand) setup() error {
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
	case "dev", "device", "devices":
		run = cmd.getDevices
	case "dep-devices":
		run = cmd.getDEPDevices
	case "dep-account":
		run = cmd.getDEPAccount
	case "dep-profiles":
		run = cmd.getDEPProfiles
	case "dep-tokens":
		run = cmd.getDepTokens
	case "blueprints":
		run = cmd.getBlueprints
	case "profiles":
		run = cmd.getProfiles
	case "users":
		run = cmd.getUsers
	case "apps":
		run = cmd.getApps
	case "vpp-apps", "vpp-app":
		run = cmd.getVPPApps
	case "vpp-licenses", "vpp-license":
		run = cmd.getVPPLicenses
	case "dep-autoassigners":
		run = cmd.getDEPAutoAssigners
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
  * dep-tokens
  * dep-devices
  * dep-account
  * dep-profiles
  * dep-autoassigners
  * users
  * profiles
  * apps
  * vpp-apps
  * vpp-licenses

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
	var (
		flFilterSerials = flagset.String("serials", "", "comma separated list of serials to search")
	)
	flagset.Usage = usageFor(flagset, "mdmctl get devices [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	out := &devicesTableOutput{w}
	out.BasicHeader()
	defer out.BasicFooter()
	ctx := context.Background()

	// convert string to []string in case a filter serial is defined
	// nil serials will return all serials from the database
	var flFilterSerialsSlice []string
	if *flFilterSerials != "" {
		flFilterSerialsSlice = strings.Split(*flFilterSerials, ",")
	}

	devices, err := cmd.devicesvc.ListDevices(ctx, device.ListDevicesOption{FilterSerial: flFilterSerialsSlice})
	if err != nil {
		return err
	}
	for _, d := range devices {
		fmt.Fprintf(out.w, "%s\t%s\t%v\t%s\n", d.UDID, d.SerialNumber, d.EnrollmentStatus, d.LastSeen)
	}
	return nil
}

const defaultmdmctlFilesPath = "mdm-files"

func (cmd *getCommand) getDepTokens(args []string) error {
	flagset := flag.NewFlagSet("dep-tokens", flag.ExitOnError)
	var (
		flFullCK        = flagset.Bool("v", false, "Display full ConsumerKey in summary list")
		flPublicKeyPath = flagset.String(
			"export-public-key",
			filepath.Join(defaultmdmctlFilesPath, "DEPPublicKey"),
			"Filename of public key to write (to be uploaded to deploy.apple.com)",
		)
		flTokenPath = flagset.String(
			"export-token",
			filepath.Join(defaultmdmctlFilesPath, "DEPOAuthToken.json"),
			"Filename to save decrypted oauth token (JSON)")
	)
	flagset.Usage = usageFor(flagset, "mdmctl get dep-tokens [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintf(w, "ConsumerKey\tAccessTokenExpiry\n")
	ctx := context.Background()
	tokens, certBytes, err := cmd.configsvc.GetDEPTokens(ctx)
	if err != nil {
		return err
	}
	var ckTrimmed string
	for _, t := range tokens {
		if len(t.ConsumerKey) > 40 && !*flFullCK {
			ckTrimmed = t.ConsumerKey[0:39] + "…"
		} else {
			ckTrimmed = t.ConsumerKey
		}
		fmt.Fprintf(w, "%s\t%s\n", ckTrimmed, t.AccessTokenExpiry.String())
	}
	w.Flush()

	if *flPublicKeyPath != "" && certBytes != nil {

		if err := os.MkdirAll(filepath.Dir(*flPublicKeyPath), 0755); err != nil {
			return errors.Wrapf(err, "create directory %s", filepath.Dir(*flPublicKeyPath))
		}

		cert, err := x509.ParseCertificate(certBytes)
		if err != nil {
			return err
		}
		err = crypto.WritePEMCertificateFile(cert, *flPublicKeyPath)
		if err != nil {
			return err
		}
		fmt.Printf("\nWrote DEP public key to: %s\n", *flPublicKeyPath)
	}

	if *flTokenPath != "" && len(tokens) > 0 {
		t := tokens[0]

		if err := os.MkdirAll(filepath.Dir(*flTokenPath), 0755); err != nil {
			return errors.Wrapf(err, "create directory %s", filepath.Dir(*flTokenPath))
		}

		tokenFile, err := os.Create(*flTokenPath)
		if err != nil {
			return err
		}
		defer tokenFile.Close()

		err = json.NewEncoder(tokenFile).Encode(t)
		if err != nil {
			return err
		}

		fmt.Printf("\nWrote DEP token JSON to: %s\n", *flTokenPath)
		if len(tokens) > 1 {
			fmt.Println("WARNING: more than one DEP token returned; only saved first")
		}
	}

	return nil
}

func (cmd *getCommand) getBlueprints(args []string) error {
	flagset := flag.NewFlagSet("blueprints", flag.ExitOnError)
	var (
		flBlueprintName = flagset.String("name", "", "name of blueprint")
		flJSONName      = flagset.String("f", "-", "filename of JSON to save to")
	)
	flagset.Usage = usageFor(flagset, "mdmctl get blueprints [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}

	ctx := context.Background()
	blueprints, err := cmd.blueprintsvc.GetBlueprints(ctx, blueprint.GetBlueprintsOption{FilterName: *flBlueprintName})
	if err != nil {
		return err
	}

	if *flBlueprintName == "" || len(blueprints) < 1 {
		w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
		fmt.Fprintf(w, "Name\tUUID\tManifests\tProfiles\tApply At\n")
		for _, bp := range blueprints {
			var applyAtStr string
			if len(bp.ApplyAt) > 0 {
				applyAtStr = strings.Join(bp.ApplyAt, ",")
			} else {
				applyAtStr = "(None)"
			}
			fmt.Fprintf(
				w,
				"%s\t%s\t%d\t%d\t%s\n",
				bp.Name,
				bp.UUID,
				len(bp.ApplicationURLs),
				len(bp.ProfileIdentifiers),
				applyAtStr,
			)
		}
		w.Flush()
	} else if *flJSONName != "" {
		bp := blueprints[0]

		var output *os.File
		{
			if *flJSONName == "-" {
				output = os.Stdout
			} else {
				var err error
				output, err = os.Create(*flJSONName)
				if err != nil {
					return err
				}
				defer output.Close()
			}
		}

		enc := json.NewEncoder(output)
		enc.SetIndent("", "  ")
		err = enc.Encode(bp)
		if err != nil {
			return err
		}

		if *flJSONName != "-" {
			fmt.Printf("wrote blueprint %s to: %s\n", *flBlueprintName, *flJSONName)
			if len(blueprints) > 1 {
				fmt.Println("WARNING: more than one Blueprint returned; only saved first")
			}
		}
	}

	return nil
}

func (cmd *getCommand) getProfiles(args []string) error {
	flagset := flag.NewFlagSet("profiles", flag.ExitOnError)
	var (
		flProfilePath = flagset.String("f", "-", "filename of profile to write")
		flIdentifier  = flagset.String("id", "", "profile Identifier")
	)
	flagset.Usage = usageFor(flagset, "mdmctl get blueprints [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}

	ctx := context.Background()
	profiles, err := cmd.profilesvc.GetProfiles(ctx, profile.GetProfilesOption{Identifier: *flIdentifier})
	if err != nil {
		return err
	}

	if *flIdentifier == "" || len(profiles) < 1 {
		w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
		fmt.Fprintf(w, "Identifier\tLength\n")
		for _, p := range profiles {
			fmt.Fprintf(
				w,
				"%s\t%d\n",
				p.Identifier,
				len(p.Mobileconfig),
			)
		}
		w.Flush()
	} else if *flProfilePath != "" {
		p := profiles[0]

		var output *os.File
		{
			if *flProfilePath == "-" {
				output = os.Stdout
			} else {
				var err error
				output, err = os.Create(*flProfilePath)
				if err != nil {
					return err
				}
				defer output.Close()
			}
		}
		_, err = output.Write([]byte(p.Mobileconfig))
		if err != nil {
			return err
		}
		if *flProfilePath != "-" {
			fmt.Printf("wrote profile id %s to: %s\n", p.Identifier, *flProfilePath)
			if len(profiles) > 1 {
				fmt.Println("WARNING: more than one Profile returned; only saved first")
			}
		}
	}
	return nil
}
