package vpp

import "github.com/pkg/errors"

type VPPAssetsSrv struct {
	TotalCount    int      `json:"totalCount"`
	Status        int      `json:"status"`
	Assets        []Asset  `json:"assets"`
	ClientContext string   `json:"clientContext"`
	UId           string   `json:"uId"`
	Location      Location `json:"location"`
}

type Asset struct {
	ProductTypeId    int    `json:"productTypeId"`
	IsIrrevocable    bool   `json:"isIrrevocable"`
	PricingParam     string `json:"pricingParam"`
	AdamIdStr        string `json:"adamIdStr"`
	ProductTypeName  string `json:"productTypeName"`
	DeviceAssignable bool   `json:"deviceAssignable"`
}

func (c *Client) GetVppAssetsSrv() (*VPPAssetsSrv, error) {
	request := map[string]interface{}{
		"sToken": c.sToken,
	}

	vPPAssetsSrvUrl := c.VPPServiceConfigSrv.GetVPPAssetsSrvUrl

	var response VPPAssetsSrv

	req, err := c.newRequest("POST", vPPAssetsSrvUrl, request)
	if err != nil {
		return nil, errors.Wrap(err, "create VPPAssetsSrv request")
	}

	err = c.do(req, &response)

	return &response, errors.Wrap(err, "VPPAssetsSrv request")
}

func (c *Client) GetPricingParamForApp(appId string) (string, error) {
	response, err := c.GetVppAssetsSrv()
	if err != nil {
		return "", errors.Wrap(err, "create VppAssetsSrv request")
	}
	var assets = response.Assets

	var pricing string
	for i := 0; i < len(assets); i++ {
		asset := assets[i]
		if asset.AdamIdStr == appId {
			pricing = asset.PricingParam
		}
	}
	return pricing, nil
}
