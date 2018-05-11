package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
)

func (cmd *getCommand) getDEPAutoAssigners(args []string) error {
	flagset := flag.NewFlagSet("dep-autoassigner", flag.ExitOnError)
	flagset.Usage = usageFor(flagset, "mdmctl get dep-autoassigner [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}

	assigners, err := cmd.depsyncsvc.GetAutoAssigners(context.TODO())
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintf(w, "Filter\tDEP Profile UUID\n")
	for _, a := range assigners {
		fmt.Fprintf(w, "%s\t%s\n", a.Filter, a.ProfileUUID)
	}
	w.Flush()

	return nil
}
