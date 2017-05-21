package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

type depProfilesTableOutput struct{ w *tabwriter.Writer }

func (out *depProfilesTableOutput) BasicHeader() {
	fmt.Fprintf(out.w, "Name\tMandatory\tRemovable\tAwaitConfigured\tSkippedItems\n")
}

func (out *depProfilesTableOutput) BasicFooter() {
	out.w.Flush()
}

func (cmd *getCommand) getDEPProfiles(args []string) error {
	flagset := flag.NewFlagSet("dep-profiles", flag.ExitOnError)
	flUUID := flagset.String("uuid", "", "DEP Profile UUID")
	flagset.Usage = usageFor(flagset, "mdmctl get dep-profiles [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}
	w := tabwriter.NewWriter(os.Stderr, 0, 4, 2, ' ', 0)
	out := &depProfilesTableOutput{w}
	out.BasicHeader()
	defer out.BasicFooter()
	ctx := context.Background()
	resp, err := cmd.list.GetDEPProfile(ctx, *flUUID)
	if err != nil {
		return err
	}
	fmt.Fprintf(out.w, "%s\t%v\t%v\t%v\t%s\n",
		resp.ProfileName,
		resp.IsMandatory,
		resp.IsMDMRemovable,
		resp.AwaitDeviceConfigured,
		strings.Join(resp.SkipSetupItems, ","),
	)

	return nil
}
