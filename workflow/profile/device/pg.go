package device

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	sq "gopkg.in/Masterminds/squirrel.v1"
)

type Postgres struct{ db *sqlx.DB }

func New(db *sqlx.DB) *Postgres {
	return &Postgres{db: db}
}

type Device struct {
	UDID        string    `db:"udid"`
	Certificate []byte    `db:"certificate"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func columns() []string {
	return []string{
		"udid",
		"certificate",
		"created_at",
		"updated_at",
	}
}

const tableName = "wf_profile_devices"

func (d *Postgres) Save(ctx context.Context, device Device) error {
	updateQuery, _, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update(tableName).
		Prefix("ON CONFLICT (udid) DO").
		Set("udid", device.UDID).
		Set("certificate", device.Certificate).
		Set("created_at", device.CreatedAt).
		Set("updated_at", device.UpdatedAt).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building update query for device save")
	}
	updateQuery = strings.Replace(updateQuery, tableName, "", -1)

	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(tableName).
		Columns(columns()...).
		Values(
			device.UDID,
			device.Certificate,
			device.CreatedAt,
			device.UpdatedAt,
		).
		Suffix(updateQuery).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building device save query")
	}

	_, err = d.db.ExecContext(ctx, query, args...)
	return errors.Wrap(err, "exec device save in pg")
}

func (d *Postgres) List(ctx context.Context) ([]Device, error) {
	// TODO add pagination
	var devices []Device
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(columns()...).
		From(tableName).
		ToSql()
	if err != nil {
		return devices, errors.Wrap(err, "building sql to list devices")
	}
	err = d.db.SelectContext(ctx, &devices, query, args...)
	if errors.Cause(err) == sql.ErrNoRows {
		return devices, notFoundErr{}
	}
	return devices, errors.Wrap(err, "finding all devices")
}

func (d *Postgres) DeviceByUDID(ctx context.Context, udid string) (Device, error) {
	var dev Device
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(columns()...).
		From(tableName).
		Where(sq.Eq{"udid": udid}).
		ToSql()
	if err != nil {
		return dev, errors.Wrap(err, "building sql")
	}

	err = d.db.QueryRowxContext(ctx, query, args...).StructScan(&dev)
	if errors.Cause(err) == sql.ErrNoRows {
		return dev, notFoundErr{}
	}
	return dev, errors.Wrap(err, "finding device by udid")
}

type notFoundErr struct{}

func (e notFoundErr) Error() string  { return "device not found" }
func (e notFoundErr) NotFound() bool { return true }
