package remove

import (
	"context"

	"github.com/go-kit/kit/log"
)

type Service interface {
	BlockDevice(ctx context.Context, udid string) error
	UnblockDevice(ctx context.Context, udid string) error
}

type Store interface {
	Save(context.Context, *Device) error
	DeviceByUDID(context.Context, string) (*Device, error)
	Delete(context.Context, string) error
}

type RemoveService struct {
	store Store
}

func New(store Store) (*RemoveService, error) {
	return &RemoveService{store: store}, nil
}

type Middleware func(next Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return logmw{logger: logger, next: next}
	}
}

type logmw struct {
	logger log.Logger
	next   Service
}