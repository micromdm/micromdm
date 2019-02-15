package vpp

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"

	"github.com/micromdm/micromdm/pkg/httputil"
	"github.com/micromdm/micromdm/vpp"
)

func (svc *VPPService) ManageVPPLicensesByAdamIdSrv(ctx context.Context, options vpp.ManageVPPLicensesByAdamIdSrvOptions) (*vpp.ManageVPPLicensesByAdamIdSrv, error) {
	if svc.client == nil {
		return nil, errors.New("VPP not configured yet. add a VPP token to enable VPP")
	}
	return svc.client.ManageVPPLicensesByAdamIdSrv(options)
}

type manageVPPLicensesByAdamIdSrvResponse struct {
	*vpp.ManageVPPLicensesByAdamIdSrv
	Err error `json:"err,omitempty"`
}

func (r manageVPPLicensesByAdamIdSrvResponse) Failed() error { return r.Err }

func decodeManageVPPLicensesByAdamIdSrvRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req vpp.ManageVPPLicensesByAdamIdSrvOptions

	switch r.Method {
	case "PUT":
		err := httputil.DecodeJSONRequest(r, &req)
		return req, err
	default:
		return req, nil
	}
}

func decodeManageVPPLicensesByAdamIdSrvResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp manageVPPLicensesByAdamIdSrvResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeManageVPPLicensesByAdamIdSrvEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(vpp.ManageVPPLicensesByAdamIdSrvOptions)
		srv, err := svc.ManageVPPLicensesByAdamIdSrv(ctx, req)
		return manageVPPLicensesByAdamIdSrvResponse{ManageVPPLicensesByAdamIdSrv: srv, Err: err}, nil
	}
}

func (e Endpoints) ManageVPPLicensesByAdamIdSrv(ctx context.Context, options vpp.ManageVPPLicensesByAdamIdSrvOptions) (*vpp.ManageVPPLicensesByAdamIdSrv, error) {
	response, err := e.ManageVPPLicensesByAdamIdSrvEndpoint(ctx, options)
	if err != nil {
		return nil, err
	}
	return response.(manageVPPLicensesByAdamIdSrvResponse).ManageVPPLicensesByAdamIdSrv, response.(manageVPPLicensesByAdamIdSrvResponse).Err
}
