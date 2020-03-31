package blueprint

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"

	mdmsvc "github.com/micromdm/micromdm/mdm"
	"github.com/micromdm/micromdm/mdm/mdm"
	"github.com/micromdm/micromdm/platform/command"
	"github.com/micromdm/micromdm/platform/device"
	"github.com/micromdm/micromdm/platform/profile"
	"github.com/micromdm/micromdm/platform/pubsub"
	"github.com/micromdm/micromdm/platform/user"
)

type BlueprintWorkerStore interface {
	BlueprintsByApplyAt(ctx context.Context, action string) ([]Blueprint, error)
}

type UserStore interface {
	User(ctx context.Context, uuid string) (*user.User, error)
}

type ProfileStore interface {
	ProfileById(ctx context.Context, id string) (*profile.Profile, error)
}

func NewWorker(
	db BlueprintWorkerStore,
	userDB UserStore,
	profileDB ProfileStore,
	cmdsvc command.Service,
	sub pubsub.Subscriber,
	logger log.Logger,
) *Worker {
	return &Worker{
		db:        db,
		userDB:    userDB,
		profileDB: profileDB,
		ps:        sub,
		cmdsvc:    cmdsvc,
		logger:    logger,
	}
}

type Worker struct {
	db        BlueprintWorkerStore
	userDB    UserStore
	profileDB ProfileStore
	ps        pubsub.Subscriber
	cmdsvc    command.Service
	logger    log.Logger
}

func (w *Worker) Run(ctx context.Context) error {
	tokenUpdateEvents, err := w.ps.Subscribe(ctx, "applyAtEnroll", device.DeviceEnrolledTopic)
	if err != nil {
		return errors.Wrapf(err, "subscribing devices to %s topic", device.DeviceEnrolledTopic)
	}

	for {
		var err error
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ev := <-tokenUpdateEvents:
			err = w.handleTokenUpdateEvent(ctx, ev.Message)
		}

		if err != nil {
			level.Info(w.logger).Log(
				"msg", "handle blueprint action",
				"err", err,
			)
			continue
		}

	}
}

// TODO: See notes from here:
// https://github.com/jessepeterson/micromdm/blob/8b068ac98d06954bb3e08b1557c193007932552b/blueprint/listener.go#L73-L103
// Also see discussion here for general direction:
// https://github.com/micromdm/micromdm/pull/149
// Finally see discussion here for high-level goals:
// https://github.com/micromdm/micromdm/issues/110
func (w *Worker) handleTokenUpdateEvent(ctx context.Context, message []byte) error {
	var ev mdmsvc.CheckinEvent
	if err := mdmsvc.UnmarshalCheckinEvent(message, &ev); err != nil {
		return errors.Wrap(err, "unmarshal checkin event")
	}
	if ev.Command.UserID != "" {
		level.Debug(w.logger).Log(
			"msg", "skipping user token update in blueprint worker.",
			"device_udid", ev.Command.UDID,
			"user_id", ev.Command.UserID,
		)
		// skip UserID token updates
		return nil
	}

	bps, err := w.db.BlueprintsByApplyAt(ctx, ApplyAtEnroll)
	if err != nil {
		return errors.Wrap(err, "get blueprints by ApplyAtEnroll")
	}

	// if there are no blueprints exit early. This will ensure that DeviceConfigured is not sent.
	if len(bps) == 0 {
		level.Debug(w.logger).Log(
			"msg", "no blueprints to apply",
			"device_udid", ev.Command.UDID,
		)
		return nil
	}

	for _, bp := range bps {
		level.Debug(w.logger).Log(
			"msg", "applying blueprint",
			"device_udid", ev.Command.UDID,
			"blueprint_name", bp.Name,
		)

		if err := w.applyToDevice(ctx, bp, ev.Command.UDID); err != nil {
			return errors.Wrapf(err, "apply blueprint to udid name=%s, udid=%s", bp.Name, ev.Command.UDID)
		}
	}

	if ev.Command.AwaitingConfiguration {
		level.Debug(w.logger).Log(
			"msg", "sending DeviceConfigured at the end of blueprint",
			"device_udid", ev.Command.UDID,
		)
		_, err := w.cmdsvc.NewCommand(ctx, &mdm.CommandRequest{
			Command: &mdm.Command{RequestType: "DeviceConfigured"},
			UDID:    ev.Command.UDID,
		})
		if err != nil {
			return errors.Wrap(err, "send DeviceConfigured")
		}
	}

	return nil

}

func (w *Worker) applyToDevice(ctx context.Context, bp Blueprint, udid string) error {
	var requests []*mdm.CommandRequest
	for _, uuid := range bp.UserUUID {
		level.Debug(w.logger).Log(
			"msg", "creating mdm command request from blueprint",
			"request_type", "AccountConfiguration",
			"blueprint_name", bp.Name,
			"user_uuid", uuid,
			"device_udid", udid,
		)

		usr, err := w.userDB.User(ctx, uuid)
		if err != nil {
			level.Info(w.logger).Log(
				"msg", "get user for AccountConfiguration request",
				"blueprint_name", bp.Name,
				"user_uuid", uuid,
				"device_udid", udid,
				"err", err,
			)
			continue
		}

		requests = append(requests, &mdm.CommandRequest{
			UDID: udid,
			Command: &mdm.Command{
				RequestType: "AccountConfiguration",
				AccountConfiguration: &mdm.AccountConfiguration{
					SkipPrimarySetupAccountCreation:     bp.SkipPrimarySetupAccountCreation,
					SetPrimarySetupAccountAsRegularUser: bp.SetPrimarySetupAccountAsRegularUser,
					AutoSetupAdminAccounts: []mdm.AdminAccount{
						{
							ShortName:    usr.UserShortname,
							FullName:     usr.UserLongname,
							PasswordHash: usr.PasswordHash,
							Hidden:       usr.Hidden,
						},
					},
				},
			},
		})
	}

	for _, appURL := range bp.ApplicationURLs {
		level.Debug(w.logger).Log(
			"msg", "creating mdm command request from blueprint",
			"request_type", "InstallApplication",
			"blueprint_name", bp.Name,
			"manifest_url", appURL,
			"device_udid", udid,
		)

		requests = append(requests, &mdm.CommandRequest{
			UDID: udid,
			Command: &mdm.Command{
				RequestType: "InstallApplication",
				InstallApplication: &mdm.InstallApplication{
					ManifestURL:     &appURL,
					ManagementFlags: intPtr(1),
				},
			},
		})
	}

	for _, pid := range bp.ProfileIdentifiers {
		level.Debug(w.logger).Log(
			"msg", "creating mdm command request from blueprint",
			"request_type", "InstallProfile",
			"blueprint_name", bp.Name,
			"profile_identifier", pid,
			"device_udid", udid,
		)
		foundProfile, err := w.profileDB.ProfileById(ctx, pid)
		if err != nil {
			level.Info(w.logger).Log(
				"msg", "retrieve profile from db",
				"blueprint_name", bp.Name,
				"profile_identifier", pid,
				"is_not_found_err", profile.IsNotFound(err),
				"err", err,
			)
			continue
		}

		requests = append(requests, &mdm.CommandRequest{
			UDID: udid,
			Command: &mdm.Command{
				RequestType: "InstallProfile",
				InstallProfile: &mdm.InstallProfile{
					Payload: foundProfile.Mobileconfig,
				},
			},
		})
	}

	for _, r := range requests {
		if _, err := w.cmdsvc.NewCommand(ctx, r); err != nil {
			return errors.Wrap(err, "create new command from blueprint")
		}
	}
	return nil
}

func intPtr(i int) *int {
	return &i
}
