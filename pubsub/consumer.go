package pubsub

type Subscriber interface {
	Subscribe(name, topic string) (<-chan Event, error)
}

func (p *Inmem) Subscribe(name, topic string) (<-chan Event, error) {
	events := make(chan Event)
	sub := subscription{
		name:      name,
		topic:     topic,
		eventChan: events,
	}
	p.mtx.Lock()
	p.subscriptions[topic] = append(p.subscriptions[topic], sub)
	p.mtx.Unlock()

	return events, nil
}

func (p *Inmem) dispatch() {
	for {
		select {
		case ev := <-p.publish:
			p.mtx.Lock()
			for _, sub := range p.subscriptions[ev.Topic] {
				go func(s subscription) { s.eventChan <- ev }(sub)
			}
			p.mtx.Unlock()
		}
	}
}
