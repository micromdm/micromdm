package remove

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	RemoveBlueprintsEndpoint endpoint.Endpoint
	RemoveProfilesEndpoint   endpoint.Endpoint
}

func MakeEndpoints(svc Service) Endpoints {
	e := Endpoints{
		RemoveBlueprintsEndpoint: MakeRemoveBlueprintsEndpoint(svc),
		RemoveProfilesEndpoint:   MakeRemoveProfilesEndpoint(svc),
	}
	return e
}

func (e Endpoints) RemoveBlueprints(ctx context.Context, names []string) error {
	request := blueprintRequest{Names: names}
	resp, err := e.RemoveBlueprintsEndpoint(ctx, request)
	if err != nil {
		return err
	}
	return resp.(blueprintResponse).Err
}

func (e Endpoints) RemoveProfiles(ctx context.Context, ids []string) error {
	request := profileRequest{Identifiers: ids}
	resp, err := e.RemoveProfilesEndpoint(ctx, request)
	if err != nil {
		return err
	}
	return resp.(profileResponse).Err
}

func MakeRemoveBlueprintsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(blueprintRequest)
		err = svc.RemoveBlueprints(ctx, req.Names)
		return blueprintResponse{
			Err: err,
		}, nil
	}
}

func MakeRemoveProfilesEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(profileRequest)
		err = svc.RemoveProfiles(ctx, req.Identifiers)
		return profileResponse{
			Err: err,
		}, nil
	}
}

type blueprintRequest struct {
	Names []string `json:"names"`
}

type blueprintResponse struct {
	Err error `json:"err,omitempty"`
}

func (r blueprintResponse) error() error { return r.Err }

type profileRequest struct {
	Identifiers []string `json:"ids"`
}

type profileResponse struct {
	Err error `json:"err,omitempty"`
}

func (r profileResponse) error() error { return r.Err }
