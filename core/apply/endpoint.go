package apply

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/micromdm/blueprint"
)

type Endpoints struct {
	ApplyBlueprintEndpoint endpoint.Endpoint
	ApplyDEPTokensEndpoint endpoint.Endpoint
}

func (e Endpoints) ApplyBlueprint(ctx context.Context, bp *blueprint.Blueprint) error {
	request := blueprintRequest{Blueprint: bp}
	resp, err := e.ApplyBlueprintEndpoint(ctx, request)
	if err != nil {
		return err
	}
	return resp.(blueprintResponse).Err
}

func (e Endpoints) ApplyDEPToken(ctx context.Context, P7MContent []byte) error {
	req := depTokensRequest{P7MContent: P7MContent}
	resp, err := e.ApplyDEPTokensEndpoint(ctx, req)
	if err != nil {
		return err
	}
	return resp.(depTokensResponse).Err
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

func MakeApplyDEPTokensEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(depTokensRequest)
		err = svc.ApplyDEPToken(ctx, req.P7MContent)
		return depTokensResponse{
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

type depTokensRequest struct {
	P7MContent []byte `json:"p7m_content"`
}

type depTokensResponse struct {
	Err error `json:"err,omitempty"`
}

func (r depTokensResponse) error() error { return r.Err }
func (r blueprintResponse) error() error { return r.Err }
