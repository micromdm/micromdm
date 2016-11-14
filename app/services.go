package app

import (
	"crypto/tls"
	"io/ioutil"

	pushcertificate "github.com/RobotsAndPencils/buford/certificate"
	"github.com/RobotsAndPencils/buford/push"
	"github.com/garyburd/redigo/redis"
	"github.com/go-kit/kit/log"

	"github.com/micromdm/dep"
	"github.com/micromdm/micromdm/application"
	"github.com/micromdm/micromdm/certificate"
	"github.com/micromdm/micromdm/checkin"
	"github.com/micromdm/micromdm/command"
	cmdredis "github.com/micromdm/micromdm/command/service/redis"
	"github.com/micromdm/micromdm/connect"
	"github.com/micromdm/micromdm/device"
	"github.com/micromdm/micromdm/driver"
	"github.com/micromdm/micromdm/enroll"
	"github.com/micromdm/micromdm/management"
	"github.com/micromdm/micromdm/workflow"
)

// setupServices uses the values from the config to set up the various components
// which MicroMDM relies on.
func setupServices(config *Config, logger log.Logger) (*serviceManager, error) {
	sm := &serviceManager{Config: config, logger: logger}
	sm.createRedisPool()

	sm.setupAppDatastore()
	sm.setupDeviceDatastore()
	sm.setupWorkflowDatastore()
	sm.setupCertificateDatastore()

	sm.setupPushService()

	sm.setupCommandService()
	sm.setupManagementService()
	sm.setupCheckinService()
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

	CommandService    command.Service
	ManagementService management.Service
	CheckinService    checkin.Service
	ConnectService    connect.Service
	EnrollmentService enroll.Service

	*Config
	pool   *redis.Pool
	logger log.Logger
	err    error
}

func (s *serviceManager) setupEnrollmentService() {
	if s.err != nil {
		return
	}
	// TODO: clean up order of inputs. Maybe pass *SCEPConfig as an arg?
	// but if you do, the packages are coupled, better not.
	s.EnrollmentService, s.err = enroll.NewService(
		s.APNS.CertificatePath,
		s.APNS.PrivateKeyPass,
		s.APNS.PrivateKeyPath,
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
	var enrollmentProfile []byte
	// TODO make this optional - checkin.WithEnrollmentProfile([]byte)
	if s.DEP.Enabled {
		enrollmentProfile, s.err = ioutil.ReadFile(s.Enrollment.ProfilePath)
		if s.err != nil {
			return
		}
	}

	s.CheckinService = checkin.NewService(
		s.DeviceDatastore,
		s.ManagementService,
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
	var tlscert tls.Certificate
	var err error
	if s.APNS.PrivateKeyPath == "" {
		// note that buford does validity checks where
		// our direct certificate parsing does not
		cert, key, err := pushcertificate.Load(
			s.APNS.CertificatePath,
			s.APNS.PrivateKeyPass,
		)
		if err != nil {
			s.err = err
			return
		}
		tlscert = pushcertificate.TLS(cert, key)
	} else {
		tlscert, err = tls.LoadX509KeyPair(
			s.APNS.CertificatePath,
			s.APNS.PrivateKeyPath,
		)
		if err != nil {
			s.err = err
			return
		}
	}
	client, err := push.NewClient(tlscert)
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
	if s.err != nil {
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
	s.CommandService, s.err = cmdredis.NewCommandService(s.pool, s.logger)
}
