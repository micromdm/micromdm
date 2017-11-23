package checkin_test

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/groob/plist"
	"github.com/micromdm/mdm"

	"github.com/micromdm/micromdm/mdm/checkin"
)

var marshalTests = []string{
	"Authenticate",
	"TokenUpdate",
	"CheckOut",
}

func TestMarshalEvent(t *testing.T) {
	for _, tt := range marshalTests {
		name := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			v := checkin.NewEvent(mustLoadCommand(t, name))
			var other checkin.Event
			if buf, err := checkin.MarshalEvent(v); err != nil {
				t.Fatal(err)
			} else if err := checkin.UnmarshalEvent(buf, &other); err != nil {
				t.Fatal(err)
			} else if !reflect.DeepEqual(v, &other) {
				t.Fatalf("\nwant: %#v\n \nhave: %#v\n", v, &other)
			}
		})
	}

}

func mustLoadCommand(t *testing.T, name string) mdm.CheckinCommand {
	var payload mdm.CheckinCommand
	data, err := ioutil.ReadFile("testdata/" + name + ".plist")
	if err != nil {
		t.Fatalf("failed to open test file %q.plist, err: %s", name, err)
	}
	if err := plist.Unmarshal(data, &payload); err != nil {
		t.Fatalf("failed to unmarshal plist %q, err: %s", name, err)
	}
	return payload
}
