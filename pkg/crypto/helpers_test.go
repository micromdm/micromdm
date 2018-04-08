package crypto

import (
	"testing"
)

func TestTopicFromValidCert(t *testing.T) {
	certificate, _ := ReadPEMCertificateFile("testdata/mock_push_cert.pem")
	pushTopic, err := TopicFromCert(certificate)
	if err != nil {
		t.Fatalf("fail %s", err)
	}

	if have, want := pushTopic, "com.apple.mgmt.External.18a16429-886b-41f1-9c30-2bd04ae4fc37"; have != want {
		t.Errorf("have %s, want %s", have, want)
	}
}

func TestTopicFromInvalidCert(t *testing.T) {
	certFileNames := []string{"mock_push_cert_wrong_uid_prefix.pem", "mock_push_cert_wrong_uid_prefix.pem"}
	for _, certFileName := range certFileNames {
		certificate, _ := ReadPEMCertificateFile("testdata/" + certFileName)
		_, err := TopicFromCert(certificate)

		if err == nil {
			t.Errorf("expected error for invalid push certificate")
		}
	}
}
