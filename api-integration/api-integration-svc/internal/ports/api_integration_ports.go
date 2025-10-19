package ports

import (
	"context"

	"github.com/FrancoRebollo/api-integration-svc/internal/domain"
)

type ApiIntegrationService interface {
	CaptureEventAPI(ctx context.Context, reqCaptureEvent domain.Event) error
	ForwardRequest(req domain.ExternalAPIRequest) (domain.ExternalAPIResponse, error)
}

type ApiIntegrationRepository interface {
	CaptureEvent(ctx context.Context, reqCaptureEvent domain.Event) error
}
