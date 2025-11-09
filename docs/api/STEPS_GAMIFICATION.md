# Steps & Gamification API

API endpoints voor het tracken van stappen, behalen van achievements en leaderboards.

## Base URL

```
Development: http://localhost:8080/api
Production: https://api.dklemailservice.com/api
```

## Authentication

Alle endpoints vereisen JWT authenticatie en specifieke permissies:

```
Authorization: Bearer <your_jwt_token>
```

---

## Steps Tracking

### Update Steps

Update het aantal stappen voor een deelnemer tijdens een evenement.

**Endpoint:** `POST /api/registration/:id/steps`

**Authentication:** Vereist - `steps:write` permissie

**Request Body:**

```json
{
  "delta_steps": 1000
}
```

**Parameters:**
- `id`: Event registration ID (UUID)
- `delta_steps`: Aantal stappen om toe te voegen (kan negatief zijn)

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "participant_id": "uuid",
    "naam": "John Doe",
    "email": "john@example.com",
    "steps": 5000,
    "total_distance": 3.5,
    "target_distance": 10.0,
    "percentage_complete": 35.0,
    "achievements_unlocked": ["first_step", "5k_milestone"]
  },
  "message": "Steps updated successfully"
}
```

**WebSocket Broadcast:**

Deze endpoint broadcast real-time updates naar alle verbonden clients via WebSocket:

```json
{
  "type": "steps_update",
  "data": {
    "participant_id": "uuid",
    "steps": 5000,
    "delta": 1000,
    "total_distance": 3.5
  },
  "timestamp": "2025-01-08T14:30:00Z"
}
```

---

### Get Participant Dashboard

Haal uitgebreide dashboard informatie op voor een deelnemer.

**Endpoint:** `GET /api/registration/:id/dashboard`

**Authentication:** Vereist - `steps:read` permissie

**Response:** `200 OK`

```json
{
  "participant": {
    "id": "uuid",
    "naam": "John Doe",
    "email": "john@example.com",
    "telefoon": "+31612345678",
    "rol": "deelnemer",
    "afstand": "10km",
    "steps": 15000,
    "created_at": "2025-01-01T10:00:00Z"
  },
  "total_steps": 15000,
  "rank": 3,
  "achievements": [
    {
      "id": "uuid",
      "name": "First Step",
      "description": "Eerste stap gezet",
      "icon_url": "https://res.cloudinary.com/...",
      "unlocked_at": "2025-01-02T10:00:00Z"
    }
  ],
  "badges": [
    {
      "id": "uuid",
      "name": "5K Hero",
      "description": "5000 stappen bereikt",
      "image_url": "https://res.cloudinary.com/...",
      "awarded_at": "2025-01-05T15:30:00Z"
    }
  ]
}
```

---

### Get Total Steps

Haal het totaal aantal stappen op voor een specifiek jaar.

**Endpoint:** `GET /api/total-steps`

**Authentication:** Vereist - `steps:read` permissie

**Query Parameters:**
- `year` (optional): Jaar (default: huidig jaar)

**Response:** `200 OK`

```json
{
  "year": 2025,
  "total_steps": 1500000,
  "total_participants": 250,
  "average_steps_per_participant": 6000
}
```

---

### Get Funds Distribution

Haal de fondsen verdeling op gebaseerd op afstanden.

**Endpoint:** `GET /api/funds-distribution`

**Authentication:** Vereist - `steps:read` permissie

**Response:** `200 OK`

```json
{
  "total_funds": 50000,
  "distributions": [
    {
      "route": "5km",
      "steps": 250000,
      "percentage": 25.0,
      "amount": 12500,
      "participants": 100
    },
    {
      "route": "10km",
      "steps": 500000,
      "percentage": 50.0,
      "amount": 25000,
      "participants": 100
    },
    {
      "route": "15km",
      "steps": 250000,
      "percentage": 25.0,
      "amount": 12500,
      "participants": 50
    }
  ]
}
```

---

## Route Funds

Beheer fondsen allocatie per route.

### List Route Funds

**Endpoint:** `GET /api/route-funds`

**Authentication:** Vereist - `route_fund:read` permissie

**Response:** `200 OK`

```json
[
  {
    "id": "uuid",
    "route_name": "5km",
    "allocated_amount": 12500,
    "description": "Fonds voor 5km route",
    "year": 2025,
    "is_active": true
  }
]
```

---

### Create Route Fund

**Endpoint:** `POST /api/route-funds`

**Authentication:** Vereist - `route_fund:write` permissie

**Request Body:**

```json
{
  "route_name": "5km",
  "allocated_amount": 12500,
  "description": "Fonds voor 5km route",
  "year": 2025,
  "is_active": true
}
```

---

### Update Route Fund

**Endpoint:** `PUT /api/route-funds/:id`

**Authentication:** Vereist - `route_fund:write` permissie

---

### Delete Route Fund

**Endpoint:** `DELETE /api/route-funds/:id`

**Authentication:** Vereist - `route_fund:delete` permissie

---

## Achievements

Beheer achievements die deelnemers kunnen ontgrendelen.

### List All Achievements

**Endpoint:** `GET /api/achievements`

**Authentication:** Vereist - `achievement:read` permissie

**Response:** `200 OK`

```json
[
  {
    "id": "uuid",
    "name": "First Step",
    "description": "Eerste stap gezet",
    "icon_url": "https://res.cloudinary.com/...",
    "required_steps": 1,
    "category": "milestone",
    "points": 10,
    "is_active": true,
    "created_at": "2025-01-08T10:00:00Z"
  }
]
```

---

### Get User Achievements

**Endpoint:** `GET /api/achievements/user/:user_id`

**Authentication:** Vereist - `achievement:read` permissie

**Response:** `200 OK`

```json
{
  "user_id": "uuid",
  "achievements": [
    {
      "id": "uuid",
      "name": "First Step",
      "description": "Eerste stap gezet",
      "unlocked_at": "2025-01-02T10:00:00Z",
      "points": 10
    }
  ],
  "total_points": 150,
  "completion_percentage": 45.5
}
```

---

### Create Achievement

**Endpoint:** `POST /api/achievements`

**Authentication:** Vereist - `achievement:write` permissie

**Request Body:**

```json
{
  "name": "Marathon Master",
  "description": "Voltooi een marathon afstand",
  "icon_url": "https://res.cloudinary.com/...",
  "required_steps": 42195,
  "category": "distance",
  "points": 100,
  "is_active": true
}
```

**Response:** `201 Created`

---

### Update Achievement

**Endpoint:** `PUT /api/achievements/:id`

**Authentication:** Vereist - `achievement:write` permissie

---

### Delete Achievement

**Endpoint:** `DELETE /api/achievements/:id`

**Authentication:** Vereist - `achievement:delete` permissie

---

## Badges

Beheer badges die aan deelnemers kunnen worden toegekend.

### List All Badges

**Endpoint:** `GET /api/badges`

**Authentication:** Vereist - `badge:read` permissie

**Response:** `200 OK`

```json
[
  {
    "id": "uuid",
    "name": "5K Hero",
    "description": "5000 stappen bereikt",
    "image_url": "https://res.cloudinary.com/...",
    "tier": "bronze",
    "required_achievement_id": "uuid",
    "is_active": true,
    "created_at": "2025-01-08T10:00:00Z"
  }
]
```

---

### Get User Badges

**Endpoint:** `GET /api/badges/user/:user_id`

**Authentication:** Vereist - `badge:read` permissie

**Response:** `200 OK`

```json
{
  "user_id": "uuid",
  "badges": [
    {
      "id": "uuid",
      "name": "5K Hero",
      "description": "5000 stappen bereikt",
      "image_url": "https://res.cloudinary.com/...",
      "tier": "bronze",
      "awarded_at": "2025-01-05T15:30:00Z"
    }
  ]
}
```

---

### Create Badge

**Endpoint:** `POST /api/badges`

**Authentication:** Vereist - `badge:write` permissie

**Request Body:**

```json
{
  "name": "Marathon Legend",
  "description": "Marathon afstand voltooid",
  "image_url": "https://res.cloudinary.com/...",
  "tier": "gold",
  "required_achievement_id": "uuid",
  "is_active": true
}
```

**Badge Tiers:**
- `bronze` - Basis niveau
- `silver` - Gemiddeld niveau
- `gold` - Hoog niveau
- `platinum` - Top niveau

**Response:** `201 Created`

---

### Update Badge

**Endpoint:** `PUT /api/badges/:id`

**Authentication:** Vereist - `badge:write` permissie

---

### Delete Badge

**Endpoint:** `DELETE /api/badges/:id`

**Authentication:** Vereist - `badge:delete` permissie

---

### Award Badge

Handmatig een badge toekennen aan een gebruiker.

**Endpoint:** `POST /api/badges/:id/award`

**Authentication:** Vereist - `badge:write` permissie

**Request Body:**

```json
{
  "user_id": "uuid",
  "reason": "Exceptional performance"
}
```

**Response:** `200 OK`

```json
{
  "success": true,
  "message": "Badge awarded successfully",
  "badge": {
    "id": "uuid",
    "name": "Special Achievement",
    "awarded_at": "2025-01-08T14:30:00Z"
  }
}
```

---

### Remove Badge

Verwijder een badge van een gebruiker.

**Endpoint:** `DELETE /api/badges/:id/remove`

**Authentication:** Vereist - `badge:delete` permissie

**Request Body:**

```json
{
  "user_id": "uuid"
}
```

---

## Leaderboard

Ranglijsten voor stappen competitie.

### Get Leaderboard

**Endpoint:** `GET /api/leaderboard`

**Authentication:** Vereist - `leaderboard:read` permissie

**Query Parameters:**
- `period` (optional): `daily`, `weekly`, `monthly`, `all` (default: `all`)
- `limit` (optional): Aantal resultaten (default: 10, max: 100)
- `route` (optional): Filter op afstand (5km, 10km, 15km)

**Response:** `200 OK`

```json
{
  "period": "weekly",
  "generated_at": "2025-01-08T14:30:00Z",
  "entries": [
    {
      "rank": 1,
      "participant_id": "uuid",
      "naam": "Jane Smith",
      "steps": 25000,
      "distance": 17.5,
      "route": "10km",
      "badges_count": 5,
      "total_points": 250
    },
    {
      "rank": 2,
      "participant_id": "uuid",
      "naam": "John Doe",
      "steps": 23500,
      "distance": 16.45,
      "route": "10km",
      "badges_count": 4,
      "total_points": 230
    }
  ],
  "total_participants": 250,
  "current_user_rank": 15
}
```

---

### Get User Stats

Statistieken voor een specifieke gebruiker.

**Endpoint:** `GET /api/leaderboard/user/:user_id`

**Authentication:** Vereist - `leaderboard:read` permissie

**Response:** `200 OK`

```json
{
  "user_id": "uuid",
  "naam": "John Doe",
  "current_rank": 15,
  "total_steps": 15000,
  "total_distance": 10.5,
  "achievements_count": 8,
  "badges_count": 3,
  "total_points": 180,
  "daily_average": 750,
  "weekly_average": 5250,
  "best_day": {
    "date": "2025-01-05",
    "steps": 3000
  },
  "streak_days": 7
}
```

---

## WebSocket Real-time Updates

Steps updates worden real-time verzonden via WebSocket.

### Connect to Steps WebSocket

**Endpoint:** `GET /ws/steps?token=<jwt_token>`

**Protocol:** WebSocket

**Authentication:** JWT token via query parameter

**Connection Example:**

```javascript
const token = localStorage.getItem('access_token');
const ws = new WebSocket(`ws://localhost:8080/ws/steps?token=${token}`);

ws.onopen = () => {
  console.log('Connected to Steps WebSocket');
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  handleStepsUpdate(message);
};
```

---

### WebSocket Message Types

#### 1. Steps Update

Verzonden wanneer een deelnemer stappen update.

```json
{
  "type": "steps_update",
  "data": {
    "participant_id": "uuid",
    "naam": "John Doe",
    "steps": 5000,
    "delta": 1000,
    "total_distance": 3.5,
    "percentage": 35.0
  },
  "timestamp": "2025-01-08T14:30:00Z"
}
```

#### 2. Achievement Unlocked

Verzonden wanneer een achievement wordt ontgrendeld.

```json
{
  "type": "achievement_unlocked",
  "data": {
    "participant_id": "uuid",
    "achievement": {
      "id": "uuid",
      "name": "5K Milestone",
      "description": "5000 stappen bereikt",
      "points": 50
    }
  },
  "timestamp": "2025-01-08T14:30:00Z"
}
```

#### 3. Badge Awarded

Verzonden wanneer een badge wordt toegekend.

```json
{
  "type": "badge_awarded",
  "data": {
    "participant_id": "uuid",
    "badge": {
      "id": "uuid",
      "name": "Bronze Runner",
      "tier": "bronze",
      "image_url": "https://res.cloudinary.com/..."
    }
  },
  "timestamp": "2025-01-08T14:30:00Z"
}
```

#### 4. Leaderboard Update

Verzonden wanneer de leaderboard wijzigt.

```json
{
  "type": "leaderboard_update",
  "data": {
    "period": "daily",
    "top_3": [
      {
        "rank": 1,
        "participant_id": "uuid",
        "naam": "Jane Smith",
        "steps": 25000
      }
    ],
    "your_rank": 15
  },
  "timestamp": "2025-01-08T14:30:00Z"
}
```

---

## Calculation Logic

### Steps to Distance Conversion

De service berekent afstand op basis van stappen:

```
Gemiddelde stap lengte: 0.7 meter
Afstand (km) = (Stappen × 0.7) / 1000
```

**Voorbeeld:**
- 10,000 stappen = 7 km
- 15,000 stappen = 10.5 km

### Progress Calculation

```
Percentage = (Huidige Afstand / Doel Afstand) × 100
```

### Points System

- Achievements: 10-100 punten afhankelijk van moeilijkheid
- Daily goal: 10 punten
- Weekly goal: 50 punten
- Monthly goal: 200 punten
- Top 10 leaderboard: 25-100 punten afhankelijk van positie

---

## Gamification Features

### Automatic Achievement Detection

Het systeem controleert automatisch op nieuwe achievements wanneer stappen worden bijgewerkt:

1. **Milestone Achievements**
   - First Step (1 stap)
   - 1K Steps (1,000 stappen)
   - 5K Steps (5,000 stappen)
   - 10K Steps (10,000 stappen)
   - Marathon (42,195 stappen)

2. **Consistency Achievements**
   - 3 Day Streak
   - 7 Day Streak
   - 30 Day Streak

3. **Distance Achievements**
   - 5K Complete
   - 10K Complete
   - Half Marathon
   - Full Marathon

### Badge Tiers

Badges zijn beschikbaar in verschillende tiers:

| Tier | Moeilijkheid | Kleur | Punten |
|------|--------------|-------|--------|
| Bronze | Basis | #CD7F32 | 10-25 |
| Silver | Gemiddeld | #C0C0C0 | 25-50 |
| Gold | Moeilijk | #FFD700 | 50-100 |
| Platinum | Expert | #E5E4E2 | 100+ |

---

## Performance Considerations

### Real-time Updates

- WebSocket broadcast beperkt tot 100 updates per seconde
- Steps updates worden gebatched voor grotere groepen
- Leaderboard updates elke 5 minuten voor all-time rankings

### Caching

- Leaderboard wordt gecached in Redis (5 minuten TTL)
- Achievement checks gebruiken in-memory caching
- User stats worden gecached per gebruiker (1 minuut TTL)

---

## Error Responses

### Common Errors

```json
{
  "error": "Participant not found",
  "code": "PARTICIPANT_NOT_FOUND"
}
```

```json
{
  "error": "Invalid steps delta. Must be between -10000 and 10000",
  "code": "INVALID_STEPS_DELTA"
}
```

```json
{
  "error": "Achievement already unlocked",
  "code": "ACHIEVEMENT_DUPLICATE"
}
```

---

## Integration Examples

### React Hook for Steps Tracking

```typescript
import { useState, useEffect } from 'react';
import { useWebSocket } from './useWebSocket';

export function useStepsTracking(participantId: string) {
  const [steps, setSteps] = useState(0);
  const [distance, setDistance] = useState(0);
  const [achievements, setAchievements] = useState([]);

  const { connected, lastMessage } = useWebSocket('steps', {
    onMessage: (message) => {
      if (message.type === 'steps_update' && 
          message.data.participant_id === participantId) {
        setSteps(message.data.steps);
        setDistance(message.data.total_distance);
      }
      
      if (message.type === 'achievement_unlocked' && 
          message.data.participant_id === participantId) {
        setAchievements(prev => [...prev, message.data.achievement]);
      }
    }
  });

  const updateSteps = async (delta: number) => {
    const response = await api.post(
      `/api/registration/${participantId}/steps`,
      { delta_steps: delta }
    );
    return response.data;
  };

  return {
    steps,
    distance,
    achievements,
    updateSteps,
    connected
  };
}
```

### Vue Component Example

```vue
<template>
  <div class="steps-tracker">
    <div class="stats">
      <h2>{{ steps.toLocaleString() }} stappen</h2>
      <p>{{ distance.toFixed(2) }} km</p>
      <div class="progress">
        <div class="bar" :style="{ width: percentage + '%' }"></div>
      </div>
    </div>
    
    <div class="achievements">
      <h3>Recent Achievements</h3>
      <div v-for="achievement in recentAchievements" :key="achievement.id">
        <img :src="achievement.icon_url" :alt="achievement.name" />
        <span>{{ achievement.name }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useStepsWebSocket } from '@/composables/useStepsWebSocket';

const props = defineProps(['participantId']);
const steps = ref(0);
const distance = ref(0);
const targetDistance = ref(10);
const recentAchievements = ref([]);

const percentage = computed(() => {
  return (distance.value / targetDistance.value) * 100;
});

const { connected } = useStepsWebSocket(props.participantId, (message) => {
  if (message.type === 'steps_update') {
    steps.value = message.data.steps;
    distance.value = message.data.total_distance;
  }
  if (message.type === 'achievement_unlocked') {
    recentAchievements.value.unshift(message.data.achievement);
  }
});

onMounted(async () => {
  const response = await fetch(`/api/registration/${props.participantId}/dashboard`);
  const data = await response.json();
  steps.value = data.total_steps;
});
</script>
```

---

## Testing

### Manual Testing with cURL

```bash
# Update steps
curl -X POST http://localhost:8080/api/registration/{id}/steps \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"delta_steps": 1000}'

# Get dashboard
curl http://localhost:8080/api/registration/{id}/dashboard \
  -H "Authorization: Bearer YOUR_TOKEN"

# Get leaderboard
curl "http://localhost:8080/api/leaderboard?period=weekly&limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Get total steps
curl "http://localhost:8080/api/total-steps?year=2025" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

Voor meer informatie:
- [WebSocket Documentation](./WEBSOCKET.md)
- [Authentication](./AUTHENTICATION.md)
- [Frontend Integration](../guides/FRONTEND_INTEGRATION.md)