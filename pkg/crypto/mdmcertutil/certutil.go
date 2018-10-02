// Package mdmcertutil contains helpers for requesting MDM Push Certifificates.
// The process is described by Apple at
// https://developer.apple.com/library/content/documentation/Miscellaneous/Reference/MobileDeviceManagementProtocolRef/7-MDMVendorCSRSigningOverview/MDMVendorCSRSigningOverview.html#//apple_ref/doc/uid/TP40017387-CH6-SW4
package mdmcertutil

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/groob/plist"
	"github.com/pkg/errors"
)

// CSRConfig defines arguments required to create a new CSR.
type CSRConfig struct {
	CommonName, Country, Email string
	PrivateKeyPassword         []byte
	PrivateKeyPath, CSRPath    string
}

// CreateCSR creates a new private key and CSR, saving both as PEM encoded files.
func CreateCSR(req *CSRConfig) error {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	pemKey, err := encryptedKey(key, req.PrivateKeyPassword)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(req.PrivateKeyPath, pemKey, 0600); err != nil {
		return err
	}

	derBytes, err := NewCSR(key, strings.ToLower(req.Email), strings.ToUpper(req.Country), req.CommonName)
	if err != nil {
		return err
	}
	pemCSR := PemCSR(derBytes)
	return ioutil.WriteFile(req.CSRPath, pemCSR, 0600)
}

// The PushCertificateRequest structure required by identity.apple.com
// to create an MDM Push certificate.
type PushCertificateRequest struct {
	PushCertRequestCSR       string
	PushCertCertificateChain string
	PushCertSignature        string
}

// Encode marshals a PushCertificateRequest to an XML Plist file and returns a base64 encoded byte representation
// of the request.
func (p *PushCertificateRequest) Encode() ([]byte, error) {
	data, err := plist.MarshalIndent(p, "  ")
	if err != nil {
		return nil, errors.Wrap(err, "marshal PushCertificateRequest")
	}
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(encoded, data)
	return encoded, nil
}

const (
	wwdrIntermediaryURL = "https://developer.apple.com/certificationauthority/AppleWWDRCA.cer"
	appleRootCAURL      = "http://www.apple.com/appleca/AppleIncRootCertificate.cer"
)

// CreatePushCertificateRequest creates a request structure required by identity.apple.com.
// It requires a "MDM CSR" certificate (the vendor certificate), a push CSR (the customer specific CSR),
// and the vendor private key.
func CreatePushCertificateRequest(mdmCertPath, pushCSRPath, pKeyPath string, pKeyPass []byte) (*PushCertificateRequest, error) {
	// private key of the mdm vendor cert
	key, err := loadKeyFromFile(pKeyPath, pKeyPass)
	if err != nil {
		return nil, errors.Wrapf(err, "load private key from %s", pKeyPath)
	}

	// push csr
	csr, err := loadCSRfromFile(pushCSRPath)
	if err != nil {
		return nil, errors.Wrapf(err, "load push CSR from %s", pushCSRPath)
	}

	// csr signature
	signature, err := signPushCSR(csr.Raw, key)
	if err != nil {
		return nil, errors.Wrapf(err, "sign push CSR with private key")
	}

	// vendor cert
	mdmCertBytes, err := loadDERCertFromFile(mdmCertPath)
	if err != nil {
		return nil, errors.Wrapf(err, "load vendor certificate from path %s", mdmCertPath)
	}
	mdmPEM := pemCert(mdmCertBytes)

	// wwdr cert
	wwdrCertBytes, err := loadCertfromHTTP(wwdrIntermediaryURL)
	if err != nil {
		return nil, errors.Wrapf(err, "load WWDR certificate from %s", wwdrIntermediaryURL)
	}
	wwdrPEM := pemCert(wwdrCertBytes)

	// apple root certificate
	rootCertBytes, err := loadCertfromHTTP(appleRootCAURL)
	if err != nil {
		return nil, errors.Wrapf(err, "load root certificate from %s", appleRootCAURL)
	}
	rootPEM := pemCert(rootCertBytes)

	csrB64 := base64.StdEncoding.EncodeToString(csr.Raw)
	sig64 := base64.StdEncoding.EncodeToString(signature)
	pushReq := &PushCertificateRequest{
		PushCertRequestCSR:       csrB64,
		PushCertCertificateChain: makeCertChain(mdmPEM, wwdrPEM, rootPEM),
		PushCertSignature:        sig64,
	}
	return pushReq, nil
}

func makeCertChain(mdmPEM, wwdrPEM, rootPEM []byte) string {
	return string(mdmPEM) + string(wwdrPEM) + string(rootPEM)
}

func signPushCSR(csrData []byte, key *rsa.PrivateKey) ([]byte, error) {
	h := sha1.New()
	_, err := h.Write(csrData)
	if err != nil {
		return nil, errors.Wrap(err, "writing CSR data")
	}
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA1, h.Sum(nil))
	return signature, errors.Wrap(err, "signing push CSR")
}

const (
	csrPEMBlockType = "CERTIFICATE REQUEST"
)

// create a CSR using the same parameters as Keychain Access would produce
func NewCSR(priv *rsa.PrivateKey, email, country, cname string) ([]byte, error) {
	subj := pkix.Name{
		Country:    []string{country},
		CommonName: cname,
		ExtraNames: []pkix.AttributeTypeAndValue{{
			Type:  []int{1, 2, 840, 113549, 1, 9, 1},
			Value: email,
		}},
	}
	template := &x509.CertificateRequest{
		Subject: subj,
	}
	return x509.CreateCertificateRequest(rand.Reader, template, priv)
}

// convert DER to PEM format
func PemCSR(derBytes []byte) []byte {
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

const (
	rsaPrivateKeyPEMBlockType = "RSA PRIVATE KEY"
)

// protect an rsa key with a password
func encryptedKey(key *rsa.PrivateKey, password []byte) ([]byte, error) {
	privBytes := x509.MarshalPKCS1PrivateKey(key)
	privPEMBlock, err := x509.EncryptPEMBlock(rand.Reader, rsaPrivateKeyPEMBlockType, privBytes, password, x509.PEMCipher3DES)
	if err != nil {
		return nil, err
	}

	out := pem.EncodeToMemory(privPEMBlock)
	return out, nil
}

// load an encrypted private key from disk
func loadKeyFromFile(path string, password []byte) (*rsa.PrivateKey, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	pemBlock, _ := pem.Decode(data)
	if pemBlock == nil {
		return nil, errors.New("PEM decode failed")
	}
	if pemBlock.Type != rsaPrivateKeyPEMBlockType {
		return nil, errors.New("unmatched type or headers")
	}

	if string(password) != "" {
		b, err := x509.DecryptPEMBlock(pemBlock, password)
		if err != nil {
			return nil, err
		}
		return x509.ParsePKCS1PrivateKey(b)
	}
	return x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
}

const (
	certificatePEMBlockType = "CERTIFICATE"
)

func pemCert(derBytes []byte) []byte {
	pemBlock := &pem.Block{
		Type:    certificatePEMBlockType,
		Headers: nil,
		Bytes:   derBytes,
	}
	out := pem.EncodeToMemory(pemBlock)
	return out
}

func loadDERCertFromFile(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	crt, err := x509.ParseCertificate(data)
	if err != nil {
		return nil, err
	}
	return crt.Raw, nil
}

func loadCertfromHTTP(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "create GET request for %s", url)
	}
	req.Header.Set("Accept", "*/*") // required by Apple at some point.

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "GET request to %s", url)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got %s when trying to http.Get %s", resp.Status, url)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "reading get certificate request response body")
	}

	crt, err := x509.ParseCertificate(data)
	if err != nil {
		return nil, errors.Wrap(err, "parse wwdr intermediate certificate")
	}
	return crt.Raw, nil
}
