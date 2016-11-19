package checkin

import (
	"bytes"
	"database/sql"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DavidHuie/gomigrate"
	"github.com/go-kit/kit/log"
	"github.com/micromdm/micromdm/command"
	"github.com/micromdm/micromdm/device"
	"github.com/micromdm/micromdm/management"
	"golang.org/x/net/context"
)

var testConn string = "user=postgres password= dbname=travis_ci_test sslmode=disable"

type fixtures struct {
	db       *sql.DB
	server   *httptest.Server
	svc      Service
	devices  device.Datastore
	mgmt     management.Service
	cmd      command.Service
	profile  []byte
	logger   log.Logger
	ctx      context.Context
	testConn string
	migrator *gomigrate.Migrator
	handler  http.Handler
}

// mockMgmtService mocks the management.Service interface which is a dependency of checkin.Service
type mockMgmtService struct{ management.Service }

func (s *mockMgmtService) Push(deviceUDID string) (string, error) {
	return "", nil
}

func setup(t *testing.T) *fixtures {
	var f *fixtures = new(fixtures)
	f.ctx = context.Background()
	l := log.NewLogfmtLogger(os.Stderr)
	f.logger = log.NewContext(l).With("source", "testing")

	db, err := sql.Open("postgres", testConn)
	if err != nil {
		t.Fatalf("opening database connection: %s", err)
	}

	f.migrator, _ = gomigrate.NewMigrator(db, gomigrate.Postgres{}, "../migrations")
	if err = f.migrator.Migrate(); err != nil {
		t.Fatalf("migrating tables: %s", err)
	}

	f.devices, err = device.NewDB("postgres", testConn, f.logger)
	if err != nil {
		t.Fatalf("constructing device datastore: %s", err)
	}

	f.mgmt = &mockMgmtService{}
	f.svc = NewService(f.devices, f.mgmt)
	f.handler = ServiceHandler(f.ctx, f.svc, f.logger)
	f.server = httptest.NewServer(f.handler)

	return f
}

func teardown(f *fixtures, t *testing.T) {
	f.migrator.RollbackAll()

	//f.db.Close()
	f.server.Close()
}

func TestAuthenticate(t *testing.T) {
	f := setup(t)
	defer teardown(f, t)
	requestBody, err := ioutil.ReadFile("../testdata/responses/macos/10.11.x/authenticate.plist")
	if err != nil {
		t.Fatal(err)
	}

	client := http.DefaultClient
	theURL := f.server.URL + "/mdm/checkin"
	req, err := http.NewRequest("PUT", theURL, bytes.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	response, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != 200 {
		var body []byte
		response.Body.Read(body)
		t.Logf("response body: %v", body)
		t.Error(response.Status)
	}

	testDevices, err := f.devices.Devices()
	if err != nil {
		t.Fatal(err)
	}

	if len(testDevices) != 1 {
		t.Errorf("expected 1 device to be inserted, got: %d", len(testDevices))
	}
}

func TestTokenUpdate(t *testing.T) {
	f := setup(t)
	defer teardown(f, t)
	requestBody, err := ioutil.ReadFile("../testdata/responses/macos/10.11.x/token_update.plist")
	if err != nil {
		t.Fatal(err)
	}

	client := http.DefaultClient
	theURL := f.server.URL + "/mdm/checkin"
	req, err := http.NewRequest("PUT", theURL, bytes.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	response, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != 200 {
		var body []byte
		response.Body.Read(body)
		t.Logf("response body: %v", body)
		t.Error(response.Status)
	}

	dev, err := f.devices.GetDeviceByUDID("00000000-1111-2222-3333-444455556666", "mdm_enrolled",
		"apple_mdm_token", "apple_mdm_topic", "apple_push_magic")
	if err != nil {
		t.Fatal(err)
	}

	if dev.Enrolled != true {
		t.Error("expected device to be enrolled")
	}

	if dev.PushMagic != "00000000-1111-2222-3333-444455556666" {
		t.Error("push magic was not updated")
	}

	if dev.MDMTopic != "com.apple.mgmt.test.00000000-1111-2222-3333-444455556666" {
		t.Error("push topic was not updated")
	}
}

func TestCheckout(t *testing.T) {
	f := setup(t)
	defer teardown(f, t)
	requestBody, err := ioutil.ReadFile("../testdata/responses/macos/10.11.x/checkout.plist")
	if err != nil {
		t.Fatal(err)
	}

	client := http.DefaultClient
	theURL := f.server.URL + "/mdm/checkin"
	req, err := http.NewRequest("PUT", theURL, bytes.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	response, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != 200 {
		var body []byte
		response.Body.Read(body)
		t.Logf("response body: %v", body)
		t.Error(response.Status)
	}

	dev, err := f.devices.GetDeviceByUDID("00000000-1111-2222-3333-444455556666", "mdm_enrolled")
	if err != nil {
		t.Fatal(err)
	}

	if dev.Enrolled != false {
		t.Error("expected device to be unenrolled")
	}
}
