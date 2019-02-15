package vpp

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

	var getVPPAppsEndpoint endpoint.Endpoint
	{
		getVPPAppsEndpoint = httptransport.NewClient(
			"POST",
			httputil.CopyURL(u, "/v1/vpp/apps"),
			httputil.EncodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeGetVPPAppsResponse,
			opts...,
		).Endpoint()
	}

	var getContentMetadataEndpoint endpoint.Endpoint
	{
		getContentMetadataEndpoint = httptransport.NewClient(
			"POST",
			httputil.CopyURL(u, "/v1/vpp/metadata"),
			httputil.EncodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeGetContentMetadataResponse,
			opts...,
		).Endpoint()
	}

	var getAssetsSrvEndpoint endpoint.Endpoint
	{
		getAssetsSrvEndpoint = httptransport.NewClient(
			"POST",
			httputil.CopyURL(u, "/v1/vpp/assets"),
			httputil.EncodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeGetAssetsSrvResponse,
			opts...,
		).Endpoint()
	}

	var getLicensesSrvEndpoint endpoint.Endpoint
	{
		getLicensesSrvEndpoint = httptransport.NewClient(
			"POST",
			httputil.CopyURL(u, "/v1/vpp/licenses"),
			httputil.EncodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeGetLicensesSrvResponse,
			opts...,
		).Endpoint()
	}

	var manageVPPLicensesByAdamIdSrvEndpoint endpoint.Endpoint
	{
		manageVPPLicensesByAdamIdSrvEndpoint = httptransport.NewClient(
			"PUT",
			httputil.CopyURL(u, "/v1/vpp/licenses"),
			httputil.EncodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeManageVPPLicensesByAdamIdSrvResponse,
			opts...,
		).Endpoint()
	}

	var getServiceConfigSrvEndpoint endpoint.Endpoint
	{
		getServiceConfigSrvEndpoint = httptransport.NewClient(
			"POST",
			httputil.CopyURL(u, "/v1/vpp/serviceconfig"),
			httputil.EncodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeGetServiceConfigSrvResponse,
			opts...,
		).Endpoint()
	}

	return Endpoints{
		GetVPPAppsEndpoint:                   getVPPAppsEndpoint,
		GetContentMetadataEndpoint:           getContentMetadataEndpoint,
		GetAssetsSrvEndpoint:                 getAssetsSrvEndpoint,
		GetLicensesSrvEndpoint:               getLicensesSrvEndpoint,
		ManageVPPLicensesByAdamIdSrvEndpoint: manageVPPLicensesByAdamIdSrvEndpoint,
		GetServiceConfigSrvEndpoint:          getServiceConfigSrvEndpoint,
	}, nil
}
