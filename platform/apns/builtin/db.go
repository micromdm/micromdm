package builtin

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/platform/apns"
	"github.com/micromdm/micromdm/platform/pubsub"
)

const PushBucket = "mdm.PushInfo"

type DB struct {
	*bolt.DB
}

func NewDB(db *bolt.DB, sub pubsub.Subscriber) (*DB, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(PushBucket))
		return err
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating %s bucket", PushBucket)
	}
	datastore := &DB{
		DB: db,
	}
	return datastore, nil
}

type FireStoreDB struct {
	FirestoreCli *firestore.Client
}

func NewFirestoreDB(ctx context.Context) (*FireStoreDB, error) {
	config := &firebase.Config{
		ProjectID: "micromdm-df039", // replace with your Project ID
	}
	app, err := firebase.NewApp(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %v", err)
	}

	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating firestore client: %v", err)
	}

	return &FireStoreDB{FirestoreCli: firestoreClient}, nil
}

func (db *FireStoreDB) SaveToFirestore(ctx context.Context, a *apns.PushInfo) error {
	data := map[string]interface{}{
		"UUID":       a.MDMTopic,
		"UDID":       a.UDID,
		"DeviceInfo": a.Token,
		"CreatedAt":  firestore.ServerTimestamp,
		// Include other fields here.
	}

	_, err := db.FirestoreCli.Collection("apns").Doc(a.UDID).Set(ctx, data)
	if err != nil {
		return fmt.Errorf("adding apns to Firestore: %v", err)
	}

	return nil
}

type notFound struct {
	ResourceType string
	Message      string
}

func (e *notFound) Error() string {
	return fmt.Sprintf("not found: %s %s", e.ResourceType, e.Message)
}

func (db *DB) PushInfo(ctx context.Context, udid string) (*apns.PushInfo, error) {
	var info apns.PushInfo
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PushBucket))
		v := b.Get([]byte(udid))
		if v == nil {
			return &notFound{"PushInfo", fmt.Sprintf("udid %s", udid)}
		}
		return apns.UnmarshalPushInfo(v, &info)
	})
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (db *DB) Save(ctx context.Context, info *apns.PushInfo) error {
	tx, err := db.DB.Begin(true)
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}
	bkt := tx.Bucket([]byte(PushBucket))
	if bkt == nil {
		return fmt.Errorf("bucket %q not found!", PushBucket)
	}
	pushproto, err := apns.MarshalPushInfo(info)
	if err != nil {
		return errors.Wrap(err, "marshalling PushInfo")
	}
	key := []byte(info.UDID)
	if err := bkt.Put(key, pushproto); err != nil {
		return errors.Wrap(err, "put PushInfo to boltdb")
	}
	// Save to Firestore
	firestoreDB, err := NewFirestoreDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to create Firestore DB: %v", err)
	}

	err = firestoreDB.SaveToFirestore(ctx, info)
	if err != nil {
		return fmt.Errorf("failed to save device to Firestore: %v", err)
	}
	return tx.Commit()
}
