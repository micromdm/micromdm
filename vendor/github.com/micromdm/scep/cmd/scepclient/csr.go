package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
)

const (
	csrPEMBlockType = "CERTIFICATE REQUEST"
)

type csrOptions struct {
	cn, org, country, ou, locality, province, challenge string
	key                                                 *rsa.PrivateKey
}

func loadOrMakeCSR(path string, opts *csrOptions) (*x509.CertificateRequest, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		if os.IsExist(err) {
			return loadCSRfromFile(path)
		}
		return nil, err
	}
	defer file.Close()

	csrBytes, err := newCSR(opts.key, opts.ou, opts.locality, opts.province, opts.country, opts.cn, opts.org)
	if err != nil {
		return nil, err
	}
	pemBlock := &pem.Block{
		Type:    csrPEMBlockType,
		Headers: nil,
		Bytes:   csrBytes,
	}
	if err := pem.Encode(file, pemBlock); err != nil {
		return nil, err
	}
	return x509.ParseCertificateRequest(csrBytes)
}

// create a CSR using the same parameters as Keychain Access would produce
func newCSR(priv *rsa.PrivateKey, ou string, locality string, province string, country string, cname, org string) ([]byte, error) {
	subj := pkix.Name{
		CommonName: cname,
	}
	if len(org) > 0 {
		subj.Organization = []string{org}
	}
	if len(ou) > 0 {
		subj.OrganizationalUnit = []string{ou}
	}
	if len(province) > 0 {
		subj.Province = []string{province}
	}
	if len(locality) > 0 {
		subj.Locality = []string{locality}
	}
	if len(country) > 0 {
		subj.Country = []string{country}
	}
	template := &x509.CertificateRequest{
		Subject: subj,
	}
	return x509.CreateCertificateRequest(rand.Reader, template, priv)
}

// convert DER to PEM format
func pemCSR(derBytes []byte) []byte {
	pemBlock := &pem.Block{
		Type:    csrPEMBlockType,
		Headers: nil,
		Bytes:   derBytes,
	}
	out := pem.EncodeToMemory(pemBlock)
	return out
}

// load PEM encoded CSR from file
func loadCSRfromFile(path string) (*x509.CertificateRequest, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	pemBlock, _ := pem.Decode(data)
	if pemBlock == nil {
		return nil, errors.New("cannot find the next PEM formatted block")
	}
	if pemBlock.Type != csrPEMBlockType || len(pemBlock.Headers) != 0 {
		return nil, errors.New("unmatched type or headers")
	}
	return x509.ParseCertificateRequest(pemBlock.Bytes)
}
