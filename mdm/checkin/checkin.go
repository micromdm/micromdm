package checkin

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/micromdm/mdm"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/micromdm/micromdm/platform/pubsub"
)

// CheckinBucket is the *bolt.DB bucket where checkins are archived.
const CheckinBucket = "mdm.Checkin.ARCHIVE"

// PubSub Topics where MDM Checkin events are published to.
const (
	AuthenticateTopic = "mdm.Authenticate"
	TokenUpdateTopic  = "mdm.TokenUpdate"
	CheckoutTopic     = "mdm.CheckOut"
)

type Checkin struct {
	db        *bolt.DB
	publisher pubsub.Publisher
	archiveFn func(int64, []byte) error
}

func New(db *bolt.DB, pub pubsub.Publisher) (*Checkin, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(CheckinBucket))
		return err
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating %s bucket", CheckinBucket)
	}
	svc := Checkin{
		db:        db,
		publisher: pub,
	}
	svc.archiveFn = svc.archive
	return &svc, nil
}

func (svc *Checkin) Authenticate(ctx context.Context, cmd mdm.CheckinCommand) error {
	if cmd.MessageType != "Authenticate" {
		return fmt.Errorf("expected Authenticate, got %s MessageType", cmd.MessageType)
	}
	fmt.Println(cmd)
	return svc.archiveAndPublish(AuthenticateTopic, cmd)
}

func (svc *Checkin) TokenUpdate(ctx context.Context, cmd mdm.CheckinCommand) error {
	if cmd.MessageType != "TokenUpdate" {
		return fmt.Errorf("expected TokenUpdate, got %s MessageType", cmd.MessageType)
	}
	fmt.Println(cmd)
	return svc.archiveAndPublish(TokenUpdateTopic, cmd)
}

func (svc *Checkin) CheckOut(ctx context.Context, cmd mdm.CheckinCommand) error {
	if cmd.MessageType != "CheckOut" {
		return fmt.Errorf("expected CheckOut, but got %s MessageType", cmd.MessageType)
	}
	fmt.Println(cmd)
	return svc.archiveAndPublish(CheckoutTopic, cmd)
}

// archive events to BoltDB bucket using timestamp as key to preserve order.
func (svc *Checkin) archive(nano int64, msg []byte) error {
	tx, err := svc.db.Begin(true)
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}
	defer tx.Rollback()

	bkt := tx.Bucket([]byte(CheckinBucket))
	if bkt == nil {
		return fmt.Errorf("bucket %q not found!", CheckinBucket)
	}
	key := []byte(fmt.Sprintf("%d", nano))
	if err := bkt.Put(key, msg); err != nil {
		return errors.Wrap(err, "put checkin event to boltdb")
	}
	return tx.Commit()
}

func (svc *Checkin) archiveAndPublish(topic string, cmd mdm.CheckinCommand) error {
	event := NewEvent(cmd)
	msg, err := MarshalEvent(event)

	if err != nil {
		return errors.Wrap(err, "marshal checkin event")
	}
	if err := svc.archiveFn(event.Time.UnixNano(), msg); err != nil {
		return errors.Wrap(err, "archive checkin")
	}
	if err := svc.publisher.Publish(context.TODO(), topic, msg); err != nil {
		return errors.Wrapf(err, "publish checkin on topic: %s", topic)
	}
	return nil
}
