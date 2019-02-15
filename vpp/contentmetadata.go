package vpp

import (
	"github.com/pkg/errors"
)

type ContentMetadata struct {
	Results map[string]interface{} `json:"results"`
}

type AppData struct {
	DeviceFamilies []string `json:"deviceFamilies"`
	ID             int      `json:"id"`
	Name           string   `json:"name"`
}

type ContentMetadataOptions struct {
	SToken   string `json:"itvt"`
	ID       string `json:"id"`
	Version  string `json:"version"`
	Platform string `json:"platform"`
}

func (c *Client) GetContentMetadata(options ContentMetadataOptions) (*ContentMetadata, error) {

	contentMetadataLookupURL := c.ServiceConfigSrv.ContentMetadataLookupURL

	req, err := c.newRequest("GET", contentMetadataLookupURL, options)
	if err != nil {
		return nil, errors.Wrap(err, "create ContentMetadata request")
	}

	q := req.URL.Query()
	q.Add("p", "mdm-lockup")
	q.Add("caller", "MDM")
	q.Add("cc", "us")
	q.Add("l", "en")

	if options.ID != "" {
		q.Add("id", options.ID)
	}
	if options.Version != "" {
		q.Add("version", options.Version)
	}
	if options.Platform != "" {
		q.Add("platform", options.Platform)
	}

	req.URL.RawQuery = q.Encode()

	var response ContentMetadata

	err = c.do(req, &response)
	if err != nil {
		return nil, errors.Wrap(err, "please verify your parameters are valid")
	}

	return &response, errors.Wrap(err, "get ContentMetadata")
}
