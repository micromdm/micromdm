package appstore

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/micromdm/micromdm/pkg/httputil"
)

type Endpoints struct {
	AppUploadEndpoint endpoint.Endpoint
	ListAppsEndpoint  endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		AppUploadEndpoint: MakeUploadAppEndpiont(s),
		ListAppsEndpoint:  MakeListAppsEndpoint(s),
	}
}

func MakeHTTPHandler(e Endpoints, logger log.Logger) *mux.Router {
	r, options := httputil.NewRouter(logger)

	// POST    /v1/apps			upload an app to the server
	// GET     /v1/apps			list apps managed by the server

	r.Methods("POST").Path("/v1/apps").Handler(httptransport.NewServer(
		e.AppUploadEndpoint,
		decodeAppUploadRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	r.Methods("GET").Path("/v1/apps").Handler(httptransport.NewServer(
		e.ListAppsEndpoint,
		decodeListAppsRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	return r
}
