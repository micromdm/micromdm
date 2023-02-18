package command

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/micromdm/micromdm/mdm/mdm"
	"github.com/micromdm/micromdm/platform/command/internal/commandproto"
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
		ID:         uuid.New().String(),
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
	e.ID = pb.Id
	e.DeviceUDID = pb.DeviceUdid
	e.Time = time.Unix(0, pb.Time).UTC()
	e.Payload = &payload
	return nil
}

type RawEvent struct {
	CommandUUID string
	Time        time.Time
	DeviceUDID  string
	Payload     []byte
}

// NewRawEvent returns a RawEvent with the current time.
func NewRawEvent(cmd *RawCommand) *RawEvent {
	event := RawEvent{
		CommandUUID: cmd.CommandUUID,
		Time:        time.Now().UTC(),
		DeviceUDID:  cmd.UDID,
		Payload:     cmd.Raw,
	}
	return &event
}

// MarshalRawEvent serializes a RawEvent to a protocol buffer wire format.
func MarshalRawEvent(e *RawEvent) ([]byte, error) {
	return proto.Marshal(&commandproto.Event{
		Id:           e.CommandUUID, // Id isn't used anywhere, so it's repurposed for CommandUUID
		Time:         e.Time.UnixNano(),
		DeviceUdid:   e.DeviceUDID,
		PayloadBytes: e.Payload,
	})
}

// UnmarshalRawEvent parses a protocol buffer representation of data into
// the RawEvent.
func UnmarshalRawEvent(data []byte, e *RawEvent) error {
	var pb commandproto.Event
	if err := proto.Unmarshal(data, &pb); err != nil {
		return errors.Wrap(err, "unmarshal pb Event")
	}
	var payload mdm.CommandPayload
	if err := mdm.UnmarshalCommandPayload(pb.PayloadBytes, &payload); err != nil {
		return err
	}
	e.CommandUUID = pb.Id
	e.Time = time.Unix(0, pb.Time).UTC()
	e.DeviceUDID = pb.DeviceUdid
	e.Payload = pb.PayloadBytes
	return nil
}
