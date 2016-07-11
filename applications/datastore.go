package applications

import (
	"fmt"
	kitlog "github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
	"github.com/pkg/errors"
	"strings"
	"time"
)

// Datastore manages devices in a database
type Datastore interface {
	New(a *Application) (string, error)
	Applications(params ...interface{}) ([]Application, error)
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

// UUID is a filter that can be added as a parameter to narrow down the list of returned results
type UUID struct {
	UUID string
}

func (p UUID) where() string {
	return fmt.Sprintf("application_uuid = '%s'", p.UUID)
}

type Name struct {
	Name string
}

func (p Name) where() string {
	return fmt.Sprintf("name = '%s'", p.Name)
}

type Version struct {
	Version string
}

func (p Version) where() string {
	return fmt.Sprintf("version = '%s'", p.Version)
}

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

func (store pgStore) Applications(params ...interface{}) ([]Application, error) {
	stmt := `SELECT * FROM applications`
	stmt = addWhereFilters(stmt, "OR", params...)

	var apps []Application

	err := store.Select(&apps, stmt)
	if err != nil {
		return nil, errors.Wrap(err, "pgStore Applications")
	}
	return apps, nil
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

// whereer is for building args passed into a method which finds resources
type whereer interface {
	where() string
}

// add WHERE clause from params
func addWhereFilters(stmt string, separator string, params ...interface{}) string {
	var where []string
	for _, param := range params {
		if f, ok := param.(whereer); ok {
			where = append(where, f.where())
		}
	}

	if len(where) != 0 {
		whereFilter := strings.Join(where, " "+separator+" ")
		stmt = fmt.Sprintf("%s WHERE %s", stmt, whereFilter)
	}
	return stmt
}
