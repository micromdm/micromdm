package pubsub

import "context"

type Event struct {
	Topic   string
	Message []byte
}

type Publisher interface {
	Publish(ctx context.Context, topic string, msg []byte) error
}

type Subscriber interface {
	Subscribe(ctx context.Context, name, topic string) (<-chan Event, error)
}

type PublishSubscriber interface {
	Publisher
	Subscriber
}
