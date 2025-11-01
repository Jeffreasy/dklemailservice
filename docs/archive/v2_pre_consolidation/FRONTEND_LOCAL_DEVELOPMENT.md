# Frontend Lokale Development Guide

Deze handleiding beschrijft hoe de frontend (adminpanel) lokaal kan verbinden met de backend API tijdens development, zonder impact op productie.

## 🏗️ Architectuur Overzicht

### Productie Setup
```
Frontend (https://admin.dekoninklijkeloop.nl)
    ↓
Backend API (https://dklemailservice.onrender.com)
    ↓
Production Database (Render PostgreSQL)
```

### Lokale Development Setup
```
Frontend Lokaal (http://localhost:3000 of http://localhost:5173)
    ↓
Backend API Lokaal (http://localhost:8082)
    ↓
Local Docker Database (PostgreSQL in container op poort 5433)
```

## 🚀 Quick Start

### 1. Start Lokale Backend (Docker)

```bash
# In de backend directory (dklemailservice)
docker-compose -f docker-compose.dev.yml up -d
```

Dit start:
- **PostgreSQL**: `localhost:5433` (lokale database met test data)
- **Redis**: `localhost:6380` (cache en session management)
- **API**: `http://localhost:8082/api` (backend API)

### 2. Verifieer Backend Status

```bash
# Test health endpoint
curl http://localhost:8082/api/health

# Test root endpoint
curl http://localhost:8082/
```

## 🔧 Frontend Configuratie

### Voor React/Vite Projecten

Maak een `.env.development` bestand in je **frontend** project:

```env
# Lokale development API
VITE_API_BASE_URL=http://localhost:8082/api
VITE_API_TIMEOUT=30000
VITE_ENV=development

# WebSocket voor chat (indien nodig)
VITE_WS_URL=ws://localhost:8082/api/chat/ws
```

Voor **productie** gebruik je `.env.production`:

```env
# Productie API
VITE_API_BASE_URL=https://dklemailservice.onrender.com/api
VITE_API_TIMEOUT=30000
VITE_ENV=production

# WebSocket voor chat
VITE_WS_URL=wss://dklemailservice.onrender.com/api/chat/ws
```

### Voor Next.js Projecten

`.env.local`:
```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8082/api
NEXT_PUBLIC_ENV=development
```

`.env.production`:
```env
NEXT_PUBLIC_API_BASE_URL=https://dklemailservice.onrender.com/api
NEXT_PUBLIC_ENV=production
```

### API Client Setup (React Voorbeeld)

```typescript
// src/config/api.ts
export const API_CONFIG = {
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8082/api',
  timeout: parseInt(import.meta.env.VITE_API_TIMEOUT || '30000'),
  headers: {
    'Content-Type': 'application/json',
  },
};

// src/services/api.ts
import axios from 'axios';
import { API_CONFIG } from '../config/api';

const apiClient = axios.create(API_CONFIG);

// Add JWT token to requests
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('auth_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Handle 401 errors (token expired)
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Redirect to login
      localStorage.removeItem('auth_token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default apiClient;
```

## 🔐 Authenticatie Flow

### 1. Login Request

```typescript
// Login functie
async function login(email: string, password: string) {
  const response = await apiClient.post('/auth/login', {
    email,
    password
  });
  
  // Sla token op
  localStorage.setItem('auth_token', response.data.token);
  localStorage.setItem('refresh_token', response.data.refresh_token);
  
  return response.data;
}
```

### 2. Test Credentials (Lokaal)

**Standaard Admin Account:**
- Email: `admin@dekoninklijkeloop.nl`
- Password: Controleer de database migraties voor het gehasht wachtwoord

Om zelf een test account aan te maken:
```bash
# In de backend directory
docker exec -it dkl-postgres psql -U postgres -d dklemailservice

-- In de psql console:
INSERT INTO gebruikers (naam, email, wachtwoord_hash, rol, is_actief)
VALUES ('Lokale Admin', 'admin@localhost.dev', '$2a$10$[hash]', 'admin', true);
```

## 📡 API Endpoints Overzicht

### Publieke Endpoints (geen auth vereist)
```
GET  /api/health                    - Health check
POST /api/contact-email             - Contact formulier
POST /api/aanmelding-email          - Aanmelding formulier
POST /api/auth/login                - Login
GET  /api/albums                    - Publieke albums
GET  /api/videos                    - Publieke videos
GET  /api/sponsors                  - Publieke sponsors
```

### Admin Endpoints (JWT auth vereist)
```
GET    /api/contact                 - Lijst contactformulieren
GET    /api/aanmelding              - Lijst aanmeldingen
GET    /api/users                   - Gebruikersbeheer
GET    /api/newsletter              - Nieuwsbrief beheer
GET    /api/albums/admin            - Album beheer
POST   /api/photos                  - Foto upload
DELETE /api/videos/:id              - Video verwijderen
```

Zie [`docs/FRONTEND_API_REFERENCE.md`](./FRONTEND_API_REFERENCE.md) voor complete API documentatie.

## 🔄 Productie vs Development Switchen

### In de Frontend

**Development modus** (wijst naar localhost):
```bash
npm run dev
# of
yarn dev
```

**Productie build** (wijst naar productie API):
```bash
npm run build
# of
yarn build
```

### Best Practice: Environment Detection

```typescript
// src/config/environment.ts
export const ENV = {
  isDevelopment: import.meta.env.DEV,
  isProduction: import.meta.env.PROD,
  apiBaseURL: import.meta.env.VITE_API_BASE_URL,
};

// Gebruik in code
if (ENV.isDevelopment) {
  console.log('Running in development mode');
  console.log('API:', ENV.apiBaseURL);
}
```

## 🛡️ CORS Configuratie

De backend is **al geconfigureerd** voor lokale development:

**Toegestane Origins** (in [`main.go:307-310`](../main.go:307)):
- `http://localhost:3000` ✓ (React/Next.js default)
- `http://localhost:5173` ✓ (Vite default)
- `https://admin.dekoninklijkeloop.nl` ✓ (Productie admin)
- `https://www.dekoninklijkeloop.nl` ✓ (Productie website)

**Als je een andere poort gebruikt**, voeg deze toe aan `.env`:
```env
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173,http://localhost:4200
```

En herstart Docker:
```bash
docker-compose -f docker-compose.dev.yml down
docker-compose -f docker-compose.dev.yml up -d
```

## 📝 Development Workflow

### Dagelijkse Workflow

**1. Start Backend:**
```bash
# In backend directory
docker-compose -f docker-compose.dev.yml up -d

# Check status
docker-compose -f docker-compose.dev.yml ps

# View logs (optioneel)
docker-compose -f docker-compose.dev.yml logs -f app
```

**2. Start Frontend:**
```bash
# In frontend directory
npm run dev
# Frontend draait nu op http://localhost:5173
```

**3. Test Verbinding:**
```bash
# Test of frontend kan connecten
curl -X POST http://localhost:8082/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@dekoninklijkeloop.nl","password":"jouw_wachtwoord"}'
```

**4. Development:**
- Frontend wijzigingen worden automatisch heet herladen
- Backend wijzigingen vereisen rebuild:
  ```bash
  docker-compose -f docker-compose.dev.yml build app
  docker-compose -f docker-compose.dev.yml up -d
  ```

**5. Afsluiten:**
```bash
# Stop backend (data blijft behouden)
docker-compose -f docker-compose.dev.yml down

# Stop backend en verwijder data (clean slate)
docker-compose -f docker-compose.dev.yml down -v
```

## 🗄️ Database Management

### Lokale Database Beheren

**Verbind met lokale database:**
```bash
# Via command line
docker exec -it dkl-postgres psql -U postgres -d dklemailservice

# Of via pgAdmin/DBeaver:
Host: localhost
Port: 5433
User: postgres
Password: postgres
Database: dklemailservice
```

### Test Data

De lokale database wordt automatisch gevuld met test data via migraties:
- Test gebruikers met verschillende rollen
- Test contactformulieren
- Test aanmeldingen
- Test albums, foto's, video's

### Reset Database

Volledig verse start:
```bash
docker-compose -f docker-compose.dev.yml down -v
docker-compose -f docker-compose.dev.yml up -d
```

## 🔍 Debugging Tips

### Network Requests Debuggen

**Browser DevTools:**
```javascript
// In browser console
fetch('http://localhost:8082/api/health')
  .then(r => r.json())
  .then(console.log)
```

**CORS Errors:**
```
Access to XMLHttpRequest blocked by CORS policy
```
**Oplossing:** Check of frontend origin in `ALLOWED_ORIGINS` staat

**Connection Errors:**
```
ERR_CONNECTION_REFUSED
```
**Oplossing:** Backend is niet gestart, run `docker-compose -f docker-compose.dev.yml up -d`

### Backend Logs Bekijken

```bash
# Real-time logs
docker-compose -f docker-compose.dev.yml logs -f app

# Laatste 100 regels
docker-compose -f docker-compose.dev.yml logs --tail=100 app

# Filter op error
docker-compose -f docker-compose.dev.yml logs app | grep ERROR
```

## 🎯 Frontend Specifieke Endpoints

### Admin Dashboard Data

**Contact Overzicht:**
```typescript
// Haal alle contactformulieren op
const contacts = await apiClient.get('/contact');

// Filter op status
const newContacts = await apiClient.get('/contact/status/nieuw');
```

**Aanmelding Overzicht:**
```typescript
// Haal alle aanmeldingen op
const registrations = await apiClient.get('/aanmelding');

// Filter op rol
const volunteers = await apiClient.get('/aanmelding/rol/vrijwilliger');
```

### Content Management

**Albums:**
```typescript
// Publieke albums (voor main site)
const albums = await apiClient.get('/albums');

// Alle albums (voor admin panel)
const allAlbums = await apiClient.get('/albums/admin');

// Album aanmaken
await apiClient.post('/albums', {
  title: 'Nieuw Album',
  description: 'Beschrijving',
  visible: true
});
```

**Steps Tracking:**
```typescript
// Update participant steps
await apiClient.post(`/steps/${participantId}`, {
  steps: 10000
});

// Get participant dashboard
const dashboard = await apiClient.get(`/participant/${participantId}/dashboard`);

// Get total steps
const totalSteps = await apiClient.get('/total-steps');
```

## 🔄 Synchronisatie met Productie

### Data Synchronisatie

**Productie data naar lokaal** (LET OP: vereist access tot productie database):

```bash
# Dump productie data
pg_dump -h production-host -U username -d dklemailservice > prod_dump.sql

# Importeer in lokale database
docker exec -i dkl-postgres psql -U postgres -d dklemailservice < prod_dump.sql
```

**Lokale wijzigingen naar productie:**
- **NOOIT** direct naar productie pushen
- Altijd eerst via git en CI/CD pipeline
- Test grondig in lokale omgeving

## 🧪 Testing met Frontend

### E2E Testing

```typescript
// cypress/e2e/admin.cy.ts
describe('Admin Panel', () => {
  beforeEach(() => {
    cy.visit('http://localhost:3000/admin');
    cy.login('admin@localhost.dev', 'password'); // Custom command
  });

  it('should load contact forms', () => {
    cy.get('[data-cy=contacts]').click();
    cy.url().should('include', '/admin/contacts');
    cy.contains('Contactformulieren').should('be.visible');
  });
});
```

## 📊 Monitoring During Development

### Check Backend Status

```bash
# Health check
curl http://localhost:8082/api/health | jq

# Check services
curl http://localhost:8082/api/health | jq '.checks'
```

Expected output:
```json
{
  "status": "healthy",
  "checks": {
    "redis": {"status": true},
    "smtp": {"default": false, "registration": false},
    "templates": {"status": true}
  }
}
```

### Monitor Logs

```bash
# Follow all logs
docker-compose -f docker-compose.dev.yml logs -f

# Only app logs
docker-compose -f docker-compose.dev.yml logs -f app

# Only database logs
docker-compose -f docker-compose.dev.yml logs -f postgres
```

## 🔐 Security Considerations

### Lokaal vs Productie

**Development (Lokaal):**
- HTTP toegestaan (http://localhost:8082)
- Zwakke JWT secret voor snelheid
- Debug logging enabled
- SMTP validation disabled (dummy credentials)
- CORS permissief voor localhost

**Productie:**
- HTTPS vereist
- Sterke JWT secrets
- Minimale logging
- Volledig geconfigureerde SMTP
- Strenge CORS policy

### Secrets Management

**NOOIT** in git:
- `.env` files met echte credentials
- Productie API keys
- Database passwords
- JWT secrets

**WEL** in git:
- `.env.example` met placeholder waarden
- Documentatie van vereiste variabelen

## 🚨 Common Issues & Solutions

### Issue 1: CORS Error
```
Access-Control-Allow-Origin error
```
**Oplossing:**
1. Check of frontend draait op `localhost:3000` of `localhost:5173`
2. Zo niet, voeg je poort toe aan `ALLOWED_ORIGINS` in backend `.env`
3. Herstart backend: `docker-compose -f docker-compose.dev.yml restart app`

### Issue 2: Connection Refused
```
ERR_CONNECTION_REFUSED
```
**Oplossing:**
1. Check of backend draait: `docker-compose -f docker-compose.dev.yml ps`
2. Start backend: `docker-compose -f docker-compose.dev.yml up -d`

### Issue 3: 401 Unauthorized
```
401 Unauthorized
```
**Oplossing:**
1. Token verlopen - login opnieuw
2. Geen token - voeg `Authorization: Bearer <token>` header toe
3. Ongeldige credentials - check gebruiker in database

### Issue 4: Database Connection Failed
```
Database connection failed
```
**Oplossing:**
1. Check PostgreSQL status: `docker-compose -f docker-compose.dev.yml ps postgres`
2. Check logs: `docker-compose -f docker-compose.dev.yml logs postgres`
3. Herstart: `docker-compose -f docker-compose.dev.yml restart postgres`

## 📱 API Testing Tools

### Postman Collection

Importeer de volgende base configuration in Postman:

**Environment: Development**
```json
{
  "name": "DKL Development",
  "values": [
    {
      "key": "base_url",
      "value": "http://localhost:8082/api",
      "enabled": true
    },
    {
      "key": "auth_token",
      "value": "",
      "enabled": true
    }
  ]
}
```

**Environment: Production**
```json
{
  "name": "DKL Production",
  "values": [
    {
      "key": "base_url",
      "value": "https://dklemailservice.onrender.com/api",
      "enabled": true
    },
    {
      "key": "auth_token",
      "value": "",
      "enabled": true
    }
  ]
}
```

### Request Examples

**Login:**
```http
POST {{base_url}}/auth/login
Content-Type: application/json

{
  "email": "admin@dekoninklijkeloop.nl",
  "password": "your_password"
}
```

**Get Contacts (with auth):**
```http
GET {{base_url}}/contact
Authorization: Bearer {{auth_token}}
```

## 🔄 Hot Reload Development

### Backend Changes

Als je wijzigingen maakt in de backend code:

**Optie 1: Rebuild (langzaam maar volledig):**
```bash
docker-compose -f docker-compose.dev.yml build app
docker-compose -f docker-compose.dev.yml up -d
```

**Optie 2: Direct Go Development (snel):**
```bash
# Stop Docker backend
docker-compose -f docker-compose.dev.yml stop app

# Update .env voor lokaal draaien
# DB_HOST=localhost (wijzig van 'postgres' naar 'localhost')
# DB_PORT=5433 (wijzig van 5432 naar 5433)
# REDIS_HOST=localhost
# REDIS_PORT=6380

# Run direct met Go
go run main.go

# Backend draait nu op localhost:8080 (niet 8082!)
# Update frontend: VITE_API_BASE_URL=http://localhost:8080/api
```

### Frontend Changes

Frontend hot reload werkt automatisch met Vite/React:
```bash
npm run dev
# Wijzigingen worden automatisch herladen
```

## 📦 Docker Commands Cheat Sheet

```bash
# Status bekijken
docker-compose -f docker-compose.dev.yml ps

# Alle logs
docker-compose -f docker-compose.dev.yml logs

# Specifieke service logs
docker-compose -f docker-compose.dev.yml logs app
docker-compose -f docker-compose.dev.yml logs postgres
docker-compose -f docker-compose.dev.yml logs redis

# Service herstarten
docker-compose -f docker-compose.dev.yml restart app

# Service stoppen
docker-compose -f docker-compose.dev.yml stop app

# Service starten
docker-compose -f docker-compose.dev.yml start app

# Alles stoppen
docker-compose -f docker-compose.dev.yml down

# Alles stoppen + volumes verwijderen (fresh start)
docker-compose -f docker-compose.dev.yml down -v

# Rebuild en starten
docker-compose -f docker-compose.dev.yml up -d --build

# Shell in container
docker exec -it dkl-email-service sh
docker exec -it dkl-postgres psql -U postgres -d dklemailservice
docker exec -it dkl-redis redis-cli
```

## 🎨 Frontend Development Best Practices

### API Error Handling

```typescript
// src/utils/api-error-handler.ts
export function handleApiError(error: any) {
  if (error.response) {
    // Server responded with error
    switch (error.response.status) {
      case 400:
        return 'Ongeldige invoer';
      case 401:
        return 'Niet geautoriseerd - log opnieuw in';
      case 403:
        return 'Geen toegang';
      case 404:
        return 'Niet gevonden';
      case 500:
        return 'Server error';
      default:
        return error.response.data?.error || 'Onbekende fout';
    }
  } else if (error.request) {
    // Request made but no response
    return 'Geen verbinding met server';
  } else {
    // Error in request setup
    return error.message;
  }
}
```

### Loading States

```typescript
const [loading, setLoading] = useState(false);
const [error, setError] = useState<string | null>(null);

async function loadData() {
  setLoading(true);
  setError(null);
  
  try {
    const response = await apiClient.get('/contact');
    setData(response.data);
  } catch (err) {
    setError(handleApiError(err));
  } finally {
    setLoading(false);
  }
}
```

## 📋 Checklist voor Nieuwe Frontend Developer

- [ ] Backend draait lokaal in Docker
- [ ] Frontend environment variabelen geconfigureerd
- [ ] API client setup met auth interceptors
- [ ] Test login werkt
- [ ] CORS configuratie correct
- [ ] Error handling geïmplementeerd
- [ ] Loading states geïmplementeerd
- [ ] API documentatie gelezen

## 🆘 Support

### Backend Logs Checken

```bash
# Check of migraties succesvol waren
docker-compose -f docker-compose.dev.yml logs app | grep "Migratie"

# Check of server draait
docker-compose -f docker-compose.dev.yml logs app | grep "Server gestart"

# Check voor errors
docker-compose -f docker-compose.dev.yml logs app | grep "ERROR\|FATAL"
```

### Database Queries

```sql
-- Check gebruikers
SELECT id, naam, email, rol, is_actief FROM gebruikers;

-- Check contact formulieren
SELECT id, naam, email, status, created_at FROM contact_formulieren ORDER BY created_at DESC LIMIT 10;

-- Check aanmeldingen
SELECT id, naam, email, rol, status, created_at FROM aanmeldingen ORDER BY created_at DESC LIMIT 10;
```

## 🎓 Aanvullende Resources

- [API Reference](./FRONTEND_API_REFERENCE.md) - Complete API documentatie
- [Authentication](../api/authentication.md) - Auth flow details
- [Development Guide](./development.md) - Backend development
- [Testing Guide](./testing.md) - Test procedures

## ✨ Pro Tips

1. **Gebruik Browser Extensions:**
   - ModHeader voor custom headers
   - JSON Viewer voor mooie JSON responses
   - React Developer Tools voor React debugging

2. **API Response Caching:**
   - Implementeer caching in frontend voor betere UX
   - Gebruik React Query of SWR voor data fetching

3. **Type Safety:**
   - Genereer TypeScript types van API responses
   - Gebruik tools zoals openapi-typescript

4. **Environment Switching:**
   - Maak een UI toggle voor easy switching tussen dev/prod API
   - Handig voor quick production checks

5. **Mock Data:**
   - Voor pure frontend development, gebruik MSW (Mock Service Worker)
   - Onafhankelijk van backend status

---

**Laatste update:** November 2025
**Auteur:** DKL Development Team