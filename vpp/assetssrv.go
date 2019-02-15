package vpp

import (
	"github.com/pkg/errors"
)

// Contains information about the VPP Assets associated with a VPP account token
type AssetsSrv struct {
	TotalCount    int      `json:"totalCount"`
	Status        int      `json:"status"`
	Assets        []Asset  `json:"assets"`
	ClientContext string   `json:"clientContext"`
	UID           string   `json:"uId"`
	Location      Location `json:"location"`
}

// Contains information about VPP Assets
type Asset struct {
	ProductTypeID    int    `json:"productTypeId"`
	IsIrrevocable    bool   `json:"isIrrevocable"`
	PricingParam     string `json:"pricingParam"`
	AdamIDStr        string `json:"adamIdStr"`
	ProductTypeName  string `json:"productTypeName"`
	DeviceAssignable bool   `json:"deviceAssignable"`
	AssignedCount    int    `json:"assignedCount,omitempty"`
	AvailableCount   int    `json:"availableCount,omitempty"`
	RetiredCount     int    `json:"retiredCount,omitempty"`
	TotalCount       int    `json:"totalCount,omitempty"`
}

type AssetsSrvOptions struct {
	IncludeLicenseCounts bool   `json:"includeLicenseCounts,omitempty"`
	SToken               string `json:"sToken,omitempty"`
	PricingParam         string `json:"pricingParam,omitempty"`
	FacilitatorMemberID  string `json:"facilitatorMemberId,omitempty"`
}

// Gets information about the VPP Assets associated with a VPP Account token
func (c *Client) GetAssetsSrv(options AssetsSrvOptions) (*AssetsSrv, error) {

	if options.SToken == "" {
		options.SToken = c.VPPToken.SToken
	}

	// Get the VPPAssetsSrvURL
	VPPAssetsSrvURL := c.ServiceConfigSrv.GetVPPAssetsSrvURL

	// Create the VPPAssetsSrv request
	req, err := c.newRequest("POST", VPPAssetsSrvURL, options)
	if err != nil {
		return nil, errors.Wrap(err, "create VPPAssetsSrv request")
	}

	// Make the request
	var response AssetsSrv
	err = c.do(req, &response)

	return &response, errors.Wrap(err, "get VPPAssetsSrv request")
}

// Gets the pricing param for a particular VPP asset
func (c *Client) GetPricingParamForApp(appID string) (string, error) {
	// Get a list of assets
	response, err := c.GetAssetsSrv(AssetsSrvOptions{})
	if err != nil {
		return "", err
	}
	var assets = response.Assets

	// Find the pricing param for the asset with matching appId
	var pricing string
	for _, asset := range assets {
		if asset.AdamIDStr == appID {
			pricing = asset.PricingParam
			break
		}
	}

	// Check for err finding Pricing Param
	if pricing == "" {
		err := errors.New("Unable to find Pricing Param")
		return pricing, err
	}
	return pricing, nil
}
