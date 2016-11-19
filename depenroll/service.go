package depenroll

import (
	"errors"
	"fmt"

	"github.com/micromdm/mdm"
	"github.com/micromdm/micromdm/command"
	"github.com/micromdm/micromdm/device"
)

type Service interface {
	// EnrollDEP returns an enrollment profile during DEP Enrollment.
	EnrollDEP(udid, serial string) ([]byte, error)
}

func NewService(devices device.Datastore, commands command.Service, profile []byte) Service {
	return &service{
		devices:  devices,
		commands: commands,
		profile:  profile,
	}
}

type service struct {
	devices  device.Datastore
	commands command.Service
	profile  []byte
}

func (svc service) EnrollDEP(udid, serial string) ([]byte, error) {
	err := svc.initialSetup(udid, serial)
	if err != nil {
		// TODO: stop ignoring the error there
		fmt.Println(err)
	}
	return svc.profile, nil
}

func (svc service) initialSetup(deviceUDID, serial string) error {
	devs, err := svc.devices.Devices(device.SerialNumber{SerialNumber: serial})
	if err != nil {
		return err
	}
	if len(devs) == 0 {
		return errors.New("device not found")
	}
	dev := devs[0]
	if dev.Workflow == "" {
		// no workflow, send DeviceConfigured
		return svc.sendConfigured(deviceUDID)
	}
	return nil
}

func (svc service) sendConfigured(deviceUDID string) error {
	cmdRequest := &mdm.CommandRequest{
		UDID:        deviceUDID,
		RequestType: "DeviceConfigured",
	}
	_, err := svc.commands.NewCommand(cmdRequest)
	if err != nil {
		return err
	}
	return nil
}
