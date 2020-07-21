package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jackc/pgx/v4"
	"micromdm.io/v2/internal/frontend/account"
	"micromdm.io/v2/pkg/frontend"
	"micromdm.io/v2/pkg/log"
)

type server struct {
	ui *frontend.Server
}

func ui(f *cliFlags, logger log.Logger) (*frontend.Server, error) {
	return frontend.New(frontend.Config{
		Logger:         logger,
		SiteName:       f.siteName,
		CSRFKey:        []byte(f.csrfKey),
		CSRFCookieName: f.csrfCookieName,
		CSRFFieldName:  f.csrfFieldName,
	})
}

func setup(ctx context.Context, f *cliFlags, logger log.Logger) (*server, error) {
	uisrv, err := ui(f, logger)
	if err != nil {
		return nil, err
	}

	var db database
	switch dbDriver(f.databaseURL) {
	case "postgres":
		db.pg, err = setupPostgres(ctx, f, logger)
	case "sqlite":
		db.sq, err = setupSQLite(ctx, f, logger)
	default:
		err = fmt.Errorf("unsupported database_url value or restricted path %q", f.databaseURL)
	}

	if err != nil {
		return nil, err
	}

	srv := &server{ui: uisrv}

	account.HTTP(account.Config{
		HTTP:      srv.ui,
		UserStore: db.userdb().(account.UserStore),
	})

	return srv, nil
}

func dbDriver(dbURL string) string {
	if _, err := pgx.ParseConfig(dbURL); err == nil {
		return "postgres"
	}

	if _, err := os.Stat(dbURL); err == nil {
		return "sqlite"
	}

	if err := ioutil.WriteFile(dbURL, []byte(""), 0644); err == nil {
		os.Remove(dbURL)
		return "sqlite"
	}

	return ""
}
