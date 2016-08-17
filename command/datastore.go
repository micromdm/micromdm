package command

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
	kitlog "github.com/go-kit/kit/log"
	"github.com/groob/plist"
	"github.com/micromdm/mdm"
)

var (
	// ErrNoKey is returned if there is no key in redis
	ErrNoKey = errors.New("There is no such key in redis.")
)

// Datastore provides methods for saving and retrieving MDM commands
type Datastore interface {
	// Saves the payload in redis
	// SET CommandUUID plistData
	SavePayload(payload *mdm.Payload) error
	// Adds MDM commands to a queue in redis list
	// LPUSH deviceUDID commandUUID
	QueueCommand(deviceUDID, commandUUID string) error
	NextCommand(deviceUDID string) ([]byte, int, error)
	DeleteCommand(deviceUDID, commandUUID string) (int, error)
	Commands(deviceUDID string) ([]mdm.Payload, error)
	Find(commandUUID string) (*mdm.Payload, error)
}

//NewDB creates a Datastore
func NewDB(driver, conn string, logger kitlog.Logger) (Datastore, error) {
	var ds Datastore
	switch driver {
	case "redis":
		ds = redisDB{pool: redisPool(conn, logger)}
		return ds, nil
	default:
		return nil, errors.New("unknown driver")
	}
}

type redisDB struct {
	pool *redis.Pool
}

func (rds redisDB) SavePayload(payload *mdm.Payload) error {
	var buf bytes.Buffer
	// get connection from redis pool
	conn := rds.pool.Get()
	defer conn.Close()
	// encode payload into a plist
	err := plist.NewEncoder(&buf).Encode(payload)
	if err != nil {
		return err
	}
	// create a commandUUID key with the plist as the value
	_, err = conn.Do("set", payload.CommandUUID, buf.String())
	if err != nil {
		return err
	}
	return nil
}

func (rds redisDB) QueueCommand(deviceUDID, commandUUID string) error {
	// get connection from redis pool
	conn := rds.pool.Get()
	defer conn.Close()
	_, err := conn.Do("lpush", deviceUDID, commandUUID)
	if err != nil {
		return err
	}
	return nil
}
func (rds redisDB) NextCommand(deviceUDID string) ([]byte, int, error) {
	// get connection from redis pool
	conn := rds.pool.Get()
	defer conn.Close()
	// pop the first command
	commandUUID, err := redis.String(conn.Do("lpop", deviceUDID))
	if err != nil && err != redis.ErrNil {
		return nil, 0, err
	}
	// if the list is empty
	if err == redis.ErrNil {
		return []byte{}, 0, nil
	}
	// push the redis command back to the end of the list
	_, err = conn.Do("rpush", deviceUDID, commandUUID)
	command, err := redis.String(conn.Do("get", commandUUID))
	if err == redis.ErrNil {
		return nil, 0, ErrNoKey
	}

	// get a command list length
	total, err := redis.Int(conn.Do("llen", deviceUDID))
	if err != nil {
		return nil, 0, err
	}
	return []byte(command), total, err
}

func (rds redisDB) DeleteCommand(deviceUDID, commandUUID string) (int, error) {
	// get connection from redis pool
	conn := rds.pool.Get()
	defer conn.Close()
	// remove from list
	_, err := conn.Do("lrem", deviceUDID, 0, commandUUID)
	if err != nil {
		return 0, err
	}
	// set the key to expire in an hour
	_, err = conn.Do("expire", commandUUID, 3600)
	if err != nil {
		return 0, err
	}
	// get a command list length
	total, err := redis.Int(conn.Do("llen", deviceUDID))
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (rds redisDB) Commands(deviceUDID string) ([]mdm.Payload, error) {
	conn := rds.pool.Get()
	defer conn.Close()

	commandUuids, err := redis.Values(conn.Do("LRANGE", deviceUDID, "0", "-1"))
	if err != nil {
		return nil, err
	}

	var payloads []mdm.Payload = make([]mdm.Payload, len(commandUuids))

	for i, commandUuid := range commandUuids {
		payloadData, err := redis.Bytes(conn.Do("GET", commandUuid))
		if err != nil {
			return nil, err
		}

		if err := plist.NewDecoder(bytes.NewReader(payloadData)).Decode(&payloads[i]); err != nil {
			return nil, err
		}
	}

	return payloads, nil
}

func (rds redisDB) Find(commandUUID string) (*mdm.Payload, error) {
	conn := rds.pool.Get()
	defer conn.Close()

	payloadData, err := redis.Bytes(conn.Do("GET", commandUUID))
	if err != nil {
		return nil, err
	}

	var payload *mdm.Payload
	if err := plist.NewDecoder(bytes.NewReader(payloadData)).Decode(&payload); err != nil {
		return nil, err
	}

	return payload, nil
}

func redisPool(conn string, logger kitlog.Logger) *redis.Pool {
	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conn)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	checkRedisConn(pool, logger)
	return pool
}

func checkRedisConn(pool *redis.Pool, logger kitlog.Logger) {
	conn := pool.Get()
	defer conn.Close()

	var dbError error
	maxAttempts := 20
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		_, dbError = conn.Do("PING")
		if dbError == nil {
			break
		}
		logger.Log("msg", fmt.Sprintf("could not connect to redis: %v", dbError))
		time.Sleep(time.Duration(attempts) * time.Second)
	}
	if dbError != nil {
		logger.Log("err", dbError)
		os.Exit(1)
	}
}
