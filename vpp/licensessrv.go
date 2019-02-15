package vpp

import (
	"github.com/pkg/errors"
)

// Contains information about the VPP Licenses associated with a VPP account token
type LicensesSrv struct {
	IfModifiedSinceMillisOrig string    `json:"ifModifiedSinceMillisOrig"`
	TotalCount                int       `json:"totalCount"`
	Status                    int       `json:"status"`
	TotalBatchCount           string    `json:"totalBatchCount"`
	Licenses                  []License `json:"licenses"`
	BatchToken                string    `json:"batchToken"`
	BatchCount                int       `json:"batchCount"`
	ClientContext             string    `json:"clientContext"`
	UID                       string    `json:"uId"`
	Location                  Location  `json:"location"`
}

// Contains information about VPP Licenses
type License struct {
	LicenseID       int    `json:"licenseId"`
	ProductTypeID   int    `json:"productTypeId"`
	IsIrrevocable   bool   `json:"isIrrevocable"`
	Status          string `json:"status"`
	PricingParam    string `json:"pricingParam"`
	AdamIDStr       string `json:"adamIdStr"`
	LicenseIDStr    string `json:"licenseIdStr"`
	ProductTypeName string `json:"productTypeName"`
	AdamID          int    `json:"adamId"`
	SerialNumber    string `json:"serialNumber"`
}

// Contains LicensesSrv request options
type LicensesSrvOptions struct {
	BatchToken          string `json:"batchToken,omitempty"`
	SinceModifiedToken  string `json:"sinceModifiedToken,omitempty"`
	AdamID              string `json:"adamId,omitempty"`
	SToken              string `json:"sToken,omitempty"`
	FacilitatorMemberID string `json:"facilitatorMemberId,omitempty"`
	AssignedOnly        bool   `json:"assignedOnly,omitempty"`
	PricingParam        string `json:"pricingParam,omitempty"`
	SerialNumber        string `json:"serialNumber,omitempty"`
	UserAssignedOnly    bool   `json:"userAssignedOnly,omitempty"`
	DeviceAssignedOnly  bool   `json:"deviceAssignedOnly,omitempty"`
	OverrideIndex       int    `json:"overrideIndex,omitempty"`
}

// Gets the LicensesSrv information
func (c *Client) GetLicensesSrv(options LicensesSrvOptions) (*LicensesSrv, error) {

	if options.SToken == "" {
		options.SToken = c.VPPToken.SToken
	}

	// Get the LicensesSrvURL
	licensesSrvURL := c.ServiceConfigSrv.GetLicensesSrvURL

	// Create the LicensesSrv request
	req, err := c.newRequest("POST", licensesSrvURL, options)
	if err != nil {
		return nil, errors.Wrap(err, "create LicensesSrv request")
	}

	// Make the request
	var response LicensesSrv
	err = c.do(req, &response)

	return &response, errors.Wrap(err, "make LicensesSrv request")
}

// Checks if a particular serial is associated with an appID
func (c *Client) CheckAssignedLicense(serial string, appID string) (bool, error) {

	options := LicensesSrvOptions{
		SerialNumber: serial,
		AdamID:       appID,
	}

	// Get all licenses with serial and appID associated
	response, err := c.GetLicensesSrv(options)
	if err != nil {
		return false, err
	}

	// Check the count of licenses returned
	if response.TotalCount > 0 {
		return true, nil
	}

	return false, nil
}
