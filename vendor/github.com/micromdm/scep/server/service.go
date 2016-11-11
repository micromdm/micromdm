package scepserver

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/asn1"
	"errors"
	"math/big"
	"time"

	"github.com/micromdm/scep/depot"
	"github.com/micromdm/scep/scep"
	"golang.org/x/net/context"
)

// Service is the interface for all supported SCEP server operations.
type Service interface {
	// GetCACaps returns a list of options
	// which are supported by the server.
	GetCACaps(ctx context.Context) ([]byte, error)

	// GetCACert returns CA certificate or
	// a CA certificate chain with intermediates
	// in a PKCS#7 Degenerate Certificates format
	GetCACert(ctx context.Context) ([]byte, int, error)

	// PKIOperation handles incoming SCEP messages such as PKCSReq and
	// sends back a CertRep PKIMessag.
	PKIOperation(ctx context.Context, msg []byte) ([]byte, error)

	// GetNextCACert returns a replacement certificate or certificate chain
	// when the old one expires. The response format is a PKCS#7 Degenerate
	// Certificates type.
	GetNextCACert(ctx context.Context) ([]byte, error)
}

type service struct {
	depot             depot.Depot
	ca                []*x509.Certificate // CA cert or chain
	caKey             *rsa.PrivateKey
	caKeyPassword     []byte
	csrTemplate       *x509.Certificate
	challengePassword string
	allowRenewal      int // days before expiry, 0 to disable
	clientValidity    int // client cert validity in days
}

func (svc service) GetCACaps(ctx context.Context) ([]byte, error) {
	defaultCaps := []byte(`POSTPKIOperation`)
	return defaultCaps, nil
}

func (svc service) GetCACert(ctx context.Context) ([]byte, int, error) {
	if len(svc.ca) == 0 {
		return nil, 0, errors.New("missing CA Cert")
	}
	if len(svc.ca) == 1 {
		return svc.ca[0].Raw, 1, nil
	}
	data, err := scep.DegenerateCertificates(svc.ca)
	return data, len(svc.ca), err
}

func (svc service) PKIOperation(ctx context.Context, data []byte) ([]byte, error) {
	msg, err := scep.ParsePKIMessage(data)
	if err != nil {
		// handle err
		return nil, err
	}
	ca := svc.ca[0]
	if err := msg.DecryptPKIEnvelope(svc.ca[0], svc.caKey); err != nil {
		return nil, err
	}

	// validate challenge passwords
	if msg.MessageType == scep.PKCSReq {
		if !svc.challengePasswordMatch(msg.CSRReqMessage.ChallengePassword) {
			// handle err
			return nil, errors.New("scep challenge password does not match")
		}
	}

	csr := msg.CSRReqMessage.CSR
	id, err := generateSubjectKeyID(csr.PublicKey)
	if err != nil {
		return nil, err
	}

	serial, err := svc.depot.Serial()
	if err != nil {
		return nil, err
	}

	duration := svc.clientValidity

	// create cert template
	tmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject:      csr.Subject,
		NotBefore:    time.Now().Add(-600).UTC(),
		NotAfter:     time.Now().AddDate(0, 0, duration).UTC(),
		SubjectKeyId: id,
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
		},
	}

	certRep, err := msg.SignCSR(ca, svc.caKey, tmpl)
	if err != nil {
		return nil, err
	}

	crt := certRep.CertRepMessage.Certificate
	name := certName(crt)

	// Test if this certificate is already in the CADB, revoke if needed
	// revocation is done if the validity of the existing certificate is
	// less than allowRenewal (14 days by default)
	err = svc.depot.HasCN(name, svc.allowRenewal, crt, false)
	if err != nil {
		return nil, err
	}

	if err := svc.depot.Put(name, crt); err != nil {
		return nil, err
	}

	return certRep.Raw, nil

}

func certName(crt *x509.Certificate) string {
	if crt.Subject.CommonName != "" {
		return crt.Subject.CommonName
	}
	return string(crt.Signature)
}

func (svc service) GetNextCACert(ctx context.Context) ([]byte, error) {
	panic("not implemented")
}

func (svc service) challengePasswordMatch(pw string) bool {
	if svc.challengePassword == "" {
		// empty password, don't validate
		return true
	}
	if svc.challengePassword == pw {
		return true
	}
	return false
}

// ServiceOption is a server configuration option
type ServiceOption func(*service) error

// ChallengePassword is an optional argument to NewService
// which allows setting a preshared key for SCEP.
func ChallengePassword(pw string) ServiceOption {
	return func(s *service) error {
		s.challengePassword = pw
		return nil
	}
}

// CAKeyPassword is an optional argument to NewService for
// specifying the CA private key password.
func CAKeyPassword(pw []byte) ServiceOption {
	return func(s *service) error {
		s.caKeyPassword = pw
		return nil
	}
}

// allowRenewal sets the days before expiry which we are allowed to renew (optional)
func AllowRenewal(duration int) ServiceOption {
	return func(s *service) error {
		s.allowRenewal = duration
		return nil
	}
}

// ClientValidity sets the validity of signed client certs in days (optional parameter)
func ClientValidity(duration int) ServiceOption {
	return func(s *service) error {
		s.clientValidity = duration
		return nil
	}
}

// NewService creates a new scep service
func NewService(depot depot.Depot, opts ...ServiceOption) (Service, error) {
	s := &service{
		depot: depot,
	}
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	var err error
	s.ca, s.caKey, err = depot.CA(s.caKeyPassword)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// rsaPublicKey reflects the ASN.1 structure of a PKCS#1 public key.
type rsaPublicKey struct {
	N *big.Int
	E int
}

// GenerateSubjectKeyID generates SubjectKeyId used in Certificate
// ID is 160-bit SHA-1 hash of the value of the BIT STRING subjectPublicKey
func generateSubjectKeyID(pub crypto.PublicKey) ([]byte, error) {
	var pubBytes []byte
	var err error
	switch pub := pub.(type) {
	case *rsa.PublicKey:
		pubBytes, err = asn1.Marshal(rsaPublicKey{
			N: pub.N,
			E: pub.E,
		})
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("only RSA public key is supported")
	}

	hash := sha1.Sum(pubBytes)

	return hash[:], nil
}
