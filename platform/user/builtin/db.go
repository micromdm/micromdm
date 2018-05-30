package builtin

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/platform/user"
)

const (
	UserBucket = "mdm.Users"

	userIndexBucket = "mdm.UserIdx"
)

type DB struct {
	*bolt.DB
	logger log.Logger
}

func NewDB(db *bolt.DB) (*DB, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(userIndexBucket))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(UserBucket))
		return err
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating %s bucket", UserBucket)
	}

	datastore := &DB{
		DB: db,
	}
	return datastore, nil
}

func (db *DB) List() ([]user.User, error) {
	var users []user.User
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UserBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var u user.User
			if err := user.UnmarshalUser(v, &u); err != nil {
				return err
			}
			users = append(users, u)
		}
		return nil
	})
	return users, errors.Wrap(err, "list users")
}

func (db *DB) Save(u *user.User) error {
	tx, err := db.DB.Begin(true)
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}
	bkt := tx.Bucket([]byte(UserBucket))
	if bkt == nil {
		return fmt.Errorf("bucket %q not found!", UserBucket)
	}
	userpb, err := user.MarshalUser(u)
	if err != nil {
		return errors.Wrap(err, "marshalling user")
	}

	// store an array of indices to reference the UUID, which will be the
	// key used to store the actual user.
	indexes := []string{u.UDID, u.UserID}
	idxBucket := tx.Bucket([]byte(userIndexBucket))
	if idxBucket == nil {
		return fmt.Errorf("bucket %q not found!", userIndexBucket)
	}
	for _, idx := range indexes {
		if idx == "" {
			continue
		}
		key := []byte(idx)
		if err := idxBucket.Put(key, []byte(u.UUID)); err != nil {
			return errors.Wrap(err, "user userIdx in boltdb")
		}
	}

	key := []byte(u.UUID)
	if err := bkt.Put(key, userpb); err != nil {
		return errors.Wrap(err, "store user in boltdb")
	}
	return tx.Commit()
}

func (db *DB) User(uuid string) (*user.User, error) {
	var u user.User
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UserBucket))
		v := b.Get([]byte(uuid))
		if v == nil {
			return &notFound{"User", fmt.Sprintf("uuid %s", uuid)}
		}
		return user.UnmarshalUser(v, &u)
	})
	if err != nil {
		return nil, errors.Wrap(err, "get user by uuid from bolt")
	}
	return &u, nil
}

func (db *DB) UserByUserID(userID string) (*user.User, error) {
	var u user.User
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UserBucket))
		ib := tx.Bucket([]byte(userIndexBucket))
		idx := ib.Get([]byte(userID))
		if idx == nil {
			return &notFound{"User", fmt.Sprintf("user id %s", userID)}
		}
		v := b.Get(idx)
		if idx == nil {
			return &notFound{"User", fmt.Sprintf("uuid %s", string(idx))}
		}
		return user.UnmarshalUser(v, &u)
	})
	if err != nil {
		return nil, errors.Wrap(err, "get user by user id from bolt")
	}
	return &u, nil
}

func (db *DB) DeviceUsers(udid string) ([]user.User, error) {
	var users []user.User
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UserBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var u user.User
			if err := user.UnmarshalUser(v, &u); err != nil {
				return errors.Wrap(err, "unmarshal user for DeviceUsers")
			}
			if u.UDID == udid {
				users = append(users, u)
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "get device users")
	}
	return users, nil
}

func (db *DB) DeleteDeviceUsers(udid string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UserBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var u user.User
			if err := user.UnmarshalUser(v, &u); err != nil {
				return errors.Wrap(err, "unmarshal user for DeviceUsers")
			}
			if u.UDID != udid {
				continue
			}
			if err := b.Delete(k); err != nil {
				return errors.Wrapf(err, "delete user %s from device %s", u.UserID, udid)
			}
		}
		return nil
	})
	return errors.Wrapf(err, "delete users for UDID %s", udid)
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
