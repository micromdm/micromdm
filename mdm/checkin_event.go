package mdm

import (
	"encoding/hex"
	"time"

	"github.com/gogo/protobuf/proto"

	"github.com/micromdm/micromdm/mdm/internal/checkinproto"
)

type CheckinEvent struct {
	ID      string
	Time    time.Time
	Command CheckinCommand
	Params  map[string]string
	Raw     []byte
}

// CheckinRequest represents an MDM checkin command struct.
type CheckinCommand struct {
	// MessageType can be either Authenticate,
	// TokenUpdate or CheckOut
	MessageType string
	Topic       string
	UDID        string
	auth
	update
}

// Authenticate Message Type
type auth struct {
	OSVersion    string
	BuildVersion string
	ProductName  string
	SerialNumber string
	IMEI         string
	MEID         string
	DeviceName   string `plist:"DeviceName,omitempty"`
	Challenge    []byte `plist:"Challenge,omitempty"`
	Model        string `plist:"Model,omitpempty"`
	ModelName    string `plist:"ModelName,omitempty"`
}

// TokenUpdate Mesage Type
type update struct {
	Token                 hexData
	PushMagic             string
	UnlockToken           hexData
	AwaitingConfiguration bool
	userTokenUpdate
}

// TokenUpdate with user keys
type userTokenUpdate struct {
	UserID        string `plist:",omitempty"`
	UserLongName  string `plist:",omitempty"`
	UserShortName string `plist:",omitempty"`
	NotOnConsole  bool   `plist:",omitempty"`
}

// data decodes to []byte,
// we can then attach a string method to the type
// Tokens are encoded as Hex Strings
type hexData []byte

func (d hexData) String() string {
	return hex.EncodeToString(d)
}

// MarshalCheckinEvent serializes an event to a protocol buffer wire format.
func MarshalCheckinEvent(e *CheckinEvent) ([]byte, error) {
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
		Params:  e.Params,
		Raw:     e.Raw,
	})
}

// UnmarshalCheckinEvent parses a protocol buffer representation of data into
// the Event.
func UnmarshalCheckinEvent(data []byte, e *CheckinEvent) error {
	var pb checkinproto.Event
	if err := proto.Unmarshal(data, &pb); err != nil {
		return err
	}
	e.ID = pb.Id
	e.Time = time.Unix(0, pb.Time).UTC()
	if pb.Command == nil {
		return nil
	}
	e.Command = CheckinCommand{
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
	e.Raw = pb.GetRaw()
	e.Params = pb.GetParams()
	return nil
}
