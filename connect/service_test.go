package connect

import (
	"database/sql"
	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/micromdm/mdm"
	"github.com/micromdm/micromdm/applications"
	"github.com/micromdm/micromdm/device"
	"golang.org/x/net/context"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"os"
	"testing"
)

type MockDevices struct{}

func (md MockDevices) New(src string, d *device.Device) (string, error) {
	return "", nil
}
func (md MockDevices) GetDeviceByUDID(udid string, fields ...string) (*device.Device, error) {
	return &device.Device{
		UUID: "00000000-1111-2222-3333-444455556666",
	}, nil
}
func (md MockDevices) GetDeviceByUUID(uuid string, fields ...string) (*device.Device, error) {
	return &device.Device{}, nil
}
func (md MockDevices) Devices(params ...interface{}) ([]device.Device, error) {
	return []device.Device{
		{
			UUID: "00000000-1111-2222-3333-444455556666",
		},
	}, nil
}
func (md MockDevices) Save(msg string, dev *device.Device) error {
	return nil
}

type MockCmd struct{}

func (mc MockCmd) NewCommand(*mdm.CommandRequest) (*mdm.Payload, error) {
	return &mdm.Payload{}, nil
}
func (mc MockCmd) NextCommand(udid string) ([]byte, int, error) {
	return []byte{}, 0, nil
}
func (mc MockCmd) DeleteCommand(deviceUDID, commandUUID string) (int, error) {
	return 0, nil
}

var ctx context.Context
var appDs applications.Datastore
var db *sql.DB
var mock sqlmock.Sqlmock
var dbx *sqlx.DB
var logger log.Logger

func setupServiceTests() {
	ctx = context.Background()
	db, mock, _ = sqlmock.New()
	//if err != nil {
	//	t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	//}
	dbx = sqlx.NewDb(db, "mock")
	logger = log.NewLogfmtLogger(os.Stdout)
	appDs, _ = applications.NewDatastore(dbx, logger)
}

func teardownServiceTests() {
	dbx.Close()
}

func TestAckQueryResponses(t *testing.T) {
	setupServiceTests()
	defer teardownServiceTests()

	response := mdm.Response{
		UDID:           "00000000-1111-2222-3333-444455556666",
		Status:         "Acknowledged",
		CommandUUID:    "10000000-1111-2222-3333-444455556666",
		RequestType:    "DeviceInformation",
		QueryResponses: mdm.QueryResponses{},
	}

	mockDevices := MockDevices{}
	mockCmd := MockCmd{}

	svc := NewService(mockDevices, appDs, mockCmd)
	svc.Acknowledge(ctx, response)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAckInstalledApplicationList(t *testing.T) {
	setupServiceTests()
	defer teardownServiceTests()

	response := mdm.Response{
		UDID:        "00000000-1111-2222-3333-444455556666",
		Status:      "Acknowledged",
		CommandUUID: "10000000-1111-2222-3333-444455556666",
		RequestType: "InstalledApplicationList",
		InstalledApplicationList: []mdm.InstalledApplicationListItem{
			{
				Name:       "Wireless Network Utility",
				BundleSize: 2416111,
			},
			{
				Name:         "Keychain Access",
				Identifier:   "com.apple.keychainaccess",
				ShortVersion: "9.0",
				Version:      "9.0",
				BundleSize:   14166172,
			},
			{
				Name:       "Bundle Size Regression",
				BundleSize: 2463209237,
			},
		},
	}

	mockDevices := MockDevices{}
	mockCmd := MockCmd{}

	// Expect insert applications as new
	wifiAppUuidRow := sqlmock.NewRows([]string{"application_uuid"}).AddRow("90000000-1111-2222-3333-444455556666")
	mock.ExpectQuery("INSERT INTO applications").WithArgs("Wireless Network Utility", nil, nil, nil, 2416111, 0, nil).WillReturnRows(wifiAppUuidRow)
	kcAppUuidRow := sqlmock.NewRows([]string{"application_uuid"}).AddRow("A0000000-1111-2222-3333-444455556666")
	mock.ExpectQuery("INSERT INTO applications").WithArgs("Keychain Access", "com.apple.keychainaccess", "9.0", "9.0", 14166172, 0, nil).WillReturnRows(kcAppUuidRow)
	sizeAppUuidRow := sqlmock.NewRows([]string{"application_uuid"}).AddRow("B0000000-1111-2222-3333-444455556666")
	mock.ExpectQuery("INSERT INTO applications").WithArgs("Bundle Size Regression", nil, nil, nil, 2463209237, 0, nil).WillReturnRows(sizeAppUuidRow)

	// Expect query for device installed apps
	deviceAppsRow := sqlmock.NewRows([]string{"application_uuid", "name"}).AddRow("APP00000-1111-2222-3333-444455556666", "Mock Application")
	mock.ExpectQuery("RIGHT JOIN devices_applications").WithArgs("00000000-1111-2222-3333-444455556666").WillReturnRows(deviceAppsRow)

	mock.ExpectExec("INSERT INTO devices_applications").WithArgs("00000000-1111-2222-3333-444455556666", "90000000-1111-2222-3333-444455556666").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO devices_applications").WithArgs("00000000-1111-2222-3333-444455556666", "A0000000-1111-2222-3333-444455556666").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO devices_applications").WithArgs("00000000-1111-2222-3333-444455556666", "B0000000-1111-2222-3333-444455556666").WillReturnResult(sqlmock.NewResult(1, 1))

	svc := NewService(mockDevices, appDs, mockCmd)
	svc.Acknowledge(ctx, response)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TODO: A regression exists where a device reports the installed application list twice and apps are duplicated.
func TestAckInstalledApplicationListDuplicateRegression(t *testing.T) {

}
