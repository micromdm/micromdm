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
		flIdentifier = flagset.String("udid", "", "device UDID, optionally comma separated")
	)
	flagset.Usage = usageFor(flagset, "mdmctl remove devices [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}

	if *flIdentifier == "" {
		return errors.New("bad input: device UDID must be provided")
	}

	ctx := context.Background()
	err := cmd.devicesvc.RemoveDevices(ctx, strings.Split(*flIdentifier, ","))
	if err != nil {
		return err
	}

	fmt.Printf("removed devices(s): %s\n", *flIdentifier)

	return nil
}
