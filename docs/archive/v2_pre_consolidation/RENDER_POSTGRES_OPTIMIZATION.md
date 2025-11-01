# Render PostgreSQL Optimalisatie Guide

Specifieke instructies voor het optimaliseren van de DKL Email Service database op Render.com.

---

## 🔗 Jouw Database Configuratie

**Host**: `dpg-cva4c01c1ekc738q6q0g-a.oregon-postgres.render.com`  
**Database**: `dekoninklijkeloopdatabase`  
**User**: `dekoninklijkeloopdatabase_user`  
**Connection String**: 
```
postgresql://dekoninklijkeloopdatabase_user:I4QP3JwyCcEbn8tGl6k3ErEvjUZ9V5rB@dpg-cva4c01c1ekc738q6q0g-a/dekoninklijkeloopdatabase
```

⚠️ **BELANGRIJK**: Bewaar credentials veilig, gebruik alleen in environment variables!

---

## 🚀 V1_47 Migratie Toepassen op Render

### Optie 1: Automatisch via Applicatie (Aanbevolen)

De migratie wordt automatisch toegepast wanneer je de applicatie op Render herdeployt met de nieuwe code:

```bash
# 1. Push naar Git repository
git add database/migrations/V1_47__performance_optimizations.sql
git commit -m "Add database performance optimizations (V1_47)"
git push origin main

# 2. Render detecteert de push en redeploys automatisch
# 3. Bij startup voert de applicatie automatisch de migratie uit
```

**Verificatie via Render Logs:**
```
Zoek in de logs naar:
"Migratie wordt uitgevoerd","file":"V1_47__performance_optimizations.sql"
"Migratie succesvol uitgevoerd","file":"V1_47__performance_optimizations.sql"
```

### Optie 2: Handmatig via psql

Als je de migratie handmatig wil uitvoeren:

#### Installeer psql lokaal (Windows)

```powershell
# Via Chocolatey
choco install postgresql

# Of download van https://www.postgresql.org/download/windows/
```

#### Voer Migratie Uit

```bash
# Via psql
psql "postgresql://dekoninklijkeloopdatabase_user:I4QP3JwyCcEbn8tGl6k3ErEvjUZ9V5rB@dpg-cva4c01c1ekc738q6q0g-a.oregon-postgres.render.com/dekoninklijkeloopdatabase" -f database/migrations/V1_47__performance_optimizations.sql

# Of via redirect
psql "postgresql://dekoninklijkeloopdatabase_user:I4QP3JwyCcEbn8tGl6k3ErEvjUZ9V5rB@dpg-cva4c01c1ekc738q6q0g-a.oregon-postgres.render.com/dekoninklijkeloopdatabase" < database/migrations/V1_47__performance_optimizations.sql
```

---

## 🔍 Database Inspecteren op Render

### Verbinden met psql

```bash
# Basis verbinding
psql "postgresql://dekoninklijkeloopdatabase_user:I4QP3JwyCcEbn8tGl6k3ErEvjUZ9V5rB@dpg-cva4c01c1ekc738q6q0g-a.oregon-postgres.render.com/dekoninklijkeloopdatabase"

# Of met environment variable
export DATABASE_URL="postgresql://dekoninklijkeloopdatabase_user:I4QP3JwyCcEbn8tGl6k3ErEvjUZ9V5rB@dpg-cva4c01c1ekc738q6q0g-a.oregon-postgres.render.com/dekoninklijkeloopdatabase"
psql $DATABASE_URL
```

### Check Migratie Status

```sql
-- Laatste 10 migraties
SELECT versie, naam, toegepast 
FROM migraties 
ORDER BY toegepast DESC 
LIMIT 10;

-- Check of V1_47 is toegepast
SELECT * FROM migraties WHERE versie = '1.47.0';
```

### Check Indexes

```sql
-- Alle indexes op een tabel
SELECT 
    indexname,
    indexdef
FROM pg_indexes
WHERE schemaname = 'public'
    AND tablename = 'verzonden_emails'
ORDER BY indexname;

-- Check of nieuwe indexes bestaan
SELECT 
    schemaname,
    tablename,
    indexname
FROM pg_indexes
WHERE schemaname = 'public'
    AND indexname LIKE 'idx_%'
ORDER BY tablename, indexname;
```

### Database Statistieken

```sql
-- Tabel groottes
SELECT 
    tablename,
    pg_size_pretty(pg_total_relation_size('public.'||tablename)) AS total_size,
    pg_size_pretty(pg_relation_size('public.'||tablename)) AS table_size,
    pg_size_pretty(pg_total_relation_size('public.'||tablename) - pg_relation_size('public.'||tablename)) AS indexes_size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size('public.'||tablename) DESC
LIMIT 10;

-- Index usage
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan as scans,
    pg_size_pretty(pg_relation_size(indexrelid)) as size
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
ORDER BY idx_scan DESC
LIMIT 20;
```

---

## 🛠️ Maintenance op Render

### VACUUM ANALYZE (Wekelijks)

```bash
# Via psql
psql "$DATABASE_URL" -c "VACUUM ANALYZE;"

# Per tabel
psql "$DATABASE_URL" -c "VACUUM ANALYZE verzonden_emails;"
psql "$DATABASE_URL" -c "VACUUM ANALYZE chat_messages;"
```

### Update Statistics

```bash
# Na bulk inserts of deletes
psql "$DATABASE_URL" -c "ANALYZE;"
```

### Data Cleanup (Maandelijks)

⚠️ **BELANGRIJK: Maak eerst backup via Render Dashboard!**

Render Dashboard → Database → Manual Backups → Create Backup

```bash
# Via het data_cleanup.sql script
psql "$DATABASE_URL" < database/scripts/data_cleanup.sql
```

---

## 📊 Render-Specifieke Beperkingen

### Wat je NIET kan doen op Render Free/Starter Plan:

- ❌ Direct superuser access
- ❌ Custom postgresql.conf (managed service)
- ❌ pg_cron extension voor scheduled jobs
- ❌ File system access voor custom scripts
- ❌ Replication setup

### Wat je WEL kan doen:

- ✅ Alle DDL statements (CREATE INDEX, ALTER TABLE, etc.)
- ✅ VACUUM en ANALYZE
- ✅ Extensions (via SQL): pg_stat_statements, pgcrypto, uuid-ossp
- ✅ Functions en Triggers
- ✅ Manual backups en restores

---

## 🎯 Render Performance Tips

### 1. Enable pg_stat_statements

```sql
-- Check of al enabled
SELECT * FROM pg_available_extensions WHERE name = 'pg_stat_statements';

-- Enable als nog niet actief
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

-- Check slow queries
SELECT 
    substring(query, 1, 100) as query,
    calls,
    ROUND(total_exec_time::numeric, 2) as total_ms,
    ROUND(mean_exec_time::numeric, 2) as avg_ms
FROM pg_stat_statements
WHERE mean_exec_time > 100  -- Slower than 100ms
ORDER BY mean_exec_time DESC
LIMIT 20;
```

### 2. Connection Pooling in App

Render databases hebben connectie limieten. Zorg dat je app correct pooling gebruikt:

```go
// In config/database.go
db.SetMaxOpenConns(20)        // Max 20 connections (Render limit afhankelijk van plan)
db.SetMaxIdleConns(5)         // 5 idle connections
db.SetConnMaxLifetime(5 * time.Minute)
db.SetConnMaxIdleTime(5 * time.Minute)
```

### 3. Monitor via Render Dashboard

Ga naar Render Dashboard → Database → Metrics om te monitoren:
- CPU Usage
- Memory Usage
- Connection Count
- Disk Usage

---

## 🔄 Backup & Restore op Render

### Backup

#### Via Render Dashboard (Aanbevolen)
1. Ga naar Database in Render Dashboard
2. Klik op "Manual Backups" tab
3. Klik "Create Backup"
4. Backup is beschikbaar voor download

#### Via pg_dump
```bash
# Full backup
pg_dump "$DATABASE_URL" > backup_$(date +%Y%m%d).sql

# Compressed
pg_dump "$DATABASE_URL" | gzip > backup_$(date +%Y%m%d).sql.gz

# Schema only
pg_dump --schema-only "$DATABASE_URL" > schema_backup.sql
```

### Restore

⚠️ **DIT VERVANGT DE DATABASE!**

```bash
# Stop de applicatie in Render Dashboard eerst!

# Restore
psql "$DATABASE_URL" < backup_20251030.sql

# Of van compressed
gunzip -c backup_20251030.sql.gz | psql "$DATABASE_URL"

# Start applicatie weer
```

---

## 📈 Performance Monitoring Checklist

### Dagelijks (Automatisch via Render)
- [ ] Check Render Dashboard Metrics
- [ ] CPU < 70%
- [ ] Memory < 80%
- [ ] Disk < 80%
- [ ] Connection count < limit

### Wekelijks (Handmatig)
```sql
-- Dead tuples check
SELECT
    schemaname,
    tablename,
    n_dead_tup,
    ROUND(100 * n_dead_tup / NULLIF(n_live_tup + n_dead_tup, 0), 2) AS dead_ratio
FROM pg_stat_user_tables
WHERE n_dead_tup > 1000
ORDER BY n_dead_tup DESC
LIMIT 10;

-- Run VACUUM if dead_ratio > 10%
VACUUM ANALYZE;
```

### Maandelijks
- [ ] Review slow queries (pg_stat_statements)
- [ ] Run data cleanup script
- [ ] Check index usage
- [ ] Manual backup via Render Dashboard
- [ ] Review disk usage trends

---

## 🚨 Troubleshooting

### Problem: Connection Timeout

```bash
# Check connection
psql "$DATABASE_URL" -c "SELECT version();"

# Check connection from Render app logs
# Look for: "Database connection failed" or "timeout"
```

**Solutions:**
- Check Render database status in dashboard
- Verify connection string in environment variables
- Check IP allowlist (if configured)
- Increase connection timeout in app config

### Problem: Out of Connections

```sql
-- Check current connections
SELECT count(*) FROM pg_stat_activity;

-- Check max connections
SHOW max_connections;

-- Kill idle connections (careful!)
SELECT pg_terminate_backend(pid) 
FROM pg_stat_activity 
WHERE state = 'idle' 
  AND state_change < NOW() - INTERVAL '10 minutes'
  AND pid != pg_backend_pid();
```

**Solutions:**
- Reduce `MaxOpenConns` in app
- Implement connection pooling (PgBouncer)
- Upgrade Render plan for more connections

### Problem: Slow Queries After Migration

```sql
-- Reset query stats
SELECT pg_stat_statements_reset();

-- Wait 1 hour, then check again
SELECT 
    query,
    calls,
    mean_exec_time
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;
```

**Solutions:**
- Run `ANALYZE;` to update statistics
- Check if indexes are being used with `EXPLAIN`
- Verify V1_47 indexes were created successfully

---

## 📞 Render Support Resources

- [Render PostgreSQL Docs](https://render.com/docs/databases)
- [Render Status Page](https://status.render.com/)
- Render Dashboard: https://dashboard.render.com/

---

## 🎯 Quick Command Reference

```bash
# Set environment variable for convenience
export DATABASE_URL="postgresql://dekoninklijkeloopdatabase_user:I4QP3JwyCcEbn8tGl6k3ErEvjUZ9V5rB@dpg-cva4c01c1ekc738q6q0g-a.oregon-postgres.render.com/dekoninklijkeloopdatabase"

# Connect
psql "$DATABASE_URL"

# Run query
psql "$DATABASE_URL" -c "SELECT version();"

# Run script
psql "$DATABASE_URL" -f script.sql

# Backup
pg_dump "$DATABASE_URL" > backup.sql

# VACUUM
psql "$DATABASE_URL" -c "VACUUM ANALYZE;"

# Check migrations
psql "$DATABASE_URL" -c "SELECT versie, naam FROM migraties ORDER BY toegepast DESC LIMIT 5;"

# Check indexes
psql "$DATABASE_URL" -c "SELECT schemaname, tablename, indexname FROM pg_indexes WHERE schemaname = 'public' ORDER BY tablename;"
```

---

**Database**: Render PostgreSQL (Managed)  
**Laatst Bijgewerkt**: 30 oktober 2025  
**Optimalisatie**: V1_47 Ready  
**Status**: Production ✅