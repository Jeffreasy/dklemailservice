# DKL Email Service

## 📋 Overzicht

De DKL Email Service verzorgt betrouwbare en schaalbare e-mailcommunicatie voor de Koninklijke Loop website. De service is gebouwd met Go en biedt uitgebreide functionaliteit voor het versturen van diverse soorten e-mails met ingebouwde bescherming, monitoring en metrics.

## ✨ Functionaliteit

- **Template-gebaseerde e-mails**: HTML templates voor professionele en consistente communicatie
- **Administratieve notificaties**: Automatische meldingen bij nieuwe contactaanvragen en inschrijvingen
- **Rate limiting**: Bescherming tegen misbruik op globaal en per-IP niveau
- **Retry mechanisme**: Automatische nieuwe pogingen bij tijdelijke fouten
- **Batch verwerking**: Efficiënte verzending van grote aantallen e-mails
- **Gestructureerde logging**: JSON-logs met verschillende detailniveaus en ELK-integratie
- **Metrics tracking**: Statistieken over verzonden e-mails en prestaties
- **Prometheus integratie**: Real-time monitoring en alerting
- **Uitgebreide testsuite**: Unit en integratietests voor alle componenten

## 🏗️ Architectuur

De service is modulair en volgens moderne software-architectuurprincipes opgebouwd:

### Core Componenten

- **Email Service**: Centrale component die e-mailverzending coördineert
- **SMTP Client**: Verzorgt de daadwerkelijke communicatie met SMTP-servers
- **Rate Limiter**: Voorkomt overbelasting en misbruik
- **Email Metrics**: Verzamelt statistieken over verzonden e-mails
- **Email Batcher**: Groepeert e-mails voor efficiënte verwerking
- **Prometheus Metrics**: Real-time monitoring van alle services

## ⚙️ Configuratie

De service wordt geconfigureerd via omgevingsvariabelen:

### Vereiste instellingen
```
SMTP_HOST=smtp.example.com
SMTP_USER=username
SMTP_PASSWORD=password
SMTP_FROM=noreply@example.com
ADMIN_EMAIL=info@dekoninklijkeloop.nl        # Voor contactformulieren
REGISTRATION_EMAIL=inschrijving@dekoninklijkeloop.nl  # Voor aanmeldingen
```

### Optionele instellingen
```
PORT=8080                        # Serverpoort (standaard 8080)
SMTP_PORT=587                    # SMTP poort (standaard 587)
EMAIL_RATE_LIMIT=10              # Aantal e-mails per minuut (standaard 10)
TEMPLATE_DIR=./templates         # Map met e-mail templates 
LOG_LEVEL=info                   # Log niveau (debug, info, warn, error)
ELK_ENDPOINT=http://elk:9200     # ELK stack endpoint voor logging
ALLOWED_ORIGINS=https://www.dekoninklijkeloop.nl,https://dekoninklijkeloop.nl
```

## 🚀 Installatie en gebruik

### Vereisten
- Go 1.21 of hoger
- SMTP server toegang
- Prometheus (optioneel voor monitoring)

### Setup
1. Clone de repository
2. Kopieer `.env.example` naar `.env` en configureer:
    ```
    cp .env.example .env
    ```
3. Installeer dependencies:
    ```
    go mod download
    ```
4. Start de service:
    ```
    go run main.go
    ```

### Docker
```
docker build -t dklemailservice .
docker run -p 8080:8080 --env-file .env dklemailservice
```

## 📡 API Endpoints

### Email endpoints
- **POST /api/contact-email**: Verzendt een contact e-mail
- **POST /api/aanmelding-email**: Verzendt een aanmeldings e-mail

### Monitoring endpoints
- **GET /api/health**: Service status en gezondheidscheck
- **GET /api/metrics/email**: E-mail statistieken (verzonden, mislukt, etc.)
- **GET /api/metrics/rate-limits**: Informatie over rate limiting status
- **GET /metrics**: Prometheus metrics endpoint

## 📊 Monitoring

### Email Metrics
De service houdt interne statistieken bij via het `EmailMetrics` subsysteem:
- Aantal verzonden e-mails per type
- Aantal mislukte verzendingen
- Succespercentage

### Prometheus Metrics
Voor geavanceerde monitoring biedt de service Prometheus metrics:

```
# Belangrijkste beschikbare metrics
email_service_emails_sent_total{type="contact_email",template="contact_admin"}
email_service_emails_failed_total{type="aanmelding_email",error_type="connection"}
email_service_latency_seconds{type="generic"}
email_service_rate_limit_exceeded_total{type="contact_email",limit_type="per_user"}
email_service_active_batches
```

### Prometheus configuratie
Voeg de volgende configuratie toe aan je Prometheus `prometheus.yml`:

```
scrape_configs:
  - job_name: 'email-service'
    scrape_interval: 15s
    static_configs:
      - targets: ['localhost:8080']
```

### Grafana Dashboard
Voor visualisatie kun je Grafana gebruiken met dashboards voor:
- Email verzending per tijdseenheid
- Success/failure ratio's
- Latency distributies
- Rate limit overtredingen

## 🔒 Rate Limiting

De service implementeert geavanceerde rate limiting:

- **Globale limieten**: Beperkt het totale aantal e-mails
- **Per-IP limieten**: Voorkomt misbruik vanaf specifieke IP-adressen
- **Configureerbare tijdsperiodes**: Limieten per minuut/uur/dag

Voorbeeld configuratie:
```
// 100 contact emails per uur globaal
rateLimiter.AddLimit("contact_email", 100, time.Hour, false)    
// 5 contact emails per uur per IP
rateLimiter.AddLimit("contact_email", 5, time.Hour, true)       
```

## 📝 Logging

De service gebruikt gestructureerde JSON-logging met meerdere niveaus:

| Niveau | Beschrijving |
|--------|--------------|
| debug  | Gedetailleerde debug informatie |
| info   | Algemene operationele informatie (standaard) |
| warn   | Waarschuwingen die aandacht kunnen vereisen |
| error  | Fouten die normale werking verstoren |

```
{"niveau":"INFO","caller":"main.go:50","bericht":"DKL Email Service wordt gestart","tijd":"2023-07-01T15:04:05Z","version":"1.0.0"}
```

De logs kunnen naar een ELK stack worden gestuurd voor geavanceerde analyse.

## 🧪 Testen

```
# Alle tests uitvoeren
go test ./... -v

# Specifieke test uitvoeren
go test -v ./tests -run TestEmailRateLimiting

# Test coverage rapportage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🔄 CI/CD

De service is voorbereid voor continuous integration via GitHub Actions:
- Automatische tests bij elke pull request
- Linting en code kwaliteitscontrole
- Automatische Docker builds

## 📦 Deployment

### Render
1. Maak een nieuwe Web Service aan op Render
2. Verbind met je GitHub repository
3. Kies "Docker" als runtime
4. Configureer de environment variables
5. Deploy! 