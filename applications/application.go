package applications

import "database/sql"

type Application struct {
	UUID         string         `plist:",omitempty" json:"uuid,omitempty" db:"application_uuid"`
	Identifier   sql.NullString `plist:",omitempty" json:"identifier,omitempty" db:"identifier"`
	Version      sql.NullString `plist:",omitempty" json:"version,omitempty" db:"version"`
	ShortVersion sql.NullString `plist:",omitempty" json:"short_version,omitempty" db:"short_version"`
	Name         string         `json:"name,omitempty" db:"name"`
	BundleSize   sql.NullInt64  `plist:",omitempty" json:"bundle_size,omitempty" db:"bundle_size"`

	// The size of the app's document, library, and other folders, in bytes. Only applies to iOS
	DynamicSize sql.NullInt64 `plist:",omitempty" json:"dynamic_size,omitempty" db:"dynamic_size"`

	// iOS only.
	IsValidated sql.NullBool `plist:",omitempty" json:"is_validated,omitempty" db:"is_validated"`
}

type DeviceApplication struct {
	DeviceUUID      string `json:"device_uuid" db:"device_uuid"`
	ApplicationUUID string `json:"application_uuid" db:"application_uuid"`

	Identifier   sql.NullString `plist:",omitempty" json:"identifier,omitempty" db:"identifier"`
	Version      sql.NullString `plist:",omitempty" json:"version,omitempty" db:"version"`
	ShortVersion sql.NullString `plist:",omitempty" json:"short_version,omitempty" db:"short_version"`
	Name         string         `json:"name,omitempty" db:"name"`
	BundleSize   sql.NullInt64  `plist:",omitempty" json:"bundle_size,omitempty" db:"bundle_size"`

	// The size of the app's document, library, and other folders, in bytes. Only applies to iOS
	DynamicSize sql.NullInt64 `plist:",omitempty" json:"dynamic_size,omitempty" db:"dynamic_size"`

	// iOS only.
	IsValidated sql.NullBool `plist:",omitempty" json:"is_validated,omitempty" db:"is_validated"`
}
