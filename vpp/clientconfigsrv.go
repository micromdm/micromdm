package vpp

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"strings"
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
type GetClientConfigSrvOptions func(*getClientConfigSrvOpts) error

type getClientConfigSrvOpts struct {
	SToken        string `json:"sToken"`
	Verbose       bool   `json:"verbose,omitempty"`
	ClientContext string `json:"clientContext,omitempty"`
}

// Verbose is an optional argument that can be added to GetClientConfigSrv
func VerboseOption(verbose bool) GetClientConfigSrvOptions {
	return func(opts *getClientConfigSrvOpts) error {
		opts.Verbose = verbose
		return nil
	}
}

// ClientContext is an optional argument that can be added to GetClientConfigSrv
func ClientContextOption(context string) GetClientConfigSrvOptions {
	return func(opts *getClientConfigSrvOpts) error {
		opts.ClientContext = context
		return nil
	}
}

// Gets ClientConfigSrv information
func (c *Client) GetClientConfigSrv(opts ...GetClientConfigSrvOptions) (*ClientConfigSrv, error) {
	// Set required and optional arguments
	request := &getClientConfigSrvOpts{SToken: c.SToken}
	for _, option := range opts {
		if err := option(request); err != nil {
			return nil, err
		}
	}

	// Get the ClientConfigSrvURL
	clientConfigSrvURL := c.VPPServiceConfigSrv.ClientConfigSrvURL

	// Create the ClientConfigSrvURL request
	req, err := c.newRequest("POST", clientConfigSrvURL, request)
	if err != nil {
		return nil, errors.Wrap(err, "create ClientConfigSrv request")
	}

	// Make the request
	var response ClientConfigSrv
	err = c.do(req, &response)

	return &response, errors.Wrap(err, "make ClientConfigSrv request")
}

// Gets the appleID field along with the standard information
func (c *Client) GetClientConfigSrvVerbose() (*ClientConfigSrv, error) {
	options := VerboseOption(true)
	response, err := c.GetClientConfigSrv(options)
	if err != nil {
		return nil, errors.Wrap(err, "using verbose option")
	}

	return response, nil
}

// Gets the values that determine which mdm server is associated with a VPP account token
func (c *Client) GetClientContext() (*ClientContext, error) {
	// Get the ClientConfigSrv info
	clientConfigSrv, err := c.GetClientConfigSrv()
	if err != nil {
		return nil, errors.Wrap(err, "get ClientContext request")
	}

	// Get the ClientContext string
	var context = clientConfigSrv.ClientContext

	// Convert the string to a ClientContext type
	var clientContext ClientContext
	err = json.NewDecoder(strings.NewReader(context)).Decode(&clientContext)
	if err != nil {
		return nil, errors.Wrap(err, "decode ClientContext")
	}

	return &clientContext, nil
}

// Sets the values that determine which mdm server is associated with a VPP account token
func (c *Client) SetClientContext(serverURL string) (*ClientContext, error) {
	// Generate a UUID that is tracked to ensure VPP licenses are up to date
	uuid := uuid.NewV4().String()

	// Generate a ClientContext string with the new UUID and the current serverURL
	context := ClientContext{serverURL, uuid}
	data, err := json.Marshal(context)
	if err != nil {
		return nil, errors.Wrap(err, "create new ClientContext")
	}
	newContext := string(data)

	// Enter the new ClientContext string into the ClientConfigSrv options
	options := ClientContextOption(newContext)

	// Set the new ClientContext into the VPP account token
	response, err := c.GetClientConfigSrv(options)
	if err != nil {
		return nil, errors.Wrap(err, "set ClientContext request")
	}

	// Get the new ClientContext string
	var contextString = response.ClientContext

	// Convert the string to a ClientContext type
	var clientContext ClientContext
	err = json.NewDecoder(strings.NewReader(contextString)).Decode(&clientContext)
	if err != nil {
		return nil, errors.Wrap(err, "decode new ClientContext")
	}

	return &clientContext, nil
}
