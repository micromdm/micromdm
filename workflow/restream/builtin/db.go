package builtin

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/workflow/restream"
)

const RestreamBucket = "restream.Bucket"

type DB struct {
	*bolt.DB
}

func New(db *bolt.DB) (*DB, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(RestreamBucket))
		return err
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating %s bucket", RestreamBucket)
	}
	datastore := &DB{
		DB: db,
	}
	return datastore, nil
}

func (db *DB) Save(ctx context.Context, event restream.Event) error {
	data, err := json.Marshal(&event)
	if err != nil {
		return err
	}

	tx, err := db.DB.Begin(true)
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}
	bkt := tx.Bucket([]byte(RestreamBucket))
	if bkt == nil {
		return fmt.Errorf("bucket %q not found!", RestreamBucket)
	}
	key := []byte(event.ID)
	if err := bkt.Put(key, data); err != nil {
		return errors.Wrap(err, "put Restream event to boltdb")
	}
	return tx.Commit()
}

func (db *DB) Event(ctx context.Context, id string) (restream.Event, error) {
	var ev restream.Event
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(RestreamBucket))
		v := b.Get([]byte(id))
		if v == nil {
			return &notFound{"Restream Event", fmt.Sprintf("id %s", id)}
		}
		return json.Unmarshal(v, &ev)
	})
	return ev, err
}

type notFound struct {
	ResourceType string
	Message      string
}

func (e *notFound) Error() string {
	return fmt.Sprintf("not found: %s %s", e.ResourceType, e.Message)
}
