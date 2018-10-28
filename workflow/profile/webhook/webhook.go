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
	db            DeviceStore
	profileDB     ProfileStore
	inventoryDB   InventoryStore
	logger        log.Logger
	profileClient ProfileClient
	cabytes       []byte
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

type ProfileClient interface {
	Send(ctx context.Context, profile []byte) (uuid string, err error)
}

type InventoryStore interface {
	UpdateFromListResponse(ctx context.Context, udid string, resp inventory.ListProfilesResponse) error
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var event webhook.Event
	if err := httputil.DecodeJSONRequest(r, &event); err != nil {
		panic(err)
	}

	switch event.Topic {
	case mdm.TokenUpdateTopic, mdm.AuthenticateTopic:
		s.ensureDevice(r.Context(), event.CheckinEvent.UDID)
	case mdm.ConnectTopic:
		if err := s.updateFromWebhook(r.Context(), event.AcknowledgeEvent); err != nil {
			level.Info(s.logger).Log(
				"err", err,
			)
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

func (s *Server) ensureDevice(ctx context.Context, udid string) {
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
		panic(err)
	}
	profile, err := s.profileDB.Create(ctx, buf.Bytes())
	if err != nil {
		panic(err)
	}
	level.Debug(s.logger).Log(
		"msg", "generated scep profile",
		"profile_uuid", params.ProfileUUID,
		"device_udid", params.UDID,
		"id", profile.ID,
	)
	_, err = s.db.DeviceByUDID(ctx, udid)
	switch {
	case err == nil:
		// update last seen
		return
	case isNotFound(err):
		if err := s.newDevice(ctx, udid); err != nil {
			level.Info(s.logger).Log(
				"msg", "creating new device ",
				"err", err,
			)
		}
		return
	default:
		level.Info(s.logger).Log(
			"msg", "looking up profile device by udid",
			"err", err,
		)
		return
	}
}

func (s *Server) newDevice(ctx context.Context, udid string) error {
	now := time.Now().UTC()
	dev := device.Device{
		UDID:      udid,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.db.Save(ctx, dev); err != nil {
		return errors.Wrapf(err, "save new device with udid %s", udid)
	}
	return nil
}

func isNotFound(err error) bool {
	err = errors.Cause(err)
	e, ok := errors.Cause(err).(interface {
		error
		NotFound() bool
	})
	return ok && e.NotFound()
}
