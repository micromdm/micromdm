package checkin

import (
	"context"
	"io"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/groob/plist"
)

type HTTPHandlers struct {
	// The CheckinHandler should accept PUT requests
	CheckinHandler http.Handler
}

func MakeHTTPHandlers(ctx context.Context, endpoints Endpoints, opts ...httptransport.ServerOption) HTTPHandlers {
	h := HTTPHandlers{
		CheckinHandler: httptransport.NewServer(
			endpoints.CheckinEndpoint,
			decodeRequest,
			encodeResponse,
			opts...,
		),
	}
	return h
}

type errorer interface {
	error() error
}

func decodeRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req checkinRequest
	err := plist.NewDecoder(io.LimitReader(r.Body, 10000)).Decode(&req)
	return req, err
}

// According to the MDM Check-in protocol, the server must respond with 200 OK
// to successful Check-in requests.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		EncodeError(ctx, e.error(), w)
		return nil
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

// EncodeError is used by the HTTP transport to encode service errors in HTTP.
// The EncodeError should be passed to the Go-Kit httptransport as the
// ServerErrorEncoder to encode error responses.
// According to the MDM Check-in protocol specification, the device only needs
// a 401 (Unauthorized) response in case of failure.
func EncodeError(ctx context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
}
