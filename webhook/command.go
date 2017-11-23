package webhook

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/micromdm/micromdm/mdm/connect"
	"github.com/micromdm/micromdm/pubsub"
	"github.com/pkg/errors"
)

const contentType = "application/x-apple-aspen-mdm"

type CommandWebhook struct {
	Topic       string
	CallbackURL string
	HTTPClient  *http.Client
}

func NewCommandWebhook(httpClient *http.Client, topic, callbackURL string) (*CommandWebhook, error) {
	if topic == "" {
		return nil, errors.New("webhook: topic should not be empty")
	}

	if callbackURL == "" {
		return nil, errors.New("webhook: callbackURL should not be empty")
	}

	return &CommandWebhook{HTTPClient: httpClient, Topic: topic, CallbackURL: callbackURL}, nil
}

func (cw CommandWebhook) StartListener(sub pubsub.Subscriber) error {
	connectEvents, err := sub.Subscribe(context.TODO(), "commandWebhook", cw.Topic)
	if err != nil {
		return errors.Wrapf(err,
			"subscribing commandWebhook to %s topic", cw.Topic)
	}

	go func() {
		for {
			select {
			case event := <-connectEvents:
				var ev connect.Event
				if err := connect.UnmarshalEvent(event.Message, &ev); err != nil {
					fmt.Println(err)
					continue
				}

				_, err := cw.HTTPClient.Post(cw.CallbackURL, contentType, bytes.NewBuffer(ev.Raw))
				if err != nil {
					fmt.Printf("error sending command response: %s\n", err)
				}
			}
		}
	}()

	return nil
}
