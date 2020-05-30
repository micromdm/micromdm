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
		{
			Command: Command{
				RequestType: "EnableRemoteDesktop",
			},
		},
		{
			Command: Command{
				RequestType: "DisableRemoteDesktop",
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
		})

		t.Run(tt.Command.RequestType+"_plist", func(t *testing.T) {
			payload := CommandPayload{CommandUUID: "abcd", Command: &tt.Command}
			buf := new(bytes.Buffer)
			enc := plist.NewEncoder(buf)
			enc.Indent("  ")
			if err := enc.Encode(&payload); err != nil {
				t.Fatal(err)
			}
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
	var tests = []struct {
		name         string
		requestBytes []byte
		testFn       func(t *testing.T, parts endToEndParts)
	}{
		{
			name: "Settings_ApplicationConfiguration",
			requestBytes: []byte(
				`{"udid":"BC5E2DA4-7FB6-5E70-9928-4981680DAFBF","request_type":"Settings","settings":[{"item":"ApplicationConfiguration","identifier":"anidentifier","configuration":"PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPCFET0NUWVBFIHBsaXN0IFBVQkxJQyAiLS8vQXBwbGUvL0RURCBQTElTVCAxLjAvL0VOIiAiaHR0cDovL3d3dy5hcHBsZS5jb20vRFREcy9Qcm9wZXJ0eUxpc3QtMS4wLmR0ZCI+CjxwbGlzdCB2ZXJzaW9uPSIxLjAiPgogIDxkaWN0PgogICAgPGtleT5iYXo8L2tleT4KICAgIDxzdHJpbmc+cXV4PC9zdHJpbmc+CiAgICA8a2V5PmNvdW50PC9rZXk+CiAgICA8aW50ZWdlcj4xPC9pbnRlZ2VyPgogICAgPGtleT5mb288L2tleT4KICAgIDxzdHJpbmc+YmFyPC9zdHJpbmc+CiAgPC9kaWN0Pgo8L3BsaXN0Pgo="}]}`,
			),
			testFn: func(t *testing.T, parts endToEndParts) {
				if len(parts.req.Command.Settings.Settings) == 0 {
					t.Error("expected settings command to include at least one setting")
				}

				if len(parts.fromProto.Command.Settings.Settings) == 0 {
					t.Error("expected settings command from proto to include at least one setting")
				}

				// unmarshal plist and check that the settings in the configuration dictionary are there
				var cmd struct {
					Command struct{ Settings []map[string]interface{} }
				}
				if err := plist.Unmarshal(parts.plistData, &cmd); err != nil {
					t.Fatal(err)
				}
				setting := cmd.Command.Settings[0]["Configuration"].(map[string]interface{})
				if have, want := setting["foo"], "bar"; have != want {
					t.Errorf("have key %s, want key %s", have, want)
				}

			},
		},

		{
			name: "InstallProfile",
			requestBytes: []byte(
				`{"udid": "BC5E2DA4-7FB6-5E70-9928-4981680DAFBF", "payload":"PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPCFET0NUWVBFIHBsaXN0IFBVQkxJQyAiLS8vQXBwbGUvL0RURCBQTElTVCAxLjAvL0VOIiAiaHR0cDovL3d3dy5hcHBsZS5jb20vRFREcy9Qcm9wZXJ0eUxpc3QtMS4wLmR0ZCI+CjxwbGlzdCB2ZXJzaW9uPSIxLjAiPgo8ZGljdD4KCTxrZXk+UGF5bG9hZENvbnRlbnQ8L2tleT4KCTxhcnJheT4KCQk8ZGljdD4KCQkJPGtleT5QYXlsb2FkQ29udGVudDwva2V5PgoJCQk8ZGljdD4KCQkJCTxrZXk+Y29tLmFwcGxlLmFzc2lzdGFudC5zdXBwb3J0PC9rZXk+CgkJCQk8ZGljdD4KCQkJCQk8a2V5PkZvcmNlZDwva2V5PgoJCQkJCTxhcnJheT4KCQkJCQkJPGRpY3Q+CgkJCQkJCQk8a2V5Pm1jeF9wcmVmZXJlbmNlX3NldHRpbmdzPC9rZXk+CgkJCQkJCQk8ZGljdD4KCQkJCQkJCQk8a2V5PkFzc2lzdGFudCBFbmFibGVkPC9rZXk+CgkJCQkJCQkJPGZhbHNlLz4KCQkJCQkJCTwvZGljdD4KCQkJCQkJPC9kaWN0PgoJCQkJCTwvYXJyYXk+CgkJCQk8L2RpY3Q+CgkJCTwvZGljdD4KCQkJPGtleT5QYXlsb2FkRW5hYmxlZDwva2V5PgoJCQk8dHJ1ZS8+CgkJCTxrZXk+UGF5bG9hZElkZW50aWZpZXI8L2tleT4KCQkJPHN0cmluZz5NQ1hUb1Byb2ZpbGUuOWM3MzgwZDItNWJmZS00ZTYwLWJhZDMtMjVhZDg2ZDYxNTBkLmFsYWNhcnRlLmN1c3RvbXNldHRpbmdzLmZiOTU4ZDk2LWE0MzEtNDUzNi04NGQwLTFiZTQ4MjM4NWZiMjwvc3RyaW5nPgoJCQk8a2V5PlBheWxvYWRUeXBlPC9rZXk+CgkJCTxzdHJpbmc+Y29tLmFwcGxlLk1hbmFnZWRDbGllbnQucHJlZmVyZW5jZXM8L3N0cmluZz4KCQkJPGtleT5QYXlsb2FkVVVJRDwva2V5PgoJCQk8c3RyaW5nPmZiOTU4ZDk2LWE0MzEtNDUzNi04NGQwLTFiZTQ4MjM4NWZiMjwvc3RyaW5nPgoJCQk8a2V5PlBheWxvYWRWZXJzaW9uPC9rZXk+CgkJCTxpbnRlZ2VyPjE8L2ludGVnZXI+CgkJPC9kaWN0PgoJPC9hcnJheT4KCTxrZXk+UGF5bG9hZERlc2NyaXB0aW9uPC9rZXk+Cgk8c3RyaW5nPlN0b3BzIFNpcmkgZnJvbSBiZWluZyBlbmFibGVkLjwvc3RyaW5nPgoJPGtleT5QYXlsb2FkRGlzcGxheU5hbWU8L2tleT4KCTxzdHJpbmc+RGlzYWJsZSBTaXJpPC9zdHJpbmc+Cgk8a2V5PlBheWxvYWRJZGVudGlmaWVyPC9rZXk+Cgk8c3RyaW5nPkRpc2FibGVTaXJpPC9zdHJpbmc+Cgk8a2V5PlBheWxvYWRPcmdhbml6YXRpb248L2tleT4KCTxzdHJpbmc+PC9zdHJpbmc+Cgk8a2V5PlBheWxvYWRSZW1vdmFsRGlzYWxsb3dlZDwva2V5PgoJPHRydWUvPgoJPGtleT5QYXlsb2FkU2NvcGU8L2tleT4KCTxzdHJpbmc+U3lzdGVtPC9zdHJpbmc+Cgk8a2V5PlBheWxvYWRUeXBlPC9rZXk+Cgk8c3RyaW5nPkNvbmZpZ3VyYXRpb248L3N0cmluZz4KCTxrZXk+UGF5bG9hZFVVSUQ8L2tleT4KCTxzdHJpbmc+OWM3MzgwZDItNWJmZS00ZTYwLWJhZDMtMjVhZDg2ZDYxNTBkPC9zdHJpbmc+Cgk8a2V5PlBheWxvYWRWZXJzaW9uPC9rZXk+Cgk8aW50ZWdlcj4xPC9pbnRlZ2VyPgo8L2RpY3Q+CjwvcGxpc3Q+Cg==", "request_type": "InstallProfile"}`,
			),
			testFn: func(t *testing.T, parts endToEndParts) {
				if len(parts.req.Command.InstallProfile.Payload) == 0 {
					t.Error("InstallProfile payload is empty after json unmarshal")
				}
				if len(parts.fromProto.Command.InstallProfile.Payload) == 0 {
					t.Error("unmarshaled proto payload is missing payload")
				}
				if !bytes.Contains(parts.plistData, []byte(`PD94bWwgdm`)) {
					t.Error("marshaled plist does not contain the required payload")
				}
			},
		},

		{
			name: "InstallEnterpriseApplication",
			requestBytes: []byte(
				`{"udid":"B59A5A44-EC36-4244-AB52-C40F6100528A","request_type":"InstallEnterpriseApplication","manifest":{"items":[{"metadata":{"items":[{"bundle-version":"1.7.5","bundle-identifier":"com.myenterprise.MyAppNotMAS"}],"bundle-version":"1.1","bundle-identifier":"com.myenterprise.MyAppPackage","kind":"display-image","sizeInBytes":1234,"title":"Test Title","subtitle":"Test Subtitle"},"Assets":[{"sha256-size":1234,"sha256s":["2a8a98c146c35ce29f8b9af4cf8218d2c026058e7eb35adb4a00236997593471"],"url":"https://example.com/p.pkg","kind":"software-package","md5-size":1234,"md5s":["cfdc14fa22a79bab2a8b423daca2c076"]}]}]}}`,
			),
			testFn: func(t *testing.T, parts endToEndParts) {
				needToSee := [][]byte{
					[]byte(`cfdc14fa22a79bab2a8b423daca2c076`),
					[]byte(`https://example.com/p.pkg`),
					[]byte(`com.myenterprise.MyAppPackage`),
					[]byte(`1.1`),
					[]byte(`1234`),
					[]byte(`2a8a98c146c35ce29f8b9af4cf8218d2c026058e7eb35adb4a00236997593471`),
					[]byte(`com.myenterprise.MyAppPackage`),
					[]byte(`com.myenterprise.MyAppNotMAS`),
					[]byte(`1.7.5`),
					[]byte(`software-package`),
					[]byte(`display-image`),
					[]byte(`Test Title`),
					[]byte(`Test Subtitle`),
				}
				for _, b := range needToSee {
					if !bytes.Contains(parts.plistData, b) {
						t.Error(fmt.Sprintf("marshaled plist does not contain required bytes: '%s'", string(b)))
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parts := endToEnd(t, tt.requestBytes)
			tt.testFn(t, parts)
		})
	}
}

type endToEndParts struct {
	requestBytes []byte          // some json, our API request
	req          CommandRequest  // after unmarshal
	payload      *CommandPayload // new payload
	protoData    []byte          // stored as
	fromProto    CommandPayload  // back from proto
	plistData    []byte          // final representation
}

func endToEnd(t *testing.T, requestBytes []byte) endToEndParts {
	t.Helper()
	var (
		err   error
		parts = endToEndParts{requestBytes: requestBytes}
	)

	if err = json.Unmarshal(parts.requestBytes, &parts.req); err != nil {
		t.Fatal(err)
	}

	if parts.payload, err = NewCommandPayload(&parts.req); err != nil {
		t.Fatal(err)
	}

	if parts.protoData, err = MarshalCommandPayload(parts.payload); err != nil {
		t.Fatal(err)
	}

	if err := UnmarshalCommandPayload(parts.protoData, &parts.fromProto); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(*parts.payload, parts.fromProto) {
		t.Errorf("command from json request does not match command from proto")
	}

	if parts.plistData, err = plist.MarshalIndent(parts.fromProto, "  "); err != nil {
		t.Fatal(err)
	}

	return parts
}
