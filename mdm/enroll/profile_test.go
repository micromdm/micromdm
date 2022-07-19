package enroll

import (
	"testing"
)

func TestEnrollProfile(t *testing.T) {
	svc := new(service)
	profile, err := svc.MakeEnrollmentProfile()
	if err != nil {
		t.Fatal(err)
	}

	mdmPayloads := profile.MDMPayloads()
	if len(mdmPayloads) != 1 {
		t.Fatal("number of MDM payloads is not 1")
	}
	payloadContent := mdmPayloads[0]

	if have, want := payloadContent.AccessRights, 8191; have != want {
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
}
