package application

import (
	"database/sql"
	"github.com/go-kit/kit/log"
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

var (
	logger log.Logger
)

func setup() {
	logger = log.NewNopLogger()
}

func teardown() {

}

func TestNewDB(t *testing.T) {
	setup()
	defer teardown()

	appsDB, err := NewDB("postgres", "host=localhost", logger)

	if err != nil {
		t.Error(err)
	}

	if _, ok := appsDB.(Datastore); !ok {
		t.Log("Did not get a datastore")
		t.Fail()
	}
}
