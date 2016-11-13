package app

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/micromdm/micromdm/version"
)

func loadConfig() (*Config, error) {
	// cli flags, using environment variables as a possible default.
	var (
		// tls config
		flTLS     = flag.Bool("tls", envBool("MICROMDM_USE_TLS"), "use https")
		flTLSCert = flag.String("tls-cert", envString("MICROMDM_TLS_CERT", ""), "path to TLS certificate")
		flTLSKey  = flag.String("tls-key", envString("MICROMDM_TLS_KEY", ""), "path to TLS private key")

		// dep config
		flDEP          = flag.Bool("dep", envBool("MICROMDM_USE_DEP"), "use DEP")
		flDEPCK        = flag.String("dep-consumer-key", envString("DEP_CONSUMER_KEY", ""), "dep consumer key")
		flDEPCS        = flag.String("dep-consumer-secret", envString("DEP_CONSUMER_SECRET", ""), "dep consumer secret")
		flDEPAT        = flag.String("dep-access-token", envString("DEP_ACCESS_TOKEN", ""), "dep access token")
		flDEPAS        = flag.String("dep-access-secret", envString("DEP_ACCESS_SECRET", ""), "dep access secret")
		flDEPSim       = flag.Bool("depsim", envBool("DEP_USE_DEPSIM"), "use default depsim credentials")
		flDEPServerURL = flag.String("dep-server-url", envString("DEP_SERVER_URL", ""), "dep server url. for testing. Use blank if not running against depsim")

		// server settings
		flURL         = flag.String("url", envString("MICROMDM_URL", ""), "public url of the server")
		flAddress     = flag.String("http-address", envString("MICROMDM_HTTP_LISTEN_ADDRESS", "0.0.0.0"), "address to listen on.")
		flCORSOrigins = flag.String("cors-origin", envString("MICROMDM_CORS_ORIGINS", ""), "allowed domain for cross origin resource sharing. comma separated")
		flPkgRepo     = flag.String("pkg-repo", envString("MICROMDM_PKG_REPO", ""), "path to a folder with packages for use with InstallApplication")

		// scep config. use the cert chain from enrollment if provided.
		flSCEPURL       = flag.String("scep-url", envString("MICROMDM_SCEP_URL", ""), "scep server url. If blank, enroll profile will not use a scep payload")
		flSCEPChallenge = flag.String("scep-challenge", envString("MICROMDM_SCEP_CHALLENGE", ""), "scep server challenge")
		flSCEPSubject   = flag.String("scep-subject", envString("MICROMDM_SCEP_SUBJECT", ""), "scep request subject in microsoft string representation")

		// enrollment config.
		flEnrollment    = flag.String("enrollment-profile", envString("MICROMDM_ENROLL_PROFILE", ""), "path to enrollment profile")
		flServerCAChain = flag.String("tls-ca-cert", envString("MICROMDM_TLS_CA_CHAIN", ""), "path to CA certificate chain for trust profile")

		// postgres connection flag. will also try to load from docker link.
		flPGconn = flag.String("postgres", envString("MICROMDM_POSTGRES_CONN_URL", ""), "postgres connection url")

		// redis connection flag. will also try to load from docker link.
		flRedisConn = flag.String("redis", envString("MICROMDM_REDIS_CONN_URL", ""), "redis connection url")

		flVersion = flag.Bool("version", false, "print version information")

		// APNS config. Can be either two files or a combined .p12 (like the one exported from keychain access")
		flPushCert = flag.String("push-cert", envString("MICROMDM_PUSH_CERT", ""), "path to push certificate")
		flPushPass = flag.String("push-password", envString("MICROMDM_PUSH_PASSWORD", ""), "push certificate password")
		flPushKey  = flag.String("push-key", envString("MICROMDM_PUSH_KEY", ""), "path to push certificate private key(if not using a single .p12 file)")
	)
	flag.Parse()

	if *flVersion {
		version.PrintFull()
		os.Exit(success)
	}

	config := &Config{}
	config.loadTLS(*flTLS, *flTLSCert, *flTLSKey)
	config.loadDEPConfig(*flDEP, *flDEPSim, *flDEPCS, *flDEPCK, *flDEPAT, *flDEPAS, *flDEPServerURL)
	config.loadServerConfig(*flURL, *flAddress, *flCORSOrigins, *flPkgRepo)
	config.loadSCEPConfig(*flSCEPURL, *flSCEPChallenge, *flSCEPSubject)
	config.loadEnrollmentConfig(*flEnrollment, *flServerCAChain)
	config.loadPostgres(*flPGconn)
	config.loadRedis(*flRedisConn)
	config.loadPushConfig(*flPushCert, *flPushKey, *flPushPass)
	if config.err != nil {
		return nil, config.err
	}
	return config, nil
}

// Config holds configuration values for MicroMDM. The config values
// can be loaded from CLI flags or environment variables.
type Config struct {
	TLS        *TLSConfig
	DEP        *DEPConfig
	SCEP       *SCEPConfig
	APNS       *PushConfig
	Enrollment *EnrollmentConfig
	Redis      *RedisConfig
	Postgres   *PostgresConfig
	Server     *ServerConfig

	// the err value is part of the config struct to allow multiple
	// 'loadConfigFoo' calls in sequence, without checking if err != nil every time.
	err error
}

type TLSConfig struct {
	Enabled         bool
	CertificatePath string
	PrivateKeyPath  string
}

func (c *Config) loadTLS(enabled bool, cert, key string) {
	if c.err != nil {
		return
	}
	config := &TLSConfig{
		Enabled:         enabled,
		CertificatePath: cert,
		PrivateKeyPath:  key,
	}
	if enabled && (cert == "" || key == "") {
		c.err = errors.New("certificate or key path missing in TLS config")
		return
	}
	c.TLS = config
}

// DEPConfig holds configuration values for the DEP API.
// If UseSim is true, the default depsim values will be used instead.
// When usind depsim, a ServerURL value must be specified.
type DEPConfig struct {
	Enabled        bool
	UseSim         bool
	ServerURL      string
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}

func (c *Config) loadDEPConfig(enabled, sim bool, ck, cs, at, as, serverURL string) {
	if c.err != nil {
		return
	}
	config := &DEPConfig{
		Enabled:        enabled,
		UseSim:         sim,
		ConsumerKey:    ck,
		ConsumerSecret: cs,
		AccessToken:    at,
		AccessSecret:   as,
		ServerURL:      serverURL,
	}
	if sim {
		config.ConsumerKey = "CK_48dd68d198350f51258e885ce9a5c37ab7f98543c4a697323d75682a6c10a32501cb247e3db08105db868f73f2c972bdb6ae77112aea803b9219eb52689d42e6"
		config.ConsumerSecret = "CS_34c7b2b531a600d99a0e4edcf4a78ded79b86ef318118c2f5bcfee1b011108c32d5302df801adbe29d446eb78f02b13144e323eb9aad51c79f01e50cb45c3a68"
		config.AccessToken = "AT_927696831c59ba510cfe4ec1a69e5267c19881257d4bca2906a99d0785b785a6f6fdeb09774954fdd5e2d0ad952e3af52c6d8d2f21c924ba0caf4a031c158b89"
		config.AccessSecret = "AS_c31afd7a09691d83548489336e8ff1cb11b82b6bca13f793344496a556b1f4972eaff4dde6deb5ac9cf076fdfa97ec97699c34d515947b9cf9ed31c99dded6ba"
	}
	if sim && serverURL == "" {
		c.err = errors.New("dep-server-url must be specified when using depsim")
		return
	}
	c.DEP = config
}

// ServerConfig holds server configuration values.
// The ListenURL is used for the server listen Addr and the PublicURL
// is used for embedding the URL in profiles/notification emails etc.
// The public url is the URL used to connect to micromdm.
// TODO: PublicURL should be configured via the API/stored in DB.
type ServerConfig struct {
	PublicURL       string
	ListenURL       string
	CORSOrigins     []string
	PackageRepoPath string
}

func (c *Config) loadServerConfig(publicURL, listenURL, cors, repoPath string) {
	if c.err != nil {
		return
	}
	config := &ServerConfig{
		PublicURL:       publicURL,
		ListenURL:       listenURL,
		CORSOrigins:     strings.Split(cors, ","),
		PackageRepoPath: repoPath,
	}
	if listenURL == "" {
		config.ListenURL = "0.0.0.0:8080"
	}
	if publicURL == "" {
		config.PublicURL = "127.0.0.1:8080"
	}
	c.Server = config
}

// SCEPConfig holds configuration for the SCEP service.
// The SCEP service can be either local or remote. If the SCEP service
// is disabled, a enrollment profile with an identity certificate must be
// used instead.
type SCEPConfig struct {
	Enabled            bool
	Embedded           bool
	RemoteURL          string
	Challenge          string
	CertificateSubject string
}

func (c *Config) loadSCEPConfig(remoteURL, challenge, certSubject string) {
	if c.err != nil {
		return
	}
	config := &SCEPConfig{
		Enabled:            remoteURL != "",
		Embedded:           false, // TODO
		RemoteURL:          remoteURL,
		Challenge:          challenge,
		CertificateSubject: certSubject,
	}
	c.SCEP = config
	// FIXME add validation
}

// EnrollmentConfig holds configuration for enrollment without SCEP.
// The ProfilePath must point to a .mobileconfig file.
// An optional path to the server root certificate can also be provided. Use
// this for adding a trust profile.
type EnrollmentConfig struct {
	ProfilePath string // TODO use a template?
	CACertPath  string
}

func (c *Config) loadEnrollmentConfig(profile, certChain string) {
	if c.err != nil {
		return
	}
	config := &EnrollmentConfig{
		ProfilePath: profile,
		CACertPath:  certChain,
	}
	c.Enrollment = config
}

type PostgresConfig struct {
	Enabled    bool
	Connection string
}

func (c *Config) loadPostgres(conn string) {
	if c.err != nil {
		return
	}
	config := &PostgresConfig{
		Enabled:    true, // currently required.
		Connection: conn,
	}
	if conn == "" {
		config.fromDockerEnv()
	}
	if conn == "" {
		c.err = errors.New("must provide postgres connection string")
		return
	}
	c.Postgres = config
}

// FromDockerEnv tries to load postgres connection info from a docker link.
// The name of the link is assumed to be "postgres".
// env values taken from https://hub.docker.com/_/postgres/
func (c *PostgresConfig) fromDockerEnv() {
	host, ok := os.LookupEnv("POSTGRES_PORT_5432_TCP_ADDR")
	if !ok {
		return // don't bother with the rest.
	}
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
		sslmode = "require"
	}
	c.Connection = fmt.Sprintf("user=%v password=%v dbname=%v sslmode=%v host=%v", user, password, dbname, sslmode, host)
}

type RedisConfig struct {
	Enabled    bool
	Connection string
}

func (c *Config) loadRedis(conn string) {
	if c.err != nil {
		return
	}
	config := &RedisConfig{
		Enabled:    true, // currently required.
		Connection: conn,
	}
	if conn == "" {
		config.fromDockerEnv()
	}
	if conn == "" {
		c.err = errors.New("must provide redis connection string")
		return
	}
	c.Redis = config
}

// FromDockerEnv tries to load connection information from a linked Docker
// container. The name of the link is assumed to be "redis".
func (c *RedisConfig) fromDockerEnv() {
	host, ok := os.LookupEnv("REDIS_PORT_6379_TCP_ADDR")
	if !ok {
		return
	}
	port := os.Getenv("REDIS_PORT_6379_TCP_PORT")
	c.Connection = fmt.Sprintf("%v:%v", host, port)
}

// PushConfig holds values for connecting to APNS. The private key can be either
// a PEM or a p12 file.
type PushConfig struct {
	CertificatePath string
	PrivateKeyPath  string
	PrivateKeyPass  string
}

func (c *Config) loadPushConfig(certPath, keyPath, keyPassword string) {
	if c.err != nil {
		return
	}
	config := &PushConfig{
		CertificatePath: certPath,
		PrivateKeyPath:  keyPath,
		PrivateKeyPass:  keyPassword,
	}
	if certPath == "" {
		c.err = errors.New("must provide MDM push certificate")
		return
	}
	c.APNS = config
}
