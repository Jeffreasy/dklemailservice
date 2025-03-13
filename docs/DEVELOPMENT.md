# Development Handleiding

## Ontwikkelomgeving Setup

### Vereisten
- Go 1.21 of hoger
- Docker (optioneel, voor lokale development)
- Git
- Een code editor met Go ondersteuning (VS Code aanbevolen)
- golangci-lint voor code linting

### VS Code Extensions
- Go (ms-vscode.go)
- Go Test Explorer
- Go Outliner
- Go Doc

### Initiële Setup

1. Clone de repository:
```bash
git clone https://github.com/Jeffreasy/dklemailservice.git
cd dklemailservice
```

2. Installeer development tools:
```bash
# Installeer golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Installeer andere development dependencies
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/go-delve/delve/cmd/dlv@latest
```

3. Configureer environment:
```bash
# Kopieer environment configuratie
cp .env.example .env
```

## Code Structuur

### Directory Layout
```
.
├── docs/                   # Documentatie
│   ├── API.md            # API documentatie
│   ├── DEPLOYMENT.md     # Deployment instructies
│   ├── DEVELOPMENT.md    # Development guidelines
│   ├── MONITORING.md     # Monitoring setup
│   ├── SECURITY.md       # Security best practices
│   ├── TEMPLATES.md      # Template documentatie
│   └── TESTING.md        # Test procedures
├── handlers/              # HTTP handlers
│   ├── email_handler.go  # Email endpoints
│   ├── health_handler.go # Health check endpoint
│   └── metrics_handler.go # Metrics endpoints
├── logger/                # Logging package
│   ├── elk_writer.go     # ELK logging implementatie
│   ├── logger.go         # Logger setup
│   ├── mock_writer.go    # Mock writer voor tests
│   └── test_logger.go    # Test logger utilities
├── models/                # Data models
│   ├── aanmelding.go     # Aanmelding model
│   ├── contact.go        # Contact model
│   └── email.go          # Email model
├── services/              # Business logic
│   ├── email_batcher.go  # Email batch processing
│   ├── email_metrics.go  # Email metrics tracking
│   ├── email_service.go  # Email service implementatie
│   ├── interfaces.go     # Service interfaces
│   ├── prometheus_metrics.go # Prometheus integratie
│   ├── rate_limiter.go   # Rate limiting logica
│   └── smtp_client.go    # SMTP client implementatie
├── templates/             # Email templates
│   ├── aanmelding_admin_email.html
│   ├── aanmelding_email.html
│   ├── contact_admin_email.html
│   └── contact_email.html
├── tests/                 # Test suite
│   ├── aanmelding_handler_test.go
│   ├── elk_integration_test.go
│   ├── email_service_test.go
│   ├── mocks.go
│   └── ... (andere test bestanden)
└── main.go               # Application entrypoint
```

## Coding Standards

### Code Style
- Volg de officiële Go [style guide](https://golang.org/doc/effective_go)
- Gebruik `gofmt` voor code formatting
- Maximale functie lengte: 50 regels
- Betekenisvolle variabele namen in camelCase
- Package names: enkelvoud, kort en beschrijvend

### Naamgeving Conventies
```go
// Interfaces
type EmailSender interface {}

// Structs
type SMTPClient struct {}

// Constanten
const (
    MaxRetries = 3
    defaultTimeout = 10 * time.Second
)

// Variabelen
var (
    errInvalidInput = errors.New("invalid input")
    defaultConfig   = Config{}
)
```

### Error Handling
```go
// Correct
if err != nil {
    return fmt.Errorf("failed to send email: %w", err)
}

// Incorrect
if err != nil {
    log.Fatal(err) // Vermijd log.Fatal
}
```

### Logging
```go
// Gebruik structured logging
logger.Info("email sent",
    "type", emailType,
    "recipient", recipient,
    "duration", duration,
)

// Error logging met context
logger.Error("failed to send email",
    "type", emailType,
    "error", err,
    "retry_count", retryCount,
)
```

## Development Workflow

### Feature Development
1. Maak een nieuwe branch:
```bash
git checkout -b feature/nieuwe-feature
```

2. Implementeer de feature:
   - Schrijf eerst tests
   - Implementeer de functionaliteit
   - Voeg documentatie toe
   - Update API docs indien nodig

3. Verifieer de wijzigingen:
```bash
# Run tests
go test ./...

# Run linter
golangci-lint run

# Check coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

4. Commit de wijzigingen:
```bash
git add .
git commit -m "Voeg nieuwe feature toe

- Implementeer X functionaliteit
- Voeg tests toe
- Update documentatie"
```

### Code Review Checklist
- [ ] Tests geschreven en slagen
- [ ] Code gedocumenteerd
- [ ] Linter errors opgelost
- [ ] Error handling correct
- [ ] Logging toegevoegd
- [ ] Performance impact overwogen
- [ ] Security implications bekeken
- [ ] Breaking changes gedocumenteerd

## Debugging

### Met Delve
```bash
# Start debugger
dlv debug

# Set breakpoint
break main.go:42

# Run tot breakpoint
continue

# Inspect variables
print variableName
```

### Met Logs
```go
// Debug logging
logger.Debug("processing request",
    "method", r.Method,
    "path", r.URL.Path,
    "remote_addr", r.RemoteAddr,
)

// Trace logging voor gedetailleerde debugging
logger.Trace("smtp connection details",
    "host", smtpHost,
    "port", smtpPort,
    "tls", tlsEnabled,
)
```

### Performance Profiling
```bash
# CPU profiling
go test -cpuprofile=cpu.prof -bench=.

# Memory profiling
go test -memprofile=mem.prof -bench=.

# Analyze met pprof
go tool pprof cpu.prof
go tool pprof mem.prof
```

## Release Process

1. Version Bump
```bash
# Update version in
# - main.go
# - README.md
```

2. Create Release Branch
```bash
git checkout -b release/v1.2.0
```

3. Final Checks
```bash
# Run all tests
go test ./...

# Run linter
golangci-lint run

# Build binary
go build -ldflags="-s -w" -o dklemailservice
```

4. Create Release
```bash
git tag -a v1.2.0 -m "Release v1.2.0"
git push origin v1.2.0
```

## Best Practices

### Concurrency
```go
// Gebruik context voor cancellation
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Worker pools voor batch processing
pool := make(chan struct{}, maxWorkers)
for _, email := range emails {
    pool <- struct{}{}
    go func(e *models.Email) {
        defer func() { <-pool }()
        // Process email
    }(email)
}
```

### Resource Management
```go
// Gebruik defer voor cleanup
file, err := os.Open(filename)
if err != nil {
    return err
}
defer file.Close()

// SMTP connection management
conn, err := smtp.Dial(smtpHost)
if err != nil {
    return err
}
defer conn.Close()
```

### Configuration
```go
// Gebruik environment variables
type Config struct {
    SMTP struct {
        Host     string        `env:"SMTP_HOST,required"`
        Port     int           `env:"SMTP_PORT" default:"587"`
        Username string        `env:"SMTP_USER,required"`
        Password string        `env:"SMTP_PASSWORD,required"`
        From     string        `env:"SMTP_FROM,required"`
    }
    RateLimit struct {
        Global int `env:"GLOBAL_RATE_LIMIT" default:"100"`
        IP     int `env:"IP_RATE_LIMIT" default:"10"`
    }
    AllowedOrigins []string `env:"ALLOWED_ORIGINS" default:"https://www.dekoninklijkeloop.nl"`
}
```
