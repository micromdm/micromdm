package config

import (
	"encoding/base64"
	"encoding/json"

	"github.com/micromdm/micromdm/vpp"
)

const VPPTokenTopic = "mdm.VPPTokenAdded"

type VPPToken struct {
	UDID   string `json:"udid"`
	SToken SToken `json:"sToken"`
}

type SToken struct {
	OrgName string `json:"orgName"`
	Token   string `json:"token"`
	ExpDate string `json:"expDate"`
}

// create a VPP client from token.
func (tok VPPToken) Client(serverURL string) (*vpp.Client, error) {

	// Convert to JSON
	tokenJSON, err := json.Marshal(tok.SToken)
	if err != nil {
		return nil, err
	}
	// Encode Base64
	sToken := base64.StdEncoding.EncodeToString(tokenJSON)

	conf := vpp.VPPToken{
		UDID:   tok.UDID,
		SToken: sToken,
	}

	client, err := vpp.NewClient(conf, serverURL)
	if err != nil {
		return nil, err
	}
	return client, nil
}
