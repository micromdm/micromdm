package restream

import (
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/micromdm/micromdm/pkg/httputil"
)

type Endpoints struct {
	CommandStatusEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service, outer endpoint.Middleware, others ...endpoint.Middleware) Endpoints {
	return Endpoints{
		CommandStatusEndpoint: endpoint.Chain(outer, others...)(MakeCommandStatusEndpoint(s)),
	}
}

func RegisterHTTPHandlers(r *mux.Router, e Endpoints, options ...httptransport.ServerOption) {
	// GET    /v1/restream/:id	see the status of a command matching the uuid

	r.Methods("GET").Path("/v1/restream/{id}").Handler(httptransport.NewServer(
		e.CommandStatusEndpoint,
		decodeCommandStatusRequest,
		httputil.EncodeJSONResponse,
		options...,
	))
}
