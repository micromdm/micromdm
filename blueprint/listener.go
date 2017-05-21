package blueprint

import (
	"context"
	"fmt"

	"github.com/micromdm/mdm"
	"github.com/micromdm/micromdm/checkin"
	"github.com/micromdm/micromdm/command"
	"github.com/micromdm/micromdm/pubsub"
	"github.com/pkg/errors"
)

func ApplyToDevice(svc command.Service, blueprint *Blueprint, udid string) error {
	ctx := context.Background()
	var requests []*mdm.CommandRequest
	for _, appURL := range blueprint.ApplicationURLs {
		requests = append(requests, &mdm.CommandRequest{
			RequestType: "InstallApplication",
			UDID:        udid,
			InstallApplication: mdm.InstallApplication{
				ManifestURL:     appURL,
				ManagementFlags: 1,
			},
		})
	}

	for _, profilePayload := range blueprint.Profiles {
		requests = append(requests, &mdm.CommandRequest{
			RequestType: "InstallProfile",
			UDID:        udid,
			InstallProfile: mdm.InstallProfile{
				Payload: profilePayload,
			},
		})
	}

	for _, r := range requests {
		_, err := svc.NewCommand(ctx, r)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) StartListener(sub pubsub.Subscriber, cmdSvc command.Service) error {
	tokUpdEvents, err := sub.Subscribe("hardcode-dep", checkin.TokenUpdateTopic)
	if err != nil {
		return errors.Wrapf(err,
			"subscribing devices to %s topic", checkin.TokenUpdateTopic)
	}

	go func() {
		for {
			select {
			case event := <-tokUpdEvents:
				var ev checkin.Event
				if err := checkin.UnmarshalEvent(event.Message, &ev); err != nil {
					fmt.Println(err)
					continue
				}
				bps, err := db.BlueprintsByWhenTag("TokenUpdate")
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("found blueprints:", len(bps))
				for _, bp := range bps {
					ApplyToDevice(cmdSvc, bp, ev.Command.UDID)
				}

				// TODO: send DeviceConfigured if ev.Command.AwaitingConfiguration == true

				// TODO: Perhaps store a counter of TokenUpdates performed on
				// the client in this way we can track if this is the first
				// TokenUpdate (we can kick off DEP-initiated packages) or not
				// (plain enrollment). I.e. a "FirstTokenUpdate" message as
				// opposed to the periodic Token checkins a device will do
				// over the course of its enrolled life. Alternately can
				// create a unique pubsub event just for enrolled status.

				// TODO: revisit APNs & command queuing. Right now every
				// successfully queued command kicks off an APNs notification
				// when really we only need a single one after the batch of
				// commands are queued

				// TODO: investigate sometimes the first DEP post-user account
				// setup MDM check-in is from the *user* but this is the only
				// way to know for sending post-DEP MDM commands for older
				// OSes (10.10? 10.11?)
			}

			// TODO: to implement many more "When"-type qualifiers
			// Envision things like "KeepState" where perhaps periodic
			// checks for profiles and apps are done.

			// TODO: to implement a predicate analysis lanauge for
			// filtering events, devices, etc. See:
			// Yaml examples:
			//   https://github.com/micromdm/micromdm/issues/110
			// Golang predicate langauge research:
			//   https://macadmins.slack.com/archives/C19RTE0L9/p1493023727619223

		}
	}()

	return nil
}
