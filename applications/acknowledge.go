package applications

import (
	"fmt"
	"github.com/micromdm/mdm"
)

//
//func (svc service) ackInstalledApplicationList(req mdm.Response) error {
//	apps := req.InstalledApplicationList
//	devices, err := svc.devices.Devices(
//		device.UDID{UDID: req.UDID},
//	)
//	if err != nil {
//		return err
//	}
//
//	if len(devices) > 1 || len(devices) == 0 {
//		return errors.New("expected a single query result for device, got more or less than one.")
//	}
//	device := devices[0]
//
//	for _, app := range apps {
//
//	}
//}

func AcknowledgeInstalledApplicationListResponse(response mdm.Response, datastores map[string]interface{}) error {
	fmt.Println("InstalledApplicationListResponseHandler TODO")
	return nil
}
