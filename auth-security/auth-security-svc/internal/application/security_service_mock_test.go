package application

import (
	"context"

	"github.com/FrancoRebollo/auth-security-svc/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockSecurityRepository struct {
	mock.Mock
}

func (m *MockSecurityRepository) CreateUser(ctx context.Context, reqAltaUser domain.UserCreated) (*domain.UserCreated, error) {
	args := m.Called(ctx, reqAltaUser)

	var result *domain.UserCreated
	if args.Get(0) != nil {
		result = args.Get(0).(*domain.UserCreated)
	}

	return result, args.Error(1)
}

func (m *MockSecurityRepository) CrearCanalDigital(ctx context.Context, crearCanalDigital domain.CanalDigital, apiKey string) error {
	args := m.Called(ctx, crearCanalDigital, apiKey)
	return args.Error(0)
}

func (m *MockSecurityRepository) AccessPerson(ctx context.Context, accessPerson domain.AccessPerson, apikey string) error {
	args := m.Called(ctx, accessPerson, apikey)
	return args.Error(0)
}

func (m *MockSecurityRepository) AccessCanalDigital(ctx context.Context, accessCanaldigital domain.AccessCanalDigital, apikey string) error {
	args := m.Called(ctx, accessCanaldigital, apikey)
	return args.Error(0)
}

func (m *MockSecurityRepository) AccessApiKey(ctx context.Context, accessApiKey domain.AccessApiKey, apikey string) error {
	args := m.Called(ctx, accessApiKey, apikey)
	return args.Error(0)
}

func (m *MockSecurityRepository) AccessPersonMethodAuth(ctx context.Context, accesPersonMethodAuth domain.AccessPersonMethodAuth, apikey string) error {
	args := m.Called(ctx, accesPersonMethodAuth, apikey)
	return args.Error(0)
}

func (m *MockSecurityRepository) LoginValidations(ctx context.Context, reqLogin domain.Login) (int, *string, error) {
	args := m.Called(ctx, reqLogin)
	return args.Int(0), nil, args.Error(2)
}

func (m *MockSecurityRepository) GetAccessTokenDuration(ctx context.Context, apiKey string) (int, error) {
	args := m.Called(ctx, apiKey)
	return args.Int(0), args.Error(1)
}

func (m *MockSecurityRepository) UpsertAccessToken(ctx context.Context, requestUpsert *domain.UpsertAccessToken) error {
	args := m.Called(ctx, requestUpsert)
	return args.Error(0)
}

func (m *MockSecurityRepository) CheckLastAccessToken(ctx context.Context, token string, credentials domain.Credentials) error {
	args := m.Called(ctx, token, credentials)
	return args.Error(0)
}

func (m *MockSecurityRepository) CheckLastRefreshToken(ctx context.Context, token string, credentials domain.Credentials) error {
	args := m.Called(ctx, token, credentials)
	return args.Error(0)
}

func (m *MockSecurityRepository) CheckTokenCreation(ctx context.Context, credentials domain.Credentials) error {
	args := m.Called(ctx, credentials)
	return args.Error(0)
}

func (m *MockSecurityRepository) PersistToken(ctx context.Context, credentials domain.CredentialsToken) error {
	args := m.Called(ctx, credentials)
	return args.Error(0)
}

func (m *MockSecurityRepository) CheckApiKeyExpirada(ctx context.Context, apiKey string) (bool, error) {
	args := m.Called(ctx, apiKey)
	return args.Bool(0), args.Error(1)
}

func (m *MockSecurityRepository) CambioPasswordByLogin(ctx context.Context, loginName string, newPassword string) error {
	args := m.Called(ctx, loginName, newPassword)
	return args.Error(0)
}
