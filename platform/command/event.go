package command

import (
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/micromdm/micromdm/mdm/mdm"
	"github.com/micromdm/micromdm/platform/command/internal/commandproto"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

type Event struct {
	ID         string
	Time       time.Time
	Payload    *mdm.CommandPayload
	DeviceUDID string
}

// NewEvent returns an Event with a unique ID and the current time.
func NewEvent(payload *mdm.CommandPayload, udid string) *Event {
	event := Event{
		ID:         uuid.NewV4().String(),
		Time:       time.Now().UTC(),
		Payload:    payload,
		DeviceUDID: udid,
	}
	return &event
}

// MarshalEvent serializes an event to a protocol buffer wire format.
func MarshalEvent(e *Event) ([]byte, error) {
	payloadBytes, err := mdm.MarshalCommandPayload(e.Payload)
	if err != nil {
		return nil, err
	}
	return proto.Marshal(&commandproto.Event{
		Id:           e.ID,
		Time:         e.Time.UnixNano(),
		PayloadBytes: payloadBytes,
		DeviceUdid:   e.DeviceUDID,
	})

}

// UnmarshalEvent parses a protocol buffer representation of data into
// the Event.
func UnmarshalEvent(data []byte, e *Event) error {
	var pb commandproto.Event
	if err := proto.Unmarshal(data, &pb); err != nil {
		return errors.Wrap(err, "unmarshal pb Event")
	}
	var payload mdm.CommandPayload
	if err := mdm.UnmarshalCommandPayload(pb.PayloadBytes, &payload); err != nil {
		return err
	}
	fmt.Println(string(pb.PayloadBytes))
	fmt.Println("")
	fmt.Println(payload.Command)

	e.ID = pb.Id
	e.DeviceUDID = pb.DeviceUdid
	e.Time = time.Unix(0, pb.Time).UTC()
	e.Payload = &payload
	return nil
}
