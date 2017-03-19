package device

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/boltdb/bolt"
)

func TestSave(t *testing.T) {
	db := setupDB(t)
	dev := &Device{
		UUID:         "a-b-c-d",
		UDID:         "UDID-FOO-BAR-BAZ",
		SerialNumber: "foobarbaz",
		ProductName:  "MacBook",
	}

	if err := db.Save(dev); err != nil {
		t.Fatalf("saving device in datastore: %s", err)
	}

	byUDID, err := db.DeviceByUDID(dev.UDID)
	if err != nil {
		t.Fatalf("getting device by UDID: %s", err)
	}

	bySerial, err := db.DeviceBySerial(dev.SerialNumber)
	if err != nil {
		t.Fatalf("getting device by UDID: %s", err)
	}

	// test helper that verifies that the retrieved device is the same
	tf := func(haveDev *Device) func(t *testing.T) {
		return func(t *testing.T) {
			if have, want := haveDev.UDID, dev.UDID; have != want {
				t.Errorf("have %s, want %s", have, want)
			}

			if have, want := haveDev.UUID, dev.UUID; have != want {
				t.Errorf("have %s, want %s", have, want)
			}

			if have, want := haveDev.SerialNumber, dev.SerialNumber; have != want {
				t.Errorf("have %s, want %s", have, want)
			}

			if have, want := haveDev.ProductName, dev.ProductName; have != want {
				t.Errorf("have %s, want %s", have, want)
			}
		}
	}

	t.Run("byUDID", tf(byUDID))
	t.Run("bySerial", tf(bySerial))

}

func setupDB(t *testing.T) *DB {
	f, _ := ioutil.TempFile("", "bolt-")
	f.Close()
	os.Remove(f.Name())

	db, err := bolt.Open(f.Name(), 0777, nil)
	if err != nil {
		t.Fatalf("couldn't open bolt, err %s\n", err)
	}
	devDB, err := NewDB(db, nil)
	if err != nil {
		t.Fatalf("couldn't create device DB, err %s\n", err)
	}
	return devDB
}
