package connect

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/mdm"
)

type MDMConnectRequest struct {
	Raw         []byte
	MDMResponse mdm.Response
}

type mdmConnectResponse struct {
	payload []byte
	Err     error `plist:"error,omitempty"`
}

func (r mdmConnectResponse) error() error { return r.Err }

type Endpoints struct {
	ConnectEndpoint endpoint.Endpoint
}

func MakeConnectEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(MDMConnectRequest)
		payload, err := svc.Acknowledge(ctx, req)
		return mdmConnectResponse{payload: payload, Err: err}, nil
	}
}
