package blueprint

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/micromdm/micromdm/pkg/httputil"
)

type Endpoints struct {
	ApplyBlueprintEndpoint   endpoint.Endpoint
	GetBlueprintsEndpoint    endpoint.Endpoint
	RemoveBlueprintsEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		GetBlueprintsEndpoint:    MakeGetBlueprintsEndpoint(s),
		ApplyBlueprintEndpoint:   MakeApplyBlueprintEndpoint(s),
		RemoveBlueprintsEndpoint: MakeRemoveBlueprintsEndpoint(s),
	}
}

func MakeHTTPHandler(e Endpoints, logger log.Logger) *mux.Router {
	r, options := httputil.NewRouter(logger)

	// PUT     /v1/blueprints			create or replace a blueprint on the server
	// GET     /v1/blueprints			get a list of blueprints managed by the server
	// DELETE  /v1/blueprints			remove one or more blueprints from the server

	r.Methods("PUT").Path("/v1/blueprints").Handler(httptransport.NewServer(
		e.ApplyBlueprintEndpoint,
		decodeApplyBlueprintRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	r.Methods("GET").Path("/v1/blueprints").Handler(httptransport.NewServer(
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

	return r
}
