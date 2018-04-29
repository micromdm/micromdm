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
	}{}
	if err := json.Unmarshal(data, &request); err != nil {
		return errors.Wrap(err, "mdm: unmarshal json command request")
	}
	c.UDID = request.UDID
	c.Command = &Command{
		RequestType: request.RequestType,
	}
	return c.Command.UnmarshalJSON(data)
}

func (c *Command) UnmarshalJSON(data []byte) error {
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
	default:
		return fmt.Errorf("mdm: unknown RequestType: %s", c.RequestType)
	}
}
