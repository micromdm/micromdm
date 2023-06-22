package command

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"github.com/micromdm/micromdm/mdm"
	"github.com/pkg/errors"
)

func (svc *CommandService) ViewQueue(ctx context.Context, udid string) ([]byte, error) {
	commands, err := svc.queue.View(ctx, mdm.CheckinEvent{Command: mdm.CheckinCommand{UDID: udid}})
	if err != nil {
		return []byte{}, errors.Wrap(err, "clearing command queue")
	}
	return commands, nil
}

type viewQueueRequest struct {
	UDID string
}

type viewQueueResponse struct {
	Err           error       `json:"error,omitempty"`
	DeviceCommand interface{} `json:"device_command,omitempty"`
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
			return clearResponse{Err: err}, nil
		}

		var out interface{}
		err = json.Unmarshal(commands, &out)
		if err != nil {
			return viewQueueResponse{Err: err}, nil
		}

		resp := viewQueueResponse{
			DeviceCommand: out,
			Err:           nil,
		}
		return resp, nil
	}
}
