package driver

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/go-kit/kit/log"
)

// NewRedisPool creates a redis pool with a backoff timer. By default,
// 20 attempts will be made, with a 1 second increasing interval.
func NewRedisPool(conn string, opts ...ConnOption) (*redis.Pool, error) {
	conf := &config{
		logger:      log.NewNopLogger(),
		maxAttempts: defaultAttempts,
	}
	for _, opt := range opts {
		opt(conf)
	}

	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conn)
			if err != nil {
				return nil, err
			}
			if conf.redisPassword != "" {
				if _, err := c.Do("AUTH", conf.redisPassword); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	if err := checkRedisConn(pool, conf.maxAttempts, conf.logger); err != nil {
		return nil, err
	}
	return pool, nil
}

func checkRedisConn(pool *redis.Pool, maxAttempts int, logger log.Logger) error {
	conn := pool.Get()
	defer conn.Close()

	var dbError error
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		_, dbError = conn.Do("PING")
		if dbError == nil {
			break
		}
		sleep := time.Duration(attempts)
		logger.Log("msg", fmt.Sprintf(
			"could not connect to redis: %v, sleeping %v", dbError, sleep))
		time.Sleep(sleep * time.Second)
	}
	return dbError
}
