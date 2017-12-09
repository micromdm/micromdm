package blueprint

import (
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type Endpoints struct {
	ApplyBlueprintEndpoint   endpoint.Endpoint
	GetBlueprintsEndpoint    endpoint.Endpoint
	RemoveBlueprintsEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		GetBlueprintsEndpoint:  MakeGetBlueprintsEndpoint(s),
		ApplyBlueprintEndpoint: MakeApplyBlueprintEndpoint(s),
	}
}

func MakeHTTPHandler(e Endpoints, logger log.Logger) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
	}

	r := mux.NewRouter()

	// PUT     /v1/blueprints			create or replace a blueprint on the server
	// GET     /v1/blueprints			get a list of blueprints managed by the server
	// DELETE  /v1/blueprints			remove one or more blueprints from the server

	r.Methods("PUT").Path("/v1/blueprints").Handler(httptransport.NewServer(
		e.ApplyBlueprintEndpoint,
		decodeApplyBlueprintRequest,
		httptransport.EncodeJSONResponse,
		options...,
	))

	r.Methods("GET").Path("/v1/blueprints").Handler(httptransport.NewServer(
		e.GetBlueprintsEndpoint,
		decodeGetBlueprintsRequest,
		httptransport.EncodeJSONResponse,
		options...,
	))

	r.Methods("DELETE").Path("/v1/blueprints").Handler(httptransport.NewServer(
		e.RemoveBlueprintsEndpoint,
		decodeRemoveBlueprintsRequest,
		httptransport.EncodeJSONResponse,
		options...,
	))

	return r
}

type errorWrapper struct {
	Error string `json:"error"`
}

type errorer interface {
	error() error
}

func errorDecoder(r *http.Response) error {
	var w errorWrapper
	if err := json.NewDecoder(r.Body).Decode(&w); err != nil {
		return err
	}
	return errors.New(w.Error)
}
