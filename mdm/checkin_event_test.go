package mdm

import (
	"bytes"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/micromdm/micromdm/mdm/internal/checkinproto"
)

func TestCheckinEventDeclarativeManagementRoundTrip(t *testing.T) {
	in := &CheckinEvent{
		ID:   "event-1",
		Time: time.Now().UTC(),
		Command: CheckinCommand{
			MessageType: "DeclarativeManagement",
			UDID:        "UDID-1",
		},
		Raw: []byte("raw"),
	}
	in.Command.Endpoint = "status"
	in.Command.Data = []byte(`{"StatusItems":{}}`)

	msg, err := MarshalCheckinEvent(in)
	if err != nil {
		t.Fatal(err)
	}

	var out CheckinEvent
	if err := UnmarshalCheckinEvent(msg, &out); err != nil {
		t.Fatal(err)
	}
	if out.Command.Endpoint != "status" {
		t.Errorf("Endpoint = %q, want status", out.Command.Endpoint)
	}
	if !bytes.Equal(out.Command.Data, in.Command.Data) {
		t.Errorf("Data = %q, want %q", out.Command.Data, in.Command.Data)
	}
	if !bytes.Equal(out.Raw, in.Raw) {
		t.Errorf("Raw = %q, want %q", out.Raw, in.Raw)
	}
}

func TestUnmarshalCheckinEventDeclarativeManagementWithoutBody(t *testing.T) {
	msg, err := proto.Marshal(&checkinproto.Event{
		Id:      "event-1",
		Command: &checkinproto.Command{MessageType: "DeclarativeManagement"},
	})
	if err != nil {
		t.Fatal(err)
	}

	var out CheckinEvent
	if err := UnmarshalCheckinEvent(msg, &out); err != nil {
		t.Fatal(err)
	}
	if out.Command.Endpoint != "" || out.Command.Data != nil {
		t.Errorf("got Endpoint=%q Data=%q, want empty", out.Command.Endpoint, out.Command.Data)
	}
}
