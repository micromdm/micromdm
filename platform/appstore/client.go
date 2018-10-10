package appstore

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

	var appUploadEndpoint endpoint.Endpoint
	{
		appUploadEndpoint = httptransport.NewClient(
			"PUT",
			httputil.CopyURL(u, "/v1/apps"),
			httputil.EncodeRequestWithToken(token, encodeUploadAppRequest),
			decodeUploadAppResponse,
			opts...,
		).Endpoint()
	}

	var listAppsEndpoint endpoint.Endpoint
	{
		listAppsEndpoint = httptransport.NewClient(
			"POST",
			httputil.CopyURL(u, "/v1/apps"),
			httputil.EncodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeListAppsResponse,
			opts...,
		).Endpoint()
	}

	return Endpoints{
		AppUploadEndpoint: appUploadEndpoint,
		ListAppsEndpoint:  listAppsEndpoint,
	}, nil
}
