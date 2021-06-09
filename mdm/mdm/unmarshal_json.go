package mdm

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

func (c *CommandRequest) UnmarshalJSON(data []byte) error {
	var request = struct {
		UDID        string `json:"udid"`
		RequestType string `json:"request_type"`
		CommandUUID string `json:"command_uuid"`
	}{}
	if err := json.Unmarshal(data, &request); err != nil {
		return errors.Wrap(err, "mdm: unmarshal json command request")
	}
	c.UDID = request.UDID
	c.Command = &Command{}
	c.CommandUUID = request.CommandUUID
	return c.Command.UnmarshalJSON(data)
}

func (c *Command) UnmarshalJSON(data []byte) error {
	var request = struct {
		RequestType string `json:"request_type"`
	}{}
	if err := json.Unmarshal(data, &request); err != nil {
		return errors.Wrap(err, "mdm: unmarshal json command request")
	}
	c.RequestType = request.RequestType

	switch c.RequestType {
	case "ProfileList",
		"ProvisioningProfileList",
		"CertificateList",
		"SecurityInfo",
		"RestartDevice",
		"ShutDownDevice",
		"StopMirroring",
		"ClearRestrictionsPassword",
		"UserList",
		"LogOutUser",
		"PlayLostModeSound",
		"DisableLostMode",
		"DeviceLocation",
		"ManagedMediaList",
		"DeviceConfigured",
		"AvailableOSUpdates",
		"NSExtensionMappings",
		"OSUpdateStatus",
		"EnableRemoteDesktop",
		"DisableRemoteDesktop",
		"ActivationLockBypassCode":
		return nil
	case "InstallProfile":
		var payload InstallProfile
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.InstallProfile = &payload
		return nil
	case "RemoveProfile":
		var payload RemoveProfile
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.RemoveProfile = &payload
		return nil
	case "InstallProvisioningProfile":
		var payload InstallProvisioningProfile
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.InstallProvisioningProfile = &payload
		return nil
	case "RemoveProvisioningProfile":
		var payload RemoveProvisioningProfile
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.RemoveProvisioningProfile = &payload
		return nil
	case "InstalledApplicationList":
		var payload InstalledApplicationList
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.InstalledApplicationList = &payload
		return nil
	case "DeviceInformation":
		var payload DeviceInformation
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.DeviceInformation = &payload
		return nil
	case "DeviceLock":
		var payload DeviceLock
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.DeviceLock = &payload
		return nil
	case "ClearPasscode":
		var payload ClearPasscode
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.ClearPasscode = &payload
		return nil
	case "EraseDevice":
		var payload EraseDevice
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.EraseDevice = &payload
		return nil
	case "RequestMirroring":
		var payload RequestMirroring
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.RequestMirroring = &payload
		return nil
	case "Restrictions":
		var payload Restrictions
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.Restrictions = &payload
		return nil
	case "UnlockUserAccount":
		var payload UnlockUserAccount
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.UnlockUserAccount = &payload
		return nil
	case "DeleteUser":
		var payload DeleteUser
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.DeleteUser = &payload
		return nil
	case "EnableLostMode":
		var payload EnableLostMode
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.EnableLostMode = &payload
		return nil
	case "InstallEnterpriseApplication":
		var payload InstallEnterpriseApplication
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.InstallEnterpriseApplication = &payload
		return nil
	case "InstallApplication":
		var payload InstallApplication
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.InstallApplication = &payload
		return nil
	case "AccountConfiguration":
		var payload AccountConfiguration
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.AccountConfiguration = &payload
		return nil
	case "ApplyRedemptionCode":
		var payload ApplyRedemptionCode
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.ApplyRedemptionCode = &payload
		return nil
	case "ManagedApplicationList":
		var payload ManagedApplicationList
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.ManagedApplicationList = &payload
		return nil
	case "RemoveApplication":
		var payload RemoveApplication
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.RemoveApplication = &payload
		return nil
	case "InviteToProgram":
		var payload InviteToProgram
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.InviteToProgram = &payload
		return nil
	case "ValidateApplications":
		var payload ValidateApplications
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.ValidateApplications = &payload
		return nil
	case "InstallMedia":
		var payload InstallMedia
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.InstallMedia = &payload
		return nil
	case "RemoveMedia":
		var payload RemoveMedia
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.RemoveMedia = &payload
		return nil
	case "Settings":
		var payload Settings
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.Settings = &payload
		return nil
	case "ManagedApplicationConfiguration":
		var payload ManagedApplicationConfiguration
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.ManagedApplicationConfiguration = &payload
		return nil
	case "ManagedApplicationAttributes":
		var payload ManagedApplicationAttributes
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.ManagedApplicationAttributes = &payload
		return nil
	case "ManagedApplicationFeedback":
		var payload ManagedApplicationFeedback
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.ManagedApplicationFeedback = &payload
		return nil
	case "SetFirmwarePassword":
		var payload SetFirmwarePassword
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.SetFirmwarePassword = &payload
		return nil
	case "VerifyFirmwarePassword":
		var payload VerifyFirmwarePassword
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.VerifyFirmwarePassword = &payload
		return nil
	case "SetRecoveryLock":
		var payload SetRecoveryLock
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.SetRecoveryLock = &payload
		return nil
	case "VerifyRecoveryLock":
		var payload VerifyRecoveryLock
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.VerifyRecoveryLock = &payload
		return nil
	case "SetAutoAdminPassword":
		var payload SetAutoAdminPassword
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.SetAutoAdminPassword = &payload
		return nil
	case "ScheduleOSUpdate":
		var payload ScheduleOSUpdate
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.ScheduleOSUpdate = &payload
		return nil
	case "ScheduleOSUpdateScan":
		var payload ScheduleOSUpdateScan
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.ScheduleOSUpdateScan = &payload
		return nil
	case "ActiveNSExtensions":
		var payload ActiveNSExtensions
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.ActiveNSExtensions = &payload
		return nil
	case "RotateFileVaultKey":
		var payload RotateFileVaultKey
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.RotateFileVaultKey = &payload
		return nil
	case "SetBootstrapToken":
		var payload SetBootstrapToken
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command json", c.RequestType)
		}
		c.SetBootstrapToken = &payload
		return nil
	default:
		return fmt.Errorf("mdm: unknown RequestType: %s", c.RequestType)
	}
}
