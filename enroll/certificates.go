package enroll

import (
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"golang.org/x/crypto/pkcs12"
	"io/ioutil"
)

const certificatePEMBlockType string = "CERTIFICATE"

var oidASN1Topic = asn1.ObjectIdentifier{0, 9, 2342, 19200300, 100, 1, 1}

func GetPushTopicFromCert(certPath, certPass, keyPath string) (string, error) {
	certData, err := ioutil.ReadFile(certPath)
	if err != nil {
		return "", err
	}

	var cert *x509.Certificate
	if keyPath == "" {
		// if keyPath is empty, treat as PKCS12
		// note that buford does validity checks where
		// our direct certificate parsing does not
		_, cert, err = pkcs12.Decode(certData, certPass)
		if err != nil {
			return "", err
		}
	} else {
		pemBlock, _ := pem.Decode(certData)
		cert, err = x509.ParseCertificate(pemBlock.Bytes)
		if err != nil {
			return "", err
		}
	}

	for _, v := range cert.Subject.Names {

		if v.Type.Equal(oidASN1Topic) {
			return v.Value.(string), nil
		}
	}

	return "", errors.New("Could not find Push Topic in the provided pkcs12 bundle.")
}
