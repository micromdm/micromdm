package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

// Read-only: export SCEP CA cert+key from a MicroMDM BoltDB backup.
func main() {
	db, err := bolt.Open(os.Args[1], 0400, &bolt.Options{ReadOnly: true, Timeout: 5 * time.Second})
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("scep_certificates"))
		crt, key := b.Get([]byte("ca_certificate")), b.Get([]byte("ca_key"))
		c, err := x509.ParseCertificate(crt)
		if err != nil {
			panic(err)
		}
		k, err := x509.ParsePKCS1PrivateKey(key)
		if err != nil {
			panic(err)
		}
		os.WriteFile(os.Args[2]+"/scep-ca.pem", pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: crt}), 0600)
		os.WriteFile(os.Args[2]+"/scep-ca.key", pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: key}), 0600)
		fmt.Println("CA:", c.Subject, "notAfter:", c.NotAfter, "keybits:", k.N.BitLen())
		return nil
	})
}
