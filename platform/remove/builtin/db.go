package builtin

import (
	"context"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/micromdm/micromdm/platform/remove"
	"github.com/pkg/errors"
)

const RemoveBucket = "mdm.RemoveDevice"

type DB struct {
	*bolt.DB
}

func NewDB(db *bolt.DB) (*DB, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(RemoveBucket))
		return err
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating %s bucket", RemoveBucket)
	}
	datastore := &DB{
		DB: db,
	}
	return datastore, nil
}

func (db *DB) DeviceByUDID(ctx context.Context, udid string) (*remove.Device, error) {
	var dev remove.Device
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(RemoveBucket))
		v := b.Get([]byte(udid))
		if v == nil {
			return &notFound{"Device", fmt.Sprintf("udid %s", udid)}
		}
		return remove.UnmarshalDevice(v, &dev)
	})
	return &dev, errors.Wrap(err, "remove: get device by udid")
}

func (db *DB) Save(ctx context.Context, dev *remove.Device) error {
	tx, err := db.DB.Begin(true)
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}
	bkt := tx.Bucket([]byte(RemoveBucket))
	if bkt == nil {
		return fmt.Errorf("bucket %q not found!", RemoveBucket)
	}
	pb, err := remove.MarshalDevice(dev)
	if err != nil {
		return errors.Wrap(err, "marshalling Device")
	}
	key := []byte(dev.UDID)
	if err := bkt.Put(key, pb); err != nil {
		return errors.Wrap(err, "put device to boltdb")
	}
	return tx.Commit()
}

func (db *DB) Delete(ctx context.Context, udid string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(RemoveBucket))
		v := b.Get([]byte(udid))
		if v == nil {
			return &notFound{"Device", fmt.Sprintf("udid %s", udid)}
		}
		return b.Delete([]byte(udid))
	})
	return errors.Wrapf(err, "delete device with udid %s", udid)
}

type notFound struct {
	ResourceType string
	Message      string
}

func (e *notFound) Error() string {
	return fmt.Sprintf("not found: %s %s", e.ResourceType, e.Message)
}

func (e *notFound) NotFound() bool {
	return true
}
