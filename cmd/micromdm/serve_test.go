package main

import "testing"

func TestLoadPushCerts(t *testing.T) {
	keypath := "testdata/ProviderPrivateKey.key"
	certpath := "testdata/pushcert.pem"
	p12path := "testdata/pushcert.p12"
	keysecret := "secret"

	cfg := &server{
		APNSPrivateKeyPath:  keypath,
		APNSCertificatePath: certpath,
		APNSPrivateKeyPass:  keysecret,
	}

	// test separate key and cert
	cfg.loadPushCerts()
	if cfg.err != nil {
		t.Fatal(cfg.err)
	}

	// test p12 with secret
	cfg.APNSCertificatePath = p12path
	cfg.APNSPrivateKeyPath = ""
	cfg.loadPushCerts()
	if cfg.err != nil {
		t.Fatal(cfg.err)
	}

}
