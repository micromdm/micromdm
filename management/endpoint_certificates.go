package management

import (
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/micromdm/certificate"
	"golang.org/x/net/context"
	"net/http"
)

type listCertificatesRequest struct {
	UUID string
}

type listCertificatesResponse struct {
	certificates []certificate.Certificate `json:"certificates,omitempty"`
	Err          error                     `json:"error,omitempty"`
}

func (r listCertificatesResponse) error() error { return r.Err }

func (r listCertificatesResponse) encodeList(w http.ResponseWriter) error {
	jsn, err := json.MarshalIndent(r.certificates, "", "  ")
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(jsn)
	return nil
}

func makeCertificatesEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listCertificatesRequest)
		certs, err := svc.Certificates(req.UUID)
		if err != nil {
			return listCertificatesResponse{Err: err}, nil
		}
		return listCertificatesResponse{certificates: certs}, nil
	}
}
