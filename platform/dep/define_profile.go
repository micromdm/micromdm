package dep

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/dep"
	"github.com/micromdm/micromdm/pkg/httputil"
)

func (svc *DEPService) DefineProfile(ctx context.Context, p *dep.Profile) (*dep.ProfileResponse, error) {
	if svc.client == nil {
		return nil, errors.New("DEP not configured yet. add a DEP token to enable DEP")
	}
	return svc.client.DefineProfile(p)
}

type defineProfileRequest struct{ *dep.Profile }
type defineProfileResponse struct {
	*dep.ProfileResponse
	Err error `json:"err,omitempty"`
}

func (r defineProfileResponse) Failed() error { return r.Err }

func decodeDefineProfileRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req defineProfileRequest
	err := httputil.DecodeJSONRequest(r, &req)
	return req, err
}

func decodeDefineProfileResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp defineProfileResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeDefineProfileEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(defineProfileRequest)
		resp, err := svc.DefineProfile(ctx, req.Profile)
		return &defineProfileResponse{
			ProfileResponse: resp,
			Err:             err,
		}, nil
	}
}

func (e Endpoints) DefineProfile(ctx context.Context, p *dep.Profile) (*dep.ProfileResponse, error) {
	request := defineProfileRequest{Profile: p}
	resp, err := e.DefineProfileEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	response := resp.(defineProfileResponse)
	return response.ProfileResponse, response.Err
}
