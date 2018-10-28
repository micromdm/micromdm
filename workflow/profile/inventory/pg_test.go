package inventory

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-kit/kit/log"
	"github.com/groob/plist"
	"github.com/kolide/kit/dbutil"

	"github.com/micromdm/micromdm/pkg/id"
	"github.com/micromdm/micromdm/workflow/profile/device"

	_ "github.com/lib/pq"
)

func Test_Decode(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/response.plist")
	if err != nil {
		t.Fatal(err)
	}
	var resp ListProfilesResponse
	if err := plist.Unmarshal(data, &resp); err != nil {
		t.Fatal(err)
	}
	spew.Dump(resp)
}

func TestCrud(t *testing.T) {
	db := setup(t)
	ctx := context.Background()

	var payloads []Payload
	for i := 0; i < 5; i++ {
		payload := Payload{
			PayloadUUID:       id.New(),
			PayloadIdentifier: id.New(),
		}
		payloads = append(payloads, payload)
	}
	err := db.UpdateFromListResponse(ctx, "foo", ListProfilesResponse{
		ProfileList: payloads,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = db.UpdateFromListResponse(ctx, "foo", ListProfilesResponse{
		ProfileList: payloads,
	})
	if err != nil {
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
