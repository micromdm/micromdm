package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/platform/device"
)

func (cmd *removeCommand) removeDevices(args []string) error {
	flagset := flag.NewFlagSet("remove-devices", flag.ExitOnError)
	var (
		flIdentifiers = flagset.String("udids", "", "comma separated list of device UDIDs")
		flSerials     = flagset.String("serials", "", "comma separated list of device serials")
	)
	flagset.Usage = usageFor(flagset, "mdmctl remove devices [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}

	if *flIdentifiers == "" && *flSerials == "" {
		return errors.New("bad input: device UDID or Serial must be provided")
	}

	opts := device.RemoveDevicesOptions{}
	if *flIdentifiers != "" {
		opts.UDIDs = strings.Split(*flIdentifiers, ",")
	}
	if *flSerials != "" {
		opts.Serials = strings.Split(*flSerials, ",")
	}

	ctx := context.Background()
	err := cmd.devicesvc.RemoveDevices(ctx, opts)
	if err != nil {
		return err
	}

	fmt.Printf("removed devices(s): %s\n", strings.Join(append(opts.UDIDs, opts.Serials...), ", "))

	return nil
}
