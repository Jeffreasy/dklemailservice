# Setup Guide

Complete installation and setup guide for the DKL Email Service.

## Prerequisites

### Required Software
- **Go 1.23.0+** - [Download](https://golang.org/dl/)
- **PostgreSQL 17** - [Download](https://www.postgresql.org/download/)
- **Redis 7+** - [Download](https://redis.io/download)
- **Git** - [Download](https://git-scm.com/downloads)

### Optional Tools
- **Docker & Docker Compose** - For containerized development
- **Make** - For build automation
- **Postman** - For API testing

## Installation

### 1. Clone Repository

```bash
git clone https://github.com/yourusername/dklemailservice.git
cd dklemailservice
```

### 2. Install Dependencies

```bash
go mod download
go mod verify
```

### 3. Configure Environment

Copy the example environment file:

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
# Server Configuration
PORT=8080
ENV=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_secure_password
DB_NAME=dklemailservice
DB_SSLMODE=disable

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT Configuration
JWT_SECRET=your_very_secure_jwt_secret_key_min_32_chars
JWT_EXPIRATION=24h
JWT_REFRESH_EXPIRATION=168h

# Email Configuration - SMTP (Outgoing)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-specific-password
SMTP_FROM=noreply@dklemailservice.com

# Registration Email SMTP (Outgoing)
REGISTRATION_SMTP_HOST=smtp.gmail.com
REGISTRATION_SMTP_PORT=587
REGISTRATION_SMTP_USER=registration@example.com
REGISTRATION_SMTP_PASSWORD=your-app-specific-password
REGISTRATION_SMTP_FROM=registration@example.com

# Email Fetcher - IMAP (Incoming)
DISABLE_AUTO_EMAIL_FETCH=false
INFO_EMAIL=info@example.com
INFO_EMAIL_PASSWORD=your-password
INSCHRIJVING_EMAIL=registration@example.com
INSCHRIJVING_EMAIL_PASSWORD=your-password
IMAP_SERVER=imap.gmail.com
IMAP_PORT=993

# Email Addresses
ADMIN_EMAIL=admin@example.com
REGISTRATION_EMAIL=registration@example.com

# Cloudinary Configuration
CLOUDINARY_CLOUD_NAME=your_cloud_name
CLOUDINARY_API_KEY=your_api_key
CLOUDINARY_API_SECRET=your_api_secret

# Frontend URL (for CORS)
FRONTEND_URL=http://localhost:3000

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1h

# Logging
LOG_LEVEL=debug
LOG_FORMAT=json
```

### 4. Database Setup

#### Option A: Manual Setup

Create the database:

```bash
createdb -U postgres dklemailservice
```

Run migrations:

```bash
go run database/migrations/run_migrations.go
```

#### Option B: Using Docker

```bash
docker-compose up -d postgres redis
```

Wait for services to start, then run migrations:

```bash
go run database/migrations/run_migrations.go
```

### 5. Verify Setup

Run health check:

```bash
go run main.go
```

Visit http://localhost:8080/api/health - you should see:

```json
{
  "status": "healthy",
  "database": "connected",
  "redis": "connected",
  "timestamp": "2025-01-08T10:00:00Z"
}
```

## Development Setup

### Using Docker Compose (Recommended)

Complete development environment:

```bash
# Start all services
docker-compose -f docker-compose.dev.yml up

# Run in background
docker-compose -f docker-compose.dev.yml up -d

# View logs
docker-compose -f docker-compose.dev.yml logs -f

# Stop services
docker-compose -f docker-compose.dev.yml down
```

### Manual Development Setup

1. **Start PostgreSQL**
```bash
# Linux/macOS
sudo systemctl start postgresql

# macOS with Homebrew
brew services start postgresql

# Windows
net start postgresql-x64-15
```

2. **Start Redis**
```bash
# Linux/macOS
redis-server

# macOS with Homebrew
brew services start redis

# Windows
redis-server.exe
```

3. **Start Application**
```bash
go run main.go
```

Or use air for hot reloading:

```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

## Database Migrations

### Running Migrations

All migrations at once:
```bash
go run database/migrations/run_migrations.go
```

Using PowerShell script:
```powershell
.\database\apply_all_migrations.ps1
```

### Creating New Migrations

1. Create new migration file:
```bash
# Naming convention: V{number}__{description}.sql
touch database/migrations/V16__add_new_feature.sql
```

2. Write migration:
```sql
-- V16__add_new_feature.sql
CREATE TABLE new_table (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_new_table_name ON new_table(name);
```

3. Run migrations:
```bash
go run database/migrations/run_migrations.go
```

## Testing Setup

### Unit Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Integration Tests

```bash
# Run integration tests only
go test ./tests -tags=integration

# Run specific test
go test ./tests -run TestAuthFlow
```

### API Testing with Postman

1. Import Postman collection (if available)
2. Set environment variables:
   - `base_url`: http://localhost:8080
   - `access_token`: (obtained from login)

## IDE Setup

### Visual Studio Code

Recommended extensions:
- `golang.go` - Go language support
- `ms-azuretools.vscode-docker` - Docker support
- `humao.rest-client` - REST API testing
- `oderwat.indent-rainbow` - Better indentation

Settings (`.vscode/settings.json`):
```json
{
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "workspace",
  "go.formatTool": "goimports",
  "editor.formatOnSave": true,
  "go.useLanguageServer": true
}
```

### GoLand

1. Open project
2. Configure Go SDK (Settings → Go → GOROOT)
3. Enable Go Modules (Settings → Go → Go Modules)
4. Set environment variables in run configuration

## Troubleshooting

### Database Connection Issues

**Error**: `connection refused`

```bash
# Check PostgreSQL is running
sudo systemctl status postgresql

# Check connection
psql -U postgres -h localhost -d dklemailservice

# Verify credentials in .env file
```

**Error**: `database does not exist`

```bash
# Create database
createdb -U postgres dklemailservice

# Or using psql
psql -U postgres
CREATE DATABASE dklemailservice;
```

### Redis Connection Issues

**Error**: `connection refused`

```bash
# Start Redis
redis-server

# Check Redis is running
redis-cli ping
# Should return: PONG
```

### Port Already in Use

**Error**: `bind: address already in use`

```bash
# Find process using port 8080
# Linux/macOS
lsof -i :8080

# Windows
netstat -ano | findstr :8080

# Kill the process or change PORT in .env
```

### Migration Failures

**Error**: `migration already applied`

```bash
# Check migration status
psql -U postgres -d dklemailservice
SELECT version FROM schema_migrations;

# Rollback if needed (manual)
# Delete the problematic migration record
DELETE FROM schema_migrations WHERE version = 'V16__add_new_feature';
```

### Go Module Issues

**Error**: `cannot find package`

```bash
# Clean module cache
go clean -modcache

# Download dependencies

## Email Fetcher Configuration

De DKL Email Service kan automatisch inkomende emails ophalen van IMAP mailboxen.

### Email Fetcher Activeren

In `.env`:
```env
# Email fetcher uitschakelen (voor ontwikkeling)
DISABLE_AUTO_EMAIL_FETCH=true

# Email fetcher inschakelen (voor productie)
DISABLE_AUTO_EMAIL_FETCH=false

# IMAP configuratie
IMAP_SERVER=imap.hostnet.nl  # of imap.gmail.com voor Gmail
IMAP_PORT=993

# Email accounts met wachtwoorden
INFO_EMAIL=info@dekoninklijkeloop.nl
INFO_EMAIL_PASSWORD=your-password
INSCHRIJVING_EMAIL=inschrijving@dekoninklijkeloop.nl
INSCHRIJVING_EMAIL_PASSWORD=your-password
```

### Email Fetch Interval

Standaard worden emails elke **15 minuten** opgehaald. 

Logs bij ophalen:
```
"Automatisch ophalen van emails starten..."
"Email auto fetcher gestart","interval":"15m0s"
"IMAP Search resultaat","account":"info@...","aantal_uids":0
```

### Handmatig Emails Ophalen

Via API endpoint (vereist authenticatie):
```bash
curl -X POST http://localhost:8080/api/mail/fetch \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### IMAP Server Instellingen

**Hostnet:**
- Server: `imap.hostnet.nl`
- Port: `993`
- SSL: Required

**Gmail:**
- Server: `imap.gmail.com`
- Port: `993`
- App-specific password vereist (geen gewoon wachtwoord)

### Troubleshooting Email Fetcher

**Error**: `Invalid credentials (Failure)`
- Controleer IMAP wachtwoorden in `.env`
- Voor Gmail: gebruik App-specific password, geen normaal wachtwoord
- Controleer of IMAP_SERVER correct is (bijv. `imap.hostnet.nl`)

**Error**: `Geen nieuwe emails gevonden`
- Dit is geen error - betekent dat de inbox leeg is
- Email fetcher werkt correct

**Email fetcher draait niet:**
- Check of `DISABLE_AUTO_EMAIL_FETCH=false` in `.env`
- Restart service: `docker-compose restart app`
- Check logs: `docker-compose logs app | grep "fetcher"`

go mod download
go mod tidy

# Verify modules
go mod verify
```

### Environment Variables Not Loading

```bash
# Verify .env file exists
ls -la .env

# Check file format (no BOM, LF line endings)
# Restart application after changes

# Load manually for testing
export $(cat .env | xargs)
go run main.go
```

## Development Workflow

### Daily Workflow

1. **Pull latest changes**
```bash
git pull origin main
```

2. **Update dependencies**
```bash
go mod download
```

3. **Run migrations**
```bash
go run database/migrations/run_migrations.go
```

4. **Start services**
```bash
docker-compose -f docker-compose.dev.yml up -d
```

5. **Start application**
```bash
air  # with hot reload
# or
go run main.go
```

### Code Quality

**Format code:**
```bash
gofmt -w .
```

**Run linter:**
```bash
golangci-lint run
```

**Run tests:**
```bash
go test ./... -v
```

### Git Workflow

```bash
# Create feature branch
git checkout -b feature/new-feature

# Make changes
# ...

# Commit with meaningful message
git add .
git commit -m "feat: add new feature

- Implement feature X
- Add tests for feature X
- Update documentation"

# Push to remote
git push origin feature/new-feature

# Create pull request on GitHub
```

## Next Steps

- [Frontend Integration](./FRONTEND_INTEGRATION.md)
- [API Documentation](../api/README.md)
- [Deployment Guide](./DEPLOYMENT.md)
- [WebSocket Setup](./WEBSOCKET.md)

## Support

For issues during setup:
1. Check [Troubleshooting](#troubleshooting) section
2. Review [Common Issues](./COMMON_ISSUES.md)
3. Contact development team

---

Setup complete! Start building with `go run main.go` 🚀