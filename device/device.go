package device

import (
	"database/sql/driver"
	"errors"
	"time"

	"github.com/micromdm/dep"
	//"github.com/micromdm/mdm"
	"database/sql"
)

type JsonNullString struct {
	sql.NullString
}

func (ns *JsonNullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return []byte("\"" + ns.String + "\""), nil
	}

	return []byte("null"), nil
}

// Device represents an iOS or OS X Computer
type Device struct {
	// Primary key is UUID
	UUID         string         `json:"uuid" db:"device_uuid"`
	UDID         JsonNullString `json:"udid,omitempty"`
	SerialNumber JsonNullString `json:"serial_number,omitempty" db:"serial_number,omitempty"`
	OSVersion    string         `json:"os_version,omitempty" db:"os_version,omitempty"`
	BuildVersion string         `json:"build_version,omitempty" db:"build_version,omitempty"`
	ProductName  string         `json:"product_name,omitempty" db:"product_name,omitempty"`
	IMEI         string         `json:"imei,omitempty" db:"imei,omitempty"`
	MEID         string         `json:"meid,omitempty" db:"meid,omitempty"`
	//Apple MDM Protocol Topic
	MDMTopic               string           `json:"mdm_topic,omitempty" db:"apple_mdm_topic,omitempty"`
	PushMagic              string           `json:"push_magic,omitempty" db:"apple_push_magic,omitempty"`
	AwaitingConfiguration  bool             `json:"awaiting_configuration,omitempty" db:"awaiting_configuration,omitempty"`
	Token                  string           `json:"token,omitempty" db:"apple_mdm_token,omitempty"`
	UnlockToken            string           `json:"unlock_token,omitempty" db:"unlock_token,omitempty"`
	Enrolled               bool             `json:"enrolled,omitempty" db:"mdm_enrolled,omitempty"`
	Workflow               string           `json:"workflow_uuid,omitempty" db:"workflow_uuid,omitempty"`
	DEPDevice              bool             `json:"dep_device,omitempty" db:"dep_device,omitempty"`
	Description            string           `json:"description,omitempty" db:"description"`
	Model                  string           `json:"model,omitempty" db:"model"`
	Color                  string           `json:"color,omitempty" db:"color"`
	AssetTag               string           `json:"asset_tag,omitempty" db:"asset_tag"`
	DEPProfileStatus       DEPProfileStatus `json:"dep_profile_status,omitempty" db:"dep_profile_status"`
	DEPProfileUUID         string           `json:"dep_profile_uuid,omitempty" db:"dep_profile_uuid"`
	DEPProfileAssignTime   time.Time        `json:"dep_profile_assign_time,omitempty" db:"dep_profile_assign_time"`
	DEPProfilePushTime     time.Time        `json:"dep_profile_push_time,omitempty" db:"dep_profile_push_time"`
	DEPProfileAssignedDate time.Time        `json:"dep_profile_assigned_date,omitempty" db:"dep_profile_assigned_date"`
	DEPProfileAssignedBy   string           `json:"dep_profile_assigned_by,omitempty" db:"dep_profile_assigned_by"`
	LastCheckin            time.Time        `json:"last_checkin" db:"last_checkin"`
	DeviceName             string           `json:"device_name" db:"device_name"`
	LastQueryResponse      []byte           `json:"last_query_response" db:"last_query_response"`
}

// DEPProfileStatus is the status of the DEP Profile
// can be either "empty", "assigned", "pushed", or "removed"
type DEPProfileStatus string

// DEPProfileStatus values
const (
	EMPTY    DEPProfileStatus = "empty"
	ASSIGNED                  = "assigned"
	PUSHED                    = "pushed"
	REMOVED                   = "removed"
)

// Value implements Valuer from database/sql
func (status DEPProfileStatus) Value() (driver.Value, error) {
	if status == "" {
		return "empty", nil
	}
	return string(status), nil
}

// Scan implements Scanner from database/sql
func (status *DEPProfileStatus) Scan(value interface{}) error {
	if value == nil {
		*status = EMPTY
		return nil
	}
	if sv, err := driver.String.ConvertValue(value); err == nil {
		if v, ok := sv.(string); ok {
			*status = DEPProfileStatus(v)
			return nil
		}
	}
	// otherwise, return an error
	return errors.New("failed to scan DEPProfileStatus")

}

// NewFromDEP returns a device from DEP response
func NewFromDEP(dd dep.Device) *Device {
	var SerialNumber JsonNullString
	SerialNumber.Scan(dd.SerialNumber)

	return &Device{
		SerialNumber:           SerialNumber,
		Model:                  dd.Model,
		Description:            dd.Description,
		Color:                  dd.Color,
		AssetTag:               dd.AssetTag,
		DEPProfileStatus:       DEPProfileStatus(dd.ProfileStatus),
		DEPProfileUUID:         dd.ProfileUUID,
		DEPProfileAssignTime:   dd.ProfileAssignTime,
		DEPProfilePushTime:     dd.ProfilePushTime,
		DEPProfileAssignedDate: dd.DeviceAssignedDate,
		DEPProfileAssignedBy:   dd.DeviceAssignedBy,
	}
}
