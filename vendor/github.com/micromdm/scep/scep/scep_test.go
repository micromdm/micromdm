package scep_test

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"math/big"
	"testing"
	"time"

	"github.com/micromdm/scep/scep"
)

func testParsePKIMessage(t *testing.T, data []byte) *scep.PKIMessage {
	msg, err := scep.ParsePKIMessage(data)
	if err != nil {
		t.Fatal(err)
	}
	if msg.TransactionID == "" {
		t.Errorf("expected TransactionID attribute")
	}
	if msg.MessageType == "" {
		t.Errorf("expected MessageType attribute")
	}
	switch msg.MessageType {
	case scep.CertRep:
		if len(msg.RecipientNonce) == 0 {
			t.Errorf("expected RecipientNonce attribute")
		}
	case scep.PKCSReq, scep.UpdateReq, scep.RenewalReq:
		if len(msg.SenderNonce) == 0 {
			t.Errorf("expected SenderNonce attribute")
		}
	}
	return msg
}

func TestDecryptPKIEnvelopeCSR(t *testing.T) {
	pkcsReq := loadTestFile(t, "testdata/PKCSReq.der")
	msg := testParsePKIMessage(t, pkcsReq)
	cacert, cakey := loadCACredentials(t)
	err := msg.DecryptPKIEnvelope(cacert, cakey)
	if err != nil {
		t.Fatal(err)
	}
	if msg.CSRReqMessage.CSR == nil {
		t.Errorf("expected non-nil CSR field")
	}
}

func TestDecryptPKIEnvelopeCert(t *testing.T) {
	certRep := loadTestFile(t, "testdata/CertRep.der")
	testParsePKIMessage(t, certRep)
	// clientcert, clientkey := loadClientCredentials(t)
	// err = msg.DecryptPKIEnvelope(clientcert, clientkey)
	// if err != nil {
	// 	t.Fatal(err)
	// }
}

func TestSignCSR(t *testing.T) {
	pkcsReq := loadTestFile(t, "testdata/PKCSReq.der")
	msg := testParsePKIMessage(t, pkcsReq)
	cacert, cakey := loadCACredentials(t)
	err := msg.DecryptPKIEnvelope(cacert, cakey)
	if err != nil {
		t.Fatal(err)
	}
	csr := msg.CSRReqMessage.CSR
	id, err := GenerateSubjectKeyID(csr.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(4),
		Subject:      csr.Subject,
		NotBefore:    time.Now().Add(-600).UTC(),
		NotAfter:     time.Now().AddDate(1, 0, 0).UTC(),
		SubjectKeyId: id,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageAny,
			x509.ExtKeyUsageClientAuth,
		},
	}
	certRep, err := msg.SignCSR(cacert, cakey, tmpl)
	if err != nil {
		t.Fatal(err)
	}
	testParsePKIMessage(t, certRep.Raw)
}

func TestNewCSRRequest(t *testing.T) {
	key, err := newRSAKey(2048)
	if err != nil {
		t.Fatal(err)
	}
	derBytes, err := newCSR(key, "john.doe@example.com", "US", "com.apple.scep.2379B935-294B-4AF1-A213-9BD44A2C6688")
	if err != nil {
		t.Fatal(err)
	}
	csr, err := x509.ParseCertificateRequest(derBytes)
	if err != nil {
		t.Fatal(err)
	}
	clientcert, clientkey := loadClientCredentials(t)
	cacert, cakey := loadCACredentials(t)
	tmpl := &scep.PKIMessage{
		MessageType: scep.PKCSReq,
		Recipients:  []*x509.Certificate{cacert},
		SignerCert:  clientcert,
		SignerKey:   clientkey,
	}

	pkcsreq, err := scep.NewCSRRequest(csr, tmpl)
	if err != nil {
		t.Fatal(err)
	}
	msg := testParsePKIMessage(t, pkcsreq.Raw)
	err = msg.DecryptPKIEnvelope(cacert, cakey)
	if err != nil {
		t.Fatal(err)
	}
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

func loadTestFile(t *testing.T, path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func loadCACredentials(t *testing.T) (*x509.Certificate, *rsa.PrivateKey) {
	cert, err := loadCertFromFile("testdata/testca/ca.crt")
	if err != nil {
		t.Fatal(err)
	}
	key, err := loadKeyFromFile("testdata/testca/ca.key")
	if err != nil {
		t.Fatal(err)
	}
	return cert, key
}

func loadClientCredentials(t *testing.T) (*x509.Certificate, *rsa.PrivateKey) {
	cert, err := loadCertFromFile("testdata/testclient/client.pem")
	if err != nil {
		t.Fatal(err)
	}
	key, err := loadKeyFromFile("testdata/testclient/client.key")
	if err != nil {
		t.Fatal(err)
	}
	return cert, key
}

const (
	rsaPrivateKeyPEMBlockType = "RSA PRIVATE KEY"
	certificatePEMBlockType   = "CERTIFICATE"
)

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

// rsaPublicKey reflects the ASN.1 structure of a PKCS#1 public key.
type rsaPublicKey struct {
	N *big.Int
	E int
}

// GenerateSubjectKeyID generates SubjectKeyId used in Certificate
// ID is 160-bit SHA-1 hash of the value of the BIT STRING subjectPublicKey
func GenerateSubjectKeyID(pub crypto.PublicKey) ([]byte, error) {
	var pubBytes []byte
	var err error
	switch pub := pub.(type) {
	case *rsa.PublicKey:
		pubBytes, err = asn1.Marshal(rsaPublicKey{
			N: pub.N,
			E: pub.E,
		})
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("only RSA public key is supported")
	}

	hash := sha1.Sum(pubBytes)

	return hash[:], nil
}
