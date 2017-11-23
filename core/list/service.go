package list

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/groob/plist"
	"github.com/micromdm/dep"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/appstore"
	"github.com/micromdm/micromdm/blueprint"
	"github.com/micromdm/micromdm/dep/deptoken"
	"github.com/micromdm/micromdm/device"
	"github.com/micromdm/micromdm/profile"
	"github.com/micromdm/micromdm/pubsub"
	"github.com/micromdm/micromdm/user"
)

type ListDevicesOption struct {
	Page    int
	PerPage int

	FilterSerial []string
	FilterUDID   []string
}

type ListUsersOption struct {
	Page    int
	PerPage int

	FilterUserID []string
	FilterUDID   []string
}

type GetBlueprintsOption struct {
	FilterName string
}

type GetProfilesOption struct {
	Identifier string `json:"id"`
}

type ListAppsOption struct {
	FilterName []string `json:"filter_name"`
}

type Service interface {
	ListDevices(ctx context.Context, opt ListDevicesOption) ([]DeviceDTO, error)
	ListUsers(ctx context.Context, opt ListUsersOption) ([]user.User, error)
	GetDEPTokens(ctx context.Context) ([]deptoken.DEPToken, []byte, error)
	GetBlueprints(ctx context.Context, opt GetBlueprintsOption) ([]blueprint.Blueprint, error)
	GetProfiles(ctx context.Context, opt GetProfilesOption) ([]profile.Profile, error)
	ListApplications(ctx context.Context, opt ListAppsOption) ([]AppDTO, error)
	DEPService
}

type ListService struct {
	mtx       sync.RWMutex
	DEPClient dep.Client

	Devices    *device.DB
	Blueprints *blueprint.DB
	Profiles   *profile.DB
	Tokens     *deptoken.DB
	Apps       appstore.AppStore
	Users      *user.DB
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
	tokenAdded, err := pubsub.Subscribe(context.TODO(), "list-token-events", deptoken.DEPTokenTopic)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event := <-tokenAdded:
				var token deptoken.DEPToken
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

func (svc *ListService) ListUsers(ctx context.Context, opts ListUsersOption) ([]user.User, error) {
	u, err := svc.Users.List()
	return u, errors.Wrap(err, "list users from api request")
}

func (svc *ListService) GetDEPTokens(ctx context.Context) ([]deptoken.DEPToken, []byte, error) {
	_, cert, err := svc.Tokens.DEPKeypair()
	if err != nil {
		return nil, nil, err
	}
	var certBytes []byte
	if cert != nil {
		certBytes = cert.Raw
	}

	tokens, err := svc.Tokens.DEPTokens()
	if err != nil {
		return nil, certBytes, err
	}

	return tokens, certBytes, nil
}

func (svc *ListService) GetBlueprints(ctx context.Context, opt GetBlueprintsOption) ([]blueprint.Blueprint, error) {
	if opt.FilterName != "" {
		bp, err := svc.Blueprints.BlueprintByName(opt.FilterName)
		if err != nil {
			return nil, err
		}
		return []blueprint.Blueprint{*bp}, err
	} else {
		bps, err := svc.Blueprints.List()
		if err != nil {
			return nil, err
		}
		return bps, nil
	}
}

func (svc *ListService) GetProfiles(ctx context.Context, opt GetProfilesOption) ([]profile.Profile, error) {
	if opt.Identifier != "" {
		foundProf, err := svc.Profiles.ProfileById(opt.Identifier)
		if err != nil {
			return nil, err
		}
		return []profile.Profile{*foundProf}, nil
	} else {
		return svc.Profiles.List()
	}
}
