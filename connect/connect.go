package connect

import (
	"context"

	"github.com/micromdm/mdm"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/pubsub"
	"github.com/micromdm/micromdm/queue"
)

const ConnectTopic = "mdm.Connect"

// The ConnectService accepts responses sent to an MDM server by an enrolled
// device.
type ConnectService interface {

	// Acknowledge acknowledges a response sent by a device and returns
	// the next payload if one is available.
	Acknowledge(ctx context.Context, req MDMConnectRequest) (payload []byte, err error)
}

type connectSvc struct {
	queue Queue
	pub   pubsub.Publisher
}

type Queue interface {
	Next(context.Context, mdm.Response) (*queue.Command, error)
}

func New(queue Queue, pub pubsub.Publisher) (ConnectService, error) {
	return &connectSvc{
		queue: queue,
		pub:   pub,
	}, nil
}

func (svc *connectSvc) Acknowledge(ctx context.Context, req MDMConnectRequest) (payload []byte, err error) {
	event := NewEvent(req)
	msg, err := MarshalEvent(event)
	if err != nil {
		return nil, errors.Wrap(err, "marshal connect response to proto")
	}
	if err := svc.pub.Publish(context.TODO(), ConnectTopic, msg); err != nil {
		return nil, errors.Wrap(err, "publish connect Response on pubsub")
	}

	cmd, err := svc.queue.Next(ctx, req.MDMResponse)
	if err != nil {
		return nil, errors.Wrap(err, "calling Next with mdm response")
	}
	// next can return no errors and no payload.
	if cmd == nil {
		return nil, nil
	}
	return cmd.Payload, nil
}
