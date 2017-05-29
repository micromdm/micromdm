package list

import (
	"context"

	"github.com/micromdm/dep"
	"github.com/micromdm/micromdm/blueprint"
	"github.com/micromdm/micromdm/deptoken"
	"github.com/micromdm/micromdm/device"
	"github.com/micromdm/micromdm/profile"
)

type ListDevicesOption struct {
	Page    int
	PerPage int

	FilterSerial []string
	FilterUDID   []string
}

type GetBlueprintsOption struct {
	FilterName string
}

type GetProfilesOption struct {
	Identifier string `json:"id"`
}

type Service interface {
	ListDevices(ctx context.Context, opt ListDevicesOption) ([]DeviceDTO, error)
	GetDEPTokens(ctx context.Context) ([]deptoken.DEPToken, []byte, error)
	GetBlueprints(ctx context.Context, opt GetBlueprintsOption) ([]blueprint.Blueprint, error)
	GetProfiles(ctx context.Context, opt GetProfilesOption) ([]profile.Profile, error)
	DEPService
}

type ListService struct {
	DEPClient  dep.Client
	Devices    *device.DB
	Blueprints *blueprint.DB
	Profiles   *profile.DB
	Tokens     *deptoken.DB
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
