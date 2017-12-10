package list

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"
)

type HTTPHandlers struct {
	GetDEPAccountInfoHandler   http.Handler
	GetDEPProfileHandler       http.Handler
	GetDEPDeviceDetailsHandler http.Handler
}

func MakeHTTPHandlers(ctx context.Context, endpoints Endpoints, opts ...httptransport.ServerOption) HTTPHandlers {
	h := HTTPHandlers{
		GetDEPAccountInfoHandler: httptransport.NewServer(
			endpoints.GetDEPAccountInfoEndpoint,
			decodeDepAccountInfoRequest,
			encodeResponse,
			opts...,
		),
		GetDEPDeviceDetailsHandler: httptransport.NewServer(
			endpoints.GetDEPDeviceEndpoint,
			decodeDepDeviceDetailsRequest,
			encodeResponse,
			opts...,
		),
		GetDEPProfileHandler: httptransport.NewServer(
			endpoints.GetDEPProfileEndpoint,
			decodeDEPProfileRequest,
			encodeResponse,
			opts...,
		),
	}
	return h
}

func decodeDepAccountInfoRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeDepDeviceDetailsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request depDeviceDetailsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeDEPProfileRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request depProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
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

func DecodeDEPAccountInfoResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errorDecoder(r)
	}
	var resp depAccountInfoResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func DecodeDEPDeviceDetailsReponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errorDecoder(r)
	}
	var resp depDeviceDetailsResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func DecodeDEPProfileResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errorDecoder(r)
	}
	var resp depProfileResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}
