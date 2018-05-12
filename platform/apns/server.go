package apns

import (
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/micromdm/micromdm/pkg/httputil"
)

type Endpoints struct {
	PushEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service, outer endpoint.Middleware, others ...endpoint.Middleware) Endpoints {
	return Endpoints{
		PushEndpoint: endpoint.Chain(outer, others...)(MakePushEndpoint(s)),
	}
}

func RegisterHTTPHandlers(r *mux.Router, e Endpoints, options ...httptransport.ServerOption) {
	// GET    /push/:udid		create an APNS Push notification for a managed device or user(deprecated)
	// POST   /v1/push/:udid	create an APNS Push notification for a managed device or user

	r.Methods("GET").Path("/push/{udid}").Handler(httptransport.NewServer(
		e.PushEndpoint,
		decodePushRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	r.Methods("POST").Path("/v1/push/{udid}").Handler(httptransport.NewServer(
		e.PushEndpoint,
		decodePushRequest,
		httputil.EncodeJSONResponse,
		options...,
	))
}
