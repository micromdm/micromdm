// +build pg

package pg

import (
	"context"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/kolide/kit/dbutil"
	_ "github.com/lib/pq"
	"github.com/micromdm/micromdm/platform/apns"
)

func TestPGCrud(t *testing.T) {
	db := setup(t)
	ctx := context.Background()

	info := apns.PushInfo{
		UDID:  "UDID-foo-bar-baz",
		Token: "tok",
	}
	err := db.Save(ctx, &info)
	if err != nil {
		t.Fatal(err)
	}

	found, err := db.PushInfo(ctx, info.UDID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := found.Token, info.Token; have != want {
		t.Errorf("have %s, want %s", have, want)
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
