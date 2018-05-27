package builtin

import (
	"context"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/micromdm/micromdm/mdm"
	"github.com/micromdm/micromdm/platform/pubsub"
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

func NewDB(db *bolt.DB, pubsubSvc pubsub.PublishSubscriber, logger log.Logger) (*DB, error) {
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
		DB:     db,
		logger: logger,
	}
	if pubsubSvc == nil { // don't start the poller without pubsub.
		return datastore, nil
	}
	if err := datastore.pollCheckin(pubsubSvc); err != nil {
		return nil, err
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

func (db *DB) pollCheckin(pubsubSvc pubsub.PublishSubscriber) error {
	tokenUpdateEvents, err := pubsubSvc.Subscribe(context.TODO(), "users", mdm.TokenUpdateTopic)
	if err != nil {
		return errors.Wrapf(err,
			"subscribing devices to %s topic", mdm.TokenUpdateTopic)
	}
	go func() {
		for {
			select {
			case e := <-tokenUpdateEvents:
				event, err := unmarshalCheckin(e)
				if err != nil {
					level.Info(db.logger).Log("err", err, "msg", "unmarshal TokenUpdate event in user db")
					break
				}
				if event.Command.UserID == "" {
					break // only interested in user commands
				}
				newUser := new(user.User)
				byGUID, err := db.UserByUserID(event.Command.UserID)
				if err != nil && !isNotFound(err) {
					level.Info(db.logger).Log("err", err, "msg", "get user from DB")
					break
				}
				if err == nil && byGUID != nil {
					newUser = byGUID
				}
				if newUser.UUID == "" {
					if err := db.DeleteDeviceUsers(event.Command.UDID); err != nil {
						level.Info(db.logger).Log(
							"err", err,
							"msg", "delete existing user before creating new one",
						)
					}
					newUser.UUID = uuid.NewV4().String()
				}
				newUser.UDID = event.Command.UDID
				newUser.UserID = event.Command.UserID
				newUser.UserLongname = event.Command.UserLongName
				newUser.UserShortname = event.Command.UserShortName
				newUser.AuthToken = event.Command.Token.String()
				if err := db.Save(newUser); err != nil {
					level.Info(db.logger).Log("err", err, "msg", "update user from TokenUpdate")
					break
				}
			}
		}
	}()

	return nil
}

func unmarshalCheckin(event pubsub.Event) (mdm.CheckinEvent, error) {
	var ev mdm.CheckinEvent
	if err := mdm.UnmarshalCheckinEvent(event.Message, &ev); err != nil {
		return mdm.CheckinEvent{}, err
	}
	return ev, nil
}

func isNotFound(err error) bool {
	cause := errors.Cause(err)
	if _, ok := cause.(*notFound); ok {
		return true
	}
	return false
}
