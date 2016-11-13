package redis

import (
	"bytes"

	"github.com/garyburd/redigo/redis"
	"github.com/go-kit/kit/log"
	"github.com/groob/plist"

	"github.com/micromdm/mdm"
	"github.com/micromdm/micromdm/command"
)

// NewCommandService creates a command.Service backed by redis.
func NewCommandService(pool *redis.Pool, logger log.Logger) (Redis, error) {
	return Redis{pool: pool}, nil
}

// Redis implements command.Service
type Redis struct {
	pool *redis.Pool
}

func (rds Redis) NewCommand(request *mdm.CommandRequest) (*mdm.Payload, error) {
	// create a payload
	payload, err := mdm.NewPayload(request)
	if err != nil {
		return nil, err
	}
	// save in redis
	err = rds.SavePayload(payload)
	if err != nil {
		return nil, err
	}
	// add command to a queue in redis
	err = rds.QueueCommand(request.UDID, payload.CommandUUID)
	if err != nil {
		return nil, err
	}
	// return created payload to user
	return payload, nil
}

func (rds Redis) SavePayload(payload *mdm.Payload) error {
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

func (rds Redis) QueueCommand(deviceUDID, commandUUID string) error {
	// get connection from redis pool
	conn := rds.pool.Get()
	defer conn.Close()
	_, err := conn.Do("lpush", deviceUDID, commandUUID)
	if err != nil {
		return err
	}
	return nil
}
func (rds Redis) NextCommand(deviceUDID string) ([]byte, int, error) {
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
	cmd, err := redis.String(conn.Do("get", commandUUID))
	if err == redis.ErrNil {
		return nil, 0, command.ErrNoKey
	}

	// get a command list length
	total, err := redis.Int(conn.Do("llen", deviceUDID))
	if err != nil {
		return nil, 0, err
	}
	return []byte(cmd), total, err
}

func (rds Redis) DeleteCommand(deviceUDID, commandUUID string) (int, error) {
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

func (rds Redis) Commands(deviceUDID string) ([]mdm.Payload, error) {
	conn := rds.pool.Get()
	defer conn.Close()

	commandUUIDs, err := redis.Values(conn.Do("LRANGE", deviceUDID, "0", "-1"))
	if err != nil {
		return nil, err
	}

	// FIXME this code is going to result in an err if a command is deleted
	// mid-loop by another process.
	var payloads []mdm.Payload = make([]mdm.Payload, len(commandUUIDs))
	for i, commandUUID := range commandUUIDs {
		payloadData, err := redis.Bytes(conn.Do("GET", commandUUID))
		if err != nil {
			return nil, err
		}
		if err := plist.Unmarshal(payloadData, &payloads[i]); err != nil {
			return nil, err
		}
	}
	return payloads, nil
}

func (rds Redis) Find(commandUUID string) (*mdm.Payload, error) {
	conn := rds.pool.Get()
	defer conn.Close()

	payloadData, err := redis.Bytes(conn.Do("GET", commandUUID))
	if err != nil {
		return nil, err
	}

	var payload *mdm.Payload
	err = plist.Unmarshal(payloadData, payload)
	return payload, err
}
