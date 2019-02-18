package builtin

import (
	"encoding/json"
	"context"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/platform/dep/sync"
)

const (
	ConfigBucket     = "mdm.DEPConfig"
	AutoAssignBucket = "mdm.DEPAutoAssign"
)

type DB struct {
	*bolt.DB
}

func NewDB(db *bolt.DB) (*DB, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(ConfigBucket))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(AutoAssignBucket))
		return err
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating %s bucket", ConfigBucket)
	}
	datastore := &DB{DB: db}
	return datastore, nil
}

func (db *DB) LoadCursor(ctx context.Context) (*sync.Cursor, error) {
	var cursor = struct {
		Cursor sync.Cursor `json:"cursor"`
	}{}
	err := db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(ConfigBucket))
		v := bkt.Get([]byte("configuration"))
		if v == nil {
			return nil // TODO add notfound
		}
		err := json.Unmarshal(v, &cursor)
		return errors.Wrap(err, "unmarshal dep cursor")
	})
	if err != nil {
		return nil, errors.Wrap(err, "load cursor from bolt")
	}
	fmt.Println(cursor.Cursor)
	return &cursor.Cursor, nil
}

func (db *DB) SaveCursor(ctx context.Context, c sync.Cursor) error {
	err := db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(ConfigBucket))
		if err != nil {
			return err
		}
		// use anonymous struct to preserve old structure.
		var cursor = struct {
			Cursor sync.Cursor `json:"cursor"`
		}{Cursor: c}
		v, err := json.Marshal(&cursor)
		if err != nil {
			return err
		}
		return bkt.Put([]byte("configuration"), v)
	})
	return errors.Wrap(err, "saving dep sync cursor")
}

func (db *DB) SaveAutoAssigner(ctx context.Context, a *sync.AutoAssigner) error {
	if a.Filter != "*" {
		return errors.New("only '*' filter auto-assigners supported")
	}
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(AutoAssignBucket))
		if err != nil {
			return err
		}
		return b.Put([]byte(a.Filter), []byte(a.ProfileUUID))
	})
	return errors.Wrap(err, "saving auto-assigner")
}

func (db *DB) DeleteAutoAssigner(ctx context.Context, filter string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(AutoAssignBucket))
		if b == nil { // bucket doesn't exist yet
			return nil
		}
		return b.Delete([]byte(filter))
	})
}

func (db *DB) LoadAutoAssigners(ctx context.Context) ([]sync.AutoAssigner, error) {
	var aa []sync.AutoAssigner
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(AutoAssignBucket))
		if b == nil { // bucket doesn't exist yet
			return nil
		}

		return b.ForEach(func(k, v []byte) error {
			aa = append(aa, sync.AutoAssigner{
				Filter:      string(k),
				ProfileUUID: string(v),
			})
			return nil
		})
	})
	return aa, errors.Wrap(err, "loading auto-assigners")
}
