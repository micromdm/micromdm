package inventory

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
	ID                       string    `db:"id"`
	CommandUUID              string    `db:"command_uuid"`
	ProfileIdentifier        string    `db:"profile_identifier"`
	ProfileUUID              string    `db:"profile_uuid"`
	DeviceUDID               string    `db:"device_udid"`
	HasRemovalPasscode       bool      `db:"has_removal_passcode"`
	IsEncrypted              bool      `db:"is_encrypted"`
	PayloadDescription       string    `db:"payload_description"`
	PayloadDisplayName       string    `db:"payload_display_name"`
	PayloadOrganization      string    `db:"payload_organization"`
	PayloadRemovalDisallowed bool      `db:"payload_removal_disallowed"`
	PayloadVersion           int       `db:"payload_version"`
	CreatedByWorkflow        bool      `db:"created_by_workflow"`
	Acknowledged             bool      `db:"acknowledged"`
	AcknowledgedAt           time.Time `db:"acknowledged_at"`
	CreatedAt                time.Time `db:"created_at"`
	UpdatedAt                time.Time `db:"updated_at"`
	LastSeenOnDevice         time.Time `db:"last_seen_on_device"`
}

func columns() []string {
	return []string{
		"id",
		"command_uuid",
		"profile_identifier",
		"profile_uuid",
		"device_udid",
		"has_removal_passcode",
		"is_encrypted",
		"payload_description",
		"payload_display_name",
		"payload_organization",
		"payload_removal_disallowed",
		"payload_version",
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

type ListProfilesResponse struct {
	ProfileList []Payload
}

type Payload struct {
	HasRemovalPasscode       bool
	IsEncrypted              bool
	PayloadDescription       string
	PayloadDisplayName       string
	PayloadIdentifier        string
	PayloadOrganization      string
	PayloadRemovalDisallowed bool
	PayloadUUID              string
	PayloadVersion           int
}

/*  UpdateFromListResponse updates the known profiles on a device aquired through a ListProfiles MDM query.
First all the profiles for the udid key will be removed, and then new profiles will be added.
The action is done in a transaction, so if the insert fails the removal would be rolled back.
Known issues:
- Updating with zero payloads would fail. It should not happen since the ProfileList
response returns at least the enrollment profile of the MDM. Still, the caller should validate.
- Updating with too many payloads would fail the bulk insert. What's too many? Postgres allows up to 65535 arguments, which roughly means over 9000 profiles with the current argument list.
*/
func (d *Postgres) UpdateFromListResponse(ctx context.Context, udid string, resp ListProfilesResponse) error {
	// bulk upsert.
	now := time.Now().UTC()
	deleteQuery, deleteArgs, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Delete(tableName).
		Where(sq.Eq{"device_udid": udid}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}

	tx, err := d.db.BeginTxx(ctx, nil)
	if _, err := tx.ExecContext(ctx, deleteQuery, deleteArgs...); err != nil {
		tx.Rollback()
		return err
	}

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(tableName).
		Columns(
			"id",
			"profile_uuid",
			"device_udid",
			"profile_identifier",
			"has_removal_passcode",
			"is_encrypted",
			"payload_description",
			"payload_display_name",
			"payload_organization",
			"payload_removal_disallowed",
			"payload_version",
			"created_by_workflow",
			"created_at",
			"updated_at",
			"last_seen_on_device",
		).
		Suffix("ON CONFLICT (profile_uuid, profile_identifier, device_udid) DO UPDATE SET updated_at = ?, last_seen_on_device = ?", now, now)

	// TODO: pg only supports 65535 params, which means about 9000 profiles.
	// there's likely bigger issues before someone reaches that many profiles, but you never know
	for _, p := range resp.ProfileList {
		id := id.New()
		builder = builder.Values(
			id,
			p.PayloadUUID,
			udid,
			p.PayloadIdentifier,
			p.HasRemovalPasscode,
			p.IsEncrypted,
			p.PayloadDescription,
			p.PayloadDisplayName,
			p.PayloadOrganization,
			p.PayloadRemovalDisallowed,
			p.PayloadVersion,
			false,
			now,
			now,
			now,
		)
	}
	query, args, err := builder.ToSql()
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "building update from list repsonse sql")
	}
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "exec update device profile from list profiles save in pg")
	}
	return tx.Commit()
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
