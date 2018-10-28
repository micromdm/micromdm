package main

import (
	"bytes"
	"encoding/base64"
	"encoding/pem"
	"net/http"
	"os"

	"github.com/go-kit/kit/auth/basic"
	"github.com/go-kit/kit/log"
	"github.com/kolide/kit/dbutil"
	"github.com/kolide/kit/logutil"
	_ "github.com/lib/pq"
	scep "github.com/micromdm/scep/server"

	"github.com/micromdm/micromdm/pkg/httputil"
	"github.com/micromdm/micromdm/workflow/profile"
	"github.com/micromdm/micromdm/workflow/profile/ca"
	"github.com/micromdm/micromdm/workflow/profile/device"
	"github.com/micromdm/micromdm/workflow/profile/inventory"
	"github.com/micromdm/micromdm/workflow/profile/webhook"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)
	db, err := dbutil.OpenDBX(
		"postgres",
		"host=localhost port=5432 user=micromdm dbname=micromdm password=micromdm sslmode=disable",
		dbutil.WithLogger(logger),
		dbutil.WithMaxAttempts(3),
	)
	if err != nil {
		logutil.Fatal(logger, "err", err)
	}
	defer db.Close()

	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(10)

	ca, err := ca.New(db)
	if err != nil {
		logutil.Fatal(logger, "err", err)
	}
	cacert, _, err := ca.CA(nil)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	encoder.Write(pemCert(cacert[0].Raw))
	encoder.Write([]byte("\n"))
	scepOpts := []scep.ServiceOption{
		scep.ClientValidity(365),
	}
	scepsvc, err := scep.NewService(ca, scepOpts...)
	if err != nil {
		panic(err)
	}
	scepEndpoints := scep.MakeServerEndpoints(scepsvc)
	scepHandler := scep.MakeHTTPHandler(scepEndpoints, scepsvc, logger)

	svc := profile.New(db, profile.WithLogger(logger))
	basicAuthEndpointMiddleware := basic.AuthMiddleware("micromdm", "secret", "micromdm")
	e := profile.MakeServerEndpoints(svc, basicAuthEndpointMiddleware)
	r, options := httputil.NewRouter(logger)
	profile.RegisterHTTPHandlers(r, e, options...)

	devdb := device.New(db)
	inventorydb := inventory.New(db)
	webhookHandler := webhook.New(devdb, logger, buf.Bytes(), svc, inventorydb)

	r.Handle("/", webhookHandler)
	r.Handle("/scep", scepHandler)

	err = http.ListenAndServe(":9000", r)
	logutil.Fatal(logger, "err", err)
}

func pemCert(derBytes []byte) []byte {
	pemBlock := &pem.Block{
		Type:    certificatePEMBlockType,
		Headers: nil,
		Bytes:   derBytes,
	}
	out := pem.EncodeToMemory(pemBlock)
	return out
}

const certificatePEMBlockType = "CERTIFICATE"
