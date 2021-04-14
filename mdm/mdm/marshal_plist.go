package mdm

import (
	"fmt"

	"github.com/groob/plist"
	"github.com/pkg/errors"
)

func (c *Command) MarshalPlist() (interface{}, error) {
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
		return &struct {
			RequestType string
		}{
			RequestType: c.RequestType,
		}, nil

	case "InstallProfile":
		return &struct {
			RequestType string
			*InstallProfile
		}{
			RequestType:    c.RequestType,
			InstallProfile: c.InstallProfile,
		}, nil
	case "RemoveProfile":
		return &struct {
			RequestType string
			*RemoveProfile
		}{
			RequestType:   c.RequestType,
			RemoveProfile: c.RemoveProfile,
		}, nil
	case "InstallProvisioningProfile":
		return &struct {
			RequestType string
			*InstallProvisioningProfile
		}{
			RequestType:                c.RequestType,
			InstallProvisioningProfile: c.InstallProvisioningProfile,
		}, nil
	case "RemoveProvisioningProfile":
		return &struct {
			RequestType string
			*RemoveProvisioningProfile
		}{
			RequestType:               c.RequestType,
			RemoveProvisioningProfile: c.RemoveProvisioningProfile,
		}, nil
	case "InstalledApplicationList":
		return &struct {
			RequestType string
			*InstalledApplicationList
		}{
			RequestType:              c.RequestType,
			InstalledApplicationList: c.InstalledApplicationList,
		}, nil
	case "DeviceInformation":
		return &struct {
			RequestType string
			*DeviceInformation
		}{
			RequestType:       c.RequestType,
			DeviceInformation: c.DeviceInformation,
		}, nil
	case "DeviceLock":
		return &struct {
			RequestType string
			*DeviceLock
		}{
			RequestType: c.RequestType,
			DeviceLock:  c.DeviceLock,
		}, nil
	case "ClearPasscode":
		return &struct {
			RequestType string
			*ClearPasscode
		}{
			RequestType:   c.RequestType,
			ClearPasscode: c.ClearPasscode,
		}, nil
	case "EraseDevice":
		return &struct {
			RequestType string
			*EraseDevice
		}{
			RequestType: c.RequestType,
			EraseDevice: c.EraseDevice,
		}, nil
	case "RequestMirroring":
		return &struct {
			RequestType string
			*RequestMirroring
		}{
			RequestType:      c.RequestType,
			RequestMirroring: c.RequestMirroring,
		}, nil
	case "Restrictions":
		return &struct {
			RequestType string
			*Restrictions
		}{
			RequestType:  c.RequestType,
			Restrictions: c.Restrictions,
		}, nil
	case "UnlockUserAccount":
		return &struct {
			RequestType string
			*UnlockUserAccount
		}{
			RequestType:       c.RequestType,
			UnlockUserAccount: c.UnlockUserAccount,
		}, nil
	case "DeleteUser":
		return &struct {
			RequestType string
			*DeleteUser
		}{
			RequestType: c.RequestType,
			DeleteUser:  c.DeleteUser,
		}, nil
	case "EnableLostMode":
		return &struct {
			RequestType string
			*EnableLostMode
		}{
			RequestType:    c.RequestType,
			EnableLostMode: c.EnableLostMode,
		}, nil
	case "InstallApplication":
		return &struct {
			RequestType string
			*InstallApplication
		}{
			RequestType:        c.RequestType,
			InstallApplication: c.InstallApplication,
		}, nil
	case "InstallEnterpriseApplication":
		return &struct {
			RequestType string
			*InstallEnterpriseApplication
		}{
			RequestType:                  c.RequestType,
			InstallEnterpriseApplication: c.InstallEnterpriseApplication,
		}, nil
	case "AccountConfiguration":
		return &struct {
			RequestType string
			*AccountConfiguration
		}{
			RequestType:          c.RequestType,
			AccountConfiguration: c.AccountConfiguration,
		}, nil
	case "ApplyRedemptionCode":
		return &struct {
			RequestType string
			*ApplyRedemptionCode
		}{
			RequestType:         c.RequestType,
			ApplyRedemptionCode: c.ApplyRedemptionCode,
		}, nil
	case "ManagedApplicationList":
		return &struct {
			RequestType string
			*ManagedApplicationList
		}{
			RequestType:            c.RequestType,
			ManagedApplicationList: c.ManagedApplicationList,
		}, nil
	case "RemoveApplication":
		return &struct {
			RequestType string
			*RemoveApplication
		}{
			RequestType:       c.RequestType,
			RemoveApplication: c.RemoveApplication,
		}, nil
	case "InviteToProgram":
		return &struct {
			RequestType string
			*InviteToProgram
		}{
			RequestType:     c.RequestType,
			InviteToProgram: c.InviteToProgram,
		}, nil
	case "ValidateApplications":
		return &struct {
			RequestType string
			*ValidateApplications
		}{
			RequestType:          c.RequestType,
			ValidateApplications: c.ValidateApplications,
		}, nil
	case "InstallMedia":
		return &struct {
			RequestType string
			*InstallMedia
		}{
			RequestType:  c.RequestType,
			InstallMedia: c.InstallMedia,
		}, nil
	case "RemoveMedia":
		return &struct {
			RequestType string
			*RemoveMedia
		}{
			RequestType: c.RequestType,
			RemoveMedia: c.RemoveMedia,
		}, nil
	case "Settings":
		// convert all the data plists into the dictionary inside settings before serialization.
		for i, set := range c.Settings.Settings {
			if len(set.ConfigurationData) > 0 {
				var configuration map[string]interface{}
				if err := plist.Unmarshal(set.ConfigurationData, &configuration); err != nil {
					return nil, errors.Wrap(err, "turning the configuration data plist into a dictionary")
				}
				set.Configuration = configuration
				c.Settings.Settings[i] = set
			}
		}
		return &struct {
			RequestType string
			*Settings
		}{
			RequestType: c.RequestType,
			Settings:    c.Settings,
		}, nil
	case "ManagedApplicationConfiguration":
		return &struct {
			RequestType string
			*ManagedApplicationConfiguration
		}{
			RequestType:                     c.RequestType,
			ManagedApplicationConfiguration: c.ManagedApplicationConfiguration,
		}, nil
	case "ManagedApplicationAttributes":
		return &struct {
			RequestType string
			*ManagedApplicationAttributes
		}{
			RequestType:                  c.RequestType,
			ManagedApplicationAttributes: c.ManagedApplicationAttributes,
		}, nil
	case "ManagedApplicationFeedback":
		return &struct {
			RequestType string
			*ManagedApplicationFeedback
		}{
			RequestType:                c.RequestType,
			ManagedApplicationFeedback: c.ManagedApplicationFeedback,
		}, nil
	case "SetFirmwarePassword":
		return &struct {
			RequestType string
			*SetFirmwarePassword
		}{
			RequestType:         c.RequestType,
			SetFirmwarePassword: c.SetFirmwarePassword,
		}, nil
	case "VerifyFirmwarePassword":
		return &struct {
			RequestType string
			*VerifyFirmwarePassword
		}{
			RequestType:            c.RequestType,
			VerifyFirmwarePassword: c.VerifyFirmwarePassword,
		}, nil
	case "SetAutoAdminPassword":
		return &struct {
			RequestType string
			*SetAutoAdminPassword
		}{
			RequestType:          c.RequestType,
			SetAutoAdminPassword: c.SetAutoAdminPassword,
		}, nil
	case "ScheduleOSUpdate":
		return &struct {
			RequestType string
			*ScheduleOSUpdate
		}{
			RequestType:      c.RequestType,
			ScheduleOSUpdate: c.ScheduleOSUpdate,
		}, nil
	case "ScheduleOSUpdateScan":
		return &struct {
			RequestType string
			*ScheduleOSUpdateScan
		}{
			RequestType:          c.RequestType,
			ScheduleOSUpdateScan: c.ScheduleOSUpdateScan,
		}, nil
	case "ActiveNSExtensions":
		return &struct {
			RequestType string
			*ActiveNSExtensions
		}{
			RequestType:        c.RequestType,
			ActiveNSExtensions: c.ActiveNSExtensions,
		}, nil
	case "RotateFileVaultKey":
		return &struct {
			RequestType string
			*RotateFileVaultKey
		}{
			RequestType:        c.RequestType,
			RotateFileVaultKey: c.RotateFileVaultKey,
		}, nil
	default:
		return nil, fmt.Errorf("mdm: unknown command RequestType, %s", c.RequestType)
	}
}
