package command

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/micromdm/mdm"
	"golang.org/x/net/context"
)

func TestService_NewCommand(t *testing.T) {
	svc := setupDB(t)
	mock := &mockPublisher{}
	svc.publisher = mock
	passPublisher := func(string, []byte) error { return nil }
	failPublisher := func(string, []byte) error {
		return errors.New("failed")
	}

	tests := []struct {
		name      string
		publisher func(string, []byte) error
		request   *mdm.CommandRequest
		wantErr   bool
	}{
		{
			name:      "happy path",
			wantErr:   false,
			publisher: passPublisher,
			request: &mdm.CommandRequest{
				UDID: "foobarbaz",
				Command: mdm.Command{
					RequestType: "DeviceInformation",
					DeviceInformation: mdm.DeviceInformation{
						Queries: []string{"foo", "bar", "baz"},
					},
				},
			},
		},
		{
			name:      "publish fail",
			wantErr:   true,
			publisher: failPublisher,
			request: &mdm.CommandRequest{
				Command: mdm.Command{
					RequestType: "DeviceInformation",
				},
			},
		},
		{
			name:      "empty request",
			wantErr:   true,
			publisher: passPublisher,
		},
		{
			name:      "bad payload",
			wantErr:   true,
			publisher: passPublisher,
			request: &mdm.CommandRequest{
				UDID: "foobarbaz",
				Command: mdm.Command{
					RequestType: "DevicePropaganda",
				},
			},
		},
	}
	for _, tt := range tests {
		mock.PublishFn = tt.publisher
		_, err := svc.NewCommand(context.Background(), tt.request)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. CommandService.NewCommand() error = %v, wantErr %v",
				tt.name, err, tt.wantErr)
			continue
		}
	}
}

type mockPublisher struct {
	PublishFn func(string, []byte) error
}

func (m *mockPublisher) Publish(s string, b []byte) error {
	return m.PublishFn(s, b)
}

func setupDB(t *testing.T) *Command {
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
