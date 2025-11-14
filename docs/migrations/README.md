# Database Migrations Index

Complete overzicht van alle database migrations voor DKL Email Service.

**Last Updated:** 2025-11-10  
**Current Schema Version:** V34

---

## Migration Overview

| Version | File | Status | Risk | Documentation |
|---------|------|--------|------|---------------|
| V01-V29 | Various | ✅ Applied | 🟢 Low | See [Migration Guide](../guides/MIGRATIONS.md) |
| **V30** | [`V30__dual_registration_system_with_rbac.sql`](../../database/migrations/V30__dual_registration_system_with_rbac.sql) | ✅ Applied | 🟡 Medium | **[Complete Docs](../V30_RBAC_INTEGRATION.md)** |
| **V31** | [`V31__complete_v28_participant_refactor.sql`](../../database/migrations/V31__complete_v28_participant_refactor.sql) | ✅ Applied | 🔴 High | **[V31 Docs](V31_PARTICIPANT_REFACTOR.md)** |
| **V32** | [`V32__final_schema_alignment_fixes.sql`](../../database/migrations/V32__final_schema_alignment_fixes.sql) | ✅ Applied | 🟡 Medium | **[V32 Docs](V32_SCHEMA_CHANGES.md)** |
| **V33** | [`V33__create_auto_responses_table.sql`](../../database/migrations/V33__create_auto_responses_table.sql) | ✅ Applied | 🟢 Low | **[API Docs](../api/AUTO_RESPONSES.md)** |
| **V34** | [`V34__remove_legacy_participant_columns.sql`](../../database/migrations/V34__remove_legacy_participant_columns.sql) | ✅ Applied | 🔴🔴 Critical | **[V34 Docs](V34_BREAKING_CHANGES.md)** |

---

## Recent Migrations Deep Dive

### V30: Dual Registration System with RBAC ⭐

**Date:** 2025-11-10  
**Impact:** Major Feature Addition  
**Risk:** 🟡 Medium

**What it does:**
- Implements dual account system (temporary vs full)
- Full RBAC integration for participants
- 6 new columns in `participants` table
- New permissions and roles
- Automatic role assignment triggers

**Key Features:**
- ✅ Temporary accounts (event-only, no app access)
- ✅ Full accounts (app access, gamification, multi-year)
- ✅ Automatic `participant_user` role assignment
- ✅ Permission-based app access control

**Documentation:**
- **[V30 RBAC Integration](../V30_RBAC_INTEGRATION.md)** - 1280+ lines, complete technical docs
- **[V30 Implementation Summary](../V30_IMPLEMENTATION_SUMMARY.md)** - Status overview
- **[V30 Quick Reference](../V30_RBAC_QUICK_REFERENCE.md)** - Quick commands
- **[V30 Frontend Guide](../frontend/V30_FRONTEND_IMPLEMENTATION_GUIDE.md)** - Frontend integration

**Database Changes:**
```sql
-- New columns in participants
account_type TEXT NOT NULL DEFAULT 'temporary'
registration_year INTEGER
wachtwoord_hash TEXT
has_app_access BOOLEAN NOT NULL DEFAULT FALSE
gebruiker_id UUID REFERENCES gebruikers(id)
upgraded_to_gebruiker_id UUID REFERENCES gebruikers(id)
upgraded_at TIMESTAMPTZ
```

**New Tables:**
- `participant_upgrades` - Audit trail for account upgrades
- `participant_rbac_audit` - RBAC action logging

**Impact:**
- ✅ All existing participants become 'temporary'
- ✅ Automatic migration path to 'full' accounts
- ✅ Backward compatible

---

### V31: Complete V28 Participant Refactor 🔄

**Date:** 2025-11-10  
**Impact:** Architecture Fix  
**Risk:** 🔴 High (Data Migration)

**What it does:**
- Completes failed V28 participant refactor
- Separates person data from event data
- Moves event columns from `participants` to `event_registrations`
- Migrates existing data to new structure

**Problem it solves:**
- V28 attempted to separate concerns but failed
- Participants table had mixed person + event data
- Couldn't track multiple events per participant

**Solution:**
```
participants table = PERSON DATA ONLY
├── naam, email, telefoon (person info)
└── account_type, has_app_access (V30 fields)

event_registrations table = EVENT PARTICIPATION
├── participant_role_name (Deelnemer/Begeleider/Vrijwilliger)
├── distance_route (2.5KM/6KM/10KM/15KM)
├── steps (event-specific tracking!)
└── ondersteuning, bijzonderheden, status
```

**Documentation:**
- **[V31 Participant Refactor](V31_PARTICIPANT_REFACTOR.md)** - Complete migration guide

**Database Changes:**
- ✅ Added 7 columns to `event_registrations`
- ✅ Migrated data from `participants` to `event_registrations`
- ✅ Added performance indexes
- ⚠️ Kept legacy columns in `participants` (removed in V34)

**Impact:**
- ✅ One participant can now have multiple event registrations
- ✅ Historical event data preserved
- ⚠️ Requires code updates to query event_registrations

---

### V32: Final Schema Alignment Fixes 🔧

**Date:** 2025-11-10  
**Impact:** Schema Cleanup  
**Risk:** 🟡 Medium (Breaking Changes)

**What it does:**
- Aligns database schema with documentation
- Adds missing columns to various tables
- **Renames columns** in notification lookup tables
- Adds display_order for UI sorting
- Adds RBAC permissions for event_registrations

**Breaking Changes:**
```sql
-- notification_types
name → type (column rename)

-- notification_priority_types
name → priority (column rename)
```

**Documentation:**
- **[V32 Schema Changes](V32_SCHEMA_CHANGES.md)** - Breaking changes guide

**Database Changes:**
- ✅ Added timestamps to `participant_roles`
- ✅ Added distance_km, description to `distances`
- ⚠️ Renamed PK columns in notification lookup tables
- ✅ Added display_order to all lookup tables
- ✅ Added event_registrations permissions

**Impact:**
- ⚠️ Requires code updates for notification type queries
- ⚠️ Model structs need updating
- ✅ Better UI ordering with display_order

---

### V33: Create Auto Responses Table 📧

**Date:** 2025-11-10  
**Impact:** New Feature  
**Risk:** 🟢 Low

**What it does:**
- Creates `auto_responses` table for email auto-reply management
- Enables vacation responders and automated acknowledgments
- Time-based activation (start_date, end_date)

**Use Cases:**
- Out-of-office messages
- Event registration confirmations
- Maintenance notifications
- Business hours responders

**Documentation:**
- **[Auto Response API](../api/AUTO_RESPONSES.md)** - Complete API reference

**Database Changes:**
```sql
CREATE TABLE auto_responses (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,
    subject VARCHAR(255),
    message TEXT,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

**API Endpoints:**
- GET `/api/mail/autoresponse` - List all
- GET `/api/mail/autoresponse/:id` - Get one
- POST `/api/mail/autoresponse` - Create
- PUT `/api/mail/autoresponse/:id` - Update
- DELETE `/api/mail/autoresponse/:id` - Delete

**Impact:**
- ✅ New feature, no breaking changes
- ✅ Fully functional immediately

---

### V34: Remove Legacy Participant Columns 🗑️

**Date:** 2025-11-10  
**Impact:** Schema Cleanup  
**Risk:** 🔴🔴 CRITICAL (Destructive!)

**What it does:**
- **Permanently removes 14 legacy columns** from `participants` table
- Completes V28/V31 refactor by deleting old event-specific columns
- Drops legacy indexes
- Uses CASCADE to force removal

**⚠️ THIS IS DESTRUCTIVE:**
```sql
ALTER TABLE participants DROP COLUMN IF EXISTS
    rol, afstand, ondersteuning, bijzonderheden,
    email_verzonden, email_verzonden_op,
    behandeld_door, behandeld_op, notities,
    steps, antwoorden_count,
    participant_role_name, distance_route, status
CASCADE;
```

**Documentation:**
- **[V34 Breaking Changes](V34_BREAKING_CHANGES.md)** - ⚠️ MUST READ before deployment

**Why?**
- V28/V31 moved data to `event_registrations`
- Kept old columns for transition period
- Now safe to remove after verification

**Database Changes:**
- 🔴 Removes 14 columns permanently
- 🔴 Drops 3 legacy indexes
- ✅ Cleans up participants table

**Impact:**
- 🔴 **BREAKING:** Code using old columns will fail
- 🔴 **NO ROLLBACK** without database restore
- ✅ Cleaner schema after removal

**Pre-Flight Checklist:**
- [ ] Database backup created
- [ ] V31 data migration verified
- [ ] All code updated to use event_registrations
- [ ] Tested in staging
- [ ] Team notified

---

## Migration Flow Diagram

```
V28 (Failed)
   │
   ├─ Attempted: Separate person/event data
   └─ Result: Failed, left schema in inconsistent state
      │
      ▼
V30 (Success) - Added dual account system
   │
   ├─ Added: account_type, has_app_access, etc.
   ├─ New: RBAC integration for participants
   └─ Status: Fully documented ✅
      │
      ▼
V31 (Fix) - Completed V28 refactor
   │
   ├─ Fixed: V28 failure
   ├─ Added: Missing columns to event_registrations
   ├─ Migrated: Data from participants to event_registrations
   ├─ Kept: Legacy columns for transition
   └─ Status: NOW DOCUMENTED ✅
      │
      ▼
V32 (Alignment) - Schema cleanup
   │
   ├─ Fixed: Column naming inconsistencies
   ├─ Added: Missing columns (timestamps, display_order)
   ├─ Renamed: notification_types columns (BREAKING)
   └─ Status: NOW DOCUMENTED ✅
      │
      ▼
V33 (Feature) - Auto responses
   │
   ├─ Added: auto_responses table
   ├─ New: Email auto-reply functionality
   └─ Status: NOW DOCUMENTED ✅
      │
      ▼
V34 (Cleanup) - Remove legacy columns
   │
   ├─ Removed: 14 legacy columns from participants
   ├─ Dropped: Legacy indexes
   ├─ Result: Clean schema matching architecture
   └─ Status: NOW DOCUMENTED ✅ (with critical warnings)
```

---

## Quick Reference

### Check Applied Migrations

```sql
-- List all applied migrations
SELECT version, beschrijving, toegepast_op, success
FROM migraties
ORDER BY version;

-- Check specific migration
SELECT * FROM migraties WHERE version = 'V31__complete_v28_participant_refactor';

-- Count migrations
SELECT COUNT(*) FROM migraties;
```

### Verify Schema State

```sql
-- Participant schema (post-V34)
\d participants
-- Should NOT have: rol, afstand, steps, etc.
-- Should HAVE: account_type, has_app_access, gebruiker_id

-- Event registrations schema (post-V31)
\d event_registrations
-- Should HAVE: steps, participant_role_name, distance_route

-- Notification types (post-V32)
\d notification_types
-- Should have 'type' column (not 'name')
```

### Database Health Check

```bash
# Check migration status
docker-compose exec app psql -U postgres -d dklemailservice -c \
  "SELECT COUNT(*) as total, 
          SUM(CASE WHEN success THEN 1 ELSE 0 END) as successful,
          SUM(CASE WHEN NOT success THEN 1 ELSE 0 END) as failed
   FROM migraties;"

# Verify V30-V34 applied
docker-compose exec app psql -U postgres -d dklemailservice -c \
  "SELECT version, success FROM migraties 
   WHERE version LIKE 'V3%' ORDER BY version;"
```

---

## Migration Dependencies

```
V30 ────┐
        │
        ├──► V31 (requires participants table from V30)
        │     │
        │     └──► V32 (enhances event_registrations from V31)
        │           │
        └───────────┴──► V33 (independent - auto responses)
                    │
                    └──► V34 (requires V31 data migration complete)
```

**Key Dependencies:**
- V31 requires V30 (uses participants table)
- V34 requires V31 (must migrate data before deleting columns)
- V32 enhances V31 (adds columns to event_registrations)
- V33 is independent (can run anytime)

---

## Documentation Structure

```
docs/
├── migrations/
│   ├── README.md                          ← YOU ARE HERE
│   ├── V31_PARTICIPANT_REFACTOR.md        ← V31 deep dive
│   ├── V32_SCHEMA_CHANGES.md              ← V32 breaking changes
│   └── V34_BREAKING_CHANGES.md            ← V34 critical warnings
│
├── V30_RBAC_INTEGRATION.md                ← V30 complete docs
├── V30_*.md                                ← V30 related docs (8 files)
│
├── api/
│   └── AUTO_RESPONSES.md                  ← V33 API docs
│
├── architecture/
│   └── DATABASE.md                        ← All tables documented
│
└── guides/
    └── MIGRATIONS.md                      ← General migration guide
```

---

## Risk Classification

### 🟢 Low Risk Migrations
- **V33**: New table, no impact on existing data
- Simple additive changes

### 🟡 Medium Risk Migrations
- **V30**: Major feature but well-planned
- **V32**: Breaking changes but manageable

### 🔴 High Risk Migrations
- **V31**: Data migration between tables

### 🔴🔴 Critical Risk Migrations
- **V34**: Permanent data deletion

---

## When to Deploy

### V30 Deployment
✅ **Ready:** Fully tested and documented  
📋 **Checklist:**
- [ ] Review [V30 RBAC Integration](../V30_RBAC_INTEGRATION.md)
- [ ] Test in staging
- [ ] Deploy with frontend changes

### V31 Deployment
✅ **Ready:** Safe to deploy (already applied)  
⚠️ **Requirements:**
- Must have active event in database
- Code must use event_registrations for event data

### V32 Deployment
✅ **Ready:** Safe to deploy (already applied)  
⚠️ **Requirements:**
- Code must use 'type' not 'name' for notification lookups
- Models updated to match renamed columns

### V33 Deployment
✅ **Ready:** Safe to deploy (already applied)  
✅ **Requirements:**
- None - new feature, no dependencies

### V34 Deployment
⚠️ **CAUTION:** Only after V31 verified!  
🔴 **CRITICAL Requirements:**
- [ ] Database backup created
- [ ] V31 data migration verified successful
- [ ] All code references to deleted columns removed
- [ ] Tested in staging environment
- [ ] Team notified of deployment

---

## Rollback Information

| Migration | Rollback Difficulty | Data Loss Risk | Documented Procedure |
|-----------|---------------------|----------------|----------------------|
| V30 | Medium | Low | [V30 Docs](../V30_RBAC_INTEGRATION.md) |
| V31 | High | Medium | [V31 Docs](V31_PARTICIPANT_REFACTOR.md#rollback-procedure) |
| V32 | Medium | Low | [V32 Docs](V32_SCHEMA_CHANGES.md#rollback-procedure) |
| V33 | Easy | None | Simple DROP TABLE |
| V34 | **Impossible** | **High** | [V34 Docs](V34_BREAKING_CHANGES.md#rollback-procedure) (restore from backup only) |

⚠️ **V34 Note:** Cannot rollback without database restore - columns are permanently deleted!

---

## Common Issues

### Migration Won't Apply

**Check:**
```sql
-- Is migration already applied?
SELECT * FROM migraties WHERE version = 'VXX__migration_name';

-- Are there conflicting schema changes?
\d table_name
```

### Data Migration Failed

**For V31:**
```sql
-- Check if active event exists
SELECT * FROM events WHERE status = 'active';

-- Manually set active event if needed
UPDATE events SET status = 'active' WHERE id = '<event_id>';
```

### Code Breaking After Migration

**For V32:**
```bash
# Find references to old column names
grep -r "\.name" . | grep "notification_types"
# Update to .type

grep -r "\.name" . | grep "notification_priority_types"
# Update to .priority
```

**For V34:**
```bash
# Find references to deleted columns
grep -r "\.steps" . | grep -v "event_registrations"
grep -r "\.rol" .
grep -r "\.afstand" .
```

---

## Getting Help

### Documentation
1. Start with [Migration Guide](../guides/MIGRATIONS.md) - General procedures
2. Check specific migration docs (V31, V32, V34)
3. Review [Database Architecture](../architecture/DATABASE.md) - Schema reference

### Commands
```bash
# View migration logs
docker-compose logs app | grep -i migration

# Check database state
docker-compose exec db psql -U postgres -d dklemailservice

# Run health check
curl http://localhost:8080/api/health
```

### Support
- **Database Issues:** Check migration-specific docs above
- **Code Issues:** Review code impact sections
- **Deployment:** See [Deployment Guide](../guides/DEPLOYMENT.md)

---

## Next Migration: V35+

When creating new migrations:

1. **Read:** [Migration Guide](../guides/MIGRATIONS.md)
2. **Follow:** Naming convention `V{NUMBER}__{description}.sql`
3. **Test:** In development environment first
4. **Document:** Create migration doc if complex or breaking
5. **Review:** Get team approval for schema changes

**Template:**
```sql
-- VXX__description.sql
-- Description: What this migration does
-- Risk Level: Low/Medium/High/Critical
-- Breaking Changes: Yes/No

BEGIN;

-- Your migration SQL here

COMMIT;
```

---

**Migration System Status:** ✅ HEALTHY  
**Documentation Status:** ✅ COMPLETE  
**Production Readiness:** ✅ READY