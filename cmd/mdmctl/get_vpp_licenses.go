package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/micromdm/micromdm/vpp"
)

type vppLicenseWithSerials struct {
	AdamID  string
	Serials []string
}

type vppLicensesTableOutput struct{ w *tabwriter.Writer }

func (out *vppLicensesTableOutput) BasicHeader() {
	fmt.Fprintf(out.w, "AppID\tTokenUUID\tSerials\n")
}

func (out *vppLicensesTableOutput) BasicFooter() {
	out.w.Flush()
}

func (cmd *getCommand) getVPPLicenses(args []string) error {
	flagset := flag.NewFlagSet("vpp-licenses", flag.ExitOnError)
	var (
		flIDFilter = flagset.String("id", "", "specify the id of a vpp app for which to get associated license details")
		flSerial   = flagset.String("serial", "", "a single serial for which to get associated license details")
		flVerbose  = flagset.Bool("verbose", false, "specify -verbose to get full licenses details")
	)
	flagset.Usage = usageFor(flagset, "mdmctl get vpp-licenses [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}
	ctx := context.Background()

	var options vpp.LicensesSrvOptions
	if *flIDFilter != "" {
		options.AdamID = *flIDFilter
	}

	var serial = strings.Split(*flSerial, ",")[0]
	if *flSerial != "" {
		options.SerialNumber = serial
	}

	licensesSrv, err := cmd.vppsvc.GetLicensesSrv(ctx, options)
	if err != nil {
		return err
	}

	if *flVerbose == true {
		bytes, err := json.MarshalIndent(licensesSrv, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(bytes))
		return nil
	}

	var clientContext vpp.ClientContext
	vpp.DecodeToClientContext(licensesSrv.ClientContext, &clientContext)

	licenses := licensesSrv.Licenses

	var licensesList []vppLicenseWithSerials
	for _, license := range licenses {
		found := false
		if len(licensesList) > 0 {
			for i, saved := range licensesList {
				if license.AdamIDStr == saved.AdamID {
					found = true
					if license.SerialNumber != "" {
						saved.Serials = append(saved.Serials, license.SerialNumber)
					}
					licensesList[i] = saved
				}
			}
		}

		if found == false {
			licensesList = append(licensesList, vppLicenseWithSerials{
				AdamID:  license.AdamIDStr,
				Serials: []string{license.SerialNumber},
			})
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	out := vppLicensesTableOutput{w}
	out.BasicHeader()
	defer out.BasicFooter()
	for _, a := range licensesList {
		//fmt.Fprintf(out.w, "%s\t%s\t%s\n", a.AdamID, clientContext.GUID, a.Serials)
		fmt.Fprintf(out.w, "%s\t%s\t%s\n", a.AdamID, clientContext.GUID, a.Serials)
	}
	return nil
}
