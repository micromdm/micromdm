package enroll

import (
	"bytes"
	"crypto/x509"
	"io/ioutil"
	"strings"

	"github.com/groob/plist"
	"golang.org/x/net/context"

	"github.com/micromdm/micromdm/profile"
)

const EnrollmentProfileId string = "com.github.micromdm.micromdm.mdm"

type Service interface {
	Enroll(ctx context.Context) (profile.Mobileconfig, error)
	OTAEnroll(ctx context.Context) (Payload, error)
	OTAPhase2(ctx context.Context) (Profile, error)
	OTAPhase3(ctx context.Context) (Profile, error)
}

// profileToProfile translates an enroll.Profile struct into a profile.Profile struct.
// Note: converts the Profile to a property list.
func profileToProfile(in Profile) (*profile.Profile, error) {
	out := new(profile.Profile)
	buf := new(bytes.Buffer)

	enc := plist.NewEncoder(buf)
	enc.Indent("  ")
	err := enc.Encode(in)
	if err != nil {
		return out, err
	}

	out.Mobileconfig = buf.Bytes()
	out.Identifier, err = out.Mobileconfig.GetPayloadIdentifier()
	if err != nil {
		return out, err
	}

	return out, err
}

func NewService(pushTopic, caCertPath, scepURL, scepChallenge, url, tlsCertPath, scepSubject string, profileDB *profile.DB) (Service, error) {
	var caCert, tlsCert []byte
	var err error

	if caCertPath != "" {
		caCert, err = ioutil.ReadFile(caCertPath)

		if err != nil {
			return nil, err
		}
	}

	if tlsCertPath != "" {
		tlsCert, err = ioutil.ReadFile(tlsCertPath)

		if err != nil {
			return nil, err
		}
	}

	if scepSubject == "" {
		scepSubject = "/O=MicroMDM/CN=MicroMDM Identity (%ComputerName%)"
	}

	subjectElements := strings.Split(scepSubject, "/")
	var subject [][][]string

	for _, element := range subjectElements {
		if element == "" {
			continue
		}
		subjectKeyValue := strings.Split(element, "=")
		subject = append(subject, [][]string{[]string{subjectKeyValue[0], subjectKeyValue[1]}})
	}

	svc := &service{
		URL:           url,
		SCEPURL:       scepURL,
		SCEPSubject:   subject,
		SCEPChallenge: scepChallenge,
		Topic:         pushTopic,
		CACert:        caCert,
		TLSCert:       tlsCert,
		ProfileDB:     profileDB,
	}

	// look for existing enrollment profile, and save a new one if not found
	_, err = svc.ProfileDB.ProfileById(EnrollmentProfileId)
	if err != nil && profile.IsNotFound(err) {
		enrollmentProfile, err := svc.MakeEnrollmentProfile()
		if err != nil {
			return svc, err
		}

		profile, err := profileToProfile(enrollmentProfile)
		if err != nil {
			return svc, err
		}

		err = svc.ProfileDB.Save(profile)
		if err != nil {
			return svc, err
		}
	}

	return svc, err
}

type service struct {
	URL           string
	SCEPURL       string
	SCEPChallenge string
	SCEPSubject   [][][]string
	Topic         string // APNS Topic for MDM notifications
	CACert        []byte
	TLSCert       []byte
	ProfileDB     *profile.DB
}

func (svc service) Enroll(ctx context.Context) (profile.Mobileconfig, error) {
	// NOTE: this effectively makes the enrollment profile static. that is it
	// only gets generated once on service start _only if_ one doesn't already
	// exist. when we integrate dynamic enrollment profile generation (for
	// e.g. per-device SCEP challenges) we'll need to modify this enrollment
	// profile before delivery. treating it much like a template (when
	// templated profiles get implemented)
	profile, err := svc.ProfileDB.ProfileById(EnrollmentProfileId)
	if err != nil {
		return nil, err
	}
	return profile.Mobileconfig, nil
}

func (svc service) MakeEnrollmentProfile() (Profile, error) {
	profile := NewProfile()
	profile.PayloadIdentifier = EnrollmentProfileId
	profile.PayloadOrganization = "MicroMDM"
	profile.PayloadDisplayName = "Enrollment Profile"
	profile.PayloadDescription = "The server may alter your settings"
	profile.PayloadScope = "System"

	mdmPayload := NewPayload("com.apple.mdm")
	mdmPayload.PayloadDescription = "Enrolls with the MDM server"
	mdmPayload.PayloadOrganization = "MicroMDM"
	mdmPayload.PayloadIdentifier = "com.github.micromdm.mdm"
	mdmPayload.PayloadScope = "System"

	mdmPayloadContent := MDMPayloadContent{
		Payload:             *mdmPayload,
		AccessRights:        8191,
		CheckInURL:          svc.URL + "/mdm/checkin",
		CheckOutWhenRemoved: true,
		ServerURL:           svc.URL + "/mdm/connect",
		Topic:               svc.Topic,
		SignMessage:         true,
	}

	payloadContent := []interface{}{}

	if svc.SCEPURL != "" {
		scepContent := SCEPPayloadContent{
			URL:      svc.SCEPURL,
			Keysize:  2048,
			KeyType:  "RSA",
			KeyUsage: int(x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment),
			Name:     "Device Management Identity Certificate",
			Subject:  svc.SCEPSubject,
		}

		if svc.SCEPChallenge != "" {
			scepContent.Challenge = svc.SCEPChallenge
		}

		scepPayload := NewPayload("com.apple.security.scep")
		scepPayload.PayloadDescription = "Configures SCEP"
		scepPayload.PayloadDisplayName = "SCEP"
		scepPayload.PayloadIdentifier = "com.github.micromdm.scep"
		scepPayload.PayloadOrganization = "MicroMDM"
		scepPayload.PayloadContent = scepContent
		scepPayload.PayloadScope = "System"

		payloadContent = append(payloadContent, *scepPayload)
		mdmPayloadContent.IdentityCertificateUUID = scepPayload.PayloadUUID
	}

	payloadContent = append(payloadContent, mdmPayloadContent)

	if len(svc.CACert) > 0 {
		caPayload := NewPayload("com.apple.security.root")
		caPayload.PayloadDisplayName = "Root certificate for MicroMDM"
		caPayload.PayloadDescription = "Installs the root CA certificate for MicroMDM"
		caPayload.PayloadIdentifier = "com.github.micromdm.ssl.ca"
		caPayload.PayloadContent = svc.CACert

		payloadContent = append(payloadContent, *caPayload)
	}

	// Client needs to trust us at this point if we are using a self signed certificate.
	if len(svc.TLSCert) > 0 {
		tlsPayload := NewPayload("com.apple.security.pem")
		tlsPayload.PayloadDisplayName = "Self-signed TLS certificate for MicroMDM"
		tlsPayload.PayloadDescription = "Installs the TLS certificate for MicroMDM"
		tlsPayload.PayloadIdentifier = "com.github.micromdm.tls"
		tlsPayload.PayloadContent = svc.TLSCert

		payloadContent = append(payloadContent, *tlsPayload)
	}

	profile.PayloadContent = payloadContent

	return *profile, nil
}

// OTAEnroll returns an Over-the-Air "Profile Service" Payload for enrollment.
func (svc service) OTAEnroll(ctx context.Context) (Payload, error) {
	payload := NewPayload("Profile Service")
	payload.PayloadIdentifier = "com.github.micromdm.micromdm.ota"
	payload.PayloadDisplayName = "MicroMDM Profile Service"
	payload.PayloadDescription = "Profile Service enrollment"
	payload.PayloadOrganization = "MicroMDM"
	payload.PayloadContent = ProfileServicePayload{
		URL:              svc.URL + "/ota/phase23",
		Challenge:        "",
		DeviceAttributes: []string{"UDID", "VERSION", "PRODUCT", "SERIAL", "MEID", "IMEI"},
	}

	// yes, this is a bare Payload, not a Profile
	return *payload, nil
}

// OTAPhase2 returns a SCEP Profile for use in phase 2 of Over-the-Air enrollment.
func (svc service) OTAPhase2(ctx context.Context) (Profile, error) {
	profile := NewProfile()
	profile.PayloadIdentifier = "com.github.micromdm.micromdm.ota.phase2"
	profile.PayloadOrganization = "MicroMDM"
	profile.PayloadDisplayName = "OTA Phase 2"
	profile.PayloadScope = "System"

	scepContent := SCEPPayloadContent{
		URL: svc.SCEPURL,

		// TODO: OTA spec says: "In general, 1024-bit keys are recommended
		// because of the overhead involved in generating 2048-bit keys." We'd
		// prefer to use 2048, but deferring for now.
		Keysize: 1024,

		KeyType: "RSA",

		// OTA spec says: "Note that the MS SCEP server will only issue
		// signature or encryption, not both."
		KeyUsage: int(x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment),

		Name:    "OTA Phase 2 Certificate",
		Subject: svc.SCEPSubject,
	}

	if svc.SCEPChallenge != "" {
		scepContent.Challenge = svc.SCEPChallenge
	}

	scepPayload := NewPayload("com.apple.security.scep")
	scepPayload.PayloadDescription = "Configures SCEP"
	scepPayload.PayloadDisplayName = "SCEP"
	scepPayload.PayloadIdentifier = "com.github.micromdm.micromdm.ota.phase2.scep"
	scepPayload.PayloadOrganization = "MicroMDM"
	scepPayload.PayloadContent = scepContent
	scepPayload.PayloadScope = "System"

	profile.PayloadContent = append(profile.PayloadContent, *scepPayload)

	return *profile, nil
}

// OTAPhase3 returns a Profile for use in phase 3 of Over-the-Air profile enrollment.
// This would typically be the final or end profile of the Over-the-Air
// enrollment process. In our case this would probably be a device-specifc
// MDM enrollment payload.
// TODO: Not implemented.
func (svc service) OTAPhase3(ctx context.Context) (Profile, error) {
	return Profile{}, nil
}
