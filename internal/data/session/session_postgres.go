package session

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	"micromdm.io/v2/pkg/viewer"
)

// Postgres provides methods for creating sessions in PostgreSQL.
type Postgres struct{ db *pgxpool.Pool }

// NewPostgres creates a Postgres client.
func NewPostgres(db *pgxpool.Pool) *Postgres {
	return &Postgres{db: db}
}

func (d *Postgres) CreateSession(ctx context.Context) (*Session, error) {
	s, err := create(ctx)
	if err != nil {
		return nil, err
	}

	q := `INSERT INTO sessions ( id, user_id ) VALUES ( $1, $2 ) RETURNING "created_at", "accessed_at";`
	if err := d.db.QueryRow(ctx, q,
		s.ID,
		s.UserID,
	).Scan(&s.CreatedAt, &s.AccessedAt); err != nil {
		return nil, fmt.Errorf("store new session in postgres: %w", err)
	}

	return s, nil
}

func (d *Postgres) DestroySession(ctx context.Context) error {
	v, ok := viewer.FromContext(ctx)
	if !ok || v.SessionID == "" {
		return fmt.Errorf("session: missing valid viewer %v", v)
	}

	q := `DELETE FROM sessions WHERE id = $1;`

	if _, err := d.db.Exec(ctx, q, v.SessionID); err != nil {
		return err
	}

	return nil
}

func (d *Postgres) FindSession(ctx context.Context, id string) (*Session, error) {
	q := fmt.Sprintf(`WITH updated AS (
  UPDATE sessions SET accessed_at = now() at time zone 'utc'
  FROM (
    SELECT user_id FROM sessions
    WHERE id = $1
    LIMIT 1
    ) sub
    JOIN users u on u.id = sub.user_id
    WHERE sessions.id = $2
    RETURNING %s
  )
	SELECT %s from updated;`, // TODO: add username, full_name to session struct.
		strings.Join(prefixTable("sessions", columns()), `, `),
		strings.Join(columns(), `, `),
	)

	s := &Session{ID: id}
	if err := d.db.QueryRow(ctx, q, id, id).Scan(
		&s.ID,
		&s.UserID,
		&s.CreatedAt,
		&s.AccessedAt,
	); err != nil {
		return nil, err
	}

	return s, nil

}

func prefixTable(prefix string, ss []string) []string {
	result := make([]string, len(ss), cap(ss))
	for i, s := range ss {
		result[i] = fmt.Sprintf("%s.%s", prefix, s)
	}
	return result
}
