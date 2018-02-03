package checkin

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/mdm"
	"github.com/pkg/errors"
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
			err = svc.Authenticate(ctx, req.CheckinCommand)
		case "TokenUpdate":
			err = svc.TokenUpdate(ctx, req.CheckinCommand)
		case "CheckOut":
			err = svc.CheckOut(ctx, req.CheckinCommand)
		case "UserAuthenticate":
			err = &rejectUserAuth{}
		default:
			err = errInvalidMessageType
		}
		if err != nil {
			return checkinResponse{Err: err}, nil
		}
		return checkinResponse{}, nil
	}
}

type checkinRequest struct {
	mdm.CheckinCommand
}

type checkinResponse struct {
	Err error `plist:"error,omitempty"`
}

func (r checkinResponse) error() error { return r.Err }

type rejectUserAuth struct{}

func (e *rejectUserAuth) Error() string {
	return "reject user auth"
}

func (e *rejectUserAuth) UserAuthReject() bool {
	return true
}

func (e *rejectUserAuth) error() error { return e }

func isRejectedUserAuth(err error) bool {
	type rejectUserAuthError interface {
		error
		UserAuthReject() bool
	}

	_, ok := errors.Cause(err).(rejectUserAuthError)
	return ok
}
