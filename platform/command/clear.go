package command

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"github.com/micromdm/micromdm/mdm"
	"github.com/pkg/errors"
)

func (svc *CommandService) ClearQueue(ctx context.Context, udid string) error {
	if err := svc.queue.Clear(context.TODO(), mdm.CheckinEvent{Command: mdm.CheckinCommand{UDID: udid}}); err != nil {
		return errors.Wrap(err, "clearing command queue")
	}
	return nil
}

type clearRequest struct {
	UDID string
}

type clearResponse struct {
	Err error `json:"error,omitempty"`
}

func (r clearResponse) Failed() error   { return r.Err }
func (r clearResponse) StatusCode() int { return http.StatusOK }

func decodeClearRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return clearRequest{UDID: mux.Vars(r)["udid"]}, nil
}

// MakeClearQueueEndpoint creates an endpoint which clears device queues.
func MakeClearQueueEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(clearRequest)
		if req.UDID == "" {
			return clearResponse{Err: errEmptyRequest}, nil
		}
		if err := svc.ClearQueue(ctx, req.UDID); err != nil {
			return clearResponse{Err: err}, nil
		}
		return clearResponse{}, nil
	}
}
