package application

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/FrancoRebollo/async-messaging-svc/internal/domain"
	"github.com/FrancoRebollo/async-messaging-svc/internal/platform/config"

	"github.com/FrancoRebollo/async-messaging-svc/internal/ports"
)

type MessageService struct {
	hr   ports.MessageRepository
	rmq  ports.MessageQueue
	conf config.App
}

func NewMessageService(hr ports.MessageRepository, rmq ports.MessageQueue, conf config.App) *MessageService {
	return &MessageService{
		hr,
		rmq,
		conf,
	}
}

func (hs *MessageService) PushEventToQueueAPI(ctx context.Context, event domain.Event) error {
	fmt.Println("🧩 Iniciando transacción controlada para PushEventToQueueAPI...")

	err := hs.hr.WithTransaction(ctx, func(tx *sql.Tx) error {
		if err := hs.hr.PushEventToQueue(ctx, tx, event); err != nil {
			if errors.Is(err, domain.ErrDuplicateEvent) {
				fmt.Println("⚠️ Evento duplicado detectado, no se publicará en la cola")
				// 👉 devolvemos nil para que NO haya rollback
				return nil
			}
			return err
		}

		fmt.Println("✅ Evento persistido correctamente en DB")

		if err := hs.rmq.Publish(ctx, event); err != nil {
			fmt.Printf("❌ Error al publicar evento %s en cola: %v\n", event.ID, err)
			return err // rollback automático
		}

		fmt.Println("📨 Evento publicado en RabbitMQ correctamente")
		return nil
	})

	if err != nil {
		fmt.Println("🔻 Transacción revertida por error")
		return err
	}

	fmt.Println("✅ Transacción completada con éxito")
	return nil
}
