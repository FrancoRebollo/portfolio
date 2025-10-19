package http

import (
	"net/http"

	"github.com/FrancoRebollo/api-integration-svc/internal/adapters/in/http/dto"
	"github.com/FrancoRebollo/api-integration-svc/internal/domain"
	"github.com/FrancoRebollo/api-integration-svc/internal/platform/logger"
	"github.com/FrancoRebollo/api-integration-svc/internal/ports"

	"github.com/gin-gonic/gin"
)

type ApiIntegrationHandler struct {
	serv ports.ApiIntegrationService
}

func NewApiIntegrationHandler(serv ports.ApiIntegrationService) *ApiIntegrationHandler {
	return &ApiIntegrationHandler{
		serv,
	}
}

func (hh *ApiIntegrationHandler) CaptureEvent(c *gin.Context) {
	ctx := c.Request.Context()

	var reqCaptureEvent dto.ReqCaptureEvent

	if err := c.BindJSON(&reqCaptureEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domCaptureEvent := &domain.Event{
		EventType:    reqCaptureEvent.EventType,
		EventContent: reqCaptureEvent.EventContent,
	}

	err := hh.serv.CaptureEventAPI(ctx, *domCaptureEvent)
	if err != nil {
		logger.LoggerError().Error(err)
		errorResponse(c, err)
		return
	}

	resp := &dto.DefaultResponse{
		Message: "Event captured",
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ApiIntegrationHandler) MakeRequest(c *gin.Context) {
	var req dto.ExternalAPIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	domainReq := domain.ExternalAPIRequest{
		Method: req.Method,
		URL:    req.URL,
		Params: req.Params,
		Body:   req.Body,
	}

	resp, err := h.serv.ForwardRequest(domainReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to call external API"})
		return
	}

	c.JSON(http.StatusOK, dto.ExternalAPIResponse{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Data:       resp.Data,
	})
}
