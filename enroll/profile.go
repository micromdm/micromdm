package enroll

import (
	"time"

	"github.com/satori/go.uuid"
)

type Payload struct {
	PayloadType         string      `json:"type" db:"type"`
	PayloadVersion      int         `json:"version" db:"version"`
	PayloadIdentifier   string      `json:"identifier" db:"identifier"`
	PayloadUUID         string      `json:"uuid" db:"uuid"`
	PayloadDisplayName  string      `json:"displayname" db:"displayname"`
	PayloadDescription  string      `json:"description,omitempty" db:"description"`
	PayloadOrganization string      `json:"organization,omitempty" db:"organization"`
	PayloadScope        string      `json:"scope" db:"scope" plist:",omitempty"`
	PayloadContent      interface{} `json:"content,omitempty" plist:"PayloadContent,omitempty"`
}

type Profile struct {
	PayloadContent           []interface{}     `json:"content,omitempty" db:"content"`
	PayloadDescription       string            `json:"description,omitempty" db:"description"`
	PayloadDisplayName       string            `json:"displayname,omitempty" db:"displayname"`
	PayloadExpirationDate    *time.Time        `json:"expiration_date,omitempty" db:"expiration_date" plist:",omitempty"`
	PayloadIdentifier        string            `json:"identifier" db:"identifier"`
	PayloadOrganization      string            `json:"organization,omitempty" db:"organization"`
	PayloadUUID              string            `json:"uuid" db:"uuid"`
	PayloadRemovalDisallowed bool              `json:"removal_disallowed" db:"removal_disallowed" plist:",omitempty"`
	PayloadType              string            `json:"type" db:"type"`
	PayloadVersion           int               `json:"version" db:"version"`
	PayloadScope             string            `json:"scope" db:"scope" plist:",omitempty"`
	RemovalDate              *time.Time        `json:"removal_date" db:"removal_date" plist:"-" plist:",omitempty"`
	DurationUntilRemoval     float32           `json:"duration_until_removal" db:"duration_until_removal" plist:",omitempty"`
	ConsentText              map[string]string `json:"consent_text" db:"consent_text" plist:",omitempty"`
}

func NewProfile() *Profile {
	payloadUuid := uuid.NewV4()

	return &Profile{
		PayloadVersion: 1,
		PayloadType:    "Configuration",
		PayloadUUID:    payloadUuid.String(),
	}
}

func NewPayload(payloadType string) *Payload {
	payloadUuid := uuid.NewV4()

	return &Payload{
		PayloadVersion: 1,
		PayloadType:    payloadType,
		PayloadUUID:    payloadUuid.String(),
	}
}

type SCEPPayloadContent struct {
	CAFingerprint []byte `plist:"CAFingerprint,omitempty"` // NSData
	Challenge     string `plist:"Challenge,omitempty"`
	Keysize       int
	KeyType       string `plist:"Key Type"`
	KeyUsage      int    `plist:"Key Usage"`
	Name          string
	Subject       [][][]string `plist:"Subject,omitempty"`
	URL           string
}

// AccessRights define the management rights of the MDM server over the device.
// May not be zero. If 2 is specified, 1 must also be specified. If 128 is specified, 64 must also be specified.
type AccessRights int

const (
	// Allow inspection of installed configuration profiles.
	ProfileInspection AccessRights = 1 << iota

	// Allow installation and removal of configuration profiles.
	ProfileInstallAndRemoval

	// Allow device lock and passcode removal.
	DeviceLock

	// Allow device erase.
	DeviceErase

	// Allow query of Device Information (device capacity, serial number).
	DeviceInformationQuery

	// 	Allow query of Network Information (phone/SIM numbers, MAC addresses).
	NetworkInformationQuery

	// Allow inspection of installed provisioning profiles.
	ProvisioningProfileInspection

	//  Allow installation and removal of provisioning profiles.
	ProvisioningProfileInstallAndRemoval

	// Allow inspection of installed applications.
	ApplicationInspection

	// Allow restriction-related queries.
	RestrictionQuery

	// Allow security-related queries.
	SecurityQuery

	// Allow manipulation of settings.
	// Availability: Available in iOS 5.0 and later. Available in macOS 10.9 for certain commands.
	SettingsManipulation

	// Allow app management.
	// Availability: Available in iOS 5.0 and later. Available in macOS 10.9 for certain commands.
	AppManagement
)

func allRights() AccessRights {
	return ProfileInspection |
		ProfileInstallAndRemoval |
		DeviceLock |
		DeviceErase |
		DeviceInformationQuery |
		NetworkInformationQuery |
		ProvisioningProfileInspection |
		ProvisioningProfileInstallAndRemoval |
		ApplicationInspection |
		RestrictionQuery |
		SecurityQuery |
		SettingsManipulation |
		AppManagement
}

// TODO: Actually this is one of those non-nested payloads that doesnt respect the PayloadContent key.
type MDMPayloadContent struct {
	Payload
	AccessRights            AccessRights
	CheckInURL              string
	CheckOutWhenRemoved     bool
	IdentityCertificateUUID string
	ServerCapabilities      []string `plist:"ServerCapabilities,omitempty"`
	SignMessage             bool     `plist:"SignMessage,omitempty"`
	ServerURL               string
	Topic                   string
}

type ProfileServicePayload struct {
	URL              string
	Challenge        string `plist:",omitempty"`
	DeviceAttributes []string
}
