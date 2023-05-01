package command

import (
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/micromdm/micromdm/pkg/httputil"
)

type Endpoints struct {
	NewCommandEndpoint    endpoint.Endpoint
	NewRawCommandEndpoint endpoint.Endpoint
	ClearQueueEndpoint    endpoint.Endpoint
}

func MakeServerEndpoints(s Service, outer endpoint.Middleware, others ...endpoint.Middleware) Endpoints {
	return Endpoints{
		NewCommandEndpoint:    endpoint.Chain(outer, others...)(MakeNewCommandEndpoint(s)),
		NewRawCommandEndpoint: endpoint.Chain(outer, others...)(MakeNewRawCommandEndpoint(s)),
		ClearQueueEndpoint:    endpoint.Chain(outer, others...)(MakeClearQueueEndpoint(s)),
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

	// POST,PUT /v1/commands/udid		Add new MDM Command with raw plist to device queue.
	r.Methods("POST", "PUT").Path("/v1/commands/{udid}").Handler(httptransport.NewServer(
		e.NewRawCommandEndpoint,
		decodeNewRawCommandRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	// DELETE     /v1/commands/udid		Clear device queue.
	r.Methods("DELETE").Path("/v1/commands/{udid}").Handler(httptransport.NewServer(
		e.ClearQueueEndpoint,
		decodeClearRequest,
		httputil.EncodeJSONResponse,
		options...,
	))
}
