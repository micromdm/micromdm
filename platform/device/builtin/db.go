package builtin

import (
	"context"
	"fmt"
	"encoding/hex"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/platform/device"
)

const (
	DeviceBucket = "mdm.Devices"

	// The deviceIndexBucket index bucket stores serial number and UDID references
	// to the device uuid.
	deviceIndexBucket = "mdm.DeviceIdx"

	// The udidCertAuthBucket stores a simple mapping from UDID to
	// sha256 hash of the device identity certificate for future validation
	udidCertAuthBucket = "mdm.UDIDCertAuth"
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
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(udidCertAuthBucket))
		return err
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating %s bucket", DeviceBucket)
	}
	datastore := &DB{DB: db}
	return datastore, nil
}

func (db *DB) GetBootstrapToken(ctx context.Context, udid string) ([]byte, error) {
	opts := device.ListDevicesOption{}
	
	var udids []string
	udids = append(udids, udid)

	opts.FilterUDID = udids
	
	devices, err := db.List(ctx, opts)
	for _, d := range devices {
		btstring, _ := hex.DecodeString(string(d.BootstrapToken))
		btbytes := []byte(btstring)
		return btbytes, err
	}

	return nil, nil
}

func (db *DB) List(ctx context.Context, opt device.ListDevicesOption) ([]device.Device, error) {
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
			if len(opt.FilterSerial) == 0 && len(opt.FilterUDID) == 0 {
				devices = append(devices, dev)
				return nil
			}
			for _, fs := range opt.FilterSerial {
				if fs == dev.SerialNumber {
					devices = append(devices, dev)
				}
			}
			for _, fs := range opt.FilterUDID {
				if fs == dev.UDID {
					devices = append(devices, dev)
				}
			}			
			return nil
		})
	})
	return devices, err
}

func (db *DB) Save(ctx context.Context, dev *device.Device) error {
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

func (db *DB) DeleteByUDID(ctx context.Context, udid string) error {
	return db.deleteByIndex(udid)
}

func (db *DB) DeleteBySerial(ctx context.Context, serial string) error {
	return db.deleteByIndex(serial)
}

func (db *DB) deleteByIndex(key string) error {
	device, err := db.deviceByIndex(key)
	if err != nil {
		return err
	}

	tx, err := db.DB.Begin(true)
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}

	bkt := tx.Bucket([]byte(DeviceBucket))
	if err := bkt.Delete([]byte(device.UUID)); err != nil {
		return errors.Wrapf(err, "delete device for key %s", key)
	}

	idxBucket := tx.Bucket([]byte(deviceIndexBucket))
	if err := idxBucket.Delete([]byte(device.UDID)); err != nil {
		return errors.Wrapf(err, "delete device index for UDID %s", device.UDID)
	}
	if err := idxBucket.Delete([]byte(device.SerialNumber)); err != nil {
		return errors.Wrapf(err, "delete device index for serial %s", device.SerialNumber)
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

func (db *DB) DeviceByUDID(ctx context.Context, udid string) (*device.Device, error) {
	return db.deviceByIndex(udid)
}

func (db *DB) DeviceBySerial(ctx context.Context, serial string) (*device.Device, error) {
	return db.deviceByIndex(serial)
}

func (db *DB) deviceByIndex(key string) (*device.Device, error) {
	var dev device.Device
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DeviceBucket))
		ib := tx.Bucket([]byte(deviceIndexBucket))
		idx := ib.Get([]byte(key))
		if idx == nil {
			return &notFound{"Device", fmt.Sprintf("key %s", key)}
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

func (db *DB) SaveUDIDCertHash(udid, certHash []byte) error {
	tx, err := db.DB.Begin(true)
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}
	bkt := tx.Bucket([]byte(udidCertAuthBucket))
	if bkt == nil {
		return fmt.Errorf("bucket %q not found!", udidCertAuthBucket)
	}
	if err := bkt.Put(udid, certHash); err != nil {
		return errors.Wrap(err, "put udid cert to boltdb")
	}
	return tx.Commit()
}

func (db *DB) GetUDIDCertHash(udid []byte) ([]byte, error) {
	var certHash []byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(udidCertAuthBucket))
		if b == nil {
			return fmt.Errorf("bucket %q not found!", udidCertAuthBucket)
		}
		certHash = b.Get(udid)
		if certHash == nil {
			return &notFound{"UDID", fmt.Sprintf("udid %s", string(udid))}
		}
		return nil
	})
	return certHash, err
}
