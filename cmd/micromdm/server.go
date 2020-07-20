package main

import (
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

func setup(f *cliFlags, logger log.Logger) (*server, error) {
	uisrv, err := ui(f, logger)
	if err != nil {
		return nil, err
	}

	srv := &server{ui: uisrv}

	account.HTTP(account.Config{
		HTTP: srv.ui,
	})

	return srv, nil
}
