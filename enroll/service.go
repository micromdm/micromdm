package enroll

type Service interface {
	Enroll()
}

type service struct {
	Url           string
	SCEPUrl       string
	SCEPChallenge string
	Topic         string // APNS Topic for MDM notifications
}

func (svc service) Enroll() (Profile, error) {
	profile := NewProfile()
	profile.PayloadIdentifier = "com.github.micromdm.micromdm.mdm"
	profile.PayloadOrganization = "MicroMDM"
	profile.PayloadDisplayName = "Enrollment Profile"
	profile.PayloadDescription = "The server may alter your settings"

	scepSubject := []string{
		[]string{
			[]string{"O", "MicroMDM"},
			[]string{"CN", "MDM Identity Certificate:UDID"},
		},
	}

	scepContent := SCEPPayload{
		Challenge: svc.SCEPChallenge,
		URL:       svc.SCEPUrl,
		Keysize:   1024,
		KeyType:   "RSA",
		KeyUsage:  0,
		Name:      "Device Management Identity Certificate",
		Subject:   scepSubject,
	}

	scepPayload := NewPayload("com.apple.security.scep")
	scepPayload.PayloadDescription = "Configures SCEP"
	scepPayload.PayloadDisplayName = "SCEP"
	scepPayload.PayloadContent = scepContent

	mdmPayload := MDMPayload{
		AccessRights:            8191,
		CheckInURL:              svc.Url + "/mdm/checkin",
		CheckOutWhenRemoved:     true,
		ServerURL:               svc.Url + "/mdm/connect",
		IdentityCertificateUUID: scepPayload.PayloadUUID,
		Topic: svc.Topic,
	}

	mdmPayload.PayloadVersion = 1
	mdmPayload.PayloadType = "com.apple.mdm"
	mdmPayload.PayloadDescription = "Enrolls with the MDM server"
	mdmPayload.PayloadOrganization = "MicroMDM"

	caPayload := NewPayload("com.apple.ssl.certificate")
	caPayload.PayloadDisplayName = "Root certificate for MicroMDM"
	caPayload.PayloadDescription = "Installs the root CA certificate for MicroMDM"

	append(profile.PayloadContent, scepPayload, mdmPayload, caPayload)

	return profile
}
