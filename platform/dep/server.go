package dep

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/micromdm/micromdm/pkg/httputil"
)

type Endpoints struct {
	DefineProfileEndpoint    endpoint.Endpoint
	FetchProfileEndpoint     endpoint.Endpoint
	GetAccountInfoEndpoint   endpoint.Endpoint
	GetDeviceDetailsEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		DefineProfileEndpoint:    MakeDefineProfileEndpoint(s),
		FetchProfileEndpoint:     MakeFetchProfileEndpoint(s),
		GetAccountInfoEndpoint:   MakeGetAccountInfoEndpoint(s),
		GetDeviceDetailsEndpoint: MakeGetDeviceDetailsEndpoint(s),
	}
}

func MakeHTTPHandler(e Endpoints, logger log.Logger) *mux.Router {
	r, options := httputil.NewRouter(logger)

	// PUT		/v1/dep/profiles		define a DEP profile with mdmenrollment.apple.com
	// POST		/v1/dep/profiles		get a DEP profile given a known profile UUID
	// GET		/v1/dep/account			get information about the dep account
	// POST		/v1/dep/devices			get device details given a list of serials

	r.Methods("PUT").Path("/v1/dep/profiles").Handler(httptransport.NewServer(
		e.DefineProfileEndpoint,
		decodeDefineProfileRequest,
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

	return r
}
