package kafka

import "context"

// Handler is message processing function
type Handler[T interface{}] func(ctx context.Context, msg *T) error

// Consumer accepts messages and processes them with some Handler
type Consumer[T interface{}] interface {
	Consume(ctx context.Context, topicName string, handler Handler[T]) (err error)
	Close() error
}

// Producer publishes messages
type Producer[T interface{}] interface {
	SendMessage(msg *T) error
	Close() error
}
