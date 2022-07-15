package dep

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

	var defineProfileEndpoint endpoint.Endpoint
	{
		defineProfileEndpoint = httptransport.NewClient(
			"PUT",
			httputil.CopyURL(u, "/v1/dep/profiles"),
			httputil.EncodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeDefineProfileResponse,
			opts...,
		).Endpoint()
	}

	var assignProfileEndpoint endpoint.Endpoint
	{
		assignProfileEndpoint = httptransport.NewClient(
			"POST",
			httputil.CopyURL(u, "/v1/dep/assign"),
			httputil.EncodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeAssignProfileResponse,
			opts...,
		).Endpoint()
	}

	var removeProfileEndpoint endpoint.Endpoint
	{
		removeProfileEndpoint = httptransport.NewClient(
			"DELETE",
			httputil.CopyURL(u, "/v1/dep/profiles"),
			httputil.EncodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeRemoveProfileResponse,
			opts...,
		).Endpoint()
	}

	var fetchProfileEndpoint endpoint.Endpoint
	{
		fetchProfileEndpoint = httptransport.NewClient(
			"POST",
			httputil.CopyURL(u, "/v1/dep/profiles"),
			httputil.EncodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeFetchProfileResponse,
			opts...,
		).Endpoint()
	}

	var getAccountInfoEndpoint endpoint.Endpoint
	{
		getAccountInfoEndpoint = httptransport.NewClient(
			"GET",
			httputil.CopyURL(u, "/v1/dep/account"),
			httputil.EncodeRequestWithToken(token, httputil.EncodeEmptyRequest),
			decodeGetAccountInfoResponse,
			opts...,
		).Endpoint()
	}

	var getDeviceDetailsEndpoint endpoint.Endpoint
	{
		getDeviceDetailsEndpoint = httptransport.NewClient(
			"POST",
			httputil.CopyURL(u, "/v1/dep/devices"),
			httputil.EncodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeDeviceDetailsResponse,
			opts...,
		).Endpoint()
	}

	return Endpoints{
		DefineProfileEndpoint:    defineProfileEndpoint,
		AssignProfileEndpoint:    assignProfileEndpoint,
		RemoveProfileEndpoint:    removeProfileEndpoint,
		FetchProfileEndpoint:     fetchProfileEndpoint,
		GetAccountInfoEndpoint:   getAccountInfoEndpoint,
		GetDeviceDetailsEndpoint: getDeviceDetailsEndpoint,
	}, nil
}
