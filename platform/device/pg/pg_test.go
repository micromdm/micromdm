//go:build pg
// +build pg

package pg

import (
	"context"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/kolide/kit/dbutil"
	_ "github.com/lib/pq"
	"github.com/micromdm/micromdm/platform/device"
)

func TestPGCrud(t *testing.T) {
	db := setup(t)
	ctx := context.Background()

	// create
	dev := &device.Device{
		UUID:             "foobar",
		UDID:             "foobar",
		DEPProfileStatus: device.ASSIGNED,
		LastSeen:         time.Now().UTC(),
	}
	err := db.Save(ctx, dev)
	if err != nil {
		t.Fatal(err)
	}

	// update
	dev.DEPProfileStatus = device.PUSHED
	err = db.Save(ctx, dev)
	if err != nil {
		t.Fatal(err)
	}

	// find
	found, err := db.DeviceByUDID(ctx, dev.UDID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := found.DEPProfileStatus, dev.DEPProfileStatus; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	// list
	devices, err := db.ListDevices(ctx, device.ListDevicesOption{})
	if err != nil {
		t.Fatal(err)
	}

	// delete
	for _, dev := range devices {
		if err := db.DeleteByUDID(ctx, dev.UDID); err != nil {
			t.Fatal(err)
		}
	}
}

func setup(t *testing.T) *Postgres {
	db, err := dbutil.OpenDBX(
		"postgres",
		"host=localhost port=5432 user=micromdm dbname=micromdm_test password=micromdm sslmode=disable",
		dbutil.WithLogger(log.NewNopLogger()),
		dbutil.WithMaxAttempts(1),
	)
	if err != nil {
		t.Fatal(err)
	}

	return New(db)
}
