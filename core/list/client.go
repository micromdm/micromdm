package list

import (
	"net/url"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewClient(instance string, logger log.Logger) (Service, error) {
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}

	var listDevicesEndpoint endpoint.Endpoint
	{
		listDevicesEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/v1/devices"),
			EncodeHTTPGenericRequest,
			DecodeDevicesResponse,
		).Endpoint()
	}

	return Endpoints{
		ListDevicesEndpoint: listDevicesEndpoint,
	}, nil
}

func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}
