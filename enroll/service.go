package enroll

type Service interface {
	Enroll() (Profile, error)
}

func NewService() Service {
	scepSubject := []string{
		[]string{
			[]string{"O", "MicroMDM"},
			[]string{"CN", "MDM Identity Certificate:UDID"},
		},
	}

	return &service{
		SCEPSubject: scepSubject,
	}
}

type service struct {
	Url           string
	SCEPUrl       string
	SCEPChallenge string
	SCEPSubject   [][][]string
	Topic         string // APNS Topic for MDM notifications
}

func (svc service) Enroll() (Profile, error) {
	profile := NewProfile()
	profile.PayloadIdentifier = "com.github.micromdm.micromdm.mdm"
	profile.PayloadOrganization = "MicroMDM"
	profile.PayloadDisplayName = "Enrollment Profile"
	profile.PayloadDescription = "The server may alter your settings"

	scepContent := SCEPPayload{
		Challenge: svc.SCEPChallenge,
		URL:       svc.SCEPUrl,
		Keysize:   1024,
		KeyType:   "RSA",
		KeyUsage:  0,
		Name:      "Device Management Identity Certificate",
		Subject:   svc.SCEPSubject,
	}

	scepPayload := NewPayload("com.apple.security.scep")
	scepPayload.PayloadDescription = "Configures SCEP"
	scepPayload.PayloadDisplayName = "SCEP"
	scepPayload.PayloadContent = scepContent

	mdmPayload := MDMPayload{
		Payload{
			PayloadVersion:      1,
			PayloadType:         "com.apple.mdm",
			PayloadDescription:  "Enrolls with the MDM server",
			PayloadOrganization: "MicroMDM",
		},
		AccessRights:            8191,
		CheckInURL:              svc.Url + "/mdm/checkin",
		CheckOutWhenRemoved:     true,
		ServerURL:               svc.Url + "/mdm/connect",
		IdentityCertificateUUID: scepPayload.PayloadUUID,
		Topic: svc.Topic,
	}

	caPayload := NewPayload("com.apple.ssl.certificate")
	caPayload.PayloadDisplayName = "Root certificate for MicroMDM"
	caPayload.PayloadDescription = "Installs the root CA certificate for MicroMDM"

	append(profile.PayloadContent, scepPayload, mdmPayload, caPayload)

	return profile
}
