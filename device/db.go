package device

import (
	"context"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/micromdm/micromdm/depsync"
	"github.com/micromdm/micromdm/mdm/checkin"
	"github.com/micromdm/micromdm/mdm/connect"
	"github.com/micromdm/micromdm/pubsub"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

const (
	DeviceBucket = "mdm.Devices"

	// The deviceIndexBucket index bucket stores serial number and UDID references
	// to the device uuid.
	deviceIndexBucket = "mdm.DeviceIdx"

	DeviceEnrolledTopic = "mdm.DeviceEnrolled"
)

type DB struct {
	*bolt.DB
}

func NewDB(db *bolt.DB, pubsubSvc pubsub.PublishSubscriber) (*DB, error) {
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
	datastore := &DB{
		DB: db,
	}
	if pubsubSvc == nil { // don't start the poller without pubsub.
		return datastore, nil
	}
	if err := datastore.pollCheckin(pubsubSvc); err != nil {
		return nil, err
	}
	return datastore, nil
}

func (db *DB) List() ([]Device, error) {
	// TODO add filter/limit with ForEach
	var devices []Device
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DeviceBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var dev Device
			if err := UnmarshalDevice(v, &dev); err != nil {
				return err
			}
			devices = append(devices, dev)
		}
		return nil
	})
	return devices, err
}

func (db *DB) Save(dev *Device) error {
	tx, err := db.DB.Begin(true)
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}
	bkt := tx.Bucket([]byte(DeviceBucket))
	if bkt == nil {
		return fmt.Errorf("bucket %q not found!", DeviceBucket)
	}
	devproto, err := MarshalDevice(dev)
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

func (db *DB) DeviceByUDID(udid string) (*Device, error) {
	var dev Device
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
		return UnmarshalDevice(v, &dev)
	})
	if err != nil {
		return nil, err
	}
	return &dev, nil
}

func (db *DB) DeviceBySerial(serial string) (*Device, error) {
	var dev Device
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
		return UnmarshalDevice(v, &dev)
	})
	if err != nil {
		return nil, err
	}
	return &dev, nil
}

func isNotFound(err error) bool {
	if _, ok := err.(*notFound); ok {
		return true
	}
	return false
}

func (db *DB) pollCheckin(pubsubSvc pubsub.PublishSubscriber) error {
	authenticateEvents, err := pubsubSvc.Subscribe(context.TODO(), "devices", checkin.AuthenticateTopic)
	if err != nil {
		return errors.Wrapf(err,
			"subscribing devices to %s topic", checkin.AuthenticateTopic)
	}
	tokenUpdateEvents, err := pubsubSvc.Subscribe(context.TODO(), "devices", checkin.TokenUpdateTopic)
	if err != nil {
		return errors.Wrapf(err,
			"subscribing devices to %s topic", checkin.TokenUpdateTopic)
	}
	checkoutEvents, err := pubsubSvc.Subscribe(context.TODO(), "devices", checkin.CheckoutTopic)
	if err != nil {
		return errors.Wrapf(err,
			"subscribing devices to %s topic", checkin.CheckoutTopic)
	}
	depSyncEvents, err := pubsubSvc.Subscribe(context.TODO(), "devices", depsync.SyncTopic)
	if err != nil {
		return errors.Wrapf(err,
			"subscribing devices to %s topic", depsync.SyncTopic)
	}
	connectEvents, err := pubsubSvc.Subscribe(context.TODO(), "devices", connect.ConnectTopic)
	if err != nil {
		return errors.Wrapf(err,
			"subscribing devices to %s topic", connect.ConnectTopic)
	}
	go func() {
		for {
			select {
			case event := <-authenticateEvents:
				var ev checkin.Event
				if err := checkin.UnmarshalEvent(event.Message, &ev); err != nil {
					fmt.Println(err)
					continue
				}
				newDevice := new(Device)
				bySerial, err := db.DeviceBySerial(ev.Command.SerialNumber)
				if err == nil && bySerial != nil { // must be a DEP device
					newDevice = bySerial
				}
				if err != nil && !isNotFound(err) {
					fmt.Println(err) // some other issue is going on
					continue
				}
				_, err = db.DeviceByUDID(ev.Command.UDID)
				if err != nil && isNotFound(err) { // never checked in
					fmt.Printf("checking in new device %s\n", ev.Command.SerialNumber)
				} else if err != nil {
					fmt.Println(err)
					continue
				} else if err == nil {
					fmt.Printf("re-enrolling device %s\n", ev.Command.SerialNumber)
					newDevice.Enrolled = false
				}

				// only create new UUID on initial enrollment.
				if newDevice.UUID == "" {
					newDevice.UUID = uuid.NewV4().String()
				}
				newDevice.UDID = ev.Command.UDID
				newDevice.OSVersion = ev.Command.OSVersion
				newDevice.BuildVersion = ev.Command.BuildVersion
				newDevice.ProductName = ev.Command.ProductName
				newDevice.SerialNumber = ev.Command.SerialNumber
				newDevice.IMEI = ev.Command.IMEI
				newDevice.MEID = ev.Command.MEID
				newDevice.DeviceName = ev.Command.DeviceName
				newDevice.Model = ev.Command.Model
				newDevice.ModelName = ev.Command.ModelName
				newDevice.LastCheckin = time.Now()
				// Challenge:    ev.Command.Challenge, // FIXME: @groob why is this commented out?

				if err := db.Save(newDevice); err != nil {
					fmt.Println(err)
					continue
				}
			case event := <-tokenUpdateEvents:
				var ev checkin.Event
				if err := checkin.UnmarshalEvent(event.Message, &ev); err != nil {
					fmt.Println(err)
					continue
				}
				if ev.Command.UserID != "" {
					continue
				}
				dev, err := db.DeviceByUDID(ev.Command.UDID)
				if err != nil {
					fmt.Println(err)
					continue
				}
				dev.Token = ev.Command.Token.String()
				dev.PushMagic = ev.Command.PushMagic
				dev.UnlockToken = ev.Command.UnlockToken.String()
				dev.AwaitingConfiguration = ev.Command.AwaitingConfiguration
				dev.LastCheckin = time.Now()
				var newlyEnrolled bool = false
				if dev.Enrolled == false {
					newlyEnrolled = true
					dev.Enrolled = true
				}
				if err := db.Save(dev); err != nil {
					fmt.Println(err)
					continue
				}
				if newlyEnrolled {
					fmt.Printf("device %s enrolled\n", ev.Command.UDID)
					err := pubsubSvc.Publish(context.TODO(), DeviceEnrolledTopic, event.Message)
					if err != nil {
						fmt.Println(err)
					}
				}
			case event := <-depSyncEvents:
				var ev depsync.Event
				if err := depsync.UnmarshalEvent(event.Message, &ev); err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Printf("got %d devices from DEP\n", len(ev.Devices))
				for _, d := range ev.Devices {
					newDevice := new(Device)
					bySerial, err := db.DeviceBySerial(d.SerialNumber)
					if err == nil && bySerial != nil { // must be a DEP device
						fmt.Printf("existing device checked in from DEP: %s\n", d.SerialNumber)
						newDevice = bySerial
					}
					if err != nil && !isNotFound(err) {
						fmt.Println(err) // some other issue is going on
						continue
					}
					if newDevice.UUID == "" { // previously unknown
						newDevice.UUID = uuid.NewV4().String()
					}
					newDevice.SerialNumber = d.SerialNumber
					newDevice.Model = d.Model
					newDevice.Description = d.Description
					newDevice.Color = d.Color
					newDevice.AssetTag = d.AssetTag
					newDevice.DEPProfileStatus = DEPProfileStatus(d.ProfileStatus)
					newDevice.DEPProfileUUID = d.ProfileUUID
					newDevice.DEPProfileAssignTime = d.ProfileAssignTime
					newDevice.DEPProfileAssignedDate = d.DeviceAssignedDate
					newDevice.DEPProfileAssignedBy = d.DeviceAssignedBy
					// TODO: deal with sync fields OpType, OpDate
					if err := db.Save(newDevice); err != nil {
						fmt.Println(err)
						continue
					}
				}
			case event := <-connectEvents:
				var ev connect.Event
				if err := connect.UnmarshalEvent(event.Message, &ev); err != nil {
					fmt.Println(err)
					continue
				}
				dev, err := db.DeviceByUDID(ev.Response.UDID)
				if err != nil {
					fmt.Println(err)
					continue
				}
				dev.LastCheckin = time.Now()
				if err := db.Save(dev); err != nil {
					fmt.Println(err)
					continue
				}
			case event := <-checkoutEvents:
				var ev checkin.Event
				if err := checkin.UnmarshalEvent(event.Message, &ev); err != nil {
					fmt.Println(err)
					continue
				}
				dev, err := db.DeviceByUDID(ev.Command.UDID)
				if err != nil {
					fmt.Println(err)
					continue
				}
				dev.Enrolled = false
				dev.LastCheckin = time.Now()
				if err := db.Save(dev); err != nil {
					fmt.Println(err)
					continue
				}
			}
		}
	}()

	return nil
}
