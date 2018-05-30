package device

import (
	"crypto/x509"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/platform/device/internal/deviceproto"
)

const DeviceEnrolledTopic = "mdm.DeviceEnrolled"

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
	LastSeen               time.Time
	LastQueryResponse      []byte
	DeviceCert             *x509.Certificate
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
		DepProfileStatus:       string(dev.DEPProfileStatus),
		DepProfileUuid:         dev.DEPProfileUUID,
		DepProfileAssignTime:   timeToNano(dev.DEPProfileAssignTime),
		DepProfilePushTime:     timeToNano(dev.DEPProfilePushTime),
		DepProfileAssignedDate: timeToNano(dev.DEPProfileAssignedDate),
		DepProfileAssignedBy:   dev.DEPProfileAssignedBy,
		LastSeen:               timeToNano(dev.LastSeen),
		LastQueryResponse:      dev.LastQueryResponse,
	}
	if dev.DeviceCert != nil {
		protodev.DeviceCert = dev.DeviceCert.Raw
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
	dev.DEPProfileStatus = DEPProfileStatus(pb.GetDepProfileStatus())
	dev.DEPProfileUUID = pb.GetDepProfileUuid()
	dev.DEPProfileAssignTime = timeFromNano(pb.GetDepProfileAssignTime())
	dev.DEPProfilePushTime = timeFromNano(pb.GetDepProfilePushTime())
	dev.DEPProfileAssignedDate = timeFromNano(pb.GetDepProfileAssignedDate())
	dev.DEPProfileAssignedBy = pb.GetDepProfileAssignedBy()
	dev.LastSeen = timeFromNano(pb.GetLastSeen())
	dev.LastQueryResponse = pb.GetLastQueryResponse()
	certs, err := x509.ParseCertificates(pb.GetDeviceCert())
	if err != nil {
		return err
	}
	if len(certs) > 0 {
		dev.DeviceCert = certs[0]
	}
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
