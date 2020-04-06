package pg

import (
	"context"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/kolide/kit/dbutil"
	_ "github.com/lib/pq"
	"github.com/micromdm/micromdm/platform/pubsub/inmem"
	"github.com/micromdm/micromdm/platform/dep/sync"
)

func Test_SaveCursor(t *testing.T) {
	db := setup(t)
	ctx := context.Background()

	// create
	cursor := sync.Cursor{
		Value:     "foobar",
		CreatedAt: time.Now().UTC(),
	}
	
	err := db.SaveCursor(ctx, cursor)

	if err != nil {
		t.Fatal(err)
	}
}

func Test_LoadCursor(t *testing.T) {
	db := setup(t)
	ctx := context.Background()

	cursor, err := db.LoadCursor(ctx)

	if err != nil {
		t.Fatal(err)
	}
	
	if cursor == nil {
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

	store, err := NewDB(db, inmem.NewPubSub())
	if err != nil {
		t.Fatal(err)
	}
	
	return store
}
