package vpp

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"

	"github.com/micromdm/micromdm/pkg/httputil"
	"github.com/micromdm/micromdm/vpp"
)

func (svc *VPPService) GetContentMetadata(ctx context.Context, options vpp.ContentMetadataOptions) (*vpp.ContentMetadata, error) {
	if svc.client == nil {
		return nil, errors.New("VPP not configured yet. add a VPP token to enable VPP")
	}
	return svc.client.GetContentMetadata(options)
}

type getContentMetadataResponse struct {
	*vpp.ContentMetadata
	Err error `json:"err,omitempty"`
}

func (r getContentMetadataResponse) Failed() error { return r.Err }

func decodeGetContentMetadataRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req vpp.ContentMetadataOptions

	switch r.Method {
	case "POST":
		err := httputil.DecodeJSONRequest(r, &req)
		return req, err
	default:
		return req, nil
	}
}

func decodeGetContentMetadataResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp getContentMetadataResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeGetContentMetadataEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(vpp.ContentMetadataOptions)
		srv, err := svc.GetContentMetadata(ctx, req)
		return getContentMetadataResponse{ContentMetadata: srv, Err: err}, nil
	}
}

func (e Endpoints) GetContentMetadata(ctx context.Context, options vpp.ContentMetadataOptions) (*vpp.ContentMetadata, error) {
	response, err := e.GetContentMetadataEndpoint(ctx, options)
	if err != nil {
		return nil, err
	}
	return response.(getContentMetadataResponse).ContentMetadata, response.(getContentMetadataResponse).Err
}
