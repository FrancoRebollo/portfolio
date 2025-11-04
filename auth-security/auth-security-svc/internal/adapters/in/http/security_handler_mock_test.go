package http

import (
	"context"

	"github.com/FrancoRebollo/auth-security-svc/internal/adapters/in/http/dto"
	"github.com/FrancoRebollo/auth-security-svc/internal/domain"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/mock"
)

type mockSecurityService struct {
	mock.Mock
}

type mockTokenParser struct {
	mock.Mock
}

func (m *mockTokenParser) GetClaims(token string, tokenType string) (jwt.MapClaims, error) {
	args := m.Called(token, tokenType)
	return args.Get(0).(jwt.MapClaims), args.Error(1)
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
func (m *mockSecurityService) LoginAPI(ctx context.Context, req domain.Login) (domain.UserStatus, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(domain.UserStatus), args.Error(1)
}
func (m *mockSecurityService) ValidateJWTAPI(ctx context.Context, req string) (*domain.CheckJWT, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*domain.CheckJWT), args.Error(1)
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
