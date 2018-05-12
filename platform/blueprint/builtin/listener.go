package builtin

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/mdm/checkin"
	"github.com/micromdm/micromdm/mdm/mdm"
	"github.com/micromdm/micromdm/platform/blueprint"
	"github.com/micromdm/micromdm/platform/command"
	"github.com/micromdm/micromdm/platform/device"
	"github.com/micromdm/micromdm/platform/profile"
	"github.com/micromdm/micromdm/platform/pubsub"
)

func (db *DB) ApplyToDevice(ctx context.Context, svc command.Service, bp *blueprint.Blueprint, udid string) error {
	var requests []*mdm.CommandRequest
	for _, uuid := range bp.UserUUID {
		fmt.Println("Adding user to admin account")
		u, err := db.userDB.User(uuid)
		if err != nil {
			fmt.Printf("User UUID %s in Blueprint %s not added \n", bp.UserUUID, bp.Name)
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
						mdm.AdminAccount{
							ShortName:    u.UserShortname,
							FullName:     u.UserLongname,
							PasswordHash: u.PasswordHash,
							Hidden:       u.Hidden,
						},
					},
				},
			},
		})
	}

	for _, appURL := range bp.ApplicationURLs {
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

	for _, p := range bp.ProfileIdentifiers {
		foundProfile, err := db.profDB.ProfileById(p)
		if err != nil {
			if profile.IsNotFound(err) {
				fmt.Printf("Profile ID %s in Blueprint %s does not exist\n", p, bp.Name)
				continue
			}
			fmt.Println(err)
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
		_, err := svc.NewCommand(ctx, r)
		if err != nil {
			return errors.Wrap(err, "create new command from blueprint")
		}
	}
	return nil
}

func (db *DB) StartListener(sub pubsub.Subscriber, cmdSvc command.Service) error {
	tokenUpdateEvents, err := sub.Subscribe(context.TODO(), "applyAtEnroll", device.DeviceEnrolledTopic)
	if err != nil {
		return errors.Wrapf(err,
			"subscribing devices to %s topic", device.DeviceEnrolledTopic)
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
				bps, err := db.BlueprintsByApplyAt(blueprint.ApplyAtEnroll)
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
						Command: &mdm.Command{RequestType: "DeviceConfigured"},
						UDID:    ev.Command.UDID,
					})
					if err != nil {
						fmt.Println(errors.Wrapf(err, "sending DeviceConfigured"))
					}
				}

				// TODO: See notes from here:
				// https://github.com/jessepeterson/micromdm/blob/8b068ac98d06954bb3e08b1557c193007932552b/blueprint/listener.go#L73-L103
				// Also see discussion here for general direction:
				// https://github.com/micromdm/micromdm/pull/149
				// Finally see discussion here for high-level goals:
				// https://github.com/micromdm/micromdm/issues/110
			}

		}
	}()

	return nil
}

func intPtr(i int) *int {
	return &i
}
