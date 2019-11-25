package restream

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/mdm"
	"github.com/micromdm/micromdm/platform/command"
	"github.com/micromdm/micromdm/platform/pubsub"
)

type Event struct {
	ID       string // the command uuid
	Status   string // TODO: make an enum
	Request  []byte
	Response []byte
}

type Worker struct {
	logger log.Logger
	db     Store
	sub    pubsub.Subscriber
}

type Store interface {
	Save(ctx context.Context, event Event) error
	Event(ctx context.Context, id string) (Event, error)
}

type Option func(*Worker)

func WithLogger(logger log.Logger) Option {
	return func(w *Worker) {
		w.logger = logger
	}
}

func New(db Store, sub pubsub.Subscriber, opts ...Option) *Worker {
	worker := &Worker{
		db:     db,
		sub:    sub,
		logger: log.NewNopLogger(),
	}

	for _, optFn := range opts {
		optFn(worker)
	}

	return worker
}

func (w *Worker) Run(ctx context.Context) error {
	const subscription = "restream_worker"

	ackEvents, err := w.sub.Subscribe(ctx, subscription, mdm.ConnectTopic)
	if err != nil {
		return errors.Wrapf(err, "subscribe %s to %s", subscription, mdm.ConnectTopic)
	}
	commandEvents, err := w.sub.Subscribe(ctx, subscription, command.CommandTopic)
	if err != nil {
		return errors.Wrapf(err, "subscribe %s to %s", subscription, command.CommandTopic)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ev := <-commandEvents:
			var cmd command.Event
			if err := command.UnmarshalEvent(ev.Message, &cmd); err != nil {
				level.Info(w.logger).Log("msg", "unmarshal command event in restream", "err", err)
				continue
			}
			var event Event
			event.ID = cmd.Payload.CommandUUID
			event.Request = ev.Message
			event.Status = "QUEUED"
			if err := w.db.Save(ctx, event); err != nil {
				level.Info(w.logger).Log("msg", "save restream event", "event_id", event.ID, "err", err)
				continue
			}
		case ev := <-ackEvents:
			var resp mdm.AcknowledgeEvent
			if err := mdm.UnmarshalAcknowledgeEvent(ev.Message, &resp); err != nil {
				level.Info(w.logger).Log("msg", "unmarshal ack  event in restream", "err", err)
				continue
			}
			if resp.Response.CommandUUID == "" {
				// just an idle response, not associated with a cmd
				continue
			}
			event, err := w.db.Event(ctx, resp.Response.CommandUUID)
			if err != nil {
				level.Info(w.logger).Log("msg", "get restream event", "err", err)
				continue
			}
			event.Status = resp.Response.Status
			event.Response = ev.Message
			if err := w.db.Save(ctx, event); err != nil {
				level.Info(w.logger).Log("msg", "save updated restream event", "event_id", event.ID, "err", err)
				continue
			}
		}
	}
}
