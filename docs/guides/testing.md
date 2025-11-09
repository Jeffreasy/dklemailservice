# Testing Guide

Complete testing guide voor de DKL Email Service.

## Overview

De DKL Email Service gebruikt verschillende soorten tests om code quality en betrouwbaarheid te waarborgen:

- **Unit Tests** - Test individuele functies en methoden
- **Integration Tests** - Test component interacties
- **API Tests** - Test HTTP endpoints
- **Database Tests** - Test database operaties
- **Performance Tests** - Test onder load

---

## Test Setup

### Prerequisites

```bash
# Install CGO (required voor SQLite in tests)
# Windows: MinGW-w64 of TDM-GCC
# Linux: gcc installeren via package manager
# macOS: Xcode Command Line Tools

# Verify CGO
go env CGO_ENABLED
# Should return: 1
```

### Environment Configuration

Create `.env.test` voor test environment:

```env
# Test Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=test_password
DB_NAME=dklemailservice_test
DB_SSL_MODE=disable

# Test Email (gebruik fake SMTP)
SMTP_HOST=localhost
SMTP_PORT=1025
SMTP_USER=test
SMTP_PASSWORD=test
SMTP_FROM=test@example.com

# JWT
JWT_SECRET=test_secret_key_min_32_characters_long
```

---

## Running Tests

### Using Test Scripts (Recommended)

**Linux/macOS:**
```bash
# Alle tests
./scripts/run_tests.sh

# Met coverage
./scripts/run_tests.sh --coverage

# Specifiek package
./scripts/run_tests.sh --package=handlers
```

**Windows:**
```batch
REM Alle tests
scripts\run_tests.bat

REM Met coverage
scripts\run_tests.bat --coverage

REM Specifiek package
scripts\run_tests.bat --package=handlers
```

### Manual Test Execution

```bash
# Enable CGO
export CGO_ENABLED=1  # Linux/macOS
set CGO_ENABLED=1     # Windows

# Run all tests
go test ./... -v

# Run specific package
go test ./tests -v
go test ./handlers -v
go test ./services -v

# Run with race detection
go test -race ./...

# Run with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## Unit Tests

### Handler Tests

**Location:** `tests/*_handler_test.go`

**Example:**

```go
func TestContactHandler_HandleContactEmail(t *testing.T) {
    // Setup
    app := fiber.New()
    handler := setupTestHandler()
    
    // Test case
    req := httptest.NewRequest("POST", "/api/contact-email", 
        strings.NewReader(`{"naam":"Test","email":"test@example.com"}`))
    req.Header.Set("Content-Type", "application/json")
    
    resp, _ := app.Test(req)
    
    // Assertions
    assert.Equal(t, 200, resp.StatusCode)
}
```

**Run handler tests:**
```bash
go test ./tests -run TestContact -v
go test ./tests -run TestAuth -v
go test ./tests -run TestParticipant -v
```

---

### Service Tests

**Location:** `tests/*_service_test.go`

**Example:**

```go
func TestEmailService_SendEmail(t *testing.T) {
    // Setup mock
    mockSMTP := &MockSMTPClient{}
    service := NewEmailService(mockSMTP)
    
    // Test
    err := service.SendEmail("test@example.com", "Subject", "Body")
    
    // Verify
    assert.NoError(t, err)
    assert.Equal(t, 1, mockSMTP.SendCount)
}
```

**Run service tests:**
```bash
go test ./tests -run TestEmail -v
go test ./tests -run TestAuth -v
go test ./tests -run TestPermission -v
```

---

## Integration Tests

### Database Integration Tests

**Location:** `tests/database_*_test.go`

```bash
# Run database tests
go test ./tests -run TestDatabase -v

# Specific database tests
go test ./tests -run TestDatabaseMigrations -v
go test ./tests -run TestModelDatabaseAlignment -v
go test ./tests -run TestConsolidatedMigrations -v
```

**Test Coverage:**
- Schema migrations
- Model-database alignment
- Foreign key constraints
- Index creation
- Data integrity

---

### RBAC Integration Tests

**Location:** `tests/rbac_integration_test.go`

```bash
# Run RBAC tests
go test ./tests -run TestRBAC -v
```

**Test Coverage:**
- Permission checking
- Role assignment
- Cache invalidation
- Redis fallback

---

### ELK Integration Tests

**Location:** `tests/elk_*_test.go`

```bash
# Run ELK tests
go test ./tests -run TestELK -v
```

---

## Performance Tests

### Email Batching Tests

**Location:** `tests/email_batcher_test.go`

```bash
go test ./tests -run TestEmailBatcher -v
```

**Metrics Tested:**
- Batch processing throughput
- Memory usage
- Concurrent sends
- Error handling

---

### Auth Performance Tests

**Location:** `tests/auth_performance_test.go`

```bash
go test ./tests -run TestAuthPerformance -v -benchmem
```

**Benchmarks:**
- Login throughput
- Password hashing speed
- JWT generation/validation
- Permission checks

---

### Rate Limit Tests

**Location:** `tests/rate_limit_test.go`

```bash
go test ./tests -run TestRateLimit -v
```

---

## Security Tests

### Auth Security Tests

**Location:** `tests/auth_security_test.go`

```bash
go test ./tests -run TestAuthSecurity -v
```

**Security Checks:**
- SQL injection prevention
- XSS protection
- Password strength
- JWT validation
- Token expiration
- CSRF protection

---

### Permission Tests

**Location:** `tests/permission_middleware_test.go`

```bash
go test ./tests -run TestPermission -v
```

---

## API Endpoint Testing

### Using test.http File

**Location:** `test.http`

```http
### Login
POST http://localhost:8080/api/auth/login
Content-Type: application/json

{
  "email": "admin@example.com",
  "password": "password"
}

### Use token from previous request
@token = {{login.response.body.data.access_token}}

### Get Contact Forms
GET http://localhost:8080/api/contact
Authorization: Bearer {{token}}
```

**VS Code Extension:** REST Client

---

### Using cURL

```bash
# Health check
curl http://localhost:8080/api/health

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password"}'

# Authenticated request
curl http://localhost:8080/api/contact \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

### Using Postman

1. Import collection (indien beschikbaar)
2. Setup environment variables:
   - `base_url`: http://localhost:8080
   - `access_token`: (from login response)

3. Create requests met authorization header:
   ```
   Authorization: Bearer {{access_token}}
   ```

---

## Mock Services

### Test Mode

Use `X-Test-Mode` header to enable test mode:

```bash
curl -X POST http://localhost:8080/api/contact-email \
  -H "X-Test-Mode: true" \
  -H "Content-Type: application/json" \
  -d '{"naam":"Test","email":"test@example.com","bericht":"Test"}'
```

**Test Mode Features:**
- Emails niet echt verzonden
- Snellere execution
- Geen externe dependencies
- Deterministic results

---

### Mock Database

**Location:** `tests/auth_test_mocks.go`

```go
// Mock database for unit tests
mockDB := &MockDatabase{
    Users: make(map[string]*models.Gebruiker),
}

// Use in tests
service := NewAuthService(mockDB)
```

---

## Coverage Reports

### Generate Coverage

```bash
# Generate coverage profile
go test ./... -coverprofile=coverage.out

# View in browser
go tool cover -html=coverage.out

# Coverage percentage
go tool cover -func=coverage.out

# Coverage by package
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out | grep -E '^total:'
```

### Target Coverage

**Goals:**
- Handlers: >80%
- Services: >85%
- Models: >70%
- Overall: >80%

**Current Status:**
Check `tests/DATABASE_TEST_SUITE_SUMMARY.md` voor details.

---

## Test Best Practices

### 1. Test Structure

```go
func TestFunctionName(t *testing.T) {
    // Arrange - Setup test data
    input := "test data"
    expected := "expected result"
    
    // Act - Execute function
    result := FunctionToTest(input)
    
    // Assert - Verify results
    assert.Equal(t, expected, result)
}
```

### 2. Table-Driven Tests

```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {"valid email", "test@example.com", false},
        {"invalid email", "invalid", true},
        {"empty email", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.email)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateEmail() error = %v, wantErr %v", 
                    err, tt.wantErr)
            }
        })
    }
}
```

### 3. Mock External Services

```go
// Mock SMTP client
type MockSMTPClient struct {
    SendCount int
    LastEmail *Email
    ShouldFail bool
}

func (m *MockSMTPClient) Send(email *Email) error {
    m.SendCount++
    m.LastEmail = email
    if m.ShouldFail {
        return errors.New("mock error")
    }
    return nil
}
```

### 4. Cleanup After Tests

```go
func TestWithDatabase(t *testing.T) {
    // Setup
    db := setupTestDB()
    defer db.Close()  // Cleanup
    
    // Create test data
    user := createTestUser(db)
    defer deleteTestUser(db, user.ID)  // Cleanup
    
    // Test...
}
```

---

## Continuous Integration

### GitHub Actions Workflow

**Location:** `.github/workflows/test.yml` (example)

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:17
        env:
          POSTGRES_PASSWORD: test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
      
      redis:
        image: redis:7-alpine
        options: >-
          --health-cmd "redis-cli ping"
          --health-interval 10s
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.23.0'
      
      - name: Install dependencies
        run: go mod download
      
      - name: Run tests
        run: |
          export CGO_ENABLED=1
          go test ./... -v -race -coverprofile=coverage.out
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out
```

---

## Load Testing

### Using k6

**Install k6:**
```bash
# macOS
brew install k6

# Linux
sudo apt-get install k6

# Windows
choco install k6
```

**Load Test Script:**

```javascript
// load_test.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 10,        // 10 virtual users
  duration: '30s', // Test duration
};

export default function() {
  const res = http.get('http://localhost:8080/api/health');
  
  check(res, {
    'status is 200': (r) => r.status === 200,
    'response time < 200ms': (r) => r.timings.duration < 200,
  });
  
  sleep(1);
}
```

**Run load test:**
```bash
k6 run load_test.js

# With results output
k6 run --out json=results.json load_test.js
```

---

### API Endpoint Load Testing

```javascript
// api_load_test.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 20 },  // Ramp up to 20 users
    { duration: '1m', target: 20 },   // Stay at 20 users
    { duration: '30s', target: 0 },   // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% requests < 500ms
    http_req_failed: ['rate<0.01'],   // Error rate < 1%
  },
};

const BASE_URL = 'http://localhost:8080';
let token;

export function setup() {
  // Login to get token
  const loginRes = http.post(`${BASE_URL}/api/auth/login`, 
    JSON.stringify({
      email: 'test@example.com',
      password: 'password',
    }), {
      headers: { 'Content-Type': 'application/json' },
    }
  );
  
  const data = JSON.parse(loginRes.body);
  return { token: data.data.access_token };
}

export default function(data) {
  const headers = {
    'Authorization': `Bearer ${data.token}`,
    'Content-Type': 'application/json',
  };
  
  // Test various endpoints
  const endpoints = [
    { method: 'GET', url: '/api/health' },
    { method: 'GET', url: '/api/events' },
    { method: 'GET', url: '/api/contact' },
  ];
  
  endpoints.forEach(endpoint => {
    const res = http.request(endpoint.method, 
      `${BASE_URL}${endpoint.url}`, null, { headers });
    
    check(res, {
      [`${endpoint.url} status is 200`]: (r) => r.status === 200,
    });
  });
  
  sleep(1);
}
```

---

## Test Data Management

### Creating Test Data

```sql
-- Create test users
INSERT INTO gebruikers (email, wachtwoord_hash, naam, actief)
VALUES 
  ('test@example.com', '$2a$12$...', 'Test User', true),
  ('admin@example.com', '$2a$12$...', 'Admin User', true);

-- Assign roles
INSERT INTO user_roles (gebruiker_id, role_id)
SELECT u.id, r.id
FROM gebruikers u, roles r
WHERE u.email = 'admin@example.com' AND r.name = 'admin';
```

### Cleanup Test Data

```sql
-- Remove test data
DELETE FROM gebruikers WHERE email LIKE '%@example.com';
DELETE FROM contact_formulieren WHERE email LIKE '%@test.com';
DELETE FROM participants WHERE email LIKE '%@test.com';
```

---

## Specific Test Suites

### Database Tests

**Run migrations test:**
```bash
go test ./tests -run TestDatabaseMigrations -v
```

**Model alignment test:**
```bash
go test ./tests -run TestModelDatabaseAlignment -v
```

**Full database suite:**
```bash
go test ./tests -run TestDatabase -v
```

---

### Email Tests

```bash
# Email service tests
go test ./tests -run TestEmailService -v

# Email metrics tests
go test ./tests -run TestEmailMetrics -v

# Email batcher tests
go test ./tests -run TestEmailBatcher -v
```

---

### Authentication Tests

```bash
# Auth service tests
go test ./tests -run TestAuthService -v

# Auth handler tests
go test ./tests -run TestAuthHandler -v

# Auth security tests
go test ./tests -run TestAuthSecurity -v

# Auth performance tests
go test ./tests -run TestAuthPerformance -v -bench=.
```

---

### Permission Tests

```bash
# Permission service tests
go test ./tests -run TestPermissionService -v

# Permission middleware tests
go test ./tests -run TestPermissionMiddleware -v

# RBAC integration tests
go test ./tests -run TestRBACIntegration -v
```

---

## Debugging Tests

### Verbose Output

```bash
# Maximum verbosity
go test ./tests -v -count=1

# Show all logs
go test ./tests -v -args -test.v
```

### Run Single Test

```bash
# Run specific test function
go test ./tests -run TestContactHandler_HandleContactEmail -v

# Run test with specific pattern
go test ./tests -run "TestAuth.*Login" -v
```

### Skip Cache

```bash
# Force rerun (skip cache)
go test ./tests -count=1 -v
```

---

## Test Coverage Analysis

### By Package

```bash
# Generate coverage per package
go test ./handlers -coverprofile=handlers.out
go test ./services -coverprofile=services.out
go test ./models -coverprofile=models.out

# Combined coverage
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

### Coverage Report

```bash
# HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Open in browser
# Linux
xdg-open coverage.html

# macOS  
open coverage.html

# Windows
start coverage.html
```

### Identify Uncovered Code

```bash
# Show uncovered lines
go tool cover -func=coverage.out | grep -v "100.0%"

# Focus on critical packages
go test ./handlers -coverprofile=handlers.out
go tool cover -func=handlers.out | grep -v "100.0%"
```

---

## Mocking Strategies

### HTTP Client Mock

```go
type MockHTTPClient struct {
    DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
    if m.DoFunc != nil {
        return m.DoFunc(req)
    }
    return nil, errors.New("DoFunc not set")
}
```

### Database Mock

```go
type MockUserRepository struct {
    Users map[string]*models.Gebruiker
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.Gebruiker) error {
    m.Users[user.Email] = user
    return nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.Gebruiker, error) {
    user, exists := m.Users[email]
    if !exists {
        return nil, errors.New("user not found")
    }
    return user, nil
}
```

---

## Testing WebSocket

### Manual WebSocket Testing

```bash
# Install wscat
npm install -g wscat

# Connect to WebSocket
wscat -c "ws://localhost:8080/ws/steps?token=YOUR_JWT_TOKEN"

# Send message
> {"type":"ping"}

# Receive response
< {"type":"pong","timestamp":"2025-01-08T10:00:00Z"}
```

### Automated WebSocket Tests

```go
func TestStepsWebSocket(t *testing.T) {
    // Setup WebSocket server
    app := fiber.New()
    handler := NewStepsWebSocketHandler(...)
    handler.RegisterRoutes(app)
    
    // Create test connection
    ws, _, err := websocket.DefaultDialer.Dial(
        "ws://localhost:8080/ws/steps?token=test_token", nil)
    assert.NoError(t, err)
    defer ws.Close()
    
    // Test message
    err = ws.WriteJSON(map[string]interface{}{
        "type": "ping",
    })
    assert.NoError(t, err)
    
    // Read response
    var response map[string]interface{}
    err = ws.ReadJSON(&response)
    assert.NoError(t, err)
    assert.Equal(t, "pong", response["type"])
}
```

---

## Benchmarking

### Run Benchmarks

```bash
# All benchmarks
go test -bench=. ./...

# Specific benchmark
go test -bench=BenchmarkAuthLogin ./tests

# With memory profiling
go test -bench=. -benchmem ./...

# CPU profiling
go test -bench=. -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof
```

### Example Benchmark

```go
func BenchmarkPasswordHash(b *testing.B) {
    password := "TestPassword123!"
    
    for i := 0; i < b.N; i++ {
        bcrypt.GenerateFromPassword([]byte(password), 12)
    }
}

func BenchmarkPermissionCheck(b *testing.B) {
    service := setupPermissionService()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        service.HasPermission(ctx, userID, "contact", "read")
    }
}
```

---

## Test Environment

### Local Test Database

```bash
# Create test database
createdb -U postgres dklemailservice_test

# Run migrations on test DB
DB_NAME=dklemailservice_test go run database/migrations/run_migrations.go

# Drop test database
dropdb -U postgres dklemailservice_test
```

### Docker Test Environment

```yaml
# docker-compose.test.yml
version: '3.8'

services:
  postgres-test:
    image: postgres:17-alpine
    environment:
      POSTGRES_DB: dklemailservice_test
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: test
    ports:
      - "5433:5432"
  
  redis-test:
    image: redis:7-alpine
    ports:
      - "6380:6379"
```

```bash
# Start test environment
docker-compose -f docker-compose.test.yml up -d

# Run tests
DB_PORT=5433 REDIS_PORT=6380 go test ./... -v

# Stop test environment
docker-compose -f docker-compose.test.yml down
```

---

## Common Test Issues

### CGO Not Enabled

**Error:** `cannot find package "github.com/mattn/go-sqlite3"`

**Solution:**
```bash
export CGO_ENABLED=1  # Linux/macOS
set CGO_ENABLED=1     # Windows

# Or use test scripts which handle this automatically
./scripts/run_tests.sh
```

---

### Database Connection Failed

**Error:** `connection refused`

**Solution:**
```bash
# Check PostgreSQL is running
sudo systemctl status postgresql

# Or start via Docker
docker-compose -f docker-compose.test.yml up -d postgres-test
```

---

### Test Timeout

**Error:** `test timed out after 10m0s`

**Solution:**
```bash
# Increase timeout
go test ./... -timeout 30m -v
```

---

## Test Documentation

### Test Suite Summary

See `tests/DATABASE_TEST_SUITE_SUMMARY.md` voor:
- Test coverage overzicht
- Test status per component
- Known issues
- Improvement suggestions

---

## Next Steps

- [API Testing](../api/QUICK_REFERENCE.md)
- [Performance Monitoring](./MONITORING.md)
- [Deployment](./DEPLOYMENT.md)

---

## Complete API Test Reports (November 2025)

### Test Report: 2025-11-08 23:36 CET

**Database:** Fresh PostgreSQL 17 Docker Setup
**Status:** ✅ System Fully Operational
**Success Rate:** 84% → 100% (after fixes)

#### Executive Summary
- **Total Endpoints:** 38 tested
- **Successful:** 32 endpoints (84% initial)
- **Failed:** 6 endpoints (16% - all fixed)
- **Database:** 48 tables, all migrations successful
- **RBAC:** 83 permissions, 4 roles operational

#### Test Results by Category

**Public Endpoints:** 15/15 (100%)
- ✅ Health check, Events, Videos, Partners, Sponsors
- ✅ Photos, Albums, Radio, Program, Social Links/Embeds
- ✅ Title Sections (fixed)
- ✅ Under Construction (fixed to public)

**Authentication:** 2/2 (100%)
- ✅ Login successful with admin credentials
- ✅ Profile endpoint with 83 permissions

**Admin Endpoints:**
- ✅ User Management: 1/1
- ✅ Contact Management: 1/1
- ✅ Participant Management: 1/1
- ✅ CMS Admin: 6/6 (all working)
- ✅ Newsletter: 1/1
- ✅ Notulen: 1/1
- ✅ Mail Management: 2/2 (after route fix)
- ✅ RBAC Management: 2/2 (after route addition)
- ✅ Gamification: 3/3 (after badge SQL fix)
- ✅ Notifications: 1/1 (after route addition)

#### Performance Metrics
- **Response Times:** 100-180ms average
- **Memory Usage:** 2.8-3.9 MB (very efficient)
- **Goroutines:** 16 (stable)
- **Database:** All 51 tables operational
- **Migrations:** V01-V32 all successful

#### Issues Found and Fixed

**1. Badge SQL Parameter (SQLSTATE 42P18)** - ✅ Fixed
- Location: `repository/badge_repository.go:104`
- Solution: Query split for nil vs non-nil participantID

**2. Missing Routes** - ✅ Fixed
- Routes added: `/api/roles`, `/api/permissions`, `/api/achievements`, `/api/notifications`, `/api/title-sections`
- Location: `main.go:686-712`

**3. Under Construction Auth** - ✅ Fixed
- Made `/api/under-construction` publicly accessible
- Location: `handlers/under_construction_handler.go:37`

**4. Mail Route Ordering** - ✅ Fixed
- Fixed `/api/mail/unprocessed` 404 error
- Location: `handlers/mail_handler.go:175`

**5. V27 Migration** - ✅ Fixed
- Column naming: `name` → `type`/`priority`
- Location: `database/migrations/V27_CONSOLIDATED__normalize_status_and_type_fields.sql:191`

#### Test Scripts Available
- `test_all_fixed_endpoints.ps1` - Comprehensive endpoint verification
- `test.http` - VS Code REST Client tests
- MCP Server tools for automated testing

#### Verification Status
- ✅ All 6 failing endpoints now operational
- ✅ Service starts without errors
- ✅ All migrations successful
- ✅ RBAC fully functional
- ✅ Public endpoints: 100% success
- ✅ Auth endpoints: proper security enforced

---

**Last Updated:** 2025-11-08