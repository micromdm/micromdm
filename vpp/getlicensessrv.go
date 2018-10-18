package vpp

import "github.com/pkg/errors"

type LicensesSrv struct {
	IfModifiedSinceMillisOrig string    `json:"ifModifiedSinceMillisOrig"`
	TotalCount                int       `json:"totalCount"`
	Status                    int       `json:"status"`
	TotalBatchCount           string    `json:"totalBatchCount"`
	Licenses                  []License `json:"licenses"`
	BatchToken                string    `json:"batchToken"`
	BatchCount                int       `json:"batchCount"`
	ClientContext             string    `json:"clientContext"`
	UId                       string    `json:"uId"`
	Location                  Location  `json:"location"`
}

type License struct {
	LicenseId       int    `json:"licenseId"`
	ProductTypeId   int    `json:"productTypeId"`
	IsIrrevocable   bool   `json:"isIrrevocable"`
	Status          string `json:"status"`
	PricingParam    string `json:"pricingParam"`
	AdamIdStr       string `json:"adamIdStr"`
	LicenseIdStr    string `json:"licenseIdStr"`
	ProductTypeName string `json:"productTypeName"`
	AdamId          int    `json:"adamId"`
	SerialNumber    string `json:"serialNumber"`
}

func (c *Client) GetLicensesSrv(serial string) (*LicensesSrv, error) {
	request := map[string]interface{}{
		"sToken": c.sToken,
	}
	if serial != "" {
		request["serialNumber"] = serial
	}

	licensesSrvUrl := c.VPPServiceConfigSrv.GetLicensesSrvUrl

	var response LicensesSrv

	req, err := c.newRequest("POST", licensesSrvUrl, request)
	if err != nil {
		return nil, errors.Wrap(err, "create LicensesSrv request")
	}

	err = c.do(req, &response)

	return &response, errors.Wrap(err, "LicensesSrv request")
}

func (c *Client) CheckAssignedLicense(serial string, appId string) (bool, error) {
	response, err := c.GetLicensesSrv(serial)
	if err != nil {
		return false, errors.Wrap(err, "create LicensesSrv request")
	}
	licenses := response.Licenses

	for i := 0; i < len(licenses); i++ {
		license := licenses[i]
		id := license.AdamIdStr
		if id == appId {
			return true, nil
		}
	}
	return false, nil
}
