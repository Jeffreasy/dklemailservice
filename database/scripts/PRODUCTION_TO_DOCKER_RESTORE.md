# 🔄 Productie naar Docker Database Restore Guide

> **Versie:** 1.49  
> **Datum:** 2025-11-02  
> **Doel:** Importeer Render.com productie database naar lokale Docker

---

## 📋 Productie Database Info

**Host:** dpg-cva4c01c1ekc738q6q0g-a.oregon-postgres.render.com  
**User:** dekoninklijkeloopdatabase_user  
**Database:** dekoninklijkeloopdatabase  
**Password:** I4QP3JwyCcEbn8tGl6k3ErEvjUZ9V5rB

---

## 🚀 Methode 1: Via pgAdmin (Eenvoudigst)

### Stap 1: Connect naar Productie in pgAdmin

1. Open pgAdmin
2. Rechtermuisklik "Servers" → "Register" → "Server"
3. **General Tab:**
   - Name: `DKL Production (Render)`
4. **Connection Tab:**
   - Host: `dpg-cva4c01c1ekc738q6q0g-a.oregon-postgres.render.com`
   - Port: `5432`
   - Database: `dekoninklijkeloopdatabase`
   - Username: `dekoninklijkeloopdatabase_user`
   - Password: `I4QP3JwyCcEbn8tGl6k3ErEvjUZ9V5rB`
   - ✅ Save password
5. Klik "Save"

### Stap 2: Backup Productie Database

1. Rechtermuisklik op `dekoninklijkeloopdatabase` (onder DKL Production)
2. Selecteer "Backup..."
3. **General Tab:**
   - Filename: `C:\Temp\dkl_production_backup.sql`
   - Format: `Plain`
4. **Data Options Tab:**
   - ✅ Data
   - ✅ Blobs
5. **Query Options Tab:**
   - ✅ Use INSERT commands
   - ✅ Include CREATE DATABASE statement: **NO**
6. Klik "Backup"
7. Wacht tot backup compleet is

### Stap 3: Connect naar Docker Database in pgAdmin

1. Rechtermuisklink "Servers" → "Register" → "Server"
2. **General Tab:**
   - Name: `DKL Docker Local`
3. **Connection Tab:**
   - Host: `localhost`
   - Port: `5432` (of check met `docker port dkl-postgres`)
   - Database: `dekoninklijkeloopdatabase`
   - Username: `dekoninklijkeloopdatabase_user`
   - Password: `I4QP3JwyCcEbn8tGl6k3ErEvjUZ9V5rB`
4. Klik "Save"

### Stap 4: Restore naar Docker

1. Rechtermuisklik op `dekoninklijkeloopdatabase` (onder DKL Docker Local)
2. Selecteer "Restore..."
3. **General Tab:**
   - Filename: `C:\Temp\dkl_production_backup.sql`
   - Format: `Plain`
4. **Data Options:**
   - ✅ Clean before restore (voorzichtig!)
5. Klik "Restore"
6. Wacht tot restore compleet is

### Stap 5: Assign RBAC Roles

**Open Query Tool in Docker database en voer uit:**

```sql
-- Run de sync script
\i C:/Users/jeffrey/Desktop/Githubmains/dklemailservice/database/scripts/sync_production_users_with_rbac.sql
```

---

## 🚀 Methode 2: Via Docker Container

### Voorbereiden

```bash
# Check Docker postgres container naam
docker ps | findstr postgres

# Resultaat: dkl-postgres
```

### Stap 1: Maak Backup van Productie via Docker

```bash
# Gebruik Docker postgres client om productie te dumpen
docker run --rm \
  -e PGPASSWORD=I4QP3JwyCcEbn8tGl6k3ErEvjUZ9V5rB \
  postgres:15-alpine \
  pg_dump \
  -h dpg-cva4c01c1ekc738q6q0g-a.oregon-postgres.render.com \
  -U dekoninklijkeloopdatabase_user \
  -d dekoninklijkeloopdatabase \
  --no-owner --no-acl \
  > production_backup.sql
```

### Stap 2: Backup Huidige Docker Database (Safety!)

```bash
# Backup current Docker DB first
docker exec dkl-postgres pg_dump \
  -U dekoninklijkeloopdatabase_user \
  dekoninklijkeloopdatabase \
  > docker_backup_before_restore.sql
```

### Stap 3: Stop Applicatie Container

```bash
# Stop app zodat geen writes naar DB tijdens restore
docker stop dkl-email-service
```

### Stap 4: Drop en Recreate Database

```bash
# Drop database (LET OP: Dit verwijdert alle data!)
docker exec dkl-postgres psql -U dekoninklijkeloopdatabase_user -d postgres -c "DROP DATABASE IF EXISTS dekoninklijkeloopdatabase;"

# Create fresh database
docker exec dkl-postgres psql -U dekoninklijkeloopdatabase_user -d postgres -c "CREATE DATABASE dekoninklijkeloopdatabase OWNER dekoninklijkeloopdatabase_user;"
```

### Stap 5: Restore Productie Data

```bash
# Kopieer backup naar container
docker cp production_backup.sql dkl-postgres:/tmp/

# Restore
docker exec dkl-postgres psql \
  -U dekoninklijkeloopdatabase_user \
  -d dekoninklijkeloopdatabase \
  -f /tmp/production_backup.sql
```

### Stap 6: Assign RBAC Roles

```bash
# Kopieer sync script naar container
docker cp database/scripts/sync_production_users_with_rbac.sql dkl-postgres:/tmp/

# Voer sync uit
docker exec dkl-postgres psql \
  -U dekoninklijkeloopdatabase_user \
  -d dekoninklijkeloopdatabase \
  -f /tmp/sync_production_users_with_rbac.sql
```

### Stap 7: Verify Results

```bash
docker exec dkl-postgres psql \
  -U dekoninklijkeloopdatabase_user \
  -d dekoninklijkeloopdatabase \
  -c "
SELECT 
    g.email,
    g.rol as legacy,
    STRING_AGG(DISTINCT r.name, ', ') as rbac_roles,
    COUNT(DISTINCT p.id) as permissions
FROM gebruikers g
LEFT JOIN user_roles ur ON g.id = ur.user_id AND ur.is_active = true
LEFT JOIN roles r ON ur.role_id = r.id
LEFT JOIN role_permissions rp ON r.id = rp.role_id
LEFT JOIN permissions p ON rp.permission_id = p.id
WHERE g.email LIKE '%@dekoninklijkeloop.nl'
GROUP BY g.email, g.rol
ORDER BY g.email;
"
```

**Verwachte Output:**
```
Email                          | Legacy  | RBAC Roles | Permissions
-------------------------------+---------+------------+-------------
admin@dekoninklijkeloop.nl     | admin   | admin      | 68
jeffrey@dekoninklijkeloop.nl   | staff   | staff      | 25
Lida@dekoninklijkeloop.nl      | staff   | staff      | 25
Marieke@dekoninklijkeloop.nl   | staff   | staff      | 25
Salih@dekoninklijkeloop.nl     | staff   | staff      | 25
```

### Stap 8: Restart Applicatie

```bash
# Start app weer
docker start dkl-email-service

# Check logs
docker logs -f dkl-email-service
```

---

## 🚀 Methode 3: Windows PowerShell Script (Alles-in-één)

Sla dit op als `restore_production.ps1`:

```powershell
# Production to Docker Restore Script
$ErrorActionPreference = "Stop"

Write-Host "=== Productie Database Restore ===" -ForegroundColor Blue

# Environment
$PROD_HOST = "dpg-cva4c01c1ekc738q6q0g-a.oregon-postgres.render.com"
$PROD_USER = "dekoninklijkeloopdatabase_user"
$PROD_DB = "dekoninklijkeloopdatabase"
$PROD_PASS = "I4QP3JwyCcEbn8tGl6k3ErEvjUZ9V5rB"
$DOCKER_CONTAINER = "dkl-postgres"

Write-Host "[1/7] Stopping application..." -ForegroundColor Yellow
docker stop dkl-email-service

Write-Host "[2/7] Backing up current Docker database..." -ForegroundColor Yellow
docker exec $DOCKER_CONTAINER pg_dump -U $PROD_USER $PROD_DB > "docker_backup_$(Get-Date -Format 'yyyyMMdd_HHmmss').sql"

Write-Host "[3/7] Dumping production database..." -ForegroundColor Yellow
docker run --rm `
  -e PGPASSWORD=$PROD_PASS `
  postgres:15-alpine `
  pg_dump `
  -h $PROD_HOST `
  -U $PROD_USER `
  -d $PROD_DB `
  --no-owner --no-acl `
  > production_backup.sql

Write-Host "[4/7] Dropping current Docker database..." -ForegroundColor Yellow
docker exec $DOCKER_CONTAINER psql -U $PROD_USER -d postgres -c "DROP DATABASE IF EXISTS $PROD_DB;"
docker exec $DOCKER_CONTAINER psql -U $PROD_USER -d postgres -c "CREATE DATABASE $PROD_DB OWNER $PROD_USER;"

Write-Host "[5/7] Restoring production data to Docker..." -ForegroundColor Yellow
docker cp production_backup.sql ${DOCKER_CONTAINER}:/tmp/
docker exec $DOCKER_CONTAINER psql -U $PROD_USER -d $PROD_DB -f /tmp/production_backup.sql

Write-Host "[6/7] Assigning RBAC roles..." -ForegroundColor Yellow
docker cp database/scripts/sync_production_users_with_rbac.sql ${DOCKER_CONTAINER}:/tmp/
docker exec $DOCKER_CONTAINER psql -U $PROD_USER -d $PROD_DB -f /tmp/sync_production_users_with_rbac.sql

Write-Host "[7/7] Verifying results..." -ForegroundColor Yellow
docker exec $DOCKER_CONTAINER psql -U $PROD_USER -d $PROD_DB -c "
SELECT COUNT(*) as total_users, 
       COUNT(DISTINCT ur.user_id) as with_rbac_roles
FROM gebruikers g
LEFT JOIN user_roles ur ON g.id = ur.user_id AND ur.is_active = true
WHERE g.is_actief = true;
"

Write-Host "`n=== RESTORE COMPLETE ===" -ForegroundColor Green
Write-Host "Start application: docker start dkl-email-service" -ForegroundColor Green
```

**Voer uit:**
```powershell
powershell -ExecutionPolicy Bypass -File restore_production.ps1
```

---

## ✅ Verificatie Na Restore

```bash
# Check alle @dekoninklijkeloop.nl users hebben staff/admin role
docker exec dkl-postgres psql -U dekoninklijkeloopdatabase_user -d dekoninklijkeloopdatabase -c "
SELECT 
    email,
    STRING_AGG(r.name, ', ') as roles,
    COUNT(p.id) as permissions
FROM gebruikers g
JOIN user_roles ur ON g.id = ur.user_id AND ur.is_active = true
JOIN roles r ON ur.role_id = r.id
JOIN role_permissions rp ON r.id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
WHERE g.email LIKE '%@dekoninklijkeloop.nl'
GROUP BY g.email
ORDER BY g.email;
"
```

**Verwacht:**
- ✅ 5 @dekoninklijkeloop.nl users
- ✅ admin@: admin role (68 permissions)
- ✅ jeffrey@, Lida@, Marieke@, Salih@: staff role (25 permissions each)

---

## ⚠️ BELANGRIJK

**Voor Production Restore:**
1. ✅ Maak ALTIJD eerst een backup van huidige Docker DB
2. ✅ Stop applicatie tijdens restore (avoid data corruption)
3. ✅ Test de applicatie na restore voordat je verder gaat
4. ✅ Check logs voor errors: `docker logs dkl-email-service`

**Na Restore:**
1. Verify RBAC roles zijn toegewezen
2. Test login met verschillende users
3. Check permissions werken correct
4. Monitor logs voor errors

---

## 🐛 Troubleshooting

### "pg_dump not found"
- Gebruik Docker methode (Method 2) 
- Of gebruik pgAdmin GUI (Method 1 - Eenvoudigst!)

### "Connection refused"
- Check firewall
- Verify Render.com whitelisting
- Try via VPN

### "Database already exists"
- Add `--clean` flag to pg_dump
- Or manually DROP DATABASE eerst

### "Permission denied"
- Verify password correct
- Check user has sufficient privileges

---

## 📞 Support

Bij problemen:
1. Check Docker logs: `docker logs dkl-postgres`
2. Verify backup file created: `dir production_backup.sql`
3. Test productie connection: `docker run postgres:15-alpine psql -h dpg-cva4c01c1ekc738q6q0g-a.oregon-postgres.render.com -U dekoninklijkeloopdatabase_user -d dekoninklijkeloopdatabase -c "SELECT COUNT(*) FROM gebruikers;"`

---

**Aanbevolen Methode:** **pgAdmin GUI** (Methode 1) - Meest gebruiksvriendelijk!