package mdm

import (
	"github.com/gogo/protobuf/proto"
	"github.com/micromdm/micromdm/mdm/mdm/internal/mdmproto"
)

func commandToProto(cmd *Command) *mdmproto.Command {
	cmdproto := mdmproto.Command{
		RequestType: cmd.RequestType,
	}
	switch cmd.RequestType {
	case "InstallProfile":
		cmdproto.Request = &mdmproto.Command_InstallProfile{
			InstallProfile: &mdmproto.InstallProfile{
				Payload: cmd.InstallProfile.Payload,
			},
		}
	case "RemoveProfile":
		cmdproto.Request = &mdmproto.Command_RemoveProfile{
			RemoveProfile: &mdmproto.RemoveProfile{
				Identifier: cmd.RemoveProfile.Identifier,
			},
		}
	}
	return &cmdproto
}

func protoToCommand(pb *mdmproto.Command) *Command {
	cmd := Command{
		RequestType: pb.RequestType,
	}
	switch pb.RequestType {
	case "InstallProfile":
		cmd.InstallProfile = &InstallProfile{
			Payload: pb.GetInstallProfile().GetPayload(),
		}
	case "RemoveProfile":
		cmd.RemoveProfile = &RemoveProfile{
			Identifier: pb.GetRemoveProfile().GetIdentifier(),
		}
	}
	return &cmd
}

func MarshalCommandPayload(cmd *CommandPayload) ([]byte, error) {
	cmdproto := mdmproto.CommandPayload{
		CommandUuid: cmd.CommandUUID,
		Command:     commandToProto(cmd.Command),
	}
	return proto.Marshal(&cmdproto)
}

func UnmarshalCommandPayload(data []byte, payload *CommandPayload) error {
	var pb mdmproto.CommandPayload
	if err := proto.Unmarshal(data, &pb); err != nil {
		return err
	}
	payload.CommandUUID = pb.CommandUuid
	payload.Command = protoToCommand(pb.Command)
	return nil
}
