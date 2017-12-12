package dep

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/dep"
	"github.com/micromdm/micromdm/pkg/httputil"
)

func (svc *DEPService) GetDeviceDetails(ctx context.Context, serials []string) (*dep.DeviceDetailsResponse, error) {
	if svc.client == nil {
		return nil, errors.New("DEP not configured yet. add a DEP token to enable DEP")
	}
	return svc.client.DeviceDetails(serials)
}

type deviceDetailsRequest struct {
	Serials []string `json:"serials"`
}

type deviceDetailsResponse struct {
	*dep.DeviceDetailsResponse
	Err error `json:"err,omitempty"`
}

func (r deviceDetailsResponse) Failed() error { return r.Err }

func decodeDeviceDetailsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req deviceDetailsRequest
	err := httputil.DecodeJSONRequest(r, &req)
	return req, err
}

func decodeDeviceDetailsResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp deviceDetailsResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeGetDeviceDetailsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deviceDetailsRequest)
		details, err := svc.GetDeviceDetails(ctx, req.Serials)
		return deviceDetailsResponse{DeviceDetailsResponse: details, Err: err}, nil
	}
}

func (e Endpoints) GetDeviceDetails(ctx context.Context, serials []string) (*dep.DeviceDetailsResponse, error) {
	request := deviceDetailsRequest{Serials: serials}
	response, err := e.GetDeviceDetailsEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(deviceDetailsResponse).DeviceDetailsResponse, response.(deviceDetailsResponse).Err
}
