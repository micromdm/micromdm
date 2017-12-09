package blueprint

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

func (svc *BlueprintService) ApplyBlueprint(ctx context.Context, bp *Blueprint) error {
	return svc.store.Save(bp)
}

type applyBlueprintRequest struct {
	Blueprint *Blueprint `json:"blueprint"`
}

type applyBlueprintResponse struct {
	Err error `json:"err,omitempty"`
}

func (r applyBlueprintResponse) error() error { return r.Err }

func decodeApplyBlueprintRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var bpReq applyBlueprintRequest
	if err := json.NewDecoder(r.Body).Decode(&bpReq); err != nil {
		return nil, err
	}
	return bpReq, nil
}

func decodeApplyBlueprintResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errorDecoder(r)
	}
	var resp applyBlueprintResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func MakeApplyBlueprintEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(applyBlueprintRequest)
		err = svc.ApplyBlueprint(ctx, req.Blueprint)
		return applyBlueprintResponse{
			Err: err,
		}, nil
	}
}

func (e Endpoints) ApplyBlueprint(ctx context.Context, bp *Blueprint) error {
	request := applyBlueprintRequest{Blueprint: bp}
	resp, err := e.ApplyBlueprintEndpoint(ctx, request)
	if err != nil {
		return err
	}
	return resp.(applyBlueprintResponse).Err
}
