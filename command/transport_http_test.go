package command

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecodeRequest(t *testing.T) {
	requestData := `
{
    "request_type": "InstallApplication",
    "udid" : "564D38A0-4C3B-AD69-803B-DAC58A298191",
    "manifest_url" : "https://mdm.acme.co/repo/munkitools-3.0.0.3298.plist",
    "management_flags" : 1
}
`
	req := httptest.NewRequest("POST", "https://mdm.acme.co/v1/commands", strings.NewReader(requestData))
	request, err := decodeRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	decoded := request.(newCommandRequest)

	if have, want := decoded.RequestType, "InstallApplication"; have != want {
		t.Errorf("have %s, want %s", have, want)
	}

	if have, want := decoded.CommandRequest.InstallApplication.ManifestURL,
		"https://mdm.acme.co/repo/munkitools-3.0.0.3298.plist"; have != want {
		t.Errorf("have %s, want %s", have, want)
	}

}
