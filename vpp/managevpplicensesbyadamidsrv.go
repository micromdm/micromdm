package vpp

import (
	"github.com/pkg/errors"
)

// Contains information about a managed license
type ManageVPPLicensesByAdamIdSrv struct {
	ProductTypeID   int              `json:"productTypeId,omitempty"`
	ProductTypeName string           `json:"productTypeName,omitempty"`
	IsIrrevocable   bool             `json:"isIrrevocable,omitempty"`
	PricingParam    string           `json:"pricingParam,omitempty"`
	UID             string           `json:"uId,omitempty,omitempty"`
	AdamIdStr       string           `json:"adamIdStr,omitempty"`
	Status          int              `json:"status"`
	ClientContext   string           `json:"clientContext,omitempty"`
	Location        *Location        `json:"location,omitempty"`
	Associations    []Association    `json:"associations,omitempty"`
	Disassociations []Disassociation `json:"disassociations,omitempty"`
	ErrorMessage    string           `json:"errorMessage,omitempty"`
	ErrorNumber     int              `json:"errorNumber,omitempty"`
}

// Contains information about an app association
type Association struct {
	SerialNumber           string   `json:"serialNumber,omitempty"`
	ErrorMessage           string   `json:"errorMessage,omitempty"`
	ErrorCode              int      `json:"errorCode,omitempty"`
	ErrorNumber            int      `json:"errorNumber,omitempty"`
	LicenseIDStr           string   `json:"licenseIdStr,omitempty"`
	LicenseAlreadyAssigned *License `json:"licenseAlreadyAssigned,omitempty"`
}

// Contains information about an app disassociation
type Disassociation struct {
	SerialNumber string `json:"serialNumber,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	ErrorCode    int    `json:"errorCode,omitempty"`
	ErrorNumber  int    `json:"errorNumber,omitempty"`
	LicenseIDStr string `json:"licenseIdStr,omitempty"`
}

// Contains options to pass to the ManageVPPLicensesByAdamIdSrv
type ManageVPPLicensesByAdamIdSrvOptions struct {
	AdamIDStr                    string   `json:"adamIdStr"`
	PricingParam                 string   `json:"pricingParam"`
	AssociateClientUserIdStrs    string   `json:"associateClientUserIdStrs,omitempty"`
	AssociateSerialNumbers       []string `json:"associateSerialNumbers,omitempty"`
	DisassociateClientUserIdStrs string   `json:"disassociateClientUserIdStrs,omitempty"`
	DisassociateLicenseIdStrs    string   `json:"disassociateLicenseIdStrs,omitempty"`
	DisassociateSerialNumbers    []string `json:"disassociateSerialNumbers,omitempty"`
	notifyDisassociation         bool     `json:"notifyDisassociation,omitempty"`
	SToken                       string   `json:"sToken,omitempty"`
	facilitatorMemberId          string   `json:"facilitatorMemberId,omitempty"`
}

// Interfaces with the ManageVPPLicensesByAdamIdSrv to managed VPP licenses
func (c *Client) ManageVPPLicensesByAdamIdSrv(options ManageVPPLicensesByAdamIdSrvOptions) (*ManageVPPLicensesByAdamIdSrv, error) {

	if options.SToken == "" {
		options.SToken = c.VPPToken.SToken
	}

	if options.AdamIDStr == "" {
		return nil, errors.Wrap(nil, "must include adamIdStr")
	}

	// Get the pricing param required to manage a vpp license
	if options.PricingParam == "" {
		pricing, err := c.GetPricingParamForApp(options.AdamIDStr)
		if err != nil {
			return nil, errors.Wrap(err, "get PricingParam request")
		}
		options.PricingParam = pricing
	}

	// Get the ManageVPPLicensesByAdamIdSrvURL
	manageVPPLicensesByAdamIdSrvUrl := c.ServiceConfigSrv.ManageVPPLicensesByAdamIdSrvURL

	// Create the ManageVPPLicensesByAdamIdSrv request
	req, err := c.newRequest("POST", manageVPPLicensesByAdamIdSrvUrl, options)
	if err != nil {
		return nil, errors.Wrap(err, "create ManageVPPLicensesByAdamIdSrv request")
	}

	// Make the Request
	var response ManageVPPLicensesByAdamIdSrv
	err = c.do(req, &response)

	return &response, errors.Wrap(err, "ManageVPPLicensesByAdamIdSrv request")
}
