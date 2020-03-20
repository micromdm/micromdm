package mysql

import (
	"context"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/kolide/kit/dbutil"
	_ "github.com/go-sql-driver/mysql"
	"github.com/micromdm/micromdm/platform/pubsub/inmem"
	"github.com/micromdm/micromdm/platform/apns"
)

func TestMysqlCrud(t *testing.T) {
	db, err := setup(t)
	ctx := context.Background()
	if err != nil {
		t.Fatal(err)
	}
	
	info := apns.PushInfo{
		UDID:  "UDID-foo-bar-baz",
		Token: "tok",
	}
	err = db.Save(ctx, &info)
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
	return NewDB(db, inmem.NewPubSub())
}
