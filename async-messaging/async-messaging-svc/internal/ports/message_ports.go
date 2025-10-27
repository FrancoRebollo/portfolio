package ports

import (
	"context"
	"database/sql"

	"github.com/FrancoRebollo/async-messaging-svc/internal/domain"
)

type MessageService interface {
	PushEventToQueueAPI(ctx context.Context, reqEvent domain.Event) error
}

type MessageRepository interface {
	PushEventToQueue(ctx context.Context, tx *sql.Tx, event domain.Event) error
	GetDatabasesPing(ctx context.Context) ([]domain.Database, error)
	WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error
}

type MessageQueue interface {
	Publish(ctx context.Context, event domain.Event) error
	Consume(ctx context.Context, queue string, handler func(domain.Event)) error
}
