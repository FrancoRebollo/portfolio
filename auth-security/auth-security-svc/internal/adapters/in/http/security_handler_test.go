package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FrancoRebollo/auth-security-svc/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ------------------------ CREATE USER ------------------------

func TestCreateUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockSecurityService)
	mockParser := new(mockTokenParser)
	handler := NewSecurityHandler(mockService, mockParser)

	body := `{
		"id_persona": 123,
		"canal_digital": "APP",
		"login_name": "frebollo",
		"password": "secure123",
		"mail_persona": "frebo@example.com",
		"tel_persona": "1155555555"
	}`

	req := httptest.NewRequest(http.MethodPost, "/create-user", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	expected := &domain.UserCreated{
		IdPersona:    123,
		CanalDigital: "APP",
		LoginName:    "frebollo",
		Password:     "secure123",
		MailPersona:  "frebo@example.com",
		TePersona:    "1155555555",
	}

	mockService.On("CreateUserAPI", mock.Anything, mock.AnythingOfType("*domain.UserCreated")).Return(expected, nil)

	handler.CreateUser(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"login_name":"frebollo"`)
	mockService.AssertExpectations(t)
}

func TestCreateUser_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockSecurityService)
	mockParser := new(mockTokenParser)
	handler := NewSecurityHandler(mockService, mockParser)

	req := httptest.NewRequest(http.MethodPost, "/create-user", strings.NewReader(`invalid`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	handler.CreateUser(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ------------------------ CREATE CANAL DIGITAL ------------------------

func TestCreateCanalDigital_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockSecurityService)
	mockParser := new(mockTokenParser)
	handler := NewSecurityHandler(mockService, mockParser)

	body := `{"canal_digital":"APP"}`
	req := httptest.NewRequest(http.MethodPost, "/create-method-auth", strings.NewReader(body))
	req.Header.Set("Api-Key", "key")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	mockService.On("CrearCanalDigitalAPI", mock.Anything, mock.AnythingOfType("domain.CanalDigital"), "key").Return(nil)

	handler.CreateCanalDigital(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"Canal digital creado"`)
	mockService.AssertExpectations(t)
}

func TestCreateCanalDigital_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockSecurityService)
	mockParser := new(mockTokenParser)
	handler := NewSecurityHandler(mockService, mockParser)

	req := httptest.NewRequest(http.MethodPost, "/create-method-auth", strings.NewReader(`bad`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	handler.CreateCanalDigital(ctx)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ------------------------ ACCESS PERSON ------------------------

func TestAccessPerson_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockSecurityService)
	mockParser := new(mockTokenParser)
	handler := NewSecurityHandler(mockService, mockParser)

	body := `{"id_persona":123,"revoke":"S"}`
	req := httptest.NewRequest(http.MethodPost, "/unaccess-person", strings.NewReader(body))
	req.Header.Set("Api-Key", "key")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	mockService.On("AccessPersonAPI", mock.Anything, mock.AnythingOfType("domain.AccessPerson"), "key").Return(nil)

	handler.AccessPerson(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// ------------------------ LOGIN ------------------------

func TestLogin_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockSecurityService)
	mockParser := new(mockTokenParser)
	handler := NewSecurityHandler(mockService, mockParser)

	body := `{"username":"frebo","password":"123","canal_digital":"APP"}`
	req := httptest.NewRequest(http.MethodPost, "/sec/log-in", strings.NewReader(body))
	req.Header.Set("Api-Key", "key")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	mockService.
		On("LoginAPI", mock.Anything, mock.Anything).
		Return(domain.UserStatus{
			Username: "frebo",
			Status:   "OK",
		}, nil)

	handler.Login(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"frebo"`)
	mockService.AssertExpectations(t)
}

// ------------------------ VALIDATE JWT ------------------------

func TestValidateJWT_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockSecurityService)
	mockParser := new(mockTokenParser)
	handler := NewSecurityHandler(mockService, mockParser)

	req := httptest.NewRequest(http.MethodGet, "/sec/validate-jwt", nil)
	req.Header.Set("Authorization", "Bearer token123")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	mockParser.On("GetClaims", "token123", "ACCESS").
		Return(jwt.MapClaims{"id_persona": 77}, nil)

	mockParser.On("GetClaims", "badtoken", "ACCESS").
		Return(nil, errors.New("invalid token"))

	mockService.On("ValidateJWTAPI", mock.Anything, "token123").
		Return(&domain.CheckJWT{TokenStatus: "VALID"}, nil)

	handler.ValidateJWT(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "77") // número, no string
	mockService.AssertExpectations(t)
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockSecurityService)
	mockParser := new(mockTokenParser)
	handler := NewSecurityHandler(mockService, mockParser)

	req := httptest.NewRequest(http.MethodGet, "/sec/validate-jwt", nil)
	req.Header.Set("Authorization", "Bearer badtoken")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	mockParser.
		On("GetClaims", "badtoken", "ACCESS").
		Return(jwt.MapClaims{}, errors.New("invalid token"))

	handler.ValidateJWT(ctx)

	assert.Equal(t, http.StatusUnauthorized, w.Code) // o BadRequest si tu lógica lo indica
	assert.Contains(t, w.Body.String(), "invalid token")

	mockParser.AssertExpectations(t)
}

func TestAccessApiKey_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockSecurityService)
	mockParser := new(mockTokenParser)
	handler := NewSecurityHandler(mockService, mockParser)

	body := `{"api_key":"test-api-key","revoke":"S"}`
	req := httptest.NewRequest(http.MethodPost, "/unaccess-api-key", strings.NewReader(body))
	req.Header.Set("Api-Key", "caller-api-key")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	expected := domain.AccessApiKey{
		ApiKey: "test-api-key",
		Revoke: "S",
	}

	mockService.On("AccessApiKeyAPI", mock.Anything, expected, "caller-api-key").Return(nil)

	handler.AccessApiKey(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Revoke access to api key")
	mockService.AssertExpectations(t)
}

func TestAccessApiKey_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockSecurityService)
	mockParser := new(mockTokenParser)
	handler := NewSecurityHandler(mockService, mockParser)

	req := httptest.NewRequest(http.MethodPost, "/unaccess-api-key", strings.NewReader(`bad`))
	req.Header.Set("Api-Key", "key")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	handler.AccessApiKey(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAccessPersonMethodAuth_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockSecurityService)
	mockParser := new(mockTokenParser)
	handler := NewSecurityHandler(mockService, mockParser)

	body := `{"id_persona":1,"method_auth":"PASS","revoke":"S"}`
	req := httptest.NewRequest(http.MethodPost, "/unaccess-digital-channel-person", strings.NewReader(body))
	req.Header.Set("Api-Key", "key")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	expected := domain.AccessPersonMethodAuth{
		IdPersona:  1,
		MethodAuth: "PASS",
		Revoke:     "S",
	}

	mockService.On("AccessPersonMethodAuthAPI", mock.Anything, expected, "key").Return(nil)

	handler.AccessPerMethodAuth(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Revoke access to person by digital channel")
	mockService.AssertExpectations(t)
}

func TestAccessPersonMethodAuth_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockSecurityService)
	mockParser := new(mockTokenParser)
	handler := NewSecurityHandler(mockService, mockParser)

	req := httptest.NewRequest(http.MethodPost, "/unaccess-digital-channel-person", strings.NewReader(`invalid`))
	req.Header.Set("Api-Key", "key")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	handler.AccessPerMethodAuth(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAccessCanalDigital_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockSecurityService)
	mockParser := new(mockTokenParser)
	handler := NewSecurityHandler(mockService, mockParser)

	body := `{"canal_digital":"APP","revoke":"N"}`
	req := httptest.NewRequest(http.MethodPost, "/unaccess-digital-channel", strings.NewReader(body))
	req.Header.Set("Api-Key", "key")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	expected := domain.AccessCanalDigital{
		CanalDigital: "APP",
		Revoke:       "N",
	}

	mockService.On("AccessCanalDigitalAPI", mock.Anything, expected, "key").Return(nil)

	handler.AccessCanalDigital(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Revoke unaccess to digital channel")
	mockService.AssertExpectations(t)
}

func TestAccessCanalDigital_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockSecurityService)
	mockParser := new(mockTokenParser)
	handler := NewSecurityHandler(mockService, mockParser)

	req := httptest.NewRequest(http.MethodPost, "/unaccess-digital-channel", strings.NewReader(`invalid`))
	req.Header.Set("Api-Key", "key")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	handler.AccessCanalDigital(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
