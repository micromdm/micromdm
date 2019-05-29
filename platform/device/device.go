package device

import (
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/platform/device/internal/deviceproto"
)

const DeviceEnrolledTopic = "mdm.DeviceEnrolled"

type Device struct {
	UUID                   string           `db:"uuid" json:"uuid,omitempty"`
	UDID                   string           `db:"udid" json:"udid,omitempty"`
	SerialNumber           string           `db:"serial_number" json:"serial_number,omitempty"`
	OSVersion              string           `db:"os_version" json:"os_version,omitempty"`
	BuildVersion           string           `db:"build_version" json:"build_version,omitempty"`
	ProductName            string           `db:"product_name" json:"product_name,omitempty"`
	IMEI                   string           `db:"imei" json:"imei,omitempty"`
	MEID                   string           `db:"meid" json:"meid,omitempty"`
	PushMagic              string           `db:"push_magic" json:"push_magic,omitempty"`
	AwaitingConfiguration  bool             `db:"awaiting_configuration" json:"awaiting_configuration,omitempty"`
	Token                  string           `db:"token" json:"token,omitempty"`
	UnlockToken            string           `db:"unlock_token" json:"unlock_token,omitempty"`
	Enrolled               bool             `db:"enrolled" json:"enrolled,omitempty"`
	Description            string           `db:"description" json:"description,omitempty"`
	Model                  string           `db:"model" json:"model,omitempty"`
	ModelName              string           `db:"model_name" json:"model_name,omitempty"`
	DeviceName             string           `db:"device_name" json:"device_name,omitempty"`
	Color                  string           `db:"color" json:"color,omitempty"`
	AssetTag               string           `db:"asset_tag" json:"asset_tag,omitempty"`
	DEPProfileStatus       DEPProfileStatus `db:"dep_profile_status" json:"dep_profile_status,omitempty"`
	DEPProfileUUID         string           `db:"dep_profile_uuid" json:"dep_profile_uuid,omitempty"`
	DEPProfileAssignTime   time.Time        `db:"dep_profile_assign_time" json:"dep_profile_assign_time,omitempty"`
	DEPProfilePushTime     time.Time        `db:"dep_profile_push_time" json:"dep_profile_push_time,omitempty"`
	DEPProfileAssignedDate time.Time        `db:"dep_profile_assigned_date" json:"dep_profile_assigned_date,omitempty"`
	DEPProfileAssignedBy   string           `db:"dep_profile_assigned_by" json:"dep_profile_assigned_by,omitempty"`
	LastSeen               time.Time        `db:"last_seen" json:"last_seen,omitempty"`
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

func MarshalDevice(dev *Device) ([]byte, error) {
	protodev := deviceproto.Device{
		Uuid:                   dev.UUID,
		Udid:                   dev.UDID,
		SerialNumber:           dev.SerialNumber,
		OsVersion:              dev.OSVersion,
		BuildVersion:           dev.BuildVersion,
		ProductName:            dev.ProductName,
		Imei:                   dev.IMEI,
		Meid:                   dev.MEID,
		Token:                  dev.Token,
		PushMagic:              dev.PushMagic,
		UnlockToken:            dev.UnlockToken,
		Enrolled:               dev.Enrolled,
		AwaitingConfiguration:  dev.AwaitingConfiguration,
		DeviceName:             dev.DeviceName,
		Model:                  dev.Model,
		ModelName:              dev.ModelName,
		Description:            dev.Description,
		Color:                  dev.Color,
		AssetTag:               dev.AssetTag,
		DepProfileStatus:       string(dev.DEPProfileStatus),
		DepProfileUuid:         dev.DEPProfileUUID,
		DepProfileAssignTime:   timeToNano(dev.DEPProfileAssignTime),
		DepProfilePushTime:     timeToNano(dev.DEPProfilePushTime),
		DepProfileAssignedDate: timeToNano(dev.DEPProfileAssignedDate),
		DepProfileAssignedBy:   dev.DEPProfileAssignedBy,
		LastSeen:               timeToNano(dev.LastSeen),
	}
	return proto.Marshal(&protodev)
}

func UnmarshalDevice(data []byte, dev *Device) error {
	var pb deviceproto.Device
	if err := proto.Unmarshal(data, &pb); err != nil {
		return errors.Wrap(err, "unmarshal proto to device")
	}
	dev.UUID = pb.GetUuid()
	dev.UDID = pb.GetUdid()
	dev.SerialNumber = pb.GetSerialNumber()
	dev.OSVersion = pb.GetOsVersion()
	dev.BuildVersion = pb.GetBuildVersion()
	dev.ProductName = pb.GetProductName()
	dev.IMEI = pb.GetImei()
	dev.MEID = pb.GetMeid()
	dev.Token = pb.GetToken()
	dev.PushMagic = pb.GetPushMagic()
	dev.UnlockToken = pb.GetUnlockToken()
	dev.Enrolled = pb.GetEnrolled()
	dev.AwaitingConfiguration = pb.GetAwaitingConfiguration()
	dev.DeviceName = pb.GetDeviceName()
	dev.Model = pb.GetModel()
	dev.ModelName = pb.GetModelName()
	dev.Description = pb.GetDescription()
	dev.Color = pb.GetColor()
	dev.AssetTag = pb.GetAssetTag()
	dev.DEPProfileStatus = DEPProfileStatus(pb.GetDepProfileStatus())
	dev.DEPProfileUUID = pb.GetDepProfileUuid()
	dev.DEPProfileAssignTime = timeFromNano(pb.GetDepProfileAssignTime())
	dev.DEPProfilePushTime = timeFromNano(pb.GetDepProfilePushTime())
	dev.DEPProfileAssignedDate = timeFromNano(pb.GetDepProfileAssignedDate())
	dev.DEPProfileAssignedBy = pb.GetDepProfileAssignedBy()
	dev.LastSeen = timeFromNano(pb.GetLastSeen())
	return nil
}

func timeToNano(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}
	return t.UnixNano()
}

func timeFromNano(nano int64) time.Time {
	if nano == 0 {
		return time.Time{}
	}
	return time.Unix(0, nano).UTC()
}
