package remove

import (
	"context"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/mdm"
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

func RemoveMiddleware(store Store) mdm.Middleware {
	return func(next mdm.Service) mdm.Service {
		return &removeMiddleware{
			store: store,
			next:  next,
		}
	}
}

type removeMiddleware struct {
	store Store
	next  mdm.Service
}

func (mw removeMiddleware) Acknowledge(ctx context.Context, req mdm.AcknowledgeEvent) ([]byte, error) {
	udid := req.Response.UDID
	_, err := mw.store.DeviceByUDID(udid)
	if err != nil {
		if !isNotFound(err) {
			return nil, errors.Wrapf(err, "remove: get device by udid %s", udid)
		}
	}
	if err == nil {
		return nil, checkoutErr{}
	}
	return mw.next.Acknowledge(ctx, req)
}

func (mw removeMiddleware) Checkin(ctx context.Context, req mdm.CheckinEvent) ([]byte, error) {
	return mw.next.Checkin(ctx, req)
}

type checkoutErr struct{}

func (checkoutErr) Error() string {
	return "checkout forced by device block"
}

func (checkoutErr) Checkout() bool {
	return true
}

func isNotFound(err error) bool {
	type notFoundError interface {
		error
		NotFound() bool
	}

	_, ok := errors.Cause(err).(notFoundError)
	return ok
}
