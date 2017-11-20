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
	ListDevicesHandler         http.Handler
	GetDEPTokensHandler        http.Handler
	GetBlueprintsHandler       http.Handler
	GetProfilesHandler         http.Handler
	GetDEPAccountInfoHandler   http.Handler
	GetDEPProfileHandler       http.Handler
	GetDEPDeviceDetailsHandler http.Handler
	ListAppsHandler            http.Handler
	ListUsersHander            http.Handler
	BackupHandler              http.Handler
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
		GetProfilesHandler: httptransport.NewServer(
			endpoints.GetProfilesEndpoint,
			decodeGetProfilesRequest,
			encodeResponse,
			opts...),
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
		ListAppsHandler: httptransport.NewServer(
			endpoints.ListAppsEndpont,
			decodeListAppsRequest,
			encodeResponse,
			opts...,
		),
		ListUsersHander: httptransport.NewServer(
			endpoints.ListUserEndpoint,
			decodeListUsersRequest,
			encodeResponse,
			opts...,
		),
		BackupHandler: httptransport.NewServer(
			endpoints.BackupEndpoint,
			decodeBackupRequest,
			encodeResponse,
			opts...,
		),
	}
	return h
}

func decodeBackupRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
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

func decodeListUsersRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := userRequest{}
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

func decodeGetProfilesRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var opts GetProfilesOption
	if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
		return nil, err
	}
	req := profilesRequest{
		Opts: opts,
	}
	return req, nil
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

func decodeListAppsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request appListRequest
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

func DecodeGetProfilesResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errorDecoder(r)
	}
	var resp profilesResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
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

func DecodeListAppsResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errorDecoder(r)
	}
	var resp appListResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func DecodeListUsersResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errorDecoder(r)
	}
	var resp userResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, errors.Wrap(err, "decode user response")
}
