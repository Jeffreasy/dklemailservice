# V32: Final Schema Alignment Fixes

**Migration File:** `database/migrations/V32__final_schema_alignment_fixes.sql`  
**Date:** 2025-11-10  
**Status:** ✅ Applied  
**Risk Level:** 🟡 MEDIUM (Breaking Changes via Column Renames)

---

## Overview

V32 performs "final alignment" fixes to ensure the database schema matches the documented architecture. This migration:

1. Adds missing timestamp columns
2. Adds missing distance/description columns
3. **Renames columns in notification tables** (BREAKING)
4. Adds display_order to lookup tables
5. Adds new RBAC permissions for event_registrations

---

## What V32 Does

### Part 1: Add Timestamps to `participant_roles`

```sql
ALTER TABLE participant_roles
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP;
```

**Why?** 
- Audit trail for when roles were created/modified
- Consistency with other tables
- Tracking for future role management features

**Impact:** Low - just adds columns, no breaking changes

---

### Part 2: Add Columns to `distances`

```sql
ALTER TABLE distances
    ADD COLUMN IF NOT EXISTS distance_km NUMERIC(10,2),
    ADD COLUMN IF NOT EXISTS description TEXT;

-- Populate data
UPDATE distances
SET distance_km = CASE
    WHEN route ~ '^\d+' THEN (regexp_match(route, '^\d+'))[1]::NUMERIC
    ELSE NULL
END
WHERE distance_km IS NULL;

UPDATE distances SET description = 'Route ' || route WHERE description IS NULL;
```

**Why?**
- `distance_km`: Standardized numeric distance for calculations
- `description`: User-friendly text for UI display
- Better than parsing `route` string every time

**Example:**
| route | distance_km | description |
|-------|-------------|-------------|
| "5km" | 5.00 | "Route 5km" |
| "10km" | 10.00 | "Route 10km" |

**Impact:** Low - adds helpful data, doesn't break anything

---

### Part 3: Rename Notification Lookup Table Columns ⚠️

**THIS IS A BREAKING CHANGE!**

#### 3a. Rename `notification_types.name` → `notification_types.type`

```sql
-- Before
CREATE TABLE notification_types (
    name VARCHAR(50) PRIMARY KEY,  -- ❌ Inconsistent naming
    ...
);

-- After V32
CREATE TABLE notification_types (
    type VARCHAR(50) PRIMARY KEY,  -- ✅ Matches foreign key column
    ...
);
```

**Why?**
- `notifications` table has FK column called `type`, not `name`
- Consistency: PK column should match FK column name
- Clearer semantics: "type" is more accurate than "name"

**Breaking Change:**
```sql
-- BEFORE V32 (queries that will break):
SELECT * FROM notification_types WHERE name = 'email';

-- AFTER V32 (correct queries):
SELECT * FROM notification_types WHERE type = 'email';
```

#### 3b. Rename `notification_priority_types.name` → `notification_priority_types.priority`

```sql
-- Before
CREATE TABLE notification_priority_types (
    name VARCHAR(50) PRIMARY KEY,  -- ❌ Inconsistent
    ...
);

-- After V32
CREATE TABLE notification_priority_types (
    priority VARCHAR(50) PRIMARY KEY,  -- ✅ Matches FK
    ...
);
```

**Why?** Same reason - consistency with FK column naming.

**Migration Process:**
1. Drop existing foreign keys
2. Rename column
3. Re-add foreign keys with updated column references

```sql
-- Drop FKs
ALTER TABLE notifications DROP CONSTRAINT IF EXISTS notifications_type_key_fkey;
ALTER TABLE notifications DROP CONSTRAINT IF EXISTS notifications_priority_key_fkey;

-- Rename columns
ALTER TABLE notification_types RENAME COLUMN name TO type;
ALTER TABLE notification_priority_types RENAME COLUMN name TO priority;

-- Re-add FKs
ALTER TABLE notifications 
    ADD CONSTRAINT notifications_type_fkey 
    FOREIGN KEY (type) REFERENCES notification_types(type);

ALTER TABLE notifications 
    ADD CONSTRAINT notifications_priority_fkey 
    FOREIGN KEY (priority) REFERENCES notification_priority_types(priority);
```

---

### Part 4: Add `display_order` to Lookup Tables

```sql
ALTER TABLE contact_status_types ADD COLUMN IF NOT EXISTS display_order INTEGER;
ALTER TABLE registration_status_types ADD COLUMN IF NOT EXISTS display_order INTEGER;
ALTER TABLE email_status_types ADD COLUMN IF NOT EXISTS display_order INTEGER;
ALTER TABLE event_status_types ADD COLUMN IF NOT EXISTS display_order INTEGER;
ALTER TABLE chat_channel_types ADD COLUMN IF NOT EXISTS display_order INTEGER;
ALTER TABLE notification_types ADD COLUMN IF NOT EXISTS display_order INTEGER;
ALTER TABLE notification_priority_types ADD COLUMN IF NOT EXISTS display_order INTEGER;
```

**Why?**
- Control sort order in dropdowns/UI
- Logical ordering (e.g., "nieuw" before "afgehandeld")
- Better UX without hardcoded sorting

**Example Values:**
```sql
UPDATE contact_status_types SET display_order = 1 WHERE status = 'nieuw';
UPDATE contact_status_types SET display_order = 2 WHERE status = 'in_behandeling';
UPDATE contact_status_types SET display_order = 3 WHERE status = 'afgehandeld';
```

**Usage in API:**
```sql
SELECT * FROM contact_status_types ORDER BY display_order;
```

**Impact:** Low - just adds metadata for UI ordering

---

### Part 5: Add `event_registrations` Permissions

```sql
INSERT INTO permissions (resource, action, description) VALUES
    ('event_registrations', 'read', 'View event registrations'),
    ('event_registrations', 'write', 'Create and update event registrations'),
    ('event_registrations', 'delete', 'Delete event registrations');

-- Assign to admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r CROSS JOIN permissions p
WHERE r.name = 'admin' AND p.resource = 'event_registrations';
```

**Why?**
- `event_registrations` is a major table after V28/V31 refactor
- Needs proper RBAC permissions
- Consistency with other resources

**Impact:** Low - adds permissions, doesn't restrict existing access

---

## Breaking Changes Summary

### ⚠️ CRITICAL: Code Changes Required

#### 1. Notification Type Queries

**BEFORE V32:**
```go
// ❌ WILL FAIL after V32
db.Table("notification_types").
    Where("name = ?", "email").
    First(&notifType)
```

**AFTER V32:**
```go
// ✅ Correct
db.Table("notification_types").
    Where("type = ?", "email").
    First(&notifType)
```

#### 2. Model Struct Changes

**File:** `models/notification_type.go`

**BEFORE V32:**
```go
type NotificationType struct {
    Name        string `gorm:"column:name;primaryKey"` // ❌ Wrong column
    Description string
}
```

**AFTER V32:**
```go
type NotificationType struct {
    Type        string `gorm:"column:type;primaryKey"` // ✅ Correct
    Description string
    DisplayOrder int   `gorm:"column:display_order"`   // New field
}
```

**File:** `models/notification_priority_type.go`

**BEFORE V32:**
```go
type NotificationPriorityType struct {
    Name        string `gorm:"column:name;primaryKey"` // ❌ Wrong
    Description string
}
```

**AFTER V32:**
```go
type NotificationPriorityType struct {
    Priority    string `gorm:"column:priority;primaryKey"` // ✅ Correct
    Description string
    DisplayOrder int   `gorm:"column:display_order"`       // New field
}
```

#### 3. Repository Queries

**File:** `repository/notification_repository.go`

**BEFORE V32:**
```go
func (r *NotificationRepo) GetTypeByName(name string) (*NotificationType, error) {
    var notifType NotificationType
    err := r.db.Where("name = ?", name).First(&notifType).Error // ❌ Wrong column
    return &notifType, err
}
```

**AFTER V32:**
```go
func (r *NotificationRepo) GetTypeByType(notifType string) (*NotificationType, error) {
    var typeRecord NotificationType
    err := r.db.Where("type = ?", notifType).First(&typeRecord).Error // ✅ Correct
    return &typeRecord, err
}
```

#### 4. SQL Queries in Code

**Search for:**
```bash
grep -r "notification_types.*name" .
grep -r "notification_priority_types.*name" .
```

**Replace:**
- `notification_types.name` → `notification_types.type`
- `notification_priority_types.name` → `notification_priority_types.priority`

---

## Migration Execution

### Development

```bash
# 1. Backup
pg_dump -U postgres -d dklemailservice_dev > backup_dev_v32.sql

# 2. Apply migration
docker-compose restart app

# 3. Verify column renames
psql -U postgres -d dklemailservice_dev
\d notification_types  # Should show 'type' column, not 'name'
\d notification_priority_types  # Should show 'priority', not 'name'

# 4. Test queries
SELECT type FROM notification_types;  -- Should work
SELECT name FROM notification_types;  -- Should fail (column doesn't exist)
```

### Staging

```bash
# 1. Update code FIRST (deploy before schema change)
git checkout feature/v32-column-rename-fixes
# Update models to use 'type' instead of 'name'
git push origin staging

# 2. Then run migration
# Migration runs automatically on deploy

# 3. Verify
curl -X GET http://staging-api/api/notifications/types
```

### Production

```bash
# IMPORTANT: Code must be deployed BEFORE or WITH migration
# If code is not updated, API will break!

# 1. Backup
pg_dump -U $DB_USER -h prod-db -d dklemailservice > backup_prod_v32.sql

# 2. Deploy (code + migration together)
git push origin main

# 3. Monitor logs
tail -f /var/log/dkl/app.log | grep -i "notification"

# 4. Test notifications endpoint
curl -X GET https://api.dekoninklijkeloop.nl/api/notifications
```

---

## Verification Steps

```sql
-- 1. Migration applied?
SELECT * FROM migraties WHERE version = 'V32__final_schema_alignment_fixes';

-- 2. participant_roles has timestamps?
\d participant_roles
-- Should show: created_at, updated_at

-- 3. distances has new columns?
\d distances
-- Should show: distance_km, description

-- 4. Notification tables renamed correctly?
\d notification_types
-- Should show: type (not name), display_order

\d notification_priority_types
-- Should show: priority (not name), display_order

-- 5. Existing data preserved?
SELECT type, description FROM notification_types;
-- Should return existing notification types

-- 6. Foreign keys working?
SELECT 
    n.id,
    n.type,
    nt.description
FROM notifications n
JOIN notification_types nt ON n.type = nt.type
LIMIT 5;
-- Should work without errors

-- 7. New permissions added?
SELECT * FROM permissions WHERE resource = 'event_registrations';
-- Should return 3 permissions (read, write, delete)
```

---

## Rollback Procedure

### SQL Rollback

```sql
-- rollback_V32.sql
BEGIN;

-- 1. Revert notification_types column rename
ALTER TABLE notifications DROP CONSTRAINT IF EXISTS notifications_type_fkey;
ALTER TABLE notification_types RENAME COLUMN type TO name;
ALTER TABLE notifications 
    ADD CONSTRAINT notifications_type_name_fkey
    FOREIGN KEY (type) REFERENCES notification_types(name);

-- 2. Revert notification_priority_types rename
ALTER TABLE notifications DROP CONSTRAINT IF EXISTS notifications_priority_fkey;
ALTER TABLE notification_priority_types RENAME COLUMN priority TO name;
ALTER TABLE notifications 
    ADD CONSTRAINT notifications_priority_name_fkey
    FOREIGN KEY (priority) REFERENCES notification_priority_types(name);

-- 3. Remove added columns
ALTER TABLE participant_roles 
    DROP COLUMN IF EXISTS created_at,
    DROP COLUMN IF EXISTS updated_at;

ALTER TABLE distances 
    DROP COLUMN IF EXISTS distance_km,
    DROP COLUMN IF EXISTS description;

ALTER TABLE contact_status_types DROP COLUMN IF EXISTS display_order;
ALTER TABLE registration_status_types DROP COLUMN IF EXISTS display_order;
-- ... (drop display_order from all lookup tables)

-- 4. Remove event_registrations permissions
DELETE FROM role_permissions 
WHERE permission_id IN (
    SELECT id FROM permissions WHERE resource = 'event_registrations'
);
DELETE FROM permissions WHERE resource = 'event_registrations';

-- 5. Remove migration record
DELETE FROM migraties WHERE version = 'V32__final_schema_alignment_fixes';

COMMIT;
```

### Execute Rollback

```bash
# 1. Backup current state
pg_dump -U postgres -d dklemailservice > backup_current_v32_rollback.sql

# 2. Stop application
docker-compose down

# 3. Execute rollback
psql -U postgres -d dklemailservice -f rollback_V32.sql

# 4. Revert code changes
git revert <v32-commit-hash>
git push

# 5. Restart application
docker-compose up -d
```

---

## Troubleshooting

### Error: Column "name" does not exist

**Symptom:**
```
ERROR: column "name" does not exist
LINE 1: SELECT * FROM notification_types WHERE name = 'email'
```

**Cause:** Code not updated to use new column names

**Solution:**
```go
// Update all references:
// notification_types.name → notification_types.type
// notification_priority_types.name → notification_priority_types.priority
```

### Error: Ambiguous column name

**Symptom:**
```
ERROR: column reference "type" is ambiguous
```

**Cause:** JOIN query without table alias

**Solution:**
```sql
-- BEFORE (ambiguous)
SELECT type FROM notifications n
JOIN notification_types nt USING (type);

-- AFTER (explicit)
SELECT n.type FROM notifications n
JOIN notification_types nt ON n.type = nt.type;
```

---

## Testing Checklist

- [ ] Migration applies successfully
- [ ] All new columns present
- [ ] Column renames successful
- [ ] Foreign keys work
- [ ] Existing data preserved
- [ ] API endpoints return correct data
- [ ] No "column does not exist" errors in logs
- [ ] Notification system works
- [ ] Lookup tables sortable by display_order
- [ ] RBAC permissions work for event_registrations

---

## Related Migrations

| Migration | Relation | Description |
|-----------|----------|-------------|
| **V30** | Foundation | RBAC system that V32 extends |
| **V31** | Prerequisite | Created event_registrations table |
| **V34** | Successor | Removes legacy columns |

---

## Documentation Links

- [Database Architecture](../architecture/DATABASE.md) - Schema documentation
- [V31 Participant Refactor](V31_PARTICIPANT_REFACTOR.md) - Event registrations context
- [V30 RBAC Integration](../V30_RBAC_INTEGRATION.md) - RBAC permissions

---

## Success Criteria

V32 migration is successful when:

✅ Migration record exists  
✅ `participant_roles` has timestamps  
✅ `distances` has distance_km and description  
✅ `notification_types.type` column exists (not `name`)  
✅ `notification_priority_types.priority` column exists (not `name`)  
✅ All lookup tables have `display_order`  
✅ `event_registrations` permissions added to RBAC  
✅ Foreign keys intact and working  
✅ No "column not found" errors  
✅ API endpoints respond correctly  

---

**Last Updated:** 2025-11-10  
**Author:** Development Team  
**Review Status:** ✅ Documented Post-Implementation  
**Risk Assessment:** 🟡 MEDIUM - Breaking changes via column renames, but manageable with code updates