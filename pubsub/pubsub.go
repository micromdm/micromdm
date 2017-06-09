package pubsub

type Event struct {
	Topic   string
	Message []byte
}

type Publisher interface {
	Publish(topic string, msg []byte) error
}

type Subscriber interface {
	Subscribe(name, topic string) (<-chan Event, error)
}

type PublishSubscriber interface {
	Publisher
	Subscriber
}
