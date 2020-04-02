package queue

import (
	"context"
	
	"github.com/micromdm/micromdm/mdm"
)

const (
	DeviceCommandBucket = "mdm.DeviceCommands"

	CommandQueuedTopic = "mdm.CommandQueued"
)

type Store interface {
	Next(ctx context.Context, resp mdm.Response) ([]byte, error)
	Save(ctx context.Context, cmd *DeviceCommand) error
	DeviceCommand(ctx context.Context, udid string) (*DeviceCommand, error)
}

func New(store Store) *QueueService {
	return &QueueService{store: store}
}

type QueueService struct {
	store Store
}

func IsNotFound(err error) bool {
	type notFoundError interface {
		error
		NotFound() bool
	}

	_, ok := err.(notFoundError)
	return ok
}
