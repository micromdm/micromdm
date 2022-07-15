package command

import (
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/micromdm/micromdm/pkg/httputil"
)

type Endpoints struct {
	NewCommandEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service, outer endpoint.Middleware, others ...endpoint.Middleware) Endpoints {
	return Endpoints{
		NewCommandEndpoint: endpoint.Chain(outer, others...)(MakeNewCommandEndpoint(s)),
	}
}

func RegisterHTTPHandlers(r *mux.Router, e Endpoints, options ...httptransport.ServerOption) {
	// POST     /v1/commands		Add new MDM Command to device queue.

	r.Methods("POST").Path("/v1/commands").Handler(httptransport.NewServer(
		e.NewCommandEndpoint,
		decodeNewCommandRequest,
		httputil.EncodeJSONResponse,
		options...,
	))
}
