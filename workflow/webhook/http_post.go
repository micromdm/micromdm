package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

func postWebhookEvent(
	ctx context.Context,
	client httpClient,
	url string,
	event interface{},
) error {
	raw, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		return errors.Wrap(err, "marshal webhook event")
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(raw))
	if err != nil {
		return errors.Wrap(err, "create webhook http request")
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return errors.Wrap(err, "post webhook event to URL")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return errors.Errorf("received unexpected HTTP status %d %s", resp.StatusCode, resp.Status)
	}
	return nil
}
