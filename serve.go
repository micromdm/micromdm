package main

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/asn1"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	stdlog "log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/crypto/pkcs12"

	"github.com/RobotsAndPencils/buford/push"
	"github.com/boltdb/bolt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/micromdm/dep"
	boltdepot "github.com/micromdm/scep/depot/bolt"
	scep "github.com/micromdm/scep/server"

	"github.com/micromdm/micromdm/checkin"
	"github.com/micromdm/micromdm/command"
	"github.com/micromdm/micromdm/connect"
	"github.com/micromdm/micromdm/depsync"
	"github.com/micromdm/micromdm/device"
	"github.com/micromdm/micromdm/enroll"
	"github.com/micromdm/micromdm/pubsub"
	nanopush "github.com/micromdm/micromdm/push"
)

const configDBPath = "/var/db/micromdm"

func serve(args []string) error {
	flagset := flag.NewFlagSet("serve", flag.ExitOnError)
	var (
		flServerURL    = flagset.String("server-url", "", "public HTTPS url of your server")
		flAPNSCertPath = flagset.String("apns-cert", "mdm.p12", "path to APNS certificate")
		flAPNSKeyPass  = flagset.String("apns-password", "secret", "password for your p12 APNS cert file (if using)")
		flAPNSKeyPath  = flagset.String("apns-key", "", "path to key file if using .pem push cert")
		flTLS          = flagset.Bool("tls", true, "use https")
		flTLSCert      = flagset.String("tls-cert", "", "path to TLS certificate")
		flTLSKey       = flagset.String("tls-key", "", "path to TLS private key")
		flHTTPAddr     = flagset.String("http-addr", ":https", "http(s) listen address of mdm server. defaults to :8080 if tls is false")
		flRedirAddr    = flagset.String("redir-addr", ":http", "http redirect to https listen address")
		flHTTPDebug    = flagset.Bool("http-debug", false, "enable debug for http(dumps full request)")
		flRepoPath     = flagset.String("filerepo", "", "path to http file repo")
		flDepSim       = flagset.Bool("depsim", false, "use depsim config")
	)
	flagset.Usage = usageFor(flagset, "micromdm serve [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}

	if *flServerURL == "" {
		return errors.New("must supply -server-url")
	}

	logger := log.NewLogfmtLogger(os.Stderr)
	stdlog.SetOutput(log.NewStdlibAdapter(logger)) // force structured logs
	mainLogger := log.With(logger, "component", "main")
	mainLogger.Log("msg", "started")

	sm := &config{
		ServerPublicURL:     *flServerURL,
		APNSCertificatePath: *flAPNSCertPath,
		APNSPrivateKeyPass:  *flAPNSKeyPass,
		APNSPrivateKeyPath:  *flAPNSKeyPath,
		depsim:              *flDepSim,
	}
	if err := os.MkdirAll(configDBPath, 0755); err != nil {
		return errors.Wrapf(err, "creating config directory %s", configDBPath)
	}
	sm.setupPubSub()
	sm.setupBolt()
	sm.loadPushCerts()
	sm.setupSCEP(logger)
	sm.setupEnrollmentService()
	sm.setupCheckinService()
	sm.setupPushService()
	sm.setupCommandService()
	sm.setupCommandQueue()
	sm.setupDEPSync()
	if sm.err != nil {
		stdlog.Fatal(sm.err)
	}

	if err := hardcodeCommands(sm); err != nil {
		stdlog.Fatal(err)
	}

	_, err := device.NewDB(sm.db, sm.pubclient)
	if err != nil {
		stdlog.Fatal(err)
	}

	ctx := context.Background()
	httpLogger := log.With(logger, "transport", "http")
	var checkinEndpoint endpoint.Endpoint
	{
		checkinEndpoint = checkin.MakeCheckinEndpoint(sm.checkinService)
	}

	checkinEndpoints := checkin.Endpoints{
		CheckinEndpoint: checkinEndpoint,
	}

	checkinOpts := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(httpLogger),
		httptransport.ServerErrorEncoder(checkin.EncodeError),
	}
	checkinHandlers := checkin.MakeHTTPHandlers(ctx, checkinEndpoints, checkinOpts...)

	pushEndpoints := nanopush.Endpoints{
		PushEndpoint: nanopush.MakePushEndpoint(sm.pushService),
	}

	commandEndpoints := command.Endpoints{
		NewCommandEndpoint: command.MakeNewCommandEndpoint(sm.commandService),
	}

	connectOpts := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(httpLogger),
		httptransport.ServerErrorEncoder(connect.EncodeError),
	}
	commandHandlers := command.MakeHTTPHandlers(ctx, commandEndpoints, connectOpts...)

	var connectEndpoint endpoint.Endpoint
	{
		connectEndpoint = connect.MakeConnectEndpoint(sm.connectService)
	}
	connectEndpoints := connect.Endpoints{
		ConnectEndpoint: connectEndpoint,
	}

	connectHandlers := connect.MakeHTTPHandlers(ctx, connectEndpoints, connectOpts...)

	pushHandlers := nanopush.MakeHTTPHandlers(ctx, pushEndpoints, checkinOpts...)
	scepHandler := scep.ServiceHandler(ctx, sm.scepService, httpLogger)
	enrollHandler := enroll.ServiceHandler(ctx, sm.enrollService, httpLogger)
	r := mux.NewRouter()
	r.Handle("/mdm/checkin", checkinHandlers.CheckinHandler).Methods("PUT")
	r.Handle("/mdm/connect", connectHandlers.ConnectHandler).Methods("PUT")
	r.Handle("/mdm/enroll", enrollHandler).Methods("GET", "POST")
	r.Handle("/scep", scepHandler)
	r.Handle("/push/{udid}", pushHandlers.PushHandler)
	r.Handle("/v1/commands", commandHandlers.NewCommandHandler).Methods("POST")

	if *flRepoPath != "" {
		r.PathPrefix("/repo/").Handler(http.StripPrefix("/repo/", http.FileServer(http.Dir(*flRepoPath))))
	}

	var handler http.Handler
	if *flHTTPDebug {
		handler = debugHTTPmiddleware(r)
	} else {
		handler = r
	}
	srv := &http.Server{
		Addr:              *flHTTPAddr,
		Handler:           handler,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       10 * time.Minute,
		MaxHeaderBytes:    1 << 18, // 0.25 MB
		TLSConfig:         tlsConfig(),
	}

	srvURL, err := url.Parse(sm.ServerPublicURL)
	if err != nil {
		stdlog.Fatal(err)
	}

	errs := make(chan error, 2)
	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig // block on signal then gracefully shutdown.
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		errs <- srv.Shutdown(ctx)
	}()

	go func() {
		logger := log.With(logger, "transport", "HTTP")
		if !*flTLS {
			var httpAddr string
			if *flHTTPAddr == ":https" {
				httpAddr = ":8080"
			} else {
				httpAddr = *flHTTPAddr
			}
			logger.Log("addr", httpAddr)
			errs <- http.ListenAndServe(httpAddr, handler)
			return
		}

		tlsFromFile := (*flTLSCert != "" && *flTLSKey != "")
		if tlsFromFile {
			logger.Log("addr", srv.Addr)
			redirectTLS(*flRedirAddr, sm.ServerPublicURL)
			errs <- serveTLS(srv, *flTLSCert, *flTLSKey)
			return
		} else {
			logger.Log("addr", srv.Addr)
			redirectTLS(*flRedirAddr, sm.ServerPublicURL)
			errs <- serveACME(srv, srvURL.Hostname())
			return
		}
	}()

	mainLogger.Log("terminated", <-errs)
	return nil
}

func serveTLS(server *http.Server, certPath, keyPath string) error {
	err := server.ListenAndServeTLS(certPath, keyPath)
	return err
}

func serveACME(server *http.Server, domain string) error {
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domain),
		Cache:      autocert.DirCache("/var/db/micromdm/le-certificates"),
	}
	server.TLSConfig.GetCertificate = m.GetCertificate
	err := server.ListenAndServeTLS("", "")
	return err
}

// redirects port 80 to port 443
func redirectTLS(addr, serverUrl string) {
	// trim trailing slash to avoid double slashes
	trimmedServer := strings.TrimRight(serverUrl, "/")
	srv := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Connection", "close")
			url := trimmedServer + req.URL.String()
			http.Redirect(w, req, url, http.StatusMovedPermanently)
		}),
	}
	go func() { stdlog.Fatal(srv.ListenAndServe()) }()
}

func tlsConfig() *tls.Config {
	cfg := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}
	return cfg
}

type config struct {
	depsim              bool
	pubclient           *pubsub.Inmem
	db                  *bolt.DB
	pushCert            pushServiceCert
	ServerPublicURL     string
	SCEPChallenge       string
	APNSPrivateKeyPath  string
	APNSCertificatePath string
	APNSPrivateKeyPass  string

	PushService    *push.Service // bufford push
	pushService    *nanopush.Push
	checkinService checkin.Service
	connectService connect.ConnectService
	enrollService  enroll.Service
	scepService    scep.Service
	commandService command.Service

	err error
}

func (c *config) setupPubSub() {
	if c.err != nil {
		return
	}
	c.pubclient = pubsub.NewInmemPubsub()
}

func (c *config) setupCommandService() {
	if c.err != nil {
		return
	}
	c.commandService, c.err = command.New(c.db, c.pubclient)
}

func (c *config) setupCommandQueue() {
	if c.err != nil {
		return
	}
	q, err := connect.NewQueue(c.db, c.pubclient)
	if err != nil {
		c.err = err
		return
	}

	connSvc, err := connect.New(q)
	if err != nil {
		c.err = err
		return
	}
	c.connectService = connSvc
}

func (c *config) setupCheckinService() {
	if c.err != nil {
		return
	}
	c.checkinService, c.err = checkin.New(c.db, c.pubclient)
}

func (c *config) setupBolt() {
	if c.err != nil {
		return
	}
	dbPath := filepath.Join(configDBPath, "micromdm.db")
	c.db, c.err = bolt.Open(dbPath, 0644, nil)
	if c.err != nil {
		return
	}
}

func (c *config) loadPushCerts() {
	if c.err != nil {
		return
	}

	if c.APNSPrivateKeyPath == "" {
		var pkcs12Data []byte
		pkcs12Data, c.err = ioutil.ReadFile(c.APNSCertificatePath)
		if c.err != nil {
			return
		}
		c.pushCert.PrivateKey, c.pushCert.Certificate, c.err =
			pkcs12.Decode(pkcs12Data, c.APNSPrivateKeyPass)
		return
	}

	var pemData []byte
	pemData, c.err = ioutil.ReadFile(c.APNSCertificatePath)
	if c.err != nil {
		return
	}

	pemBlock, _ := pem.Decode(pemData)
	if pemBlock == nil {
		c.err = errors.New("invalid PEM data for cert")
		return
	}
	c.pushCert.Certificate, c.err = x509.ParseCertificate(pemBlock.Bytes)
	if c.err != nil {
		return
	}

	pemData, c.err = ioutil.ReadFile(c.APNSPrivateKeyPath)
	if c.err != nil {
		return
	}

	pemBlock, _ = pem.Decode(pemData)
	if pemBlock == nil {
		c.err = errors.New("invalid PEM data for privkey")
		return
	}
	c.pushCert.PrivateKey, c.err =
		x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
}

type pushServiceCert struct {
	*x509.Certificate
	PrivateKey interface{}
}

func (c *config) setupPushService() {
	if c.err != nil {
		return
	}
	tlsCert := tls.Certificate{
		Certificate: [][]byte{c.pushCert.Certificate.Raw},
		PrivateKey:  c.pushCert.PrivateKey,
		Leaf:        c.pushCert.Certificate,
	}
	client, err := push.NewClient(tlsCert)
	if err != nil {
		c.err = err
		return
	}
	c.PushService = &push.Service{
		Client: client,
		Host:   push.Production,
	}

	db, err := nanopush.NewDB(c.db, c.pubclient)
	if err != nil {
		c.err = err
		return
	}
	c.pushService = nanopush.New(db, c.PushService)
}

func (c *config) setupEnrollmentService() {
	if c.err != nil {
		return
	}
	pushTopic, err := topicFromCert(c.pushCert.Certificate)
	if err != nil {
		c.err = err
		return
	}
	SCEPRemoteURL := strings.TrimRight(c.ServerPublicURL, "/") + "/scep"

	var tlsCert string
	var SCEPCertificateSubject string
	// TODO: clean up order of inputs. Maybe pass *SCEPConfig as an arg?
	// but if you do, the packages are coupled, better not.
	c.enrollService, c.err = enroll.NewService(
		pushTopic,
		scepCACertName,
		SCEPRemoteURL,
		c.SCEPChallenge,
		c.ServerPublicURL,
		tlsCert,
		SCEPCertificateSubject,
	)
}

func topicFromCert(cert *x509.Certificate) (string, error) {
	var oidASN1UserID = asn1.ObjectIdentifier{0, 9, 2342, 19200300, 100, 1, 1}
	for _, v := range cert.Subject.Names {
		if v.Type.Equal(oidASN1UserID) {
			return v.Value.(string), nil
		}
	}

	return "", errors.New("could not find Push Topic (UserID OID) in certificate")
}

// TODO: refactor enroll service and remove the need to reference this cert.
// but it might be useful to keep the PEM around for anyone who will need to export
// the CA.
const scepCACertName = "/var/db/micromdm/SCEPCACert.pem"

func (c *config) depClient() (dep.Client, error) {
	if c.err != nil {
		return nil, c.err
	}
	// depsim config
	depsim := c.depsim
	var conf *dep.Config
	if depsim {
		conf = &dep.Config{
			ConsumerKey:    "CK_48dd68d198350f51258e885ce9a5c37ab7f98543c4a697323d75682a6c10a32501cb247e3db08105db868f73f2c972bdb6ae77112aea803b9219eb52689d42e6",
			ConsumerSecret: "CS_34c7b2b531a600d99a0e4edcf4a78ded79b86ef318118c2f5bcfee1b011108c32d5302df801adbe29d446eb78f02b13144e323eb9aad51c79f01e50cb45c3a68",
			AccessToken:    "AT_927696831c59ba510cfe4ec1a69e5267c19881257d4bca2906a99d0785b785a6f6fdeb09774954fdd5e2d0ad952e3af52c6d8d2f21c924ba0caf4a031c158b89",
			AccessSecret:   "AS_c31afd7a09691d83548489336e8ff1cb11b82b6bca13f793344496a556b1f4972eaff4dde6deb5ac9cf076fdfa97ec97699c34d515947b9cf9ed31c99dded6ba",
		}
	}

	// try getting the oauth config from bolt

	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(depTokenBucket))
		if b == nil {
			return nil
		}
		_, v := b.Cursor().First()
		if v == nil {
			return nil
		}
		var token DEPTokenJSON
		err := json.Unmarshal(v, &token)
		if err != nil {
			return err
		}
		conf = new(dep.Config)
		conf.ConsumerSecret = token.ConsumerSecret
		conf.ConsumerKey = token.ConsumerKey
		conf.AccessSecret = token.AccessSecret
		conf.AccessToken = token.AccessToken
		// TODO handle expiration.
		return nil
	})
	if err != nil {
		return nil, err
	}

	if conf == nil {
		return nil, nil
	}

	depServerURL := "https://mdmenrollment.apple.com"
	if depsim {
		depServerURL = "http://dep.micromdm.io:9000"
	}
	client, err := dep.NewClient(conf, dep.ServerURL(depServerURL))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *config) setupDEPSync() {
	if c.err != nil {
		return
	}

	client, err := c.depClient()
	if err != nil {
		c.err = err
		return
	}
	if client == nil {
		fmt.Println("no DEP server configured. skipping device sync from DEP.")
		return
	}

	_, c.err = depsync.New(client, c.pubclient)
	if err != nil {
		return
	}
}

func (c *config) setupSCEP(logger log.Logger) {
	if c.err != nil {
		return
	}

	depot, err := boltdepot.NewBoltDepot(c.db)
	if err != nil {
		c.err = err
		return
	}

	key, err := depot.CreateOrLoadKey(2048)
	if err != nil {
		c.err = err
		return
	}

	caCert, err := depot.CreateOrLoadCA(key, 5, "MicroMDM", "US")
	if err != nil {
		c.err = err
		return
	}

	c.err = savePEMCert(scepCACertName, caCert)
	if c.err != nil {
		return
	}

	opts := []scep.ServiceOption{
		scep.ClientValidity(365),
	}
	c.scepService, c.err = scep.NewService(depot, opts...)
	if c.err == nil {
		c.scepService = scep.NewLoggingService(logger, c.scepService)
	}
}

func savePEMKey(path string, key *rsa.PrivateKey) error {
	keyOutput, err := os.OpenFile(path,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600,
	)
	if err != nil {
		return err
	}
	defer keyOutput.Close()

	return pem.Encode(
		keyOutput,
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		})
}

func savePEMCert(path string, cert *x509.Certificate) error {
	certOutput, err := os.Create(path)
	if err != nil {
		return err
	}
	defer certOutput.Close()

	return pem.Encode(
		certOutput,
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		})
}

func debugHTTPmiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := io.TeeReader(r.Body, os.Stderr)
		r.Body = ioutil.NopCloser(body)
		out, err := httputil.DumpRequest(r, true)
		if err != nil {
			stdlog.Println(err)
		}
		fmt.Println("")
		fmt.Println(string(out))
		fmt.Println("")
		next.ServeHTTP(w, r)
	})
}
