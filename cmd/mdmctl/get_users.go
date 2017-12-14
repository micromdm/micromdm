package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/core/list"
)

type usersTableOutput struct{ w *tabwriter.Writer }

func (out *usersTableOutput) BasicHeader() {
	fmt.Fprintf(out.w, "UUID\tUDID\tUserID\tUserShortName\tUserLongName\n")
}

func (out *usersTableOutput) BasicFooter() {
	out.w.Flush()
}

func (cmd *getCommand) getUsers(args []string) error {
	flagset := flag.NewFlagSet("users", flag.ExitOnError)
	flagset.Usage = usageFor(flagset, "mdmctl get users [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	out := &usersTableOutput{w}
	out.BasicHeader()
	defer out.BasicFooter()

	users, err := cmd.list.ListUsers(context.TODO(), list.ListUsersOption{})
	if err != nil {
		return errors.Wrap(err, "list users")
	}

	for _, u := range users {
		fmt.Fprintf(out.w, "%s\t%s\t%s\t%v\t%s\n", u.UUID, u.UDID, u.UserID, u.UserShortname, u.UserLongname)
	}

	return nil
}
