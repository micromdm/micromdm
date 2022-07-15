package pg

import (
	"context"
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	sq "gopkg.in/Masterminds/squirrel.v1"

	"github.com/micromdm/micromdm/platform/device"
)

type Postgres struct{ db *sqlx.DB }

func New(db *sqlx.DB) *Postgres {
	return &Postgres{db: db}
}

func columns() []string {
	return []string{
		"uuid",
		"udid",
		"serial_number",
		"os_version",
		"build_version",
		"product_name",
		"imei",
		"meid",
		"push_magic",
		"awaiting_configuration",
		"token",
		"unlock_token",
		"enrolled",
		"description",
		"model",
		"model_name",
		"device_name",
		"color",
		"asset_tag",
		"dep_profile_status",
		"dep_profile_uuid",
		"dep_profile_assign_time",
		"dep_profile_push_time",
		"dep_profile_assigned_date",
		"dep_profile_assigned_by",
		"last_seen",
	}
}

const tableName = "devices"

func (d *Postgres) Save(ctx context.Context, device *device.Device) error {
	updateQuery, _, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update(tableName).
		Prefix("ON CONFLICT (uuid) DO").
		Set("uuid", device.UUID).
		Set("udid", device.UDID).
		Set("serial_number", device.SerialNumber).
		Set("os_version", device.OSVersion).
		Set("build_version", device.BuildVersion).
		Set("product_name", device.ProductName).
		Set("imei", device.IMEI).
		Set("meid", device.MEID).
		Set("push_magic", device.PushMagic).
		Set("awaiting_configuration", device.AwaitingConfiguration).
		Set("token", device.Token).
		Set("unlock_token", device.UnlockToken).
		Set("enrolled", device.Enrolled).
		Set("description", device.Description).
		Set("model", device.Model).
		Set("model_name", device.ModelName).
		Set("device_name", device.DeviceName).
		Set("color", device.Color).
		Set("asset_tag", device.AssetTag).
		Set("dep_profile_status", device.DEPProfileStatus).
		Set("dep_profile_uuid", device.DEPProfileUUID).
		Set("dep_profile_assign_time", device.DEPProfileAssignTime).
		Set("dep_profile_push_time", device.DEPProfilePushTime).
		Set("dep_profile_assigned_date", device.DEPProfileAssignedDate).
		Set("dep_profile_assigned_by", device.DEPProfileAssignedBy).
		Set("last_seen", device.LastSeen).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building update query for device save")
	}
	updateQuery = strings.Replace(updateQuery, tableName, "", -1)

	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(tableName).
		Columns(columns()...).
		Values(
			device.UUID,
			device.UDID,
			device.SerialNumber,
			device.OSVersion,
			device.BuildVersion,
			device.ProductName,
			device.IMEI,
			device.MEID,
			device.PushMagic,
			device.AwaitingConfiguration,
			device.Token,
			device.UnlockToken,
			device.Enrolled,
			device.Description,
			device.Model,
			device.ModelName,
			device.DeviceName,
			device.Color,
			device.AssetTag,
			device.DEPProfileStatus,
			device.DEPProfileUUID,
			device.DEPProfileAssignTime,
			device.DEPProfilePushTime,
			device.DEPProfileAssignedDate,
			device.DEPProfileAssignedBy,
			device.LastSeen,
		).
		Suffix(updateQuery).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building device save query")
	}

	_, err = d.db.ExecContext(ctx, query, args...)
	return errors.Wrap(err, "exec device save in pg")
}

func (d *Postgres) DeviceByUDID(ctx context.Context, udid string) (*device.Device, error) {
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(columns()...).
		From(tableName).
		Where(sq.Eq{"udid": udid}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building sql")
	}

	var dev device.Device
	err = d.db.QueryRowxContext(ctx, query, args...).StructScan(&dev)
	if errors.Cause(err) == sql.ErrNoRows {
		return nil, deviceNotFoundErr{}
	}
	return &dev, errors.Wrap(err, "finding device by udid")
}

func (d *Postgres) DeviceBySerial(ctx context.Context, serial string) (*device.Device, error) {
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(columns()...).
		From(tableName).
		Where(sq.Eq{"serial_number": serial}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building sql")
	}

	var dev device.Device
	err = d.db.QueryRowxContext(ctx, query, args...).StructScan(&dev)
	if errors.Cause(err) == sql.ErrNoRows {
		return nil, deviceNotFoundErr{}
	}
	return &dev, errors.Wrap(err, "finding device by serial")
}

func (d *Postgres) ListDevices(ctx context.Context, opt device.ListDevicesOption) ([]device.Device, error) {
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(columns()...).
		From(tableName).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building sql")
	}
	var list []device.Device
	err = d.db.SelectContext(ctx, &list, query, args...)
	return list, errors.Wrap(err, "list devices")
}

func (d *Postgres) DeleteByUDID(ctx context.Context, udid string) error {
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Delete(tableName).
		Where(sq.Eq{"udid": udid}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}
	_, err = d.db.ExecContext(ctx, query, args...)
	return errors.Wrap(err, "delete device by udid")
}

func (d *Postgres) DeleteBySerial(ctx context.Context, serial string) error {
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Delete(tableName).
		Where(sq.Eq{"serial_number": serial}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}
	_, err = d.db.ExecContext(ctx, query, args...)
	return errors.Wrap(err, "delete device by serial_number")
}

type deviceNotFoundErr struct{}

func (e deviceNotFoundErr) Error() string {
	return "device not found"
}

func (e deviceNotFoundErr) NotFound() bool {
	return true
}
