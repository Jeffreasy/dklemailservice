# Steps API Reference

Complete API referentie voor de stappen tracking functionaliteit in de DKL Email Service.

## Overzicht

De Steps API stelt deelnemers in staat om hun dagelijkse stappen bij te houden en hun voortgang te monitoren. Het systeem berekent automatisch toegewezen fondsen gebaseerd op de afstand die ze lopen.

## Endpoints

### POST /api/steps/:id

Werkt het aantal stappen bij voor een specifieke deelnemer (delta update).

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "steps": 1500
}
```

**Validatie:**
- `steps`: Integer, kan positief of negatief zijn (delta waarde)
- `id`: Geldige UUID van een deelnemer

**Response (200 OK):**
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "naam": "John Doe",
    "email": "john@example.com",
    "afstand": "10 KM",
    "steps": 2500,
    "created_at": "2024-03-20T15:04:05Z",
    "updated_at": "2024-03-20T16:00:00Z"
}
```

**Response (400 Bad Request):**
```json
{
    "error": "ID is verplicht"
}
```

**Response (403 Forbidden):**
```json
{
    "error": "Unauthorized"
}
```

**Response (404 Not Found):**
```json
{
    "error": "Deelnemer niet gevonden"
}
```

**Permissions:** `steps:write`

**Implementatie:** [`handlers/steps_handler.go:65`](../../handlers/steps_handler.go:65)

### GET /api/participant/:id/dashboard

Haalt dashboard informatie op voor een specifieke deelnemer.

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**
```json
{
    "steps": 2500,
    "route": "10 KM",
    "allocatedFunds": 75
}
```

**Response (404 Not Found):**
```json
{
    "error": "Deelnemer niet gevonden"
}
```

**Permissions:** `steps:read`

**Implementatie:** [`handlers/steps_handler.go:111`](../../handlers/steps_handler.go:111)

### GET /api/total-steps

Haalt het totaal aantal stappen op voor een specifiek jaar.

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Query Parameters:**
- `year` (optioneel): Jaar om op te filteren (default: 2025)

**Response (200 OK):**
```json
{
    "total_steps": 150000
}
```

**Permissions:** `steps:read`

**Implementatie:** [`handlers/steps_handler.go:149`](../../handlers/steps_handler.go:149)

### GET /api/funds-distribution

Haalt de fondsverdeling op over alle afstanden.

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**
```json
{
    "totalX": 350,
    "routes": {
        "6 KM": 50,
        "10 KM": 75,
        "15 KM": 100,
        "20 KM": 125
    }
}
```

**Permissions:** `steps:read`

**Implementatie:** [`handlers/steps_handler.go:185`](../../handlers/steps_handler.go:185)

### GET /api/steps/admin/route-funds

Haalt alle route fondsallocaties op voor beheer (admin only).

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**
```json
[
    {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "route": "6 KM",
        "amount": 50,
        "created_at": "2024-03-20T15:04:05Z",
        "updated_at": "2024-03-20T15:04:05Z"
    },
    {
        "id": "550e8400-e29b-41d4-a716-446655440001",
        "route": "10 KM",
        "amount": 75,
        "created_at": "2024-03-20T15:04:05Z",
        "updated_at": "2024-03-20T15:04:05Z"
    }
]
```

**Permissions:** `steps:write`

**Implementatie:** [`handlers/steps_handler.go:200`](../../handlers/steps_handler.go:200)

### POST /api/steps/admin/route-funds

Maakt een nieuwe route fondsallocatie aan (admin only).

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "route": "5 KM",
    "amount": 40
}
```

**Response (201 Created):**
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440002",
    "route": "5 KM",
    "amount": 40,
    "created_at": "2024-03-20T15:04:05Z",
    "updated_at": "2024-03-20T15:04:05Z"
}
```

**Permissions:** `steps:write`

**Implementatie:** [`handlers/steps_handler.go:220`](../../handlers/steps_handler.go:220)

### PUT /api/steps/admin/route-funds/{route}

Werkt een route fondsallocatie bij (admin only).

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "route": "6 KM",
    "amount": 60
}
```

**Response (200 OK):**
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "route": "6 KM",
    "amount": 60,
    "created_at": "2024-03-20T15:04:05Z",
    "updated_at": "2024-03-20T16:00:00Z"
}
```

**Permissions:** `steps:write`

**Implementatie:** [`handlers/steps_handler.go:270`](../../handlers/steps_handler.go:270)

### DELETE /api/steps/admin/route-funds/{route}

Verwijdert een route fondsallocatie (admin only).

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Route fund succesvol verwijderd"
}
```

**Permissions:** `steps:write`

**Implementatie:** [`handlers/steps_handler.go:320`](../../handlers/steps_handler.go:320)

## Business Logic

### Fondsverdeling per Afstand

Het systeem kent automatisch fondsen toe gebaseerd op de gekozen afstand. Deze bedragen zijn **configureerbaar** via de admin API.

**Standaard bedragen (kan aangepast worden):**

| Afstand | Toegewezen Fonds |
|---------|------------------|
| 6 KM    | €50             |
| 10 KM   | €75             |
| 15 KM   | €100            |
| 20 KM   | €125            |

**Admin beheer van fondsbedragen:**
- `GET /api/steps/admin/route-funds` - Alle fondsbedragen ophalen
- `POST /api/steps/admin/route-funds` - Nieuw fondsbedrag toevoegen
- `PUT /api/steps/admin/route-funds/{route}` - Fondsbedrag bijwerken
- `DELETE /api/steps/admin/route-funds/{route}` - Fondsbedrag verwijderen

### Stappen Updates

- Stappen worden altijd als delta toegevoegd (niet overschreven)
- Negatieve waarden worden geaccepteerd maar kunnen niet leiden tot negatieve totaal stappen
- Minimum totaal stappen = 0

### Fondsverdeling Berekening

De fondsverdeling kan op twee manieren worden berekend:

1. **Gelijke verdeling:** Totaal bedrag gelijk verdeeld over alle routes
2. **Proportionele verdeling:** Gebaseerd op aantal deelnemers per route

## RBAC Permissions

### Vereiste Permissions

| Endpoint | Permission | Rol |
|----------|------------|-----|
| POST /api/steps/:id | `steps:write` | Admin, Staff |
| GET /api/participant/:id/dashboard | `steps:read` | Admin, Staff, Deelnemer |
| GET /api/total-steps | `steps:read` | Admin, Staff |
| GET /api/funds-distribution | `steps:read` | Admin, Staff |

### Permission Setup

Permissions worden automatisch aangemaakt via migratie `V1_45__add_steps_permissions.sql`.

## Database Schema

### Aanmeldingen Tabel

```sql
ALTER TABLE aanmeldingen ADD COLUMN steps INTEGER DEFAULT 0;
```

**Nieuwe Kolom:**
- `steps`: INTEGER, DEFAULT 0, NOT NULL

## Error Handling

### HTTP Status Codes

| Status | Betekenis |
|--------|-----------|
| 200 | Succes |
| 400 | Ongeldige request (verkeerd ID formaat, ontbrekende parameters) |
| 401 | Niet geauthenticeerd |
| 403 | Geen toestemming |
| 404 | Deelnemer niet gevonden |
| 500 | Server fout |

### Error Responses

```json
{
    "error": "Beschrijvende foutmelding"
}
```

## Rate Limiting

Alle endpoints zijn onderhevig aan de algemene rate limiting regels van de applicatie.

## Testing

### Test Mode

Steps endpoints ondersteunen test mode via de `X-Test-Mode` header.

### Voorbeeld Test Data

```javascript
// Update stappen voor deelnemer
const response = await fetch('/api/steps/550e8400-e29b-41d4-a716-446655440000', {
    method: 'POST',
    headers: {
        'Authorization': 'Bearer ' + token,
        'Content-Type': 'application/json'
    },
    body: JSON.stringify({
        steps: 1000
    })
});

const result = await response.json();
console.log('Nieuwe totaal stappen:', result.steps);
```

## Frontend Integration

### React Hook Voorbeeld

```jsx
import { useState, useEffect } from 'react';

function ParticipantDashboard({ participantId }) {
    const [dashboard, setDashboard] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        fetchDashboard();
    }, [participantId]);

    const fetchDashboard = async () => {
        try {
            const response = await fetch(`/api/participant/${participantId}/dashboard`, {
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                }
            });

            if (response.ok) {
                const data = await response.json();
                setDashboard(data);
            }
        } catch (error) {
            console.error('Fout bij ophalen dashboard:', error);
        } finally {
            setLoading(false);
        }
    };

    const updateSteps = async (deltaSteps) => {
        try {
            const response = await fetch(`/api/steps/${participantId}`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ steps: deltaSteps })
            });

            if (response.ok) {
                // Refresh dashboard data
                await fetchDashboard();
            }
        } catch (error) {
            console.error('Fout bij bijwerken stappen:', error);
        }
    };

    if (loading) return <div>Loading...</div>;

    return (
        <div className="dashboard">
            <h2>Dashboard</h2>
            <p>Stappen: {dashboard.steps}</p>
            <p>Afstand: {dashboard.route}</p>
            <p>Toegewezen Fonds: €{dashboard.allocatedFunds}</p>

            <button onClick={() => updateSteps(500)}>
                +500 Stappen
            </button>
            <button onClick={() => updateSteps(-100)}>
                -100 Stappen
            </button>
        </div>
    );
}
```

## Monitoring

### Metrics

Steps endpoints worden opgenomen in de algemene applicatie metrics:

- Request/response tijden
- Error rates
- Rate limiting events

### Logging

Alle stappen updates worden gelogd met:
- Deelnemer ID
- Delta waarde
- Nieuwe totaal
- Timestamp
- Gebruiker die de update uitvoerde

## Security Considerations

### Authentication
- Alle endpoints vereisen geldige JWT tokens
- Tokens worden gevalideerd op elke request

### Authorization
- Fine-grained permissions via RBAC systeem
- Deelnemers kunnen alleen hun eigen dashboard bekijken
- Admin/Staff hebben volledige toegang

### Input Validation
- UUID validatie voor participant IDs
- Integer validatie voor steps waarden
- SQL injection preventie via parameterized queries

## Performance

### Database Indexes
- Index op `steps` kolom voor aggregatie queries
- Index op `afstand` kolom voor fondsverdeling berekeningen

### Caching
- Fondsverdeling kan worden gecached voor betere performance
- Redis caching voor veelgebruikte dashboard data

### Optimizations
- Batch updates voor bulk operaties
- Lazy loading voor grote datasets
- Pagination voor lijst endpoints

## Troubleshooting

### Common Issues

**401 Unauthorized**
- Controleer JWT token geldigheid
- Controleer token expiration
- Controleer Authorization header formaat

**403 Forbidden**
- Controleer user permissions
- Controleer role assignments
- Controleer RBAC configuratie

**404 Not Found**
- Controleer participant ID formaat (moet UUID zijn)
- Controleer of deelnemer bestaat in database

**500 Internal Server Error**
- Controleer database connectie
- Controleer service dependencies
- Bekijk server logs voor details

### Debug Headers

Voor debugging kunnen de volgende headers worden gebruikt:

```http
X-Debug-Mode: true
X-Request-ID: uuid-v4
```

## Migration Guide

### Database Migration

1. Run migration `V1_44__add_steps_to_aanmeldingen.sql`
2. Run migration `V1_45__add_steps_permissions.sql`
3. Restart application

### API Changes

Nieuwe endpoints zijn backwards compatible en breken bestaande functionaliteit niet.

### Frontend Updates

1. Add steps tracking UI components
2. Implement dashboard views
3. Add step update functionality
4. Update permission checks

## Database Schema

### Aanmeldingen Tabel

```sql
ALTER TABLE aanmeldingen ADD COLUMN steps INTEGER DEFAULT 0;
```

**Nieuwe Kolom:**
- `steps`: INTEGER, DEFAULT 0, NOT NULL

### Route Funds Tabel

```sql
CREATE TABLE route_funds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    route VARCHAR(50) NOT NULL UNIQUE,
    amount INTEGER NOT NULL CHECK (amount >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

**Kolommen:**
- `id`: UUID, Primary Key
- `route`: VARCHAR(50), Unique, Not Null (bijv. "6 KM", "10 KM")
- `amount`: INTEGER, Not Null, >= 0 (bedrag in euro's)
- `created_at`: TIMESTAMP, Auto-create
- `updated_at`: TIMESTAMP, Auto-update

## Future Enhancements

### Geplande Features

- **Step History**: Gedetailleerde geschiedenis van stappen updates
- **Goals**: Persoonlijke doelen instellen
- **Leaderboards**: Competitie tussen deelnemers
- **Gamification**: Badges en achievements
- **Mobile App**: Native mobile ondersteuning
- **Wearable Integration**: Directe integratie met fitness trackers
- **Bulk Import**: Excel/CSV import van stappen data
- **Step Validation**: Automatische validatie van ingevoerde stappen
- **Reporting**: Uitgebreide rapportages voor sponsors

### API Extensions

- `GET /api/steps/:id/history` - Stappen geschiedenis
- `POST /api/steps/:id/goal` - Doel instellen
- `GET /api/leaderboard` - Leaderboard data
- `GET /api/steps/:id/stats` - Gedetailleerde statistieken

## Support

Voor vragen of problemen met de Steps API:

1. Controleer deze documentatie
2. Bekijk server logs voor foutmeldingen
3. Controleer database connectie en migrations
4. Neem contact op met development team

---

**Version:** 1.0.0
**Last Updated:** 2024-03-20
**API Base:** `/api`