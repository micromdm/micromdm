package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"time"
)

func RandomCertificateSerialNumber() (*big.Int, error) {
	limit := new(big.Int).Lsh(big.NewInt(1), 128)
	return rand.Int(rand.Reader, limit)
}

func SimpleSelfSignedRSAKeypair(cn string, days int) (key *rsa.PrivateKey, cert *x509.Certificate, err error) {
	key, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return
	}

	serialNumber, err := RandomCertificateSerialNumber()
	if err != nil {
		return
	}
	timeNow := time.Now()
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: cn,
		},
		NotBefore:             timeNow,
		NotAfter:              timeNow.Add(time.Duration(days) * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{cn},
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return
	}
	cert, err = x509.ParseCertificate(certBytes)
	if err != nil {
		return
	}

	return
}

func WritePEMCertificateFile(cert *x509.Certificate, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return pem.Encode(
		file,
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		})
}

func WritePEMRSAKeyFile(key *rsa.PrivateKey, path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0700)
	if err != nil {
		return err
	}
	defer file.Close()

	return pem.Encode(
		file,
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		})
}
