package scepserver_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	kitlog "github.com/go-kit/kit/log"
	"golang.org/x/net/context"

	"github.com/micromdm/scep/server"
)

func TestCACaps(t *testing.T) {
	server, _ := newServer(t)
	defer server.Close()
	url := server.URL + "/scep?operation=GetCACaps"
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Error("expected", http.StatusOK, "got", resp.StatusCode)
	}
}

func TestPKIOperation(t *testing.T) {
	server, _ := newServer(t)
	defer server.Close()
	pkcsreq := loadTestFile(t, "../scep/testdata/PKCSReq.der")
	body := bytes.NewReader(pkcsreq)
	url := server.URL + "/scep?operation=PKIOperation"
	resp, err := http.Post(url, "", body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Error("expected", http.StatusOK, "got", resp.StatusCode)
	}
}

func newServer(t *testing.T, opts ...scepserver.ServiceOption) (*httptest.Server, scepserver.Service) {
	var err error
	var depot scepserver.Depot // cert storage
	{
		depot, err = scepserver.NewFileDepot("../scep/testdata/testca")
		if err != nil {
			t.Fatal(err)
		}
	}
	var svc scepserver.Service // scep service
	{
		svc, err = scepserver.NewService(depot, opts...)
		if err != nil {
			t.Fatal(err)
		}
	}
	ctx := context.Background()
	logger := kitlog.NewNopLogger()
	handler := scepserver.ServiceHandler(ctx, svc, logger)
	server := httptest.NewServer(handler)
	return server, svc
}

/* helpers */
const (
	rsaPrivateKeyPEMBlockType = "RSA PRIVATE KEY"
	certificatePEMBlockType   = "CERTIFICATE"
)

func loadTestFile(t *testing.T, path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// create a new RSA private key
func newRSAKey(bits int) (*rsa.PrivateKey, error) {
	private, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	return private, nil
}

// create a CSR using the same parameters as Keychain Access would produce
func newCSR(priv *rsa.PrivateKey, email, country, cname string) ([]byte, error) {
	subj := pkix.Name{
		Country:    []string{country},
		CommonName: cname,
		ExtraNames: []pkix.AttributeTypeAndValue{pkix.AttributeTypeAndValue{
			Type:  []int{1, 2, 840, 113549, 1, 9, 1},
			Value: email,
		}},
	}
	template := &x509.CertificateRequest{
		Subject: subj,
	}
	return x509.CreateCertificateRequest(rand.Reader, template, priv)
}
func loadCACredentials(t *testing.T) (*x509.Certificate, *rsa.PrivateKey) {
	cert, err := loadCertFromFile("../scep/testdata/testca/ca.crt")
	if err != nil {
		t.Fatal(err)
	}
	key, err := loadKeyFromFile("../scep/testdata/testca/ca.key")
	if err != nil {
		t.Fatal(err)
	}
	return cert, key
}

func loadClientCredentials(t *testing.T) (*x509.Certificate, *rsa.PrivateKey) {
	cert, err := loadCertFromFile("../scep/testdata/testclient/client.pem")
	if err != nil {
		t.Fatal(err)
	}
	key, err := loadKeyFromFile("../scep/testdata/testclient/client.key")
	if err != nil {
		t.Fatal(err)
	}
	return cert, key
}

func loadCertFromFile(path string) (*x509.Certificate, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	pemBlock, _ := pem.Decode(data)
	if pemBlock == nil {
		return nil, errors.New("PEM decode failed")
	}
	if pemBlock.Type != certificatePEMBlockType {
		return nil, errors.New("unmatched type or headers")
	}
	return x509.ParseCertificate(pemBlock.Bytes)
}

// load an encrypted private key from disk
func loadKeyFromFile(path string) (*rsa.PrivateKey, error) {
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

	// testca key has a password
	if len(pemBlock.Headers) > 0 {
		password := []byte("")
		b, err := x509.DecryptPEMBlock(pemBlock, password)
		if err != nil {
			return nil, err
		}
		return x509.ParsePKCS1PrivateKey(b)
	}

	return x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
}
