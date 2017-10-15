package push

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/RobotsAndPencils/buford/payload"
	"github.com/RobotsAndPencils/buford/push"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/micromdm/micromdm/config"
	"github.com/micromdm/micromdm/pubsub"
	"github.com/micromdm/micromdm/queue"
)

type Push struct {
	db       *DB
	start    chan struct{}
	provider PushCertificateProvider

	mu      sync.RWMutex
	pushsvc *push.Service
}

type PushCertificateProvider interface {
	PushCertificate() (*tls.Certificate, error)
}

type Option func(*Push)

func WithPushService(svc *push.Service) Option {
	return func(p *Push) {
		p.pushsvc = svc
	}
}

func New(db *DB, provider PushCertificateProvider, sub pubsub.Subscriber, opts ...Option) (*Push, error) {
	pushSvc := Push{
		db:       db,
		provider: provider,
		start:    make(chan struct{}),
	}
	for _, opt := range opts {
		opt(&pushSvc)
	}
	// if there is no push service, the push certificate hasn't been provided.
	// start a goroutine that delays the run of this service.
	if err := updateClient(&pushSvc, sub); err != nil {
		return nil, errors.Wrap(err, "wait for push service config")
	}

	if err := pushSvc.startQueuedSubscriber(sub); err != nil {
		return &pushSvc, err
	}
	return &pushSvc, nil
}

func (svc *Push) startQueuedSubscriber(sub pubsub.Subscriber) error {
	commandQueuedEvents, err := sub.Subscribe(context.TODO(), "push-info", queue.CommandQueuedTopic)
	if err != nil {
		return errors.Wrapf(err,
			"subscribing push to %s topic", queue.CommandQueuedTopic)
	}
	go func() {
		if svc.pushsvc == nil {
			log.Println("push: waiting for push certificate before enabling APNS service provider")
			<-svc.start
			log.Println("push: service started")
		}
		for {
			select {
			case event := <-commandQueuedEvents:
				cq, err := queue.UnmarshalQueuedCommand(event.Message)
				if err != nil {
					fmt.Println(err)
					continue
				}
				_, err = svc.Push(context.TODO(), cq.DeviceUDID)
				if err != nil {
					fmt.Println(err)
					continue
				}
			}
		}
	}()

	return nil
}

func updateClient(svc *Push, sub pubsub.Subscriber) error {
	configEvents, err := sub.Subscribe(context.TODO(), "push-server-configs", config.ConfigTopic)
	if err != nil {
		return errors.Wrap(err, "update push service client")
	}
	go func() {
		for {
			select {
			case <-configEvents:
				pushsvc, err := NewPushService(svc.provider)
				if err != nil {
					log.Println("push: could not get push certificate %s", err)
					continue
				}
				svc.mu.Lock()
				svc.pushsvc = pushsvc
				svc.mu.Unlock()
				go func() { svc.start <- struct{}{} }() // unblock queue
			}
		}
	}()
	return nil
}

func NewPushService(provider PushCertificateProvider) (*push.Service, error) {
	cert, err := provider.PushCertificate()
	if err != nil {
		return nil, errors.Wrap(err, "get push certificate from store")
	}

	client, err := push.NewClient(*cert)
	if err != nil {
		return nil, errors.Wrap(err, "create push service client")
	}

	svc := push.NewService(client, push.Production)
	return svc, nil
}

func (svc *Push) Push(ctx context.Context, deviceUDID string) (string, error) {
	info, err := svc.db.PushInfo(deviceUDID)
	if err != nil {
		return "", errors.Wrap(err, "retrieving PushInfo by UDID")
	}

	p := payload.MDM{Token: info.PushMagic}
	valid := push.IsDeviceTokenValid(info.Token)
	if !valid {
		return "", errors.New("invalid push token")
	}
	jsonPayload, err := json.Marshal(p)
	if err != nil {
		return "", errors.Wrap(err, "marshalling push notification payload")
	}
	result, err := svc.pushsvc.Push(info.Token, nil, jsonPayload)
	if err != nil && strings.HasSuffix(err.Error(), "remote error: tls: internal error") {
		// TODO: yuck, error substring searching. see:
		// https://github.com/micromdm/micromdm/issues/150
		return result, errors.Wrap(err, "push error: possibly expired or invalid APNs certificate")
	}
	return result, err
}
