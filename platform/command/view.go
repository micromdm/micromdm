package command

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"github.com/micromdm/micromdm/mdm"
	"github.com/pkg/errors"
)

func (svc *CommandService) ViewQueue(ctx context.Context, udid string) ([]*mdm.Command, error) {
	commands, err := svc.queue.ViewQueue(
		ctx,
		mdm.CheckinEvent{Command: mdm.CheckinCommand{UDID: udid}},
	)
	if err != nil {
		return nil, errors.Wrap(err, "clearing command queue")
	}
	return commands, nil
}

type viewQueueRequest struct {
	UDID string
}

type viewQueueResponse struct {
	Err      error          `json:"error,omitempty"`
	Commands []*mdm.Command `json:"commands,omitempty"`
}

func (r viewQueueResponse) Failed() error   { return r.Err }
func (r viewQueueResponse) StatusCode() int { return http.StatusOK }

func decodeViewQueueRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return viewQueueRequest{UDID: mux.Vars(r)["udid"]}, nil
}

// MakeViewQueueEndpoint creates an endpoint which views device queues.
func MakeViewQueueEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(viewQueueRequest)
		if req.UDID == "" {
			return viewQueueResponse{Err: errEmptyRequest}, nil
		}
		commands, err := svc.ViewQueue(ctx, req.UDID)
		if err != nil {
			return viewQueueResponse{Err: err}, nil
		}

		return viewQueueResponse{Commands: commands}, nil
	}
}
