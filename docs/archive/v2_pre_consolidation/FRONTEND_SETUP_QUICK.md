# Frontend Setup - Quick Reference

Snelle handleiding voor het verbinden van de frontend met de lokale backend.

## 🎯 In 3 Stappen Klaar

### 1️⃣ Start Backend Lokaal

```bash
# In backend directory (dklemailservice)
docker-compose -f docker-compose.dev.yml up -d
```

✅ Backend draait nu op: **http://localhost:8082**

### 2️⃣ Configureer Frontend

Maak in je **frontend** project een `.env.development` bestand:

**Voor Vite/React:**
```env
VITE_API_BASE_URL=http://localhost:8082/api
VITE_WS_URL=ws://localhost:8082/api/chat/ws
```

**Voor Next.js:**
```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8082/api
NEXT_PUBLIC_WS_URL=ws://localhost:8082/api/chat/ws
```

**Voor Create React App:**
```env
REACT_APP_API_BASE_URL=http://localhost:8082/api
REACT_APP_WS_URL=ws://localhost:8082/api/chat/ws
```

### 3️⃣ Update API Client

**src/config/api.ts** of **src/utils/api.ts**:

```typescript
import axios from 'axios';

// Haal base URL uit environment
const baseURL = import.meta.env.VITE_API_BASE_URL || 
                process.env.NEXT_PUBLIC_API_BASE_URL ||
                process.env.REACT_APP_API_BASE_URL ||
                'http://localhost:8082/api';

// Maak axios instance
export const apiClient = axios.create({
  baseURL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Voeg auth token toe aan requests
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('auth_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Handle expired tokens
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('auth_token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);
```

## ✅ Verifieer Setup

### Test 1: Backend Health

```bash
curl http://localhost:8082/api/health
```

Verwacht: `{"status":"healthy",...}`

### Test 2: Login Test

```typescript
// In browser console of component
const response = await fetch('http://localhost:8082/api/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: 'admin@dekoninklijkeloop.nl',
    password: 'jouw_wachtwoord'
  })
});

const data = await response.json();
console.log(data); // Moet token bevatten
```

### Test 3: Authenticated Request

```typescript
// Na login
const token = localStorage.getItem('auth_token');

const response = await fetch('http://localhost:8082/api/contact', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
});

const contacts = await response.json();
console.log(contacts); // Moet lijst van contacts tonen
```

## 🔄 Switchen tussen Local en Productie

### Optie 1: Via NPM Scripts (Aanbevolen)

**package.json:**
```json
{
  "scripts": {
    "dev": "vite",
    "dev:prod-api": "vite --mode production-api",
    "build": "vite build",
    "build:prod": "vite build --mode production"
  }
}
```

Dan:
```bash
# Lokale backend
npm run dev

# Productie backend (voor testing)
npm run dev:prod-api
```

### Optie 2: Environment Files

Maak meerdere env files:

- `.env.development` → Lokale backend
- `.env.production` → Productie backend
- `.env.staging` → Staging backend (indien aanwezig)

Vite/Next.js kiest automatisch de juiste op basis van `NODE_ENV`.

### Optie 3: Runtime Toggle (Advanced)

```typescript
// src/config/api-config.ts
export const API_ENDPOINTS = {
  local: 'http://localhost:8082/api',
  production: 'https://api.dekoninklijkeloop.nl/api',
  staging: 'https://staging-api.dekoninklijkeloop.nl/api', // indien beschikbaar
};

// Lees uit env of localStorage
const savedEndpoint = localStorage.getItem('api_endpoint');
const envEndpoint = import.meta.env.VITE_API_BASE_URL;

export const currentEndpoint = savedEndpoint || envEndpoint || API_ENDPOINTS.local;

// UI component voor switchen
function ApiSwitcher() {
  const [endpoint, setEndpoint] = useState(currentEndpoint);
  
  const switchAPI = (newEndpoint: string) => {
    localStorage.setItem('api_endpoint', newEndpoint);
    setEndpoint(newEndpoint);
    window.location.reload(); // Reload om nieuwe endpoint te gebruiken
  };
  
  return (
    <select value={endpoint} onChange={(e) => switchAPI(e.target.value)}>
      <option value={API_ENDPOINTS.local}>Local (8082)</option>
      <option value={API_ENDPOINTS.production}>Production</option>
    </select>
  );
}
```

## 🔑 Belangrijke Endpoints voor Admin Panel

### Dashboard Data
```typescript
// Contactformulieren
GET /api/contact                       // Alle contactformulieren
GET /api/contact/status/nieuw          // Nieuwe contactformulieren
GET /api/contact/:id                   // Specifiek contactformulier
POST /api/contact/:id/antwoord         // Antwoord toevoegen

// Aanmeldingen
GET /api/aanmelding                    // Alle aanmeldingen
GET /api/aanmelding/rol/vrijwilliger   // Filter op rol
GET /api/aanmelding/:id                // Specifieke aanmelding
POST /api/aanmelding/:id/antwoord      // Antwoord toevoegen

// Gebruikers
GET /api/users                         // Alle gebruikers
GET /api/users/:id                     // Specifieke gebruiker
POST /api/users/:id/roles              // Rol toewijzen

// Nieuwsbrieven
GET /api/newsletter                    // Alle nieuwsbrieven
POST /api/newsletter                   // Nieuwe nieuwsbrief
POST /api/newsletter/:id/send          // Nieuwsbrief verzenden
```

### Content Management
```typescript
// Albums
GET /api/albums/admin                  // Alle albums
POST /api/albums                       // Album aanmaken
PUT /api/albums/:id                    // Album updaten
DELETE /api/albums/:id                 // Album verwijderen

// Photos
GET /api/photos/admin                  // Alle foto's
POST /api/photos                       // Foto uploaden
GET /api/photos?year=2024              // Filter op jaar

// Videos
GET /api/videos/admin                  // Alle videos
POST /api/videos                       // Video toevoegen

// Sponsors
GET /api/sponsors/admin                // Alle sponsors
POST /api/sponsors                     // Sponsor toevoegen
```

### Steps Tracking (Nieuw!)
```typescript
// Deelnemer dashboard
GET /api/participant/:id/dashboard     // Dashboard data voor deelnemer

// Steps updaten
POST /api/steps/:id                    // Update steps voor deelnemer
Body: { "steps": 10000 }

// Totalen ophalen
GET /api/total-steps                   // Totaal aantal steps dit jaar
GET /api/funds-distribution            // Fondsen verdeling
```

## 🐛 Debug Checklist

Gebruik deze checklist als iets niet werkt:

### Backend Issues
```bash
# 1. Is backend actief?
docker-compose -f docker-compose.dev.yml ps
# Verwacht: dkl-email-service status "Up"

# 2. Draaien alle services?
docker-compose -f docker-compose.dev.yml ps
# Verwacht: app, postgres, redis allemaal "Up" of "healthy"

# 3. Zijn er errors?
docker-compose -f docker-compose.dev.yml logs app | grep "ERROR\|FATAL"

# 4. Is de poort juist?
curl http://localhost:8082/api/health
# Verwacht: JSON response met status
```

### Frontend Issues
```bash
# 1. Juiste API URL?
console.log(import.meta.env.VITE_API_BASE_URL)
# Verwacht: "http://localhost:8082/api"

# 2. CORS errors in console?
# Check of frontend poort in ALLOWED_ORIGINS staat

# 3. Network tab in DevTools
# Bekijk request/response voor details
```

### Authentication Issues
```bash
# 1. Token aanwezig?
console.log(localStorage.getItem('auth_token'))

# 2. Token geldig?
# Decode JWT op https://jwt.io

# 3. Login werkt?
# Test login endpoint met Postman eerst
```

## 💡 Tips

1. **Browser Extension: ModHeader**
   - Voeg automatisch auth headers toe voor testing

2. **React DevTools**
   - Bekijk component state en props
   - Debug waarom data niet laadt

3. **Network Throttling**
   - Test hoe app reageert op langzame connectie
   - Belangrijk voor loading states

4. **Clear Storage**
   ```javascript
   // Als je vreemde auth issues hebt
   localStorage.clear();
   sessionStorage.clear();
   ```

5. **Proxy Alternative (indien CORS issues)**
   ```javascript
   // vite.config.ts
   export default defineConfig({
     server: {
       proxy: {
         '/api': {
           target: 'http://localhost:8082',
           changeOrigin: true,
         }
       }
     }
   })
   ```

## 📞 Hulp Nodig?

1. Check [`FRONTEND_LOCAL_DEVELOPMENT.md`](./FRONTEND_LOCAL_DEVELOPMENT.md) voor uitgebreide info
2. Bekijk [`FRONTEND_API_REFERENCE.md`](./FRONTEND_API_REFERENCE.md) voor API details
3. Check backend logs met `docker-compose logs`
4. Test endpoints met Postman/curl eerst