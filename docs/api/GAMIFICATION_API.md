# Gamification API Reference

Complete API referentie voor de gamification functionaliteit (badges, achievements en leaderboard) in de DKL Email Service.

## Overzicht

De Gamification API stelt deelnemers in staat om badges te verdienen, hun rank te bekijken, en te competeren op de leaderboard. Admins kunnen badges beheren en toekennen aan deelnemers.

---

## Public Endpoints

### GET /api/leaderboard

Haalt de leaderboard op met rankings van alle deelnemers.

**Query Parameters:**
- `page` (optioneel): Paginanummer (default: 1)
- `limit` (optioneel): Aantal entries per pagina (default: 50, max: 100)
- `route` (optioneel): Filter op specifieke route (bijv. "10 KM")
- `year` (optioneel): Filter op jaar (bijv. 2026)
- `min_steps` (optioneel): Minimum aantal stappen
- `top_n` (optioneel): Toon alleen top N deelnemers

**Response (200 OK):**
```json
{
  "entries": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "naam": "John Doe",
      "route": "10 KM",
      "steps": 15000,
      "achievement_points": 250,
      "total_score": 15250,
      "rank": 1,
      "badge_count": 5,
      "joined_at": "2026-01-15T10:00:00Z"
    }
  ],
  "total_entries": 120,
  "current_page": 1,
  "total_pages": 3,
  "limit": 50
}
```

**Implementatie:** [`handlers/gamification_handler.go:56`](../../handlers/gamification_handler.go:56)

---

### GET /api/badges

Haalt alle actieve badges op die deelnemers kunnen verdienen.

**Headers:**
```http
Authorization: Bearer <jwt-token> (optioneel)
```

**Response (200 OK):**
```json
[
  {
    "id": "badge-uuid-1",
    "name": "First Steps",
    "description": "Je eerste 1000 stappen gezet",
    "icon_url": "/icons/badges/first-steps.svg",
    "criteria": {
      "min_steps": 1000
    },
    "points": 10,
    "is_active": true,
    "display_order": 1,
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-01-01T00:00:00Z",
    "earned_count": 45,
    "last_earned_at": "2026-01-20T15:30:00Z",
    "earned_by_current_user": true
  }
]
```

**Implementatie:** [`handlers/gamification_handler.go:105`](../../handlers/gamification_handler.go:105)

---

## Participant Endpoints

### GET /api/participant/achievements

Haalt alle achievements op van de ingelogde deelnemer.

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**
```json
{
  "participant_id": "550e8400-e29b-41d4-a716-446655440000",
  "participant_name": "John Doe",
  "total_badges": 5,
  "total_points": 285,
  "achievements": [
    {
      "id": "achievement-uuid-1",
      "participant_id": "550e8400-e29b-41d4-a716-446655440000",
      "earned_at": "2026-01-15T14:30:00Z",
      "badge": {
        "id": "badge-uuid-1",
        "name": "First Steps",
        "description": "Je eerste 1000 stappen gezet",
        "icon_url": "/icons/badges/first-steps.svg",
        "points": 10
      }
    }
  ],
  "available_badges": [
    {
      "id": "badge-uuid-5",
      "name": "Marathon Walker",
      "description": "42195 stappen bereikt",
      "icon_url": "/icons/badges/marathon.svg",
      "points": 500,
      "earned_by_current_user": false
    }
  ]
}
```

**Permissions:** `achievements:read`

**Implementatie:** [`handlers/gamification_handler.go:130`](../../handlers/gamification_handler.go:130)

---

### GET /api/participant/rank

Haalt rank informatie op van de ingelogde deelnemer.

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**
```json
{
  "participant_id": "550e8400-e29b-41d4-a716-446655440000",
  "naam": "John Doe",
  "rank": 5,
  "total_score": 12450,
  "steps": 12000,
  "achievement_points": 450,
  "badge_count": 8,
  "above_me": {
    "id": "other-uuid",
    "naam": "Jane Smith",
    "rank": 4,
    "total_score": 12600
  },
  "below_me": {
    "id": "another-uuid",
    "naam": "Bob Johnson",
    "rank": 6,
    "total_score": 12300
  }
}
```

**Permissions:** `leaderboard:read`

**Implementatie:** [`handlers/gamification_handler.go:149`](../../handlers/gamification_handler.go:149)

---

### GET /api/participant/:id/achievements

Haalt achievements op voor een specifieke deelnemer (admin/staff toegang).

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**URL Parameters:**
- `id`: UUID van de deelnemer

**Response:** Zelfde als `/api/participant/achievements`

**Permissions:** `achievements:read`

**Implementatie:** [`handlers/gamification_handler.go:167`](../../handlers/gamification_handler.go:167)

---

## Admin Badge Management

### POST /api/admin/badges

Maakt een nieuwe badge aan (admin only).

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Super Walker",
  "description": "100000 stappen bereikt",
  "icon_url": "/icons/badges/super-walker.svg",
  "criteria": {
    "min_steps": 100000
  },
  "points": 1000,
  "is_active": true,
  "display_order": 10
}
```

**Response (201 Created):**
```json
{
  "id": "new-badge-uuid",
  "name": "Super Walker",
  "description": "100000 stappen bereikt",
  "icon_url": "/icons/badges/super-walker.svg",
  "criteria": {
    "min_steps": 100000
  },
  "points": 1000,
  "is_active": true,
  "display_order": 10,
  "created_at": "2026-01-20T10:00:00Z",
  "updated_at": "2026-01-20T10:00:00Z"
}
```

**Permissions:** `badges:write`

**Implementatie:** [`handlers/gamification_handler.go:190`](../../handlers/gamification_handler.go:190)

---

### GET /api/admin/badges

Haalt alle badges op inclusief inactieve (admin only).

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Query Parameters:**
- `active_only` (optioneel): Toon alleen actieve badges (default: false)

**Response (200 OK):**
```json
[
  {
    "id": "badge-uuid-1",
    "name": "First Steps",
    "description": "Je eerste 1000 stappen gezet",
    "icon_url": "/icons/badges/first-steps.svg",
    "criteria": {
      "min_steps": 1000
    },
    "points": 10,
    "is_active": true,
    "display_order": 1,
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-01-01T00:00:00Z"
  }
]
```

**Permissions:** `badges:write`

**Implementatie:** [`handlers/gamification_handler.go:230`](../../handlers/gamification_handler.go:230)

---

### GET /api/admin/badges/:id

Haalt een specifieke badge op (admin only).

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**URL Parameters:**
- `id`: UUID van de badge

**Response (200 OK):** Zie POST response

**Permissions:** `badges:write`

**Implementatie:** [`handlers/gamification_handler.go:247`](../../handlers/gamification_handler.go:247)

---

### PUT /api/admin/badges/:id

Werkt een badge bij (admin only).

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**URL Parameters:**
- `id`: UUID van de badge

**Request Body:** Zie POST request body

**Response (200 OK):** Zie POST response

**Permissions:** `badges:write`

**Implementatie:** [`handlers/gamification_handler.go:267`](../../handlers/gamification_handler.go:267)

---

### DELETE /api/admin/badges/:id

Verwijdert een badge (admin only).

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**URL Parameters:**
- `id`: UUID van de badge

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Badge succesvol verwijderd"
}
```

**Permissions:** `badges:write`

**Implementatie:** [`handlers/gamification_handler.go:313`](../../handlers/gamification_handler.go:313)

---

## Admin Achievement Management

### POST /api/admin/achievements/award

Kent een badge toe aan een deelnemer (admin/staff).

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "participant_id": "550e8400-e29b-41d4-a716-446655440000",
  "badge_id": "badge-uuid-1"
}
```

**Response (201 Created):**
```json
{
  "id": "achievement-uuid",
  "participant_id": "550e8400-e29b-41d4-a716-446655440000",
  "earned_at": "2026-01-20T15:00:00Z",
  "badge": {
    "id": "badge-uuid-1",
    "name": "First Steps",
    "description": "Je eerste 1000 stappen gezet",
    "icon_url": "/icons/badges/first-steps.svg",
    "points": 10
  }
}
```

**Permissions:** `achievements:write`

**Implementatie:** [`handlers/gamification_handler.go:339`](../../handlers/gamification_handler.go:339)

---

### DELETE /api/admin/achievements/remove

Verwijdert een badge van een deelnemer (admin/staff).

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "participant_id": "550e8400-e29b-41d4-a716-446655440000",
  "badge_id": "badge-uuid-1"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Achievement succesvol verwijderd"
}
```

**Permissions:** `achievements:write`

**Implementatie:** [`handlers/gamification_handler.go:365`](../../handlers/gamification_handler.go:365)

---

## Badge Criteria

Badges kunnen verschillende criteria hebben die bepalen wanneer ze worden verdiend:

### Beschikbare Criteria

```typescript
interface BadgeCriteria {
  min_steps?: number;           // Minimum aantal stappen
  min_days?: number;            // Minimum aantal dagen actief
  consecutive_days?: number;    // Opeenvolgende dagen met activiteit
  route?: string;               // Specifieke route (bijv. "20 KM")
  early_participant?: boolean;  // Binnen eerste 50 deelnemers
  has_team?: boolean;          // Lid van een team
}
```

### Voorbeelden

**Steps Badge:**
```json
{
  "name": "10K Master",
  "criteria": {
    "min_steps": 10000
  },
  "points": 100
}
```

**Route Badge:**
```json
{
  "name": "Distance Hero",
  "criteria": {
    "route": "20 KM"
  },
  "points": 150
}
```

**Early Bird Badge:**
```json
{
  "name": "Early Bird",
  "criteria": {
    "early_participant": true
  },
  "points": 25
}
```

---

## Scoring System

De leaderboard ranking wordt bepaald door:

```
Total Score = Steps + Achievement Points
```

- **Steps**: Aantal gelopen stappen
- **Achievement Points**: Som van punten van alle verdiende badges
- **Rank**: Positie op leaderboard (1 = hoogste score)

---

## RBAC Permissions

| Endpoint | Permission | Rol |
|----------|------------|-----|
| GET /api/leaderboard | Geen | Public |
| GET /api/badges | Geen | Public |
| GET /api/participant/achievements | `achievements:read` | Deelnemer, Staff, Admin |
| GET /api/participant/rank | `leaderboard:read` | Deelnemer, Staff, Admin |
| GET /api/participant/:id/achievements | `achievements:read` | Staff, Admin |
| POST /api/admin/badges | `badges:write` | Admin |
| PUT /api/admin/badges/:id | `badges:write` | Admin |
| DELETE /api/admin/badges/:id | `badges:write` | Admin |
| POST /api/admin/achievements/award | `achievements:write` | Staff, Admin |
| DELETE /api/admin/achievements/remove | `achievements:write` | Staff, Admin |

---

## Database Schema

### Badges Table
```sql
CREATE TABLE badges (
    id UUID PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT NOT NULL,
    icon_url VARCHAR(500),
    criteria JSONB NOT NULL,
    points INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
);
```

### Participant Achievements Table
```sql
CREATE TABLE participant_achievements (
    id UUID PRIMARY KEY,
    participant_id UUID NOT NULL,
    badge_id UUID NOT NULL,
    earned_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(participant_id, badge_id)
);
```

---

## Frontend Integration

### React Hook Example

```tsx
import { useState, useEffect } from 'react';

function LeaderboardComponent() {
    const [leaderboard, setLeaderboard] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        fetchLeaderboard();
    }, []);

    const fetchLeaderboard = async () => {
        try {
            const response = await fetch('/api/leaderboard?limit=10');
            if (response.ok) {
                const data = await response.json();
                setLeaderboard(data);
            }
        } catch (error) {
            console.error('Error fetching leaderboard:', error);
        } finally {
            setLoading(false);
        }
    };

    if (loading) return <div>Loading...</div>;

    return (
        <div className="leaderboard">
            <h2>Top 10 Deelnemers</h2>
            <table>
                <thead>
                    <tr>
                        <th>Rank</th>
                        <th>Naam</th>
                        <th>Steps</th>
                        <th>Badges</th>
                        <th>Total Score</th>
                    </tr>
                </thead>
                <tbody>
                    {leaderboard?.entries.map(entry => (
                        <tr key={entry.id}>
                            <td>#{entry.rank}</td>
                            <td>{entry.naam}</td>
                            <td>{entry.steps.toLocaleString()}</td>
                            <td>{entry.badge_count}</td>
                            <td>{entry.total_score.toLocaleString()}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
}
```

---

## Error Handling

### HTTP Status Codes

| Status | Betekenis |
|--------|-----------|
| 200 | Succes |
| 201 | Resource aangemaakt |
| 400 | Ongeldige request |
| 401 | Niet geauthenticeerd |
| 403 | Geen toestemming |
| 404 | Resource niet gevonden |
| 409 | Conflict (bijv. duplicate badge name) |
| 500 | Server fout |

---

## Troubleshooting

### Common Issues

**400 Bad Request bij badge create**
- Controleer of alle verplichte velden zijn ingevuld
- Controleer of de naam uniek is
- Controleer of criteria valid JSON is

**403 Forbidden**
- Controleer user permissions (badges:write voor admin endpoints)
- Controleer role assignments

**409 Conflict bij award badge**
- Deelnemer heeft badge al
- Gebruik GET /api/participant/:id/achievements om te checken

---

**Version:** 1.0.0  
**Last Updated:** 2026-01-02  
**API Base:** `/api`