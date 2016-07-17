package applications

import (
	"database/sql"
	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

const MockUUID string = "ABCD-EFGH-IJKL"
const MockName string = "Mock Application"

var appFixtures []Application = []Application{
	{ // Normal macOS Application
		UUID:         "aba03d9d-6d80-4b96-bc6d-04233bddb26d",
		Name:         "Keychain Access",
		Identifier:   sql.NullString{"com.apple.keychainaccess", true},
		ShortVersion: sql.NullString{"9.0", true},
		Version:      sql.NullString{"9.0", true},
		BundleSize:   sql.NullInt64{14166172, true},
	},
	{ // macOS Application with no versioning available
		UUID:       "ddacb35f-6a6a-42eb-8ce6-aad55da5a237",
		Name:       "unetbootin",
		Identifier: sql.NullString{"com.yourcompany.unetbootin", true},
		BundleSize: sql.NullInt64{22292686, true},
	},
	{ // macOS Application with no bundle size available
		UUID:         "cdd950a1-f596-4777-a63a-9839c28e4d48",
		Name:         "FileMerge",
		Identifier:   sql.NullString{"com.apple.FileMerge", true},
		ShortVersion: sql.NullString{"2.9.1", true},
		Version:      sql.NullString{"2.9.1", true},
	},
	{ // macOS Application with no bundle identifier
		UUID:       "84c78174-8331-4cf3-98c3-a4b1434617e5",
		Name:       "Wireless Network Utility",
		BundleSize: sql.NullInt64{2416111, true},
	},
}

var logger log.Logger = log.NewNopLogger()

//func TestNewDB(t *testing.T) {
//	var logger logger.Logger = logger.NewNopLogger()
//	appsDB, err := NewDB("postgres", "host=localhost", logger)
//
//	if err != nil {
//		t.Error(err)
//	}
//
//	if _, ok := appsDB.(Datastore); !ok {
//		t.Log("Did not get a datastore")
//		t.Fail()
//	}
//}

func TestNewDatastore(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbx := sqlx.NewDb(db, "mock")
	defer dbx.Close()

	if _, err := NewDatastore(dbx, logger); err != nil {
		t.Error(err)
	}
}

func TestNewApplication(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbx := sqlx.NewDb(db, "mock")
	defer dbx.Close()

	appsDs, err := NewDatastore(dbx, logger)
	if err != nil {
		t.Error(err)
	}

	for _, fixture := range appFixtures {
		newRow := sqlmock.NewRows([]string{"application_uuid"}).AddRow(MockUUID)
		mock.ExpectQuery("INSERT INTO applications").WithArgs(
			fixture.Name,
			fixture.Identifier,
			fixture.ShortVersion,
			fixture.Version,
			fixture.BundleSize,
			nil,
			nil,
		).WillReturnRows(newRow)

		appUuid, err := appsDs.New(&fixture)
		if err != nil {
			t.Error(err)
		}

		if appUuid != MockUUID {
			t.Errorf("inserting a mock application did not return the mock uuid, got: %s", appUuid)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestApplications(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbx := sqlx.NewDb(db, "mock")
	defer dbx.Close()

	appsDs, err := NewDatastore(dbx, logger)
	if err != nil {
		t.Error(err)
	}

	mock.ExpectQuery(`SELECT .* FROM applications`).WillReturnRows(
		sqlmock.NewRows([]string{"application_uuid"}),
	)

	if _, err := appsDs.Applications(); err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestApplicationsWhereUUID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbx := sqlx.NewDb(db, "mock")
	defer dbx.Close()

	appsDs, err := NewDatastore(dbx, logger)
	if err != nil {
		t.Error(err)
	}

	mockRow := sqlmock.NewRows([]string{"application_uuid"}).AddRow(MockUUID)
	mock.ExpectQuery(`WHERE application_uuid =`).WillReturnRows(mockRow)

	apps, err := appsDs.Applications(UUID{MockUUID})
	if err != nil {
		t.Error(err)
	}

	if len(apps) != 1 {
		t.Fatalf("unexpected number of results returned: %d", len(apps))
	}

	if apps[0].UUID != MockUUID {
		t.Errorf("unexpected application uuid when querying: %s", apps[0].UUID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestApplicationsWhereName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbx := sqlx.NewDb(db, "mock")
	defer dbx.Close()

	appsDs, err := NewDatastore(dbx, logger)
	if err != nil {
		t.Error(err)
	}

	mockRow := sqlmock.NewRows([]string{"application_uuid", "name"}).AddRow(MockUUID, MockName)
	mock.ExpectQuery(`WHERE name =`).WillReturnRows(mockRow)

	apps, err := appsDs.Applications(Name{MockName})
	if err != nil {
		t.Error(err)
	}

	if len(apps) != 1 {
		t.Fatalf("unexpected number of results returned: %d", len(apps))
	}

	if apps[0].UUID != MockUUID {
		t.Errorf("unexpected application uuid when querying: %s", apps[0].UUID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetApplicationsByDeviceUUID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbx := sqlx.NewDb(db, "mock")
	defer dbx.Close()

	appsDs, err := NewDatastore(dbx, logger)
	if err != nil {
		t.Error(err)
	}

	mockRow := sqlmock.NewRows([]string{"application_uuid", "name"}).AddRow(MockUUID, MockName)
	mock.ExpectQuery(`WHERE devices_applications.device_uuid=`).WithArgs(MockUUID).WillReturnRows(mockRow)

	if _, err := appsDs.GetApplicationsByDeviceUUID(MockUUID); err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSaveApplicationByDeviceUUID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbx := sqlx.NewDb(db, "mock")
	defer dbx.Close()

	appsDs, err := NewDatastore(dbx, logger)
	if err != nil {
		t.Error(err)
	}

	for _, fixture := range appFixtures {
		mock.ExpectExec("INSERT INTO devices_applications").WithArgs(MockUUID, fixture.UUID).WillReturnResult(sqlmock.NewResult(1, 1))
		if err := appsDs.SaveApplicationByDeviceUUID(MockUUID, &fixture); err != nil {
			t.Error(err)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
