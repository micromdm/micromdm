package builtin

import (
	"context"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/mdm/checkin"
	"github.com/micromdm/micromdm/platform/apns"
	"github.com/micromdm/micromdm/platform/pubsub"
)

const PushBucket = "mdm.PushInfo"

type DB struct {
	*bolt.DB
}

func NewDB(db *bolt.DB, sub pubsub.Subscriber) (*DB, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(PushBucket))
		return err
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating %s bucket", PushBucket)
	}
	datastore := &DB{
		DB: db,
	}
	if err := datastore.pollCheckin(sub); err != nil {
		return nil, err
	}
	return datastore, nil
}

type notFound struct {
	ResourceType string
	Message      string
}

func (e *notFound) Error() string {
	return fmt.Sprintf("not found: %s %s", e.ResourceType, e.Message)
}

func (db *DB) PushInfo(udid string) (*apns.PushInfo, error) {
	var info apns.PushInfo
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PushBucket))
		v := b.Get([]byte(udid))
		if v == nil {
			return &notFound{"PushInfo", fmt.Sprintf("udid %s", udid)}
		}
		return apns.UnmarshalPushInfo(v, &info)
	})
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (db *DB) Save(info *apns.PushInfo) error {
	tx, err := db.DB.Begin(true)
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}
	bkt := tx.Bucket([]byte(PushBucket))
	if bkt == nil {
		return fmt.Errorf("bucket %q not found!", PushBucket)
	}
	pushproto, err := apns.MarshalPushInfo(info)
	if err != nil {
		return errors.Wrap(err, "marshalling PushInfo")
	}
	key := []byte(info.UDID)
	if err := bkt.Put(key, pushproto); err != nil {
		return errors.Wrap(err, "put PushInfo to boltdb")
	}
	return tx.Commit()
}

func (db *DB) pollCheckin(sub pubsub.Subscriber) error {
	tokenUpdateEvents, err := sub.Subscribe(context.TODO(), "push-info", checkin.TokenUpdateTopic)
	if err != nil {
		return errors.Wrapf(err,
			"subscribing push to %s topic", checkin.TokenUpdateTopic)
	}
	go func() {
		for {
			select {
			case event := <-tokenUpdateEvents:
				var ev checkin.Event
				if err := checkin.UnmarshalEvent(event.Message, &ev); err != nil {
					fmt.Println(err)
					continue
				}
				info := apns.PushInfo{
					UDID:      ev.Command.UDID,
					Token:     ev.Command.Token.String(),
					PushMagic: ev.Command.PushMagic,
					MDMTopic:  ev.Command.Topic,
				}
				if ev.Command.UserID != "" {
					// use the GUID if this is a user TokenUpdate.
					info.UDID = ev.Command.UserID
				}
				if err := db.Save(&info); err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Printf("updated pushinfo for udid %s\n", info.UDID)
			}
		}
	}()

	return nil
}
