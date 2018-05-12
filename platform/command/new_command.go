package command

import (
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/micromdm/micromdm/mdm/mdm"
	"github.com/micromdm/micromdm/pkg/httputil"
)

const (
	// CommandTopic is a PubSub topic that events are published to.
	CommandTopic = "mdm.Command"
)

func (svc *CommandService) NewCommand(ctx context.Context, request *mdm.CommandRequest) (*mdm.CommandPayload, error) {
	if request == nil {
		return nil, errors.New("empty CommandRequest")
	}
	payload, err := mdm.NewCommandPayload(request)
	if err != nil {
		return nil, errors.Wrap(err, "creating mdm payload")
	}
	event := NewEvent(payload, request.UDID)
	msg, err := MarshalEvent(event)
	if err != nil {
		return nil, errors.Wrap(err, "marshalling mdm command event")
	}
	if err := svc.publisher.Publish(context.TODO(), CommandTopic, msg); err != nil {
		return nil, errors.Wrapf(err, "publish mdm command on topic: %s", CommandTopic)
	}
	return payload, nil
}

type newCommandRequest struct {
	mdm.CommandRequest
}

type newCommandResponse struct {
	Payload *mdm.CommandPayload `json:"payload,omitempty"`
	Err     error               `json:"error,omitempty"`
}

func (r newCommandResponse) Failed() error   { return r.Err }
func (r newCommandResponse) StatusCode() int { return http.StatusCreated }

func decodeNewCommandRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req newCommandRequest
	err := httputil.DecodeJSONRequest(r, &req)
	return req, err
}

var errEmptyRequest = errors.New("request must contain UDID of the device")

// MakeNewCommandEndpoint creates an endpoint which creates new MDM Commands.
func MakeNewCommandEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(newCommandRequest)
		if req.UDID == "" || req.RequestType == "" {
			return newCommandResponse{Err: errEmptyRequest}, nil
		}
		payload, err := svc.NewCommand(ctx, &req.CommandRequest)
		if err != nil {
			return newCommandResponse{Err: err}, nil
		}
		return newCommandResponse{Payload: payload}, nil
	}
}
