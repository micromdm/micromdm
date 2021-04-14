package mdm

import (
	"fmt"

	"github.com/pkg/errors"
)

func (c *Command) UnmarshalPlist(unmarshal func(i interface{}) error) error {
	var requestType = struct {
		RequestType string
	}{}
	if err := unmarshal(&requestType); err != nil {
		return errors.Wrap(err, "mdm: unmarshal request type")
	}
	c.RequestType = requestType.RequestType

	switch requestType.RequestType {
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
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.InstallProfile = &payload
		return nil
	case "RemoveProfile":
		var payload RemoveProfile
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.RemoveProfile = &payload
		return nil
	case "InstallProvisioningProfile":
		var payload InstallProvisioningProfile
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.InstallProvisioningProfile = &payload
		return nil
	case "RemoveProvisioningProfile":
		var payload RemoveProvisioningProfile
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.RemoveProvisioningProfile = &payload
		return nil
	case "InstalledApplicationList":
		var payload InstalledApplicationList
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.InstalledApplicationList = &payload
		return nil
	case "DeviceInformation":
		var payload DeviceInformation
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.DeviceInformation = &payload
		return nil
	case "DeviceLock":
		var payload DeviceLock
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.DeviceLock = &payload
		return nil
	case "ClearPasscode":
		var payload ClearPasscode
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.ClearPasscode = &payload
		return nil
	case "EraseDevice":
		var payload EraseDevice
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.EraseDevice = &payload
		return nil
	case "RequestMirroring":
		var payload RequestMirroring
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.RequestMirroring = &payload
		return nil
	case "Restrictions":
		var payload Restrictions
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.Restrictions = &payload
		return nil
	case "UnlockUserAccount":
		var payload UnlockUserAccount
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.UnlockUserAccount = &payload
		return nil
	case "DeleteUser":
		var payload DeleteUser
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.DeleteUser = &payload
		return nil
	case "EnableLostMode":
		var payload EnableLostMode
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.EnableLostMode = &payload
		return nil
	case "InstallEnterpriseApplication":
		var payload InstallEnterpriseApplication
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.InstallEnterpriseApplication = &payload
		return nil
	case "InstallApplication":
		var payload InstallApplication
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.InstallApplication = &payload
		return nil
	case "AccountConfiguration":
		var payload AccountConfiguration
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.AccountConfiguration = &payload
		return nil
	case "ApplyRedemptionCode":
		var payload ApplyRedemptionCode
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.ApplyRedemptionCode = &payload
		return nil
	case "ManagedApplicationList":
		var payload ManagedApplicationList
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.ManagedApplicationList = &payload
		return nil
	case "RemoveApplication":
		var payload RemoveApplication
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.RemoveApplication = &payload
		return nil
	case "InviteToProgram":
		var payload InviteToProgram
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.InviteToProgram = &payload
		return nil
	case "ValidateApplications":
		var payload ValidateApplications
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.ValidateApplications = &payload
		return nil
	case "InstallMedia":
		var payload InstallMedia
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.InstallMedia = &payload
		return nil
	case "RemoveMedia":
		var payload RemoveMedia
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.RemoveMedia = &payload
		return nil
	case "Settings":
		var payload Settings
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.Settings = &payload
		return nil
	case "ManagedApplicationConfiguration":
		var payload ManagedApplicationConfiguration
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.ManagedApplicationConfiguration = &payload
		return nil
	case "ManagedApplicationAttributes":
		var payload ManagedApplicationAttributes
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.ManagedApplicationAttributes = &payload
		return nil
	case "ManagedApplicationFeedback":
		var payload ManagedApplicationFeedback
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.ManagedApplicationFeedback = &payload
		return nil
	case "SetFirmwarePassword":
		var payload SetFirmwarePassword
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.SetFirmwarePassword = &payload
		return nil
	case "VerifyFirmwarePassword":
		var payload VerifyFirmwarePassword
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.VerifyFirmwarePassword = &payload
		return nil
	case "SetAutoAdminPassword":
		var payload SetAutoAdminPassword
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.SetAutoAdminPassword = &payload
		return nil
	case "ScheduleOSUpdate":
		var payload ScheduleOSUpdate
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.ScheduleOSUpdate = &payload
		return nil
	case "ScheduleOSUpdateScan":
		var payload ScheduleOSUpdateScan
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.ScheduleOSUpdateScan = &payload
		return nil
	case "ActiveNSExtensions":
		var payload ActiveNSExtensions
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.ActiveNSExtensions = &payload
		return nil
	case "RotateFileVaultKey":
		var payload RotateFileVaultKey
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.RotateFileVaultKey = &payload
		return nil
	default:
		return fmt.Errorf("mdm: unknown RequestType: %s", requestType.RequestType)
	}
}
