package pubsub

import "context"

type (
	EventType string

	Event[T any] struct {
		Type    EventType
		Payload T
	}
)

type Subscriber[T any] interface {
	Subscribe(context.Context) <-chan Event[T]
}

const (
	CreatedEvent EventType = "created"
	UpdatedEvent EventType = "updated"
	DeletedEvent EventType = "deleted"
)
