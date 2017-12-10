package config

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/micromdm/micromdm/pkg/httputil"
)

type Endpoints struct {
	ApplyDEPTokensEndpoint      endpoint.Endpoint
	SavePushCertificateEndpoint endpoint.Endpoint
	GetDEPTokensEndpoint        endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		ApplyDEPTokensEndpoint:      MakeApplyDEPTokensEndpoint(s),
		SavePushCertificateEndpoint: MakeSavePushCertificateEndpoint(s),
		GetDEPTokensEndpoint:        MakeGetDEPTokensEndpoint(s),
	}
}

func MakeHTTPHandler(e Endpoints, logger log.Logger) *mux.Router {
	r, options := httputil.NewRouter(logger)

	// PUT     /v1/config/certificate		create or replace the MDM Push Certificate
	// PUT     /v1/dep-tokens				create or replace a DEP OAuth token
	// GET     /v1/dep-tokens				get the OAuth Token used for the DEP client

	r.Methods("PUT").Path("/v1/config/certificate").Handler(httptransport.NewServer(
		e.SavePushCertificateEndpoint,
		decodeSavePushCertificateRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	r.Methods("PUT").Path("/v1/dep-tokens").Handler(httptransport.NewServer(
		e.ApplyDEPTokensEndpoint,
		decodeApplyDEPTokensRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	r.Methods("GET").Path("/v1/dep-tokens").Handler(httptransport.NewServer(
		e.GetDEPTokensEndpoint,
		decodeGetDEPTokensRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	return r
}
