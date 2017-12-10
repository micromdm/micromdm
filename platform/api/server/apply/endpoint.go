package apply

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/dep"
)

type Endpoints struct {
	DefineDEPProfileEndpoint endpoint.Endpoint
}

func (e Endpoints) DefineDEPProfile(ctx context.Context, p *dep.Profile) (*dep.ProfileResponse, error) {
	request := depProfileRequest{Profile: p}
	resp, err := e.DefineDEPProfileEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	response := resp.(depProfileResponse)
	return response.ProfileResponse, response.Err
}

func MakeDefineDEPProfile(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(depProfileRequest)
		resp, err := svc.DefineDEPProfile(ctx, req.Profile)
		return &depProfileResponse{
			ProfileResponse: resp,
			Err:             err,
		}, nil
	}
}

type depProfileRequest struct{ *dep.Profile }
type depProfileResponse struct {
	*dep.ProfileResponse
	Err error `json:"err,omitempty"`
}

func (r *depProfileResponse) error() error { return r.Err }
