package session

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
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
