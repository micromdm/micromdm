package vpp

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"

	"github.com/micromdm/micromdm/pkg/httputil"
	"github.com/micromdm/micromdm/vpp"
)

func (svc *VPPService) GetVPPApps(ctx context.Context, ids []string) (*vpp.VPPAppResponse, error) {
	if svc.client == nil {
		return nil, errors.New("VPP not configured yet. add a VPP token to enable VPP")
	}
	return svc.client.GetVPPApps(ids...)
}

type vppAppsRequest struct {
	IDs []string `json:"ids"`
}

type getVPPAppsResponse struct {
	*vpp.VPPAppResponse
	Err error `json:"err,omitempty"`
}

func (r getVPPAppsResponse) Failed() error { return r.Err }

func decodeGetVPPAppsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req vppAppsRequest

	switch r.Method {
	case "POST":
		err := httputil.DecodeJSONRequest(r, &req)
		return req, err
	default:
		return req, nil
	}
}

func decodeGetVPPAppsResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp getVPPAppsResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeGetVPPAppsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(vppAppsRequest)
		apps, err := svc.GetVPPApps(ctx, req.IDs)
		return getVPPAppsResponse{VPPAppResponse: apps, Err: err}, nil
	}
}

func (e Endpoints) GetVPPApps(ctx context.Context, ids []string) (*vpp.VPPAppResponse, error) {
	request := vppAppsRequest{IDs: ids}
	response, err := e.GetVPPAppsEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(getVPPAppsResponse).VPPAppResponse, response.(getVPPAppsResponse).Err
}
