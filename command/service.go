package command

import "github.com/micromdm/mdm"

// Service defines methods for managing MDM commands
type Service interface {
	NewCommand(*mdm.CommandRequest) (*mdm.Payload, error)
	NextCommand(udid string) ([]byte, int, error)
	DeleteCommand(deviceUDID, commandUUID string) (int, error)
	Commands(deviceUDID string) ([]mdm.Payload, error)
}

// NewService returns a new command service
func NewService(ds Datastore) Service {
	return &service{
		db: ds,
	}
}

type service struct {
	db Datastore
}

func (svc service) NewCommand(request *mdm.CommandRequest) (*mdm.Payload, error) {
	// create a payload
	payload, err := mdm.NewPayload(request)
	if err != nil {
		return nil, err
	}
	// save in redis
	err = svc.db.SavePayload(payload)
	if err != nil {
		return nil, err
	}
	// add command to a queue in redis
	err = svc.db.QueueCommand(request.UDID, payload.CommandUUID)
	if err != nil {
		return nil, err
	}
	// return created payload to user
	return payload, nil
}

// NextCommand returns an MDM Payload from a list of queued payloads
func (svc service) NextCommand(udid string) ([]byte, int, error) {
	return svc.db.NextCommand(udid)
}

// DeleteCommand returns an MDM Payload from a list of queued payloads
func (svc service) DeleteCommand(deviceUDID, commandUUID string) (int, error) {
	return svc.db.DeleteCommand(deviceUDID, commandUUID)
}

func (svc service) Commands(deviceUDID string) ([]mdm.Payload, error) {
	return svc.db.Commands(deviceUDID)
}
