package mdm

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/base64"
	"net/http/httptest"
	"testing"

	"github.com/micromdm/micromdm/pkg/crypto"

	"go.mozilla.org/pkcs7"
)

// imitate a Mdm-Signature header
func mdmSignRequest(body []byte) (*x509.Certificate, string, error) {
	key, cert, err := crypto.SimpleSelfSignedRSAKeypair("test", 365)
	if err != nil {
		return nil, "", err
	}

	sd, err := pkcs7.NewSignedData(body)
	if err != nil {
		return nil, "", err
	}

	sd.AddSigner(cert, key, pkcs7.SignerInfoConfig{})
	sd.Detach()
	sig, err := sd.Finish()
	if err != nil {
		return nil, "", err
	}

	return cert, base64.StdEncoding.EncodeToString(sig), nil
}

func Test_mdmMdmSignatureHeader(t *testing.T) {
	// test that url values from checkin and acknowledge requests are passed to the event.
	req := httptest.NewRequest("GET", "/mdm/checkin", bytes.NewReader([]byte(sampleCheckinRequest)))
	cert, b64sig, err := mdmSignRequest([]byte(sampleCheckinRequest))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Mdm-Signature", b64sig)

	ctx := context.Background()
	ctx = (verifier{PKCS7Verifier: &crypto.PKCS7Verifier{}}).populateDeviceCertificateFromSignRequestHeader(ctx, req)

	reqcert, err := DeviceCertificateFromContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if 0 != bytes.Compare(cert.Raw, reqcert.Raw) {
		t.Error("certificate mismatch")
	}
}
