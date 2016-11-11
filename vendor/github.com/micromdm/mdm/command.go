package mdm

import (
	"fmt"

	"github.com/satori/go.uuid"
)

// CommandRequest represents an MDM command request
type CommandRequest struct {
	// Included with every payload
	RequestType string `json:"request_type"`
	UDID        string `json:"udid"`
	// DeviceInformation request
	Queries []string `json:"queries,omitempty"`
	InstallApplication
	AccountConfiguration
	ScheduleOSUpdateScan
	InstallProfile
	RemoveProfile
	InstallProvisioningProfile
	RemoveProvisioningProfile
	InstalledApplicationList
	DeviceLock
	ClearPasscode
	EraseDevice
	RequestMirroring
	Restrictions
	DeleteUser
	EnableLostMode
}

// Payload is an MDM payload
type Payload struct {
	CommandUUID string
	Command     *command
}

type command struct {
	RequestType string `json:"request_type"`
	DeviceInformation
	InstallApplication
	InstallProfile
	InstalledApplicationList
	AccountConfiguration
	ScheduleOSUpdateScan
}

// The following commands are in the order provided by the apple documentation.

// InstallProfile is an InstallProfile MDM Command
type InstallProfile struct {
	Payload []byte `plist:",omitempty" json:"payload,omitempty"`
}

type RemoveProfile struct {
	Identifier string `plist:",omitempty" json:"identifier,omitempty"`
}

type InstallProvisioningProfile struct {
	ProvisioningProfile []byte `plist:",omitempty" json:"provisioning_profile,omitempty"`
}

type RemoveProvisioningProfile struct {
	UUID string `plist:",omitempty" json:"uuid,omitempty"`
}

type InstalledApplicationList struct {
	Identifiers     []string `plist:",omitempty" json:"identifiers,omitempty"`
	ManagedAppsOnly bool     `plist:",omitempty" json:"managed_apps_only,omitempty"`
}

// DeviceInformation is a DeviceInformation MDM Command
type DeviceInformation struct {
	Queries []string `plist:",omitempty" json:"queries,omitempty"`
}

type DeviceLock struct {
	PIN         string `json:"pin,omitempty"`
	Message     string `plist:",omitempty" json:"message,omitempty"`
	PhoneNumber string `plist:",omitempty" json:"phone_number,omitempty"`
}

type ClearPasscode struct {
	UnlockToken []byte `plist:",omitempty" json:"unlock_token,omitempty"`
}

type EraseDevice struct {
	PIN string `plist:",omitempty" json:"pin,omitempty"`
}

type RequestMirroring struct {
	DestinationName     string `plist:",omitempty" json:"destination_name,omitempty"`
	DestinationDeviceID string `plist:",omitempty" json:"destination_device_id,omitempty"`
	ScanTime            string `plist:",omitempty" json:"scan_time,omitempty"`
	Password            string `plist:",omitempty" json:"password,omitempty"`
}

type Restrictions struct {
	ProfileRestrictions bool `plist:",omitempty" json:"restrictions,omitempty"`
}

type DeleteUser struct {
	UserName      string `plist:",omitempty" json:"user_name,omitempty"`
	ForceDeletion bool   `plist:",omitempty" json:"force_deletion,omitempty"`
}

type EnableLostMode struct {
	Message     string `plist:",omitempty" json:"message,omitempty"`
	PhoneNumber string `plist:",omitempty" json:"phone_number,omitempty"`
	Footnote    string `plist:",omitempty" json:"footnote,omitempty"`
}

// InstallApplication is an InstallApplication MDM Command
type InstallApplication struct {
	ITunesStoreID         int                       `plist:"iTunesStoreID,omitempty" json:"itunes_store_id,omitempty"`
	Identifier            string                    `plist:",omitempty" json:"identifier,omitempty"`
	ManifestURL           string                    `plist:",omitempty" json:"manifest_url,omitempty"`
	ManagementFlags       int                       `plist:",omitempty" json:"management_flags,omitempty"`
	NotManaged            bool                      `plist:",omitempty" json:"not_managed,omitempty"`
	ChangeManagementState string                    `plist:",omitempty" json:"change_management_state,omitempty"`
	Options               InstallApplicationOptions `plist:",omitempty" json:"options,omitempty"`
	// TODO: add remaining optional fields
}

type InstallApplicationConfiguration struct {
	// TODO: managed app config
}

type InstallApplicationOptions struct {
	NotManaged     bool `plist:",omitempty" json:"not_managed,omitempty"`
	PurchaseMethod int  `plist:",omitempty" json:"purchase_method,omitempty"`
}

type ApplyRedemptionCode struct {
	Identifier     string `plist:",omitempty" json:"identifier,omitempty"`
	RedemptionCode string `plist:",omitempty" json:"redemption_code,omitempty"`
}

type ManagedApplicationList struct {
	Identifiers []string `plist:",omitempty" json:"identifiers,omitempty"`
}

type RemoveApplication struct {
	Identifier string `json:"identifier,omitempty"`
}

type InviteToProgram struct {
	ProgramID     string `json:"program_id,omitempty"`
	InvitationURL string `json:"invitation_url,omitempty"`
}

type ValidateApplications struct {
	Identifiers []string `plist:",omitempty" json:"identifiers,omitempty"`
}

type InstallMedia struct {
	ITunesStoreID int    `plist:"iTunesStoreID,omitempty" json:"itunes_store_id,omitempty"`
	MediaURL      string `plist:",omitempty" json:"media_url,omitempty"`
	MediaType     string `json:"media_type"`
	// TODO: media url fields
}

type RemoveMedia struct {
	MediaType     string `json:"media_type"`
	ITunesStoreID int    `plist:"iTunesStoreID,omitempty" json:"itunes_store_id,omitempty"`
	PersistentID  string `plist:",omitempty" json:"persistent_id,omitempty"`
}

type Settings struct {
	Settings []Setting `json:"settings"`
}

type ManagedApplicationConfiguration struct {
	Identifiers []string `plist:",omitempty" json:"identifiers,omitempty"`
}

type ApplicationConfiguration struct {
	Identifier    string            `json:"identifier,omitempty"`
	Configuration map[string]string `plist:",omitempty" json:"configuration,omitempty"` // TODO: string map is temporary
}

type ManagedApplicationAttributes struct {
	Identifiers []string `plist:",omitempty" json:"identifiers,omitempty"`
}

type ManagedApplicationFeedback struct {
	Identifiers    []string `plist:",omitempty" json:"identifiers,omitempty"`
	DeleteFeedback bool     `plist:",omitempty" json:"delete_feedback,omitempty"`
}

// AccountConfiguration is an MDM command to create a primary user on OS X
// It allows skipping the UI to set up a user.
type AccountConfiguration struct {
	SkipPrimarySetupAccountCreation     bool           `plist:",omitempty" json:"skip_primary_setup_account_creation,omitempty"`
	SetPrimarySetupAccountAsRegularUser bool           `plist:",omitempty" json:"skip_primary_setup_account_as_regular_user,omitempty"`
	AutoSetupAdminAccounts              []AdminAccount `plist:",omitempty" json:"auto_setup_admin_accounts,omitempty"`
}

// AdminAccount is the configuration for the
// Admin account created during Setup Assistant
type AdminAccount struct {
	ShortName    string `plist:"shortName" json:"short_name"`
	FullName     string `plist:"fullName,omitempty" json:"full_name,omitempty"`
	PasswordHash data   `plist:"passwordHash" json:"password_hash"`
	Hidden       bool   `plist:"hidden,omitempty" json:"hidden,omitempty"`
}

type OSUpdate struct {
	ProductKey    string `json:"product_key"`
	InstallAction string `json:"install_action"`
}

// ScheduleOSUpdate runs update(s) immediately
type ScheduleOSUpdate struct {
	Updates []OSUpdate `plist:",omitempty" json:"updates,omitempty"`
}

// ScheduleOSUpdateScan schedules a (background) OS SoftwareUpdate check
type ScheduleOSUpdateScan struct {
	Force bool `plist:",omitempty" json:"force,omitempty"`
}

type data []byte

func newPayload(requestType string) *Payload {
	u := uuid.NewV4()
	return &Payload{u.String(),
		&command{RequestType: requestType}}
}

// NewPayload creates an MDM Payload
func NewPayload(request *CommandRequest) (*Payload, error) {
	requestType := request.RequestType
	payload := newPayload(requestType)
	switch requestType {
	case "DeviceInformation":
		payload.Command.DeviceInformation.Queries = request.Queries
	case "ScheduleOSUpdateScan":
		payload.Command.ScheduleOSUpdateScan = request.ScheduleOSUpdateScan
	case "ProfileList",
		"ProvisioningProfileList",
		"CertificateList",
		"SecurityInfo",
		"StopMirroring",
		"ClearRestrictionsPassword",
		"UsersList",
		"LogOutUser",
		"DisableLostMode",
		"DeviceLocation",
		"ManagedMediaList",
		"OSUpdateStatus",
		"DeviceConfigured",
		"AvailableOSUpdates",
		"Restrictions":
		return payload, nil
	case "InstallApplication":
		payload.Command.InstallApplication = request.InstallApplication
	case "InstallProfile":
		payload.Command.InstallProfile = request.InstallProfile
	case "AccountConfiguration":
		payload.Command.AccountConfiguration = request.AccountConfiguration
	case "InstalledApplicationList":
		payload.Command.InstalledApplicationList = request.InstalledApplicationList
	default:
		return nil, fmt.Errorf("Unsupported MDM RequestType %v", requestType)
	}
	return payload, nil
}
