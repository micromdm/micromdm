package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/crypto/mdmcertutil"
)

type mdmcertCommand struct{}

func (cmd *mdmcertCommand) Usage() error {
	const usageText = `
Create new MDM Push Certificate.
This utility helps obtain a MDM Push Certificate using the Apple Developer MDM CSR option in the enterprise developer portal.

First you must create a vendor CSR which you will upload to the enterprise developer portal and get a signed MDM Vendor certificate. Use the MDM-CSR option in the dev portal when creating the certificate.
The MDM Vendor certificate is required in order to obtain the MDM push certificate.

Next, create a push CSR. This step generates a CSR required to get a signed a push certificate. 
Once you created the push CSR, you mush sign the push CSR with the MDM Vendor Certificate, and get a push certificate request file.

Once generated, upload the push certificate request to https://identity.apple.com. 
Use the push private key and the push cert you got from identity.apple.com in your MDM server.

Commands:
    vendor
    push

Examples:
    mdmctl mdmcert vendor -password=secret -country=US -email=foo@gmail.com
`
	fmt.Println(usageText)
	return nil

}

func (cmd *mdmcertCommand) Run(args []string) error {
	if len(args) < 1 {
		cmd.Usage()
		os.Exit(1)
	}

	var run func([]string) error
	switch strings.ToLower(args[0]) {
	case "vendor":
		run = cmd.runVendor
	case "push":
		run = cmd.runPush
	default:
		cmd.Usage()
		os.Exit(1)
	}

	return run(args[1:])
}

const (
	pushCSRFilename                   = "PushCertificateRequest.csr"
	pushCertificatePrivateKeyFilename = "PushCertificatePrivateKey.key"
	vendorPKeyFilename                = "VendorPrivateKey.key"
	vendorCSRFilename                 = "VendorCertificateRequest.csr"
	pushRequestFilename               = "PushCertificateRequest"
	mdmcertdir                        = "mdm-certificates"
)

func (cmd *mdmcertCommand) runVendor(args []string) error {
	flagset := flag.NewFlagSet("vendor", flag.ExitOnError)
	flagset.Usage = usageFor(flagset, "mdmctl mdmcert vendor [flags]")
	var (
		flSign           = flagset.Bool("sign", false, "Signs a user CSR with the MDM vendor certificate.")
		flEmail          = flagset.String("email", "", "Email address to use in CSR Subject.")
		flCountry        = flagset.String("country", "US", "Two letter country code for the CSR Subject(example: US).")
		flCN             = flagset.String("cn", "micromdm-vendor", "CommonName for the CSR Subject.")
		flPKeyPass       = flagset.String("password", "", "Password to encrypt/read the RSA key.")
		flVendorcertPath = flagset.String("cert", "mdm.cer", "Path to the MDM Vendor certificate from dev portal.")
		flPushCSRPath    = flagset.String("push-csr", filepath.Join(mdmcertdir, pushCSRFilename), "Path to the user CSR(required for the -sign step).")
		flKeyPath        = flagset.String("private-key", filepath.Join(mdmcertdir, vendorPKeyFilename), "Path to the vendor private key. A new RSA key will be created at this path.")
		flCSRPath        = flagset.String("out", filepath.Join(mdmcertdir, vendorCSRFilename), "Path to save the MDM Vendor CSR.")
	)

	if err := flagset.Parse(args); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(*flCSRPath), 0755); err != nil {
		errors.Wrapf(err, "create directory %s", filepath.Dir(*flCSRPath))
	}

	password := []byte(*flPKeyPass)
	if *flSign {
		request, err := mdmcertutil.CreatePushCertificateRequest(
			*flVendorcertPath,
			*flPushCSRPath,
			*flKeyPath,
			password,
		)
		if err != nil {
			return errors.Wrap(err, "signing push certificate request with vendor private key")
		}
		encoded, err := request.Encode()
		if err != nil {
			return errors.Wrap(err, "encode base64 push certificate request")
		}
		err = ioutil.WriteFile(filepath.Join(mdmcertdir, pushRequestFilename), encoded, 0600)
		return errors.Wrap(err, "write PushCertificateRequest to file file")
	}

	if err := checkCSRFlags(*flCN, *flCountry, *flEmail, password); err != nil {
		return errors.Wrap(err, "Private key password, CN, Email, and country code must be specified when creating a CSR.")
	}

	request := &mdmcertutil.CSRConfig{
		CommonName:         *flCN,
		Country:            *flCountry,
		Email:              *flEmail,
		PrivateKeyPassword: password,
		PrivateKeyPath:     *flKeyPath,
		CSRPath:            *flCSRPath,
	}

	err := mdmcertutil.CreateCSR(request)
	return errors.Wrap(err, "creating MDM vendor CSR")
}

func (cmd *mdmcertCommand) runPush(args []string) error {
	flagset := flag.NewFlagSet("push", flag.ExitOnError)
	flagset.Usage = usageFor(flagset, "mdmctl mdmcert push [flags]")
	var (
		flEmail    = flagset.String("email", "", "Email address to use in CSR Subject.")
		flCountry  = flagset.String("country", "US", "Two letter country code for the CSR Subject(Example: US).")
		flCN       = flagset.String("cn", "micromdm-user", "CommonName for the CSR Subject.")
		flPKeyPass = flagset.String("password", "", "Password to encrypt/read the RSA key.")
		flKeyPath  = flagset.String("private-key", filepath.Join(mdmcertdir, pushCertificatePrivateKeyFilename), "Path to the push certificate private key. A new RSA key will be created at this path.")

		flCSRPath = flagset.String("out", filepath.Join(mdmcertdir, pushCSRFilename), "Path to save the MDM Push Certificate request.")
	)

	if err := flagset.Parse(args); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(*flCSRPath), 0755); err != nil {
		errors.Wrapf(err, "create directory %s", filepath.Dir(*flCSRPath))
	}

	password := []byte(*flPKeyPass)
	if err := checkCSRFlags(*flCN, *flCountry, *flEmail, password); err != nil {
		return errors.Wrap(err, "Private key password, CN, Email, and country code must be specified when creating a CSR.")
	}

	request := &mdmcertutil.CSRConfig{
		CommonName:         *flCN,
		Country:            *flCountry,
		Email:              *flEmail,
		PrivateKeyPassword: password,
		PrivateKeyPath:     *flKeyPath,
		CSRPath:            *flCSRPath,
	}

	err := mdmcertutil.CreateCSR(request)
	return errors.Wrap(err, "creating MDM Push certificate request.")
}

func checkCSRFlags(cname, country, email string, password []byte) error {
	if cname == "" {
		return errors.New("cn flag not specified")
	}
	if email == "" {
		return errors.New("email flag not specified")
	}
	if country == "" {
		return errors.New("country flag not specified")
	}
	if len(password) == 0 {
		return errors.New("private key password empty")
	}
	if len(country) != 2 {
		return errors.New("must be a two letter country code")
	}
	return nil
}
