package webhook

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/micromdm/micromdm/platform/pubsub"
	"github.com/pkg/errors"
)

type EnrollWebhook struct {
	Topic       string
	CallbackURL string
	HTTPClient  *http.Client
}

func NewEnrollWebhook(httpClient *http.Client, topic, callbackURL string) (*EnrollWebhook, error) {
	if topic == "" {
		return nil, errors.New("webhook: topic should not be empty")
	}

	if callbackURL == "" {
		return nil, errors.New("webhook: callbackURL should not be empty")
	}

	return &EnrollWebhook{HTTPClient: httpClient, Topic: topic, CallbackURL: callbackURL}, nil
}

func (cw EnrollWebhook) StartListener(sub pubsub.Subscriber) error {
	enrollEvent, err := sub.Subscribe(context.TODO(), "enrollWebhook", cw.Topic)
	if err != nil {
		return errors.Wrapf(err,
			"subscribing enrollWebhook to %s topic", cw.Topic)
	}

	go func() {
		for {

			select {
			case event := <-enrollEvent:
				_, err := cw.HTTPClient.Post(cw.CallbackURL, "application/json", bytes.NewBuffer(event.Message))
				if err != nil {
					fmt.Printf("error sending enroll callback: %s\n", err)
				}
			}
		}
	}()

	return nil
}
