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

- **Authenticatie & Autorisatie** (Geïmplementeerd)
  - JWT-gebaseerde authenticatie voor beveiligde endpoints
  - Gebruikersbeheer met rollen (admin, gebruiker)
  - Wachtwoord hashing met bcrypt
  - Login rate limiting voor beveiliging
  - Beveiligde wachtwoord reset functionaliteit
  - HTTP-only cookies voor token opslag
  - Middleware voor rol-gebaseerde toegangscontrole

- **Contact & Aanmelding Beheer** (Geïmplementeerd)
  - Beheer van contactformulieren (lijst, details, bijwerken, verwijderen)
  - Beheer van aanmeldingen (lijst, details, bijwerken, verwijderen)
  - Antwoorden toevoegen aan contactformulieren en aanmeldingen
  - Filteren op status (contactformulieren) en rol (aanmeldingen)
  - Automatische email notificaties bij antwoorden
  - Status tracking van contactformulieren en aanmeldingen
  - Notities toevoegen voor interne communicatie

- **Chat Systeem** (Geïmplementeerd)
  - Real-time messaging met WebSocket ondersteuning
  - Publieke, private en directe chat kanalen
  - Channel beheer met rollen (owner, admin, member)
  - Bericht reacties en bewerkingen
  - Gebruikers aanwezigheid en online status
  - Typing indicators en read receipts
  - Channel deelnemers beheer

- **Role-Based Access Control (RBAC)** (Geïmplementeerd)
  - Granulaire permissie systeem gebaseerd op rollen
  - Redis-gebaseerde caching voor optimale performance
  - Dynamische rol- en permissie toewijzing
  - Resource-based toegang controle (contact, aanmelding, newsletter, etc.)
  - Permission inheritance via rol hiërarchie
  - Audit logging van permissie wijzigingen

- **Nieuwsbrief Beheer** (Geïmplementeerd)
  - Automatische nieuwsbrief generatie van RSS feeds
  - HTML template ondersteuning voor professionele layouts
  - Batch verzending met rate limiting
  - Subscriber beheer met opt-in/opt-out functionaliteit
  - Verzendgeschiedenis en statistieken
  - Meerdere nieuwsbronnen ondersteuning

- **Gebruikersbeheer** (Geïmplementeerd)
  - Uitgebreide gebruikersprofielen
  - Rol en permissie beheer
  - Wachtwoord reset en account recovery
  - Gebruikersactiviteit monitoring
  - Bulk operaties voor gebruikersbeheer

- **Beveiliging & Stabiliteit**
  - Rate limiting per IP en globaal voor spam preventie
  - CORS beveiliging met configureerbare origins
  - Graceful shutdown met cleanup van resources
  - Retry mechanisme voor failed emails met exponentiële backoff
  - Input validatie en sanitization
  - Secure SMTP configuratie met TLS support
  - XSS preventie in email templates

- **Monitoring & Observability** (Geïmplementeerd)
  - Prometheus metrics voor real-time monitoring
  - ELK logging integratie voor centrale log aggregatie
  - Gedetailleerde email metrics per template en type
  - Health check endpoints met uitgebreide status informatie
  - Performance metrics voor email verzending
  - Rate limit statistieken
  - API key authenticatie voor metrics endpoints
  - Error tracking en reporting

- **Performance**
  - Email batching voor efficiënte bulk verzending
  - Configureerbare rate limits per email type
  - Efficiënte template caching met auto-reload
  - Non-blocking email verzending met goroutines
  - Connection pooling voor SMTP verbindingen
  - Optimale resource utilizatie
  - Automatische cleanup van oude data
  - Automatische email ophaling met configureerbaar interval (Geïmplementeerd)

## 📋 Vereisten

- Go 1.21 of hoger
- SMTP server voor email verzending
  - Ondersteuning voor TLS
  - Voldoende verzendlimieten voor verwacht volume
- PostgreSQL 12 of hoger voor persistente opslag
  - Gebruiker met CREATE/ALTER/INSERT/UPDATE/DELETE rechten
  - Voldoende opslagruimte voor verwacht datavolume
- (Optioneel) SQLite voor lokale ontwikkeling en tests (vereist CGO)
- (Optioneel) Redis voor caching en real-time features
  - Redis 6.x of hoger
  - Vereist voor RBAC permissie caching en chat presence
  - Verbeterd performance voor rate limiting
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

# Redis Configuratie
# Redis is VEREIST voor:
# - RBAC permissie caching (optimale performance)
# - Chat presence en typing indicators
# - Rate limiting (Redis-backed voor productie)
# - Session management
REDIS_ENABLED=true
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
# Alternative: Use REDIS_URL for cloud providers (e.g., Render)
# REDIS_URL=redis://username:password@host:port/db

# Automatische Email Ophaling
EMAIL_FETCH_INTERVAL=15
DISABLE_AUTO_EMAIL_FETCH=false
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

#### Health & Monitoring (Geïmplementeerd)
- `GET /api/health` - Health check met uitgebreide service status
- `GET /api/metrics/email` - Gedetailleerde email statistieken (vereist API key)
- `GET /api/metrics/rate-limits` - Rate limit status en statistieken (vereist API key)
- `GET /metrics` - Prometheus metrics endpoint (vereist API key)

#### Email Verzending (Geïmplementeerd)
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

#### Authenticatie (Geïmplementeerd)
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

#### Contact Beheer (Geïmplementeerd)
- `GET /api/contact` - Lijst van contactformulieren ophalen
- `GET /api/contact/:id` - Details van een specifiek contactformulier ophalen
- `PUT /api/contact/:id` - Contactformulier bijwerken (status, notities)
- `DELETE /api/contact/:id` - Contactformulier verwijderen
- `POST /api/contact/:id/antwoord` - Antwoord toevoegen aan contactformulier
- `GET /api/contact/status/:status` - Contactformulieren filteren op status

#### Aanmelding Beheer (Geïmplementeerd)
- `GET /api/aanmelding` - Lijst van aanmeldingen ophalen
- `GET /api/aanmelding/:id` - Details van een specifieke aanmelding ophalen
- `PUT /api/aanmelding/:id` - Aanmelding bijwerken (status, notities)
- `DELETE /api/aanmelding/:id` - Aanmelding verwijderen
- `POST /api/aanmelding/:id/antwoord` - Antwoord toevoegen aan aanmelding
- `GET /api/aanmelding/rol/:rol` - Aanmeldingen filteren op rol

#### Mail Beheer (Geïmplementeerd)
- `GET /api/mail` - Lijst van inkomende emails ophalen
- `GET /api/mail/:id` - Details van een specifieke email ophalen
- `PUT /api/mail/:id/processed` - Email markeren als verwerkt
- `DELETE /api/mail/:id` - Email verwijderen
- `POST /api/mail/fetch` - Handmatig nieuwe emails ophalen
- `GET /api/mail/unprocessed` - Lijst van onverwerkte emails ophalen
- `GET /api/mail/account/:type` - Emails filteren op account type (info, inschrijving)

#### Chat Systeem (Geïmplementeerd)
- `GET /api/chat/channels` - Lijst van gebruikers kanalen ophalen
- `GET /api/chat/channels/:id/participants` - Deelnemers van een kanaal ophalen
- `GET /api/chat/public-channels` - Publieke kanalen ophalen
- `POST /api/chat/direct` - Direct kanaal aanmaken tussen gebruikers
- `POST /api/chat/channels` - Nieuw kanaal aanmaken
- `POST /api/chat/channels/:id/join` - Deelnemen aan kanaal
- `POST /api/chat/channels/:id/leave` - Kanaal verlaten
- `GET /api/chat/channels/:channel_id/messages` - Berichten van kanaal ophalen
- `POST /api/chat/channels/:channel_id/messages` - Bericht verzenden
- `PUT /api/chat/messages/:id` - Bericht bewerken
- `DELETE /api/chat/messages/:id` - Bericht verwijderen
- `POST /api/chat/messages/:id/reactions` - Reactie toevoegen aan bericht
- `DELETE /api/chat/messages/:id/reactions/:emoji` - Reactie verwijderen
- `PUT /api/chat/presence` - Aanwezigheid status bijwerken
- `GET /api/chat/online-users` - Online gebruikers ophalen
- `GET /api/chat/users` - Gebruikers lijst voor direct chat
- `GET /api/chat/ws/:channel_id` - WebSocket verbinding voor real-time chat

#### Nieuwsbrief Beheer (Geïmplementeerd)
- `GET /api/newsletter` - Lijst van nieuwsbrieven ophalen
- `POST /api/newsletter` - Nieuwe nieuwsbrief aanmaken
- `GET /api/newsletter/:id` - Specifieke nieuwsbrief ophalen
- `PUT /api/newsletter/:id` - Nieuwsbrief bijwerken
- `DELETE /api/newsletter/:id` - Nieuwsbrief verwijderen
- `POST /api/newsletter/:id/send` - Nieuwsbrief verzenden

#### Gebruikersbeheer (Geïmplementeerd)
- `GET /api/users` - Lijst van gebruikers ophalen
- `GET /api/users/:id` - Gebruikersdetails ophalen
- `PUT /api/users/:id` - Gebruiker bijwerken
- `DELETE /api/users/:id` - Gebruiker verwijderen
- `POST /api/users/:id/roles` - Rol toewijzen aan gebruiker
- `DELETE /api/users/:id/roles/:role_id` - Rol verwijderen van gebruiker

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
    depends_on:
      - redis
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/api/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

### Docker Development Setup met Redis

Voor lokale ontwikkeling met volledige Redis ondersteuning:

```bash
# Start alle services (PostgreSQL, Redis, App)
docker-compose -f docker-compose.dev.yml up -d

# Check container status
docker-compose -f docker-compose.dev.yml ps

# View logs
docker logs dkl-email-service --tail 100

# Stop services
docker-compose -f docker-compose.dev.yml down
```

De development setup bevat:
- **PostgreSQL**: Database op poort 5433 (host) → 5432 (container)
- **Redis**: Cache op poort 6380 (host) → 6379 (container)
- **App**: Service op poort 8082 (host) → 8080 (container)

#### Redis Configuratie in Docker

Redis is standaard ingeschakeld in de Docker development setup:

```yaml
environment:
  REDIS_ENABLED: "true"
  REDIS_HOST: redis
  REDIS_PORT: 6379
  REDIS_PASSWORD: ""
  REDIS_DB: 0
```

Verifieer Redis verbinding:
```bash
# Test Redis in container
docker exec dkl-redis redis-cli ping
# Verwachte output: PONG

# Check health endpoint met Redis status
curl http://localhost:8082/api/health | jq '.checks.redis'
# Verwachte output: {"status": true}
```

    command: redis-server --requirepass your_redis_password
    volumes:
      - redis_data:/data
    restart: always

volumes:
  redis_data:
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

## 📬 Automatische Email Ophaling

De service biedt automatische ophaling van inkomende emails via de `EmailAutoFetcher` component:

### Functionaliteit
- Periodiek ophalen van emails uit geconfigureerde mailboxen
- Automatisch starten bij applicatie-opstart (configureerbaar)
- Voorkomen van duplicaten door UID-controle
- Graceful shutdown bij applicatie-afsluiting
- Thread-safe operatie met concurrency controle

### Configuratie
De automatische email ophaling kan worden geconfigureerd via de volgende omgevingsvariabelen:

```env
# Automatische email ophaling configuratie
EMAIL_FETCH_INTERVAL=15     # Interval in minuten tussen ophaal-operaties (standaard: 15)
DISABLE_AUTO_EMAIL_FETCH=false # Zet op "true" om automatisch ophalen uit te schakelen
```

### Email Accounts
De service is geconfigureerd om emails op te halen van:
- info@dekoninklijkeloop.nl
- inschrijving@dekoninklijkeloop.nl

Inkomende emails worden automatisch gesorteerd en opgeslagen in de database, en zijn beschikbaar via de beveiligde `/api/mail` endpoints.

## 🛠 Architectuur

De service volgt een modulaire architectuur met de volgende componenten:

### Core Components
- `handlers/` - HTTP request handlers
  - `email_handler.go` - Email verzending endpoints
  - `health_handler.go` - Health check endpoint
  - `metrics_handler.go` - Metrics endpoints
  - `contact_handler.go` - Contact formulier beheer endpoints
  - `aanmelding_handler.go` - Aanmelding beheer endpoints
  - `auth_handler.go` - Authenticatie endpoints
  - `chat_handler.go` - Chat systeem endpoints met WebSocket ondersteuning
  - `newsletter_handler.go` - Nieuwsbrief beheer endpoints
  - `user_handler.go` - Gebruikersbeheer endpoints
  - `mail_handler.go` - Inkomende email beheer endpoints

- `services/` - Business logic
  - `email_service.go` - Email verzending logica
  - `smtp_client.go` - SMTP communicatie
  - `rate_limiter.go` - Rate limiting
  - `email_batcher.go` - Batch processing
  - `email_metrics.go` - Metrics tracking
  - `prometheus_metrics.go` - Prometheus integratie
  - `email_auto_fetcher.go` - Automatische email ophaling
  - `mail_fetcher.go` - IMAP communicatie voor inkomende emails
  - `chat_service.go` - Chat kanaal en bericht beheer
  - `websocket_hub.go` - WebSocket hub voor real-time messaging
  - `auth_service.go` - Gebruikersauthenticatie en autorisatie
  - `permission_service.go` - RBAC permissie controle
  - `newsletter_service.go` - Nieuwsbrief generatie en verzending
  - `newsletter_sender.go` - Batch nieuwsbrief verzending
  - `telegram_bot_service.go` - Telegram bot integratie

- `models/` - Data structuren
  - `email.go` - Email gerelateerde structs
  - `contact.go` - Contact formulier model
  - `aanmelding.go` - Aanmelding formulier model
  - `chat_channel.go` - Chat kanaal model
  - `chat_message.go` - Chat bericht model
  - `chat_channel_participant.go` - Chat deelnemer model
  - `chat_message_reaction.go` - Bericht reactie model
  - `chat_user_presence.go` - Gebruikers aanwezigheid model
  - `newsletter.go` - Nieuwsbrief model
  - `role_rbac.go` - RBAC rol en permissie modellen
  - `gebruiker.go` - Gebruikersmodel

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

### Handler Implementaties

#### Contact Handler
De `ContactHandler` biedt een volledige implementatie voor het beheren van contactformulieren:
- **ListContactFormulieren**: Haalt een gepagineerde lijst van contactformulieren op
- **GetContactFormulier**: Haalt details van een specifiek contactformulier op, inclusief antwoorden
- **UpdateContactFormulier**: Werkt een contactformulier bij (status, notities)
- **DeleteContactFormulier**: Verwijdert een contactformulier
- **AddContactAntwoord**: Voegt een antwoord toe aan een contactformulier en stuurt een email naar de indiener
- **GetContactFormulierenByStatus**: Filtert contactformulieren op status (nieuw, in_behandeling, beantwoord, gesloten)

Alle endpoints zijn beveiligd met JWT authenticatie en vereisen admin rechten.

#### Aanmelding Handler
De `AanmeldingHandler` biedt een volledige implementatie voor het beheren van aanmeldingen:
- **ListAanmeldingen**: Haalt een gepagineerde lijst van aanmeldingen op
- **GetAanmelding**: Haalt details van een specifieke aanmelding op, inclusief antwoorden
- **UpdateAanmelding**: Werkt een aanmelding bij (status, notities)
- **DeleteAanmelding**: Verwijdert een aanmelding
- **AddAanmeldingAntwoord**: Voegt een antwoord toe aan een aanmelding en stuurt een email naar de indiener
- **GetAanmeldingenByRol**: Filtert aanmeldingen op rol (vrijwilliger, deelnemer, etc.)

Alle endpoints zijn beveiligd met JWT authenticatie en vereisen admin rechten.

## 🗄️ Database Architectuur

De applicatie is uitgebreid met een robuuste PostgreSQL database integratie voor het persistent opslaan van gegevens. Deze integratie maakt gebruik van GORM als ORM (Object-Relational Mapping) framework en implementeert het Repository Pattern voor een schone scheiding van verantwoordelijkheden.

### Database Modellen (Geïmplementeerd)

De volgende modellen zijn geïmplementeerd:

- **ContactFormulier**: Slaat contactformulier gegevens op met velden voor naam, email, bericht, status en behandeling.
- **ContactAntwoord**: Houdt antwoorden op contactformulieren bij, gekoppeld via een één-op-veel relatie.
- **Aanmelding**: Registreert aanmeldingen voor het evenement met persoonlijke gegevens en voorkeuren.
- **AanmeldingAntwoord**: Bewaart antwoorden op aanmeldingen, gekoppeld via een één-op-veel relatie.
- **EmailTemplate**: Slaat email templates op voor hergebruik en consistentie in communicatie.
- **VerzondEmail**: Houdt een log bij van alle verzonden emails voor tracking en auditing.
- **Gebruiker**: Beheert gebruikersaccounts voor administratieve toegang tot het systeem.
- **Migratie**: Houdt database migraties bij om schema-wijzigingen gecontroleerd uit te voeren.

### Repository Pattern (Geïmplementeerd)

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

## 🔐 Role-Based Access Control (RBAC) Systeem

De DKL Email Service gebruikt een geavanceerd Role-Based Access Control (RBAC) systeem voor gedetailleerde toegangscontrole. Het systeem is opgebouwd uit rollen, permissies en gebruikersrelaties.

### 🔑 Kernconcepten

#### Rollen
Rollen definiëren de verantwoordelijkheden van gebruikers in het systeem:

- **admin**: Volledige systeemtoegang
- **staff**: Ondersteunend personeel met beperkte beheerrechten
- **user**: Standaard gebruiker met basis chat permissies
- **owner**: Chat kanaal eigenaar met volledig kanaalbeheer
- **chat_admin**: Chat moderator met moderatierechten
- **member**: Chat lid met basis toegang
- **deelnemer**: Evenement deelnemer (geen speciale permissies)
- **begeleider**: Evenement begeleider (geen speciale permissies)
- **vrijwilliger**: Evenement vrijwilliger (geen speciale permissies)

#### Permissies
Permissies volgen het patroon `resource:action`:

**Contact Management:**
- `contact:read` - Contactformulieren bekijken
- `contact:write` - Contactformulieren bewerken
- `contact:delete` - Contactformulieren verwijderen

**Aanmeldingen:**
- `aanmelding:read` - Aanmeldingen bekijken
- `aanmelding:write` - Aanmeldingen bewerken
- `aanmelding:delete` - Aanmeldingen verwijderen

**Nieuwsbrieven:**
- `newsletter:read` - Nieuwsbrieven bekijken
- `newsletter:write` - Nieuwsbrieven aanmaken/bewerken
- `newsletter:send` - Nieuwsbrieven verzenden
- `newsletter:delete` - Nieuwsbrieven verwijderen

**Email Management:**
- `email:read` - Inkomende emails bekijken
- `email:write` - Emails bewerken
- `email:delete` - Emails verwijderen
- `email:fetch` - Nieuwe emails ophalen

**Gebruikersbeheer:**
- `user:read` - Gebruikers bekijken
- `user:write` - Gebruikers aanmaken/bewerken
- `user:delete` - Gebruikers verwijderen
- `user:manage_roles` - Gebruikersrollen beheren

**Chat:**
- `chat:read` - Chat kanalen bekijken
- `chat:write` - Berichten verzenden
- `chat:manage_channel` - Kanalen beheren
- `chat:moderate` - Berichten modereren

**Admin Email:**
- `admin_email:send` - Emails verzenden namens admin

### 🗄️ Database Structuur

Het RBAC systeem gebruikt de volgende tabellen:

- **`roles`**: Systeemrollen met beschrijvingen
- **`permissions`**: Granulaire permissies (resource:action)
- **`role_permissions`**: Koppeling tussen rollen en permissies
- **`user_roles`**: Koppeling tussen gebruikers en rollen
- **`gebruikers`**: Uitgebreide gebruikersinformatie

### ⚡ Performance Optimalisatie

- **Redis Caching**: Permissies worden gecached voor snelle toegang (vereist Redis)
- **Database Views**: `user_permissions` view voor efficiënte queries
- **Indexed Queries**: Geoptimaliseerde indexen voor snelle lookups

**⚠️ Belangrijk**: Redis is vereist voor optimale RBAC performance. Zonder Redis werken permissie controles nog steeds maar langzamer via database queries.

### 🔧 Beheer API

#### Rollen Beheer
```http
GET /api/rbac/roles          # Lijst van alle rollen
POST /api/rbac/roles         # Nieuwe rol aanmaken
PUT /api/rbac/roles/:id      # Rol bijwerken
DELETE /api/rbac/roles/:id   # Rol verwijderen
```

#### Permissies Beheer
```http
GET /api/rbac/permissions           # Lijst van alle permissies
POST /api/rbac/roles/:id/permissions # Permissie toewijzen aan rol
DELETE /api/rbac/roles/:id/permissions/:pid # Permissie verwijderen
```

#### Gebruikersrollen Beheer
```http
GET /api/users/:id/roles      # Gebruikersrollen ophalen
POST /api/users/:id/roles     # Rol toewijzen aan gebruiker
DELETE /api/users/:id/roles/:rid # Rol verwijderen van gebruiker
```

### 🔍 Permissie Controle

Het systeem controleert permissies op meerdere niveaus:

1. **JWT Token Validatie**: Basis authenticatie
2. **Rol Controle**: Gebruiker heeft juiste rol
3. **Permissie Controle**: Specifieke resource:action permissie
4. **Context Controle**: Additionele context (eigenaar, etc.)

### 📊 Audit Logging

Alle RBAC operaties worden gelogd voor audit doeleinden:

- Rol toewijzingen/verwijderingen
- Permissie wijzigingen
- Gebruiker acties
- Systeem configuratie wijzigingen

## 💬 Chat Systeem

De DKL Email Service bevat een volledig geïntegreerd real-time chat systeem voor interne communicatie tussen beheerders en gebruikers.

### ✨ Functionaliteiten

#### Real-time Messaging
- **WebSocket Gebaseerd**: Directe berichtenuitwisseling zonder polling
- **Multi-channel Ondersteuning**: Publieke, private en directe kanalen
- **Typing Indicators**: Zie wanneer anderen typen
- **Read Receipts**: Bevestiging dat berichten zijn gelezen
- **Message Reactions**: Emoji reacties op berichten
- **Message Editing**: Berichten bewerken na verzending
- **Presence Status**: Online/offline/away/busy status

#### Channel Types
- **Publieke Kanalen**: Voor algemene communicatie
- **Private Kanalen**: Beperkte toegang gebaseerd op rollen
- **Direct Messages**: Persoonlijke gesprekken tussen gebruikers

#### Gebruikersbeheer
- **Channel Rollen**: Owner, Admin, Member met verschillende permissies
- **Participant Management**: Gebruikers toevoegen/verwijderen
- **Moderatie Tools**: Berichten verwijderen, gebruikers muten

### 🏗️ Architectuur

#### Core Componenten
- **`ChatHandler`**: HTTP endpoints voor chat beheer
- **`ChatService`**: Business logic voor chat operaties
- **`WebSocketHub`**: Real-time bericht distributie
- **Database Models**: Chat kanalen, berichten, deelnemers, reacties

#### Database Structuur
- **`chat_channels`**: Kanaal informatie en instellingen
- **`chat_channel_participants`**: Deelnemers en hun rollen
- **`chat_messages`**: Berichten met metadata
- **`chat_message_reactions`**: Emoji reacties op berichten
- **`chat_user_presence`**: Online status van gebruikers

### 🔌 WebSocket API

#### Verbinding
```javascript
const ws = new WebSocket('ws://localhost:8080/api/chat/ws/channel_id');
```

#### Bericht Format
```json
{
  "type": "message",
  "channel_id": "channel-uuid",
  "user_id": "user-uuid",
  "content": "Hallo allemaal!",
  "timestamp": "2024-03-20T10:30:00Z"
}
```

#### Event Types
- `message`: Nieuw bericht
- `user_joined`: Gebruiker joint kanaal
- `user_left`: Gebruiker verlaat kanaal
- `typing_start`: Gebruiker begint te typen
- `typing_stop`: Gebruiker stopt met typen
- `presence_update`: Aanwezigheid status wijziging
- `reaction_added`: Reactie toegevoegd
- `message_edited`: Bericht bewerkt

### 🔒 Beveiliging

#### Authenticatie
- JWT token vereist voor WebSocket verbinding
- Channel-specifieke toegang controle
- RBAC integratie voor kanaal permissies

#### Autorisatie
- Channel deelnemers controle
- Message ownership verificatie
- Moderatie permissies controle

### 📊 Performance

#### Optimalisaties
- **Connection Pooling**: Efficiënt WebSocket beheer
- **Message Batching**: Bulk bericht verzending
- **Database Indexing**: Snelle queries op kanalen en berichten
- **Redis Caching**: Gebruiker presence en session data (vereist Redis)

**⚠️ Belangrijk**: Redis is vereist voor chat presence, typing indicators en optimale performance. Zonder Redis werken basis chat functies maar ontbreken real-time presence features.

#### Schaalbaarheid
- Horizontale scaling met Redis pub/sub
- Load balancing ondersteuning
- Connection limits per gebruiker

### 🛠️ Beheer API

#### Channel Management
```http
GET /api/chat/channels              # Gebruikers kanalen
POST /api/chat/channels             # Nieuw kanaal aanmaken
POST /api/chat/channels/:id/join    # Deelnemen aan kanaal
POST /api/chat/channels/:id/leave   # Kanaal verlaten
```

#### Message Management
```http
GET /api/chat/channels/:id/messages     # Berichten ophalen
POST /api/chat/channels/:id/messages    # Bericht verzenden
PUT /api/chat/messages/:id              # Bericht bewerken
DELETE /api/chat/messages/:id           # Bericht verwijderen
```

#### Presence & Status
```http
PUT /api/chat/presence              # Status bijwerken
GET /api/chat/online-users          # Online gebruikers
```

### 📱 Frontend Integratie

#### React Hook Voorbeeld
```javascript
const useChat = (channelId) => {
  const [messages, setMessages] = useState([]);
  const [ws, setWs] = useState(null);

  useEffect(() => {
    const websocket = new WebSocket(`/api/chat/ws/${channelId}`);

    websocket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data.type === 'message') {
        setMessages(prev => [...prev, data]);
      }
    };

    setWs(websocket);
    return () => websocket.close();
  }, [channelId]);

  const sendMessage = (content) => {
    ws.send(JSON.stringify({
      type: 'message',
      content: content
    }));
  };

  return { messages, sendMessage };
};
```

### 🔧 Configuratie

#### Omgevingsvariabelen
```env
# Chat systeem configuratie
CHAT_MAX_CONNECTIONS_PER_USER=5
CHAT_MESSAGE_HISTORY_LIMIT=1000
CHAT_PRESENCE_TIMEOUT=300
CHAT_TYPING_TIMEOUT=10
```

### 📈 Monitoring

#### Metrics
- Actieve WebSocket verbindingen
- Berichten per seconde
- Channel gebruik statistieken
- Gebruikers presence data
- Connection errors en timeouts

## 📰 Nieuwsbrief Systeem

De DKL Email Service bevat een geautomatiseerd nieuwsbrief systeem dat nieuws verzamelt van RSS feeds en professioneel opgemaakte emails verstuurt naar subscribers.

### ✨ Functionaliteiten

#### Automatische Content Generatie
- **RSS Feed Integratie**: Automatische verzameling van nieuws van geconfigureerde bronnen
- **Content Filtering**: Intelligente filtering op relevantie en categorie
- **Deduplicatie**: Voorkoming van dubbele content
- **Content Summarization**: Automatische samenvatting van lange artikelen

#### Email Templates
- **Responsive Design**: Professionele HTML templates die werken op alle apparaten
- **Brand Consistentie**: Gebruik van DKL huisstijl en kleuren
- **Dynamic Content**: Personalisatie gebaseerd op subscriber gegevens
- **Fallback Plaintext**: Automatische generatie van tekstversies

#### Subscriber Management
- **Opt-in/Oopt-out**: GDPR compliant subscriber beheer
- **Segmentatie**: Groepering van subscribers op interesses
- **Bounce Handling**: Automatische verwijdering van ongeldige email adressen
- **Analytics**: Open rates, click rates en engagement metrics

### 🏗️ Architectuur

#### Core Componenten
- **`NewsletterService`**: Hoofd service voor nieuwsbrief generatie
- **`NewsletterSender`**: Batch verzending van nieuwsbrieven
- **`NewsletterFetcher`**: RSS feed processing
- **`NewsletterFormatter`**: HTML template rendering
- **`NewsletterProcessor`**: Content filtering en verwerking

#### Database Structuur
- **`newsletters`**: Opgeslagen nieuwsbrieven met verzendstatus
- **`gebruikers.newsletter_subscribed`**: Subscriber status per gebruiker

### 🔄 Workflow

#### Automatische Generatie
1. **Feed Monitoring**: Periodieke controle van RSS feeds
2. **Content Extraction**: Relevant nieuws ophalen en filteren
3. **Template Rendering**: HTML email genereren
4. **Batch Queue**: Nieuwsbrief toevoegen aan verzend queue

#### Handmatige Beheer
1. **Draft Aanmaken**: Admin maakt concept nieuwsbrief
2. **Content Review**: Controle van content en formatting
3. **Test Verzending**: Preview naar test subscribers
4. **Bulk Verzending**: Productie verzending naar alle subscribers

### 📧 Email Templates

#### Template Structuur
```html
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{.Subject}}</title>
</head>
<body>
  <div class="header">
    <img src="{{.LogoUrl}}" alt="DKL Logo">
    <h1>{{.Title}}</h1>
  </div>

  <div class="content">
    {{range .Items}}
    <article class="news-item">
      <h2>{{.Title}}</h2>
      <p class="summary">{{.Description}}</p>
      <a href="{{.Link}}" class="read-more">Lees meer</a>
    </article>
    {{end}}
  </div>

  <div class="footer">
    <p>Ontvang je deze email niet graag? <a href="{{.UnsubscribeUrl}}">Afmelden</a></p>
  </div>
</body>
</html>
```

#### Template Variabelen
- `{{.Subject}}`: Email onderwerp
- `{{.Title}}`: Hoofdtitel van de nieuwsbrief
- `{{.Items}}`: Array van nieuws items
- `{{.LogoUrl}}`: URL naar logo
- `{{.UnsubscribeUrl}}`: Afmeld link

### 🔧 Configuratie

#### RSS Feed Configuratie
```env
# Nieuwsbrief configuratie
ENABLE_NEWSLETTER=true
NEWSLETTER_SOURCES=https://example.com/feed1.xml,https://example.com/feed2.xml
NEWSLETTER_CHECK_INTERVAL=1h
NEWSLETTER_MAX_ITEMS_PER_EMAIL=10
NEWSLETTER_TEMPLATE_PATH=./templates/newsletter.html
```

#### SMTP Configuratie
```env
# Nieuwsbrief SMTP (kan anders zijn dan reguliere emails)
NEWSLETTER_SMTP_HOST=smtp.newsletter-provider.com
NEWSLETTER_SMTP_PORT=587
NEWSLETTER_SMTP_USER=newsletter@dkl.nl
NEWSLETTER_SMTP_PASSWORD=secure_password
```

### 📊 Analytics & Monitoring

#### Verzend Metrics
- **Delivery Rate**: Percentage succesvol bezorgde emails
- **Open Rate**: Percentage geopende emails
- **Click Rate**: Percentage aangeklikte links
- **Bounce Rate**: Percentage bounced emails
- **Unsubscribe Rate**: Percentage afmeldingen

#### Content Metrics
- **Engagement Score**: Gecombineerde metric van opens/clicks
- **Popular Topics**: Meest gelezen onderwerpen
- **Optimal Send Time**: Beste verzendtijd gebaseerd op opens

### 🔒 Beveiliging & Compliance

#### GDPR Compliance
- **Explicit Consent**: Alleen verzending naar expliciete subscribers
- **Easy Unsubscribe**: Directe afmeld links in elke email
- **Data Minimization**: Alleen noodzakelijke gegevens opslaan
- **Audit Logging**: Alle subscriber acties loggen

#### Email Best Practices
- **SPF/DKIM/DMARC**: Email authenticatie setup
- **List Cleaning**: Regelmatige verwijdering van inactieve subscribers
- **Rate Limiting**: Gecontroleerde verzend snelheid
- **Monitoring**: Bounce en complaint monitoring

### 🛠️ Beheer API

#### Nieuwsbrief Beheer
```http
GET /api/newsletter              # Lijst van nieuwsbrieven
POST /api/newsletter             # Nieuwe nieuwsbrief aanmaken
GET /api/newsletter/:id          # Specifieke nieuwsbrief ophalen
PUT /api/newsletter/:id          # Nieuwsbrief bijwerken
DELETE /api/newsletter/:id       # Nieuwsbrief verwijderen
POST /api/newsletter/:id/send    # Nieuwsbrief verzenden
```

#### Subscriber Beheer
```http
GET /api/newsletter/subscribers         # Subscriber lijst
POST /api/newsletter/subscribers        # Subscriber toevoegen
DELETE /api/newsletter/subscribers/:id  # Subscriber verwijderen
```

#### Analytics
```http
GET /api/newsletter/:id/stats    # Verzend statistieken
GET /api/newsletter/analytics    # Algemene analytics
```

### 📈 Performance Optimalisatie

#### Batch Processing
- **Queue-based Sending**: Asynchrone verwerking van grote lijsten
- **Rate Limiting**: SMTP provider limits respecteren
- **Retry Logic**: Automatische retry bij tijdelijke fouten
- **Progress Tracking**: Realtime verzend voortgang

#### Content Caching
- **Template Caching**: Voorverwerkte templates
- **Feed Caching**: Vermindering van externe API calls
- **Image Optimization**: Geoptimaliseerde afbeeldingen voor email

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
- `API.md` - API documentatie (Bijgewerkt met Contact, Aanmelding, Chat en Nieuwsbrief endpoints)
- `DEPLOYMENT.md` - Deployment instructies
- `DEVELOPMENT.md` - Development guidelines
- `MONITORING.md` - Monitoring setup
- `SECURITY.md` - Security best practices
- `TEMPLATES.md` - Template documentatie
- `TESTING.md` - Test procedures
- `AUTH.md` - Authenticatie documentatie
- `RBAC_FRONTEND.md` - RBAC frontend integratie documentatie
- `WFC_INTEGRATION.md` - Whisky for Charity integratie documentatie

## 📮 Telegram Bot Service

De DKL Email Service bevat een geïntegreerde Telegram bot die beheerders in staat stelt om contactformulieren en aanmeldingen direct vanuit Telegram te bekijken en beheren. De bot biedt realtime notificaties en maakt het mogelijk om snel te reageren op nieuwe contactverzoeken en aanmeldingen.

### Functionaliteiten

- **Realtime Notificaties**: Ontvang direct notificaties bij nieuwe contactformulieren en aanmeldingen
- **Commando Interface**: Gebruik eenvoudige commando's om informatie op te vragen
- **Status Monitoring**: Bekijk de huidige status van de service en statistieken
- **Beveiligde Toegang**: Alleen geautoriseerde beheerders hebben toegang tot de bot

### Configuratie

De Telegram bot kan worden ingeschakeld en geconfigureerd via de volgende omgevingsvariabelen:

```env
# Telegram Bot configuratie
ENABLE_TELEGRAM_BOT=true        # Zet op "true" om de Telegram bot in te schakelen
TELEGRAM_BOT_TOKEN=your_bot_token  # Token verkregen van BotFather (@BotFather)
TELEGRAM_CHAT_ID=your_chat_id      # Chat ID voor het ontvangen van berichten
NOTIFICATION_MIN_PRIORITY=medium   # Minimale prioriteit voor notificaties (low, medium, high, critical)
NOTIFICATION_THROTTLE=15m          # Throttle tijd voor vergelijkbare notificaties
```

### Instellen van de Bot

1. **Bot Aanmaken**:
   - Start een chat met de [BotFather](https://t.me/botfather) op Telegram
   - Gebruik het commando `/newbot` om een nieuwe bot aan te maken
   - Volg de instructies en noteer het bot token

2. **Chat ID Verkrijgen**:
   - Start een chat met je nieuwe bot
   - Stuur een bericht naar de bot
   - Gebruik de API om je chat ID te vinden: `https://api.telegram.org/bot<YourBOTToken>/getUpdates`
   - Noteer het `chat_id` uit de JSON respons

3. **Configuratie Toepassen**:
   - Voeg het bot token en chat ID toe aan je omgevingsvariabelen
   - Herstart de service om de wijzigingen toe te passen

### Beschikbare Commando's

| Commando        | Beschrijving                            |
|-----------------|----------------------------------------|
| `/start`        | Start de interactie met de bot         |
| `/help`         | Toon alle beschikbare commando's       |
| `/contact`      | Toon recente contactformulieren        |
| `/contactnew`   | Toon nieuwe contactformulieren         |
| `/aanmelding`   | Toon recente aanmeldingen              |
| `/aanmeldingnew`| Toon onverwerkte aanmeldingen          |
| `/status`       | Toon service status en statistieken    |

### Migreren naar een Nieuwe Chat

Als je de bot wilt verplaatsen naar een andere chat of groep, volg dan deze stappen:

1. Voeg de bot toe aan de nieuwe groep
2. Verkrijg het nieuwe chat ID
3. Update de `TELEGRAM_CHAT_ID` omgevingsvariabele
4. Herstart de service of gebruik het API endpoint om notificaties opnieuw te verwerken:

```
POST /api/v1/notifications/reprocess-all
Authorization: Bearer <your_jwt_token>
```

### Probleemoplossing

- **Bot Reageert Niet**: Controleer of `ENABLE_TELEGRAM_BOT` is ingesteld op "true"
- **Conflictfout**: Als je een "Conflict: terminated by other getUpdates request" fout ziet, betekent dit dat er meerdere instanties van de bot actief zijn. Voer een harde herstart uit van de service.
- **Geen Notificaties**: Controleer of `NOTIFICATION_MIN_PRIORITY` niet te hoog is ingesteld

## 🥃 Whisky for Charity Integratie

De DKL Email Service bevat een gespecialiseerde module voor het afhandelen van orders voor het Whisky for Charity platform. Deze module is volledig gescheiden van de hoofdfunctionaliteit en heeft zijn eigen configuratie, templates en API endpoints.

### Features

- **Dedicated API Endpoint**: Beveiligd endpoint (`/api/wfc/order-email`) voor het versturen van orderbevestigingen
- **API Key Authenticatie**: Beveiligde toegang met een specifieke API key voor WFC integratie
- **Gescheiden Email Templates**: Speciaal ontworpen email templates voor WFC orders
- **Dual Email Flow**: Automatisch verzenden van zowel klantenbevestigingen als admin notificaties
- **Separate SMTP Configuratie**: Aparte SMTP-instellingen voor het Whisky for Charity domein
- **Geen Telegram Logging**: WFC-emails worden niet gelogd naar het Telegram kanaal voor privacy

### Architectuur

De WFC module is opgebouwd uit de volgende componenten:

#### Models
- `WFCOrder`: Representeert een Whisky for Charity bestelling met klantgegevens en items
- `WFCOrderItem`: Bevat gegevens over individuele items in een bestelling
- `WFCOrderEmailData`: Container voor email-gerelateerde data
- `WFCOrderRequest`: Definieert de verwachte structuur voor inkomende API requests

#### Handlers
- `WFCOrderHandler`: Verwerkt inkomende order requests en stuurt zowel klant- als adminemails
- `WFCAPIKeyMiddleware`: Middleware voor API key-gebaseerde authenticatie

#### Services
- `SendWFCOrderEmail`: Specifieke service-methode voor het verzenden van WFC order emails
- `SendWFC`: Gespecialiseerde methode in de SMTP client voor het gebruik van WFC-specifieke SMTP configuratie

#### Templates
- `wfc_order_confirmation.html`: Template voor klantenbevestigingen
- `wfc_order_admin.html`: Template voor admin notificaties met uitgebreide orderdetails

### API Documentatie

#### Order Email API Endpoint

- **URL**: `/api/wfc/order-email`
- **Methode**: `POST`
- **Authenticatie**: API Key (via `X-API-Key` header)
- **Content-Type**: `application/json`

**Request Body**:
```json
{
  "order_id": "string",              // Verplicht: Unieke order ID
  "customer_name": "string",         // Verplicht: Naam van de klant
  "customer_email": "string",        // Verplicht: Email van de klant
  "customer_address": "string",      // Optioneel: Adres van de klant
  "customer_city": "string",         // Optioneel: Stad van de klant
  "customer_postal": "string",       // Optioneel: Postcode van de klant
  "customer_country": "string",      // Optioneel: Land van de klant
  "total_amount": 0.00,              // Verplicht: Totaalbedrag van de order
  "items": [                         // Verplicht: Array van bestelde items
    {
      "id": "string",                // Unieke ID van het item
      "order_id": "string",          // Order ID waartoe dit item behoort
      "product_id": "string",        // Product ID referentie
      "product_name": "string",      // Naam van het product
      "quantity": 0,                 // Aantal stuks
      "price": 0.00                  // Prijs per stuk
    }
  ],
  "notify_admin": true,              // Optioneel: Niet meer gebruikt (beide emails worden altijd verstuurd)
  "template_type": "string"          // Optioneel: Type template om te gebruiken
}
```

**Voorbeeld Response (Success)**:
```json
{
  "success": true,
  "customer_email_sent": true,
  "admin_email_sent": true,
  "order_id": "order_123456"
}
```

**Voorbeeld Response (Error)**:
```json
{
  "error": "Missing required fields"
}
```

**Status Codes**:
- `200 OK`: Verzoek succesvol verwerkt
- `400 Bad Request`: Ongeldige input of ontbrekende velden
- `401 Unauthorized`: Ontbrekende of ongeldige API key
- `500 Internal Server Error`: Server error bij het verwerken van het verzoek

### Configuratie

De Whisky for Charity integratie wordt geconfigureerd via de volgende omgevingsvariabelen:

```env
# WFC SMTP Configuratie
WFC_SMTP_HOST=mail.example.com        // SMTP server voor WFC emails
WFC_SMTP_PORT=465                     // SMTP poort (meestal 465 voor SSL)
WFC_SMTP_USER=noreply@example.com     // SMTP gebruikersnaam
WFC_SMTP_PASSWORD=your_smtp_password  // SMTP wachtwoord
WFC_SMTP_FROM=noreply@example.com     // Afzender email adres

# WFC API Beveiliging
WFC_API_KEY=your_api_key              // API key voor WFC endpoint authenticatie

# WFC Admin Configuratie
WFC_ADMIN_EMAIL=admin@example.com     // Admin email adres voor ordernotificaties
WFC_SITE_URL=https://example.com      // Site URL voor links in emails
```

### Email Templates

De WFC module gebruikt twee gespecialiseerde email templates:

#### Customer Order Confirmation

Dit template wordt gebruikt voor het verzenden van orderbevestigingen naar klanten. Het bevat:
- Order details (ID, datum, totaalbedrag)
- Itemlijst met productnamen, aantallen en prijzen
- Verzendgegevens
- Links naar de orderpagina op de website

#### Admin Order Notification

Dit template wordt gebruikt voor het informeren van admins over nieuwe orders. Het bevat:
- Gedetailleerde orderinformatie
- Uitgebreide klantgegevens inclusief contactinformatie
- Volledige itemlijst met product IDs, prijzen en totalen
- Responsieve layout met geformatteerde secties

Beide templates ondersteunen de volgende template functies:
- `multiply`: Voor het berekenen van totaalbedragen per item (prijs × aantal)
- `currentYear`: Voor het dynamisch weergeven van het huidige jaar in copyright notices

### Integratie met Frontend

Het WFC order email endpoint is ontworpen om naadloos te integreren met het Whisky for Charity webplatform. De frontend maakt een HTTP POST request naar het endpoint met de ordergegevens nadat een bestelling is geplaatst.

De frontend moet de volgende stappen volgen:
1. Verzamel alle benodigde ordergegevens (klantinfo, items, bedragen)
2. Formateer de data volgens het verwachte request formaat (snake_case JSON)
3. Voeg de API key toe aan de `X-API-Key` header
4. Verstuur de POST request naar het `/api/wfc/order-email` endpoint
5. Verwerk de response om te bevestigen dat emails succesvol zijn verzonden

### Veiligheid en Privacy

De WFC module is ontworpen met speciale aandacht voor veiligheid en privacy:

1. **API Key Authenticatie**: Alle requests vereisen een geldige API key
2. **Gescheiden SMTP Configuratie**: WFC emails worden verzonden via een specifieke SMTP server
3. **Geen Telegram Logging**: WFC ordergegevens worden niet doorgestuurd naar Telegram voor privacy
4. **TLS Encryptie**: Alle SMTP communicatie is versleuteld met TLS
5. **Rate Limiting**: Standaard rate limiting wordt toegepast om misbruik te voorkomen

### Foutafhandeling

De WFC module bevat robuuste foutafhandeling:

1. **Request Validatie**: Controleert of alle vereiste velden aanwezig zijn
2. **Template Validatie**: Valideert templates voordat ze worden gebruikt
3. **SMTP Error Handling**: Vangt SMTP fouten af en logt deze voor troubleshooting
4. **Partial Success Handling**: Bij gedeeltelijk succes (bijv. alleen admin email verzonden) wordt dit in de response aangegeven
5. **Rate Limit Errors**: Bij overschrijding van rate limits wordt een duidelijke foutmelding gegeven

### Lokaal Testen

Voor lokaal testen van de WFC integratie:

1. Configureer de WFC omgevingsvariabelen in je `.env` bestand
2. Start de applicatie met `go run main.go`
3. Gebruik een tool zoals Postman of cURL om requests te sturen naar het endpoint:

```bash
curl -X POST http://localhost:8080/api/wfc/order-email \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your_api_key" \
  -d '{
    "order_id": "test_123",
    "customer_name": "Test User",
    "customer_email": "test@example.com",
    "total_amount": 150.00,
    "items": [
      {
        "id": "item_1",
        "order_id": "test_123",
        "product_id": "prod_1",
        "product_name": "Test Whisky",
        "quantity": 2,
        "price": 75.00
      }
    ]
  }'
```

## 🔄 Recente Updates

### Oktober 2025
- Toegevoegd: Uitgebreide README documentatie met alle nieuwe features
- Toegevoegd: Chat systeem met real-time WebSocket messaging
- Toegevoegd: Role-Based Access Control (RBAC) systeem met granulaire permissies
- Toegevoegd: Automatisch nieuwsbrief systeem met RSS feed integratie
- Toegevoegd: Uitgebreid gebruikersbeheer met rol- en permissiebeheer
- Verbeterd: Database architectuur met nieuwe tabellen voor chat, RBAC en nieuwsbrieven
- Toegevoegd: WebSocket ondersteuning voor real-time chat functionaliteit
- Verbeterd: API endpoints voor chat, gebruikersbeheer en nieuwsbrief beheer
- Toegevoegd: Redis caching voor RBAC permissies en chat presence
- Verbeterd: Service architectuur met nieuwe services voor chat en RBAC

### Maart 2025
- Toegevoegd: Whisky for Charity (WFC) integratie met API endpoints, modellen en email templates
- Verbeterd: WFC admin email template met uitgebreide orderdetails en responsieve layout
- Toegevoegd: Dual-email flow voor WFC orders (klant en admin notificaties)
- Toegevoegd: Gescheiden SMTP configuratie voor Whisky for Charity emails
- Toegevoegd: API key authenticatie voor WFC endpoints
- Toegevoegd: API documentatie voor WFC integratie
- Toegevoegd: API key authenticatie voor metrics endpoints
- Verbeterd: Metrics endpoints voor email en rate limiting statistieken
- Toegevoegd: Prometheus metrics endpoint
- Verbeterd: Test scripts voor API endpoints
- Toegevoegd: Automatische API tests met PowerShell script
- Verbeterd: Documentatie voor metrics en authenticatie
- Toegevoegd: Volledige implementatie van Contact Beheer endpoints
- Toegevoegd: Volledige implementatie van Aanmelding Beheer endpoints
- Verbeterd: Repository Pattern implementatie voor data toegang
- Toegevoegd: EmailAutoFetcher voor automatisch ophalen van inkomende emails
- Verbeterd: Integratie van automatische email ophaling in ServiceFactory
- Toegevoegd: Configuratie opties voor email ophaal interval en aan/uit zetten
- Verbeterd: Graceful shutdown met correcte afsluiting van achtergrondprocessen
- Toegevoegd: Test ondersteuning voor mail endpoints in test_api_light.ps1
- Toegevoegd: Telegram Bot Service voor het bekijken van contactformulieren en aanmeldingen
- Toegevoegd: Commando interface voor het opvragen van informatie via Telegram
- Toegevoegd: Testgegevens migraties voor het testen van contactformulieren en aanmeldingen
- Verbeterd: Probleemoplossing voor Telegram bot configuratie
- Toegevoegd: API endpoint voor het opnieuw verwerken van notificaties
- Toegevoegd: Test commando voor het verificeren van de Telegram bot configuratie