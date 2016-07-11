package device

import (
	"encoding/json"
	"errors"
	"github.com/micromdm/mdm"
	"time"
)

// Acknowledge Queries sent with DeviceInformation command
func AcknowledgeDeviceInformationResponse(req mdm.Response, datastores map[string]interface{}) error {
	store, found := datastores["devices"]
	if !found {
		return errors.New("Do not have access to datastore for saving device information")
	}

	devicesStore, ok := store.(Datastore)
	if !ok {
		return errors.New("could not acknowledge device information because the given datastore isnt a device datastore.")
	}

	devices, err := devicesStore.Devices(
		SerialNumber{SerialNumber: req.QueryResponses.SerialNumber},
		UDID{UDID: req.UDID},
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

	var serialNumber JsonNullString
	serialNumber.Scan(req.QueryResponses.SerialNumber)

	existing.ProductName = req.QueryResponses.ProductName
	existing.BuildVersion = req.QueryResponses.BuildVersion
	existing.DeviceName = req.QueryResponses.DeviceName
	existing.IMEI = req.QueryResponses.IMEI
	existing.MEID = req.QueryResponses.MEID
	existing.Model = req.QueryResponses.Model
	existing.OSVersion = req.QueryResponses.OSVersion
	existing.SerialNumber = serialNumber

	return devicesStore.Save("queryResponses", &existing)
}
