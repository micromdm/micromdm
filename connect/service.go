package connect

import (
	"encoding/json"
	"github.com/micromdm/mdm"
	apps "github.com/micromdm/micromdm/applications"
	"github.com/micromdm/micromdm/certificates"
	"github.com/micromdm/micromdm/command"
	"github.com/micromdm/micromdm/device"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"time"
)

// Service defines methods for an MDM service
type Service interface {
	Acknowledge(ctx context.Context, req mdm.Response) (int, error)
	NextCommand(ctx context.Context, req mdm.Response) ([]byte, int, error)
	FailCommand(ctx context.Context, req mdm.Response) (int, error)
}

// NewService creates a mdm service
func NewService(devices device.Datastore, cs command.Service) Service {
	return &service{
		commands: cs,
		devices:  devices,
	}
}

type service struct {
	devices  device.Datastore
	commands command.Service
}

// Acknowledge a response from a device.
// NOTE: IOS devices do not always include the key `RequestType` in their response. Only the presence of the
// result key can be used to identify the response (or the command UUID)
func (svc service) Acknowledge(ctx context.Context, req mdm.Response) (int, error) {
	requestPayload, err := svc.commands.Find(req.CommandUUID)

	switch requestPayload.Command.RequestType {
	case "DeviceInformation":
		if err := svc.ackQueryResponses(req); err != nil {
			return 0, err
		}
	case "InstalledApplicationList":
		if err := svc.ackInstalledApplicationList(req); err != nil {
			return 0, err
		}
	case "CertificateList":
		if err := svc.ackCertificateList(req); err != nil {
			return 0, err
		}
	default:
		// Need to handle the absence of RequestType in IOS8 devices
		if req.QueryResponses.UDID != "" {
			if err := svc.ackQueryResponses(req); err != nil {
				return 0, err
			}
		}
		// Unhandled MDM client response
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

// Acknowledge Queries sent with DeviceInformation command
func (svc service) ackQueryResponses(req mdm.Response) error {
	devices, err := svc.devices.Devices(
		device.SerialNumber{SerialNumber: req.QueryResponses.SerialNumber},
		device.UDID{UDID: req.UDID},
	)

	if err != nil {
		return errors.Wrap(err, "ackQueryResponses fetching device")
	}

	if len(devices) == 0 {
		return errors.New("no enrolled device matches the one responding")
	}

	if len(devices) > 1 {
		return fmt.Errorf("expected a single device for udid: %s, serial number: %s, but got more than one.", req.UDID, req.QueryResponses.SerialNumber)
	}

	existing := devices[0]

	now := time.Now()
	existing.LastCheckin = &now
	existing.LastQueryResponse, err = json.Marshal(req.QueryResponses)

	if err != nil {
		return err
	}

	var serialNumber device.JsonNullString
	serialNumber.Scan(req.QueryResponses.SerialNumber)

	existing.ProductName = req.QueryResponses.ProductName
	existing.BuildVersion = req.QueryResponses.BuildVersion
	existing.DeviceName = req.QueryResponses.DeviceName
	existing.IMEI = req.QueryResponses.IMEI
	existing.MEID = req.QueryResponses.MEID
	existing.Model = req.QueryResponses.Model
	existing.OSVersion = req.QueryResponses.OSVersion
	existing.SerialNumber = serialNumber

	return svc.devices.Save("queryResponses", &existing)
}
