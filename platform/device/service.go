package device

import (
	"context"
)

type Service interface {
	ListDevices(ctx context.Context, opt ListDevicesOption) ([]DeviceDTO, error)
}

type Store interface {
	List(opt ListDevicesOption) ([]Device, error)
}

type DeviceService struct {
	store Store
}

func New(store Store) *DeviceService {
	return &DeviceService{store: store}
}
