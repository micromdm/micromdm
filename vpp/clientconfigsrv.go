package vpp

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"strings"
)

type ClientContext struct {
	HostName string `json:"hostname"`
	Guid     string `json:"guid"`
}

type Location struct {
	LocationName string `json:"locationName"`
	LocationId   int    `json:"locationId"`
}

type ClientConfigSrvResponse struct {
	ClientContext      interface{} `json:"clientContext"`
	AppleId            string      `json:"appleId"`
	OrganizationIdHash string      `json:"organizationIdHash"`
	Status             int         `json:"status"`
	OrganizationId     int         `json:"organizationId"`
	UId                string      `json:"uId"`
	CountryCode        string      `json:"countryCode"`
	Location           Location    `json:"location"`
	ApnToken           string      `json:"apnToken"`
	Email              string      `json:"email"`
}

type ClientConfigSrv struct {
	ClientContext      ClientContext `json:"clientContext"`
	AppleId            string        `json:"appleId"`
	OrganizationIdHash string        `json:"organizationIdHash"`
	Status             int           `json:"status"`
	OrganizationId     int           `json:"organizationId"`
	UId                string        `json:"uId"`
	CountryCode        string        `json:"countryCode"`
	Location           Location      `json:"location"`
	ApnToken           string        `json:"apnToken"`
	Email              string        `json:"email"`
}

func (c *Client) GetClientConfigSrv() (*ClientConfigSrv, error) {
	request := map[string]interface{}{
		"sToken": c.sToken,
	}

	clientConfigSrvUrl := c.VPPServiceConfigSrv.ClientConfigSrvUrl
	var response ClientConfigSrvResponse
	req, err := c.newRequest("POST", clientConfigSrvUrl, request)
	if err != nil {
		return nil, errors.Wrap(err, "create ClientConfigSrv request")
	}

	err = c.do(req, &response)

	var context = response.ClientContext
	var clientContext ClientContext
	json.NewDecoder(strings.NewReader(context.(string))).Decode(&clientContext)

	config := ClientConfigSrv{
		ClientContext:      clientContext,
		AppleId:            response.AppleId,
		OrganizationIdHash: response.OrganizationIdHash,
		Status:             response.Status,
		OrganizationId:     response.OrganizationId,
		UId:                response.UId,
		CountryCode:        response.CountryCode,
		Location:           response.Location,
		ApnToken:           response.ApnToken,
		Email:              response.Email,
	}
	return &config, errors.Wrap(err, "ClientConfigSrv request")
}

func (c *Client) GetClientContext() (*ClientContext, error) {
	clientConfigSrv, err := c.GetClientConfigSrv()
	if err != nil {
		return nil, errors.Wrap(err, "create ClientConfigSrv request")
	}

	var context = clientConfigSrv.ClientContext
	return &context, nil
}

func (c *Client) SetClientContext(serverURL string) (*ClientContext, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, errors.Wrap(err, "create uuid")
	}

	newContext := map[string]interface{}{
		"sToken":        c.sToken,
		"clientContext": fmt.Sprintf("{\"hostname\":\"%s\",\"guid\":\"%s\"}", serverURL, uuid),
	}

	clientConfigSrvUrl := c.VPPServiceConfigSrv.ClientConfigSrvUrl

	var response ClientConfigSrvResponse

	req, err := c.newRequest("POST", clientConfigSrvUrl, newContext)
	if err != nil {
		return nil, errors.Wrap(err, "create Client Context ClientConfigSrv request")
	}

	err = c.do(req, &response)

	var context = response.ClientContext
	var clientContext ClientContext
	json.NewDecoder(strings.NewReader(context.(string))).Decode(&clientContext)

	return &clientContext, errors.Wrap(err, "set Client Context ClientConfigSrv request")
}
