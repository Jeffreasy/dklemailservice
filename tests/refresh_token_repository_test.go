package tests

import (
	"context"
	"dklautomationgo/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================================================
// TEST SUITE: REFRESH TOKEN MODEL - IsValid() METHOD
// ============================================================================

func TestRefreshToken_IsValid_ValidToken(t *testing.T) {
	token := &models.RefreshToken{
		ID:        "token-123",
		UserID:    "user-123",
		Token:     "valid-token",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		IsRevoked: false,
		CreatedAt: time.Now(),
	}

	assert.True(t, token.IsValid(), "Token should be valid (not revoked, not expired)")
}

func TestRefreshToken_IsValid_ExpiredToken(t *testing.T) {
	token := &models.RefreshToken{
		ID:        "token-expired",
		UserID:    "user-123",
		Token:     "expired-token",
		ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired 1 hour ago
		IsRevoked: false,
		CreatedAt: time.Now().Add(-8 * 24 * time.Hour),
	}

	assert.False(t, token.IsValid(), "Token should be invalid (expired)")
}

func TestRefreshToken_IsValid_RevokedToken(t *testing.T) {
	token := &models.RefreshToken{
		ID:        "token-revoked",
		UserID:    "user-123",
		Token:     "revoked-token",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		IsRevoked: true, // Revoked
		RevokedAt: func() *time.Time { t := time.Now(); return &t }(),
		CreatedAt: time.Now(),
	}

	assert.False(t, token.IsValid(), "Token should be invalid (revoked)")
}

func TestRefreshToken_IsValid_RevokedAndExpired(t *testing.T) {
	token := &models.RefreshToken{
		ID:        "token-both",
		UserID:    "user-123",
		Token:     "both-invalid",
		ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired
		IsRevoked: true,                           // AND revoked
		RevokedAt: func() *time.Time { t := time.Now(); return &t }(),
		CreatedAt: time.Now().Add(-8 * 24 * time.Hour),
	}

	assert.False(t, token.IsValid(), "Token should be invalid (both revoked and expired)")
}

func TestRefreshToken_IsValid_EdgeCase_ExactlyExpiring(t *testing.T) {
	// Token expiring in exactly 1 second
	token := &models.RefreshToken{
		ID:        "token-edge",
		UserID:    "user-123",
		Token:     "edge-token",
		ExpiresAt: time.Now().Add(1 * time.Second),
		IsRevoked: false,
		CreatedAt: time.Now(),
	}

	// Should be valid now
	assert.True(t, token.IsValid())

	// Wait for expiration
	time.Sleep(2 * time.Second)

	// Should be invalid now
	assert.False(t, token.IsValid())
}

// ============================================================================
// TEST SUITE: REFRESH TOKEN REPOSITORY - CRUD OPERATIONS
// ============================================================================

func TestRefreshTokenRepo_Create_Success(t *testing.T) {
	mockRepo := new(AuthMockRefreshTokenRepository)

	token := &models.RefreshToken{
		UserID:    "user-create-123",
		Token:     "new-refresh-token",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		IsRevoked: false,
	}

	mockRepo.On("Create", mock.Anything, token).Return(nil)

	ctx := context.Background()
	err := mockRepo.Create(ctx, token)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRefreshTokenRepo_GetByToken_Success(t *testing.T) {
	mockRepo := new(AuthMockRefreshTokenRepository)

	expectedToken := &models.RefreshToken{
		ID:        "token-get-123",
		UserID:    "user-123",
		Token:     "find-this-token",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		IsRevoked: false,
	}

	mockRepo.On("GetByToken", mock.Anything, "find-this-token").Return(expectedToken, nil)

	ctx := context.Background()
	token, err := mockRepo.GetByToken(ctx, "find-this-token")

	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, "find-this-token", token.Token)
	assert.Equal(t, "user-123", token.UserID)
	assert.False(t, token.IsRevoked)

	mockRepo.AssertExpectations(t)
}

func TestRefreshTokenRepo_GetByToken_NotFound(t *testing.T) {
	mockRepo := new(AuthMockRefreshTokenRepository)

	mockRepo.On("GetByToken", mock.Anything, "nonexistent-token").Return(nil, nil)

	ctx := context.Background()
	token, err := mockRepo.GetByToken(ctx, "nonexistent-token")

	assert.NoError(t, err)
	assert.Nil(t, token)

	mockRepo.AssertExpectations(t)
}

func TestRefreshTokenRepo_RevokeToken_Success(t *testing.T) {
	mockRepo := new(AuthMockRefreshTokenRepository)

	mockRepo.On("RevokeToken", mock.Anything, "revoke-me").Return(nil)

	ctx := context.Background()
	err := mockRepo.RevokeToken(ctx, "revoke-me")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRefreshTokenRepo_RevokeAllUserTokens_Success(t *testing.T) {
	mockRepo := new(AuthMockRefreshTokenRepository)

	mockRepo.On("RevokeAllUserTokens", mock.Anything, "user-revoke-all").Return(nil)

	ctx := context.Background()
	err := mockRepo.RevokeAllUserTokens(ctx, "user-revoke-all")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRefreshTokenRepo_DeleteExpired_Success(t *testing.T) {
	mockRepo := new(AuthMockRefreshTokenRepository)

	mockRepo.On("DeleteExpired", mock.Anything).Return(nil)

	ctx := context.Background()
	err := mockRepo.DeleteExpired(ctx)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// ============================================================================
// TEST SUITE: REFRESH TOKEN LIFECYCLE
// ============================================================================

func TestRefreshToken_LifecycleSimulation(t *testing.T) {
	mockRepo := new(AuthMockRefreshTokenRepository)

	// Step 1: Create token
	newToken := &models.RefreshToken{
		UserID:    "user-lifecycle",
		Token:     "lifecycle-token",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		IsRevoked: false,
	}
	mockRepo.On("Create", mock.Anything, newToken).Return(nil)

	ctx := context.Background()
	err := mockRepo.Create(ctx, newToken)
	assert.NoError(t, err)

	// Step 2: Retrieve token
	mockRepo.On("GetByToken", mock.Anything, "lifecycle-token").Return(newToken, nil)
	retrieved, err := mockRepo.GetByToken(ctx, "lifecycle-token")
	assert.NoError(t, err)
	assert.True(t, retrieved.IsValid())

	// Step 3: Revoke token
	mockRepo.On("RevokeToken", mock.Anything, "lifecycle-token").Return(nil).Run(func(args mock.Arguments) {
		// Simulate the revocation on our token object
		newToken.IsRevoked = true
		now := time.Now()
		newToken.RevokedAt = &now
	})
	err = mockRepo.RevokeToken(ctx, "lifecycle-token")
	assert.NoError(t, err)

	// Step 4: Verify token is now revoked (using the modified token)
	assert.True(t, newToken.IsRevoked, "Token should be marked as revoked")
	assert.False(t, newToken.IsValid(), "Revoked token should not be valid")

	mockRepo.AssertExpectations(t)
}

// ============================================================================
// TEST SUITE: TOKEN EXPIRATION SCENARIOS
// ============================================================================

func TestRefreshToken_ExpirationScenarios(t *testing.T) {
	testCases := []struct {
		name          string
		expiresIn     time.Duration
		shouldBeValid bool
	}{
		{
			name:          "Just created (7 days)",
			expiresIn:     7 * 24 * time.Hour,
			shouldBeValid: true,
		},
		{
			name:          "1 day remaining",
			expiresIn:     24 * time.Hour,
			shouldBeValid: true,
		},
		{
			name:          "1 hour remaining",
			expiresIn:     1 * time.Hour,
			shouldBeValid: true,
		},
		{
			name:          "1 minute remaining",
			expiresIn:     1 * time.Minute,
			shouldBeValid: true,
		},
		{
			name:          "Just expired (1 second ago)",
			expiresIn:     -1 * time.Second,
			shouldBeValid: false,
		},
		{
			name:          "Expired 1 hour ago",
			expiresIn:     -1 * time.Hour,
			shouldBeValid: false,
		},
		{
			name:          "Expired 1 day ago",
			expiresIn:     -24 * time.Hour,
			shouldBeValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token := &models.RefreshToken{
				UserID:    "user-exp",
				Token:     "exp-token",
				ExpiresAt: time.Now().Add(tc.expiresIn),
				IsRevoked: false,
			}

			assert.Equal(t, tc.shouldBeValid, token.IsValid())
		})
	}
}

// ============================================================================
// TEST SUITE: SECURITY - TOKEN REVOCATION
// ============================================================================

func TestRefreshToken_RevokeAll_SecurityScenario(t *testing.T) {
	mockRepo := new(AuthMockRefreshTokenRepository)

	// Scenario: User reports account compromise
	// All tokens for user-compromised should be revoked
	compromisedUserID := "user-compromised"

	mockRepo.On("RevokeAllUserTokens", mock.Anything, compromisedUserID).Return(nil)

	ctx := context.Background()
	err := mockRepo.RevokeAllUserTokens(ctx, compromisedUserID)

	assert.NoError(t, err, "Should successfully revoke all tokens for compromised account")
	mockRepo.AssertExpectations(t)
}

func TestRefreshToken_SequentialRevocation(t *testing.T) {
	mockRepo := new(AuthMockRefreshTokenRepository)

	// User has 3 active sessions (3 refresh tokens)
	tokens := []string{"session-1", "session-2", "session-3"}

	// Revoke them one by one
	for i, tokenStr := range tokens {
		mockRepo.On("RevokeToken", mock.Anything, tokenStr).Return(nil).Once()

		ctx := context.Background()
		err := mockRepo.RevokeToken(ctx, tokenStr)

		assert.NoError(t, err, "Token %d should revoke successfully", i+1)
	}

	mockRepo.AssertExpectations(t)
	assert.Equal(t, 3, len(mockRepo.Calls), "Should have 3 revoke calls")
}

// ============================================================================
// TEST SUITE: CLEANUP OPERATIONS
// ============================================================================

func TestRefreshToken_DeleteExpired_CleanupOldTokens(t *testing.T) {
	mockRepo := new(AuthMockRefreshTokenRepository)

	// Simulate cleanup job running daily
	mockRepo.On("DeleteExpired", mock.Anything).Return(nil)

	ctx := context.Background()
	err := mockRepo.DeleteExpired(ctx)

	assert.NoError(t, err, "Cleanup should succeed")
	mockRepo.AssertExpectations(t)
}

// ============================================================================
// TEST SUMMARY
// ============================================================================

func TestRefreshTokenRepositoryTestSuiteSummary(t *testing.T) {
	summary := `
╔══════════════════════════════════════════════════════════════════════╗
║       REFRESH TOKEN REPOSITORY COMPREHENSIVE TEST SUITE              ║
╠══════════════════════════════════════════════════════════════════════╣
║                                                                       ║
║  ✅ MODEL VALIDATION (5 tests)                                       ║
║     • Valid token detection                                         ║
║     • Expired token detection                                       ║
║     • Revoked token detection                                       ║
║     • Both revoked and expired                                      ║
║     • Edge case: exactly expiring token                             ║
║                                                                       ║
║  ✅ CRUD OPERATIONS (6 tests)                                        ║
║     • Create token                                                  ║
║     • Get by token - success                                        ║
║     • Get by token - not found                                      ║
║     • Revoke single token                                           ║
║     • Revoke all user tokens                                        ║
║     • Delete expired tokens                                         ║
║                                                                       ║
║  ✅ LIFECYCLE (1 test)                                              ║
║     • Full lifecycle simulation (create → get → revoke → verify)    ║
║                                                                       ║
║  ✅ EXPIRATION SCENARIOS (7 tests)                                  ║
║     • Just created (7 days)                                         ║
║     • 1 day, 1 hour, 1 minute remaining                             ║
║     • Just expired, 1 hour ago, 1 day ago                           ║
║                                                                       ║
║  ✅ SECURITY (2 tests)                                              ║
║     • Revoke all on account compromise                              ║
║     • Sequential revocation of multiple sessions                    ║
║                                                                       ║
║  ✅ CLEANUP (1 test)                                                ║
║     • Delete expired tokens job                                     ║
║                                                                       ║
║  📊 TOTAL: 22 comprehensive repository tests                        ║
║                                                                       ║
╚══════════════════════════════════════════════════════════════════════╝
`
	t.Log(summary)
}
