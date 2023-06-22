// Package command provides utilities for creating MDM Payloads.
package command

import (
	mdmsvc "github.com/micromdm/micromdm/mdm"
	"github.com/micromdm/micromdm/mdm/mdm"
	"github.com/micromdm/micromdm/platform/pubsub"
	"golang.org/x/net/context"
)

type Service interface {
	NewCommand(context.Context, *mdm.CommandRequest) (*mdm.CommandPayload, error)
	NewRawCommand(context.Context, *RawCommand) error
	ClearQueue(ctx context.Context, udid string) error
	ViewQueue(ctx context.Context, udid string) ([]byte, error)
}

// Queue is an MDM Command Queue.
type Queue interface {
	Clear(context.Context, mdmsvc.CheckinEvent) error
	ViewQueue(context.Context, mdmsvc.CheckinEvent) ([]byte, error)
}

type CommandService struct {
	publisher pubsub.Publisher
	queue     Queue
}

func New(pub pubsub.Publisher, queue Queue) (*CommandService, error) {
	svc := CommandService{
		publisher: pub,
		queue:     queue,
	}
	return &svc, nil
}
