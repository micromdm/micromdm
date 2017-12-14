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

const noUUIDText = `The DEP API does not support listing profiles. 
A UUID flag must be specified. To get currently assigned profile UUIDs run
	mdmctl get dep-devices -serials=serial1,serial2,serial3
The output of the dep-devices response will contain the profile UUIDs.
`

func (cmd *getCommand) getDEPProfiles(args []string) error {
	flagset := flag.NewFlagSet("dep-profiles", flag.ExitOnError)
	var (
		flProfilePath = flagset.String("f", "", "filename of DEP profile to apply")
		flUUID        = flagset.String("uuid", "", "DEP Profile UUID(required)")
	)
	flagset.Usage = usageFor(flagset, "mdmctl get dep-profiles [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}

	if *flUUID == "" {
		fmt.Println(noUUIDText)
		flagset.Usage()
		os.Exit(1)
	}

	ctx := context.Background()
	resp, err := cmd.list.GetDEPProfile(ctx, *flUUID)
	if err != nil {
		return err
	}

	if *flProfilePath == "" {
		w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
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
