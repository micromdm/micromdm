package sync

import (
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

type config struct {
	*bolt.DB
	Cursor cursor `json:"cursor"`
}

type cursor struct {
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

type AutoAssigner struct {
	Filter      string `json:"filter"`
	ProfileUUID string `json:"profile_uuid"`
}

func (cfg *config) Save() error {
	err := cfg.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(ConfigBucket))
		if err != nil {
			return err
		}
		v, err := json.Marshal(cfg)
		if err != nil {
			return err
		}
		return bkt.Put([]byte("configuration"), v)
	})
	return errors.Wrap(err, "saving dep sync cursor")
}

func (cfg *config) saveAutoAssigner(assigner *AutoAssigner) error {
	if assigner.Filter != "*" {
		return errors.New("only '*' filter auto-assigners supported")
	}
	err := cfg.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(AutoAssignBucket))
		if err != nil {
			return err
		}
		return b.Put([]byte(assigner.Filter), []byte(assigner.ProfileUUID))
	})
	return errors.Wrap(err, "saving auto-assigner")
}

func (cfg *config) loadAutoAssigners() ([]*AutoAssigner, error) {
	assigners := []*AutoAssigner{}
	err := cfg.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(AutoAssignBucket))
		if b == nil { // bucket doesn't exist yet
			return nil
		}

		return b.ForEach(func(k, v []byte) error {
			assigners = append(assigners, &AutoAssigner{
				Filter:      string(k),
				ProfileUUID: string(v),
			})
			return nil
		})
	})
	return assigners, errors.Wrap(err, "loading auto-assigners")
}

func (cfg *config) deleteAutoAssigner(filter string) error {
	return cfg.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(AutoAssignBucket))
		if b == nil { // bucket doesn't exist yet
			return nil
		}
		return b.Delete([]byte(filter))
	})
}

func LoadConfig(db *bolt.DB) (*config, error) {
	conf := config{DB: db}
	err := db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(ConfigBucket))
		if err != nil {
			return err
		}

		v := bkt.Get([]byte("configuration"))
		if v == nil {
			return nil
		}
		if err := json.Unmarshal(v, &conf); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
