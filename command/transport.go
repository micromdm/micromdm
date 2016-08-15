package command

import (
	"encoding/json"
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

// ServiceHandler returns an HTTP Handler for the command service
func ServiceHandler(ctx context.Context, svc Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	newCommandHandler := kithttp.NewServer(
		ctx,
		makeNewCommandEndpoint(svc),
		decodeNewCommandRequest,
		encodeResponse,
		opts...,
	)
	nextCommandHandler := kithttp.NewServer(
		ctx,
		makeNextCommandEndpoint(svc),
		decodeNextCommandRequest,
		encodeResponse,
		opts...,
	)
	deleteCommandHandler := kithttp.NewServer(
		ctx,
		makeDeleteCommandEndpoint(svc),
		decodeDeleteCommandRequest,
		encodeResponse,
		opts...,
	)
	getCommandsHandler := kithttp.NewServer(
		ctx,
		makeGetCommandsEndpoint(svc),
		decodeGetCommandsRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/mdm/commands/{udid}", getCommandsHandler).Methods("GET")
	r.Handle("/mdm/commands", newCommandHandler).Methods("POST")
	r.Handle("/mdm/commands/{udid}/next", nextCommandHandler).Methods("GET")
	r.Handle("/mdm/commands/{udid}/{uuid}", deleteCommandHandler).Methods("DELETE")

	return r
}

func decodeNewCommandRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request newCommandRequest
	err := json.NewDecoder(r.Body).Decode(&request.CommandRequest)
	return request, err
}

func decodeNextCommandRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	udid, ok := vars["udid"]
	if !ok {
		return nil, errBadRouting
	}
	var request nextCommandRequest
	request.UDID = udid
	return request, nil
}

func decodeDeleteCommandRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	udid, ok := vars["udid"]
	if !ok {
		return nil, errBadRouting
	}
	uuid, ok := vars["uuid"]
	if !ok {
		return nil, errBadRouting
	}
	var request deleteCommandRequest
	request.UDID = udid
	request.UUID = uuid
	return request, nil
}

func decodeGetCommandsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	udid, ok := vars["udid"]
	if !ok {
		return nil, errBadRouting
	}

	var request getCommandsRequest
	request.UDID = udid
	return request, nil
}

type errorer interface {
	error() error
}

type statuser interface {
	status() int
}

type listEncoder interface {
	encodeList(w http.ResponseWriter) error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// for success responses
	if e, ok := response.(statuser); ok {
		w.WriteHeader(e.status())
		if e.status() == http.StatusNoContent {
			return nil
		}
	}

	// check if this is a collection
	if e, ok := response.(listEncoder); ok {
		return e.encodeList(w)

	}
	return json.NewEncoder(w).Encode(response)
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	// unwrap if the error is wrapped by kit http in it's own error type
	if httperr, ok := err.(kithttp.Error); ok {
		err = httperr.Err
	}

	switch err {
	// case ErrNotFound:
	// 	w.WriteHeader(http.StatusNotFound)
	// case errEmptyRequest, errBadUUID:
	// 	w.WriteHeader(http.StatusBadRequest)
	// case workflow.ErrExists:
	// 	w.WriteHeader(http.StatusConflict)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
