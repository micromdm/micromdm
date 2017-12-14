package config

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/micromdm/pkg/httputil"
	"github.com/pkg/errors"
)

func (svc *ConfigService) SavePushCertificate(ctx context.Context, cert, key []byte) error {
	err := svc.store.SavePushCertificate(cert, key)
	return errors.Wrap(err, "save push certificate")
}

type saveRequest struct {
	Cert []byte `json:"cert"`
	Key  []byte `json:"key"`
}

type saveResponse struct {
	Err error `json:"err,omitempty"`
}

func (r saveResponse) Failed() error { return r.Err }

func decodeSavePushCertificateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req saveRequest
	err := httputil.DecodeJSONRequest(r, &req)
	return req, err
}

func decodeSavePushCertificateResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp saveResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeSavePushCertificateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(saveRequest)
		err = svc.SavePushCertificate(ctx, req.Cert, req.Key)
		return saveResponse{Err: err}, nil
	}
}

func (e Endpoints) SavePushCertificate(ctx context.Context, cert, key []byte) error {
	request := saveRequest{
		Cert: cert,
		Key:  key,
	}

	response, err := e.SavePushCertificateEndpoint(ctx, request)
	if err != nil {
		return err
	}

	return response.(saveResponse).Err
}
