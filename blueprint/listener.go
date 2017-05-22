package blueprint

import (
	"context"
	"fmt"

	"github.com/micromdm/mdm"
	"github.com/micromdm/micromdm/checkin"
	"github.com/micromdm/micromdm/command"
	"github.com/micromdm/micromdm/profile"
	"github.com/micromdm/micromdm/pubsub"
	"github.com/pkg/errors"
)

func (db *DB) ApplyToDevice(ctx context.Context, svc command.Service, bp *Blueprint, udid string) error {

	var requests []*mdm.CommandRequest
	for _, appURL := range bp.ApplicationURLs {
		requests = append(requests, &mdm.CommandRequest{
			RequestType: "InstallApplication",
			UDID:        udid,
			InstallApplication: mdm.InstallApplication{
				ManifestURL:     appURL,
				ManagementFlags: 1,
			},
		})
	}

	for _, p := range bp.ProfileIdentifiers {
		foundProfile, err := db.profDB.ProfileById(p)
		if err != nil {
			if profile.IsNotFound(err) {
				fmt.Printf("Profile ID %s in Blueprint %s does not exist\n", p, bp.Name)
				continue
			}
			fmt.Println(errors.Wrap(err, "fetching profile"))
			continue
		}

		requests = append(requests, &mdm.CommandRequest{
			RequestType: "InstallProfile",
			UDID:        udid,
			InstallProfile: mdm.InstallProfile{
				Payload: foundProfile.Mobileconfig,
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
	tokenUpdateEvents, err := sub.Subscribe("applyAtEnroll", checkin.TokenUpdateTopic)
	if err != nil {
		return errors.Wrapf(err,
			"subscribing devices to %s topic", checkin.TokenUpdateTopic)
	}

	go func() {
		for {
			select {
			case event := <-tokenUpdateEvents:
				var ev checkin.Event
				if err := checkin.UnmarshalEvent(event.Message, &ev); err != nil {
					fmt.Println(err)
					continue
				}
				if ev.Command.UserID != "" {
					// skip UserID token updates
					continue
				}
				fmt.Println("UserID", ev.Command.UserID)
				bps, err := db.BlueprintsByApplyAt("enroll")
				if err != nil {
					fmt.Println(err)
					continue
				}
				ctx := context.Background()
				for _, bp := range bps {
					fmt.Printf("applying blueprint %s to %s\n", bp.Name, ev.Command.UDID)
					err := db.ApplyToDevice(ctx, cmdSvc, bp, ev.Command.UDID)
					if err != nil {
						fmt.Println(err)
					}
				}

				if ev.Command.AwaitingConfiguration {
					_, err := cmdSvc.NewCommand(ctx, &mdm.CommandRequest{
						RequestType: "DeviceConfigured",
						UDID:        ev.Command.UDID,
					})
					if err != nil {
						fmt.Println(errors.Wrapf(err, "sending DeviceConfigured"))
					}
				}

				// TODO: See notes from here:
				// https://github.com/jessepeterson/micromdm/blob/8b068ac98d06954bb3e08b1557c193007932552b/blueprint/listener.go#L73-L103
			}

		}
	}()

	return nil
}
