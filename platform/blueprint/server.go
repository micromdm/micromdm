package blueprint

import (
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/micromdm/micromdm/pkg/httputil"
)

type Endpoints struct {
	ApplyBlueprintEndpoint   endpoint.Endpoint
	GetBlueprintsEndpoint    endpoint.Endpoint
	RemoveBlueprintsEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service, outer endpoint.Middleware, others ...endpoint.Middleware) Endpoints {
	return Endpoints{
		GetBlueprintsEndpoint:    endpoint.Chain(outer, others...)(MakeGetBlueprintsEndpoint(s)),
		ApplyBlueprintEndpoint:   endpoint.Chain(outer, others...)(MakeApplyBlueprintEndpoint(s)),
		RemoveBlueprintsEndpoint: endpoint.Chain(outer, others...)(MakeRemoveBlueprintsEndpoint(s)),
	}
}

func RegisterHTTPHandlers(r *mux.Router, e Endpoints, options ...httptransport.ServerOption) {
	// PUT     /v1/blueprints			create or replace a blueprint on the server
	// POST    /v1/blueprints			get a list of blueprints managed by the server
	// DELETE  /v1/blueprints			remove one or more blueprints from the server

	r.Methods("PUT").Path("/v1/blueprints").Handler(httptransport.NewServer(
		e.ApplyBlueprintEndpoint,
		decodeApplyBlueprintRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	r.Methods("POST").Path("/v1/blueprints").Handler(httptransport.NewServer(
		e.GetBlueprintsEndpoint,
		decodeGetBlueprintsRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	r.Methods("DELETE").Path("/v1/blueprints").Handler(httptransport.NewServer(
		e.RemoveBlueprintsEndpoint,
		decodeRemoveBlueprintsRequest,
		httputil.EncodeJSONResponse,
		options...,
	))
}
