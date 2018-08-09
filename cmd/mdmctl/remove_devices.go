package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

func (cmd *removeCommand) removeDevices(args []string) error {
	flagset := flag.NewFlagSet("remove-devices", flag.ExitOnError)
	var (
		flIdentifiers = flagset.String("udids", "", "comma separated list of device UDIDs")
	)
	flagset.Usage = usageFor(flagset, "mdmctl remove devices [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}

	if *flIdentifiers == "" {
		return errors.New("bad input: device UDID must be provided")
	}

	ctx := context.Background()
	err := cmd.devicesvc.RemoveDevices(ctx, strings.Split(*flIdentifiers, ","))
	if err != nil {
		return err
	}

	fmt.Printf("removed devices(s): %s\n", *flIdentifiers)

	return nil
}
