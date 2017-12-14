package device

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/micromdm/micromdm/pkg/httputil"
)

type Endpoints struct {
	ListDevicesEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		ListDevicesEndpoint: MakeListDevicesEndpoint(s),
	}
}

func MakeHTTPHandler(e Endpoints, logger log.Logger) *mux.Router {
	r, options := httputil.NewRouter(logger)

	// GET     /v1/devices		get a list of devices managed by the server

	r.Methods("GET").Path("/v1/devices").Handler(httptransport.NewServer(
		e.ListDevicesEndpoint,
		decodeListDevicesRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	return r
}
