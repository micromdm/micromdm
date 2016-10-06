package command

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/micromdm/mdm"
	"os"
	"testing"
)

type datastoreFixtures struct {
	ds     Datastore
	logger log.Logger
}

func setup() (datastoreFixtures, error) {
	logger := log.NewLogfmtLogger(os.Stdout)
	commandsDb, err := NewDB("redis", "localhost", logger)
	if err != nil {
		return nil, err
	}

	return datastoreFixtures{ds: commandsDb, logger: logger}
}

func teardown() {

}

func TestService_Commands(t *testing.T) {
	fixtures, err := setup()
	defer teardown()
	if err != nil {
		t.Errorf("error making new datastore: %v", err)
	}

	var commands []mdm.Payload
	commands, err = fixtures.ds.Commands("ABCDEF")
	if err != nil {
		t.Errorf("datastore.Commands returned error: %v", err)
	}

	fmt.Printf("%v", commands)
}
