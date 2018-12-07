package main

import "testing"

func TestLoadPushCerts(t *testing.T) {
	keypath := "testdata/ProviderPrivateKey.key"
	certpath := "testdata/pushcert.pem"
	p12path := "testdata/pushcert.p12"
	keysecret := "secret"

	_, _, err := loadPushCerts(certpath, keypath, keysecret)
	if err != nil {
		t.Errorf("failed to load PEM push certs with err %s", err)
	}

	// try to load from p12
	_, _, err = loadPushCerts(p12path, "", keysecret)
	if err != nil {
		t.Errorf("failed to load p12 push certs with err %s", err)
	}
}
