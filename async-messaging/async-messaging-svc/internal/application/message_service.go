package application

import (
	"context"

	"github.com/FrancoRebollo/async-messaging-svc/internal/adapters/out/rabbitmq"
	"github.com/FrancoRebollo/async-messaging-svc/internal/domain"
	"github.com/FrancoRebollo/async-messaging-svc/internal/platform/config"

	"github.com/FrancoRebollo/async-messaging-svc/internal/ports"
)

type MessageService struct {
	hr   ports.MessageRepository
	rmq  *rabbitmq.RabbitMQAdapter
	conf config.App
}

func NewMessageService(hr ports.MessageRepository, rmq *rabbitmq.RabbitMQAdapter, conf config.App) *MessageService {
	return &MessageService{
		hr,
		rmq,
		conf,
	}
}

func (hs *MessageService) PushEventToQueueAPI(ctx context.Context, reqEvent domain.Event) error {
	if err := hs.hr.PushEventToQueue(ctx, reqEvent); err != nil {
		return err
	}

	return nil
}
