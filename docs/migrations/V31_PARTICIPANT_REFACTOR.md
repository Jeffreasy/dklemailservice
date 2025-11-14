# V31: Complete V28 Participant Refactor

**Migration File:** `database/migrations/V31__complete_v28_participant_refactor.sql`  
**Date:** 2025-11-10  
**Status:** ✅ Applied  
**Risk Level:** 🔴 HIGH (Data Migration)

---

## Overview

V31 completes the participant/event registration data separation that was initiated in V28 but failed. This migration is critical because it:

1. **Fixes a failed migration** - V28 didn't complete successfully
2. **Separates concerns** - Moves event-specific data from `participants` to `event_registrations`
3. **Prevents data duplication** - One participant can have multiple event registrations
4. **Improves scalability** - Better architecture for multi-year events

---

## Why Was V28 Needed?

### The Problem (Pre-V28)

**Original Design Flaw:**
```
participants table = PERSON DATA + EVENT DATA (mixed!)
├── naam, email, telefoon (person data) ✅
├── rol, afstand (event-specific) ❌
├── steps, ondersteuning (event-specific) ❌
└── bijzonderheden, status (event-specific) ❌
```

**Issues:**
- A participant could only register for ONE event per year
- Event-specific data (role, distance, steps) was stored at person level
- No historical data - previous year data was lost
- Participant could NOT change role/distance between events

### The Solution (V28/V31)

**New Architecture:**
```
participants table = PERSON DATA ONLY
├── id, naam, email, telefoon
├── account_type (V30: temporary/full)
└── has_app_access (V30: boolean)

event_registrations table = EVENT PARTICIPATION
├── participant_id → FK to participants
├── event_id → FK to events
├── participant_role_name (Deelnemer, Begeleider, Vrijwilliger)
├── distance_route (2.5KM, 6KM, 10KM, 15KM)
├── steps (real-time tracking per event!)
├── ondersteuning, bijzonderheden
└── status, tracking_status
```

**Benefits:**
- ✅ One participant = multiple event registrations
- ✅ Historical data preserved per event
- ✅ Can have different roles/distances per event
- ✅ Steps tracked separately per event

---

## What V31 Does

### Part 1: Schema Changes

#### Added Columns to `event_registrations`

```sql
ALTER TABLE event_registrations
    ADD COLUMN IF NOT EXISTS steps INTEGER DEFAULT 0,
    ADD COLUMN IF NOT EXISTS ondersteuning TEXT,
    ADD COLUMN IF NOT EXISTS test_mode BOOLEAN DEFAULT false,
    ADD COLUMN IF NOT EXISTS email_verzonden BOOLEAN DEFAULT false,
    ADD COLUMN IF NOT EXISTS email_verzonden_op TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS behandeld_door TEXT,
    ADD COLUMN IF NOT EXISTS behandeld_op TIMESTAMPTZ;
```

**Why These Columns?**
- `steps` - Steps are EVENT-SPECIFIC (Jeffrey has 50k steps in DKL 2026, 75k in DKL 2027)
- `ondersteuning` - Support needs can differ per event
- `test_mode` - Testing flag per registration
- `email_verzonden*` - Email notification tracking per registration
- `behandeld_*` - Admin processing tracking per registration

#### Handle `current_steps` → `steps` Rename

V28 may have created `current_steps` instead of `steps`. V31 handles this:

```sql
-- If both exist: merge and drop current_steps
-- If only current_steps exists: rename to steps
```

### Part 2: Data Migration

**Critical:** Migrates existing participant data to event registrations.

```sql
INSERT INTO event_registrations (
    event_id,                    -- Active event
    participant_id,              -- From participants.id
    registered_at,               -- From participants.created_at
    steps,                       -- From participants.steps
    ondersteuning,              -- From participants.ondersteuning
    bijzonderheden,             -- From participants.bijzonderheden
    status,                      -- From participants.status
    distance_route,              -- From participants.distance_route
    participant_role_name,       -- From participants.participant_role_name
    email_verzonden,            -- From participants.email_verzonden
    -- ... etc
)
SELECT ... FROM participants p
WHERE NOT EXISTS (
    SELECT 1 FROM event_registrations er
    WHERE er.participant_id = p.id 
      AND er.event_id = <active_event_id>
);
```

**Safety:** Only migrates if:
- Active event exists (via `get_active_event()` function)
- No duplicate registration exists already

### Part 3: Indexes

```sql
CREATE INDEX IF NOT EXISTS idx_event_registrations_participant_role 
    ON event_registrations(participant_role_name);

CREATE INDEX IF NOT EXISTS idx_event_registrations_distance 
    ON event_registrations(distance_route);

CREATE INDEX IF NOT EXISTS idx_event_registrations_steps 
    ON event_registrations(steps) WHERE steps > 0;
```

**Why:** Optimize queries for:
- Finding participants by role
- Filtering by distance
- Leaderboards (steps > 0)

### Part 4: Legacy Columns (Kept for Now)

V31 **DOES NOT** drop old columns from `participants`:
- `afstand`, `rol`, `ondersteuning`, `bijzonderheden`
- `steps`, `status`, `email_verzonden`, etc.

**Why?** Backward compatibility during transition period.

**Next Step:** V34 will remove these (see [`V34_BREAKING_CHANGES.md`](V34_BREAKING_CHANGES.md))

---

## Migration Flow Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                    V31 MIGRATION FLOW                        │
└─────────────────────────────────────────────────────────────┘

1. CHECK ACTIVE EVENT
   ├─ get_active_event() returns event_id? 
   │  ├─ YES → Continue
   │  └─ NO → Migration skips data migration (only schema)
   │
2. ADD MISSING COLUMNS
   ├─ event_registrations gets: steps, ondersteuning, etc.
   ├─ Handle current_steps rename
   └─ Columns are now aligned with participants
   │
3. DATA MIGRATION
   ├─ For each participant without event_registration:
   │  ├─ Copy person data → participants (already there)
   │  ├─ Copy event data → event_registrations (NEW!)
   │  └─ Link via participant_id
   │
4. ADD PERFORMANCE INDEXES
   ├─ idx_event_registrations_participant_role
   ├─ idx_event_registrations_distance
   └─ idx_event_registrations_steps
   │
5. KEEP LEGACY COLUMNS (for now)
   └─ Old columns still exist in participants for backward compat
```

---

## Verification Steps

### After Migration

```sql
-- 1. Check migration applied
SELECT * FROM migraties WHERE version = 'V31__complete_v28_participant_refactor';

-- 2. Verify new columns exist
\d event_registrations
-- Should show: steps, ondersteuning, test_mode, etc.

-- 3. Check data migration count
SELECT COUNT(*) FROM event_registrations;
SELECT COUNT(*) FROM participants;
-- Should have registrations for all participants (if active event exists)

-- 4. Verify data integrity
SELECT 
    p.id,
    p.naam,
    p.email,
    er.participant_role_name,
    er.distance_route,
    er.steps
FROM participants p
LEFT JOIN event_registrations er ON er.participant_id = p.id
WHERE p.account_type = 'temporary'  -- Check temporary accounts
LIMIT 10;

-- 5. Check indexes
\di event_registrations*
-- Should show new indexes
```

---

## Rollback Procedure

⚠️ **WARNING:** Rollback will **LOSE DATA** migrated to event_registrations!

### Step 1: Backup First

```bash
pg_dump -U postgres -d dklemailservice > backup_before_v31_rollback.sql
```

### Step 2: Rollback SQL

```sql
-- rollback_V31.sql
BEGIN;

-- Drop added columns (THIS DELETES DATA!)
ALTER TABLE event_registrations
    DROP COLUMN IF EXISTS steps,
    DROP COLUMN IF EXISTS ondersteuning,
    DROP COLUMN IF EXISTS test_mode,
    DROP COLUMN IF EXISTS email_verzonden,
    DROP COLUMN IF EXISTS email_verzonden_op,
    DROP COLUMN IF EXISTS behandeld_door,
    DROP COLUMN IF EXISTS behandeld_op;

-- Drop indexes
DROP INDEX IF EXISTS idx_event_registrations_participant_role;
DROP INDEX IF EXISTS idx_event_registrations_distance;
DROP INDEX IF EXISTS idx_event_registrations_steps;

-- Remove migration record
DELETE FROM migraties WHERE version = 'V31__complete_v28_participant_refactor';

COMMIT;
```

### Step 3: Execute Rollback

```bash
psql -U postgres -d dklemailservice -f rollback_V31.sql
```

---

## Code Impact

### Models Updated

**Before V31:**
```go
// Participant had event-specific fields
type Participant struct {
    ID     uuid.UUID
    Naam   string
    Email  string
    Steps  int        // ❌ Should be in event_registrations
    Rol    string     // ❌ Should be in event_registrations
}
```

**After V31:**
```go
// Participant = person data only
type Participant struct {
    ID          uuid.UUID
    Naam        string
    Email       string
    AccountType string  // V30: temporary/full
}

// EventRegistration = event participation
type EventRegistration struct {
    ID                   uuid.UUID
    EventID              uuid.UUID
    ParticipantID        uuid.UUID
    Steps                int        // ✅ Event-specific
    ParticipantRoleName  string     // ✅ Event-specific
    DistanceRoute        string     // ✅ Event-specific
}
```

### Repository Changes

**Query Pattern Change:**

```go
// BEFORE V31: Get participant with steps
participant, err := repo.GetParticipantByID(id)
steps := participant.Steps  // ❌ Wrong - steps are event-specific

// AFTER V31: Get participant with event registration
participant, err := participantRepo.GetByID(id)
registration, err := eventRegRepo.GetByParticipantAndEvent(id, eventID)
steps := registration.Steps  // ✅ Correct - steps per event
```

### WebSocket Steps Update

**Before V31:**
```go
// Update participant.steps
db.Model(&Participant{}).Where("id = ?", participantID).Update("steps", newSteps)
```

**After V31:**
```go
// Update event_registrations.steps
db.Model(&EventRegistration{}).
    Where("participant_id = ? AND event_id = ?", participantID, eventID).
    Update("steps", newSteps)
```

---

## Testing Checklist

- [ ] Migration applies without errors
- [ ] All participants have event_registrations (if active event exists)
- [ ] Data counts match (participants count = event_registrations count for active event)
- [ ] Steps are preserved correctly
- [ ] Indexes are created
- [ ] WebSocket steps updates work correctly
- [ ] API endpoints return correct data structure
- [ ] Frontend displays event-specific data correctly

---

## Known Issues & Warnings

### Issue 1: No Active Event

**Symptom:** Migration runs but doesn't migrate data

```
[V31] No active event found. Cannot migrate participant data.
[V31] You can manually create event_registrations or set an active event.
```

**Solution:**
```sql
-- Check if event exists
SELECT * FROM events WHERE status = 'active';

-- If no active event, either:
-- 1. Set an event as active
UPDATE events SET status = 'active' WHERE id = '<event_id>';

-- 2. Or manually migrate later when event is created
```

### Issue 2: Duplicate current_steps

**Symptom:** Both `steps` and `current_steps` columns exist

**Resolution:** V31 handles this automatically:
```sql
UPDATE event_registrations 
SET steps = COALESCE(current_steps, steps, 0)
WHERE steps = 0 AND current_steps > 0;

ALTER TABLE event_registrations DROP COLUMN current_steps;
```

---

## Related Migrations

| Migration | Relation | Description |
|-----------|----------|-------------|
| **V28** | Predecessor | Original participant refactor (failed) |
| **V30** | Parallel | Added account_type system to participants |
| **V32** | Related | Added missing columns to supporting tables |
| **V34** | Successor | Removes legacy columns from participants |

---

## Architecture Documentation

For complete architecture details, see:
- [`docs/architecture/DATABASE.md`](../architecture/DATABASE.md) - Database schema
- [`docs/V30_DATABASE_TABLES_UITLEG.md`](../V30_DATABASE_TABLES_UITLEG.md) - Table explanations
- [`docs/api/EVENTS.md`](../api/EVENTS.md) - Event registrations API

---

## Success Criteria

V31 migration is successful when:

✅ Migration record exists in `migraties` table  
✅ All new columns present in `event_registrations`  
✅ Indexes created successfully  
✅ Data migrated (if active event exists)  
✅ No data loss (all participants have registrations)  
✅ Application starts without errors  
✅ WebSocket steps updates work  
✅ API returns event_registrations data correctly

---

## Support

For issues with V31 migration:
1. Check application logs for migration errors
2. Verify database schema with `\d event_registrations`
3. Check migration status: `SELECT * FROM migraties WHERE version LIKE '%V31%'`
4. Review this document for troubleshooting steps

**Emergency Rollback:** See "Rollback Procedure" section above

---

**Last Updated:** 2025-11-10  
**Author:** Development Team  
**Review Status:** ✅ Documented Post-Implementation