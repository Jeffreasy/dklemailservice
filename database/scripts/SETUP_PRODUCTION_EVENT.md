# Setup Production Event - Render Deployment

**Status:** ⏳ Pending  
**Database:** Render PostgreSQL Production  
**Datum:** 2025-11-03

---

## 📋 Overzicht

De event tracking backend is volledig geïmplementeerd en getest op local Docker. Nu moet het actieve test event ook op **Render production** worden aangemaakt.

---

## 🎯 Wat Moet Gebeuren

Render production database heeft nu:
- ✅ V1.53 migratie (events tabellen)
- ✅ Permissions in RBAC  
- ✅ API endpoints actief
- ❌ **GEEN test event data** (seed event is alleen voorbeeld)

**Je moet:** Een ACTIEF event aanmaken met Dronten coordinaten.

---

## 🚀 Deployment Opties

### Optie 1: Via Render Shell (AANBEVOLEN)

```bash
# 1. SSH naar Render
# Via Render Dashboard: 
#   - Ga naar je service
#   - Click "Shell" tab
#   - Wacht tot shell connect

# 2. Connect naar PostgreSQL
psql $DATABASE_URL

# 3. Kopieer en plak deze SQL:
DELETE FROM events 
WHERE name LIKE '%Test%' 
OR (name = 'De Koninklijke Loop 2025' AND description LIKE '%Test%');

INSERT INTO events (
    name,
    description,
    start_time,
    end_time,
    status,
    is_active,
    geofences,
    event_config
) VALUES (
    'De Koninklijke Loop 2025 - GPS Test',
    'Test event voor GPS geofencing in Dronten',
    '2025-01-01 00:00:00+00',
    '2025-12-31 23:59:59+00',
    'active',
    true,
    '[
        {
            "type": "start",
            "lat": 52.5185,
            "long": 5.7220,
            "radius": 500,
            "name": "Start - Dronten Spiegelstraat"
        },
        {
            "type": "checkpoint",
            "lat": 52.5228,
            "long": 5.7306,
            "radius": 300,
            "name": "Checkpoint Noord"
        },
        {
            "type": "checkpoint",
            "lat": 52.5142,
            "long": 5.7134,
            "radius": 300,
            "name": "Checkpoint Zuid"
        },
        {
            "type": "finish",
            "lat": 52.5185,
            "long": 5.7220,
            "radius": 500,
            "name": "Finish - Dronten Spiegelstraat"
        }
    ]'::jsonb,
    '{
        "minStepsInterval": 10,
        "requireGeofenceCheckin": true,
        "distanceThreshold": 100,
        "accuracyLevel": "balanced"
    }'::jsonb
);

# 4. Verify
SELECT 
    id,
    name,
    status,
    is_active,
    jsonb_array_length(geofences) as geofence_count
FROM events
WHERE is_active = true;

# 5. Exit
\q
```

### Optie 2: Via API (VEREIST ADMIN LOGIN)

**Stap 1: Login als Admin**
```bash
curl -X POST https://dklemailservice.onrender.com/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@dekoninklijkeloop.nl",
    "wachtwoord": "Bootje@12"
  }'

# Response bevat:
# {"token": "eyJ..."}
# Kopieer de token!
```

**Stap 2: Maak Event Aan**
```bash
# Vervang <TOKEN> met je token van stap 1
curl -X POST https://dklemailservice.onrender.com/api/events \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "De Koninklijke Loop 2025 - GPS Test",
    "description": "Test event voor GPS geofencing in Dronten",
    "start_time": "2025-01-01T00:00:00Z",
    "end_time": "2025-12-31T23:59:59Z",
    "status": "active",
    "is_active": true,
    "geofences": [
      {
        "type": "start",
        "lat": 52.5185,
        "long": 5.7220,
        "radius": 500,
        "name": "Start - Dronten Spiegelstraat"
      },
      {
        "type": "checkpoint",
        "lat": 52.5228,
        "long": 5.7306,
        "radius": 300,
        "name": "Checkpoint Noord"
      },
      {
        "type": "checkpoint",
        "lat": 52.5142,
        "long": 5.7134,
        "radius": 300,
        "name": "Checkpoint Zuid"
      },
      {
        "type": "finish",
        "lat": 52.5185,
        "long": 5.7220,
        "radius": 500,
        "name": "Finish - Dronten Spiegelstraat"
      }
    ],
    "event_config": {
      "minStepsInterval": 10,
      "requireGeofenceCheckin": true,
      "distanceThreshold": 100,
      "accuracyLevel": "balanced"
    }
  }'
```

**Stap 3: Verify**
```bash
curl https://dklemailservice.onrender.com/api/events/active

# Expected: 200 OK met event data
# Als 404: Event niet aangemaakt of niet actief
```

### Optie 3: Via SQL File Upload (Render Dashboard)

1. Download [`create_test_event_dronten.sql`](create_test_event_dronten.sql)
2. Ga naar Render Dashboard → Database → SQL Editor
3. Upload en execute het script
4. Verify met: `SELECT * FROM events WHERE is_active = true;`

---

## ✅ Verification Checklist

Na het aanmaken van het event, verify:

```bash
# 1. Test actief event endpoint
curl https://dklemailservice.onrender.com/api/events/active

# Expected response:
# - Status: 200 OK (niet 404)
# - name: "De Koninklijke Loop 2025 - GPS Test"
# - status: "active"
# - is_active: true
# - geofences: array met 4 items
# - start_time: in verleden
# - end_time: in toekomst

# 2. Test alle events endpoint
curl https://dklemailservice.onrender.com/api/events

# Expected: Array met minimaal 1 event

# 3. Test specifiek event
curl https://dklemailservice.onrender.com/api/events/4fb1dd30-1465-4a0d-b5e0-6ae814109182

# Expected: 200 OK met complete event data
```

---

## 📍 Event Configuratie

### Geofence Locaties (Dronten)

**Start/Finish Locatie:**
- **Adres:** Spiegelstraat 6, Dronten
- **Coordinaten:** 52.5185, 5.7220
- **Radius:** 500 meter (groot gebied voor makkelijk testen)

**Checkpoint Noord:**
- **Coordinaten:** 52.5228, 5.7306
- **Radius:** 300 meter

**Checkpoint Zuid:**
- **Coordinaten:** 52.5142, 5.7134
- **Radius:** 300 meter

### Event Timing

- **Start:** 2025-01-01 00:00:00 UTC (in verleden = altijd actief)
- **End:** 2025-12-31 23:59:59 UTC (hele jaar beschikbaar)
- **Status:** `active`
- **Is Active:** `true`

Dit zorgt dat het event **altijd beschikbaar is voor testing**.

---

## 🧪 Testing na Deployment

### 1. Backend Verificatie

```bash
# Production endpoint
curl https://dklemailservice.onrender.com/api/events/active

# Should return JSON met:
{
  "id": "...",
  "name": "De Koninklijke Loop 2025 - GPS Test",
  "geofences": [
    { "type": "start", "lat": 52.5185, "long": 5.7220, "radius": 500 }
  ],
  "status": "active",
  "is_active": true
}
```

### 2. Mobile App Testing

**In React Native app:**
```typescript
// App zal automatisch event ophalen
const event = await fetch('https://dklemailservice.onrender.com/api/events/active')
    .then(r => r.json());

// Should log:
console.log('Event:', event.name);
// "De Koninklijke Loop 2025 - GPS Test"

console.log('Geofences:', event.geofences.length);
// 4
```

### 3. GPS Emulator Test

**iOS Simulator:**
```
1. Open app
2. Enable geofencing
3. Debug > Simulate Location > Custom Location
4. Lat: 52.5185, Long: 5.7220
5. Status should show: "✓ Binnen Gebied"
```

**Android Emulator:**
```
1. Open app  
2. Enable geofencing
3. Extended Controls > Location
4. Lat: 52.5185, Long: 5.7220
5. Click "Send"
6. Status should show: "✓ Binnen Gebied"
```

---

## 🔧 Troubleshooting

### "404 Not Found" na Script

**Oorzaak:** Event bestaat maar is niet "actief"

**Fix:**
```sql
-- Check event status
SELECT id, name, status, is_active, start_time, end_time
FROM events
ORDER BY created_at DESC;

-- Activate event
UPDATE events
SET is_active = true, status = 'active'
WHERE name LIKE '%GPS Test%';
```

### "Event in Toekomst"

**Oorzaak:** start_time > NOW

**Fix:**
```sql
UPDATE events
SET start_time = '2025-01-01 00:00:00+00'
WHERE name LIKE '%GPS Test%';
```

### "Multiple Events Found"

**Oorzaak:** Seed event + test event beide actief

**Fix:**
```sql
-- Deactiveer oude events, houd alleen test event
UPDATE events
SET is_active = false
WHERE name NOT LIKE '%GPS Test%';
```

---

## 📊 Success Metrics

Na deployment moet je zien:

**✅ Local Docker:**
```bash
curl http://localhost:8082/api/events/active
# 200 OK ✅

curl http://localhost:8082/api/events
# Array met 2-3 events ✅
```

**✅ Render Production:**
```bash
curl https://dklemailservice.onrender.com/api/events/active
# 200 OK ✅ (na Optie 1 of 2 uitvoeren)

# Response bevat:
# - Dronten coordinaten (52.5185, 5.7220)
# - 4 geofences
# - status: "active"
# - is_active: true
```

**✅ Mobile App:**
```
- Open app als deelnemer
- Enable "Event Locatie Tracking"
- See event: "De Koninklijke Loop 2025 - GPS Test"
- See geofences: 4 stuks
- See status calculated with distance
```

---

## 🎉 Na Succesvolle Setup

Je hebt dan:
1. ✅ Backend API volledig werkend
2. ✅ Database met actief test event
3. ✅ Geofences met Dronten coordinaten
4. ✅ Ready voor mobile app GPS testing

**Next Step:** Test de mobile app op physical device bij Spiegelstraat 6, Dronten!

---

**Script Location:** [`database/scripts/create_test_event_dronten.sql`](create_test_event_dronten.sql)  
**Deployment:** Run Optie 1 (Render Shell) of Optie 2 (API)  
**Verify:** `curl https://dklemailservice.onrender.com/api/events/active` → 200 OK