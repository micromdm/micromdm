package mdm

import (
	"bytes"
	"context"
	"net/http/httptest"
	"testing"
)

func Test_decodeCheckinRequest(t *testing.T) {
	// test that url values from checkin and acknowledge requests are passed to the event.
	req := httptest.NewRequest("GET", "/mdm/checkin?id=1111", bytes.NewReader([]byte(sampleCheckinRequest)))
	dec := &requestDecoder{}
	resp, err := dec.decodeCheckinRequest(context.Background(), req)
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
