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

//func TestNewDB(t *testing.T) {
//	var log log.Logger = log.NewNopLogger()
//	appsDB, err := NewDB("postgres", "host=localhost", log)
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
	var log log.Logger = log.NewNopLogger()
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbx := sqlx.NewDb(db, "mock")
	defer dbx.Close()

	if _, err := NewDatastore(dbx, log); err != nil {
		t.Error(err)
	}
}

func TestNewApplication(t *testing.T) {
	var log log.Logger = log.NewNopLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbx := sqlx.NewDb(db, "mock")
	defer dbx.Close()

	appsDs, err := NewDatastore(dbx, log)
	if err != nil {
		t.Error(err)
	}

	// macOS style: no DynamicSize, no IsValidated
	fixture := Application{
		Name:         "Keychain Access",
		Identifier:   sql.NullString{"com.apple.keychainaccess", true},
		ShortVersion: sql.NullString{"9.0", true},
		Version:      sql.NullString{"9.0", true},
		BundleSize:   sql.NullInt64{14166172, true},
	}

	newRow := sqlmock.NewRows([]string{"application_uuid"}).AddRow(MockUUID)
	mock.ExpectQuery("INSERT INTO applications").WithArgs(
		fixture.Name,
		fixture.Identifier.String,
		fixture.ShortVersion.String,
		fixture.Version.String,
		fixture.BundleSize.Int64,
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

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestApplications(t *testing.T) {
	var log log.Logger = log.NewNopLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbx := sqlx.NewDb(db, "mock")
	defer dbx.Close()

	appsDs, err := NewDatastore(dbx, log)
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
	var log log.Logger = log.NewNopLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbx := sqlx.NewDb(db, "mock")
	defer dbx.Close()

	appsDs, err := NewDatastore(dbx, log)
	if err != nil {
		t.Error(err)
	}

	mockRow := sqlmock.NewRows([]string{"application_uuid"}).AddRow(MockUUID)
	mock.ExpectQuery(`WHERE application_uuid =`).WithArgs(MockUUID).WillReturnRows(mockRow)

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
	var log log.Logger = log.NewNopLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbx := sqlx.NewDb(db, "mock")
	defer dbx.Close()

	appsDs, err := NewDatastore(dbx, log)
	if err != nil {
		t.Error(err)
	}

	mockRow := sqlmock.NewRows([]string{"application_uuid", "name"}).AddRow(MockUUID, MockName)
	mock.ExpectQuery(`WHERE name =`).WithArgs(MockName).WillReturnRows(mockRow)

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
	var log log.Logger = log.NewNopLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbx := sqlx.NewDb(db, "mock")
	defer dbx.Close()

	appsDs, err := NewDatastore(dbx, log)
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
	var log log.Logger = log.NewNopLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbx := sqlx.NewDb(db, "mock")
	defer dbx.Close()

	appsDs, err := NewDatastore(dbx, log)
	if err != nil {
		t.Error(err)
	}

	fixture := Application{
		Name:         "Keychain Access",
		Identifier:   sql.NullString{"com.apple.keychainaccess", true},
		ShortVersion: sql.NullString{"9.0", true},
		Version:      sql.NullString{"9.0", true},
		BundleSize:   sql.NullInt64{14166172, true},
	}
	if err := appsDs.SaveApplicationByDeviceUUID(MockUUID, &fixture); err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
