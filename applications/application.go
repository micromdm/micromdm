package applications

import "database/sql"

type Application struct {
	Identifier   sql.NullString `plist:",omitempty" json:"identifier,omitempty"`
	Version      sql.NullString `plist:",omitempty" json:"version,omitempty"`
	ShortVersion sql.NullString `plist:",omitempty" json:"short_version,omitempty"`
	Name         string         `json:"name,omitempty"`
	BundleSize   sql.NullInt64  `plist:",omitempty" json:"bundle_size,omitempty"`

	// The size of the app's document, library, and other folders, in bytes.
	DynamicSize sql.NullInt64 `plist:",omitempty" json:"dynamic_size,omitempty"`

	IsValidated sql.NullBool `plist:",omitempty" json:"is_validated,omitempty"`
}
