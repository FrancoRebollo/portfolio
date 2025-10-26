package ports

import (
	"context"

	"github.com/FrancoRebollo/async-messaging-svc/internal/domain"
)

type MessageService interface {
	PushEventToQueueAPI(ctx context.Context, reqEvent domain.Event) error
}

type MessageRepository interface {
	PushEventToQueue(ctx context.Context, reqEvent domain.Event) error
	GetDatabasesPing(ctx context.Context) ([]domain.Database, error)
}

type MessageQueue interface {
	PushEventToQueue(ctx context.Context) (*domain.Event, error)
	PullEventFromQueue(ctx context.Context) ([]domain.Event, error)
}
