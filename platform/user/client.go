package user

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

	var applyUserEndpoint endpoint.Endpoint
	{
		applyUserEndpoint = httptransport.NewClient(
			"PUT",
			httputil.CopyURL(u, "/v1/users"),
			httputil.EncodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeApplyUserResponse,
			opts...,
		).Endpoint()
	}

	var listUsersEndpoint endpoint.Endpoint
	{
		listUsersEndpoint = httptransport.NewClient(
			"POST",
			httputil.CopyURL(u, "/v1/users"),
			httputil.EncodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeListUsersResponse,
			opts...,
		).Endpoint()
	}

	return Endpoints{
		ApplyUserEndpoint: applyUserEndpoint,
		ListUsersEndpoint: listUsersEndpoint,
	}, nil
}
