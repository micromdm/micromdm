package remove

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
	BlueprintHandler http.Handler
	ProfileHandler   http.Handler
}

func MakeHTTPHandlers(ctx context.Context, endpoint Endpoints, opts ...httptransport.ServerOption) HTTPHandlers {
	h := HTTPHandlers{
		BlueprintHandler: httptransport.NewServer(
			endpoint.RemoveBlueprintsEndpoint,
			decodeBlueprintRequest,
			encodeResponse,
			opts...,
		),
		ProfileHandler: httptransport.NewServer(
			endpoint.RemoveProfilesEndpoint,
			decodeProfileRequest,
			encodeResponse,
			opts...,
		),
	}
	return h
}

func decodeBlueprintRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var bpReq blueprintRequest
	if err := json.NewDecoder(r.Body).Decode(&bpReq); err != nil {
		return nil, err
	}
	return bpReq, nil
}

func decodeProfileRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req profileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
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

func DecodeBlueprintResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errorDecoder(r)
	}
	var resp blueprintResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func DecodeProfileResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errorDecoder(r)
	}
	var resp profileResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}
