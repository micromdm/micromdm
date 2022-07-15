package server

import (
	"context"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/micromdm/micromdm/dep"
	"github.com/micromdm/micromdm/mdm"
	"github.com/micromdm/micromdm/mdm/enroll"
	"github.com/micromdm/micromdm/platform/apns"
	apnsbuiltin "github.com/micromdm/micromdm/platform/apns/builtin"
	"github.com/micromdm/micromdm/platform/command"
	"github.com/micromdm/micromdm/platform/config"
	configbuiltin "github.com/micromdm/micromdm/platform/config/builtin"
	"github.com/micromdm/micromdm/platform/dep/sync"
	syncbuiltin "github.com/micromdm/micromdm/platform/dep/sync/builtin"
	"github.com/micromdm/micromdm/platform/device"
	devicebuiltin "github.com/micromdm/micromdm/platform/device/builtin"
	"github.com/micromdm/micromdm/platform/profile"
	profilebuiltin "github.com/micromdm/micromdm/platform/profile/builtin"
	"github.com/micromdm/micromdm/platform/pubsub"
	"github.com/micromdm/micromdm/platform/pubsub/inmem"
	"github.com/micromdm/micromdm/platform/queue"
	queueinmem "github.com/micromdm/micromdm/platform/queue/inmem"
	block "github.com/micromdm/micromdm/platform/remove"
	blockbuiltin "github.com/micromdm/micromdm/platform/remove/builtin"
	"github.com/micromdm/micromdm/workflow/webhook"

	"github.com/boltdb/bolt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/micromdm/scep/v2/challenge"
	boltchallenge "github.com/micromdm/scep/v2/challenge/bolt"
	"github.com/micromdm/scep/v2/depot"
	boltdepot "github.com/micromdm/scep/v2/depot/bolt"
	scep "github.com/micromdm/scep/v2/server"
	"github.com/pkg/errors"
)

type Server struct {
	ConfigPath             string
	Depsim                 string
	PubClient              pubsub.PublishSubscriber
	DB                     *bolt.DB
	ServerPublicURL        string
	SCEPChallenge          string
	SCEPClientValidity     int
	TLSCertPath            string
	TLSCACertPath          string
	SCEPDepot              depot.Depot
	UseDynSCEPChallenge    bool
	GenDynSCEPChallenge    bool
	SCEPChallengeDepot     challenge.Store
	ProfileDB              profile.Store
	ConfigDB               config.Store
	RemoveDB               block.Store
	CommandWebhookURL      string
	DEPClient              *dep.Client
	SyncDB                 *syncbuiltin.DB
	NoCmdHistory           bool
	ValidateSCEPIssuer     bool
	ValidateSCEPExpiration bool
	UDIDCertAuthWarnOnly   bool
	Queue                  string

	APNSPushService apns.Service
	CommandService  command.Service
	MDMService      mdm.Service
	EnrollService   enroll.Service
	SCEPService     scep.Service
	ConfigService   config.Service

	WebhooksHTTPClient *http.Client
}

func (c *Server) Setup(logger log.Logger) error {
	if err := c.setupPubSub(); err != nil {
		return err
	}

	if err := c.setupBolt(); err != nil {
		return err
	}

	if err := c.setupRemoveService(); err != nil {
		return err
	}

	if err := c.setupConfigStore(); err != nil {
		return err
	}

	if err := c.setupSCEP(logger); err != nil {
		return err
	}

	if err := c.setupPushService(logger); err != nil {
		return err
	}

	if err := c.setupCommandService(); err != nil {
		return err
	}

	if err := c.setupWebhooks(logger); err != nil {
		return err
	}

	if err := c.setupCommandQueue(logger); err != nil {
		return err
	}

	if err := c.setupDepClient(); err != nil {
		return err
	}

	if err := c.setupProfileDB(); err != nil {
		return err
	}

	err := c.setupEnrollmentService()

	return err
}

func (c *Server) setupProfileDB() error {
	profileDB, err := profilebuiltin.NewDB(c.DB)
	if err != nil {
		return err
	}
	c.ProfileDB = profileDB
	return nil
}

func (c *Server) setupPubSub() error {
	c.PubClient = inmem.NewPubSub()
	return nil
}

func (c *Server) setupWebhooks(logger log.Logger) error {
	if c.CommandWebhookURL == "" {
		return nil
	}

	ctx := context.Background()
	ww := webhook.New(c.CommandWebhookURL, c.PubClient, webhook.WithLogger(logger), webhook.WithHTTPClient(c.WebhooksHTTPClient))
	go ww.Run(ctx)
	return nil
}

func (c *Server) setupRemoveService() error {
	removeDB, err := blockbuiltin.NewDB(c.DB)
	if err != nil {
		return err
	}
	c.RemoveDB = removeDB
	return nil
}

func (c *Server) setupCommandService() error {
	commandService, err := command.New(c.PubClient)
	if err != nil {
		return err
	}
	c.CommandService = commandService
	return nil
}

func (c *Server) setupCommandQueue(logger log.Logger) error {
	var q mdm.Queue
	switch c.Queue {
	case "inmem":
		q = queueinmem.New(c.PubClient, logger)
	case "builtin":
		opts := []queue.Option{queue.WithLogger(logger)}
		if c.NoCmdHistory {
			opts = append(opts, queue.WithoutHistory())
		}
		var err error
		q, err = queue.NewQueue(c.DB, c.PubClient, opts...)
		if err != nil {
			return err
		}
	case "":
		return errors.New("empty command queue type")
	default:
		return fmt.Errorf("invalid command queue type: %s", c.Queue)
	}

	devDB, err := devicebuiltin.NewDB(c.DB)
	if err != nil {
		return errors.Wrap(err, "new device db")
	}

	var mdmService mdm.Service
	{
		svc := mdm.NewService(c.PubClient, q, devDB)
		mdmService = svc
		mdmService = block.RemoveMiddleware(c.RemoveDB)(mdmService)

		udidauthLogger := log.With(logger, "component", "udidcertauth")
		mdmService = device.UDIDCertAuthMiddleware(devDB, udidauthLogger, c.UDIDCertAuthWarnOnly)(mdmService)

		verifycertLogger := log.With(logger, "component", "verifycert")
		mdmService = VerifyCertificateMiddleware(c.ValidateSCEPIssuer, c.ValidateSCEPExpiration, c.SCEPDepot, verifycertLogger)(mdmService)
	}
	c.MDMService = mdmService

	return nil
}

func (c *Server) setupBolt() error {
	dbPath := filepath.Join(c.ConfigPath, "micromdm.db")
	db, err := bolt.Open(dbPath, 0644, &bolt.Options{Timeout: time.Second})
	if err != nil {
		return errors.Wrap(err, "opening boltdb")
	}
	c.DB = db

	return nil
}

type pushServiceCert struct {
	*x509.Certificate
	PrivateKey interface{}
}

func (c *Server) setupConfigStore() error {
	db, err := configbuiltin.NewDB(c.DB, c.PubClient)
	if err != nil {
		return err
	}
	c.ConfigDB = db
	c.ConfigService = config.New(db)

	return nil
}

func (c *Server) setupPushService(logger log.Logger) error {
	db, err := apnsbuiltin.NewDB(c.DB, c.PubClient)
	if err != nil {
		return err
	}

	service, err := apns.New(db, c.ConfigDB, c.PubClient)
	if err != nil {
		return errors.Wrap(err, "starting micromdm push service")
	}
	c.APNSPushService = apns.LoggingMiddleware(
		log.With(level.Info(logger), "component", "apns"),
	)(service)

	pushinfoWorker := apns.NewWorker(db, c.PubClient, logger)
	go pushinfoWorker.Run(context.Background())

	return nil
}

func (c *Server) setupEnrollmentService() error {
	var (
		SCEPCertificateSubject string
		err                    error
	)

	chalStore := c.SCEPChallengeDepot
	if !c.GenDynSCEPChallenge {
		chalStore = nil
	}

	// TODO: clean up order of inputs. Maybe pass *SCEPConfig as an arg?
	// but if you do, the packages are coupled, better not.
	c.EnrollService, err = enroll.NewService(
		c.ConfigDB,
		c.PubClient,
		c.ServerPublicURL+"/scep",
		c.SCEPChallenge,
		c.ServerPublicURL,
		c.TLSCertPath,
		SCEPCertificateSubject,
		c.TLSCACertPath,
		c.ProfileDB,
		chalStore,
	)
	return errors.Wrap(err, "setting up enrollment service")
}

func (c *Server) setupDepClient() error {
	var (
		conf           dep.OAuthParameters
		depsim         = c.Depsim
		hasTokenConfig bool
		opts           []dep.Option
	)

	// try getting the oauth config from bolt
	tokens, err := c.ConfigDB.DEPTokens()
	if err != nil {
		return err
	}
	if len(tokens) >= 1 {
		hasTokenConfig = true
		conf.ConsumerSecret = tokens[0].ConsumerSecret
		conf.ConsumerKey = tokens[0].ConsumerKey
		conf.AccessSecret = tokens[0].AccessSecret
		conf.AccessToken = tokens[0].AccessToken
		// TODO: handle expiration
	}

	// override with depsim keys if specified on CLI
	if depsim != "" {
		hasTokenConfig = true
		conf = dep.OAuthParameters{
			ConsumerKey:    "CK_48dd68d198350f51258e885ce9a5c37ab7f98543c4a697323d75682a6c10a32501cb247e3db08105db868f73f2c972bdb6ae77112aea803b9219eb52689d42e6",
			ConsumerSecret: "CS_34c7b2b531a600d99a0e4edcf4a78ded79b86ef318118c2f5bcfee1b011108c32d5302df801adbe29d446eb78f02b13144e323eb9aad51c79f01e50cb45c3a68",
			AccessToken:    "AT_927696831c59ba510cfe4ec1a69e5267c19881257d4bca2906a99d0785b785a6f6fdeb09774954fdd5e2d0ad952e3af52c6d8d2f21c924ba0caf4a031c158b89",
			AccessSecret:   "AS_c31afd7a09691d83548489336e8ff1cb11b82b6bca13f793344496a556b1f4972eaff4dde6deb5ac9cf076fdfa97ec97699c34d515947b9cf9ed31c99dded6ba",
		}
		depsimurl, err := url.Parse(depsim)
		if err != nil {
			return err
		}
		opts = append(opts, dep.WithServerURL(depsimurl))
	}

	if !hasTokenConfig {
		return nil
	}

	c.DEPClient = dep.NewClient(conf, opts...)
	return nil
}

func (c *Server) CreateDEPSyncer(logger log.Logger) (sync.Syncer, error) {
	client := c.DEPClient
	opts := []sync.Option{
		sync.WithLogger(log.With(logger, "component", "depsync")),
	}
	if client != nil {
		opts = append(opts, sync.WithClient(client))
	}

	syncdb, err := syncbuiltin.NewDB(c.DB)
	if err != nil {
		return nil, err
	}
	c.SyncDB = syncdb

	var syncer sync.Syncer
	syncer, err = sync.NewWatcher(c.SyncDB, c.PubClient, opts...)
	if err != nil {
		return nil, err
	}
	return syncer, nil
}

func (c *Server) setupSCEP(logger log.Logger) error {
	svcBoltDepot, err := boltdepot.NewBoltDepot(c.DB)
	if err != nil {
		return err
	}
	c.SCEPDepot = svcBoltDepot

	key, err := svcBoltDepot.CreateOrLoadKey(2048)
	if err != nil {
		return err
	}

	crt, err := svcBoltDepot.CreateOrLoadCA(key, 5, "MicroMDM", "US")
	if err != nil {
		return err
	}

	var signer scep.CSRSigner = depot.NewSigner(
		c.SCEPDepot,
		depot.WithAllowRenewalDays(0),
		depot.WithValidityDays(c.SCEPClientValidity),
	)
	if c.UseDynSCEPChallenge {
		c.SCEPChallengeDepot, err = boltchallenge.NewBoltDepot(c.DB)
		if err != nil {
			return err
		}
		signer = challenge.Middleware(c.SCEPChallengeDepot, signer)
	} else {
		signer = scep.ChallengeMiddleware(c.SCEPChallenge, signer)
	}

	c.SCEPService, err = scep.NewService(crt, key, signer)
	if err != nil {
		return err
	}
	c.SCEPService = scep.NewLoggingService(logger, c.SCEPService)

	return nil
}
