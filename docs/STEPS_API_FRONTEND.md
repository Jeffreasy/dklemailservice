# Steps API - Frontend Documentatie

## Overzicht

Het Steps systeem is een stappen-tracking en fondsenallocatie systeem voor wandelevenementen. Deelnemers kunnen hun gewandelde stappen registreren, hun voortgang volgen, en zien hoeveel fondsen er zijn toegewezen aan hun gekozen route.

## Kernfunctionaliteiten

1. **Stappen Tracking**: Deelnemers kunnen hun stappen bijwerken en volgen
2. **Dashboard**: Persoonlijk dashboard met stappen, route en toegewezen fondsen
3. **Totalen**: Overzicht van totaal aantal stappen van alle deelnemers
4. **Fondsverdeling**: Inzicht in hoe fondsen verdeeld zijn over routes
5. **Admin Beheer**: Beheer van route fondsallocaties (alleen admin/staff)

---

## Data Modellen

### Aanmelding (Participant)
```typescript
interface Aanmelding {
  id: string;                    // UUID
  naam: string;                  // Naam deelnemer
  email: string;                 // Email adres
  telefoon?: string;             // Telefoonnummer
  rol?: string;                  // Rol (bijv. "speler", "begeleider")
  afstand: string;               // Gekozen route ("6 KM", "10 KM", "15 KM", "20 KM")
  ondersteuning?: string;        // Benodigde ondersteuning
  bijzonderheden?: string;       // Extra opmerkingen
  terms: boolean;                // Akkoord met voorwaarden
  status: string;                // Status ("nieuw", "bevestigd", etc.)
  steps: number;                 // Aantal gewandelde stappen (default: 0)
  gebruiker_id?: string;         // Link naar gebruikersaccount
  created_at: string;            // Timestamp van aanmelding
  updated_at: string;            // Laatste update timestamp
  test_mode: boolean;            // Test modus indicator
}
```

### RouteFund
```typescript
interface RouteFund {
  id: string;                    // UUID
  route: string;                 // Route naam ("6 KM", "10 KM", "15 KM", "20 KM")
  amount: number;                // Toegewezen bedrag in euro's
  created_at: string;            // Timestamp
  updated_at: string;            // Laatste update timestamp
}
```

### Dashboard Response
```typescript
interface ParticipantDashboard {
  steps: number;                 // Aantal gewandelde stappen
  route: string;                 // Gekozen route
  allocatedFunds: number;        // Toegewezen fondsen in euro's
  naam: string;                  // Naam deelnemer
  email: string;                 // Email deelnemer
}
```

---

## API Endpoints

### Basis URL
```
Base URL: https://your-api-domain.com/api
```

Alle endpoints vereisen authenticatie via Bearer Token in de Authorization header.

---

## 1. Stappen Bijwerken

### Voor Ingelogde Gebruiker (Eigen Stappen)
```http
POST /api/steps
Authorization: Bearer {token}
Content-Type: application/json

{
  "steps": 1000  // Delta: aantal stappen om toe te voegen (kan negatief zijn)
}
```

**Response:**
```json
{
  "id": "uuid",
  "naam": "Jan Jansen",
  "email": "jan@example.com",
  "afstand": "10 KM",
  "steps": 5000,  // Nieuwe totaal
  "created_at": "2025-01-15T10:00:00Z",
  "updated_at": "2025-01-15T12:00:00Z"
  // ... andere velden
}
```

**Logica:**
- De `steps` parameter is een **delta** (toevoeging)
- Als je 1000 stuurt en de gebruiker had 4000, krijgt hij 5000
- Negatieve waarden zijn toegestaan (om correcties te maken)
- Het systeem voorkomt automatisch dat stappen onder 0 komen
- Gebruikt `gebruiker_id` uit JWT token om automatisch juiste deelnemer te vinden

**Frontend Voorbeeld:**
```typescript
async function updateMySteps(deltaSteps: number) {
  const response = await fetch('/api/steps', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${getToken()}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ steps: deltaSteps })
  });
  
  if (!response.ok) throw new Error('Failed to update steps');
  return await response.json();
}

// Gebruik:
await updateMySteps(500);  // Voeg 500 stappen toe
```

---

### Voor Admin/Staff (Specifieke Deelnemer)
```http
POST /api/steps/{participant_id}
Authorization: Bearer {admin_token}
Content-Type: application/json

{
  "steps": 1000
}
```

**Vereiste Permissie:** `steps:write`

---

## 2. Dashboard Ophalen

### Voor Ingelogde Gebruiker (Eigen Dashboard)
```http
GET /api/participant/dashboard
Authorization: Bearer {token}
```

**Response:**
```json
{
  "steps": 5000,
  "route": "10 KM",
  "allocatedFunds": 75,
  "naam": "Jan Jansen",
  "email": "jan@example.com"
}
```

**Frontend Voorbeeld:**
```typescript
interface DashboardData {
  steps: number;
  route: string;
  allocatedFunds: number;
  naam: string;
  email: string;
}

async function getMyDashboard(): Promise<DashboardData> {
  const response = await fetch('/api/participant/dashboard', {
    headers: {
      'Authorization': `Bearer ${getToken()}`
    }
  });
  
  if (!response.ok) throw new Error('Failed to fetch dashboard');
  return await response.json();
}
```

**UI Suggesties:**
```tsx
function ParticipantDashboard() {
  const [dashboard, setDashboard] = useState<DashboardData | null>(null);
  
  useEffect(() => {
    getMyDashboard().then(setDashboard);
  }, []);
  
  if (!dashboard) return <Loading />;
  
  return (
    <div className="dashboard">
      <h2>Welkom, {dashboard.naam}!</h2>
      <div className="stats">
        <StatCard 
          title="Stappen" 
          value={dashboard.steps.toLocaleString()} 
          icon="🚶"
        />
        <StatCard 
          title="Route" 
          value={dashboard.route} 
          icon="🗺️"
        />
        <StatCard 
          title="Toegewezen Fondsen" 
          value={`€${dashboard.allocatedFunds}`} 
          icon="💰"
        />
      </div>
    </div>
  );
}
```

---

### Voor Admin/Staff (Specifieke Deelnemer Dashboard)
```http
GET /api/participant/{participant_id}/dashboard
Authorization: Bearer {admin_token}
```

**Vereiste Permissie:** `steps:read`

---

## 3. Totaal Stappen (Alle Deelnemers)

```http
GET /api/total-steps?year=2025
Authorization: Bearer {token}
```

**Query Parameters:**
- `year` (optional): Jaar voor filtering (default: 2025)

**Response:**
```json
{
  "total_steps": 125000
}
```

**Vereiste Permissie:** `steps:read_total`

**Frontend Voorbeeld:**
```typescript
async function getTotalSteps(year: number = 2025): Promise<number> {
  const response = await fetch(`/api/total-steps?year=${year}`, {
    headers: {
      'Authorization': `Bearer ${getToken()}`
    }
  });
  
  if (!response.ok) throw new Error('Failed to fetch total steps');
  const data = await response.json();
  return data.total_steps;
}

// Animated counter component
function TotalStepsCounter() {
  const [total, setTotal] = useState(0);
  
  useEffect(() => {
    getTotalSteps().then(setTotal);
    
    // Refresh elke 30 seconden
    const interval = setInterval(() => {
      getTotalSteps().then(setTotal);
    }, 30000);
    
    return () => clearInterval(interval);
  }, []);
  
  return (
    <div className="total-steps">
      <h3>Totaal Gewandelde Stappen</h3>
      <div className="counter">{total.toLocaleString()}</div>
      <p className="subtitle">door alle deelnemers samen</p>
    </div>
  );
}
```

---

## 4. Fondsverdeling

```http
GET /api/funds-distribution
Authorization: Bearer {token}
```

**Response:**
```json
{
  "totalX": 10000,  // Totaal beschikbaar bedrag
  "routes": {
    "6 KM": 2500,   // Bedrag voor 6 KM route
    "10 KM": 2500,  // Bedrag voor 10 KM route
    "15 KM": 2500,  // Bedrag voor 15 KM route
    "20 KM": 2500   // Bedrag voor 20 KM route
  }
}
```

**Vereiste Permissie:** `steps:read`

**Logica:**
- Verdeling is **proportioneel** gebaseerd op aantal deelnemers per route
- Als er 40 deelnemers zijn (10 per route), krijgt elke route 25% van totaal
- Als één route meer deelnemers heeft, krijgt die meer fondsen

**Frontend Voorbeeld:**
```typescript
interface FundsDistribution {
  totalX: number;
  routes: Record<string, number>;
}

async function getFundsDistribution(): Promise<FundsDistribution> {
  const response = await fetch('/api/funds-distribution', {
    headers: {
      'Authorization': `Bearer ${getToken()}`
    }
  });
  
  if (!response.ok) throw new Error('Failed to fetch distribution');
  return await response.json();
}

// Visualisatie met chart
function FundsChart() {
  const [distribution, setDistribution] = useState<FundsDistribution | null>(null);
  
  useEffect(() => {
    getFundsDistribution().then(setDistribution);
  }, []);
  
  if (!distribution) return <Loading />;
  
  return (
    <div className="funds-chart">
      <h3>Fondsverdeling per Route</h3>
      <p>Totaal: €{distribution.totalX.toLocaleString()}</p>
      
      <div className="chart">
        {Object.entries(distribution.routes).map(([route, amount]) => (
          <div key={route} className="bar">
            <div className="label">{route}</div>
            <div 
              className="bar-fill" 
              style={{ 
                width: `${(amount / distribution.totalX) * 100}%` 
              }}
            >
              €{amount}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
```

---

## 5. Admin Endpoints - Route Funds Beheer

### Alle Route Funds Ophalen
```http
GET /api/steps/admin/route-funds
Authorization: Bearer {admin_token}
```

**Response:**
```json
[
  {
    "id": "uuid-1",
    "route": "6 KM",
    "amount": 50,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  },
  {
    "id": "uuid-2",
    "route": "10 KM",
    "amount": 75,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
  // etc...
]
```

**Vereiste Permissie:** `steps:write`

---

### Nieuwe Route Fund Aanmaken
```http
POST /api/steps/admin/route-funds
Authorization: Bearer {admin_token}
Content-Type: application/json

{
  "route": "25 KM",
  "amount": 150
}
```

**Validatie:**
- `route`: verplicht, uniek
- `amount`: verplicht, >= 0

**Response:** Status 201 + created RouteFund object

---

### Route Fund Bijwerken
```http
PUT /api/steps/admin/route-funds/{route}
Authorization: Bearer {admin_token}
Content-Type: application/json

{
  "amount": 100
}
```

**Response:** Updated RouteFund object

---

### Route Fund Verwijderen
```http
DELETE /api/steps/admin/route-funds/{route}
Authorization: Bearer {admin_token}
```

**Response:**
```json
{
  "success": true,
  "message": "Route fund succesvol verwijderd"
}
```

---

### Admin Panel Voorbeeld
```typescript
function RouteFundsAdmin() {
  const [funds, setFunds] = useState<RouteFund[]>([]);
  
  async function loadFunds() {
    const response = await fetch('/api/steps/admin/route-funds', {
      headers: { 'Authorization': `Bearer ${getToken()}` }
    });
    setFunds(await response.json());
  }
  
  async function updateFund(route: string, amount: number) {
    await fetch(`/api/steps/admin/route-funds/${route}`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${getToken()}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ amount })
    });
    loadFunds(); // Refresh
  }
  
  useEffect(() => { loadFunds(); }, []);
  
  return (
    <div className="admin-funds">
      <h2>Route Fondsen Beheer</h2>
      <table>
        <thead>
          <tr>
            <th>Route</th>
            <th>Bedrag (€)</th>
            <th>Acties</th>
          </tr>
        </thead>
        <tbody>
          {funds.map(fund => (
            <tr key={fund.id}>
              <td>{fund.route}</td>
              <td>
                <input 
                  type="number" 
                  value={fund.amount}
                  onChange={(e) => updateFund(fund.route, +e.target.value)}
                />
              </td>
              <td>
                <button onClick={() => deleteFund(fund.route)}>
                  Verwijderen
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
```

---

## Authenticatie & Autorisatie

### JWT Token
Alle endpoints vereisen een geldig JWT token in de Authorization header:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### User Context
Bij authenticatie worden deze velden automatisch uit het token gehaald:
- `userID`: UUID van de ingelogde gebruiker
- Dit wordt gebruikt om automatisch de juiste `Aanmelding` te vinden via `gebruiker_id`

### Permissies

| Endpoint | Permissie | Beschrijving |
|----------|-----------|--------------|
| `POST /api/steps` | `steps:write` | Eigen stappen bijwerken |
| `POST /api/steps/:id` | `steps:write` | Admin: stappen van anderen bijwerken |
| `GET /api/participant/dashboard` | `steps:read` | Eigen dashboard bekijken |
| `GET /api/participant/:id/dashboard` | `steps:read` | Admin: dashboard van anderen bekijken |
| `GET /api/total-steps` | `steps:read_total` | Totaal stappen zien (alle gebruikers) |
| `GET /api/funds-distribution` | `steps:read` | Fondsverdeling zien |
| `GET /api/steps/admin/route-funds` | `steps:write` | Admin: route funds beheren |
| `POST /api/steps/admin/route-funds` | `steps:write` | Admin: route fund aanmaken |
| `PUT /api/steps/admin/route-funds/:route` | `steps:write` | Admin: route fund bijwerken |
| `DELETE /api/steps/admin/route-funds/:route` | `steps:write` | Admin: route fund verwijderen |

---

## Error Handling

### Standaard Error Response
```json
{
  "error": "Beschrijving van de fout"
}
```

### HTTP Status Codes

| Code | Betekenis | Mogelijke Oorzaken |
|------|-----------|-------------------|
| 200 | OK | Succesvolle request |
| 201 | Created | Resource succesvol aangemaakt |
| 400 | Bad Request | Ongeldige input data |
| 401 | Unauthorized | Geen of ongeldig token |
| 403 | Forbidden | Geen permissie voor deze actie |
| 404 | Not Found | Resource niet gevonden |
| 500 | Internal Server Error | Server fout |

### Frontend Error Handling
```typescript
async function apiRequest<T>(
  url: string, 
  options?: RequestInit
): Promise<T> {
  const response = await fetch(url, {
    ...options,
    headers: {
      'Authorization': `Bearer ${getToken()}`,
      'Content-Type': 'application/json',
      ...options?.headers
    }
  });
  
  if (!response.ok) {
    const error = await response.json();
    
    switch (response.status) {
      case 401:
        // Redirect naar login
        window.location.href = '/login';
        break;
      case 403:
        throw new Error('Geen toegang tot deze functie');
      case 404:
        throw new Error('Niet gevonden');
      default:
        throw new Error(error.error || 'Er ging iets mis');
    }
  }
  
  return await response.json();
}
```

---

## Complete Frontend Integration Voorbeeld

```typescript
// api/steps.ts
import { apiRequest } from './utils';

export interface StepsDashboard {
  steps: number;
  route: string;
  allocatedFunds: number;
  naam: string;
  email: string;
}

export const stepsApi = {
  // Eigen stappen bijwerken
  updateMySteps: (deltaSteps: number) =>
    apiRequest<Aanmelding>('/api/steps', {
      method: 'POST',
      body: JSON.stringify({ steps: deltaSteps })
    }),
  
  // Dashboard ophalen
  getMyDashboard: () =>
    apiRequest<StepsDashboard>('/api/participant/dashboard'),
  
  // Totaal stappen
  getTotalSteps: (year: number = 2025) =>
    apiRequest<{ total_steps: number }>(`/api/total-steps?year=${year}`),
  
  // Fondsverdeling
  getFundsDistribution: () =>
    apiRequest<{ totalX: number; routes: Record<string, number> }>(
      '/api/funds-distribution'
    ),
  
  // Admin functies
  admin: {
    getRouteFunds: () =>
      apiRequest<RouteFund[]>('/api/steps/admin/route-funds'),
    
    createRouteFund: (route: string, amount: number) =>
      apiRequest<RouteFund>('/api/steps/admin/route-funds', {
        method: 'POST',
        body: JSON.stringify({ route, amount })
      }),
    
    updateRouteFund: (route: string, amount: number) =>
      apiRequest<RouteFund>(`/api/steps/admin/route-funds/${route}`, {
        method: 'PUT',
        body: JSON.stringify({ amount })
      }),
    
    deleteRouteFund: (route: string) =>
      apiRequest(`/api/steps/admin/route-funds/${route}`, {
        method: 'DELETE'
      })
  }
};

// components/StepsTracker.tsx
function StepsTracker() {
  const [dashboard, setDashboard] = useState<StepsDashboard | null>(null);
  const [totalSteps, setTotalSteps] = useState(0);
  const [inputSteps, setInputSteps] = useState('');
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    loadData();
    const interval = setInterval(loadData, 60000); // Refresh elke minuut
    return () => clearInterval(interval);
  }, []);

  async function loadData() {
    const [dashboardData, total] = await Promise.all([
      stepsApi.getMyDashboard(),
      stepsApi.getTotalSteps()
    ]);
    setDashboard(dashboardData);
    setTotalSteps(total.total_steps);
  }

  async function addSteps() {
    if (!inputSteps) return;
    
    setLoading(true);
    try {
      await stepsApi.updateMySteps(parseInt(inputSteps));
      setInputSteps('');
      await loadData(); // Refresh data
      toast.success('Stappen toegevoegd!');
    } catch (error) {
      toast.error(error.message);
    } finally {
      setLoading(false);
    }
  }

  if (!dashboard) return <Spinner />;

  return (
    <div className="steps-tracker">
      <section className="personal-stats">
        <h2>Jouw Statistieken</h2>
        <div className="stats-grid">
          <StatCard 
            icon="🚶" 
            title="Jouw Stappen" 
            value={dashboard.steps.toLocaleString()}
            subtitle={`Route: ${dashboard.route}`}
          />
          <StatCard 
            icon="💰" 
            title="Toegewezen Fondsen" 
            value={`€${dashboard.allocatedFunds}`}
            subtitle="Voor jouw route"
          />
          <StatCard 
            icon="🌍" 
            title="Totaal Stappen" 
            value={totalSteps.toLocaleString()}
            subtitle="Alle deelnemers"
          />
        </div>
      </section>

      <section className="add-steps">
        <h3>Stappen Toevoegen</h3>
        <div className="input-group">
          <input
            type="number"
            value={inputSteps}
            onChange={(e) => setInputSteps(e.target.value)}
            placeholder="Aantal stappen..."
            min="0"
          />
          <button 
            onClick={addSteps} 
            disabled={loading || !inputSteps}
          >
            {loading ? 'Bezig...' : 'Toevoegen'}
          </button>
        </div>
        <p className="hint">
          Tip: Gebruik je fitness tracker om je dagelijkse stappen te zien
        </p>
      </section>

      <section className="progress">
        <h3>Voortgang</h3>
        <ProgressBar 
          current={dashboard.steps} 
          goal={100000} 
          label="Persoonlijk doel"
        />
      </section>
    </div>
  );
}

export default StepsTracker;
```

---

## Database Schema

### Aanmeldingen Tabel
```sql
CREATE TABLE aanmeldingen (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  naam VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  telefoon VARCHAR(50),
  rol VARCHAR(100),
  afstand VARCHAR(50),  -- Route: "6 KM", "10 KM", "15 KM", "20 KM"
  ondersteuning TEXT,
  bijzonderheden TEXT,
  terms BOOLEAN NOT NULL,
  status VARCHAR(50) DEFAULT 'nieuw',
  steps INTEGER DEFAULT 0,  -- Stappen counter
  gebruiker_id UUID REFERENCES gebruikers(id),  -- Link naar user account
  test_mode BOOLEAN DEFAULT false
);

CREATE INDEX idx_aanmeldingen_email ON aanmeldingen(email);
CREATE INDEX idx_aanmeldingen_status ON aanmeldingen(status);
CREATE INDEX idx_aanmeldingen_gebruiker_id ON aanmeldingen(gebruiker_id);
```

### Route Funds Tabel
```sql
CREATE TABLE route_funds (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  route VARCHAR(50) NOT NULL UNIQUE,  -- "6 KM", "10 KM", etc.
  amount INTEGER NOT NULL CHECK (amount >= 0),  -- Bedrag in euro's
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_route_funds_route ON route_funds(route);

-- Default waarden
INSERT INTO route_funds (route, amount) VALUES
  ('6 KM', 50),
  ('10 KM', 75),
  ('15 KM', 100),
  ('20 KM', 125);
```

---

## Business Logica

### Stappen Berekening
```go
// Delta toevoegen aan bestaande stappen
newSteps := participant.Steps + deltaSteps

// Voorkom negatieve stappen
if newSteps < 0 {
    newSteps = 0
}

participant.Steps = newSteps
```

### Fondsen Allocatie
```go
// Bereken toegewezen fondsen voor een route
func CalculateAllocatedFunds(route string) int {
    // Haal route fund op uit database
    routeFund := GetRouteFundByRoute(route)
    return routeFund.Amount
}

// Standaard waarden als route niet gevonden:
// - 6 KM:  €50
// - 10 KM: €75
// - 15 KM: €100
// - 20 KM: €125
```

### Proportionele Verdeling
```go
// Verdeel fondsen proportioneel over routes
func GetFundsDistributionProportional() {
    totalParticipants := CountAllParticipants()
    
    for each route {
        participantCount := CountParticipantsForRoute(route)
        routeFund := GetRouteFundAmount(route)
        
        // Proportioneel berekenen
        allocation := (routeFund * participantCount) / totalParticipants
    }
}
```

**Voorbeeld:**
- Totaal fondsen: €10,000
- Totaal deelnemers: 100
- Route "10 KM": 30 deelnemers
- Route "10 KM" fondsallocatie: €75 per deelnemer
- Berekening: (75 * 30) / 100 = €22.50 per deelnemer voor deze route

---

## Testing

### Test Endpoints met cURL

```bash
# Get dashboard
curl -X GET "https://api.example.com/api/participant/dashboard" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Update steps
curl -X POST "https://api.example.com/api/steps" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"steps": 1000}'

# Get total steps
curl -X GET "https://api.example.com/api/total-steps?year=2025" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Get funds distribution
curl -X GET "https://api.example.com/api/funds-distribution" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## Best Practices

### 1. Real-time Updates
```typescript
// Gebruik websockets of polling voor live updates
const useRealTimeSteps = () => {
  const [totalSteps, setTotalSteps] = useState(0);
  
  useEffect(() => {
    // Initial load
    stepsApi.getTotalSteps().then(data => setTotalSteps(data.total_steps));
    
    // Poll elke 30 seconden voor updates
    const interval = setInterval(async () => {
      const data = await stepsApi.getTotalSteps();
      setTotalSteps(data.total_steps);
    }, 30000);
    
    return () => clearInterval(interval);
  }, []);
  
  return totalSteps;
};
```

### 2. Optimistic Updates
```typescript
async function addStepsOptimistic(delta: number) {
  // Update UI immediate
  setDashboard(prev => ({
    ...prev!,
    steps: prev!.steps + delta
  }));
  
  try {
    // Send to server
    await stepsApi.updateMySteps(delta);
  } catch (error) {
    // Rollback on error
    setDashboard(prev => ({
      ...prev!,
      steps: prev!.steps - delta
    }));
    toast.error('Kon stappen niet opslaan');
  }
}
```

### 3. Caching
```typescript
// Cache dashboard data
const useDashboard = () => {
  return useQuery('dashboard', stepsApi.getMyDashboard, {
    staleTime: 60000, // 1 minuut
    cacheTime: 300000 // 5 minuten
  });
};
```

### 4. Error Boundaries
```typescript
class StepsErrorBoundary extends React.Component {
  componentDidCatch(error: Error) {
    if (error.message.includes('401')) {
      // Redirect naar login
      window.location.href = '/login';
    } else {
      // Show error message
      toast.error('Er ging iets mis met het laden van je stappen');
    }
  }
  
  render() {
    return this.props.children;
  }
}
```

---

## Veelgestelde Vragen

### Q: Kan ik stappen verwijderen?
**A:** Ja, stuur een negatieve delta. Bijvoorbeeld `{"steps": -500}` haalt 500 stappen af.

### Q: Wat gebeurt er als ik meer stappen aftrek dan ik heb?
**A:** Het systeem voorkomt automatisch dat stappen onder 0 komen. Als je 100 stappen hebt en 200 probeert af te halen, wordt het 0.

### Q: Hoe worden fondsen berekend?
**A:** Fondsen worden proportioneel verdeeld gebaseerd op aantal deelnemers per route. Routes met meer deelnemers krijgen een groter aandeel van het totale budget.

### Q: Kan ik de route van een deelnemer wijzigen?
**A:** De route staat vast in de `aanmeldingen` tabel. Om dit te wijzigen moet je de `Aanmelding` record updaten (via andere API).

### Q: Worden stappen automatisch gesynchroniseerd met fitness trackers?
**A:** Nee, de API biedt alleen de backend. De frontend moet integreren met fitness apps (Strava, Garmin, etc.) en vervolgens stappen versturen naar deze API.

### Q: Is er een maximum aantal stappen?
**A:** Nee, er is geen maximum. Het veld is een integer die tot ~2 miljard kan.

---

## Support & Contact

Voor vragen of problemen met de Steps API:
- **Backend Team**: backend@example.com
- **API Documentatie**: https://api.example.com/docs
- **Status Page**: https://status.example.com

---

## Changelog

### v1.0 (2025-01-15)
- ✅ Basis stappen tracking
- ✅ Dashboard functionaliteit
- ✅ Fondsverdeling
- ✅ Admin route funds beheer
- ✅ JWT authenticatie
- ✅ RBAC permissies

### Toekomstige Features
- 🔄 Websocket real-time updates
- 🔄 Fitness tracker integraties (Strava, Garmin)
- 🔄 Leaderboard
- 🔄 Team/groep functionaliteit
- 🔄 Export naar CSV/Excel
- 🔄 Grafische voortgangsrapportages