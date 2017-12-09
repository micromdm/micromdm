package blueprint

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

func (svc *BlueprintService) RemoveBlueprints(ctx context.Context, names []string) error {
	for _, name := range names {
		err := svc.store.Delete(name)
		if err != nil {
			return err
		}
	}
	return nil
}

type removeBlueprintsRequest struct {
	Names []string `json:"names"`
}

type removeBlueprintsResponse struct {
	Err error `json:"err,omitempty"`
}

func (r removeBlueprintsResponse) error() error { return r.Err }

func decodeRemoveBlueprintsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req removeBlueprintsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeRemoveBlueprintsResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errorDecoder(r)
	}
	var resp removeBlueprintsResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func MakeRemoveBlueprintsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(removeBlueprintsRequest)
		err = svc.RemoveBlueprints(ctx, req.Names)
		return removeBlueprintsResponse{
			Err: err,
		}, nil
	}
}

func (e Endpoints) RemoveBlueprints(ctx context.Context, names []string) error {
	request := removeBlueprintsRequest{Names: names}
	resp, err := e.RemoveBlueprintsEndpoint(ctx, request)
	if err != nil {
		return err
	}
	return resp.(removeBlueprintsResponse).Err
}
