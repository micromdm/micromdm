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
	SaveApplicationByDeviceUUID(deviceUUID string, app *Application) error
}

type pgStore struct {
	*sqlx.DB
	logger kitlog.Logger
}

func NewDatastore(connection *sqlx.DB, logger kitlog.Logger) (Datastore, error) {
	return pgStore{DB: connection, logger: logger}, nil
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
		VALUES ($0, $1, $2, $3, $4, $5, $6)
		ON CONFLICT (name, version) DO UPDATE SET
			identifier=$1,
			short_version=$2,
			bundle_size=$4,
			dynamic_size=$5,
			is_validated=$6
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
	stmt := `SELECT * FROM applications`
	stmt = addWhereFilters(stmt, "OR", params...)

	var apps []Application

	err := store.Select(&apps, stmt)
	if err != nil {
		return nil, errors.Wrap(err, "pgStore Applications")
	}
	return apps, nil
}

// Retrieve only applications which are installed on the given device.
func (store pgStore) GetApplicationsByDeviceUUID(deviceUUID string) (*[]Application, error) {
	var apps []Application
	query := `SELECT * FROM applications
	RIGHT JOIN devices_applications ON applications.application_uuid = devices_applications.application_uuid
	WHERE devices_applications.device_uuid=$1`

	err := store.Select(&apps, query, deviceUUID)

	if err != nil {
		return nil, err
	}

	return &apps, nil
}

// Associate the given applications with the given device uuid by inserting into `device_applications`.
func (store pgStore) SaveApplicationByDeviceUUID(deviceUUID string, app *Application) error {
	stmt := `INSERT INTO devices_applications (
		device_uuid, application_uuid
		) VALUES ($1, $2)`

	_, err := store.Exec(stmt, deviceUUID, app.UUID)
	return err
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

// boolean operators are applied to where conditions which are part of a whereClauseGroup
type booleanOperator string

const (
	OR  = "OR"
	AND = "AND"
)

type whereClauseGroup struct {
	Operator booleanOperator
	Clauses  []whereClause
}

// Get a string representing the where clause
// Second return value is an array of arguments to give to db.Exec etc.
func (cg whereClauseGroup) String() (string, []string) {
	var clauses []string
	var values []string = make([]string, len(cg.Clauses))

	for i, c := range cg.Clauses {
		c.Placeholder = fmt.Sprintf("$%d", i)
		clauses = append(clauses, c.String())
		values = append(values, c.Value)
	}

	return strings.Join(clauses, string(cg.Operator)), values
}

// Struct representation of a where clause. Does not deal with field name escaping or any inference of the value.
// I.E Do your own quoting.
type whereClause struct {
	Operator    string
	Field       string
	Value       string
	Placeholder string
}

func (c whereClause) String() string {
	return fmt.Sprintf(`%s %s %s`, c.Field, c.Operator, c.Value)
}

func Where(field string, operator string, value string) whereClause {
	return whereClause{
		Operator:    operator,
		Field:       field,
		Value:       value,
		Placeholder: "$1",
	}
}

func WhereAnd(clauses ...whereClause) whereClauseGroup {
	return whereClauseGroup{
		Operator: "AND",
		Clauses:  clauses,
	}
}

func WhereOr(clauses ...whereClause) whereClauseGroup {
	return whereClauseGroup{
		Operator: "OR",
		Clauses:  clauses,
	}
}
