package blueprint

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

	var applyBlueprintEndpoint endpoint.Endpoint
	{
		applyBlueprintEndpoint = httptransport.NewClient(
			"PUT",
			httputil.CopyURL(u, "/v1/blueprints"),
			httputil.EncodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeApplyBlueprintResponse,
			opts...,
		).Endpoint()
	}

	var getBlueprintsEndpoint endpoint.Endpoint
	{
		getBlueprintsEndpoint = httptransport.NewClient(
			"GET",
			httputil.CopyURL(u, "/v1/blueprints"),
			httputil.EncodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeGetBlueprintsResponse,
			opts...,
		).Endpoint()
	}

	var removeBlueprintsEndpoint endpoint.Endpoint
	{
		removeBlueprintsEndpoint = httptransport.NewClient(
			"DELETE",
			httputil.CopyURL(u, "/v1/blueprints"),
			httputil.EncodeRequestWithToken(token, httptransport.EncodeJSONRequest),
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
