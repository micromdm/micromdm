package app

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/micromdm/scep/server"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/cors"
	"golang.org/x/net/context"

	"github.com/micromdm/micromdm/checkin"
	"github.com/micromdm/micromdm/command"
	"github.com/micromdm/micromdm/connect"
	"github.com/micromdm/micromdm/depenroll"
	"github.com/micromdm/micromdm/enroll"
	"github.com/micromdm/micromdm/management"
	"github.com/micromdm/micromdm/version"
)

func makeHTTPHandler(logger log.Logger, sm *serviceManager) http.Handler {
	httpLogger := log.NewContext(logger).With("component", "http")
	ctx := context.Background()

	mgmtHandler := management.ServiceHandler(ctx, sm.ManagementService, httpLogger)
	commandHandler := command.ServiceHandler(ctx, sm.CommandService, httpLogger)
	checkinHandler := checkin.ServiceHandler(ctx, sm.CheckinService, httpLogger)
	connectHandler := connect.ServiceHandler(ctx, sm.ConnectService, httpLogger)
	enrollHandler := enroll.ServiceHandler(ctx, sm.EnrollmentService, httpLogger)
	depenrollHandler := depenroll.ServiceHandler(ctx, sm.DEPEnrollmentService, httpLogger)

	var handler http.Handler
	mux := http.NewServeMux()
	mux.Handle("/management/v1/", mgmtHandler)
	mux.Handle("/mdm/commands", commandHandler)
	mux.Handle("/mdm/commands/", commandHandler)
	mux.Handle("/mdm/checkin", checkinHandler)
	mux.Handle("/mdm/connect", connectHandler)
	mux.Handle("/mdm/enroll", enrollHandler)
	mux.Handle("/mdm/enroll/dep", depenrollHandler)
	mux.Handle("/_metrics", prometheus.Handler())
	mux.Handle("/_version", version.Handler())
	if sm.Server.PackageRepoPath != "" {
		pkgrepoHandler := http.StripPrefix("/repo/",
			http.FileServer(http.Dir(sm.Server.PackageRepoPath)),
		)
		mux.Handle("/repo/", pkgrepoHandler)
	}
	handler = mux

	if len(sm.Server.CORSOrigins) > 0 {
		handler = cors.New(cors.Options{
			AllowedOrigins:   sm.Server.CORSOrigins,
			AllowCredentials: true,
			AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE"},
		}).Handler(mux)
	}

	return handler
}

func serveHTTP(logger log.Logger, h http.Handler, tlsEnabled bool, httpAddr, keyPath, certPath string) error {
	if tlsEnabled {
		if err := verifyTLSCerts(certPath, keyPath); err != nil {
			return err
		}

		logger.Log("msg", "serving https", "addr", httpAddr)
		return http.ListenAndServeTLS(httpAddr, certPath, keyPath, h)
	} else {
		logger.Log("msg", "serving http", "addr", httpAddr)
		return http.ListenAndServe(httpAddr, h)
	}
}

func serveSCEP(logger log.Logger, sm *serviceManager) error {
	httpLogger := log.NewContext(logger).With("component", "scep-http")
	if sm.SCEPService == nil {
		return nil
	}
	ctx := context.TODO()
	scepHandler := scepserver.ServiceHandler(ctx, sm.SCEPService, httpLogger)
	// FIXME allow specifying port.
	logger.Log("msg", "serving SCEP http", "addr", "0.0.0.0:2016")
	return http.ListenAndServe(":2016", scepHandler)
}

func verifyTLSCerts(certPath, keyPath string) error {
	chain, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return fmt.Errorf("serve: failed to load TLS cert or private key: %s", err)
	}

	cert, err := x509.ParseCertificate(chain.Certificate[0]) // Leaf is always the first entry
	if err != nil {
		return fmt.Errorf("server: error parsing TLS certificate: %s", err)
	}

	if _, err := cert.Verify(x509.VerifyOptions{}); err != nil {
		switch e := err.(type) {
		case x509.CertificateInvalidError:
			switch e.Reason {
			case x509.Expired:
				return fmt.Errorf("server certificate has expired: %s", err)
			default:
				return err
			}
		}
	}
	return nil
}
