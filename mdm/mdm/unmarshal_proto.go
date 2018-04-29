package mdm

import (
	"github.com/gogo/protobuf/proto"
	"github.com/micromdm/micromdm/mdm/mdm/internal/mdmproto"
)

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
	case "InstallProvisioningProfile":
		cmd.InstallProvisioningProfile = &InstallProvisioningProfile{
			ProvisioningProfile: pb.GetInstallProvisioningProfile().GetProvisioningProfile(),
		}
	case "RemoveProvisioningProfile":
		cmd.RemoveProvisioningProfile = &RemoveProvisioningProfile{
			UUID: pb.GetRemoveProfisioningProfile().GetUuid(),
		}
	case "InstalledApplicationList":
		pbcmd := pb.GetInstalledApplicationList()
		cmd.InstalledApplicationList = &InstalledApplicationList{
			Identifiers:     pbcmd.GetIdentifiers(),
			ManagedAppsOnly: pbcmd.GetManagedAppsOnly(),
		}
	case "DeviceInformation":
		cmd.DeviceInformation = &DeviceInformation{
			Queries: pb.GetDeviceInformation().GetQueries(),
		}
	case "DeviceLock":
		pbc := pb.GetDeviceLock()
		cmd.DeviceLock = &DeviceLock{
			PIN:         pbc.GetPin(),
			Message:     pbc.GetMessage(),
			PhoneNumber: pbc.GetPhoneNumber(),
		}
	case "ClearPasscode":
		pbc := pb.GetClearPasscode()
		cmd.ClearPasscode = &ClearPasscode{
			UnlockToken: pbc.GetUnlockToken(),
		}
	case "EraseDevice":
		pbc := pb.GetEraseDevice()
		cmd.EraseDevice = &EraseDevice{
			PIN:                    pbc.GetPin(),
			PreserveDataPlan:       pbc.GetPreserveDataPlan(),
			DisallowProximitySetup: pbc.GetDisallowProximitySetup(),
		}
	case "RequestMirroring":
		pbc := pb.GetRequestMirroring()
		cmd.RequestMirroring = &RequestMirroring{
			DestinationName:     pbc.GetDestinationName(),
			DestinationDeviceID: pbc.GetDestinationDeviceId(),
			ScanTime:            pbc.GetScanTime(),
			Password:            pbc.GetPassword(),
		}
	case "Restrictions":
		pbc := pb.GetRestrictions()
		cmd.Restrictions = &Restrictions{
			ProfileRestrictions: pbc.GetProfileRestrictions(),
		}
	case "UnlockUserAccount":
		pbc := pb.GetUnlockUserAccount()
		cmd.UnlockUserAccount = &UnlockUserAccount{
			UserName: pbc.GetUsername(),
		}
	}
	return &cmd
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
