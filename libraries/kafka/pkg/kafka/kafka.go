package kafka

import "context"

// Handler is message processing function
type Handler func(ctx context.Context, msg interface{}) error

// Consumer accepts messages and processes them with some Handler
type Consumer interface {
	Consume(ctx context.Context, topicName string, handler Handler) (err error)
	Close() error
}

// Producer publishes messages
type Producer interface {
	SendMessage(msg interface{}) error
	Close() error
}
