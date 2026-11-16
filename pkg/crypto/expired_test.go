package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"testing"
	"time"

	"github.com/smallstep/pkcs7"
)

type testPKI struct {
	ca     *x509.Certificate
	caKey  *rsa.PrivateKey
	device *x509.Certificate
	devKey *rsa.PrivateKey
}

func newCert(t *testing.T, tmpl *x509.Certificate, parent *x509.Certificate, pub *rsa.PublicKey, signer *rsa.PrivateKey) *x509.Certificate {
	t.Helper()
	der, err := x509.CreateCertificate(rand.Reader, tmpl, parent, pub, signer)
	if err != nil {
		t.Fatal(err)
	}
	c, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// newPKI makes a SCEP-like CA and a device certificate valid from notBefore to notAfter.
func newPKI(t *testing.T, notBefore, notAfter time.Time) testPKI {
	t.Helper()
	caKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	caTmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{Organization: []string{"MicroMDM"}, OrganizationalUnit: []string{"MICROMDM SCEP CA"}},
		NotBefore:             time.Now().Add(-5 * 365 * 24 * time.Hour),
		NotAfter:              time.Now().Add(5 * 365 * 24 * time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
	}
	ca := newCert(t, caTmpl, caTmpl, &caKey.PublicKey, caKey)
	devKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	devTmpl := &x509.Certificate{
		SerialNumber: big.NewInt(42),
		Subject:      pkix.Name{Organization: []string{"Kidsnanny"}, CommonName: "Kidsnanny"},
		NotBefore:    notBefore,
		NotAfter:     notAfter,
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
	}
	return testPKI{ca: ca, caKey: caKey, device: newCert(t, devTmpl, ca, &devKey.PublicKey, caKey), devKey: devKey}
}

// signDetached produces what iOS puts in Mdm-Signature: a detached CMS SignedData with signingTime=now.
func signDetached(t *testing.T, body []byte, cert *x509.Certificate, key *rsa.PrivateKey) []byte {
	t.Helper()
	sd, err := pkcs7.NewSignedData(body)
	if err != nil {
		t.Fatal(err)
	}
	sd.SetDigestAlgorithm(pkcs7.OIDDigestAlgorithmSHA256)
	if err := sd.AddSigner(cert, key, pkcs7.SignerInfoConfig{}); err != nil {
		t.Fatal(err)
	}
	sd.Detach()
	sig, err := sd.Finish()
	if err != nil {
		t.Fatal(err)
	}
	return sig
}

func parsed(t *testing.T, sig, body []byte) *pkcs7.PKCS7 {
	t.Helper()
	p7, err := pkcs7.Parse(sig)
	if err != nil {
		t.Fatal(err)
	}
	p7.Content = body
	return p7
}

func libVerify(t *testing.T, sig, body []byte) error {
	t.Helper()
	return (&PKCS7Verifier{}).Verify(parsed(t, sig, body))
}

var body = []byte(`<?xml version="1.0"?><plist version="1.0"><dict><key>Status</key><string>Idle</string><key>UDID</key><string>00008110-00067519010BA01E</string></dict></plist>`)

func expiredPKI(t *testing.T) testPKI {
	return newPKI(t, time.Now().Add(-400*24*time.Hour), time.Now().Add(-35*24*time.Hour))
}

func TestExpiredSignerIsRejectedByLibraryAndAcceptedByPolicy(t *testing.T) {
	pki := expiredPKI(t)
	sig := signDetached(t, body, pki.device, pki.devKey)

	err := libVerify(t, sig, body)
	if !IsSignerValidityError(err) {
		t.Fatalf("expected the pkcs7 signing-time validity error, got %v", err)
	}
	policy := &ExpiredDeviceCertPolicy{Until: time.Now().Add(time.Hour), Issuer: pki.ca}
	if cert, err := policy.VerifyExpiredSigner(parsed(t, sig, body), time.Now()); err != nil || cert.SerialNumber.Int64() != 42 {
		t.Fatalf("policy should accept a valid signature from an expired CA-issued cert: %v", err)
	}
}

func TestPolicyRejects(t *testing.T) {
	pki := expiredPKI(t)
	sig := signDetached(t, body, pki.device, pki.devKey)
	now := time.Now()
	active := &ExpiredDeviceCertPolicy{Until: now.Add(time.Hour), Issuer: pki.ca}

	other := expiredPKI(t) // a different CA
	notExpired := newPKI(t, now.Add(-time.Hour), now.Add(300*24*time.Hour))
	notExpiredSig := signDetached(t, body, notExpired.device, notExpired.devKey)
	notYetValid := newPKI(t, now.Add(24*time.Hour), now.Add(48*time.Hour))
	notYetValidSig := signDetached(t, body, notYetValid.device, notYetValid.devKey)
	// signed by the expired device cert's identity, but with an attacker's key
	attackerKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	forged := signDetached(t, body, pki.device, attackerKey)

	cases := []struct {
		name   string
		policy *ExpiredDeviceCertPolicy
		p7     *pkcs7.PKCS7
	}{
		{"nil policy", nil, parsed(t, sig, body)},
		{"policy ended", &ExpiredDeviceCertPolicy{Until: now.Add(-time.Minute), Issuer: pki.ca}, parsed(t, sig, body)},
		{"no issuer", &ExpiredDeviceCertPolicy{Until: now.Add(time.Hour)}, parsed(t, sig, body)},
		{"tampered body", active, parsed(t, sig, append([]byte("x"), body...))},
		{"cert from another CA", &ExpiredDeviceCertPolicy{Until: now.Add(time.Hour), Issuer: other.ca}, parsed(t, sig, body)},
		{"cert not expired", &ExpiredDeviceCertPolicy{Until: now.Add(time.Hour), Issuer: notExpired.ca}, parsed(t, notExpiredSig, body)},
		{"cert not yet valid", &ExpiredDeviceCertPolicy{Until: now.Add(time.Hour), Issuer: notYetValid.ca}, parsed(t, notYetValidSig, body)},
		{"signature made with another key", active, parsed(t, forged, body)},
		{"nil signature", active, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := c.policy.VerifyExpiredSigner(c.p7, now); err == nil {
				t.Fatal("expected rejection")
			}
		})
	}
}

func TestValidSignerStillGoesThroughLibrary(t *testing.T) {
	pki := newPKI(t, time.Now().Add(-time.Hour), time.Now().Add(300*24*time.Hour))
	sig := signDetached(t, body, pki.device, pki.devKey)
	if err := libVerify(t, sig, body); err != nil {
		t.Fatalf("valid cert must pass the normal path: %v", err)
	}
	if err := VerifyIgnoringSignerValidity(parsed(t, sig, body), pki.device); err != nil {
		t.Fatalf("manual verifier must agree with the library for valid signatures: %v", err)
	}
}
