package enroll

import (
	"testing"
)

func TestEnrollProfile(t *testing.T) {
	svc := new(service)
	profile, err := svc.MakeEnrollmentProfile("")
	if err != nil {
		t.Fatal(err)
	}

	var payloadContent MDMPayloadContent
	for _, payload := range profile.PayloadContent {
		if c, ok := payload.(MDMPayloadContent); ok {
			payloadContent = c
		}
	}

	if have, want := payloadContent.AccessRights, AccessRights(8191); have != want {
		t.Errorf("have %d, want %d", have, want)
	}

	var hasPerUserConnections bool
	for _, cap := range payloadContent.ServerCapabilities {
		if cap == perUserConnections {
			hasPerUserConnections = true
		}
	}

	if have, want := hasPerUserConnections, true; have != want {
		t.Errorf("missing ServerCapabilities: macOS enrollment profile requires %s", perUserConnections)
	}

	if have, want := payloadContent.CheckInURL, "/mdm/checkin"; have != want {
		t.Errorf("have checkin URL %s, want %s", have, want)
	}
}

func TestEnrollProfileWithId(t *testing.T) {
	id := "03c6b007-abba-46a8-8811-a639dbcfca81"

	svc := new(service)
	profile, err := svc.MakeEnrollmentProfile(id)
	if err != nil {
		t.Fatal(err)
	}

	var payloadContent MDMPayloadContent
	for _, payload := range profile.PayloadContent {
		if c, ok := payload.(MDMPayloadContent); ok {
			payloadContent = c
		}
	}

	if have, want := payloadContent.CheckInURL, "/mdm/checkin?id="+id; have != want {
		t.Errorf("have checkin URL %s, want %s", have, want)
	}
}