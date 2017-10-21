package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/groob/plist"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/crypto/password"
	"github.com/micromdm/micromdm/user"
)

func (cmd *applyCommand) applyUser(args []string) error {
	flagset := flag.NewFlagSet("users", flag.ExitOnError)
	var (
		flUserManifest = flagset.String("f", "", "Path to user manifest")
		flTemplate     = flagset.Bool("template", false, "Print a JSON example of a user manifest.")
		flPassword     = flagset.String("password", "", "Password of the user. Only required when creating a new user.")
	)
	flagset.Usage = usageFor(flagset, "mdmctl apply users [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}

	if *flTemplate {
		printUserTemplate()
		return nil
	}

	manifestData, err := ioutil.ReadFile(*flUserManifest)
	if err != nil {
		return errors.Wrap(err, "read user manifest file")
	}

	var manifest user.User
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return errors.Wrap(err, "unmarshal user manifest")
	}

	if manifest.UUID == "" && *flPassword == "" {
		return errors.New("password argument must be provided when creating a user")
	}

	if *flPassword != "" {
		salted, err := password.SaltedSHA512PBKDF2(*flPassword)
		if err != nil {
			return errors.Wrap(err, "salting plaintext password")
		}
		hashDict := struct {
			SaltedSHA512PBKDF2 password.SaltedSHA512PBKDF2Dictionary `plist:"SALTED-SHA512-PBKDF2"`
		}{
			SaltedSHA512PBKDF2: salted,
		}
		hashPlist, err := plist.Marshal(hashDict)
		if err != nil {
			return errors.Wrap(err, "marshal salted password to plist")
		}
		manifest.PasswordHash = hashPlist
	}

	usr, err := cmd.applysvc.ApplyUser(context.TODO(), manifest)
	if err != nil {
		return errors.Wrap(err, "apply user with mdmctl")
	}

	f, err := os.OpenFile(*flUserManifest, os.O_RDWR, os.ModeAppend)
	if err != nil {
		return errors.Wrapf(err, "open user manifest %s", *flUserManifest)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(usr); err != nil {
		return errors.Wrap(err, "encode user manifest")
	}
	return nil

}

func printUserTemplate() {
	jsn := `
{
  "user_shortname": "",
  "user_longname": "",
  "hidden": false
}
`
	fmt.Println(jsn)
}
