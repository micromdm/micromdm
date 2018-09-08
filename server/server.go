package server

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/RobotsAndPencils/buford/push"
	"github.com/boltdb/bolt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	boltdepot "github.com/micromdm/scep/depot/bolt"
	scep "github.com/micromdm/scep/server"
	"github.com/pkg/errors"
	"golang.org/x/crypto/pkcs12"

	"github.com/micromdm/micromdm/dep"
	"github.com/micromdm/micromdm/mdm"
	"github.com/micromdm/micromdm/mdm/enroll"
	"github.com/micromdm/micromdm/pkg/crypto"
	"github.com/micromdm/micromdm/platform/apns"
	apnsbuiltin "github.com/micromdm/micromdm/platform/apns/builtin"
	"github.com/micromdm/micromdm/platform/command"
	"github.com/micromdm/micromdm/platform/config"
	configbuiltin "github.com/micromdm/micromdm/platform/config/builtin"
	"github.com/micromdm/micromdm/platform/dep/sync"
	"github.com/micromdm/micromdm/platform/profile"
	profilebuiltin "github.com/micromdm/micromdm/platform/profile/builtin"
	"github.com/micromdm/micromdm/platform/pubsub"
	"github.com/micromdm/micromdm/platform/pubsub/inmem"
	"github.com/micromdm/micromdm/platform/queue"
	block "github.com/micromdm/micromdm/platform/remove"
	blockbuiltin "github.com/micromdm/micromdm/platform/remove/builtin"
	"github.com/micromdm/micromdm/workflow/webhook"
)

type Server struct {
	ConfigPath          string
	Depsim              string
	PubClient           pubsub.PublishSubscriber
	DB                  *bolt.DB
	PushCert            pushServiceCert
	ServerPublicURL     string
	SCEPChallenge       string
	APNSPrivateKeyPath  string
	APNSCertificatePath string
	APNSPrivateKeyPass  string
	TLSCertPath         string
	SCEPDepot           *boltdepot.Depot
	ProfileDB           profile.Store
	ConfigDB            config.Store
	RemoveDB            block.Store
	CommandWebhookURL   string
	DEPClient           *dep.Client

	PushService     *push.Service // bufford push
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

	if err := c.loadPushCerts(); err != nil {
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

	if err := c.setupEnrollmentService(); err != nil {
		return err
	}

	return nil
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
	q, err := queue.NewQueue(c.DB, c.PubClient)
	if err != nil {
		return err
	}

	var mdmService mdm.Service
	{
		svc := mdm.NewService(c.PubClient, q)
		mdmService = svc
		mdmService = block.RemoveMiddleware(c.RemoveDB)(mdmService)
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

func (c *Server) loadPushCerts() error {
	if c.APNSCertificatePath == "" && c.APNSPrivateKeyPass == "" && c.APNSPrivateKeyPath == "" {
		// this is optional, config could also be provided with mdmctl
		return nil
	}
	var err error

	if c.APNSPrivateKeyPath == "" {
		var pkcs12Data []byte

		pkcs12Data, err = ioutil.ReadFile(c.APNSCertificatePath)
		if err != nil {
			return err
		}
		c.PushCert.PrivateKey, c.PushCert.Certificate, err =
			pkcs12.Decode(pkcs12Data, c.APNSPrivateKeyPass)
		return err
	}

	c.PushCert.Certificate, err = crypto.ReadPEMCertificateFile(c.APNSCertificatePath)
	if err != nil {
		return err
	}

	var pemData []byte
	pemData, err = ioutil.ReadFile(c.APNSPrivateKeyPath)
	if err != nil {
		return err
	}

	pkeyBlock := new(bytes.Buffer)
	pemBlock, _ := pem.Decode(pemData)
	if pemBlock == nil {
		return errors.New("invalid PEM data for privkey")
	}

	if x509.IsEncryptedPEMBlock(pemBlock) {
		b, err := x509.DecryptPEMBlock(pemBlock, []byte(c.APNSPrivateKeyPass))
		if err != nil {
			return fmt.Errorf("decrypting DES private key %s", err)
		}
		pkeyBlock.Write(b)
	} else {
		pkeyBlock.Write(pemBlock.Bytes)
	}

	priv, err := x509.ParsePKCS1PrivateKey(pkeyBlock.Bytes())
	if err != nil {
		return fmt.Errorf("parsing pkcs1 private key: %s", err)
	}
	c.PushCert.PrivateKey = priv

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
	var opts []apns.Option
	{
		cert, _ := c.ConfigDB.PushCertificate()
		if c.PushCert.Certificate != nil && cert == nil {
			cert = &tls.Certificate{
				Certificate: [][]byte{c.PushCert.Certificate.Raw},
				PrivateKey:  c.PushCert.PrivateKey,
				Leaf:        c.PushCert.Certificate,
			}
		}
		if cert == nil {
			goto after
		}
		client, err := push.NewClient(*cert)
		if err != nil {
			return err
		}
		svc := push.NewService(client, push.Production)
		opts = append(opts, apns.WithPushService(svc))
	}
after:

	db, err := apnsbuiltin.NewDB(c.DB, c.PubClient)
	if err != nil {
		return err
	}

	service, err := apns.New(db, c.ConfigDB, c.PubClient, opts...)
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
	var topicProvider enroll.TopicProvider
	var err error
	if c.PushCert.Certificate != nil {
		pushTopic, err := crypto.TopicFromCert(c.PushCert.Certificate)
		if err != nil {
			return errors.Wrap(err, "get apns topic from certificate")
		}
		topicProvider = staticTopicProvider{topic: pushTopic}
	} else {
		topicProvider = c.ConfigDB
	}

	var SCEPCertificateSubject string
	// TODO: clean up order of inputs. Maybe pass *SCEPConfig as an arg?
	// but if you do, the packages are coupled, better not.
	c.EnrollService, err = enroll.NewService(
		topicProvider,
		c.PubClient,
		c.ServerPublicURL+"/scep",
		c.SCEPChallenge,
		c.ServerPublicURL,
		c.TLSCertPath,
		SCEPCertificateSubject,
		c.ProfileDB,
	)
	if err != nil {
		return err
	}

	return nil
}

// if the apns-cert flags are specified this provider will be used in the enroll service.
type staticTopicProvider struct{ topic string }

func (p staticTopicProvider) PushTopic() (string, error) {
	return p.topic, nil
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

	var syncer sync.Syncer
	syncer, err := sync.New(c.PubClient, c.DB, logger, opts...)
	if err != nil {
		return nil, err
	}
	return syncer, nil
}

func (c *Server) setupSCEP(logger log.Logger) error {
	depot, err := boltdepot.NewBoltDepot(c.DB)
	if err != nil {
		return err
	}

	key, err := depot.CreateOrLoadKey(2048)
	if err != nil {
		return err
	}

	_, err = depot.CreateOrLoadCA(key, 5, "MicroMDM", "US")
	if err != nil {
		return err
	}

	opts := []scep.ServiceOption{
		scep.ClientValidity(365),
		scep.ChallengePassword(c.SCEPChallenge),
	}
	c.SCEPDepot = depot
	c.SCEPService, err = scep.NewService(depot, opts...)
	if err != nil {
		return err
	}
	c.SCEPService = scep.NewLoggingService(logger, c.SCEPService)

	return nil
}
