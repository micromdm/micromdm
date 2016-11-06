package dep

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/micromdm/dep"
	"github.com/micromdm/micromdm/device"
)

type Writer interface {
	Start(deviceChan <-chan dep.Device)
	Write(*dep.Device) error
}

type writer struct {
	datastore device.Datastore
	logger    log.Logger
}

func NewWriter(datastore device.Datastore, logger log.Logger) Writer {
	return &writer{
		datastore: datastore,
		logger:    logger,
	}
}

func (w *writer) Start(deviceChan <-chan dep.Device) {
	for dev := range deviceChan {
		//var opType string = dev.OpType
		//var opDate time.Time = dev.OpDate

		switch dev.OpType {
		case "added":
			w.logger.Log("level", "debug", "msg", "writing added device to database")
			deviceUUID, err := w.Write(&dev)
			if err != nil {
				w.logger.Log("level", "error", "msg", fmt.Sprintf("Failed to write DEP device: %s", err))
			}
			w.logger.Log("level", "debug", "msg", fmt.Sprintf("wrote device with UUID: %s", deviceUUID))
		case "modified":
		case "deleted":
			w.logger.Log("level", "debug", "msg", "removing dep device from database")
			if err := w.Delete(&dev); err != nil {
				w.logger.Log("level", "error", "msg", fmt.Sprintf("Failed to delete DEP device: %s", err))
			}
		default:
			w.logger.Log("level", "debug", "msg", "writing fetched device to database")
			deviceUUID, err := w.Write(&dev)
			if err != nil {
				w.logger.Log("level", "error", "msg", fmt.Sprintf("Failed to write DEP device: %s", err))
			}
			w.logger.Log("level", "debug", "msg", fmt.Sprintf("wrote device with UUID: %s", deviceUUID))
		}

	}
}

func (w *writer) Write(dev *dep.Device) (string, error) {
	dvc := &device.Device{
		SerialNumber: dev.SerialNumber,
		Model:        dev.Model,
		Description:  dev.Description,
		Color:        dev.Color,
		AssetTag:     dev.AssetTag,
	}

	deviceUUID, err := w.datastore.New("fetch", dvc)
	if err != nil {
		return "", err
	}

	return deviceUUID, nil
}

func (w *writer) Delete(dev *dep.Device) error {

}
