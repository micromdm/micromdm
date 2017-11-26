package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/pkg/errors"
)

func (cmd *removeCommand) removeBlock(args []string) error {
	flagset := flag.NewFlagSet("unblock", flag.ExitOnError)
	var (
		flUDID = flagset.String("udid", "", "UDID of device to unblock")
	)
	flagset.Usage = usageFor(flagset, "mdmctl remove block [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}

	if *flUDID == "" {
		flagset.Usage()
		return errors.New("bad input: must provide a device UDID to unblock.")
	}

	ctx := context.Background()
	if err := cmd.remove.UnblockDevice(ctx, *flUDID); err != nil {
		return err
	}

	fmt.Println("success")

	return nil
}
