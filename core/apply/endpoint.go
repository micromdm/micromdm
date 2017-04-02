package apply

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/micromdm/blueprint"
)

type Endpoints struct {
	ApplyBlueprintEndpoint endpoint.Endpoint
}

func (e Endpoints) ApplyBlueprint(ctx context.Context, bp *blueprint.Blueprint) error {
	request := blueprintRequest{Blueprint: bp}
	resp, err := e.ApplyBlueprintEndpoint(ctx, request)
	if err != nil {
		return err
	}
	return resp.(blueprintResponse).Err
}

func MakeApplyBlueprintEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(blueprintRequest)
		err = svc.ApplyBlueprint(ctx, req.Blueprint)
		return blueprintResponse{
			Err: err,
		}, nil
	}
}

type blueprintRequest struct {
	Blueprint *blueprint.Blueprint `json:"blueprint"`
}

type blueprintResponse struct {
	Err error `json:"err,omitempty"`
}
