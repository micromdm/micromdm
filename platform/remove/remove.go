package remove

import (
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/platform/remove/internal/removeproto"
)

type Device struct {
	UDID string `json:"udid"`
}

func MarshalDevice(dev *Device) ([]byte, error) {
	protodev := removeproto.Device{
		Udid: dev.UDID,
	}
	return proto.Marshal(&protodev)
}

func UnmarshalDevice(data []byte, dev *Device) error {
	var pb removeproto.Device
	if err := proto.Unmarshal(data, &pb); err != nil {
		return errors.Wrap(err, "remove: unmarshal proto to device")
	}
	dev.UDID = pb.GetUdid()
	return nil
}
