package bolt

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"math/big"

	"github.com/boltdb/bolt"
)

// Depot implements a SCEP certifiacte store using boltdb.
// https://github.com/boltdb/bolt
type Depot struct {
	*bolt.DB
}

const certBucket = "scep_certificates"

// NewBoltDepot creates a depot.Depot backed by BoltDB.
func NewBoltDepot(db *bolt.DB) (*Depot, error) {
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
		cert, err := x509.ParseCertificate(caCert)
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
	return chain, key, err
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
		return bucket.Put([]byte(name), crt.Raw)
	})
	if err != nil {
		return err
	}
	return db.incrementSerial(serial)
}

func (db *Depot) Serial() (*big.Int, error) {
	s := big.NewInt(2)
	if !db.hasKey([]byte("serial")) {
		if err := db.writeSerial(s); err != nil {
			return nil, err
		}
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
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (db *Depot) writeSerial(s *big.Int) error {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(certBucket))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found!", certBucket)
		}
		return bucket.Put([]byte("serial"), []byte(s.Bytes()))
	})
	return err
}

func (db *Depot) hasKey(name []byte) bool {
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
		return bucket.Put([]byte("serial"), []byte(serial.Bytes()))
	})
	return err
}

func (db *Depot) HasCN(cn string, allowTime int, cert *x509.Certificate, revokeOldCertificate bool) error {
	// FIXME: not implemented.
	return nil
}
