package http

import (
	"net/http"
	"time"

	"github.com/FrancoRebollo/async-messaging-svc/internal/adapters/in/http/dto"
	"github.com/FrancoRebollo/async-messaging-svc/internal/domain"
	"github.com/FrancoRebollo/async-messaging-svc/internal/platform/logger"
	"github.com/FrancoRebollo/async-messaging-svc/internal/ports"
	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	serv ports.MessageService
}

func NewMessageHandler(serv ports.MessageService) *MessageHandler {
	return &MessageHandler{
		serv,
	}
}

// GetHealthcheck verifica el estado del servicio
// @Summary Verifica el estado del servicio
// @Description Devuelve un JSON indicando si el servicio está activo.
// @Tags healthcheck
// @Accept json
// @Produce json
// @Success 200 {object} domain.Healthcheck "Estado del servicio"
// @Failure 400 {object} domain.HealthcheckError "Bad Request"
// @Failure 401 {object} domain.HealthcheckError "Unauthorized"
// @Failure 404 {object} domain.HealthcheckError "Not found"
// @Failure 409 {object} domain.HealthcheckError "Conflict"
// @Failure 500 {object} domain.HealthcheckError "Internal Server Error"
// @Failure 503 {object} domain.HealthcheckError "Service Unavailable"
// @Failure 504 {object} domain.HealthcheckError "Timeout"
// @Security BearerAuth
// @Router /api/healthcheck/ [get]
func (hh *MessageHandler) PushEventToQueue(c *gin.Context) {
	ctx := c.Request.Context()

	var reqPushEvent dto.RequestPushEvent
	if err := c.BindJSON(&reqPushEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domainEvent := domain.Event{
		EventId:      reqPushEvent.EventId,
		EventOrigin:  reqPushEvent.EventOrigin,
		EventDestiny: reqPushEvent.EventDestiny,
		EventType:    reqPushEvent.EventType,
		Payload:      reqPushEvent.Payload,
		Status:       "PENDING",  // por defecto
		CreatedAt:    time.Now(), // 🕒 asigna la fecha/hora actual
	}

	err := hh.serv.PushEventToQueueAPI(ctx, domainEvent)
	if err != nil {
		logger.LoggerError().Error(err)
		errorResponse(c, err)
		return
	}

	responseDefault := dto.ResponseDefault{
		Message: "Mensajo encolado exitosamente",
	}

	c.JSON(http.StatusOK, responseDefault)
}
