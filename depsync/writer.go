package depsync

import (
	"database/sql"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/micromdm/dep"
	"github.com/micromdm/micromdm/device"
)

type Writer interface {
	Start(deviceChan <-chan dep.Device)
	Write(*dep.Device) (string, error)
}

type writer struct {
	datastore device.Datastore
	logger    log.Logger
	done      <-chan struct{}
}

func NewWriter(datastore device.Datastore, logger log.Logger, done <-chan struct{}) Writer {
	return &writer{
		datastore: datastore,
		logger:    logger,
		done:      done,
	}
}

func (w *writer) Start(DEPDevices <-chan dep.Device) {
	var dev dep.Device

	select {
	case dev = <-DEPDevices:
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
		default:
			w.logger.Log("level", "debug", "msg", "writing fetched device to database")
			deviceUUID, err := w.Write(&dev)
			if err != nil {
				w.logger.Log("level", "error", "msg", fmt.Sprintf("Failed to write DEP device: %s", err))
			}
			w.logger.Log("level", "debug", "msg", fmt.Sprintf("wrote device with UUID: %s", deviceUUID))
		}
	case <-w.done:
		w.logger.Log("level", "info", "stopping DEP device writer")
		return
	}
}

func (w *writer) Write(dev *dep.Device) (string, error) {
	dvc := &device.Device{
		SerialNumber: device.JsonNullString{sql.NullString{dev.SerialNumber, true}},
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
