package enroll

import (
	"github.com/satori/go.uuid"
	"time"
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
	ConsentText              map[string]string `json:"consent_text" db:"consent_text" plist:"omitempty"`
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

// TODO: Actually this is one of those non-nested payloads that doesnt respect the PayloadContent key.
type MDMPayloadContent struct {
	Payload
	AccessRights            int
	CheckInURL              string
	CheckOutWhenRemoved     bool
	IdentityCertificateUUID string
	ServerCapabilities      []string `plist:"ServerCapabilities,omitempty"`
	ServerURL               string
	Topic                   string
}
