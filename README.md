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

- `services/` - Business logic
  - `email_service.go` - Email verzending logica
  - `smtp_client.go` - SMTP communicatie
  - `rate_limiter.go` - Rate limiting
  - `email_batcher.go` - Batch processing
  - `email_metrics.go` - Metrics tracking
  - `prometheus_metrics.go` - Prometheus integratie
  - `email_auto_fetcher.go` - Automatische email ophaling
  - `mail_fetcher.go` - IMAP communicatie voor inkomende emails

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
- `API.md` - API documentatie (Bijgewerkt met Contact en Aanmelding Beheer endpoints)
- `DEPLOYMENT.md` - Deployment instructies
- `DEVELOPMENT.md` - Development guidelines
- `MONITORING.md` - Monitoring setup
- `SECURITY.md` - Security best practices
- `TEMPLATES.md` - Template documentatie
- `TESTING.md` - Test procedures
- `AUTH.md` - Authenticatie documentatie (Nieuw)
- `CONTACT_AANMELDING.md` - Contact en Aanmelding Beheer documentatie (Nieuw)

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