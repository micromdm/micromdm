package applications

import (
	"fmt"
	kitlog "github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
	"github.com/pkg/errors"
	"time"
)

// Datastore manages devices in a database
type Datastore interface {
	GetApplicationsByDeviceUUID(deviceUUID string) (*[]Application, error)
}

type pgStore struct {
	*sqlx.DB
}

func NewDB(driver, conn string, logger kitlog.Logger) (Datastore, error) {
	switch driver {
	case "postgres":
		db, err := sqlx.Open(driver, conn)
		if err != nil {
			return nil, errors.Wrap(err, "device datastore")
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
			return nil, errors.Wrap(dbError, "device datastore")
		}
		return pgStore{DB: db}, nil
	default:
		return nil, errors.New("unknown driver")
	}
}

func (store pgStore) New(src string, a *Application) (string, error) {
	return "", nil
}

func (store pgStore) GetApplicationsByDeviceUUID(deviceUUID string) (*[]Application, error) {
	apps := []Application{}
	query := `SELECT * FROM applications
	RIGHT JOIN devices_applications ON applications.application_uuid = devices_applications.application_uuid
	WHERE devices_applications.device_uuid=$1`

	err := store.Select(&apps, query, deviceUUID)

	if err != nil {
		return nil, err
	}

	return &apps, nil
}
