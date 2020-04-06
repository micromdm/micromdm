package pg

import (
	"context"
	"strings"
	"database/sql"
	
	"github.com/pkg/errors"
	"github.com/jmoiron/sqlx"
	sq "gopkg.in/Masterminds/squirrel.v1"

	"github.com/micromdm/micromdm/platform/remove"
)

type Postgres struct{ db *sqlx.DB }
const tableName = "remove_device"

func NewDB(db *sqlx.DB) (*Postgres, error) {
	
	_,err := db.Exec(`CREATE TABLE IF NOT EXISTS remove_device (
		    udid VARCHAR(40) PRIMARY KEY
		);`)
		
	if err != nil {
	   return nil, errors.Wrap(err, "creating push_info table failed")
	}
	
	return &Postgres{db: db}, nil
}

func columns() []string {
	return []string{
		"udid",
	}
}

func (d *Postgres) DeviceByUDID(ctx context.Context, udid string) (*remove.Device, error) {
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select(columns()...).
		From(tableName).
		Where(sq.Eq{"udid": udid}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building sql")
	}

	var dev remove.Device
	
	err = d.db.QueryRowxContext(ctx, query, args...).StructScan(&dev)
	if errors.Cause(err) == sql.ErrNoRows {
		return nil, removeDeviceNotFoundErr{}
	}
	return &dev, errors.Wrap(err, "finding remove device by udid")
}

func (d *Postgres) Save(ctx context.Context, dev *remove.Device) error {
	updateQuery, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Update(tableName).
		Prefix("ON CONFLICT (udid) DO").
		Set("udid", dev.UDID).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building update query for push_info save")
	}
	updateQuery = strings.Replace(updateQuery, tableName, "", -1)

	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Insert(tableName).
		Columns(columns()...).
		Values(
			dev.UDID,
		).
		Suffix(updateQuery).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building push_info save query")
	}

	_, err = d.db.ExecContext(ctx, query, args...)
	return errors.Wrap(err, "exec remove_device save in pg")
}

func (d *Postgres) Delete(ctx context.Context, udid string) error {
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Delete(tableName).
		Where(sq.Eq{"udid": udid}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}
	_, err = d.db.ExecContext(ctx, query, args...)
	return errors.Wrap(err, "delete remove device by udid")
}

type removeDeviceNotFoundErr struct{}

func (e removeDeviceNotFoundErr) Error() string {
	return "remove device not found"
}

func (e removeDeviceNotFoundErr) NotFound() bool {
	return true
}
