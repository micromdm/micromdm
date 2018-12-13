package mysql

import (
	"context"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/kolide/kit/dbutil"
	_ "github.com/go-sql-driver/mysql"
	//"github.com/micromdm/micromdm/platform/dep/sync"
)

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
	
	store, _ := NewDB(db)
	return store
}
