package remove

import (
	"context"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewClient(instance string, logger log.Logger, token string, opts ...httptransport.ClientOption) (Service, error) {
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}

	var unblockDeviceEndpoint endpoint.Endpoint
	{
		unblockDeviceEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, ""), //modified by encodeRequestFunc
			encodeRequestWithToken(token, encodeUnblockDeviceRequest),
			DecodeUnblockDeviceResponse,
			opts...,
		).Endpoint()
	}

	return Endpoints{
		UnblockDeviceEndpoint: unblockDeviceEndpoint,
	}, nil
}

func encodeRequestWithToken(token string, next httptransport.EncodeRequestFunc) httptransport.EncodeRequestFunc {
	return func(ctx context.Context, r *http.Request, request interface{}) error {
		r.SetBasicAuth("micromdm", token)
		return next(ctx, r, request)
	}
}

func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}
