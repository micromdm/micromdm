package blueprint

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/micromdm/pkg/httputil"
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

func (r removeBlueprintsResponse) Failed() error { return r.Err }

func decodeRemoveBlueprintsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req removeBlueprintsRequest
	err := httputil.DecodeJSONRequest(r, &req)
	return req, err
}

func decodeRemoveBlueprintsResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp removeBlueprintsResponse
	err := httputil.DecodeJSONResponse(r, &resp)
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
