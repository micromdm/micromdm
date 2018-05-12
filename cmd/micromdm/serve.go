package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	stdlog "log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/RobotsAndPencils/buford/push"
	"github.com/boltdb/bolt"
	"github.com/fullsailor/pkcs7"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/groob/finalizer/logutil"
	"github.com/micromdm/dep"
	"github.com/micromdm/go4/env"
	"github.com/micromdm/go4/httputil"
	"github.com/micromdm/go4/version"
	boltdepot "github.com/micromdm/scep/depot/bolt"
	scep "github.com/micromdm/scep/server"
	"github.com/pkg/errors"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/crypto/pkcs12"

	"github.com/micromdm/micromdm/dep/depsync"
	"github.com/micromdm/micromdm/mdm/checkin"
	"github.com/micromdm/micromdm/mdm/connect"
	"github.com/micromdm/micromdm/mdm/enroll"
	"github.com/micromdm/micromdm/pkg/crypto"
	"github.com/micromdm/micromdm/platform/apns"
	apnsbuiltin "github.com/micromdm/micromdm/platform/apns/builtin"
	"github.com/micromdm/micromdm/platform/appstore"
	appsbuiltin "github.com/micromdm/micromdm/platform/appstore/builtin"
	"github.com/micromdm/micromdm/platform/blueprint"
	blueprintbuiltin "github.com/micromdm/micromdm/platform/blueprint/builtin"
	"github.com/micromdm/micromdm/platform/command"
	"github.com/micromdm/micromdm/platform/config"
	configbuiltin "github.com/micromdm/micromdm/platform/config/builtin"
	depapi "github.com/micromdm/micromdm/platform/dep"
	"github.com/micromdm/micromdm/platform/device"
	devicebuiltin "github.com/micromdm/micromdm/platform/device/builtin"
	"github.com/micromdm/micromdm/platform/profile"
	profilebuiltin "github.com/micromdm/micromdm/platform/profile/builtin"
	"github.com/micromdm/micromdm/platform/pubsub"
	"github.com/micromdm/micromdm/platform/pubsub/inmem"
	"github.com/micromdm/micromdm/platform/queue"
	block "github.com/micromdm/micromdm/platform/remove"
	blockbuiltin "github.com/micromdm/micromdm/platform/remove/builtin"
	"github.com/micromdm/micromdm/platform/user"
	userbuiltin "github.com/micromdm/micromdm/platform/user/builtin"
	"github.com/micromdm/micromdm/workflow/webhook"
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
		flConfigPath        = flagset.String("config-path", "/var/db/micromdm", "path to configuration directory")
		flServerURL         = flagset.String("server-url", "", "public HTTPS url of your server")
		flAPIKey            = flagset.String("api-key", env.String("MICROMDM_API_KEY", ""), "API Token for mdmctl command")
		flAPNSCertPath      = flagset.String("apns-cert", "", "path to APNS certificate")
		flAPNSKeyPass       = flagset.String("apns-password", env.String("MICROMDM_APNS_KEY_PASSWORD", ""), "password for your p12 APNS cert file (if using)")
		flAPNSKeyPath       = flagset.String("apns-key", "", "path to key file if using .pem push cert")
		flTLS               = flagset.Bool("tls", true, "use https")
		flTLSCert           = flagset.String("tls-cert", "", "path to TLS certificate")
		flTLSKey            = flagset.String("tls-key", "", "path to TLS private key")
		flHTTPAddr          = flagset.String("http-addr", ":https", "http(s) listen address of mdm server. defaults to :8080 if tls is false")
		flHTTPDebug         = flagset.Bool("http-debug", false, "enable debug for http(dumps full request)")
		flRepoPath          = flagset.String("filerepo", "", "path to http file repo")
		flDepSim            = flagset.String("depsim", "", "use depsim URL")
		flExamples          = flagset.Bool("examples", false, "prints some example usage")
		flCommandWebhookURL = flagset.String("command-webhook-url", "", "URL to send command responses as raw plists.")
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
	sm := &server{
		configPath:          *flConfigPath,
		ServerPublicURL:     strings.TrimRight(*flServerURL, "/"),
		APNSCertificatePath: *flAPNSCertPath,
		APNSPrivateKeyPass:  *flAPNSKeyPass,
		APNSPrivateKeyPath:  *flAPNSKeyPath,
		depsim:              *flDepSim,
		tlsCertPath:         *flTLSCert,
		CommandWebhookURL:   *flCommandWebhookURL,

		webhooksHTTPClient: &http.Client{Timeout: time.Second * 30},

		// TODO: we have a static SCEP challenge password here to prevent
		// being prompted for the SCEP challenge which happens in a "normal"
		// (non-DEP) enrollment. While security is not improved it is at least
		// no less secure and prevents a useless dialog from showing.
		SCEPChallenge: "micromdm",
	}

	sm.setupPubSub()
	sm.setupBolt()
	sm.setupRemoveService()
	sm.setupConfigStore()
	sm.loadPushCerts()
	sm.setupSCEP(logger)
	sm.setupCheckinService()
	sm.setupPushService(logger)
	sm.setupCommandService()
	sm.setupWebhooks()
	sm.setupCommandQueue(logger)
	sm.setupDepClient()
	syncer := sm.setupDEPSync(logger)
	if sm.err != nil {
		stdlog.Fatal(sm.err)
	}

	removeService, err := block.New(sm.removeDB)
	if err != nil {
		stdlog.Fatal(err)
	}

	devDB, err := devicebuiltin.NewDB(sm.db, sm.pubclient)
	if err != nil {
		stdlog.Fatal(err)
	}

	userDB, err := userbuiltin.NewDB(sm.db, sm.pubclient, log.With(logger, "component", "user db"))
	if err != nil {
		stdlog.Fatal(err)
	}

	sm.profileDB, err = profilebuiltin.NewDB(sm.db)
	if err != nil {
		stdlog.Fatal(err)
	}

	sm.setupEnrollmentService()
	if sm.err != nil {
		stdlog.Fatalf("enrollment service: %s", sm.err)
	}

	bpDB, err := blueprintbuiltin.NewDB(sm.db, sm.profileDB, userDB)
	if err != nil {
		stdlog.Fatal(err)
	}

	if err := bpDB.StartListener(sm.pubclient, sm.commandService); err != nil {
		stdlog.Fatal(err)
	}

	sm.startWebhooks()

	ctx := context.Background()
	httpLogger := log.With(logger, "transport", "http")

	var checkinHandlers checkin.HTTPHandlers
	{
		e := checkin.Endpoints{
			CheckinEndpoint: checkin.MakeCheckinEndpoint(sm.checkinService),
		}
		opts := []httptransport.ServerOption{
			httptransport.ServerErrorLogger(httpLogger),
			httptransport.ServerErrorEncoder(checkin.EncodeError),
		}
		checkinHandlers = checkin.MakeHTTPHandlers(ctx, e, opts...)
	}

	commandEndpoints := command.MakeServerEndpoints(sm.commandService)

	connectOpts := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(httpLogger),
		httptransport.ServerErrorEncoder(connect.EncodeError),
	}

	var connectEndpoint endpoint.Endpoint
	{
		connectEndpoint = connect.MakeConnectEndpoint(sm.connectService)
	}
	connectEndpoints := connect.Endpoints{
		ConnectEndpoint: connectEndpoint,
	}

	dc := sm.depClient
	appDB := &appsbuiltin.Repo{Path: *flRepoPath}

	var profilesvc profile.Service
	{
		profilesvc = profile.New(sm.profileDB)
	}

	profileEndpoints := profile.MakeServerEndpoints(profilesvc)

	var blueprintsvc blueprint.Service
	{
		blueprintsvc = blueprint.New(bpDB)
	}

	blueprintEndpoints := blueprint.MakeServerEndpoints(blueprintsvc)
	blockEndpoints := block.MakeServerEndpoints(removeService)

	var usersvc user.Service
	{
		usersvc = user.New(userDB)
	}

	userEndpoints := user.MakeServerEndpoints(usersvc)

	var configsvc config.Service
	{
		configsvc = config.New(sm.configDB)
	}

	configEndpoints := config.MakeServerEndpoints(configsvc)

	var appsvc appstore.Service
	{
		appsvc = appstore.New(appDB)
	}

	appEndpoints := appstore.MakeServerEndpoints(appsvc)

	var devicesvc device.Service
	{
		devicesvc = device.New(devDB)
	}
	deviceEndpoints := device.MakeServerEndpoints(devicesvc)

	var depsvc depapi.Service
	{
		depsvc = depapi.New(dc, sm.pubclient)
	}
	depEndpoints := depapi.MakeServerEndpoints(depsvc)

	apnsEndpoints := apns.MakeServerEndpoints(sm.pushService)

	connectHandlers := connect.MakeHTTPHandlers(ctx, connectEndpoints, connectOpts...)

	scepHandler := scep.ServiceHandler(ctx, sm.scepService, httpLogger)
	enrollHandlers := enroll.MakeHTTPHandlers(ctx, enroll.MakeServerEndpoints(sm.enrollService, sm.scepDepot), httptransport.ServerErrorLogger(httpLogger))

	depsyncEndpoints := depsync.MakeServerEndpoints(depsync.NewService(syncer))

	r := mux.NewRouter()
	r.Handle("/version", version.Handler())
	r.Handle("/mdm/checkin", mdmAuthSignMessageMiddleware(sm.scepDepot, checkinHandlers.CheckinHandler)).Methods("PUT")
	r.Handle("/mdm/connect", mdmAuthSignMessageMiddleware(sm.scepDepot, connectHandlers.ConnectHandler)).Methods("PUT")
	r.Handle("/mdm/enroll", enrollHandlers.EnrollHandler).Methods("GET", "POST")
	r.Handle("/ota/enroll", enrollHandlers.OTAEnrollHandler)
	r.Handle("/ota/phase23", enrollHandlers.OTAPhase2Phase3Handler).Methods("POST")
	r.Handle("/scep", scepHandler)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, homePage)
	})

	profilesHandler := profile.MakeHTTPHandler(profileEndpoints, logger)
	blueprintsHandler := blueprint.MakeHTTPHandler(blueprintEndpoints, logger)
	blockhandler := block.MakeHTTPHandler(blockEndpoints, logger)
	userHandler := user.MakeHTTPHandler(userEndpoints, logger)
	configHandler := config.MakeHTTPHandler(configEndpoints, logger)
	appsHandler := appstore.MakeHTTPHandler(appEndpoints, logger)
	deviceHandler := device.MakeHTTPHandler(deviceEndpoints, logger)
	depHandlers := depapi.MakeHTTPHandler(depEndpoints, logger)
	apnsHandlers := apns.MakeHTTPHandler(apnsEndpoints, logger)
	depsyncHandlers := depsync.MakeHTTPHandler(depsyncEndpoints, logger)
	commandHandler := command.MakeHTTPHandler(commandEndpoints, logger)

	// API commands. Only handled if the user provides an api key.
	if *flAPIKey != "" {
		r.Handle("/v1/profiles", apiAuthMiddleware(*flAPIKey, profilesHandler))
		r.Handle("/v1/blueprints", apiAuthMiddleware(*flAPIKey, blueprintsHandler))
		r.Handle("/v1/users", apiAuthMiddleware(*flAPIKey, userHandler))
		r.Handle("/v1/apps", apiAuthMiddleware(*flAPIKey, appsHandler))
		r.Handle("/v1/devices/{udid}/block", apiAuthMiddleware(*flAPIKey, blockhandler))
		r.Handle("/v1/devices/{udid}/unblock", apiAuthMiddleware(*flAPIKey, blockhandler))
		r.Handle("/v1/devices", apiAuthMiddleware(*flAPIKey, deviceHandler))
		r.Handle("/v1/dep-tokens", apiAuthMiddleware(*flAPIKey, configHandler))
		r.Handle("/v1/dep-tokens", apiAuthMiddleware(*flAPIKey, configHandler))
		r.Handle("/v1/config/certificate", apiAuthMiddleware(*flAPIKey, configHandler))
		r.Handle("/v1/dep/devices", apiAuthMiddleware(*flAPIKey, depHandlers))
		r.Handle("/v1/dep/account", apiAuthMiddleware(*flAPIKey, depHandlers))
		r.Handle("/v1/dep/profiles", apiAuthMiddleware(*flAPIKey, depHandlers))
		r.Handle("/v1/dep/syncnow", apiAuthMiddleware(*flAPIKey, depsyncHandlers))
		r.Handle("/v1/dep/autoassigners", apiAuthMiddleware(*flAPIKey, depsyncHandlers))
		r.Handle("/v1/commands", apiAuthMiddleware(*flAPIKey, commandHandler))
		r.Handle("/push/{udid}", apiAuthMiddleware(*flAPIKey, apnsHandlers))
	} else {
		mainLogger.Log("msg", "no api key specified")
	}

	if *flRepoPath != "" {
		if _, err := os.Stat(*flRepoPath); os.IsNotExist(err) {
			stdlog.Fatal(err)
		}
		r.PathPrefix("/repo/").Handler(http.StripPrefix("/repo/", http.FileServer(http.Dir(*flRepoPath))))
	}

	var handler http.Handler
	if *flHTTPDebug {
		handler = httputil.HTTPDebugMiddleware(os.Stdout, true, logger.Log)(r)
	} else {
		handler = r
	}
	handler = logutil.NewHTTPLogger(httpLogger).Middleware(handler)

	srvURL, err := url.Parse(sm.ServerPublicURL)
	if err != nil {
		return errors.Wrapf(err, "parsing serverURL %q", sm.ServerPublicURL)
	}

	serveOpts := serveOptions(
		handler,
		*flHTTPAddr,
		srvURL.Hostname(),
		logger,
		*flTLSCert,
		*flTLSKey,
		sm.configPath,
		*flTLS,
	)
	err = httputil.ListenAndServe(serveOpts...)
	return errors.Wrap(err, "calling ListenAndServe")
}

// serveOptions configures the []httputil.Options for ListenAndServe
func serveOptions(
	handler http.Handler,
	addr string,
	hostname string,
	logger log.Logger,
	certPath string,
	keyPath string,
	configPath string,
	tls bool,
) []httputil.Option {
	tlsFromFile := (certPath != "" && keyPath != "")
	serveOpts := []httputil.Option{
		httputil.WithACMEHosts([]string{hostname}),
		httputil.WithLogger(logger),
		httputil.WithHTTPHandler(handler),
	}
	if tlsFromFile {
		serveOpts = append(serveOpts, httputil.WithKeyPair(certPath, keyPath))
	}
	if !tls && addr == ":https" {
		serveOpts = append(serveOpts, httputil.WithAddress(":8080"))
	}
	if tls {
		serveOpts = append(serveOpts, httputil.WithAutocertCache(autocert.DirCache(filepath.Join(configPath, "le-certificates"))))
	}
	if addr != ":https" {
		serveOpts = append(serveOpts, httputil.WithAddress(addr))
	}
	return serveOpts
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

type server struct {
	configPath          string
	depsim              string
	pubclient           pubsub.PublishSubscriber
	db                  *bolt.DB
	pushCert            pushServiceCert
	ServerPublicURL     string
	SCEPChallenge       string
	APNSPrivateKeyPath  string
	APNSCertificatePath string
	APNSPrivateKeyPass  string
	tlsCertPath         string
	scepDepot           *boltdepot.Depot
	profileDB           profile.Store
	configDB            config.Store
	removeDB            block.Store
	CommandWebhookURL   string
	depClient           dep.Client

	// TODO: refactor enroll service and remove the need to reference
	// this on-disk cert. but it might be useful to keep the PEM
	// around for anyone who will need to export the CA.
	scepCACertPath string

	PushService    *push.Service // bufford push
	pushService    apns.Service
	checkinService checkin.Service
	connectService connect.Service
	enrollService  enroll.Service
	scepService    scep.Service
	commandService command.Service
	configService  config.Service

	responseWebhook    *webhook.CommandWebhook
	webhooksHTTPClient *http.Client

	err error
}

func (c *server) setupPubSub() {
	if c.err != nil {
		return
	}
	c.pubclient = inmem.NewPubSub()
}

func (c *server) setupCommandService() {
	if c.err != nil {
		return
	}
	c.commandService, c.err = command.New(c.pubclient)
}

func (c *server) setupWebhooks() {
	if c.err != nil {
		return
	}

	if c.CommandWebhookURL == "" {
		return
	}

	h, err := webhook.NewCommandWebhook(c.webhooksHTTPClient, connect.ConnectTopic, c.CommandWebhookURL)
	if err != nil {
		c.err = err
		return
	}

	c.responseWebhook = h
}

func (c *server) startWebhooks() {
	if c.err != nil {
		return
	}

	if c.responseWebhook != nil {
		c.responseWebhook.StartListener(c.pubclient)
	}
}

func (c *server) setupRemoveService() {
	if c.err != nil {
		return
	}
	removeDB, err := blockbuiltin.NewDB(c.db)
	if err != nil {
		c.err = err
		return
	}
	c.removeDB = removeDB
}

func (c *server) setupCommandQueue(logger log.Logger) {
	if c.err != nil {
		return
	}
	q, err := queue.NewQueue(c.db, c.pubclient)
	if err != nil {
		c.err = err
		return
	}

	var connectService connect.Service
	{
		svc, err := connect.New(q, c.pubclient)
		if err != nil {
			c.err = err
			return
		}
		connectService = svc
		connectService = connect.LoggingMiddleware(log.With(level.Info(logger), "component", "connect"))(svc)
		connectService = block.RemoveMiddleware(c.removeDB)(connectService)
	}
	c.connectService = connectService
}

func (c *server) setupCheckinService() {
	if c.err != nil {
		return
	}
	c.checkinService, c.err = checkin.New(c.db, c.pubclient)
}

func (c *server) setupBolt() {
	if c.err != nil {
		return
	}
	dbPath := filepath.Join(c.configPath, "micromdm.db")
	db, err := bolt.Open(dbPath, 0644, &bolt.Options{Timeout: time.Second})
	if err != nil {
		c.err = errors.Wrap(err, "opening boltdb")
		return
	}
	c.db = db
}

func (c *server) loadPushCerts() {
	if c.APNSCertificatePath == "" && c.APNSPrivateKeyPass == "" && c.APNSPrivateKeyPath == "" {
		// this is optional, config could also be provided with mdmctl
		return
	}
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

	c.pushCert.Certificate, c.err = crypto.ReadPEMCertificateFile(c.APNSCertificatePath)
	if c.err != nil {
		return
	}

	var pemData []byte
	pemData, c.err = ioutil.ReadFile(c.APNSPrivateKeyPath)
	if c.err != nil {
		return
	}

	pkeyBlock := new(bytes.Buffer)
	pemBlock, _ := pem.Decode(pemData)
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

func (c *server) setupConfigStore() {
	if c.err != nil {
		return
	}
	db, err := configbuiltin.NewDB(c.db, c.pubclient)
	if err != nil {
		c.err = err
		return
	}
	c.configDB = db
	c.configService = config.New(db)

}

func (c *server) setupPushService(logger log.Logger) {
	if c.err != nil {
		return
	}

	var opts []apns.Option
	{
		cert, _ := c.configDB.PushCertificate()
		if c.pushCert.Certificate != nil && cert == nil {
			cert = &tls.Certificate{
				Certificate: [][]byte{c.pushCert.Certificate.Raw},
				PrivateKey:  c.pushCert.PrivateKey,
				Leaf:        c.pushCert.Certificate,
			}
		}
		if cert == nil {
			goto after
		}
		client, err := push.NewClient(*cert)
		if err != nil {
			c.err = err
			return
		}
		svc := push.NewService(client, push.Production)
		opts = append(opts, apns.WithPushService(svc))
	}
after:

	db, err := apnsbuiltin.NewDB(c.db, c.pubclient)
	if err != nil {
		c.err = err
		return
	}

	service, err := apns.New(db, c.configDB, c.pubclient, opts...)
	if err != nil {
		c.err = errors.Wrap(err, "starting micromdm push service")
		return
	}
	c.pushService = apns.LoggingMiddleware(
		log.With(level.Info(logger), "component", "apns"),
	)(service)
}

func (c *server) setupEnrollmentService() {
	if c.err != nil {
		return
	}

	var topicProvider enroll.TopicProvider
	if c.pushCert.Certificate != nil {
		pushTopic, err := crypto.TopicFromCert(c.pushCert.Certificate)
		if err != nil {
			c.err = errors.Wrap(err, "get apns topic from certificate")
			return
		}
		topicProvider = staticTopicProvider{topic: pushTopic}
	} else {
		topicProvider = c.configDB
	}

	var SCEPCertificateSubject string
	// TODO: clean up order of inputs. Maybe pass *SCEPConfig as an arg?
	// but if you do, the packages are coupled, better not.
	c.enrollService, c.err = enroll.NewService(
		topicProvider,
		c.pubclient,
		c.scepCACertPath,
		c.ServerPublicURL+"/scep",
		c.SCEPChallenge,
		c.ServerPublicURL,
		c.tlsCertPath,
		SCEPCertificateSubject,
		c.profileDB,
	)
}

// if the apns-cert flags are specified this provider will be used in the enroll service.
type staticTopicProvider struct{ topic string }

func (p staticTopicProvider) PushTopic() (string, error) {
	return p.topic, nil
}

func (c *server) setupDepClient() (dep.Client, error) {
	if c.err != nil {
		return nil, c.err
	}
	// depsim config
	depsim := c.depsim
	var conf *dep.Config

	// try getting the oauth config from bolt
	tokens, err := c.configDB.DEPTokens()
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
	if depsim != "" {
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
	if depsim != "" {
		depServerURL = depsim
	}
	client, err := dep.NewClient(conf, dep.ServerURL(depServerURL))
	if err != nil {
		return nil, err
	}

	c.depClient = client

	return client, nil
}

func (c *server) setupDEPSync(logger log.Logger) depsync.Syncer {
	if c.err != nil {
		return nil
	}

	client := c.depClient
	opts := []depsync.Option{
		depsync.WithLogger(log.With(logger, "component", "depsync")),
	}
	if client != nil {
		opts = append(opts, depsync.WithClient(client))
	}

	var syncer depsync.Syncer
	syncer, c.err = depsync.New(c.pubclient, c.db, logger, opts...)
	if c.err != nil {
		return nil
	}
	return syncer
}

func (c *server) setupSCEP(logger log.Logger) {
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

func apiAuthMiddleware(token string, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, password, ok := r.BasicAuth()
		if !ok || password != token {
			w.Header().Set("WWW-Authenticate", `Basic realm="micromdm"`)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "you need to log in"}`))
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
