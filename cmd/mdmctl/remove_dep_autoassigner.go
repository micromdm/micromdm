package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/pkg/errors"
)

func (cmd *removeCommand) removeDEPAutoAssigner(args []string) error {
	flagset := flag.NewFlagSet("dep-autoassigner", flag.ExitOnError)
	var (
		flFilter = flagset.String("filter", "*", "filter string (only '*' supported right now)")
	)
	flagset.Usage = usageFor(flagset, "mdmctl remove dep-autoassigner [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}

	if *flFilter == "" {
		return errors.New("bad input: must provide -filter")
	}

	err := cmd.depsyncsvc.RemoveAutoAssigner(context.TODO(), *flFilter)
	if err != nil {
		return err
	}

	fmt.Printf("removed DEP profile associated with filter '%s'\n", *flFilter)
	return nil
}
