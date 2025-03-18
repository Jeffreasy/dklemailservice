# Testing Handleiding

## Overzicht

Deze handleiding beschrijft de test setup en procedures voor de DKL Email Service, inclusief:
- Unit tests
- Integration tests
- Performance tests
- Mocking
- Test coverage

## Test Setup

### Vereisten
- Go 1.21 of hoger
- `testify` voor assertions
- Docker voor integration tests

### Test Directory Structuur
```
tests/
├── aanmelding_handler_test.go   # Aanmelding handler tests
├── elk_integration_test.go      # ELK Stack integration tests
├── elk_writer_test.go          # ELK writer unit tests
├── email_batcher_test.go       # Email batcher unit tests
├── email_metrics_test.go       # Email metrics unit tests
├── email_service_metrics_test.go # Email service metrics tests
├── email_service_test.go       # Email service unit tests
├── handler_test.go             # Generic handler tests
├── logger_test.go             # Logger unit tests
├── mocks.go                   # Mock implementations
├── rate_limit_test.go         # Rate limiter tests
├── smtp_client_test.go        # SMTP client tests
├── template_test.go           # Template tests
└── test_helper.go             # Test utilities
```

## Unit Tests

### Handler Tests
```go
// tests/aanmelding_handler_test.go
func TestAanmeldingHandler(t *testing.T) {
    tests := []struct {
        name       string
        input      models.AanmeldingRequest
        mockSetup  func(*mocks.MockEmailService)
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
            },
            mockSetup: func(m *mocks.MockEmailService) {
                m.EXPECT().
                    SendAanmeldingEmail(gomock.Any()).
                    Return(nil)
            },
            wantStatus: http.StatusOK,
            wantErr:    false,
        },
        // ... meer test cases
    }
    // ... test implementatie
}
```

### Service Tests
```go
// tests/email_service_test.go
func TestEmailService_SendContactEmail(t *testing.T) {
    tests := []struct {
        name      string
        input     models.ContactEmailData
        mockSetup func(*mocks.MockSMTPClient)
        wantErr   bool
    }{
        // ... test cases
    }
    // ... test implementatie
}
```

### Integration Tests
```go
// tests/elk_integration_test.go
func TestELKIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    // Setup ELK client
    elkWriter := logger.NewELKWriter(logger.ELKConfig{
        Endpoint: "http://localhost:9200",
        Index:    "test-logs",
    })

    // Test log writing
    err := elkWriter.Write([]byte(`{"level":"info","message":"test"}`))
    require.NoError(t, err)

    // Verify log entry
    // ... verificatie logica
}
```

### Performance Tests
```go
// tests/email_service_test.go
func BenchmarkEmailService_SendBatch(b *testing.B) {
    service := services.NewEmailService(
        &mocks.MockSMTPClient{},
        services.EmailServiceConfig{
            BatchSize: 10,
            BatchInterval: time.Second,
        },
    )

    emails := generateTestEmails(100)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        err := service.SendBatch(emails)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

## Mocking

### Mock Implementaties
```go
// tests/mocks.go
type MockSMTPClient struct {
    mock.Mock
}

func (m *MockSMTPClient) Send(email *models.Email) error {
    args := m.Called(email)
    return args.Error(0)
}

type MockEmailService struct {
    mock.Mock
}

func (m *MockEmailService) SendContactEmail(data models.ContactEmailData) error {
    args := m.Called(data)
    return args.Error(0)
}
```

## Test Coverage

### Coverage Doelen
- Handlers: 90% coverage
- Services: 95% coverage
- Models: 100% coverage
- Templates: 100% coverage

### Coverage Rapport Genereren
```bash
# Run tests with coverage
go test -coverprofile=coverage.out ./...

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html

# View coverage stats
go tool cover -func=coverage.out
```

## Test Best Practices

### Code Organization
1. Gebruik table-driven tests
2. Groepeer gerelateerde tests
3. Gebruik subtests voor variaties
4. Isoleer test dependencies
5. Gebruik meaningful test names

### Test Data
1. Gebruik constants voor test data
2. Implementeer test helpers
3. Cleanup test resources
4. Gebruik realistic test data
5. Test edge cases

### Assertions
1. Gebruik testify assertions
2. Check error types
3. Verify side effects
4. Test timeouts
5. Validate error messages

### Mocking
1. Mock external dependencies
2. Gebruik interface mocking
3. Verify mock calls
4. Test error scenarios
5. Mock time dependencies

## Debugging Tests

### Common Issues
1. Race Conditions
```bash
# Run tests with race detector
go test -race ./...
```

2. Timeouts
```go
// Set test timeout
func TestLongRunning(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping long running test")
    }
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // Test implementation
}
```

3. Resource Leaks
```go
// Use cleanup function
func TestWithResources(t *testing.T) {
    resource := setupTestResource()
    t.Cleanup(func() {
        resource.Close()
    })
    
    // Test implementation
}
```

### Test Logging
```go
// Enable verbose logging
func TestWithLogging(t *testing.T) {
    if testing.Verbose() {
        t.Logf("Running test with config: %+v", config)
    }
    
    // Test implementation
    
    if testing.Verbose() {
        t.Logf("Test result: %+v", result)
    }
}

## API Endpoint Tests

De volgende PowerShell scripts zijn beschikbaar om de API endpoints te testen:

- `test_api_light.ps1`: Test alle API endpoints zonder stress test
- `test_api_full.ps1`: Test alle API endpoints inclusief stress test
- `stress_test.ps1`: Voert alleen de stress test uit
- `check_mail_logs.ps1`: Analyseert email verzend logs

### Parameters voor test scripts

`test_api_light.ps1` en `test_api_full.ps1` accepteren de volgende parameters:

```
-BaseUrl <string>       : API base URL (default: http://localhost:8080)
-DetailedHealth         : Geeft gedetailleerde health check informatie
-TestAdminEndpoints     : Test admin endpoints (vereist admin API key)
-SkipStressTest         : Slaat de stress test over (alleen bij test_api_full.ps1)
-TestMetrics            : Test metrics endpoints (vereist admin API key)
-AdminAPIKey <string>   : Admin API key voor beveiligde endpoints
-TestMailEndpoints      : Test email functionaliteit endpoints
-IncludeSecuredEndpoints: Test beveiligde email endpoints (vereist JWT token)
-JWTToken <string>      : JWT token voor beveiligde endpoints
```

### Test voorbeelden

```powershell
# Basis API test
./test_api_light.ps1

# API test met gedetailleerde health check
./test_api_light.ps1 -DetailedHealth

# Test inclusief mail endpoints
./test_api_light.ps1 -TestMailEndpoints

# Test inclusief beveiligde mail endpoints
./test_api_light.ps1 -TestMailEndpoints -IncludeSecuredEndpoints -JWTToken "your-jwt-token"

# Test van productie API
./test_api_light.ps1 -BaseUrl "https://api.dekoninklijkeloop.nl"
```

## EmailAutoFetcher Tests

De EmailAutoFetcher component verwerkt inkomende emails via IMAP en slaat deze op in de database. De volgende test procedures zijn beschikbaar:

### Unit Tests

Unit tests voor de EmailAutoFetcher zijn beschikbaar in:
- `services/email_auto_fetcher_test.go`

Run specifieke tests met:

```bash
go test -v ./services -run TestEmailAutoFetcher
```

### Integratie Tests

Voor integratie tests met echte mail servers:

1. Configureer test credentials in een `.env.test` bestand:

```
TEST_EMAIL_HOST=imap.example.com
TEST_EMAIL_PORT=993
TEST_EMAIL_USER=testuser@example.com
TEST_EMAIL_PASSWORD=testpassword
TEST_EMAIL_USE_TLS=true
```

2. Run integratie tests:

```bash
go test -v ./integration -run TestEmailFetching
```

### Handmatige Tests

1. Configureer EmailAutoFetcher voor snel testen:

```
# In .env file
EMAIL_FETCH_INTERVAL=1
EMAIL_INFO_HOST=imap.example.com
EMAIL_INFO_PORT=993
EMAIL_INFO_USERNAME=info@dekoninklijkeloop.nl
EMAIL_INFO_PASSWORD=your-password
EMAIL_INFO_USE_TLS=true
```

2. Start de applicatie en monitor logs:

```bash
go run main.go | grep "EmailAutoFetcher"
```

3. Stuur een test email naar het geconfigureerde adres en verifieer dat deze wordt opgehaald.

4. Test de handmatige fetch API:

```bash
# Login en JWT token verkrijgen
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}' | jq -r '.token')

# Handmatig emails fetchen
curl -X POST http://localhost:8080/api/mail/fetch \
  -H "Authorization: Bearer $TOKEN"

# Bekijk ontvangen emails
curl -X GET http://localhost:8080/api/mail \
  -H "Authorization: Bearer $TOKEN"
```

### Test cases voor EmailAutoFetcher

| Test Case | Beschrijving | Verwacht Resultaat |
|-----------|--------------|-------------------|
| Automatisch ophalen | Start de applicatie met EmailAutoFetcher ingeschakeld | Logs tonen periodieke fetch operaties volgens ingesteld interval |
| Duplicate detection | Verstuur dezelfde email meerdere keren | Email wordt slechts één keer opgeslagen in database |
| Error recovery | Simuleer een tijdelijke IMAP verbindingsfout | EmailAutoFetcher herstelt en probeert opnieuw bij volgende interval |
| Graceful shutdown | Stop applicatie tijdens fetch operatie | Fetch operatie wordt veilig afgebroken |
| Configuratie wijziging | Wijzig interval via env variabele | EmailAutoFetcher past interval aan |
| Handmatige fetch | Trigger handmatige fetch via API | Nieuwe emails worden onmiddellijk opgehaald |

## Logging

// ... existing code ... 