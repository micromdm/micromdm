package vpp

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"

	"github.com/micromdm/micromdm/pkg/httputil"
	"github.com/micromdm/micromdm/vpp"
)

type GetAssetsSrvOptions struct {
	AssetsSrvOptions vpp.AssetsSrvOptions
}

func (svc *VPPService) GetAssetsSrv(ctx context.Context, options vpp.AssetsSrvOptions) (*vpp.AssetsSrv, error) {
	if svc.client == nil {
		return nil, errors.New("VPP not configured yet. add a VPP token to enable VPP")
	}
	return svc.client.GetAssetsSrv(options)
}

type getAssetsSrvResponse struct {
	*vpp.AssetsSrv
	Err error `json:"err,omitempty"`
}

func (r getAssetsSrvResponse) Failed() error { return r.Err }

func decodeGetAssetsSrvRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req vpp.AssetsSrvOptions

	switch r.Method {
	case "POST":
		err := httputil.DecodeJSONRequest(r, &req)
		return req, err
	default:
		return req, nil
	}
}

func decodeGetAssetsSrvResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp getAssetsSrvResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeGetAssetsSrvEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(vpp.AssetsSrvOptions)
		srv, err := svc.GetAssetsSrv(ctx, req)
		return getAssetsSrvResponse{AssetsSrv: srv, Err: err}, nil
	}
}

func (e Endpoints) GetAssetsSrv(ctx context.Context, options vpp.AssetsSrvOptions) (*vpp.AssetsSrv, error) {
	response, err := e.GetAssetsSrvEndpoint(ctx, options)
	if err != nil {
		return nil, err
	}
	return response.(getAssetsSrvResponse).AssetsSrv, response.(getAssetsSrvResponse).Err
}
