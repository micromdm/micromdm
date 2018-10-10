package appstore

import (
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/micromdm/micromdm/pkg/httputil"
)

type Endpoints struct {
	AppUploadEndpoint endpoint.Endpoint
	ListAppsEndpoint  endpoint.Endpoint
}

func MakeServerEndpoints(s Service, outer endpoint.Middleware, others ...endpoint.Middleware) Endpoints {
	return Endpoints{
		AppUploadEndpoint: endpoint.Chain(outer, others...)(MakeUploadAppEndpiont(s)),
		ListAppsEndpoint:  endpoint.Chain(outer, others...)(MakeListAppsEndpoint(s)),
	}
}

func RegisterHTTPHandlers(r *mux.Router, e Endpoints, options ...httptransport.ServerOption) {
	// PUT		/v1/apps			upload an app to the server
	// POST   /v1/apps			list apps managed by the server

	r.Methods("PUT").Path("/v1/apps").Handler(httptransport.NewServer(
		e.AppUploadEndpoint,
		decodeAppUploadRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	r.Methods("POST").Path("/v1/apps").Handler(httptransport.NewServer(
		e.ListAppsEndpoint,
		decodeListAppsRequest,
		httputil.EncodeJSONResponse,
		options...,
	))
}
