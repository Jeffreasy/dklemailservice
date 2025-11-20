# Events API

API endpoints voor het beheren van evenementen en registraties.

## Base URL

```
Development: http://localhost:8080/api
Production: https://api.dklemailservice.com/api
```

## Authentication

Admin endpoints vereisen JWT authenticatie en permissies:

```
Authorization: Bearer <your_jwt_token>
```

---

## Events Management

### List Events (Public)

Haal alle events op (publiek toegankelijk).

**Endpoint:** `GET /api/events`

**Authentication:** Niet vereist

**Query Parameters:**
- `status` (optional): Filter op status (`upcoming`, `active`, `completed`, `cancelled`)
- `is_active` (optional): Filter op actieve events (true/false)

**Response:** `200 OK`

```json
[
  {
    "id": "uuid",
    "name": "De Koninklijke Loop 2025",
    "description": "Jaarlijks hardloop evenement voor KWF",
    "start_time": "2025-06-15T09:00:00Z",
    "end_time": "2025-06-15T15:00:00Z",
    "status": "upcoming",
    "status_description": "Binnenkort",
    "geofences": [
      {
        "type": "start",
        "lat": 52.370216,
        "long": 4.895168,
        "radius": 50,
        "name": "Start Lijn"
      },
      {
        "type": "checkpoint",
        "lat": 52.375216,
        "long": 4.900168,
        "radius": 30,
        "name": "Checkpoint 1"
      },
      {
        "type": "finish",
        "lat": 52.380216,
        "long": 4.905168,
        "radius": 50,
        "name": "Finish"
      }
    ],
    "event_config": {
      "max_participants": 500,
      "registration_open": true,
      "allow_team_registration": true,
      "require_medical_certificate": false
    },
    "is_active": true
  }
]
```

---

### Get Active Event (Public)

Haal het huidige actieve event op.

**Endpoint:** `GET /api/events/active`

**Authentication:** Niet vereist

**Response:** `200 OK`

Retourneert hetzelfde format als List Events, maar alleen het actieve event.

**Response wanneer geen actief event:**

```json
{
  "message": "No active event found"
}
```

---

### Get Event Details (Public)

**Endpoint:** `GET /api/events/:id`

**Authentication:** Niet vereist

**Response:** `200 OK`

Retourneert event details inclusief registratie statistieken.

---

### Create Event (Admin)

**Endpoint:** `POST /api/events`

**Authentication:** Vereist - `event:write` permissie

**Request Body:**

```json
{
  "name": "De Koninklijke Loop 2025",
  "description": "Jaarlijks hardloop evenement voor KWF",
  "start_time": "2025-06-15T09:00:00Z",
  "end_time": "2025-06-15T15:00:00Z",
  "status": "upcoming",
  "geofences": [
    {
      "type": "start",
      "lat": 52.370216,
      "long": 4.895168,
      "radius": 50,
      "name": "Start Lijn"
    }
  ],
  "event_config": {
    "max_participants": 500,
    "registration_open": true
  },
  "is_active": true
}
```

**Validation:**
- `name`: Verplicht, max 255 karakters
- `start_time`: Verplicht, moet in de toekomst zijn
- `geofences`: Verplicht, minimaal 1 start en 1 finish geofence
- `status`: Optioneel, default `upcoming`

**Response:** `201 Created`

---

### Update Event (Admin)

**Endpoint:** `PUT /api/events/:id`

**Authentication:** Vereist - `event:write` permissie

**Request Body:**

```json
{
  "name": "Updated Event Name",
  "status": "active",
  "is_active": true
}
```

**Response:** `200 OK`

---

### Delete Event (Admin)

**Endpoint:** `DELETE /api/events/:id`

**Authentication:** Vereist - `event:delete` permissie

**Response:** `200 OK`

```json
{
  "success": true,
  "message": "Event deleted successfully"
}
```

**Note:** Cascade delete verwijdert ook alle gerelateerde registraties.

---

## Event Registrations

### Get Event Registrations

Haal alle registraties voor een specifiek event op.

**Endpoint:** `GET /api/events/:id/registrations`

**Authentication:** Vereist - `event:read` permissie

**Query Parameters:**
- `status` (optional): Filter op registratie status
- `tracking_status` (optional): Filter op tracking status
- `role` (optional): Filter op deelnemer rol
- `route` (optional): Filter op afstand route

**Response:** `200 OK`

```json
{
  "event_id": "uuid",
  "event_name": "De Koninklijke Loop 2025",
  "total_registrations": 150,
  "registrations": [
    {
      "id": "uuid",
      "event_id": "uuid",
      "participant_id": "uuid",
      "participant_name": "John Doe",
      "status": "confirmed",
      "status_description": "Bevestigd",
      "tracking_status": "registered",
      "distance_route": "10km",
      "participant_role": "deelnemer",
      "participant_role_description": "Evenement deelnemer",
      "registered_at": "2025-01-01T10:00:00Z",
      "steps": 0,
      "total_distance": 0
    }
  ]
}
```

---

### Get Registration Details

**Endpoint:** `GET /api/registration/:id`

**Authentication:** Vereist - `registration:read` permissie

**Response:** `200 OK`

```json
{
  "id": "uuid",
  "event_id": "uuid",
  "event_name": "De Koninklijke Loop 2025",
  "participant_id": "uuid",
  "participant_name": "John Doe",
  "status": "confirmed",
  "status_description": "Bevestigd",
  "tracking_status": "in_progress",
  "distance_route": "10km",
  "participant_role": "deelnemer",
  "registered_at": "2025-01-01T10:00:00Z",
  "check_in_time": "2025-06-15T08:45:00Z",
  "start_time": "2025-06-15T09:05:00Z",
  "finish_time": null,
  "steps": 5000,
  "total_distance": 3.5,
  "last_location_update": "2025-06-15T09:30:00Z"
}
```

---

### Update Registration Status

**Endpoint:** `PUT /api/registration/:id`

**Authentication:** Vereist - `registration:write` permissie

**Request Body:**

```json
{
  "status": "confirmed",
  "tracking_status": "checked_in",
  "notes": "Deelnemer is aanwezig"
}
```

**Available Statuses:**
- Registration Status: `pending`, `confirmed`, `cancelled`, `completed`
- Tracking Status: `registered`, `checked_in`, `started`, `in_progress`, `finished`, `dnf`

**Response:** `200 OK`

---

### Filter Registrations by Role

**Endpoint:** `GET /api/registration/rol/:rol`

**Authentication:** Vereist - `registration:read` permissie

**Parameters:**
- `rol`: Participant role (`deelnemer`, `vrijwilliger`, `begeleider`)

**Response:** `200 OK`

Retourneert array van registraties gefilterd op rol.

---

## Geofencing

### Geofence Types

Events gebruiken geofences voor locatie tracking tijdens het evenement.

**Geofence Structure:**

```json
{
  "type": "start|checkpoint|finish",
  "lat": 52.370216,
  "long": 4.895168,
  "radius": 50,
  "name": "Checkpoint Name"
}
```

**Types:**
- `start`: Start locatie (verplicht)
- `checkpoint`: Tussenstation (optioneel, meerdere mogelijk)
- `finish`: Finish locatie (verplicht)

**Radius:**
- Opgegeven in meters
- Aanbevolen: 30-100 meter voor checkpoints, 50-200 meter voor start/finish

---

### Location Updates During Event

**Endpoint:** `POST /api/registration/:id/location`

**Authentication:** Vereist - `registration:write` permissie

**Request Body:**

```json
{
  "lat": 52.370216,
  "long": 4.895168,
  "accuracy": 10,
  "timestamp": "2025-06-15T09:30:00Z"
}
```

**Response:** `200 OK`

```json
{
  "success": true,
  "geofence_triggered": true,
  "geofence_type": "checkpoint",
  "geofence_name": "Checkpoint 1",
  "tracking_status": "in_progress",
  "total_distance": 3.5
}
```

**Automatic Status Updates:**
- Bij start geofence: status → `started`
- Bij checkpoint: distance increment
- Bij finish geofence: status → `finished`

---

## Event Status Lifecycle

### Status Flow

```
upcoming → active → completed
    ↓
cancelled (mogelijk op elk moment)
```

**Status Descriptions:**
- `upcoming`: Event is nog niet begonnen
- `active`: Event is momenteel gaande
- `completed`: Event is afgelopen
- `cancelled`: Event is geannuleerd

---

## Event Configuration

### EventConfig Object

Het `event_config` JSONB veld ondersteunt flexibele configuratie:

```json
{
  "max_participants": 500,
  "registration_open": true,
  "allow_team_registration": true,
  "require_medical_certificate": false,
  "early_bird_deadline": "2025-05-01T00:00:00Z",
  "registration_fee": 25.00,
  "age_restrictions": {
    "min_age": 16,
    "max_age": null
  },
  "routes": [
    {
      "name": "5km",
      "distance": 5.0,
      "max_participants": 200
    },
    {
      "name": "10km",
      "distance": 10.0,
      "max_participants": 200
    },
    {
      "name": "15km",
      "distance": 15.0,
      "max_participants": 100
    }
  ]
}
```

---

## Statistics & Reporting

### Event Statistics

**Endpoint:** `GET /api/events/:id/stats`

**Authentication:** Vereist - `event:read` permissie

**Response:** `200 OK`

```json
{
  "event_id": "uuid",
  "event_name": "De Koninklijke Loop 2025",
  "total_registrations": 150,
  "confirmed_count": 140,
  "cancelled_count": 5,
  "pending_count": 5,
  "checked_in_count": 120,
  "started_count": 115,
  "finished_count": 105,
  "dnf_count": 10,
  "total_steps": 1500000,
  "total_distance": 1050,
  "average_steps_per_participant": 10000,
  "by_route": {
    "5km": {
      "count": 50,
      "steps": 250000
    },
    "10km": {
      "count": 70,
      "steps": 700000
    },
    "15km": {
      "count": 30,
      "steps": 550000
    }
  },
  "by_role": {
    "deelnemer": 140,
    "vrijwilliger": 8,
    "begeleider": 2
  }
}
```

---

## Integration Examples

### React Event Registration

```typescript
import { useState } from 'react';
import api from './api';

export function EventRegistration({ eventId }: { eventId: string }) {
  const [loading, setLoading] = useState(false);

  const handleRegister = async (data: RegistrationData) => {
    setLoading(true);
    try {
      // Create participant and registration
      const response = await api.post('/api/register', {
        naam: data.naam,
        email: data.email,
        telefoon: data.telefoon,
        rol: 'deelnemer',
        afstand: '10km',
        event_id: eventId
      });

      alert('Registratie succesvol!');
      return response.data;
    } catch (error) {
      console.error('Registration failed:', error);
      alert('Registratie mislukt');
    } finally {
      setLoading(false);
    }
  };

  return (
    // Registration form UI
  );
}
```

---

### Vue Event Dashboard

```vue
<template>
  <div class="event-dashboard">
    <div class="event-header">
      <h1>{{ event.name }}</h1>
      <span class="status-badge" :class="event.status">
        {{ event.status_description }}
      </span>
    </div>

    <div class="stats-grid">
      <div class="stat-card">
        <h3>Registraties</h3>
        <p class="stat-value">{{ stats.total_registrations }}</p>
      </div>
      <div class="stat-card">
        <h3>Totaal Stappen</h3>
        <p class="stat-value">{{ stats.total_steps.toLocaleString() }}</p>
      </div>
      <div class="stat-card">
        <h3>Afstand</h3>
        <p class="stat-value">{{ stats.total_distance }}km</p>
      </div>
    </div>

    <div class="registrations-list">
      <h2>Deelnemers</h2>
      <table>
        <thead>
          <tr>
            <th>Naam</th>
            <th>Route</th>
            <th>Status</th>
            <th>Stappen</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="reg in registrations" :key="reg.id">
            <td>{{ reg.participant_name }}</td>
            <td>{{ reg.distance_route }}</td>
            <td>{{ reg.status_description }}</td>
            <td>{{ reg.steps.toLocaleString() }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import api from '@/services/api';

const props = defineProps(['eventId']);
const event = ref({});
const stats = ref({});
const registrations = ref([]);

onMounted(async () => {
  // Fetch event details
  const eventRes = await api.get(`/api/events/${props.eventId}`);
  event.value = eventRes.data;

  // Fetch stats
  const statsRes = await api.get(`/api/events/${props.eventId}/stats`);
  stats.value = statsRes.data;

  // Fetch registrations
  const regRes = await api.get(`/api/events/${props.eventId}/registrations`);
  registrations.value = regRes.data.registrations;
});
</script>
```

---

## Geofence Integration

### Checking Geofence Proximity

```typescript
interface Geofence {
  type: string;
  lat: number;
  long: number;
  radius: number;
  name: string;
}

interface Location {
  lat: number;
  long: number;
}

function isWithinGeofence(location: Location, geofence: Geofence): boolean {
  const R = 6371e3; // Earth radius in meters
  const φ1 = location.lat * Math.PI / 180;
  const φ2 = geofence.lat * Math.PI / 180;
  const Δφ = (geofence.lat - location.lat) * Math.PI / 180;
  const Δλ = (geofence.long - location.long) * Math.PI / 180;

  const a = Math.sin(Δφ/2) * Math.sin(Δφ/2) +
            Math.cos(φ1) * Math.cos(φ2) *
            Math.sin(Δλ/2) * Math.sin(Δλ/2);
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a));

  const distance = R * c; // Distance in meters

  return distance <= geofence.radius;
}

// Usage
const currentLocation = { lat: 52.370216, long: 4.895168 };
const startGeofence = event.geofences.find(g => g.type === 'start');

if (isWithinGeofence(currentLocation, startGeofence)) {
  // User is at start line
  await updateTrackingStatus(registrationId, 'started');
}
```

---

### Mobile Geolocation Integration

```typescript
// React Native / Mobile web geolocation
import { useEffect, useState } from 'react';

export function useGeolocation(registrationId: string) {
  const [location, setLocation] = useState<Location | null>(null);
  const [watching, setWatching] = useState(false);

  useEffect(() => {
    if (!watching) return;

    const watchId = navigator.geolocation.watchPosition(
      (position) => {
        const newLocation = {
          lat: position.coords.latitude,
          long: position.coords.longitude,
          accuracy: position.coords.accuracy
        };
        
        setLocation(newLocation);

        // Send location update to backend
        api.post(`/api/registration/${registrationId}/location`, newLocation)
          .catch(err => console.error('Location update failed:', err));
      },
      (error) => {
        console.error('Geolocation error:', error);
      },
      {
        enableHighAccuracy: true,
        timeout: 5000,
        maximumAge: 0
      }
    );

    return () => {
      navigator.geolocation.clearWatch(watchId);
    };
  }, [watching, registrationId]);

  return {
    location,
    startTracking: () => setWatching(true),
    stopTracking: () => setWatching(false)
  };
}
```

---

## Status Constants

### Event Status

```typescript
export enum EventStatus {
  UPCOMING = 'upcoming',
  ACTIVE = 'active',
  COMPLETED = 'completed',
  CANCELLED = 'cancelled'
}
```

### Registration Status

```typescript
export enum RegistrationStatus {
  PENDING = 'pending',
  CONFIRMED = 'confirmed',
  CANCELLED = 'cancelled',
  COMPLETED = 'completed'
}
```

### Tracking Status

```typescript
export enum TrackingStatus {
  REGISTERED = 'registered',
  CHECKED_IN = 'checked_in',
  STARTED = 'started',
  IN_PROGRESS = 'in_progress',
  FINISHED = 'finished',
  DNF = 'dnf' // Did Not Finish
}
```

---

## Best Practices

### Event Creation

1. **Geofences**: Minimaal definiëren van start en finish geofences
2. **Radius**: Gebruik voldoende radius voor betrouwbare detectie (50-100m)
3. **Configuration**: Stel `max_participants` in om overboekingen te voorkomen
4. **Timing**: Zet events op `active` status 30 minuten voor start

### Registration Management

1. **Confirmation**: Bevestig registraties binnen 24 uur
2. **Communication**: Verstuur emails bij statuswijzigingen
3. **Tracking**: Start tracking alleen voor confirmed registrations
4. **Privacy**: Respecteer GDPR - sta deelnemers toe data te verwijderen

### Real-time Tracking

1. **Frequency**: Update locatie elke 30-60 seconden tijdens event
2. **Battery**: Gebruik low-power geolocation mode waar mogelijk
3. **Offline**: Cache location updates en sync wanneer verbinding herstelt
4. **Accuracy**: Filter out hoge accuracy values (>100m) voor betere data

---

## Error Codes

| Code | HTTP Status | Beschrijving |
|------|-------------|--------------|
| `EVENT_NOT_FOUND` | 404 | Event niet gevonden |
| `EVENT_FULL` | 409 | Event heeft maximum deelnemers bereikt |
| `REGISTRATION_CLOSED` | 409 | Registratie is gesloten |
| `DUPLICATE_REGISTRATION` | 409 | Gebruiker is al geregistreerd |
| `INVALID_GEOFENCE` | 400 | Ongeldige geofence configuratie |
| `INVALID_STATUS_TRANSITION` | 400 | Ongeldige status overgang |

---

Voor meer informatie:
- [Steps & Gamification](./STEPS_GAMIFICATION.md)
- [Authentication](./AUTHENTICATION.md)
- [API Overview](./README.md)