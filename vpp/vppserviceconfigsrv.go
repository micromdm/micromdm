package vpp

import "github.com/pkg/errors"

type VPPServiceConfigSrv struct {
	EditUserSrvUrl                   string      `json:"editUserSrvUrl"`
	DisassociateLicenseSrvUrl        string      `json:"disassociateLicenseSrvUrl"`
	ContentMetadataLookupUrl         string      `json:"contentMetadataLookupUrl"`
	ClientConfigSrvUrl               string      `json:"clientConfigSrvUrl"`
	GetUserSrvUrl                    string      `json:"getUserSrvUrl"`
	GetUsersSrvUrl                   string      `json:"getUsersSrvUrl"`
	GetLicensesSrvUrl                string      `json:"getLicensesSrvUrl"`
	GetVPPAssetsSrvUrl               string      `json:"getVPPAssetsSrvUrl"`
	VppWebsiteUrl                    string      `json:"vppWebsiteUrl"`
	InvitationEmailUrl               string      `json:"invitationEmailUrl"`
	RetireUserSrvUrl                 string      `json:"retireUserSrvUrl"`
	AssociateLicenseSrvUrl           string      `json:"associateLicenseSrvUrl"`
	ManageVPPLicensesByAdamIdSrvUrl  string      `json:"manageVPPLicensesByAdamIdSrvUrl"`
	RegisterUserSrvUrl               string      `json:"registerUserSrvUrl"`
	MaxBatchAssociateLicenseCount    int         `json:"maxBatchAssociateLicenseCount"`
	MaxBatchDisassociateLicenseCount int         `json:"maxBatchDisassociateLicenseCount"`
	Status                           int         `json:"status"`
	ErrorCodes                       interface{} `json:"errorCodes"`
}

func (c *Client) GetVPPServiceConfigSrv() (*VPPServiceConfigSrv, error) {
	var response VPPServiceConfigSrv
	req, err := c.newRequest("GET", c.baseURL.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "create VPPServiceConfigSrv request")
	}

	err = c.do(req, &response)
	return &response, errors.Wrap(err, "VPPServiceConfigSrv request")
}
