package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"crypto/x509"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/micromdm/micromdm/core/list"
	"github.com/micromdm/micromdm/crypto"
)

type getCommand struct {
	config *ClientConfig
	list   list.Service
}

func (cmd *getCommand) setup() error {
	cfg, err := LoadClientConfig()
	if err != nil {
		return err
	}
	cmd.config = cfg
	logger := log.NewLogfmtLogger(os.Stderr)
	listsvc, err := list.NewClient(cfg.ServerURL, logger, cfg.APIToken, httptransport.SetClient(skipVerifyHTTPClient(cmd.config.SkipVerify)))
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
	case "dep-tokens":
		run = cmd.getDepTokens
	case "blueprints":
		run = cmd.getBlueprints
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
	flagset.Usage = usageFor(flagset, "mdmctl get devices [flags]")
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

func (cmd *getCommand) getDepTokens(args []string) error {
	flagset := flag.NewFlagSet("dep-tokens", flag.ExitOnError)
	var (
		flFullCK        = flagset.Bool("v", false, "display full ConsumerKey in summary list")
		flPublicKeyPath = flagset.String("export-public-key", "", "filename of public key to write (to be uploaded to deploy.apple.com)")
		flTokenPath     = flagset.String("export-token", "", "filename to save decrypted oauth token (JSON)")
	)
	flagset.Usage = usageFor(flagset, "mdmctl get dep-tokens [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintf(w, "ConsumerKey\tAccessTokenExpiry\n")
	ctx := context.Background()
	tokens, certBytes, err := cmd.list.GetDEPTokens(ctx)
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
		flJSONName      = flagset.String("json", "", "file name of JSON to save for a single result")
	)
	flagset.Usage = usageFor(flagset, "mdmctl get blueprints [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}

	ctx := context.Background()
	blueprints, err := cmd.list.GetBlueprints(ctx, list.GetBlueprintsOption{FilterName: *flBlueprintName})
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintf(w, "Name\tUUID\tManifests\tProfiles\n")
	for _, bp := range blueprints {
		fmt.Fprintf(
			w,
			"%s\t%s\t%d\t%d\n",
			bp.Name,
			bp.UUID,
			len(bp.ApplicationURLs),
			len(bp.Profiles),
		)
	}
	w.Flush()

	if *flJSONName != "" && len(blueprints) > 0 {
		bp := blueprints[0]

		bpFile, err := os.Create(*flJSONName)
		if err != nil {
			return err
		}
		defer bpFile.Close()

		enc := json.NewEncoder(bpFile)
		enc.SetIndent("", "  ")
		err = enc.Encode(bp)
		if err != nil {
			return err
		}

		fmt.Printf("\nWrote Blueprint to: %s\n", *flJSONName)
		if len(blueprints) > 1 {
			fmt.Println("WARNING: more than one Blueprint returned; only saved first")
		}
	}

	return nil
}
