package command

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/micromdm/flow/server/id"
	"github.com/pkg/errors"
	sq "gopkg.in/Masterminds/squirrel.v1"
)

type DeviceProfile struct {
	ID                string    `db:"id"`
	CommandUUID       string    `db:"command_uuid"`
	ProfileID         string    `db:"profile_id"`
	ProfileUUID       string    `db:"profile_uuid"`
	DeviceUDID        string    `db:"device_udid"`
	CreatedByWorkflow bool      `db:"created_by_workflow"`
	Acknowledged      bool      `db:"acknowledged"`
	AcknowledgedAt    time.Time `db:"acknowledged_at"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
	LastSeenOnDevice  time.Time `db:"last_seen_on_device"`
}

func columns() []string {
	return []string{
		"id",
		"command_uuid",
		"profile_id",
		"profile_uuid",
		"device_udid",
		"created_by_workflow",
		"acknowledged",
		"acknowledged_at",
		"created_at",
		"updated_at",
		"last_seen_on_device",
	}
}

type Postgres struct{ db *sqlx.DB }

func New(db *sqlx.DB) *Postgres {
	return &Postgres{db: db}
}

const tableName = "wf_device_profiles"

func (d *Postgres) UpdateFromList(ctx context.Context, udid, uuid string) error {
	id, now := id.New(), time.Now().UTC()
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(tableName).
		Columns(
			"id",
			"profile_uuid",
			"device_udid",
			"created_by_workflow",
			"created_at",
			"updated_at",
			"last_seen_on_device",
		).
		Values(
			id,
			uuid,
			udid,
			false,
			now,
			now,
			now,
		).
		Suffix("ON CONFLICT (profile_uuid, device_udid) DO UPDATE SET updated_at = ?, last_seen_on_device = ?", now, now).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building device save query")
	}

	_, err = d.db.ExecContext(ctx, query, args...)
	return errors.Wrap(err, "exec device save in pg")
}

func (d *Postgres) DeviceProfileByUUID(ctx context.Context, deviceUDID, uuid string) (DeviceProfile, error) {
	var dev DeviceProfile
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(columns()...).
		From(tableName).
		Where(
			sq.And{
				sq.Eq{"profile_uuid": uuid},
				sq.Eq{"device_udid": deviceUDID},
			},
		).
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

func (e notFoundErr) Error() string  { return "device profile not found" }
func (e notFoundErr) NotFound() bool { return true }

func isNotFound(err error) bool {
	err = errors.Cause(err)
	e, ok := errors.Cause(err).(interface {
		error
		NotFound() bool
	})
	return ok && e.NotFound()
}
