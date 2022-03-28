package builtin

import (
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/micromdm/micromdm/pkg/crypto"
	"github.com/micromdm/micromdm/platform/config"
)

const (
	depTokenBucket = "mdm.DEPToken"
)

func (db *DB) AddDEPToken(consumerKey string, json []byte) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(depTokenBucket))
		if err != nil {
			return err
		}
		err = b.Put([]byte("last_added"), []byte(consumerKey))
		if err != nil {
			return err
		}
		return b.Put([]byte(consumerKey), json)
	})
	if err != nil {
		return err
	}
	err = db.Publisher.Publish(context.TODO(), config.DEPTokenTopic, json)
	return err
}

func (db *DB) DEPTokens() ([]config.DEPToken, error) {
	var result []config.DEPToken
	var lastAdded string
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(depTokenBucket))
		if b == nil {
			return nil
		}
		lastAdded = string(b.Get([]byte("last_added")))
		c := b.Cursor()

		prefix := []byte("CK_")
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			var depToken config.DEPToken
			err := json.Unmarshal(v, &depToken)
			if err != nil {
				// TODO: log problematic DEP token, or remove altogether?
				continue
			}
			if lastAdded == depToken.ConsumerKey {
				// the server merely takes the first DEP token. let's
				// make sure the most recent one is the one that's
				// returned first
				result = append([]config.DEPToken{depToken}, result...)
			} else {
				result = append(result, depToken)
			}
		}
		return nil
	})
	return result, err
}

func (db *DB) DEPKeypair() (key *rsa.PrivateKey, cert *x509.Certificate, err error) {
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

func generateAndStoreDEPKeypair(db *DB) (key *rsa.PrivateKey, cert *x509.Certificate, err error) {
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
