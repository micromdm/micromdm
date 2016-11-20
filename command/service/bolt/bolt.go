package bolt

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/groob/plist"
	"github.com/micromdm/mdm"
)

type Service struct {
	*bolt.DB
}

func NewCommandService(db *bolt.DB) *Service {
	return &Service{db}
}

const commandBucket = "commands"

func (svc *Service) NewCommand(request *mdm.CommandRequest) (*mdm.Payload, error) {
	payload, err := mdm.NewPayload(request)
	if err != nil {
		return nil, err
	}
	data, err := plist.Marshal(payload)
	if err != nil {
		return nil, err
	}
	// save Key, Payload
	err = svc.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(commandBucket))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found!", commandBucket)
		}
		return bucket.Put([]byte(payload.CommandUUID), data)
	})
	if err != nil {
		return nil, err
	}

	// queue command in device bucket
	err = svc.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(request.UDID))
		if err != nil {
			return fmt.Errorf("create device bucket: %s", err)
		}
		return bucket.Put([]byte(payload.CommandUUID), data)
	})
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func (svc *Service) NextCommand(deviceUDID string) ([]byte, int, error) {
	var (
		next  []byte
		total int
	)
	err := svc.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(deviceUDID))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found!", []byte(deviceUDID))
		}
		_, next = bucket.Cursor().First()
		total = bucket.Stats().KeyN
		return nil
	})
	if err != nil {
		return nil, 0, err
	}
	return next, total, nil
}

func (svc *Service) DeleteCommand(deviceUDID string, commandUUID string) (int, error) {
	var total int
	// delete the command from the command bucket.
	err := svc.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(commandBucket))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found!", commandBucket)
		}
		return bucket.Delete([]byte(commandUUID))
	})
	if err != nil {
		return 0, err
	}

	// delete it from the device queue.
	err = svc.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(deviceUDID))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found!", deviceUDID)
		}
		total = bucket.Stats().KeyN - 1
		return bucket.Delete([]byte(commandUUID))
	})
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (svc *Service) Commands(deviceUDID string) ([]mdm.Payload, error) {
	var payloads []mdm.Payload
	err := svc.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(deviceUDID))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found!", deviceUDID)
		}
		err := bucket.ForEach(func(k, v []byte) error {
			var payload mdm.Payload
			if err := plist.Unmarshal(v, &payload); err != nil {
				return err
			}
			payloads = append(payloads, payload)
			return nil
		})
		return err
	})
	if err != nil {
		return nil, err
	}
	return payloads, nil
}

func (svc *Service) Find(commandUUID string) (*mdm.Payload, error) {
	var payload mdm.Payload
	err := svc.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(commandBucket))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found!", commandBucket)
		}
		data := bucket.Get([]byte(commandUUID))
		if data == nil {
			return fmt.Errorf("command %q not found!", commandUUID)
		}
		if err := plist.Unmarshal(data, &payload); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &payload, nil
}
