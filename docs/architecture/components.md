# System Components

## Overzicht

De DKL Email Service bestaat uit verschillende componenten die samenwerken om een complete email management oplossing te bieden.

## Core Components

### 1. Main Application ([`main.go`](../../main.go:1))

Het entry point van de applicatie dat alle componenten initialiseert en configureert.

**Verantwoordelijkheden:**
- Environment validatie
- Database initialisatie en migraties
- Service factory setup
- Route registratie
- Graceful shutdown handling

**Belangrijke Functies:**
```go
func ValidateEnv() error                    // Valideert omgevingsvariabelen
func initializeMailFetcher() *MailFetcher  // Configureert email fetcher
```

### 2. Handlers Layer

HTTP request handlers die de API endpoints implementeren.

#### Email Handler ([`handlers/email_handler.go`](../../handlers/email_handler.go:1))

Verwerkt contact formulieren en aanmeldingen.

**Endpoints:**
- `POST /api/contact-email` - Contact formulier verwerking
- `POST /api/aanmelding-email` - Aanmelding verwerking

**Features:**
- Test mode detectie via header, locals of body parameter
- Automatische notificaties via Telegram
- Email validatie en sanitization
- Database opslag van aanmeldingen

```go
type EmailHandler struct {
    emailService        EmailServiceInterface
    notificationService services.NotificationService
    aanmeldingRepo      repository.AanmeldingRepository
}
```

#### Auth Handler ([`handlers/auth_handler.go`](../../handlers/auth_handler.go:1))

Beheert authenticatie en autorisatie.

**Endpoints:**
- `POST /api/auth/login` - Gebruiker login met JWT
- `POST /api/auth/logout` - Gebruiker logout
- `POST /api/auth/refresh` - Token refresh
- `GET /api/auth/profile` - Gebruikersprofiel ophalen
- `POST /api/auth/reset-password` - Wachtwoord wijzigen

**Features:**
- JWT token generatie en validatie
- Refresh token mechanisme
- Rate limiting op login pogingen
- RBAC permissies integratie

```go
type AuthHandler struct {
    authService       services.AuthService
    permissionService services.PermissionService
    rateLimiter       services.RateLimiterService
}
```

#### Contact Handler ([`handlers/contact_handler.go`](../../handlers/contact_handler.go:1))

Beheert contact formulieren in het admin panel.

**Endpoints:**
- `GET /api/contact` - Lijst van contact formulieren
- `GET /api/contact/:id` - Specifiek contact formulier
- `PUT /api/contact/:id` - Contact formulier bijwerken
- `DELETE /api/contact/:id` - Contact formulier verwijderen
- `POST /api/contact/:id/antwoord` - Antwoord toevoegen
- `GET /api/contact/status/:status` - Filteren op status

#### Aanmelding Handler ([`handlers/aanmelding_handler.go`](../../handlers/aanmelding_handler.go:1))

Beheert aanmeldingen in het admin panel.

**Endpoints:**
- `GET /api/aanmelding` - Lijst van aanmeldingen
- `GET /api/aanmelding/:id` - Specifieke aanmelding
- `PUT /api/aanmelding/:id` - Aanmelding bijwerken
- `DELETE /api/aanmelding/:id` - Aanmelding verwijderen
- `POST /api/aanmelding/:id/antwoord` - Antwoord toevoegen
- `GET /api/aanmelding/rol/:rol` - Filteren op rol

#### Mail Handler ([`handlers/mail_handler.go`](../../handlers/mail_handler.go:1))

Beheert inkomende emails opgehaald door EmailAutoFetcher.

**Endpoints:**
- `GET /api/mail` - Lijst van inkomende emails
- `GET /api/mail/:id` - Specifieke email details
- `PUT /api/mail/:id/processed` - Markeer als verwerkt
- `DELETE /api/mail/:id` - Email verwijderen
- `POST /api/mail/fetch` - Handmatig emails ophalen
- `GET /api/mail/unprocessed` - Onverwerkte emails
- `GET /api/mail/account/:type` - Filter op account type

#### Chat Handler ([`handlers/chat_handler.go`](../../handlers/chat_handler.go:1))

Beheert real-time chat functionaliteit via WebSockets.

**Endpoints:**
- `GET /api/chat/channels` - Lijst van chat kanalen
- `POST /api/chat/channels` - Nieuw kanaal aanmaken
- `GET /api/chat/channels/:id/messages` - Berichten ophalen
- `POST /api/chat/channels/:id/messages` - Bericht versturen
- `GET /ws/chat/:channelId` - WebSocket connectie

### 3. Services Layer

Business logic en externe integraties.

#### Email Service ([`services/email_service.go`](../../services/email_service.go:1))

Centrale service voor email verzending.

**Functionaliteit:**
- Template-based email rendering
- Multi-SMTP configuratie (standaard, registratie, WFC)
- Rate limiting integratie
- Metrics tracking
- Test mode support
- Excluded email addresses voor testing

**Belangrijke Methoden:**
```go
func (s *EmailService) SendContactEmail(data *models.ContactEmailData) error
func (s *EmailService) SendAanmeldingEmail(data *models.AanmeldingEmailData) error
func (s *EmailService) SendWFCOrderEmail(data *models.WFCOrderEmailData) error
func (s *EmailService) SendEmail(to, subject, body string, fromAddress ...string) error
```

**Templates:**
- `contact_email.html` - Gebruiker bevestiging
- `contact_admin_email.html` - Admin notificatie
- `aanmelding_email.html` - Aanmelding bevestiging
- `aanmelding_admin_email.html` - Admin notificatie
- `wfc_order_confirmation.html` - WFC order bevestiging
- `wfc_order_admin.html` - WFC admin notificatie
- `newsletter.html` - Newsletter template

#### SMTP Client ([`services/smtp_client.go`](../../services/smtp_client.go:1))

Beheert SMTP verbindingen en email verzending.

**Configuraties:**
- Standaard SMTP (contact emails)
- Registratie SMTP (aanmelding emails)
- WFC SMTP (Whisky for Charity)

**Features:**
- TLS/SSL support
- Connection pooling
- Retry mechanisme
- Test mode (geen echte verzending)

#### Auth Service ([`services/auth_service.go`](../../services/auth_service.go:1))

Authenticatie en gebruikersbeheer.

**Functionaliteit:**
- JWT token generatie en validatie
- Refresh token mechanisme
- Wachtwoord hashing met bcrypt
- Gebruiker CRUD operaties
- RBAC integratie

#### Email Auto Fetcher ([`services/email_auto_fetcher.go`](../../services/email_auto_fetcher.go:1))

Automatisch ophalen van inkomende emails via IMAP.

**Features:**
- Periodiek ophalen (configureerbaar interval)
- Multi-account support (info, inschrijving)
- Duplicate detectie
- Thread-safe operaties
- Graceful shutdown

**Configuratie:**
```go
EMAIL_FETCH_INTERVAL=15              // Minuten tussen fetches
DISABLE_AUTO_EMAIL_FETCH=false       // Uitschakelen indien true
INFO_EMAIL=info@example.com          // Info account
INFO_EMAIL_PASSWORD=password         // Info wachtwoord
INSCHRIJVING_EMAIL=reg@example.com   // Registratie account
INSCHRIJVING_EMAIL_PASSWORD=password // Registratie wachtwoord
```

#### Mail Fetcher ([`services/mail_fetcher.go`](../../services/mail_fetcher.go:1))

IMAP client voor email ophaling.

**Functionaliteit:**
- IMAP verbinding beheer
- Email parsing (headers, body, attachments)
- Multi-account configuratie
- Error handling en retry logic

#### Notification Service ([`services/notification_service.go`](../../services/notification_service.go:1))

Telegram notificaties voor belangrijke events.

**Features:**
- Prioriteit levels (low, medium, high, critical)
- Throttling om spam te voorkomen
- Verschillende notificatie types
- Emoji's voor visuele feedback

#### Chat Service ([`services/chat_service.go`](../../services/chat_service.go:1))

Real-time chat functionaliteit.

**Features:**
- Kanaal beheer
- Bericht verzending en ontvangst
- Gebruiker presence tracking
- Message reactions
- WebSocket integratie via Hub

#### WebSocket Hub ([`services/websocket_hub.go`](../../services/websocket_hub.go:1))

Centrale hub voor WebSocket connecties.

**Functionaliteit:**
- Client registratie en deregistratie
- Broadcast berichten naar kanalen
- Connection management
- Concurrent-safe operaties

#### Newsletter Service ([`services/newsletter_service.go`](../../services/newsletter_service.go:1))

Geautomatiseerde newsletter generatie en verzending.

**Components:**
- Newsletter Fetcher - Haalt content op
- Newsletter Formatter - Formatteert content
- Newsletter Processor - Verwerkt en genereert
- Newsletter Sender - Verzendt naar subscribers

#### Rate Limiter ([`services/rate_limiter.go`](../../services/rate_limiter.go:1))

Rate limiting voor API endpoints.

**Features:**
- Redis-backed voor distributed rate limiting
- Sliding window algoritme
- Per-endpoint configuratie
- IP-based en global limits

#### Permission Service ([`services/permission_service.go`](../../services/permission_service.go:1))

RBAC permission management.

**Functionaliteit:**
- Gebruiker permissies ophalen
- Rol toewijzing
- Permission checks
- Resource-based access control

### 4. Repository Layer

Data access layer voor database operaties.

**Repositories:**
- `GebruikerRepository` - Gebruikersbeheer
- `ContactRepository` - Contact formulieren
- `AanmeldingRepository` - Aanmeldingen
- `IncomingEmailRepository` - Inkomende emails
- `ChatChannelRepository` - Chat kanalen
- `ChatMessageRepository` - Chat berichten
- `NotificationRepository` - Notificaties
- `NewsletterRepository` - Nieuwsbrieven
- `PermissionRepository` - RBAC permissies
- `RoleRepository` - RBAC rollen

**Pattern:**
Alle repositories implementeren een interface voor testbaarheid en volgen het Repository pattern.

### 5. Models Layer

Data models en structuren.

**Belangrijke Models:**
- [`Gebruiker`](../../models/gebruiker.go:11) - Gebruiker account
- [`ContactFormulier`](../../models/contact.go:1) - Contact formulier
- [`Aanmelding`](../../models/aanmelding.go:1) - Aanmelding
- [`IncomingEmail`](../../models/incoming_email.go:1) - Inkomende email
- [`ChatChannel`](../../models/chat_channel.go:1) - Chat kanaal
- [`ChatMessage`](../../models/chat_message.go:1) - Chat bericht
- [`Notification`](../../models/notification.go:1) - Notificatie
- [`Newsletter`](../../models/newsletter.go:1) - Nieuwsbrief
- [`RBACRole`](../../models/role_rbac.go:1) - RBAC rol
- [`Permission`](../../models/role_rbac.go:1) - RBAC permissie

### 6. Middleware

Request processing middleware.

**Middleware Components:**
- `AuthMiddleware` - JWT authenticatie verificatie
- `AdminMiddleware` - Admin rol verificatie
- `PermissionMiddleware` - RBAC permission checks
- `RateLimitMiddleware` - Rate limiting
- `TestModeMiddleware` - Test mode detectie
- `CORSMiddleware` - CORS configuratie

### 7. Database Layer

Database configuratie en migraties.

**Components:**
- [`config/database.go`](../../config/database.go:1) - Database configuratie
- [`database/migrations.go`](../../database/migrations.go:1) - Migratie manager
- [`database/migrations/`](../../database/migrations/) - SQL migratie bestanden

**Features:**
- Automatische migraties bij startup
- Versie tracking
- Seed data voor initiële setup
- PostgreSQL met GORM ORM

### 8. Logger

Structured logging systeem.

**Features:**
- Multiple log levels (DEBUG, INFO, WARN, ERROR, FATAL)
- Structured logging met key-value pairs
- ELK stack integratie
- Console en file output
- Request/response logging

### 9. Metrics & Monitoring

**Components:**
- Prometheus metrics export
- Email metrics tracking
- Rate limit metrics
- System health checks
- Performance monitoring

## Component Interacties

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────────┐
│         Fiber HTTP Server           │
│  ┌──────────────────────────────┐  │
│  │      Middleware Layer        │  │
│  │  - Auth                      │  │
│  │  - CORS                      │  │
│  │  - Rate Limiting             │  │
│  │  - Test Mode                 │  │
│  └──────────────────────────────┘  │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│         Handlers Layer              │
│  - EmailHandler                     │
│  - AuthHandler                      │
│  - ContactHandler                   │
│  - AanmeldingHandler                │
│  - MailHandler                      │
│  - ChatHandler                      │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│         Services Layer              │
│  - EmailService                     │
│  - AuthService                      │
│  - NotificationService              │
│  - ChatService                      │
│  - PermissionService                │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│       Repository Layer              │
│  - Database Access                  │
│  - CRUD Operations                  │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│         PostgreSQL Database         │
└─────────────────────────────────────┘
```

## External Integrations

- **SMTP Servers** - Email verzending
- **IMAP Servers** - Email ophaling
- **Redis** - Rate limiting en caching
- **Telegram Bot API** - Notificaties
- **Prometheus** - Metrics export
- **ELK Stack** - Centralized logging

## Configuration

Alle componenten worden geconfigureerd via environment variabelen. Zie [`.env.example`](../../.env.example) voor een complete lijst.

## Deployment

De applicatie kan worden gedeployed als:
- Standalone binary
- Docker container
- Cloud service (Render, Heroku, etc.)

Zie [Deployment Guide](../guides/deployment.md) voor details.