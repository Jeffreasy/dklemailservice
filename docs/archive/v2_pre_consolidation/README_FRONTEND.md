# 📚 Frontend Documentatie Index

Welkom bij de DKL Email Service backend API documentatie voor frontend developers.

## 🚀 Quick Start

**Je backend draait al lokaal!** Volg deze 3 simpele stappen:

1. **📖 Lees:** [`FRONTEND_SETUP_QUICK.md`](./FRONTEND_SETUP_QUICK.md) (5 minuten)
2. **📋 Kopieer:** [`frontend-env-example.txt`](./frontend-env-example.txt) naar je frontend project
3. **💻 Code:** Gebruik [`frontend-api-client-example.ts`](./frontend-api-client-example.ts) als basis

## 📑 Documentatie Overzicht

### Voor Beginners
| Document | Tijd | Beschrijving |
|----------|------|--------------|
| [`FRONTEND_CONNECTION_SUMMARY.md`](./FRONTEND_CONNECTION_SUMMARY.md) | 2 min | Overzicht en status |
| [`FRONTEND_SETUP_QUICK.md`](./FRONTEND_SETUP_QUICK.md) | 5 min | 3-stappen setup |

### Uitgebreide Guides
| Document | Tijd | Beschrijving |
|----------|------|--------------|
| [`FRONTEND_LOCAL_DEVELOPMENT.md`](./FRONTEND_LOCAL_DEVELOPMENT.md) | 15 min | Complete development guide |
| [`FRONTEND_API_REFERENCE.md`](./FRONTEND_API_REFERENCE.md) | 20 min | Alle API endpoints |

### Code Voorbeelden
| Bestand | Type | Beschrijving |
|---------|------|--------------|
| [`frontend-env-example.txt`](./frontend-env-example.txt) | Config | Environment variabelen |
| [`frontend-api-client-example.ts`](./frontend-api-client-example.ts) | Code | TypeScript API client |

## 🎯 Belangrijkste URLs

### Backend API
```
Lokaal:    http://localhost:8082/api
Productie: https://dklemailservice.onrender.com/api
```

### Services (Lokaal)
```
API:        localhost:8082
PostgreSQL: localhost:5433
Redis:      localhost:6380
```

## 🔑 Essentiele Info

### Environment Setup (Frontend)

**Development (.env.development):**
```env
VITE_API_BASE_URL=http://localhost:8082/api
```

**Production (.env.production):**
```env
VITE_API_BASE_URL=https://dklemailservice.onrender.com/api
```

### Authenticatie Flow

```typescript
// 1. Login
const response = await api.auth.login(email, password);
// Response: { token, refresh_token, user }

// 2. Token wordt automatisch toegevoegd aan requests
const contacts = await api.contacts.list();
// Header: Authorization: Bearer <token>
```

### Populaire Endpoints

**Dashboard Data:**
```
GET /api/contact           - Contactformulieren
GET /api/aanmelding        - Aanmeldingen  
GET /api/users             - Gebruikers
```

**Content Management:**
```
GET /api/albums/admin      - Albums
GET /api/photos/admin      - Foto's
GET /api/videos/admin      - Video's
```

**Nieuw - Steps Tracking:**
```
POST /api/steps/:id        - Update steps
GET  /api/total-steps      - Totaal steps
```

## 📊 Current Status

### Backend Status: ✅ RUNNING

```bash
# Verify
curl http://localhost:8082/api/health

# Services
✅ PostgreSQL - Up and healthy
✅ Redis - Up and healthy  
✅ API - Up and running
```

### Features Beschikbaar

- ✅ Authenticatie (JWT)
- ✅ Contact & Aanmelding beheer
- ✅ Album & Foto management
- ✅ Video & Sponsor management
- ✅ Steps tracking (NIEUW!)
- ✅ User management
- ✅ Newsletter system
- ✅ Chat system
- ✅ RBAC permissions

## 🔄 Development vs Production

### Zo Wissel Je

**Optie 1: Environment Files (Aanbevolen)**
```bash
# Development
npm run dev
# Gebruikt .env.development → localhost:8082

# Production build
npm run build
# Gebruikt .env.production → dklemailservice.onrender.com
```

**Optie 2: NPM Scripts**
```json
{
  "scripts": {
    "dev": "vite",
    "dev:prod": "vite --mode production"
  }
}
```

### Waarom Dit Veilig Is

1. **Gescheiden Databases**
   - Lokaal: Eigen Docker database met test data
   - Productie: Render database met echte data
   - Geen overlap!

2. **Gescheiden URLs**
   - `.env.development` → localhost
   - `.env.production` → render.com
   - Framework kiest automatisch juiste

3. **Geen Impact op Productie**
   - Lokale wijzigingen raken productie niet
   - Test alles veilig lokaal
   - Deploy via git/CI-CD naar productie

## 🛠️ Development Commands

### Backend (Docker)
```bash
# Start
docker-compose -f docker-compose.dev.yml up -d

# Status
docker-compose -f docker-compose.dev.yml ps

# Logs
docker-compose -f docker-compose.dev.yml logs -f app

# Stop
docker-compose -f docker-compose.dev.yml down

# Reset (fresh database)
docker-compose -f docker-compose.dev.yml down -v
docker-compose -f docker-compose.dev.yml up -d
```

### Frontend
```bash
# Development (wijst naar localhost:8082)
npm run dev

# Production build (wijst naar render.com)
npm run build
```

## ✅ Checklist voor Nieuwe Frontend Dev

- [ ] Backend draait lokaal (docker-compose up -d)
- [ ] Frontend .env.development aangemaakt
- [ ] API client geconfigureerd
- [ ] Test login werkt
- [ ] Kan data ophalen van backend
- [ ] Wijzigingen worden niet naar productie gepusht

## 🎓 Volgende Stappen

1. **Lees** [`FRONTEND_SETUP_QUICK.md`](./FRONTEND_SETUP_QUICK.md) - 5 minuten
2. **Kopieer** [`frontend-api-client-example.ts`](./frontend-api-client-example.ts) naar je project
3. **Test** login en data ophalen
4. **Ontwikkel** met confidence - productie blijft veilig!

## 💡 Pro Tips

**Tip 1: API Status Check**
```bash
# Check of backend bereikbaar is
curl http://localhost:8082/api/health
```

**Tip 2: Browser Console Test**
```javascript
// Test API vanaf browser
fetch('http://localhost:8082/api/health')
  .then(r => r.json())
  .then(console.log)
```

**Tip 3: Check Environment**
```javascript
// Welke API URL gebruikt mijn app?
console.log(import.meta.env.VITE_API_BASE_URL)
```

**Tip 4: CORS Issues?**
- Check of je frontend draait op poort 3000 of 5173
- Andere poort? Voeg toe aan backend ALLOWED_ORIGINS

**Tip 5: Database Reset**
```bash
# Verse start met nieuwe test data
docker-compose -f docker-compose.dev.yml down -v
docker-compose -f docker-compose.dev.yml up -d
```

## 📞 Support

**Backend Logs Checken:**
```bash
docker-compose -f docker-compose.dev.yml logs app | grep ERROR
```

**Database Bekijken:**
```bash
docker exec -it dkl-postgres psql -U postgres -d dklemailservice
```

**API Testen:**
- Postman collection beschikbaar
- Curl voorbeelden in documentatie
- Health endpoint: `http://localhost:8082/api/health`

---

**Last Updated:** 1 November 2025  
**Backend Version:** 1.1.0  
**Status:** ✅ Lokaal draaiend op :8082  
**Production:** https://dklemailservice.onrender.com