# Deployment Guide

Complete deployment handleiding voor de DKL Email Service.

## Vereisten

### Software
- **Go**: 1.21 of hoger
- **Docker**: 20.10 of hoger (optioneel)
- **PostgreSQL**: 13 of hoger
- **Redis**: 6.0 of hoger (optioneel, voor rate limiting)

### Externe Services
- **SMTP Server**: Voor email verzending
- **IMAP Server**: Voor email ophaling
- **SSL Certificaten**: Voor HTTPS (productie)

## Environment Configuratie

### Basis Configuratie

Kopieer het voorbeeld bestand:
```bash
cp .env.example .env
```

### Verplichte Variabelen

**Database Configuratie:**
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-secure-password
DB_NAME=dklemailservice
DB_SSL_MODE=disable  # 'require' voor productie
```

**Implementatie:** [`config/database.go:26`](../../config/database.go:26)

**SMTP Configuratie (Standaard):**
```bash
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=noreply@dekoninklijkeloop.nl
SMTP_PASSWORD=your-smtp-password
SMTP_FROM=noreply@dekoninklijkeloop.nl
SMTP_TLS_ENABLED=true
SMTP_TIMEOUT=10s
```

**SMTP Configuratie (Registratie):**
```bash
REGISTRATION_SMTP_HOST=smtp.example.com
REGISTRATION_SMTP_PORT=587
REGISTRATION_SMTP_USER=inschrijving@dekoninklijkeloop.nl
REGISTRATION_SMTP_PASSWORD=your-smtp-password
REGISTRATION_SMTP_FROM=inschrijving@dekoninklijkeloop.nl
REGISTRATION_SMTP_TLS_ENABLED=true
REGISTRATION_SMTP_TIMEOUT=10s
```

**Email Adressen:**
```bash
ADMIN_EMAIL=info@dekoninklijkeloop.nl
REGISTRATION_EMAIL=inschrijving@dekoninklijkeloop.nl
```

**JWT Configuratie:**
```bash
JWT_SECRET=your-very-secure-secret-key-change-in-production
JWT_TOKEN_EXPIRY=20m
```

**Validatie:** [`main.go:28`](../../main.go:28)

### Optionele Variabelen

**Email Auto Fetcher:**
```bash
DISABLE_AUTO_EMAIL_FETCH=false
EMAIL_FETCH_INTERVAL=15  # Minuten

# Info account
INFO_EMAIL=info@dekoninklijkeloop.nl
INFO_EMAIL_PASSWORD=your-password

# Inschrijving account
INSCHRIJVING_EMAIL=inschrijving@dekoninklijkeloop.nl
INSCHRIJVING_EMAIL_PASSWORD=your-password

# IMAP Server
IMAP_SERVER=mail.hostnet.nl
IMAP_PORT=993
```

**Implementatie:** [`main.go:642`](../../main.go:642)

**Whisky for Charity (WFC):**
```bash
WFC_SMTP_HOST=arg-plplcl14.argewebhosting.nl
WFC_SMTP_PORT=465
WFC_SMTP_USER=noreply@whiskyforcharity.com
WFC_SMTP_PASSWORD=your-password
WFC_SMTP_FROM=noreply@whiskyforcharity.com
WFC_SMTP_SSL=true
```

**Redis:**
```bash
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

**Telegram Notificaties:**
```bash
ENABLE_NOTIFICATIONS=true
TELEGRAM_BOT_TOKEN=your-bot-token
TELEGRAM_CHAT_ID=your-chat-id
NOTIFICATION_THROTTLE=15m
NOTIFICATION_MIN_PRIORITY=medium
```

**Newsletter:**
```bash
ENABLE_NEWSLETTER=true
NEWSLETTER_SOURCES=https://example.com/rss
NEWSLETTER_FETCH_INTERVAL=24h
```

**Monitoring:**
```bash
LOG_LEVEL=info  # debug, info, warn, error
ELK_ENDPOINT=http://elk:9200
PROMETHEUS_ENABLED=true
```

**CORS:**
```bash
ALLOWED_ORIGINS=https://www.dekoninklijkeloop.nl,https://dekoninklijkeloop.nl,https://admin.dekoninklijkeloop.nl
```

## Lokale Development

### 1. Database Setup

**PostgreSQL installeren:**
```bash
# Windows (via Chocolatey)
choco install postgresql

# macOS
brew install postgresql

# Linux
sudo apt-get install postgresql
```

**Database aanmaken:**
```sql
CREATE DATABASE dklemailservice;
CREATE USER dkluser WITH PASSWORD 'your-password';
GRANT ALL PRIVILEGES ON DATABASE dklemailservice TO dkluser;
```

### 2. Redis Setup (Optioneel)

**Redis installeren:**
```bash
# Windows (via Chocolatey)
choco install redis-64

# macOS
brew install redis

# Linux
sudo apt-get install redis-server
```

**Redis starten:**
```bash
redis-server
```

### 3. Applicatie Starten

**Dependencies installeren:**
```bash
go mod download
```

**Applicatie starten:**
```bash
go run main.go
```

**Output:**
```
INFO DKL Email Service wordt gestart version=1.0.0
INFO Database configuratie geladen host=localhost port=5432
INFO Database verbinding succesvol
INFO Migraties uitgevoerd
INFO Email auto fetcher gestart
INFO Server gestart port=8080
```

**Implementatie:** [`main.go:109`](../../main.go:109)

## Docker Deployment

### Dockerfile

De service gebruikt een multi-stage build voor optimale image grootte.

**Implementatie:** [`Dockerfile`](../../Dockerfile:1)

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main -ldflags="-w -s" .

# Final stage
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates

EXPOSE 8080

CMD ["./main"]
```

### Build Image

```bash
# Development build
docker build -t dklemailservice:dev .

# Production build
docker build -t dklemailservice:latest .
```

### Run Container

**Met environment file:**
```bash
docker run -d \
  --name dklemailservice \
  --restart=always \
  -p 8080:8080 \
  --env-file .env \
  dklemailservice:latest
```

**Met individuele variabelen:**
```bash
docker run -d \
  --name dklemailservice \
  --restart=always \
  -p 8080:8080 \
  -e DB_HOST=postgres \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=password \
  -e DB_NAME=dklemailservice \
  -e SMTP_HOST=smtp.example.com \
  -e SMTP_USER=user@example.com \
  -e SMTP_PASSWORD=password \
  dklemailservice:latest
```

### Docker Compose

**Implementatie:** [`docker-compose.yml`](../../docker-compose.yml:1)

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    env_file:
      - .env
    depends_on:
      - postgres
      - redis
    restart: always

  postgres:
    image: postgres:13-alpine
    environment:
      POSTGRES_DB: dklemailservice
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  redis:
    image: redis:6-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

volumes:
  postgres_data:
  redis_data:
```

**Starten:**
```bash
docker-compose up -d
```

**Logs bekijken:**
```bash
docker-compose logs -f app
```

**Stoppen:**
```bash
docker-compose down
```

## Cloud Deployment (Render)

### Render Configuratie

**Implementatie:** [`render.yaml`](../../render.yaml:1)

```yaml
services:
  - type: web
    name: dklautomatie-backend
    runtime: docker
    dockerfilePath: ./Dockerfile
    envVars:
      - key: DB_HOST
        sync: false
      - key: DB_PORT
        sync: false
      - key: DB_USER
        sync: false
      - key: DB_PASSWORD
        sync: false
      - key: DB_NAME
        sync: false
      - key: SMTP_HOST
        sync: false
      - key: SMTP_USER
        sync: false
      - key: SMTP_PASSWORD
        sync: false
      - key: JWT_SECRET
        sync: false
      - key: ALLOWED_ORIGINS
        sync: false
    healthCheckPath: /api/health
```

### Deployment Steps

1. **Connect Repository:**
   - Link GitHub repository in Render dashboard
   - Selecteer branch (main/production)

2. **Configure Environment:**
   - Voeg alle environment variabelen toe in Render dashboard
   - Gebruik Render's PostgreSQL database

3. **Deploy:**
   - Render bouwt en deployt automatisch bij push
   - Monitor deployment logs in dashboard

### Render PostgreSQL

**Productie Configuratie:** [`config/database.go:29`](../../config/database.go:29)

```go
if appEnv := os.Getenv("APP_ENV"); appEnv == "prod" {
    // Render PostgreSQL configuratie
    possibleHosts := []string{
        "dpg-cva4c01c1ekc738q6q0g-a",                            // Interne hostname
        "dpg-cva4c01c1ekc738q6q0g-a.oregon-postgres.render.com", // Externe hostname
    }
    
    // Probeer beide hostnamen
    for _, host := range possibleHosts {
        db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err == nil {
            return db, nil
        }
    }
}
```

## Production Checklist

### Pre-Deployment

- [ ] Alle environment variabelen geconfigureerd
- [ ] Database migraties getest
- [ ] SSL certificaten geïnstalleerd
- [ ] CORS origins correct ingesteld
- [ ] Rate limiting geconfigureerd
- [ ] JWT secret ingesteld (sterk wachtwoord)
- [ ] SMTP credentials gevalideerd
- [ ] IMAP credentials gevalideerd (indien auto-fetch enabled)
- [ ] Logging niveau ingesteld (info of warn)
- [ ] Health check endpoint getest
- [ ] Backup strategie bepaald

### Post-Deployment

- [ ] Health check monitoren
- [ ] Logs controleren op errors
- [ ] Email verzending testen
- [ ] Email ophaling testen (indien enabled)
- [ ] Authenticatie testen
- [ ] Rate limiting verificeren
- [ ] Metrics endpoint controleren
- [ ] Telegram notificaties testen (indien enabled)

## Database Migraties

### Automatische Migraties

Migraties worden automatisch uitgevoerd bij startup.

**Implementatie:** [`main.go:196`](../../main.go:196)

```go
migrationManager := database.NewMigrationManager(db, repoFactory.Migratie)
if err := migrationManager.MigrateDatabase(); err != nil {
    logger.Fatal("Database migratie fout", "error", err)
}

if err := migrationManager.SeedDatabase(); err != nil {
    logger.Fatal("Database seeding fout", "error", err)
}
```

### Migratie Bestanden

Locatie: [`database/migrations/`](../../database/migrations/)

**Naming Convention:**
```
V{version}__{description}.sql
```

**Voorbeelden:**
- `V1_20__create_rbac_tables.sql`
- `V1_21__seed_rbac_data.sql`
- `V1_28__add_refresh_tokens.sql`

### Handmatige Migratie

Voor productie updates zonder restart:

```bash
# Via psql
psql -h localhost -U postgres -d dklemailservice -f database/migrations/V1_XX__new_migration.sql
```

## Monitoring Setup

### Health Checks

**Endpoint:** `/api/health`

**Implementatie:** [`handlers/health_handler.go`](../../handlers/health_handler.go:1)

**Monitoring Tools:**
- Render: Automatische health checks
- Uptime Robot: Externe monitoring
- Prometheus: Metrics scraping

### Prometheus

**Metrics Endpoint:** `/metrics`

**Scrape Configuratie:**
```yaml
scrape_configs:
  - job_name: 'dklemailservice'
    scrape_interval: 15s
    static_configs:
      - targets: ['api.dekoninklijkeloop.nl:8080']
    metrics_path: '/metrics'
    scheme: 'https'
```

### Logging

**Log Levels:**
- `debug` - Gedetailleerde debugging informatie
- `info` - Algemene informatie (aanbevolen voor productie)
- `warn` - Waarschuwingen
- `error` - Errors
- `fatal` - Kritieke fouten (applicatie stopt)

**ELK Stack Integratie:**
```bash
ELK_ENDPOINT=http://elk:9200
```

**Implementatie:** [`main.go:152`](../../main.go:152)

```go
if elkEndpoint := os.Getenv("ELK_ENDPOINT"); elkEndpoint != "" {
    logger.SetupELK(logger.ELKConfig{
        Endpoint:      elkEndpoint,
        BatchSize:     100,
        FlushInterval: 5 * time.Second,
        AppName:       "dklemailservice",
        Environment:   os.Getenv("ENVIRONMENT"),
    })
}
```

## Email Auto Fetcher Deployment

### Configuratie

**Single Instance:**
```bash
DISABLE_AUTO_EMAIL_FETCH=false
EMAIL_FETCH_INTERVAL=15
```

**Multi-Instance (Kubernetes/Load Balancer):**

**Master Instance:**
```bash
DISABLE_AUTO_EMAIL_FETCH=false
EMAIL_FETCH_INTERVAL=15
```

**Worker Instances:**
```bash
DISABLE_AUTO_EMAIL_FETCH=true
```

**Implementatie:** [`main.go:274`](../../main.go:274)

```go
if os.Getenv("DISABLE_AUTO_EMAIL_FETCH") != "true" {
    logger.Info("Automatisch ophalen van emails starten...")
    serviceFactory.EmailAutoFetcher.Start()
    logger.Info("Automatische email fetcher gestart")
} else {
    logger.Info("Automatisch ophalen van emails is uitgeschakeld")
}
```

### IMAP Accounts

**Info Account:**
```bash
INFO_EMAIL=info@dekoninklijkeloop.nl
INFO_EMAIL_PASSWORD=your-password
```

**Inschrijving Account:**
```bash
INSCHRIJVING_EMAIL=inschrijving@dekoninklijkeloop.nl
INSCHRIJVING_EMAIL_PASSWORD=your-password
```

**Server Configuratie:**
```bash
IMAP_SERVER=mail.hostnet.nl
IMAP_PORT=993
```

## Graceful Shutdown

De applicatie ondersteunt graceful shutdown voor veilig afsluiten.

**Implementatie:** [`main.go:599`](../../main.go:599)

```go
// Wacht op interrupt signaal
stop := make(chan os.Signal, 1)
signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
<-stop

logger.Info("Server wordt afgesloten...")

// Stop services
if serviceFactory.EmailBatcher != nil {
    serviceFactory.EmailBatcher.Shutdown()
}

if serviceFactory.EmailAutoFetcher != nil && serviceFactory.EmailAutoFetcher.IsRunning() {
    serviceFactory.EmailAutoFetcher.Stop()
}

if serviceFactory.NewsletterService != nil {
    serviceFactory.NewsletterService.Stop()
}

if rateLimiter != nil {
    rateLimiter.Shutdown()
}

// Log laatste metrics
serviceFactory.EmailMetrics.LogMetrics()

// Sluit log writers
logger.CloseWriters()

// Shutdown Fiber
app.Shutdown()
```

## Security Hardening

### SSL/TLS

**Productie:**
```bash
TLS_ENABLED=true
TLS_CERT_FILE=/path/to/cert.pem
TLS_KEY_FILE=/path/to/key.pem
```

**Let's Encrypt (aanbevolen):**
```bash
certbot certonly --standalone -d api.dekoninklijkeloop.nl
```

### Firewall

**Poorten:**
- 8080: HTTP/HTTPS (applicatie)
- 5432: PostgreSQL (alleen intern)
- 6379: Redis (alleen intern)

**iptables Voorbeeld:**
```bash
# Allow HTTP/HTTPS
sudo iptables -A INPUT -p tcp --dport 8080 -j ACCEPT

# Block direct database access
sudo iptables -A INPUT -p tcp --dport 5432 -s 127.0.0.1 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 5432 -j DROP
```

### Environment Security

**Nooit committen:**
- `.env` bestanden
- SSL certificaten
- Database credentials
- API keys

**Gebruik:**
- Environment variabelen in productie
- Secrets management (Render, AWS Secrets Manager)
- Encrypted backups

## Backup & Recovery

### Database Backup

**Automatische Backup:**
```bash
# Cron job (dagelijks om 2:00 AM)
0 2 * * * pg_dump -h localhost -U postgres dklemailservice > /backups/dkl_$(date +\%Y\%m\%d).sql
```

**Handmatige Backup:**
```bash
pg_dump -h localhost -U postgres dklemailservice > backup.sql
```

**Restore:**
```bash
psql -h localhost -U postgres dklemailservice < backup.sql
```

### Configuration Backup

Backup belangrijke bestanden:
- `.env` (encrypted)
- SSL certificaten
- Email templates
- Migratie bestanden

## Troubleshooting

### Database Connection Issues

**Check connectie:**
```bash
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME
```

**Logs bekijken:**
```bash
# Docker
docker logs dklemailservice

# Systemd
journalctl -u dklemailservice -f
```

### SMTP Issues

**Test SMTP verbinding:**
```bash
telnet $SMTP_HOST $SMTP_PORT
```

**Test met curl:**
```bash
curl -v --url "smtp://$SMTP_HOST:$SMTP_PORT" \
  --mail-from "test@example.com" \
  --mail-rcpt "recipient@example.com" \
  --user "$SMTP_USER:$SMTP_PASSWORD"
```

### Email Auto Fetcher Issues

**Check IMAP verbinding:**
```bash
telnet $IMAP_SERVER $IMAP_PORT
```

**Handmatig fetch triggeren:**
```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@dekoninklijkeloop.nl","wachtwoord":"password"}' | jq -r '.token')

curl -X POST http://localhost:8080/api/mail/fetch \
  -H "Authorization: Bearer $TOKEN"
```

### Rate Limiting Issues

**Check rate limits:**
```bash
curl -H "X-API-Key: $ADMIN_API_KEY" \
  https://api.dekoninklijkeloop.nl/api/metrics/rate-limits
```

**Reset rate limits (Redis):**
```bash
redis-cli FLUSHDB
```

## Performance Tuning

### Database Connection Pool

**Configuratie:** [`config/database.go:210`](../../config/database.go:210)

```go
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

### Email Batching

```bash
EMAIL_BATCH_SIZE=50
EMAIL_BATCH_INTERVAL=15m
MAX_CONCURRENT_SENDS=10
```

### Template Caching

Templates worden bij startup geladen en in memory gecached.

**Implementatie:** [`services/email_service.go:45`](../../services/email_service.go:45)

## Scaling

### Horizontal Scaling

**Load Balancer Setup:**
- Meerdere app instances
- Shared PostgreSQL database
- Shared Redis voor rate limiting
- Eén instance voor EmailAutoFetcher

**Nginx Load Balancer:**
```nginx
upstream dkl_backend {
    server app1:8080;
    server app2:8080;
    server app3:8080;
}

server {
    listen 443 ssl;
    server_name api.dekoninklijkeloop.nl;
    
    location / {
        proxy_pass http://dkl_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### Vertical Scaling

**Resource Requirements:**
- **Minimum:** 512MB RAM, 1 CPU
- **Aanbevolen:** 1GB RAM, 2 CPU
- **High Traffic:** 2GB RAM, 4 CPU

## Zie Ook

- [Development Guide](./development.md) - Ontwikkelomgeving
- [Security Guide](./security.md) - Security best practices
- [Monitoring Guide](./monitoring.md) - Monitoring setup
- [Testing Guide](./testing.md) - Test procedures