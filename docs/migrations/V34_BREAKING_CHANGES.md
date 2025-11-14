# V34: Remove Legacy Participant Columns

**Migration File:** `database/migrations/V34__remove_legacy_participant_columns.sql`  
**Date:** 2025-11-10  
**Status:** ✅ Applied  
**Risk Level:** 🔴🔴 CRITICAL (DESTRUCTIVE - Permanent Data Removal)

---

## ⚠️ CRITICAL WARNING

**THIS IS A DESTRUCTIVE MIGRATION!**

V34 **permanently deletes** multiple columns from the `participants` table. These columns were moved to `event_registrations` in V28/V31, but the actual column deletion happens here.

### What Gets Deleted

```sql
ALTER TABLE participants 
DROP COLUMN IF EXISTS rol CASCADE,
DROP COLUMN IF EXISTS afstand CASCADE,
DROP COLUMN IF EXISTS ondersteuning CASCADE,
DROP COLUMN IF EXISTS bijzonderheden CASCADE,
DROP COLUMN IF EXISTS email_verzonden CASCADE,
DROP COLUMN IF EXISTS email_verzonden_op CASCADE,
DROP COLUMN IF EXISTS behandeld_door CASCADE,
DROP COLUMN IF EXISTS behandeld_op CASCADE,
DROP COLUMN IF EXISTS notities CASCADE,
DROP COLUMN IF EXISTS steps CASCADE,
DROP COLUMN IF EXISTS antwoorden_count CASCADE,
DROP COLUMN IF EXISTS participant_role_name CASCADE,
DROP COLUMN IF EXISTS distance_route CASCADE,
DROP COLUMN IF EXISTS status CASCADE;
```

**14 columns permanently removed** with `CASCADE` (force delete dependencies).

---

## Why V34 is Needed

### Background

After V28/V31 refactored the data architecture, these columns became redundant:

**Old Architecture (Pre-V28):**
```sql
participants = PERSON + EVENT DATA (mixed)
├── naam, email          ← Person data ✅
├── rol, afstand         ← Event data ❌ (wrong place!)
├── steps, ondersteuning ← Event data ❌ (wrong place!)
└── status, bijzonderheden ← Event data ❌ (wrong place!)
```

**New Architecture (Post-V31):**
```sql
participants = PERSON DATA ONLY
├── id, naam, email, telefoon
├── account_type, has_app_access (V30)
└── gebruiker_id, upgraded_at

event_registrations = EVENT DATA
├── participant_role_name ← Moved here ✅
├── distance_route        ← Moved here ✅
├── steps                 ← Moved here ✅
├── ondersteuning         ← Moved here ✅
└── bijzonderheden        ← Moved here ✅
```

### Why Keep Them Until V34?

**Transition Period:**
- V28/V31 copied data but kept old columns
- Allowed for gradual code migration
- Provided rollback safety
- Enabled verification of data migration

**Now (V34):**
- Data migration complete and verified
- All code updated to use event_registrations
- Safe to remove legacy columns
- Clean up database schema

---

## What V34 Does

### Step 1: Drop Legacy Indexes

```sql
DROP INDEX IF EXISTS idx_aanmeldingen_rol CASCADE;
DROP INDEX IF EXISTS idx_aanmeldingen_afstand CASCADE;
DROP INDEX IF EXISTS idx_aanmeldingen_status CASCADE;
```

**Why First?** Indexes must be dropped before columns, otherwise column drop will fail.

### Step 2: Remove Legacy Columns

```sql
ALTER TABLE participants 
DROP COLUMN IF EXISTS rol CASCADE,
DROP COLUMN IF EXISTS afstand CASCADE,
-- ... (14 columns total)
```

**CASCADE Effect:**
- Automatically drops any views using these columns
- Drops any triggers referencing these columns
- Drops any constraints involving these columns

### Step 3: Update Table Comment

```sql
COMMENT ON TABLE participants IS 
'V34: Cleaned up - removed legacy event-specific columns. 
Person data only + V30 account type fields.';
```

---

## Before Running V34

### Pre-Migration Checklist

**MANDATORY BEFORE V34:**

- [ ] **Backup database** - This is NON-NEGOTIABLE!
  ```bash
  pg_dump -U postgres -d dklemailservice > backup_before_v34.sql
  ```

- [ ] **Verify V31 succeeded** - Data must be in event_registrations
  ```sql
  SELECT COUNT(*) FROM event_registrations; -- Should have registrations
  ```

- [ ] **Check code references** - No code should use deleted columns
  ```bash
  # Search codebase for column references
  grep -r "\.rol" .
  grep -r "\.afstand" .
  grep -r "\.steps" . | grep -v "event_registrations"
  ```

- [ ] **Test in staging first** - NEVER run directly in production!

- [ ] **Schedule maintenance window** - Brief downtime recommended

- [ ] **Notify team** - Breaking changes ahead!

---

## Code Impact Analysis

### Models That Change

**File:** `models/participant.go`

**BEFORE V34:**
```go
type Participant struct {
    // Person data
    ID    uuid.UUID `gorm:"type:uuid;primaryKey"`
    Naam  string    `gorm:"type:text;not null"`
    Email string    `gorm:"type:text;not null"`
    
    // V30 fields
    AccountType   string    `gorm:"type:text"`
    HasAppAccess  bool      `gorm:"type:boolean"`
    
    // Legacy event fields (WILL BE DELETED IN V34!)
    Rol           string    `gorm:"type:text"` // ❌ DELETE THIS
    Afstand       string    `gorm:"type:text"` // ❌ DELETE THIS
    Steps         int       `gorm:"type:integer"` // ❌ DELETE THIS
    Ondersteuning string    `gorm:"type:text"` // ❌ DELETE THIS
    Status        string    `gorm:"type:text"` // ❌ DELETE THIS
}
```

**AFTER V34:**
```go
type Participant struct {
    // Person data
    ID    uuid.UUID `gorm:"type:uuid;primaryKey"`
    Naam  string    `gorm:"type:text;not null"`
    Email string    `gorm:"type:text;not null"`
    
    // V30 fields
    AccountType     string    `gorm:"type:text;not null;default:'temporary'"`
    HasAppAccess    bool      `gorm:"type:boolean;not null;default:false"`
    RegistrationYear int      `gorm:"type:integer"`
    WachtwoordHash  string    `gorm:"type:text"`
    GebruikerID     *uuid.UUID `gorm:"type:uuid"`
    
    // Metadata
    Terms     bool      `gorm:"type:boolean;not null;default:false"`
    TestMode  bool      `gorm:"type:boolean;not null;default:false"`
    CreatedAt time.Time `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
    UpdatedAt time.Time `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
}
// No more event-specific fields! ✅
```

### Repository Changes Required

**File:** `repository/participant_repository.go`

**BEFORE V34:**
```go
// participantColumns includes legacy columns
var participantColumns = []string{
    "id", "naam", "email",
    "rol",      // ❌ Will fail after V34
    "afstand",  // ❌ Will fail after V34
    "steps",    // ❌ Will fail after V34
}

func (r *ParticipantRepository) GetByID(id uuid.UUID) (*models.Participant, error) {
    var participant models.Participant
    err := r.db.Select(participantColumns).First(&participant, id).Error
    // ERROR after V34: column "rol" does not exist
    return &participant, err
}
```

**AFTER V34:**
```go
// participantColumns - only person data + V30 fields
var participantColumns = []string{
    "id", "naam", "email", "telefoon",
    "account_type", "registration_year", 
    "wachtwoord_hash", "has_app_access",
    "gebruiker_id", "upgraded_to_gebruiker_id", "upgraded_at",
    "terms", "test_mode", "created_at", "updated_at",
}

func (r *ParticipantRepository) GetByID(id uuid.UUID) (*models.Participant, error) {
    var participant models.Participant
    err := r.db.Select(participantColumns).First(&participant, id).Error
    return &participant, err // ✅ Works - no legacy columns
}

// If you need event data, join with event_registrations:
func (r *ParticipantRepository) GetWithEventData(
    participantID, eventID uuid.UUID,
) (*ParticipantWithEvent, error) {
    var result ParticipantWithEvent
    err := r.db.Table("participants p").
        Select(`
            p.*,
            er.participant_role_name,
            er.distance_route,
            er.steps,
            er.ondersteuning
        `).
        Joins("LEFT JOIN event_registrations er ON er.participant_id = p.id").
        Where("p.id = ? AND er.event_id = ?", participantID, eventID).
        Scan(&result).Error
    return &result, err
}
```

### Handler Changes Required

**File:** `handlers/participant_handler.go`

**BEFORE V34:**
```go
func (h *ParticipantHandler) GetParticipant(c *gin.Context) {
    participant, err := h.repo.GetByID(id)
    
    // Return participant with steps (WRONG - steps are event-specific!)
    c.JSON(200, gin.H{
        "id":    participant.ID,
        "naam":  participant.Naam,
        "steps": participant.Steps, // ❌ Fails after V34
    })
}
```

**AFTER V34:**
```go
func (h *ParticipantHandler) GetParticipant(c *gin.Context) {
    participant, err := h.participantRepo.GetByID(id)
    
    // Get event registration separately
    eventID := getCurrentEventID()
    registration, _ := h.eventRegRepo.GetByParticipantAndEvent(
        participant.ID, eventID,
    )
    
    c.JSON(200, gin.H{
        "id":    participant.ID,
        "naam":  participant.Naam,
        "email": participant.Email,
        "accountType": participant.AccountType,
        "eventData": gin.H{ // ✅ Correct structure
            "steps": registration.Steps,
            "role":  registration.ParticipantRoleName,
            "distance": registration.DistanceRoute,
        },
    })
}
```

---

## Migration Execution

### Development Environment

```bash
# 1. Backup
pg_dump -U postgres -d dklemailservice_dev > backup_dev_v34.sql

# 2. Run migration
docker-compose restart app
# Or manually:
psql -U postgres -d dklemailservice_dev -f database/migrations/V34__remove_legacy_participant_columns.sql

# 3. Verify
psql -U postgres -d dklemailservice_dev
\d participants  # Check schema
```

### Staging Environment

```bash
# 1. Backup
pg_dump -U $DB_USER -h staging-db -d dklemailservice > backup_staging_v34.sql

# 2. Deploy with migration
git push origin staging

# 3. Verify application logs
kubectl logs -f deployment/dkl-email-service

# 4. Test API endpoints
curl -X GET http://staging-api/api/participants/123
```

### Production Environment

```bash
# 1. Schedule maintenance window
# Notify users: "Brief maintenance from 02:00-02:15 UTC"

# 2. Backup (CRITICAL!)
pg_dump -U $DB_USER -h prod-db -d dklemailservice > backup_prod_v34_$(date +%Y%m%d_%H%M%S).sql

# 3. Deploy
git push origin main

# 4. Monitor
# Watch logs for errors
tail -f /var/log/dkl/app.log | grep -i error

# 5. Verify
psql -U $DB_USER -h prod-db -d dklemailservice
SELECT * FROM migraties WHERE version LIKE '%V34%';
\d participants  # Verify columns removed

# 6. Smoke test
curl -X GET https://api.dekoninklijkeloop.nl/api/health
curl -X GET https://api.dekoninklijkeloop.nl/api/participants/test-id
```

---

## Verification Steps

### Post-Migration Checks

```sql
-- 1. Migration applied?
SELECT * FROM migraties WHERE version = 'V34__remove_legacy_participant_columns';

-- 2. Columns actually removed?
\d participants
-- Should NOT see: rol, afstand, steps, ondersteuning, etc.

-- 3. Indexes removed?
\di idx_aanmeldingen*
-- Should return 0 results

-- 4. Table comment updated?
SELECT obj_description('participants'::regclass);
-- Should mention "V34: Cleaned up"

-- 5. Application starts?
-- Check logs for startup errors

-- 6. API endpoints work?
-- Test participant endpoints
```

### Functional Testing

```bash
# Test participant creation
curl -X POST http://localhost:8080/api/participants \
  -H "Content-Type: application/json" \
  -d '{"naam":"Test User","email":"test@example.com"}'

# Test participant retrieval
curl -X GET http://localhost:8080/api/participants/123

# Test event registration
curl -X POST http://localhost:8080/api/event-registrations \
  -H "Content-Type: application/json" \
  -d '{
    "participant_id":"123",
    "event_id":"456",
    "participant_role_name":"Deelnemer",
    "distance_route":"10KM"
  }'

# Test steps update (WebSocket or API)
# Should update event_registrations.steps, not participants.steps
```

---

## Rollback Procedure

### ⚠️ IMPORTANT: Rollback Limitations

**You CANNOT fully rollback V34** because:
1. Columns are dropped with CASCADE
2. Data was already migrated to event_registrations
3. Original column data is gone

**Best Option:** Restore from backup.

### Emergency Rollback

```bash
# 1. Stop application immediately
docker-compose down
# Or: kubectl scale deployment dkl-email-service --replicas=0

# 2. Restore database from backup
psql -U postgres -d dklemailservice < backup_before_v34.sql

# 3. Remove V34 migration record
psql -U postgres -d dklemailservice -c "
DELETE FROM migraties WHERE version = 'V34__remove_legacy_participant_columns';
"

# 4. Revert code to pre-V34 version
git revert <v34-commit-hash>
git push origin main

# 5. Restart application
docker-compose up -d
```

---

## Troubleshooting

### Error: Column does not exist

**Symptom:**
```
ERROR: column "steps" of relation "participants" does not exist
```

**Cause:** Code still references legacy columns

**Solution:**
```bash
# 1. Find all references
grep -rn "\.steps" handlers/ models/ repository/

# 2. Update to use event_registrations
# See "Code Impact Analysis" section above

# 3. Redeploy fixed code
```

### Error: View depends on column

**Symptom:**
```
ERROR: cannot drop column rol because other objects depend on it
DETAIL: view v_participant_summary depends on column rol
```

**Solution:** V34 uses CASCADE to handle this, but if it fails:
```sql
-- Manually drop dependent views first
DROP VIEW IF EXISTS v_participant_summary CASCADE;
DROP VIEW IF EXISTS v_participant_stats CASCADE;

-- Then re-run V34 migration
```

### Error: Application won't start

**Symptom:** Application crashes on startup after V34

**Diagnosis:**
```bash
# Check application logs
docker-compose logs app | grep -i "error\|fatal"

# Look for GORM errors
docker-compose logs app | grep -i "column.*does not exist"
```

**Solution:**
1. Identify which code references removed columns
2. Update model structs
3. Update repository column lists
4. Redeploy

---

## Success Criteria

V34 migration is successful when:

✅ Migration record exists in `migraties` table  
✅ All 14 legacy columns removed from `participants`  
✅ Legacy indexes dropped  
✅ Table comment updated  
✅ Application starts without errors  
✅ All API endpoints return 200 OK  
✅ Participant CRUD operations work  
✅ Event registration operations work  
✅ WebSocket steps updates work  
✅ No references to removed columns in logs  

---

## Related Migrations

| Migration | Relation | Description |
|-----------|----------|-------------|
| **V28** | Foundation | Started participant refactor (failed) |
| **V31** | Prerequisite | Completed data migration to event_registrations |
| **V32** | Related | Added supporting columns to lookup tables |
| **V30** | Parallel | Added account_type system |

---

## Documentation Links

- [V31 Participant Refactor](V31_PARTICIPANT_REFACTOR.md) - Data migration details
- [Database Architecture](../architecture/DATABASE.md) - New schema
- [V30 Database Tables](../V30_DATABASE_TABLES_UITLEG.md) - Table explanations
- [Migration Guide](../guides/MIGRATIONS.md) - General migration procedures

---

## Final Notes

### Why This Approach?

**Gradual Migration Strategy:**
1. V28: Attempted refactor (failed)
2. V31: Fixed and completed refactor, kept old columns
3. V34: Remove old columns after verification

**Benefits:**
- ✅ Safe transition period
- ✅ Rollback capability during transition
- ✅ Time to update all code
- ✅ Verification of data migration

**Alternative (Not Chosen):**
- ❌ Drop columns immediately in V28/V31
- ❌ Higher risk
- ❌ No transition period

---

**Last Updated:** 2025-11-10  
**Author:** Development Team  
**Review Status:** ✅ Documented Post-Implementation  
**Risk Assessment:** 🔴🔴 CRITICAL - Permanent data structure change