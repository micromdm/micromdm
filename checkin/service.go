package checkin

import (
	"time"

	"github.com/micromdm/mdm"
	"github.com/micromdm/micromdm/device"
	"github.com/micromdm/micromdm/management"
)

// Service defines methods for and MDM Checkin service
type Service interface {
	Authenticate(mdm.CheckinCommand) error
	TokenUpdate(mdm.CheckinCommand) error
	Checkout(mdm.CheckinCommand) error
}

// NewService creates a checkin service.
func NewService(devices device.Datastore, ms management.Service) Service {
	return &service{
		devices: devices,
		mgmt:    ms,
	}
}

type service struct {
	devices device.Datastore
	mgmt    management.Service
}

func (svc service) Authenticate(cmd mdm.CheckinCommand) error {
	var udid, serialNumber device.JsonNullString

	if err := udid.Scan(cmd.UDID); err != nil {
		return err
	}

	if err := serialNumber.Scan(cmd.SerialNumber); err != nil {
		return err
	}

	dev := &device.Device{
		UDID:         udid,
		SerialNumber: serialNumber,
		OSVersion:    cmd.OSVersion,
		BuildVersion: cmd.BuildVersion,
		ProductName:  cmd.ProductName,
		IMEI:         cmd.IMEI,
		MEID:         cmd.MEID,
		MDMTopic:     cmd.Topic,
		Model:        cmd.Model,
		DeviceName:   cmd.DeviceName,
		LastCheckin:  time.Now().UTC(),
	}

	_, err := svc.devices.New("authenticate", dev)
	return err
}

func (svc service) TokenUpdate(cmd mdm.CheckinCommand) error {
	if cmd.UserID != "" {
		// don't handle user updates for now
		return nil
	}
	token := cmd.Token.String()
	unlockToken := cmd.UnlockToken.String()
	existing, err := svc.devices.GetDeviceByUDID(cmd.UDID, []string{"device_uuid"}...)
	if err != nil {
		return err
	}
	existing.Token = token
	existing.MDMTopic = cmd.Topic
	existing.PushMagic = cmd.PushMagic
	existing.UnlockToken = unlockToken
	existing.AwaitingConfiguration = cmd.AwaitingConfiguration
	existing.Enrolled = true
	existing.LastCheckin = time.Now().UTC()

	err = svc.devices.Save("tokenUpdate", existing)
	if err != nil {
		return err
	}
	// trigger a push notification
	svc.mgmt.Push(cmd.UDID)
	return nil
}

func (svc service) Checkout(cmd mdm.CheckinCommand) error {
	existing, err := svc.devices.GetDeviceByUDID(cmd.UDID, []string{"device_uuid"}...)
	if err != nil {
		return err
	}
	existing.Enrolled = false
	err = svc.devices.Save("checkout", existing)
	if err != nil {
		return err
	}
	return nil
}
