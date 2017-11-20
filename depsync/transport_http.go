package depsync

import (
	"context"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
)

type HTTPHandlers struct {
	RefreshHandler http.Handler
}

func MakeHTTPHandlers(ctx context.Context, endpoints Endpoints, opts ...httptransport.ServerOption) HTTPHandlers {
	return HTTPHandlers{
		RefreshHandler: httptransport.NewServer(
			endpoints.RefreshEndpoint,
			decodeEmptyRequest,
			encodeEmptyResponse,
			opts...,
		),
	}
}

func decodeEmptyRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func encodeEmptyResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return nil
}
