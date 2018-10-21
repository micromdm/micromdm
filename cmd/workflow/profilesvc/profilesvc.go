package main

import (
	"net/http"
	"os"

	"github.com/go-kit/kit/auth/basic"
	"github.com/go-kit/kit/log"
	"github.com/kolide/kit/dbutil"
	"github.com/kolide/kit/logutil"

	"github.com/micromdm/flow/server/httputil"
	"github.com/micromdm/micromdm/workflow/profile"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)
	db, err := dbutil.OpenDBX(
		"postgres",
		"host=localhost port=5432 user=micromdm dbname=micromdm password=micromdm sslmode=disable",
		dbutil.WithLogger(logger),
		dbutil.WithMaxAttempts(3),
	)
	if err != nil {
		logutil.Fatal(logger, "err", err)
	}
	defer db.Close()

	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(10)

	svc := profile.New(db, profile.WithLogger(logger))
	basicAuthEndpointMiddleware := basic.AuthMiddleware("micromdm", "secret", "micromdm")
	e := profile.MakeServerEndpoints(svc, basicAuthEndpointMiddleware)
	r, options := httputil.NewRouter(logger)
	profile.RegisterHTTPHandlers(r, e, options...)

	err = http.ListenAndServe(":9000", r)
	logutil.Fatal(logger, "err", err)
}
