package list

import (
	"bytes"
	"context"
	"encoding/json"

	"crypto/rsa"
	"crypto/x509"

	"github.com/boltdb/bolt"

	"github.com/micromdm/micromdm/blueprint"
	"github.com/micromdm/micromdm/crypto"
	"github.com/micromdm/micromdm/device"
	"github.com/micromdm/micromdm/profile"
)

type ListDevicesOption struct {
	Page    int
	PerPage int

	FilterSerial []string
	FilterUDID   []string
}

type GetBlueprintsOption struct {
	FilterName string
}

type GetProfilesOption struct {
	Identifier string `json:"id"`
}

type Service interface {
	ListDevices(ctx context.Context, opt ListDevicesOption) ([]DeviceDTO, error)
	GetDEPTokens(ctx context.Context) ([]DEPToken, []byte, error)
	GetBlueprints(ctx context.Context, opt GetBlueprintsOption) ([]blueprint.Blueprint, error)
	GetProfiles(ctx context.Context, opt GetProfilesOption) ([]profile.Profile, error)
}

type ListService struct {
	Devices    *device.DB
	Blueprints *blueprint.DB
	Profiles   *profile.DB
	DB         *bolt.DB // TODO: replace with reference to DEP token svc/pkg
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
func GetDEPTokens(db *bolt.DB) ([]DEPToken, error) {
	var result []DEPToken
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(depTokenBucket))
		if b == nil {
			return nil
		}
		c := b.Cursor()

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
	key, cert, err = crypto.SimpleSelfSignedRSAKeypair("micromdm-dep-token", 365)
	if err != nil {
		return
	}

	pkBytes := x509.MarshalPKCS1PrivateKey(key)
	certBytes := cert.Raw

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(depTokenBucket))
		if err != nil {
			return err
		}
		err = b.Put([]byte("key"), pkBytes)
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
		if b == nil {
			return nil
		}
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

	tokens, err := GetDEPTokens(svc.DB)
	if err != nil {
		return nil, certBytes, err
	}

	return tokens, certBytes, nil
}

func (svc *ListService) GetBlueprints(ctx context.Context, opt GetBlueprintsOption) ([]blueprint.Blueprint, error) {
	if opt.FilterName != "" {
		bp, err := svc.Blueprints.BlueprintByName(opt.FilterName)
		if err != nil {
			return nil, err
		}
		return []blueprint.Blueprint{*bp}, err
	} else {
		bps, err := svc.Blueprints.List()
		if err != nil {
			return nil, err
		}
		return bps, nil
	}
}

func (svc *ListService) GetProfiles(ctx context.Context, opt GetProfilesOption) ([]profile.Profile, error) {
	if opt.Identifier != "" {
		foundProf, err := svc.Profiles.ProfileById(opt.Identifier)
		if err != nil {
			return nil, err
		}
		return []profile.Profile{*foundProf}, nil
	} else {
		return svc.Profiles.List()
	}
}
