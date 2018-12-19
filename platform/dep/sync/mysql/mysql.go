package mysql

import (
	"context"
	"strings"
	"database/sql"
	"time"
	"fmt"

	"github.com/pkg/errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	sq "gopkg.in/Masterminds/squirrel.v1"

	"github.com/micromdm/micromdm/platform/dep/sync"
)

type Mysql struct{ db *sqlx.DB }

func NewDB(db *sqlx.DB) (*Mysql, error) {
	// Required for TIMESTAMP DEFAULT 0
	_,err := db.Exec(`SET sql_mode = '';`)

	_,err = db.Exec(`CREATE TABLE IF NOT EXISTS cursors (
		    value VARCHAR(128) PRIMARY KEY,
		    created_at TIMESTAMP DEFAULT 0
		);`)
	if err != nil {
	   return nil, errors.Wrap(err, "creating cursors sql table failed")
	}
	
	return &Mysql{db: db}, nil
}

func columns() []string {
	return []string{
		"value",
		"created_at",
	}
}

const tableName = "cursors"

func (d *Mysql) LoadCursor(ctx context.Context) (*sync.Cursor, error) {
	
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
		Select(columns()...).
		From(tableName).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building sql")
	}

	var cursor = struct {
		Cursor sync.Cursor `json:"cursor"`
	}{}
	
	err = d.db.QueryRowxContext(ctx, query, args...).StructScan(&cursor.Cursor)
	if errors.Cause(err) == sql.ErrNoRows {
		return &sync.Cursor{}, nil
		return nil, cursorNotFoundErr{}
	}
	return &cursor.Cursor, errors.Wrap(err, "loading cursor")
}

func (d *Mysql) SaveCursor(ctx context.Context, cursor sync.Cursor) error {
	// Make sure we take the time offset into account for "zero" dates	
	t := time.Now()
	_, offset := t.Zone()
	var min_timestamp_sec int64 = int64(offset) * 60 * 60 * 24
	
	if (cursor.CreatedAt.IsZero() || cursor.CreatedAt.Unix() < min_timestamp_sec) {
		cursor.CreatedAt = time.Unix(min_timestamp_sec, 0)
	}
	
	updateQuery, args_update, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
		Update(tableName).
		Prefix("ON DUPLICATE KEY").
		Set("value", cursor.Value).
		Set("created_at", cursor.CreatedAt).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building update query for cursor save")
	}
	
	// MySql Convention
	// Replace "ON DUPLICATE KEY UPDATE TABLE_NAME SET" to "ON DUPLICATE KEY UPDATE"
	updateQuery = strings.Replace(updateQuery, tableName+" SET ", "", -1)

	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
		Insert(tableName).
		Columns(columns()...).
		Values(
			cursor.Value,
			cursor.CreatedAt,
		).
		Suffix(updateQuery).
		ToSql()
	
	var all_args = append(args, args_update...)
	
	if err != nil {
		return errors.Wrap(err, "building cursor save query")
	}
	
	_, err = d.db.ExecContext(ctx, query, all_args...)
	
	return errors.Wrap(err, "exec cursor save in mysql")
}

func (d *Mysql) SaveAutoAssigner(ctx context.Context, a *sync.AutoAssigner) error {
	fmt.Println("SaveAutoAssigner")
	// TODO
	return nil
}

func (d *Mysql) DeleteAutoAssigner(ctx context.Context, filter string) error {
	fmt.Println("DeleteAutoAssigner")
	// TODO
	return nil
}

func (d *Mysql) LoadAutoAssigners(ctx context.Context) ([]sync.AutoAssigner, error) {
	fmt.Println("LoadAutoAssigners")
	var aa []sync.AutoAssigner
	return aa, nil
}





type cursorNotFoundErr struct{}

func (e cursorNotFoundErr) Error() string {
	return "cursor not found"
}

func (e cursorNotFoundErr) NotFound() bool {
	return true
}