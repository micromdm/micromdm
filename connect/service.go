package connect

import (
	"encoding/json"
	"fmt"
	"github.com/micromdm/mdm"
	"github.com/micromdm/micromdm/applications"
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
func NewService(devices device.Datastore, apps applications.Datastore, cs command.Service) Service {
	return &service{
		commands: cs,
		devices:  devices,
		apps:     apps,
	}
}

type service struct {
	devices  device.Datastore
	apps     applications.Datastore
	commands command.Service
}

func (svc service) Acknowledge(ctx context.Context, req mdm.Response) (int, error) {
	switch req.RequestType {
	case "DeviceInformation":
		if err := svc.ackQueryResponses(req); err != nil {
			return 0, err
		}
	default:
		// Need to handle the absence of RequestType in IOS8 devices
		if req.QueryResponses.UDID != "" {
			if err := svc.ackQueryResponses(req); err != nil {
				return 0, err
			}
		}
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
		return err
	}

	if len(devices) > 1 {
		return errors.New("expected a single query result for device, got more than one.")
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

// Acknowledge a response to `InstalledApplicationList`.
func (svc service) ackInstalledApplicationList(req mdm.Response) error {
	device, err := svc.devices.GetDeviceByUDID(req.UDID)
	if err != nil {
		return err
	}

	deviceApps, err := svc.apps.GetApplicationsByDeviceUUID(device.UUID)
	if err != nil {
		return err
	}

	// Any installed applications that are already represented in the applications datastore should be skipped.
	var updated []applications.Application = make([]applications.Application, len(req.InstalledApplicationList))

skip:
	for _, ackApp := range req.InstalledApplicationList {
		for _, app := range *deviceApps {
			if app.Name == ackApp.Name && app.Version == ackApp.Version {
				continue skip
			}
		}

		updated = append(updated, ackApp)
	}

	if len(updated) == 0 {
		return nil
	}

	// Determine applications which we have no record of at all, then insert them (find or create).
	for _, newApp := range updated {
		existing, err := svc.apps.Applications(applications.Name{newApp.Name}, applications.Version{newApp.Version})
		if err != nil {
			return err
		}

		switch {
		case len(existing) > 1:
			return fmt.Errorf("expected a single application match for application name: %s, got %d results", newApp.Name, len(existing))
		case len(existing) == 0: // No record exists and therefore both the application row and device association must be created.
			appUuid, err := svc.apps.New(newApp)
			if err != nil {
				return err
			}

			newApp.UUID = appUuid
		}

		// For both len(existing) == 0 and len(existing) == 1, the row must be inserted for devices_applications.
		if err := svc.apps.SaveApplicationByDeviceUUID(device.UUID, newApp.UUID); err != nil {
			return err
		}
	}

	return nil
}
