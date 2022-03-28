package vpp

import (
	"bytes"
	"encoding/json"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"github.com/pkg/errors"
)

var version = "dev"

const (
	serverURL = "https://your.server.com" // This needs to be modified to be imported from server

	defaultBaseURL               = "https://vpp.itunes.apple.com/WebObjects/MZFinance.woa/wa/VPPServiceConfigSrv"
	mediaType                    = "application/json;charset=UTF8"
	XServerProtocolVersionHeader = "X-Server-Protocol-Version"
	XServerProtocolVersion       = "3"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Contains the sToken string used to authenticate to the various VPP services
// Contains the return VPPServiceConfigSrv information
type Client struct {
	UDID								string
	SToken              string
	VPPServiceConfigSrv *VPPServiceConfigSrv
	UserAgent           string
	Client              HTTPClient
	BaseURL             *url.URL
}

type VPPToken struct {
	UDID    string `json:"udid"`
	SToken  []byte `json:"sToken"`
}

func NewClient(vppToken VPPToken) (*Client, error) {

	sToken := base64.StdEncoding.EncodeToString(vppToken.SToken)

	baseURL, _ := url.Parse(defaultBaseURL)
	c := Client{
		UDID:      vppToken.UDID,
		SToken:    sToken,
		UserAgent: path.Join("micromdm", version),
		Client:    http.DefaultClient,
		BaseURL:   baseURL,
	}

	// Get VPPServiceConfigSrv Data
	VPPServiceConfigSrv, err := c.GetVPPServiceConfigSrv()
	if err != nil {
		return nil, errors.Wrap(err, "create VPPServiceConfigSrv request")
	}
	c.VPPServiceConfigSrv = VPPServiceConfigSrv

	// Set Client Context If Needed
	context, err := c.GetClientContext()
	if err != nil {
		return nil, errors.Wrap(err, "GetClientContext request")
	}

	if context.HostName != serverURL {
		_, err := c.SetClientContext(serverURL)
		if err != nil {
			return nil, errors.Wrap(err, "SetClientContext request")
		}
	}

	return &c, nil
}

func (c *Client) newRequest(method, URLStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(URLStr)
	if err != nil {
		return nil, errors.Wrapf(err, "parse vpp request url %s", URLStr)
	}

	u := c.BaseURL.ResolveReference(rel)
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, errors.Wrap(err, "encode http body for VPP request")
		}
	}

	req, err := http.NewRequest(method, u.String(), &buf)
	if err != nil {
		return nil, errors.Wrapf(err, "creating %s request to vpp %s", method, u.String())
	}

	req.Header.Add("User-Agent", c.UserAgent)
	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add(XServerProtocolVersionHeader, XServerProtocolVersion)
	return req, nil
}

func (c *Client) do(req *http.Request, into interface{}) error {
	resp, err := c.Client.Do(req)
	if err != nil {
		return errors.Wrap(err, "perform vpp request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.Errorf("unexpected vpp response. status=%d VPP API Error: %s", resp.StatusCode, string(body))
	}

	err = json.NewDecoder(resp.Body).Decode(into)

	return errors.Wrap(err, "decode VPP response body")
}
