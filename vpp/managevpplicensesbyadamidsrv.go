package vpp

import "github.com/pkg/errors"

type ManageVPPLicensesByAdamIdSrv struct {
	ProductTypeId   int           `json:"productTypeId,omitempty"`
	ProductTypeName string        `json:"productTypeName,omitempty"`
	IsIrrevocable   bool          `json:"isIrrevocable,omitempty"`
	PricingParam    string        `json:"pricingParam,omitempty"`
	UId             string        `json:"uId,omitempty,omitempty"`
	AdamIdStr       string        `json:"adamIdStr,omitempty"`
	Status          int           `json:"status"`
	ClientContext   string        `json:"clientContext,omitempty"`
	Location        *Location     `json:"location,omitempty"`
	Associations    []Association `json:"associations,omitempty"`
	ErrorMessage    string        `json:"errorMessage,omitempty"`
	ErrorNumber     int           `json:"errorNumber,omitempty"`
}

type Association struct {
	SerialNumber           string   `json:"serialNumber"`
	ErrorMessage           string   `json:"errorMessage,omitempty"`
	ErrorCode              int      `json:"errorCode,omitempty"`
	ErrorNumber            int      `json:"errorNumber,omitempty"`
	LicenseIdStr           string   `json:"licenseIdStr,omitempty"`
	LicenseAlreadyAssigned *License `json:"licenseAlreadyAssigned,omitempty"`
}

func (c *Client) AssociateSerialsToApp(appId string, serials []string) (*ManageVPPLicensesByAdamIdSrv, error) {
	//func (c *Client) AssociateSerialsToApp(appId string, serials []string) (interface{}, error) {
	options := map[string]interface{}{
		"associateSerialNumbers": serials,
	}

	response, err := c.ManageVPPLicensesByAdamIdSrv(appId, options)
	return &response, err
	//return response, errors.Wrap(err, "ManageVPPLicensesByAdamIdSrv request")
}

func (c *Client) DisassociateSerialsToApp(appId string, serials []string) (*ManageVPPLicensesByAdamIdSrv, error) {
	//func (c *Client) AssociateSerialsToApp(appId string, serials []string) (interface{}, error) {
	options := map[string]interface{}{
		"disassociateSerialNumbers": serials,
	}

	response, err := c.ManageVPPLicensesByAdamIdSrv(appId, options)
	return &response, err
	//return response, errors.Wrap(err, "ManageVPPLicensesByAdamIdSrv request")
}

func (c *Client) ManageVPPLicensesByAdamIdSrv(appId string, options map[string]interface{}) (ManageVPPLicensesByAdamIdSrv, error) {
	//func (c *Client) ManageVPPLicensesByAdamIdSrv(appId string, serials []string) (interface{}, error) {
	pricing, err := c.GetPricingParamForApp(appId)
	if err != nil {
		return ManageVPPLicensesByAdamIdSrv{}, errors.Wrap(err, "get PricingParam request")
	}

	request := map[string]interface{}{
		"sToken":       c.sToken,
		"adamIdStr":    appId,
		"pricingParam": pricing,
	}
	for k, v := range options {
		request[k] = v
	}

	manageVPPLicensesByAdamIdSrvUrl := c.VPPServiceConfigSrv.ManageVPPLicensesByAdamIdSrvUrl

	var response ManageVPPLicensesByAdamIdSrv
	//var response interface{}

	req, err := c.newRequest("POST", manageVPPLicensesByAdamIdSrvUrl, request)
	if err != nil {
		return ManageVPPLicensesByAdamIdSrv{}, errors.Wrap(err, "create ManageVPPLicensesByAdamIdSrv request")
	}

	err = c.do(req, &response)

	return response, errors.Wrap(err, "ManageVPPLicensesByAdamIdSrv request")
	//return response, errors.Wrap(err, "ManageVPPLicensesByAdamIdSrv request")
}
