// Package driver provides utilities for creating reliablie connections to
// external services like redis and sql databases.
package driver

import "github.com/go-kit/kit/log"

const defaultAttempts = 20

// ConnOption is a driver connection option.
type ConnOption func(c *config)

type config struct {
	logger        log.Logger
	maxAttempts   int
	redisPassword string
}

// Logger adds a logger to the connection config.
func Logger(logger log.Logger) ConnOption {
	return func(c *config) {
		c.logger = logger
	}
}

// WithAttemtps t allows overriding the default 20 attempts for creating
// a driver connection.
func WithAttemtps(n int) ConnOption {
	return func(c *config) {
		c.maxAttempts = n
	}
}

// WithPassword adds an AUTH check when creating a redis pool.
func WithPassword(password string) ConnOption {
	return func(c *config) {
		c.redisPassword = password
	}
}
