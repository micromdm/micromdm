package profile

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

func (svc *ProfileService) GetProfiles(ctx context.Context, opt GetProfilesOption) ([]Profile, error) {
	if opt.Identifier != "" {
		foundProf, err := svc.store.ProfileById(opt.Identifier)
		if err != nil {
			return nil, err
		}
		return []Profile{*foundProf}, nil
	} else {
		return svc.store.List()
	}
}

type getProfilesRequest struct{ Opts GetProfilesOption }

type getProfilesResponse struct {
	Profiles []Profile `json:"profiles"`
	Err      error     `json:"err,omitempty"`
}

func (r getProfilesResponse) error() error { return r.Err }

func decodeGetProfilesRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var opts GetProfilesOption
	if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
		return nil, err
	}
	req := getProfilesRequest{
		Opts: opts,
	}
	return req, nil
}

func decodeGetProfilesResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errorDecoder(r)
	}
	var resp getProfilesResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func MakeGetProfilesEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getProfilesRequest)
		profiles, err := svc.GetProfiles(ctx, req.Opts)
		return getProfilesResponse{
			Profiles: profiles,
			Err:      err,
		}, nil
	}
}

func (e Endpoints) GetProfiles(ctx context.Context, opt GetProfilesOption) ([]Profile, error) {
	request := getProfilesRequest{opt}
	response, err := e.GetProfilesEndpoint(ctx, request.Opts)
	if err != nil {
		return nil, err
	}
	return response.(getProfilesResponse).Profiles, response.(getProfilesResponse).Err
}
