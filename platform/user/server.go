package user

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/micromdm/micromdm/pkg/httputil"
)

type Endpoints struct {
	ApplyUserEndpoint endpoint.Endpoint
	ListUsersEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		ApplyUserEndpoint: MakeApplyUserEndpoint(s),
		ListUsersEndpoint: MakeListUsersEndpoint(s),
	}
}

func MakeHTTPHandler(e Endpoints, logger log.Logger) *mux.Router {
	r, options := httputil.NewRouter(logger)

	// PUT     /v1/users		create or replace an user
	// GET     /v1/users		get a list of users managed by the server

	r.Methods("PUT").Path("/v1/users").Handler(httptransport.NewServer(
		e.ApplyUserEndpoint,
		decodeApplyUserRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	r.Methods("GET").Path("/v1/users").Handler(httptransport.NewServer(
		e.ListUsersEndpoint,
		decodeListUsersRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	return r
}
