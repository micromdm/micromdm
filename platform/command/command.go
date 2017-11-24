package command

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/micromdm/mdm"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/micromdm/micromdm/platform/pubsub"
)

const (

	// CommandBucket is the *bolt.DB bucket where commands are archived.
	CommandBucket = "mdm.Command.ARCHIVE"

	// CommandTopic is a PubSub topic that events are published to.
	CommandTopic = "mdm.Command"
)

type Command struct {
	db        *bolt.DB
	publisher pubsub.Publisher
	archiveFn func(int64, []byte) error
}

func New(db *bolt.DB, pub pubsub.Publisher) (*Command, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(CommandBucket))
		return err
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating %s bucket", CommandBucket)
	}
	svc := Command{
		db:        db,
		publisher: pub,
	}
	svc.archiveFn = svc.archive
	return &svc, nil
}

func (svc *Command) NewCommand(ctx context.Context, request *mdm.CommandRequest) (*mdm.Payload, error) {
	if request == nil {
		return nil, errors.New("empty CommandRequest")
	}
	payload, err := mdm.NewPayload(request)
	if err != nil {
		return nil, errors.Wrap(err, "creating mdm payload")
	}
	event := NewEvent(*payload, request.UDID)
	msg, err := MarshalEvent(event)
	if err != nil {
		return nil, errors.Wrap(err, "marshalling mdm command event")
	}
	if err := svc.archive(event.Time.UnixNano(), msg); err != nil {
		return nil, errors.Wrap(err, "archive mdm command")
	}
	if err := svc.publisher.Publish(context.TODO(), CommandTopic, msg); err != nil {
		return nil, errors.Wrapf(err, "publish mdm command on topic: %s", CommandTopic)
	}
	return payload, nil
}

// archive events to BoltDB bucket using timestamp as key to preserve order.
func (svc *Command) archive(nano int64, msg []byte) error {
	tx, err := svc.db.Begin(true)
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}
	defer tx.Rollback()

	bkt := tx.Bucket([]byte(CommandBucket))
	if bkt == nil {
		return fmt.Errorf("bucket %q not found!", CommandBucket)
	}
	key := []byte(fmt.Sprintf("%d", nano))
	if err := bkt.Put(key, msg); err != nil {
		return errors.Wrap(err, "put command event to boltdb")
	}
	return tx.Commit()
}
