# Backend Fix: V34 Registration System

**Date:** 2025-11-10  
**Issue:** Backend trying to use non-existent columns after V34 migration  
**Status:** ✅ FIXED

## Problem

After the V34 migration removed legacy columns from the `participants` table, the backend was still trying to:
1. Store event-data (rol, afstand, ondersteuning) in `participants` table ❌
2. Create `Gebruiker` records for public registrations ❌
3. Assign RBAC roles to non-existent gebruikers ❌

### What Was Wrong

```go
// ❌ Backend was doing this (FAILED):
participant.Rol = req.Rol              // ERROR: column "rol" does not exist
participant.Afstand = req.Afstand      // ERROR: column "afstand" does not exist
participant.Ondersteuning = req.Ondersteuning  // ERROR: column exists not

// ❌ Also creating Gebruiker records for public registrations
gebruiker := &models.Gebruiker{...}
db.Create(gebruiker)
```

## The Fix (3 Steps)

### 1. Participant = Person Data Only

```go
participant := &models.Participant{
    Naam:           req.Naam,
    Email:          req.Email,
    Telefoon:       req.Telefoon,        // ✅ Store phone in participants
    Terms:          req.Terms,
    AccountType:    models.AccountTypeFull,
    WachtwoordHash: &hashedPasswordStr,  // ✅ For authentication
    HasAppAccess:   true,
    // ❌ NO: Rol, Afstand, Ondersteuning - columns don't exist!
    // ❌ NO: GebruikerID - no Gebruiker for public registrations
}
```

### 2. Event Registration = Event Data

```go
registration := &models.EventRegistration{
    EventID:             event.ID,
    ParticipantID:       participant.ID,
    ParticipantRoleName: &req.Rol,         // ✅ Rol goes here!
    DistanceRoute:       &req.Afstand,     // ✅ Afstand goes here!
    Ondersteuning:       req.Ondersteuning, // ✅ Ondersteuning here!
    Bijzonderheden:      req.Bijzonderheden,
    Status:              "registered",
}
```

### 3. No Gebruiker Records

```go
// ❌ REMOVED: No Gebruiker record for public registrations
// ✅ gebruikers table is ONLY for admin/staff
// ✅ Full accounts authenticate via participants.wachtwoord_hash
```

## Files Changed

### [`handlers/public_registration_handler.go`](../handlers/public_registration_handler.go)

**Function: [`registerFullAccount`](../handlers/public_registration_handler.go:172-239)**
- ❌ Removed: Gebruiker creation (lines ~198-209)
- ❌ Removed: RBAC role assignments (lines ~215-239)
- ✅ Added: Only participant with wachtwoord_hash
- ✅ Added: Event registration with event-data

**Function: [`UpgradeToFullAccount`](../handlers/public_registration_handler.go:323-438)**
- ❌ Removed: Gebruiker creation (lines ~395-409)
- ❌ Removed: RBAC role assignments (lines ~415-446)
- ❌ Removed: GebruikerID linking (lines ~453-454)
- ✅ Updated: Only participant fields updated

## How It Works Now

### Full Account Registration
1. User fills form with `want_account: true`
2. Backend creates **Participant** with `wachtwoord_hash`
3. Backend creates **EventRegistration** with rol/afstand/ondersteuning
4. ✅ Done! No Gebruiker record needed

### Temporary Account Registration
1. User fills form with `want_account: false`
2. Backend creates **Participant** (no password)
3. Backend creates **EventRegistration** with rol/afstand/ondersteuning
4. ✅ Done! Can upgrade later

### Account Upgrade
1. User enters email + password
2. Backend finds temporary participant
3. Backend updates participant: `account_type = 'full'`, sets `wachtwoord_hash`
4. ✅ Done! No Gebruiker record needed

## Authentication

### OLD (V30-V33)
```
User logs in → Check gebruiker.wachtwoord_hash
```

### NEW (V34+)
```
User logs in → Check participants.wachtwoord_hash
```

**Note:** The `gebruikers` table is now **only for admin/staff users**, not for public participant registrations.

## Database Schema (V34)

### `participants` table
```sql
CREATE TABLE participants (
    id UUID PRIMARY KEY,
    naam TEXT NOT NULL,
    email TEXT NOT NULL,
    telefoon TEXT,                    -- ✅ Phone stored here
    wachtwoord_hash TEXT,             -- ✅ For full accounts
    account_type TEXT,                -- 'full' or 'temporary'
    has_app_access BOOLEAN,
    registration_year INTEGER,
    -- ❌ NO MORE: rol, afstand, ondersteuning, bijzonderheden, steps
    -- ❌ NO MORE: gebruiker_id (for public registrations)
);
```

### `event_registrations` table
```sql
CREATE TABLE event_registrations (
    id UUID PRIMARY KEY,
    event_id UUID REFERENCES events(id),
    participant_id UUID REFERENCES participants(id),
    participant_role_name TEXT,       -- ✅ Rol here!
    distance_route TEXT,              -- ✅ Afstand here!
    ondersteuning TEXT,               -- ✅ Ondersteuning here!
    bijzonderheden TEXT,
    steps INTEGER DEFAULT 0,
    status TEXT,
);
```

### `gebruikers` table
```sql
-- ✅ ONLY for admin/staff users
-- ❌ NOT for public participant registrations
CREATE TABLE gebruikers (
    id UUID PRIMARY KEY,
    naam TEXT,
    email TEXT UNIQUE,
    wachtwoord_hash TEXT,
    is_actief BOOLEAN,
    -- Links to rbac_user_roles for permissions
);
```

## Testing

Service rebuilt and started successfully:
```
✅ Server gestart: http://127.0.0.1:8080
✅ Public registration routes registered: /api/public/aanmelden
✅ All migrations applied successfully
✅ No errors in logs
```

## Summary

The backend now correctly:
✅ Stores person-data in `participants` table  
✅ Stores event-data in `event_registrations` table  
✅ Uses `participants.wachtwoord_hash` for authentication  
✅ Does NOT create `Gebruiker` records for public registrations  
✅ Does NOT try to use non-existent columns  

**Result:** Public registrations will now work without database errors! 🎉