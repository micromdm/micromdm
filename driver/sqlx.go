package driver

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
)

// NewSQLxDB creates a *sqlx.DB with a backoff timer. By default 20 attempts
// will be made, with a 1 second increasing interval.
func NewSQLxDB(driver, conn string, opts ...ConnOption) (*sqlx.DB, error) {
	conf := &config{
		logger:      log.NewNopLogger(),
		maxAttempts: defaultAttempts,
	}
	for _, opt := range opts {
		opt(conf)
	}
	db, err := sqlx.Open(driver, conn)
	if err != nil {
		return nil, err
	}

	var dbError error
	for attempts := 1; attempts <= conf.maxAttempts; attempts++ {
		dbError = db.Ping()
		if dbError == nil {
			break
		}
		sleep := time.Duration(attempts)
		conf.logger.Log("msg", fmt.Sprintf(
			"could not connect to %s: %v, sleeping %v", driver, dbError, sleep))
		time.Sleep(sleep * time.Second)
	}
	if dbError != nil {
		return nil, err
	}
	return db, nil
}
