# DKL Email Service - Database Quick Reference

⚡ **Snelle referentie voor dagelijks database beheer**

---

## 🚀 Snelstart Commando's

### Database Status Checken
```bash
# Via MCP Server
use_mcp_tool dkl-email-service service_status

# Direct via Docker
docker exec dkl-postgres psql -U postgres -d dklemailservice -c "SELECT version();"
```

### Migraties Draaien
```bash
# Automatisch via applicatie startup
docker-compose -f docker-compose.dev.yml up -d

# Handmatig nieuwe migratie toevoegen
# 1. Maak bestand: database/migrations/V1_XX__description.sql
# 2. Herstart de applicatie
docker-compose -f docker-compose.dev.yml restart app
```

---

## 📊 Monitoring Queries

### Tabel Groottes
```sql
SELECT 
    tablename,
    pg_size_pretty(pg_total_relation_size('public.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size('public.'||tablename) DESC
LIMIT 10;
```

### Index Gebruik
```sql
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan as scans,
    pg_size_pretty(pg_relation_size(indexrelid)) as size
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
ORDER BY idx_scan ASC
LIMIT 20;
```

### Langzame Queries (vereist pg_stat_statements)
```sql
SELECT 
    substring(query, 1, 100) as query_snippet,
    calls,
    ROUND(total_exec_time::numeric, 2) as total_ms,
    ROUND(mean_exec_time::numeric, 2) as avg_ms,
    ROUND(max_exec_time::numeric, 2) as max_ms
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 20;
```

### Actieve Verbindingen
```sql
SELECT 
    datname,
    count(*) as connections,
    state
FROM pg_stat_activity
WHERE datname = 'dklemailservice'
GROUP BY datname, state;
```

---

## 🔧 Maintenance Taken

### Dagelijks
```sql
-- Niet nodig - auto-vacuum regelt dit
```

### Wekelijks
```sql
-- Vacuuming voor performance
VACUUM ANALYZE;
```

### Maandelijks
```sql
-- Reindex hele database
REINDEX DATABASE dklemailservice;

-- Of specifieke grote tabellen
REINDEX TABLE verzonden_emails;
REINDEX TABLE chat_messages;
```

### Data Cleanup
```sql
-- Verwijder verlopen refresh tokens (ouder dan 30 dagen)
DELETE FROM refresh_tokens 
WHERE expires_at < NOW() - INTERVAL '30 days';

-- Archiveer oude verzonden emails (ouder dan 1 jaar)
-- OPGELET: Maak eerst backup!
DELETE FROM verzonden_emails 
WHERE verzonden_op < NOW() - INTERVAL '1 year';
```

---

## 🔍 Veelgebruikte Queries

### Gebruikers
```sql
-- Actieve gebruikers met rollen
SELECT 
    g.id,
    g.naam,
    g.email,
    r.name as role,
    g.laatste_login
FROM gebruikers g
LEFT JOIN roles r ON g.role_id = r.id
WHERE g.is_actief = TRUE
ORDER BY g.laatste_login DESC NULLS LAST;

-- Gebruikers met meerdere rollen (via RBAC)
SELECT 
    g.naam,
    g.email,
    array_agg(r.name) as roles
FROM gebruikers g
JOIN user_roles ur ON g.id = ur.user_id
JOIN roles r ON ur.role_id = r.id
WHERE ur.is_active = TRUE
GROUP BY g.id, g.naam, g.email;
```

### Contact Formulieren
```sql
-- Onbeantwoorde contact formulieren
SELECT 
    id,
    naam,
    email,
    LEFT(bericht, 100) as bericht_preview,
    status,
    created_at
FROM contact_formulieren
WHERE beantwoord = FALSE
ORDER BY created_at DESC
LIMIT 20;

-- Contact formulieren per status
SELECT 
    status,
    COUNT(*) as aantal,
    COUNT(*) FILTER (WHERE beantwoord = TRUE) as beantwoord
FROM contact_formulieren
GROUP BY status
ORDER BY aantal DESC;
```

### Aanmeldingen
```sql
-- Recente aanmeldingen per afstand
SELECT 
    afstand,
    COUNT(*) as aantal,
    AVG(steps) as gem_steps
FROM aanmeldingen
WHERE created_at > NOW() - INTERVAL '30 days'
GROUP BY afstand
ORDER BY aantal DESC;

-- Top deelnemers (meeste steps)
SELECT 
    naam,
    email,
    afstand,
    steps
FROM aanmeldingen
WHERE steps > 0
ORDER BY steps DESC
LIMIT 10;
```

### Email Tracking
```sql
-- Email verzend statistieken (laatste 7 dagen)
SELECT 
    DATE(verzonden_op) as datum,
    status,
    COUNT(*) as aantal
FROM verzonden_emails
WHERE verzonden_op > NOW() - INTERVAL '7 days'
GROUP BY DATE(verzonden_op), status
ORDER BY datum DESC, status;

-- Failed emails met details
SELECT 
    id,
    ontvanger,
    onderwerp,
    status,
    fout_bericht,
    verzonden_op
FROM verzonden_emails
WHERE status = 'failed'
ORDER BY verzonden_op DESC
LIMIT 20;
```

### Chat Statistieken
```sql
-- Meest actieve channels
SELECT 
    c.name,
    c.type,
    COUNT(m.id) as berichten_aantal,
    COUNT(DISTINCT m.user_id) as unieke_gebruikers
FROM chat_channels c
LEFT JOIN chat_messages m ON c.id = m.channel_id
WHERE c.is_active = TRUE
GROUP BY c.id, c.name, c.type
ORDER BY berichten_aantal DESC
LIMIT 10;

-- Online gebruikers
SELECT 
    g.naam,
    g.email,
    up.status,
    up.last_seen
FROM chat_user_presence up
JOIN gebruikers g ON up.user_id = g.id
WHERE up.status IN ('online', 'away')
ORDER BY up.last_seen DESC;
```

---

## 🛡️ Security Checks

### Gebruikers zonder rollen
```sql
SELECT 
    id,
    naam,
    email,
    rol as legacy_rol,
    role_id
FROM gebruikers
WHERE role_id IS NULL
  AND is_actief = TRUE;
```

### Verouderde refresh tokens
```sql
SELECT 
    COUNT(*) as verlopen_tokens
FROM refresh_tokens
WHERE expires_at < NOW()
  AND is_revoked = FALSE;
```

### Inactieve gebruikers (niet ingelogd in 90 dagen)
```sql
SELECT 
    naam,
    email,
    laatste_login,
    AGE(NOW(), laatste_login) as inactief_sinds
FROM gebruikers
WHERE laatste_login < NOW() - INTERVAL '90 days'
  AND is_actief = TRUE
ORDER BY laatste_login ASC;
```

---

## 🔥 Performance Checks

### Tabellen zonder indexes op foreign keys
```sql
-- Dit zou leeg moeten zijn na V1_47 migratie
SELECT
    tc.table_name,
    kcu.column_name,
    ccu.table_name AS foreign_table_name,
    ccu.column_name AS foreign_column_name
FROM information_schema.table_constraints AS tc
JOIN information_schema.key_column_usage AS kcu
    ON tc.constraint_name = kcu.constraint_name
JOIN information_schema.constraint_column_usage AS ccu
    ON ccu.constraint_name = tc.constraint_name
WHERE tc.constraint_type = 'FOREIGN KEY'
    AND tc.table_schema = 'public'
    AND NOT EXISTS (
        SELECT 1
        FROM pg_indexes
        WHERE schemaname = 'public'
            AND tablename = tc.table_name
            AND indexdef LIKE '%' || kcu.column_name || '%'
    );
```

### Ongebruikte indexes
```sql
SELECT
    schemaname,
    tablename,
    indexname,
    idx_scan,
    pg_size_pretty(pg_relation_size(indexrelid)) as size
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
    AND idx_scan = 0
    AND indexrelname NOT LIKE '%_pkey'
ORDER BY pg_relation_size(indexrelid) DESC;
```

### Bloated tables (te veel dode rows)
```sql
SELECT
    schemaname,
    tablename,
    n_dead_tup,
    n_live_tup,
    ROUND(100 * n_dead_tup / NULLIF(n_live_tup + n_dead_tup, 0), 2) as dead_ratio
FROM pg_stat_user_tables
WHERE n_dead_tup > 1000
ORDER BY n_dead_tup DESC;
```

---

## 🆘 Troubleshooting

### Database verbinding mislukt
```bash
# Check of container draait
docker ps | grep dkl-postgres

# Check logs
docker logs dkl-postgres --tail 50

# Herstart container
docker-compose -f docker-compose.dev.yml restart postgres

# Wacht op healthy status
docker-compose -f docker-compose.dev.yml ps postgres
```

### Migratie gefaald
```bash
# Check huidige migratie versie
docker exec dkl-postgres psql -U postgres -d dklemailservice -c "SELECT * FROM migraties ORDER BY toegepast DESC LIMIT 5;"

# Check applicatie logs
docker logs dkl-email-service --tail 100

# Rollback laatste migratie (VOORZICHTIG!)
# Dit moet handmatig per migratie uitgevoerd worden
```

### Performance problemen
```sql
-- Check current queries
SELECT 
    pid,
    now() - query_start as duration,
    state,
    LEFT(query, 100) as query_snippet
FROM pg_stat_activity
WHERE state != 'idle'
    AND query NOT LIKE '%pg_stat_activity%'
ORDER BY duration DESC;

-- Kill long-running query (VOORZICHTIG!)
-- SELECT pg_terminate_backend(PID);
```

### Disk space issues
```bash
# Check disk usage in container
docker exec dkl-postgres df -h /var/lib/postgresql/data

# Check largest tables
docker exec dkl-postgres psql -U postgres -d dklemailservice -c "
SELECT 
    tablename,
    pg_size_pretty(pg_total_relation_size('public.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size('public.'||tablename) DESC
LIMIT 5;"
```

---

## 📦 Backup & Restore

### Backup maken
```bash
# Full database dump
docker exec dkl-postgres pg_dump -U postgres dklemailservice > backup_$(date +%Y%m%d_%H%M%S).sql

# Alleen schema
docker exec dkl-postgres pg_dump -U postgres --schema-only dklemailservice > schema_backup.sql

# Alleen data
docker exec dkl-postgres pg_dump -U postgres --data-only dklemailservice > data_backup.sql

# Specifieke tabel
docker exec dkl-postgres pg_dump -U postgres -t gebruikers dklemailservice > gebruikers_backup.sql
```

### Restore
```bash
# WAARSCHUWING: Dit vervangt de huidige database!

# Stop applicatie
docker-compose -f docker-compose.dev.yml stop app

# Restore database
docker exec -i dkl-postgres psql -U postgres -d dklemailservice < backup_20251030.sql

# Start applicatie
docker-compose -f docker-compose.dev.yml start app
```

---

## 🎯 Optimalisatie Checklist

### Na V1_47 Migratie
- [ ] Run `ANALYZE;` om statistics bij te werken
- [ ] Check index sizes met monitoring query
- [ ] Test dashboard performance (contact formulieren, aanmeldingen)
- [ ] Monitor slow query log voor 24 uur
- [ ] Document performance improvement metrics

### Wekelijks
- [ ] Check tabel groottes
- [ ] Review failed emails
- [ ] Check ongebruikte indexes
- [ ] Run VACUUM ANALYZE

### Maandelijks
- [ ] Full database backup
- [ ] Review en cleanup oude data
- [ ] REINDEX grote tabellen
- [ ] Update documentation indien schema changes

---

## 📚 Handige Links

- [Volledige Database Analyse](./DATABASE_ANALYSIS.md)
- [PostgreSQL Official Docs](https://www.postgresql.org/docs/15/)
- [Performance Tips](https://wiki.postgresql.org/wiki/Performance_Optimization)
- [Index Strategies](https://www.postgresql.org/docs/15/indexes.html)

---

**Laatst bijgewerkt**: 30 oktober 2025  
**Versie**: 1.0