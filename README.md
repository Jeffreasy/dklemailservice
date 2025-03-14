# DKL Email Service

Een robuuste en schaalbare email service voor De Koninklijke Loop, geschreven in Go. Deze service verzorgt alle email communicatie voor het evenement, inclusief aanmeldingen, contactformulieren en administratieve notificaties.

## 🌟 Functionaliteiten

- **Email Afhandeling**
  - Contactformulier emails met automatische bevestigingen
  - Aanmeldingsformulier emails met gepersonaliseerde content
  - Automatische bevestigingsmails met event-specifieke informatie
  - Admin notificaties voor nieuwe aanmeldingen en contactverzoeken
  - Ondersteuning voor HTML templates met dynamische content
  - Fallback naar plaintext voor betere deliverability

- **Authenticatie & Autorisatie** (Nieuw)
  - JWT-gebaseerde authenticatie voor beveiligde endpoints
  - Gebruikersbeheer met rollen (admin, gebruiker)
  - Wachtwoord hashing met bcrypt
  - Login rate limiting voor beveiliging
  - Beveiligde wachtwoord reset functionaliteit
  - HTTP-only cookies voor token opslag
  - Middleware voor rol-gebaseerde toegangscontrole

- **Beveiliging & Stabiliteit**
  - Rate limiting per IP en globaal voor spam preventie
  - CORS beveiliging met configureerbare origins
  - Graceful shutdown met cleanup van resources
  - Retry mechanisme voor failed emails met exponentiële backoff
  - Input validatie en sanitization
  - Secure SMTP configuratie met TLS support
  - XSS preventie in email templates

- **Monitoring & Observability**
  - Prometheus metrics voor real-time monitoring
  - ELK logging integratie voor centrale log aggregatie
  - Gedetailleerde email metrics per template en type
  - Health check endpoints met uitgebreide status informatie
  - Performance metrics voor email verzending
  - Rate limit statistieken
  - Error tracking en reporting

- **Performance**
  - Email batching voor efficiënte bulk verzending
  - Configureerbare rate limits per email type
  - Efficiënte template caching met auto-reload
  - Non-blocking email verzending met goroutines
  - Connection pooling voor SMTP verbindingen
  - Optimale resource utilizatie
  - Automatische cleanup van oude data

## 📋 Vereisten

- Go 1.21 of hoger
- SMTP server voor email verzending
  - Ondersteuning voor TLS
  - Voldoende verzendlimieten voor verwacht volume
- PostgreSQL 12 of hoger voor persistente opslag
  - Gebruiker met CREATE/ALTER/INSERT/UPDATE/DELETE rechten
  - Voldoende opslagruimte voor verwacht datavolume
- (Optioneel) SQLite voor lokale ontwikkeling en tests (vereist CGO)
- (Optioneel) ELK stack voor logging
  - Elasticsearch 7.x of hoger
  - Logstash voor log processing
  - Kibana voor visualisatie
- (Optioneel) Prometheus voor metrics
  - Prometheus server
  - Grafana voor dashboards

## 🚀 Installatie

1. Clone de repository:
```bash
git clone https://github.com/Jeffreasy/dklemailservice.git
cd dklemailservice
```

2. Installeer dependencies:
```bash
go mod download
go mod verify
```

3. Kopieer het voorbeeld configuratie bestand:
```bash
cp .env.example .env
```

4. Configureer de omgevingsvariabelen in `.env`:
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

# Database configuratie
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=dklemailservice
DB_SSL_MODE=disable

# JWT configuratie
JWT_SECRET=your_jwt_secret_key
JWT_TOKEN_EXPIRY=24h

# Rate Limiting
GLOBAL_RATE_LIMIT=1000
IP_RATE_LIMIT=50
RATE_LIMIT_WINDOW=1h
LOGIN_LIMIT_COUNT=5
LOGIN_LIMIT_PERIOD=300

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

## 🏃‍♂️ Gebruik

### Service Starten

Development mode:
```bash
go run main.go
```

Production mode:
```bash
go build -ldflags="-s -w" -o dklemailservice
./dklemailservice
```

### Tests Uitvoeren

Voor het uitvoeren van tests is CGO vereist vanwege SQLite afhankelijkheid in de tests. Gebruik de meegeleverde scripts:

#### Linux/macOS:
```bash
# Voer alle tests uit
./scripts/run_tests.sh

# Voer tests uit met coverage rapport
./scripts/run_tests.sh --coverage
```

#### Windows:
```batch
# Voer alle tests uit
scripts\run_tests.bat

# Voer tests uit met coverage rapport
scripts\run_tests.bat --coverage
```

Handmatig CGO inschakelen:
```bash
# Zet CGO aan voor SQLite ondersteuning
export CGO_ENABLED=1  # Linux/macOS
set CGO_ENABLED=1     # Windows

# Voer tests uit
go test ./tests/... -v
```

### API Endpoints

#### Health & Monitoring
- `GET /api/health` - Health check met uitgebreide service status
- `GET /api/metrics/email` - Gedetailleerde email statistieken
- `GET /api/metrics/rate-limits` - Rate limit status en statistieken
- `GET /metrics` - Prometheus metrics endpoint

#### Email Verzending
- `POST /api/contact-email` - Verstuur contact formulier
  ```json
  {
    "naam": "string",
    "email": "string",
    "bericht": "string",
    "privacy_akkoord": true
  }
  ```
- `POST /api/aanmelding-email` - Verstuur aanmelding formulier
  ```json
  {
    "naam": "string",
    "email": "string",
    "telefoon": "string",
    "rol": "string",
    "afstand": "string",
    "ondersteuning": "string",
    "bijzonderheden": "string",
    "terms": true
  }
  ```

#### Authenticatie (Nieuw)
- `POST /api/auth/login` - Gebruiker inloggen
  ```json
  {
    "email": "string",
    "wachtwoord": "string"
  }
  ```
- `POST /api/auth/logout` - Gebruiker uitloggen
- `GET /api/auth/profile` - Gebruikersprofiel ophalen (vereist authenticatie)
- `POST /api/auth/reset-password` - Wachtwoord wijzigen (vereist authenticatie)
  ```json
  {
    "huidig_wachtwoord": "string",
    "nieuw_wachtwoord": "string"
  }
  ```

#### Contact Beheer (Nieuw)
- `GET /api/contact` - Lijst van contactformulieren ophalen
- `GET /api/contact/:id` - Details van een specifiek contactformulier ophalen
- `PUT /api/contact/:id` - Contactformulier bijwerken (status, notities)
- `DELETE /api/contact/:id` - Contactformulier verwijderen
- `POST /api/contact/:id/antwoord` - Antwoord toevoegen aan contactformulier
- `GET /api/contact/status/:status` - Contactformulieren filteren op status

#### Aanmelding Beheer (Nieuw)
- `GET /api/aanmelding` - Lijst van aanmeldingen ophalen
- `GET /api/aanmelding/:id` - Details van een specifieke aanmelding ophalen
- `PUT /api/aanmelding/:id` - Aanmelding bijwerken (status, notities)
- `DELETE /api/aanmelding/:id` - Aanmelding verwijderen
- `POST /api/aanmelding/:id/antwoord` - Antwoord toevoegen aan aanmelding
- `GET /api/aanmelding/rol/:rol` - Aanmeldingen filteren op rol

### Docker

Build de Docker image:
```bash
# Development build
docker build -t dklemailservice:dev .

# Production build
docker build --build-arg GO_ENV=production -t dklemailservice:latest .
```

Run de container:
```bash
# Development
docker run -p 8080:8080 --env-file .env dklemailservice:dev

# Production
docker run -d --restart=always -p 8080:8080 --env-file .env dklemailservice:latest
```

Docker Compose setup:
```yaml
version: '3.8'
services:
  emailservice:
    build: .
    ports:
      - "8080:8080"
    env_file: .env
    restart: always
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/api/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

## 📊 Monitoring

### Prometheus Metrics

De service exporteert de volgende metrics:

#### Email Metrics
- `email_sent_total{type="contact|aanmelding",template="admin|user"}` - Aantal verzonden emails
- `email_failed_total{type="contact|aanmelding",error="smtp|template|validation"}` - Aantal gefaalde emails
- `email_latency_seconds{type="contact|aanmelding"}` - Email verzend latency
- `email_batch_size{type="contact|aanmelding"}` - Huidige batch grootte
- `email_template_render_duration_seconds` - Template render tijd

#### Rate Limiting
- `rate_limit_exceeded_total{type="ip|global"}` - Rate limit overschrijdingen
- `rate_limit_remaining{type="ip|global"}` - Resterende requests
- `rate_limit_reset_seconds` - Tijd tot rate limit reset

#### System Metrics
- `go_goroutines` - Aantal actieve goroutines
- `go_memory_alloc_bytes` - Geheugengebruik
- `process_cpu_seconds_total` - CPU gebruik

### Grafana Dashboard

Een voorgedefinieerd Grafana dashboard is beschikbaar in `./dashboards/email-service.json` met:
- Email verzending statistieken
- Rate limiting overzicht
- Performance metrics
- Error tracking
- System resource gebruik

### Logging

Logs worden geschreven in JSON formaat en kunnen worden doorgestuurd naar ELK:

#### Log Levels
- `DEBUG` - Gedetailleerde debug informatie
  - Template rendering details
  - SMTP communicatie
  - Rate limit checks
- `INFO` - Algemene operationele informatie
  - Email verzendingen
  - Service start/stop
  - Configuratie wijzigingen
- `WARN` - Waarschuwingen
  - Rate limit overschrijdingen
  - Template parsing issues
  - Connectie timeouts
- `ERROR` - Fouten die aandacht vereisen
  - SMTP fouten
  - Template fouten
  - Validatie fouten
- `FATAL` - Kritieke fouten die de service stoppen
  - Configuratie fouten
  - Port binding fouten
  - Database connectie fouten

#### Log Format
```json
{
  "level": "info",
  "timestamp": "2024-03-20T15:04:05Z",
  "caller": "email_service.go:42",
  "message": "Email verzonden",
  "email_type": "contact",
  "template": "admin",
  "duration_ms": 150,
  "success": true
}
```

## 🧪 Testen

### Unit Tests
```bash
# Run alle tests
go test ./... -v

# Run specifieke test package
go test ./services -v
go test ./handlers -v

# Run met race condition detection
go test -race ./...
```

### Coverage Tests
```bash
# Generate coverage report
go test ./... -coverprofile=coverage.out

# View coverage in browser
go tool cover -html=coverage.out

# Check coverage percentage
go tool cover -func=coverage.out
```

### Integration Tests
```bash
# Run integration tests
go test ./tests -tags=integration

# Run specific integration test
go test ./tests -run TestEmailFlow -tags=integration
```

### Load Tests
```bash
# Install k6
go install go.k6.io/k6@latest

# Run load tests
k6 run ./tests/load/email_load_test.js
```

## 📝 Email Templates

De service gebruikt HTML templates voor emails met de volgende features:
- Responsive design voor mobile devices
- Toegankelijk voor screen readers
- Ondersteuning voor verschillende email clients
- Dynamische content injectie
- Fallback plaintext versies

### Template Locaties
- `templates/contact_email.html` - Bevestiging voor contactformulier
- `templates/contact_admin_email.html` - Admin notificatie voor contactformulier
- `templates/aanmelding_email.html` - Bevestiging voor aanmelding
- `templates/aanmelding_admin_email.html` - Admin notificatie voor aanmelding

### Template Data
Beschikbare variabelen in templates:

#### Contact Templates
```go
type ContactData struct {
    Naam    string
    Email   string
    Bericht string
}
```

#### Aanmelding Templates
```go
type AanmeldingData struct {
    Naam           string
    Email          string
    Telefoon       string
    Rol            string
    Afstand        string
    Ondersteuning  string
    Bijzonderheden string
}
```

## 🔒 Rate Limiting

### Standaard Limieten
- Contact emails:
  - 100 emails per uur globaal
  - 5 emails per uur per IP
- Aanmelding emails:
  - 200 emails per uur globaal
  - 10 emails per uur per IP

### Configuratie
Rate limits kunnen worden aangepast via environment variables of runtime configuratie:

```go
rateLimiter.AddLimit("contact_email", 100, time.Hour, false)    // Globaal
rateLimiter.AddLimit("contact_email", 5, time.Hour, true)       // Per IP
rateLimiter.AddLimit("aanmelding_email", 200, time.Hour, false) // Globaal
rateLimiter.AddLimit("aanmelding_email", 10, time.Hour, true)   // Per IP
```

## 🛠 Architectuur

De service volgt een modulaire architectuur met de volgende componenten:

### Core Components
- `handlers/` - HTTP request handlers
  - `email_handler.go` - Email verzending endpoints
  - `health_handler.go` - Health check endpoint
  - `metrics_handler.go` - Metrics endpoints

- `services/` - Business logic
  - `email_service.go` - Email verzending logica
  - `smtp_client.go` - SMTP communicatie
  - `rate_limiter.go` - Rate limiting
  - `email_batcher.go` - Batch processing
  - `email_metrics.go` - Metrics tracking
  - `prometheus_metrics.go` - Prometheus integratie

- `models/` - Data structuren
  - `email.go` - Email gerelateerde structs
  - `contact.go` - Contact formulier model
  - `aanmelding.go` - Aanmelding formulier model

- `logger/` - Logging configuratie
  - `logger.go` - Logger setup
  - `elk_writer.go` - ELK integratie

- `templates/` - Email templates
  - HTML templates
  - Partials voor herbruikbare componenten

- `tests/` - Test suites
  - Unit tests
  - Integration tests
  - Load tests
  - Mocks

### Design Patterns
- Repository pattern voor data access
- Factory pattern voor service instantiatie
- Strategy pattern voor email verzending
- Observer pattern voor metrics
- Builder pattern voor email constructie

## 🗄️ Database Architectuur

De applicatie is uitgebreid met een robuuste PostgreSQL database integratie voor het persistent opslaan van gegevens. Deze integratie maakt gebruik van GORM als ORM (Object-Relational Mapping) framework en implementeert het Repository Pattern voor een schone scheiding van verantwoordelijkheden.

### Database Modellen

De volgende modellen zijn geïmplementeerd:

- **ContactFormulier**: Slaat contactformulier gegevens op met velden voor naam, email, bericht, status en behandeling.
- **ContactAntwoord**: Houdt antwoorden op contactformulieren bij, gekoppeld via een één-op-veel relatie.
- **Aanmelding**: Registreert aanmeldingen voor het evenement met persoonlijke gegevens en voorkeuren.
- **AanmeldingAntwoord**: Bewaart antwoorden op aanmeldingen, gekoppeld via een één-op-veel relatie.
- **EmailTemplate**: Slaat email templates op voor hergebruik en consistentie in communicatie.
- **VerzondEmail**: Houdt een log bij van alle verzonden emails voor tracking en auditing.
- **Gebruiker**: Beheert gebruikersaccounts voor administratieve toegang tot het systeem.
- **Migratie**: Houdt database migraties bij om schema-wijzigingen gecontroleerd uit te voeren.

### Repository Pattern

De applicatie implementeert het Repository Pattern voor data-toegang:

- **Basisrepository**: `PostgresRepository` biedt gemeenschappelijke functionaliteit zoals foutafhandeling en timeouts.
- **Gespecialiseerde repositories**: Voor elk model is er een specifieke repository die CRUD-operaties implementeert.
- **Repository Factory**: Centraliseert de creatie van repositories en zorgt voor eenvoudige dependency injection.

### Database Configuratie

De database verbinding wordt geconfigureerd via environment variables:

```env
# Database configuratie
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=dklemailservice
DB_SSL_MODE=disable
```

### Migraties

De applicatie ondersteunt automatische database migraties bij het opstarten:

1. Schema migraties: Creëert en update tabellen op basis van de gedefinieerde modellen.
2. Data seeding: Vult de database met initiële gegevens zoals standaard email templates.

### Mock Database voor Tests

Voor tests is een mock database implementatie beschikbaar die geen externe database vereist:

- In-memory opslag voor alle modellen
- Volledige implementatie van repository interfaces
- Automatische fallback naar mock database als SQLite niet beschikbaar is (CGO uitgeschakeld)

## 🐳 Docker Multi-stage Builds

De Dockerfile is verbeterd met multi-stage builds die twee versies van de applicatie produceren:

1. **Productie binary** (CGO uitgeschakeld):
   - Kleinere, statisch gelinkte binary
   - Betere performance en veiligheid
   - Geen ondersteuning voor SQLite (gebruikt PostgreSQL in productie)

2. **Ontwikkeling/test binary** (CGO ingeschakeld):
   - Ondersteunt SQLite voor lokale ontwikkeling en tests
   - Bevat debugging informatie
   - Geschikt voor testen met database-afhankelijke tests

De te gebruiken binary wordt bepaald door de `APP_ENV` environment variable:
```bash
# Productie mode (standaard)
docker run -e APP_ENV=prod dklemailservice

# Ontwikkeling/test mode met SQLite ondersteuning
docker run -e APP_ENV=dev dklemailservice
```

### Concurrency

De service gebruikt goroutines voor non-blocking operations en channels voor communicatie. Mutex voor thread-safe operations en context voor cancellation worden gebruikt om te voorkomen dat er race conditions optreden. Worker pools voor batch processing zorgen voor efficiënte bulk verzending.

## 👥 Bijdragen

1. Fork de repository
2. Maak een feature branch
```bash
git checkout -b feature/mijn-feature
```
3. Commit je wijzigingen
```bash
git commit -m 'Voeg nieuwe feature toe'
```
4. Push naar de branch
```bash
git push origin feature/mijn-feature
```
5. Open een Pull Request

### Development Guidelines
- Volg Go best practices en idioms
- Schrijf tests voor nieuwe functionaliteit
- Update documentatie waar nodig
- Voeg relevante logging toe
- Zorg voor adequate error handling
- Valideer input data
- Overweeg performance implicaties

## 📄 Licentie

Copyright (c) 2024 De Koninklijke Loop. Alle rechten voorbehouden.

Deze software is eigendom van De Koninklijke Loop en mag niet worden gebruikt, gekopieerd, gemodificeerd of gedistribueerd zonder uitdrukkelijke toestemming van De Koninklijke Loop.

## 📚 Documentatie

Uitgebreide documentatie is beschikbaar in de `/docs` directory:
- `API.md` - API documentatie
- `DEPLOYMENT.md` - Deployment instructies
- `DEVELOPMENT.md` - Development guidelines
- `MONITORING.md` - Monitoring setup
- `SECURITY.md` - Security best practices
- `TEMPLATES.md` - Template documentatie
- `TESTING.md` - Test procedures
- `AUTH.md` - Authenticatie documentatie (Nieuw)