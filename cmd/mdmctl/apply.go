package main

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/pkcs12"

	"github.com/micromdm/micromdm/pkg/crypto/profileutil"
	"github.com/micromdm/micromdm/platform/blueprint"
	"github.com/micromdm/micromdm/platform/profile"
)

type applyCommand struct {
	config *ServerConfig
	*remoteServices
}

func (cmd *applyCommand) setup() error {
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
	case "block":
		run = cmd.applyBlock
	case "users":
		run = cmd.applyUser
	case "dep-autoassigner":
		run = cmd.applyDEPAutoAssigner
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
  * dep-autoassigner
  * app
  * block

Examples:
  # Apply a Blueprint.
  mdmctl apply blueprints -f /path/to/blueprint.json

  # Apply a DEP Profile.
  mdmctl apply dep-profiles -f /path/to/dep-profile.json

`
	fmt.Print(applyUsage)
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
			UUID:               uuid.New().String(),
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
		err = cmd.blueprintsvc.ApplyBlueprint(ctx, &blpt)
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
	err = cmd.configsvc.ApplyDEPToken(ctx, p7mBytes)
	if err != nil {
		return err
	}
	fmt.Println("imported DEP token")
	return nil
}

func (cmd *applyCommand) applyBlock(args []string) error {
	flagset := flag.NewFlagSet("block", flag.ExitOnError)
	var (
		flUDID = flagset.String("udid", "", "UDID of a device to block.")
	)
	flagset.Usage = usageFor(flagset, "mdmctl apply block [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}
	if *flUDID == "" {
		flagset.Usage()
		return errors.New("bad input: must provide a device UDID to block.")
	}
	if err := cmd.blocksvc.BlockDevice(context.Background(), *flUDID); err != nil {
		return err
	}

	// trigger a push
	u, err := url.Parse(cmd.config.ServerURL)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	u.Path = "/push/" + url.QueryEscape(*flUDID)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	req.SetBasicAuth("micromdm", cmd.config.APIToken)
	skipVerifyHTTPClient(cmd.config.SkipVerify).Do(req)
	return nil
}

func loadSigningKey(keyPass, keyPath, certPath string) (crypto.PrivateKey, *x509.Certificate, error) {
	certData, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, nil, err
	}

	isP12 := filepath.Ext(certPath) == ".p12"
	if isP12 {
		pkey, cert, err := pkcs12.Decode(certData, keyPass)
		return pkey, cert, errors.Wrap(err, "decode p12 contents")
	}

	keyData, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, nil, errors.Wrap(err, "read key from file")
	}

	keyDataBlock, _ := pem.Decode(keyData)
	if keyDataBlock == nil {
		return nil, nil, errors.Errorf("invalid PEM data for private key %s", keyPath)
	}
	var pemKeyData []byte
	if x509.IsEncryptedPEMBlock(keyDataBlock) {
		b, err := x509.DecryptPEMBlock(keyDataBlock, []byte(keyPass))
		if err != nil {
			return nil, nil, fmt.Errorf("decrypting DES private key %s", err)
		}
		pemKeyData = b
	} else {
		pemKeyData = keyDataBlock.Bytes
	}

	priv, err := x509.ParsePKCS1PrivateKey(pemKeyData)
	if err != nil {
		return nil, nil, errors.Wrap(err, "parse private key")
	}

	pub, _ := pem.Decode(certData)
	if pub == nil {
		return nil, nil, errors.Errorf("invalid PEM data for certificate %q", certPath)
	}

	cert, err := x509.ParseCertificate(pub.Bytes)
	if err != nil {
		return nil, nil, errors.Wrap(err, "parse PEM certificate data")
	}

	return priv, cert, nil
}

func (cmd *applyCommand) applyProfile(args []string) error {
	flagset := flag.NewFlagSet("profiles", flag.ExitOnError)
	var (
		flProfilePath = flagset.String("f", "", "Path to profile payload.")
		flSign        = flagset.Bool("sign", false, "Sign the profile. Requires key and certificate path.")
		flOut         = flagset.String("out", "", "Output path for signed profile(optional).")
		flKeyPass     = flagset.String("password", "", "Password to encrypt/read the signing key(optional) or p12 file.")
		flKeyPath     = flagset.String("private-key", "", "Path to the signing private key. Don't use with p12 file.")
		flCertPath    = flagset.String("cert", "", "Path to the signing certificate or p12 file.")
	)
	flagset.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n",
			`Upload profiles to the server.

Uploaded profiles can also be specified in a blueprint, which will be applied on device enrollment.
This command can also be used to replace the enrollment profile.
Profiles can be signed before upload.

Examples

  # Upload a mobileconfig
  mdmctl apply profiles -f /path/to/profile.mobileconfig

  # Sign and upload
  mdmctl apply profiles -f /path/to/profile.mobileconfig -private-key key.pem -cert certificate.pem -password secret -sign

  # Sign and save to local directory instead of uploading
  # Use "-out -" print the output to stdout instead of a file.
  mdmctl apply profiles -f /path/to/profile.mobileconfig -private-key key.pem -cert certificate.pem -password secret -sign -out signed.mobileconfig

`)
		usageFor(flagset, "mdmctl apply profiles [flags]")()
	}
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

	if *flSign {
		priv, pub, err := loadSigningKey(*flKeyPass, *flKeyPath, *flCertPath)
		if err != nil {
			return errors.Wrap(err, "loading signing certificate and private key")
		}
		signed, err := profileutil.Sign(priv, pub, profileBytes)
		if err != nil {
			return errors.Wrap(err, "signing profile with the specified key")
		}

		if *flOut == "-" { // print to stdout and return
			_, err = os.Stdout.Write(signed)
			return err
		} else if *flOut != "" { // write to file and return
			return ioutil.WriteFile(*flOut, signed, 0644)
		}

		profileBytes = signed
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
	err = cmd.profilesvc.ApplyProfile(ctx, &p)
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
