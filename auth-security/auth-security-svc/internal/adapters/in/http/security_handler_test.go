package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FrancoRebollo/auth-security-svc/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockSecurityService)
	handler := NewSecurityHandler(mockService)

	// Request simulada
	requestBody := `{
		"id_persona": 123,
		"canal_digital": "APP",
		"login_name": "frebollo",
		"password": "secure123",
		"mail_persona": "frebo@example.com",
		"tel_persona": "1155555555"
	}`

	req := httptest.NewRequest(http.MethodPost, "/create-user", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Response recorder
	w := httptest.NewRecorder()

	// Crear contexto de Gin
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	// Setear lo que el mock debe devolver
	expectedUser := &domain.UserCreated{
		IdPersona:    123,
		CanalDigital: "APP",
		LoginName:    "frebollo",
		Password:     "secure123",
		MailPersona:  "frebo@example.com",
		TePersona:    "1155555555",
	}

	mockService.
		On("CreateUserAPI", mock.Anything, mock.AnythingOfType("*domain.UserCreated")).
		Return(expectedUser, nil)

	// Ejecutar handler
	handler.CreateUser(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"login_name":"frebollo"`)
	mockService.AssertExpectations(t)
}

func TestCreateUser_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockSecurityService)
	handler := NewSecurityHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/create-user", strings.NewReader(`invalid json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	handler.CreateUser(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error"`)
}

func TestCreateCanalDigital_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockSecurityService)
	handler := NewSecurityHandler(mockService)

	reqBody := `{"canal_digital":"USER_PASSWORD"}`
	req := httptest.NewRequest(http.MethodPost, "/create-method-auth", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", "example-api-key")
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	// Expectation del mock
	mockService.
		On(
			"CrearCanalDigitalAPI",
			mock.Anything,
			mock.AnythingOfType("domain.CanalDigital"),
			"example-api-key",
		).
		Return(nil) // ✅ caso feliz: NO error

	// Act
	handler.CreateCanalDigital(ctx)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"Canal digital creado"`)

	mockService.AssertExpectations(t)
}

func TestCreateCanalDigital_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockSecurityService)
	handler := NewSecurityHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/create-method-auth", strings.NewReader(`invalid json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	handler.CreateCanalDigital(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error"`)
}

func TestAccessPerson_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockSecurityService)
	handler := NewSecurityHandler(mockService)

	reqBody := `{
		"id_persona": 123,
		"revoke": "S"
	}`

	req := httptest.NewRequest(http.MethodPost, "/unaccess-person", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", "example-api-key")
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	// Expectation del mock
	mockService.
		On(
			"AccessPersonAPI",
			mock.Anything,
			mock.AnythingOfType("domain.AccessPerson"),
			"example-api-key",
		).
		Return(nil) // ✅ caso feliz: NO error

	// Act
	handler.AccessPerson(ctx)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Condition(t, func() bool {
		body := w.Body.String()
		return strings.Contains(body, "Revoke access to person") ||
			strings.Contains(body, "Revoke unaccess to person")
	}, "Response body should contain one of the expected messages")

	mockService.AssertExpectations(t)
}

func TestAccessPerson_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockSecurityService)
	handler := NewSecurityHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/unaccess-person", strings.NewReader(`invalid json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	handler.AccessPerson(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error"`)
}

func TestAccessCanalDigital_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockSecurityService)
	handler := NewSecurityHandler(mockService)

	reqBody := `{
		"id_persona": 123,
		"revoke": "S"
	}`

	req := httptest.NewRequest(http.MethodPost, "/unaccess-digital-channel", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", "example-api-key")
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	// Expectation del mock
	mockService.
		On(
			"AccessCanalDigitalAPI",
			mock.Anything,
			mock.AnythingOfType("domain.AccessCanalDigital"),
			"example-api-key",
		).
		Return(nil) // ✅ caso feliz: NO error

	// Act
	handler.AccessCanalDigital(ctx)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Condition(t, func() bool {
		body := w.Body.String()
		return strings.Contains(body, "Revoke access to digital channel") ||
			strings.Contains(body, "Revoke unaccess to digital channel")
	}, "Response body should contain one of the expected messages")

	mockService.AssertExpectations(t)
}

func TestAccessCanalDigital_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockSecurityService)
	handler := NewSecurityHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/unaccess-digital-channel", strings.NewReader(`invalid json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	handler.AccessCanalDigital(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error"`)
}

func TestAccessPerMethodAuth_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockSecurityService)
	handler := NewSecurityHandler(mockService)

	reqBody := `{
		"id_persona": 123,
		"method_auth": "USER_PASSWORD",
		"revoke": "S"
	}`

	req := httptest.NewRequest(http.MethodPost, "/unaccess-digital-channel-person", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", "example-api-key")
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	// Expectation del mock
	mockService.
		On(
			"AccessPersonMethodAuthAPI",
			mock.Anything,
			mock.AnythingOfType("domain.AccessPersonMethodAuth"),
			"example-api-key",
		).
		Return(nil) // ✅ caso feliz: NO error

	// Act
	handler.AccessPerMethodAuth(ctx)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Condition(t, func() bool {
		body := w.Body.String()
		return strings.Contains(body, "Revoke access to person by digital channel") ||
			strings.Contains(body, "Revoke unaccess to person by digital channel")
	}, "Response body should contain one of the expected messages")

	mockService.AssertExpectations(t)
}

func TestAccessPerMethodAuth_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockSecurityService)
	handler := NewSecurityHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/unaccess-digital-channel-person", strings.NewReader(`invalid json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	handler.AccessPerMethodAuth(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error"`)
}

func TestAcessApiKey_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockSecurityService)
	handler := NewSecurityHandler(mockService)

	reqBody := `{
		"api_key": "asdf",
		"revoke": "S"
	}`

	req := httptest.NewRequest(http.MethodPost, "/unaccess-api-key", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", "example-api-key")
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	// Expectation del mock
	mockService.
		On(
			"AccessApiKeyAPI",
			mock.Anything,
			mock.AnythingOfType("domain.AccessApiKey"),
			"example-api-key",
		).
		Return(nil) // ✅ caso feliz: NO error

	// Act
	handler.AccessApiKey(ctx)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Condition(t, func() bool {
		body := w.Body.String()
		return strings.Contains(body, "Revoke access to api key") ||
			strings.Contains(body, "Revoke unaccess to api key")
	}, "Response body should contain one of the expected messages")

	mockService.AssertExpectations(t)
}

func TestAcessApiKey_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockSecurityService)
	handler := NewSecurityHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/unaccess-api-key", strings.NewReader(`invalid json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	handler.AccessApiKey(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error"`)
}

/*
func Test<HandlerMethod>_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockService := new(mockSecurityService)
	handler := NewSecurityHandler(mockService)

	reqBody := `<json aquí>`
	req := httptest.NewRequest(http.MethodPost, "/endpoint", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	expected := ... // lo que debe devolver el servicio
	mockService.On("Metodo", mock.Anything, ...).Return(expected, nil)

	// Act
	handler.<HandlerMethod>(ctx)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `<campo esperado>`)
	mockService.AssertExpectations(t)
}

*/
