package depsync

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/micromdm/dep"
	"github.com/micromdm/micromdm/deptoken"
	"github.com/micromdm/micromdm/pubsub"
	"github.com/pkg/errors"
)

const (
	SyncTopic    = "mdm.DepSync"
	ConfigBucket = "mdm.DEPConfig"
	syncDuration = 30 * time.Minute
)

type Syncer interface {
	SyncNow()
}

type watcher struct {
	mtx    sync.RWMutex
	logger log.Logger
	client dep.Client

	publisher pubsub.Publisher
	conf      *config
	startSync chan bool
	syncNow   chan bool
}

type cursor struct {
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

// A cursor is valid for a week.
func (c cursor) Valid() bool {
	expiration := time.Now().Add(24 * 7 * time.Hour)
	if c.CreatedAt.After(expiration) {
		return false
	}
	return true
}

type Option func(*watcher)

func WithClient(client dep.Client) Option {
	return func(w *watcher) {
		w.client = client
	}
}

func WithLogger(logger log.Logger) Option {
	return func(w *watcher) {
		w.logger = logger
	}
}

func New(pub pubsub.PublishSubscriber, db *bolt.DB, logger log.Logger, opts ...Option) (Syncer, error) {
	conf, err := LoadConfig(db)
	if err != nil {
		return nil, err
	}
	if conf.Cursor.Valid() {
		level.Info(logger).Log("msg", "loaded DEP config", "cursor", conf.Cursor.Value)
	} else {
		conf.Cursor.Value = ""
	}

	sync := &watcher{
		publisher: pub,
		conf:      conf,
		startSync: make(chan bool),
		syncNow:   make(chan bool),
	}

	// apply our supplied options
	for _, opt := range opts {
		opt(sync)
	}

	// if no logger option has been set use the null logger
	if sync.logger == nil {
		sync.logger = log.NewNopLogger()
	}

	if err := sync.updateClient(pub); err != nil {
		return nil, err
	}

	saveCursor := func() {
		if err := conf.Save(); err != nil {
			level.Info(logger).Log("err", err, "msg", "saving cursor")
			return
		}
		level.Info(logger).Log("msg", "saved DEP config", "cursor", conf.Cursor.Value)
	}

	go func() {
		defer saveCursor()
		if sync.client == nil {
			// block until we have a DEP client to start sync process
			level.Info(logger).Log("msg", "waiting for DEP token to be added before starting sync")
			<-sync.startSync
		}
		if err := sync.Run(); err != nil {
			level.Info(logger).Log("err", err, "msg", "DEP watcher failed")
		}
	}()
	return sync, nil
}

func (w *watcher) updateClient(pubsub pubsub.Subscriber) error {
	tokenAdded, err := pubsub.Subscribe(context.TODO(), "token-events", deptoken.DEPTokenTopic)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event := <-tokenAdded:
				var token deptoken.DEPToken
				if err := json.Unmarshal(event.Message, &token); err != nil {
					level.Info(w.logger).Log("err", err, "msg", "unmarshalling tokenAdd to token")
					continue
				}

				client, err := token.Client()
				if err != nil {
					level.Info(w.logger).Log("err", err, "msg", "creating new DEP client")
					continue
				}

				w.mtx.Lock()
				w.client = client
				w.mtx.Unlock()
				go func() { w.startSync <- true }() // unblock Run
			}
		}
	}()
	return nil
}

func (w *watcher) SyncNow() {
	w.syncNow <- true
}

// TODO this needs to be a proper error in the micromdm/dep package.
func isCursorExhausted(err error) bool {
	return strings.Contains(err.Error(), "EXHAUSTED_CURSOR")
}

func isCursorExpired(err error) bool {
	return strings.Contains(err.Error(), "EXPIRED_CURSOR")
}

func (w *watcher) Run() error {
	ticker := time.NewTicker(syncDuration).C
FETCH:
	for {
		resp, err := w.client.FetchDevices(dep.Limit(100), dep.Cursor(w.conf.Cursor.Value))
		if err != nil && isCursorExhausted(err) {
			goto SYNC
		} else if err != nil {
			return err
		}
		level.Info(w.logger).Log("msg", "DEP fetch", "more", resp.MoreToFollow, "cursor", resp.Cursor, "fetched", resp.FetchedUntil)
		w.conf.Cursor = cursor{Value: resp.Cursor, CreatedAt: time.Now()}
		if err := w.conf.Save(); err != nil {
			return errors.Wrap(err, "saving cursor from fetch")
		}
		e := NewEvent(resp.Devices)
		data, err := MarshalEvent(e)
		if err != nil {
			return err
		}
		if err := w.publisher.Publish(context.TODO(), SyncTopic, data); err != nil {
			return err
		}
		if !resp.MoreToFollow {
			goto SYNC
		}
	}

SYNC:
	for {
		resp, err := w.client.SyncDevices(w.conf.Cursor.Value, dep.Cursor(w.conf.Cursor.Value))
		if err != nil && isCursorExpired(err) {
			w.conf.Cursor.Value = ""
			goto FETCH
		} else if err != nil {
			return err
		}
		if len(resp.Devices) != 0 {
			level.Info(w.logger).Log("msg", "DEP sync", "more", resp.MoreToFollow, "cursor", resp.Cursor, "fetched", resp.FetchedUntil)
		}
		w.conf.Cursor = cursor{Value: resp.Cursor, CreatedAt: time.Now()}
		if err := w.conf.Save(); err != nil {
			return errors.Wrap(err, "saving cursor from sync")
		}
		if len(resp.Devices) > 0 {
			e := NewEvent(resp.Devices)
			data, err := MarshalEvent(e)
			if err != nil {
				return err
			}
			if err := w.publisher.Publish(context.TODO(), SyncTopic, data); err != nil {
				return err
			}
		}
		if !resp.MoreToFollow {
			select {
			case <-ticker:
			case <-w.syncNow:
				level.Info(w.logger).Log("msg", "explicit DEP sync requested")
			}
		}
	}
}
