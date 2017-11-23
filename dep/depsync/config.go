package depsync

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

type config struct {
	*bolt.DB
	Cursor cursor `json:"cursor"`
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
