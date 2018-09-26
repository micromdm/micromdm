package dep

import "github.com/pkg/errors"

const accountBasePath = "account"

type LimitDict struct {
	Default int `json:"default"`
	Maximum int `json:"maximum"`
}

type URL struct {
	URI        string    `json:"uri"`
	HTTPMethod []string  `json:"http_method"`
	Limit      LimitDict `json:"limit"`
}

type Account struct {
	ServerName    string `json:"server_name"`
	ServerUUID    string `json:"server_uuid"`
	AdminID       string `json:"admin_id"`
	FacilitatorID string `json:"facilitator_id,omitempty"` //deprecated
	OrgName       string `json:"org_name"`
	OrgEmail      string `json:"org_email"`
	OrgPhone      string `json:"org_phone"`
	OrgAddress    string `json:"org_address"`
	OrgID         string `json"org_id"`
	OrgIDHash     string `json"org_id_hash"`
	URLs          []URL  `json:"urls"`
	OrgType       string `json:"org_type"`
	OrgVersion    string `json:"org_version"`
}

func (c *Client) Account() (*Account, error) {
	var account Account
	req, err := c.newRequest("GET", accountBasePath, nil)
	if err != nil {
		return nil, errors.Wrap(err, "create account request")
	}

	err = c.do(req, &account)
	return &account, errors.Wrap(err, "account request")
}
