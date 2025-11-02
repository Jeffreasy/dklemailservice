# 🚀 Frontend Quickstart

**Backend API draait NU lokaal en is klaar voor connectie!**

## ✅ Status Check

```bash
Backend API:  http://localhost:8082     ✓ RUNNING
PostgreSQL:   localhost:5433            ✓ HEALTHY
Redis:        localhost:6380            ✓ HEALTHY
```

## 📝 Voor Je Frontend Team

### Stap 1: Kopieer Environment Config

Maak in je **FRONTEND** project `.env.development`:

```env
VITE_API_BASE_URL=http://localhost:8082/api
```

### Stap 2: Gebruik API Client

```typescript
import axios from 'axios';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 30000,
});

// Login voorbeeld
const login = async (email: string, password: string) => {
  const res = await api.post('/auth/login', { email, password });
  localStorage.setItem('auth_token', res.data.token);
  return res.data;
};

// Fetch data met auth
api.interceptors.request.use(config => {
  const token = localStorage.getItem('auth_token');
  if (token) config.headers.Authorization = `Bearer ${token}`;
  return config;
});

// Gebruik
const contacts = await api.get('/contact');
```

### Stap 3: Start Frontend

```bash
npm run dev
```

**Klaar!** Frontend verbindt nu met lokale backend.

## 🔄 Productie vs Development

### Development (Lokaal)
```env
# .env.development
VITE_API_BASE_URL=http://localhost:8082/api
```
👉 Test veilig lokaal, geen impact op productie

### Production (Live)
```env
# .env.production  
VITE_API_BASE_URL=https://dklemailservice.onrender.com/api
```
👉 Wijst naar echte productie backend

## 📚 Volledige Documentatie

Alles staat in [`docs/`](./docs/) directory:

- **Quick:** [`docs/FRONTEND_SETUP_QUICK.md`](./docs/FRONTEND_SETUP_QUICK.md)
- **Complete:** [`docs/FRONTEND_LOCAL_DEVELOPMENT.md`](./docs/FRONTEND_LOCAL_DEVELOPMENT.md)
- **API:** [`docs/FRONTEND_API_REFERENCE.md`](./docs/FRONTEND_API_REFERENCE.md)
- **Code:** [`docs/frontend-api-client-example.ts`](./docs/frontend-api-client-example.ts)

## 🎯 Belangrijkste Endpoints

```typescript
// Auth
POST /api/auth/login              // Login
GET  /api/auth/profile            // Get user info

// Admin Data
GET  /api/contact                 // Contact formulieren
GET  /api/aanmelding              // Aanmeldingen
GET  /api/users                   // Gebruikers

// Content
GET  /api/albums/admin            // Albums
GET  /api/photos/admin            // Foto's
GET  /api/videos/admin            // Video's

// Steps (Nieuw!)
POST /api/steps/:id               // Update steps
GET  /api/total-steps             // Totalen
```

## ✨ Nieuwe Features

- 🎯 **Steps Tracking** - Deelnemers kunnen stappen registreren
- 💰 **Route Funds** - Fondsen verdeling per route
- ⚡ **Database Optimalisaties** - Snellere queries

## 🐛 Troubleshoot

**CORS Error?**
→ Check of frontend op poort 3000 of 5173 draait

**Connection Refused?**
→ `docker-compose -f docker-compose.dev.yml up -d`

**401 Unauthorized?**
→ Login opnieuw, token expired

---

**🎉 Je bent klaar om te ontwikkelen!**

Productie blijft veilig - alle wijzigingen zijn lokaal.