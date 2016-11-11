package mdm

import "time"

// Response is an MDM Command Response
type Response struct {
	UDID                     string
	UserID                   *string `json:"user_id,omitempty" plist:"UserID,omitempty"`
	Status                   string
	CommandUUID              string
	RequestType              string                           `json:"request_type,omitempty" plist:",omitempty"`
	ErrorChain               []ErrorChainItem                 `json:"error_chain" plist:",omitempty"`
	QueryResponses           QueryResponses                   `json:"query_responses,omitempty" plist:",omitempty"`
	SecurityInfo             SecurityInfo                     `json:"security_info,omitempty" plist:",omitempty"`
	CertificateList          CertificateList                  `json:"certificate_list,omitempty" plist:",omitempty"`
	InstalledApplicationList InstalledApplicationListResponse `json:"installed_application_list,omitempty" plist:",omitempty"`
}

type ProvisioningProfileListItem struct {
	Name       string    `plist:",omitempty" json:"name,omitempty"`
	UUID       string    `plist:",omitempty" json:"uuid,omitempty"`
	ExpiryDate time.Time `plist:",omitempty" json:"expiry_date,omitempty"`
}

type ProvisioningProfileListResponse []ProvisioningProfileListItem

type CertificateListItem struct {
	CommonName string `json:"common_name"`
	Data       []byte `json:"data"`
	IsIdentity bool   `json:"is_identity"`
}

// CertificateList is the CertificateList MDM Command Response
type CertificateList []CertificateListItem

type InstalledApplicationListItem struct {
	Identifier   string `plist:",omitempty" json:"identifier,omitempty"`
	Version      string `plist:",omitempty" json:"version,omitempty"`
	ShortVersion string `plist:",omitempty" json:"short_version,omitempty"`
	Name         string `json:"name,omitempty"`
	BundleSize   uint32 `plist:",omitempty" json:"bundle_size,omitempty"`
	DynamicSize  uint32 `plist:",omitempty" json:"dynamic_size,omitempty"`
	IsValidated  bool   `plist:",omitempty" json:"is_validated,omitempty"`
}

type InstalledApplicationListResponse []InstalledApplicationListItem

// CommonQueryResponses has a list of query responses common to all device types
type CommonQueryResponses struct {
	UDID                  string            `json:"udid"`
	Languages             []string          `json:"languages,omitempty"`              // ATV 6+
	Locales               []string          `json:"locales,omitempty"`                // ATV 6+
	DeviceID              string            `json:"device_id"`                        // ATV 6+
	OrganizationInfo      map[string]string `json:"organization_info,omitempty"`      // IOS 7+
	LastCloudBackupDate   string            `json:"last_cloud_backup_date,omitempty"` // IOS 8+
	AwaitingConfiguration bool              `json:"awaiting_configuration"`           // IOS 9+
	// iTunes
	ITunesStoreAccountIsActive bool   `json:"itunes_store_account_is_active"` // IOS 7+ OSX 10.9+
	ITunesStoreAccountHash     string `json:"itunes_store_account_hash"`      // IOS 8+ OSX 10.10+

	// Device
	DeviceName                    string  `json:"device_name"`
	OSVersion                     string  `json:"os_version"`
	BuildVersion                  string  `json:"build_version"`
	ModelName                     string  `json:"model_name"`
	Model                         string  `json:"model"`
	ProductName                   string  `json:"product_name"`
	SerialNumber                  string  `json:"serial_number"`
	DeviceCapacity                float32 `json:"device_capacity"`
	AvailableDeviceCapacity       float32 `json:"available_device_capacity"`
	BatteryLevel                  float32 `json:"battery_level,omitempty"`           // IOS 5+
	CellularTechnology            int     `json:"cellular_technology,omitempty"`     // IOS 4+
	IsSupervised                  bool    `json:"is_supervised"`                     // IOS 6+
	IsDeviceLocatorServiceEnabled bool    `json:"is_device_locator_service_enabled"` // IOS 7+
	IsActivationLockEnabled       bool    `json:"is_activation_lock_enabled"`        // IOS 7+ OSX 10.9+
	IsDoNotDisturbInEffect        bool    `json:"is_dnd_in_effect"`                  // IOS 7+
	EASDeviceIdentifier           string  `json:"eas_device_identifier"`             // IOS 7 OSX 10.9
	IsCloudBackupEnabled          bool    `json:"is_cloud_backup_enabled"`           // IOS 7.1

	// Network
	BluetoothMAC string   `json:"bluetooth_mac"`
	WiFiMAC      string   `json:"wifi_mac"`
	EthernetMACs []string `json:"ethernet_macs"` // Surprisingly works in IOS
}

// AtvQueryResponses contains AppleTV QueryResponses
type AtvQueryResponses struct {
}

// IosQueryResponses contains iOS QueryResponses
type IosQueryResponses struct {
	IMEI                 string `json:"imei"`
	MEID                 string `json:"meid"`
	ModemFirmwareVersion string `json:"modem_firmware_version"`
	IsMDMLostModeEnabled bool   `json:"is_mdm_lost_mode_enabled,omitempty"` // IOS 9.3
	MaximumResidentUsers int    `json:"maximum_resident_users"`             // IOS 9.3

	// Network
	ICCID                    string `json:"iccid,omitempty"` // IOS
	CurrentCarrierNetwork    string `json:"current_carrier_network,omitempty"`
	SIMCarrierNetwork        string `json:"sim_carrier_network,omitempty"`
	SubscriberCarrierNetwork string `json:"subscriber_carrier_network,omitempty"`
	CarrierSettingsVersion   string `json:"carrier_settings_version,omitempty"`
	PhoneNumber              string `json:"phone_number,omitempty"`
	VoiceRoamingEnabled      bool   `json:"voice_roaming_enabled,omitempty"`
	DataRoamingEnabled       bool   `json:"data_roaming_enabled,omitempty"`
	IsRoaming                bool   `json:"is_roaming,omitempty"`
	PersonalHotspotEnabled   bool   `json:"personal_hotspot_enabled,omitempty"`
	SubscriberMCC            string `json:"subscriber_mcc,omitempty"`
	SubscriberMNC            string `json:"subscriber_mnc,omitempty"`
	CurrentMCC               string `json:"current_mcc,omitempty"`
	CurrentMNC               string `json:"current_mnc,omitempty"`
}

// OSUpdateSettingsResponse contains information about macOS update settings.
type OSUpdateSettingsResponse struct {
	AutoCheckEnabled                bool      `json:"auto_check_enabled"`
	AutomaticAppInstallationEnabled bool      `json:"automatic_app_installation_enabled"`
	AutomaticOSInstallationEnabled  bool      `json:"automatic_os_installation_enabled"`
	AutomaticSecurityUpdatesEnabled bool      `json:"automatic_security_updates_enabled"`
	BackgroundDownloadEnabled       bool      `json:"background_download_enabled"`
	CatalogURL                      string    `json:"catalog_url"`
	IsDefaultCatalog                bool      `json:"is_default_catalog"`
	PerformPeriodicCheck            bool      `json:"perform_periodic_check"`
	PreviousScanDate                time.Time `json:"previous_scan_date"`
	PreviousScanResult              int       `json:"previous_scan_result"`
}

// MacosQueryResponses contains macOS queryResponses
type MacosQueryResponses struct {
	OSUpdateSettings   OSUpdateSettingsResponse // OSX 10.11+
	LocalHostName      string                   `json:"local_host_name,omitempty"` // OSX 10.11
	HostName           string                   `json:"host_name,omitempty"`       // OSX 10.11
	ActiveManagedUsers []string                 `json:"active_managed_users"`      // OSX 10.11
}

// QueryResponses is a DeviceInformation MDM Command Response
type QueryResponses struct {
	CommonQueryResponses
	MacosQueryResponses
	IosQueryResponses
	AtvQueryResponses
}

// SecurityInfo is the SecurityInfo MDM Command Response
type SecurityInfo struct {
	FDEEnabled                     bool `json:"fde_enabled,omitempty"` // OSX
	FDEHasPersonalRecoveryKey      bool `json:"fde_has_personal_recovery_key,omitempty"`
	FDEHasInstitutionalRecoveryKey bool `json:"fde_has_institutional_recovery_key,omitempty"`

	HardwareEncryptionCaps        int  `json:"hardware_encryption_caps,omitempty"` // iOS
	PasscodeCompliant             bool `json:"passcode_compliant,omitempty"`
	PasscodeCompliantWithProfiles bool `json:"passcode_compliant_with_profiles,omitempty"`
	PasscodePresent               bool `json:"passcode_present,omitempty"`
}

type RequestMirroringResponse struct {
	MirroringResult string `json:"mirroring_result,omitempty"`
}

//type GlobalRestrictions struct {
//	RestrictedBool map[string]bool `plist:"restrictedBool,omitempty" json:"restricted_bool,omitempty"`
//	RestrictedValue map[string]int `plist:"restrictedValue,omitempty" json:"restricted_value,omitempty"`
//	Intersection map[string]string `plist:"intersection,omitempty" json:"intersection,omitempty"` // TODO: not actually string values
//	Union map[string]string `plist:"union,omitempty" json:"union,omitempty"` // TODO: not actually string values
//}

type UsersListItem struct {
	UserName      string `json:"user_name,omitempty"`
	HasDataToSync bool   `json:"has_data_to_sync,omitempty"`
	DataQuota     int    `json:"data_quota,omitempty"`
	DataUsed      int    `json:"data_used,omitempty"`
}

type UsersListResponse []UsersListItem

// Represents a single error in the error chain response
type ErrorChainItem struct {
	ErrorCode            int    `json:"error_code,omitempty"`
	ErrorDomain          string `json:"error_domain,omitempty"`
	LocalizedDescription string `json:"localized_description,omitempty"`
	USEnglishDescription string `json:"us_english_description,omitempty"`
}

type ErrorChain []ErrorChainItem
