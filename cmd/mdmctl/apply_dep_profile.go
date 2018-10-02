package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/dep"
	"github.com/micromdm/micromdm/pkg/crypto"
	"github.com/micromdm/micromdm/platform/dep/sync"
)

func certificatesFromURL(serverURL string, insecure bool) ([]*x509.Certificate, error) {
	urlParsed, err := url.Parse(serverURL)
	if err != nil {
		return nil, err
	}
	addr := urlParsed.Host
	if urlParsed.Port() == "" {
		addr += ":443"
	}
	conn, err := tls.Dial("tcp", addr, &tls.Config{InsecureSkipVerify: insecure})
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return conn.ConnectionState().PeerCertificates, nil
}

func (cmd *applyCommand) applyDEPProfile(args []string) error {
	flagset := flag.NewFlagSet("dep-profiles", flag.ExitOnError)
	var (
		flProfilePath = flagset.String("f", "", "filename of DEP profile to apply")
		flTemplate    = flagset.Bool("template", false, "print a JSON example of a DEP profile")
		flAnchorFile  = flagset.String("anchor", "", "filename of PEM cert(s) to add to anchor certs in template")
		flUseServer   = flagset.Bool("use-server-cert", false, "use the server cert(s) to add to anchor certs in template")
		flFilter      = flagset.String("filter", "", "set the auto-assign filter to for the defined profile")
	)
	flagset.Usage = usageFor(flagset, "mdmctl apply dep-profiles [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}

	if *flTemplate {
		var anchorCerts []*x509.Certificate
		if *flAnchorFile != "" {
			certs, err := crypto.ReadPEMCertificatesFile(*flAnchorFile)
			if err != nil {
				return err
			}
			anchorCerts = append(anchorCerts, certs...)
		}
		if *flUseServer {
			certs, err := certificatesFromURL(cmd.config.ServerURL, cmd.config.SkipVerify)
			if err != nil {
				return err
			}
			anchorCerts = append(anchorCerts, certs...)
		}
		return printDEPProfileTemplate(anchorCerts)
	}

	if *flProfilePath == "" {
		flagset.Usage()
		return errors.New("bad input: must provide -f or -template parameter ")
	}

	var output *os.File
	{
		if *flProfilePath == "-" {
			output = os.Stdin
		} else {
			var err error
			output, err = os.Open(*flProfilePath)
			if err != nil {
				return err
			}
			defer output.Close()
		}
	}

	var profile dep.Profile
	if err := json.NewDecoder(output).Decode(&profile); err != nil {
		return errors.Wrap(err, "decode DEP Profile JSON")
	}

	resp, err := cmd.depsvc.DefineProfile(context.TODO(), &profile)
	if err != nil {
		return errors.Wrap(err, "define dep profile")
	}

	// TODO: it would be nice to encode back a profile that save the
	// UUID for future reference.
	fmt.Printf("Defined DEP Profile with UUID %s\n", resp.ProfileUUID)

	if *flFilter != "" {
		assigner := sync.AutoAssigner{Filter: *flFilter, ProfileUUID: resp.ProfileUUID}
		err := cmd.depsyncsvc.ApplyAutoAssigner(context.TODO(), &assigner)
		if err != nil {
			return errors.Wrap(err, "set auto-assigner")
		}
		fmt.Printf("Saved auto-assign filter '%s' for this DEP profile\n", *flFilter)
	}

	return nil
}

func printDEPProfileTemplate(anchorCerts []*x509.Certificate) error {
	var anchorCertStr string = "[]"

	// convert certificates into base64 encoded strings
	// json.Marshal does this for us for byte[] arrays
	if len(anchorCerts) > 0 {
		var certs [][]byte
		for _, cert := range anchorCerts {
			certs = append(certs, cert.Raw)
		}
		jsonBytes, err := json.Marshal(certs)
		if err != nil {
			return nil
		}
		anchorCertStr = string(jsonBytes)
	}

	resp := fmt.Sprintf(`{
  "profile_name": "(Required) Human readable name",
  "url": "https://mymdm.example.org/mdm/enroll",
  "allow_pairing": true,
  "is_supervised": false,
  "is_multi_user": false,
  "is_mandatory": false,
  "await_device_configured": false,
  "is_mdm_removable": true,
  "support_phone_number": "(Optional) +1 408 555 1010",
  "support_email_address": "(Optional) support@example.com",
  "org_magic": "(Optional)",
  "anchor_certs": %s,
  "supervising_host_certs": [],
  "skip_setup_items": ["AppleID", "Android"],
  "department": "(Optional) support@example.com",
  "devices": ["SERIAL1","SERIAL2"]
}`, anchorCertStr)
	fmt.Println(resp)
	return nil
}
