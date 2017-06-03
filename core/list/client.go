package list

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

	var listDevicesEndpoint endpoint.Endpoint
	{
		listDevicesEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/v1/devices"),
			encodeRequestWithToken(token, EncodeHTTPGenericRequest),
			DecodeDevicesResponse,
			opts...,
		).Endpoint()
	}
	var getDEPTokensEndpoint endpoint.Endpoint
	{
		getDEPTokensEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/v1/dep-tokens"),
			encodeRequestWithToken(token, EncodeHTTPGenericRequest),
			DecodeGetDEPTokensResponse,
			opts...,
		).Endpoint()
	}
	var getBlueprintsEndpoint endpoint.Endpoint
	{
		getBlueprintsEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/v1/blueprints"),
			encodeRequestWithToken(token, EncodeHTTPGenericRequest),
			DecodeGetBlueprintsResponse,
			opts...,
		).Endpoint()
	}
	var getProfilesEndpoint endpoint.Endpoint
	{
		getProfilesEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/v1/profiles"),
			encodeRequestWithToken(token, EncodeHTTPGenericRequest),
			DecodeGetProfilesResponse,
			opts...,
		).Endpoint()
	}

	var getDEPAccountInfoEndpoint endpoint.Endpoint
	{
		getDEPAccountInfoEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/v1/dep/account"),
			encodeRequestWithToken(token, EncodeHTTPGenericRequest),
			DecodeDEPAccountInfoResponse,
			opts...,
		).Endpoint()
	}

	var getDEPDeviceDetailsEndpoint endpoint.Endpoint
	{
		getDEPDeviceDetailsEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/v1/dep/devices"),
			encodeRequestWithToken(token, EncodeHTTPGenericRequest),
			DecodeDEPDeviceDetailsReponse,
			opts...,
		).Endpoint()
	}

	var getDEPProfilesEndpoint endpoint.Endpoint
	{
		getDEPProfilesEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/v1/dep/profiles"),
			encodeRequestWithToken(token, EncodeHTTPGenericRequest),
			DecodeDEPProfileResponse,
			opts...,
		).Endpoint()
	}

	var listAppsEndpoint endpoint.Endpoint
	{
		listAppsEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/v1/apps"),
			encodeRequestWithToken(token, EncodeHTTPGenericRequest),
			DecodeListAppsResponse,
			opts...,
		).Endpoint()
	}

	return Endpoints{
		ListDevicesEndpoint:       listDevicesEndpoint,
		GetDEPTokensEndpoint:      getDEPTokensEndpoint,
		GetBlueprintsEndpoint:     getBlueprintsEndpoint,
		GetProfilesEndpoint:       getProfilesEndpoint,
		GetDEPAccountInfoEndpoint: getDEPAccountInfoEndpoint,
		GetDEPDeviceEndpoint:      getDEPDeviceDetailsEndpoint,
		GetDEPProfileEndpoint:     getDEPProfilesEndpoint,
		ListAppsEndpont:           listAppsEndpoint,
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
