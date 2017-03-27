package list

import (
	"context"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewClient(instance string, logger log.Logger, token string) (Service, error) {
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}

	var listDevicesEndpoint endpoint.Endpoint
	{
		listDevicesEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/v1/devices"),
			encodeRequestWithToken(token, EncodeHTTPGenericRequest),
			DecodeDevicesResponse,
		).Endpoint()
	}

	return Endpoints{
		ListDevicesEndpoint: listDevicesEndpoint,
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
