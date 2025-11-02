# 🔍 RBAC Database Verification Guide

> **Versie:** V1.49  
> **Datum:** 2025-11-02

---

## 📋 Overzicht

Dit document beschrijft hoe je de RBAC database tables kunt verifiëren in verschillende environments.

**Verification Script:** [`verify_rbac_tables.sql`](verify_rbac_tables.sql)

---

## 🚀 Uitvoeren van Verificatie

### Option 1: Docker Container (Development)

```bash
# 1. Vind de database container naam
docker ps | grep postgres

# 2. Voer verificatie script uit
docker exec -i <postgres_container_name> psql -U <username> -d <dbname> < database/scripts/verify_rbac_tables.sql

# Voorbeeld:
docker exec -i dklemailservice-db-1 psql -U dkluser -d dklemailservice < database/scripts/verify_rbac_tables.sql
```

### Option 2: Direct psql Connection (Production)

```bash
# Met environment variables
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f database/scripts/verify_rbac_tables.sql

# Of direct:
psql -h localhost -p 5432 -U dkluser -d dklemailservice -f database/scripts/verify_rbac_tables.sql
```

### Option 3: Via docker-compose exec

```bash
# Als je docker-compose gebruikt
docker-compose exec db psql -U dkluser -d dklemailservice -f /database/scripts/verify_rbac_tables.sql
```

### Option 4: Copy-Paste in psql Shell

```bash
# Start psql shell
docker exec -it <postgres_container> psql -U <username> <dbname>

# Dan copy-paste de inhoud van verify_rbac_tables.sql
```

---

## ✅ Verwachte Output

### 1. Table Existence Check
```
1. CHECKING TABLE EXISTENCE...

 table_name        | status
-------------------+-----------
 gebruikers        | ✓ EXISTS
 permissions       | ✓ EXISTS
 refresh_tokens    | ✓ EXISTS
 role_permissions  | ✓ EXISTS
 roles             | ✓ EXISTS
 user_roles        | ✓ EXISTS
```

### 2. System Roles Verification
```
5. VERIFYING SYSTEM ROLES...

 expected_name  | status     | description
----------------+------------+------------------------------------------
 admin          | ✓ EXISTS   | Volledige beheerder met toegang tot...
 staff          | ✓ EXISTS   | Ondersteunend personeel met beperkte...
 user           | ✓ EXISTS   | Standaard gebruiker
 owner          | ✓ EXISTS   | Chat kanaal eigenaar
 chat_admin     | ✓ EXISTS   | Chat kanaal beheerder
 member         | ✓ EXISTS   | Chat kanaal lid
 deelnemer      | ✓ EXISTS   | Evenement deelnemer
 begeleider     | ✓ EXISTS   | Evenement begeleider
 vrijwilliger   | ✓ EXISTS   | Evenement vrijwilliger
```

### 3. Permissions Count
```
6. PERMISSIONS PER ROLE...

 role_name    | is_system_role | permission_count | permissions
--------------+----------------+------------------+-------------
 admin        | t              | 58               | admin:access, contact:read, ...
 staff        | t              | 12               | staff:access, contact:read, ...
 owner        | t              | 4                | chat:read, chat:write, ...
 chat_admin   | t              | 3                | chat:read, chat:write, chat:moderate
 member       | t              | 2                | chat:read, chat:write
 user         | t              | 2                | chat:read, chat:write
 deelnemer    | t              | 0                | (event categorization only)
 begeleider   | t              | 0                | (event categorization only)
 vrijwilliger | t              | 0                | (event categorization only)
```

### 4. Legacy vs RBAC Status
```
8. LEGACY VS RBAC ROLE COMPARISON...

 email                          | legacy_role | rbac_roles      | status
--------------------------------+-------------+-----------------+--------------
 admin@dekoninklijkeloop.nl     | admin       | admin           | ✓ MIGRATED
 jeffrey@dekoninklijkeloop.nl   | staff       | staff           | ✓ MIGRATED
 user@dekoninklijkeloop.nl      | user        | user            | ✓ MIGRATED
```

### 5. Summary Statistics
```
13. SUMMARY STATISTICS...

 check_item                           | actual  | status
--------------------------------------+---------+---------
 9 System Roles Expected              | 9       | ✓ PASS
 At least 25 Permissions Expected     | 58      | ✓ PASS
 All Active Users Have RBAC Roles     | 5 of 5  | ✓ PASS
 Roles Table Populated                | 9       | ✓ PASS
 Permissions Table Populated          | 58      | ✓ PASS
 Valid Refresh Tokens                 | 3       | ✓ INFO
```

---

## 🔴 Problematische Cases

### ✗ FAIL Status

Als je `✗ FAIL` ziet bij een check, betekent dit een **kritiek probleem**:

| Issue | Mogelijke Oorzaak | Oplossing |
|-------|------------------|-----------|
| Missing tables | Migraties niet uitgevoerd | Run migrations V1.20, V1.21, V1.22, V1.28 |
| < 9 system roles | Seed data niet geladen | Run migration V1.21 |
| < 25 permissions | Permissions migrations missing | Run migrations V1.24-V1.45 |
| System roles not system | Data corruption | UPDATE roles SET is_system_role = true WHERE name IN (...) |

### ⚠ WARNING Status

Als je `⚠ WARNING` ziet:

| Issue | Mogelijke Oorzaak | Oplossing |
|-------|------------------|-----------|
| Users without RBAC roles | Nieuwe users na migratie | Assign roles via API/admin panel |
| MISMATCH between legacy and RBAC | Manual changes | Sync via admin panel |
| Unassigned permissions | New permissions added | Assign to appropriate roles |

### ✓ INFO Status

Dit zijn informatieve checks, geen problemen.

---

## 🛠️ Common Fixes

### Fix 1: Assign Missing Roles

```sql
-- Assign 'user' role to users without any roles
INSERT INTO user_roles (user_id, role_id, is_active)
SELECT 
    g.id,
    r.id,
    true
FROM gebruikers g
CROSS JOIN roles r
WHERE r.name = 'user' 
  AND r.is_system_role = true
  AND NOT EXISTS (
      SELECT 1 FROM user_roles ur 
      WHERE ur.user_id = g.id AND ur.is_active = true
  )
ON CONFLICT (user_id, role_id) DO NOTHING;
```

### Fix 2: Clean Up Expired Tokens

```sql
-- Remove expired and revoked refresh tokens (older than 30 days)
DELETE FROM refresh_tokens 
WHERE (is_revoked = true OR expires_at < NOW() - INTERVAL '30 days');
```

### Fix 3: Deactivate Expired Roles

```sql
-- Deactivate user roles that have expired
UPDATE user_roles 
SET is_active = false
WHERE expires_at IS NOT NULL 
  AND expires_at <= NOW()
  AND is_active = true;
```

### Fix 4: Sync Legacy Role to RBAC

```sql
-- Voor één specifieke gebruiker
INSERT INTO user_roles (user_id, role_id, is_active)
SELECT 
    g.id,
    r.id,
    true
FROM gebruikers g
CROSS JOIN roles r
WHERE g.email = 'user@example.com'
  AND r.name = 'admin' -- of andere role
ON CONFLICT (user_id, role_id) DO UPDATE 
SET is_active = true;
```

---

## 📊 Monitoring Queries

### Daily Checks

```sql
-- 1. Active users without permissions
SELECT COUNT(*) as users_without_permissions
FROM gebruikers g
WHERE g.is_actief = true
  AND NOT EXISTS (
      SELECT 1 FROM user_roles ur
      JOIN role_permissions rp ON ur.role_id = rp.role_id
      WHERE ur.user_id = g.id AND ur.is_active = true
  );

-- 2. Orphaned records count
SELECT 
    (SELECT COUNT(*) FROM user_roles ur 
     LEFT JOIN gebruikers g ON ur.user_id = g.id 
     WHERE g.id IS NULL) as orphaned_user_roles,
    (SELECT COUNT(*) FROM refresh_tokens rt 
     LEFT JOIN gebruikers g ON rt.user_id = g.id 
     WHERE g.id IS NULL) as orphaned_tokens;

-- 3. Permission distribution
SELECT 
    r.name,
    COUNT(DISTINCT ur.user_id) as active_users
FROM roles r
LEFT JOIN user_roles ur ON r.id = ur.role_id AND ur.is_active = true
GROUP BY r.name
ORDER BY active_users DESC;
```

### Weekly Checks

```sql
-- 1. Refresh token usage
SELECT 
    DATE(created_at) as date,
    COUNT(*) as tokens_created,
    COUNT(*) FILTER (WHERE is_revoked = true) as tokens_revoked
FROM refresh_tokens
WHERE created_at > NOW() - INTERVAL '7 days'
GROUP BY DATE(created_at)
ORDER BY date DESC;

-- 2. Role assignment trends
SELECT 
    r.name as role_name,
    DATE(ur.assigned_at) as date,
    COUNT(*) as assignments
FROM user_roles ur
JOIN roles r ON ur.role_id = r.id
WHERE ur.assigned_at > NOW() - INTERVAL '7 days'
GROUP BY r.name, DATE(ur.assigned_at)
ORDER BY date DESC, role_name;
```

---

## 🎯 Expected Database State

### Roles Table
- **Total:** 9+ (9 system + custom roles)
- **System Roles:** 9 (admin, staff, user, owner, chat_admin, member, deelnemer, begeleider, vrijwilliger)
- **All marked:** `is_system_role = true`

### Permissions Table
- **Total:** 58+ system permissions
- **Resources:** 19 (admin, staff, contact, aanmelding, user, photo, album, video, partner, sponsor, radio_recording, program_schedule, social_embed, social_link, under_construction, newsletter, email, admin_email, chat)
- **All marked:** `is_system_permission = true`

### Role_Permissions Table
- **Admin role:** 58 permissions (ALL)
- **Staff role:** ~12 permissions (read-only)
- **Chat roles:** 2-4 chat permissions
- **Event roles:** 0 permissions (categorization only)

### User_Roles Table
- **All active users:** Should have at least one active role
- **is_active:** true for current assignments
- **expires_at:** NULL for permanent assignments

### Refresh_Tokens Table
- **Valid tokens:** expires_at > NOW() AND is_revoked = false
- **Cleanup:** Old revoked/expired tokens should be cleaned periodically

---

## 🔧 Troubleshooting

### Problem: "Table does not exist"

**Cause:** Migrations not run

**Solution:**
```bash
# Check which migrations are applied
SELECT versie, naam FROM migraties ORDER BY versie;

# If V1.20, V1.21, V1.22 missing, run migrations
go run main.go  # Will auto-run migrations on startup
```

### Problem: "Users without RBAC roles"

**Cause:** New users created after migration, or migration failed

**Solution:**
```sql
-- Run the user-role sync from V1_22 migration manually
INSERT INTO user_roles (user_id, role_id, is_active)
SELECT g.id, r.id, true
FROM gebruikers g
CROSS JOIN roles r
WHERE r.name = 'user' 
  AND NOT EXISTS (SELECT 1 FROM user_roles WHERE user_id = g.id)
ON CONFLICT DO NOTHING;
```

### Problem: "Permission count mismatch"

**Cause:** Migrations V1.24-V1.45 not all applied

**Solution:**
```bash
# Check migration status
SELECT versie, naam FROM migraties 
WHERE versie LIKE 'V1.%' 
ORDER BY versie;

# Latest should be V1.48 or higher
```

---

## 📈 Performance Checks

Run these to verify indexes are working:

```sql
-- 1. Check index usage
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan as index_scans,
    idx_tup_read as tuples_read,
    idx_tup_fetch as tuples_fetched
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
  AND tablename IN ('roles', 'permissions', 'role_permissions', 'user_roles')
ORDER BY idx_scan DESC;

-- 2. Check table sizes
SELECT 
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size,
    pg_total_relation_size(schemaname||'.'||tablename) as bytes
FROM pg_tables
WHERE schemaname = 'public'
  AND tablename IN ('roles', 'permissions', 'role_permissions', 'user_roles', 'refresh_tokens')
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

---

## 🎯 Quick Verification (One-Liner)

```bash
# Docker environment
docker exec -i $(docker ps --filter "name=postgres" --format "{{.Names}}" | head -1) psql -U dkluser -d dklemailservice -c "
SELECT 
    (SELECT COUNT(*) FROM roles WHERE is_system_role = true) as system_roles,
    (SELECT COUNT(*) FROM permissions WHERE is_system_permission = true) as system_permissions,
    (SELECT COUNT(DISTINCT user_id) FROM user_roles WHERE is_active = true) as users_with_roles,
    (SELECT COUNT(*) FROM gebruikers WHERE is_actief = true) as active_users,
    CASE 
        WHEN (SELECT COUNT(DISTINCT user_id) FROM user_roles WHERE is_active = true) >= 
             (SELECT COUNT(*) FROM gebruikers WHERE is_actief = true)
        THEN '✓ ALL USERS HAVE ROLES'
        ELSE '✗ SOME USERS MISSING ROLES'
    END as status;
"
```

**Expected Output:**
```
 system_roles | system_permissions | users_with_roles | active_users | status
--------------+--------------------+------------------+--------------+-------------------------
           9  |                 58 |              5   |           5  | ✓ ALL USERS HAVE ROLES
```

---

## 📝 Automated Health Check

Voeg deze check toe aan je monitoring:

```bash
#!/bin/bash
# rbac_health_check.sh

DB_CONTAINER="dklemailservice-db-1"
DB_USER="dkluser"
DB_NAME="dklemailservice"

# Run verification
docker exec -i $DB_CONTAINER psql -U $DB_USER -d $DB_NAME -f /path/to/verify_rbac_tables.sql

# Check for failures
if grep -q "✗ FAIL" verification_output.log; then
    echo "❌ RBAC verification failed!"
    # Send alert
    exit 1
else
    echo "✅ RBAC verification passed"
    exit 0
fi
```

---

## 🔍 Deep Dive Queries

### Check Specific User Permissions

```sql
SELECT 
    g.email,
    r.name as role_name,
    p.resource || ':' || p.action as permission,
    ur.is_active,
    ur.expires_at
FROM gebruikers g
JOIN user_roles ur ON g.id = ur.user_id
JOIN roles r ON ur.role_id = r.id
JOIN role_permissions rp ON r.id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
WHERE g.email = 'admin@dekoninklijkeloop.nl'
  AND ur.is_active = true
ORDER BY p.resource, p.action;
```

### Find Users with Specific Permission

```sql
-- Wie heeft user:manage_roles permission?
SELECT DISTINCT
    g.email,
    g.naam,
    r.name as role_name
FROM gebruikers g
JOIN user_roles ur ON g.id = ur.user_id
JOIN roles r ON ur.role_id = r.id
JOIN role_permissions rp ON r.id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
WHERE p.resource = 'user' 
  AND p.action = 'manage_roles'
  AND ur.is_active = true
ORDER BY g.email;
```

### Audit Recent Role Changes

```sql
-- Recent role assignments (last 7 days)
SELECT 
    g.email as user_email,
    r.name as role_name,
    ur.assigned_at,
    ur.is_active,
    g2.email as assigned_by_email
FROM user_roles ur
JOIN gebruikers g ON ur.user_id = g.id
JOIN roles r ON ur.role_id = r.id
LEFT JOIN gebruikers g2 ON ur.assigned_by = g2.id
WHERE ur.assigned_at > NOW() - INTERVAL '7 days'
ORDER BY ur.assigned_at DESC;
```

---

## 📅 Maintenance Schedule

### Daily
- [ ] Run quick verification one-liner
- [ ] Check for users without roles
- [ ] Monitor orphaned records count

### Weekly
- [ ] Run full verification script
- [ ] Clean up expired refresh tokens (>30 days old)
- [ ] Review role assignment trends
- [ ] Audit permission changes

### Monthly
- [ ] Full RBAC audit
- [ ] Performance check (index usage)
- [ ] Review custom roles (non-system)
- [ ] Database backup verification

---

## 🚨 Emergency Procedures

### If RBAC System Fails

1. **Check migrations:**
   ```sql
   SELECT * FROM migraties WHERE versie LIKE '1.2%' ORDER BY versie;
   ```

2. **Verify admin access:**
   ```sql
   SELECT * FROM user_permissions WHERE email LIKE '%admin%';
   ```

3. **Emergency admin grant:**
   ```sql
   -- Grant admin role to specific user (EMERGENCY ONLY)
   INSERT INTO user_roles (user_id, role_id, is_active)
   SELECT g.id, r.id, true
   FROM gebruikers g
   CROSS JOIN roles r
   WHERE g.email = 'emergency@admin.com'
     AND r.name = 'admin'
   ON CONFLICT (user_id, role_id) DO UPDATE 
   SET is_active = true;
   ```

4. **Clear Redis cache:**
   ```bash
   redis-cli FLUSHDB  # Clears permission cache
   ```

---

## ✅ Pre-Deployment Checklist

Voer deze checks uit **voor deployment** naar productie:

- [ ] Full verification script succesvol (geen ✗ FAIL)
- [ ] Alle system roles aanwezig (9 roles)
- [ ] Alle permissions aanwezig (58+)
- [ ] Admin user heeft admin role
- [ ] Geen orphaned records
- [ ] Alle indexes present
- [ ] Foreign key constraints intact
- [ ] Geen duplicate user-role assignments
- [ ] Migration versie is V1.48+

---

## 📞 Support

Bij problemen met RBAC database verificatie:

1. Check logs: `docker-compose logs db`
2. Review migrations: `SELECT * FROM migraties ORDER BY versie`
3. Consult: [`RBAC_SECURITY_AUDIT.md`](../../docs/RBAC_SECURITY_AUDIT.md)
4. Emergency: Use emergency admin grant procedure

---

**Laatste Update:** 2025-11-02  
**Script Versie:** 1.49  
**Compatibility:** PostgreSQL 12+