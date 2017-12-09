package profile

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

type Endpoints struct {
	ApplyProfileEndpoint   endpoint.Endpoint
	GetProfilesEndpoint    endpoint.Endpoint
	RemoveProfilesEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		ApplyProfileEndpoint:   MakeApplyProfileEndpoint(s),
		GetProfilesEndpoint:    MakeGetProfilesEndpoint(s),
		RemoveProfilesEndpoint: MakeRemoveProfilesEndpoint(s),
	}
}

func MakeHTTPHandler(e Endpoints, logger log.Logger) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
	}

	r := mux.NewRouter()

	// GET     /v1/profiles		get a list of profiles managed by the server
	// PUT     /v1/profiles		create or replace a profile on the server
	// DELETE  /v1/profiles		remove one or more profiles from the server

	r.Methods("GET").Path("/v1/profiles").Handler(httptransport.NewServer(
		e.GetProfilesEndpoint,
		decodeGetProfilesRequest,
		httptransport.EncodeJSONResponse,
		options...,
	))

	r.Methods("PUT").Path("/v1/profiles").Handler(httptransport.NewServer(
		e.ApplyProfileEndpoint,
		decodeApplyProfileRequest,
		httptransport.EncodeJSONResponse,
		options...,
	))

	r.Methods("DELETE").Path("/v1/profiles").Handler(httptransport.NewServer(
		e.RemoveProfilesEndpoint,
		decodeRemoveProfilesRequest,
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
