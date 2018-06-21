package device

import (
	"context"
)

type Service interface {
	ListDevices(ctx context.Context, opt ListDevicesOption) ([]DeviceDTO, error)
	RemoveDevices(ctx context.Context, udids []string) error
}

type Store interface {
	List(opt ListDevicesOption) ([]Device, error)
	Delete(udid string) error
}

type DeviceService struct {
	store Store
}

func New(store Store) *DeviceService {
	return &DeviceService{store: store}
}
