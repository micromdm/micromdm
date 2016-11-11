package mdm

// All of these settings are changed by sending the `Settings` command.

type Setting struct {
	Item string `json:"item"`
}

type VoiceRoamingSetting struct {
	Setting
	Enabled bool `json:"enabled"`
}

type PersonalHotspotSetting struct {
	Setting
	Enabled bool `json:"enabled"`
}

type WallpaperSetting struct {
	Setting
	Image []byte `json:"image"`
	Where int    `json:"where"`
}

type DataRoamingSetting struct {
	Setting
	Enabled bool `json:"enabled"`
}

type ApplicationAttributesSetting struct {
	Setting
	Identifier string            `json:"identifier"`
	Attributes map[string]string `plist:",omitempty" json:"attributes,omitempty"`
}

type DeviceNameSetting struct {
	Setting
	DeviceName string `json:"device_name"`
}

type MDMOptions struct {
	ActivationLockAllowedWhileSupervised bool `json:"activation_lock_allowed_while_supervised"`
}

type MDMOptionsSetting struct {
	Setting
	MDMOptions MDMOptions `json:"mdm_options"`
}

type MaximumResidentUsersSetting struct {
	MaximumResidentUsers int `json:"maximum_resident_users"`
}
