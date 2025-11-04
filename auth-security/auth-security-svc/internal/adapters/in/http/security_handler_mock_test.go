package http

import (
	"context"

	"github.com/FrancoRebollo/auth-security-svc/internal/adapters/in/http/dto"
	"github.com/FrancoRebollo/auth-security-svc/internal/domain"
	"github.com/stretchr/testify/mock"
)

type mockSecurityService struct {
	mock.Mock
}

func (m *mockSecurityService) CreateUserAPI(ctx context.Context, req *domain.UserCreated) (*domain.UserCreated, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*domain.UserCreated), args.Error(1)
}

func (m *mockSecurityService) CrearCanalDigitalAPI(ctx context.Context, req domain.CanalDigital, apiKey string) error {
	args := m.Called(ctx, req, apiKey)
	return args.Error(0)
}

func (m *mockSecurityService) AccessPersonAPI(ctx context.Context, req domain.AccessPerson, apiKey string) error {
	args := m.Called(ctx, req, apiKey)
	return args.Error(0)
}

func (m *mockSecurityService) AccessCanalDigitalAPI(ctx context.Context, req domain.AccessCanalDigital, apiKey string) error {
	args := m.Called(ctx, req, apiKey)
	return args.Error(0)
}

func (m *mockSecurityService) AccessApiKeyAPI(ctx context.Context, req domain.AccessApiKey, apiKey string) error {
	args := m.Called(ctx, req, apiKey)
	return args.Error(0)
}
func (m *mockSecurityService) AccessPersonMethodAuthAPI(ctx context.Context, req domain.AccessPersonMethodAuth, apiKey string) error {
	args := m.Called(ctx, req, apiKey)
	return args.Error(0)
}
func (m *mockSecurityService) LoginAPI(context.Context, domain.Login) (domain.UserStatus, error) {
	return domain.UserStatus{}, nil
}
func (m *mockSecurityService) ValidateJWTAPI(context.Context, string) (*domain.CheckJWT, error) {
	return &domain.CheckJWT{}, nil
}
func (m *mockSecurityService) GetJWTAPI(context.Context, string, string) (string, error) {
	return "", nil
}
func (m *mockSecurityService) CheckApiKeyExpiradaAPI(context.Context, string) (bool, error) {
	return false, nil
}
func (m *mockSecurityService) RecuperacionPasswordAPI(context.Context, dto.ReqRecoveryPasswordDos) error {
	return nil
}
