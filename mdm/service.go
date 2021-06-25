package mdm

import (
	"context"

	"github.com/micromdm/micromdm/platform/pubsub"
)

// Service describes the core MDM protocol interactions with clients.
type Service interface {
	// Checkin is called for all checkin messages (such as Authenticate
	// and TokenUpdate). Checkin messages return a response to the
	// client in payload.
	Checkin(ctx context.Context, event CheckinEvent) (payload []byte, err error)
	// Acknowledge is called when a client reports a command result and
	// fetches the next command. Acknowledge calls return a command in
	// payload.
	Acknowledge(ctx context.Context, event AcknowledgeEvent) (payload []byte, err error)
}

type Middleware func(Service) Service

const (
	ConnectTopic      = "mdm.Connect"
	AuthenticateTopic = "mdm.Authenticate"
	TokenUpdateTopic  = "mdm.TokenUpdate"
	CheckoutTopic     = "mdm.CheckOut"
)

// Queue is an MDM Command Queue.
type Queue interface {
	Next(context.Context, Response) ([]byte, error)
	Clear(context.Context, CheckinEvent) error
}

type MDMService struct {
	pub   pubsub.Publisher
	queue Queue
}

func NewService(pub pubsub.Publisher, queue Queue) *MDMService {
	return &MDMService{
		pub:   pub,
		queue: queue,
	}
}
