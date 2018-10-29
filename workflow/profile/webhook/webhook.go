package webhook

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/groob/plist"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/micromdm/micromdm/mdm"
	"github.com/micromdm/micromdm/pkg/httputil"
	"github.com/micromdm/micromdm/workflow/profile"
	"github.com/micromdm/micromdm/workflow/profile/device"
	"github.com/micromdm/micromdm/workflow/profile/inventory"
	"github.com/micromdm/micromdm/workflow/webhook"
)

type Server struct {
	db          DeviceStore
	profileDB   ProfileStore
	inventoryDB InventoryStore
	logger      log.Logger
	cabytes     []byte
}

func New(
	db DeviceStore,
	logger log.Logger,
	cabytes []byte,
	profileDB ProfileStore,
	inventoryDB InventoryStore,
) *Server {
	return &Server{
		db:          db,
		profileDB:   profileDB,
		inventoryDB: inventoryDB,
		logger:      logger,
		cabytes:     cabytes,
	}
}

type DeviceStore interface {
	Save(ctx context.Context, d device.Device) error
	DeviceByUDID(ctx context.Context, udid string) (device.Device, error)
}

type ProfileStore interface {
	Create(ctx context.Context, payload []byte) (profile.Profile, error)
}

type InventoryStore interface {
	UpdateFromListResponse(ctx context.Context, udid string, resp inventory.ListProfilesResponse) error
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var event webhook.Event
	if err := httputil.DecodeJSONRequest(r, &event); err != nil {
		level.Info(s.logger).Log("err", err)
		return
	}

	switch event.Topic {
	case mdm.TokenUpdateTopic, mdm.AuthenticateTopic:
		if err := s.ensureDevice(r.Context(), event.CheckinEvent.UDID); err != nil {
			level.Info(s.logger).Log("msg", "ensure device from checkin", "err", err)
			return
		}
	case mdm.ConnectTopic:
		if err := s.ensureDevice(r.Context(), event.AcknowledgeEvent.UDID); err != nil {
			level.Info(s.logger).Log("msg", "ensure device from connect", "err", err)
			return
		}
		if err := s.updateFromWebhook(r.Context(), event.AcknowledgeEvent); err != nil {
			level.Info(s.logger).Log("err", err)
			return
		}
	case mdm.CheckoutTopic:
		spew.Dump(event.CheckinEvent)
	}
}

func (s *Server) updateFromWebhook(ctx context.Context, event *webhook.AcknowledgeEvent) error {
	level.Debug(s.logger).Log(
		"msg", "handling ack event",
		"command_uuid", event.CommandUUID,
		"device_udid", event.UDID,
		"status", event.Status,
	)
	var resp inventory.ListProfilesResponse
	if err := plist.Unmarshal(event.RawPayload, &resp); err != nil {
		return errors.Wrap(err, "unmarshal raw payload plist")
	}
	switch {
	case len(resp.ProfileList) > 0:
		return s.inventoryDB.UpdateFromListResponse(ctx, event.UDID, resp)
	default:
		// not what we are looking for, just update last seen
	}
	return nil
}

func (s *Server) ensureDevice(ctx context.Context, udid string) error {
	_, err := s.db.DeviceByUDID(ctx, udid)
	switch {
	case err == nil:
		// update last seen
		return nil
	case isNotFound(err):
		level.Debug(s.logger).Log("msg", "creating new device", "udid", udid)
		err = s.newDevice(ctx, udid)
		return errors.Wrapf(err, "create new device for udid %s", udid)
	default:
		return errors.Wrapf(err, "ensuring device exists for udid %s", udid)
	}
}

func (s *Server) scepProfileForDevice(ctx context.Context, udid string) error {
	params := struct {
		UDID              string
		ProfileUUID       string
		SCEPPayloadUUID   string
		CACertPayloadUUID string
		CAPEM             string
	}{
		UDID:              udid,
		ProfileUUID:       uuid.NewV4().String(),
		SCEPPayloadUUID:   uuid.NewV4().String(),
		CACertPayloadUUID: uuid.NewV4().String(),
		CAPEM:             string(s.cabytes),
	}
	var buf bytes.Buffer
	if err := tmplStr.Execute(&buf, params); err != nil {
		return errors.Wrapf(err, "generate scep profile for udid %s", udid)
	}
	profile, err := s.profileDB.Create(ctx, buf.Bytes())
	if err != nil {
		return errors.Wrapf(err, "create scep profile for udid %s", udid)
	}
	level.Debug(s.logger).Log(
		"msg", "generated scep profile",
		"profile_uuid", params.ProfileUUID,
		"device_udid", params.UDID,
		"id", profile.ID,
	)
	return nil
}

func (s *Server) newDevice(ctx context.Context, udid string) error {
	now := time.Now().UTC()
	dev := device.Device{
		UDID:      udid,
		CreatedAt: now,
		UpdatedAt: now,
	}
	err := s.db.Save(ctx, dev)
	return errors.Wrapf(err, "save new device with udid %s", udid)
}

func isNotFound(err error) bool {
	err = errors.Cause(err)
	e, ok := errors.Cause(err).(interface {
		error
		NotFound() bool
	})
	return ok && e.NotFound()
}
