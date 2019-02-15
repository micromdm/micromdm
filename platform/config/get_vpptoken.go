package config

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/micromdm/pkg/httputil"
)

func (svc *ConfigService) GetVPPTokens(ctx context.Context) ([]VPPToken, error) {
	tokens, err := svc.store.VPPTokens()
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

type getVPPTokenResponse struct {
	VPPTokens []VPPToken `json:"vpp_tokens"`
	Err       error      `json:"err,omitempty"`
}

func (r getVPPTokenResponse) Failed() error { return r.Err }

func decodeGetVPPTokensRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeGetVPPTokensResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp getVPPTokenResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeGetVPPTokensEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tokens, err := svc.GetVPPTokens(ctx)
		return getVPPTokenResponse{
			VPPTokens: tokens,
			Err:       err,
		}, nil
	}
}

func (e Endpoints) GetVPPTokens(ctx context.Context) ([]VPPToken, error) {
	resp, err := e.GetVPPTokensEndpoint(ctx, nil)
	if err != nil {
		return nil, err
	}
	return resp.(getVPPTokenResponse).VPPTokens, nil
}
