package builtin

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"

	firebase "firebase.google.com/go/v4"
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

func (db *FireStoreDB) SaveToFirestore(ctx context.Context, d *device.Device) error {
	data := map[string]interface{}{
		"UUID":         d.UUID,
		"UDID":         d.UDID,
		"DeviceInfo":   d,
		"SerialNumber": d.SerialNumber,
		"CreatedAt":    firestore.ServerTimestamp,
		// Include other fields here.
	}

	_, err := db.FirestoreCli.Collection("devices").Doc(d.UDID).Set(ctx, data)
	if err != nil {
		return fmt.Errorf("adding user to Firestore: %v", err)
	}

	return nil
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

// GetBootstrapToken returns the Bootstrap Token for the device by udid
func (db *DB) GetBootstrapToken(ctx context.Context, udid string) ([]byte, error) {
	d, err := db.DeviceByUDID(ctx, udid)
	if err != nil {
		return nil, errors.Wrap(err, "lookup device by uuid")
	}
	return d.BootstrapToken, nil
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
	// Save to Firestore
	firestoreDB, err := NewFirestoreDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to create Firestore DB: %v", err)
	}

	err = firestoreDB.SaveToFirestore(ctx, dev)
	if err != nil {
		return fmt.Errorf("failed to save device to Firestore: %v", err)
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
