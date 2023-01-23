package crypto

import (
	"crypto/x509"
	"encoding/pem"
	"testing"
)

func TestCmpAppleDeviceCAPublicKey(t *testing.T) {
	block, _ := pem.Decode([]byte(appleiPhoneDeviceCAPEM))
	if block == nil || block.Type != "CERTIFICATE" {
		t.Fatalf("appleiPhoneDeviceCAPEM: invalid PEM block")
	}
	parent, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("appleiPhoneDeviceCAPEM: err parsing: %v", err)
	}

	if !appleiPhoneDeviceCAPublicKey.Equal(parent.PublicKey) {
		t.Error("appleiPhoneDeviceCAPEM: public keys not equal")
	}
}

func TestVerifyFromAppleDeviceCA(t *testing.T) {
	// extracted from /Library/Keychains/apsd.keychain
	cert, _ := ReadPEMCertificateFile("testdata/appleca_signed.pem")

	if err := VerifyFromAppleDeviceCA(cert); err != nil {
		t.Errorf("received error verifying Apple Device CA: %v", err)
	}

	// test invalidly-signed cert
	// reuse cert so we don't need to generate one. This should fail since it's not self-signed.
	block, _ := pem.Decode([]byte(appleiPhoneDeviceCAPEM))
	if block == nil || block.Type != "CERTIFICATE" {
		t.Fatalf("appleiPhoneDeviceCAPEM: invalid PEM block")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("appleiPhoneDeviceCAPEM: err parsing: %v", err)
	}

	if err := VerifyFromAppleDeviceCA(cert); err == nil {
		t.Error("expected error verifying non-apple-device-signed cert")
	}
}
