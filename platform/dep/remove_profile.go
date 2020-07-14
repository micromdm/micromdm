package dep

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"

	// "github.com/micromdm/micromdm/dep"
	"github.com/micromdm/micromdm/pkg/httputil"
)

func (svc *DEPService) RemoveProfile(ctx context.Context, serials ...string) (map[string]string, error) {
	if svc.client == nil {
		return nil, errors.New("DEP not configured yet. add a DEP token to enable DEP")
	}
	return svc.client.RemoveProfile(serials...)
}

type removeProfileRequest struct {
	Serials []string `json:"serials"`
}

type removeProfileResponse struct {
	Serials map[string]string `json:"serials"`
	Err     error             `json:"err,omitempty"`
}

func (r removeProfileResponse) Failed() error { return r.Err }

func decodeRemoveProfileRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req removeProfileRequest
	err := httputil.DecodeJSONRequest(r, &req)
	return req, err
}

func decodeRemoveProfileResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp removeProfileResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeRemoveProfileEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(removeProfileRequest)
		resp, err := svc.RemoveProfile(ctx, req.Serials...)
		return &removeProfileResponse{
			Serials: resp,
			Err:     err,
		}, nil
	}
}

func (e Endpoints) RemoveProfile(ctx context.Context, serials ...string) (map[string]string, error) {
	request := removeProfileRequest{Serials: serials}
	resp, err := e.RemoveProfileEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	response := resp.(removeProfileResponse)
	return response.Serials, response.Err
}
