package connect

import (
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

	RegisterAckHandler(requestType string, handler func(req mdm.Response, datastores map[string]interface{}) error, datastores map[string]interface{})
	FindAckHandler(requestType string) (func(req mdm.Response) error, bool)
	ExecAckHandler(requestType string, req mdm.Response) error
}

// NewService creates a mdm service
func NewService(devices device.Datastore, cs command.Service) Service {
	return &service{
		commands: cs,
		devices:  devices,
	}
}

type ackHandler struct {
	requestType string
	handler     func(req mdm.Response) error
}

type service struct {
	devices  device.Datastore
	commands command.Service
	handlers []ackHandler
}

func (svc service) Acknowledge(ctx context.Context, req mdm.Response) (int, error) {

	err := svc.ExecAckHandler(req.RequestType, req)

	//// Need to handle the absence of RequestType in IOS8 devices
	//if req.QueryResponses.UDID != "" {
	//	if err := svc.ackQueryResponses(req); err != nil {
	//		return 0, err
	//	}
	//}

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

func (svc service) RegisterAckHandler(requestType string, handler func(req mdm.Response, datastores map[string]interface{}) error, datastores map[string]interface{}) {
	datastoreInjectedHandler := func(req mdm.Response) error {
		return handler(req, datastores)
	}
	svc.handlers = append(svc.handlers, ackHandler{requestType, datastoreInjectedHandler})
}

// If not found, second return variable is false
func (svc service) FindAckHandler(requestType string) (func(req mdm.Response) error, bool) {
	for _, h := range svc.handlers {
		if h.requestType == requestType {
			return h.handler, true
		}
	}

	return nil, false
}

func (svc service) ExecAckHandler(requestType string, req mdm.Response) error {
	handler, found := svc.FindAckHandler(requestType)
	if !found {
		return errors.Errorf("There is no registered handler for the response type: %s", requestType)
	}

	return handler(req)
}
