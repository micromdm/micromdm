package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"crawshaw.io/sqlite"
	"crawshaw.io/sqlite/sqlitex"
	"github.com/jackc/pgx/v4/pgxpool"

	"micromdm.io/v2/internal/data/user"
	"micromdm.io/v2/pkg/log"
)

// Wrap all the database types here and create access methods.
// The caller can assert a minimal interface for each data package.
//
// Example:
//	account.Config{UserStore: db.userdb().(account.UserStore)}
//
// Code generation (or eventual generics) can make this cleaner.
type database struct {
	sq *sqlitedb
	pg *postgresdb
}

func (db *database) userdb() interface{} {
	if db.sq != nil {
		return db.sq.userdb
	}
	return db.pg.userdb
}

type sqlitedb struct {
	userdb *user.SQLite
}

func setupSQLite(ctx context.Context, f *cliFlags, logger log.Logger) (*sqlitedb, error) {
	conn, err := sqlite.OpenConn(f.databaseURL, 0)
	if err != nil {
		return nil, fmt.Errorf("open sqlite dbfile: %s", err)
	}

	if err := sqliteInit(conn); err != nil {
		conn.Close()
		return nil, fmt.Errorf("init sqlite: %s", err)
	}

	if err := conn.Close(); err != nil {
		return nil, fmt.Errorf("sqlite init close: %s", err)
	}

	pool, err := sqlitex.Open(f.databaseURL, 0, 24)
	if err != nil {
		return nil, fmt.Errorf("create sqlite pool: %s", err)
	}

	db := &sqlitedb{
		userdb: user.NewSQLite(pool),
	}

	log.Debug(logger).Log("msg", "connected to db", "backend", "sqlite")
	return db, nil
}

func sqliteInit(conn *sqlite.Conn) error {
	if err := sqlitex.ExecTransient(conn, "PRAGMA journal_mode=WAL;", nil); err != nil {
		return err
	}

	if err := sqlitex.ExecTransient(conn, "PRAGMA cache_size = -50000;", nil); err != nil {
		return err
	}

	var (
		migrations []byte
		script     = "internal/data/migrations/sqlite/initial_tables.sql"
	)

	// only load and run migrations script if the path exists.
	// temporary workaround for tests which will have to be
	// replaced once the schema needs to exist for the tests.
	if _, err := os.Stat(script); err == nil {
		migrations, err = ioutil.ReadFile(script)
		if err != nil {
			return err
		}
	}

	if err := sqlitex.ExecScript(conn, string(migrations)); err != nil {
		return err
	}

	return nil
}

type postgresdb struct {
	userdb *user.Postgres
}

func setupPostgres(ctx context.Context, f *cliFlags, logger log.Logger) (*postgresdb, error) {
	//dbpool, err := pgxpool.Connect(ctx, "host=localhost port=5432 user=app dbname=app password=secret sslmode=disable")
	dbpool, err := pgxpool.Connect(ctx, f.databaseURL)
	if err != nil {
		return nil, err
	}

	// temporary until migrations are set up
	migrations, err := ioutil.ReadFile("internal/data/migrations/postgres/initial_tables.sql")
	if _, err := dbpool.Exec(ctx, string(migrations)); err != nil {
		dbpool.Close()
		return nil, err
	}

	db := &postgresdb{
		userdb: user.NewPostgres(dbpool),
	}

	log.Debug(logger).Log("msg", "connected to db", "backend", "postgres")

	return db, nil
}
