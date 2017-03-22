package list

import (
	"context"

	"github.com/micromdm/micromdm/device"
)

type ListDevicesOption struct {
	Page    int
	PerPage int

	FilterSerial []string
	FilterUDID   []string
}

type Service interface {
	ListDevices(ctx context.Context, opt ListDevicesOption) ([]DeviceDTO, error)
}

type ListService struct {
	Devices *device.DB
}

func (svc *ListService) ListDevices(ctx context.Context, opt ListDevicesOption) ([]DeviceDTO, error) {
	devices, err := svc.Devices.List()
	dto := []DeviceDTO{}
	for _, d := range devices {
		dto = append(dto, DeviceDTO{
			SerialNumber:     d.SerialNumber,
			UDID:             d.UDID,
			EnrollmentStatus: d.Enrolled,
			LastSeen:         d.LastCheckin,
		})
	}
	return dto, err
}
