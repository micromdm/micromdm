package mysql

import (
	"context"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/kolide/kit/dbutil"
	_ "github.com/go-sql-driver/mysql"
	"github.com/micromdm/micromdm/platform/device"
)

func TestMysqlCrud(t *testing.T) {
	db := setup(t)
	ctx := context.Background()

	// create
	dev := &device.Device{
		UUID:             "foobar",
		UDID:             "foobar",
		DEPProfileStatus: device.ASSIGNED,
		LastSeen:         time.Now().UTC(),
		DEPProfileAssignTime: time.Unix(60*60*24,0).UTC(),
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

func setup(t *testing.T) *Mysql {
	// https://stackoverflow.com/a/23550874/464016
	db, err := dbutil.OpenDBX(
		"mysql",
		"micromdm:micromdm@tcp(127.0.0.1:3306)/micromdm_test?parseTime=true",
		//"host=127.0.0.1 port=3306 user=micromdm dbname=micromdm_test password=micromdm sslmode=disable",
		dbutil.WithLogger(log.NewNopLogger()),
		dbutil.WithMaxAttempts(1),
	)
	
	if err != nil {
		t.Fatal(err)
	}

	return NewDB(db)
}
