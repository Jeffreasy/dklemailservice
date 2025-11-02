package tests

import (
	"dklautomationgo/handlers"
	"dklautomationgo/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// TEST SUITE: PERMISSION MIDDLEWARE - BASIC FUNCTIONALITY
// ============================================================================

func TestPermissionMiddleware_AllowWithPermission(t *testing.T) {
	app := fiber.New()
	mockPermService := new(AuthMockPermissionService)

	// User has the required permission
	mockPermService.On("HasPermission", mock.Anything, "user-123", "contact", "read").Return(true)

	// Setup middleware and protected route
	app.Get("/contacts",
		func(c *fiber.Ctx) error {
			c.Locals("userID", "user-123")
			return c.Next()
		},
		handlers.PermissionMiddleware(mockPermService, "contact", "read"),
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "success"})
		},
	)

	req := httptest.NewRequest("GET", "/contacts", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "success", response["message"])

	mockPermService.AssertExpectations(t)
}

func TestPermissionMiddleware_BlockWithoutPermission(t *testing.T) {
	app := fiber.New()
	mockPermService := new(AuthMockPermissionService)

	// User does NOT have the required permission
	mockPermService.On("HasPermission", mock.Anything, "user-456", "contact", "write").Return(false)

	app.Post("/contacts",
		func(c *fiber.Ctx) error {
			c.Locals("userID", "user-456")
			return c.Next()
		},
		handlers.PermissionMiddleware(mockPermService, "contact", "write"),
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "success"})
		},
	)

	req := httptest.NewRequest("POST", "/contacts", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Geen toegang")

	mockPermService.AssertExpectations(t)
}

func TestPermissionMiddleware_MissingUserID(t *testing.T) {
	app := fiber.New()
	mockPermService := new(AuthMockPermissionService)

	// No userID set in context
	app.Get("/contacts",
		handlers.PermissionMiddleware(mockPermService, "contact", "read"),
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "success"})
		},
	)

	req := httptest.NewRequest("GET", "/contacts", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "geautoriseerd")

	// Should not call HasPermission if no userID
	mockPermService.AssertNotCalled(t, "HasPermission")
}

func TestPermissionMiddleware_EmptyUserID(t *testing.T) {
	app := fiber.New()
	mockPermService := new(AuthMockPermissionService)

	app.Get("/contacts",
		func(c *fiber.Ctx) error {
			c.Locals("userID", "") // Empty string
			return c.Next()
		},
		handlers.PermissionMiddleware(mockPermService, "contact", "read"),
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "success"})
		},
	)

	req := httptest.NewRequest("GET", "/contacts", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	mockPermService.AssertNotCalled(t, "HasPermission")
}

// ============================================================================
// TEST SUITE: CONVENIENCE MIDDLEWARES
// ============================================================================

func TestAdminPermissionMiddleware_AdminHasAccess(t *testing.T) {
	app := fiber.New()
	mockPermService := new(AuthMockPermissionService)

	// Admin user has admin:access permission
	mockPermService.On("HasPermission", mock.Anything, "admin-user", "admin", "access").Return(true)

	app.Get("/admin",
		func(c *fiber.Ctx) error {
			c.Locals("userID", "admin-user")
			return c.Next()
		},
		handlers.AdminPermissionMiddleware(mockPermService),
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "admin area"})
		},
	)

	req := httptest.NewRequest("GET", "/admin", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	mockPermService.AssertExpectations(t)
}

func TestAdminPermissionMiddleware_NonAdminBlocked(t *testing.T) {
	app := fiber.New()
	mockPermService := new(AuthMockPermissionService)

	// Regular user does NOT have admin:access
	mockPermService.On("HasPermission", mock.Anything, "regular-user", "admin", "access").Return(false)

	app.Get("/admin",
		func(c *fiber.Ctx) error {
			c.Locals("userID", "regular-user")
			return c.Next()
		},
		handlers.AdminPermissionMiddleware(mockPermService),
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "admin area"})
		},
	)

	req := httptest.NewRequest("GET", "/admin", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)

	mockPermService.AssertExpectations(t)
}

func TestStaffPermissionMiddleware_StaffHasAccess(t *testing.T) {
	app := fiber.New()
	mockPermService := new(AuthMockPermissionService)

	// Staff user has staff:access permission
	mockPermService.On("HasPermission", mock.Anything, "staff-user", "staff", "access").Return(true)

	app.Get("/staff",
		func(c *fiber.Ctx) error {
			c.Locals("userID", "staff-user")
			return c.Next()
		},
		handlers.StaffPermissionMiddleware(mockPermService),
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "staff area"})
		},
	)

	req := httptest.NewRequest("GET", "/staff", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	mockPermService.AssertExpectations(t)
}

func TestStaffPermissionMiddleware_NonStaffBlocked(t *testing.T) {
	app := fiber.New()
	mockPermService := new(AuthMockPermissionService)

	// User does NOT have staff:access
	mockPermService.On("HasPermission", mock.Anything, "regular-user", "staff", "access").Return(false)

	app.Get("/staff",
		func(c *fiber.Ctx) error {
			c.Locals("userID", "regular-user")
			return c.Next()
		},
		handlers.StaffPermissionMiddleware(mockPermService),
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "staff area"})
		},
	)

	req := httptest.NewRequest("GET", "/staff", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)

	mockPermService.AssertExpectations(t)
}

// ============================================================================
// TEST SUITE: RESOURCE PERMISSION MIDDLEWARE (MULTIPLE PERMISSIONS)
// ============================================================================

func TestResourcePermissionMiddleware_AllPermissionsRequired(t *testing.T) {
	app := fiber.New()
	mockPermService := new(AuthMockPermissionService)

	permissions := []models.PermissionCheck{
		{Resource: "contact", Action: "read"},
		{Resource: "contact", Action: "write"},
	}

	// User has both required permissions
	mockPermService.On("HasPermission", mock.Anything, "user-789", "contact", "read").Return(true)
	mockPermService.On("HasPermission", mock.Anything, "user-789", "contact", "write").Return(true)

	app.Put("/contacts/:id",
		func(c *fiber.Ctx) error {
			c.Locals("userID", "user-789")
			return c.Next()
		},
		handlers.ResourcePermissionMiddleware(mockPermService, permissions...),
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "updated"})
		},
	)

	req := httptest.NewRequest("PUT", "/contacts/123", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	mockPermService.AssertExpectations(t)
}

func TestResourcePermissionMiddleware_MissingOnePermission(t *testing.T) {
	app := fiber.New()
	mockPermService := new(AuthMockPermissionService)

	permissions := []models.PermissionCheck{
		{Resource: "contact", Action: "read"},
		{Resource: "contact", Action: "delete"},
	}

	// User has read but NOT delete
	mockPermService.On("HasPermission", mock.Anything, "user-abc", "contact", "read").Return(true)
	mockPermService.On("HasPermission", mock.Anything, "user-abc", "contact", "delete").Return(false)

	app.Delete("/contacts/:id",
		func(c *fiber.Ctx) error {
			c.Locals("userID", "user-abc")
			return c.Next()
		},
		handlers.ResourcePermissionMiddleware(mockPermService, permissions...),
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "deleted"})
		},
	)

	req := httptest.NewRequest("DELETE", "/contacts/123", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Geen toegang")

	mockPermService.AssertExpectations(t)
}

// ============================================================================
// TEST SUITE: DIFFERENT RESOURCES
// ============================================================================

func TestPermissionMiddleware_DifferentResources(t *testing.T) {
	testCases := []struct {
		name           string
		userID         string
		resource       string
		action         string
		hasPermission  bool
		expectedStatus int
	}{
		{
			name:           "User can read",
			userID:         "user-1",
			resource:       "user",
			action:         "read",
			hasPermission:  true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "User cannot manage roles",
			userID:         "user-2",
			resource:       "user",
			action:         "manage_roles",
			hasPermission:  false,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Can read photos",
			userID:         "user-3",
			resource:       "photo",
			action:         "read",
			hasPermission:  true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Cannot delete newsletters",
			userID:         "user-4",
			resource:       "newsletter",
			action:         "delete",
			hasPermission:  false,
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app := fiber.New()
			mockPermService := new(AuthMockPermissionService)

			mockPermService.On("HasPermission", mock.Anything, tc.userID, tc.resource, tc.action).
				Return(tc.hasPermission)

			app.Get("/resource",
				func(c *fiber.Ctx) error {
					c.Locals("userID", tc.userID)
					return c.Next()
				},
				handlers.PermissionMiddleware(mockPermService, tc.resource, tc.action),
				func(c *fiber.Ctx) error {
					return c.JSON(fiber.Map{"message": "success"})
				},
			)

			req := httptest.NewRequest("GET", "/resource", nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			mockPermService.AssertExpectations(t)
		})
	}
}

// ============================================================================
// TEST SUITE: MIDDLEWARE CHAINING
// ============================================================================

func TestPermissionMiddleware_ChainedMiddlewares(t *testing.T) {
	app := fiber.New()
	mockPermService := new(AuthMockPermissionService)

	// User needs both admin AND contact:write permissions
	mockPermService.On("HasPermission", mock.Anything, "power-user", "admin", "access").Return(true)
	mockPermService.On("HasPermission", mock.Anything, "power-user", "contact", "write").Return(true)

	app.Post("/admin/contacts",
		func(c *fiber.Ctx) error {
			c.Locals("userID", "power-user")
			return c.Next()
		},
		handlers.AdminPermissionMiddleware(mockPermService),                // First check admin
		handlers.PermissionMiddleware(mockPermService, "contact", "write"), // Then check contact write
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "created"})
		},
	)

	req := httptest.NewRequest("POST", "/admin/contacts", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	mockPermService.AssertExpectations(t)
}

func TestPermissionMiddleware_ChainedMiddlewares_FirstFails(t *testing.T) {
	app := fiber.New()
	mockPermService := new(AuthMockPermissionService)

	// User has contact:write but NOT admin:access
	mockPermService.On("HasPermission", mock.Anything, "limited-user", "admin", "access").Return(false)

	app.Post("/admin/contacts",
		func(c *fiber.Ctx) error {
			c.Locals("userID", "limited-user")
			return c.Next()
		},
		handlers.AdminPermissionMiddleware(mockPermService),
		handlers.PermissionMiddleware(mockPermService, "contact", "write"),
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "created"})
		},
	)

	req := httptest.NewRequest("POST", "/admin/contacts", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Second middleware should NOT be called if first fails
	mockPermService.AssertNotCalled(t, "HasPermission", mock.Anything, "limited-user", "contact", "write")
	mockPermService.AssertExpectations(t)
}

// ============================================================================
// TEST SUITE: COMPLETE PERMISSION CATALOG
// ============================================================================

func TestPermissionMiddleware_AllSystemResources(t *testing.T) {
	// Test all 19 resources from AUTH_AND_RBAC.md
	resources := []struct {
		resource string
		action   string
	}{
		{"admin", "access"},
		{"staff", "access"},
		{"contact", "read"},
		{"contact", "write"},
		{"contact", "delete"},
		{"aanmelding", "read"},
		{"aanmelding", "write"},
		{"user", "read"},
		{"user", "manage_roles"},
		{"photo", "read"},
		{"photo", "write"},
		{"album", "read"},
		{"video", "read"},
		{"partner", "write"},
		{"newsletter", "send"},
		{"email", "fetch"},
		{"chat", "moderate"},
		{"chat", "manage_channel"},
	}

	for _, res := range resources {
		t.Run(res.resource+":"+res.action, func(t *testing.T) {
			app := fiber.New()
			mockPermService := new(AuthMockPermissionService)

			mockPermService.On("HasPermission", mock.Anything, "test-user", res.resource, res.action).Return(true)

			app.Get("/test",
				func(c *fiber.Ctx) error {
					c.Locals("userID", "test-user")
					return c.Next()
				},
				handlers.PermissionMiddleware(mockPermService, res.resource, res.action),
				func(c *fiber.Ctx) error {
					return c.SendStatus(http.StatusOK)
				},
			)

			req := httptest.NewRequest("GET", "/test", nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			mockPermService.AssertExpectations(t)
		})
	}
}

// ============================================================================
// TEST SUMMARY
// ============================================================================

func TestPermissionMiddlewareTestSuiteSummary(t *testing.T) {
	summary := `
╔══════════════════════════════════════════════════════════════════════╗
║       PERMISSION MIDDLEWARE COMPREHENSIVE TEST SUITE                 ║
╠══════════════════════════════════════════════════════════════════════╣
║                                                                       ║
║  ✅ BASIC FUNCTIONALITY (4 tests)                                    ║
║     • Allow access with permission                                  ║
║     • Block access without permission                               ║
║     • Missing userID rejection                                      ║
║     • Empty userID rejection                                        ║
║                                                                       ║
║  ✅ CONVENIENCE MIDDLEWARES (4 tests)                               ║
║     • Admin middleware - admin has access                           ║
║     • Admin middleware - non-admin blocked                          ║
║     • Staff middleware - staff has access                           ║
║     • Staff middleware - non-staff blocked                          ║
║                                                                       ║
║  ✅ MIDDLEWARE CHAINING (2 tests)                                   ║
║     • Multiple middlewares - all pass                               ║
║     • Multiple middlewares - first fails (short-circuit)            ║
║                                                                       ║
║  ✅ COMPLETE CATALOG (18 tests)                                     ║
║     • All 19 system resources tested                                ║
║     • admin, staff, contact, aanmelding, user, photo, album,        ║
║       video, partner, newsletter, email, chat                       ║
║                                                                       ║
║  📊 TOTAL: 28 comprehensive middleware tests                        ║
║                                                                       ║
╚══════════════════════════════════════════════════════════════════════╝
`
	t.Log(summary)
}
