package connect

import (
	"context"
	"fmt"

	"github.com/micromdm/mdm"
	"github.com/micromdm/micromdm/queue"
)

// The ConnectService accepts responses sent to an MDM server by an enrolled
// device.
type ConnectService interface {

	// Acknowledge acknowledges a response sent by a device and returns
	// the next payload if one is available.
	Acknowledge(ctx context.Context, req mdm.Response) (payload []byte, err error)
}

type connectSvc struct {
	queue Queue
}

type Queue interface {
	Next(context.Context, mdm.Response) (*queue.Command, error)
}

func New(queue Queue) (ConnectService, error) {
	return &connectSvc{
		queue: queue,
	}, nil
}

func (svc *connectSvc) Acknowledge(ctx context.Context, req mdm.Response) (payload []byte, err error) {
	fmt.Printf("connected udid=%s type=%s, status=%s\n", req.UDID, req.RequestType, req.Status)
	cmd, err := svc.queue.Next(ctx, req)
	if err != nil {
		return nil, err
	}
	// next can return no errors and no payload.
	if cmd == nil {
		return nil, nil
	}
	return cmd.Payload, nil
}
