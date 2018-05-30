package builtin

import (
	"context"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/micromdm/micromdm/dep/depsync"
	"github.com/micromdm/micromdm/mdm"
	"github.com/micromdm/micromdm/platform/device"
	"github.com/micromdm/micromdm/platform/pubsub"
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

func IsNotFound(err error) bool {
	if _, ok := err.(*notFound); ok {
		return true
	}
	return false
}

func (db *DB) pollCheckin(pubsubSvc pubsub.PublishSubscriber) error {
	authenticateEvents, err := pubsubSvc.Subscribe(context.TODO(), "devices", mdm.AuthenticateTopic)
	if err != nil {
		return errors.Wrapf(err,
			"subscribing devices to %s topic", mdm.AuthenticateTopic)
	}
	tokenUpdateEvents, err := pubsubSvc.Subscribe(context.TODO(), "devices", mdm.TokenUpdateTopic)
	if err != nil {
		return errors.Wrapf(err,
			"subscribing devices to %s topic", mdm.TokenUpdateTopic)
	}
	checkoutEvents, err := pubsubSvc.Subscribe(context.TODO(), "devices", mdm.CheckoutTopic)
	if err != nil {
		return errors.Wrapf(err,
			"subscribing devices to %s topic", mdm.CheckoutTopic)
	}
	depSyncEvents, err := pubsubSvc.Subscribe(context.TODO(), "devices", depsync.SyncTopic)
	if err != nil {
		return errors.Wrapf(err,
			"subscribing devices to %s topic", depsync.SyncTopic)
	}
	connectEvents, err := pubsubSvc.Subscribe(context.TODO(), "devices", mdm.ConnectTopic)
	if err != nil {
		return errors.Wrapf(err,
			"subscribing devices to %s topic", mdm.ConnectTopic)
	}
	go func() {
		for {
			select {
			case event := <-authenticateEvents:
				var ev mdm.CheckinEvent
				if err := mdm.UnmarshalCheckinEvent(event.Message, &ev); err != nil {
					fmt.Println(err)
					continue
				}
				newDevice := new(device.Device)
				bySerial, err := db.DeviceBySerial(ev.Command.SerialNumber)
				if err == nil && bySerial != nil { // must be a DEP device
					newDevice = bySerial
				}
				if err != nil && !IsNotFound(err) {
					fmt.Println(err) // some other issue is going on
					continue
				}
				byUDID, err := db.DeviceByUDID(ev.Command.UDID)
				if err != nil && IsNotFound(err) { // never checked in
					fmt.Printf("checking in new device %s\n", ev.Command.SerialNumber)
				} else if err != nil {
					fmt.Println(err)
					continue
				} else if err == nil {
					fmt.Printf("re-enrolling device %s\n", ev.Command.SerialNumber)
					newDevice = byUDID
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
				newDevice.LastSeen = time.Now()
				// the Authenticate message is the only place where we
				// should unconditionally set the cert in the DB
				// technically this is duplicated because the middleware
				// will save a cert to a device that has a nil cert
				// but at some point we'll turn that off
				newDevice.DeviceCert = ev.DeviceCert
				// Challenge:    ev.Command.Challenge, // FIXME: @groob why is this commented out?

				if err := db.Save(newDevice); err != nil {
					fmt.Println(err)
					continue
				}
			case event := <-tokenUpdateEvents:
				var ev mdm.CheckinEvent
				if err := mdm.UnmarshalCheckinEvent(event.Message, &ev); err != nil {
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
				dev.LastSeen = time.Now()
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
					err := pubsubSvc.Publish(context.TODO(), device.DeviceEnrolledTopic, event.Message)
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
					updDevice, updDeviceErr := db.DeviceBySerial(d.SerialNumber)
					if updDeviceErr != nil && !IsNotFound(updDeviceErr) {
						fmt.Printf("error getting device %s: %s\n", d.SerialNumber, err)
						continue
					}

					if updDeviceErr != nil && IsNotFound(updDeviceErr) {
						updDevice = new(device.Device)
						if d.OpType == "modified" {
							fmt.Printf("warning: no existing device for DEP device update: %s\n", d.SerialNumber)
						}
					}
					if updDeviceErr == nil && d.OpType == "added" {
						// consider issuing this warning if op_type == "" as well.
						// in that case it's likely the device came from a DEP
						// fetch (vs. a sync) which could be a re-fetch of devices
						fmt.Printf("warning: re-adding existing DEP device: %s\n", d.SerialNumber)
					}
					if d.OpType == "deleted" {
						fmt.Printf("warning: DEP device unassigned: %s\n", d.SerialNumber)
					}

					if updDevice.UUID == "" {
						// generate UUID for any device that doesn't have one
						updDevice.UUID = uuid.NewV4().String()
					}

					updDevice.SerialNumber = d.SerialNumber
					updDevice.Model = d.Model
					updDevice.Description = d.Description
					updDevice.Color = d.Color
					updDevice.AssetTag = d.AssetTag
					updDevice.DEPProfileStatus = device.DEPProfileStatus(d.ProfileStatus)
					updDevice.DEPProfileUUID = d.ProfileUUID
					updDevice.DEPProfileAssignTime = d.ProfileAssignTime
					updDevice.DEPProfileAssignedDate = d.DeviceAssignedDate
					updDevice.DEPProfileAssignedBy = d.DeviceAssignedBy
					// TODO: support profile_push_time, os, device_family, op_date

					if err := db.Save(updDevice); err != nil {
						fmt.Println(err)
						continue
					}
				}
			case event := <-connectEvents:
				var ev mdm.AcknowledgeEvent
				if err := mdm.UnmarshalAcknowledgeEvent(event.Message, &ev); err != nil {
					fmt.Println(err)
					continue
				}
				dev, err := db.DeviceByUDID(ev.Response.UDID)
				if err != nil {
					fmt.Println(err)
					continue
				}
				dev.LastSeen = time.Now()
				if err := db.Save(dev); err != nil {
					fmt.Println(err)
					continue
				}
			case event := <-checkoutEvents:
				var ev mdm.CheckinEvent
				if err := mdm.UnmarshalCheckinEvent(event.Message, &ev); err != nil {
					fmt.Println(err)
					continue
				}
				dev, err := db.DeviceByUDID(ev.Command.UDID)
				if err != nil {
					fmt.Println(err)
					continue
				}
				dev.Enrolled = false
				dev.LastSeen = time.Now()
				if err := db.Save(dev); err != nil {
					fmt.Println(err)
					continue
				}
			}
		}
	}()

	return nil
}
