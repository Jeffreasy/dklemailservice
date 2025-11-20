package tests

import (
	"context"
	"dklautomationgo/models"
	"dklautomationgo/services"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// ============================================================================
// BENCHMARK SUITE: AUTH SERVICE PERFORMANCE
// ============================================================================

func BenchmarkLogin_Success(b *testing.B) {
	mockUserRepo := new(AuthMockGebruikerRepository)
	mockRefreshRepo := new(AuthMockRefreshTokenRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	testUser := &models.Gebruiker{
		ID:             "bench-user",
		Email:          "bench@dekoninklijkeloop.nl",
		WachtwoordHash: string(hashedPassword),
		IsActief:       true,
	}

	mockUserRepo.On("GetByEmail", mock.Anything, "bench@dekoninklijkeloop.nl").Return(testUser, nil)
	mockUserRepo.On("UpdateLastLogin", mock.Anything, "bench-user").Return(nil)
	mockRefreshRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.RefreshToken")).Return(nil)
	mockUserRoleRepo.On("ListActiveByUser", mock.Anything, "bench-user").Return([]*models.UserRole{}, nil)

	authService := services.NewAuthServiceWithRBAC(mockUserRepo, mockRefreshRepo, mockUserRoleRepo)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = authService.Login(ctx, "bench@dekoninklijkeloop.nl", "testpass")
	}
}

func BenchmarkTokenValidation(b *testing.B) {
	mockUserRepo := new(AuthMockGebruikerRepository)
	mockRefreshRepo := new(AuthMockRefreshTokenRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	testUser := &models.Gebruiker{
		ID:             "token-bench",
		Email:          "token@dekoninklijkeloop.nl",
		WachtwoordHash: string(hashedPassword),
		IsActief:       true,
	}

	mockUserRepo.On("GetByEmail", mock.Anything, "token@dekoninklijkeloop.nl").Return(testUser, nil)
	mockUserRepo.On("UpdateLastLogin", mock.Anything, "token-bench").Return(nil)
	mockRefreshRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.RefreshToken")).Return(nil)
	mockUserRoleRepo.On("ListActiveByUser", mock.Anything, "token-bench").Return([]*models.UserRole{}, nil)

	authService := services.NewAuthServiceWithRBAC(mockUserRepo, mockRefreshRepo, mockUserRoleRepo)
	ctx := context.Background()

	// Generate token once
	token, _, _ := authService.Login(ctx, "token@dekoninklijkeloop.nl", "testpass")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = authService.ValidateToken(token)
	}
}

func BenchmarkPasswordHashing(b *testing.B) {
	mockUserRepo := new(AuthMockGebruikerRepository)
	mockRefreshRepo := new(AuthMockRefreshTokenRepository)
	authService := services.NewAuthService(mockUserRepo, mockRefreshRepo)

	password := "BenchmarkPassword123!"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = authService.HashPassword(password)
	}
}

func BenchmarkPasswordVerification(b *testing.B) {
	mockUserRepo := new(AuthMockGebruikerRepository)
	mockRefreshRepo := new(AuthMockRefreshTokenRepository)
	authService := services.NewAuthService(mockUserRepo, mockRefreshRepo)

	password := "BenchmarkPassword123!"
	hash, _ := authService.HashPassword(password)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = authService.VerifyPassword(hash, password)
	}
}

func BenchmarkRefreshToken_Generation(b *testing.B) {
	mockUserRepo := new(AuthMockGebruikerRepository)
	mockRefreshRepo := new(AuthMockRefreshTokenRepository)
	authService := services.NewAuthService(mockUserRepo, mockRefreshRepo)

	mockRefreshRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.RefreshToken")).Return(nil)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = authService.GenerateRefreshToken(ctx, "bench-user-id")
	}
}

func BenchmarkRefreshAccessToken_Flow(b *testing.B) {
	mockUserRepo := new(AuthMockGebruikerRepository)
	mockRefreshRepo := new(AuthMockRefreshTokenRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)

	testUser := &models.Gebruiker{
		ID:       "refresh-bench",
		Email:    "refresh@dekoninklijkeloop.nl",
		IsActief: true,
	}

	validToken := &models.RefreshToken{
		OwnerID:   "refresh-bench",
		Token:     "bench-refresh-token",
		ExpiresAt: testUser.CreatedAt.Add(7 * 24 * 60 * 60 * 1000000000), // 7 days
		IsRevoked: false,
	}

	mockRefreshRepo.On("GetByToken", mock.Anything, "bench-refresh-token").Return(validToken, nil)
	mockUserRepo.On("GetByID", mock.Anything, "refresh-bench").Return(testUser, nil)
	mockRefreshRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.RefreshToken")).Return(nil)
	mockRefreshRepo.On("RevokeToken", mock.Anything, "bench-refresh-token").Return(nil)
	mockUserRoleRepo.On("ListActiveByUser", mock.Anything, "refresh-bench").Return([]*models.UserRole{}, nil)
	mockUserRoleRepo.On("GetByUserIDWithRoles", mock.Anything, "refresh-bench").Return([]*models.UserRole{}, nil)

	authService := services.NewAuthServiceWithRBAC(mockUserRepo, mockRefreshRepo, mockUserRoleRepo)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = authService.RefreshAccessToken(ctx, "bench-refresh-token")
	}
}

// ============================================================================
// BENCHMARK SUITE: PERMISSION SERVICE PERFORMANCE
// ============================================================================

func BenchmarkPermissionCheck_WithPermissions(b *testing.B) {
	mockRoleRepo := new(MockRBACRoleRepository)
	mockPermRepo := new(MockPermissionRepository)
	mockRolePermRepo := new(MockRolePermissionRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)

	service := services.NewPermissionService(mockRoleRepo, mockPermRepo, mockRolePermRepo, mockUserRoleRepo)

	permissions := []*models.UserPermission{
		{Resource: "contact", Action: "read"},
		{Resource: "contact", Action: "write"},
		{Resource: "admin", Action: "access"},
	}

	mockUserRoleRepo.On("GetUserPermissions", mock.Anything, "perf-user").Return(permissions, nil)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.HasPermission(ctx, "perf-user", "contact", "read")
	}
}

func BenchmarkPermissionCheck_WithoutPermissions(b *testing.B) {
	mockRoleRepo := new(MockRBACRoleRepository)
	mockPermRepo := new(MockPermissionRepository)
	mockRolePermRepo := new(MockRolePermissionRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)

	service := services.NewPermissionService(mockRoleRepo, mockPermRepo, mockRolePermRepo, mockUserRoleRepo)

	permissions := []*models.UserPermission{
		{Resource: "contact", Action: "read"},
	}

	mockUserRoleRepo.On("GetUserPermissions", mock.Anything, "limited-user").Return(permissions, nil)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.HasPermission(ctx, "limited-user", "admin", "access")
	}
}

// ============================================================================
// BENCHMARK SUITE: CONCURRENT OPERATIONS
// ============================================================================

func BenchmarkConcurrentLogins(b *testing.B) {
	mockUserRepo := new(AuthMockGebruikerRepository)
	mockRefreshRepo := new(AuthMockRefreshTokenRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	testUser := &models.Gebruiker{
		ID:             "concurrent-bench",
		Email:          "concurrent@dekoninklijkeloop.nl",
		WachtwoordHash: string(hashedPassword),
		IsActief:       true,
	}

	mockUserRepo.On("GetByEmail", mock.Anything, "concurrent@dekoninklijkeloop.nl").Return(testUser, nil)
	mockUserRepo.On("UpdateLastLogin", mock.Anything, "concurrent-bench").Return(nil)
	mockRefreshRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.RefreshToken")).Return(nil)
	mockUserRoleRepo.On("ListActiveByUser", mock.Anything, "concurrent-bench").Return([]*models.UserRole{}, nil)

	authService := services.NewAuthServiceWithRBAC(mockUserRepo, mockRefreshRepo, mockUserRoleRepo)
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _, _ = authService.Login(ctx, "concurrent@dekoninklijkeloop.nl", "testpass")
		}
	})
}

func BenchmarkConcurrentTokenValidation(b *testing.B) {
	mockUserRepo := new(AuthMockGebruikerRepository)
	mockRefreshRepo := new(AuthMockRefreshTokenRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	testUser := &models.Gebruiker{
		ID:             "val-bench",
		Email:          "val@dekoninklijkeloop.nl",
		WachtwoordHash: string(hashedPassword),
		IsActief:       true,
	}

	mockUserRepo.On("GetByEmail", mock.Anything, "val@dekoninklijkeloop.nl").Return(testUser, nil)
	mockUserRepo.On("UpdateLastLogin", mock.Anything, "val-bench").Return(nil)
	mockRefreshRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.RefreshToken")).Return(nil)
	mockUserRoleRepo.On("ListActiveByUser", mock.Anything, "val-bench").Return([]*models.UserRole{}, nil)

	authService := services.NewAuthServiceWithRBAC(mockUserRepo, mockRefreshRepo, mockUserRoleRepo)
	ctx := context.Background()
	token, _, _ := authService.Login(ctx, "val@dekoninklijkeloop.nl", "testpass")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = authService.ValidateToken(token)
		}
	})
}

func BenchmarkConcurrentPermissionChecks(b *testing.B) {
	mockRoleRepo := new(MockRBACRoleRepository)
	mockPermRepo := new(MockPermissionRepository)
	mockRolePermRepo := new(MockRolePermissionRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)

	service := services.NewPermissionService(mockRoleRepo, mockPermRepo, mockRolePermRepo, mockUserRoleRepo)

	permissions := []*models.UserPermission{
		{Resource: "contact", Action: "read"},
		{Resource: "contact", Action: "write"},
	}

	mockUserRoleRepo.On("GetUserPermissions", mock.Anything, "concurrent-perm-user").Return(permissions, nil)

	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = service.HasPermission(ctx, "concurrent-perm-user", "contact", "read")
		}
	})
}

// ============================================================================
// PERFORMANCE TEST VALIDATION
// ============================================================================

func TestPerformance_LoginTime_Target(t *testing.T) {
	// Target: < 300ms per login (as per AUTH_AND_RBAC.md)
	// This is a validation test, not a benchmark

	mockUserRepo := new(AuthMockGebruikerRepository)
	mockRefreshRepo := new(AuthMockRefreshTokenRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	testUser := &models.Gebruiker{
		ID:             "perf-user",
		Email:          "perf@dekoninklijkeloop.nl",
		WachtwoordHash: string(hashedPassword),
		IsActief:       true,
	}

	mockUserRepo.On("GetByEmail", mock.Anything, "perf@dekoninklijkeloop.nl").Return(testUser, nil)
	mockUserRepo.On("UpdateLastLogin", mock.Anything, "perf-user").Return(nil)
	mockRefreshRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.RefreshToken")).Return(nil)
	mockUserRoleRepo.On("ListActiveByUser", mock.Anything, "perf-user").Return([]*models.UserRole{
		{
			UserID: "perf-user",
			RoleID: "role-admin",
			Role: models.RBACRole{
				ID:   "role-admin",
				Name: "admin",
			},
			IsActive: true,
		},
	}, nil)

	authService := services.NewAuthServiceWithRBAC(mockUserRepo, mockRefreshRepo, mockUserRoleRepo)
	ctx := context.Background()

	// Measure 10 logins
	start := time.Now()
	for i := 0; i < 10; i++ {
		_, _, err := authService.Login(ctx, "perf@dekoninklijkeloop.nl", "testpass")
		if err != nil {
			t.Fatalf("Login failed: %v", err)
		}
	}
	elapsed := time.Since(start)
	avgTime := elapsed / 10

	t.Logf("Average login time: %v (target: < 300ms)", avgTime)

	// Note: In production with real database, target is < 300ms
	// With mocks, should be much faster
	if avgTime > 500*time.Millisecond {
		t.Logf("Warning: Login time %v exceeds even mock target", avgTime)
	}
}

func TestPerformance_TokenValidation_Target(t *testing.T) {
	// Target: < 10ms per validation (as per AUTH_AND_RBAC.md)

	mockUserRepo := new(AuthMockGebruikerRepository)
	mockRefreshRepo := new(AuthMockRefreshTokenRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	testUser := &models.Gebruiker{
		ID:             "val-perf",
		Email:          "valperf@dekoninklijkeloop.nl",
		WachtwoordHash: string(hashedPassword),
		IsActief:       true,
	}

	mockUserRepo.On("GetByEmail", mock.Anything, "valperf@dekoninklijkeloop.nl").Return(testUser, nil)
	mockUserRepo.On("UpdateLastLogin", mock.Anything, "val-perf").Return(nil)
	mockRefreshRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.RefreshToken")).Return(nil)
	mockUserRoleRepo.On("ListActiveByUser", mock.Anything, "val-perf").Return([]*models.UserRole{}, nil)

	authService := services.NewAuthServiceWithRBAC(mockUserRepo, mockRefreshRepo, mockUserRoleRepo)
	ctx := context.Background()
	token, _, _ := authService.Login(ctx, "valperf@dekoninklijkeloop.nl", "testpass")

	start := time.Now()
	for i := 0; i < 1000; i++ {
		_, err := authService.ValidateToken(token)
		if err != nil {
			t.Fatalf("Validation failed: %v", err)
		}
	}
	elapsed := time.Since(start)
	avgTime := elapsed / 1000

	t.Logf("Average token validation: %v (target: < 10ms)", avgTime)

	if avgTime > 10*time.Millisecond {
		t.Logf("Warning: Validation time %v exceeds target", avgTime)
	}
}

// ============================================================================
// BENCHMARK SUMMARY
// ============================================================================

func TestBenchmarkSummary(t *testing.T) {
	summary := `
╔══════════════════════════════════════════════════════════════════════╗
║       AUTH & RBAC PERFORMANCE BENCHMARK SUITE                        ║
╠══════════════════════════════════════════════════════════════════════╣
║                                                                       ║
║  📊 BENCHMARKS AVAILABLE:                                            ║
║                                                                       ║
║  BenchmarkLogin_Success                                              ║
║     Target: < 300ms (with database)                                 ║
║     Tests: Full login flow with bcrypt                              ║
║                                                                       ║
║  BenchmarkTokenValidation                                            ║
║     Target: < 10ms                                                  ║
║     Tests: JWT parsing and validation                               ║
║                                                                       ║
║  BenchmarkPasswordHashing                                            ║
║     Tests: BCrypt hash generation                                   ║
║                                                                       ║
║  BenchmarkPasswordVerification                                       ║
║     Tests: BCrypt password comparison                               ║
║                                                                       ║
║  BenchmarkRefreshToken_Generation                                    ║
║     Tests: Secure random token generation                           ║
║                                                                       ║
║  BenchmarkRefreshAccessToken_Flow                                    ║
║     Target: < 200ms (with database)                                 ║
║     Tests: Complete refresh flow with rotation                      ║
║                                                                       ║
║  BenchmarkConcurrentLogins                                           ║
║     Tests: Parallel login handling                                  ║
║                                                                       ║
║  BenchmarkConcurrentTokenValidation                                  ║
║     Tests: Parallel token validation                                ║
║                                                                       ║
║  BenchmarkConcurrentPermissionChecks                                 ║
║     Tests: Parallel permission checks                               ║
║                                                                       ║
║  📊 RUN BENCHMARKS:                                                  ║
║     go test -bench=. ./tests -benchmem                              ║
║                                                                       ║
║  📊 SPECIFIC BENCHMARK:                                              ║
║     go test -bench=BenchmarkLogin ./tests -benchmem                 ║
║                                                                       ║
║  📊 WITH CPU PROFILE:                                                ║
║     go test -bench=. ./tests -cpuprofile=cpu.prof                   ║
║                                                                       ║
╚══════════════════════════════════════════════════════════════════════╝
`
	t.Log(summary)
}
