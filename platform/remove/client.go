package remove

import (
	"net/url"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/micromdm/micromdm/pkg/httputil"
)

func NewHTTPClient(instance, token string, logger log.Logger, opts ...httptransport.ClientOption) (Service, error) {
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}

	var blockDeviceEndpoint endpoint.Endpoint
	{
		blockDeviceEndpoint = httptransport.NewClient(
			"POST",
			httputil.CopyURL(u, ""), // empty path, modified by the encodeRequest func
			httputil.EncodeRequestWithToken(token, encodeBlockDeviceRequest),
			decodeBlockDeviceResponse,
			opts...,
		).Endpoint()
	}

	var unblockDeviceEndpoint endpoint.Endpoint
	{
		unblockDeviceEndpoint = httptransport.NewClient(
			"POST",
			httputil.CopyURL(u, ""), //modified by encodeRequestFunc
			httputil.EncodeRequestWithToken(token, encodeUnblockDeviceRequest),
			decodeUnblockDeviceResponse,
			opts...,
		).Endpoint()
	}

	return Endpoints{
		BlockDeviceEndpoint:   blockDeviceEndpoint,
		UnblockDeviceEndpoint: unblockDeviceEndpoint,
	}, nil
}
