package user

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/micromdm/micromdm/mdm"
	"github.com/micromdm/micromdm/platform/pubsub"
)

type WorkerStore interface {
	Save(*User) error
	DeleteDeviceUsers(udid string) error
	UserByUserID(userID string) (*User, error)
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
	const subscription = "user_worker"
	tokenUpdateEvents, err := w.sub.Subscribe(ctx, subscription, mdm.TokenUpdateTopic)
	if err != nil {
		return errors.Wrapf(err, "subcribe %s to %s", subscription, mdm.TokenUpdateTopic)
	}

	for {
		var err error
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ev := <-tokenUpdateEvents:
			err = w.updateUserFromTokenUpdate(ctx, ev.Message)
		}

		if err != nil {
			level.Info(w.logger).Log(
				"msg", "update user from event",
				"err", err,
			)
			continue
		}
	}
}

func (w *Worker) updateUserFromTokenUpdate(ctx context.Context, message []byte) error {
	var ev mdm.CheckinEvent
	if err := mdm.UnmarshalCheckinEvent(message, &ev); err != nil {
		return errors.Wrap(err, "unmarshal checkin event for user worker")
	}

	if ev.Command.UserID == "" {
		// only process user events
		return nil
	}

	usr, err := getOrCreateUser(w.db, ev.Command.UserID)
	if err != nil {
		return err
	}
	if usr.UUID == "" {
		if err := w.db.DeleteDeviceUsers(ev.Command.UDID); err != nil {
			return errors.Wrapf(
				err,
				"delete users for device %s before re-creating. user_id=%s",
				ev.Command.UDID,
				ev.Command.UserID,
			)
			usr.UUID = uuid.NewV4().String()
		}
	}

	usr.UDID = ev.Command.UDID
	usr.UserID = ev.Command.UserID
	usr.UserLongname = ev.Command.UserLongName
	usr.UserShortname = ev.Command.UserShortName
	usr.AuthToken = ev.Command.Token.String()
	err = w.db.Save(usr)
	return errors.Wrapf(err, "saving user %s to device %s", ev.Command.UserID, ev.Command.UDID)
}

func getOrCreateUser(db WorkerStore, userID string) (*User, error) {
	byGUID, err := db.UserByUserID(userID)
	if err == nil {
		return byGUID, nil
	}
	if err != nil && !isNotFound(err) {
		return nil, errors.Wrap(err, "get user by ID")
	}

	usr := new(User)
	return usr, nil
}

func isNotFound(err error) bool {
	err = errors.Cause(err)
	type notFoundErr interface {
		error
		NotFound() bool
	}

	e, ok := err.(notFoundErr)
	return ok && e.NotFound()
}
