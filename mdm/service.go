package mdm

import (
	"context"

	"github.com/micromdm/micromdm/platform/pubsub"
)

type Service interface {
	Checkin(ctx context.Context, event CheckinEvent) error
	Acknowledge(ctx context.Context, event AcknowledgeEvent) (payload []byte, err error)
}

type Middleware func(Service) Service

const (
	ConnectTopic           = "mdm.Connect"
	AuthenticateTopic      = "mdm.Authenticate"
	TokenUpdateTopic       = "mdm.TokenUpdate"
	CheckoutTopic          = "mdm.CheckOut"
	GetBootstrapTokenTopic = "mdm.GetBootstrapToken"
	SetBootstrapTokenTopic = "mdm.SetBootstrapToken"
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
