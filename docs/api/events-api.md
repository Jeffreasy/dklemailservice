# Events API Reference

Complete API referentie voor de event tracking functionaliteit in de DKL Email Service.

## Overzicht

De Events API stelt de applicatie in staat om loopwedstrijd events te beheren met real-time geofencing en tracking capabilities. Het systeem ondersteunt meerdere events, geofence configuratie, en participant tracking.

## Base URL

```
Production: https://dklemailservice.onrender.com
Development: https://dklemailservice.onrender.com
```

## Endpoints

### GET /api/events/:id

Haalt gedetailleerde informatie op voor een specifiek event, inclusief starttijd en geofence configuratie.

**URL Parameters:**
- `id`: Event ID (UUID)

**Response (200 OK):**
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "De Koninklijke Loop 2025",
    "description": "De jaarlijkse Koninklijke Loop hardloopevenement",
    "start_time": "2025-05-16T09:00:00Z",
    "end_time": "2025-05-16T16:00:00Z",
    "status": "upcoming",
    "geofences": [
        {
            "type": "start",
            "lat": 52.0907,
            "long": 5.1214,
            "radius": 50,
            "name": "Start Locatie"
        },
        {
            "type": "checkpoint",
            "lat": 52.0950,
            "long": 5.1300,
            "radius": 30,
            "name": "Checkpoint 5KM"
        },
        {
            "type": "finish",
            "lat": 52.0907,
            "long": 5.1214,
            "radius": 50,
            "name": "Finish Locatie"
        }
    ],
    "event_config": {
        "minStepsInterval": 10,
        "requireGeofenceCheckin": true,
        "distanceThreshold": 100,
        "accuracyLevel": "balanced"
    },
    "is_active": true
}
```

**Response (404 Not Found):**
```json
{
    "error": "Event niet gevonden",
    "code": "NOT_FOUND"
}
```

**Permissions:** PUBLIC (geen authenticatie vereist)

**Implementatie:** [`handlers/event_handler.go:58`](../../handlers/event_handler.go:58)

---

### GET /api/events/active

Haalt het eerstvolgende actieve event op.

**Response (200 OK):**
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "De Koninklijke Loop 2025",
    "start_time": "2025-05-16T09:00:00Z",
    "status": "upcoming",
    "geofences": [...],
    "is_active": true
}
```

**Response (404 Not Found):**
```json
{
    "error": "Geen actief event gevonden",
    "code": "NOT_FOUND"
}
```

**Permissions:** PUBLIC

**Implementatie:** [`handlers/event_handler.go:95`](../../handlers/event_handler.go:95)

---

### GET /api/events

Haalt lijst van alle events op met optionele filtering.

**Query Parameters:**
- `limit` (optioneel): Aantal resultaten (default: 50)
- `offset` (optioneel): Offset voor paginatie (default: 0)
- `active_only` (optioneel): Alleen actieve events (default: false)

**Voorbeelden:**
```bash
# Alle events
GET /api/events

# Alleen actieve events
GET /api/events?active_only=true

# Met paginatie
GET /api/events?limit=10&offset=20
```

**Response (200 OK):**
```json
[
    {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "De Koninklijke Loop 2025",
        "start_time": "2025-05-16T09:00:00Z",
        "status": "upcoming",
        "is_active": true
    },
    {
        "id": "550e8400-e29b-41d4-a716-446655440001",
        "name": "De Koninklijke Loop 2024",
        "start_time": "2024-05-17T09:00:00Z",
        "status": "completed",
        "is_active": false
    }
]
```

**Permissions:** PUBLIC

**Implementatie:** [`handlers/event_handler.go:125`](../../handlers/event_handler.go:125)

---

### POST /api/events

Maakt een nieuw event aan (admin only).

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "name": "De Koninklijke Loop 2026",
    "description": "Volgende editie van De Koninklijke Loop",
    "start_time": "2026-05-15T09:00:00Z",
    "end_time": "2026-05-15T16:00:00Z",
    "status": "upcoming",
    "geofences": [
        {
            "type": "start",
            "lat": 52.0907,
            "long": 5.1214,
            "radius": 50,
            "name": "Start Locatie"
        },
        {
            "type": "finish",
            "lat": 52.0907,
            "long": 5.1214,
            "radius": 50,
            "name": "Finish Locatie"
        }
    ],
    "event_config": {
        "minStepsInterval": 10,
        "requireGeofenceCheckin": true
    },
    "is_active": true
}
```

**Validatie:**
- `name`: Required, niet leeg
- `start_time`: Required, ISO 8601 format
- `geofences`: Required, minimaal 1 geofence
- `end_time`: Optioneel, ISO 8601 format
- `status`: Optioneel, moet één van: `upcoming`, `active`, `completed`, `cancelled`

**Response (201 Created):**
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440002",
    "name": "De Koninklijke Loop 2026",
    "start_time": "2026-05-15T09:00:00Z",
    "status": "upcoming",
    "geofences": [...],
    "is_active": true
}
```

**Permissions:** `events:write`

**Implementatie:** [`handlers/event_handler.go:171`](../../handlers/event_handler.go:171)

---

### PUT /api/events/:id

Werkt een bestaand event bij (admin only).

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "name": "De Koninklijke Loop 2025 - UPDATED",
    "status": "active",
    "geofences": [
        {
            "type": "start",
            "lat": 52.0910,
            "long": 5.1215,
            "radius": 60
        }
    ]
}
```

**Note:** Alle velden zijn optioneel. Alleen opgegeven velden worden bijgewerkt.

**Response (200 OK):**
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "De Koninklijke Loop 2025 - UPDATED",
    "start_time": "2025-05-16T09:00:00Z",
    "status": "active",
    "geofences": [...],
    "is_active": true
}
```

**Permissions:** `events:write`

**Implementatie:** [`handlers/event_handler.go:282`](../../handlers/event_handler.go:282)

---

### DELETE /api/events/:id

Verwijdert een event (admin only).

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Event succesvol verwijderd"
}
```

**Permissions:** `events:write`

**Implementatie:** [`handlers/event_handler.go:380`](../../handlers/event_handler.go:380)

---

### GET /api/events/:id/participants

Haalt alle participants voor een event op (authenticated).

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**
```json
[
    {
        "id": "650e8400-e29b-41d4-a716-446655440000",
        "event_id": "550e8400-e29b-41d4-a716-446655440000",
        "event_name": "De Koninklijke Loop 2025",
        "participant_id": "750e8400-e29b-41d4-a716-446655440000",
        "participant_name": "John Doe",
        "tracking_status": "registered",
        "registered_at": "2025-03-01T10:00:00Z",
        "total_distance": 0.0,
        "current_steps": 0
    }
]
```

**Permissions:** `events:read`

**Implementatie:** [`handlers/event_handler.go:433`](../../handlers/event_handler.go:433)

---

## Data Models

### Geofence Object

```typescript
interface Geofence {
    type: "start" | "checkpoint" | "finish";
    lat: number;      // Latitude (WGS84)
    long: number;     // Longitude (WGS84)
    radius: number;   // Radius in meters
    name?: string;    // Optional display name
}
```

**Geofence Types:**
- `start`: Start locatie waar event begint
- `checkpoint`: Tussenpunt langs de route
- `finish`: Finish locatie waar event eindigt

**Voorbeeld:**
```json
{
    "type": "start",
    "lat": 52.0907,
    "long": 5.1214,
    "radius": 50,
    "name": "Start - Utrechtseweg"
}
```

### Event Config Object

```typescript
interface EventConfig {
    minStepsInterval?: number;        // Min seconden tussen step updates
    requireGeofenceCheckin?: boolean; // Moet participant in start geofence zijn?
    distanceThreshold?: number;       // Max meters van route afwijking
    accuracyLevel?: "high" | "balanced" | "low"; // GPS accuracy level
    maxParticipants?: number;         // Maximum aantal deelnemers
    autoStart?: boolean;              // Auto-start tracking bij geofence entry
}
```

### Event Status

| Status | Beschrijving |
|--------|--------------|
| `upcoming` | Event is gepland maar nog niet gestart |
| `active` | Event is momenteel actief |
| `completed` | Event is afgelopen |
| `cancelled` | Event is geannuleerd |

### Tracking Status

| Status | Beschrijving |
|--------|--------------|
| `registered` | Participant is geregistreerd |
| `checked_in` | Participant heeft ingecheckt |
| `started` | Participant is gestart met lopen |
| `in_progress` | Participant is onderweg |
| `finished` | Participant heeft gefinisht |
| `dnf` | Did Not Finish |

---

## Business Logic

### Geofence Validatie

Het systeem valideert of een participant binnen een geofence is via de database functie `is_in_geofence()`:

```sql
SELECT is_in_geofence(
    52.0908,  -- participant latitude
    5.1215,   -- participant longitude
    '{"lat": 52.0907, "long": 5.1214, "radius": 50}'::jsonb
);
```

**Algoritme:**
- Gebruikt Haversine formule voor afstandsberekening
- Retourneert `true` als participant binnen radius is
- Accuracy: ~1 meter voor kleine afstanden (<1km)

### Event Lifecycle

```
upcoming → active → completed
    ↓
cancelled (kan op elk moment)
```

**Automatische Status Updates:**
- `upcoming` → `active`: Wanneer start_time is bereikt
- `active` → `completed`: Wanneer end_time is bereikt (of handmatig)

---

## RBAC Permissions

### Vereiste Permissions

| Endpoint | Permission | Rollen |
|----------|-----------|--------|
| GET /api/events | - | PUBLIC |
| GET /api/events/active | - | PUBLIC |
| GET /api/events/:id | - | PUBLIC |
| POST /api/events | `events:write` | Admin |
| PUT /api/events/:id | `events:write` | Admin |
| DELETE /api/events/:id | `events:write` | Admin |
| GET /api/events/:id/participants | `events:read` | Admin, Staff |

### Permission Setup

Permissions worden automatisch aangemaakt via migratie [`V1_53__add_events_table.sql`](../../database/migrations/V1_53__add_events_table.sql).

**Toegewezen aan:**
- **Admin**: `events:read`, `events:write`, `events:delete`, `event_tracking:read`, `event_tracking:write`
- **Staff**: `events:read`, `event_tracking:read`
- **Deelnemer**: `events:read`, `event_tracking:read`, `event_tracking:write`

---

## Database Schema

### Events Tabel

```sql
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE,
    status VARCHAR(50) NOT NULL DEFAULT 'upcoming',
    geofences JSONB NOT NULL DEFAULT '[]',
    event_config JSONB DEFAULT '{}',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by UUID
);
```

### Event Participants Tabel

```sql
CREATE TABLE event_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    participant_id UUID NOT NULL REFERENCES aanmeldingen(id) ON DELETE CASCADE,
    registered_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    check_in_time TIMESTAMP WITH TIME ZONE,
    start_time TIMESTAMP WITH TIME ZONE,
    finish_time TIMESTAMP WITH TIME ZONE,
    tracking_status VARCHAR(50) DEFAULT 'registered',
    last_location_update TIMESTAMP WITH TIME ZONE,
    total_distance DECIMAL(10,2) DEFAULT 0,
    current_steps INT DEFAULT 0,
    UNIQUE(event_id, participant_id)
);
```

---

## Frontend Integration

### React/TypeScript Voorbeeld

```typescript
import { useState, useEffect } from 'react';

interface Event {
    id: string;
    name: string;
    start_time: string;
    geofences: Geofence[];
    event_config: EventConfig;
}

interface Geofence {
    type: 'start' | 'checkpoint' | 'finish';
    lat: number;
    long: number;
    radius: number;
    name?: string;
}

interface EventConfig {
    minStepsInterval?: number;
    requireGeofenceCheckin?: boolean;
    distanceThreshold?: number;
    accuracyLevel?: 'high' | 'balanced' | 'low';
}

// Hook voor event data ophalen
function useEventDetails(eventId?: string) {
    const [event, setEvent] = useState<Event | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchEvent = async () => {
            try {
                // Haal actief event op als geen ID opgegeven
                const url = eventId 
                    ? `/api/events/${eventId}` 
                    : '/api/events/active';
                
                const response = await fetch(url);
                
                if (response.ok) {
                    const data = await response.json();
                    setEvent(data);
                } else if (response.status === 404) {
                    setError('Geen event gevonden');
                } else {
                    setError('Fout bij ophalen event');
                }
            } catch (err) {
                setError('Netwerk fout');
            } finally {
                setLoading(false);
            }
        };

        fetchEvent();
    }, [eventId]);

    return { event, loading, error };
}

// Component voorbeeld
function EventDetails({ eventId }: { eventId?: string }) {
    const { event, loading, error } = useEventDetails(eventId);

    if (loading) return <div>Loading...</div>;
    if (error) return <div>Error: {error}</div>;
    if (!event) return <div>Geen event gevonden</div>;

    return (
        <div className="event-details">
            <h2>{event.name}</h2>
            <p>Start: {new Date(event.start_time).toLocaleString('nl-NL')}</p>
            
            <h3>Geofences:</h3>
            <ul>
                {event.geofences.map((fence, i) => (
                    <li key={i}>
                        {fence.name || fence.type}: 
                        Lat {fence.lat}, Long {fence.long}, 
                        Radius {fence.radius}m
                    </li>
                ))}
            </ul>
            
            <h3>Configuratie:</h3>
            <pre>{JSON.stringify(event.event_config, null, 2)}</pre>
        </div>
    );
}
```

### React Native/Expo Integration

Voor mobile apps die geofencing willen implementeren:

```typescript
// In je mobile app (aparte repository)
import { useEventDetails } from './hooks/useEventDetails';
import * as Location from 'expo-location';

function EventTrackingScreen() {
    const { event } = useEventDetails(); // Fetch van API
    const [inGeofence, setInGeofence] = useState(false);
    
    useEffect(() => {
        if (!event) return;
        
        // Check if current location is in any geofence
        Location.watchPositionAsync(
            { accuracy: Location.Accuracy.Balanced },
            (location) => {
                const isIn = event.geofences.some(fence => 
                    isWithinRadius(
                        location.coords.latitude,
                        location.coords.longitude,
                        fence.lat,
                        fence.long,
                        fence.radius
                    )
                );
                setInGeofence(isIn);
            }
        );
    }, [event]);
    
    return (
        <View>
            <Text>Event: {event?.name}</Text>
            <Text>In Geofence: {inGeofence ? 'Ja' : 'Nee'}</Text>
        </View>
    );
}

// Helper functie voor distance berekening
function isWithinRadius(
    lat1: number, 
    lon1: number, 
    lat2: number, 
    lon2: number, 
    radius: number
): boolean {
    const R = 6371e3; // Earth radius in meters
    const φ1 = lat1 * Math.PI/180;
    const φ2 = lat2 * Math.PI/180;
    const Δφ = (lat2-lat1) * Math.PI/180;
    const Δλ = (lon2-lon1) * Math.PI/180;

    const a = Math.sin(Δφ/2) * Math.sin(Δφ/2) +
              Math.cos(φ1) * Math.cos(φ2) *
              Math.sin(Δλ/2) * Math.sin(Δλ/2);
    const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a));

    const distance = R * c; // Distance in meters
    return distance <= radius;
}
```

---

## Error Handling

### HTTP Status Codes

| Status | Betekenis |
|--------|-----------|
| 200 | Succes |
| 201 | Created (nieuw event aangemaakt) |
| 400 | Ongeldige request (validatie fout) |
| 401 | Niet geauthenticeerd |
| 403 | Geen toestemming (permission denied) |
| 404 | Event niet gevonden |
| 500 | Server fout |

### Error Response Format

```json
{
    "error": "Beschrijvende foutmelding",
    "code": "ERROR_CODE"
}
```

**Error Codes:**
- `MISSING_ID`: Event ID ontbreekt
- `MISSING_NAME`: Event naam ontbreekt
- `MISSING_START_TIME`: Start tijd ontbreekt
- `MISSING_GEOFENCES`: Geen geofences opgegeven
- `INVALID_REQUEST`: Request data kan niet worden ge-parsed
- `INVALID_TIME_FORMAT`: Tijd is niet in ISO 8601 format
- `NOT_FOUND`: Event niet gevonden
- `INTERNAL_ERROR`: Server fout

---

## Use Cases

### 1. Mobile App - Haal Actief Event Op

```typescript
// Fetch active event voor tracking
const response = await fetch('https://api.dekoninklijkeloop.nl/api/events/active');
const event = await response.json();

// Start geofence monitoring met event.geofences
// Check event.start_time om te weten wanneer te starten
```

### 2. Admin Panel - Nieuw Event Aanmaken

```typescript
const newEvent = {
    name: "De Koninklijke Loop 2026",
    start_time: "2026-05-15T09:00:00Z",
    geofences: [
        { type: "start", lat: 52.0907, long: 5.1214, radius: 50 }
    ]
};

const response = await fetch('/api/events', {
    method: 'POST',
    headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
    },
    body: JSON.stringify(newEvent)
});
```

### 3. Dashboard - Toon Event Participants

```typescript
const response = await fetch(`/api/events/${eventId}/participants`, {
    headers: { 'Authorization': `Bearer ${token}` }
});
const participants = await response.json();

// Toon lijst van participants met tracking status
```

---

## Integration met Steps API

De Events API integreert naadloos met de bestaande Steps API:

```typescript
// 1. Haal event details op
const event = await fetch('/api/events/active').then(r => r.json());

// 2. Check geofences via mobile app
// 3. Update steps wanneer in geofence
await fetch('/api/steps', {
    method: 'POST',
    headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
    },
    body: JSON.stringify({ steps: 1500 })
});

// 4. Dashboard toont gecombineerde data
const dashboard = await fetch('/api/participant/dashboard', {
    headers: { 'Authorization': `Bearer ${token}` }
}).then(r => r.json());
```

---

## Security Considerations

### Authentication
- Public endpoints: GET /api/events, GET /api/events/:id, GET /api/events/active
- Protected endpoints: Vereisen geldige JWT token
- Admin endpoints: Vereisen `events:write` permission

### Authorization
- RBAC systeem voor fine-grained access control
- Public kan events lezen (voor mobile app integration)
- Alleen admin kan events aanmaken/wijzigen
- Staff en admin kunnen participants bekijken

### Input Validation
- UUID validatie voor IDs
- ISO 8601 validatie voor tijden
- Geofence coordinaten validatie
- SQL injection preventie via parameterized queries

---

## Performance

### Database Indexes
- Index op `status` voor snelle filtering
- Index op `is_active` voor actieve events
- Index op `start_time` voor tijd-gebaseerde queries
- GIN index op `geofences` JSONB voor geofence queries

### Caching Strategy
- Event details kunnen worden gecached (TTL: 5 minuten)
- Active event kan worden gecached (TTL: 1 minuut)
- Invalideer cache bij event updates

### Optimization Tips
- Gebruik `active_only=true` parameter waar mogelijk
- Cache event details in mobile app
- Poll active event endpoint niet vaker dan 1x per minuut

---

## Testing

### cURL Voorbeelden

**Haal actief event op (Production):**
```bash
curl -X GET https://dklemailservice.onrender.com/api/events/active
```

**Haal actief event op (Local Docker):**
```bash
curl -X GET http://localhost:8082/api/events/active
```

**Haal specifiek event op:**
```bash
curl -X GET https://dklemailservice.onrender.com/api/events/550e8400-e29b-41d4-a716-446655440000
```

**Maak nieuw event aan (admin, Production):**
```bash
curl -X POST https://dklemailservice.onrender.com/api/events \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Event",
    "start_time": "2025-06-01T09:00:00Z",
    "geofences": [
      {
        "type": "start",
        "lat": 52.0907,
        "long": 5.1214,
        "radius": 50
      }
    ]
  }'
```

### Postman Collection

Importeer de test.http file in Postman voor complete test suite.

---

## Troubleshooting

### Common Issues

**404 - Geen actief event gevonden**
- Controleer of er een event is met `is_active = true` en `status IN ('upcoming', 'active')`
- Check database: `SELECT * FROM events WHERE is_active = true;`

**400 - Invalid time format**
- Gebruik ISO 8601 format: `2025-05-16T09:00:00Z`
- Include timezone (Z voor UTC of +01:00 voor Amsterdam)

**403 - Permission denied**
- Controleer user roles: `SELECT * FROM user_roles WHERE user_id = ?`
- Controleer permissions: `SELECT * FROM role_permissions WHERE role_id = ?`

### Debug Queries

```sql
-- Check events
SELECT id, name, status, start_time, is_active FROM events ORDER BY start_time DESC;

-- Check geofences voor event
SELECT name, geofences FROM events WHERE id = 'event-id-here';

-- Check participants voor event
SELECT * FROM event_participants_view WHERE event_id = 'event-id-here';

-- Vind actief event
SELECT get_active_event();
```

---

## Migration Guide

### Database Migration

1. Migratie wordt automatisch uitgevoerd bij applicatie start
2. Standaard event wordt aangemaakt: "De Koninklijke Loop 2025"
3. Permissions worden toegevoegd aan RBAC systeem

```bash
# Check migratie status
SELECT * FROM migraties WHERE filename LIKE '%V1_53%';

# Manual migration (indien nodig)
psql -d dkl_db -f database/migrations/V1_53__add_events_table.sql
```

### API Changes

**Nieuwe endpoints (backwards compatible):**
- `/api/events` - Event lijst
- `/api/events/active` - Actief event
- `/api/events/:id` - Event details
- `/api/events/:id/participants` - Participants lijst

**Geen breaking changes** - Bestaande functionaliteit blijft werken.

---

## WebSocket Integration

Events kunnen worden gecombineerd met WebSocket updates voor real-time tracking:

```typescript
// 1. Haal event op via REST
const event = await fetch('/api/events/active').then(r => r.json());

// 2. Connect WebSocket voor real-time updates
const ws = new WebSocket('ws://localhost:8082/ws/steps');
// Production: wss://dklemailservice.onrender.com/ws/steps

ws.onmessage = (msg) => {
    const update = JSON.parse(msg.data);
    
    if (update.type === 'step_update') {
        // Participant heeft stappen geüpdatet
        console.log(`${update.naam}: +${update.delta} stappen`);
    }
};

// 3. Subscribe to updates
ws.send(JSON.stringify({
    type: 'subscribe',
    channels: ['step_updates', 'total_updates']
}));
```

Zie [`docs/WEBSOCKET_INTEGRATION_GUIDE.md`](../WEBSOCKET_INTEGRATION_GUIDE.md) voor volledige WebSocket documentatie.

---

## Future Enhancements

### Geplande Features

- **Live Tracking**: Real-time locatie tracking tijdens events
- **Geofence Notifications**: Push notifications bij geofence entry/exit
- **Route Visualization**: Kaart weergave van participant routes
- **Event Analytics**: Gedetailleerde statistieken per event
- **Weather Integration**: Weer voorspelling voor event datum
- **Photo Checkpoints**: Foto verificatie bij checkpoints
- **Team Events**: Team-gebaseerde events en competitie

### API Extensions

- `POST /api/events/:id/register` - Participant registratie
- `POST /api/events/:id/checkin` - Check-in bij event start
- `POST /api/events/:id/location` - Locatie update tijdens event
- `GET /api/events/:id/live` - Live tracking data
- `GET /api/events/:id/analytics` - Event statistieken

---

## Support & Resources

### Documentatie
- [Steps API](./steps-api.md) - Steps tracking API
- [WebSocket Guide](../WEBSOCKET_INTEGRATION_GUIDE.md) - Real-time updates
- [RBAC Reference](../AUTH_AND_RBAC.md) - Permissions systeem
- [Database Schema](../DATABASE_REFERENCE.md) - Database documentatie

### Contact
Voor vragen of problemen:
1. Controleer deze documentatie
2. Bekijk server logs voor foutmeldingen
3. Test met cURL voor API debugging
4. Neem contact op met development team

---

**Version:** 1.0.0  
**Last Updated:** 2025-11-03  
**API Base:** `/api`  
**Migration:** V1_53