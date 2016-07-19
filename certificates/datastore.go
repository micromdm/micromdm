package certificates

import (
	"fmt"
	kitlog "github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
	"github.com/pkg/errors"
	"strings"
	"time"
)

var (
	insertCertificateStmt = `INSERT INTO devices_certificates (
		device_uuid,
		common_name,
		data,
		is_identity
	) VALUES ($1, $2, $3, $4)
	RETURNING certificate_uuid;`

	selectCertificatesStmt = `SELECT
		certificate_uuid,
		device_uuid
		common_name,
		data,
		is_identity
		FROM certificates`

	selectCertificatesByDeviceStmt = `SELECT
		certificate_uuid,
		certificates.device_uuid device_uuid
		common_name,
		data,
		is_identity
		FROM certificates
		INNER JOIN devices ON certificates.device_uuid = devices.device_uuid
		WHERE devices.udid = $1`
)

// This Datastore manages a list of certificates assigned to devices.
type Datastore interface {
	New(crt *Certificate) (string, error)
	Certificates(params ...interface{}) ([]Certificate, error)
	GetCertificatesByDeviceUDID(udid string) ([]Certificate, error)
}

type pgStore struct {
	*sqlx.DB
}

func NewDB(driver, conn string, logger kitlog.Logger) (Datastore, error) {
	switch driver {
	case "postgres":
		db, err := sqlx.Open(driver, conn)
		if err != nil {
			return nil, errors.Wrap(err, "certificates datastore")
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

func (store pgStore) New(c *Certificate) (string, error) {
	if err := store.QueryRow(insertCertificateStmt, c.DeviceUUID, c.CommonName, "", c.IsIdentity).Scan(&c.UUID); err != nil {
		return "", err
	}

	return c.UUID, nil
}

func (store pgStore) Certificates(params ...interface{}) ([]Certificate, error) {
	stmt := selectCertificatesStmt
	stmt = addWhereFilters(stmt, "OR", params...)
	var certificates []Certificate
	err := store.Select(&certificates, stmt)
	if err != nil {
		return nil, errors.Wrap(err, "pgStore Certificates")
	}
	return certificates, nil
}

func (store pgStore) GetCertificatesByDeviceUDID(udid string) ([]Certificate, error) {
	var certificates []Certificate
	err := store.Select(&certificates, selectCertificatesByDeviceStmt, udid)
	if err != nil {
		return nil, errors.Wrap(err, "pgStore GetCertificatesByDeviceUDID")
	}
	return certificates, nil
}

// UUID is a filter that can be added as a parameter to narrow down the list of returned results
type UUID struct {
	UUID string
}

func (p UUID) where() string {
	return fmt.Sprintf("certificate_uuid = '%s'", p.UUID)
}

// Filter by a device uuid
type DeviceUUID struct {
	UUID string
}

func (p DeviceUUID) where() string {
	return fmt.Sprintf("device_uuid = '%s'", p.UUID)
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
