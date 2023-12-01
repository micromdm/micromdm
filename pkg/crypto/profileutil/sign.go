// Package profileutil signs configuration profiles.
package profileutil

import (
	"crypto"
	"crypto/x509"

	"github.com/pkg/errors"
	"github.com/smallstep/pkcs7"
)

// Sign takes an unsigned payload and signs it with the provided private key and certificate.
func Sign(key crypto.PrivateKey, cert *x509.Certificate, mobileconfig []byte) ([]byte, error) {
	sd, err := pkcs7.NewSignedData(mobileconfig)
	if err != nil {
		return nil, errors.Wrap(err, "create signed data for mobileconfig")
	}

	if err := sd.AddSigner(cert, key, pkcs7.SignerInfoConfig{}); err != nil {
		return nil, errors.Wrap(err, "add crypto signer to mobileconfig signed data")
	}

	signedMobileconfig, err := sd.Finish()
	return signedMobileconfig, errors.Wrap(err, "complete mobileconfig signing")
}
