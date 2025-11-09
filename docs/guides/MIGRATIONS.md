# Database Migrations Guide

Complete guide voor het beheren van database schema wijzigingen via migraties.

## Overview

De DKL Email Service gebruikt versioned SQL migraties voor database schema management. Dit zorgt voor:

- **Reproduceerbare deployments** - Schema wijzigingen zijn gecontroleerd en getest
- **Version control** - Alle wijzigingen worden bijgehouden in Git
- **Rollback capability** - Mogelijkheid om terug te draaien bij problemen
- **Team collaboration** - Gestandaardiseerd proces voor schema wijzigingen

---

## Migration System

### Migration Manager

**Location:** `database/migrations.go`

De `MigrationManager` handelt het volledige migratie proces af:

```go
type MigrationManager struct {
    DB              *gorm.DB
    MigratieRepo    repository.MigratieRepository
    MigrationDir    string
}
```

**Features:**
- Automatic schema versioning
- Transaction-based migrations
- Rollback on failure
- Migration history tracking

---

## Migration Files

### Location

```
database/migrations/
├── V01__initial_schema.sql
├── V02__seed_data.sql
├── V03__add_test_data.sql
├── V04__create_chat_tables.sql
├── V05__add_newsletter_tables.sql
├── V06__create_rbac_tables.sql
├── V07__seed_rbac_tables.sql
...
```

### Naming Convention

```
V{VERSION}__{DESCRIPTION}.sql
```

**Rules:**
- `VERSION`: Sequentieel nummer (01, 02, 03, etc.)
- `DESCRIPTION`: Snake_case beschrijving van wijziging
- Extension: `.sql`

**Examples:**
- `V01__initial_schema.sql` - Initial database setup
- `V15__add_user_email_index.sql` - Add index
- `V20__create_events_table.sql` - New table
- `V25__alter_participants_add_steps.sql` - Column addition

---

## Running Migrations

### Automatic (Application Startup)

Migraties worden automatisch uitgevoerd bij het starten van de applicatie:

```go
// main.go
migrationManager := database.NewMigrationManager(db, repoFactory.Migratie)
if err := migrationManager.MigrateDatabase(); err != nil {
    logger.Fatal("Database migratie fout", "error", err)
}
```

**Benefits:**
- Consistent state
- No manual intervention
- Automatic on deploy

---

### Manual Execution

#### Using Go Command

```bash
# Run migrations programmatically
go run database/migrations/run_migrations.go
```

#### Using PowerShell Script

```powershell
# Windows
.\database\apply_all_migrations.ps1

# With specific database
.\database\apply_all_migrations.ps1 -DBName "dklemailservice_test"
```

#### Using Shell Script

```bash
# Linux/macOS
./database/apply_all_migrations.sh

# With environment variables
DB_NAME=dklemailservice_test ./database/apply_all_migrations.sh
```

---

## Creating New Migrations

### Step 1: Create Migration File

```bash
# Determine next version number
ls database/migrations/ | sort | tail -1
# Last file: V30__something.sql

# Create new migration
touch database/migrations/V31__add_new_feature.sql
```

### Step 2: Write Migration SQL

```sql
-- V31__add_new_feature.sql
-- Description: Add email verification timestamp to gebruikers table

-- Add column
ALTER TABLE gebruikers 
ADD COLUMN email_verified_at TIMESTAMP;

-- Add index
CREATE INDEX idx_gebruikers_verified ON gebruikers(email_verified_at)
WHERE email_verified_at IS NOT NULL;

-- Update existing records (optional)
UPDATE gebruikers 
SET email_verified_at = created_at 
WHERE geverifieerd = true;

-- Add comment
COMMENT ON COLUMN gebruikers.email_verified_at IS 'Timestamp when email was verified';
```

### Step 3: Test Migration

```bash
# Test on development database
DB_NAME=dklemailservice_dev go run database/migrations/run_migrations.go

# Verify schema changes
psql -U postgres -d dklemailservice_dev
\d gebruikers
```

### Step 4: Commit Migration

```bash
git add database/migrations/V31__add_new_feature.sql
git commit -m "feat: add email verification timestamp

- Add email_verified_at column to gebruikers
- Create index for verification queries
- Update existing verified users"
```

---

## Migration Types

### Schema Migrations

#### Create Table

```sql
-- V32__create_feedback_table.sql
CREATE TABLE IF NOT EXISTS feedback (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES gebruikers(id) ON DELETE CASCADE,
    rating INTEGER NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_feedback_user ON feedback(user_id);
CREATE INDEX idx_feedback_rating ON feedback(rating);
CREATE INDEX idx_feedback_created ON feedback(created_at);
```

#### Alter Table

```sql
-- V33__add_feedback_category.sql
ALTER TABLE feedback
ADD COLUMN category VARCHAR(100) NOT NULL DEFAULT 'general';

CREATE INDEX idx_feedback_category ON feedback(category);
```

#### Drop Table

```sql
-- V34__drop_legacy_table.sql
DROP TABLE IF EXISTS legacy_data CASCADE;
```

---

### Data Migrations

#### Insert Data

```sql
-- V35__seed_feedback_categories.sql
INSERT INTO feedback_categories (name, description) VALUES
    ('bug_report', 'Bug rapportage'),
    ('feature_request', 'Feature verzoek'),
    ('general', 'Algemene feedback'),
    ('complaint', 'Klacht')
ON CONFLICT (name) DO NOTHING;
```

#### Update Data

```sql
-- V36__normalize_email_addresses.sql
UPDATE gebruikers
SET email = LOWER(TRIM(email))
WHERE email != LOWER(TRIM(email));
```

#### Data Migration with Transformation

```sql
-- V37__migrate_old_status_format.sql
UPDATE events
SET status = CASE 
    WHEN old_status = 'open' THEN 'upcoming'
    WHEN old_status = 'running' THEN 'active'
    WHEN old_status = 'closed' THEN 'completed'
    ELSE 'cancelled'
END;

-- Drop old column
ALTER TABLE events DROP COLUMN IF EXISTS old_status;
```

---

## Advanced Migration Patterns

### Add Column with Backfill

```sql
-- V38__add_total_points_column.sql

-- Step 1: Add column (nullable first)
ALTER TABLE participants
ADD COLUMN total_points INTEGER;

-- Step 2: Backfill existing data
UPDATE participants p
SET total_points = (
    SELECT COALESCE(SUM(a.points), 0)
    FROM user_achievements ua
    JOIN achievements a ON ua.achievement_id = a.id
    WHERE ua.user_id::text = p.id  -- Assuming user_id mapping
);

-- Step 3: Make NOT NULL with default
ALTER TABLE participants
ALTER COLUMN total_points SET DEFAULT 0,
ALTER COLUMN total_points SET NOT NULL;

-- Step 4: Add index
CREATE INDEX idx_participants_points ON participants(total_points);
```

---

### Rename Column Safely

```sql
-- V39__rename_column_with_backward_compatibility.sql

-- Step 1: Add new column
ALTER TABLE fotos
ADD COLUMN title VARCHAR(255);

-- Step 2: Copy data
UPDATE fotos
SET title = titel;

-- Step 3: Create view for backward compatibility
CREATE OR REPLACE VIEW fotos_legacy AS
SELECT 
    id,
    title as titel,  -- Map new to old
    afbeelding_url,
    visible
FROM fotos;

-- Step 4: In future migration, drop old column and view
-- ALTER TABLE fotos DROP COLUMN titel;
-- DROP VIEW fotos_legacy;
```

---

### Complex Foreign Key Changes

```sql
-- V40__update_foreign_keys_with_cascade.sql

-- Step 1: Drop existing foreign key
ALTER TABLE chat_messages
DROP CONSTRAINT IF EXISTS chat_messages_user_id_fkey;

-- Step 2: Add new foreign key with CASCADE
ALTER TABLE chat_messages
ADD CONSTRAINT chat_messages_user_id_fkey
FOREIGN KEY (user_id) 
REFERENCES gebruikers(id) 
ON DELETE CASCADE;
```

---

## Migration Tracking

### Migrations Table

```sql
CREATE TABLE IF NOT EXISTS migraties (
    id SERIAL PRIMARY KEY,
    version VARCHAR(50) UNIQUE NOT NULL,
    beschrijving TEXT,
    toegepast_op TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    success BOOLEAN DEFAULT true
);
```

### Check Migration Status

```sql
-- View applied migrations
SELECT version, beschrijving, toegepast_op, success
FROM migraties
ORDER BY version;

-- Check if specific migration is applied
SELECT EXISTS(
    SELECT 1 FROM migraties 
    WHERE version = 'V31__add_new_feature'
) as is_applied;

-- Find failed migrations
SELECT * FROM migraties WHERE success = false;
```

---

## Rollback Procedures

### Manual Rollback

**⚠️ Warning:** Rollbacks kunnen data verlies veroorzaken. Maak altijd eerst een backup!

#### Step 1: Backup Database

```bash
# Create backup before rollback
pg_dump -U postgres -d dklemailservice > backup_before_rollback.sql
```

#### Step 2: Write Rollback SQL

```sql
-- rollback_V31.sql
-- Rollback for V31__add_new_feature

-- Remove column
ALTER TABLE gebruikers 
DROP COLUMN IF EXISTS email_verified_at;

-- Remove from migrations table
DELETE FROM migraties 
WHERE version = 'V31__add_new_feature';
```

#### Step 3: Execute Rollback

```bash
psql -U postgres -d dklemailservice -f rollback_V31.sql
```

---

### Automated Rollback (Future)

Planning for automated rollback migrations:

```
database/migrations/
├── V31__add_feature.sql      (up migration)
└── V31__add_feature.down.sql (down migration)
```

---

## Testing Migrations

### Test on Development Database

```bash
# 1. Create test database
createdb -U postgres dklemailservice_test

# 2. Run migrations
DB_NAME=dklemailservice_test go run database/migrations/run_migrations.go

# 3. Verify schema
psql -U postgres -d dklemailservice_test
\dt  # List tables
\d table_name  # Describe table

# 4. Run application tests
DB_NAME=dklemailservice_test go test ./tests -v

# 5. Cleanup
dropdb -U postgres dklemailservice_test
```

---

### Migration Test Script

```bash
#!/bin/bash
# test_migration.sh

set -e

DB_NAME="dklemailservice_test"

echo "Creating test database..."
createdb -U postgres $DB_NAME

echo "Running migrations..."
DB_NAME=$DB_NAME go run database/migrations/run_migrations.go

echo "Running tests..."
DB_NAME=$DB_NAME go test ./tests -v

echo "Cleanup..."
dropdb -U postgres $DB_NAME

echo "Migration test passed!"
```

---

## Production Migration Strategy

### Pre-Deployment Checklist

- [ ] Migration tested on development database
- [ ] Migration tested on staging database  
- [ ] All tests passing
- [ ] Backup created
- [ ] Rollback plan documented
- [ ] Team notified of schema changes
- [ ] Maintenance window scheduled (indien nodig)

---

### Zero-Downtime Migrations

#### Approach 1: Add Column (Safe)

```sql
-- V41__add_optional_column.sql

-- Step 1: Add column as nullable
ALTER TABLE users
ADD COLUMN phone_verified BOOLEAN;

-- Step 2: Backfill in background (if needed)
-- Can be done after deployment

-- Step 3: In next migration, make NOT NULL if needed
```

#### Approach 2: Remove Column (Safe)

```sql
-- Migration 1: V42__deprecate_column.sql
-- Mark column as deprecated in code comments
-- Deploy code that doesn't use column

-- Wait for deployment

-- Migration 2: V43__drop_deprecated_column.sql
-- Now safe to drop
ALTER TABLE users DROP COLUMN IF EXISTS deprecated_field;
```

#### Approach 3: Rename Table (Requires Downtime)

```sql
-- V44__rename_table_with_view.sql

-- Step 1: Rename table
ALTER TABLE old_table_name 
RENAME TO new_table_name;

-- Step 2: Create view with old name
CREATE VIEW old_table_name AS
SELECT * FROM new_table_name;

-- Step 3: Update application code to use new name
-- Deploy

-- Step 4: Drop view in future migration
```

---

## Migration Best Practices

### 1. Idempotent Migrations

Migrations should be repeatable:

```sql
-- Good: Safe to run multiple times
CREATE TABLE IF NOT EXISTS new_table (...);
ALTER TABLE users ADD COLUMN IF NOT EXISTS new_field VARCHAR(255);
DROP TABLE IF EXISTS old_table;

-- Bad: Fails on second run
CREATE TABLE new_table (...);
ALTER TABLE users ADD COLUMN new_field VARCHAR(255);
```

### 2. Use Transactions

```sql
-- V45__complex_migration.sql
BEGIN;

-- Multiple related changes
ALTER TABLE users ADD COLUMN status VARCHAR(50);
UPDATE users SET status = 'active' WHERE actief = true;
ALTER TABLE users DROP COLUMN actief;

COMMIT;
-- Auto-rollback on error
```

### 3. Add Indexes Concurrently

```sql
-- Avoid table locks in production
CREATE INDEX CONCURRENTLY idx_users_email ON users(email);

-- Or in migration with note
-- Note: This may lock table, run during maintenance window
CREATE INDEX idx_users_email ON users(email);
```

### 4. Validate Data Before Constraints

```sql
-- V46__add_email_constraint.sql

-- Step 1: Fix invalid data first
UPDATE users 
SET email = LOWER(TRIM(email))
WHERE email ~ '[A-Z]' OR email != TRIM(email);

-- Step 2: Add constraint
ALTER TABLE users
ADD CONSTRAINT users_email_valid 
CHECK (email ~* '^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$');
```

### 5. Document Complex Migrations

```sql
-- V47__complex_schema_refactor.sql
/*
 * Migration: Refactor event registration system
 * 
 * Changes:
 * - Create new event_registrations table
 * - Migrate data from old aanmeldingen table
 * - Add foreign keys to events and participants
 * - Create indexes for performance
 * 
 * Rollback: See rollback_V47.sql
 * 
 * Author: Development Team
 * Date: 2025-01-08
 * Ticket: TASK-123
 */

BEGIN;

-- Create new table
CREATE TABLE event_registrations (
    ...
);

-- Migrate data
INSERT INTO event_registrations (...)
SELECT ... FROM aanmeldingen;

-- Verify migration
DO $$
BEGIN
    IF (SELECT COUNT(*) FROM event_registrations) != 
       (SELECT COUNT(*) FROM aanmeldingen) THEN
        RAISE EXCEPTION 'Data migration count mismatch';
    END IF;
END $$;

COMMIT;
```

---

## Common Migration Patterns

### Add Column with Default

```sql
-- Safe pattern
ALTER TABLE table_name
ADD COLUMN new_column VARCHAR(255) DEFAULT 'default_value';

-- Remove default after backfill (optional)
ALTER TABLE table_name
ALTER COLUMN new_column DROP DEFAULT;
```

---

### Change Column Type

```sql
-- Step 1: Add new column
ALTER TABLE users ADD COLUMN age_new INTEGER;

-- Step 2: Migrate data with validation
UPDATE users 
SET age_new = CAST(age_old AS INTEGER)
WHERE age_old ~ '^\d+$';

-- Step 3: Handle nulls
UPDATE users 
SET age_new = 0 
WHERE age_new IS NULL;

-- Step 4: Drop old, rename new
ALTER TABLE users DROP COLUMN age_old;
ALTER TABLE users RENAME COLUMN age_new TO age;
```

---

### Add Foreign Key to Existing Table

```sql
-- V48__add_created_by_foreign_key.sql

-- Step 1: Ensure referential integrity
UPDATE albums 
SET created_by = NULL 
WHERE created_by NOT IN (SELECT id FROM gebruikers);

-- Step 2: Add foreign key
ALTER TABLE albums
ADD CONSTRAINT albums_created_by_fkey
FOREIGN KEY (created_by) 
REFERENCES gebruikers(id) 
ON DELETE SET NULL;
```

---

### Create Lookup Table

```sql
-- V49__create_status_lookup.sql

-- Step 1: Create lookup table
CREATE TABLE status_types (
    status VARCHAR(50) PRIMARY KEY,
    description TEXT NOT NULL,
    color VARCHAR(7),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Step 2: Populate with existing values
INSERT INTO status_types (status, description)
SELECT DISTINCT status, status
FROM events
WHERE status IS NOT NULL;

-- Step 3: Add foreign key
ALTER TABLE events
ADD CONSTRAINT events_status_fkey
FOREIGN KEY (status)
REFERENCES status_types(status);
```

---

## Troubleshooting

### Migration Failed

**Symptom:** Migration script errors tijdens uitvoering

**Diagnosis:**
```sql
-- Check migration status
SELECT * FROM migraties 
WHERE success = false 
ORDER BY toegepast_op DESC;
```

**Solution:**
```bash
# 1. Review error in logs
# 2. Fix migration SQL
# 3. Mark as not applied
psql -U postgres -d dklemailservice
DELETE FROM migraties WHERE version = 'VXX__failed_migration';

# 4. Rerun migration
go run database/migrations/run_migrations.go
```

---

### Duplicate Migration

**Symptom:** Migration already applied error

**Solution:**
```sql
-- Check if migration is recorded
SELECT * FROM migraties WHERE version = 'VXX__migration';

-- If incorrectly recorded, remove
DELETE FROM migraties WHERE version = 'VXX__migration';
```

---

### Lock Timeout

**Symptom:** Could not obtain lock on table

**Solution:**
```sql
-- Increase lock timeout
SET lock_timeout = '10s';

-- Or run during maintenance window
-- Or use CONCURRENTLY for indexes
```

---

## Production Deployment

### Deployment Process

```bash
# 1. Backup production database
pg_dump -U $DB_USER -h $DB_HOST -d $DB_NAME > production_backup.sql

# 2. Run migrations on staging first
DB_HOST=staging.db.com go run database/migrations/run_migrations.go

# 3. Test staging thoroughly
# Run smoke tests, check logs

# 4. Deploy to production
git push origin main  # Triggers deployment

# 5. Migrations run automatically on startup
# Monitor logs for errors

# 6. Verify schema
psql -U $DB_USER -h $DB_HOST -d $DB_NAME
SELECT * FROM migraties ORDER BY toegepast_op DESC LIMIT 5;
```

---

### Rollback in Production

**Only if migration causes critical issues!**

```bash
# 1. Stop application
# Prevent further migrations

# 2. Restore from backup
psql -U $DB_USER -h $DB_HOST -d $DB_NAME < production_backup.sql

# 3. Redeploy previous version
git revert <migration_commit>
git push origin main

# 4. Monitor recovery
# Check application logs and health endpoints
```

---

## Migration Monitoring

### Health Check

```sql
-- Count pending migrations
SELECT COUNT(*) as pending
FROM (
    SELECT name FROM information_schema.tables 
    WHERE table_schema = 'public'
) t
WHERE t.name NOT IN (
    SELECT version FROM migraties
);
```

### Recent Migrations

```sql
-- Last 10 applied migrations
SELECT 
    version,
    beschrijving,
    toegepast_op,
    success,
    CASE WHEN success THEN '✅' ELSE '❌' END as status
FROM migraties
ORDER BY toegepast_op DESC
LIMIT 10;
```

### Find Long-Running Queries

```sql
-- Check for blocking migrations
SELECT 
    pid,
    now() - pg_stat_activity.query_start AS duration,
    query
FROM pg_stat_activity
WHERE state = 'active'
AND query NOT LIKE '%pg_stat_activity%'
ORDER BY duration DESC;
```

---

## Database Schema Versioning

### Current Schema Version

```bash
# Get latest migration version
psql -U postgres -d dklemailservice -c \
  "SELECT version FROM migraties ORDER BY toegepast_op DESC LIMIT 1;"
```

### Schema Documentation

Maintain `database/SCHEMA_VERSION.md`:

```markdown
# Schema Version History

## Current Version: V50

### V50 (2025-01-08)
- Added email verification tracking
- Created feedback system tables
- Optimized indexes for queries

### V49 (2025-01-07)
- Created status lookup tables
- Normalized event status
...
```

---

## Development Workflow

### Creating Feature Branch Migration

```bash
# 1. Create feature branch
git checkout -b feature/user-verification

# 2. Create migration
touch database/migrations/V51__add_user_verification.sql

# 3. Write migration
# Edit V51__add_user_verification.sql

# 4. Test locally
go run database/migrations/run_migrations.go

# 5. Commit
git add database/migrations/V51__add_user_verification.sql
git commit -m "feat: add user email verification"

# 6. Push and create PR
git push origin feature/user-verification
```

---

### Merging Migrations

**Conflict Resolution:**

Als twee branches beide migraties toevoegen:

```bash
# Branch A: V50__feature_a.sql
# Branch B: V50__feature_b.sql (conflict!)

# Resolution:
# 1. Rename one migration
git mv database/migrations/V50__feature_b.sql \
      database/migrations/V51__feature_b.sql

# 2. Update migration references if any

# 3. Commit resolution
git add database/migrations/V51__feature_b.sql
git commit -m "fix: resolve migration version conflict"
```

---

## Best Practices Summary

### ✅ Do

- Use sequential version numbers
- Write idempotent migrations
- Add descriptive comments
- Test on development first
- Create backups before production
- Use transactions for complex changes
- Add proper indexes
- Validate data before constraints

### ❌ Don't

- Modify existing migration files
- Skip version numbers
- Run untested migrations in production
- Make breaking changes without transition period
- Delete old migrations
- Use hardcoded values (use parameters)
- Ignore foreign key constraints
- Forget to add indexes

---

## Tools & Scripts

### Migration Utilities

**Check pending migrations:**
```bash
# List files not in migraties table
psql -U postgres -d dklemailservice -c "
SELECT file 
FROM unnest(ARRAY[$(ls database/migrations/*.sql | sed 's/.*\///' | sed 's/.sql//' | paste -sd,)]) AS file
WHERE file NOT IN (SELECT version FROM migraties);
"
```

**Generate migration template:**
```bash
#!/bin/bash
# scripts/create_migration.sh

NEXT_VERSION=$(ls database/migrations/ | tail -1 | grep -oP 'V\d+' | grep -oP '\d+' | awk '{print $1+1}')
DESCRIPTION=$1

if [ -z "$DESCRIPTION" ]; then
    echo "Usage: $0 <description>"
    exit 1
fi

FILENAME="database/migrations/V${NEXT_VERSION}__${DESCRIPTION}.sql"

cat > $FILENAME << EOF
-- V${NEXT_VERSION}__${DESCRIPTION}.sql
-- Description: TODO

BEGIN;

-- Add your migration SQL here

COMMIT;
EOF

echo "Created: $FILENAME"
```

**Usage:**
```bash
./scripts/create_migration.sh add_user_preferences
# Creates: database/migrations/V51__add_user_preferences.sql
```

---

## Related Documentation

- [Database Architecture](../architecture/DATABASE.md)
- [Setup Guide](./SETUP.md)
- [Deployment Guide](./DEPLOYMENT.md)
- [Testing Guide](./TESTING.md)

---

**Last Updated:** 2025-01-08