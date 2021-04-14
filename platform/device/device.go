package device

import (
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/platform/device/internal/deviceproto"
)

const DeviceEnrolledTopic = "mdm.DeviceEnrolled"

type Device struct {
	UUID                   string           `db:"uuid"`
	UDID                   string           `db:"udid"`
	SerialNumber           string           `db:"serial_number"`
	OSVersion              string           `db:"os_version"`
	BuildVersion           string           `db:"build_version"`
	ProductName            string           `db:"product_name"`
	IMEI                   string           `db:"imei"`
	MEID                   string           `db:"meid"`
	PushMagic              string           `db:"push_magic"`
	AwaitingConfiguration  bool             `db:"awaiting_configuration"`
	Token                  string           `db:"token"`
	UnlockToken            string           `db:"unlock_token"`
	Enrolled               bool             `db:"enrolled"`
	Description            string           `db:"description"`
	Model                  string           `db:"model"`
	ModelName              string           `db:"model_name"`
	DeviceName             string           `db:"device_name"`
	Color                  string           `db:"color"`
	AssetTag               string           `db:"asset_tag"`
	DEPProfileStatus       DEPProfileStatus `db:"dep_profile_status"`
	DEPProfileUUID         string           `db:"dep_profile_uuid"`
	DEPProfileAssignTime   time.Time        `db:"dep_profile_assign_time"`
	DEPProfilePushTime     time.Time        `db:"dep_profile_push_time"`
	DEPProfileAssignedDate time.Time        `db:"dep_profile_assigned_date"`
	DEPProfileAssignedBy   string           `db:"dep_profile_assigned_by"`
	LastSeen               time.Time        `db:"last_seen"`
	BootstrapToken         string           `db:"bootstrap_token"`
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
		BootstrapToken:         dev.BootstrapToken,
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
	dev.BootstrapToken = pb.GetBootstrapToken()
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
