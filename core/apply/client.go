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
			DecodeBlueprintRequest,
			opts...,
		).Endpoint()
	}
	var applyDEPTokensEndpoint endpoint.Endpoint
	{
		applyDEPTokensEndpoint = httptransport.NewClient(
			"PUT",
			copyURL(u, "/v1/dep-tokens"),
			encodeRequestWithToken(token, EncodeHTTPGenericRequest),
			DecodeDEPTokensRequest,
			opts...,
		).Endpoint()
	}

	return Endpoints{
		ApplyBlueprintEndpoint: applyBlueprintEndpoint,
		ApplyDEPTokensEndpoint: applyDEPTokensEndpoint,
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
