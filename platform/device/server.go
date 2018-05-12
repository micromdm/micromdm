package device

import (
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/micromdm/micromdm/pkg/httputil"
)

type Endpoints struct {
	ListDevicesEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service, outer endpoint.Middleware, others ...endpoint.Middleware) Endpoints {
	return Endpoints{
		ListDevicesEndpoint: endpoint.Chain(outer, others...)(MakeListDevicesEndpoint(s)),
	}
}

func RegisterHTTPHandlers(r *mux.Router, e Endpoints, options ...httptransport.ServerOption) {
	// GET     /v1/devices		get a list of devices managed by the server

	r.Methods("GET").Path("/v1/devices").Handler(httptransport.NewServer(
		e.ListDevicesEndpoint,
		decodeListDevicesRequest,
		httputil.EncodeJSONResponse,
		options...,
	))
}
