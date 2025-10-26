package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/FrancoRebollo/async-messaging-svc/internal/domain"

	"github.com/streadway/amqp"
)

type RabbitMQAdapter struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

// ✅ Constructor del adaptador (será inyectado al servicio)
func NewRabbitMQAdapter(amqpURL, queueName string) (*RabbitMQAdapter, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open RabbitMQ channel: %w", err)
	}

	// Declaración de la cola (si no existe, la crea)
	q, err := ch.QueueDeclare(
		queueName, // nombre
		true,      // durable: persiste tras reinicio
		false,     // auto-delete
		false,     // exclusive
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	fmt.Printf("✅ Connected to RabbitMQ and queue '%s' ready.\n", q.Name)

	return &RabbitMQAdapter{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

// ✅ PushEventToQueue: publica un mensaje en la cola
func (r *RabbitMQAdapter) PushEventToQueue(ctx context.Context) (*domain.Event, error) {
	event, ok := ctx.Value("event").(*domain.Event)
	if !ok || event == nil {
		return nil, fmt.Errorf("no event found in context")
	}

	body, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event: %w", err)
	}

	err = r.channel.Publish(
		"",           // exchange (directo por defecto)
		r.queue.Name, // routing key (nombre de la cola)
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Timestamp:   time.Now(),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to publish message: %w", err)
	}

	fmt.Printf("📨 Event %s published to queue '%s'\n", event.EventId, r.queue.Name)
	return event, nil
}

// ✅ PullEventToQueue: consume mensajes desde la cola
func (r *RabbitMQAdapter) PullEventToQueue(ctx context.Context) ([]domain.Event, error) {
	msgs, err := r.channel.Consume(
		r.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to consume messages: %w", err)
	}

	var events []domain.Event

	// Canal concurrente para procesar mensajes
	for {
		select {
		case d := <-msgs:
			var e domain.Event
			if err := json.Unmarshal(d.Body, &e); err != nil {
				fmt.Printf("⚠️ Error unmarshalling event: %v\n", err)
				continue
			}
			fmt.Printf("📥 Received event: %s from %s\n", e.EventId, e.EventOrigin)
			events = append(events, e)

		case <-ctx.Done():
			fmt.Println("🛑 Consumer stopped by context")
			return events, nil
		}
	}
}

// ✅ Cierre de conexión ordenado
func (r *RabbitMQAdapter) Close() {
	if r.channel != nil {
		_ = r.channel.Close()
	}
	if r.conn != nil {
		_ = r.conn.Close()
	}
	fmt.Println("🧹 RabbitMQ connection closed.")
}
