package mdm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/groob/plist"
)

/*
Last Verified against PDF: April 7 2018
Fully tested commands (with a 10.13.4 mac):
 - ProfileList
 - InstallProfile
 - RemoveProfile
 - ProvisioningProfileList
 - CertificateList

--
Not tested end to end but checked against pdf:
  - InstallProvisioningProfile
  - RemoveProvisioningProfile
*/

func TestMarshalCommand(t *testing.T) {
	var tests = []struct {
		Command Command
	}{
		{
			Command: Command{
				RequestType: "ProfileList",
			},
		},
		{
			Command: Command{
				RequestType: "InstallProfile",
				InstallProfile: &InstallProfile{
					Payload: []byte("foobarbaz"),
				},
			},
		},
		{
			Command: Command{
				RequestType: "RemoveProfile",
				RemoveProfile: &RemoveProfile{
					Identifier: "foobarbaz",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Command.RequestType+"_json", func(t *testing.T) {
			payload := CommandPayload{CommandUUID: "abcd", Command: &tt.Command}
			buf := new(bytes.Buffer)
			enc := json.NewEncoder(buf)
			enc.SetIndent("", "  ")
			if err := enc.Encode(&payload); err != nil {
				t.Fatal(err)
			}
			fmt.Println(buf.String())
		})

		t.Run(tt.Command.RequestType+"_plist", func(t *testing.T) {
			payload := CommandPayload{CommandUUID: "abcd", Command: &tt.Command}
			buf := new(bytes.Buffer)
			enc := plist.NewEncoder(buf)
			enc.Indent("  ")
			if err := enc.Encode(&payload); err != nil {
				t.Fatal(err)
			}
			fmt.Println(buf.String())
		})
	}
}

func TestUnmarshalCommandPayload(t *testing.T) {
	var tests = []struct {
		RequestType string
	}{
		{RequestType: "InstallProfile"},
	}

	for _, tt := range tests {
		t.Run(tt.RequestType+"_json", func(t *testing.T) {
			filename := fmt.Sprintf("%s.json", tt.RequestType)
			data := mustLoadFile(t, filename)
			var payload CommandPayload
			testCommandUnmarshal(t, tt.RequestType, json.Unmarshal, data, &payload)
		})

		t.Run(tt.RequestType+"_plist", func(t *testing.T) {
			filename := fmt.Sprintf("%s.plist", tt.RequestType)
			data := mustLoadFile(t, filename)
			var payload CommandPayload
			testCommandUnmarshal(t, tt.RequestType, plist.Unmarshal, data, &payload)
		})
	}
}

func mustLoadFile(t *testing.T, filename string) []byte {
	t.Helper()
	data, err := ioutil.ReadFile(filepath.Join("testdata", filename))
	if err != nil {
		t.Fatalf("could not read test file %s: %s", filename, err)
	}
	return data
}

type unmarshalFunc func([]byte, interface{}) error

func testCommandUnmarshal(
	t *testing.T,
	requestType string,
	unmarshal unmarshalFunc,
	data []byte,
	payload *CommandPayload,
) {
	t.Helper()
	if err := unmarshal(data, payload); err != nil {
		t.Fatalf("unmarshal command type %s: %s", requestType, err)
	}

	if payload.CommandUUID == "" {
		t.Errorf("missing CommandUUID")
	}

	if have, want := payload.Command.RequestType, requestType; have != want {
		t.Errorf("have %s, want %s", have, want)
	}
}

func TestEndToEnd(t *testing.T) {
	// given an request that came over http as JSON
	requestBytes := []byte(`{"udid": "BC5E2DA4-7FB6-5E70-9928-4981680DAFBF", "payload":"PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPCFET0NUWVBFIHBsaXN0IFBVQkxJQyAiLS8vQXBwbGUvL0RURCBQTElTVCAxLjAvL0VOIiAiaHR0cDovL3d3dy5hcHBsZS5jb20vRFREcy9Qcm9wZXJ0eUxpc3QtMS4wLmR0ZCI+CjxwbGlzdCB2ZXJzaW9uPSIxLjAiPgo8ZGljdD4KCTxrZXk+UGF5bG9hZENvbnRlbnQ8L2tleT4KCTxhcnJheT4KCQk8ZGljdD4KCQkJPGtleT5QYXlsb2FkQ29udGVudDwva2V5PgoJCQk8ZGljdD4KCQkJCTxrZXk+Y29tLmFwcGxlLmFzc2lzdGFudC5zdXBwb3J0PC9rZXk+CgkJCQk8ZGljdD4KCQkJCQk8a2V5PkZvcmNlZDwva2V5PgoJCQkJCTxhcnJheT4KCQkJCQkJPGRpY3Q+CgkJCQkJCQk8a2V5Pm1jeF9wcmVmZXJlbmNlX3NldHRpbmdzPC9rZXk+CgkJCQkJCQk8ZGljdD4KCQkJCQkJCQk8a2V5PkFzc2lzdGFudCBFbmFibGVkPC9rZXk+CgkJCQkJCQkJPGZhbHNlLz4KCQkJCQkJCTwvZGljdD4KCQkJCQkJPC9kaWN0PgoJCQkJCTwvYXJyYXk+CgkJCQk8L2RpY3Q+CgkJCTwvZGljdD4KCQkJPGtleT5QYXlsb2FkRW5hYmxlZDwva2V5PgoJCQk8dHJ1ZS8+CgkJCTxrZXk+UGF5bG9hZElkZW50aWZpZXI8L2tleT4KCQkJPHN0cmluZz5NQ1hUb1Byb2ZpbGUuOWM3MzgwZDItNWJmZS00ZTYwLWJhZDMtMjVhZDg2ZDYxNTBkLmFsYWNhcnRlLmN1c3RvbXNldHRpbmdzLmZiOTU4ZDk2LWE0MzEtNDUzNi04NGQwLTFiZTQ4MjM4NWZiMjwvc3RyaW5nPgoJCQk8a2V5PlBheWxvYWRUeXBlPC9rZXk+CgkJCTxzdHJpbmc+Y29tLmFwcGxlLk1hbmFnZWRDbGllbnQucHJlZmVyZW5jZXM8L3N0cmluZz4KCQkJPGtleT5QYXlsb2FkVVVJRDwva2V5PgoJCQk8c3RyaW5nPmZiOTU4ZDk2LWE0MzEtNDUzNi04NGQwLTFiZTQ4MjM4NWZiMjwvc3RyaW5nPgoJCQk8a2V5PlBheWxvYWRWZXJzaW9uPC9rZXk+CgkJCTxpbnRlZ2VyPjE8L2ludGVnZXI+CgkJPC9kaWN0PgoJPC9hcnJheT4KCTxrZXk+UGF5bG9hZERlc2NyaXB0aW9uPC9rZXk+Cgk8c3RyaW5nPlN0b3BzIFNpcmkgZnJvbSBiZWluZyBlbmFibGVkLjwvc3RyaW5nPgoJPGtleT5QYXlsb2FkRGlzcGxheU5hbWU8L2tleT4KCTxzdHJpbmc+RGlzYWJsZSBTaXJpPC9zdHJpbmc+Cgk8a2V5PlBheWxvYWRJZGVudGlmaWVyPC9rZXk+Cgk8c3RyaW5nPkRpc2FibGVTaXJpPC9zdHJpbmc+Cgk8a2V5PlBheWxvYWRPcmdhbml6YXRpb248L2tleT4KCTxzdHJpbmc+PC9zdHJpbmc+Cgk8a2V5PlBheWxvYWRSZW1vdmFsRGlzYWxsb3dlZDwva2V5PgoJPHRydWUvPgoJPGtleT5QYXlsb2FkU2NvcGU8L2tleT4KCTxzdHJpbmc+U3lzdGVtPC9zdHJpbmc+Cgk8a2V5PlBheWxvYWRUeXBlPC9rZXk+Cgk8c3RyaW5nPkNvbmZpZ3VyYXRpb248L3N0cmluZz4KCTxrZXk+UGF5bG9hZFVVSUQ8L2tleT4KCTxzdHJpbmc+OWM3MzgwZDItNWJmZS00ZTYwLWJhZDMtMjVhZDg2ZDYxNTBkPC9zdHJpbmc+Cgk8a2V5PlBheWxvYWRWZXJzaW9uPC9rZXk+Cgk8aW50ZWdlcj4xPC9pbnRlZ2VyPgo8L2RpY3Q+CjwvcGxpc3Q+Cg==", "request_type": "InstallProfile"}`)

	// unmarshal the request into a go structure
	var req CommandRequest
	if err := json.Unmarshal(requestBytes, &req); err != nil {
		t.Fatal(err)
	}
	if len(req.Command.InstallProfile.Payload) == 0 {
		t.Fatal("InstallProfile payload is empty after json unmarshal")
	}

	// create a payload from the request
	payload, err := NewCommandPayload(&req)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(payload.Command.InstallProfile.Payload, req.Command.InstallProfile.Payload) {
		t.Fatal("mdm payload and request do not have the same payload data")
	}

	// marshal to proto (for storage)
	data, err := MarshalCommandPayload(payload)
	if err != nil {
		t.Fatal(err)
	}

	// unmarshal the proto back into go (from storage)
	var newPayload CommandPayload
	err = UnmarshalCommandPayload(data, &newPayload)
	if err != nil {
		t.Fatal(err)
	}
	if len(newPayload.Command.InstallProfile.Payload) == 0 {
		t.Fatal("unmarshaled proto payload is missing payload")
	}

	// marshal it into a plist to send to the device
	pd, err := plist.MarshalIndent(newPayload, "  ")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(pd, []byte(`PD94bWwgdm`)) {
		t.Fatal("marshaled plist does not contain the required payload")
	}
}
