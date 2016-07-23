package enroll

import "io/ioutil"

type Service interface {
	Enroll() (Profile, error)
}

func NewService(pushCertPath string, pushCertPass string, caCertPath string, url string, scepUrl string) (Service, error) {
	pushTopic, err := GetPushTopicFromPKCS12(pushCertPath, pushCertPass)
	if err != nil {
		return nil, err
	}

	var caCert []byte

	if caCertPath != "" {
		caCert, err = ioutil.ReadFile(caCertPath)

		if err != nil {
			return nil, err
		}
	}

	scepSubject := [][][]string{
		[][]string{
			[]string{"O", "MicroMDM"},
			[]string{"CN", "MDM Identity Certificate:UDID"},
		},
	}

	return &service{
		Url:         url,
		SCEPUrl:     scepUrl,
		SCEPSubject: scepSubject,
		Topic:       pushTopic,
		CACert:      caCert,
	}, nil
}

type service struct {
	Url           string
	SCEPUrl       string
	SCEPChallenge string
	SCEPSubject   [][][]string
	Topic         string // APNS Topic for MDM notifications
	CACert        []byte
}

func (svc service) Enroll() (Profile, error) {
	profile := NewProfile()
	profile.PayloadIdentifier = "com.github.micromdm.micromdm.mdm"
	profile.PayloadOrganization = "MicroMDM"
	profile.PayloadDisplayName = "Enrollment Profile"
	profile.PayloadDescription = "The server may alter your settings"

	scepContent := SCEPPayloadContent{
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
	scepPayload.PayloadIdentifier = "com.github.micromdm.scep"
	scepPayload.PayloadOrganization = "MicroMDM"
	scepPayload.PayloadContent = scepContent

	mdmPayload := NewPayload("com.apple.mdm")
	mdmPayload.PayloadDescription = "Enrolls with the MDM server"
	mdmPayload.PayloadOrganization = "MicroMDM"
	mdmPayload.PayloadIdentifier = "com.github.micromdm.mdm"

	mdmPayloadContent := MDMPayloadContent{
		Payload:                 *mdmPayload,
		AccessRights:            8191,
		CheckInURL:              svc.Url + "/mdm/checkin",
		CheckOutWhenRemoved:     true,
		ServerURL:               svc.Url + "/mdm/connect",
		IdentityCertificateUUID: scepPayload.PayloadUUID,
		Topic: svc.Topic,
	}

	if len(svc.CACert) > 0 {
		caPayload := NewPayload("com.apple.ssl.certificate")
		caPayload.PayloadDisplayName = "Root certificate for MicroMDM"
		caPayload.PayloadDescription = "Installs the root CA certificate for MicroMDM"
		caPayload.PayloadContent = svc.CACert

		profile.PayloadContent = []interface{}{*scepPayload, mdmPayloadContent, *caPayload}
	} else {
		profile.PayloadContent = []interface{}{*scepPayload, mdmPayloadContent}
	}

	return *profile, nil
}
