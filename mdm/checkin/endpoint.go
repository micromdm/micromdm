package checkin

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/mdm"
)

// errInvalidMessageType is an invalid checking command.
var errInvalidMessageType = errors.New("invalid message type")

type Endpoints struct {
	CheckinEndpoint endpoint.Endpoint
}

func MakeCheckinEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(checkinRequest)
		var err error
		switch req.MessageType {
		case "Authenticate":
			err = svc.Authenticate(ctx, req.CheckinCommand, req.id)
		case "TokenUpdate":
			err = svc.TokenUpdate(ctx, req.CheckinCommand, req.id)
		case "CheckOut":
			err = svc.CheckOut(ctx, req.CheckinCommand, req.id)
		default:
			return checkinResponse{Err: errInvalidMessageType}, nil
		}
		if err != nil {
			return checkinResponse{Err: err}, nil
		}
		return checkinResponse{}, nil
	}
}

type checkinRequest struct {
	mdm.CheckinCommand
	id string
}

type checkinResponse struct {
	Err error `plist:"error,omitempty"`
}

func (r checkinResponse) error() error { return r.Err }
