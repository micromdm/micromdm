package remove

import "context"

type Service interface {
	BlockDevice(ctx context.Context, udid string) error
	UnblockDevice(ctx context.Context, udid string) error
}

type Store interface {
	Save(*Device) error
	DeviceByUDID(string) (*Device, error)
	Delete(string) error
}

type RemoveService struct {
	store Store
}

func New(store Store) (*RemoveService, error) {
	return &RemoveService{store: store}, nil
}
