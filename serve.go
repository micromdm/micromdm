package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
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

	"github.com/fullsailor/pkcs7"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/crypto/pkcs12"

	"github.com/RobotsAndPencils/buford/push"
	"github.com/boltdb/bolt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/micromdm/dep"
	boltdepot "github.com/micromdm/scep/depot/bolt"
	scep "github.com/micromdm/scep/server"

	"github.com/micromdm/micromdm/blueprint"
	"github.com/micromdm/micromdm/checkin"
	"github.com/micromdm/micromdm/command"
	"github.com/micromdm/micromdm/connect"
	"github.com/micromdm/micromdm/core/apply"
	"github.com/micromdm/micromdm/core/list"
	"github.com/micromdm/micromdm/crypto"
	"github.com/micromdm/micromdm/depsync"
	"github.com/micromdm/micromdm/device"
	"github.com/micromdm/micromdm/enroll"
	"github.com/micromdm/micromdm/profile"
	"github.com/micromdm/micromdm/pubsub"
	nanopush "github.com/micromdm/micromdm/push"
	"github.com/micromdm/micromdm/queue"
)

const homePage = `<!doctype html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>MicroMDM</title>
	<style>
		body {
			font-family: -apple-system, BlinkMacSystemFont, sans-serif;
		}
	</style>
</head>
<body>
	<h3>Welcome to MicroMDM!</h3>
	<p><a href="mdm/enroll">Enroll a device</a></p>
</body>
</html>
`

func serve(args []string) error {
	flagset := flag.NewFlagSet("serve", flag.ExitOnError)
	var (
		flConfigPath   = flagset.String("config-path", "/var/db/micromdm", "path to configuration directory")
		flServerURL    = flagset.String("server-url", "", "public HTTPS url of your server")
		flAPIKey       = flagset.String("api-key", "", "API Token for mdmctl command")
		flAPNSCertPath = flagset.String("apns-cert", "", "path to APNS certificate")
		flAPNSKeyPass  = flagset.String("apns-password", "", "password for your p12 APNS cert file (if using)")
		flAPNSKeyPath  = flagset.String("apns-key", "", "path to key file if using .pem push cert")
		flTLS          = flagset.Bool("tls", true, "use https")
		flTLSCert      = flagset.String("tls-cert", "", "path to TLS certificate")
		flTLSKey       = flagset.String("tls-key", "", "path to TLS private key")
		flHTTPAddr     = flagset.String("http-addr", ":https", "http(s) listen address of mdm server. defaults to :8080 if tls is false")
		flRedirAddr    = flagset.String("redir-addr", ":http", "http redirect to https listen address")
		flHTTPDebug    = flagset.Bool("http-debug", false, "enable debug for http(dumps full request)")
		flRepoPath     = flagset.String("filerepo", "", "path to http file repo")
		flDepSim       = flagset.Bool("depsim", false, "use depsim config")
		flExamples     = flagset.Bool("examples", false, "prints some example usage")
	)
	flagset.Usage = usageFor(flagset, "micromdm serve [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}

	if *flExamples {
		printExamples()
		return nil
	}

	if *flServerURL == "" {
		return errors.New("must supply -server-url")
	}
	if !strings.HasPrefix(*flServerURL, "https://") {
		return errors.New("-server-url must begin with https://")
	}

	logger := log.NewLogfmtLogger(os.Stderr)
	stdlog.SetOutput(log.NewStdlibAdapter(logger)) // force structured logs
	mainLogger := log.With(logger, "component", "main")
	mainLogger.Log("msg", "started")

	if err := os.MkdirAll(*flConfigPath, 0755); err != nil {
		return errors.Wrapf(err, "creating config directory %s", *flConfigPath)
	}
	sm := &config{
		configPath:          *flConfigPath,
		ServerPublicURL:     strings.TrimRight(*flServerURL, "/"),
		APNSCertificatePath: *flAPNSCertPath,
		APNSPrivateKeyPass:  *flAPNSKeyPass,
		APNSPrivateKeyPath:  *flAPNSKeyPath,
		depsim:              *flDepSim,
		tlsCertPath:         *flTLSCert,

		// TODO: we have a static SCEP challenge password here to prevent
		// being prompted for the SCEP challenge which happens in a "normal"
		// (non-DEP) enrollment. While security is not improved it is at least
		// no less secure and prevents a useless dialog from showing.
		SCEPChallenge: "micromdm",
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

	devDB, err := device.NewDB(sm.db, sm.pubclient)
	if err != nil {
		stdlog.Fatal(err)
	}

	profDB, err := profile.NewDB(sm.db)
	if err != nil {
		stdlog.Fatal(err)
	}

	bpDB, err := blueprint.NewDB(sm.db, profDB)
	if err != nil {
		stdlog.Fatal(err)
	}

	if err := bpDB.StartListener(sm.pubclient, sm.commandService); err != nil {
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

	dc, err := sm.depClient()
	if err != nil {
		stdlog.Fatalf("creating DEP client %s\n", err)
	}
	var listsvc list.Service
	{
		listsvc = &list.ListService{DEPClient: dc, Devices: devDB, DB: sm.db, Blueprints: bpDB, Profiles: profDB}
	}
	var listDevicesEndpoint endpoint.Endpoint
	{
		listDevicesEndpoint = list.MakeListDevicesEndpoint(listsvc)

	}
	listEndpoints := list.Endpoints{
		ListDevicesEndpoint:       listDevicesEndpoint,
		GetDEPTokensEndpoint:      list.MakeGetDEPTokensEndpoint(listsvc),
		GetBlueprintsEndpoint:     list.MakeGetBlueprintsEndpoint(listsvc),
		GetProfilesEndpoint:       list.MakeGetProfilesEndpoint(listsvc),
		GetDEPAccountInfoEndpoint: list.MakeGetDEPAccountInfoEndpoint(listsvc),
		GetDEPProfileEndpoint:     list.MakeGetDEPProfileEndpoint(listsvc),
		GetDEPDeviceEndpoint:      list.MakeGetDEPDeviceDetailsEndpoint(listsvc),
	}

	var applysvc apply.Service
	{
		applysvc = &apply.ApplyService{DEPClient: dc, Blueprints: bpDB, DB: sm.db, Profiles: profDB}
	}

	var applyBlueprintEndpoint endpoint.Endpoint
	{
		applyBlueprintEndpoint = apply.MakeApplyBlueprintEndpoint(applysvc)
	}

	var applyProfileEndpoint endpoint.Endpoint
	{
		applyProfileEndpoint = apply.MakeApplyProfileEndpoint(applysvc)
	}

	applyEndpoints := apply.Endpoints{
		ApplyBlueprintEndpoint: applyBlueprintEndpoint,
		ApplyDEPTokensEndpoint: apply.MakeApplyDEPTokensEndpoint(applysvc),
		ApplyProfileEndpoint:   applyProfileEndpoint,
	}

	applyAPIHandlers := apply.MakeHTTPHandlers(ctx, applyEndpoints, connectOpts...)

	listAPIHandlers := list.MakeHTTPHandlers(ctx, listEndpoints, connectOpts...)

	connectHandlers := connect.MakeHTTPHandlers(ctx, connectEndpoints, connectOpts...)

	pushHandlers := nanopush.MakeHTTPHandlers(ctx, pushEndpoints, checkinOpts...)
	scepHandler := scep.ServiceHandler(ctx, sm.scepService, httpLogger)
	enrollHandlers := enroll.MakeHTTPHandlers(ctx, enroll.MakeServerEndpoints(sm.enrollService, sm.scepDepot), httptransport.ServerErrorLogger(httpLogger))
	r := mux.NewRouter()
	r.Handle("/mdm/checkin", mdmAuthSignMessageMiddleware(sm.scepDepot, checkinHandlers.CheckinHandler)).Methods("PUT")
	r.Handle("/mdm/connect", mdmAuthSignMessageMiddleware(sm.scepDepot, connectHandlers.ConnectHandler)).Methods("PUT")
	r.Handle("/mdm/enroll", enrollHandlers.EnrollHandler).Methods("GET", "POST")
	r.Handle("/ota/enroll", enrollHandlers.OTAEnrollHandler)
	r.Handle("/ota/phase23", enrollHandlers.OTAPhase2Phase3Handler).Methods("POST")
	r.Handle("/scep", scepHandler)
	r.Handle("/push/{udid}", pushHandlers.PushHandler)
	r.Handle("/v1/commands", commandHandlers.NewCommandHandler).Methods("POST")
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, homePage)
	})

	// API commands. Only handled if the user provides an api key.
	if *flAPIKey != "" {
		r.Handle("/v1/devices", apiAuthMiddleware(*flAPIKey, listAPIHandlers.ListDevicesHandler)).Methods("GET")
		r.Handle("/v1/dep-tokens", apiAuthMiddleware(*flAPIKey, listAPIHandlers.GetDEPTokensHandler)).Methods("GET")
		r.Handle("/v1/dep-tokens", apiAuthMiddleware(*flAPIKey, applyAPIHandlers.DEPTokensHandler)).Methods("PUT")
		r.Handle("/v1/blueprints", apiAuthMiddleware(*flAPIKey, listAPIHandlers.GetBlueprintsHandler)).Methods("GET")
		r.Handle("/v1/blueprints", apiAuthMiddleware(*flAPIKey, applyAPIHandlers.BlueprintHandler)).Methods("PUT")
		r.Handle("/v1/profiles", apiAuthMiddleware(*flAPIKey, listAPIHandlers.GetProfilesHandler)).Methods("GET")
		r.Handle("/v1/profiles", apiAuthMiddleware(*flAPIKey, applyAPIHandlers.ProfileHandler)).Methods("PUT")
		r.Handle("/v1/dep/devices", apiAuthMiddleware(*flAPIKey, listAPIHandlers.GetDEPDeviceDetailsHandler)).Methods("GET")
		r.Handle("/v1/dep/account", apiAuthMiddleware(*flAPIKey, listAPIHandlers.GetDEPAccountInfoHandler)).Methods("GET")
		r.Handle("/v1/dep/profiles", apiAuthMiddleware(*flAPIKey, listAPIHandlers.GetDEPProfileHander)).Methods("GET")
		r.Handle("/v1/dep/profiles", apiAuthMiddleware(*flAPIKey, applyAPIHandlers.DefineDEPProfileHandler)).Methods("POST")
	}

	if *flRepoPath != "" {
		if _, err := os.Stat(*flRepoPath); os.IsNotExist(err) {
			stdlog.Fatal(err)
		}
		r.PathPrefix("/repo/").Handler(http.StripPrefix("/repo/", http.FileServer(http.Dir(*flRepoPath))))
	}

	var handler http.Handler
	if *flHTTPDebug {
		handler = debugHTTPmiddleware(r)
	} else {
		handler = r
	}
	handler = handlers.CombinedLoggingHandler(os.Stdout, handler)
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
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
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
			errs <- serveACME(srv, srvURL.Hostname(), filepath.Join(sm.configPath, "le-certificates"))
			return
		}
	}()

	mainLogger.Log("terminated", <-errs)
	return nil
}

func printExamples() {
	const exampleText = `
		Quickstart:
		sudo micromdm serve -apns-cert /path/to/mdm_push_cert.p12 -apns-password=password_for_p12 -server-url=https://my-server-url

		Using self-signed certs:
		*Note, -apns flags are still required!*
		sudo micromdm serve -tls-cert=/path/to/server.crt -tls-key=/path/to/server.key

		`
	fmt.Println(exampleText)
}

func serveTLS(server *http.Server, certPath, keyPath string) error {
	err := server.ListenAndServeTLS(certPath, keyPath)
	return err
}

func serveACME(server *http.Server, domain, cachePath string) error {
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domain),
		Cache:      autocert.DirCache(cachePath),
	}
	server.TLSConfig.GetCertificate = m.GetCertificate
	err := server.ListenAndServeTLS("", "")
	return err
}

// redirects port 80 to port 443
func redirectTLS(addr, serverUrl string) {
	srv := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Connection", "close")
			url := serverUrl + req.URL.String()
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
	configPath          string
	depsim              bool
	pubclient           *pubsub.Inmem
	db                  *bolt.DB
	pushCert            pushServiceCert
	ServerPublicURL     string
	SCEPChallenge       string
	APNSPrivateKeyPath  string
	APNSCertificatePath string
	APNSPrivateKeyPass  string
	tlsCertPath         string
	scepDepot           *boltdepot.Depot

	// TODO: refactor enroll service and remove the need to reference
	// this on-disk cert. but it might be useful to keep the PEM
	// around for anyone who will need to export the CA.
	scepCACertPath string

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
	q, err := queue.NewQueue(c.db, c.pubclient)
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
	dbPath := filepath.Join(c.configPath, "micromdm.db")
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

	pkeyBlock := new(bytes.Buffer)
	pemBlock, _ = pem.Decode(pemData)
	if pemBlock == nil {
		c.err = errors.New("invalid PEM data for privkey")
		return
	}

	if x509.IsEncryptedPEMBlock(pemBlock) {
		b, err := x509.DecryptPEMBlock(pemBlock, []byte(c.APNSPrivateKeyPass))
		if err != nil {
			c.err = fmt.Errorf("decrypting DES private key %s", err)
			return
		}
		pkeyBlock.Write(b)
	} else {
		pkeyBlock.Write(pemBlock.Bytes)
	}

	priv, err := x509.ParsePKCS1PrivateKey(pkeyBlock.Bytes())
	if err != nil {
		c.err = fmt.Errorf("parsing pkcs1 private key: %s", err)
		return
	}
	c.pushCert.PrivateKey = priv
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
	c.pushService, err = nanopush.New(db, c.PushService, c.pubclient)
	if err != nil {
		c.err = err
		return
	}
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

	var SCEPCertificateSubject string
	// TODO: clean up order of inputs. Maybe pass *SCEPConfig as an arg?
	// but if you do, the packages are coupled, better not.
	c.enrollService, c.err = enroll.NewService(
		pushTopic,
		c.scepCACertPath,
		c.ServerPublicURL+"/scep",
		c.SCEPChallenge,
		c.ServerPublicURL,
		c.tlsCertPath,
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

func (c *config) depClient() (dep.Client, error) {
	if c.err != nil {
		return nil, c.err
	}
	// depsim config
	depsim := c.depsim
	var conf *dep.Config

	// try getting the oauth config from bolt
	tokens, err := list.GetDEPTokens(c.db)
	if err != nil {
		return nil, err
	}
	if len(tokens) >= 1 {
		conf = new(dep.Config)
		conf.ConsumerSecret = tokens[0].ConsumerSecret
		conf.ConsumerKey = tokens[0].ConsumerKey
		conf.AccessSecret = tokens[0].AccessSecret
		conf.AccessToken = tokens[0].AccessToken
		// TODO: handle expiration
	}

	// override with depsim keys if specified on CLI
	if depsim {
		conf = &dep.Config{
			ConsumerKey:    "CK_48dd68d198350f51258e885ce9a5c37ab7f98543c4a697323d75682a6c10a32501cb247e3db08105db868f73f2c972bdb6ae77112aea803b9219eb52689d42e6",
			ConsumerSecret: "CS_34c7b2b531a600d99a0e4edcf4a78ded79b86ef318118c2f5bcfee1b011108c32d5302df801adbe29d446eb78f02b13144e323eb9aad51c79f01e50cb45c3a68",
			AccessToken:    "AT_927696831c59ba510cfe4ec1a69e5267c19881257d4bca2906a99d0785b785a6f6fdeb09774954fdd5e2d0ad952e3af52c6d8d2f21c924ba0caf4a031c158b89",
			AccessSecret:   "AS_c31afd7a09691d83548489336e8ff1cb11b82b6bca13f793344496a556b1f4972eaff4dde6deb5ac9cf076fdfa97ec97699c34d515947b9cf9ed31c99dded6ba",
		}
	}

	if conf == nil {
		return nil, nil
	}

	depServerURL := "https://mdmenrollment.apple.com"
	if depsim {
		// TODO: support supplied depsim URL
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

	_, c.err = depsync.New(client, c.pubclient, c.db)
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

	c.scepCACertPath = filepath.Join(c.configPath, "SCEPCACert.pem")

	c.err = crypto.WritePEMCertificateFile(caCert, c.scepCACertPath)
	if c.err != nil {
		return
	}

	opts := []scep.ServiceOption{
		scep.ClientValidity(365),
		scep.ChallengePassword(c.SCEPChallenge),
	}
	c.scepDepot = depot
	c.scepService, c.err = scep.NewService(depot, opts...)
	if c.err == nil {
		c.scepService = scep.NewLoggingService(logger, c.scepService)
	}
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

// TODO: move to separate package/library
func mdmAuthSignMessageMiddleware(db *boltdepot.Depot, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b64sig := r.Header.Get("Mdm-Signature")
		if b64sig == "" {
			http.Error(w, "Signature missing", http.StatusBadRequest)
			return
		}
		sig, err := base64.StdEncoding.DecodeString(b64sig)
		if err != nil {
			http.Error(w, "Signature decoding error", http.StatusBadRequest)
			return
		}
		p7, err := pkcs7.Parse(sig)
		if err != nil {
			http.Error(w, "Signature parsing error", http.StatusBadRequest)
			return
		}
		bodyBuf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Problem reading request", http.StatusInternalServerError)
			return
		}

		// the signed data is the HTTP body message
		p7.Content = bodyBuf

		// reassign body to our already-read buffer
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBuf))
		// TODO: r.Body.Close() as we've ReadAll()'d it?

		err = p7.Verify()
		if err != nil {
			http.Error(w, "Signature verification error", http.StatusBadRequest)
			return
		}

		cert := p7.GetOnlySigner()
		if cert == nil {
			http.Error(w, "Invalid signer", http.StatusBadRequest)
			return
		}

		hasCN, err := HasCN(db, cert.Subject.CommonName, 0, cert, false)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Unable to validate signature", http.StatusInternalServerError)
			return
		}
		if !hasCN {
			fmt.Println("Unauthorized client signature from:", cert.Subject.CommonName)
			// NOTE: We're not returning 401 Unauthorized to avoid unenrolling a device
			// this may change in the future
			http.Error(w, "Unauthorized", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func basicAuth(password string) string {
	const authUsername = "micromdm"
	auth := authUsername + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func apiAuthMiddleware(token string, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, password, ok := r.BasicAuth()
		if !ok || password != token {
			w.Header().Set("WWW-Authenticate", `Basic realm="micromdm"`)
			http.Error(w, `{"error": "you need to log in"}`, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)

	}
}

// implement HasCN function that belongs in micromdm/scep/depot/bolt
//   note: added bool return, different from micromdm/scep interface
func HasCN(db *boltdepot.Depot, cn string, allowTime int, cert *x509.Certificate, revokeOldCertificate bool) (bool, error) {
	// TODO: implement allowTime
	// TODO: implement revocation
	if cert == nil {
		return false, errors.New("nil certificate provided")
	}
	var hasCN bool
	err := db.View(func(tx *bolt.Tx) error {
		// TODO: "scep_certificates" is internal const in micromdm/scep
		bucket := tx.Bucket([]byte("scep_certificates"))
		certKey := []byte(cert.Subject.CommonName + "." + cert.SerialNumber.String())
		certCandidate := bucket.Get(certKey)
		if certCandidate != nil {
			hasCN = bytes.Compare(certCandidate, cert.Raw) == 0
		}
		return nil
	})
	return hasCN, err
}
