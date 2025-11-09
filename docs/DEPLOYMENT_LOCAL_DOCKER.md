# Lokale Docker Deployment - Instructies

**Document Versie:** 1.0  
**Datum:** 2025-01-08  

---

## 🚨 Huidige Status

De backend fix voor `/api/participant/emails` is succesvol toegepast en ge-build in Docker.

**Probleem:** Container crasht vanwege ontbrekende environment variabelen.

---

## ⚙️ Vereiste Configuratie

### Stap 1: Maak .env bestand

Kopieer `.env.example` naar `.env`:

```bash
cp .env.example .env
```

### Stap 2: Vul minimale vereiste variabelen in

Edit `.env` en vul **minimaal** deze variabelen in:

```bash
# Database
DB_HOST=dkl-postgres
DB_PORT=5432
DB_USER=dkl_user
DB_PASSWORD=yourpassword
DB_NAME=dklemailservice
DB_SSL_MODE=disable

# JWT
JWT_SECRET=your-secret-key-change-this-in-production

# SMTP (Algemeen)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=your-email@gmail.com

# SMTP (Registratie)
REGISTRATION_SMTP_HOST=smtp.gmail.com
REGISTRATION_SMTP_PORT=587
REGISTRATION_SMTP_USER=your-email@gmail.com
REGISTRATION_SMTP_PASSWORD=your-app-password
REGISTRATION_SMTP_FROM=your-email@gmail.com

# Email Adressen
ADMIN_EMAIL=admin@example.com
REGISTRATION_EMAIL=registratie@example.com

# Optioneel: Email Fetcher (kan leeg blijven voor testing)
INFO_EMAIL=
INFO_EMAIL_PASSWORD=
INSCHRIJVING_EMAIL=
INSCHRIJVING_EMAIL_PASSWORD=
DISABLE_AUTO_EMAIL_FETCH=true

# Optioneel: Telegram (kan leeg blijven)
TELEGRAM_BOT_TOKEN=
TELEGRAM_CHAT_ID=

# Optioneel: Cloudinary (kan leeg blijven)
CLOUDINARY_CLOUD_NAME=
CLOUDINARY_API_KEY=
CLOUDINARY_API_SECRET=
CLOUDINARY_UPLOAD_PRESET=

# Optioneel: Newsletter (kan leeg blijven)
NEWSLETTER_SOURCES=
ENABLE_NEWSLETTER=false
```

### Stap 3: Herstart Docker Containers

```bash
docker-compose down
docker-compose up -d
```

### Stap 4: Verifieer dat service draait

```bash
# Check container status
docker-compose ps

# Check logs
docker-compose logs app --tail=50

# Test health endpoint
curl http://localhost:8080/api/health
```

---

## 🧪 Test de Fix

### Stap 1: Login en Verkrijg Token

```bash
# Login (vervang met jouw credentials)
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "yourpassword"
  }'

# Response bevat: {"token": "eyJ..."}
# Kopieer de token value
```

### Stap 2: Test Participant Emails Endpoint

```bash
# Test de gefixte route
curl -X GET http://localhost:8080/api/participant/emails \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json"

# Verwachte response (200 OK):
{
  "participant_emails": ["email1@example.com", "email2@example.com"],
  "system_emails": ["admin@example.com", "registratie@example.com"],
  "all_emails": [...],
  "counts": {
    "participants": 2,
    "system": 2,
    "total": 4
  }
}
```

### Stap 3: Verifieer dat /:id Route Nog Steeds Werkt

```bash
# Test met een echte UUID (na het aanmaken van een participant)
curl -X GET http://localhost:8080/api/participant/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json"

# Verwacht: 200 OK met participant data OF 404 Not Found
```

### Stap 4: Test Under Construction Endpoint

```bash
# Test zonder auth (public endpoint)
curl -X GET http://localhost:8080/api/under-construction/active \
  -H "Content-Type: application/json"

# Verwacht: 404 Not Found met "No active under construction found"
# Dit is CORRECT gedrag (geen bug!)
```

---

## 📊 Verwachte Resultaten

### ✅ Succesvolle Deployment

```bash
# docker-compose ps output:
NAME                IMAGE                 COMMAND                   SERVICE    STATUS
dkl-email-service   dklemailservice-app   "/bin/sh -c 'if [ \"$…"   app        Up (healthy)
dkl-postgres        postgres:17           "docker-entrypoint.s…"    postgres   Up (healthy)
dkl-redis           redis:7-alpine        "docker-entrypoint.s…"    redis      Up (healthy)
```

```bash
# Health check output:
{
  "status": "healthy",
  "database": "connected",
  "redis": "connected",
  "version": "1.1.0",
  "timestamp": "2025-01-08T23:00:00Z"
}
```

### ✅ Gefixte Endpoint Werkt

```bash
# GET /api/participant/emails response:
HTTP/1.1 200 OK
Content-Type: application/json

{
  "participant_emails": [...],
  "system_emails": [...],
  "all_emails": [...],
  "counts": {
    "participants": 5,
    "system": 2,
    "total": 7
  }
}
```

---

## 🔧 Troubleshooting

### Container Crasht Direct Na Start

**Symptoom:**
```
dkl-email-service   Restarting (1) X seconds ago
```

**Oorzaak:** Ontbrekende environment variabelen

**Oplossing:**
1. Check logs: `docker-compose logs app --tail=50`
2. Zoek naar: `"ontbrekende omgevingsvariabele: X"`
3. Voeg ontbrekende variabele toe aan `.env`
4. Restart: `docker-compose restart app`

### Database Connection Error

**Symptoom:**
```
"Database initialisatie fout"
```

**Oplossing:**
1. Check of postgres container draait: `docker-compose ps`
2. Verify DB credentials in `.env`
3. Wait for postgres health check: `docker-compose up -d` (wait ~10 seconds)

### Port Already in Use

**Symptoom:**
```
Error: bind: address already in use
```

**Oplossing:**
```bash
# Stop andere services op port 8080
# Op Windows:
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# Of wijzig port in docker-compose.yml:
ports:
  - "8081:8080"  # Use 8081 instead
```

### Cannot Connect to Email Server

**Symptoom:**
```
"Fout bij verzenden email"
```

**Oplossing:**
Dit is verwacht als je geen echte SMTP credentials hebt.
- Voor local testing: gebruik test mode in API calls
- Voor echte emails: configureer Gmail App Password of andere SMTP

---

## 📝 Quick Reference Commands

```bash
# Start services
docker-compose up -d

# Stop services
docker-compose down

# Rebuild after code changes
docker-compose up --build -d

# View logs
docker-compose logs -f app

# Check status
docker-compose ps

# Execute command in container
docker-compose exec app /bin/sh

# Restart specific service
docker-compose restart app

# Remove all containers and volumes
docker-compose down -v
```

---

## 🎯 Verificatie Checklist

- [ ] `.env` bestand aangemaakt met alle vereiste variabelen
- [ ] Docker containers draaien zonder crashes
- [ ] Health endpoint returnt 200 OK
- [ ] `/api/participant/emails` returnt 200 OK (met auth token)
- [ ] `/api/participant/:id` werkt nog steeds correct
- [ ] `/api/under-construction/active` returnt 404 (expected behavior)
- [ ] Backend logs tonen geen FATAL errors

---

## 📚 Related Documentation

- [Backend Fix Details](./BACKEND_FIX_PARTICIPANT_EMAILS.md)
- [Backend API Errors Analysis](./api/BACKEND_API_ERRORS.md)
- [Frontend Guide](./frontend/BACKEND_API_ERRORS_FRONTEND.md)
- [Setup Guide](./guides/SETUP.md)

---

## 💡 Next Steps

Na succesvolle local deployment:

1. ✅ Test alle endpoints met Postman/curl
2. 📝 Update frontend team dat fix beschikbaar is
3. 🚀 Deploy naar staging environment
4. ✅ Run integration tests
5. 🎉 Deploy naar production
