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
	"github.com/pkg/errors"

	conf "github.com/micromdm/micromdm/platform/config"
	"github.com/micromdm/micromdm/platform/pubsub"
)

const (
	SyncTopic        = "mdm.DepSync"
	ConfigBucket     = "mdm.DEPConfig"
	AutoAssignBucket = "mdm.DEPAutoAssign"

	syncDuration        = 30 * time.Minute
	cursorValidDuration = 7 * 24 * time.Hour
)

type Syncer interface {
	SyncNow()
	GetConfig() *config // TODO: #302
}

type AutoAssigner struct {
	Filter      string `json:"filter"`
	ProfileUUID string `json:"profile_uuid"`
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
	expiration := time.Now().Add(cursorValidDuration)
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
		logger:    log.NewNopLogger(),
		publisher: pub,
		conf:      conf,
		startSync: make(chan bool),
		syncNow:   make(chan bool),
	}

	for _, opt := range opts {
		opt(sync)
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
	tokenAdded, err := pubsub.Subscribe(context.TODO(), "token-events", conf.DEPTokenTopic)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event := <-tokenAdded:
				var token conf.DEPToken
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
	if w.client == nil {
		level.Info(w.logger).Log("msg", "waiting for DEP token to be added before starting sync")
		return
	}
	w.syncNow <- true
}

func (w *watcher) GetConfig() *config {
	return w.conf
}

// TODO this needs to be a proper error in the micromdm/dep package.
func isCursorExhausted(err error) bool {
	return strings.Contains(err.Error(), "EXHAUSTED_CURSOR")
}

func isCursorExpired(err error) bool {
	return strings.Contains(err.Error(), "EXPIRED_CURSOR")
}

// Process DEP messages and pull out filter-matching serial numbers
// associated to profile UUIDs for auto-assignment.
func (w *watcher) filteredAutoAssignments(devices []dep.Device) (map[string][]string, error) {
	// load auto-assigners every run to make sure we get the latest set of
	// auto-assigner profile UUIDs/filters. Note this makes every *watcher
	// (i.e. every DEP sync instance) share the current DB set of auto-
	// assigners. perhaps to refactor to be more separated.
	assigners, err := w.conf.loadAutoAssigners()
	if err != nil {
		return nil, err
	}
	assigned := make(map[string][]string)
	// skip looping over serials if we have no autoassigners
	if len(assigners) < 1 {
		return assigned, nil
	}
	for _, d := range devices {
		// only process DEP "added" OpType messages
		if d.OpType != "added" {
			continue
		}
		// filter our devices by our assigner filters and get list of
		// which devices are to be assigned to which profiles
		for _, assigner := range assigners {
			if assigner.Filter == "*" { // only supported filter type right now
				if serials, ok := assigned[assigner.ProfileUUID]; ok {
					assigned[assigner.ProfileUUID] = append(serials, d.SerialNumber)
				} else {
					assigned[assigner.ProfileUUID] = []string{d.SerialNumber}
				}
			}
		}

	}
	return assigned, nil
}

func (w *watcher) processAutoAssign(devices []dep.Device) error {
	assignments, err := w.filteredAutoAssignments(devices)
	if err != nil {
		return err
	}

	for profileUUID, serials := range assignments {
		resp, err := w.client.AssignProfile(profileUUID, serials)
		if err != nil {
			level.Info(w.logger).Log(
				"err", err,
				"msg", "auto-assign error assigning serials to profile",
				"profile", profileUUID,
			)
			continue
		}
		// count our results for logging
		resultCounts := map[string]int{
			"SUCCESS":        0,
			"NOT_ACCESSIBLE": 0,
			"FAILED":         0,
		}
		for _, result := range resp.Devices {
			if ct, ok := resultCounts[result]; ok {
				// NOTE: we're logging _only_ the above pre-defined result types
				resultCounts[result] = ct + 1
			}
		}
		// TODO: alternate strategy is to log all failed devices
		// TODO: handle/requeue failed devices?
		level.Info(w.logger).Log(
			"msg", "DEP auto-assigned",
			"profile", profileUUID,
			"success", resultCounts["SUCCESS"],
			"not_accessible", resultCounts["NOT_ACCESSIBLE"],
			"failed", resultCounts["FAILED"],
		)
	}

	return nil
}

func (w *watcher) publishAndProcessDevices(devices []dep.Device) error {
	e := NewEvent(devices)
	data, err := MarshalEvent(e)
	if err != nil {
		return err
	}
	err = w.publisher.Publish(context.TODO(), SyncTopic, data)
	if err != nil {
		return err
	}

	// TODO: instead of directly kicking off the auto-assigner process
	// consider placing a subscriber on the DEP pubsub topic. The same
	// information gets marshalled but it allows us the future
	// flexibility to separate out that component if we desired.
	go func() {
		err := w.processAutoAssign(devices)
		if err != nil {
			level.Info(w.logger).Log("err", err, "msg", "auto-assign error")
		}
	}()
	return nil
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
		level.Info(w.logger).Log(
			"msg", "DEP fetch",
			"more", resp.MoreToFollow,
			"cursor", resp.Cursor,
			"fetched", resp.FetchedUntil,
			"devices", len(resp.Devices),
		)
		w.conf.Cursor = cursor{Value: resp.Cursor, CreatedAt: time.Now()}
		if err := w.conf.Save(); err != nil {
			return errors.Wrap(err, "saving cursor from fetch")
		}
		if err := w.publishAndProcessDevices(resp.Devices); err != nil {
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
		level.Info(w.logger).Log(
			"msg", "DEP sync",
			"more", resp.MoreToFollow,
			"cursor", resp.Cursor,
			"fetched", resp.FetchedUntil,
			"devices", len(resp.Devices),
		)
		w.conf.Cursor = cursor{Value: resp.Cursor, CreatedAt: time.Now()}
		if err := w.conf.Save(); err != nil {
			return errors.Wrap(err, "saving cursor from sync")
		}
		if err := w.publishAndProcessDevices(resp.Devices); err != nil {
			return err
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
