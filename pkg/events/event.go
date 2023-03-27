package events

import (
	"github.com/google/uuid"
	"sync"
	"time"
)

type Event struct {
	EID        uuid.UUID   `json:"eid"`
	OccurredAt time.Time   `json:"occurred_at"`
	Name       string      `json:"name"`
	Content    interface{} `json:"content"`
}

func NewEvent(Name string, Content interface{}) *Event {
	return &Event{
		EID:        uuid.New(),
		OccurredAt: time.Now(),
		Name:       Name,
		Content:    Content,
	}
}

type EventPublisherInterface interface {
	Register(event Event) EventPublisherInterface
	Publish()
}

type Producer interface {
	Send(event Event, wg *sync.WaitGroup) error
}
