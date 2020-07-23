package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Postgres provides methods for creating users in PostgreSQL.
type Postgres struct{ db *pgxpool.Pool }

// NewPostgres creates a Postgres client.
func NewPostgres(db *pgxpool.Pool) *Postgres {
	return &Postgres{db: db}
}

// CreateUser creates a new user in the database.
func (d *Postgres) CreateUser(ctx context.Context, username, email, password string) (*User, error) {
	u, err := create(username, email, password)
	if err != nil {
		return nil, fmt.Errorf("create postgres user: %w", err)
	}

	q :=
		`INSERT INTO users (
			id, username, email, password, salt, confirmation_hash
		) VALUES (
			$1, $2, $3, $4, $5, $6
		) RETURNING "created_at", "updated_at";`

	if err := d.db.QueryRow(ctx, q,
		u.ID,
		u.Username,
		u.Email,
		u.Password,
		u.Salt,
		u.ConfirmationHash,
	).Scan(&u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, fmt.Errorf("store created user in postgres: %w", checkPostgres(err))
	}

	return u, nil
}

func (d *Postgres) ConfirmUser(ctx context.Context, confirmation string) error {
	q :=
		`UPDATE users
		SET confirmation_hash = NULL
		WHERE confirmation_hash = $1;`

	if tag, err := d.db.Exec(ctx, q, confirmation); err != nil {
		return fmt.Errorf("set postgres confirmation_hash to NULL: %w", err)
	} else if tag.RowsAffected() == 0 {
		return errors.New("unknown confirmation_hash in postgres")
	}

	return nil
}

func checkPostgres(err error) error {
	var dbErr *pgconn.PgError
	if !errors.As(err, &dbErr) {
		return err
	}

	switch dbErr.Code {
	case "23514":
		if kv, ok := constraints[dbErr.ConstraintName]; ok {
			return Error{invalid: kv}
		}
	}

	return err
}
