package command

import (
	"context"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/kolide/kit/dbutil"
	_ "github.com/lib/pq"
	"github.com/micromdm/micromdm/workflow/profile/device"
)

func TestCrud(t *testing.T) {
	db := setup(t)
	ctx := context.Background()

	if err := db.UpdateFromList(ctx, "foo", "bar"); err != nil {
		t.Fatal(err)
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
	devdb := device.New(db)
	devdb.Save(context.TODO(), device.Device{
		UDID: "foo",
	})

	return New(db)
}
