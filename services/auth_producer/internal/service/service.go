package service

import "context"

// UserCreatorService produces new users
type UserCreatorService interface {
	Create(ctx context.Context) error
	ProduceService
}

// ProduceService is some producer of events
type ProduceService interface {
	RunProducer(ctx context.Context) error
}
