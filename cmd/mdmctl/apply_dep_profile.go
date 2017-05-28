package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/micromdm/dep"
	"github.com/pkg/errors"
)

func (cmd *applyCommand) applyDEPProfile(args []string) error {
	flagset := flag.NewFlagSet("dep-profiles", flag.ExitOnError)
	var (
		flProfilePath = flagset.String("f", "", "filename of DEP profile to apply")
	)
	flagset.Usage = usageFor(flagset, "mdmctl apply dep-profiles [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}
	if *flProfilePath == "" {
		flagset.Usage()
		return errors.New("bad input: must provide -f parameter")
	}

	pf, err := os.Open(*flProfilePath)
	if err != nil {
		return errors.Wrap(err, "opening DEP profile file")
	}
	defer pf.Close()

	var profile dep.Profile
	if err := json.NewDecoder(pf).Decode(&profile); err != nil {
		return errors.Wrap(err, "decode DEP Profile JSON")
	}

	resp, err := cmd.applysvc.DefineDEPProfile(context.TODO(), &profile)
	if err != nil {
		return errors.Wrap(err, "define dep profile")
	}

	// TODO: it would be nice to encode back a profile that save the
	// UUID for future reference.
	fmt.Printf("Defined DEP Profile with UUID %s\n", resp.ProfileUUID)
	return nil
}
