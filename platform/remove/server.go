package remove

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

type Endpoints struct {
	BlockDeviceEndpoint   endpoint.Endpoint
	UnblockDeviceEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		BlockDeviceEndpoint:   MakeBlockDeviceEndpoint(s),
		UnblockDeviceEndpoint: MakeUnblockDeviceEndpoint(s),
	}
}

func MakeHTTPHandler(e Endpoints, logger log.Logger) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
	}

	r := mux.NewRouter()

	// POST		/v1/devices/:udid/block			force a device to unenroll next time it connects
	// POST		/v1/devices/:udid/unblock		allow a blocked device to enroll again

	r.Methods("POST").Path("/v1/devices/{udid}/block").Handler(httptransport.NewServer(
		e.BlockDeviceEndpoint,
		decodeBlockDeviceRequest,
		httptransport.EncodeJSONResponse,
		options...,
	))

	r.Methods("POST").Path("/v1/devices/{udid}/unblock").Handler(httptransport.NewServer(
		e.UnblockDeviceEndpoint,
		decodeUnblockDeviceRequest,
		httptransport.EncodeJSONResponse,
		options...,
	))

	return r

}

type errorWrapper struct {
	Error string `json:"error"`
}

type errorer interface {
	error() error
}

func errorDecoder(r *http.Response) error {
	var w errorWrapper
	if err := json.NewDecoder(r.Body).Decode(&w); err != nil {
		return err
	}
	return errors.New(w.Error)
}
