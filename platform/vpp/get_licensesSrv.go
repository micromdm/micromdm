package vpp

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"

	"github.com/micromdm/micromdm/pkg/httputil"
	"github.com/micromdm/micromdm/vpp"
)

func (svc *VPPService) GetLicensesSrv(ctx context.Context, options vpp.LicensesSrvOptions) (*vpp.LicensesSrv, error) {
	if svc.client == nil {
		return nil, errors.New("VPP not configured yet. add a VPP token to enable VPP")
	}
	return svc.client.GetLicensesSrv(options)
}

type getLicensesSrvResponse struct {
	*vpp.LicensesSrv
	Err error `json:"err,omitempty"`
}

func (r getLicensesSrvResponse) Failed() error { return r.Err }

func decodeGetLicensesSrvRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req vpp.LicensesSrvOptions

	switch r.Method {
	case "POST":
		err := httputil.DecodeJSONRequest(r, &req)
		return req, err
	default:
		return req, nil
	}
}

func decodeGetLicensesSrvResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp getLicensesSrvResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeGetLicensesSrvEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(vpp.LicensesSrvOptions)
		srv, err := svc.GetLicensesSrv(ctx, req)
		return getLicensesSrvResponse{LicensesSrv: srv, Err: err}, nil
	}
}

func (e Endpoints) GetLicensesSrv(ctx context.Context, options vpp.LicensesSrvOptions) (*vpp.LicensesSrv, error) {
	response, err := e.GetLicensesSrvEndpoint(ctx, options)
	if err != nil {
		return nil, err
	}
	return response.(getLicensesSrvResponse).LicensesSrv, response.(getLicensesSrvResponse).Err
}
