# Testing Guide

Complete test handleiding voor de DKL Email Service.

## Test Strategie

### Test Pyramid

```
        /\
       /  \      E2E Tests (5%)
      /____\
     /      \    Integration Tests (15%)
    /________\
   /          \  Unit Tests (80%)
  /____________\
```

**Focus:**
- **80% Unit Tests** - Snelle, geïsoleerde tests
- **15% Integration Tests** - Database, SMTP, Redis
- **5% E2E Tests** - Complete flows

## Test Setup

### Dependencies

```bash
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/mock
go get github.com/stretchr/testify/require
```

### Test Directory

**Structuur:**
```
tests/
├── aanmelding_handler_test.go
├── contact_handler_test.go
├── email_service_test.go
├── auth_service_test.go
├── rate_limiter_test.go
├── mail_handler_test.go
├── elk_integration_test.go
├── mocks.go
└── test_helpers.go
```

## Unit Tests

### Handler Tests

**Voorbeeld:** [`tests/aanmelding_handler_test.go`](../../tests/aanmelding_handler_test.go:1)

```go
func TestAanmeldingHandler(t *testing.T) {
    tests := []struct {
        name       string
        input      models.AanmeldingRequest
        mockSetup  func(*MockEmailService)
        wantStatus int
        wantErr    bool
    }{
        {
            name: "valid request",
            input: models.AanmeldingRequest{
                Naam:    "Test User",
                Email:   "test@example.com",
                Telefoon: "0612345678",
                Rol:     "loper",
                Afstand: "5km",
                Terms:   true,
            },
            mockSetup: func(m *MockEmailService) {
                m.On("SendAanmeldingEmail", mock.Anything).Return(nil)
            },
            wantStatus: http.StatusOK,
            wantErr:    false,
        },
        {
            name: "missing naam",
            input: models.AanmeldingRequest{
                Email: "test@example.com",
                Terms: true,
            },
            wantStatus: http.StatusBadRequest,
            wantErr:    true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup
            mockEmail := &MockEmailService{}
            if tt.mockSetup != nil {
                tt.mockSetup(mockEmail)
            }
            
            handler := handlers.NewEmailHandler(mockEmail, nil, nil)
            
            // Create request
            body, _ := json.Marshal(tt.input)
            req := httptest.NewRequest("POST", "/api/aanmelding-email", bytes.NewReader(body))
            req.Header.Set("Content-Type", "application/json")
            
            // Execute
            resp, err := app.Test(req)
            
            // Assert
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
            assert.Equal(t, tt.wantStatus, resp.StatusCode)
            
            if tt.mockSetup != nil {
                mockEmail.AssertExpectations(t)
            }
        })
    }
}
```

### Service Tests

**Voorbeeld:** [`tests/email_service_test.go:12`](../../tests/email_service_test.go:12)

```go
func TestEmailService_SendEmail(t *testing.T) {
    // Setup
    reg := NewTestRegistry()
    mockSMTP := &mockSMTP{}
    emailMetrics := services.NewEmailMetrics(time.Hour)
    prometheusMetrics := services.NewPrometheusMetricsWithRegistry(reg)
    rateLimiter := services.NewRateLimiter(prometheusMetrics)
    
    emailService := services.NewEmailService(
        mockSMTP, 
        emailMetrics, 
        rateLimiter, 
        prometheusMetrics,
    )
    
    t.Run("Succesvolle verzending", func(t *testing.T) {
        mockSMTP.On("Send", mock.Anything).Return(nil)
        
        err := emailService.SendEmail("test@example.com", "Test", "Body")
        
        assert.NoError(t, err)
        mockSMTP.mutex.Lock()
        assert.True(t, mockSMTP.SendCalled)
        assert.Equal(t, "test@example.com", mockSMTP.LastTo)
        mockSMTP.mutex.Unlock()
    })
    
    t.Run("SMTP fout", func(t *testing.T) {
        mockSMTP.SetShouldFail(true)
        
        err := emailService.SendEmail("test@example.com", "Test", "Body")
        
        assert.Error(t, err)
        mockSMTP.SetShouldFail(false)
    })
    
    t.Run("Rate limit overschreden", func(t *testing.T) {
        rateLimiter.AddLimit("email_generic", 1, time.Hour, false)
        
        // Eerste email moet lukken
        err := emailService.SendEmail("test@example.com", "Test", "Body")
        assert.NoError(t, err)
        
        // Tweede email moet falen
        err = emailService.SendEmail("test@example.com", "Test", "Body")
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "rate limit")
    })
}
```

## Mock Implementations

### SMTP Mock

**Implementatie:** [`tests/mocks.go:108`](../../tests/mocks.go:108)

```go
type mockSMTP struct {
    mock.Mock
    mutex              sync.Mutex
    SendCalled         bool
    SendRegCalled      bool
    SendWFCCalled      bool
    SendWithFromCalled bool
    LastFrom           string
    LastTo             string
    LastSubject        string
    LastBody           string
    SendError          error
}

func (m *mockSMTP) Send(msg *services.EmailMessage) error {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    m.SendCalled = true
    m.LastTo = msg.To
    m.LastSubject = msg.Subject
    m.LastBody = msg.Body
    if m.SendError != nil {
        return m.SendError
    }
    return nil
}

func (m *mockSMTP) SetShouldFail(shouldFail bool) {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    if shouldFail {
        m.SendError = fmt.Errorf("mock SMTP error")
    } else {
        m.SendError = nil
    }
}
```

### Auth Service Mock

**Implementatie:** [`tests/mocks.go:561`](../../tests/mocks.go:561)

```go
type MockAuthService struct {
    mock.Mock
}

func (m *MockAuthService) Login(ctx context.Context, email, wachtwoord string) (string, string, error) {
    args := m.Called(ctx, email, wachtwoord)
    return args.String(0), args.String(1), args.Error(2)
}

func (m *MockAuthService) ValidateToken(token string) (string, error) {
    args := m.Called(token)
    return args.String(0), args.Error(1)
}
```

### Notification Service Mock

**Implementatie:** [`tests/mocks.go:419`](../../tests/mocks.go:419)

```go
type MockNotificationService struct {
    mock.Mock
}

func (m *MockNotificationService) CreateNotification(
    ctx context.Context,
    notificationType models.NotificationType,
    priority models.NotificationPriority,
    title, message string,
) (*models.Notification, error) {
    args := m.Called(ctx, notificationType, priority, title, message)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.Notification), args.Error(1)
}
```

## Integration Tests

### Database Integration

```go
func TestDatabaseIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    // Setup test database
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    // Test repository
    repo := repository.NewGebruikerRepository(db)
    
    t.Run("Create user", func(t *testing.T) {
        gebruiker := &models.Gebruiker{
            Naam:  "Test User",
            Email: "test@example.com",
            Rol:   "gebruiker",
        }
        
        err := repo.Create(context.Background(), gebruiker)
        assert.NoError(t, err)
        assert.NotEmpty(t, gebruiker.ID)
    })
    
    t.Run("Get user by email", func(t *testing.T) {
        gebruiker, err := repo.GetByEmail(context.Background(), "test@example.com")
        assert.NoError(t, err)
        assert.NotNil(t, gebruiker)
        assert.Equal(t, "Test User", gebruiker.Naam)
    })
}
```

### ELK Integration

**Implementatie:** [`tests/elk_integration_test.go`](../../tests/elk_integration_test.go:1)

```go
func TestELKIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    elkWriter := logger.NewELKWriter(logger.ELKConfig{
        Endpoint: "http://localhost:9200",
        Index:    "test-logs",
    })
    
    err := elkWriter.Write([]byte(`{"level":"info","message":"test"}`))
    require.NoError(t, err)
    
    // Verify log entry in Elasticsearch
    // ... verificatie logica
}
```

## Test Coverage

### Run Tests

**Alle Tests:**
```bash
go test ./...
```

**Met Verbose Output:**
```bash
go test -v ./...
```

**Specifiek Package:**
```bash
go test -v ./handlers
go test -v ./services
```

**Specifieke Test:**
```bash
go test -v ./tests -run TestEmailService_SendEmail
```

### Coverage Report

**Generate Coverage:**
```bash
go test -coverprofile=coverage.out ./...
```

**HTML Report:**
```bash
go tool cover -html=coverage.out -o coverage.html
```

**Coverage per Function:**
```bash
go tool cover -func=coverage.out
```

**Coverage per Package:**
```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep -E '^total:|^dklautomationgo'
```

### Coverage Goals

| Package | Target | Current |
|---------|--------|---------|
| handlers | 90% | - |
| services | 95% | - |
| models | 100% | - |
| repository | 90% | - |
| middleware | 95% | - |

## API Testing

### PowerShell Scripts

**Locatie:** [`testscripts/`](../../testscripts/)

**Beschikbare Scripts:**
- `test_api_light.ps1` - Basis API tests
- `test_api_full.ps1` - Complete API tests met stress test
- `stress_test.ps1` - Alleen stress test
- `check_mail_logs.ps1` - Email log analyse

### Test API Light

**Gebruik:**
```powershell
# Basis test
.\testscripts\test_api_light.ps1

# Met gedetailleerde health check
.\testscripts\test_api_light.ps1 -DetailedHealth

# Test mail endpoints
.\testscripts\test_api_light.ps1 -TestMailEndpoints

# Test beveiligde endpoints
.\testscripts\test_api_light.ps1 -TestMailEndpoints -IncludeSecuredEndpoints -JWTToken "your-token"

# Test productie API
.\testscripts\test_api_light.ps1 -BaseUrl "https://api.dekoninklijkeloop.nl"
```

**Parameters:**
- `-BaseUrl` - API base URL (default: http://localhost:8080)
- `-DetailedHealth` - Gedetailleerde health check
- `-TestMailEndpoints` - Test email endpoints
- `-IncludeSecuredEndpoints` - Test beveiligde endpoints
- `-JWTToken` - JWT token voor authenticatie
- `-AdminAPIKey` - Admin API key voor metrics

### Manual API Testing

**cURL Examples:**

```bash
# Health check
curl http://localhost:8080/api/health

# Contact email (test mode)
curl -X POST http://localhost:8080/api/contact-email \
  -H "Content-Type: application/json" \
  -H "X-Test-Mode: true" \
  -d '{
    "naam": "Test User",
    "email": "test@example.com",
    "bericht": "Test bericht",
    "privacy_akkoord": true
  }'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@dekoninklijkeloop.nl",
    "wachtwoord": "your-password"
  }'

# Get profile (met token)
TOKEN="your-jwt-token"
curl http://localhost:8080/api/auth/profile \
  -H "Authorization: Bearer $TOKEN"
```

## Performance Testing

### Benchmark Tests

**Email Service Benchmark:**
```go
func BenchmarkEmailService_SendBatch(b *testing.B) {
    mockSMTP := &mockSMTP{}
    mockSMTP.On("Send", mock.Anything).Return(nil)
    
    service := setupEmailService(mockSMTP)
    emails := generateTestEmails(100)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        for _, email := range emails {
            service.SendEmail(email.To, email.Subject, email.Body)
        }
    }
}
```

**Run Benchmarks:**
```bash
# Alle benchmarks
go test -bench=. ./...

# Met memory profiling
go test -bench=. -benchmem ./...

# Specifieke benchmark
go test -bench=BenchmarkEmailService ./services
```

### Load Testing

**Apache Bench:**
```bash
# 1000 requests, 10 concurrent
ab -n 1000 -c 10 -H "Content-Type: application/json" \
  -p test_data.json \
  http://localhost:8080/api/contact-email
```

**Vegeta:**
```bash
# Install
go install github.com/tsenart/vegeta@latest

# Load test
echo "POST http://localhost:8080/api/contact-email" | \
  vegeta attack -duration=30s -rate=50 | \
  vegeta report
```

## Test Mode

### Activatie

**Via Header:**
```http
X-Test-Mode: true
```

**Via Body:**
```json
{
    "test_mode": true,
    "naam": "Test User",
    "email": "test@example.com"
}
```

**Implementatie:** [`handlers/middleware.go:174`](../../handlers/middleware.go:174)

```go
func TestModeMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Check header
        if testMode := c.Get("X-Test-Mode"); testMode == "true" {
            c.Locals("test_mode", true)
        }
        
        // Check query parameter
        if testMode := c.Query("test_mode"); testMode == "true" {
            c.Locals("test_mode", true)
        }
        
        return c.Next()
    }
}
```

### Test Mode Gedrag

**Email Handler:** [`handlers/email_handler.go:128`](../../handlers/email_handler.go:128)

```go
// In testmodus sturen we geen echte emails
if testMode {
    logger.Info("Test modus: Geen admin email verzonden", "admin_email", adminEmail)
} else {
    if err := h.emailService.SendContactEmail(adminEmailData); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Fout bij het verzenden van de email: " + err.Error(),
        })
    }
}
```

**Response:**
```json
{
    "success": true,
    "message": "[TEST MODE] Je bericht is verwerkt (geen echte email verzonden).",
    "test_mode": true
}
```

## Email Auto Fetcher Testing

### Unit Tests

```go
func TestEmailAutoFetcher_Start(t *testing.T) {
    mockFetcher := &MockMailFetcher{}
    mockRepo := &MockIncomingEmailRepository{}
    
    autoFetcher := services.NewEmailAutoFetcher(mockFetcher, mockRepo)
    
    // Start fetcher
    autoFetcher.Start()
    assert.True(t, autoFetcher.IsRunning())
    
    // Stop fetcher
    autoFetcher.Stop()
    assert.False(t, autoFetcher.IsRunning())
}
```

### Integration Tests

**Test met Echte IMAP:**
```bash
# Configureer test credentials
export TEST_EMAIL_HOST=imap.example.com
export TEST_EMAIL_PORT=993
export TEST_EMAIL_USER=test@example.com
export TEST_EMAIL_PASSWORD=password

# Run integration tests
go test -v ./integration -run TestEmailFetching
```

### Manual Testing

**1. Configureer Snelle Interval:**
```bash
EMAIL_FETCH_INTERVAL=1  # 1 minuut voor testing
```

**2. Start Applicatie:**
```bash
go run main.go | grep "EmailAutoFetcher"
```

**3. Monitor Logs:**
```
INFO EmailAutoFetcher: Starting email fetch operation
INFO EmailAutoFetcher: Fetching emails from info@dekoninklijkeloop.nl
INFO EmailAutoFetcher: Found 5 emails, 3 new emails saved
```

**4. Handmatige Fetch:**
```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@dekoninklijkeloop.nl","wachtwoord":"password"}' | jq -r '.token')

curl -X POST http://localhost:8080/api/mail/fetch \
  -H "Authorization: Bearer $TOKEN"
```

## Race Detection

### Run met Race Detector

```bash
# Tests
go test -race ./...

# Applicatie
go run -race main.go
```

### Common Race Conditions

**Voorbeeld:**
```go
// ❌ Race condition
type Counter struct {
    count int
}

func (c *Counter) Increment() {
    c.count++  // Not thread-safe
}

// ✅ Fixed met mutex
type Counter struct {
    mu    sync.Mutex
    count int
}

func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}
```

## Test Helpers

### Setup Functions

```go
func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)
    
    // Run migrations
    db.AutoMigrate(&models.Gebruiker{}, &models.ContactFormulier{})
    
    return db
}

func cleanupTestDB(t *testing.T, db *gorm.DB) {
    sqlDB, _ := db.DB()
    sqlDB.Close()
}

func generateTestEmails(count int) []services.EmailMessage {
    emails := make([]services.EmailMessage, count)
    for i := 0; i < count; i++ {
        emails[i] = services.EmailMessage{
            To:      fmt.Sprintf("test%d@example.com", i),
            Subject: fmt.Sprintf("Test %d", i),
            Body:    "Test body",
        }
    }
    return emails
}
```

### Test Data Generators

```go
func createTestGebruiker() *models.Gebruiker {
    return &models.Gebruiker{
        Naam:     "Test User",
        Email:    "test@example.com",
        Rol:      "gebruiker",
        IsActief: true,
    }
}

func createTestContact() *models.ContactFormulier {
    return &models.ContactFormulier{
        Naam:           "Test Contact",
        Email:          "contact@example.com",
        Bericht:        "Test bericht",
        PrivacyAkkoord: true,
        Status:         "nieuw",
    }
}
```

## Continuous Integration

### GitHub Actions

**.github/workflows/test.yml:**
```yaml
name: Tests

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:13
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: dklemailservice_test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432
      
      redis:
        image: redis:6
        options: >-
          --health-cmd "redis-cli ping"
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 6379:6379
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Install dependencies
      run: go mod download
    
    - name: Run tests
      run: go test -v -race -coverprofile=coverage.out ./...
      env:
        DB_HOST: localhost
        DB_PORT: 5432
        DB_USER: postgres
        DB_PASSWORD: postgres
        DB_NAME: dklemailservice_test
        REDIS_HOST: localhost
        REDIS_PORT: 6379
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        files: ./coverage.out
```

## Test Best Practices

### Table-Driven Tests

```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {"valid email", "test@example.com", false},
        {"no @", "testexample.com", true},
        {"no domain", "test@", true},
        {"empty", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateEmail(tt.email)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Subtests

```go
func TestEmailService(t *testing.T) {
    service := setupEmailService()
    
    t.Run("SendContactEmail", func(t *testing.T) {
        t.Run("to admin", func(t *testing.T) {
            // Test admin email
        })
        
        t.Run("to user", func(t *testing.T) {
            // Test user email
        })
    })
    
    t.Run("SendAanmeldingEmail", func(t *testing.T) {
        // Test aanmelding email
    })
}
```

### Test Cleanup

```go
func TestWithCleanup(t *testing.T) {
    resource := setupResource()
    
    t.Cleanup(func() {
        resource.Close()
    })
    
    // Test implementation
}
```

### Test Timeouts

```go
func TestLongRunning(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping long running test")
    }
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // Test implementation with context
}
```

## Debugging Tests

### Verbose Logging

```go
func TestWithLogging(t *testing.T) {
    if testing.Verbose() {
        t.Logf("Running test with config: %+v", config)
    }
    
    // Test implementation
    
    if testing.Verbose() {
        t.Logf("Test result: %+v", result)
    }
}
```

**Run met Verbose:**
```bash
go test -v ./...
```

### Test Failures

**Skip Tests:**
```bash
go test -v ./... -skip=TestIntegration
```

**Run Only Failed:**
```bash
go test -v ./... -run=TestFailed
```

## Test Documentation

### Test Comments

```go
// TestEmailService_SendEmail tests the email sending functionality
// It covers:
// - Successful email sending
// - SMTP errors
// - Rate limiting
// - Template rendering
func TestEmailService_SendEmail(t *testing.T) {
    // ...
}
```

### Test Coverage Report

**Generate Report:**
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**CI Integration:**
```bash
# Upload to Codecov
bash <(curl -s https://codecov.io/bash)
```

## Troubleshooting

### Common Issues

**1. Import Cycles:**
```
package X imports package Y
package Y imports package X
```
**Oplossing:** Gebruik interfaces, verplaats naar shared package

**2. Race Conditions:**
```bash
go test -race ./...
```
**Oplossing:** Gebruik mutexes, channels, of atomic operations

**3. Flaky Tests:**
```go
// ❌ Time-dependent test
time.Sleep(100 * time.Millisecond)
assert.True(t, condition)

// ✅ Use Eventually
assert.Eventually(t, func() bool {
    return condition
}, 5*time.Second, 100*time.Millisecond)
```

**4. Database Locks:**
```go
// ✅ Use transactions properly
tx := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

// ... test operations ...

tx.Commit()
```

## Zie Ook

- [Development Guide](./development.md) - Ontwikkelomgeving
- [API Documentation](../api/rest-api.md) - API referentie
- [Deployment Guide](./deployment.md) - Productie deployment