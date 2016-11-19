package app

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"io/ioutil"

	pushcertificate "github.com/RobotsAndPencils/buford/certificate"
	"github.com/RobotsAndPencils/buford/push"
	"github.com/boltdb/bolt"
	"github.com/garyburd/redigo/redis"
	"github.com/go-kit/kit/log"
	"github.com/micromdm/dep"
	nsq "github.com/nsqio/go-nsq"
	"github.com/nsqio/nsq/nsqd"

	"github.com/micromdm/micromdm/application"
	"github.com/micromdm/micromdm/certificate"
	"github.com/micromdm/micromdm/checkin"
	"github.com/micromdm/micromdm/checkin/archive"
	eventcheckin "github.com/micromdm/micromdm/checkin/service/event"
	"github.com/micromdm/micromdm/command"
	cmdbolt "github.com/micromdm/micromdm/command/service/bolt"
	cmdredis "github.com/micromdm/micromdm/command/service/redis"
	"github.com/micromdm/micromdm/connect"
	"github.com/micromdm/micromdm/depenroll"
	"github.com/micromdm/micromdm/device"
	"github.com/micromdm/micromdm/driver"
	"github.com/micromdm/micromdm/enroll"
	"github.com/micromdm/micromdm/event"
	"github.com/micromdm/micromdm/management"
	"github.com/micromdm/micromdm/workflow"

	"golang.org/x/crypto/pkcs12"
)

// setupServices uses the values from the config to set up the various components
// which MicroMDM relies on.
func setupServices(config *Config, logger log.Logger) (*serviceManager, error) {
	sm := &serviceManager{Config: config, logger: logger}
	sm.createRedisPool()
	sm.setupEventPublisher()
	sm.setupBoltStores()

	stop := make(chan bool)
	go func() {
		err := archive.Persist(stop, sm.commandStore)
		if err != nil {
			panic(err)
		}
	}()

	sm.setupAppDatastore()
	sm.setupDeviceDatastore()
	sm.setupWorkflowDatastore()
	sm.setupCertificateDatastore()

	sm.loadPushCerts()
	sm.setupPushService()

	sm.setupCommandService()
	sm.setupManagementService()
	sm.setupCheckinService()
	sm.setupDEPEnrollmentService()
	sm.setupConnectService()
	sm.setupEnrollmentService()
	if sm.err != nil {
		return nil, sm.err
	}
	return sm, nil
}

// serviceManager knows how to setup the independent components which make up
// MicroMDM, mainly Datastores and Services.
type serviceManager struct {
	CertificateDatastore certificate.Datastore
	DeviceDatastore      device.Datastore
	WorkflowDatastore    workflow.Datastore
	ApplicationDatastore application.Datastore

	PushService *push.Service
	pushServiceCert
	Publisher *event.Publisher

	CommandService       command.Service
	ManagementService    management.Service
	CheckinService       checkin.Service
	ConnectService       connect.Service
	EnrollmentService    enroll.Service
	DEPEnrollmentService depenroll.Service

	*Config
	pool         *redis.Pool
	commandStore *bolt.DB
	logger       log.Logger
	err          error
}

type pushServiceCert struct {
	*x509.Certificate
	PrivateKey interface{}
}

func (s *serviceManager) setupBoltStores() {
	if s.err != nil {
		return
	}
	s.commandStore, s.err = bolt.Open("db.bolt", 0777, nil)
}

func (s *serviceManager) loadPushCerts() {
	if s.err != nil {
		return
	}

	if s.APNS.PrivateKeyPath == "" {
		var pkcs12Data []byte
		pkcs12Data, s.err = ioutil.ReadFile(s.APNS.CertificatePath)
		if s.err != nil {
			return
		}
		s.pushServiceCert.PrivateKey, s.pushServiceCert.Certificate, s.err =
			pkcs12.Decode(pkcs12Data, s.APNS.PrivateKeyPass)
		return
	}

	var pemData []byte
	pemData, s.err = ioutil.ReadFile(s.APNS.CertificatePath)
	if s.err != nil {
		return
	}

	pemBlock, _ := pem.Decode(pemData)
	if pemBlock == nil {
		s.err = errors.New("invalid PEM data for cert")
		return
	}
	s.pushServiceCert.Certificate, s.err = x509.ParseCertificate(pemBlock.Bytes)
	if s.err != nil {
		return
	}

	pemData, s.err = ioutil.ReadFile(s.APNS.PrivateKeyPath)
	if s.err != nil {
		return
	}

	pemBlock, _ = pem.Decode(pemData)
	if pemBlock == nil {
		s.err = errors.New("invalid PEM data for privkey")
		return
	}
	s.pushServiceCert.PrivateKey, s.err =
		x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
}

var oidASN1UserID = asn1.ObjectIdentifier{0, 9, 2342, 19200300, 100, 1, 1}

func topicFromCert(cert *x509.Certificate) (string, error) {
	for _, v := range cert.Subject.Names {
		if v.Type.Equal(oidASN1UserID) {
			return v.Value.(string), nil
		}
	}

	return "", errors.New("Could not find Push Topic (UserID OID) in certificate")
}

func (s *serviceManager) setupEnrollmentService() {
	if s.err != nil {
		return
	}
	pushTopic, err := topicFromCert(s.pushServiceCert.Certificate)
	if err != nil {
		s.err = err
		return
	}
	// TODO: clean up order of inputs. Maybe pass *SCEPConfig as an arg?
	// but if you do, the packages are coupled, better not.
	s.EnrollmentService, s.err = enroll.NewService(
		pushTopic,
		s.Enrollment.CACertPath,
		s.SCEP.RemoteURL,
		s.SCEP.Challenge,
		s.Server.PublicURL,
		s.TLS.CertificatePath,
		s.SCEP.CertificateSubject,
	)
}

func (s *serviceManager) setupConnectService() {
	if s.err != nil {
		return
	}
	s.ConnectService = connect.NewService(
		s.DeviceDatastore,
		s.ApplicationDatastore,
		s.CertificateDatastore,
		s.CommandService,
	)
}

func (s *serviceManager) setupCheckinService() {
	if s.err != nil {
		return
	}

	/*
		s.CheckinService = direct.NewService(
			s.DeviceDatastore,
			s.ManagementService,
		)
	*/
	s.CheckinService = eventcheckin.NewCheckinService(s.Publisher)

}

func (s *serviceManager) setupEventPublisher() {
	if s.err != nil {
		return
	}
	done := make(chan bool)
	go func() {
		opts := nsqd.NewOptions()
		nsqd := nsqd.New(opts)
		nsqd.Main()

		// wait until we are told to continue and exit
		<-done
		nsqd.Exit()
	}()

	cfg := nsq.NewConfig()
	producer, err := nsq.NewProducer("localhost:4150", cfg)
	if err != nil {
		s.err = err
	}
	s.Publisher = &event.Publisher{Producer: producer}
}

func (s *serviceManager) setupDEPEnrollmentService() {
	if s.err != nil || !s.DEP.Enabled {
		return
	}
	// TODO make this optional - depenroll.WithEnrollmentProfile([]byte)
	var enrollmentProfile []byte
	if s.Enrollment.ProfilePath != "" {
		enrollmentProfile, s.err = ioutil.ReadFile(s.Enrollment.ProfilePath)
		if s.err != nil {
			return
		}
	}
	s.DEPEnrollmentService = depenroll.NewService(
		s.DeviceDatastore,
		s.CommandService,
		enrollmentProfile,
	)
}

func (s *serviceManager) setupManagementService() {
	if s.err != nil {
		return
	}
	var client dep.Client
	if s.DEP.Enabled {
		var opts []func(*dep.Config) error
		config := &dep.Config{
			ConsumerKey:    s.DEP.ConsumerKey,
			ConsumerSecret: s.DEP.ConsumerSecret,
			AccessToken:    s.DEP.AccessToken,
			AccessSecret:   s.DEP.AccessSecret,
		}
		if s.DEP.UseSim {
			opts = append(opts, dep.ServerURL(s.DEP.ServerURL))
		}
		client, s.err = dep.NewClient(config, opts...)
		if s.err != nil {
			return
		}
	}
	s.ManagementService = management.NewService(
		s.DeviceDatastore,
		s.WorkflowDatastore,
		client,
		s.PushService,
		s.ApplicationDatastore,
		s.CertificateDatastore,
	)
}

func (s *serviceManager) setupPushService() {
	if s.err != nil {
		return
	}
	client, err := push.NewClient(pushcertificate.TLS(
		s.pushServiceCert.Certificate,
		s.pushServiceCert.PrivateKey.(*rsa.PrivateKey),
	))
	if err != nil {
		s.err = err
		return
	}
	s.PushService = &push.Service{
		Client: client,
		Host:   push.Production,
	}

}

func (s *serviceManager) setupCertificateDatastore() {
	if s.err != nil {
		return
	}
	db, err := certificate.NewDB("postgres", s.Postgres.Connection, s.logger)
	if err != nil {
		s.err = err
		return
	}
	s.CertificateDatastore = db
}

func (s *serviceManager) setupAppDatastore() {
	if s.err != nil {
		return
	}
	db, err := application.NewDB("postgres", s.Postgres.Connection, s.logger)
	if err != nil {
		s.err = err
		return
	}
	s.ApplicationDatastore = db
}

func (s *serviceManager) setupWorkflowDatastore() {
	if s.err != nil {
		return
	}
	db, err := workflow.NewDB("postgres", s.Postgres.Connection, s.logger)
	if err != nil {
		s.err = err
		return
	}
	s.WorkflowDatastore = db
}

func (s *serviceManager) setupDeviceDatastore() {
	if s.err != nil {
		return
	}
	db, err := device.NewDB("postgres", s.Postgres.Connection, s.logger)
	if err != nil {
		s.err = err
		return
	}
	s.DeviceDatastore = db
}

func (s *serviceManager) createRedisPool() {
	if s.err != nil || !s.Redis.Enabled {
		return
	}
	opts := []driver.ConnOption{driver.Logger(s.logger)}
	if s.Redis.Password != "" {
		opts = append(opts, driver.WithPassword(s.Redis.Password))
	}
	s.pool, s.err = driver.NewRedisPool(s.Redis.Connection, opts...)
}

func (s *serviceManager) setupCommandService() {
	if s.err != nil {
		return
	}
	if s.Redis.Enabled {
		s.CommandService, s.err = cmdredis.NewCommandService(s.pool, s.logger)
		return
	}
	s.CommandService = cmdbolt.NewCommandService(s.commandStore)
}
