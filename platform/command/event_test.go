package command_test

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/groob/plist"
	"github.com/micromdm/mdm"

	"github.com/micromdm/micromdm/platform/command"
)

var marshalTests = []string{
	"DeviceInformation",
	"DeviceInformation_empty_queries",
	"InstallProfile",
	"Settings_hostname_devicename",
}

func TestMarshalEvent(t *testing.T) {
	for _, tt := range marshalTests {
		name := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			v := command.NewEvent(mustLoadPayload(t, name), name)
			var other command.Event
			if buf, err := command.MarshalEvent(v); err != nil {
				t.Fatal(err)
			} else if err := command.UnmarshalEvent(buf, &other); err != nil {
				t.Fatal(err)
			} else if !reflect.DeepEqual(v, &other) {
				t.Logf("\nwant: %#v\n, \nhave: %#v\n", v.Payload.Command, other.Payload.Command)
				t.Fatalf("\nwant: %#v\n \nhave: %#v\n", v, other)
			}
		})
	}
}

func BenchmarkMarshalProto(b *testing.B) {
	for _, tt := range marshalTests {
		v := command.NewEvent(mustLoadPayload(&testing.T{}, tt), tt)
		for n := 0; n < b.N; n++ {
			var other command.Event
			if buf, err := command.MarshalEvent(v); err != nil {
				b.Fatal(err)
			} else if err := command.UnmarshalEvent(buf, &other); err != nil {
				b.Fatal(err)
			} else if !reflect.DeepEqual(v, &other) {
				b.Fatalf("\nwant: %#v\n \nhave: %#v\n", v, other)
			}
		}
	}
}

func BenchmarkMarshalJSON(b *testing.B) {
	for _, tt := range marshalTests {
		v := command.NewEvent(mustLoadPayload(&testing.T{}, tt), tt)
		for n := 0; n < b.N; n++ {
			var other command.Event
			if buf, err := json.Marshal(&v); err != nil {
				b.Fatal(err)
			} else if err := json.Unmarshal(buf, &other); err != nil {
				b.Fatal(err)
			} else if !reflect.DeepEqual(v, &other) {
				b.Fatalf("\nwant: %#v\n \nhave: %#v\n", v, other)
			}
		}
	}
}

func mustLoadPayload(t *testing.T, name string) mdm.Payload {
	var payload mdm.Payload
	data, err := ioutil.ReadFile("testdata/" + name + ".plist")
	if err != nil {
		t.Fatalf("failed to open test file %q.plist, err: %s", name, err)
	}
	if err := plist.Unmarshal(data, &payload); err != nil {
		t.Fatalf("failed to unmarshal plist %q, err: %s", name, err)
	}
	return payload
}
