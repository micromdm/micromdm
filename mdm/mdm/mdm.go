package mdm

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

type CommandRequest struct {
	UDID string `json:"udid"`
	*Command
}

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

type CommandPayload struct {
	CommandUUID string   `json:"command_uuid"`
	Command     *Command `json:"command"`
}

func NewCommandPayload(request *CommandRequest) (*CommandPayload, error) {
	payload := &CommandPayload{
		CommandUUID: uuid.NewV4().String(),
		Command:     request.Command,
	}
	return payload, nil
}

type Command struct {
	RequestType                     string `json:"request_type"`
	InstallProfile                  *InstallProfile
	RemoveProfile                   *RemoveProfile
	InstallProvisioningProfile      *InstallProvisioningProfile
	RemoveProvisioningProfile       *RemoveProvisioningProfile
	InstalledApplicationList        *InstalledApplicationList
	DeviceInformation               *DeviceInformation
	DeviceLock                      *DeviceLock
	ClearPasscode                   *ClearPasscode
	EraseDevice                     *EraseDevice
	RequestMirroring                *RequestMirroring
	Restrictions                    *Restrictions
	UnlockUserAccount               *UnlockUserAccount
	DeleteUser                      *DeleteUser
	EnableLostMode                  *EnableLostMode
	InstallApplication              *InstallApplication
	AccountConfiguration            *AccountConfiguration
	ApplyRedemptionCode             *ApplyRedemptionCode
	ManagedApplicationList          *ManagedApplicationList
	RemoveApplication               *RemoveApplication
	InviteToProgram                 *InviteToProgram
	ValidateApplications            *ValidateApplications
	InstallMedia                    *InstallMedia
	Settings                        *Settings
	ManagedApplicationConfiguration *ManagedApplicationConfiguration
	ManagedApplicationAttributes    *ManagedApplicationAttributes
	ManagedApplicationFeedback      *ManagedApplicationFeedback
	SetFirmwarePassword             *SetFirmwarePassword
	VerifyFirmwarePassword          *VerifyFirmwarePassword
	SetAutoAdminPassword            *SetAutoAdminPassword
	ScheduleOSUpdate                *ScheduleOSUpdate
	ScheduleOSUpdateScan            *ScheduleOSUpdateScan
	ActiveNSExtensions              *ActiveNSExtensions
	RotateFileVaultKey              *RotateFileVaultKey
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
		"TODO_remove":
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
	case "InstallApplication":
		var payload InstallApplication
		if err := unmarshal(&payload); err != nil {
			return errors.Wrapf(err, "mdm: unmarshal %s command plist", requestType.RequestType)
		}
		c.InstallApplication = &payload
		return nil
	default:
		return fmt.Errorf("mdm: unknown RequestType: %s", requestType.RequestType)
	}
}

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
		"TODO_remove":
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
	case "AccountConfiguration":
		var x = struct {
			RequestType string `json:"request_type"`
			*AccountConfiguration
		}{
			RequestType:          c.RequestType,
			AccountConfiguration: c.AccountConfiguration,
		}
		return json.Marshal(&x)
	default:
		return nil, fmt.Errorf("mdm: unknown RequestType: %s", c.RequestType)
	}
}

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

// InstallProfile is an InstallProfile MDM Command
type InstallProfile struct {
	Payload []byte `json:"payload,omitempty"`
}

type RemoveProfile struct {
	Identifier string `json:"identifier,omitempty"`
}

type InstallProvisioningProfile struct {
	ProvisioningProfile []byte `plist:",omitempty" json:"provisioning_profile,omitempty"`
}

type RemoveProvisioningProfile struct {
	UUID string `json:"uuid"`
}

type InstalledApplicationList struct {
	Identifiers     []string `plist:",omitempty" json:"identifiers,omitempty"`
	ManagedAppsOnly bool     `plist:",omitempty" json:"managed_appd_only,omitempty"`
}

type DeviceInformation struct {
	Queries []string `plist:",omitempty" json:"queries,omitempty"`
}

type DeviceLock struct {
	PIN         string `plist:",omitempty" json:"pin"`
	Message     string `plist:",omitempty" json:"message,omitempty"`
	PhoneNumber string `plist:",omitempty" json:"phone_number,omitempty"`
}

type ClearPasscode struct {
	UnlockToken []byte `json:"unlock_token"`
}

type EraseDevice struct {
	PIN                    string `json:"pin"`
	PreserveDataPlan       bool   `plist:",omitempty" json:"preserve_data_plan,omitempty"`
	DisallowProximitySetup bool   `plist:",omitempty" json:"disallow_proximity_setup,omitempty"`
}

type RequestMirroring struct {
	DestinationName     string `plist:",omitempty" json:"destination_name,omitempty"`
	DestinationDeviceID string `plist:",omitempty" json:"destination_device_id,omitempty"`
	ScanTime            string `plist:",omitempty" json:"scan_time,omitempty"`
	Password            string `plist:",omitempty" json:"password,omitempty"`
}

type Restrictions struct {
	ProfileRestrictions bool `json:"profile_restrictions"`
}

type UnlockUserAccount struct {
	UserName string `json:"username"`
}

type DeleteUser struct {
	UserName      string `plist:",omitempty" json:"username,omitempty"`
	ForceDeletion bool   `plist:",omitempty" json:"force_deletion,omitempty"`
}

type EnableLostMode struct {
	Message     string `plist:",omitempty" json:"message,omitempty"`
	PhoneNumber string `plist:",omitempty" json:"phone_number,omitempty"`
	Footnote    string `plist:",omitempty" json:"footnote,omitempty"`
}

type InstallApplication struct {
	ITunesStoreID         *int64                           `plist:"iTunesStoreID,omitempty" json:"itunes_store_id,omitempty"`
	Identifier            *string                          `plist:",omitempty" json:"identifier,omitempty"`
	ManagementFlags       *int                             `plist:",omitempty" json:"management_flags,omitempty"`
	ChangeManagementState *string                          `plist:",omitempty" json:"change_management_state,omitempty"`
	ManifestURL           *string                          `plist:",omitempty" json:"manifest_url,omitempty"`
	Options               *InstallApplicationOptions       `plist:",omitempty" json:"options,omitempty"`
	Configuration         *InstallApplicationConfiguration `plist:",omitempty" json:"configuration,omitempty"`
	Attributes            *InstallApplicationAttributes    `plist:",omitempty" json:"attributes,omitempty"`
}

type InstallApplicationOptions struct {
	PurchaseMethod int64 `plist:",omitempty" json:"purchase_method,omitempty"`
}

type InstallApplicationConfiguration struct{}
type InstallApplicationAttributes struct{}

type AccountConfiguration struct {
	SkipPrimarySetupAccountCreation     bool           `plist:",omitempty" json:"skip_primary_setup_account_creation,omitempty"`
	SetPrimarySetupAccountAsRegularUser bool           `plist:",omitempty" json:"skip_primary_setup_account_as_regular_user,omitempty"`
	AutoSetupAdminAccounts              []AdminAccount `plist:",omitempty" json:"auto_setup_admin_accounts,omitempty"`
}

type AdminAccount struct {
	ShortName    string `plist:"shortName" json:"short_name"`
	FullName     string `plist:"fullName,omitempty" json:"full_name,omitempty"`
	PasswordHash []byte `plist:"passwordHash" json:"password_hash"`
	Hidden       bool   `plist:"hidden,omitempty" json:"hidden,omitempty"`
}

type ApplyRedemptionCode struct {
	Identifier     string `plist:",omitempty" json:"identifier,omitempty"`
	RedemptionCode string `plist:",omitempty" json:"redemption_code,omitempty"`
}

type ManagedApplicationList struct {
	Identifiers []string `plist:",omitempty" json:"identifiers,omitempty"`
}

type RemoveApplication struct {
	Identifier string `plist:",omitempty" json:"identifier,omitempty"`
}

type InviteToProgram struct {
	ProgramID     string `plist:",omitempty" json"program_id,omitempty"`
	InvitationURL string `plist:",omitempty" json:"invitation_url,omitempty"`
}

type ValidateApplications struct {
	Identifiers []string `plist:",omitempty" json:"identifiers,omitempty"`
}

type InstallMedia struct {
	ITunesStoreID *int64 `plist:"iTunesStoreID,omitempty" json:"itunes_store_id,omitempty"`
	MediaURL      string `plist:",omitempty" json:"media_url,omitempty"`
	MediaType     string `plist:",omitempty" json:"media_type,omitempty"`
}

type RemoveMedia struct {
	ITunesStoreID *int64 `plist:"iTunesStoreID,omitempty" json:"itunes_store_id,omitempty"`
	MediaType     string `plist:",omitempty" json:"media_type,omitempty"`
	PersistentID  string `plist:",omitempty" json:"persistent_id,omitempty"`
}

type Settings struct {
	Settings []Setting `plist:",omitempty" json:"settings,omitempty"`
}

type Setting struct {
	Item                    string                 `json:"item"`
	Enabled                 *bool                  `plist:",omitempty" json:"enabled,omitempty"`
	DeviceName              *string                `plist:",omitempty" json:"device_name,omitempty"`
	HostName                *string                `plist:",omitempty" json:"hostname,omitempty"`
	Identifier              *string                `plist:",omitempty" json:"identifier"`
	Attributes              map[string]string      `plist:",omitempty" json:"attributes,omitempty"`
	Image                   []byte                 `plist:",omitempty" json:"image,omitempty"`
	Where                   *int                   `plist:",omitempty" json:"where,omitempty"`
	MDMOptions              map[string]interface{} `plist:",omitempty" json:"mdm_options,omitempty"`
	PasscodeLockGracePeriod *int                   `plist:",omitempty" json:"passcode_lock_grace_period,omitempty"`
	MaximumResidentUsers    *int                   `plist:",omitempty" json:"maximum_resident_users,omitempty"`
}

type ManagedApplicationConfiguration struct {
	Identifiers []string `plist:",omitempty" json:"identifiers,omitempty"`
}

type ManagedApplicationAttributes struct {
	Identifiers []string `plist:",omitempty" json:"identifiers,omitempty"`
}

type ManagedApplicationFeedback struct {
	Identifiers    []string `plist:",omitempty" json:"identifiers,omitempty"`
	DeleteFeedback bool     `plist:",omitempty" json:"delete_feedback,omitempty"`
}

type SetFirmwarePassword struct {
	CurrentPassword string `plist:",omitempty" json:"current_password,omitempty"`
	NewPassword     string `plist:",omitempty" json:"new_password,omitempty"`
	AllowOroms      bool   `plist:",omitempty" json:"allow_oroms,omitempty"`
}

type VerifyFirmwarePassword struct {
	Password string `plist:",omitempty" json:"password,omitempty"`
}

type SetAutoAdminPassword struct {
	GUID         string `plist:",omitempty" json:"guid,omitempty"`
	PasswordHash []byte `plist:"passwordHash" json:"password_hash"`
}

type OSUpdate struct {
	ProductKey    string `json:"product_key"`
	InstallAction string `json:"install_action"`
}

type ScheduleOSUpdate struct {
	Updates []OSUpdate `plist:",omitempty" json:"updates,omitempty"`
}

type ScheduleOSUpdateScan struct {
	Force bool `plist:",omitempty" json:"force,omitempty"`
}

type ActiveNSExtensions struct {
	FilterExtensionPoints []string `plist:",omitempty json:"filter_extensions_points,omitempty"`
}

type RotateFileVaultKey struct {
	KeyType                    string          `plist:",omitempty" json:"key_type,omitempty"`
	FileVaultUnlock            FileVaultUnlock `plist:",omitempty" json:"filevault_unlock,omitempty"`
	NewCertificate             []byte          `plist:",omitempty" json:"new_certificate,omitempty"`
	ReplyEncryptionCertificate []byte          `plist:",omitempty" json:"reply_encryption_certificate,omitempty"`
}

type FileVaultUnlock struct {
	Password                 string `plist:",omitempty" json:"password,omitempty"`
	PrivateKeyExport         []byte `plist:",omitempty" json:"private_key_export,omitempty"`
	PrivateKeyExportPassword string `plist:",omitempty" json:"private_key_export_password,omitempty"`
}
