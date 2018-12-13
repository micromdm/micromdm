package mysql

import (
	"context"
	"strings"
	"database/sql"
	
	"fmt"

	"github.com/pkg/errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	sq "gopkg.in/Masterminds/squirrel.v1"

	"github.com/micromdm/micromdm/platform/apns"
	"github.com/micromdm/micromdm/platform/pubsub"
)

type Mysql struct{ db *sqlx.DB }

func NewDB(db *sqlx.DB, sub pubsub.Subscriber) (*Mysql, error) {
	
	// Required for TIMESTAMP DEFAULT 0
	_,err := db.Exec(`SET sql_mode = '';`)
	
	_,err = db.Exec(`CREATE TABLE IF NOT EXISTS push_info (
		    udid VARCHAR(40) PRIMARY KEY,
		    token TEXT DEFAULT '',
		    push_magic TEXT DEFAULT '',
		    mdm_topic TEXT DEFAULT ''
		);`)
		
	if err != nil {
	   return nil, errors.Wrap(err, "creating push_info table failed")
	}
	
	return &Mysql{db: db}, nil
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

func (d *Mysql) Save(ctx context.Context, i *apns.PushInfo) error {
	
	fmt.Println("apns.Save")
	
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
	
	fmt.Println(query)
	fmt.Println(args)

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