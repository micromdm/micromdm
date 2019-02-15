package vpp

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"

	"github.com/micromdm/micromdm/pkg/httputil"
	"github.com/micromdm/micromdm/vpp"
)

func (svc *VPPService) GetServiceConfigSrv(ctx context.Context, options vpp.ServiceConfigSrvOptions) (*vpp.ServiceConfigSrv, error) {
	if svc.client == nil {
		return nil, errors.New("VPP not configured yet. add a VPP token to enable VPP")
	}
	return svc.client.GetServiceConfigSrv(options)
}

//type getVPPServiceConfigSrvRequest struct{}
type getServiceConfigSrvResponse struct {
	*vpp.ServiceConfigSrv
	Err error `json:"err,omitempty"`
}

func (r getServiceConfigSrvResponse) Failed() error { return r.Err }

func decodeGetServiceConfigSrvRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	var req vpp.ServiceConfigSrvOptions

	switch r.Method {
	case "POST":
		err := httputil.DecodeJSONRequest(r, &req)
		return req, err
	default:
		return req, nil
	}
}

func decodeGetServiceConfigSrvResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp getServiceConfigSrvResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeGetServiceConfigSrvEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(vpp.ServiceConfigSrvOptions)
		configsrv, err := svc.GetServiceConfigSrv(ctx, req)
		return getServiceConfigSrvResponse{ServiceConfigSrv: configsrv, Err: err}, nil
	}
}

func (e Endpoints) GetServiceConfigSrv(ctx context.Context, options vpp.ServiceConfigSrvOptions) (*vpp.ServiceConfigSrv, error) {
	response, err := e.GetServiceConfigSrvEndpoint(ctx, options)
	if err != nil {
		return nil, err
	}
	return response.(getServiceConfigSrvResponse).ServiceConfigSrv, response.(getServiceConfigSrvResponse).Err
}
