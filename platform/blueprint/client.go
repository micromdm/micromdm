package blueprint

import (
	"context"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHTTPClient(instance, token string, logger log.Logger, opts ...httptransport.ClientOption) (Service, error) {
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}

	var applyBlueprintEndpoint endpoint.Endpoint
	{
		applyBlueprintEndpoint = httptransport.NewClient(
			"PUT",
			copyURL(u, "/v1/blueprints"),
			encodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeApplyBlueprintResponse,
			opts...,
		).Endpoint()
	}

	var getBlueprintsEndpoint endpoint.Endpoint
	{
		getBlueprintsEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/v1/blueprints"),
			encodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeGetBlueprintsResponse,
			opts...,
		).Endpoint()
	}

	var removeBlueprintsEndpoint endpoint.Endpoint
	{
		removeBlueprintsEndpoint = httptransport.NewClient(
			"DELETE",
			copyURL(u, "/v1/blueprints"),
			encodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeRemoveBlueprintsResponse,
			opts...,
		).Endpoint()
	}

	return Endpoints{
		ApplyBlueprintEndpoint:   applyBlueprintEndpoint,
		GetBlueprintsEndpoint:    getBlueprintsEndpoint,
		RemoveBlueprintsEndpoint: removeBlueprintsEndpoint,
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
