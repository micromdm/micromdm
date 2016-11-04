package connect

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/DavidHuie/gomigrate"
	"github.com/go-kit/kit/log"
	"github.com/micromdm/mdm"
	"github.com/micromdm/micromdm/application"
	"github.com/micromdm/micromdm/certificate"
	"github.com/micromdm/micromdm/command"
	"github.com/micromdm/micromdm/device"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// mockCommandService implements the command.Service interface and returns only the single command struct given to it
// instead of querying a live redis instance.
type mockCommandService struct {
	t           *testing.T
	MockCommand *mdm.CommandRequest
}

func (svc mockCommandService) NewCommand(cmd *mdm.CommandRequest) (*mdm.Payload, error) {
	svc.MockCommand = cmd
	return nil, nil
}

func (svc mockCommandService) NextCommand(udid string) ([]byte, int, error) {
	return nil, 0, nil
}

func (svc mockCommandService) DeleteCommand(deviceUDID, commandUUID string) (int, error) {
	svc.t.Logf("Deleting mock command with UUID %s", commandUUID)
	return 1, nil
}

func (svc mockCommandService) Commands(deviceUDID string) ([]mdm.Payload, error) {
	return []mdm.Payload{}, nil
}

func (svc mockCommandService) Find(commandUUID string) (*mdm.Payload, error) {
	svc.t.Logf("Returning mock response finding command with UUID %s", commandUUID)

	payload, _ := mdm.NewPayload(svc.MockCommand)
	payload.CommandUUID = commandUUID

	return payload, nil
}

type connectFixtures struct {
	db         *sql.DB
	server     *httptest.Server
	svc        Service
	devices    device.Datastore
	apps       application.Datastore
	certs      certificate.Datastore
	cs         command.Service
	logger     log.Logger
	deviceUUID string
	migrator   *gomigrate.Migrator
}

func setup(t *testing.T, cmd *mdm.CommandRequest) *connectFixtures {
	ctx := context.Background()
	l := log.NewLogfmtLogger(os.Stderr)
	logger := log.NewContext(l).With("source", "testing")

	var (
		err      error
		testConn string = "user=postgres password= dbname=travis_ci_test sslmode=disable"
		devices  device.Datastore
		apps     application.Datastore
		certs    certificate.Datastore
		cs       command.Service
	)

	db, err := sql.Open("postgres", testConn)
	if err != nil {
		t.Fatal(err)
	}
	migrator, _ := gomigrate.NewMigrator(db, gomigrate.Postgres{}, "../migrations")
	if err = migrator.Migrate(); err != nil {
		t.Fatalf("migrating tables: %s", err)
	}

	devices, err = device.NewDB("postgres", testConn, logger)
	if err != nil {
		t.Fatal(err)
	}

	apps, err = application.NewDB("postgres", testConn, logger)
	if err != nil {
		t.Fatal(err)
	}

	certs, err = certificate.NewDB("postgres", testConn, logger)
	if err != nil {
		t.Fatal(err)
	}

	cs = mockCommandService{t, cmd}

	d := &device.Device{
		UDID:         device.JsonNullString{sql.NullString{"00000000-1111-2222-3333-444455556666", true}},
		MDMTopic:     "mdmtopic",
		OSVersion:    "10.11",
		BuildVersion: "10G1000",
		ProductName:  "Mock Product",
		SerialNumber: device.JsonNullString{sql.NullString{"11111111", true}},
		Model:        "MockModel",
	}

	deviceUUID, err := devices.New("authenticate", d)
	if err != nil {
		t.Fatalf("creating fixture device: %s", err)
	}
	t.Logf("created mock device with UUID %v", deviceUUID)

	svc := NewService(devices, apps, certs, cs)
	handler := ServiceHandler(ctx, svc, logger)
	server := httptest.NewServer(handler)

	return &connectFixtures{
		db:         db,
		server:     server,
		svc:        svc,
		devices:    devices,
		apps:       apps,
		certs:      certs,
		cs:         cs,
		logger:     logger,
		deviceUUID: deviceUUID,
		migrator:   migrator,
	}
}

func teardown(fixtures *connectFixtures) {
	defer fixtures.server.Close()
	defer fixtures.db.Close()

	fixtures.migrator.RollbackAll()
}

func TestAcknowledgeDeviceInformation(t *testing.T) {
	cmd := mdm.CommandRequest{
		UDID:        "00000000-1111-2222-3333-444455556666",
		RequestType: "DeviceInformation",
	}

	fixtures := setup(t, &cmd)
	defer teardown(fixtures)

	requestBody, err := ioutil.ReadFile("../testdata/responses/macos/10.11.x/device_information.plist")
	if err != nil {
		t.Fatal(err)
	}

	client := http.DefaultClient
	theURL := fixtures.server.URL + "/mdm/connect"
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
}

func TestAcknowledgeInstalledApplicationList(t *testing.T) {
	cmd := mdm.CommandRequest{
		UDID:        "00000000-1111-2222-3333-444455556666",
		RequestType: "InstalledApplicationList",
	}

	fixtures := setup(t, &cmd)
	defer teardown(fixtures)

	requestBody, err := ioutil.ReadFile("../testdata/responses/macos/10.11.x/installed_application_list.plist")
	if err != nil {
		t.Fatal(err)
	}

	client := http.DefaultClient
	theURL := fixtures.server.URL + "/mdm/connect"
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

	var count int
	err = fixtures.db.QueryRow("SELECT COUNT(*) FROM devices_applications;").Scan(&count)
	if err != nil {
		t.Fatal(err)
	}

	if count != 3 {
		t.Error("expected number of inserted applications to be 3")
	}
}

func TestAcknowledgeCertificateList(t *testing.T) {
	cmd := mdm.CommandRequest{
		UDID:        "00000000-1111-2222-3333-444455556666",
		RequestType: "CertificateList",
	}

	fixtures := setup(t, &cmd)
	defer teardown(fixtures)

	requestBody, err := ioutil.ReadFile("../testdata/responses/macos/10.11.x/certificate_list.plist")
	if err != nil {
		t.Fatal(err)
	}

	client := http.DefaultClient
	theURL := fixtures.server.URL + "/mdm/connect"
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
}
