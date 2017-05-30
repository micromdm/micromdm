package main

import (
	"context"
	"encoding/json"
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
	flProfilePath := flagset.String("f", "", "filename of DEP profile to apply")
	flUUID := flagset.String("uuid", "", "DEP Profile UUID")
	flagset.Usage = usageFor(flagset, "mdmctl get dep-profiles [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}

	ctx := context.Background()
	resp, err := cmd.list.GetDEPProfile(ctx, *flUUID)
	if err != nil {
		return err
	}

	if *flProfilePath == "" {
		w := tabwriter.NewWriter(os.Stderr, 0, 4, 2, ' ', 0)
		out := &depProfilesTableOutput{w}
		out.BasicHeader()
		defer out.BasicFooter()

		fmt.Fprintf(out.w, "%s\t%v\t%v\t%v\t%s\n",
			resp.ProfileName,
			resp.IsMandatory,
			resp.IsMDMRemovable,
			resp.AwaitDeviceConfigured,
			strings.Join(resp.SkipSetupItems, ","),
		)
	} else {
		var output *os.File
		{
			if *flProfilePath == "-" {
				output = os.Stdout
			} else {
				var err error
				output, err = os.Create(*flProfilePath)
				if err != nil {
					return err
				}
				defer output.Close()
			}
		}

		// TODO: perhaps we want to store the raw DEP profile for storage
		// as we may have problems with default/non-default values getting
		// omitted in the marshalled JSON
		enc := json.NewEncoder(output)
		enc.SetIndent("", "  ")
		err = enc.Encode(resp)
		if err != nil {
			return err
		}

		if *flProfilePath != "-" {
			fmt.Printf("wrote DEP profile %s to: %s\n", *flUUID, *flProfilePath)
		}
	}

	return nil
}
