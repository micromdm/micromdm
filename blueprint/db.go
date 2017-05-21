package blueprint

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

const (
	BlueprintBucket      = "mdm.Blueprint"
	blueprintIndexBucket = "mdm.BlueprintIdx"
)

type DB struct {
	*bolt.DB
}

func NewDB(db *bolt.DB) (*DB, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(blueprintIndexBucket))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(BlueprintBucket))
		return err
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating %s bucket", BlueprintBucket)
	}
	datastore := &DB{
		DB: db,
	}
	return datastore, nil
}

func (db *DB) Create(bp *Blueprint) error {
	_, err := db.BlueprintByName(bp.Name)
	if err != nil && isNotFound(err) {
		return errors.New("blueprint must have a unique name")
	} else if err != nil {
		return errors.Wrap(err, "creating blueprint")
	}

	return db.Save(bp)
}

func (db *DB) List() ([]Blueprint, error) {
	// TODO add filter/limit with ForEach
	var blueprints []Blueprint
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlueprintBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var bp Blueprint
			if err := UnmarshalBlueprint(v, &bp); err != nil {
				return err
			}
			blueprints = append(blueprints, bp)
		}
		return nil
	})
	return blueprints, err
}

func (db *DB) Save(bp *Blueprint) error {
	if bp.Name == "" || bp.UUID == "" {
		return errors.New("cannot Save: blueprint must have Name and UUID")
	}
	tx, err := db.DB.Begin(true)
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}
	bkt := tx.Bucket([]byte(BlueprintBucket))
	if bkt == nil {
		return fmt.Errorf("bucket %q not found!", BlueprintBucket)
	}
	bpproto, err := MarshalBlueprint(bp)
	if err != nil {
		return errors.Wrap(err, "marshalling blueprint")
	}
	indexes := []string{bp.UUID, bp.Name}
	idxBucket := tx.Bucket([]byte(blueprintIndexBucket))
	if idxBucket == nil {
		return fmt.Errorf("bucket %q not found!", idxBucket)
	}
	for _, idx := range indexes {
		if idx == "" {
			continue
		}
		key := []byte(idx)
		if err := idxBucket.Put(key, []byte(bp.UUID)); err != nil {
			return errors.Wrap(err, "put blueprint idx to boltdb")
		}
	}

	key := []byte(bp.UUID)
	if err := bkt.Put(key, bpproto); err != nil {
		return errors.Wrap(err, "put blueprint to boltdb")
	}
	return tx.Commit()
}

func (db *DB) BlueprintByName(name string) (*Blueprint, error) {
	var bp Blueprint
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlueprintBucket))
		ib := tx.Bucket([]byte(blueprintIndexBucket))
		idx := ib.Get([]byte(name))
		if idx == nil {
			return &notFound{"Blueprint", fmt.Sprintf("name %s", name)}
		}
		v := b.Get(idx)
		if idx == nil {
			return &notFound{"Blueprint", fmt.Sprintf("uuid %s", string(idx))}
		}
		return UnmarshalBlueprint(v, &bp)
	})
	if err != nil {
		return nil, err
	}
	return &bp, nil
}

func (db *DB) BlueprintsByWhenTag(when string) ([]*Blueprint, error) {
	var bps []*Blueprint
	err := db.View(func(tx *bolt.Tx) error {
		// TODO: SUPER inefficient. demo purposes only. right thing to do is
		// probably to make an index on the When field and iterate over that
		// somehow.
		b := tx.Bucket([]byte(BlueprintBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var bp Blueprint
			err := UnmarshalBlueprint(v, &bp)
			if err != nil {
				// skipping for now
				fmt.Println("could not unmarshal blueprint")
				continue
			}
			if bp.When == when {
				bps = append(bps, &bp)
			}
		}
		return nil
	})
	return bps, err
}

type notFound struct {
	ResourceType string
	Message      string
}

func (e *notFound) Error() string {
	return fmt.Sprintf("not found: %s %s", e.ResourceType, e.Message)
}

func isNotFound(err error) bool {
	if _, ok := err.(*notFound); ok {
		return true
	}
	return false
}
