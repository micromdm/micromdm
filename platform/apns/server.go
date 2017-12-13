package apns

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/micromdm/micromdm/pkg/httputil"
)

type Endpoints struct {
	PushEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		PushEndpoint: MakePushEndpoint(s),
	}
}

func MakeHTTPHandler(e Endpoints, logger log.Logger) *mux.Router {
	r, options := httputil.NewRouter(logger)

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

	return r
}
