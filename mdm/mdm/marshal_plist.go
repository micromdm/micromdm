package mdm

import "fmt"

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
		"TODO_remove":
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
	case "AccountConfiguration":
		return &struct {
			RequestType string
			*AccountConfiguration
		}{
			RequestType:          c.RequestType,
			AccountConfiguration: c.AccountConfiguration,
		}, nil
	default:
		return nil, fmt.Errorf("mdm: unknown command RequestType, %s", c.RequestType)
	}
}
