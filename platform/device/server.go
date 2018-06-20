package device

import (
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/micromdm/micromdm/pkg/httputil"
)

type Endpoints struct {
	ListDevicesEndpoint   endpoint.Endpoint
	RemoveDevicesEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service, outer endpoint.Middleware, others ...endpoint.Middleware) Endpoints {
	return Endpoints{
		ListDevicesEndpoint:   endpoint.Chain(outer, others...)(MakeListDevicesEndpoint(s)),
		RemoveDevicesEndpoint: endpoint.Chain(outer, others...)(MakeRemoveDevicesEndpoint(s)),
	}
}

func RegisterHTTPHandlers(r *mux.Router, e Endpoints, options ...httptransport.ServerOption) {
	// GET     /v1/devices		get a list of devices managed by the server
	// DELETE  /v1/devices		remove one or more devices from the server

	r.Methods("GET").Path("/v1/devices").Handler(httptransport.NewServer(
		e.ListDevicesEndpoint,
		decodeListDevicesRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	r.Methods("DELETE").Path("/v1/devices").Handler(httptransport.NewServer(
		e.RemoveDevicesEndpoint,
		decodeRemoveDevicesRequest,
		httputil.EncodeJSONResponse,
		options...,
	))
}
