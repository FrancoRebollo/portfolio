package application

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/FrancoRebollo/auth-security-svc/internal/adapters/in/http/dto"
	"github.com/FrancoRebollo/auth-security-svc/internal/domain"
	"github.com/FrancoRebollo/auth-security-svc/internal/platform/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildService(repo *MockSecurityRepository) SecurityService {
	mockCfg := &config.App{
		Name: "TEST", Environment: "TEST", Client: "TEST", Version: "0.0.1",
		FechaStartUp: time.Now().Format("02/01/2006 15:04:05"),
	}
	return *NewSecurityService(repo, *mockCfg)
}
func TestCreateUserAPI_Success(t *testing.T) {
	repo := new(MockSecurityRepository)
	service := buildService(repo)

	user := &domain.UserCreated{IdPersona: 1, CanalDigital: "USER_PASSWORD", LoginName: "frebollo"}
	repo.On("CreateUser", mock.Anything, *user).Return(user, nil)

	result, err := service.CreateUserAPI(context.Background(), user)

	assert.NoError(t, err)
	assert.Equal(t, user, result)
	repo.AssertExpectations(t)
}

func TestCreateUserAPI_Error(t *testing.T) {
	repo := new(MockSecurityRepository)
	service := buildService(repo)

	user := &domain.UserCreated{IdPersona: 1, CanalDigital: "USER_PASSWORD", LoginName: "frebollo"}
	repo.On("CreateUser", mock.Anything, *user).Return(nil, errors.New("db error"))

	result, err := service.CreateUserAPI(context.Background(), user)

	assert.Nil(t, result)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestCrearCanalDigitalAPI_Success(t *testing.T) {
	repo := new(MockSecurityRepository)
	service := buildService(repo)

	cd := domain.CanalDigital{CanalDigital: "APP"}
	repo.On("CrearCanalDigital", mock.Anything, cd, "apikey123").Return(nil)

	err := service.CrearCanalDigitalAPI(context.Background(), cd, "apikey123")

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestAccessPersonAPI_Success(t *testing.T) {
	repo := new(MockSecurityRepository)
	service := buildService(repo)
	input := domain.AccessPerson{IdPersona: 1, Revoke: "false"}
	repo.On("AccessPerson", mock.Anything, input, "key").Return(nil)

	err := service.AccessPersonAPI(context.Background(), input, "key")
	assert.NoError(t, err)
}

func TestAccessCanalDigitalAPI_Success(t *testing.T) {
	repo := new(MockSecurityRepository)
	service := buildService(repo)
	input := domain.AccessCanalDigital{CanalDigital: "APP", Revoke: "false"}
	repo.On("AccessCanalDigital", mock.Anything, input, "key").Return(nil)

	err := service.AccessCanalDigitalAPI(context.Background(), input, "key")
	assert.NoError(t, err)
}

func TestAccessApiKeyAPI_Success(t *testing.T) {
	repo := new(MockSecurityRepository)
	service := buildService(repo)
	input := domain.AccessApiKey{ApiKey: "key", FechaVigencia: "2025-12-01", Revoke: "false"}
	repo.On("AccessApiKey", mock.Anything, input, "key").Return(nil)

	err := service.AccessApiKeyAPI(context.Background(), input, "key")
	assert.NoError(t, err)
}

func TestAccessPersonMethodAuthAPI_Success(t *testing.T) {
	repo := new(MockSecurityRepository)
	service := buildService(repo)
	input := domain.AccessPersonMethodAuth{IdPersona: 1, MethodAuth: "PASS", Revoke: "false"}
	repo.On("AccessPersonMethodAuth", mock.Anything, input, "key").Return(nil)

	err := service.AccessPersonMethodAuthAPI(context.Background(), input, "key")
	assert.NoError(t, err)
}

func TestLoginAPI_Success(t *testing.T) {
	_ = os.Setenv("JWT_ACCESS_SEED", "1234567890123456")  // AES-128
	_ = os.Setenv("JWT_REFRESH_SEED", "6543210987654321") // AES-128
	_ = os.Setenv("REF_TOKEN_DURATION", "30")             // AES-128

	repo := new(MockSecurityRepository)
	service := buildService(repo)

	input := domain.Login{Username: "user", Password: "pass", ApiKey: "key", CanalDigital: "APP"}

	repo.On("LoginValidations", mock.Anything, input).Return(1, nil, nil)
	repo.On("GetAccessTokenDuration", mock.Anything, "key").Return(3600, nil)
	repo.On("UpsertAccessToken", mock.Anything, mock.AnythingOfType("*domain.UpsertAccessToken")).Return(nil)
	repo.On("PersistToken", mock.Anything, mock.AnythingOfType("domain.CredentialsToken")).Return(nil)

	res, err := service.LoginAPI(context.Background(), input)

	assert.NoError(t, err)
	assert.Equal(t, "user", res.Username)
}

func TestCheckApiKeyExpiradaAPI_Success(t *testing.T) {
	repo := new(MockSecurityRepository)
	service := buildService(repo)
	repo.On("CheckApiKeyExpirada", mock.Anything, "key").Return(true, nil)

	expired, err := service.CheckApiKeyExpiradaAPI(context.Background(), "key")
	assert.NoError(t, err)
	assert.True(t, expired)
}

func TestRecuperacionPasswordAPI_Success(t *testing.T) {
	repo := new(MockSecurityRepository)
	service := buildService(repo)
	repo.On("CambioPasswordByLogin", mock.Anything, "frebo", mock.Anything).Return(nil)

	dto := dto.ReqRecoveryPasswordDos{LoginName: "frebo"}
	err := service.RecuperacionPasswordAPI(context.Background(), dto)
	assert.NoError(t, err)
}
