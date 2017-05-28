package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
)

type depaccountTableOutput struct{ w *tabwriter.Writer }

func (out *depaccountTableOutput) BasicHeader() {
	fmt.Fprintf(out.w, "OrgName\tOrgPhone\tOrgEmail\tServerName\n")
}

func (out *depaccountTableOutput) BasicFooter() {
	out.w.Flush()
}

func (cmd *getCommand) getDEPAccount(args []string) error {
	flagset := flag.NewFlagSet("dep-account", flag.ExitOnError)
	flagset.Usage = usageFor(flagset, "mdmctl get dep-account [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}
	w := tabwriter.NewWriter(os.Stderr, 0, 4, 2, ' ', 0)
	out := &depaccountTableOutput{w}
	out.BasicHeader()
	defer out.BasicFooter()
	ctx := context.Background()
	resp, err := cmd.list.GetDEPAccountInfo(ctx)
	if err != nil {
		return err
	}
	fmt.Fprintf(out.w, "%s\t%s\t%s\t%s\n", resp.OrgName, resp.OrgPhone, resp.OrgEmail, resp.ServerName)
	return nil
}
