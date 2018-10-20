package vpp

import "github.com/pkg/errors"

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

// Options for the LicensesSrv
type GetLicensesSrvOptions struct {
	SToken       string `json:"sToken"`
	SerialNumber string `json:"serialNumber,omitempty"`
}

// Gets the LicensesSrv information
func (c *Client) GetLicensesSrv(options GetLicensesSrvOptions) (*LicensesSrv, error) {
	// Sends the sToken string
	options.SToken = c.SToken

	// Get the LicensesSrvURL
	licensesSrvURL := c.VPPServiceConfigSrv.GetLicensesSrvURL

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

// Gets licenses with specified serial associated
func (c *Client) GetLicensesForSerial(serial string) ([]License, error) {
	options := GetLicensesSrvOptions{
		SerialNumber: serial,
	}

	response, err := c.GetLicensesSrv(options)
	if err != nil {
		return nil, err
	}
	licenses := response.Licenses
	return licenses, err
}

// Checks if a particular serial is associated with an appID
func (c *Client) CheckAssignedLicense(serial string, appID string) (bool, error) {
	// Get all licenses with serial associated
	licenses, err := c.GetLicensesForSerial(serial)
	if err != nil {
		return false, err
	}

	// Check for the particular appID
	for _, lic := range licenses {
		if lic.AdamIDStr == appID {
			return true, nil
		}
	}
	return false, nil
}
