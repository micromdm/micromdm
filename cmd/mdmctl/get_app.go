package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"text/tabwriter"

	"github.com/micromdm/micromdm/platform/api/server/list"
)

type appsTableOutput struct{ w *tabwriter.Writer }

func (out *appsTableOutput) BasicHeader() {
	fmt.Fprintf(out.w, "Name\tManifestURL\n")
}

func (out *appsTableOutput) BasicFooter() {
	out.w.Flush()
}

func (cmd *getCommand) getApps(args []string) error {
	flagset := flag.NewFlagSet("apps", flag.ExitOnError)
	var (
		flNameFilter = flagset.String("name", "", "specify the name of the app to get full details")
		flOutputPath = flagset.String("f", "-", "path to save file to. defaults to stdout.")
	)
	flagset.Usage = usageFor(flagset, "mdmctl get apps [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}
	ctx := context.Background()
	apps, err := cmd.list.ListApplications(ctx, list.ListAppsOption{
		FilterName: []string{*flNameFilter},
	})
	if err != nil {
		return err
	}
	if *flNameFilter != "" && (len(apps) != 0) {
		payload := apps[0].Payload
		if *flOutputPath == "-" {
			fmt.Println(string(payload))
			return nil
		}
		return ioutil.WriteFile(*flOutputPath, payload, 0644)
	}
	rURL, err := repoURL(cmd.config.ServerURL)
	if err != nil {
		return err
	}
	w := tabwriter.NewWriter(os.Stderr, 0, 4, 2, ' ', 0)
	out := appsTableOutput{w}
	out.BasicHeader()
	defer out.BasicFooter()
	for _, a := range apps {
		manifestURL := rURL + "/" + a.Name
		fmt.Fprintf(out.w, "%s\t%s\n", a.Name, manifestURL)
	}
	return nil
}
