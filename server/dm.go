package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const enrollmentIDHeader = "X-Enrollment-ID"

type DeclarativeManagementHTTPCaller struct {
	url    *url.URL
	client *http.Client
}

// NewDeclarativeManagementHTTPCaller creates a new DeclarativeManagementHTTPCaller
func NewDeclarativeManagementHTTPCaller(urlPrefix string, client *http.Client) (*DeclarativeManagementHTTPCaller, error) {
	url, err := url.Parse(urlPrefix)
	return &DeclarativeManagementHTTPCaller{url: url, client: client}, err
}

type HTTPStatusError struct {
	Err    error
	Status int
}

func (e HTTPStatusError) Error() string {
	return e.Err.Error()
}

func (e HTTPStatusError) StatusCode() int {
	return e.Status
}

// DeclarativeManagement calls out to an HTTP URL to handle the actual Declarative Management protocol
func (c *DeclarativeManagementHTTPCaller) DeclarativeManagement(ctx context.Context, id, endpoint string, data []byte) ([]byte, error) {
	if c.url == nil {
		return nil, errors.New("missing URL")
	}
	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("parsing endpoint URL: %w", err)
	}
	u := c.url.ResolveReference(endpointURL)
	method := http.MethodGet
	if len(data) > 0 {
		method = http.MethodPut
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set(enrollmentIDHeader, id)
	if len(data) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		// return the same HTTP status with a Go-kit StatusCoder
		return bodyBytes, HTTPStatusError{
			Err:    fmt.Errorf("unexpected HTTP status: %d", resp.StatusCode),
			Status: resp.StatusCode,
		}
	}
	return bodyBytes, nil
}
