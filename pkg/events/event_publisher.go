package events

import (
	"sync"
)

type EventPublisher struct {
	Producer Producer
	Events   map[string]Event
}

func NewEventPublisher(producer Producer) *EventPublisher {
	return &EventPublisher{
		Producer: producer,
		Events:   make(map[string]Event),
	}
}

func (k *EventPublisher) Register(event Event) EventPublisherInterface {
	k.Events[event.Name] = event
	return k
}

func (k *EventPublisher) Publish() {
	wg := &sync.WaitGroup{}
	for _, event := range k.Events {
		wg.Add(1)
		go func(event Event, wg *sync.WaitGroup) {
			err := k.Producer.Send(event, wg)
			if err == nil {
				delete(k.Events, event.Name)
			}
		}(event, wg)
	}
	wg.Wait()
}
