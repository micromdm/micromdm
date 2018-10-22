package vpp

import "github.com/pkg/errors"

// Contains the most recent data from Apple for configuring vpp
type VPPServiceConfigSrv struct {
	EditUserSrvURL                   string  `json:"editUserSrvUrl"`
	DisassociateLicenseSrvURL        string  `json:"disassociateLicenseSrvUrl"`
	ContentMetadataLookupURL         string  `json:"contentMetadataLookupUrl"`
	ClientConfigSrvURL               string  `json:"clientConfigSrvUrl"`
	GetUserSrvURL                    string  `json:"getUserSrvUrl"`
	GetUsersSrvURL                   string  `json:"getUsersSrvUrl"`
	GetLicensesSrvURL                string  `json:"getLicensesSrvUrl"`
	GetVPPAssetsSrvURL               string  `json:"getVPPAssetsSrvUrl"`
	VppWebsiteURL                    string  `json:"vppWebsiteUrl"`
	InvitationEmailURL               string  `json:"invitationEmailUrl"`
	RetireUserSrvURL                 string  `json:"retireUserSrvUrl"`
	AssociateLicenseSrvURL           string  `json:"associateLicenseSrvUrl"`
	ManageVPPLicensesByAdamIdSrvURL  string  `json:"manageVPPLicensesByAdamIdSrvUrl"`
	RegisterUserSrvURL               string  `json:"registerUserSrvUrl"`
	MaxBatchAssociateLicenseCount    int     `json:"maxBatchAssociateLicenseCount"`
	MaxBatchDisassociateLicenseCount int     `json:"maxBatchDisassociateLicenseCount"`
	Status                           int     `json:"status"`
	ErrorCodes                       []Error `json:"errorCodes"`
}

type Error struct {
	ErrorMessage string `json:"errorMessage"`
	ErrorNumber  int    `json:"errorNumber"`
}

func (c *Client) GetVPPServiceConfigSrv() (*VPPServiceConfigSrv, error) {
	var response VPPServiceConfigSrv
	req, err := c.newRequest("GET", c.BaseURL.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "create VPPServiceConfigSrv request")
	}

	err = c.do(req, &response)
	return &response, errors.Wrap(err, "VPPServiceConfigSrv request")
}
