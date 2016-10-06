package certificates

type Certificate struct {
	UUID       string `db:"certificate_uuid" json:"uuid"`
	DeviceUUID string `db:"device_uuid" json:"device_uuid"`
	Data       []byte `db:"data" json:"data,omitempty"`
	CommonName string `db:"common_name" json:"common_name,omitempty"`
	IsIdentity bool   `db:"is_identity" json:"is_identity"`
}
