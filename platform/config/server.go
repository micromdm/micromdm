package config

import (
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/micromdm/micromdm/pkg/httputil"
)

type Endpoints struct {
	SavePushCertificateEndpoint endpoint.Endpoint
	GetPushCertificateEndpoint  endpoint.Endpoint
	ApplyDEPTokensEndpoint      endpoint.Endpoint
	GetDEPTokensEndpoint        endpoint.Endpoint
}

func MakeServerEndpoints(s Service, outer endpoint.Middleware, others ...endpoint.Middleware) Endpoints {
	return Endpoints{
		SavePushCertificateEndpoint: endpoint.Chain(outer, others...)(MakeSavePushCertificateEndpoint(s)),
		GetPushCertificateEndpoint:  endpoint.Chain(outer, others...)(MakeGetPushCertificateEndpoint(s)),
		ApplyDEPTokensEndpoint:      endpoint.Chain(outer, others...)(MakeApplyDEPTokensEndpoint(s)),
		GetDEPTokensEndpoint:        endpoint.Chain(outer, others...)(MakeGetDEPTokensEndpoint(s)),
	}
}

func RegisterHTTPHandlers(r *mux.Router, e Endpoints, options ...httptransport.ServerOption) {
	// PUT     /v1/config/certificate		create or replace the MDM Push Certificate
	// GET     /v1/config/certificate		retrieve the MDM Push Certificate
	// PUT     /v1/dep-tokens				create or replace a DEP OAuth token
	// GET     /v1/dep-tokens				get the OAuth Token used for the DEP client

	r.Methods("PUT").Path("/v1/config/certificate").Handler(httptransport.NewServer(
		e.SavePushCertificateEndpoint,
		decodeSavePushCertificateRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	r.Methods("GET").Path("/v1/config/certificate").Handler(httptransport.NewServer(
		e.GetPushCertificateEndpoint,
		decodeGetPushCertificateRequest,
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
}
