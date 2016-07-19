package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"database/sql"
	"github.com/DavidHuie/gomigrate"
	"github.com/RobotsAndPencils/buford/certificate"
	"github.com/RobotsAndPencils/buford/push"
	"github.com/go-kit/kit/log"
	"github.com/micromdm/dep"
	"github.com/micromdm/micromdm/applications"
	"github.com/micromdm/micromdm/certificates"
	"github.com/micromdm/micromdm/checkin"
	"github.com/micromdm/micromdm/command"
	"github.com/micromdm/micromdm/connect"
	"github.com/micromdm/micromdm/device"
	"github.com/micromdm/micromdm/enroll"
	"github.com/micromdm/micromdm/management"
	"github.com/micromdm/micromdm/workflow"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/rs/cors"
	"golang.org/x/net/context"
	"time"
)

var (
	// Version info
	Version = "unreleased"
	gitHash = "unknown"
)

func main() {
	ctx := context.Background()
	logger := log.NewLogfmtLogger(os.Stderr)

	//flags
	var (
		flUrl          = flag.String("url", envString("MICROMDM_URL", ""), "public facing url")
		flPort         = flag.String("port", envString("MICROMDM_HTTP_LISTEN_PORT", ""), "port to listen on")
		flTLS          = flag.Bool("tls", envBool("MICROMDM_USE_TLS"), "use https")
		flTLSCert      = flag.String("tls-cert", envString("MICROMDM_TLS_CERT", ""), "path to TLS certificate")
		flTLSKey       = flag.String("tls-key", envString("MICROMDM_TLS_KEY", ""), "path to TLS private key")
		flTLSCACert    = flag.String("tls-ca-cert", envString("MICROMDM_TLS_CA_CERT", ""), "path to CA certificate")
		flScepUrl      = flag.String("scep-url", envString("MICROMDM_SCEP_URL", ""), "url of SCEP server")
		flPGconn       = flag.String("postgres", envString("MICROMDM_POSTGRES_CONN_URL", ""), "postgres connection url")
		flRedisconn    = flag.String("redis", envString("MICROMDM_REDIS_CONN_URL", ""), "redis connection url")
		flVersion      = flag.Bool("version", false, "print version information")
		flPushCert     = flag.String("push-cert", envString("MICROMDM_PUSH_CERT", ""), "path to push certificate")
		flPushPass     = flag.String("push-pass", envString("MICROMDM_PUSH_PASS", ""), "push certificate password")
		flEnrollment   = flag.String("profile", envString("MICROMDM_ENROLL_PROFILE", ""), "path to enrollment profile")
		flDEPCK        = flag.String("dep-consumer-key", envString("DEP_CONSUMER_KEY", ""), "dep consumer key")
		flDEPCS        = flag.String("dep-consumer-secret", envString("DEP_CONSUMER_SECRET", ""), "dep consumer secret")
		flDEPAT        = flag.String("dep-access-token", envString("DEP_ACCESS_TOKEN", ""), "dep access token")
		flDEPAS        = flag.String("dep-access-secret", envString("DEP_ACCESS_SECRET", ""), "dep access secret")
		flDEPsim       = flag.Bool("depsim", envBool("DEP_USE_DEPSIM"), "use default depsim credentials")
		flDEPServerURL = flag.String("dep-server-url", envString("DEP_SERVER_URL", ""), "dep server url. for testing. Use blank if not running against depsim")
		flPkgRepo      = flag.String("pkg-repo", envString("MICROMDM_PKG_REPO", ""), "path to pkg repo")
		flCorsOrigin   = flag.String("cors-origin", envString("MICROMDM_CORS_ORIGIN", ""), "allowed domain for cross origin resource sharing")
	)

	// set tls to true by default. let user set it to false
	*flTLS = true
	flag.Parse()

	// -version flag
	if *flVersion {
		fmt.Printf("micromdm - Version %s\n", Version)
		fmt.Printf("Git Hash - %s\n", gitHash)
		os.Exit(0)
	}

	// check port flag
	// if none is provided, default to 80 or 443
	if *flPort == "" {
		port := defaultPort(*flTLS)
		logger.Log("msg", fmt.Sprintf("No port flag specified. Using %v by default", port))
		*flPort = port
	}

	if *flEnrollment == "" {
		logger.Log("err", "must set path to enrollment profile")
		os.Exit(1)
	}
	enrollmentProfile, err := ioutil.ReadFile(*flEnrollment)
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}

	// check cert and key if -tls=true
	if *flTLS {
		if err := checkTLSFlags(*flTLSKey, *flTLSCert); err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
	}

	pgHostAddr := os.Getenv("POSTGRES_PORT_5432_TCP_ADDR")
	if *flPGconn == "" && pgHostAddr != "" {
		*flPGconn = getPGConnFromENV(logger, pgHostAddr)
	}

	// check database connection
	if *flPGconn == "" {
		logger.Log("err", "database connection url not specified")
		os.Exit(1)
	}
	if checkEmptyArgs(*flPushCert, *flPushPass) {
		logger.Log("err", "must specify push cert path and password")
		os.Exit(1)
	}

	pushSvc, err := pushService(*flPushCert, *flPushPass)
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}

	// Run migrations
	db, err := sql.Open("postgres", *flPGconn)
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}
	var dbError error
	maxAttempts := 20
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		dbError = db.Ping()
		if dbError == nil {
			break
		}
		logger.Log("msg", fmt.Sprintf("could not connect to postgres: %v", dbError))
		time.Sleep(time.Duration(attempts) * time.Second)
	}
	if dbError != nil {
		logger.Log("err", dbError)
		os.Exit(1)
	}

	migrator, _ := gomigrate.NewMigrator(db, gomigrate.Postgres{}, "./migrations")
	migrationErr := migrator.Migrate()

	if migrationErr != nil {
		logger.Log("err", err)
		os.Exit(1)
	}

	workflowDB, err := workflow.NewDB(
		"postgres",
		*flPGconn,
		logger,
	)
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}

	deviceDB, err := device.NewDB(
		"postgres",
		*flPGconn,
		logger,
	)
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}

	redisHostAddr := os.Getenv("REDIS_PORT_6379_TCP_ADDR")
	if *flRedisconn == "" && redisHostAddr != "" {
		*flRedisconn = getRedisConnFromENV(redisHostAddr)
	}

	// check database connection
	if *flRedisconn == "" {
		logger.Log("err", "database connection url not specified")
		os.Exit(1)
	}

	commandDB, err := command.NewDB("redis", *flRedisconn, logger)
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}

	dc := depClient(logger, *flDEPCK, *flDEPCS, *flDEPAT, *flDEPAS, *flDEPServerURL, *flDEPsim)
	mgmtSvc := management.NewService(deviceDB, workflowDB, dc, pushSvc)
	commandSvc := command.NewService(commandDB)
	checkinSvc := checkin.NewService(deviceDB, mgmtSvc, commandSvc, enrollmentProfile)
	connectSvc := connect.NewService(deviceDB, commandSvc)

	httpLogger := log.NewContext(logger).With("component", "http")
	managementHandler := management.ServiceHandler(ctx, mgmtSvc, httpLogger)
	commandHandler := command.ServiceHandler(ctx, commandSvc, httpLogger)
	checkinHandler := checkin.ServiceHandler(ctx, checkinSvc, httpLogger)
	connectHandler := connect.ServiceHandler(ctx, connectSvc, httpLogger)

	mux := http.NewServeMux()

	mux.Handle("/management/v1/", managementHandler)
	mux.Handle("/mdm/commands", commandHandler)
	mux.Handle("/mdm/commands/", commandHandler)
	mux.Handle("/mdm/checkin", checkinHandler)
	mux.Handle("/mdm/connect", connectHandler)

	if checkEmptyArgs(*flUrl, *flScepUrl) {
		logger.Log("warn", "Enrollment endpoint /mdm/enroll will be disabled because you did not specify flags for the external URL or SCEP URL")
	} else {
		enrollSvc, _ := enroll.NewService(*flPushCert, *flPushPass, *flTLSCACert, *flUrl, *flScepUrl)
		enrollHandler := enroll.ServiceHandler(ctx, enrollSvc, httpLogger)
		mux.Handle("/mdm/enroll", enrollHandler)
	}

	if *flPkgRepo != "" {
		mux.Handle("/repo/", http.StripPrefix("/repo/", http.FileServer(http.Dir(*flPkgRepo))))
	}

	if *flCorsOrigin != "" {
		c := cors.New(cors.Options{
			AllowedOrigins:   []string{*flCorsOrigin},
			AllowCredentials: true,
		})

		corsHandler := c.Handler(mux)
		http.Handle("/", corsHandler)
	} else {
		http.Handle("/", mux)
	}

	http.Handle("/metrics", stdprometheus.Handler())

	serve(logger, *flTLS, *flPort, *flTLSKey, *flTLSCert)
}

func depClient(logger log.Logger, consumerKey, consumerSecret, accessToken, accessSecret, serverURL string, depsim bool) dep.Client {
	depsimDefault := &dep.Config{
		ConsumerKey:    "CK_48dd68d198350f51258e885ce9a5c37ab7f98543c4a697323d75682a6c10a32501cb247e3db08105db868f73f2c972bdb6ae77112aea803b9219eb52689d42e6",
		ConsumerSecret: "CS_34c7b2b531a600d99a0e4edcf4a78ded79b86ef318118c2f5bcfee1b011108c32d5302df801adbe29d446eb78f02b13144e323eb9aad51c79f01e50cb45c3a68",
		AccessToken:    "AT_927696831c59ba510cfe4ec1a69e5267c19881257d4bca2906a99d0785b785a6f6fdeb09774954fdd5e2d0ad952e3af52c6d8d2f21c924ba0caf4a031c158b89",
		AccessSecret:   "AS_c31afd7a09691d83548489336e8ff1cb11b82b6bca13f793344496a556b1f4972eaff4dde6deb5ac9cf076fdfa97ec97699c34d515947b9cf9ed31c99dded6ba",
	}
	var config *dep.Config
	if depsim {
		config = depsimDefault
	} else {
		if checkEmptyArgs(consumerKey, consumerSecret, accessToken, accessSecret) {
			logger.Log("err", "must specify DEP server credentials")
			logger.Log("ConsumerKey", consumerKey, "ConsumerSecret", consumerSecret, "AccessToken", accessToken, "AccessSecret", accessSecret)
			os.Exit(1)
		}
		config = &dep.Config{
			ConsumerKey:    consumerKey,
			ConsumerSecret: consumerSecret,
			AccessToken:    accessToken,
			AccessSecret:   accessSecret,
		}
	}
	var client dep.Client
	var err error
	if serverURL != "" {
		client, err = dep.NewClient(config, dep.ServerURL(serverURL))
	} else {
		client, err = dep.NewClient(config)
	}
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)

	}

	return client
}

func pushService(certPath, password string) (*push.Service, error) {
	cert, key, err := certificate.Load(certPath, password)
	if err != nil {
		return nil, err
	}
	client, err := push.NewClient(certificate.TLS(cert, key))
	if err != nil {
		return nil, err
	}
	service := &push.Service{
		Client: client,
		Host:   push.Production,
	}

	return service, nil
}

func checkEmptyArgs(args ...string) bool {
	for _, arg := range args {
		if arg == "" {
			return true
		}
	}
	return false
}

// choose http or https
func serve(logger log.Logger, tls bool, port, key, cert string) {
	portStr := fmt.Sprintf(":%v", port)
	if tls {
		logger.Log("msg", "HTTPs", "addr", port)
		logger.Log("err", http.ListenAndServeTLS(portStr, cert, key, nil))
	} else {
		logger.Log("msg", "HTTP", "addr", port)
		logger.Log("err", http.ListenAndServe(portStr, nil))
	}
}

func envString(key, def string) string {
	if env := os.Getenv(key); env != "" {
		return env
	}
	return def
}

func envBool(key string) bool {
	if env := os.Getenv(key); env == "true" {
		return true
	}
	return false
}

func checkTLSFlags(key, cert string) error {
	if key == "" || cert == "" {
		return errors.New("You must provide a valid path to a TLS cert and key")
	}
	return nil
}

func defaultPort(tls bool) string {
	if tls {
		return "443"
	}
	return "80"
}

// use this in docker container
func getPGConnFromENV(logger log.Logger, host string) string {
	user := os.Getenv("POSTGRES_ENV_POSTGRES_USER")
	if user == "" {
		user = "postgres"
	}
	dbname := os.Getenv("POSTGRES_ENV_POSTGRES_DB")
	if dbname == "" {
		dbname = user //same defaults as the docker pgcontainer
	}
	password := os.Getenv("POSTGRES_ENV_POSTGRES_PASSWORD")
	if password == "" {
		password = "postgres"
	}
	sslmode := os.Getenv("POSTGRES_ENV_SSLMODE")
	if sslmode == "" {
		logger.Log("msg", "POSTGRES_ENV_SSLMODE not specified, using 'require' by default")
		sslmode = "require"
	}
	conn := fmt.Sprintf("user=%v password=%v dbname=%v sslmode=%v host=%v", user, password, dbname, sslmode, host)
	return conn
}

func getRedisConnFromENV(host string) string {
	port := os.Getenv("REDIS_PORT_6379_TCP_PORT")
	conn := fmt.Sprintf("%v:%v", host, port)
	return conn
}
