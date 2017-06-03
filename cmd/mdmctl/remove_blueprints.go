package main

import (
	"context"
	"flag"
	"fmt"
	"strings"
)

func (cmd *removeCommand) removeBlueprints(args []string) error {
	flagset := flag.NewFlagSet("remove-blueprints", flag.ExitOnError)
	var (
		flBlueprintName = flagset.String("name", "", "name of blueprint, optionally comma separated")
	)
	flagset.Usage = usageFor(flagset, "mdmctl remove blueprints [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}

	ctx := context.Background()
	err := cmd.remove.RemoveBlueprints(ctx, strings.Split(*flBlueprintName, ","))
	if err != nil {
		return err
	}

	fmt.Printf("removed blueprint(s): %s\n", *flBlueprintName)

	return nil
}
