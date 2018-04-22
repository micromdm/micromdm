package webhook

import (
	"context"
	"fmt"
	"net/http"

	"github.com/micromdm/micromdm/platform/pubsub"
	"github.com/pkg/errors"
	"bytes"
)

//const contentType = "application/x-apple-aspen-mdm"

type JSONWebhook struct {
	Topic       string
	CallbackURL string
	HTTPClient  *http.Client
}

func NewJSONWebhook(httpClient *http.Client, topic, callbackURL string) (*JSONWebhook, error) {
	if topic == "" {
		return nil, errors.New("webhook: topic should not be empty")
	}

	if callbackURL == "" {
		return nil, errors.New("webhook: callbackURL should not be empty")
	}

	return &JSONWebhook{HTTPClient: httpClient, Topic: topic, CallbackURL: callbackURL}, nil
}

func (cw JSONWebhook) StartListener(sub pubsub.Subscriber) error {
	events, err := sub.Subscribe(context.TODO(), "webhook", cw.Topic)
	if err != nil {
		return errors.Wrapf(err,
			"subscribing JSON webhook to %s topic", cw.Topic)
	}

	go func() {
		for {
			select {
			case event := <-events:
				_, err := cw.HTTPClient.Post(cw.CallbackURL, "application/json", bytes.NewBuffer(event.Message))
				if err != nil {
					fmt.Printf("error sending checking callback: %s\n", err)
				}
			}
		}
	}()

	return nil
}
