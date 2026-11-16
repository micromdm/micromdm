package main

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/boltdb/bolt"
	"github.com/micromdm/micromdm/platform/device"
)

func main() {
	db, err := bolt.Open(os.Args[1], 0400, &bolt.Options{ReadOnly: true, Timeout: 5 * time.Second})
	if err != nil {
		panic(err)
	}
	defer db.Close()
	now := time.Now()
	certs := map[string]*x509.Certificate{}
	type row struct {
		udid, serial, model, lastSeen string
		enrolled                      bool
		days, queued                  int
		notAfter                      string
		left                          int
	}
	var rows []row
	err = db.View(func(tx *bolt.Tx) error {
		tx.Bucket([]byte("scep_certificates")).ForEach(func(k, v []byte) error {
			if c, err := x509.ParseCertificate(v); err == nil {
				h := sha256.Sum256(v)
				certs[string(h[:])] = c
			}
			return nil
		})
		auth := tx.Bucket([]byte("mdm.UDIDCertAuth"))
		q := tx.Bucket([]byte("mdm.DeviceCommands"))
		for _, b := range []string{"mdm.Devices", "mdm.UDIDCertAuth", "mdm.DeviceCommands", "scep_certificates", "mdm.CommandHistory"} {
			if bk := tx.Bucket([]byte(b)); bk != nil {
				s := bk.Stats()
				fmt.Fprintf(os.Stderr, "bucket %-22s keys=%d leafInuse=%dMB\n", b, s.KeyN, (s.LeafInuse+s.BranchInuse)/1e6)
			}
		}
		return tx.Bucket([]byte("mdm.Devices")).ForEach(func(k, v []byte) error {
			var d device.Device
			if err := device.UnmarshalDevice(v, &d); err != nil {
				return nil
			}
			r := row{udid: d.UDID, serial: d.SerialNumber, model: d.ProductName, enrolled: d.Enrolled,
				lastSeen: d.LastSeen.Format("2006-01-02"), days: int(now.Sub(d.LastSeen).Hours() / 24), left: -99999}
			if h := auth.Get([]byte(d.UDID)); h != nil {
				if c, ok := certs[string(h)]; ok {
					r.notAfter = c.NotAfter.Format("2006-01-02")
					r.left = int(c.NotAfter.Sub(now).Hours() / 24)
				}
			}
			if q != nil {
				if raw := q.Get([]byte(d.UDID)); raw != nil {
					r.queued = len(raw) // serialized queue size in bytes
				}
			}
			rows = append(rows, r)
			return nil
		})
	})
	if err != nil {
		panic(err)
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].days < rows[j].days })
	w := csv.NewWriter(os.Stdout)
	w.Write([]string{"udid", "serial", "model", "enrolled", "last_seen", "days_since_seen", "cert_not_after", "cert_days_left", "queue_bytes"})
	for _, r := range rows {
		left := ""
		if r.left != -99999 {
			left = fmt.Sprint(r.left)
		}
		w.Write([]string{r.udid, r.serial, r.model, fmt.Sprint(r.enrolled), r.lastSeen, fmt.Sprint(r.days), r.notAfter, left, fmt.Sprint(r.queued)})
	}
	w.Flush()
}
