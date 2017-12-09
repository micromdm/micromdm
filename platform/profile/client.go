package profile

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

	var applyProfileEndpoint endpoint.Endpoint
	{
		applyProfileEndpoint = httptransport.NewClient(
			"PUT",
			copyURL(u, "/v1/profiles"),
			encodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeApplyProfileResponse,
			opts...,
		).Endpoint()
	}

	var getProfilesEndpoint endpoint.Endpoint
	{
		getProfilesEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/v1/profiles"),
			encodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeGetProfilesResponse,
			opts...,
		).Endpoint()
	}

	var removeProfilesEndpoint endpoint.Endpoint
	{
		removeProfilesEndpoint = httptransport.NewClient(
			"DELETE",
			copyURL(u, "/v1/profiles"),
			encodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeRemoveProfileResponse,
			opts...,
		).Endpoint()
	}

	return Endpoints{
		ApplyProfileEndpoint:   applyProfileEndpoint,
		GetProfilesEndpoint:    getProfilesEndpoint,
		RemoveProfilesEndpoint: removeProfilesEndpoint,
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
