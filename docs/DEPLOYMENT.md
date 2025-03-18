# Deployment Handleiding

## Vereisten

- Docker
- Go 1.21 of hoger (voor lokale development)
- SMTP server toegang voor email verzending
- SSL certificaten voor HTTPS
- (Optioneel) Prometheus voor metrics
- (Optioneel) ELK stack voor logging

## Lokale Development Setup

1. Clone de repository:
```bash
git clone https://github.com/Jeffreasy/dklemailservice.git
cd dklemailservice
```

2. Kopieer het voorbeeld configuratie bestand:
```bash
cp .env.example .env
```

3. Configureer de omgevingsvariabelen in `.env`:
```env
# Algemene SMTP configuratie
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=user@example.com
SMTP_PASSWORD=your_password
SMTP_FROM=noreply@example.com
SMTP_TLS_ENABLED=true
SMTP_TIMEOUT=10s

# Registratie SMTP configuratie
REGISTRATION_SMTP_HOST=smtp.example.com
REGISTRATION_SMTP_PORT=587
REGISTRATION_SMTP_USER=registration@example.com
REGISTRATION_SMTP_PASSWORD=your_password
REGISTRATION_SMTP_FROM=registration@example.com
REGISTRATION_SMTP_TLS_ENABLED=true
REGISTRATION_SMTP_TIMEOUT=10s

# Email adressen
ADMIN_EMAIL=admin@example.com
REGISTRATION_EMAIL=registration@example.com

# Rate Limiting
GLOBAL_RATE_LIMIT=1000
IP_RATE_LIMIT=50
RATE_LIMIT_WINDOW=1h

# Monitoring & Logging
LOG_LEVEL=info
LOG_FORMAT=json
ELK_ENDPOINT=http://elk:9200
ELK_INDEX=dkl-emails
ELK_BATCH_SIZE=100
PROMETHEUS_ENABLED=true

# Security
ALLOWED_ORIGINS=https://www.dekoninklijkeloop.nl,https://dekoninklijkeloop.nl
TLS_ENABLED=true
TLS_CERT_FILE=./certs/server.crt
TLS_KEY_FILE=./certs/server.key

# Performance
EMAIL_BATCH_SIZE=50
EMAIL_BATCH_INTERVAL=15m
TEMPLATE_RELOAD_INTERVAL=1h
MAX_CONCURRENT_SENDS=10
```

4. Start de applicatie:
```bash
go run main.go
```

## Docker Deployment

### Build

1. Build de Docker image:
```bash
# Development build
docker build -t dklemailservice:dev .

# Production build
docker build -t dklemailservice:latest .
```

De Dockerfile gebruikt een multi-stage build voor een geoptimaliseerde productie image:
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -o main -ldflags="-w -s" .

# Final stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy templates directory
COPY --from=builder /app/templates ./templates

# Expose port
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
```

### Run

Start de container met de juiste environment variables:
```bash
docker run -d \
  --name dklemailservice \
  --restart=always \
  -p 8080:8080 \
  --env-file .env \
  -v $(pwd)/certs:/app/certs \
  dklemailservice:latest
```

## Cloud Deployment (Render)

De service is geconfigureerd voor deployment op Render via `render.yaml`:

```yaml
services:
  - type: web
    name: dklautomatie-backend
    runtime: docker
    dockerfilePath: ./Dockerfile
    envVars:
      - key: ADMIN_EMAIL
        sync: false
      - key: SMTP_HOST
        sync: false
      - key: SMTP_PORT
        sync: false
      - key: SMTP_USER
        sync: false
      - key: SMTP_PASSWORD
        sync: false
      - key: ALLOWED_ORIGINS
        sync: false
      - key: SMTP_FROM
        sync: false
    healthCheckPath: /api/health
    buildCommand: docker build -t dklautomatie-backend .
    startCommand: docker run -p 8080:8080 dklautomatie-backend
```

## Monitoring Setup

### Prometheus

1. Configureer Prometheus target:
```yaml
scrape_configs:
  - job_name: 'dklemailservice'
    scrape_interval: 15s
    static_configs:
      - targets: ['localhost:8080']
```

2. Beschikbare metrics:
- `email_sent_total{type="contact|aanmelding"}`
- `email_failed_total{type="contact|aanmelding"}`
- `email_latency_seconds`
- `rate_limit_exceeded_total`

### ELK Stack

Als `ELK_ENDPOINT` is geconfigureerd, worden logs automatisch naar ELK verstuurd met de volgende configuratie:
- Index: `dkl-emails-YYYY.MM.dd`
- Batch grootte: 100 logs
- Flush interval: 5 seconden

## Health Checks

De service biedt een health check endpoint op `/api/health` die wordt gebruikt voor:
- Docker health checks
- Render health monitoring
- Load balancer checks

## Backup & Recovery

De service is stateless, dus er zijn geen database backups nodig. Belangrijke items om te backuppen zijn:
- SSL certificaten
- Environment configuratie
- Email templates

## Security Checklist

Voor deployment:
- [ ] Alle environment variables geconfigureerd
- [ ] SSL certificaten geïnstalleerd
- [ ] CORS origins correct ingesteld
- [ ] Rate limiting geconfigureerd
- [ ] Admin API key ingesteld
- [ ] Logging niveau ingesteld
- [ ] Firewall regels geconfigureerd

## Troubleshooting

### Common Issues

1. SMTP Connectie:
```bash
# Test SMTP verbinding
nc -zv $SMTP_HOST $SMTP_PORT
```

2. Rate Limiting:
```bash
# Check rate limit status
curl -H "X-API-Key: $ADMIN_API_KEY" https://api.dekoninklijkeloop.nl/api/metrics/rate-limits
```

3. Template Issues:
```bash
# Verifieer template directory
docker exec dklemailservice ls -la /app/templates
```

### Logs

Bekijk container logs:
```bash
# Live logs
docker logs -f dklemailservice

# Laatste 100 logs
docker logs --tail 100 dklemailservice
```

### Environment Variables

Configureer de volgende omgevingsvariabelen voor een productie deployment:

```bash
# Server configuratie
PORT=8080
ENV=production
LOG_LEVEL=info
BASE_URL=https://api.dekoninklijkeloop.nl

# Database configuratie
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-secure-password
DB_NAME=dklemailservice
DB_SSL_MODE=require

# SMTP configuratie
SMTP_HOST=smtp.provider.com
SMTP_PORT=587
SMTP_USERNAME=noreply@dekoninklijkeloop.nl
SMTP_PASSWORD=your-smtp-password
SMTP_USE_TLS=true
SMTP_FROM=noreply@dekoninklijkeloop.nl

# IMAP configuratie voor EmailAutoFetcher
EMAIL_FETCH_INTERVAL=15  # in minuten
DISABLE_AUTO_EMAIL_FETCH=false  # true om auto fetch uit te schakelen

# Info email account
EMAIL_INFO_HOST=imap.provider.com
EMAIL_INFO_PORT=993
EMAIL_INFO_USERNAME=info@dekoninklijkeloop.nl
EMAIL_INFO_PASSWORD=your-email-password
EMAIL_INFO_USE_TLS=true

# Inschrijving email account
EMAIL_INSCHRIJVING_HOST=imap.provider.com
EMAIL_INSCHRIJVING_PORT=993
EMAIL_INSCHRIJVING_USERNAME=inschrijving@dekoninklijkeloop.nl
EMAIL_INSCHRIJVING_PASSWORD=your-email-password
EMAIL_INSCHRIJVING_USE_TLS=true

# API security
JWT_SECRET=your-jwt-secret-key
API_KEY=your-admin-api-key
CORS_ALLOWED_ORIGINS=https://dekoninklijkeloop.nl,https://www.dekoninklijkeloop.nl
```

### Email Auto Fetcher Configuratie

De EmailAutoFetcher service is verantwoordelijk voor het automatisch ophalen van inkomende emails. Voor een robuuste productie-inzet volg deze best practices:

#### 1. Interval configuratie

Kies een geschikte interval voor email ophaal operaties:

```bash
# Standaard: 15 minuten
EMAIL_FETCH_INTERVAL=15

# Voor drukke mail accounts: elke 5 minuten
EMAIL_FETCH_INTERVAL=5

# Voor accounts met weinig verkeer: elk uur
EMAIL_FETCH_INTERVAL=60
```

#### 2. Mail account configuratie

Configureer de IMAP accounts voor elke mailbox die gemonitord moet worden:

```bash
# Info mailbox
EMAIL_INFO_HOST=imap.provider.com
EMAIL_INFO_PORT=993
EMAIL_INFO_USERNAME=info@dekoninklijkeloop.nl
EMAIL_INFO_PASSWORD=your-email-password
EMAIL_INFO_USE_TLS=true

# Aanmeldingen mailbox
EMAIL_INSCHRIJVING_HOST=imap.provider.com
EMAIL_INSCHRIJVING_PORT=993
EMAIL_INSCHRIJVING_USERNAME=inschrijving@dekoninklijkeloop.nl
EMAIL_INSCHRIJVING_PASSWORD=your-email-password
EMAIL_INSCHRIJVING_USE_TLS=true
```

#### 3. Uitschakelen in multi-instance omgevingen

Als meerdere instances van de applicatie worden uitgevoerd (bijv. in een Kubernetes cluster), configureer dan één designated instance voor email ophalen:

```bash
# Op designated master instance
DISABLE_AUTO_EMAIL_FETCH=false

# Op alle andere instances
DISABLE_AUTO_EMAIL_FETCH=true
```

#### 4. Monitoring configuratie

Voeg specifieke monitoring toe voor de EmailAutoFetcher:

```bash
# Log level voor EmailAutoFetcher (optioneel, overschrijft globale LOG_LEVEL)
EMAIL_FETCHER_LOG_LEVEL=info

# Alert drempelwaarde voor mislukte fetch operaties
EMAIL_FETCHER_ALERT_THRESHOLD=3
```

### Cron Job (Alternatief voor Auto Fetcher)

Als een externe trigger gewenst is in plaats van de ingebouwde EmailAutoFetcher, kan een cron job worden gebruikt:

```bash
# In crontab:
*/15 * * * * curl -X POST https://api.dekoninklijkeloop.nl/api/mail/fetch -H "Authorization: Bearer $JWT_TOKEN" > /dev/null 2>&1
```

Verkrijg de JWT token via:

```bash
JWT_TOKEN=$(curl -s -X POST https://api.dekoninklijkeloop.nl/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}' | jq -r '.token')
```

### Health Checks

// ... existing code ... 