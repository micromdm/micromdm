package bolt

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/boltdb/bolt"
)

// Depot implements a SCEP certifiacte store using boltdb.
// https://github.com/boltdb/bolt
type Depot struct {
	*bolt.DB
}

const (
	certBucket = "scep_certificates"
)

// NewBoltDepot creates a depot.Depot backed by BoltDB.
func NewBoltDepot(db *bolt.DB) (*Depot, error) {
	fmt.Println("--- NEW Bolt Depot for SCEP")
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(certBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &Depot{db}, nil
}

func (db *Depot) CA(pass []byte) ([]*x509.Certificate, *rsa.PrivateKey, error) {
	chain := []*x509.Certificate{}
	var key *rsa.PrivateKey
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(certBucket))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found!", certBucket)
		}
		// get ca_certificate
		caCert := bucket.Get([]byte("ca_certificate"))
		if caCert == nil {
			return fmt.Errorf("no ca_certificate in bucket")
		}
		// we need to make a copy of the byte slice because the asn.Unmarshal
		// method called by ParseCertificate will retain a reference to the original.
		// The slice should no longer be referenced once the BoltDB transaction is closed.
		caCertBytes := append([]byte(nil), caCert...)
		cert, err := x509.ParseCertificate(caCertBytes)
		if err != nil {
			return err
		}
		chain = append(chain, cert)

		// get ca_key
		caKey := bucket.Get([]byte("ca_key"))
		if caKey == nil {
			return fmt.Errorf("no ca_key in bucket")
		}
		key, err = x509.ParsePKCS1PrivateKey(caKey)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return chain, key, nil
}

func (db *Depot) Put(cn string, crt *x509.Certificate) error {
	if crt == nil || crt.Raw == nil {
		return fmt.Errorf("%q does not specify a valid certificate for storage", cn)
	}
	serial, err := db.Serial()
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(certBucket))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found!", certBucket)
		}
		name := cn + "." + serial.String()
		fmt.Println("--- scep/builtin/db.go Put")
		fmt.Println(name)
		fmt.Println(crt.Raw)
		return bucket.Put([]byte(name), crt.Raw)
	})
	if err != nil {
		return err
	}
	return db.incrementSerial(serial)
}

func (db *Depot) Serial() (*big.Int, error) {
	fmt.Println("--- scep/builtin/db.go Serial()")
		
	s := big.NewInt(2)
	if !db.hasKey([]byte("serial")) {
		if err := db.writeSerial(s); err != nil {
			return nil, err
		}
		fmt.Println(s)
		return s, nil
	}
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(certBucket))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found!", certBucket)
		}
		k := bucket.Get([]byte("serial"))
		if k == nil {
			return fmt.Errorf("key %q not found", "serial")
		}
		s = s.SetBytes(k)
		fmt.Println(s)
		return nil
	})
	if err != nil {
		return nil, err
	}
	fmt.Println(s)
	return s, nil
}

func (db *Depot) writeSerial(s *big.Int) error {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(certBucket))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found!", certBucket)
		}
		fmt.Println("--- scep/builtin/db.go writeSerial")
		fmt.Println("serial")
		fmt.Println(s.Bytes())
		return bucket.Put([]byte("serial"), []byte(s.Bytes()))
	})
	return err
}

func (db *Depot) hasKey(name []byte) bool {
	fmt.Println("--- scep/builtin/db.go hasKey")
	fmt.Println(name)
	
	var present bool
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(certBucket))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found!", certBucket)
		}
		k := bucket.Get([]byte("serial"))
		if k != nil {
			present = true
		}
		return nil
	})
	return present
}

func (db *Depot) incrementSerial(s *big.Int) error {
	serial := s.Add(s, big.NewInt(1))
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(certBucket))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found!", certBucket)
		}
		fmt.Println("--- scep/builtin/db.go incrementSerial")
		fmt.Println("serial")
		fmt.Println(serial.Bytes())
		return bucket.Put([]byte("serial"), []byte(serial.Bytes()))
	})
	return err
}

func (db *Depot) HasCN(cn string, allowTime int, cert *x509.Certificate, revokeOldCertificate bool) (bool, error) {
	fmt.Println("--- scep/builtin/db.go HasCN")
	fmt.Println(cn)
	fmt.Println(cert.Subject.CommonName)
	
	// TODO: implement allowTime
	// TODO: implement revocation
	if cert == nil {
		return false, errors.New("nil certificate provided")
	}
	var hasCN bool
	err := db.View(func(tx *bolt.Tx) error {
		// TODO: "scep_certificates" is internal const in micromdm/scep
		curs := tx.Bucket([]byte("scep_certificates")).Cursor()
		prefix := []byte(cert.Subject.CommonName)
		for k, v := curs.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = curs.Next() {
			if bytes.Compare(v, cert.Raw) == 0 {
				hasCN = true
				return nil
			}
		}

		return nil
	})
	return hasCN, err
}

func (db *Depot) CreateOrLoadKey(bits int) (*rsa.PrivateKey, error) {
	var (
		key *rsa.PrivateKey
		err error
	)
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(certBucket))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found!", certBucket)
		}
		priv := bucket.Get([]byte("ca_key"))
		if priv == nil {
			return nil
		}
		key, err = x509.ParsePKCS1PrivateKey(priv)
		return err
	})
	if err != nil {
		return nil, err
	}
	if key != nil {
		return key, nil
	}
	key, err = rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(certBucket))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found!", certBucket)
		}
		fmt.Println("--- scep/builtin/db.go CreateOrLoadKey")
		fmt.Println("ca_key")
		fmt.Println(x509.MarshalPKCS1PrivateKey(key))
		return bucket.Put([]byte("ca_key"), x509.MarshalPKCS1PrivateKey(key))
	})
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (db *Depot) CreateOrLoadCA(key *rsa.PrivateKey, years int, org, country string) (*x509.Certificate, error) {
	var (
		cert *x509.Certificate
		err  error
	)
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(certBucket))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found!", certBucket)
		}
		caCert := bucket.Get([]byte("ca_certificate"))
		if caCert == nil {
			return nil
		}
		cert, err = x509.ParseCertificate(caCert)
		return err
	})
	if err != nil {
		return nil, err
	}
	if cert != nil {
		return cert, nil
	}

	subject := pkix.Name{
		Country:            []string{country},
		Organization:       []string{org},
		OrganizationalUnit: []string{"MICROMDM SCEP CA"},
		Locality:           nil,
		Province:           nil,
		StreetAddress:      nil,
		PostalCode:         nil,
		SerialNumber:       "",
		CommonName:         org,
	}

	subjectKeyID, err := generateSubjectKeyID(&key.PublicKey)
	if err != nil {
		return nil, err
	}

	authTemplate := x509.Certificate{
		SerialNumber:       big.NewInt(1),
		Subject:            subject,
		NotBefore:          time.Now().Add(-600).UTC(),
		NotAfter:           time.Now().AddDate(years, 0, 0).UTC(),
		KeyUsage:           x509.KeyUsageCertSign,
		ExtKeyUsage:        nil,
		UnknownExtKeyUsage: nil,

		BasicConstraintsValid: true,
		IsCA:                        true,
		MaxPathLen:                  0,
		SubjectKeyId:                subjectKeyID,
		DNSNames:                    nil,
		PermittedDNSDomainsCritical: false,
		PermittedDNSDomains:         nil,
	}

	crtBytes, err := x509.CreateCertificate(rand.Reader, &authTemplate, &authTemplate, &key.PublicKey, key)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(certBucket))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found!", certBucket)
		}
		fmt.Println("--- 	scep/builtin/db.go CreateOrLoadKey")
		fmt.Println("ca_certificate")
		fmt.Println(crtBytes)
		return bucket.Put([]byte("ca_certificate"), crtBytes)
	})
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(crtBytes)
}

// rsaPublicKey reflects the ASN.1 structure of a PKCS#1 public key.
type rsaPublicKey struct {
	N *big.Int
	E int
}

// GenerateSubjectKeyID generates SubjectKeyId used in Certificate
// ID is 160-bit SHA-1 hash of the value of the BIT STRING subjectPublicKey
func generateSubjectKeyID(pub crypto.PublicKey) ([]byte, error) {
	fmt.Println("--- scep/builtin/db.go generateSubjectKeyID")
	fmt.Println(pub)
	
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