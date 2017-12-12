package dep

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/dep"
	"github.com/micromdm/micromdm/pkg/httputil"
)

func (svc *DEPService) FetchProfile(ctx context.Context, uuid string) (*dep.Profile, error) {
	if svc.client == nil {
		return nil, errors.New("DEP not configured yet. add a DEP token to enable DEP")
	}
	return svc.client.FetchProfile(uuid)
}

type fetchProfileRequest struct {
	UUID string `json:"uuid"`
}

type fetchProfileResponse struct {
	*dep.Profile
	Err error `json:"err,omitempty"`
}

func (r fetchProfileResponse) Failed() error { return r.Err }

func decodeFetchProfileRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req fetchProfileRequest
	err := httputil.DecodeJSONRequest(r, &req)
	return req, err
}

func decodeFetchProfileResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp fetchProfileResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeFetchProfileEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(fetchProfileRequest)
		profile, err := svc.FetchProfile(ctx, req.UUID)
		return fetchProfileResponse{Profile: profile, Err: err}, nil
	}
}

func (e Endpoints) FetchProfile(ctx context.Context, uuid string) (*dep.Profile, error) {
	request := fetchProfileRequest{UUID: uuid}
	response, err := e.FetchProfileEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(fetchProfileResponse).Profile, response.(fetchProfileResponse).Err
}
