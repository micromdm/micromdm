package queue

import (
	"time"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/micromdm/micromdm/platform/queue/internal/devicecommandproto"
)

type Command struct {
	UUID    string
	Payload []byte

	CreatedAt    time.Time
	LastSentAt   time.Time
	Acknowledged time.Time

	TimesSent int

	LastStatus     string
	FailureMessage []byte
}

type DeviceCommand struct {
	DeviceUDID string
	Commands   []Command

	// These are going to scale great. We'll have to see.
	Completed []Command
	Failed    []Command
	NotNow    []Command
}

func MarshalDeviceCommand(c *DeviceCommand) ([]byte, error) {
	protoc := devicecommandproto.DeviceCommand{
		DeviceUdid: c.DeviceUDID,
	}

	// TODO add helper here to reduce copy/pasted boilerplate.
	for _, command := range c.Commands {
		protoc.Commands = append(protoc.Commands, &devicecommandproto.Command{
			Uuid:         command.UUID,
			Payload:      command.Payload,
			CreatedAt:    command.CreatedAt.UnixNano(),
			LastSentAt:   command.LastSentAt.UnixNano(),
			Acknowledged: command.Acknowledged.UnixNano(),

			TimesSent: int64(command.TimesSent),

			LastStatus:     command.LastStatus,
			FailureMessage: command.FailureMessage,
		})
	}

	for _, command := range c.Completed {
		protoc.Completed = append(protoc.Completed, &devicecommandproto.Command{
			Uuid:         command.UUID,
			Payload:      command.Payload,
			CreatedAt:    command.CreatedAt.UnixNano(),
			LastSentAt:   command.LastSentAt.UnixNano(),
			Acknowledged: command.Acknowledged.UnixNano(),

			TimesSent: int64(command.TimesSent),

			LastStatus:     command.LastStatus,
			FailureMessage: command.FailureMessage,
		})
	}

	for _, command := range c.Failed {
		protoc.Failed = append(protoc.Failed, &devicecommandproto.Command{
			Uuid:         command.UUID,
			Payload:      command.Payload,
			CreatedAt:    command.CreatedAt.UnixNano(),
			LastSentAt:   command.LastSentAt.UnixNano(),
			Acknowledged: command.Acknowledged.UnixNano(),

			TimesSent: int64(command.TimesSent),

			LastStatus:     command.LastStatus,
			FailureMessage: command.FailureMessage,
		})
	}

	for _, command := range c.NotNow {
		protoc.NotNow = append(protoc.NotNow, &devicecommandproto.Command{
			Uuid:         command.UUID,
			Payload:      command.Payload,
			CreatedAt:    command.CreatedAt.UnixNano(),
			LastSentAt:   command.LastSentAt.UnixNano(),
			Acknowledged: command.Acknowledged.UnixNano(),

			TimesSent: int64(command.TimesSent),

			LastStatus:     command.LastStatus,
			FailureMessage: command.FailureMessage,
		})
	}
	return proto.Marshal(&protoc)
}

func UnmarshalDeviceCommand(data []byte, c *DeviceCommand) error {
	var pb devicecommandproto.DeviceCommand
	if err := proto.Unmarshal(data, &pb); err != nil {
		return errors.Wrap(err, "unmarshal proto to DeviceCommand")
	}
	c.DeviceUDID = pb.GetDeviceUdid()
	protoCommands := pb.GetCommands()
	protoCommandsCompleted := pb.GetCompleted()
	protoCommandsFailed := pb.GetFailed()
	protoCommandsNotNow := pb.GetNotNow()
	for _, command := range protoCommands {
		c.Commands = append(c.Commands, Command{
			UUID:         command.GetUuid(),
			Payload:      command.GetPayload(),
			CreatedAt:    time.Unix(0, command.GetCreatedAt()).UTC(),
			LastSentAt:   time.Unix(0, command.GetLastSentAt()).UTC(),
			Acknowledged: time.Unix(0, command.GetAcknowledged()).UTC(),

			TimesSent: int(command.TimesSent),

			LastStatus:     command.LastStatus,
			FailureMessage: command.FailureMessage,
		})
	}

	for _, command := range protoCommandsCompleted {
		c.Completed = append(c.Completed, Command{
			UUID:         command.GetUuid(),
			Payload:      command.GetPayload(),
			CreatedAt:    time.Unix(0, command.GetCreatedAt()).UTC(),
			LastSentAt:   time.Unix(0, command.GetLastSentAt()).UTC(),
			Acknowledged: time.Unix(0, command.GetAcknowledged()).UTC(),

			TimesSent: int(command.TimesSent),

			LastStatus:     command.LastStatus,
			FailureMessage: command.FailureMessage,
		})
	}

	for _, command := range protoCommandsFailed {
		c.Failed = append(c.Failed, Command{
			UUID:         command.GetUuid(),
			Payload:      command.GetPayload(),
			CreatedAt:    time.Unix(0, command.GetCreatedAt()).UTC(),
			LastSentAt:   time.Unix(0, command.GetLastSentAt()).UTC(),
			Acknowledged: time.Unix(0, command.GetAcknowledged()).UTC(),

			TimesSent: int(command.TimesSent),

			LastStatus:     command.LastStatus,
			FailureMessage: command.FailureMessage,
		})
	}

	for _, command := range protoCommandsNotNow {
		c.NotNow = append(c.NotNow, Command{
			UUID:         command.GetUuid(),
			Payload:      command.GetPayload(),
			CreatedAt:    time.Unix(0, command.GetCreatedAt()).UTC(),
			LastSentAt:   time.Unix(0, command.GetLastSentAt()).UTC(),
			Acknowledged: time.Unix(0, command.GetAcknowledged()).UTC(),

			TimesSent: int(command.TimesSent),

			LastStatus:     command.LastStatus,
			FailureMessage: command.FailureMessage,
		})
	}
	return nil
}
