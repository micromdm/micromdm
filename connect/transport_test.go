package connect

import (
	"bytes"
	"database/sql"
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

type mockCommandService struct {
	t *testing.T
}

func (svc mockCommandService) NewCommand(*mdm.CommandRequest) (*mdm.Payload, error) {
	return nil, nil
}

func (svc mockCommandService) NextCommand(udid string) ([]byte, int, error) {
	return nil, 0, nil
}

func (svc mockCommandService) DeleteCommand(deviceUDID, commandUUID string) (int, error) {
	return 1, nil
}

func (svc mockCommandService) Commands(deviceUDID string) ([]mdm.Payload, error) {
	return []mdm.Payload{}, nil
}

func (svc mockCommandService) Find(commandUUID string) (*mdm.Payload, error) {
	svc.t.Logf("Returning mock response finding command with UUID %s", commandUUID)
	cmd := mdm.CommandRequest{
		RequestType: "InstalledApplicationList",
	}
	payload, _ := mdm.NewPayload(&cmd)
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
}

func setup(t *testing.T) *connectFixtures {
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
	migrationErr := migrator.Migrate()
	if migrationErr != nil {
		t.Fatal(err)
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

	cs = mockCommandService{t}

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
		t.Fatal(err)
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
	}
}

func teardown(fixtures *connectFixtures) {
	defer fixtures.db.Close()
	defer fixtures.server.Close()
	drop := `
	DROP TABLE IF EXISTS devices;
	DROP INDEX IF EXISTS devices.serial_idx;
	DROP INDEX IF EXISTS devices.udid_idx;
	DROP TABLE IF EXISTS workflow_profile;
	DROP TABLE IF EXISTS workflow_workflow;
	DROP TABLE IF EXISTS workflows;
	DROP TABLE IF EXISTS profiles;
	DROP TABLE IF EXISTS applications;
	DROP TABLE IF EXISTS devices_applications;
	DROP TABLE IF EXISTS devices_certificates;
	`
	fixtures.db.Exec(drop)
}

func TestInstalledApplicationListResponse(t *testing.T) {
	fixtures := setup(t)
	defer teardown(fixtures)

	// create the faux command in the command service because connect will search for a match
	cmd := mdm.CommandRequest{
		UDID:        "00000000-1111-2222-3333-444455556666",
		RequestType: "InstalledApplicationList",
	}
	fixtures.cs.NewCommand(&cmd)

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

	t.Log("asserting correct number of applications have been inserted")
	var count int
	err = fixtures.db.QueryRow("SELECT COUNT(*) FROM devices_applications;").Scan(&count)
	if err != nil {
		t.Fatal(err)
	}

	if count != 3 {
		t.Fail()
	}
}
