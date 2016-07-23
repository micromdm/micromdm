package applications

import (
	"fmt"
	kitlog "github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
	"github.com/pkg/errors"
	"time"
)

// This Datastore manages a list of known applications, and their relationship to devices.
type Datastore interface {
	New(a *Application) (string, error)
	Applications(params ...interface{}) ([]Application, error)
	GetApplicationsByDeviceUUID(deviceUUID string) ([]Application, error)
	SaveApplicationByDeviceUUID(deviceUUID string, app *Application) error
}

type pgStore struct {
	*sqlx.DB
}

func NewDB(driver, conn string, logger kitlog.Logger) (Datastore, error) {
	switch driver {
	case "postgres":
		db, err := sqlx.Open(driver, conn)
		if err != nil {
			return nil, errors.Wrap(err, "applications datastore")
		}
		var dbError error
		maxAttempts := 20
		for attempts := 1; attempts <= maxAttempts; attempts++ {
			dbError = db.Ping()
			if dbError == nil {
				break
			}
			logger.Log("msg", fmt.Sprintf("could not connect to postgres: %v", dbError))
			time.Sleep(time.Duration(attempts) * time.Second)
		}
		if dbError != nil {
			return nil, errors.Wrap(dbError, "applications datastore")
		}
		return pgStore{DB: db}, nil
	default:
		return nil, errors.New("unknown driver")
	}
}

// This function inserts a new application into the applications table.
// Applications are uniquely identifier by both their name and their long form version because some do not have
// identifiers, and some do not have short versions.
func (store pgStore) New(a *Application) (string, error) {
	err := store.QueryRow(
		`INSERT INTO applications (
			name,
			identifier,
			short_version,
			version,
			bundle_size,
			dynamic_size,
			is_validated
			)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (name, version) DO UPDATE SET
			identifier=$2,
			short_version=$3,
			bundle_size=$5,
			dynamic_size=$6,
			is_validated=$7
		RETURNING application_uuid;`,
		a.Name,
		a.Identifier,
		a.ShortVersion,
		a.Version,
		a.BundleSize,
		a.DynamicSize,
		a.IsValidated,
	).Scan(&a.UUID)

	if err != nil {
		return "", err
	}

	return a.UUID, nil
}

// Retrieve a list of applications
func (store pgStore) Applications(params ...interface{}) ([]Application, error) {
	stmt := `SELECT
		application_uuid,
		name,
		identifier,
		short_version,
		version,
		bundle_size,
		dynamic_size,
		is_validated
	FROM applications`
	stmt = addWhereFilters(stmt, "OR", params...)

	var apps []Application

	err := store.Select(&apps, stmt)
	if err != nil {
		return nil, errors.Wrap(err, "pgStore Applications")
	}
	return apps, nil
}

// Retrieve only applications which are installed on the given device.
func (store pgStore) GetApplicationsByDeviceUUID(deviceUUID string) ([]Application, error) {
	if len(deviceUUID) == 0 {
		return nil, errors.New("empty uuid supplied to GetApplicationsByDeviceUUID")
	}

	var apps []Application
	query := `SELECT
		applications.application_uuid AS application_uuid,
		name,
		identifier,
		short_version,
		version,
		bundle_size,
		dynamic_size,
		is_validated
	FROM applications
	RIGHT JOIN devices_applications ON applications.application_uuid = devices_applications.application_uuid
	WHERE devices_applications.device_uuid=$1`

	err := store.Select(&apps, query, deviceUUID)
	if err != nil {
		return nil, err
	}

	return apps, nil
}

// Associate the given applications with the given device uuid by inserting into `device_applications`.
func (store pgStore) SaveApplicationByDeviceUUID(deviceUUID string, app *Application) error {
	if deviceUUID == "" {
		return errors.New("empty uuid supplied to SaveApplicationByDeviceUUID for deviceUUID")
	}

	if app.UUID == "" {
		return errors.New("empty uuid supplied to SaveApplicationByDeviceUUID for application")
	}

	stmt := `INSERT INTO devices_applications (
		device_uuid, application_uuid
		) VALUES ($1, $2)`

	_, err := store.Exec(stmt, deviceUUID, app.UUID)
	return err
}
