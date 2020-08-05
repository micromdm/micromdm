package user

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"crawshaw.io/sqlite"
	"crawshaw.io/sqlite/sqlitex"
)

// SQLite provides methods for creating users in SQLite.
type SQLite struct{ db *sqlitex.Pool }

// NewSQLite creates a SQLite client.
func NewSQLite(db *sqlitex.Pool) *SQLite {
	return &SQLite{db: db}
}

// CreateUser creates a new user in the database.
func (d *SQLite) CreateUser(ctx context.Context, username, email, password string) (*User, error) {
	u, err := create(username, email, password)
	if err != nil {
		return nil, err
	}

	conn := d.db.Get(ctx)
	if conn == nil {
		return nil, context.Canceled
	}
	defer d.db.Put(conn)

	stmt := conn.Prep(
		`INSERT INTO users (
			id, username, email, password, salt, confirmation_hash
		) VALUES (
			$id, $username, $email, $password, $salt, $confirmationHash);`)
	stmt.SetText("$id", u.ID)
	stmt.SetText("$username", u.Username)
	stmt.SetText("$email", u.Email)
	stmt.SetBytes("$password", u.Password)
	stmt.SetBytes("$salt", u.Salt)
	stmt.SetText("$confirmationHash", *u.ConfirmationHash)
	if _, err := stmt.Step(); err != nil {
		return nil, fmt.Errorf("store created user in sqlite: %w", checkSqlite(err))
	}

	stmt = conn.Prep(fmt.Sprintf(
		`SELECT %s FROM users WHERE id = $id;`, strings.Join(columns(), `, `)))
	stmt.SetText("$id", u.ID)
	if found, err := stmt.Step(); err != nil {
		return nil, err
	} else if !found {
		return nil, fmt.Errorf("user (email %q) not found in sqlite", email)
	}

	usr, err := sqliteUser(stmt)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (d *SQLite) ConfirmUser(ctx context.Context, confirmation string) error {
	conn := d.db.Get(ctx)
	if conn == nil {
		return context.Canceled
	}
	defer d.db.Put(conn)

	stmt := conn.Prep(
		`UPDATE users 
		 SET confirmation_hash = NULL 
		 WHERE confirmation_hash = $confirmationHash;`)

	stmt.SetText("$confirmationHash", confirmation)
	if _, err := stmt.Step(); err != nil {
		return fmt.Errorf("set sqlite confirmation_hash to NULL: %w", err)
	}

	if conn.Changes() == 0 {
		return errors.New("unknown confirmation_hash in sqlite")
	}

	return nil
}

func (d *SQLite) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	conn := d.db.Get(ctx)
	if conn == nil {
		return nil, context.Canceled
	}
	defer d.db.Put(conn)

	if email == "" {
		return nil, Error{invalid: constraints["chk_email_not_empty"]}
	}

	stmt := conn.Prep(fmt.Sprintf(
		`SELECT %s FROM users WHERE email = $email;`, strings.Join(columns(), `, `)))
	stmt.SetText("$email", email)
	if found, err := stmt.Step(); err != nil {
		return nil, err
	} else if !found {
		return nil, fmt.Errorf("user (email %q) not found in sqlite", email)
	}

	usr, err := sqliteUser(stmt)
	if err != nil {
		return nil, err
	}
	return usr, stmt.Reset()
}

func sqliteUser(stmt *sqlite.Stmt) (*User, error) {
	var (
		u   = new(User)
		err error
	)

	u.ID = stmt.GetText("id")
	u.Email = stmt.GetText("email")
	u.Password = []byte(stmt.GetText("password"))
	u.Salt = []byte(stmt.GetText("salt"))

	// TODO: do we really need the pointer?
	if hash := stmt.GetText("confirmation_hash"); hash != "" {
		u.ConfirmationHash = &hash
	}

	u.CreatedAt, err = time.Parse("2006-01-02 15:04:05", stmt.GetText("created_at"))
	if err != nil {
		return nil, err
	}

	u.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", stmt.GetText("updated_at"))
	if err != nil {
		return nil, err
	}

	return u, stmt.Reset()
}

func checkSqlite(err error) error {
	var sErr sqlite.Error
	if !errors.As(err, &sErr) {
		return err
	}

	switch sErr.Code {
	case sqlite.SQLITE_CONSTRAINT_CHECK:
		c := strings.Split(sErr.Msg, ": ")
		if kv, ok := constraints[c[len(c)-1]]; ok {
			return Error{invalid: kv}
		}
	}

	return fmt.Errorf("user: %w", err)
}
