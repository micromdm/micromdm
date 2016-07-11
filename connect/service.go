package connect

import (
	"fmt"
	"github.com/micromdm/mdm"
	"github.com/micromdm/micromdm/command"
	"github.com/micromdm/micromdm/device"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// Service defines methods for an MDM service
type Service interface {
	Acknowledge(ctx context.Context, req mdm.Response) (int, error)
	NextCommand(ctx context.Context, req mdm.Response) ([]byte, int, error)
	FailCommand(ctx context.Context, req mdm.Response) (int, error)

	RegisterAckHandler(predicate func(req mdm.Response) bool, handler func(req mdm.Response, datastores map[string]interface{}) error, datastores map[string]interface{})
	FindAckHandler(req mdm.Response) (func(req mdm.Response) error, bool)
	ExecAckHandler(req mdm.Response) error
}

// NewService creates a mdm service
func NewService(devices device.Datastore, cs command.Service) Service {
	return &service{
		commands: cs,
		devices:  devices,
		handlers: []ackHandler{},
	}
}

type ackHandler struct {
	predicate func(req mdm.Response) bool
	handler   func(req mdm.Response) error
}

type service struct {
	devices  device.Datastore
	commands command.Service
	handlers []ackHandler
}

func (svc service) Acknowledge(ctx context.Context, req mdm.Response) (int, error) {
	err := svc.ExecAckHandler(req)
	if err != nil {
		return 0, err
	}

	total, err := svc.commands.DeleteCommand(req.UDID, req.CommandUUID)
	if err != nil {
		return total, err
	}
	if total == 0 {
		total, err = svc.checkRequeue(req.UDID)
		if err != nil {
			return total, err
		}
		return total, nil
	}
	return total, nil
}

func (svc service) NextCommand(ctx context.Context, req mdm.Response) ([]byte, int, error) {
	return svc.commands.NextCommand(req.UDID)
}

func (svc service) FailCommand(ctx context.Context, req mdm.Response) (int, error) {
	return svc.commands.DeleteCommand(req.UDID, req.CommandUUID)
}

func (svc service) checkRequeue(deviceUDID string) (int, error) {
	existing, err := svc.devices.GetDeviceByUDID(deviceUDID, []string{"awaiting_configuration"}...)
	if err != nil {
		return 0, errors.Wrap(err, "check and requeue")
	}
	if existing.AwaitingConfiguration {
		cmdRequest := &mdm.CommandRequest{
			UDID:        deviceUDID,
			RequestType: "DeviceConfigured",
		}
		_, err := svc.commands.NewCommand(cmdRequest)
		if err != nil {
			return 0, err
		}
		return 1, nil
	}
	return 0, nil
}

// Register a handler function for a given request, include datastore dependencies as a map.
func (svc *service) RegisterAckHandler(predicate func(req mdm.Response) bool, handler func(req mdm.Response, datastores map[string]interface{}) error, datastores map[string]interface{}) {
	datastoreInjectedHandler := func(req mdm.Response) error {
		return handler(req, datastores)
	}
	newHandler := ackHandler{predicate: predicate, handler: datastoreInjectedHandler}
	svc.handlers = append(svc.handlers, newHandler)
}

// Find a handler function which is registered to deal with the RequestType
func (svc service) FindAckHandler(req mdm.Response) (func(req mdm.Response) error, bool) {
	for i, h := range svc.handlers {
		fmt.Println(i)
		if h.predicate(req) {
			return h.handler, true
		}
	}

	return nil, false
}

// Execute any registered handler function which matches the given RequestType
func (svc service) ExecAckHandler(req mdm.Response) error {
	handler, found := svc.FindAckHandler(req)
	if !found {
		return errors.New("There is no registered handler for the response.")
	}

	return handler(req)
}
