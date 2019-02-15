package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	"github.com/micromdm/micromdm/vpp"
	"github.com/pkg/errors"
)

func (cmd *removeCommand) removeVPPApp(args []string) error {
	flagset := flag.NewFlagSet("vpp-apps", flag.ExitOnError)
	var (
		flAppID   = flagset.String("id", "", "specify the id of the vpp app from which to disassociate serials")
		flSerials = flagset.String("serials", "", "comma separated list of serials to disassociate from the vpp app")
		flVerbose = flagset.Bool("verbose", false, "specify -verbose to get full disassociation details")
	)
	flagset.Usage = usageFor(flagset, "mdmctl remove vpp-apps [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}
	ctx := context.Background()

	if *flAppID == "" || *flSerials == "" {
		flagset.Usage()
		return errors.New("bad input: must provide an app id and at least one serial")
	}

	serials := strings.Split(*flSerials, ",")
	manage, err := cmd.vppsvc.ManageVPPLicensesByAdamIdSrv(ctx, vpp.ManageVPPLicensesByAdamIdSrvOptions{
		AdamIDStr:                 *flAppID,
		DisassociateSerialNumbers: serials,
	})
	if err != nil {
		return err
	}

	if *flVerbose == false {
		if manage.Status == 0 {
			fmt.Printf("Success: Disssociated license for %s, from %s.\n\n", manage.AdamIdStr, serials)
			fmt.Printf("To confirm, run 'mdmctl get vpp-licenses -id %s'\n", manage.AdamIdStr)
			fmt.Println("The license counts in 'mdmctl get vpp-apps' will be updated in a few minutes.")
			return nil
		} else {
			fmt.Println("Error: One or more license removals failed.")
		}
	}

	bytes, err := json.MarshalIndent(manage, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))
	return nil
}
