package connect

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/groob/plist"
	"github.com/micromdm/mdm"
)

type HTTPHandlers struct {
	ConnectHandler http.Handler
}

func MakeHTTPHandlers(ctx context.Context, endpoints Endpoints, opts ...httptransport.ServerOption) HTTPHandlers {
	h := HTTPHandlers{
		ConnectHandler: httptransport.NewServer(
			endpoints.ConnectEndpoint,
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
	var res mdm.Response

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	err = plist.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	req := MDMConnectRequest{MDMResponse: res, Raw: body}
	return req, nil
}

// According to the MDM Check-in protocol, the server must respond with 200 OK
// to successful Check-in requests.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		EncodeError(ctx, e.error(), w)
		return nil
	}

	resp := response.(mdmConnectResponse)

	w.WriteHeader(http.StatusOK)
	w.Write(resp.payload)
	return nil
}

// EncodeError is used by the HTTP transport to encode service errors in HTTP.
// The EncodeError should be passed to the Go-Kit httptransport as the
// ServerErrorEncoder to encode error responses.
func EncodeError(ctx context.Context, err error, w http.ResponseWriter) {
	fmt.Printf("connect error: %s\n", err)
	w.WriteHeader(http.StatusInternalServerError)
}
