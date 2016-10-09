package management

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/micromdm/dep"
	"github.com/micromdm/micromdm/device"
	"github.com/micromdm/micromdm/workflow"
	"golang.org/x/net/context"
)

func TestListDevices(t *testing.T) {
	server, svc := newServer(t)
	defer teardown()
	defer server.Close()

	fetchDEPDevices(t, server, svc)
	testListDevicesHTTP(t, svc, server, http.StatusOK)
}
func TestShowDevice(t *testing.T) {
	server, svc := newServer(t)
	defer teardown()
	defer server.Close()

	fetchDEPDevices(t, server, svc)
	devices := testListDevicesHTTP(t, svc, server, http.StatusOK)

	for _, d := range devices {
		var dev device.Device
		testGetHTTPInto("devices", t, svc, server, d.UUID, http.StatusOK, &dev)
	}
}

func TestAddWorkflowWithProfiles(t *testing.T) {
	server, svc := newServer(t)
	defer teardown()
	defer server.Close()

	profileData := []byte(`{
    "payload_identifier": "com.micromdm.example2",
    "data" : "fooProfile"
	}`)

	pfBody := testAddHTTP("profiles", t, svc, server, profileData, http.StatusCreated)
	var pf workflow.Profile
	err := json.NewDecoder(pfBody).Decode(&pf)
	if err != nil {
		t.Fatal(err)
	}

	wf := workflow.Workflow{
		Name:     "test_workflow",
		Profiles: []workflow.Profile{pf},
	}
	wfData, err := json.Marshal(wf)
	if err != nil {
		t.Fatal(err)
	}

	wfBody := testAddHTTP("workflows", t, svc, server, wfData, http.StatusCreated)

	var returned workflow.Workflow
	err = json.NewDecoder(wfBody).Decode(&returned)
	if err != nil {
		t.Fatal(err)
	}

	if !returned.HasProfile(pf.PayloadIdentifier) {
		t.Fatal("returned workflow must have profile")
	}
}

// create some workflows to be used by tests
func addWorkflows(t *testing.T, server *httptest.Server, svc Service) []workflow.Workflow {
	profileData := []byte(`{
    "payload_identifier": "com.micromdm.example2",
    "data" : "fooProfile"
	}`)

	pfBody := testAddHTTP("profiles", t, svc, server, profileData, http.StatusCreated)
	var pf workflow.Profile
	err := json.NewDecoder(pfBody).Decode(&pf)
	if err != nil {
		t.Fatal(err)
	}

	wf := workflow.Workflow{
		Name:     "test_workflow",
		Profiles: []workflow.Profile{pf},
	}
	wfData, err := json.Marshal(wf)
	if err != nil {
		t.Fatal(err)
	}

	wfBody := testAddHTTP("workflows", t, svc, server, wfData, http.StatusCreated)

	var returned workflow.Workflow
	err = json.NewDecoder(wfBody).Decode(&returned)
	if err != nil {
		t.Fatal(err)
	}
	return []workflow.Workflow{returned}
}

func TestAddWorkflow(t *testing.T) {
	server, svc := newServer(t)
	defer teardown()
	defer server.Close()

	workflowData := []byte(`{
    "name": "testWorkflow"
	}`)

	var addTests = []struct {
		in       []byte
		expected int
	}{
		{
			in:       nil,
			expected: http.StatusBadRequest,
		},
		{
			in:       workflowData,
			expected: http.StatusCreated,
		},
		{
			in:       workflowData,
			expected: http.StatusConflict,
		},
	}

	for _, tt := range addTests {
		testAddHTTP("workflows", t, svc, server, tt.in, tt.expected)
	}
}

func TestDeleteProfile(t *testing.T) {
	server, svc := newServer(t)
	defer teardown()
	defer server.Close()
	testDeleteHTTP(t, svc, server, "foo", http.StatusBadRequest)
	testDeleteHTTP(t, svc, server, "65fd8c1d-20bc-4342-9bb6-258120d6b124", http.StatusNotFound)

	profileData := []byte(`{
    "payload_identifier": "com.micromdm.example2",
    "data" : "fooProfile"
	}`)

	testAddHTTP("profiles", t, svc, server, profileData, http.StatusCreated)
	profiles := testListHTTP(t, svc, server, http.StatusOK)
	for _, p := range profiles {
		testDeleteHTTP(t, svc, server, p.UUID, http.StatusNoContent)
	}
}

func TestShowProfile(t *testing.T) {
	server, svc := newServer(t)
	defer teardown()
	defer server.Close()

	profileData := []byte(`{
    "payload_identifier": "com.micromdm.example2",
    "data" : "fooProfile"
	}`)

	testGetHTTP(t, svc, server, "foo", http.StatusBadRequest)
	testGetHTTP(t, svc, server, "036d339c-4fe4-4d6e-a051-65fafbec8c93", http.StatusNotFound)

	testAddHTTP("profiles", t, svc, server, profileData, http.StatusCreated)
	profiles := testListHTTP(t, svc, server, http.StatusOK)
	for _, p := range profiles {
		returned := testGetHTTP(t, svc, server, p.UUID, http.StatusOK)
		// check that we get what we pass in
		if returned.PayloadIdentifier != p.PayloadIdentifier {
			t.Fatal("expected", p.PayloadIdentifier, "got", returned.PayloadIdentifier)
		}
	}

}

func testDeleteHTTP(t *testing.T, svc Service, server *httptest.Server, uuid string, expectedStatus int) {
	client := http.DefaultClient
	theURL := server.URL + "/management/v1/profiles" + "/" + uuid
	req, err := http.NewRequest("DELETE", theURL, nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != expectedStatus {
		io.Copy(os.Stdout, resp.Body)
		t.Fatal("expected", expectedStatus, "got", resp.StatusCode)
	}
}

func testGetHTTPInto(endpoint string, t *testing.T, svc Service, server *httptest.Server, uuid string, expectedStatus int, into interface{}) {
	client := http.DefaultClient
	theURL := server.URL + "/management/v1/" + endpoint + "/" + uuid
	resp, err := client.Get(theURL)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != expectedStatus {
		io.Copy(os.Stdout, resp.Body)
		t.Fatal("expected", expectedStatus, "got", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(into); err != nil {
		t.Log("failed to decode profiles from GET response")
		t.Fatal(err)
	}
}

func testGetHTTP(t *testing.T, svc Service, server *httptest.Server, uuid string, expectedStatus int) *workflow.Profile {
	client := http.DefaultClient
	theURL := server.URL + "/management/v1/profiles" + "/" + uuid
	resp, err := client.Get(theURL)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != expectedStatus {
		io.Copy(os.Stdout, resp.Body)
		t.Fatal("expected", expectedStatus, "got", resp.StatusCode)
	}

	// test decoding the result into a struct
	var profile workflow.Profile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		t.Log("failed to decode profiles from GET response")
		t.Fatal(err)
	}
	return &profile
}

func TestListWorkflows(t *testing.T) {
	server, svc := newServer(t)
	defer teardown()
	defer server.Close()
	addWorkflows(t, server, svc)
	testListWorkflowsHTTP(t, svc, server, http.StatusOK)
}

func TestListProfiles(t *testing.T) {
	server, svc := newServer(t)
	defer teardown()
	defer server.Close()

	profileData := []byte(`{
    "payload_identifier": "com.micromdm.example2",
    "data" : "fooProfile"
	}`)

	testAddHTTP("profiles", t, svc, server, profileData, http.StatusCreated)
	testListHTTP(t, svc, server, http.StatusOK)

}

func TestAddProfile(t *testing.T) {
	server, svc := newServer(t)
	defer teardown()
	defer server.Close()

	profileData := []byte(`{
    "payload_identifier": "com.micromdm.example2",
    "data" : "fooProfile"
	}`)

	var addTests = []struct {
		in       []byte
		expected int
	}{
		{
			in:       nil,
			expected: http.StatusBadRequest,
		},
		{
			in:       profileData,
			expected: http.StatusCreated,
		},
		{
			in:       profileData,
			expected: http.StatusConflict,
		},
	}

	for _, tt := range addTests {
		testAddHTTP("profiles", t, svc, server, tt.in, tt.expected)
	}
}

var testConn = "user=micromdm password=micromdm dbname=micromdm sslmode=disable"

func newServer(t *testing.T) (*httptest.Server, Service) {
	ctx := context.Background()
	l := log.NewLogfmtLogger(os.Stderr)
	logger := log.NewContext(l).With("source", "testing")
	ds, err := device.NewDB("postgres", testConn, logger)
	if err != nil {
		t.Fatal(err)
	}

	ps, err := workflow.NewDB("postgres", testConn, logger)
	if err != nil {
		t.Fatal(err)
	}
	// dep client
	config := &dep.Config{
		ConsumerKey:    "CK_48dd68d198350f51258e885ce9a5c37ab7f98543c4a697323d75682a6c10a32501cb247e3db08105db868f73f2c972bdb6ae77112aea803b9219eb52689d42e6",
		ConsumerSecret: "CS_34c7b2b531a600d99a0e4edcf4a78ded79b86ef318118c2f5bcfee1b011108c32d5302df801adbe29d446eb78f02b13144e323eb9aad51c79f01e50cb45c3a68",
		AccessToken:    "AT_927696831c59ba510cfe4ec1a69e5267c19881257d4bca2906a99d0785b785a6f6fdeb09774954fdd5e2d0ad952e3af52c6d8d2f21c924ba0caf4a031c158b89",
		AccessSecret:   "AS_c31afd7a09691d83548489336e8ff1cb11b82b6bca13f793344496a556b1f4972eaff4dde6deb5ac9cf076fdfa97ec97699c34d515947b9cf9ed31c99dded6ba",
	}

	dc, err := dep.NewClient(config, dep.ServerURL("http://localhost:9000"))
	if err != nil {
		t.Fatal(err)
	}

	as, err := applications.NewDB("postgres", testConn, logger)
	if err != nil {
		t.Fatal(err)
	}

	svc := NewService(ds, ps, dc, nil, as)
	handler := ServiceHandler(ctx, svc, logger)
	server := httptest.NewServer(handler)
	return server, svc
}

func testListDevicesHTTP(t *testing.T, svc Service, server *httptest.Server, expectedStatus int) []device.Device {
	client := http.DefaultClient
	theURL := server.URL + "/management/v1/devices"
	resp, err := client.Get(theURL)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != expectedStatus {
		io.Copy(os.Stdout, resp.Body)
		t.Fatal("expected", expectedStatus, "got", resp.StatusCode)
	}

	// test decoding the result into a struct
	var devices []device.Device
	if err := json.NewDecoder(resp.Body).Decode(&devices); err != nil {
		t.Log("failed to decode profiles from list response")
		t.Fatal(err)
	}
	return devices
}

func testListWorkflowsHTTP(t *testing.T, svc Service, server *httptest.Server, expectedStatus int) []workflow.Workflow {
	client := http.DefaultClient
	theURL := server.URL + "/management/v1/workflows"
	resp, err := client.Get(theURL)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != expectedStatus {
		io.Copy(os.Stdout, resp.Body)
		t.Fatal("expected", expectedStatus, "got", resp.StatusCode)
	}

	// test decoding the result into a struct
	var workflows []workflow.Workflow
	if err := json.NewDecoder(resp.Body).Decode(&workflows); err != nil {
		t.Log("failed to decode profiles from list response")
		t.Fatal(err)
	}
	return workflows
}

func testListHTTP(t *testing.T, svc Service, server *httptest.Server, expectedStatus int) []workflow.Profile {
	client := http.DefaultClient
	theURL := server.URL + "/management/v1/profiles"
	resp, err := client.Get(theURL)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != expectedStatus {
		io.Copy(os.Stdout, resp.Body)
		t.Fatal("expected", expectedStatus, "got", resp.StatusCode)
	}

	// test decoding the result into a struct
	var profiles []workflow.Profile
	if err := json.NewDecoder(resp.Body).Decode(&profiles); err != nil {
		t.Log("failed to decode profiles from list response")
		t.Fatal(err)
	}
	return profiles
}

func testAddHTTP(endpoint string, t *testing.T, svc Service, server *httptest.Server, data []byte, expectedStatus int) io.Reader {
	body := &nopCloser{bytes.NewBuffer(data)}

	client := http.DefaultClient
	theURL := server.URL + "/management/v1/" + endpoint
	resp, err := client.Post(theURL, "application/json", body)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != expectedStatus {
		io.Copy(os.Stdout, resp.Body)
		t.Fatal("expected", expectedStatus, "got", resp.StatusCode)
	}

	return resp.Body
}

func TestFetchDEPDevices(t *testing.T) {
	ctx := context.Background()
	logger := log.NewLogfmtLogger(os.Stderr)
	ds, err := device.NewDB("postgres", testConn, logger)
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	// Auth configuration
	config := &dep.Config{
		ConsumerKey:    "CK_48dd68d198350f51258e885ce9a5c37ab7f98543c4a697323d75682a6c10a32501cb247e3db08105db868f73f2c972bdb6ae77112aea803b9219eb52689d42e6",
		ConsumerSecret: "CS_34c7b2b531a600d99a0e4edcf4a78ded79b86ef318118c2f5bcfee1b011108c32d5302df801adbe29d446eb78f02b13144e323eb9aad51c79f01e50cb45c3a68",
		AccessToken:    "AT_927696831c59ba510cfe4ec1a69e5267c19881257d4bca2906a99d0785b785a6f6fdeb09774954fdd5e2d0ad952e3af52c6d8d2f21c924ba0caf4a031c158b89",
		AccessSecret:   "AS_c31afd7a09691d83548489336e8ff1cb11b82b6bca13f793344496a556b1f4972eaff4dde6deb5ac9cf076fdfa97ec97699c34d515947b9cf9ed31c99dded6ba",
	}

	dc, err := dep.NewClient(config, dep.ServerURL("http://localhost:9000"))
	if err != nil {
		t.Fatal(err)
	}
	svc := NewService(ds, nil, dc, nil, nil)
	handler := ServiceHandler(ctx, svc, logger)
	server := httptest.NewServer(handler)
	defer server.Close()

	client := http.DefaultClient
	theURL := server.URL + "/management/v1/devices/fetch"
	resp, err := client.Post(theURL, "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		io.Copy(os.Stdout, resp.Body)
		t.Fatal("expected", http.StatusOK, "got", resp.StatusCode)
	}
}

func fetchDEPDevices(t *testing.T, server *httptest.Server, svc Service) *http.Response {
	client := http.DefaultClient
	theURL := server.URL + "/management/v1/devices/fetch"
	resp, err := client.Post(theURL, "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		io.Copy(os.Stdout, resp.Body)
		t.Fatal("expected", http.StatusOK, "got", resp.StatusCode)
	}

	return resp

}

// a face io.ReadCloser for constructing request Body
type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func teardown() {
	db, err := sqlx.Open("postgres", testConn)
	if err != nil {
		panic(err)
	}

	drop := `
	DROP TABLE IF EXISTS device_workflow;
	DROP TABLE IF EXISTS devices;
	DROP INDEX IF EXISTS devices.serial_idx;
	DROP INDEX IF EXISTS devices.udid_idx;
	DROP TABLE IF EXISTS workflow_profile;
	DROP TABLE IF EXISTS workflow_workflow;
	DROP TABLE IF EXISTS workflows;
	DROP TABLE IF EXISTS profiles;
	`
	db.MustExec(drop)
	defer db.Close()
}
