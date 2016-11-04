package checkin

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/DavidHuie/gomigrate"
	"github.com/go-kit/kit/log"
	"github.com/micromdm/micromdm/command"
	"github.com/micromdm/micromdm/device"
	"github.com/micromdm/micromdm/management"
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
