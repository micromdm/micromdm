package apns

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/mdm"
	"github.com/micromdm/micromdm/platform/pubsub"
)

type WorkerStore interface {
	Save(context.Context, *PushInfo) error
}

type Worker struct {
	db     WorkerStore
	sub    pubsub.Subscriber
	logger log.Logger
}

func NewWorker(db WorkerStore, subscriber pubsub.Subscriber, logger log.Logger) *Worker {
	return &Worker{
		db:     db,
		sub:    subscriber,
		logger: logger,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	const subscription = "pushinfo_worker"
	tokenUpdateEvents, err := w.sub.Subscribe(ctx, subscription, mdm.TokenUpdateTopic)
	if err != nil {
		return errors.Wrapf(err,
			"subscribing %s to %s topic", subscription, mdm.TokenUpdateTopic)
	}

	for {
		var err error
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event := <-tokenUpdateEvents:
			err = w.updatePushInfoFromTokenUpdate(ctx, event.Message)
		}
		if err != nil {
			level.Info(w.logger).Log(
				"msg", "update pushinfo from event",
				"err", err,
			)
			continue
		}
	}
}

func (w *Worker) updatePushInfoFromTokenUpdate(ctx context.Context, message []byte) error {
	var ev mdm.CheckinEvent
	if err := mdm.UnmarshalCheckinEvent(message, &ev); err != nil {
		return errors.Wrap(err, "unmarshal pushinfo event")
	}
	info := PushInfo{
		UDID:      ev.Command.UDID,
		Token:     ev.Command.Token.String(),
		PushMagic: ev.Command.PushMagic,
		MDMTopic:  ev.Command.Topic,
	}
	if ev.Command.UserID != "" {
		// use the GUID if this is a user TokenUpdate.
		info.UDID = ev.Command.UserID
	}
	err := w.db.Save(ctx, &info)
	return errors.Wrapf(err, "saving pushinfo for udid=%s", info.UDID)
}
