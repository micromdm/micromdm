package mdm

import (
	"encoding/json"
	"fmt"
)

func (c *Command) MarshalJSON() ([]byte, error) {
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
		var x = struct {
			RequestType string `json:"request_type"`
		}{
			RequestType: c.RequestType,
		}
		return json.Marshal(&x)
	case "InstallProfile":
		var x = struct {
			RequestType string `json:"request_type"`
			*InstallProfile
		}{
			RequestType:    c.RequestType,
			InstallProfile: c.InstallProfile,
		}
		return json.Marshal(&x)
	case "RemoveProfile":
		var x = struct {
			RequestType string `json:"request_type"`
			*RemoveProfile
		}{
			RequestType:   c.RequestType,
			RemoveProfile: c.RemoveProfile,
		}
		return json.Marshal(&x)
	case "InstallProvisioningProfile":
		var x = struct {
			RequestType string `json:"request_type"`
			*InstallProvisioningProfile
		}{
			RequestType:                c.RequestType,
			InstallProvisioningProfile: c.InstallProvisioningProfile,
		}
		return json.Marshal(&x)
	case "RemoveProvisioningProfile":
		var x = struct {
			RequestType string `json:"request_type"`
			*RemoveProvisioningProfile
		}{
			RequestType:               c.RequestType,
			RemoveProvisioningProfile: c.RemoveProvisioningProfile,
		}
		return json.Marshal(&x)
	case "InstalledApplicationList":
		var x = struct {
			RequestType string `json:"request_type"`
			*InstalledApplicationList
		}{
			RequestType:              c.RequestType,
			InstalledApplicationList: c.InstalledApplicationList,
		}
		return json.Marshal(&x)
	case "DeviceInformation":
		var x = struct {
			RequestType string `json:"request_type"`
			*DeviceInformation
		}{
			RequestType:       c.RequestType,
			DeviceInformation: c.DeviceInformation,
		}
		return json.Marshal(&x)
	case "DeviceLock":
		var x = struct {
			RequestType string `json:"request_type"`
			*DeviceLock
		}{
			RequestType: c.RequestType,
			DeviceLock:  c.DeviceLock,
		}
		return json.Marshal(&x)
	case "ClearPasscode":
		var x = struct {
			RequestType string `json:"request_type"`
			*ClearPasscode
		}{
			RequestType:   c.RequestType,
			ClearPasscode: c.ClearPasscode,
		}
		return json.Marshal(&x)
	case "EraseDevice":
		var x = struct {
			RequestType string `json:"request_type"`
			*EraseDevice
		}{
			RequestType: c.RequestType,
			EraseDevice: c.EraseDevice,
		}
		return json.Marshal(&x)
	case "RequestMirroring":
		var x = struct {
			RequestType string `json:"request_type"`
			*RequestMirroring
		}{
			RequestType:      c.RequestType,
			RequestMirroring: c.RequestMirroring,
		}
		return json.Marshal(&x)
	case "Restrictions":
		var x = struct {
			RequestType string `json:"request_type"`
			*Restrictions
		}{
			RequestType:  c.RequestType,
			Restrictions: c.Restrictions,
		}
		return json.Marshal(&x)
	case "UnlockUserAccount":
		var x = struct {
			RequestType string `json:"request_type"`
			*UnlockUserAccount
		}{
			RequestType:       c.RequestType,
			UnlockUserAccount: c.UnlockUserAccount,
		}
		return json.Marshal(&x)
	case "DeleteUser":
		var x = struct {
			RequestType string `json:"request_type"`
			*DeleteUser
		}{
			RequestType: c.RequestType,
			DeleteUser:  c.DeleteUser,
		}
		return json.Marshal(&x)
	case "EnableLostMode":
		var x = struct {
			RequestType string `json:"request_type"`
			*EnableLostMode
		}{
			RequestType:    c.RequestType,
			EnableLostMode: c.EnableLostMode,
		}
		return json.Marshal(&x)
	case "InstallApplication":
		var x = struct {
			RequestType string `json:"request_type"`
			*InstallApplication
		}{
			RequestType:        c.RequestType,
			InstallApplication: c.InstallApplication,
		}
		return json.Marshal(&x)
	case "InstallEnterpriseApplication":
		var x = struct {
			RequestType string `json:"request_type"`
			*InstallEnterpriseApplication
		}{
			RequestType:                  c.RequestType,
			InstallEnterpriseApplication: c.InstallEnterpriseApplication,
		}
		return json.Marshal(&x)
	case "AccountConfiguration":
		var x = struct {
			RequestType string `json:"request_type"`
			*AccountConfiguration
		}{
			RequestType:          c.RequestType,
			AccountConfiguration: c.AccountConfiguration,
		}
		return json.Marshal(&x)
	case "ApplyRedemptionCode":
		var x = struct {
			RequestType string `json:"request_type"`
			*ApplyRedemptionCode
		}{
			RequestType:         c.RequestType,
			ApplyRedemptionCode: c.ApplyRedemptionCode,
		}
		return json.Marshal(&x)
	case "ManagedApplicationList":
		var x = struct {
			RequestType string `json:"request_type"`
			*ManagedApplicationList
		}{
			RequestType:            c.RequestType,
			ManagedApplicationList: c.ManagedApplicationList,
		}
		return json.Marshal(&x)
	case "RemoveApplication":
		var x = struct {
			RequestType string `json:"request_type"`
			*RemoveApplication
		}{
			RequestType:       c.RequestType,
			RemoveApplication: c.RemoveApplication,
		}
		return json.Marshal(&x)
	case "InviteToProgram":
		var x = struct {
			RequestType string `json:"request_type"`
			*InviteToProgram
		}{
			RequestType:     c.RequestType,
			InviteToProgram: c.InviteToProgram,
		}
		return json.Marshal(&x)
	case "ValidateApplications":
		var x = struct {
			RequestType string `json:"request_type"`
			*ValidateApplications
		}{
			RequestType:          c.RequestType,
			ValidateApplications: c.ValidateApplications,
		}
		return json.Marshal(&x)
	case "InstallMedia":
		var x = struct {
			RequestType string `json:"request_type"`
			*InstallMedia
		}{
			RequestType:  c.RequestType,
			InstallMedia: c.InstallMedia,
		}
		return json.Marshal(&x)
	case "RemoveMedia":
		var x = struct {
			RequestType string `json:"request_type"`
			*RemoveMedia
		}{
			RequestType: c.RequestType,
			RemoveMedia: c.RemoveMedia,
		}
		return json.Marshal(&x)
	case "Settings":
		var x = struct {
			RequestType string `json:"request_type"`
			*Settings
		}{
			RequestType: c.RequestType,
			Settings:    c.Settings,
		}
		return json.Marshal(&x)
	case "ManagedApplicationConfiguration":
		var x = struct {
			RequestType string `json:"request_type"`
			*ManagedApplicationConfiguration
		}{
			RequestType:                     c.RequestType,
			ManagedApplicationConfiguration: c.ManagedApplicationConfiguration,
		}
		return json.Marshal(&x)
	case "ManagedApplicationAttributes":
		var x = struct {
			RequestType string `json:"request_type"`
			*ManagedApplicationAttributes
		}{
			RequestType:                  c.RequestType,
			ManagedApplicationAttributes: c.ManagedApplicationAttributes,
		}
		return json.Marshal(&x)
	case "ManagedApplicationFeedback":
		var x = struct {
			RequestType string `json:"request_type"`
			*ManagedApplicationFeedback
		}{
			RequestType:                c.RequestType,
			ManagedApplicationFeedback: c.ManagedApplicationFeedback,
		}
		return json.Marshal(&x)
	case "SetFirmwarePassword":
		var x = struct {
			RequestType string `json:"request_type"`
			*SetFirmwarePassword
		}{
			RequestType:         c.RequestType,
			SetFirmwarePassword: c.SetFirmwarePassword,
		}
		return json.Marshal(&x)
	case "VerifyFirmwarePassword":
		var x = struct {
			RequestType string `json:"request_type"`
			*VerifyFirmwarePassword
		}{
			RequestType:            c.RequestType,
			VerifyFirmwarePassword: c.VerifyFirmwarePassword,
		}
		return json.Marshal(&x)
	case "SetAutoAdminPassword":
		var x = struct {
			RequestType string `json:"request_type"`
			*SetAutoAdminPassword
		}{
			RequestType:          c.RequestType,
			SetAutoAdminPassword: c.SetAutoAdminPassword,
		}
		return json.Marshal(&x)
	case "ScheduleOSUpdate":
		var x = struct {
			RequestType string `json:"request_type"`
			*ScheduleOSUpdate
		}{
			RequestType:      c.RequestType,
			ScheduleOSUpdate: c.ScheduleOSUpdate,
		}
		return json.Marshal(&x)
	case "ScheduleOSUpdateScan":
		var x = struct {
			RequestType string `json:"request_type"`
			*ScheduleOSUpdateScan
		}{
			RequestType:          c.RequestType,
			ScheduleOSUpdateScan: c.ScheduleOSUpdateScan,
		}
		return json.Marshal(&x)
	case "ActiveNSExtensions":
		var x = struct {
			RequestType string `json:"request_type"`
			*ActiveNSExtensions
		}{
			RequestType:        c.RequestType,
			ActiveNSExtensions: c.ActiveNSExtensions,
		}
		return json.Marshal(&x)
	case "RotateFileVaultKey":
		var x = struct {
			RequestType string `json:"request_type"`
			*RotateFileVaultKey
		}{
			RequestType:        c.RequestType,
			RotateFileVaultKey: c.RotateFileVaultKey,
		}
		return json.Marshal(&x)
	default:
		return nil, fmt.Errorf("mdm: unknown RequestType: %s", c.RequestType)
	}
}
