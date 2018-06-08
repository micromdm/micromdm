package mdm

import (
	"bytes"
	"context"
	"encoding/base64"
	"net/http/httptest"
	"testing"

	"github.com/fullsailor/pkcs7"
	"github.com/micromdm/micromdm/pkg/crypto"
)

// immitate an Mdm-Signature header
func mdmSignRequest(body []byte) (string, error) {
	key, cert, err := crypto.SimpleSelfSignedRSAKeypair("test", 365)
	if err != nil {
		return "", err
	}

	sd, err := pkcs7.NewSignedData(body)
	if err != nil {
		return "", err
	}

	sd.AddSigner(cert, key, pkcs7.SignerInfoConfig{})
	sd.Detach()
	sig, err := sd.Finish()
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(sig), nil
}

func Test_decodeCheckinRequest(t *testing.T) {
	// test that url values from checkin and acknowledge requests are passed to the event.
	req := httptest.NewRequest("GET", "/mdm/checkin?id=1111", bytes.NewReader([]byte(sampleCheckinRequest)))
	b64sig, err := mdmSignRequest([]byte(sampleCheckinRequest))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Mdm-Signature", b64sig)
	resp, err := decodeCheckinRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	response := resp.(checkinRequest)
	if have, want := response.Event.Params["id"], "1111"; have != want {
		t.Errorf("have %s, want %s", have, want)
	}
}

const sampleCheckinRequest = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>AwaitingConfiguration</key>
	<false/>
	<key>MessageType</key>
	<string>TokenUpdate</string>
	<key>PushMagic</key>
	<string>AB62BB8A-7757-4130-94CC-CC8C5333D481</string>
	<key>Token</key>
	<data>
	YWJjZGUK
	</data>
	<key>Topic</key>
	<string>com.apple.mgmt.External.80bb2169-e864-4685-9a96-faa734f0b978</string>
	<key>UDID</key>
	<string>BC5E2DA4-7FB6-5E70-9928-4981680DAFBF</string>
</dict>
</plist>`
