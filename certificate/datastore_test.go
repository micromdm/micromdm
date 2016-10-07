package certificate

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

var (
	db   *sql.DB
	dbx  *sqlx.DB
	mock sqlmock.Sqlmock
	err  error
)

func setup() {
	db, mock, err = sqlmock.New()
	if err != nil {
		panic("an error was not expected when opening a stub database connection")
	}
	dbx = sqlx.NewDb(db, "mock")
}

func teardown() {
	dbx.Close()
}

func NewDatastore(connection *sqlx.DB) (Datastore, error) {
	return pgStore{DB: connection}, nil
}

func TestNewDatastore(t *testing.T) {
	setup()
	defer teardown()

	_, err := NewDatastore(dbx)
	if err != nil {
		t.Fatal(err)
	}
}
