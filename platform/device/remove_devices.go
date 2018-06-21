package device

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/micromdm/pkg/httputil"
)

func (svc *DeviceService) RemoveDevices(ctx context.Context, udids []string) error {
	for _, udid := range udids {
		err := svc.store.Delete(udid)
		if err != nil {
			return err
		}
	}

	return nil
}

type removeDevicesRequest struct {
	Identifiers []string `json:"udids"`
}

type removeDevicesResponse struct {
	Err error `json:"err,omitempty"`
}

func (r removeDevicesResponse) Failed() error { return r.Err }

func decodeRemoveDevicesRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req removeDevicesRequest
	err := httputil.DecodeJSONRequest(r, &req)
	return req, err
}

func decodeRemoveDevicesResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp removeDevicesResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeRemoveDevicesEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(removeDevicesRequest)
		err = svc.RemoveDevices(ctx, req.Identifiers)
		return removeDevicesResponse{
			Err: err,
		}, nil
	}
}

func (e Endpoints) RemoveDevices(ctx context.Context, ids []string) error {
	request := removeDevicesRequest{Identifiers: ids}
	resp, err := e.RemoveDevicesEndpoint(ctx, request)
	if err != nil {
		return err
	}
	return resp.(removeDevicesResponse).Err
}
