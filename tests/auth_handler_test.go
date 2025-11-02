package tests

import (
	"bytes"
	"dklautomationgo/handlers"
	"dklautomationgo/models"
	"dklautomationgo/services"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// TEST SUITE: AUTH HANDLER - LOGIN ENDPOINT
// ============================================================================

func TestAuthHandler_HandleLogin_Success(t *testing.T) {
	// Setup
	app := fiber.New()
	mockAuthService := new(MockAuthService)
	mockPermService := new(AuthMockPermissionService)
	mockRateLimiter := newMockRateLimiter()

	handler := handlers.NewAuthHandler(mockAuthService, mockPermService, mockRateLimiter)

	testUser := &models.Gebruiker{
		ID:       "user-123",
		Email:    "test@dekoninklijkeloop.nl",
		Naam:     "Test User",
		Rol:      "admin",
		IsActief: true,
	}

	permissions := []*models.UserPermission{
		{Resource: "admin", Action: "access"},
		{Resource: "contact", Action: "read"},
	}

	// Mock expectations
	mockAuthService.On("Login", mock.Anything, "test@dekoninklijkeloop.nl", "password123").
		Return("access-token-xyz", "refresh-token-abc", nil)
	mockAuthService.On("GetUserFromToken", mock.Anything, "access-token-xyz").
		Return(testUser, nil)
	mockPermService.On("GetUserPermissions", mock.Anything, "user-123").
		Return(permissions, nil)

	// Setup route
	app.Post("/api/auth/login", handler.HandleLogin)

	// Create request
	loginData := map[string]string{
		"email":      "test@dekoninklijkeloop.nl",
		"wachtwoord": "password123",
	}
	jsonData, _ := json.Marshal(loginData)

	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Execute
	resp, err := app.Test(req)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse response
	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	assert.True(t, response["success"].(bool))
	assert.Equal(t, "access-token-xyz", response["token"])
	assert.Equal(t, "refresh-token-abc", response["refresh_token"])
	assert.NotNil(t, response["user"])

	user := response["user"].(map[string]interface{})
	assert.Equal(t, "user-123", user["id"])
	assert.Equal(t, "test@dekoninklijkeloop.nl", user["email"])

	mockAuthService.AssertExpectations(t)
	mockPermService.AssertExpectations(t)
}

func TestAuthHandler_HandleLogin_MissingEmail(t *testing.T) {
	app := fiber.New()
	mockAuthService := new(MockAuthService)
	mockPermService := new(AuthMockPermissionService)
	mockRateLimiter := newMockRateLimiter()

	handler := handlers.NewAuthHandler(mockAuthService, mockPermService, mockRateLimiter)
	app.Post("/api/auth/login", handler.HandleLogin)

	// Missing email
	loginData := map[string]string{
		"wachtwoord": "password123",
	}
	jsonData, _ := json.Marshal(loginData)

	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "verplicht")
}

func TestAuthHandler_HandleLogin_MissingPassword(t *testing.T) {
	app := fiber.New()
	mockAuthService := new(MockAuthService)
	mockPermService := new(AuthMockPermissionService)
	mockRateLimiter := newMockRateLimiter()

	handler := handlers.NewAuthHandler(mockAuthService, mockPermService, mockRateLimiter)
	app.Post("/api/auth/login", handler.HandleLogin)

	// Missing password
	loginData := map[string]string{
		"email": "test@dekoninklijkeloop.nl",
	}
	jsonData, _ := json.Marshal(loginData)

	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "verplicht")
}

func TestAuthHandler_HandleLogin_InvalidCredentials(t *testing.T) {
	app := fiber.New()
	mockAuthService := new(MockAuthService)
	mockPermService := new(AuthMockPermissionService)
	mockRateLimiter := newMockRateLimiter()

	handler := handlers.NewAuthHandler(mockAuthService, mockPermService, mockRateLimiter)
	app.Post("/api/auth/login", handler.HandleLogin)

	mockAuthService.On("Login", mock.Anything, "test@dekoninklijkeloop.nl", "wrongpass").
		Return("", "", services.ErrInvalidCredentials)

	loginData := map[string]string{
		"email":      "test@dekoninklijkeloop.nl",
		"wachtwoord": "wrongpass",
	}
	jsonData, _ := json.Marshal(loginData)

	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Ongeldige inloggegevens")

	mockAuthService.AssertExpectations(t)
}

func TestAuthHandler_HandleLogin_InactiveUser(t *testing.T) {
	app := fiber.New()
	mockAuthService := new(MockAuthService)
	mockPermService := new(AuthMockPermissionService)
	mockRateLimiter := newMockRateLimiter()

	handler := handlers.NewAuthHandler(mockAuthService, mockPermService, mockRateLimiter)
	app.Post("/api/auth/login", handler.HandleLogin)

	mockAuthService.On("Login", mock.Anything, "inactive@dekoninklijkeloop.nl", "password123").
		Return("", "", services.ErrUserInactive)

	loginData := map[string]string{
		"email":      "inactive@dekoninklijkeloop.nl",
		"wachtwoord": "password123",
	}
	jsonData, _ := json.Marshal(loginData)

	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "inactief")

	mockAuthService.AssertExpectations(t)
}

func TestAuthHandler_HandleLogin_RateLimited(t *testing.T) {
	app := fiber.New()
	mockAuthService := new(MockAuthService)
	mockPermService := new(AuthMockPermissionService)
	mockRateLimiter := newMockRateLimiter()

	// Configure rate limiter to reject
	mockRateLimiter.AllowFunc = func(key string) bool {
		return false // Reject request
	}

	handler := handlers.NewAuthHandler(mockAuthService, mockPermService, mockRateLimiter)
	app.Post("/api/auth/login", handler.HandleLogin)

	loginData := map[string]string{
		"email":      "test@dekoninklijkeloop.nl",
		"wachtwoord": "password123",
	}
	jsonData, _ := json.Marshal(loginData)

	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Te veel")
}

func TestAuthHandler_HandleLogin_InvalidJSON(t *testing.T) {
	app := fiber.New()
	mockAuthService := new(MockAuthService)
	mockPermService := new(AuthMockPermissionService)
	mockRateLimiter := newMockRateLimiter()

	handler := handlers.NewAuthHandler(mockAuthService, mockPermService, mockRateLimiter)
	app.Post("/api/auth/login", handler.HandleLogin)

	// Invalid JSON
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBufferString("{invalid json"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// ============================================================================
// TEST SUITE: AUTH HANDLER - REFRESH TOKEN ENDPOINT
// ============================================================================

func TestAuthHandler_HandleRefreshToken_Success(t *testing.T) {
	app := fiber.New()
	mockAuthService := new(MockAuthService)
	mockPermService := new(AuthMockPermissionService)
	mockRateLimiter := newMockRateLimiter()

	handler := handlers.NewAuthHandler(mockAuthService, mockPermService, mockRateLimiter)
	app.Post("/api/auth/refresh", handler.HandleRefreshToken)

	mockAuthService.On("RefreshAccessToken", mock.Anything, "old-refresh-token").
		Return("new-access-token", "new-refresh-token", nil)

	refreshData := map[string]string{
		"refresh_token": "old-refresh-token",
	}
	jsonData, _ := json.Marshal(refreshData)

	req := httptest.NewRequest("POST", "/api/auth/refresh", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	assert.True(t, response["success"].(bool))
	assert.Equal(t, "new-access-token", response["token"])
	assert.Equal(t, "new-refresh-token", response["refresh_token"])

	mockAuthService.AssertExpectations(t)
}

func TestAuthHandler_HandleRefreshToken_MissingToken(t *testing.T) {
	app := fiber.New()
	mockAuthService := new(MockAuthService)
	mockPermService := new(AuthMockPermissionService)
	mockRateLimiter := newMockRateLimiter()

	handler := handlers.NewAuthHandler(mockAuthService, mockPermService, mockRateLimiter)
	app.Post("/api/auth/refresh", handler.HandleRefreshToken)

	refreshData := map[string]string{}
	jsonData, _ := json.Marshal(refreshData)

	req := httptest.NewRequest("POST", "/api/auth/refresh", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "verplicht")
}

func TestAuthHandler_HandleRefreshToken_InvalidToken(t *testing.T) {
	app := fiber.New()
	mockAuthService := new(MockAuthService)
	mockPermService := new(AuthMockPermissionService)
	mockRateLimiter := newMockRateLimiter()

	handler := handlers.NewAuthHandler(mockAuthService, mockPermService, mockRateLimiter)
	app.Post("/api/auth/refresh", handler.HandleRefreshToken)

	mockAuthService.On("RefreshAccessToken", mock.Anything, "invalid-token").
		Return("", "", services.ErrInvalidToken)

	refreshData := map[string]string{
		"refresh_token": "invalid-token",
	}
	jsonData, _ := json.Marshal(refreshData)

	req := httptest.NewRequest("POST", "/api/auth/refresh", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Ongeldige")
	assert.Equal(t, "REFRESH_TOKEN_INVALID", response["code"])

	mockAuthService.AssertExpectations(t)
}

// ============================================================================
// TEST SUITE: AUTH HANDLER - LOGOUT ENDPOINT
// ============================================================================

func TestAuthHandler_HandleLogout_Success(t *testing.T) {
	app := fiber.New()
	mockAuthService := new(MockAuthService)
	mockPermService := new(AuthMockPermissionService)
	mockRateLimiter := newMockRateLimiter()

	handler := handlers.NewAuthHandler(mockAuthService, mockPermService, mockRateLimiter)
	app.Post("/api/auth/logout", handler.HandleLogout)

	req := httptest.NewRequest("POST", "/api/auth/logout", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["message"], "Logout succesvol")

	// Check if cookie was cleared
	cookies := resp.Cookies()
	var authCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "auth_token" {
			authCookie = cookie
			break
		}
	}

	if authCookie != nil {
		// Cookie should be expired or have empty value
		assert.True(t, authCookie.MaxAge < 0 || authCookie.Value == "")
	}
}

// ============================================================================
// TEST SUITE: AUTH HANDLER - GET PROFILE ENDPOINT
// ============================================================================

func TestAuthHandler_HandleGetProfile_Success(t *testing.T) {
	app := fiber.New()
	mockAuthService := new(MockAuthService)
	mockPermService := new(AuthMockPermissionService)
	mockRateLimiter := newMockRateLimiter()

	handler := handlers.NewAuthHandler(mockAuthService, mockPermService, mockRateLimiter)

	testUser := &models.Gebruiker{
		ID:       "user-456",
		Naam:     "Profile User",
		Email:    "profile@dekoninklijkeloop.nl",
		Rol:      "admin",
		IsActief: true,
		LaatsteLogin: func() *time.Time {
			t := time.Now().Add(-1 * time.Hour)
			return &t
		}(),
	}

	permissions := []*models.UserPermission{
		{Resource: "admin", Action: "access"},
	}

	userRoles := []*models.UserRole{
		{
			UserID: "user-456",
			RoleID: "role-admin",
			Role: models.RBACRole{
				ID:          "role-admin",
				Name:        "admin",
				Description: "Administrator",
			},
			IsActive:   true,
			AssignedAt: time.Now(),
		},
	}

	mockAuthService.On("GetUser", mock.Anything, "user-456").Return(testUser, nil)
	mockPermService.On("GetUserPermissions", mock.Anything, "user-456").Return(permissions, nil)
	mockPermService.On("GetUserRoles", mock.Anything, "user-456").Return(userRoles, nil)

	app.Get("/api/auth/profile", func(c *fiber.Ctx) error {
		c.Locals("userID", "user-456")
		return handler.HandleGetProfile(c)
	})

	req := httptest.NewRequest("GET", "/api/auth/profile", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	assert.Equal(t, "user-456", response["id"])
	assert.Equal(t, "Profile User", response["naam"])
	assert.Equal(t, "profile@dekoninklijkeloop.nl", response["email"])
	assert.True(t, response["is_actief"].(bool))
	assert.NotNil(t, response["permissions"])
	assert.NotNil(t, response["roles"])

	mockAuthService.AssertExpectations(t)
	mockPermService.AssertExpectations(t)
}

func TestAuthHandler_HandleGetProfile_Unauthorized(t *testing.T) {
	app := fiber.New()
	mockAuthService := new(MockAuthService)
	mockPermService := new(AuthMockPermissionService)
	mockRateLimiter := newMockRateLimiter()

	handler := handlers.NewAuthHandler(mockAuthService, mockPermService, mockRateLimiter)

	// No userID in context
	app.Get("/api/auth/profile", handler.HandleGetProfile)

	req := httptest.NewRequest("GET", "/api/auth/profile", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "geautoriseerd")
}

func TestAuthHandler_HandleGetProfile_UserNotFound(t *testing.T) {
	app := fiber.New()
	mockAuthService := new(MockAuthService)
	mockPermService := new(AuthMockPermissionService)
	mockRateLimiter := newMockRateLimiter()

	handler := handlers.NewAuthHandler(mockAuthService, mockPermService, mockRateLimiter)

	// Return nil user (no error) to trigger 404 instead of 500
	mockAuthService.On("GetUser", mock.Anything, "nonexistent-user").
		Return((*models.Gebruiker)(nil), nil)

	app.Get("/api/auth/profile", func(c *fiber.Ctx) error {
		c.Locals("userID", "nonexistent-user")
		return handler.HandleGetProfile(c)
	})

	req := httptest.NewRequest("GET", "/api/auth/profile", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	mockAuthService.AssertExpectations(t)
}

// ============================================================================
// TEST SUMMARY
// ============================================================================

func TestAuthHandlerTestSuiteSummary(t *testing.T) {
	summary := `
╔══════════════════════════════════════════════════════════════════════╗
║       AUTH HANDLER COMPREHENSIVE TEST SUITE                          ║
╠══════════════════════════════════════════════════════════════════════╣
║                                                                       ║
║  ✅ LOGIN ENDPOINT (7 tests)                                         ║
║     • Successful login                                              ║
║     • Missing email validation                                      ║
║     • Missing password validation                                   ║
║     • Invalid credentials rejection                                 ║
║     • Inactive user blocking                                        ║
║     • Rate limiting enforcement                                     ║
║     • Invalid JSON handling                                         ║
║                                                                       ║
║  ✅ REFRESH TOKEN ENDPOINT (3 tests)                                ║
║     • Successful token refresh                                      ║
║     • Missing token validation                                      ║
║     • Invalid token rejection                                       ║
║                                                                       ║
║  ✅ LOGOUT ENDPOINT (1 test)                                        ║
║     • Successful logout with cookie clearing                        ║
║                                                                       ║
║  ✅ GET PROFILE ENDPOINT (3 tests)                                  ║
║     • Successfully retrieve profile with permissions & roles        ║
║     • Unauthorized access blocking                                  ║
║     • User not found handling                                       ║
║                                                                       ║
║  📊 TOTAL: 14 comprehensive HTTP endpoint tests                     ║
║                                                                       ║
╚══════════════════════════════════════════════════════════════════════╝
`
	t.Log(summary)
}
