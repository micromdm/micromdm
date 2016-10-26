package connect

import (
	"database/sql"
	"github.com/DavidHuie/gomigrate"
	"github.com/go-kit/kit/log"
	"github.com/micromdm/micromdm/application"
	"github.com/micromdm/micromdm/certificate"
	"github.com/micromdm/micromdm/command"
	"github.com/micromdm/micromdm/device"
	"golang.org/x/net/context"
	"net/http/httptest"
	"os"
	"testing"
)

type connectFixtures struct {
	db      *sql.DB
	server  *httptest.Server
	svc     Service
	devices device.Datastore
	apps    application.Datastore
	certs   certificate.Datastore
	cs      command.Service
}

func setup(t *testing.T) *connectFixtures {
	ctx := context.Background()
	l := log.NewLogfmtLogger(os.Stderr)
	logger := log.NewContext(l).With("source", "testing")

	var (
		err           error
		testConn      string = "user=postgres password= dbname=travis_ci_test sslmode=disable"
		testRedisConn string = ":6379"
		devices       device.Datastore
		apps          application.Datastore
		certs         certificate.Datastore
		cs            command.Service
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

	commandDB, err := command.NewDB("redis", testRedisConn, logger)
	if err != nil {
		t.Fatal(err)
	}
	cs = command.NewService(commandDB)

	svc := NewService(devices, apps, certs, cs)
	handler := ServiceHandler(ctx, svc, logger)
	server := httptest.NewServer(handler)

	return &connectFixtures{
		db:      db,
		server:  server,
		svc:     svc,
		devices: devices,
		apps:    apps,
		certs:   certs,
		cs:      cs,
	}
}

func teardown(fixtures *connectFixtures) {
	defer fixtures.db.Close()
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

func TestDeviceInformationResponse(t *testing.T) {
	fixtures := setup(t)
	defer teardown(fixtures)

}
