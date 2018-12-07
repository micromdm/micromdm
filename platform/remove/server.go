package remove

import (
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/micromdm/micromdm/pkg/httputil"
)

type Endpoints struct {
	BlockDeviceEndpoint   endpoint.Endpoint
	UnblockDeviceEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service, outer endpoint.Middleware, others ...endpoint.Middleware) Endpoints {
	return Endpoints{
		BlockDeviceEndpoint:   endpoint.Chain(outer, others...)(MakeBlockDeviceEndpoint(s)),
		UnblockDeviceEndpoint: endpoint.Chain(outer, others...)(MakeUnblockDeviceEndpoint(s)),
	}
}

func RegisterHTTPHandlers(r *mux.Router, e Endpoints, options ...httptransport.ServerOption) {
	// POST		/v1/devices/:udid/block			force a device to unenroll next time it connects
	// POST		/v1/devices/:udid/unblock		allow a blocked device to enroll again

	r.Methods("POST").Path("/v1/devices/{udid}/block").Handler(httptransport.NewServer(
		e.BlockDeviceEndpoint,
		decodeBlockDeviceRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	r.Methods("POST").Path("/v1/devices/{udid}/unblock").Handler(httptransport.NewServer(
		e.UnblockDeviceEndpoint,
		decodeUnblockDeviceRequest,
		httputil.EncodeJSONResponse,
		options...,
	))
}
