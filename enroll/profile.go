package enroll

import "time"

type Payload struct {
	PayloadType         string      `json:"type" db:"type"`
	PayloadVersion      int         `json:"version" db:"version"`
	PayloadIdentifier   string      `json:"identifier" db:"identifier"`
	PayloadUUID         string      `json:"uuid" db:"uuid"`
	PayloadDisplayName  string      `json:"displayname" db:"displayname"`
	PayloadDescription  string      `json:"description,omitempty" db:"description"`
	PayloadOrganization string      `json:"organization,omitempty" db:"organization"`
	PayloadContent      interface{} `json:"content,omitempty"`
}

type Profile struct {
	PayloadContent           []interface{}     `json:"content,omitempty" db:"content"`
	PayloadDescription       string            `json:"description,omitempty" db:"description"`
	PayloadDisplayName       string            `json:"displayname,omitempty" db:"displayname"`
	PayloadExpirationDate    *time.Time        `json:"expiration_date,omitempty" db:"expiration_date" plist:"omitempty"`
	PayloadIdentifier        string            `json:"identifier" db:"identifier"`
	PayloadOrganization      string            `json:"organization,omitempty" db:"organization"`
	PayloadUUID              string            `json:"uuid" db:"uuid"`
	PayloadRemovalDisallowed bool              `json:"removal_disallowed" db:"removal_disallowed" plist:"omitempty"`
	PayloadType              string            `json:"type" db:"type"`
	PayloadVersion           int               `json:"version" db:"version"`
	PayloadScope             string            `json:"scope" db:"scope" plist:"omitempty"`
	RemovalDate              *time.Time        `json:"removal_date" db:"removal_date" plist:"-" plist:"omitempty"`
	DurationUntilRemoval     float32           `json:"duration_until_removal" db:"duration_until_removal" plist:"omitempty"`
	ConsentText              map[string]string `json:"consent_text" db:"consent_text" plist:"omitempty"`
}

func NewProfile() Profile {
	return &Profile{
		PayloadVersion: 1,
		PayloadType:    "Configuration",
	}
}

type ProfileServicePayload struct {
	URL              string
	DeviceAttributes []string
	Challenge        string
}
