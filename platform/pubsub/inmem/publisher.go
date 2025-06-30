package inmem

import (
	"context"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/micromdm/micromdm/platform/pubsub"
)

func NewPubSub(logger log.Logger) *Inmem {
	publish := make(chan pubsub.Event)
	subscriptions := make(map[string][]subscription)
	inmem := &Inmem{
		publish:       publish,
		subscriptions: subscriptions,
		logger:        logger,
	}
	go inmem.dispatch()
	return inmem
}

type Inmem struct {
	mtx           sync.RWMutex
	subscriptions map[string][]subscription

	publish chan pubsub.Event
	logger  log.Logger
}

type subscription struct {
	name      string
	topic     string
	eventChan chan<- pubsub.Event
}

func (p *Inmem) Publish(_ context.Context, topic string, msg []byte) error {
	event := pubsub.Event{Topic: topic, Message: msg}
	go func() { p.publish <- event }()
	return nil
}
