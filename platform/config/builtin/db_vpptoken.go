package builtin

import (
	"context"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/micromdm/micromdm/platform/config"
)

const (
	vppTokenBucket = "mdm.VPPToken"
)

func (db *DB) AddVPPToken(sToken string, json []byte) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(vppTokenBucket))
		if err != nil {
			return err
		}
		return b.Put([]byte(sToken), json)
	})
	if err != nil {
		return err
	}
	err = db.Publisher.Publish(context.TODO(), config.VPPTokenTopic, json)
	return err
}

func (db *DB) DeleteVPPToken(sToken string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(vppTokenBucket))
		return b.Delete([]byte(sToken))
	})
	if err != nil {
		return err
	}
	err = db.Publisher.Publish(context.TODO(), config.VPPTokenTopic, []byte(sToken))
	return err
}

func (db *DB) VPPTokens() ([]config.VPPToken, error) {
	var result []config.VPPToken
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(vppTokenBucket))
		if b == nil {
			return nil
		}
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var vppToken config.VPPToken
			err := json.Unmarshal(v, &vppToken)
			if err != nil {
				// TODO: log problematic VPP token, or remove altogether?
				continue
			}
			result = append(result, vppToken)
		}
		return nil
	})
	return result, err
}
