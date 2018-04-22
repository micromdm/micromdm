package checkin

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/groob/plist"
	"github.com/micromdm/mdm"
)

func setupDB(t *testing.T) *Checkin {
	f, _ := ioutil.TempFile("", "bolt-")
	f.Close()
	os.Remove(f.Name())

	db, err := bolt.Open(f.Name(), 0777, nil)
	if err != nil {
		t.Fatalf("couldn't open bolt, err %s\n", err)
	}
	svc, err := New(db, nil)
	if err != nil {
		t.Fatalf("couldn't create service, err %s\n", err)
	}
	return svc
}

func mustLoadCommand(t *testing.T, name string) mdm.CheckinCommand {
	var payload mdm.CheckinCommand
	data, err := ioutil.ReadFile("testdata/" + name + ".plist")
	if err != nil {
		t.Fatalf("failed to open test file %q.plist, err: %s", name, err)
	}
	if err := plist.Unmarshal(data, &payload); err != nil {
		t.Fatalf("failed to unmarshal plist %q, err: %s", name, err)
	}
	return payload
}

type mockPublisher struct {
	Invoked   bool
	PublishFn func(string, []byte) error
}

func (m *mockPublisher) Publish(ctx context.Context, s string, b []byte) error {
	m.Invoked = true
	return m.PublishFn(s, b)
}

var passPublisher = func(string, []byte) error { return nil }
var failPublisher = func(string, []byte) error {
	return errors.New("failed")
}

// archiveFunc is the function signature for archiving events in BoltDB.
type archiveFunc func(int64, []byte) error

// override the timestamp with a custom value when saving to BoltDB.
func archiveAt(timestamp int64, svc *Checkin) archiveFunc {
	return func(nano int64, event []byte) error {
		return svc.archive(timestamp, event)
	}
}

func archiveFail() archiveFunc {
	return func(nano int64, event []byte) error {
		return errors.New("archive failed")
	}
}

// load a specific event from the bolt bucket.
func loadEvent(t *testing.T, db *bolt.DB, nano int64) *Event {
	var event Event
	err := db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(CheckinBucket))
		if bkt == nil {
			return fmt.Errorf("no such bucket: CheckinBucket")
		}
		key := []byte(fmt.Sprintf("%d", nano))
		ev := bkt.Get(key)
		if ev == nil {
			return fmt.Errorf("no event at %d timestamp", nano)
		}
		return UnmarshalEvent(ev, &event)
	})
	if err != nil {
		t.Fatalf("error loading event: err = %q", err)
	}
	return &event
}

func TestService_Authenticate(t *testing.T) {
	svc := setupDB(t)
	mock := &mockPublisher{}
	svc.publisher = mock
	tests := []struct {
		name      string
		publisher func(string, []byte) error
		archiveFn archiveFunc
		request   mdm.CheckinCommand
		timestamp int64
		wantErr   bool
	}{
		{
			name:      "happy_path",
			publisher: passPublisher,
			request:   mustLoadCommand(t, "Authenticate"),
			archiveFn: archiveAt(1111, svc),
			timestamp: 1111,
		},
		{
			name:      "archive_fail",
			publisher: passPublisher,
			request:   mustLoadCommand(t, "Authenticate"),
			archiveFn: archiveFail(),
			wantErr:   true,
		},
		{
			name:      "publisher_fail",
			publisher: failPublisher,
			request:   mustLoadCommand(t, "Authenticate"),
			archiveFn: svc.archive,
			wantErr:   true,
		},
		{
			name:      "messageType_fail",
			publisher: passPublisher,
			request:   mustLoadCommand(t, "CheckOut"),
			archiveFn: svc.archive,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.PublishFn = tt.publisher
			mock.Invoked = false
			svc.archiveFn = tt.archiveFn
			err := svc.Authenticate(context.Background(), tt.request, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("%q. Authenticate error = %v, wantErr %v",
					tt.name, err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			event := loadEvent(t, svc.db, tt.timestamp)
			if !reflect.DeepEqual(event.Command, tt.request) {
				t.Errorf("\nwant: %#v\n,\nhave: %#v\n", tt.request, event.Command)
			}

			if !mock.Invoked {
				t.Errorf("publisher not invoked")
			}
		})
	}
}

func TestService_TokenUpdate(t *testing.T) {
	svc := setupDB(t)
	mock := &mockPublisher{}
	svc.publisher = mock
	tests := []struct {
		name      string
		publisher func(string, []byte) error
		archiveFn archiveFunc
		request   mdm.CheckinCommand
		timestamp int64
		wantErr   bool
	}{
		{
			name:      "happy_path",
			publisher: passPublisher,
			request:   mustLoadCommand(t, "TokenUpdate"),
			archiveFn: archiveAt(2222, svc),
			timestamp: 2222,
		},
		{
			name:      "archive_fail",
			publisher: passPublisher,
			request:   mustLoadCommand(t, "TokenUpdate"),
			archiveFn: archiveFail(),
			wantErr:   true,
		},
		{
			name:      "publisher_fail",
			publisher: failPublisher,
			request:   mustLoadCommand(t, "TokenUpdate"),
			archiveFn: svc.archive,
			wantErr:   true,
		},
		{
			name:      "messageType_fail",
			publisher: passPublisher,
			request:   mustLoadCommand(t, "CheckOut"),
			archiveFn: svc.archive,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.PublishFn = tt.publisher
			mock.Invoked = false
			svc.archiveFn = tt.archiveFn
			err := svc.TokenUpdate(context.Background(), tt.request, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("%q. TokenUpdate error = %v, wantErr %v",
					tt.name, err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			event := loadEvent(t, svc.db, tt.timestamp)
			if !reflect.DeepEqual(event.Command, tt.request) {
				t.Errorf("\nwant: %#v\n,\nhave: %#v\n", tt.request, event.Command)
			}

			if !mock.Invoked {
				t.Errorf("publisher not invoked")
			}
		})
	}
}

func TestService_CheckOut(t *testing.T) {
	svc := setupDB(t)
	mock := &mockPublisher{}
	svc.publisher = mock
	tests := []struct {
		name      string
		publisher func(string, []byte) error
		archiveFn archiveFunc
		request   mdm.CheckinCommand
		timestamp int64
		wantErr   bool
	}{
		{
			name:      "happy_path",
			publisher: passPublisher,
			request:   mustLoadCommand(t, "CheckOut"),
			archiveFn: archiveAt(1111, svc),
			timestamp: 1111,
		},
		{
			name:      "archive_fail",
			publisher: passPublisher,
			request:   mustLoadCommand(t, "CheckOut"),
			archiveFn: archiveFail(),
			wantErr:   true,
		},
		{
			name:      "publisher_fail",
			publisher: failPublisher,
			request:   mustLoadCommand(t, "CheckOut"),
			archiveFn: svc.archive,
			wantErr:   true,
		},
		{
			name:      "messageType_fail",
			publisher: passPublisher,
			request:   mustLoadCommand(t, "Authenticate"),
			archiveFn: svc.archive,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.PublishFn = tt.publisher
			mock.Invoked = false
			svc.archiveFn = tt.archiveFn
			err := svc.CheckOut(context.Background(), tt.request, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("%q. CheckOut error = %v, wantErr %v",
					tt.name, err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			event := loadEvent(t, svc.db, tt.timestamp)
			if !reflect.DeepEqual(event.Command, tt.request) {
				t.Errorf("\nwant: %#v\n,\nhave: %#v\n", tt.request, event.Command)
			}

			if !mock.Invoked {
				t.Errorf("publisher not invoked")
			}
		})
	}
}
