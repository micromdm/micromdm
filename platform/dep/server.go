package dep

import (
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/micromdm/micromdm/pkg/httputil"
)

type Endpoints struct {
	DefineProfileEndpoint    endpoint.Endpoint
	FetchProfileEndpoint     endpoint.Endpoint
	GetAccountInfoEndpoint   endpoint.Endpoint
	GetDeviceDetailsEndpoint endpoint.Endpoint
	AssignProfileEndpoint    endpoint.Endpoint
}

func MakeServerEndpoints(s Service, outer endpoint.Middleware, others ...endpoint.Middleware) Endpoints {
	return Endpoints{
		AssignProfileEndpoint:    endpoint.Chain(outer, others...)(MakeAssignProfileEndpoint(s)),
		DefineProfileEndpoint:    endpoint.Chain(outer, others...)(MakeDefineProfileEndpoint(s)),
		FetchProfileEndpoint:     endpoint.Chain(outer, others...)(MakeFetchProfileEndpoint(s)),
		GetAccountInfoEndpoint:   endpoint.Chain(outer, others...)(MakeGetAccountInfoEndpoint(s)),
		GetDeviceDetailsEndpoint: endpoint.Chain(outer, others...)(MakeGetDeviceDetailsEndpoint(s)),
	}
}

func RegisterHTTPHandlers(r *mux.Router, e Endpoints, options ...httptransport.ServerOption) {
	// PUT		/v1/dep/profiles		define a DEP profile with mdmenrollment.apple.com
	// POST		/v1/dep/profiles		get a DEP profile given a known profile UUID
	// GET		/v1/dep/account			get information about the dep account
	// POST		/v1/dep/devices			get device details given a list of serials
	// POST		/v1/dep/assign			assign a specific profile ID to one or more serial

	r.Methods("PUT").Path("/v1/dep/profiles").Handler(httptransport.NewServer(
		e.DefineProfileEndpoint,
		decodeDefineProfileRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	r.Methods("POST").Path("/v1/dep/assign").Handler(httptransport.NewServer(
		e.AssignProfileEndpoint,
		decodeAssignProfileRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	r.Methods("POST").Path("/v1/dep/profiles").Handler(httptransport.NewServer(
		e.FetchProfileEndpoint,
		decodeFetchProfileRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	r.Methods("GET").Path("/v1/dep/account").Handler(httptransport.NewServer(
		e.GetAccountInfoEndpoint,
		decodeGetAccountInfoRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	r.Methods("POST").Path("/v1/dep/devices").Handler(httptransport.NewServer(
		e.GetDeviceDetailsEndpoint,
		decodeDeviceDetailsRequest,
		httputil.EncodeJSONResponse,
		options...,
	))
}
