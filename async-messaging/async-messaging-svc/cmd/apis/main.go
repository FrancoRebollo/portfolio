// cmd/apis/main.go
package main

import (
	"fmt"
	"os"

	httpin "github.com/FrancoRebollo/async-messaging-svc/internal/adapters/in/http"
	pg "github.com/FrancoRebollo/async-messaging-svc/internal/adapters/out/postgres"
	"github.com/FrancoRebollo/async-messaging-svc/internal/adapters/out/rabbitmq"
	"github.com/FrancoRebollo/async-messaging-svc/internal/application"
	"github.com/FrancoRebollo/async-messaging-svc/internal/platform/config"
	"github.com/FrancoRebollo/async-messaging-svc/internal/platform/logger"
)

func main() {
	// 1) Configuración
	cfg, err := config.GetGlobalConfiguration()
	if err != nil {
		logger.LoggerError().Error(err)
		os.Exit(1)
	}

	// 2) Conexiones a bases de datos (según cfg.DB[*].Connection)
	var dbPostgres *pg.PostgresDB

	for _, conf := range cfg.DB {
		switch conf.Connection {
		case "POSTGRES":
			dbPostgres, err = pg.GetInstance(conf)
			if err != nil {
				logger.LoggerError().Errorf("Error conectando a Postgres: %s", err)
				os.Exit(1)
			}
		}
	}

	if dbPostgres != nil {
		logger.LoggerInfo().Info("Conexión a Postgres exitosa")
	}

	// 3) Repositorios (adapters out)
	versionRepository := pg.NewVersionRepository(*dbPostgres)
	healthcheckRepository := pg.NewHealthcheckRepository(dbPostgres)
	messageRepository := pg.NewMessageRepository(dbPostgres)

	// 4) RabbitMQ queues (adapters out)
	amqpURL := os.Getenv("RABBITMQ_URL")
	rabbitMQAdapter, err := rabbitmq.NewRabbitMQAdapter(amqpURL, "user.events")

	if err != nil {
		logger.LoggerError().Errorf("No se pudo iniciar cola de mensajeria: %s", err.Error())
		os.Exit(1)
	}

	// 5) Servicios (application)
	versionService := application.NewVersionService(versionRepository, *cfg.App)
	healthcheckService := application.NewHealthcheckService(healthcheckRepository, *cfg.App)
	messageService := application.NewMessageService(messageRepository, rabbitMQAdapter, *cfg.App)

	// 6) Handlers (adapters in/http)
	versionHandler := httpin.NewVersionHandler(versionService)
	healthcheckHandler := httpin.NewHealthcheckHandler(healthcheckService)
	messageHandler := httpin.NewMessageHandler(messageService)

	// 7) Router
	rt, err := httpin.NewRouter(cfg.HTTP, versionHandler, *healthcheckHandler, *messageHandler)
	if err != nil {
		fmt.Println(err)
	}

	// 8) Server
	address := fmt.Sprintf("%s:%s", cfg.HTTP.Url, cfg.HTTP.Port)
	if err := rt.Listen(address); err != nil {
		logger.LoggerError().Errorf("No se pudo iniciar el servidor: %s", err.Error())
		os.Exit(1)
	}
}
