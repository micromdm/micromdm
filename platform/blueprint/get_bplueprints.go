package blueprint

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

func (svc *BlueprintService) GetBlueprints(ctx context.Context, opt GetBlueprintsOption) ([]Blueprint, error) {
	if opt.FilterName != "" {
		bp, err := svc.store.BlueprintByName(opt.FilterName)
		if err != nil {
			return nil, err
		}
		return []Blueprint{*bp}, err
	} else {
		bps, err := svc.store.List()
		if err != nil {
			return nil, err
		}
		return bps, nil
	}
}

type getBlueprintsRequest struct{ Opts GetBlueprintsOption }
type getBlueprintsResponse struct {
	Blueprints []Blueprint `json:"blueprints"`
	Err        error       `json:"err,omitempty"`
}

func (r getBlueprintsResponse) error() error { return r.Err }

func decodeGetBlueprintsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var opts GetBlueprintsOption
	if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
		return nil, err
	}
	req := getBlueprintsRequest{
		Opts: opts,
	}
	return req, nil
}

func decodeGetBlueprintsResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errorDecoder(r)
	}
	var resp getBlueprintsResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func MakeGetBlueprintsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getBlueprintsRequest)
		blueprints, err := svc.GetBlueprints(ctx, req.Opts)
		return getBlueprintsResponse{
			Blueprints: blueprints,
			Err:        err,
		}, nil
	}
}

func (e Endpoints) GetBlueprints(ctx context.Context, opt GetBlueprintsOption) ([]Blueprint, error) {
	request := getBlueprintsRequest{opt}
	response, err := e.GetBlueprintsEndpoint(ctx, request.Opts)
	if err != nil {
		return nil, err
	}
	return response.(getBlueprintsResponse).Blueprints, response.(getBlueprintsResponse).Err
}
