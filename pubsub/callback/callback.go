package callback

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/micromdm/micromdm/pubsub"
)

type CallbackPublisher struct {
	url  string
	next pubsub.Publisher
}

func (p *CallbackPublisher) Publish(ctx context.Context, topic string, msg []byte) error {

	// send to callback url
	var buf bytes.Buffer
	var event = struct {
		Topic   string `json:"topic"`
		Message []byte `json:"message"`
	}{
		Topic:   topic,
		Message: msg,
	}
	// TODO error handling
	json.NewEncoder(&buf).Encode(&event)
	http.Post(p.url, "application/json", &buf)

	// call the curent publisher as normal
	return p.next.Publish(ctx, topic, msg)
}

func NewCallBackPublisher(next pubsub.Publisher, url string) *CallbackPublisher {
	return &CallbackPublisher{
		url:  url,
		next: next,
	}
}
