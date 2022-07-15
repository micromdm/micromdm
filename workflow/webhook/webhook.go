package webhook

import (
	"context"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/mdm"
	"github.com/micromdm/micromdm/platform/pubsub"
)

type Event struct {
	Topic     string    `json:"topic"`
	EventID   string    `json:"event_id"`
	CreatedAt time.Time `json:"created_at"`

	AcknowledgeEvent *AcknowledgeEvent `json:"acknowledge_event,omitempty"`
	CheckinEvent     *CheckinEvent     `json:"checkin_event,omitempty"`
}

type Worker struct {
	logger log.Logger
	url    string
	client *http.Client
	sub    pubsub.Subscriber
}

type Option func(*Worker)

func WithLogger(logger log.Logger) Option {
	return func(w *Worker) {
		w.logger = logger
	}
}

func WithHTTPClient(client *http.Client) Option {
	return func(w *Worker) {
		w.client = client
	}
}

func New(url string, sub pubsub.Subscriber, opts ...Option) *Worker {
	worker := &Worker{
		url:    url,
		sub:    sub,
		logger: log.NewNopLogger(),
		client: http.DefaultClient,
	}

	for _, optFn := range opts {
		optFn(worker)
	}

	return worker
}

func (w *Worker) Run(ctx context.Context) error {
	const subscription = "webhook_worker"

	ackEvents, err := w.sub.Subscribe(ctx, subscription, mdm.ConnectTopic)
	if err != nil {
		return errors.Wrapf(err, "subscribe %s to %s", subscription, mdm.ConnectTopic)
	}

	authenticateEvents, err := w.sub.Subscribe(ctx, subscription, mdm.AuthenticateTopic)
	if err != nil {
		return errors.Wrapf(err, "subscribe %s to %s", subscription, mdm.AuthenticateTopic)
	}

	tokenUpdateEvents, err := w.sub.Subscribe(ctx, subscription, mdm.TokenUpdateTopic)
	if err != nil {
		return errors.Wrapf(err, "subscribe %s to %s", subscription, mdm.TokenUpdateTopic)
	}

	checkoutEvents, err := w.sub.Subscribe(ctx, subscription, mdm.CheckoutTopic)
	if err != nil {
		return errors.Wrapf(err, "subscribe %s to %s", subscription, mdm.CheckoutTopic)
	}

	getBootstrapTokenEvents, err := w.sub.Subscribe(ctx, subscription, mdm.GetBootstrapTokenTopic)
	if err != nil {
		return errors.Wrapf(err, "subscribe %s to %s", subscription, mdm.GetBootstrapTokenTopic)
	}

	setBootstrapTokenEvents, err := w.sub.Subscribe(ctx, subscription, mdm.SetBootstrapTokenTopic)
	if err != nil {
		return errors.Wrapf(err, "subscribe %s to %s", subscription, mdm.SetBootstrapTokenTopic)
	}

	for {
		var (
			event *Event
			err   error
		)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ev := <-ackEvents:
			event, err = acknowledgeEvent(ev.Topic, ev.Message)
		case ev := <-authenticateEvents:
			event, err = checkinEvent(ev.Topic, ev.Message)
		case ev := <-tokenUpdateEvents:
			event, err = checkinEvent(ev.Topic, ev.Message)
		case ev := <-checkoutEvents:
			event, err = checkinEvent(ev.Topic, ev.Message)
		case ev := <-getBootstrapTokenEvents:
			event, err = checkinEvent(ev.Topic, ev.Message)
		case ev := <-setBootstrapTokenEvents:
			event, err = checkinEvent(ev.Topic, ev.Message)
		}

		if err != nil {
			level.Info(w.logger).Log(
				"msg", "create webhook event",
				"err", err,
			)
			continue
		}

		if err := postWebhookEvent(ctx, w.client, w.url, event); err != nil {
			level.Info(w.logger).Log(
				"msg", "post webhook event",
				"err", err,
			)
			continue
		}

	}
}
