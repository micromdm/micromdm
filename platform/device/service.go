package device

import (
	"context"
)

type RemoveDevicesOptions struct {
	UDIDs   []string `json:"udids"`
	Serials []string `json:"serials"`
}

type Service interface {
	ListDevices(ctx context.Context, opt ListDevicesOption) ([]DeviceDTO, error)
	RemoveDevices(ctx context.Context, opt RemoveDevicesOptions) error
}

type Store interface {
	List(ctx context.Context, opt ListDevicesOption) ([]Device, error)
	DeleteByUDID(ctx context.Context, udid string) error
	DeleteBySerial(ctx context.Context, serial string) error
}

type DeviceService struct {
	store Store
}

func New(store Store) *DeviceService {
	return &DeviceService{store: store}
}
