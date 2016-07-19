package management

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/micromdm/certificates"
	"golang.org/x/net/context"
)

type certificatesRequest struct {
	UUID string
}

type certificatesResponse struct {
	certificates []certificates.Certificate
	Err          error `json:"error,omitempty"`
}

func makeCertificatesEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(certificatesRequest)
		certs, err := svc.Certificates(req.UUID)
		if err != nil {
			return certificatesResponse{Err: err}, nil
		}
		return certificatesResponse{certificates: certs}, nil
	}
}
