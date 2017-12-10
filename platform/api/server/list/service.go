package list

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/groob/plist"
	"github.com/micromdm/dep"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/platform/appstore"
	"github.com/micromdm/micromdm/platform/config"
	"github.com/micromdm/micromdm/platform/device"
	"github.com/micromdm/micromdm/platform/pubsub"
)

type ListDevicesOption struct {
	Page    int
	PerPage int

	FilterSerial []string
	FilterUDID   []string
}

type ListAppsOption struct {
	FilterName []string `json:"filter_name"`
}

type Service interface {
	ListDevices(ctx context.Context, opt ListDevicesOption) ([]DeviceDTO, error)
	ListApplications(ctx context.Context, opt ListAppsOption) ([]AppDTO, error)
	DEPService
}

type ListService struct {
	mtx       sync.RWMutex
	DEPClient dep.Client

	Devices *device.DB
	Apps    appstore.AppStore
}

func (svc *ListService) ListApplications(ctx context.Context, opts ListAppsOption) ([]AppDTO, error) {
	var filter string
	if len(opts.FilterName) == 1 {
		filter = opts.FilterName[0]
	}
	apps, err := svc.Apps.Apps(filter)
	if err != nil {
		return nil, err
	}
	var appList []AppDTO
	for name, app := range apps {
		payload, err := plist.MarshalIndent(&app, "  ")
		if err != nil {
			return nil, errors.Wrap(err, "create dto payload")
		}
		appList = append(appList, AppDTO{
			Name:    name,
			Payload: payload,
		})
	}
	return appList, nil
}

func (svc *ListService) WatchTokenUpdates(pubsub pubsub.Subscriber) error {
	tokenAdded, err := pubsub.Subscribe(context.TODO(), "list-token-events", config.DEPTokenTopic)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event := <-tokenAdded:
				var token config.DEPToken
				if err := json.Unmarshal(event.Message, &token); err != nil {
					log.Printf("unmarshalling tokenAdded to token: %s\n", err)
					continue
				}

				client, err := token.Client()
				if err != nil {
					log.Printf("creating new DEP client: %s\n", err)
					continue
				}

				svc.mtx.Lock()
				svc.DEPClient = client
				svc.mtx.Unlock()
			}
		}
	}()

	return nil
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
