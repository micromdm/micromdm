package vpp

import "github.com/pkg/errors"

// Contains information about a managed license
type ManageVPPLicensesByAdamIdSrv struct {
	ProductTypeID   int           `json:"productTypeId,omitempty"`
	ProductTypeName string        `json:"productTypeName,omitempty"`
	IsIrrevocable   bool          `json:"isIrrevocable,omitempty"`
	PricingParam    string        `json:"pricingParam,omitempty"`
	UID             string        `json:"uId,omitempty,omitempty"`
	AdamIdStr       string        `json:"adamIdStr,omitempty"`
	Status          int           `json:"status"`
	ClientContext   string        `json:"clientContext,omitempty"`
	Location        *Location     `json:"location,omitempty"`
	Associations    []Association `json:"associations,omitempty"`
	ErrorMessage    string        `json:"errorMessage,omitempty"`
	ErrorNumber     int           `json:"errorNumber,omitempty"`
}

// Contains information about an app association
type Association struct {
	SerialNumber           string   `json:"serialNumber"`
	ErrorMessage           string   `json:"errorMessage,omitempty"`
	ErrorCode              int      `json:"errorCode,omitempty"`
	ErrorNumber            int      `json:"errorNumber,omitempty"`
	LicenseIDStr           string   `json:"licenseIdStr,omitempty"`
	LicenseAlreadyAssigned *License `json:"licenseAlreadyAssigned,omitempty"`
}

// Contains options to pass to the ManageVPPLicensesByAdamIdSrv
type ManageVPPLicensesByAdamIdSrvOptions struct {
	SToken                    string   `json:"sToken"`
	AdamIDStr                 string   `json:"adamIdStr"`
	PricingParam              string   `json:"pricingParam"`
	AssociateSerialNumbers    []string `json:"associateSerialNumbers,omitempty"`
	DisassociateSerialNumbers []string `json:"disassociateSerialNumbers,omitempty"`
}

// Associates a list of serials to a VPP app license
func (c *Client) AssociateSerialsToApp(appID string, serials []string) (*ManageVPPLicensesByAdamIdSrv, error) {
	options := ManageVPPLicensesByAdamIdSrvOptions{
		AssociateSerialNumbers: serials,
	}

	response, err := c.ManageVPPLicensesByAdamIdSrv(appID, options)
	return &response, err
}

// Disssociates a list of serials to a VPP app license
func (c *Client) DisassociateSerialsToApp(appID string, serials []string) (*ManageVPPLicensesByAdamIdSrv, error) {
	options := ManageVPPLicensesByAdamIdSrvOptions{
		DisassociateSerialNumbers: serials,
	}

	response, err := c.ManageVPPLicensesByAdamIdSrv(appID, options)
	return &response, err
}

// Interfaces with the ManageVPPLicensesByAdamIdSrv to managed VPP licenses
func (c *Client) ManageVPPLicensesByAdamIdSrv(appID string, options ManageVPPLicensesByAdamIdSrvOptions) (ManageVPPLicensesByAdamIdSrv, error) {
	options.SToken = c.SToken
	options.AdamIDStr = appID

	// Get the pricing param required to manage a vpp license
	pricing, err := c.GetPricingParamForApp(appID)
	if err != nil {
		return ManageVPPLicensesByAdamIdSrv{}, errors.Wrap(err, "get PricingParam request")
	}
	options.PricingParam = pricing

	// Get the ManageVPPLicensesByAdamIdSrvURL
	manageVPPLicensesByAdamIdSrvUrl := c.VPPServiceConfigSrv.ManageVPPLicensesByAdamIdSrvURL

	// Create the ManageVPPLicensesByAdamIdSrv request
	req, err := c.newRequest("POST", manageVPPLicensesByAdamIdSrvUrl, options)
	if err != nil {
		return ManageVPPLicensesByAdamIdSrv{}, errors.Wrap(err, "create ManageVPPLicensesByAdamIdSrv request")
	}

	// Make the Request
	var response ManageVPPLicensesByAdamIdSrv
	err = c.do(req, &response)

	return response, errors.Wrap(err, "ManageVPPLicensesByAdamIdSrv request")
}
