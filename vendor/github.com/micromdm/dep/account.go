package dep

const (
	accountBasePath = "account"
)

// AccountService communicates with the DEP Account Details endpoint
/*
	account, err := client.Account()
*/
type AccountService interface {
	Account() (*Account, error)
}

type accountService struct {
	client *depClient
}

// Account is a DEP account
type Account struct {
	ServerName    string   `json:"server_name"`
	ServerUUID    string   `json:"server_uuid"`
	AdminID       string   `json:"admin_id"`
	FacilitatorID string   `json:"facilitator_id,omitempty"` //deprecated
	OrgName       string   `json:"org_name"`
	OrgEmail      string   `json:"org_email"`
	OrgPhone      string   `json:"org_phone"`
	OrgAddress    string   `json:"org_address"`
	URLs          []string `json:"urls"`
}

// Account returns account details
func (s accountService) Account() (*Account, error) {
	var account Account
	req, err := s.client.NewRequest("GET", accountBasePath, nil)
	if err != nil {
		return nil, err
	}
	err = s.client.Do(req, &account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}
