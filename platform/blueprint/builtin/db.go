package builtin

import (
	"context"
	"fmt"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/platform/blueprint"
	"github.com/micromdm/micromdm/platform/profile"
)

const (
	BlueprintBucket      = "mdm.Blueprint"
	blueprintIndexBucket = "mdm.BlueprintIdx"
)

type DB struct {
	*bolt.DB
	profDB profile.Store
}

func NewDB(
	db *bolt.DB,
	profDB profile.Store,
) (*DB, error) {
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
		DB:     db,
		profDB: profDB,
	}
	return datastore, nil
}

func (db *DB) List() ([]blueprint.Blueprint, error) {
	// TODO add filter/limit with ForEach
	var blueprints []blueprint.Blueprint
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlueprintBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var bp blueprint.Blueprint
			if err := blueprint.UnmarshalBlueprint(v, &bp); err != nil {
				return err
			}
			blueprints = append(blueprints, bp)
		}
		return nil
	})
	return blueprints, err
}

func (db *DB) Save(bp *blueprint.Blueprint) error {
	ctx := context.TODO()
	err := bp.Verify()
	if err != nil {
		return err
	}
	check_bp, err := db.BlueprintByName(bp.Name)
	if err != nil && !isNotFound(err) {
		return err
	}
	if err == nil && bp.UUID != check_bp.UUID {
		return fmt.Errorf("Blueprint not saved: same name %s exists", bp.Name)
	}
	// verify that each Profile ID represents a profile we know about
	for _, p := range bp.ProfileIdentifiers {
		if _, err := db.profDB.ProfileById(ctx, p); err != nil {
			if profile.IsNotFound(err) {
				return fmt.Errorf("Profile ID %s in Blueprint %s does not exist", p, bp.Name)
			}
			return errors.Wrap(err, "fetching profile")
		}
	}
	tx, err := db.DB.Begin(true)
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}
	bkt := tx.Bucket([]byte(BlueprintBucket))
	if bkt == nil {
		return fmt.Errorf("bucket %q not found!", BlueprintBucket)
	}
	bpproto, err := blueprint.MarshalBlueprint(bp)
	if err != nil {
		return errors.Wrap(err, "marshalling blueprint")
	}
	idxBucket := tx.Bucket([]byte(blueprintIndexBucket))
	if idxBucket == nil {
		return fmt.Errorf("bucket %v not found!", idxBucket)
	}
	key := []byte(bp.Name)
	if err := idxBucket.Put(key, []byte(bp.UUID)); err != nil {
		return errors.Wrap(err, "put blueprint idx to boltdb")
	}

	key = []byte(bp.UUID)
	if err := bkt.Put(key, bpproto); err != nil {
		return errors.Wrap(err, "put blueprint to boltdb")
	}
	return tx.Commit()
}

func (db *DB) BlueprintByName(name string) (*blueprint.Blueprint, error) {
	var bp blueprint.Blueprint
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
		return blueprint.UnmarshalBlueprint(v, &bp)
	})
	if err != nil {
		return nil, err
	}
	return &bp, nil
}

func (db *DB) BlueprintsByApplyAt(ctx context.Context, name string) ([]blueprint.Blueprint, error) {
	var bps []blueprint.Blueprint
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlueprintBucket))
		c := b.Cursor()
		// TODO: fix this to use an index of ApplyAt strings mapping to
		// an array of Blueprints or other more efficient means. Looping
		// over every blueprint is quite inefficient!
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var bp blueprint.Blueprint
			err := blueprint.UnmarshalBlueprint(v, &bp)
			if err != nil {
				fmt.Println("could not Unmarshal Blueprint")
				continue
			}
			for _, n := range bp.ApplyAt {
				if strings.ToLower(n) == strings.ToLower(name) {
					bps = append(bps, bp)
					break
				}
			}
		}
		return nil
	})
	return bps, err
}

func (db *DB) Delete(name string) error {
	bp, err := db.BlueprintByName(name)
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		// TODO: reformulate into a transaction?
		b := tx.Bucket([]byte(BlueprintBucket))
		i := tx.Bucket([]byte(blueprintIndexBucket))
		err := i.Delete([]byte(bp.Name))
		if err != nil {
			return err
		}
		err = b.Delete([]byte(bp.UUID))
		if err != nil {
			return err
		}
		return nil
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

func isNotFound(err error) bool {
	if _, ok := err.(*notFound); ok {
		return true
	}
	return false
}
