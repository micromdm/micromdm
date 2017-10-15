package config

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	SavePushCertificateEndpoint endpoint.Endpoint
}

type saveRequest struct {
	Cert []byte `json:"cert"`
	Key  []byte `json:"key"`
}

type saveResponse struct {
	Err error
}

func (r saveResponse) error() error { return r.Err }

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

func MakeSavePushCertificateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(saveRequest)
		err = svc.SavePushCertificate(ctx, req.Cert, req.Key)
		return saveResponse{Err: err}, nil
	}
}
