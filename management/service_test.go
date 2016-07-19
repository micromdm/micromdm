package management

import "testing"

func svcSetup() {

}

func svcTearDown() {

}

func TestService_InstalledApps(t *testing.T) {
	svcSetup()
	defer svcTearDown()

	svc := NewService(nil, nil, nil, nil, nil)
	_, err := svc.InstalledApps("00000000-1111-2222-3333-444455556666")
	if err != nil {
		t.Fatal(err)
	}

}
