package applications

type Application struct {
	Identifier   string `plist:",omitempty" json:"identifier,omitempty"`
	Version      string `plist:",omitempty" json:"version,omitempty"`
	ShortVersion string `plist:",omitempty" json:"short_version,omitempty"`
	Name         string `json:"name,omitempty"`
	BundleSize   int    `plist:",omitempty" json:"bundle_size,omitempty"`

	// The size of the app's document, library, and other folders, in bytes.
	DynamicSize int `plist:",omitempty" json:"dynamic_size,omitempty"`

	IsValidated bool `plist:",omitempty" json:"is_validated,omitempty"`
}
