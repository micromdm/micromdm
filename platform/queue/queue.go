// Package queue implements a boldDB backed queue for MDM Commands.
package queue

import (
	"context"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/groob/plist"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/mdm"
	"github.com/micromdm/micromdm/platform/command"
	"github.com/micromdm/micromdm/platform/pubsub"
)

const (
	DeviceCommandBucket = "mdm.DeviceCommands"

	CommandQueuedTopic = "mdm.CommandQueued"
)

type Store struct {
	*bolt.DB
}

func (db *Store) Next(ctx context.Context, resp mdm.Response) ([]byte, error) {
	cmd, err := db.nextCommand(ctx, resp)
	if err != nil {
		return nil, err
	}
	if cmd == nil {
		return nil, nil
	}
	return cmd.Payload, nil
}

func (db *Store) nextCommand(ctx context.Context, resp mdm.Response) (*Command, error) {
	udid := resp.UDID
	if resp.UserID != nil {
		// use the user id for user level commands
		udid = *resp.UserID
	}
	dc, err := db.DeviceCommand(udid)
	if err != nil {
		if isNotFound(err) {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "get device command from queue, udid: %s", resp.UDID)
	}

	var cmd *Command
	switch resp.Status {
	case "NotNow":
		// We will try this command later when the device is not
		// responding with NotNow
		x, a := cut(dc.Commands, resp.CommandUUID)
		dc.Commands = a
		if x == nil {
			break
		}
		dc.NotNow = append(dc.NotNow, *x)

	case "Acknowledged":
		// move to completed, send next
		x, a := cut(dc.Commands, resp.CommandUUID)
		dc.Commands = a
		if x == nil {
			break
		}
		dc.Completed = append(dc.Completed, *x)
	case "Error":
		// move to failed, send next
		x, a := cut(dc.Commands, resp.CommandUUID)
		dc.Commands = a
		if x == nil { // must've already bin ackd
			break
		}
		dc.Failed = append(dc.Failed, *x)

	case "CommandFormatError":
		// move to failed
		x, a := cut(dc.Commands, resp.CommandUUID)
		dc.Commands = a
		if x == nil {
			break
		}
		dc.Failed = append(dc.Failed, *x)

	case "Idle":
		// will send next command below

	default:
		return nil, fmt.Errorf("unknown response status: %s", resp.Status)
	}

	// pop the first command from the queue and add it to the end.
	// If the regular queue is empty, send a command that got
	// refused with NotNow before.
	cmd, dc.Commands = popFirst(dc.Commands)
	if cmd != nil {
		dc.Commands = append(dc.Commands, *cmd)
	} else if resp.Status != "NotNow" {
		cmd, dc.NotNow = popFirst(dc.NotNow)
		if cmd != nil {
			dc.Commands = append(dc.Commands, *cmd)
		}
	}

	if err := db.Save(dc); err != nil {
		return nil, err
	}

	return cmd, nil
}

func popFirst(all []Command) (*Command, []Command) {
	if len(all) == 0 {
		return nil, all
	}
	first := all[0]
	all = append(all[:0], all[1:]...)
	return &first, all
}

func cut(all []Command, uuid string) (*Command, []Command) {
	for i, cmd := range all {
		if cmd.UUID == uuid {
			all = append(all[:i], all[i+1:]...)
			return &cmd, all
		}
	}
	return nil, all
}

func NewQueue(db *bolt.DB, pubsub pubsub.PublishSubscriber) (*Store, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(DeviceCommandBucket))
		return err
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating %s bucket", DeviceCommandBucket)
	}
	datastore := &Store{DB: db}
	if err := datastore.pollCommands(pubsub); err != nil {
		return nil, err
	}
	return datastore, nil
}

func (db *Store) Save(cmd *DeviceCommand) error {
	tx, err := db.DB.Begin(true)
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}
	bkt := tx.Bucket([]byte(DeviceCommandBucket))
	if bkt == nil {
		return fmt.Errorf("bucket %q not found!", DeviceCommandBucket)
	}
	devproto, err := MarshalDeviceCommand(cmd)
	if err != nil {
		return errors.Wrap(err, "marshalling DeviceCommand")
	}
	key := []byte(cmd.DeviceUDID)
	if err := bkt.Put(key, devproto); err != nil {
		return errors.Wrap(err, "put DeviceCommand to boltdb")
	}
	return tx.Commit()
}

func (db *Store) DeviceCommand(udid string) (*DeviceCommand, error) {
	var dev DeviceCommand
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DeviceCommandBucket))
		v := b.Get([]byte(udid))
		if v == nil {
			return &notFound{"DeviceCommand", fmt.Sprintf("udid %s", udid)}
		}
		return UnmarshalDeviceCommand(v, &dev)
	})
	if err != nil {
		return nil, err
	}
	return &dev, nil
}

type notFound struct {
	ResourceType string
	Message      string
}

func (e *notFound) Error() string {
	return fmt.Sprintf("not found: %s %s", e.ResourceType, e.Message)
}

func (db *Store) pollCommands(pubsub pubsub.PublishSubscriber) error {
	commandEvents, err := pubsub.Subscribe(context.TODO(), "command-queue", command.CommandTopic)
	if err != nil {
		return errors.Wrapf(err,
			"subscribing push to %s topic", command.CommandTopic)
	}
	go func() {
		for {
			select {
			case event := <-commandEvents:
				var ev command.Event
				if err := command.UnmarshalEvent(event.Message, &ev); err != nil {
					fmt.Println(err)
					continue
				}

				cmd := new(DeviceCommand)
				cmd.DeviceUDID = ev.DeviceUDID
				byUDID, err := db.DeviceCommand(ev.DeviceUDID)
				if err == nil && byUDID != nil {
					cmd = byUDID
				}
				newPayload, err := plist.Marshal(ev.Payload)
				if err != nil {
					fmt.Println(errors.Wrap(err, "marshal event payload"))
					continue
				}
				newCmd := Command{
					UUID:    ev.Payload.CommandUUID,
					Payload: newPayload,
				}
				cmd.Commands = append(cmd.Commands, newCmd)
				if err := db.Save(cmd); err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Printf("queued event for device: %s\n", ev.DeviceUDID)

				cq := new(QueueCommandQueued)
				cq.DeviceUDID = ev.DeviceUDID
				cq.CommandUUID = ev.Payload.CommandUUID

				msgBytes, err := MarshalQueuedCommand(cq)
				if err != nil {
					fmt.Println(err)
					continue
				}

				err = pubsub.Publish(context.TODO(), CommandQueuedTopic, msgBytes)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}()

	return nil
}

func isNotFound(err error) bool {
	if _, ok := err.(*notFound); ok {
		return true
	}
	return false
}
