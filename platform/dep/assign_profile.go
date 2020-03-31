package dep

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"

	"github.com/micromdm/micromdm/dep"
	"github.com/micromdm/micromdm/pkg/httputil"
)

func (svc *DEPService) AssignProfile(ctx context.Context, id string, serials ...string) (*dep.ProfileResponse, error) {
	if svc.client == nil {
		return nil, errors.New("DEP not configured yet. add a DEP token to enable DEP")
	}
	return svc.client.AssignProfile(id, serials...)
}

type assignProfileRequest struct {
	ID      string   `json:"id"`
	Serials []string `json:"serials"`
}

type assignProfileResponse struct {
	*dep.ProfileResponse
	Err error `json:"err,omitempty"`
}

func (r assignProfileResponse) Failed() error { return r.Err }

func decodeAssignProfileRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req assignProfileRequest
	err := httputil.DecodeJSONRequest(r, &req)
	return req, err
}

func decodeAssignProfileResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp assignProfileResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeAssignProfileEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(assignProfileRequest)
		resp, err := svc.AssignProfile(ctx, req.ID, req.Serials...)
		return &assignProfileResponse{
			ProfileResponse: resp,
			Err:             err,
		}, nil
	}
}

func (e Endpoints) AssignProfile(ctx context.Context, id string, serials ...string) (*dep.ProfileResponse, error) {
	request := assignProfileRequest{ID: id, Serials: serials}
	resp, err := e.AssignProfileEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	response := resp.(assignProfileResponse)
	return response.ProfileResponse, response.Err
}
