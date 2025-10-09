# Development Guide

Complete ontwikkelhandleiding voor de DKL Email Service.

## Ontwikkelomgeving Setup

### Vereisten

**Software:**
- Go 1.21 of hoger
- PostgreSQL 13+
- Redis 6.0+ (optioneel)
- Git
- Docker (optioneel)

**Aanbevolen Tools:**
- VS Code met Go extensie
- Postman of Insomnia (API testing)
- pgAdmin of DBeaver (database management)
- Redis Commander (Redis management)

### VS Code Extensions

```json
{
    "recommendations": [
        "golang.go",
        "ms-vscode.go",
        "hbenl.vscode-test-explorer",
        "766b.go-outliner",
        "lukehoban.go-outline"
    ]
}
```

### Initiële Setup

**1. Clone Repository:**
```bash
git clone https://github.com/Jeffreasy/dklemailservice.git
cd dklemailservice
```

**2. Installeer Development Tools:**
```bash
# golangci-lint voor code linting
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# goimports voor import management
go install golang.org/x/tools/cmd/goimports@latest

# delve voor debugging
go install github.com/go-delve/delve/cmd/dlv@latest

# air voor hot reload (optioneel)
go install github.com/cosmtrek/air@latest
```

**3. Dependencies Installeren:**
```bash
go mod download
```

**4. Environment Configuratie:**
```bash
cp .env.example .env
# Bewerk .env met jouw lokale configuratie
```

**5. Database Setup:**
```sql
CREATE DATABASE dklemailservice;
CREATE USER dkluser WITH PASSWORD 'dev-password';
GRANT ALL PRIVILEGES ON DATABASE dklemailservice TO dkluser;
```

**6. Start Applicatie:**
```bash
go run main.go
```

## Project Structuur

```
dklemailservice/
├── api/                    # API utilities
│   ├── jwt_middleware.go   # JWT middleware
│   └── telegram_bot_handler.go
├── config/                 # Configuratie
│   ├── database.go         # Database config
│   └── redis.go           # Redis config
├── database/              # Database layer
│   ├── migrations.go      # Migratie manager
│   └── migrations/        # SQL migraties
├── docs/                  # Documentatie
│   ├── architecture/      # Architectuur docs
│   ├── api/              # API docs
│   ├── guides/           # Handleidingen
│   └── reports/          # Rapporten
├── handlers/             # HTTP handlers
│   ├── auth_handler.go   # Authenticatie
│   ├── email_handler.go  # Email endpoints
│   ├── contact_handler.go
│   ├── aanmelding_handler.go
│   ├── mail_handler.go
│   ├── chat_handler.go
│   └── middleware.go     # Middleware
├── logger/               # Logging
│   ├── logger.go
│   └── elk_writer.go
├── models/               # Data models
│   ├── gebruiker.go
│   ├── contact.go
│   ├── aanmelding.go
│   ├── incoming_email.go
│   ├── chat_*.go
│   └── role_rbac.go
├── repository/           # Data access layer
│   ├── interfaces.go
│   ├── factory.go
│   └── *_repository.go
├── services/            # Business logic
│   ├── email_service.go
│   ├── auth_service.go
│   ├── smtp_client.go
│   ├── rate_limiter.go
│   ├── email_auto_fetcher.go
│   ├── chat_service.go
│   └── notification_service.go
├── templates/           # Email templates
│   ├── contact_email.html
│   ├── aanmelding_email.html
│   └── newsletter.html
├── tests/              # Test suite
│   ├── *_test.go
│   └── mocks.go
├── main.go            # Entry point
├── go.mod
├── go.sum
├── Dockerfile
└── docker-compose.yml
```

## Coding Standards

### Go Style Guide

Volg de officiële [Effective Go](https://golang.org/doc/effective_go) richtlijnen.

**Formatting:**
```bash
# Format code
go fmt ./...

# Organize imports
goimports -w .

# Lint code
golangci-lint run
```

### Naming Conventions

**Packages:**
```go
package handlers  // Enkelvoud, lowercase
package services
package models
```

**Interfaces:**
```go
type EmailService interface {
    SendEmail(to, subject, body string) error
}

type AuthService interface {
    Login(ctx context.Context, email, password string) (string, error)
}
```

**Structs:**
```go
type EmailHandler struct {
    emailService EmailService
    authService  AuthService
}

type SMTPClient struct {
    host     string
    port     int
    username string
}
```

**Functions:**
```go
// Public functions: PascalCase
func NewEmailService() *EmailService

// Private functions: camelCase
func validateEmail(email string) error
```

**Constants:**
```go
const (
    MaxRetries     = 3
    DefaultTimeout = 10 * time.Second
)
```

**Variables:**
```go
var (
    ErrInvalidEmail = errors.New("invalid email")
    defaultConfig   = Config{}
)
```

### Error Handling

**Best Practices:**

```go
// ✅ Correct: Wrap errors met context
if err != nil {
    return fmt.Errorf("failed to send email: %w", err)
}

// ✅ Correct: Custom error types
var ErrInvalidCredentials = errors.New("ongeldige inloggegevens")

// ❌ Incorrect: Panic in production code
if err != nil {
    panic(err)  // Alleen in init() of tests
}

// ❌ Incorrect: Ignore errors
_ = doSomething()  // Vermijd dit
```

**Implementatie Voorbeeld:** [`services/auth_service.go:21`](../../services/auth_service.go:21)

```go
var (
    ErrInvalidCredentials = errors.New("ongeldige inloggegevens")
    ErrUserInactive = errors.New("gebruiker is inactief")
    ErrInvalidToken = errors.New("ongeldig token")
    ErrUserNotFound = errors.New("gebruiker niet gevonden")
)
```

### Logging

**Structured Logging:**

```go
// ✅ Correct: Structured logging met key-value pairs
logger.Info("email sent",
    "type", emailType,
    "recipient", recipient,
    "duration", duration,
)

// ✅ Correct: Error logging met context
logger.Error("failed to send email",
    "type", emailType,
    "error", err,
    "retry_count", retryCount,
)

// ❌ Incorrect: String concatenation
logger.Info("Email sent to " + recipient)
```

**Log Levels:**
- `Debug` - Gedetailleerde debugging info
- `Info` - Algemene informatie
- `Warn` - Waarschuwingen
- `Error` - Errors die handling vereisen
- `Fatal` - Kritieke fouten (applicatie stopt)

**Implementatie:** [`logger/logger.go`](../../logger/logger.go:1)

## Development Workflow

### Feature Development

**1. Maak Feature Branch:**
```bash
git checkout -b feature/nieuwe-feature
```

**2. Implementeer Feature:**

**Test-Driven Development:**
```go
// 1. Schrijf eerst de test
func TestNewFeature(t *testing.T) {
    // Arrange
    service := NewService()
    
    // Act
    result, err := service.DoSomething()
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}

// 2. Implementeer de functionaliteit
func (s *Service) DoSomething() (Result, error) {
    // Implementation
}

// 3. Refactor indien nodig
```

**3. Verifieer Wijzigingen:**
```bash
# Run tests
go test ./...

# Run tests met coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run linter
golangci-lint run

# Format code
go fmt ./...
goimports -w .
```

**4. Commit:**
```bash
git add .
git commit -m "feat: voeg nieuwe feature toe

- Implementeer X functionaliteit
- Voeg tests toe (coverage: 95%)
- Update documentatie"
```

### Code Review Checklist

- [ ] Tests geschreven en slagen (min 90% coverage)
- [ ] Code gedocumenteerd met comments
- [ ] Linter errors opgelost
- [ ] Error handling correct geïmplementeerd
- [ ] Logging toegevoegd op juiste niveau
- [ ] Performance impact overwogen
- [ ] Security implications bekeken
- [ ] Breaking changes gedocumenteerd
- [ ] API documentatie bijgewerkt

## Testing

### Unit Tests

**Test Structuur:**

```go
func TestEmailService_SendEmail(t *testing.T) {
    // Setup
    mockSMTP := &mockSMTP{}
    emailMetrics := services.NewEmailMetrics(time.Hour)
    prometheusMetrics := services.NewPrometheusMetrics()
    rateLimiter := services.NewRateLimiter(prometheusMetrics)
    
    emailService := services.NewEmailService(
        mockSMTP, 
        emailMetrics, 
        rateLimiter, 
        prometheusMetrics,
    )
    
    // Test cases
    t.Run("Succesvolle verzending", func(t *testing.T) {
        err := emailService.SendEmail("test@example.com", "Test", "Body")
        assert.NoError(t, err)
        assert.True(t, mockSMTP.SendCalled)
    })
    
    t.Run("SMTP fout", func(t *testing.T) {
        mockSMTP.SetShouldFail(true)
        err := emailService.SendEmail("test@example.com", "Test", "Body")
        assert.Error(t, err)
    })
}
```

**Implementatie:** [`tests/email_service_test.go:12`](../../tests/email_service_test.go:12)

### Mocking

**Mock Implementaties:** [`tests/mocks.go`](../../tests/mocks.go:1)

**SMTP Mock:**
```go
type mockSMTP struct {
    mock.Mock
    mutex       sync.Mutex
    SendCalled  bool
    LastTo      string
    LastSubject string
    LastBody    string
    SendError   error
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
```

**Auth Service Mock:**
```go
type MockAuthService struct {
    mock.Mock
}

func (m *MockAuthService) Login(ctx context.Context, email, wachtwoord string) (string, string, error) {
    args := m.Called(ctx, email, wachtwoord)
    return args.String(0), args.String(1), args.Error(2)
}
```

### Integration Tests

**Database Tests:**
```go
func TestDatabaseIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    // Setup test database
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    // Test repository operations
    repo := repository.NewGebruikerRepository(db)
    gebruiker := &models.Gebruiker{
        Naam:  "Test User",
        Email: "test@example.com",
    }
    
    err := repo.Create(context.Background(), gebruiker)
    assert.NoError(t, err)
    assert.NotEmpty(t, gebruiker.ID)
}
```

### Test Coverage

**Run Tests met Coverage:**
```bash
# Alle tests
go test ./... -v

# Met coverage
go test -coverprofile=coverage.out ./...

# HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Coverage per package
go tool cover -func=coverage.out
```

**Coverage Doelen:**
- Handlers: 90%+
- Services: 95%+
- Models: 100%
- Repositories: 90%+

### Benchmark Tests

**Performance Testing:**
```go
func BenchmarkEmailService_SendBatch(b *testing.B) {
    service := setupEmailService()
    emails := generateTestEmails(100)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        service.SendBatch(emails)
    }
}
```

**Run Benchmarks:**
```bash
go test -bench=. -benchmem ./...
```

## Debugging

### Delve Debugger

**Start Debugger:**
```bash
dlv debug
```

**Breakpoints:**
```bash
# Set breakpoint
break main.go:42
break handlers.EmailHandler.HandleContactEmail

# Run tot breakpoint
continue

# Inspect variables
print variableName
locals
args

# Step through code
next
step
stepout
```

### VS Code Debugging

**launch.json:**
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}",
            "env": {
                "LOG_LEVEL": "debug"
            },
            "args": []
        }
    ]
}
```

### Log-Based Debugging

**Debug Logging:**
```go
logger.Debug("processing request",
    "method", r.Method,
    "path", r.URL.Path,
    "remote_addr", r.RemoteAddr,
    "headers", r.Header,
)
```

**Implementatie:** [`logger/logger.go`](../../logger/logger.go:1)

## Code Quality

### Linting

**golangci-lint Configuratie:**
```yaml
# .golangci.yml
linters:
  enable:
    - gofmt
    - govet
    - errcheck
    - staticcheck
    - unused
    - gosimple
    - structcheck
    - varcheck
    - ineffassign
    - deadcode
    - typecheck

linters-settings:
  errcheck:
    check-blank: true
  govet:
    check-shadowing: true
```

**Run Linter:**
```bash
golangci-lint run
```

### Code Formatting

**Automatisch Formatteren:**
```bash
# Format all files
go fmt ./...

# Organize imports
goimports -w .
```

### Pre-commit Hooks

**Git Hook (.git/hooks/pre-commit):**
```bash
#!/bin/sh

# Format code
go fmt ./...

# Run linter
golangci-lint run

# Run tests
go test ./...

# Check if any changes
if ! git diff --quiet; then
    echo "Code was formatted, please review and commit again"
    exit 1
fi
```

## Database Development

### Migraties Maken

**Nieuwe Migratie:**
```bash
# Maak nieuw bestand
touch database/migrations/V1_XX__description.sql
```

**Migratie Template:**
```sql
-- V1_XX__description.sql

-- Up Migration
CREATE TABLE IF NOT EXISTS new_table (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Down Migration (optioneel, in comments)
-- DROP TABLE IF EXISTS new_table;
```

**Migratie Manager:** [`database/migrations.go`](../../database/migrations.go:1)

### Database Seeding

**Seed Data:**
```sql
-- database/migrations/002_seed_data.sql
INSERT INTO gebruikers (naam, email, wachtwoord_hash, rol, is_actief)
VALUES 
    ('Admin', 'admin@dekoninklijkeloop.nl', '$2a$10$...', 'admin', true),
    ('Staff', 'staff@dekoninklijkeloop.nl', '$2a$10$...', 'staff', true);
```

## API Development

### Nieuwe Endpoint Toevoegen

**1. Definieer Model:**
```go
// models/new_model.go
type NewModel struct {
    ID        string    `json:"id" gorm:"primaryKey;type:uuid"`
    Name      string    `json:"name" gorm:"not null"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
```

**2. Maak Repository:**
```go
// repository/new_model_repository.go
type NewModelRepository interface {
    Create(ctx context.Context, model *NewModel) error
    GetByID(ctx context.Context, id string) (*NewModel, error)
    List(ctx context.Context, limit, offset int) ([]*NewModel, error)
}

type newModelRepositoryImpl struct {
    db *gorm.DB
}

func NewNewModelRepository(db *gorm.DB) NewModelRepository {
    return &newModelRepositoryImpl{db: db}
}
```

**3. Maak Service:**
```go
// services/new_service.go
type NewService interface {
    CreateItem(ctx context.Context, item *NewModel) error
}

type newServiceImpl struct {
    repo NewModelRepository
}

func NewNewService(repo NewModelRepository) NewService {
    return &newServiceImpl{repo: repo}
}
```

**4. Maak Handler:**
```go
// handlers/new_handler.go
type NewHandler struct {
    service     NewService
    authService AuthService
}

func (h *NewHandler) HandleCreate(c *fiber.Ctx) error {
    var request NewModel
    if err := c.BodyParser(&request); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request",
        })
    }
    
    if err := h.service.CreateItem(c.Context(), &request); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    return c.JSON(fiber.Map{
        "success": true,
        "data": request,
    })
}
```

**5. Registreer Routes:**
```go
// main.go
newHandler := handlers.NewNewHandler(serviceFactory.NewService, serviceFactory.AuthService)
api.Post("/new", 
    handlers.AuthMiddleware(serviceFactory.AuthService),
    newHandler.HandleCreate,
)
```

## Testing Strategies

### Table-Driven Tests

**Implementatie:**
```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {
            name:    "valid email",
            email:   "test@example.com",
            wantErr: false,
        },
        {
            name:    "invalid email - no @",
            email:   "testexample.com",
            wantErr: true,
        },
        {
            name:    "invalid email - no domain",
            email:   "test@",
            wantErr: true,
        },
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

### Mock Setup

**Voorbeeld:** [`tests/mocks.go:108`](../../tests/mocks.go:108)

```go
type mockSMTP struct {
    mock.Mock
    mutex       sync.Mutex
    SendCalled  bool
    LastTo      string
    LastSubject string
    SendError   error
}

func (m *mockSMTP) Send(msg *services.EmailMessage) error {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    m.SendCalled = true
    m.LastTo = msg.To
    if m.SendError != nil {
        return m.SendError
    }
    return nil
}
```

**Gebruik in Tests:**
```go
func TestWithMock(t *testing.T) {
    mockSMTP := &mockSMTP{}
    mockSMTP.On("Send", mock.Anything).Return(nil)
    
    service := NewEmailService(mockSMTP, ...)
    err := service.SendEmail("test@example.com", "Subject", "Body")
    
    assert.NoError(t, err)
    mockSMTP.AssertExpectations(t)
}
```

## Performance Profiling

### CPU Profiling

```bash
# Run met CPU profiling
go test -cpuprofile=cpu.prof -bench=.

# Analyze
go tool pprof cpu.prof
```

**In pprof:**
```
(pprof) top10
(pprof) list FunctionName
(pprof) web
```

### Memory Profiling

```bash
# Run met memory profiling
go test -memprofile=mem.prof -bench=.

# Analyze
go tool pprof mem.prof
```

### Race Detection

```bash
# Run tests met race detector
go test -race ./...

# Run applicatie met race detector
go run -race main.go
```

## Hot Reload Development

### Met Air

**Installeer Air:**
```bash
go install github.com/cosmtrek/air@latest
```

**.air.toml:**
```toml
root = "."
tmp_dir = "tmp"

[build]
  cmd = "go build -o ./tmp/main ."
  bin = "tmp/main"
  include_ext = ["go", "html"]
  exclude_dir = ["tmp", "vendor"]
  delay = 1000

[log]
  time = true
```

**Start met Hot Reload:**
```bash
air
```

## Environment Management

### Multiple Environments

**.env.development:**
```bash
LOG_LEVEL=debug
DB_HOST=localhost
SMTP_HOST=smtp.mailtrap.io
```

**.env.production:**
```bash
LOG_LEVEL=info
DB_HOST=production-db
SMTP_HOST=smtp.sendgrid.net
```

**Load Specific Environment:**
```bash
# Development
cp .env.development .env
go run main.go

# Production
cp .env.production .env
go run main.go
```

## Troubleshooting

### Common Issues

**1. Import Cycle:**
```
package X imports package Y
package Y imports package X
```

**Oplossing:** Gebruik interfaces of herstructureer packages

**2. Race Conditions:**
```bash
go test -race ./...
```

**Oplossing:** Gebruik mutexes of channels

**3. Memory Leaks:**
```bash
go test -memprofile=mem.prof
go tool pprof mem.prof
```

**Oplossing:** Check goroutine leaks, defer cleanup

## Best Practices

### Concurrency

**Worker Pools:**
```go
// Implementatie voorbeeld
pool := make(chan struct{}, maxWorkers)
for _, item := range items {
    pool <- struct{}{}
    go func(i Item) {
        defer func() { <-pool }()
        processItem(i)
    }(item)
}
```

**Context Usage:**
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

result, err := doSomethingWithContext(ctx)
```

### Resource Management

**Defer Cleanup:**
```go
file, err := os.Open(filename)
if err != nil {
    return err
}
defer file.Close()

// Database transactions
tx := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()
```

### Configuration

**Environment Variables:**
```go
type Config struct {
    SMTP struct {
        Host     string `env:"SMTP_HOST,required"`
        Port     int    `env:"SMTP_PORT" default:"587"`
        Username string `env:"SMTP_USER,required"`
        Password string `env:"SMTP_PASSWORD,required"`
    }
}
```

## Zie Ook

- [Testing Guide](./testing.md) - Uitgebreide test procedures
- [Deployment Guide](./deployment.md) - Productie deployment
- [Security Guide](./security.md) - Security best practices
- [API Documentation](../api/rest-api.md) - API referentie