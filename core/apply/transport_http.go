package apply

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
	BlueprintHandler        http.Handler
	DEPTokensHandler        http.Handler
	ProfileHandler          http.Handler
	DefineDEPProfileHandler http.Handler
	AppUploadHander         http.Handler
}

func MakeHTTPHandlers(ctx context.Context, endpoints Endpoints, opts ...httptransport.ServerOption) HTTPHandlers {
	h := HTTPHandlers{
		BlueprintHandler: httptransport.NewServer(
			endpoints.ApplyBlueprintEndpoint,
			decodeBlueprintRequest,
			encodeResponse,
			opts...,
		),
		DEPTokensHandler: httptransport.NewServer(
			endpoints.ApplyDEPTokensEndpoint,
			decodeDEPTokensRequest,
			encodeResponse,
			opts...,
		),
		ProfileHandler: httptransport.NewServer(
			endpoints.ApplyProfileEndpoint,
			decodeProfileRequest,
			encodeResponse,
			opts...,
		),
		DefineDEPProfileHandler: httptransport.NewServer(
			endpoints.DefineDEPProfileEndpoint,
			decodeDEPProfileRequest,
			encodeResponse,
			opts...,
		),
		AppUploadHander: httptransport.NewServer(
			endpoints.AppUploadEndpoint,
			decodeAppUploadRequest,
			encodeResponse,
			opts...,
		),
	}
	return h
}

func decodeDEPTokensRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req depTokensRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
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

func decodeDEPProfileRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req depProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeAppUploadRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	appManifestFilename := r.FormValue("app_manifest_filename")
	manifestFile, _, err := r.FormFile("app_manifest_filedata")
	if err != nil && err != http.ErrMissingFile {
		return nil, errors.Wrap(err, "manifest file")
	}
	pkgFilename := r.FormValue("pkg_name")
	pkgFile, _, err := r.FormFile("pkg_filedata")
	if err != nil && err != http.ErrMissingFile {
		return nil, err
	}

	return appUploadRequest{
		ManifestName: appManifestFilename,
		ManifestFile: manifestFile,
		PKGFilename:  pkgFilename,
		PKGFile:      pkgFile,
	}, nil
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

func DecodeDEPTokensResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errorDecoder(r)
	}
	var resp depTokensResponse
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

func DecodeDEPProfileResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errorDecoder(r)
	}
	var resp depProfileResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}
