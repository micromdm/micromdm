package list

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
)

type HTTPHandlers struct {
	ListDevicesHandler   http.Handler
	GetDEPTokensHandler  http.Handler
	GetBlueprintsHandler http.Handler
}

func MakeHTTPHandlers(ctx context.Context, endpoints Endpoints, opts ...httptransport.ServerOption) HTTPHandlers {
	h := HTTPHandlers{
		ListDevicesHandler: httptransport.NewServer(
			endpoints.ListDevicesEndpoint,
			decodeListDevicesRequest,
			encodeResponse,
			opts...,
		),
		GetDEPTokensHandler: httptransport.NewServer(
			endpoints.GetDEPTokensEndpoint,
			decodeGetDEPTokensRequest,
			encodeResponse,
			opts...),
		GetBlueprintsHandler: httptransport.NewServer(
			endpoints.GetBlueprintsEndpoint,
			decodeGetBlueprintsRequest,
			encodeResponse,
			opts...),
	}
	return h
}

func decodeGetDEPTokensRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeListDevicesRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := devicesRequest{
		Opts: ListDevicesOption{},
	}
	return req, nil
}

func decodeGetBlueprintsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var opts GetBlueprintsOption
	if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
		return nil, err
	}
	req := blueprintsRequest{
		Opts: opts,
	}
	return req, nil
}

func errorDecoder(r *http.Response) error {
	var w errorWrapper
	if err := json.NewDecoder(r.Body).Decode(&w); err != nil {
		return err
	}
	return errors.New(w.Error)
}

type errorWrapper struct {
	Error string `json:"error"`
}

type errorer interface {
	error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		EncodeError(ctx, e.error(), w)
		return nil
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", " ")
	return enc.Encode(response)
}

func EncodeError(ctx context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	enc := json.NewEncoder(w)
	enc.SetIndent("", " ")
	enc.Encode(errorWrapper{Error: err.Error()})
}

// EncodeHTTPGenericRequest is a transport/http.EncodeRequestFunc that
// JSON-encodes any request to the request body. Primarily useful in a client.
func EncodeHTTPGenericRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func DecodeDevicesResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errorDecoder(r)
	}
	var resp devicesResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func DecodeGetDEPTokensResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errorDecoder(r)
	}
	var resp depTokenResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func DecodeGetBlueprintsResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errorDecoder(r)
	}
	var resp blueprintsResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}
