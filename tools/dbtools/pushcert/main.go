package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/micromdm/micromdm/platform/config"
)

// Command pushcert is read-only: it extracts the current push cert+key from a BoltDB backup into PEM files.
func main() {
	db, err := bolt.Open(os.Args[1], 0400, &bolt.Options{ReadOnly: true, Timeout: 5 * time.Second})
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var conf config.ServerConfig
	err = db.View(func(tx *bolt.Tx) error {
		return config.UnmarshalServerConfig(tx.Bucket([]byte("mdm.ServerConfig")).Get([]byte("config")), &conf)
	})
	if err != nil {
		panic(err)
	}
	der := conf.PushCertificate
	if b, _ := pem.Decode(der); b != nil {
		der = b.Bytes
	}
	c, err := x509.ParseCertificate(der)
	if err != nil {
		panic(err)
	}
	fmt.Println("cert:", c.Subject.String(), "notAfter:", c.NotAfter)
	os.WriteFile(os.Args[2]+".pem", pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), 0600)
	blk, _ := pem.Decode(conf.PrivateKey)
	if blk != nil {
		os.WriteFile(os.Args[2]+".key", conf.PrivateKey, 0600)
	} else {
		os.WriteFile(os.Args[2]+".key", pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: conf.PrivateKey}), 0600)
	}
}
