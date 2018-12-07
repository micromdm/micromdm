package user

import (
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/micromdm/micromdm/pkg/httputil"
)

type Endpoints struct {
	ApplyUserEndpoint endpoint.Endpoint
	ListUsersEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service, outer endpoint.Middleware, others ...endpoint.Middleware) Endpoints {
	return Endpoints{
		ApplyUserEndpoint: endpoint.Chain(outer, others...)(MakeApplyUserEndpoint(s)),
		ListUsersEndpoint: endpoint.Chain(outer, others...)(MakeListUsersEndpoint(s)),
	}
}

func RegisterHTTPHandlers(r *mux.Router, e Endpoints, options ...httptransport.ServerOption) {
	// PUT     /v1/users		create or replace an user
	// POST    /v1/users		get a list of users managed by the server

	r.Methods("PUT").Path("/v1/users").Handler(httptransport.NewServer(
		e.ApplyUserEndpoint,
		decodeApplyUserRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	r.Methods("POST").Path("/v1/users").Handler(httptransport.NewServer(
		e.ListUsersEndpoint,
		decodeListUsersRequest,
		httputil.EncodeJSONResponse,
		options...,
	))
}
