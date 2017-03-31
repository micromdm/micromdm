package queue

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/micromdm/mdm"
)

func TestNext_Error(t *testing.T) {
	store, teardown := setupDB(t)
	defer teardown()

	dc := &DeviceCommand{DeviceUDID: "TestDevice"}
	dc.Commands = append(dc.Commands, Command{UUID: "xCmd"})
	dc.Commands = append(dc.Commands, Command{UUID: "yCmd"})
	dc.Commands = append(dc.Commands, Command{UUID: "zCmd"})
	if err := store.Save(dc); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	resp := mdm.Response{
		UDID:        dc.DeviceUDID,
		CommandUUID: "xCmd",
		Status:      "Error",
	}
	for range dc.Commands {
		cmd, err := store.Next(ctx, resp)
		if err != nil {
			t.Fatalf("expected nil, but got err: %s", err)
		}
		if cmd == nil {
			t.Fatal("expected cmd but got nil")
		}

		if have, errd := cmd.UUID, resp.CommandUUID; have == errd {
			t.Error("got back command which previously failed")
		}
	}
}

func TestNext_NotNow(t *testing.T) {
	store, teardown := setupDB(t)
	defer teardown()

	dc := &DeviceCommand{DeviceUDID: "TestDevice"}
	dc.Commands = append(dc.Commands, Command{UUID: "xCmd"})
	dc.Commands = append(dc.Commands, Command{UUID: "yCmd"})
	dc.Commands = append(dc.Commands, Command{UUID: "zCmd"})
	if err := store.Save(dc); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	tf := func(t *testing.T) {
		resp := mdm.Response{
			UDID:        dc.DeviceUDID,
			CommandUUID: "xCmd",
			Status:      "NotNow",
		}
		for i := 0; i <= 10; i++ {
			cmd, err := store.Next(ctx, resp)
			if err != nil {
				t.Fatalf("expected nil, but got err: %s", err)
			}
			if cmd == nil {
				continue
			}
			if have, notNow := cmd.UUID, "xCmd"; have == notNow {
				t.Error("got back a notnowed command.")
			}
		}
	}

	t.Run("withManyCommands", tf)
	dc.Commands = []Command{{UUID: "xCmd"}}
	if err := store.Save(dc); err != nil {
		t.Fatal(err)
	}
	t.Run("withOneCommand", tf)
}

func TestNext_Idle(t *testing.T) {
	store, teardown := setupDB(t)
	defer teardown()

	dc := &DeviceCommand{DeviceUDID: "TestDevice"}
	dc.Commands = append(dc.Commands, Command{UUID: "xCmd"})
	dc.Commands = append(dc.Commands, Command{UUID: "yCmd"})
	dc.Commands = append(dc.Commands, Command{UUID: "zCmd"})
	if err := store.Save(dc); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	resp := mdm.Response{
		UDID:        dc.DeviceUDID,
		CommandUUID: "xCmd",
		Status:      "Idle",
	}
	for i, _ := range dc.Commands {
		cmd, err := store.Next(ctx, resp)
		if err != nil {
			t.Fatalf("expected nil, but got err: %s", err)
		}
		if cmd == nil {
			t.Fatal("expected cmd but got nil")
		}

		if have, want := cmd.UUID, dc.Commands[i].UUID; have != want {
			t.Errorf("have %s, want %s, index %d", have, want, i)
		}
	}
}

func TestNext_zeroCommands(t *testing.T) {
	store, teardown := setupDB(t)
	defer teardown()

	dc := &DeviceCommand{DeviceUDID: "TestDevice"}
	if err := store.Save(dc); err != nil {
		t.Fatal(err)
	}

	var allStatuses = []string{
		"Acknowledged",
		"NotNow",
	}

	ctx := context.Background()
	for _, s := range allStatuses {
		t.Run(s, func(t *testing.T) {
			resp := mdm.Response{CommandUUID: s, Status: s}
			cmd, err := store.Next(ctx, resp)
			if err != nil {
				t.Errorf("expected nil, but got err: %s", err)
			}
			if cmd != nil {
				t.Errorf("expected nil cmd but got %s", cmd.UUID)
			}
		})
	}

}

func setupDB(t *testing.T) (*Store, func()) {
	f, _ := ioutil.TempFile("", "bolt-")
	teardown := func() {
		f.Close()
		os.Remove(f.Name())
	}

	db, err := bolt.Open(f.Name(), 0777, nil)
	if err != nil {
		t.Fatalf("couldn't open bolt, err %s\n", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(DeviceCommandBucket))
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	store := &Store{db}
	return store, teardown
}
