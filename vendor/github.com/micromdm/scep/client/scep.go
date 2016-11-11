package scepclient

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/micromdm/scep/server"
	"golang.org/x/net/context"
)

// Client implements the SCEP service and extra methods
type Client interface {
	scepserver.Service
	Supports(string) bool
}
type client struct {
	getRemote    endpoint.Endpoint
	postRemote   endpoint.Endpoint
	capabilities []byte
}

func (c *client) Supports(cap string) bool {
	if len(c.capabilities) == 0 {
		ctx := context.Background()
		// try to retrieve caps
		c.GetCACaps(ctx)
	}
	return bytes.Contains(c.capabilities, []byte(cap))
}

// NewClient returns a SCEP service that's backed by the provided Endpoint
func NewClient(baseURL string) Client {
	scepURL, _ := url.Parse(baseURL)
	httpc := http.DefaultClient
	return &client{
		getRemote: httptransport.NewClient(
			"GET",
			scepURL,
			scepserver.EncodeSCEPRequest,
			scepserver.DecodeSCEPResponse,
			httptransport.SetClient(httpc),
		).Endpoint(),
		postRemote: httptransport.NewClient(
			"POST",
			scepURL,
			scepserver.EncodeSCEPRequest,
			scepserver.DecodeSCEPResponse,
			httptransport.SetClient(httpc),
		).Endpoint(),
	}
}

func (c *client) GetCACaps(ctx context.Context) ([]byte, error) {
	request := scepserver.SCEPRequest{
		Operation: "GetCACaps",
	}
	reply, err := c.getRemote(ctx, request)
	if err != nil {
		return nil, err
	}
	r := reply.(scepserver.SCEPResponse)
	c.capabilities = r.Data
	return r.Data, nil
}

func (c *client) GetCACert(ctx context.Context) ([]byte, int, error) {
	request := scepserver.SCEPRequest{
		Operation: "GetCACert",
	}
	reply, err := c.getRemote(ctx, request)
	if err != nil {
		return nil, 0, err
	}
	r := reply.(scepserver.SCEPResponse)
	return r.Data, r.CACertNum, nil
}

func (c *client) PKIOperation(ctx context.Context, data []byte) ([]byte, error) {
	request := scepserver.SCEPRequest{
		Operation: "PKIOperation",
		Message:   data,
	}
	if c.Supports("POSTPKIOperation") {
		reply, err := c.postRemote(ctx, request)
		if err != nil {
			return nil, err
		}
		r := reply.(scepserver.SCEPResponse)
		return r.Data, nil
	}
	return nil, errors.New("no POSTPKIOperation support")
}

func (c *client) GetNextCACert(ctx context.Context) ([]byte, error) {
	panic("not implemented")
}
