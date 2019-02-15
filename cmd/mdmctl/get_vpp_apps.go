package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/pkg/errors"
)

type VPPAppData struct {
	Name           string   `json:"name"`
	DeviceFamilies []string `json:"deviceFamilies"`
}

type vppAppsTableOutput struct{ w *tabwriter.Writer }

func (out *vppAppsTableOutput) BasicHeader() {
	fmt.Fprintf(out.w, "ID\tName\tDeviceFamilies\tLicenses\tAssigned\tAvailable\tTokenUUID\n")
}

func (out *vppAppsTableOutput) BasicFooter() {
	out.w.Flush()
}

func (cmd *getCommand) getVPPApps(args []string) error {
	flagset := flag.NewFlagSet("vpp-apps", flag.ExitOnError)
	var (
		flIDFilter = flagset.String("id", "", "comma separated list of vpp app ids to filter results")
		flVerbose  = flagset.Bool("verbose", false, "specify -verbose to get full apps details")
	)
	flagset.Usage = usageFor(flagset, "mdmctl get vpp-apps [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}
	ctx := context.Background()

	ids := strings.Split(*flIDFilter, ",")

	response, err := cmd.vppsvc.GetVPPApps(ctx, ids)
	if err != nil {
		return err
	}

	if *flVerbose == true {
		bytes, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(bytes))
		return nil
	}

	apps := response.VPPAppsList

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	out := vppAppsTableOutput{w}
	out.BasicHeader()
	defer out.BasicFooter()
	for _, a := range apps {
		results, err := getAppData(a.Metadata.Results[a.Asset.AdamIDStr])
		if err != nil {
			return err
		}

		fmt.Fprintf(out.w, "%s\t%s\t%s\t%d\t%d\t%d\t%s\n", a.Asset.AdamIDStr, results.Name, results.DeviceFamilies, a.Asset.TotalCount, a.Asset.AssignedCount, a.Asset.AvailableCount, a.ClientContext.GUID)
	}
	return nil
}

func getAppData(metadata interface{}) (*VPPAppData, error) {

	bytes, err := json.Marshal(metadata)
	if err != nil {
		return nil, errors.Wrap(err, "get vpp app data")
	}

	var response VPPAppData
	err = json.Unmarshal(bytes, &response)
	return &response, errors.Wrap(err, "get VPPAppData")
}
