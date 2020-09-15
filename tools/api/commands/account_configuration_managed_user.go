package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/groob/plist"
	"github.com/micromdm/micromdm/mdm/mdm"
	"github.com/micromdm/micromdm/pkg/crypto/password"
)

func main() {
	var (
		flPassword = flag.String("password", "groob", "Password of the user. Only required when creating a new user.")
	)
	flag.Parse()

	salted, err := password.SaltedSHA512PBKDF2(*flPassword)
	if err != nil {
		log.Fatal(err)
	}

	hashDict := struct {
		SaltedSHA512PBKDF2 password.SaltedSHA512PBKDF2Dictionary `plist:"SALTED-SHA512-PBKDF2"`
	}{
		SaltedSHA512PBKDF2: salted,
	}

	hashPlist, err := plist.Marshal(hashDict)
	if err != nil {
		log.Fatal(err)
	}

	cmd := &mdm.CommandRequest{
		Command: &mdm.Command{
			RequestType: "AccountConfiguration",
			AccountConfiguration: &mdm.AccountConfiguration{
				SkipPrimarySetupAccountCreation: true,
				ManagedLocalUserShortName:       "groob",
				AutoSetupAdminAccounts: []mdm.AdminAccount{
					{
						ShortName:    "groob",
						FullName:     "Victor Vrantchan",
						PasswordHash: hashPlist,
					},
				},
			},
		},
	}

	out, err := json.MarshalIndent(cmd, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
}
