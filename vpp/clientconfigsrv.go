package vpp

import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
)

// Contains information that associates your particular mdm server to a VPP account token
type ClientContext struct {
	HostName string `json:"hostname"`
	GUID     string `json:"guid"`
}

// Contains location information associated with a VPP account token
type Location struct {
	LocationName string `json:"locationName"`
	LocationID   int    `json:"locationId"`
}

// Contains org information associated with a VPP account token
type ClientConfigSrv struct {
	ClientContext      string   `json:"clientContext"`
	AppleID            string   `json:"appleId,omitempty"`
	OrganizationIDHash string   `json:"organizationIdHash"`
	Status             int      `json:"status"`
	OrganizationID     int      `json:"organizationId"`
	UID                string   `json:"uId"`
	CountryCode        string   `json:"countryCode"`
	Location           Location `json:"location"`
	APNToken           string   `json:"apnToken"`
	Email              string   `json:"email"`
}

// These specify options for the ClientConfigSrv
type ClientConfigSrvOptions struct {
	SToken        string `json:"sToken,omitempty"`
	Verbose       bool   `json:"verbose,omitempty"`
	ClientContext string `json:"clientContext,omitempty"`
}

// Gets ClientConfigSrv information
func (c *Client) GetClientConfigSrv(options ClientConfigSrvOptions) (*ClientConfigSrv, error) {

	// Set required and optional arguments
	if options.SToken == "" {
		options.SToken = c.VPPToken.SToken
	}

	// Get the ClientConfigSrvURL
	clientConfigSrvURL := c.ServiceConfigSrv.ClientConfigSrvURL

	// Create the ClientConfigSrvURL request
	req, err := c.newRequest("POST", clientConfigSrvURL, options)
	if err != nil {
		return nil, errors.Wrap(err, "create ClientConfigSrv request")
	}

	// Make the request
	var response ClientConfigSrv
	err = c.do(req, &response)

	return &response, errors.Wrap(err, "make ClientConfigSrv request")
}

// Configure ClientContext for all VPPTokens
func (c *Client) ConfigureClientContext(options ClientConfigSrvOptions) error {

	// Set Client Context If Needed
	context, err := c.GetClientContext(options)
	if err != nil {
		return errors.Wrap(err, "GetClientContext request")
	}

	if context.HostName != c.ServerPublicURL {
		_, err := c.SetClientContext(options)
		if err != nil {
			return errors.Wrap(err, "SetClientContext HostName request")
		}
	} else if context.GUID != c.VPPToken.UDID {
		_, err := c.SetClientContext(options)
		if err != nil {
			return errors.Wrap(err, "SetClientContext GUID request")
		}
	}

	return nil
}

// Gets the values that determine which mdm server is associated with a VPP account token
func (c *Client) GetClientContext(options ClientConfigSrvOptions) (*ClientContext, error) {

	// Get the ClientConfigSrv info
	clientConfigSrv, err := c.GetClientConfigSrv(options)
	if err != nil {
		return nil, errors.Wrap(err, "get ClientContext request")
	}

	// Get the ClientContext string
	var clientContext ClientContext
	var context = clientConfigSrv.ClientContext
	if context != "" {
		// Convert the string to a ClientContext type
		err = DecodeToClientContext(context, &clientContext)
		if err != nil {
			return nil, errors.Wrap(err, "get ClientContext")
		}
	} else {
		clientContext.HostName = ""
		clientContext.GUID = ""
	}

	return &clientContext, nil
}

// Sets the values that determine which mdm server is associated with a VPP account token
func (c *Client) SetClientContext(options ClientConfigSrvOptions) (*ClientContext, error) {

	// Generate a ClientContext string with the UDID and the current ServerPublicURL
	context := ClientContext{c.ServerPublicURL, c.VPPToken.UDID}
	data, err := json.Marshal(context)
	if err != nil {
		return nil, errors.Wrap(err, "create new ClientContext")
	}
	newContext := string(data)

	// Enter the new ClientContext string into the ClientConfigSrv options
	options.ClientContext = newContext

	// Set the new ClientContext into the VPP account token
	response, err := c.GetClientConfigSrv(options)
	if err != nil {
		return nil, errors.Wrap(err, "set ClientContext request")
	}

	// Convert the string to a ClientContext type
	var clientContext ClientContext
	err = DecodeToClientContext(response.ClientContext, &clientContext)
	if err != nil {
		return nil, errors.Wrap(err, "set ClientContext")
	}

	return &clientContext, nil
}

func DecodeToClientContext(contextString string, context *ClientContext) error {
	// Convert the string to a ClientContext type
	err := json.NewDecoder(strings.NewReader(contextString)).Decode(&context)
	if err != nil {
		return errors.Wrap(err, "decode ClientContext")
	}
	return nil
}
