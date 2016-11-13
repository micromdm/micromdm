// Package command manages an MDM Command queue for enrolled devices.
package command

import (
	"errors"

	"github.com/micromdm/mdm"
)

// Service defines methods for managing MDM commands in a queue.
type Service interface {

	// NewCommand turns an MDM Command Request into a MDM payload.
	NewCommand(*mdm.CommandRequest) (*mdm.Payload, error)

	// NextCommand retrieves the next command in a device's queue.
	NextCommand(udid string) ([]byte, int, error)

	// DeleteCommand deletes a previously queued command from a device's queue.
	DeleteCommand(deviceUDID, commandUUID string) (int, error)

	// Commands returns all the commands in a device's current queue.
	Commands(deviceUDID string) ([]mdm.Payload, error)

	// Find returns a previously queued command.
	Find(commandUUID string) (*mdm.Payload, error)
}

// TODO change this error to a type/interface
// ErrNoKey is returned if there is no key in redis
var ErrNoKey = errors.New("there is no such key in redis.")
