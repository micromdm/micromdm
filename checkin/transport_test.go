package checkin

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/DavidHuie/gomigrate"
	"github.com/go-kit/kit/log"
	"github.com/micromdm/micromdm/application"
	"github.com/micromdm/micromdm/certificate"
	"github.com/micromdm/micromdm/command"
	"github.com/micromdm/micromdm/device"
	"github.com/micromdm/micromdm/management"
	"github.com/micromdm/micromdm/workflow"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	db       *sql.DB
	server   *httptest.Server
	svc      Service
	devices  device.Datastore
	mgmt     management.Service
	cmd      command.Service
	profile  []byte = []byte{}
	logger   log.Logger
	ctx      context.Context
	testConn string = "user=postgres password= dbname=travis_ci_test sslmode=disable"
	migrator *gomigrate.Migrator
)

// mockMgmtService mocks the management.Service interface which is a dependency of checkin.Service
type mockMgmtService struct{}

func (s *mockMgmtService) AddProfile(prf *workflow.Profile) (*workflow.Profile, error) {
	return nil, nil
}

func (s *mockMgmtService) Profiles() ([]workflow.Profile, error) {
	return []workflow.Profile{}, nil
}

func (s *mockMgmtService) Profile(uuid string) (*workflow.Profile, error) {
	return nil, nil
}

func (s *mockMgmtService) DeleteProfile(uuid string) error {
	return nil
}

func (s *mockMgmtService) AddWorkflow(wf *workflow.Workflow) (*workflow.Workflow, error) {
	return nil, nil
}

func (s *mockMgmtService) Workflows() ([]workflow.Workflow, error) {
	return nil, nil
}

func (s *mockMgmtService) Devices() ([]device.Device, error) {
	return nil, nil
}

func (s *mockMgmtService) Device(uuid string) (*device.Device, error) {
	return nil, nil
}

func (s *mockMgmtService) InstalledApps(deviceUUID string) ([]application.Application, error) {
	return nil, nil
}

func (s *mockMgmtService) Certificates(deviceUUID string) ([]certificate.Certificate, error) {
	return nil, nil
}

func (s *mockMgmtService) AssignWorkflow(deviceUUID, workflowUUID string) error {
	return nil
}

func (s *mockMgmtService) Push(deviceUDID string) (string, error) {
	return "", nil
}

func (s *mockMgmtService) FetchDEPDevices() error {
	return nil
}

func TestMain(m *testing.M) {
	ctx = context.Background()
	l := log.NewLogfmtLogger(os.Stderr)
	logger = log.NewContext(l).With("source", "testing")

	db, err := sql.Open("postgres", testConn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	migrator, _ = gomigrate.NewMigrator(db, gomigrate.Postgres{}, "../migrations")
	if err = migrator.Migrate(); err != nil {
		panic(fmt.Sprintf("migrating tables: %s", err))
	}

	devices, err = device.NewDB("postgres", testConn, logger)
	if err != nil {
		panic(err)
	}

	mgmt = &mockMgmtService{}

	svc = NewService(devices, mgmt, cmd, profile)
	handler := ServiceHandler(ctx, svc, logger)
	server = httptest.NewServer(handler)
	defer server.Close()

	status := m.Run()
	migrator.RollbackAll()

	os.Exit(status)
}

func TestAuthenticate(t *testing.T) {
	requestBody, err := ioutil.ReadFile("../testdata/responses/macos/10.11.x/authenticate.plist")
	if err != nil {
		t.Fatal(err)
	}

	client := http.DefaultClient
	theURL := server.URL + "/mdm/checkin"
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

	testDevices, err := devices.Devices()
	if err != nil {
		t.Fatal(err)
	}

	if len(testDevices) != 1 {
		t.Errorf("expected 1 device to be inserted, got: %d", len(testDevices))
	}
}

func TestTokenUpdate(t *testing.T) {
	requestBody, err := ioutil.ReadFile("../testdata/responses/macos/10.11.x/token_update.plist")
	if err != nil {
		t.Fatal(err)
	}

	client := http.DefaultClient
	theURL := server.URL + "/mdm/checkin"
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

	dev, err := devices.GetDeviceByUDID("00000000-1111-2222-3333-444455556666", "mdm_enrolled",
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
	requestBody, err := ioutil.ReadFile("../testdata/responses/macos/10.11.x/checkout.plist")
	if err != nil {
		t.Fatal(err)
	}

	client := http.DefaultClient
	theURL := server.URL + "/mdm/checkin"
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

	dev, err := devices.GetDeviceByUDID("00000000-1111-2222-3333-444455556666", "mdm_enrolled")
	if err != nil {
		t.Fatal(err)
	}

	if dev.Enrolled != false {
		t.Error("expected device to be unenrolled")
	}
}
