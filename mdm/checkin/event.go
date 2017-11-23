package checkin

import (
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/micromdm/mdm"
	uuid "github.com/satori/go.uuid"

	"github.com/micromdm/micromdm/mdm/checkin/internal/checkinproto"
)

type Event struct {
	ID      string
	Time    time.Time
	Command mdm.CheckinCommand
}

// NewEvent returns an Event with a unique ID and the current time.
func NewEvent(cmd mdm.CheckinCommand) *Event {
	event := Event{
		ID:      uuid.NewV4().String(),
		Time:    time.Now().UTC(),
		Command: cmd,
	}
	return &event
}

// MarshalEvent serializes an event to a protocol buffer wire format.
func MarshalEvent(e *Event) ([]byte, error) {
	command := &checkinproto.Command{
		MessageType: e.Command.MessageType,
		Topic:       e.Command.Topic,
		Udid:        e.Command.UDID,
	}
	switch e.Command.MessageType {
	case "Authenticate":
		command.Authenticate = &checkinproto.Authenticate{
			OsVersion:    e.Command.OSVersion,
			BuildVersion: e.Command.BuildVersion,
			SerialNumber: e.Command.SerialNumber,
			Imei:         e.Command.IMEI,
			Meid:         e.Command.MEID,
			DeviceName:   e.Command.DeviceName,
			Challenge:    e.Command.Challenge,
			Model:        e.Command.Model,
			ModelName:    e.Command.ModelName,
			ProductName:  e.Command.ProductName,
		}
	case "TokenUpdate":
		command.TokenUpdate = &checkinproto.TokenUpdate{
			Token:                 e.Command.Token,
			PushMagic:             e.Command.PushMagic,
			UnlockToken:           e.Command.UnlockToken,
			AwaitingConfiguration: e.Command.AwaitingConfiguration,
			UserId:                e.Command.UserID,
			UserLongName:          e.Command.UserLongName,
			UserShortName:         e.Command.UserShortName,
			NotOnConsole:          e.Command.NotOnConsole,
		}
	}
	return proto.Marshal(&checkinproto.Event{
		Id:      e.ID,
		Time:    e.Time.UnixNano(),
		Command: command,
	})
}

// UnmarshalEvent parses a protocol buffer representation of data into
// the Event.
func UnmarshalEvent(data []byte, e *Event) error {
	var pb checkinproto.Event
	if err := proto.Unmarshal(data, &pb); err != nil {
		return err
	}
	e.ID = pb.Id
	e.Time = time.Unix(0, pb.Time).UTC()
	if pb.Command == nil {
		return nil
	}
	e.Command = mdm.CheckinCommand{
		MessageType: pb.Command.MessageType,
		Topic:       pb.Command.Topic,
		UDID:        pb.Command.Udid,
	}
	switch pb.Command.MessageType {
	case "Authenticate":
		e.Command.OSVersion = pb.Command.Authenticate.OsVersion
		e.Command.BuildVersion = pb.Command.Authenticate.BuildVersion
		e.Command.SerialNumber = pb.Command.Authenticate.SerialNumber
		e.Command.IMEI = pb.Command.Authenticate.Imei
		e.Command.MEID = pb.Command.Authenticate.Meid
		e.Command.DeviceName = pb.Command.Authenticate.DeviceName
		e.Command.Challenge = pb.Command.Authenticate.Challenge
		e.Command.Model = pb.Command.Authenticate.Model
		e.Command.ModelName = pb.Command.Authenticate.ModelName
		e.Command.ProductName = pb.Command.Authenticate.ProductName
	case "TokenUpdate":
		e.Command.Token = pb.Command.TokenUpdate.Token
		e.Command.PushMagic = pb.Command.TokenUpdate.PushMagic
		e.Command.UnlockToken = pb.Command.TokenUpdate.UnlockToken
		e.Command.AwaitingConfiguration = pb.Command.TokenUpdate.AwaitingConfiguration
		e.Command.UserID = pb.Command.TokenUpdate.UserId
		e.Command.UserLongName = pb.Command.TokenUpdate.UserLongName
		e.Command.UserShortName = pb.Command.TokenUpdate.UserShortName
		e.Command.NotOnConsole = pb.Command.TokenUpdate.NotOnConsole
	}
	return nil
}
