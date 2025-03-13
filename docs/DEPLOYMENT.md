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