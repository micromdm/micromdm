package apply

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

	var applyBlueprintEndpoint endpoint.Endpoint
	{
		applyBlueprintEndpoint = httptransport.NewClient(
			"PUT",
			copyURL(u, "/v1/blueprints"),
			encodeRequestWithToken(token, EncodeHTTPGenericRequest),
			DecodeBlueprintResponse,
			opts...,
		).Endpoint()
	}
	var applyDEPTokensEndpoint endpoint.Endpoint
	{
		applyDEPTokensEndpoint = httptransport.NewClient(
			"PUT",
			copyURL(u, "/v1/dep-tokens"),
			encodeRequestWithToken(token, EncodeHTTPGenericRequest),
			DecodeDEPTokensResponse,
			opts...,
		).Endpoint()
	}
	var applyProfileEndpoint endpoint.Endpoint
	{
		applyProfileEndpoint = httptransport.NewClient(
			"PUT",
			copyURL(u, "/v1/profiles"),
			encodeRequestWithToken(token, EncodeHTTPGenericRequest),
			DecodeProfileResponse,
			opts...,
		).Endpoint()
	}

	var defineDEPProfileEndpoint endpoint.Endpoint
	{
		defineDEPProfileEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/v1/dep/profiles"),
			encodeRequestWithToken(token, EncodeHTTPGenericRequest),
			DecodeDEPProfileResponse,
			opts...,
		).Endpoint()
	}

	var uploadAppEndpoint endpoint.Endpoint
	{
		uploadAppEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/v1/apps"),
			encodeRequestWithToken(token, EncodeUploadAppRequest),
			DecodeUploadAppResponse,
			opts...,
		).Endpoint()
	}

	var applyUserEndpoint endpoint.Endpoint
	{
		applyUserEndpoint = httptransport.NewClient(
			"PUT",
			copyURL(u, "/v1/users"),
			encodeRequestWithToken(token, EncodeHTTPGenericRequest),
			DecodeApplyUserResponse,
			opts...,
		).Endpoint()
	}

	var blockDeviceEndpoint endpoint.Endpoint
	{
		blockDeviceEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, ""), // empty path, modified by the encodeRequest func
			encodeRequestWithToken(token, encodeBlockDeviceRequest),
			DecodeBlockDeviceResponse,
			opts...,
		).Endpoint()
	}

	return Endpoints{
		ApplyBlueprintEndpoint:   applyBlueprintEndpoint,
		ApplyDEPTokensEndpoint:   applyDEPTokensEndpoint,
		ApplyProfileEndpoint:     applyProfileEndpoint,
		DefineDEPProfileEndpoint: defineDEPProfileEndpoint,
		AppUploadEndpoint:        uploadAppEndpoint,
		ApplyUserEndpoint:        applyUserEndpoint,
		BlockDeviceEndpoint:      blockDeviceEndpoint,
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
