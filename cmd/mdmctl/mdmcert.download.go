package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/fullsailor/pkcs7"
	"github.com/go-kit/kit/log"
	"github.com/micromdm/micromdm/pkg/crypto"
	"github.com/micromdm/micromdm/pkg/crypto/mdmcertutil"
	"github.com/pkg/errors"
)

const (
	mdmcertRequestURL = "https://mdmcert.download/api/v1/signrequest"
	// see
	// https://github.com/jessepeterson/commandment/blob/1352b51ba6697260d1111eccc3a5a0b5b9af60d0/commandment/mdmcert.py#L23-L28
	mdmcertAPIKey = "f847aea2ba06b41264d587b229e2712c89b1490a1208b7ff1aafab5bb40d47bc"
)

// format of a signing request to mdmcert.download
type signRequest struct {
	CSR     string `json:"csr"` // base64 encoded PEM CSR
	Email   string `json:"email"`
	Key     string `json:"key"`     // server key from above
	Encrypt string `json:"encrypt"` // mdmcert pki cert
}

type mdmcertDownloadCommand struct {
	*remoteServices
}

func (cmd *mdmcertDownloadCommand) setup() error {
	logger := log.NewLogfmtLogger(os.Stderr)
	remote, err := setupClient(logger)
	if err != nil {
		return err
	}
	cmd.remoteServices = remote
	return nil
}

func (cmd *mdmcertDownloadCommand) Usage() error {
	const usageText = `
Request new MDM Push Certificate from https://mdmcert.download
This utility helps obtain an MDM Push Certificate using the service
at mdmcert.download.

First we'll generate the initial request (which also generates a private key):

	mdmctl mdmcert.download -new -email=cool.mdm.admin@example.org

This will output the private key into the file mdmcert.download.key.
Then, after you check your email and download the request file you just
need to decrypt the push certificate request:

	mdmctl mdmcert.download -decrypt=~/Downloads/mdm_signed_request.20171122_094910_220.plist.b64.p7

This will output the push certificate request to mdmcert.download.req.
Upload this file to https://identity.apple.com and download the signed
certificate. Then use the 'mdmctl mdmcert upload' command to upload it,
(and the above private key) into MicroMDM.

`
	fmt.Println(usageText)
	return nil

}

func (cmd *mdmcertDownloadCommand) Run(args []string) error {
	flagset := flag.NewFlagSet("mdmcert.download", flag.ExitOnError)
	flagset.Usage = usageFor(flagset, "mdmctl mdmcert.download [flags]")
	var (
		flNew       = flagset.Bool("new", false, "Generates a new privkey and uploads new MDM request")
		flDecrypt   = flagset.String("decrypt", "", "Decrypts and mdmcert.download push certificate request")
		flEmail     = flagset.String("email", "", "Email address to use in mdmcert request & CSR Subject")
		flCountry   = flagset.String("country", "US", "Two letter country code for the CSR Subject (example: US).")
		flCN        = flagset.String("cn", "mdm-push", "CommonName for the CSR Subject.")
		flCertPath  = flagset.String("pki-cert", "mdmcert.download.pki.crt", "Path for generated MDMCert pki exchange certificate")
		flKeyPath   = flagset.String("pki-private-key", "mdmcert.download.pki.key", "Path for generated MDMCert pki exchange private key")
		flPKeyPass  = flagset.String("pki-password", "", "Password to encrypt/read the RSA key.")
		flCCSRPath  = flagset.String("push-csr", "mdmcert.download.push.csr", "Path for generated Push Certificate CSR")
		flCReqPath  = flagset.String("push-req", "mdmcert.download.push.req", "Path for generated Push Certificate Request")
		flCKeyPath  = flagset.String("push-private-key", "mdmcert.download.push.key", "Path to the generated Push Cert private key")
		flCPKeyPass = flagset.String("push-password", "", "Password to encrypt/read the push RSA key.")
	)

	if err := flagset.Parse(args); err != nil {
		cmd.Usage()
		return err
	}

	// neither flag was used
	if !*flNew && *flDecrypt == "" {
		cmd.Usage()
		return errors.New("bad input: must either use -new or -decrypt")
	}

	// both flags used
	if *flNew && (*flDecrypt != "") {
		// cmd.Usage()
		return errors.New("bad input: can't use both -new and -decrypt")
	}

	if *flNew {
		if *flEmail == "" {
			return errors.New("bad input: must provide -email")
		}

		paths := []string{*flCertPath, *flKeyPath, *flCCSRPath, *flCKeyPath}
		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				return fmt.Errorf("file already exists: %s", path)
			}
		}

		pkiKey, pkiCert, err := crypto.SimpleSelfSignedRSAKeypair("mdmcert.download", 365)
		if err != nil {
			return errors.Wrap(err, "could not create PKI keypair")
		}

		pemBlock := &pem.Block{
			Type:    "CERTIFICATE",
			Headers: nil,
			Bytes:   pkiCert.Raw,
		}
		pemPkiCert := pem.EncodeToMemory(pemBlock)

		if err := crypto.WritePEMCertificateFile(pkiCert, *flCertPath); err != nil {
			return errors.Wrap(err, "could not write PKI cert")
		}

		if *flPKeyPass != "" {
			err = crypto.WriteEncryptedPEMRSAKeyFile(pkiKey, []byte(*flPKeyPass), *flKeyPath)
		} else {
			err = crypto.WritePEMRSAKeyFile(pkiKey, *flKeyPath)
		}
		if err != nil {
			return errors.Wrap(err, "could not write private key")
		}

		pushKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return errors.Wrap(err, "could not generate push private key")
		}

		if *flCPKeyPass != "" {
			err = crypto.WriteEncryptedPEMRSAKeyFile(pushKey, []byte(*flCPKeyPass), *flCKeyPath)
		} else {
			err = crypto.WritePEMRSAKeyFile(pushKey, *flCKeyPath)
		}
		if err != nil {
			return errors.Wrap(err, "could not write push private key")
		}

		derBytes, err := mdmcertutil.NewCSR(pushKey, *flEmail, *flCountry, *flCN)
		if err != nil {
			return errors.Wrap(err, "could not generate push CSR")
		}
		pemCSR := mdmcertutil.PemCSR(derBytes)
		// Do we even need to write-out the CSR?
		err = ioutil.WriteFile(*flCCSRPath, pemCSR, 0600)
		if err != nil {
			return errors.Wrap(err, "could not write PEM file")
		}

		sign := newMdmcertDownloadSignRequest(*flEmail, pemCSR, pemPkiCert)
		req, err := sign.HTTPRequest()
		if err != nil {
			return errors.Wrap(err, "could not create http request")
		}
		err = sendMdmcertDownloadRequest(http.DefaultClient, req)
		if err != nil {
			return errors.Wrap(err, "error sending http request")
		}

		fmt.Print("Request successfully sent to mdmcert.download. Your CSR should now\n" +
			"be signed. Check your email for next steps. Then use the -decrypt option\n" +
			"to extract the CSR request which will then be uploaded to Apple.\n")

	} else { // -decrypt switch
		if _, err := os.Stat(*flCReqPath); err == nil {
			return fmt.Errorf("file already exists: %s", *flCReqPath)
		}
		hexBytes, err := ioutil.ReadFile(*flDecrypt)
		if err != nil {
			return errors.Wrap(err, "reading encrypted file")
		}
		pkcsBytes, err := hex.DecodeString(string(hexBytes))
		if err != nil {
			return errors.Wrap(err, "error decoding hex")
		}
		pkiCert, err := crypto.ReadPEMCertificateFile(*flCertPath)
		if err != nil {
			return errors.Wrap(err, "reading PKI certificate")
		}
		var pkiKey *rsa.PrivateKey
		if *flPKeyPass != "" {
			pkiKey, err = crypto.ReadEncryptedPEMRSAKeyFile(*flKeyPath, []byte(*flPKeyPass))
		} else {
			pkiKey, err = crypto.ReadPEMRSAKeyFile(*flKeyPath)
		}
		if err != nil {
			return errors.Wrap(err, "reading PKI private key")
		}
		ioutil.WriteFile("/tmp/fubar.p7", pkcsBytes, 0666)
		p7, err := pkcs7.Parse(pkcsBytes)
		if err != nil {
			return errors.Wrap(err, "parsing mdmcert PKCS7 response")
		}
		// fmt.Println(p7)
		content, err := p7.Decrypt(pkiCert, pkiKey)
		if err != nil {
			return errors.Wrap(err, "decrypting mdmcert PKCS7 response")
		}
		err = ioutil.WriteFile(*flCReqPath, content, 0666)
		if err != nil {
			return errors.Wrap(err, "writing Push Request response")
		}

		fmt.Printf("Successfully able to decrypt the MDM Push Certificate request! Please upload\n"+
			"the file '%s' to Apple by visiting https://identity.apple.com\n"+
			"Once your Push Certificate is signed by Apple you can download it\n"+
			"and import it into MicroMDM using the `mdmctl mdmcert upload` command\n", *flCReqPath)
	}

	return nil
}

func newMdmcertDownloadSignRequest(email string, pemCSR []byte, serverCertificate []byte) *signRequest {
	encodedCSR := base64.StdEncoding.EncodeToString(pemCSR)
	encodedServerCert := base64.StdEncoding.EncodeToString(serverCertificate)
	return &signRequest{
		CSR:     encodedCSR,
		Email:   email,
		Key:     mdmcertAPIKey,
		Encrypt: encodedServerCert,
	}
}

func (sign *signRequest) HTTPRequest() (*http.Request, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(sign); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", mdmcertRequestURL, ioutil.NopCloser(buf))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "micromdm/certhelper")
	return req, nil
}

func sendMdmcertDownloadRequest(client *http.Client, req *http.Request) error {
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received bad status from mdmcert.download. status=%q", resp.Status)
	}
	var jsn = struct {
		Result string
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&jsn); err != nil {
		return err
	}
	if jsn.Result != "success" {
		return fmt.Errorf("got unexpected result body: %q\n", jsn.Result)
	}
	return nil
}
