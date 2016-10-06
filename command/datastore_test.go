package command

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/micromdm/mdm"
	"os"
	"testing"
)

var (
	ds     Datastore
	logger log.Logger
	err    error
)

//func setup() {
//	l = log.NewLogfmtLogger(os.Stdout)
//	d, err := NewDB("redis", "localhost", l)
//	if err != nil {
//		panic("Test set up failed")
//	}
//}
//
//func teardown() {
//
//}

func TestService_Commands(t *testing.T) {
	logger = log.NewLogfmtLogger(os.Stdout)
	ds, err = NewDB("redis", "localhost:6379", logger)
	if err != nil {
		t.Errorf("error making new datastore: %v", err)
	}

	var commands []mdm.Payload
	commands, err = ds.Commands("ABCDEF")
	if err != nil {
		t.Errorf("datastore.Commands returned error: %v", err)
	}

	fmt.Printf("%v", commands)
}
