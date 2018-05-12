package profile

import (
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/micromdm/micromdm/pkg/httputil"
)

type Endpoints struct {
	ApplyProfileEndpoint   endpoint.Endpoint
	GetProfilesEndpoint    endpoint.Endpoint
	RemoveProfilesEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service, outer endpoint.Middleware, others ...endpoint.Middleware) Endpoints {
	return Endpoints{
		ApplyProfileEndpoint:   endpoint.Chain(outer, others...)(MakeApplyProfileEndpoint(s)),
		GetProfilesEndpoint:    endpoint.Chain(outer, others...)(MakeGetProfilesEndpoint(s)),
		RemoveProfilesEndpoint: endpoint.Chain(outer, others...)(MakeRemoveProfilesEndpoint(s)),
	}
}

func RegisterHTTPHandlers(r *mux.Router, e Endpoints, options ...httptransport.ServerOption) {
	// GET     /v1/profiles		get a list of profiles managed by the server
	// PUT     /v1/profiles		create or replace a profile on the server
	// DELETE  /v1/profiles		remove one or more profiles from the server

	r.Methods("GET").Path("/v1/profiles").Handler(httptransport.NewServer(
		e.GetProfilesEndpoint,
		decodeGetProfilesRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	r.Methods("PUT").Path("/v1/profiles").Handler(httptransport.NewServer(
		e.ApplyProfileEndpoint,
		decodeApplyProfileRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	r.Methods("DELETE").Path("/v1/profiles").Handler(httptransport.NewServer(
		e.RemoveProfilesEndpoint,
		decodeRemoveProfilesRequest,
		httputil.EncodeJSONResponse,
		options...,
	))
}
