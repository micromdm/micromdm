package mysql

import (
	"context"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/kolide/kit/dbutil"
	_ "github.com/go-sql-driver/mysql"
	"github.com/micromdm/micromdm/platform/dep/sync"
)

func Test_SaveCursor(t *testing.T) {
	db, err := setup(t)
	if err != nil {
		t.Fatal(err)
	}
	
	ctx := context.Background()

	// create
	cursor := sync.Cursor{
		Value:     "foobar",
		CreatedAt: time.Now().UTC(),
	}
	
	err = db.SaveCursor(ctx, cursor)

	if err != nil {
		t.Fatal(err)
	}
}

func Test_LoadCursor(t *testing.T) {
	db, err := setup(t)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	cursor, err := db.LoadCursor(ctx)

	if err != nil {
		t.Fatal(err)
	}
	
	if cursor == nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) (*Mysql, error) {
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
