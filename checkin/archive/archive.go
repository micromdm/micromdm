package archive

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/groob/plist"
	nsq "github.com/nsqio/go-nsq"

	"github.com/micromdm/micromdm/event"
)

const (
	// channel is the NSQ channel to listen on.
	channel = "archiver"

	// checkinBucket is the name of the *bolt.DB bucket
	// to archive events in.
	checkinBucket = "checkin"
)

func Persist(stop <-chan bool, boltDB *bolt.DB) error {
	cfg := nsq.NewConfig()
	db, err := newBoltStore(boltDB)
	if err != nil {
		return err
	}
	authenticate, err := nsq.NewConsumer("mdm.Authenticate", channel, cfg)
	if err != nil {
		return err
	}
	tokenUpdate, err := nsq.NewConsumer("mdm.TokenUpdate", channel, cfg)
	if err != nil {
		return err
	}
	checkout, err := nsq.NewConsumer("mdm.Checkout", channel, cfg)
	if err != nil {
		return err
	}
	errs := make(chan error)
	go func() {
		authenticate.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
			return db.SaveEvent(m.Body)
		}))

		if err := authenticate.ConnectToNSQD("localhost:4150"); err != nil {
			errs <- err
		}
	}()
	go func() {
		tokenUpdate.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
			return db.SaveEvent(m.Body)
		}))

		if err := tokenUpdate.ConnectToNSQD("localhost:4150"); err != nil {
			errs <- err
		}
	}()
	go func() {
		checkout.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
			return db.SaveEvent(m.Body)
		}))

		if err := checkout.ConnectToNSQD("localhost:4150"); err != nil {
			errs <- err
		}
	}()

	for {
		select {
		case err := <-errs:
			return err

		case <-stop:
			return nil
		}
	}
}

type store struct {
	*bolt.DB
}

func newBoltStore(db *bolt.DB) (*store, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(checkinBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &store{db}, nil
}

func (db *store) SaveEvent(ev []byte) error {
	var checkin event.CheckinEvent
	if err := plist.Unmarshal(ev, &checkin); err != nil {
		return err
	}
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(checkinBucket))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found!", checkinBucket)
		}
		k := fmt.Sprintf("%d", checkin.Time)
		return bucket.Put([]byte(k), ev)
	})
	return err
}
