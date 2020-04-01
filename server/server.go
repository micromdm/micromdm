package server

import (
    "fmt"
    "reflect"
	
	"context"
	"crypto/x509"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	scepstore "github.com/micromdm/micromdm/platform/scep"
	scepbuiltin "github.com/micromdm/micromdm/platform/scep/builtin"
	scepmysql "github.com/micromdm/micromdm/platform/scep/mysql"
	challengestore "github.com/micromdm/scep/challenge/bolt"
	scep "github.com/micromdm/scep/server"
	"github.com/pkg/errors"
	
	"github.com/jmoiron/sqlx"
	"github.com/kolide/kit/dbutil"

	"github.com/micromdm/micromdm/dep"
	"github.com/micromdm/micromdm/mdm"
	"github.com/micromdm/micromdm/mdm/enroll"
	"github.com/micromdm/micromdm/platform/apns"
	apnsbuiltin "github.com/micromdm/micromdm/platform/apns/builtin"
	apnsmysql "github.com/micromdm/micromdm/platform/apns/mysql"
	"github.com/micromdm/micromdm/platform/command"
	"github.com/micromdm/micromdm/platform/config"
	configbuiltin "github.com/micromdm/micromdm/platform/config/builtin"
	configmysql "github.com/micromdm/micromdm/platform/config/mysql"
	
	"github.com/micromdm/micromdm/platform/dep/sync"
	syncbuiltin "github.com/micromdm/micromdm/platform/dep/sync/builtin"
	syncmysql "github.com/micromdm/micromdm/platform/dep/sync/mysql"
	
	"github.com/micromdm/micromdm/platform/device"
	devicebuiltin "github.com/micromdm/micromdm/platform/device/builtin"
	devicemysql "github.com/micromdm/micromdm/platform/device/mysql"
	"github.com/micromdm/micromdm/platform/profile"
	profilebuiltin "github.com/micromdm/micromdm/platform/profile/builtin"
	profilemysql "github.com/micromdm/micromdm/platform/profile/mysql"
	"github.com/micromdm/micromdm/platform/pubsub"
	"github.com/micromdm/micromdm/platform/pubsub/inmem"
	
	"github.com/micromdm/micromdm/platform/queue"
	queueBuiltin "github.com/micromdm/micromdm/platform/queue/builtin"
	queueMysql "github.com/micromdm/micromdm/platform/queue/mysql"
	
	
	block "github.com/micromdm/micromdm/platform/remove"
	blockbuiltin "github.com/micromdm/micromdm/platform/remove/builtin"
	blockmysql "github.com/micromdm/micromdm/platform/remove/mysql"
	"github.com/micromdm/micromdm/workflow/webhook"
)

type Server struct {
	ConfigPath        string
	Depsim            string
	PubClient         pubsub.PublishSubscriber
	DB                *bolt.DB
	ServerPublicURL   string
	SCEPChallenge     string
	SCEPClientValidity int
	TLSCertPath       string
	
	UseDynSCEPChallenge bool
	GenDynSCEPChallenge bool
	SCEPChallengeDepot  *challengestore.Depot
	SCEPBuiltin	      *scepbuiltin.Depot
	SCEPMysqlDB       *scepmysql.Depot
	
	SCEPDB			  scepstore.Store
	ProfileDB         profile.Store
	ConfigDB          config.Store
	RemoveDB          block.Store
	CommandWebhookURL string
	DEPClient         *dep.Client
	NoCmdHistory       bool
	
	SyncBuiltin       *syncbuiltin.DB
	SyncMysqlDB       *syncmysql.Mysql
	
	MysqlDB			  *sqlx.DB
	MysqlUsername     string
	MysqlPassword     string
	MysqlDatabase     string
	MysqlHost 	      string
	MysqlPort 	      string

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

	if err := c.setupMysql(); err != nil {
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
	var profileDB profile.Store
	var err error
	
	// If Mysql is set up, use Mysql, else use Bolt as Fallback
	if c.MysqlDB != nil {
		profileDB, err = profilemysql.NewDB(c.MysqlDB)
	} else {
		profileDB, err = profilebuiltin.NewDB(c.DB)
	}
	
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
	var removeDB block.Store
	var err error

	// If Mysql is set up, use Mysql, else use Bolt as Fallback
	if c.MysqlDB != nil {
		removeDB, err = blockmysql.NewDB(c.MysqlDB)
	} else {
		removeDB, err = blockbuiltin.NewDB(c.DB)
	}
	
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
	var q queue.Store
	var devDB device.UDIDCertAuthStore
	var err error

	// If Mysql is set up, use Mysql, else use Bolt as Fallback
	if c.MysqlDB != nil {
		opts := []queueMysql.Option{queueMysql.WithLogger(logger)}
		if c.NoCmdHistory {
			opts = append(opts, queueMysql.WithoutHistory())
		}

		q, err = queueMysql.NewQueue(c.MysqlDB, c.PubClient, opts...)
		if err != nil {
			return err
		}

		devDB, err = devicemysql.NewDB(c.MysqlDB)
	} else {
		opts := []queueBuiltin.Option{queueBuiltin.WithLogger(logger)}
		if c.NoCmdHistory {
			opts = append(opts, queueBuiltin.WithoutHistory())
		}
		
		q, err = queueBuiltin.NewQueue(c.DB, c.PubClient, opts...)
		if err != nil {
			return err
		}

		devDB, err = devicebuiltin.NewDB(c.DB)
	}

	
	if err != nil {
		return errors.Wrap(err, "new device db")
	}

	var mdmService mdm.Service
	{
		svc := mdm.NewService(c.PubClient, q)
		mdmService = svc
		mdmService = block.RemoveMiddleware(c.RemoveDB)(mdmService)

		udidauthLogger := log.With(logger, "component", "udidcertauth")
		mdmService = device.UDIDCertAuthMiddleware(devDB, udidauthLogger)(mdmService)
	
		verifycertLogger := log.With(logger, "component", "verifycert")
		mdmService = VerifyCertificateMiddleware(c.SCEPDB, verifycertLogger)(mdmService)
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

func (c *Server) setupMysql() error {
	
	if (c.MysqlUsername == "" || c.MysqlPassword == "" || c.MysqlHost == "" || c.MysqlPort == "" || c.MysqlDatabase == "") {
		return nil
	}

	
	db, err := dbutil.OpenDBX(
		"mysql",
//		"micromdm:micromdm@tcp(127.0.0.1:3306)/micromdm_test?parseTime=true",
		c.MysqlUsername+":"+c.MysqlPassword+"@tcp("+c.MysqlHost+":"+c.MysqlPort+")/"+c.MysqlDatabase+"?parseTime=true",
		dbutil.WithLogger(log.NewNopLogger()),
		dbutil.WithMaxAttempts(1),
	)
	if err != nil {
		return errors.Wrap(err, "opening mysql")
	}
	c.MysqlDB = db
	
	// Set the number of open and idle connection to a maximum total of 2.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	
	return nil
}

type pushServiceCert struct {
	*x509.Certificate
	PrivateKey interface{}
}

func (c *Server) setupConfigStore() error {
	var db config.Store
	var err error

	// If Mysql is set up, use Mysql, else use Bolt as Fallback
	if c.MysqlDB != nil {
		db, err = configmysql.NewDB(c.MysqlDB, c.PubClient)
	} else {
		db, err = configbuiltin.NewDB(c.DB, c.PubClient)
	}
	
	if err != nil {
		return err
	}
	c.ConfigDB = db
	c.ConfigService = config.New(db)

	return nil
}

func (c *Server) setupPushService(logger log.Logger) error {
	var service apns.Service

	// Code duplication is needed unfortunately...
	// workerStore may be of different types. (Mysql vs. Bolt)
	// If Mysql is set up, use Mysql, else use Bolt as Fallback
	if c.MysqlDB != nil {
		workerStore, err := apnsmysql.NewDB(c.MysqlDB, c.PubClient)
		if err != nil {
			return err
		}
	
		service, err := apns.New(workerStore, c.ConfigDB, c.PubClient)

		if err != nil {
			return errors.Wrap(err, "starting micromdm push service")
		}
		c.APNSPushService = apns.LoggingMiddleware(
			log.With(level.Info(logger), "component", "apns"),
		)(service)
		
		fmt.Println(reflect.TypeOf(workerStore))

		pushinfoWorker := apns.NewWorker(workerStore, c.PubClient, logger)
		go pushinfoWorker.Run(context.Background())
		
	} else {
		workerStore, err := apnsbuiltin.NewDB(c.DB, c.PubClient)
		if err != nil {
			return err
		}
		service, err = apns.New(workerStore, c.ConfigDB, c.PubClient)

		if err != nil {
			return errors.Wrap(err, "starting micromdm push service")
		}
		c.APNSPushService = apns.LoggingMiddleware(
			log.With(level.Info(logger), "component", "apns"),
		)(service)

		pushinfoWorker := apns.NewWorker(workerStore, c.PubClient, logger)
		go pushinfoWorker.Run(context.Background())
	}

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
	ctx := context.Background()
	tokens, err := c.ConfigDB.DEPTokens(ctx)
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
	
	
	var err error
	
	if c.MysqlDB != nil {
		syncdb, err := syncmysql.NewDB(c.MysqlDB)
		if err != nil {
			return nil, err
		}

		c.SyncMysqlDB = syncdb
	} else {
		syncdb, err := syncbuiltin.NewDB(c.DB)
		if err != nil {
			return nil, err
		}

		c.SyncBuiltin = syncdb
	}
	
	var syncer sync.Syncer
	if c.SyncMysqlDB != nil {
		syncer, err = sync.NewWatcher(c.SyncMysqlDB, c.PubClient, opts...)
	} else {
		syncer, err = sync.NewWatcher(c.SyncBuiltin, c.PubClient, opts...)
	}
	
	if err != nil {
		return nil, err
	}
	return syncer, nil
}

func (c *Server) setupSCEP(logger log.Logger) error {
	var store scepstore.Store
	var err error
	
	opts := []scep.ServiceOption{
		scep.ClientValidity(c.SCEPClientValidity),
	}
	
	var scepChalOpt scep.ServiceOption
	if c.UseDynSCEPChallenge {
		c.SCEPChallengeDepot, err = challengestore.NewBoltDepot(c.DB)
		if err != nil {
			return err
		}
		scepChalOpt = scep.WithDynamicChallenges(c.SCEPChallengeDepot)
	} else {
		scepChalOpt = scep.ChallengePassword(c.SCEPChallenge)
	}
	opts = append(opts, scepChalOpt)
	
	// If Mysql is set up, use Mysql, else use Bolt as Fallback
	if c.MysqlDB != nil {
		c.SCEPMysqlDB, err = scepmysql.NewDB(c.MysqlDB)
		store = c.SCEPMysqlDB
		
		key, err := store.CreateOrLoadKey(2048)
		if err != nil {
			return err
		}
	
		fmt.Println("CreateOrLoadCA")
		_, err = store.CreateOrLoadCA(key, 5, "MicroMDM", "US")
		if err != nil {
			return err
		}
	
		c.SCEPService, err = scep.NewService(store, opts...)
	} else {
		c.SCEPBuiltin, err = scepbuiltin.NewDB(c.DB)
		store = c.SCEPBuiltin
		
		key, err := store.CreateOrLoadKey(2048)
		if err != nil {
			return err
		}
	
		fmt.Println("CreateOrLoadCA")
		_, err = store.CreateOrLoadCA(key, 5, "MicroMDM", "US")
		if err != nil {
			return err
		}
		
		c.SCEPService, err = scep.NewService(store, opts...)
	}
	
	if err != nil {
		return err
	}

	c.SCEPDB = store
	c.SCEPService = scep.NewLoggingService(logger, c.SCEPService)

	return nil
}
