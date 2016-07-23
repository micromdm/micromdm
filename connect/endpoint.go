package connect

import (
	"errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/mdm"
	"golang.org/x/net/context"
)

// errInvalidMessageType is an invalid checking command
var errInvalidMessageType = errors.New("Invalid MessageType")

type mdmConnectRequest struct {
	mdm.Response
}

type mdmConnectResponse struct {
	payload []byte
	Err     error `plist:"error,omitempty"`
}

func (r mdmConnectResponse) error() error { return r.Err }

func makeConnectEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(mdmConnectRequest)
		if req.UserID != nil {
			// don't handle user
			return mdmConnectResponse{}, nil
		}
		var err error
		switch req.Status {
		case "Acknowledged":
			total, err := svc.Acknowledge(ctx, req.Response)
			if err != nil {
				return mdmConnectResponse{Err: err}, nil
			}
			if req.RequestType == "DeviceConfigured" {
				return mdmConnectResponse{}, nil
			}
			if total != 0 {
				next, _, err := svc.NextCommand(ctx, req.Response)
				if err != nil {
					return mdmConnectResponse{Err: err}, nil
				}
				return mdmConnectResponse{payload: next}, nil
			}
		case "Idle":
			next, total, err := svc.NextCommand(ctx, req.Response)
			if err != nil {
				return mdmConnectResponse{Err: err}, nil
			}
			if total == 0 {
				return mdmConnectResponse{}, nil
			}
			return mdmConnectResponse{payload: next}, nil
		case "Error":
			total, err := svc.FailCommand(ctx, req.Response)
			if err != nil {
				return mdmConnectResponse{Err: err}, nil
			}
			// TODO: Deal with command failures
			if total != 0 {
				next, _, err := svc.NextCommand(ctx, req.Response)
				if err != nil {
					return mdmConnectResponse{Err: err}, nil
				}
				return mdmConnectResponse{payload: next}, nil
			}

		default:
			return mdmConnectResponse{Err: errInvalidMessageType}, nil
		}
		if err != nil {
			return mdmConnectResponse{Err: err}, nil
		}
		return mdmConnectResponse{}, nil
	}
}
