package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/micromdm/micromdm/platform/dep/sync"
	"github.com/pkg/errors"
)

func (cmd *applyCommand) applyDEPAutoAssigner(args []string) error {
	flagset := flag.NewFlagSet("dep-autoassigner", flag.ExitOnError)
	var (
		flFilter      = flagset.String("filter", "*", "filter string (only '*' supported right now)")
		flProfileUUID = flagset.String("uuid", "", "DEP profile UUID to set")
	)
	flagset.Usage = usageFor(flagset, "mdmctl apply dep-autoassigner [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}

	if *flFilter == "" || *flProfileUUID == "" {
		return errors.New("bad input: must provide both -filter and -uuid")
	}

	assigner := sync.AutoAssigner{Filter: *flFilter, ProfileUUID: *flProfileUUID}

	err := cmd.depsyncsvc.ApplyAutoAssigner(context.TODO(), &assigner)
	if err != nil {
		return err
	}

	fmt.Printf("saved auto-assign filter '%s' to DEP profile UUID '%s'\n", assigner.Filter, assigner.ProfileUUID)
	fmt.Println("newly added DEP devices will be auto-assigned to the above profile UUID")
	return nil
}
