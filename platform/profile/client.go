package profile

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

	var applyProfileEndpoint endpoint.Endpoint
	{
		applyProfileEndpoint = httptransport.NewClient(
			"PUT",
			httputil.CopyURL(u, "/v1/profiles"),
			httputil.EncodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeApplyProfileResponse,
			opts...,
		).Endpoint()
	}

	var getProfilesEndpoint endpoint.Endpoint
	{
		getProfilesEndpoint = httptransport.NewClient(
			"GET",
			httputil.CopyURL(u, "/v1/profiles"),
			httputil.EncodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeGetProfilesResponse,
			opts...,
		).Endpoint()
	}

	var removeProfilesEndpoint endpoint.Endpoint
	{
		removeProfilesEndpoint = httptransport.NewClient(
			"DELETE",
			httputil.CopyURL(u, "/v1/profiles"),
			httputil.EncodeRequestWithToken(token, httptransport.EncodeJSONRequest),
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
