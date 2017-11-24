package command

import (
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/micromdm/mdm"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/micromdm/micromdm/platform/command/internal/commandproto"
)

type Event struct {
	ID         string
	Time       time.Time
	Payload    mdm.Payload
	DeviceUDID string
}

// NewEvent returns an Event with a unique ID and the current time.
func NewEvent(cmd mdm.Payload, udid string) *Event {
	event := Event{
		ID:         uuid.NewV4().String(),
		Time:       time.Now().UTC(),
		Payload:    cmd,
		DeviceUDID: udid,
	}
	return &event
}

// MarshalEvent serializes an event to a protocol buffer wire format.
func MarshalEvent(e *Event) ([]byte, error) {
	payload := &commandproto.Payload{
		CommandUuid: e.Payload.CommandUUID,
	}
	if e.Payload.Command != nil {
		payload.Command = &commandproto.Command{
			RequestType: e.Payload.Command.RequestType,
		}
	}
	switch e.Payload.Command.RequestType {
	case "DeviceLock":
		payload.Command.DeviceLock = &commandproto.DeviceLock{
			Pin:         e.Payload.Command.DeviceLock.PIN,
			Message:     e.Payload.Command.DeviceLock.Message,
			PhoneNumber: e.Payload.Command.DeviceLock.PhoneNumber,
		}
	case "EraseDevice":
		payload.Command.EraseDevice = &commandproto.EraseDevice{
			Pin: e.Payload.Command.EraseDevice.PIN,
		}
	case "DeleteUser":
		payload.Command.DeleteUser = &commandproto.DeleteUser{
			Username:      e.Payload.Command.DeleteUser.UserName,
			ForceDeletion: e.Payload.Command.DeleteUser.ForceDeletion,
		}
	case "ScheduleOSUpdateScan":
		payload.Command.ScheduleOsUpdateScan = &commandproto.ScheduleOSUpdateScan{
			Force: e.Payload.Command.ScheduleOSUpdateScan.Force,
		}
	case "ScheduleOSUpdate":
		p := e.Payload.Command.ScheduleOSUpdate
		var updates []*commandproto.OSUpdate
		for _, update := range p.Updates {
			updates = append(updates, &commandproto.OSUpdate{
				ProductKey:    update.ProductKey,
				InstallAction: update.InstallAction,
			})
		}
		payload.Command.ScheduleOsUpdate = &commandproto.ScheduleOSUpdate{
			Updates: updates,
		}
	case "AccountConfiguration":
		p := e.Payload.Command.AccountConfiguration
		payload.Command.AccountConfiguration = &commandproto.AccountConfiguration{
			SkipPrimarySetupAccountCreation:     p.SkipPrimarySetupAccountCreation,
			SetPrimarySetupAccountAsRegularUser: p.SetPrimarySetupAccountAsRegularUser,
		}
		for _, account := range p.AutoSetupAdminAccounts {
			payload.Command.AccountConfiguration.AutoSetupAdminAccounts = append(
				payload.Command.AccountConfiguration.AutoSetupAdminAccounts, &commandproto.AutoSetupAdminAccounts{
					ShortName:    account.ShortName,
					FullName:     account.FullName,
					PasswordHash: account.PasswordHash,
					Hidden:       account.Hidden,
				})
		}
	case "DeviceInformation":
		payload.Command.DeviceInformation = &commandproto.DeviceInformation{
			Queries: e.Payload.Command.DeviceInformation.Queries,
		}
	case "InstallProfile":
		payload.Command.InstallProfile = &commandproto.InstallProfile{
			Payload: e.Payload.Command.InstallProfile.Payload,
		}
	case "RemoveProfile":
		payload.Command.RemoveProfile = &commandproto.RemoveProfile{
			Identifier: e.Payload.Command.RemoveProfile.Identifier,
		}
	case "InstallApplication":
		cmd := e.Payload.Command.InstallApplication
		payload.Command.InstallApplication = &commandproto.InstallApplication{
			ItunesStoreId:         int64(cmd.ITunesStoreID),
			Identifier:            cmd.Identifier,
			ManifestUrl:           cmd.ManifestURL,
			ManagementFlags:       int64(cmd.ManagementFlags),
			NotManaged:            cmd.NotManaged,
			ChangeManagementState: cmd.ChangeManagementState,
		}
	case "Settings":
		cmd := e.Payload.Command.Settings
		var settings []*commandproto.Setting
		for _, s := range cmd.Settings {
			protoSetting := &commandproto.Setting{
				Item: s.Item,
			}
			if s.DeviceName != nil {
				protoSetting.DeviceName = &commandproto.DeviceNameSetting{
					DeviceName: *s.DeviceName,
				}
			}

			if s.HostName != nil {
				protoSetting.Hostname = &commandproto.HostnameSetting{
					Hostname: *s.HostName,
				}
			}
			settings = append(settings, protoSetting)
		}
		payload.Command.Settings = &commandproto.Settings{Settings: settings}
	}
	return proto.Marshal(&commandproto.Event{
		Id:         e.ID,
		Time:       e.Time.UnixNano(),
		Payload:    payload,
		DeviceUdid: e.DeviceUDID,
	})

}

// UnmarshalEvent parses a protocol buffer representation of data into
// the Event.
func UnmarshalEvent(data []byte, e *Event) error {
	var pb commandproto.Event
	if err := proto.Unmarshal(data, &pb); err != nil {
		return errors.Wrap(err, "unmarshal pb Event")
	}
	e.ID = pb.Id
	e.DeviceUDID = pb.DeviceUdid
	e.Time = time.Unix(0, pb.Time).UTC()
	if pb.Payload == nil {
		return nil
	}
	e.Payload = mdm.Payload{
		CommandUUID: pb.Payload.CommandUuid,
	}
	if pb.Payload.Command == nil {
		return nil
	}
	e.Payload.Command = &mdm.Command{
		RequestType: pb.Payload.Command.RequestType,
	}
	switch pb.Payload.Command.RequestType {
	case "DeviceLock":
		cmd := pb.Payload.Command.GetDeviceLock()
		e.Payload.Command.DeviceLock = mdm.DeviceLock{
			PIN:         cmd.GetPin(),
			Message:     cmd.GetMessage(),
			PhoneNumber: cmd.GetPhoneNumber(),
		}
	case "EraseDevice":
		cmd := pb.Payload.Command.GetEraseDevice()
		e.Payload.Command.EraseDevice = mdm.EraseDevice{
			PIN: cmd.GetPin(),
		}
	case "DeleteUser":
		cmd := pb.Payload.Command.GetDeleteUser()
		e.Payload.Command.DeleteUser = mdm.DeleteUser{
			UserName:      cmd.GetUsername(),
			ForceDeletion: cmd.GetForceDeletion(),
		}
	case "ScheduleOSUpdateScan":
		cmd := pb.Payload.Command.GetScheduleOsUpdateScan()
		e.Payload.Command.ScheduleOSUpdateScan = mdm.ScheduleOSUpdateScan{
			Force: cmd.GetForce(),
		}
	case "ScheduleOSUpdate":
		cmd := pb.Payload.Command.GetScheduleOsUpdate()
		var updates []mdm.OSUpdate
		for _, update := range cmd.GetUpdates() {
			updates = append(updates, mdm.OSUpdate{
				ProductKey:    update.GetProductKey(),
				InstallAction: update.GetInstallAction(),
			})
		}
		e.Payload.Command.ScheduleOSUpdate = mdm.ScheduleOSUpdate{
			Updates: updates,
		}
	case "AccountConfiguration":
		cmd := pb.Payload.Command.GetAccountConfiguration()
		e.Payload.Command.AccountConfiguration = mdm.AccountConfiguration{
			SkipPrimarySetupAccountCreation:     cmd.GetSkipPrimarySetupAccountCreation(),
			SetPrimarySetupAccountAsRegularUser: cmd.GetSetPrimarySetupAccountAsRegularUser(),
		}
		for _, account := range cmd.GetAutoSetupAdminAccounts() {
			e.Payload.Command.AccountConfiguration.AutoSetupAdminAccounts = append(e.Payload.Command.AutoSetupAdminAccounts, mdm.AdminAccount{
				ShortName:    account.GetShortName(),
				FullName:     account.GetFullName(),
				PasswordHash: account.GetPasswordHash(),
				Hidden:       account.GetHidden(),
			})
		}
	case "DeviceInformation":
		e.Payload.Command.DeviceInformation = mdm.DeviceInformation{
			Queries: pb.Payload.Command.DeviceInformation.Queries,
		}
	case "InstallProfile":
		e.Payload.Command.InstallProfile = mdm.InstallProfile{
			Payload: pb.Payload.Command.InstallProfile.Payload,
		}
	case "RemoveProfile":
		e.Payload.Command.RemoveProfile = mdm.RemoveProfile{
			Identifier: pb.Payload.Command.RemoveProfile.Identifier,
		}
	case "InstallApplication":
		cmd := pb.Payload.Command.GetInstallApplication()
		e.Payload.Command.InstallApplication = mdm.InstallApplication{
			ITunesStoreID:         int(cmd.GetItunesStoreId()),
			Identifier:            cmd.GetIdentifier(),
			ManifestURL:           cmd.GetManifestUrl(),
			ManagementFlags:       int(cmd.GetManagementFlags()),
			ChangeManagementState: cmd.GetChangeManagementState(),
		}
	case "Settings":
		cmd := pb.Payload.Command.GetSettings()
		var settings []mdm.Setting
		for _, s := range cmd.GetSettings() {
			mdmSetting := mdm.Setting{
				Item: s.GetItem(),
			}

			if s.GetDeviceName() != nil {
				mdmSetting.DeviceName = stringPtr(s.GetDeviceName().GetDeviceName())
			}

			if s.GetHostname() != nil {
				mdmSetting.HostName = stringPtr(s.GetHostname().GetHostname())
			}

			settings = append(settings, mdmSetting)
		}
		e.Payload.Command.Settings = mdm.Settings{Settings: settings}
	}
	return nil
}

func stringPtr(s string) *string {
	return &s
}
