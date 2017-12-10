package profile

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/micromdm/pkg/httputil"
)

func (svc *ProfileService) ApplyProfile(ctx context.Context, p *Profile) error {
	return svc.store.Save(p)
}

type applyProfileRequest struct {
	Profile *Profile `json:"profile"`
}

type applyProfileResponse struct {
	Err error `json:"err,omitempty"`
}

func (r applyProfileResponse) Failed() error { return r.Err }

func decodeApplyProfileRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req applyProfileRequest
	err := httputil.DecodeJSONRequest(r, &req)
	return req, err
}

func decodeApplyProfileResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp applyProfileResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeApplyProfileEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(applyProfileRequest)
		err = svc.ApplyProfile(ctx, req.Profile)
		return applyProfileResponse{
			Err: err,
		}, nil
	}
}

func (e Endpoints) ApplyProfile(ctx context.Context, p *Profile) error {
	request := applyProfileRequest{Profile: p}
	resp, err := e.ApplyProfileEndpoint(ctx, request)
	if err != nil {
		return err
	}
	return resp.(applyProfileResponse).Err
}
