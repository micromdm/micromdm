package session

import (
	"context"
	"errors"
	"fmt"
	"time"

	"crawshaw.io/sqlite/sqlitex"
	"micromdm.io/v2/pkg/viewer"
)

// SQLite provides methods for creating users in SQLite.
type SQLite struct{ db *sqlitex.Pool }

// NewSQLite creates a SQLite client.
func NewSQLite(db *sqlitex.Pool) *SQLite {
	return &SQLite{db: db}
}

func (d *SQLite) CreateSession(ctx context.Context) (*Session, error) {
	s, err := create(ctx)
	if err != nil {
		return nil, err
	}

	conn := d.db.Get(ctx)
	if conn == nil {
		return nil, context.Canceled
	}
	defer d.db.Put(conn)

	stmt := conn.Prep(`INSERT INTO sessions ( id, user_id ) VALUES ( $id, $user_id );`)
	stmt.SetText("$id", s.ID)
	stmt.SetText("$user_id", s.UserID)
	if _, err := stmt.Step(); err != nil {
		return nil, fmt.Errorf("store new session in sqlite: %w", err)
	}

	stmt = conn.Prep(`SELECT created_at, accessed_at FROM sessions WHERE id = $id`)
	stmt.SetText("$id", s.ID)
	if _, err := stmt.Step(); err != nil {
		return nil, err
	}

	s.CreatedAt, err = time.Parse("2006-01-02 15:04:05", stmt.GetText("created_at"))
	if err != nil {
		return nil, err
	}

	s.AccessedAt, err = time.Parse("2006-01-02 15:04:05", stmt.GetText("accessed_at"))
	if err != nil {
		return nil, err
	}

	return s, stmt.Reset()
}

func (d *SQLite) DestroySession(ctx context.Context) error {
	v, ok := viewer.FromContext(ctx)
	if !ok || v.SessionID == "" {
		return fmt.Errorf("session: missing valid viewer %v", v)
	}

	conn := d.db.Get(ctx)
	if conn == nil {
		return context.Canceled
	}
	defer d.db.Put(conn)

	stmt := conn.Prep(`DELETE FROM sessions WHERE id = $id;`)
	stmt.SetText("$id", v.SessionID)
	if _, err := stmt.Step(); err != nil {
		return err
	}
	return nil
}

func (d *SQLite) FindSession(ctx context.Context, id string) (*Session, error) {
	conn := d.db.Get(ctx)
	if conn == nil {
		return nil, context.Canceled
	}
	defer d.db.Put(conn)

	stmt := conn.Prep(
		`UPDATE sessions
		SET accessed_at = CURRENT_TIMESTAMP
		WHERE id = $id;`)
	stmt.SetText("$id", id)
	if _, err := stmt.Step(); err != nil {
		return nil, err
	}

	if conn.Changes() == 0 {
		return nil, errors.New("unknown session.id in sqlite")
	}

	stmt = conn.Prep(`SELECT created_at, accessed_at, user_id FROM sessions WHERE id = $id`)
	stmt.SetText("$id", id)
	if _, err := stmt.Step(); err != nil {
		return nil, err
	}

	s := Session{ID: id}
	var err error
	s.CreatedAt, err = time.Parse("2006-01-02 15:04:05", stmt.GetText("created_at"))
	if err != nil {
		return nil, err
	}

	s.AccessedAt, err = time.Parse("2006-01-02 15:04:05", stmt.GetText("accessed_at"))
	if err != nil {
		return nil, err
	}

	s.UserID = stmt.GetText("user_id")
	return &s, stmt.Reset()
}
