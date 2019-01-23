package config

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/micromdm/pkg/httputil"
)

func (svc *ConfigService) GetDEPTokens(ctx context.Context) ([]DEPToken, []byte, error) {
	_, cert, err := svc.store.DEPKeypair(ctx)
	if err != nil {
		return nil, nil, err
	}
	var certBytes []byte
	if cert != nil {
		certBytes = cert.Raw
	}

	tokens, err := svc.store.DEPTokens(ctx)
	if err != nil {
		return nil, certBytes, err
	}

	return tokens, certBytes, nil
}

type getDEPTokenResponse struct {
	DEPTokens []DEPToken `json:"dep_tokens"`
	DEPPubKey []byte     `json:"public_key"`
	Err       error      `json:"err,omitempty"`
}

func (r getDEPTokenResponse) Failed() error { return r.Err }

func decodeGetDEPTokensRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeGetDEPTokensResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp getDEPTokenResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeGetDEPTokensEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tokens, pubkey, err := svc.GetDEPTokens(ctx)
		return getDEPTokenResponse{
			DEPTokens: tokens,
			DEPPubKey: pubkey,
			Err:       err,
		}, nil
	}
}

func (e Endpoints) GetDEPTokens(ctx context.Context) ([]DEPToken, []byte, error) {
	resp, err := e.GetDEPTokensEndpoint(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	return resp.(getDEPTokenResponse).DEPTokens, resp.(getDEPTokenResponse).DEPPubKey, nil
}
