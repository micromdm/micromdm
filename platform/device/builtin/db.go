package builtin

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/platform/device"
)

const (
	DeviceBucket = "mdm.Devices"

	// The deviceIndexBucket index bucket stores serial number and UDID references
	// to the device uuid.
	deviceIndexBucket = "mdm.DeviceIdx"
)

type DB struct {
	*bolt.DB
}

func NewDB(db *bolt.DB) (*DB, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(deviceIndexBucket))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(DeviceBucket))
		return err
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating %s bucket", DeviceBucket)
	}
	datastore := &DB{DB: db}
	return datastore, nil
}

func (db *DB) List(opt device.ListDevicesOption) ([]device.Device, error) {
	var devices []device.Device
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DeviceBucket))
		// TODO optimize by implemting Seek() and bytes.HasPrefix() so we don't
		// hit all keys in the database if we dont have to.
		return b.ForEach(func(k, v []byte) error {
			var dev device.Device
			if err := device.UnmarshalDevice(v, &dev); err != nil {
				return err
			}
			if len(opt.FilterSerial) == 0 {
				devices = append(devices, dev)
				return nil
			}
			for _, fs := range opt.FilterSerial {
				if fs == dev.SerialNumber {
					devices = append(devices, dev)
				}
			}
			return nil
		})
	})
	return devices, err
}

func (db *DB) Save(dev *device.Device) error {
	tx, err := db.DB.Begin(true)
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}
	bkt := tx.Bucket([]byte(DeviceBucket))
	if bkt == nil {
		return fmt.Errorf("bucket %q not found!", DeviceBucket)
	}
	devproto, err := device.MarshalDevice(dev)
	if err != nil {
		return errors.Wrap(err, "marshalling device")
	}

	// store an array of indices to reference the UUID, which will be the
	// key used to store the actual device.
	indexes := []string{dev.UDID, dev.SerialNumber}
	idxBucket := tx.Bucket([]byte(deviceIndexBucket))
	if idxBucket == nil {
		return fmt.Errorf("bucket %q not found!", deviceIndexBucket)
	}
	for _, idx := range indexes {
		if idx == "" {
			continue
		}
		key := []byte(idx)
		if err := idxBucket.Put(key, []byte(dev.UUID)); err != nil {
			return errors.Wrap(err, "put device to boltdb")
		}
	}

	key := []byte(dev.UUID)
	if err := bkt.Put(key, devproto); err != nil {
		return errors.Wrap(err, "put device to boltdb")
	}
	return tx.Commit()
}

type notFound struct {
	ResourceType string
	Message      string
}

func (e *notFound) Error() string {
	return fmt.Sprintf("not found: %s %s", e.ResourceType, e.Message)
}

func (e *notFound) NotFound() bool {
	return true
}

func (db *DB) DeviceByUDID(udid string) (*device.Device, error) {
	var dev device.Device
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DeviceBucket))
		ib := tx.Bucket([]byte(deviceIndexBucket))
		idx := ib.Get([]byte(udid))
		if idx == nil {
			return &notFound{"Device", fmt.Sprintf("udid %s", udid)}
		}
		v := b.Get(idx)
		if idx == nil {
			return &notFound{"Device", fmt.Sprintf("uuid %s", string(idx))}
		}
		return device.UnmarshalDevice(v, &dev)
	})
	if err != nil {
		return nil, err
	}
	return &dev, nil
}

func (db *DB) DeviceBySerial(serial string) (*device.Device, error) {
	var dev device.Device
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DeviceBucket))
		ib := tx.Bucket([]byte(deviceIndexBucket))
		idx := ib.Get([]byte(serial))
		if idx == nil {
			return &notFound{"Device", fmt.Sprintf("serial %s", serial)}
		}
		v := b.Get(idx)
		if idx == nil {
			return &notFound{"Device", fmt.Sprintf("uuid %s", string(idx))}
		}
		return device.UnmarshalDevice(v, &dev)
	})
	if err != nil {
		return nil, err
	}
	return &dev, nil
}
