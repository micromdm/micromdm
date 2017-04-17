package list

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"

	"github.com/boltdb/bolt"

	"github.com/micromdm/micromdm/device"
)

type ListDevicesOption struct {
	Page    int
	PerPage int

	FilterSerial []string
	FilterUDID   []string
}

type Service interface {
	ListDevices(ctx context.Context, opt ListDevicesOption) ([]DeviceDTO, error)
	GetDEPTokens(ctx context.Context) ([]DEPToken, []byte, error)
}

type ListService struct {
	Devices *device.DB
	DB      *bolt.DB // TODO: replace with reference to DEP token svc/pkg
}

func (svc *ListService) ListDevices(ctx context.Context, opt ListDevicesOption) ([]DeviceDTO, error) {
	devices, err := svc.Devices.List()
	dto := []DeviceDTO{}
	for _, d := range devices {
		dto = append(dto, DeviceDTO{
			SerialNumber:     d.SerialNumber,
			UDID:             d.UDID,
			EnrollmentStatus: d.Enrolled,
			LastSeen:         d.LastCheckin,
		})
	}
	return dto, err
}

// TODO: move into seperate svc/pkg
const (
	depTokenBucket = "mdm.DEPToken"
)

// TODO: move into seperate svc/pkg
func getDEPTokens(db *bolt.DB) ([]DEPToken, error) {
	var result []DEPToken
	err := db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(depTokenBucket)).Cursor()

		prefix := []byte("CK_")
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			var depToken DEPToken
			err := json.Unmarshal(v, &depToken)
			if err != nil {
				// TODO: log problematic DEP token, or remove altogether?
				continue
			}
			result = append(result, depToken)
		}
		return nil
	})
	return result, err
}

// TODO: move into seperate svc/pkg
func generateAndStoreDEPKeypair(db *bolt.DB) (key *rsa.PrivateKey, cert *x509.Certificate, err error) {
	key, cert, err = SimpleSelfSignedRSAKeypair("micromdm-dep-token", 365)
	if err != nil {
		return
	}

	pkBytes := x509.MarshalPKCS1PrivateKey(key)
	certBytes := cert.Raw

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(depTokenBucket))
		err := b.Put([]byte("key"), pkBytes)
		if err != nil {
			return err
		}
		err = b.Put([]byte("certificate"), certBytes)
		if err != nil {
			return err
		}
		return nil
	})

	return
}

// TODO: move into seperate svc/pkg
func GetDEPKeypair(db *bolt.DB) (key *rsa.PrivateKey, cert *x509.Certificate, err error) {
	var keyBytes, certBytes []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(depTokenBucket))
		keyBytes = b.Get([]byte("key"))
		certBytes = b.Get([]byte("certificate"))
		return nil
	})
	if err != nil {
		return
	}
	if keyBytes == nil || certBytes == nil {
		// if there is no certificate or private key then generate
		key, cert, err = generateAndStoreDEPKeypair(db)
	} else {
		key, err = x509.ParsePKCS1PrivateKey(keyBytes)
		if err != nil {
			return
		}
		cert, err = x509.ParseCertificate(certBytes)
		if err != nil {
			return
		}
	}
	return
}

func (svc *ListService) GetDEPTokens(ctx context.Context) ([]DEPToken, []byte, error) {
	_, cert, err := GetDEPKeypair(svc.DB)
	if err != nil {
		return nil, nil, err
	}
	var certBytes []byte
	if cert != nil {
		certBytes = cert.Raw
	}

	tokens, err := getDEPTokens(svc.DB)
	if err != nil {
		return nil, certBytes, err
	}

	return tokens, certBytes, nil
}

// TODO: move into crypto package
func RandomCertificateSerialNumber() (*big.Int, error) {
	limit := new(big.Int).Lsh(big.NewInt(1), 128)
	return rand.Int(rand.Reader, limit)
}

// TODO: move into crypto package
func SimpleSelfSignedRSAKeypair(cn string, days int) (key *rsa.PrivateKey, cert *x509.Certificate, err error) {
	key, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return
	}

	serialNumber, err := RandomCertificateSerialNumber()
	if err != nil {
		return
	}
	timeNow := time.Now()
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: cn,
		},
		NotBefore:             timeNow,
		NotAfter:              timeNow.Add(time.Duration(days) * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return
	}
	cert, err = x509.ParseCertificate(certBytes)
	if err != nil {
		return
	}

	return
}
