package config

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/micromdm/pkg/httputil"
	"github.com/pkg/errors"
)

func (svc *ConfigService) GetPushCertificate(ctx context.Context) ([]byte, error) {
	cert, err := svc.store.GetPushCertificate(ctx)
	if err != nil {
		return cert, errors.Wrap(err, "get push certificate")
	}
	return cert, nil
}

type getResponse struct {
	Cert []byte `json:"cert"`
	Err  error  `json:"err,omitempty"`
}

func (r getResponse) Failed() error { return r.Err }

func decodeGetPushCertificateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeGetPushCertificateResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp getResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeGetPushCertificateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		cert, err := svc.GetPushCertificate(ctx)
		return getResponse{Err: err, Cert: cert}, nil
	}
}

func (e Endpoints) GetPushCertificate(ctx context.Context) ([]byte, error) {
	response, err := e.GetPushCertificateEndpoint(ctx, nil)
	if err != nil {
		return nil, err
	}

	return response.(getResponse).Cert, response.(getResponse).Err
}
