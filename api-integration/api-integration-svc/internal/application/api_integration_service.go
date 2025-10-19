package application

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/FrancoRebollo/api-integration-svc/internal/platform/config"

	"github.com/FrancoRebollo/api-integration-svc/internal/domain"
	"github.com/FrancoRebollo/api-integration-svc/internal/ports"
)

type ApiIntegrationService struct {
	hr         ports.ApiIntegrationRepository
	conf       config.App
	httpClient *http.Client
}

func NewApiIntegrationService(hr ports.ApiIntegrationRepository, conf config.App, httpClient *http.Client) *ApiIntegrationService {
	return &ApiIntegrationService{
		hr,
		conf,
		httpClient,
	}
}

func (hs *ApiIntegrationService) CaptureEventAPI(ctx context.Context, reqCaptureEvent domain.Event) error {

	if err := hs.hr.CaptureEvent(ctx, reqCaptureEvent); err != nil {
		return err
	}
	return nil
}

func (s *ApiIntegrationService) ForwardRequest(req domain.ExternalAPIRequest) (domain.ExternalAPIResponse, error) {
	fullURL := req.URL

	// Armar query params para GET
	if req.Method == "GET" && len(req.Params) > 0 {
		query := url.Values{}
		for k, v := range req.Params {
			query.Add(k, v)
		}
		fullURL = fmt.Sprintf("%s?%s", fullURL, query.Encode())
	}

	// Armar body si es POST
	var body io.Reader
	if req.Method == "POST" && req.Body != nil {
		jsonBody, _ := json.Marshal(req.Body)
		body = bytes.NewBuffer(jsonBody)
	}

	httpReq, err := http.NewRequest(req.Method, fullURL, body)
	if err != nil {
		return domain.ExternalAPIResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return domain.ExternalAPIResponse{}, err
	}
	defer resp.Body.Close()

	var result any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return domain.ExternalAPIResponse{}, err
	}

	return domain.ExternalAPIResponse{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Data:       result,
	}, nil
}
