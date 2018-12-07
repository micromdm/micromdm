package mysql

import (
	"context"
	"strings"
	"database/sql"
	
	//"fmt"

	"github.com/pkg/errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	sq "gopkg.in/Masterminds/squirrel.v1"

	"github.com/micromdm/micromdm/platform/config"
	"github.com/micromdm/micromdm/platform/pubsub"
)

type Mysql struct{ db *sqlx.DB }

func NewDB(db *sqlx.DB) *Mysql {
	return &Mysql{db: db}
}

func columns() []string {
	return []string{
		"push_certificate",
		"private_key",
	}
}

const tableName = "server_config"

func (d *Mysql) SavePushCertificate(cert, key []byte) error {
	updateQuery, args_update, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
		Update(tableName).
		Prefix("ON DUPLICATE KEY").
		Set("config_id", 1).
		Set("push_certificate", cert).
		Set("private_key", key).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building update query for server_config save")
	}
	// MySql Convention
	// Replace "ON DUPLICATE KEY UPDATE TABLE_NAME SET" to "ON DUPLICATE KEY UPDATE"
	updateQuery = strings.Replace(updateQuery, tableName+" SET ", "", -1)
	
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
		Insert(tableName).
		Columns(columns()...).
		Values(
			1,
			device.PushCertificate,
			device.PrivateKey,
		).
		Suffix(updateQuery).
		ToSql()
	
	var all_args = append(args, args_update...)
	
	if err != nil {
		return errors.Wrap(err, "building server_config save query")
	}
	
	_, err = d.db.ExecContext(ctx, query, all_args...)
	
	return errors.Wrap(err, "exec server_config save in mysql")
}

func (d *Mysql) serverConfig() (*config.ServerConfig, error) {
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
		Select(columns()...).
		From(tableName).
		Where(sq.Eq{"config_id": 1}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building sql")
	}

	var config config.ServerConfig
	
	err = d.db.QueryRowxContext(ctx, query, args...).StructScan(&config)
	if errors.Cause(err) == sql.ErrNoRows {
		return nil, deviceNotFoundErr{}
	}
	return &config, errors.Wrap(err, "finding config by config_id")
}













func (d *Mysql) Save(ctx context.Context, i *apns.PushInfo) error {
	updateQuery, args_update, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
		Update(tableName).
		Prefix("ON DUPLICATE KEY").
		Set("udid", i.UDID).
		Set("push_magic", i.PushMagic).
		Set("token", i.Token).
		Set("mdm_topic", i.MDMTopic).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building update query for push_info save")
	}
	updateQuery = strings.Replace(updateQuery, tableName+" SET ", "", -1)

	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
		Insert(tableName).
		Columns(columns()...).
		Values(
			i.UDID,
			i.PushMagic,
			i.Token,
			i.MDMTopic,
		).
		Suffix(updateQuery).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building push_info save query")
	}

	var all_args = append(args, args_update...)
	_, err = d.db.ExecContext(ctx, query, all_args...)
	return errors.Wrap(err, "exec push_info save in pg")
}

func (d *Mysql) PushInfo(ctx context.Context, udid string) (*apns.PushInfo, error) {
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
		Select(columns()...).
		From(tableName).
		Where(sq.Eq{"udid": udid}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building sql")
	}

	var i apns.PushInfo
	err = d.db.QueryRowxContext(ctx, query, args...).StructScan(&i)
	if errors.Cause(err) == sql.ErrNoRows {
		return nil, pushInfoNotFoundErr{}
	}
	return &i, errors.Wrap(err, "finding push_info by udid")
}

type pushInfoNotFoundErr struct{}

func (e pushInfoNotFoundErr) Error() string  { return "push_info not found" }
func (e pushInfoNotFoundErr) NotFound() bool { return true }
