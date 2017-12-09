package builtin

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/micromdm/micromdm/platform/profile"
	"github.com/pkg/errors"
)

const (
	ProfileBucket = "mdm.Profile"
)

type DB struct {
	*bolt.DB
}

func NewDB(db *bolt.DB) (*DB, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(ProfileBucket))
		return err
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating %s bucket", ProfileBucket)
	}
	datastore := &DB{
		DB: db,
	}
	return datastore, nil
}

func (db *DB) List() ([]profile.Profile, error) {
	// TODO add filter/limit with ForEach
	var list []profile.Profile
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ProfileBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var p profile.Profile
			if err := profile.UnmarshalProfile(v, &p); err != nil {
				return err
			}
			list = append(list, p)
		}
		return nil
	})
	return list, err
}

func (db *DB) Save(p *profile.Profile) error {
	err := p.Validate()
	if err != nil {
		return err
	}
	tx, err := db.DB.Begin(true)
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}
	bkt := tx.Bucket([]byte(ProfileBucket))
	if bkt == nil {
		return fmt.Errorf("bucket %q not found!", ProfileBucket)
	}
	pproto, err := profile.MarshalProfile(p)
	if err != nil {
		return errors.Wrap(err, "marshalling profile")
	}
	if err := bkt.Put([]byte(p.Identifier), pproto); err != nil {
		return errors.Wrap(err, "put profile to boltdb")
	}
	return tx.Commit()
}

func (db *DB) ProfileById(id string) (*profile.Profile, error) {
	var p profile.Profile
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ProfileBucket))
		v := b.Get([]byte(id))
		if v == nil {
			return &notFound{"Profile", fmt.Sprintf("id %s", id)}
		}
		return profile.UnmarshalProfile(v, &p)
	})
	return &p, err
}

func (db *DB) Delete(id string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ProfileBucket))
		v := b.Get([]byte(id))
		if v == nil {
			return &notFound{"Profile", fmt.Sprintf("id %s", id)}
		}
		return b.Delete([]byte(id))
	})
	return err
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
