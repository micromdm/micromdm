package profile

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/micromdm/pkg/httputil"
)

func (svc *ProfileService) RemoveProfiles(ctx context.Context, ids []string) error {
	for _, id := range ids {
		err := svc.store.Delete(id)
		if err != nil {
			return err
		}
	}
	return nil
}

type removeProfileRequest struct {
	Identifiers []string `json:"ids"`
}

type removeProfileResponse struct {
	Err error `json:"err,omitempty"`
}

func (r removeProfileResponse) Failed() error { return r.Err }

func decodeRemoveProfilesRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req removeProfileRequest
	err := httputil.DecodeJSONRequest(r, &req)
	return req, err
}

func decodeRemoveProfileResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp removeProfileResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeRemoveProfilesEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(removeProfileRequest)
		err = svc.RemoveProfiles(ctx, req.Identifiers)
		return removeProfileResponse{
			Err: err,
		}, nil
	}
}

func (e Endpoints) RemoveProfiles(ctx context.Context, ids []string) error {
	request := removeProfileRequest{Identifiers: ids}
	resp, err := e.RemoveProfilesEndpoint(ctx, request)
	if err != nil {
		return err
	}
	return resp.(removeProfileResponse).Err
}
