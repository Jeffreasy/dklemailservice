package tests

import (
	"context"
	"dklautomationgo/models"
	"dklautomationgo/services"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// ============================================================================
// MOCK REPOSITORIES FOR AUTH SERVICE
// ============================================================================

type AuthMockGebruikerRepository struct {
	mock.Mock
}

func (m *AuthMockGebruikerRepository) Create(ctx context.Context, gebruiker *models.Gebruiker) error {
	args := m.Called(ctx, gebruiker)
	return args.Error(0)
}

func (m *AuthMockGebruikerRepository) GetByID(ctx context.Context, id string) (*models.Gebruiker, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Gebruiker), args.Error(1)
}

func (m *AuthMockGebruikerRepository) GetByEmail(ctx context.Context, email string) (*models.Gebruiker, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Gebruiker), args.Error(1)
}

func (m *AuthMockGebruikerRepository) Update(ctx context.Context, gebruiker *models.Gebruiker) error {
	args := m.Called(ctx, gebruiker)
	return args.Error(0)
}

func (m *AuthMockGebruikerRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *AuthMockGebruikerRepository) List(ctx context.Context, limit, offset int) ([]*models.Gebruiker, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Gebruiker), args.Error(1)
}

func (m *AuthMockGebruikerRepository) UpdateLastLogin(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *AuthMockGebruikerRepository) GetNewsletterSubscribers(ctx context.Context) ([]*models.Gebruiker, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Gebruiker), args.Error(1)
}

type AuthMockRefreshTokenRepository struct {
	mock.Mock
}

func (m *AuthMockRefreshTokenRepository) Create(ctx context.Context, token *models.RefreshToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *AuthMockRefreshTokenRepository) GetByToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RefreshToken), args.Error(1)
}

func (m *AuthMockRefreshTokenRepository) RevokeToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *AuthMockRefreshTokenRepository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *AuthMockRefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// Use the existing MockUserRoleRepository from permission_service_test.go

// ============================================================================
// TEST SUITE: AUTH SERVICE - LOGIN FLOW
// ============================================================================

func TestAuthService_Login_Success(t *testing.T) {
	// Setup
	mockUserRepo := new(AuthMockGebruikerRepository)
	mockRefreshRepo := new(AuthMockRefreshTokenRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)

	// Create test user with bcrypt hashed password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass123"), bcrypt.DefaultCost)
	testUser := &models.Gebruiker{
		ID:             "user-123",
		Email:          "test@dekoninklijkeloop.nl",
		Naam:           "Test User",
		WachtwoordHash: string(hashedPassword),
		Rol:            "admin",
		IsActief:       true,
		LaatsteLogin:   nil,
	}

	// Mock expectations
	mockUserRepo.On("GetByEmail", mock.Anything, "test@dekoninklijkeloop.nl").Return(testUser, nil)
	mockUserRepo.On("UpdateLastLogin", mock.Anything, "user-123").Return(nil)
	mockRefreshRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.RefreshToken")).Return(nil)
	mockUserRoleRepo.On("ListActiveByUser", mock.Anything, "user-123").Return([]*models.UserRole{
		{
			UserID: "user-123",
			RoleID: "role-admin",
			Role: models.RBACRole{
				ID:   "role-admin",
				Name: "admin",
			},
			IsActive: true,
		},
	}, nil)

	// Create service
	authService := services.NewAuthServiceWithRBAC(mockUserRepo, mockRefreshRepo, mockUserRoleRepo)

	// Execute
	ctx := context.Background()
	accessToken, refreshToken, err := authService.Login(ctx, "test@dekoninklijkeloop.nl", "testpass123")

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)

	// Verify JWT token contains correct claims
	token, err := jwt.ParseWithClaims(accessToken, &services.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("default_jwt_secret_change_in_production"), nil
	})
	require.NoError(t, err)
	claims := token.Claims.(*services.JWTClaims)
	assert.Equal(t, "test@dekoninklijkeloop.nl", claims.Email)
	assert.Equal(t, "admin", claims.Role)
	assert.Equal(t, []string{"admin"}, claims.Roles)
	assert.True(t, claims.RBACActive)
	assert.Equal(t, "user-123", claims.Subject)

	mockUserRepo.AssertExpectations(t)
	mockRefreshRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	mockUserRepo := new(AuthMockGebruikerRepository)
	mockRefreshRepo := new(AuthMockRefreshTokenRepository)

	mockUserRepo.On("GetByEmail", mock.Anything, "invalid@dekoninklijkeloop.nl").Return(nil, nil)

	authService := services.NewAuthService(mockUserRepo, mockRefreshRepo)

	ctx := context.Background()
	accessToken, refreshToken, err := authService.Login(ctx, "invalid@dekoninklijkeloop.nl", "wrongpass")

	assert.Error(t, err)
	assert.Equal(t, services.ErrInvalidCredentials, err)
	assert.Empty(t, accessToken)
	assert.Empty(t, refreshToken)

	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_Login_InactiveUser(t *testing.T) {
	mockUserRepo := new(AuthMockGebruikerRepository)
	mockRefreshRepo := new(AuthMockRefreshTokenRepository)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass123"), bcrypt.DefaultCost)
	inactiveUser := &models.Gebruiker{
		ID:             "user-456",
		Email:          "inactive@dekoninklijkeloop.nl",
		WachtwoordHash: string(hashedPassword),
		IsActief:       false,
	}

	mockUserRepo.On("GetByEmail", mock.Anything, "inactive@dekoninklijkeloop.nl").Return(inactiveUser, nil)

	authService := services.NewAuthService(mockUserRepo, mockRefreshRepo)

	ctx := context.Background()
	accessToken, refreshToken, err := authService.Login(ctx, "inactive@dekoninklijkeloop.nl", "testpass123")

	assert.Error(t, err)
	assert.Equal(t, services.ErrUserInactive, err)
	assert.Empty(t, accessToken)
	assert.Empty(t, refreshToken)

	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	mockUserRepo := new(AuthMockGebruikerRepository)
	mockRefreshRepo := new(AuthMockRefreshTokenRepository)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpass"), bcrypt.DefaultCost)
	testUser := &models.Gebruiker{
		ID:             "user-789",
		Email:          "test@dekoninklijkeloop.nl",
		WachtwoordHash: string(hashedPassword),
		IsActief:       true,
	}

	mockUserRepo.On("GetByEmail", mock.Anything, "test@dekoninklijkeloop.nl").Return(testUser, nil)

	authService := services.NewAuthService(mockUserRepo, mockRefreshRepo)

	ctx := context.Background()
	accessToken, refreshToken, err := authService.Login(ctx, "test@dekoninklijkeloop.nl", "wrongpass")

	assert.Error(t, err)
	assert.Equal(t, services.ErrInvalidCredentials, err)
	assert.Empty(t, accessToken)
	assert.Empty(t, refreshToken)

	mockUserRepo.AssertExpectations(t)
}

// ============================================================================
// TEST SUITE: JWT TOKEN VALIDATION
// ============================================================================

func TestAuthService_ValidateToken_Success(t *testing.T) {
	mockUserRepo := new(AuthMockGebruikerRepository)
	mockRefreshRepo := new(AuthMockRefreshTokenRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass123"), bcrypt.DefaultCost)
	testUser := &models.Gebruiker{
		ID:             "user-valid-123",
		Email:          "valid@dekoninklijkeloop.nl",
		WachtwoordHash: string(hashedPassword),
		IsActief:       true,
	}

	mockUserRepo.On("GetByEmail", mock.Anything, "valid@dekoninklijkeloop.nl").Return(testUser, nil)
	mockUserRepo.On("UpdateLastLogin", mock.Anything, "user-valid-123").Return(nil)
	mockRefreshRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.RefreshToken")).Return(nil)
	mockUserRoleRepo.On("ListActiveByUser", mock.Anything, "user-valid-123").Return([]*models.UserRole{}, nil)

	authService := services.NewAuthServiceWithRBAC(mockUserRepo, mockRefreshRepo, mockUserRoleRepo)

	// Generate token via login
	ctx := context.Background()
	token, _, err := authService.Login(ctx, "valid@dekoninklijkeloop.nl", "testpass123")
	require.NoError(t, err)

	// Validate the token
	userID, err := authService.ValidateToken(token)

	assert.NoError(t, err)
	assert.Equal(t, "user-valid-123", userID)

	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_ValidateToken_InvalidFormat(t *testing.T) {
	mockUserRepo := new(AuthMockGebruikerRepository)
	mockRefreshRepo := new(AuthMockRefreshTokenRepository)

	authService := services.NewAuthService(mockUserRepo, mockRefreshRepo)

	testCases := []struct {
		name  string
		token string
	}{
		{"Empty token", ""},
		{"Invalid format", "invalid.token.format"},
		{"Malformed JWT", "eyJhbGc.invalid"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userID, err := authService.ValidateToken(tc.token)
			assert.Error(t, err)
			assert.Equal(t, services.ErrInvalidToken, err)
			assert.Empty(t, userID)
		})
	}
}

func TestAuthService_ValidateToken_ExpiredToken(t *testing.T) {
	// Create an expired token manually
	claims := services.JWTClaims{
		Email: "expired@dekoninklijkeloop.nl",
		Role:  "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			Subject:   "user-expired",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte("default_jwt_secret_change_in_production"))

	mockUserRepo := new(AuthMockGebruikerRepository)
	mockRefreshRepo := new(AuthMockRefreshTokenRepository)

	authService := services.NewAuthService(mockUserRepo, mockRefreshRepo)

	userID, err := authService.ValidateToken(signedToken)

	assert.Error(t, err)
	assert.Equal(t, services.ErrInvalidToken, err)
	assert.Empty(t, userID)
}

func TestAuthService_ValidateToken_WithBearerPrefix(t *testing.T) {
	mockUserRepo := new(AuthMockGebruikerRepository)
	mockRefreshRepo := new(AuthMockRefreshTokenRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass123"), bcrypt.DefaultCost)
	testUser := &models.Gebruiker{
		ID:             "user-bearer-123",
		Email:          "bearer@dekoninklijkeloop.nl",
		WachtwoordHash: string(hashedPassword),
		IsActief:       true,
	}

	mockUserRepo.On("GetByEmail", mock.Anything, "bearer@dekoninklijkeloop.nl").Return(testUser, nil)
	mockUserRepo.On("UpdateLastLogin", mock.Anything, "user-bearer-123").Return(nil)
	mockRefreshRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.RefreshToken")).Return(nil)
	mockUserRoleRepo.On("ListActiveByUser", mock.Anything, "user-bearer-123").Return([]*models.UserRole{}, nil)

	authService := services.NewAuthServiceWithRBAC(mockUserRepo, mockRefreshRepo, mockUserRoleRepo)

	ctx := context.Background()
	token, _, err := authService.Login(ctx, "bearer@dekoninklijkeloop.nl", "testpass123")
	require.NoError(t, err)

	// Test with Bearer prefix
	bearerToken := "Bearer " + token
	userID, err := authService.ValidateToken(bearerToken)

	assert.NoError(t, err)
	assert.Equal(t, "user-bearer-123", userID)

	mockUserRepo.AssertExpectations(t)
}

// ============================================================================
// TEST SUITE: REFRESH TOKEN FLOW
// ============================================================================

func TestAuthService_RefreshAccessToken_Success(t *testing.T) {
	mockUserRepo := new(AuthMockGebruikerRepository)
	mockRefreshRepo := new(AuthMockRefreshTokenRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)

	testUser := &models.Gebruiker{
		ID:       "user-refresh-123",
		Email:    "refresh@dekoninklijkeloop.nl",
		Naam:     "Refresh User",
		Rol:      "user",
		IsActief: true,
	}

	validRefreshToken := &models.RefreshToken{
		ID:        "token-123",
		UserID:    "user-refresh-123",
		Token:     "valid-refresh-token-xyz",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		IsRevoked: false,
	}

	mockRefreshRepo.On("GetByToken", mock.Anything, "valid-refresh-token-xyz").Return(validRefreshToken, nil)
	mockUserRepo.On("GetByID", mock.Anything, "user-refresh-123").Return(testUser, nil)
	mockRefreshRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.RefreshToken")).Return(nil)
	mockRefreshRepo.On("RevokeToken", mock.Anything, "valid-refresh-token-xyz").Return(nil)
	mockUserRoleRepo.On("ListActiveByUser", mock.Anything, "user-refresh-123").Return([]*models.UserRole{}, nil)

	authService := services.NewAuthServiceWithRBAC(mockUserRepo, mockRefreshRepo, mockUserRoleRepo)

	ctx := context.Background()
	newAccessToken, newRefreshToken, err := authService.RefreshAccessToken(ctx, "valid-refresh-token-xyz")

	assert.NoError(t, err)
	assert.NotEmpty(t, newAccessToken)
	assert.NotEmpty(t, newRefreshToken)
	assert.NotEqual(t, "valid-refresh-token-xyz", newRefreshToken, "New refresh token should be different (token rotation)")

	mockRefreshRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_RefreshAccessToken_InvalidToken(t *testing.T) {
	mockUserRepo := new(AuthMockGebruikerRepository)
	mockRefreshRepo := new(AuthMockRefreshTokenRepository)

	mockRefreshRepo.On("GetByToken", mock.Anything, "invalid-token").Return(nil, nil)

	authService := services.NewAuthService(mockUserRepo, mockRefreshRepo)

	ctx := context.Background()
	accessToken, refreshToken, err := authService.RefreshAccessToken(ctx, "invalid-token")

	assert.Error(t, err)
	assert.Empty(t, accessToken)
	assert.Empty(t, refreshToken)

	mockRefreshRepo.AssertExpectations(t)
}

func TestAuthService_RefreshAccessToken_ExpiredToken(t *testing.T) {
	mockUserRepo := new(AuthMockGebruikerRepository)
	mockRefreshRepo := new(AuthMockRefreshTokenRepository)

	expiredRefreshToken := &models.RefreshToken{
		ID:        "token-expired",
		UserID:    "user-123",
		Token:     "expired-token",
		ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired 1 hour ago
		IsRevoked: false,
	}

	mockRefreshRepo.On("GetByToken", mock.Anything, "expired-token").Return(expiredRefreshToken, nil)

	authService := services.NewAuthService(mockUserRepo, mockRefreshRepo)

	ctx := context.Background()
	accessToken, refreshToken, err := authService.RefreshAccessToken(ctx, "expired-token")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ongeldige of verlopen")
	assert.Empty(t, accessToken)
	assert.Empty(t, refreshToken)

	mockRefreshRepo.AssertExpectations(t)
}

func TestAuthService_RefreshAccessToken_InactiveUser(t *testing.T) {
	mockUserRepo := new(AuthMockGebruikerRepository)
	mockRefreshRepo := new(AuthMockRefreshTokenRepository)

	inactiveUser := &models.Gebruiker{
		ID:       "user-inactive-456",
		Email:    "inactive@dekoninklijkeloop.nl",
		IsActief: false,
	}

	validRefreshToken := &models.RefreshToken{
		UserID:    "user-inactive-456",
		Token:     "valid-but-user-inactive",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		IsRevoked: false,
	}

	mockRefreshRepo.On("GetByToken", mock.Anything, "valid-but-user-inactive").Return(validRefreshToken, nil)
	mockUserRepo.On("GetByID", mock.Anything, "user-inactive-456").Return(inactiveUser, nil)

	authService := services.NewAuthService(mockUserRepo, mockRefreshRepo)

	ctx := context.Background()
	accessToken, refreshToken, err := authService.RefreshAccessToken(ctx, "valid-but-user-inactive")

	assert.Error(t, err)
	assert.Equal(t, services.ErrUserInactive, err)
	assert.Empty(t, accessToken)
	assert.Empty(t, refreshToken)

	mockRefreshRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

// ============================================================================
// TEST SUITE: PASSWORD OPERATIONS
// ============================================================================

func TestAuthService_HashPassword_Success(t *testing.T) {
	mockUserRepo := new(AuthMockGebruikerRepository)
	mockRefreshRepo := new(AuthMockRefreshTokenRepository)

	authService := services.NewAuthService(mockUserRepo, mockRefreshRepo)

	password := "SuperSecurePassword123!"
	hash, err := authService.HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)
	assert.True(t, strings.HasPrefix(hash, "$2a$") || strings.HasPrefix(hash, "$2b$"), "Should be bcrypt hash")
}

func TestAuthService_VerifyPassword_Success(t *testing.T) {
	mockUserRepo := new(AuthMockGebruikerRepository)
	mockRefreshRepo := new(AuthMockRefreshTokenRepository)

	authService := services.NewAuthService(mockUserRepo, mockRefreshRepo)

	password := "TestPassword456"
	hash, err := authService.HashPassword(password)
	require.NoError(t, err)

	isValid := authService.VerifyPassword(hash, password)
	assert.True(t, isValid)
}

func TestAuthService_VerifyPassword_WrongPassword(t *testing.T) {
	mockUserRepo := new(AuthMockGebruikerRepository)
	mockRefreshRepo := new(AuthMockRefreshTokenRepository)

	authService := services.NewAuthService(mockUserRepo, mockRefreshRepo)

	password := "CorrectPassword"
	hash, err := authService.HashPassword(password)
	require.NoError(t, err)

	isValid := authService.VerifyPassword(hash, "WrongPassword")
	assert.False(t, isValid)
}

func TestAuthService_ResetPassword_Success(t *testing.T) {
	mockUserRepo := new(AuthMockGebruikerRepository)
	mockRefreshRepo := new(AuthMockRefreshTokenRepository)

	existingUser := &models.Gebruiker{
		ID:    "user-reset-123",
		Email: "reset@dekoninklijkeloop.nl",
	}

	mockUserRepo.On("GetByEmail", mock.Anything, "reset@dekoninklijkeloop.nl").Return(existingUser, nil)
	mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Gebruiker")).Return(nil)

	authService := services.NewAuthService(mockUserRepo, mockRefreshRepo)

	ctx := context.Background()
	err := authService.ResetPassword(ctx, "reset@dekoninklijkeloop.nl", "NewSecurePassword123")

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_ResetPassword_UserNotFound(t *testing.T) {
	mockUserRepo := new(AuthMockGebruikerRepository)
	mockRefreshRepo := new(AuthMockRefreshTokenRepository)

	mockUserRepo.On("GetByEmail", mock.Anything, "nonexistent@dekoninklijkeloop.nl").Return(nil, nil)

	authService := services.NewAuthService(mockUserRepo, mockRefreshRepo)

	ctx := context.Background()
	err := authService.ResetPassword(ctx, "nonexistent@dekoninklijkeloop.nl", "NewPassword")

	assert.Error(t, err)
	assert.Equal(t, services.ErrUserNotFound, err)

	mockUserRepo.AssertExpectations(t)
}

// ===========================================================================
// TEST SUITE: RBAC INTEGRATION
// ============================================================================

func TestAuthService_Login_WithRBACRoles(t *testing.T) {
	mockUserRepo := new(AuthMockGebruikerRepository)
	mockRefreshRepo := new(AuthMockRefreshTokenRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass123"), bcrypt.DefaultCost)
	testUser := &models.Gebruiker{
		ID:             "user-rbac-123",
		Email:          "rbac@dekoninklijkeloop.nl",
		WachtwoordHash: string(hashedPassword),
		Rol:            "admin",
		IsActief:       true,
	}

	userRoles := []*models.UserRole{
		{
			UserID: "user-rbac-123",
			RoleID: "role-admin",
			Role: models.RBACRole{
				ID:          "role-admin",
				Name:        "admin",
				Description: "Administrator role",
			},
			IsActive: true,
		},
		{
			UserID: "user-rbac-123",
			RoleID: "role-staff",
			Role: models.RBACRole{
				ID:          "role-staff",
				Name:        "staff",
				Description: "Staff role",
			},
			IsActive: true,
		},
	}

	mockUserRepo.On("GetByEmail", mock.Anything, "rbac@dekoninklijkeloop.nl").Return(testUser, nil)
	mockUserRepo.On("UpdateLastLogin", mock.Anything, "user-rbac-123").Return(nil)
	mockRefreshRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.RefreshToken")).Return(nil)
	mockUserRoleRepo.On("ListActiveByUser", mock.Anything, "user-rbac-123").Return(userRoles, nil)

	authService := services.NewAuthServiceWithRBAC(mockUserRepo, mockRefreshRepo, mockUserRoleRepo)

	ctx := context.Background()
	accessToken, _, err := authService.Login(ctx, "rbac@dekoninklijkeloop.nl", "testpass123")

	require.NoError(t, err)
	assert.NotEmpty(t, accessToken)

	// Verify JWT contains RBAC roles
	token, err := jwt.ParseWithClaims(accessToken, &services.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("default_jwt_secret_change_in_production"), nil
	})
	require.NoError(t, err)

	claims := token.Claims.(*services.JWTClaims)
	assert.Equal(t, []string{"admin", "staff"}, claims.Roles)
	assert.True(t, claims.RBACActive)

	mockUserRepo.AssertExpectations(t)
	mockUserRoleRepo.AssertExpectations(t)
}

func TestAuthService_Login_WithoutRBAC_FallbackToLegacy(t *testing.T) {
	mockUserRepo := new(AuthMockGebruikerRepository)
	mockRefreshRepo := new(AuthMockRefreshTokenRepository)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass123"), bcrypt.DefaultCost)
	testUser := &models.Gebruiker{
		ID:             "user-legacy-123",
		Email:          "legacy@dekoninklijkeloop.nl",
		WachtwoordHash: string(hashedPassword),
		Rol:            "user",
		IsActief:       true,
	}

	mockUserRepo.On("GetByEmail", mock.Anything, "legacy@dekoninklijkeloop.nl").Return(testUser, nil)
	mockUserRepo.On("UpdateLastLogin", mock.Anything, "user-legacy-123").Return(nil)
	mockRefreshRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.RefreshToken")).Return(nil)

	// Create without RBAC repo (legacy mode)
	authService := services.NewAuthService(mockUserRepo, mockRefreshRepo)

	ctx := context.Background()
	accessToken, _, err := authService.Login(ctx, "legacy@dekoninklijkeloop.nl", "testpass123")

	require.NoError(t, err)

	// Verify JWT has legacy role but no RBAC
	token, err := jwt.ParseWithClaims(accessToken, &services.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("default_jwt_secret_change_in_production"), nil
	})
	require.NoError(t, err)

	claims := token.Claims.(*services.JWTClaims)
	assert.Equal(t, "user", claims.Role) // Legacy field
	assert.Empty(t, claims.Roles)        // No RBAC roles
	assert.False(t, claims.RBACActive)

	mockUserRepo.AssertExpectations(t)
}

// ============================================================================
// TEST SUMMARY
// ============================================================================

func TestAuthServiceTestSuiteSummary(t *testing.T) {
	summary := `
╔══════════════════════════════════════════════════════════════════════╗
║       AUTH SERVICE COMPREHENSIVE TEST SUITE                          ║
╠══════════════════════════════════════════════════════════════════════╣
║                                                                       ║
║  ✅ LOGIN FLOW (4 tests)                                            ║
║     •  Successful login with valid credentials                      ║
║     • Invalid credentials rejection                                 ║
║     • Inactive user blocking                                        ║
║     • Wrong password rejection                                      ║
║                                                                       ║
║  ✅ JWT TOKEN VALIDATION (4 tests)                                  ║
║     • Token generation and validation                               ║
║     • Invalid format detection                                      ║
║     • Expired token rejection                                       ║
║     • Bearer prefix handling                                        ║
║                                                                       ║
║  ✅ REFRESH TOKEN FLOW (4 tests)                                    ║
║     • Successful token refresh with rotation                        ║
║     • Invalid token rejection                                       ║
║     • Expired token rejection                                       ║
║     • Inactive user blocking during refresh                         ║
║                                                                       ║
║  ✅ PASSWORD OPERATIONS (4 tests)                                   ║
║     • Password hashing (bcrypt)                                     ║
║     • Password verification                                         ║
║     • Password reset flow                                           ║
║     • User not found handling                                       ║
║                                                                       ║
║  ✅ RBAC INTEGRATION (2 tests)                                      ║
║     • Multi-role JWT generation                                     ║
║     • Legacy fallback when RBAC unavailable                         ║
║                                                                       ║
║  📊 TOTAL: 18 comprehensive test cases                              ║
║                                                                       ║
╚══════════════════════════════════════════════════════════════════════╝
`
	t.Log(summary)
}
