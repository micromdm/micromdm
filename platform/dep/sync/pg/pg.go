package pg

import (
	"context"
	"strings"
	"database/sql"
	"time"
	"fmt"

	"github.com/pkg/errors"
	"github.com/jmoiron/sqlx"
	sq "gopkg.in/Masterminds/squirrel.v1"

	"github.com/micromdm/micromdm/platform/dep/sync"
	"github.com/micromdm/micromdm/platform/pubsub"
)

type Postgres struct{ db *sqlx.DB }

func NewDB(db *sqlx.DB, sub pubsub.Subscriber) (*Postgres, error) {
	// Required for TIMESTAMP DEFAULT 0
	_,err := db.Exec(`SET sql_mode = '';`)

	_,err = db.Exec(`CREATE TABLE IF NOT EXISTS cursors (
		    value VARCHAR(128) PRIMARY KEY,
		    created_at TIMESTAMPTZ DEFAULT '1970-01-01 00:00:00+00'
		);`)
	if err != nil {
	   return nil, errors.Wrap(err, "creating cursors sql table failed")
	}
	
	_,err = db.Exec(`CREATE TABLE IF NOT EXISTS dep_auto_assign (
		    profile_uuid VARCHAR(128) PRIMARY KEY,
		    filter TEXT DEFAULT NULL
		);`)
	if err != nil {
	   return nil, errors.Wrap(err, "creating cursors sql table failed")
	}
	
	store := &Postgres{db: db}
	return store, err
}

func columns() []string {
	return []string{
		"value",
		"created_at",
	}
}

const tableName = "cursors"

func (d *Postgres) LoadCursor(ctx context.Context) (*sync.Cursor, error) {
	
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
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

func (d *Postgres) SaveCursor(ctx context.Context, cursor sync.Cursor) error {
	// Make sure we take the time offset into account for "zero" dates	
	t := time.Now()
	_, offset := t.Zone()
	// Don't multiply by zero
	if (offset <= 0) {
		offset = 1
	}
	var min_timestamp_sec int64 = int64(offset) * 60 * 60 * 24
	
	if (cursor.CreatedAt.IsZero() || cursor.CreatedAt.Unix() < min_timestamp_sec) {
		cursor.CreatedAt = time.Unix(min_timestamp_sec, 0)
	}
	
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Insert(tableName).
		Columns(columns()...).
		Values(
			cursor.Value,
			cursor.CreatedAt,
		).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building cursor save query")
	}

	_, err = d.db.ExecContext(ctx, query, args...)
	return errors.Wrap(err, "exec cursor save in pg")
}

func (d *Postgres) SaveAutoAssigner(ctx context.Context, a *sync.AutoAssigner) error {
	fmt.Println("SaveAutoAssigner")
	if a.Filter != "*" {
		return errors.New("only '*' filter auto-assigners supported")
	}
	
	updateQuery, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Update("dep_auto_assign").
		Prefix("ON CONFLICT (profile_uuid) DO").
		Set("profile_uuid", a.ProfileUUID).
		Set("filter", a.Filter).
		ToSql()
		
	if err != nil {
		return errors.Wrap(err, "building update query for cursor save")
	}
	
	updateQuery = strings.Replace(updateQuery, "dep_auto_assign SET ", "", -1)

	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
		Insert("dep_auto_assign").
		Columns("profile_uuid", "filter").
		Values(
			a.ProfileUUID,
			a.Filter,
		).
		Suffix(updateQuery).
		ToSql()
	
	var all_args = append(args, args...)
	
	if err != nil {
		return errors.Wrap(err, "building cursor save query")
	}
	
	_, err = d.db.ExecContext(ctx, query, all_args...)
	
	return errors.Wrap(err, "exec cursor save in pg")
}

func (d *Postgres) DeleteAutoAssigner(ctx context.Context, filter string) error {
	fmt.Println("DeleteAutoAssigner")
	// TODO
	return nil
}

func (d *Postgres) LoadAutoAssigners(ctx context.Context) ([]sync.AutoAssigner, error) {
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