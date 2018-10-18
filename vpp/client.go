package vpp

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/pkg/errors"
)

const (
	serverURL = "https://your.server.com" // This needs to be modified to be imported from server
	version   = ""                        // This needs to be modified to be imported from server

	defaultBaseURL               = "https://vpp.itunes.apple.com/WebObjects/MZFinance.woa/wa/VPPServiceConfigSrv"
	mediaType                    = "application/json;charset=UTF8"
	XServerProtocolVersionHeader = "X-Server-Protocol-Version"
	XServerProtocolVersion       = "3"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	sToken              string
	VPPServiceConfigSrv *VPPServiceConfigSrv
	userAgent           string
	client              HTTPClient
	baseURL             *url.URL
}

func NewClient(sToken string) (*Client, error) {
	baseURL, _ := url.Parse(defaultBaseURL)
	c := Client{
		sToken:    sToken,
		userAgent: path.Join("micromdm", version),
		client:    http.DefaultClient,
		baseURL:   baseURL,
	}

	// Get VPPServiceConfigSrv Data
	vppServiceConfigSrv, err := c.GetVPPServiceConfigSrv()
	if err != nil {
		return nil, errors.Wrap(err, "create VPPServiceConfigSrv request")
	}
	c.VPPServiceConfigSrv = vppServiceConfigSrv

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

func (c *Client) newRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, errors.Wrapf(err, "parse vpp request url %s", urlStr)
	}

	u := c.baseURL.ResolveReference(rel)
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

	req.Header.Add("User-Agent", c.userAgent)
	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add(XServerProtocolVersionHeader, XServerProtocolVersion)
	return req, nil
}

func (c *Client) do(req *http.Request, into interface{}) error {
	resp, err := c.client.Do(req)
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
