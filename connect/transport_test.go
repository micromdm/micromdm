package connect

import (
	"bytes"
	"github.com/go-kit/kit/log"
	"github.com/micromdm/micromdm/applications"
	"github.com/micromdm/micromdm/command"
	"github.com/micromdm/micromdm/device"
	"golang.org/x/net/context"
	"net/http"
	"net/http/httptest"
	"testing"
)

var m *http.ServeMux
var respRec *httptest.ResponseRecorder
var ctx context.Context
var handler http.Handler
var kitService Service
var logger log.Logger

// TODO: Mock services
var deviceDB device.Datastore
var appsDB applications.Datastore
var commandSvc command.Service

func setup() {
	ctx = context.Background()
	logger = log.NewNopLogger()
	m = http.NewServeMux()
	kitService = NewService(deviceDB, appsDB, commandSvc)
	handler = ServiceHandler(ctx, kitService, logger)

	m.Handle("/mdm/connect", handler)
	respRec = httptest.NewRecorder()
	http.Handle("/", m)
}

func tearDown() {}

func TestServiceHandler(t *testing.T) {
	setup()
	defer tearDown()

	fixture := "<plist></plist>"
	reader := bytes.NewReader([]byte(fixture))

	req, err := http.NewRequest("PUT", "/mdm/connect", reader)
	if err != nil {
		t.Fatal(err)
	}

	m.ServeHTTP(respRec, req)
}
