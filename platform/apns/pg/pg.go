package pg

import (
	"context"
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	sq "gopkg.in/Masterminds/squirrel.v1"

	"github.com/micromdm/micromdm/platform/apns"
)

type Postgres struct{ db *sqlx.DB }

func New(db *sqlx.DB) *Postgres {
	return &Postgres{db: db}
}

func columns() []string {
	return []string{
		"udid",
		"push_magic",
		"token",
		"mdm_topic",
	}
}

const tableName = "push_info"

func (d *Postgres) Save(ctx context.Context, i *apns.PushInfo) error {
	updateQuery, _, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update(tableName).
		Prefix("ON CONFLICT (udid) DO").
		Set("udid", i.UDID).
		Set("push_magic", i.PushMagic).
		Set("token", i.Token).
		Set("mdm_topic", i.MDMTopic).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building update query for push_info save")
	}
	updateQuery = strings.Replace(updateQuery, tableName, "", -1)

	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
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

	_, err = d.db.ExecContext(ctx, query, args...)
	return errors.Wrap(err, "exec push_info save in pg")
}

func (d *Postgres) PushInfo(ctx context.Context, udid string) (*apns.PushInfo, error) {
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
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
