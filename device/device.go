package device

import (
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/micromdm/micromdm/device/internal/deviceproto"
	"github.com/pkg/errors"
)

type Device struct {
	UUID                   string
	UDID                   string
	SerialNumber           string
	OSVersion              string
	BuildVersion           string
	ProductName            string
	IMEI                   string
	MEID                   string
	MDMTopic               string
	PushMagic              string
	AwaitingConfiguration  bool
	Token                  string
	UnlockToken            string
	Enrolled               bool
	DEPDevice              bool
	Description            string
	Model                  string
	ModelName              string
	DeviceName             string
	Color                  string
	AssetTag               string
	DEPProfileStatus       DEPProfileStatus
	DEPProfileUUID         string
	DEPProfileAssignTime   time.Time
	DEPProfilePushTime     time.Time
	DEPProfileAssignedDate time.Time
	DEPProfileAssignedBy   string
	LastCheckin            time.Time
	LastQueryResponse      []byte
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
		MdmTopic:               dev.MDMTopic,
		UnlockToken:            dev.UnlockToken,
		Enrolled:               dev.Enrolled,
		AwaitingConfiguration:  dev.AwaitingConfiguration,
		DeviceName:             dev.DeviceName,
		Model:                  dev.Model,
		ModelName:              dev.ModelName,
		Description:            dev.Description,
		Color:                  dev.Color,
		AssetTag:               dev.AssetTag,
		DepDevice:              dev.DEPDevice,
		DepProfileStatus:       dev.DEPProfileStatus,
		DepProfileUuid:         dev.DEPProfileUUID,
		DepAssignTime:          dev.DEPAssignTime,
		DepPushTime:            dev.DEPPushTime,
		DepProfileAssignedDate: dev.DEPProfileAssignedDate,
		DEPProfileAssignedBy:   dev.DEPProfileAssignedBy,
		LastCheckin:            dev.LastCheckin,
		LastQueryResponse:      dev.LastQueryResponse,
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
	dev.MDMTopic = pb.GetMdmTopic()
	dev.UnlockToken = pb.GetUnlockToken()
	dev.Enrolled = pb.GetEnrolled()
	dev.AwaitingConfiguration = pb.GetAwaitingConfiguration()
	dev.DeviceName = pb.GetDeviceName()
	dev.Model = pb.GetModel()
	dev.ModelName = pb.GetModelName()
	dev.Description = pb.GetDescription()
	dev.Color = pb.GetColor()
	dev.AssetTag = pb.GetAssetTag()
	dev.DEPDevice = pb.GetDepDevice()
	dev.DEPProfileStatus = pb.GetDepProfileStatus()
	dev.DEPProfileUUID = pb.GetDepProfileUuid()
	dev.DEPAssignTime = pb.GetDepAssignTime()
	dev.DEPPushTime = pb.GetDepPushTime()
	dev.DEPProfileAssignedDate = pb.GetDepProfileAssignedDate()
	dev.DEPProfileAssignedBy = pb.GetDepProfileAssignedBy()
	dev.LastCheckin = pb.GetLastCheckin()
	dev.LastQueryResponse = pb.GetLastQueryResponse()
	return nil
}
