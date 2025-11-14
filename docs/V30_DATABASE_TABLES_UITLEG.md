# V30 Database Tables - Heldere Uitleg

## 🎯 Overzicht: 3 Belangrijkste Tables

Bij registratie worden **3 tables** gebruikt:

1. **`participants`** = De Persoon (WIE)
2. **`event_registrations`** = De Deelname (HOE & WANNEER)
3. **`gebruikers`** = Voor App Login (alleen full accounts)

---

## 📋 Table 1: `participants` - DE PERSOON

**Doel:** Wie is deze persoon? Basis persoonsgegeven + account type.

### Wat staat erin (per persoon):

```sql
CREATE TABLE participants (
  -- Identificatie
  id UUID PRIMARY KEY,
  naam TEXT NOT NULL,
  email TEXT NOT NULL,
  telefoon TEXT,
  
  -- V30: Account Type Info
  account_type TEXT NOT NULL,           -- 'full' of 'temporary'
  wachtwoord_hash TEXT,                 -- Alleen als account_type = 'full'
  has_app_access BOOLEAN DEFAULT false, -- true = kan app gebruiken
  
  -- Voor temporary accounts
  registration_year INTEGER,            -- Bijv. 2026 (geldig voor dat jaar)
  
  -- Links naar andere tables
  gebruiker_id UUID,                    -- Link naar gebruikers (bij full account)
  
  -- Timestamps
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ
);
```

### Voorbeeld Data:

```
ID: "abc-123"
Naam: "Jeffrey Lavente"
Email: "jeffrey@gmail.com"
Account Type: "full"
Wachtwoord Hash: "$2a$10$..." 
Has App Access: true
Gebruiker ID: "xyz-789"
```

### Wat staat ER NIET in:
- ❌ GEEN rol (Deelnemer/Begeleider)  
- ❌ GEEN afstand (2.5KM, 6KM, etc)
- ❌ GEEN ondersteuningskeuze
- ❌ GEEN event-specifieke info

**Waarom niet?** → Omdat een persoon bij meerdere events kan deelnemen in verschillende rollen!

---

## 📋 Table 2: `event_registrations` - DE DEELNAME

**Doel:** HOE neemt deze persoon deel aan een specifiek evenement?

### Wat staat erin (per deelname):

```sql
CREATE TABLE event_registrations (
  -- Identificatie
  id UUID PRIMARY KEY,
  event_id UUID NOT NULL,          -- Welk evenement? (DKL 2026, DKL 2027, etc)
  participant_id UUID NOT NULL,    -- Wie neemt deel? (link naar participants)
  
  -- Event-specifieke keuzes
  participant_role_name TEXT,      -- ROL: "Deelnemer", "Begeleider", "Vrijwilliger"
  distance_route TEXT,             -- AFSTAND: "2.5 KM", "6 KM", "10 KM", "15 KM"
  ondersteuning TEXT,              -- "Ja", "Nee", "Anders"
  bijzonderheden TEXT,             -- Extra info bij ondersteuning
  
  -- Status tracking
  status TEXT DEFAULT 'registered', -- "registered", "paid", "confirmed", etc
  
  -- Timestamps
  registered_at TIMESTAMPTZ,
  check_in_time TIMESTAMPTZ,
  start_time TIMESTAMPTZ,
  finish_time TIMESTAMPTZ,
  
  -- Steps tracking
  steps INTEGER DEFAULT 0
);
```

### Voorbeeld Data:

```
ID: "reg-456"
Event ID: "dkl-2026"
Participant ID: "abc-123" (← Dit is Jeffrey)

Participant Role Name: "Deelnemer"     ← HIER staat de ROL
Distance Route: "2.5 KM"                ← HIER staat de AFSTAND
Ondersteuning: "Nee"                    ← HIER staat ONDERSTEUNING
Bijzonderheden: ""
Status: "registered"
Steps: 0
```

---

## 📋 Table 3: `gebruikers` - APP LOGIN

**Doel:** Alleen voor mensen die de DKL Step App kunnen gebruiken (full accounts).

### Wat staat erin:

```sql
CREATE TABLE gebruikers (
  id UUID PRIMARY KEY,
  naam TEXT NOT NULL,
  email TEXT NOT NULL,
  wachtwoord_hash TEXT NOT NULL,  -- Voor login in app
  is_actief BOOLEAN DEFAULT true
);
```

### Voorbeeld Data:

```
ID: "xyz-789"
Naam: "Jeffrey Lavente"
Email: "jeffrey@gmail.com"
Wachtwoord Hash: "$2a$10$..."
Is Actief: true
```

**Wanneer wordt dit aangemaakt?**
- ✅ Bij registratie met `want_account = true` (full account)
- ❌ NIET bij `want_account = false` (temporary account)

---

## 🔄 COMPLETE REGISTRATIE FLOW

### Scenario A: FULL ACCOUNT (met app toegang)

**Formulier input:**
```
Naam: Jeffrey
Email: jeffrey@gmail.com
Telefoon: 06-12345678
Rol: Deelnemer
Afstand: 2.5 KM
Ondersteuning: Nee
Want Account: ✅ JA
Wachtwoord: geheim123
```

**Backend maakt aan:**

1. **`participants` record:**
```
{
  naam: "Jeffrey",
  email: "jeffrey@gmail.com", 
  telefoon: "06-12345678",
  account_type: "full",
  wachtwoord_hash: "$2a$10...",
  has_app_access: true
}
```

2. **`gebruikers` record:** (voor app login)
```
{
  naam: "Jeffrey",
  email: "jeffrey@gmail.com",
  wachtwoord_hash: "$2a$10..."
}
```

3. **`event_registrations` record:** (event deelname)
```
{
  event_id: "dkl-2026",
  participant_id: [ID van Jeffrey],
  participant_role_name: "Deelnemer",    ← Rol zit HIER
  distance_route: "2.5 KM",              ← Afstand zit HIER
  ondersteuning: "Nee",                  ← Ondersteuning zit HIER
  bijzonderheden: ""
}
```

---

### Scenario B: TEMPORARY ACCOUNT (geen app)

**Formulier input:**
```
Naam: Hans
Email: hans@gmail.com
Rol: Vrijwilliger
Afstand: 10 KM
Ondersteuning: Anders
Bijzonderheden: "Ik help bij de finish"
Want Account: ❌ NEE (geen wachtwoord)
```

**Backend maakt aan:**

1. **`participants` record:**
```
{
  naam: "Hans",
  email: "hans@gmail.com",
  account_type: "temporary",       ← Temporary!
  registration_year: 2026,         ← Alleen geldig voor 2026
  wachtwoord_hash: NULL,           ← GEEN wachtwoord
  has_app_access: false,           ← KAN NIET app gebruiken
  gebruiker_id: NULL               ← GEEN gebruiker link
}
```

2. **GEEN `gebruikers` record** (omdat temporary)

3. **`event_registrations` record:**
```
{
  event_id: "dkl-2026",
  participant_id: [ID van Hans],
  participant_role_name: "Vrijwilliger", ← Rol zit HIER
  distance_route: "10 KM",               ← Afstand zit HIER  
  ondersteuning: "Anders",               ← Ondersteuning zit HIER
  bijzonderheden: "Ik help bij finish"   ← Bijzonderheden zit HIER
}
```

---

## 🔍 WAT WAS HET PROBLEEM?

### Voor V34 Migration:

---

## 👟 WAAR WORDEN STAPPEN BIJGEHOUDEN?

**Antwoord:** In de `event_registrations` table!

### Waarom niet in `participants`?

Omdat een persoon bij meerdere events kan deelnemen:
```
Jeffrey's stappen:
- DKL 2026: 50.000 stappen  ← event_registrations record 1
- DKL 2027: 75.000 stappen  ← event_registrations record 2
```

### Hoe werkt het?

**In `event_registrations` table:**
```sql
{
  id: "reg-2026-jeffrey",
  participant_id: "jeffrey-id",
  event_id: "dkl-2026",
  
  participant_role_name: "Deelnemer",
  distance_route: "2.5 KM",
  
  steps: 50000,              ← HIER staan de stappen!
  
  registered_at: "2025-11-10",
  start_time: "2026-05-04 09:00",
  finish_time: "2026-05-04 11:30"
}
```

### WebSocket Updates:

Wanneer de app stappen update via `/api/ws/steps`:

```javascript
// App stuurt:
{
  "participant_id": "jeffrey-id",
  "steps": 1000  // Nieuwe stappen count
}

// Backend update:
UPDATE event_registrations 
SET steps = 1000
WHERE participant_id = 'jeffrey-id' 
  AND event_id = [active_event_id];
```

### Leaderboard Query:

```sql
-- Top 10 deelnemers voor DKL 2026:
SELECT p.naam, er.steps, er.distance_route
FROM event_registrations er
JOIN participants p ON p.id = er.participant_id  
WHERE er.event_id = 'dkl-2026'
ORDER BY er.steps DESC
LIMIT 10;
```

**Resultaat:**
```
Naam            | Steps   | Afstand
----------------|---------|--------
Jeffrey         | 75,000  | 2.5 KM
Hans            | 50,000  | 10 KM
Maria           | 45,000  | 6 KM
```

### Samenvatting Stappen:

| Vraag | Table | Kolom |
|-------|-------|-------|
| Hoeveel stappen heeft Jeffrey bij DKL 2026? | `event_registrations` | `steps` |
| Wie heeft de meeste stappen bij dit event? | `event_registrations` | ORDER BY `steps` DESC |
| Update stappen real-time via app | `event_registrations` | UPDATE `steps` |

**De stappen zijn EVENT-SPECIFIEK, dus ze horen bij de event_registrations table, niet bij de participants table!**
```
participants table had:
- naam ✅
- email ✅
- rol ❌ (hoort niet hier!)
- afstand ❌ (hoort niet hier!)
- ondersteuning ❌ (hoort niet hier!)
- bijzonderheden ❌ (hoort niet hier!)

event_registrations table had:
- participant_role_name ✅
- distance_route ✅
- ondersteuning ✅
- bijzonderheden ✅
```

**Probleem:** Rol/Afstand stond in BEIDE tables! (dubbel opgeslagen)

### Na V34 Migration:
```
participants table heeft:
- naam ✅
- email ✅
- account_type ✅
- wachtwoord_hash ✅
- has_app_access ✅
❌ GEEN rol/afstand/ondersteuning meer!

event_registrations table heeft:
- participant_role_name ✅
- distance_route ✅
- ondersteuning ✅
- bijzonderheden ✅
```

**Opgelost:** Data staat nu alleen op de JUISTE plek!

---

## 📖 Samenvatting

| Vraag | Table | Veld |
|-------|-------|------|
| Wie is deze persoon? | `participants` | `naam`, `email` |
| Kan deze persoon de app gebruiken? | `participants` | `has_app_access` |
| Full of temporary account? | `participants` | `account_type` |
| In welke ROL neemt iemand deel? | `event_registrations` | `participant_role_name` |
| Welke AFSTAND loopt iemand? | `event_registrations` | `distance_route` |
| Heeft iemand ONDERSTEUNING nodig? | `event_registrations` | `ondersteuning` |
| Wat zijn de BIJZONDERHEDEN? | `event_registrations` | `bijzonderheden` |
| Kan iemand inloggen in de app? | `gebruikers` | `email`, `wachtwoord_hash` |

**Simpel gezegd:**
- `participants` = Wie ben je + welk soort account
- `event_registrations` = Hoe doe je mee aan dit evenement  
- `gebruikers` = Login voor de app (alleen full accounts)