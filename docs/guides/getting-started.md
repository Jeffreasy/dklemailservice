# Getting Started

Snelstart gids voor de DKL Email Service.

## Quick Start

### 1. Clone Repository

```bash
git clone https://github.com/Jeffreasy/dklemailservice.git
cd dklemailservice
```

### 2. Configureer Environment

```bash
# Kopieer voorbeeld configuratie
cp .env.example .env

# Bewerk .env met jouw instellingen
nano .env
```

**Minimale Configuratie:**
```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-password
DB_NAME=dklemailservice
DB_SSL_MODE=disable

# SMTP (Standaard)
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=noreply@dekoninklijkeloop.nl
SMTP_PASSWORD=your-smtp-password
SMTP_FROM=noreply@dekoninklijkeloop.nl

# SMTP (Registratie)
REGISTRATION_SMTP_HOST=smtp.example.com
REGISTRATION_SMTP_PORT=587
REGISTRATION_SMTP_USER=inschrijving@dekoninklijkeloop.nl
REGISTRATION_SMTP_PASSWORD=your-smtp-password
REGISTRATION_SMTP_FROM=inschrijving@dekoninklijkeloop.nl

# Email Adressen
ADMIN_EMAIL=info@dekoninklijkeloop.nl
REGISTRATION_EMAIL=inschrijving@dekoninklijkeloop.nl

# JWT
JWT_SECRET=change-this-to-a-secure-random-string
```

### 3. Setup Database

**PostgreSQL Installeren:**
```bash
# Windows (Chocolatey)
choco install postgresql

# macOS
brew install postgresql

# Linux
sudo apt-get install postgresql
```

**Database Aanmaken:**
```bash
# Start PostgreSQL
# Windows: Start via Services
# macOS/Linux: 
sudo service postgresql start

# Maak database
psql -U postgres
```

```sql
CREATE DATABASE dklemailservice;
CREATE USER dkluser WITH PASSWORD 'dev-password';
GRANT ALL PRIVILEGES ON DATABASE dklemailservice TO dkluser;
\q
```

**Update .env:**
```bash
DB_USER=dkluser
DB_PASSWORD=dev-password
```

### 4. Installeer Dependencies

```bash
go mod download
```

### 5. Start Applicatie

```bash
go run main.go
```

**Verwachte Output:**
```
INFO DKL Email Service wordt gestart version=1.0.0
INFO Database configuratie geladen host=localhost port=5432
INFO Database verbinding succesvol
INFO Migraties uitgevoerd
INFO Server gestart port=8080
```

### 6. Test de Service

**Health Check:**
```bash
curl http://localhost:8080/api/health
```

**Response:**
```json
{
    "status": "healthy",
    "version": "1.0.0",
    "services": {
        "database": "connected",
        "email_service": "operational"
    }
}
```

## Eerste Stappen

### 1. Test Contact Email (Test Mode)

```bash
curl -X POST http://localhost:8080/api/contact-email \
  -H "Content-Type: application/json" \
  -H "X-Test-Mode: true" \
  -d '{
    "naam": "Test User",
    "email": "test@example.com",
    "bericht": "Dit is een test bericht",
    "privacy_akkoord": true
  }'
```

**Response:**
```json
{
    "success": true,
    "message": "[TEST MODE] Je bericht is verwerkt (geen echte email verzonden).",
    "test_mode": true
}
```

### 2. Login als Admin

**Default Credentials:**
- Email: `admin@dekoninklijkeloop.nl`
- Wachtwoord: Zie database seed data

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@dekoninklijkeloop.nl",
    "wachtwoord": "your-admin-password"
  }'
```

**Response:**
```json
{
    "success": true,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "dGhpcyBpcyBhIHJlZnJlc2ggdG9rZW4...",
    "user": {
        "id": "uuid",
        "email": "admin@dekoninklijkeloop.nl",
        "naam": "Admin",
        "rol": "admin",
        "permissions": [...]
    }
}
```

### 3. Bekijk Contact Formulieren

```bash
TOKEN="your-jwt-token-from-login"

curl http://localhost:8080/api/contact \
  -H "Authorization: Bearer $TOKEN"
```

### 4. Test Email Auto Fetcher (Optioneel)

**Configureer IMAP:**
```bash
# In .env
INFO_EMAIL=info@dekoninklijkeloop.nl
INFO_EMAIL_PASSWORD=your-password
IMAP_SERVER=mail.hostnet.nl
IMAP_PORT=993
EMAIL_FETCH_INTERVAL=1  # 1 minuut voor testing
```

**Handmatig Triggeren:**
```bash
curl -X POST http://localhost:8080/api/mail/fetch \
  -H "Authorization: Bearer $TOKEN"
```

## Project Structuur Begrijpen

### Belangrijkste Directories

```
dklemailservice/
├── handlers/      # HTTP request handlers
├── services/      # Business logic
├── repository/    # Database access
├── models/        # Data structures
├── config/        # Configuratie
├── database/      # Migraties
├── templates/     # Email templates
└── tests/         # Test suite
```

### Code Flow

```
Request → Middleware → Handler → Service → Repository → Database
                                    ↓
                                 SMTP/IMAP
```

**Voorbeeld Flow:**

1. **Request:** `POST /api/contact-email`
2. **Middleware:** CORS, Test Mode detectie
3. **Handler:** [`handlers/email_handler.go:43`](../../handlers/email_handler.go:43)
4. **Service:** [`services/email_service.go:129`](../../services/email_service.go:129)
5. **SMTP:** Email verzending
6. **Response:** Success/Error

## Development Workflow

### 1. Maak Feature Branch

```bash
git checkout -b feature/nieuwe-feature
```

### 2. Implementeer Feature

**Voorbeeld: Nieuwe Email Template**

**a. Maak Template:**
```html
<!-- templates/nieuwe_email.html -->
<!DOCTYPE html>
<html>
<head>
    <title>{{ .Subject }}</title>
</head>
<body>
    <h1>Hallo {{ .Naam }}</h1>
    <p>{{ .Bericht }}</p>
</body>
</html>
```

**b. Update Service:**
```go
// services/email_service.go
templateFiles := []string{
    "contact_admin_email",
    "contact_email",
    "nieuwe_email",  // Voeg toe
}
```

**c. Maak Handler Method:**
```go
func (s *EmailService) SendNieuweEmail(data *NieuweEmailData) error {
    return s.sendEmailWithTemplate("nieuwe_email", data.To, data.Subject, data)
}
```

### 3. Test Feature

```bash
# Run tests
go test ./...

# Test handmatig
curl -X POST http://localhost:8080/api/nieuwe-email \
  -H "Content-Type: application/json" \
  -H "X-Test-Mode: true" \
  -d '{"naam": "Test", "bericht": "Test"}'
```

### 4. Commit & Push

```bash
git add .
git commit -m "feat: voeg nieuwe email template toe"
git push origin feature/nieuwe-feature
```

## Common Tasks

### Nieuwe Gebruiker Aanmaken

**Via Database:**
```sql
INSERT INTO gebruikers (naam, email, wachtwoord_hash, rol, is_actief)
VALUES (
    'Nieuwe User',
    'user@example.com',
    '$2a$10$...', -- Gebruik bcrypt hash
    'gebruiker',
    true
);
```

**Via API (Admin Required):**
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "naam": "Nieuwe User",
    "email": "user@example.com",
    "wachtwoord": "password123",
    "rol": "gebruiker"
  }'
```

### Email Template Wijzigen

**1. Bewerk Template:**
```bash
nano templates/contact_email.html
```

**2. Herstart Service:**
```bash
# Templates worden bij startup geladen
go run main.go
```

**3. Test:**
```bash
curl -X POST http://localhost:8080/api/contact-email \
  -H "X-Test-Mode: true" \
  -d '{"naam":"Test","email":"test@example.com","bericht":"Test","privacy_akkoord":true}'
```

### Database Migratie Toevoegen

**1. Maak Migratie Bestand:**
```bash
touch database/migrations/V1_XX__description.sql
```

**2. Schrijf Migratie:**
```sql
-- V1_XX__description.sql
CREATE TABLE IF NOT EXISTS nieuwe_tabel (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    naam VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO migraties (versie, naam, toegepast)
VALUES ('1.XX.0', 'Description', CURRENT_TIMESTAMP)
ON CONFLICT (versie) DO NOTHING;
```

**3. Herstart Service:**
```bash
# Migraties worden automatisch uitgevoerd
go run main.go
```

## Troubleshooting

### Service Start Niet

**Check Logs:**
```bash
# Laatste 50 regels
tail -50 server.log

# Live logs
tail -f server.log
```

**Common Issues:**
- Database connectie mislukt
- Ontbrekende environment variabelen
- Poort al in gebruik
- SMTP credentials ongeldig

### Database Connectie Fout

**Test Connectie:**
```bash
psql -h localhost -p 5432 -U postgres -d dklemailservice
```

**Check Environment:**
```bash
echo $DB_HOST
echo $DB_PORT
echo $DB_USER
```

### SMTP Fout

**Test SMTP:**
```bash
telnet smtp.example.com 587
```

**Check Credentials:**
```bash
echo $SMTP_USER
echo $SMTP_HOST
```

## Volgende Stappen

### Leer Meer

1. **[Development Guide](./development.md)** - Uitgebreide development setup
2. **[API Documentation](../api/rest-api.md)** - Complete API referentie
3. **[Architecture](../architecture/components.md)** - Systeem architectuur
4. **[Testing Guide](./testing.md)** - Test procedures
5. **[Deployment Guide](./deployment.md)** - Productie deployment

### Experimenteer

**Test Endpoints:**
- Contact formulier: `POST /api/contact-email`
- Aanmelding: `POST /api/aanmelding-email`
- Login: `POST /api/auth/login`
- Health: `GET /api/health`

**Gebruik Test Mode:**
```http
X-Test-Mode: true
```

**Bekijk Metrics:**
```bash
curl http://localhost:8080/metrics
```

### Join Development

**Contribute:**
1. Fork repository
2. Maak feature branch
3. Implementeer feature met tests
4. Submit pull request

**Code Review:**
- Volg coding standards
- Schrijf tests (90%+ coverage)
- Update documentatie
- Test lokaal

## Resources

### Documentation

- [README](../../README.md) - Project overview
- [API Docs](../api/) - API referentie
- [Architecture](../architecture/) - Technische architectuur
- [Guides](../guides/) - Handleidingen

### External Links

- [Go Documentation](https://golang.org/doc/)
- [Fiber Framework](https://docs.gofiber.io/)
- [GORM](https://gorm.io/docs/)
- [PostgreSQL](https://www.postgresql.org/docs/)

### Support

- GitHub Issues: [Report bugs](https://github.com/Jeffreasy/dklemailservice/issues)
- Email: admin@dekoninklijkeloop.nl

## Zie Ook

- [Development Guide](./development.md) - Uitgebreide development setup
- [Deployment Guide](./deployment.md) - Productie deployment
- [Testing Guide](./testing.md) - Test procedures
- [API Documentation](../api/rest-api.md) - API referentie