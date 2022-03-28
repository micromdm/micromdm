package config

import (
	"context"
	"fmt"
	"net/http"
	"encoding/json"
	"strings"

	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/micromdm/pkg/httputil"
)

func (svc *ConfigService) RemoveVPPToken(ctx context.Context, token_Content []byte) error {

	// Save the data to a SToken
	var sToken SToken
	err := json.NewDecoder(strings.NewReader(string(token_Content))).Decode(&sToken)
	if err != nil {
		return err
	}

	err = svc.store.DeleteVPPToken(string(sToken.Token))
	if err != nil {
		return err
	}
	fmt.Println("removed VPP token with Token", string(sToken.Token))
	return nil
}

type removeVPPTokenRequest struct {
	Content []byte `json:"sToken"`
}

type removeVPPTokenResponse struct {
	Err error `json:"err,omitempty"`
}

func (r removeVPPTokenResponse) Failed() error { return r.Err }

func decodeRemoveVPPTokensRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req removeVPPTokenRequest
	err := httputil.DecodeJSONRequest(r, &req)
	return req, err
}

func decodeRemoveVPPTokensResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp removeVPPTokenResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeRemoveVPPTokensEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(removeVPPTokenRequest)
		err = svc.RemoveVPPToken(ctx, req.Content)
		return removeVPPTokenResponse{
			Err: err,
		}, nil
	}
}

func (e Endpoints) RemoveVPPToken(ctx context.Context, Content []byte) error {
	req := removeVPPTokenRequest{Content: Content}
	resp, err := e.RemoveVPPTokensEndpoint(ctx, req)
	if err != nil {
		return err
	}
	return resp.(removeVPPTokenResponse).Err
}
