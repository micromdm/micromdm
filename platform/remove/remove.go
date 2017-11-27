package remove

import (
	"context"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/mdm/connect"
	"github.com/micromdm/micromdm/platform/remove/internal/removeproto"
)

type Service interface {
	BlockDevice(ctx context.Context, udid string) error
	UnblockDevice(ctx context.Context, udid string) error
}

type RemoveService struct {
	db *DB
}

func NewService(db *DB) (*RemoveService, error) {
	return &RemoveService{db: db}, nil
}

func (svc *RemoveService) BlockDevice(ctx context.Context, udid string) error {
	return svc.db.Save(&Device{UDID: udid})
}

func (svc *RemoveService) UnblockDevice(ctx context.Context, udid string) error {
	return svc.db.Delete(udid)
}

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

func RemoveMiddleware(db *DB) connect.Middleware {
	return func(next connect.Service) connect.Service {
		return &removeMiddleware{
			db:   db,
			next: next,
		}
	}
}

type removeMiddleware struct {
	db   *DB
	next connect.Service
}

func (mw removeMiddleware) Acknowledge(ctx context.Context, req connect.MDMConnectRequest) ([]byte, error) {
	udid := req.MDMResponse.UDID
	_, err := mw.db.DeviceByUDID(udid)
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

type checkoutErr struct{}

func (checkoutErr) Error() string {
	return "checkout forced by device block"
}

func (checkoutErr) Checkout() bool {
	return true
}
