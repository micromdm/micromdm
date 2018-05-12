package command

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/micromdm/micromdm/pkg/httputil"
)

type Endpoints struct {
	NewCommandEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		NewCommandEndpoint: MakeNewCommandEndpoint(s),
	}
}

func MakeHTTPHandler(e Endpoints, logger log.Logger) *mux.Router {
	r, options := httputil.NewRouter(logger)

	// POST     /v1/commands		Add new MDM Command to device queue.

	r.Methods("POST").Path("/v1/commands").Handler(httptransport.NewServer(
		e.NewCommandEndpoint,
		decodeNewCommandRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	return r
}
