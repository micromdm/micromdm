package server

import "testing"

func TestLoadPushCerts(t *testing.T) {
	keypath := "testdata/ProviderPrivateKey.key"
	certpath := "testdata/pushcert.pem"
	p12path := "testdata/pushcert.p12"
	keysecret := "secret"

	s := &Server{
		APNSPrivateKeyPath:  keypath,
		APNSCertificatePath: certpath,
		APNSPrivateKeyPass:  keysecret,
	}

	// test separate key and cert
	err := s.loadPushCerts()
	if err != nil {
		t.Fatal(err)
	}

	// test p12 with secret
	s.APNSCertificatePath = p12path
	s.APNSPrivateKeyPath = ""
	err = s.loadPushCerts()
	if err != nil {
		t.Fatal(err)
	}

}
