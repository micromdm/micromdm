package mdm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/groob/plist"
)

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
