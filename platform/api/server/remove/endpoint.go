package remove

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	UnblockDeviceEndpoint endpoint.Endpoint
}

func MakeEndpoints(svc Service) Endpoints {
	e := Endpoints{
		UnblockDeviceEndpoint: MakeUnblockDeviceEndpoint(svc),
	}
	return e
}

func (e Endpoints) UnblockDevice(ctx context.Context, udid string) error {
	request := unblockDeviceRequest{UDID: udid}
	resp, err := e.UnblockDeviceEndpoint(ctx, request)
	if err != nil {
		return err
	}
	return resp.(unblockDeviceResponse).Err
}

func MakeUnblockDeviceEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(unblockDeviceRequest)
		err = svc.UnblockDevice(ctx, req.UDID)
		return unblockDeviceResponse{
			Err: err,
		}, nil
	}
}

type unblockDeviceRequest struct {
	UDID string
}

type unblockDeviceResponse struct {
	Err error `json:"err,omitempty"`
}

func (r unblockDeviceResponse) error() error { return r.Err }

type blueprintRequest struct {
	Names []string `json:"names"`
}

type blueprintResponse struct {
	Err error `json:"err,omitempty"`
}

func (r blueprintResponse) error() error { return r.Err }

type profileRequest struct {
	Identifiers []string `json:"ids"`
}

type profileResponse struct {
	Err error `json:"err,omitempty"`
}

func (r profileResponse) error() error { return r.Err }
