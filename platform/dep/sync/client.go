package sync

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

	var syncNowEndpoint endpoint.Endpoint
	{
		syncNowEndpoint = httptransport.NewClient(
			"PUT",
			httputil.CopyURL(u, "/v1/dep/autoassign"),
			httputil.EncodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeEmptyResponse,
			opts...,
		).Endpoint()
	}

	var applyAutoAssignerEndpoint endpoint.Endpoint
	{
		applyAutoAssignerEndpoint = httptransport.NewClient(
			"POST",
			httputil.CopyURL(u, "/v1/dep/autoassigners"),
			httputil.EncodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeApplyAutoAssignerResponse,
			opts...,
		).Endpoint()
	}

	var getAutoAssignersEndpoint endpoint.Endpoint
	{
		getAutoAssignersEndpoint = httptransport.NewClient(
			"GET",
			httputil.CopyURL(u, "/v1/dep/autoassigners"),
			httputil.EncodeRequestWithToken(token, httputil.EncodeEmptyRequest),
			decodeGetAutoAssignersResponse,
			opts...,
		).Endpoint()
	}

	var removeAutoAssignerEndpoint endpoint.Endpoint
	{
		removeAutoAssignerEndpoint = httptransport.NewClient(
			"DELETE",
			httputil.CopyURL(u, "/v1/dep/autoassigners"),
			httputil.EncodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeRemoveAutoAssignerResponse,
			opts...,
		).Endpoint()
	}

	return Endpoints{
		SyncNowEndpoint:            syncNowEndpoint,
		ApplyAutoAssignerEndpoint:  applyAutoAssignerEndpoint,
		GetAutoAssignersEndpoint:   getAutoAssignersEndpoint,
		RemoveAutoAssignerEndpoint: removeAutoAssignerEndpoint,
	}, nil
}
