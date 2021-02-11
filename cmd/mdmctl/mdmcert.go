package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	"golang.org/x/crypto/pkcs12"

	"github.com/micromdm/micromdm/pkg/crypto"
	"github.com/micromdm/micromdm/pkg/crypto/mdmcertutil"
)

type mdmcertCommand struct {
	*remoteServices
}

func (cmd *mdmcertCommand) setup() error {
	logger := log.NewLogfmtLogger(os.Stderr)
	remote, err := setupClient(logger)
	if err != nil {
		return err
	}
	cmd.remoteServices = remote
	return nil
}

func (cmd *mdmcertCommand) Usage() error {
	const usageText = `
Create new MDM Push Certificate.
This utility helps obtain an MDM Push Certificate using the Apple Developer MDM CSR option in the enterprise developer portal.

First you must create a vendor CSR which you will upload to the enterprise developer portal and get a signed MDM Vendor certificate. Use the MDM-CSR option in the dev portal when creating the certificate.
The MDM Vendor certificate is required in order to obtain the MDM push certificate. After you complete the MDM-CSR step, copy the downloaded file to the same folder as the private key. By default this will be
mdm-certificates/

    mdmctl mdmcert vendor -password=secret -country=US -email=admin@acme.co

Next, create a push CSR. This step generates a CSR required to get a signed a push certificate.

	mdmctl mdmcert push -password=secret -country=US -email=admin@acme.co

Once you created the push CSR, you mush sign the push CSR with the MDM Vendor Certificate, and get a push certificate request file.

    mdmctl mdmcert vendor -sign -cert=./mdm-certificates/mdm.cer -password=secret

Once generated, upload the PushCertificateRequest.plist file to https://identity.apple.com to obtain your MDM Push Certificate.
Use the push private key and the push cert you got from identity.apple.com in your MDM server.

Commands:
    vendor
    push
    upload
`
	fmt.Println(usageText)
	return nil

}

func (cmd *mdmcertCommand) Run(args []string) error {
	if len(args) < 1 {
		cmd.Usage()
		os.Exit(1)
	}

	if err := cmd.setup(); err != nil {
		return err
	}

	var run func([]string) error
	switch strings.ToLower(args[0]) {
	case "vendor":
		run = cmd.runVendor
	case "push":
		run = cmd.runPush
	case "upload":
		run = cmd.runUpload
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
	pushRequestFilename               = "PushCertificateRequest.plist"
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
		flVendorcertPath = flagset.String("cert", filepath.Join(mdmcertdir, "mdm.cer"), "Path to the MDM Vendor certificate from dev portal.")
		flPushCSRPath    = flagset.String("push-csr", filepath.Join(mdmcertdir, pushCSRFilename), "Path to the user CSR(required for the -sign step).")
		flKeyPath        = flagset.String("private-key", filepath.Join(mdmcertdir, vendorPKeyFilename), "Path to the vendor private key. A new RSA key will be created at this path.")
		flCSRPath        = flagset.String("out", filepath.Join(mdmcertdir, vendorCSRFilename), "Path to save the MDM Vendor CSR.")
	)

	if err := flagset.Parse(args); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(*flCSRPath), 0755); err != nil {
		return errors.Wrapf(err, "create directory %s", filepath.Dir(*flCSRPath))
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
		return errors.Wrapf(err, "write %s to file", pushRequestFilename)
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
		flLocalOnly= flagset.Bool("local-only",false,"No server configuration required.")

		flCSRPath = flagset.String("out", filepath.Join(mdmcertdir, pushCSRFilename), "Path to save the MDM Push Certificate request.")
	)

	if err := flagset.Parse(args); err != nil {
		return err
	}

    if !*flLocalOnly {
        if err := cmd.setup(); err != nil {
    		return err
   		}
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

func (cmd *mdmcertCommand) runUpload(args []string) error {
	flagset := flag.NewFlagSet("upload", flag.ExitOnError)
	flagset.Usage = usageFor(flagset, "mdmctl mdmcert upload [flags]")
	var (
		flKeyPass  = flagset.String("password", "", "Password to encrypt/read the RSA key.")
		flKeyPath  = flagset.String("private-key", filepath.Join(mdmcertdir, pushCertificatePrivateKeyFilename), "Path to the push certificate private key.")
		flCertPath = flagset.String("cert", "", "Path to the MDM Push Certificate.")
	)
	if err := flagset.Parse(args); err != nil {
		return err
	}

	cert, key, err := loadPushCerts(*flCertPath, *flKeyPath, *flKeyPass)
	if err != nil {
		return errors.Wrap(err, "load push certificate")
	}

	if err := cmd.configsvc.SavePushCertificate(context.Background(), cert, key); err != nil {
		return errors.Wrap(err, "upload push certificate and key to server")
	}

	return nil
}

func loadPushCerts(certPath, keyPath, keyPass string) (cert, key []byte, err error) {
	isP12 := filepath.Ext(certPath) == ".p12"
	if isP12 {
		pkcs12Data, err := ioutil.ReadFile(certPath)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "read p12 path %s", certPath)
		}
		pkeyi, certificate, err := pkcs12.Decode(pkcs12Data, keyPass)
		if err != nil {
			return nil, nil, errors.Wrap(err, "decode pkcs12 file")
		}
		pkey, ok := pkeyi.(*rsa.PrivateKey)
		if !ok {
			return nil, nil, errors.New("private key not a valid rsa key")
		}

		pemKey := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(pkey),
		})

		pemCert := pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: certificate.Raw,
		})
		return pemCert, pemKey, nil
	}

	keyData, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "read push certificate private key at path %s", keyPath)
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
		return nil, nil, errors.Wrapf(err, "parse push certiificate private key %s", keyPath)
	}

	pemKey := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	})

	certificate, err := crypto.ReadPEMCertificateFile(certPath)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "read push certificate from pem file %s", certPath)
	}

	pemCert := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certificate.Raw,
	})

	return pemCert, pemKey, nil
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
