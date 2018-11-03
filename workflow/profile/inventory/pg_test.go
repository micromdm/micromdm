package inventory

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/groob/plist"
	"github.com/kolide/kit/dbutil"

	"github.com/micromdm/micromdm/pkg/id"
	"github.com/micromdm/micromdm/workflow/profile/device"

	_ "github.com/lib/pq"
)

const testDeviceIdentifier = "foo"

func Test_Decode(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/response.plist")
	if err != nil {
		t.Fatal(err)
	}
	var resp ListProfilesResponse
	if err := plist.Unmarshal(data, &resp); err != nil {
		t.Fatal(err)
	}

	if have, want := len(resp.ProfileList), 1; have != want {
		t.Errorf("expected %d, got %d profiles in response", want, have)
	}
}

func TestUpdateFromList(t *testing.T) {
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

	updatePayloads(t, db, payloads)

	// remove one and do it again
	removeID := payloads[1].PayloadUUID
	payloads = append(payloads[:1], payloads[1+1:]...)
	updatePayloads(t, db, payloads)

	// check that it's gone
	_, err := db.DeviceProfileByUUID(ctx, testDeviceIdentifier, removeID)
	if have, want := isNotFound(err), true; have != want {
		t.Errorf("expected a not found error, got %v", err)
	}
}

func updatePayloads(t *testing.T, db *Postgres, payloads []Payload) {
	t.Helper()
	ctx := context.Background()
	if err := db.UpdateFromListResponse(
		ctx,
		testDeviceIdentifier,
		ListProfilesResponse{
			ProfileList: payloads,
		},
	); err != nil {
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
		UDID: testDeviceIdentifier,
	})

	return New(db)
}
