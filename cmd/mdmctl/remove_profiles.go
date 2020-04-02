package main

import (
	"context"
	"flag"
	"fmt"
	"strings"
)

func (cmd *removeCommand) removeProfiles(args []string) error {
	flagset := flag.NewFlagSet("remove-profiles", flag.ExitOnError)
	var (
		flIdentifier = flagset.String("id", "", "profile identifier, optionally comma-separated")
	)
	flagset.Usage = usageFor(flagset, "mdmctl remove profiles [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}

	ctx := context.Background()
	err := cmd.profilesvc.RemoveProfiles(ctx, strings.Split(*flIdentifier, ","))
	if err != nil {
		return err
	}

	fmt.Printf("removed profile(s): %s\n", *flIdentifier)

	return nil
}
