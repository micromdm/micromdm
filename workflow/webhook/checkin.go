package webhook

import (
	"context"
	"fmt"
	"net/http"

	"github.com/micromdm/micromdm/platform/pubsub"
	"github.com/pkg/errors"
	"github.com/micromdm/micromdm/mdm/checkin"
	"encoding/json"
	"bytes"
)

//const contentType = "application/x-apple-aspen-mdm"

type CheckinWebhook struct {
	Topic       string
	EventType   string
	CallbackURL string
	HTTPClient  *http.Client
}

func NewCheckinWebhook(httpClient *http.Client, name string, topic, callbackURL string) (*CheckinWebhook, error) {
	if topic == "" {
		return nil, errors.New("webhook: topic should not be empty")
	}

	if callbackURL == "" {
		return nil, errors.New("webhook: callbackURL should not be empty")
	}

	return &CheckinWebhook{HTTPClient: httpClient, EventType: name, Topic: topic, CallbackURL: callbackURL}, nil
}

func (cw CheckinWebhook) StartListener(sub pubsub.Subscriber) error {
	checkinEvents, err := sub.Subscribe(context.TODO(), "webhook", cw.Topic)
	if err != nil {
		return errors.Wrapf(err,
			"subscribing %s webhook to %s topic", cw.EventType, cw.Topic)
	}

	go func() {
		for {
			select {
			case event := <-checkinEvents:
				var ev checkin.Event
				if err := checkin.UnmarshalEvent(event.Message, &ev); err != nil {
					fmt.Println(err)
					continue
				}
				messageData := map[string]string{"event": cw.EventType, "udid": ev.Command.UDID}
				messageJson, _ := json.Marshal(messageData)

				_, err := cw.HTTPClient.Post(cw.CallbackURL, "application/json", bytes.NewBuffer(messageJson))
				if err != nil {
					fmt.Printf("error sending checking callback: %s\n", err)
				}
			}
		}
	}()

	return nil
}
