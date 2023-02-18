// Package command provides utilities for creating MDM Payloads.
package command

import (
	"github.com/micromdm/micromdm/mdm/mdm"
	"github.com/micromdm/micromdm/platform/pubsub"
	"golang.org/x/net/context"
)

type Service interface {
	NewCommand(context.Context, *mdm.CommandRequest) (*mdm.CommandPayload, error)
	NewRawCommand(context.Context, *RawCommand) error
}

type CommandService struct {
	publisher pubsub.Publisher
}

func New(pub pubsub.Publisher) (*CommandService, error) {
	svc := CommandService{
		publisher: pub,
	}
	return &svc, nil
}
