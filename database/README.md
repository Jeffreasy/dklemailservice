# DKL Email Service - Database Optimalisatie Pakket

Complete database analyse, optimalisaties en onderhoudsscripts voor de DKL Email Service PostgreSQL database.

---

## 📚 Documentatie Overzicht

| Document | Beschrijving | Status |
|----------|-------------|---------|
| **[DATABASE_ANALYSIS.md](../docs/DATABASE_ANALYSIS.md)** | Complete database analyse met 33 tabellen, indexes en optimalisatie aanbevelingen | ✅ Complete |
| **[DATABASE_QUICK_REFERENCE.md](../docs/DATABASE_QUICK_REFERENCE.md)** | Snelle referentie voor dagelijks database beheer | ✅ Complete |
| **[POSTGRESQL_CONFIGURATION.md](../docs/POSTGRESQL_CONFIGURATION.md)** | PostgreSQL configuratie optimalisaties | ✅ Complete |

---

## 🚀 Snelstart: Optimalisaties Toepassen

### Stap 1: Herstart Applicatie (Automatische Migratie)

De optimalisaties worden automatisch toegepast bij herstart:

```bash
cd /path/to/dklemailservice
docker-compose -f docker-compose.dev.yml restart app
```

**Wat gebeurt er:**
- V1_47 migratie wordt automatisch uitgevoerd
- 30+ nieuwe indexes worden aangemaakt
- Foreign key indexes
- Compound indexes voor dashboard queries
- Full-text search indexes
- Partial indexes voor gefilterde queries

### Stap 2: Verify Migratie Status

```bash
docker logs dkl-email-service --tail 50 | grep -i "migratie"
```

Zoek naar: `"Migratie succesvol uitgevoerd","file":"V1_47__performance_optimizations.sql"`

### Stap 3: Update Database Statistieken

```bash
docker exec dkl-postgres psql -U postgres -d dklemailservice -c "ANALYZE;"
```

---

## 📋 Beschikbare Scripts

### Maintenance Scripts

| Script | Doel | Frequentie | Commando |
|--------|------|-----------|----------|
| **[vacuum_analyze.sql](scripts/vacuum_analyze.sql)** | VACUUM ANALYZE + statistieken rapportage | Wekelijks | Zie hieronder |
| **[data_cleanup.sql](scripts/data_cleanup.sql)** | Cleanup oude data en archivering | Maandelijks | Zie hieronder |
| **[setup_partitioning.sql](scripts/setup_partitioning.sql)** | Table partitioning voor grote tabellen | Eenmalig | ⚠️ Downtime vereist |

### Vacuum Analyze (Wekelijks)

```bash
# Via Docker
docker exec dkl-postgres psql -U postgres -d dklemailservice < database/scripts/vacuum_analyze.sql

# Of via psql
docker exec -i dkl-postgres psql -U postgres -d dklemailservice < database/scripts/vacuum_analyze.sql
```

### Data Cleanup (Maandelijks)

⚠️ **BELANGRIJK: Maak eerst een backup!**

```bash
# Backup maken
docker exec dkl-postgres pg_dump -U postgres dklemailservice > backup_$(date +%Y%m%d).sql

# Cleanup uitvoeren
docker exec -i dkl-postgres psql -U postgres -d dklemailservice < database/scripts/data_cleanup.sql
```

### Table Partitioning (Optioneel)

⚠️ **WAARSCHUWING: Vereist downtime!**

Dit script partitioneert grote tabellen (verzonden_emails, chat_messages) voor betere performance:

```bash
# Stop applicatie
docker-compose -f docker-compose.dev.yml stop app

# Backup maken
docker exec dkl-postgres pg_dump -U postgres dklemailservice > backup_before_partitioning.sql

# Partitioning toepassen
docker exec -i dkl-postgres psql -U postgres -d dklemailservice < database/scripts/setup_partitioning.sql

# Start applicatie
docker-compose -f docker-compose.dev.yml start app
```

---

## 📊 Migraties Overzicht

### V1_47: Performance Optimizations

**Status**: ✅ Klaar voor deployment  
**Impact**: HIGH - Significante performance verbetering  
**Downtime**: Geen (indexes worden online gebouwd)

**Toegevoegde Indexes:**

#### Foreign Key Indexes (Kritiek!)
- `idx_gebruikers_role_id` - RBAC lookups
- `idx_aanmeldingen_gebruiker_id` - User registrations
- `idx_verzonden_emails_contact_id` - Email tracking
- `idx_verzonden_emails_aanmelding_id` - Email tracking
- `idx_verzonden_emails_template_id` - Template usage
- `idx_contact_antwoorden_contact_id` - Responses
- `idx_aanmelding_antwoorden_aanmelding_id` - Responses

#### Compound Indexes (Dashboard Performance)
- `idx_contact_formulieren_status_created` - Unanswered forms
- `idx_aanmeldingen_status_created` - Registration dashboard
- `idx_verzonden_emails_status_tijd` - Email status tracking

#### Full-Text Search Indexes
- `idx_contact_formulieren_fts` - Contact form zoeken
- `idx_aanmeldingen_fts` - Registration zoeken  
- `idx_chat_messages_fts` - Chat message zoeken

#### Partial Indexes (Filtered Queries)
- `idx_verzonden_emails_errors` - Failed emails only
- `idx_incoming_emails_processing` - Unprocessed queue
- `idx_contact_formulieren_nieuw` - New submissions
- `idx_chat_participants_active` - Active participants

**Verwachte Performance Verbetering:**
- JOIN operations: **50-90% sneller**
- Dashboard queries: **60-80% sneller**
- Search operations: **90-95% sneller**
- Filtered queries: **70-85% sneller**

---

## 🎯 Database Schema Overzicht

### Domein Verdeling

| Domein | Tabellen | Hoofdfunctie |
|--------|----------|--------------|
| **Core Email & Users** | 9 | Gebruikers, contact forms, aanmeldingen, emails |
| **Chat System** | 5 | Real-time chat, channels, messages, reactions |
| **RBAC** | 4 | Roles, permissions, user-role mapping |
| **Authentication** | 1 | Refresh tokens voor JWT |
| **Content Management** | 14 | Photos, videos, albums, sponsors, partners |

**Totaal**: 33 tabellen, 46 migraties

### Kritieke Tabellen (Require Monitoring)

| Tabel | Type | Growth Rate | Monitoring |
|-------|------|-------------|------------|
| `verzonden_emails` | High Traffic | Hoog | Dagelijks |
| `chat_messages` | High Traffic | Hoog | Dagelijks |
| `incoming_emails` | Moderate | Medium | Wekelijks |
| `contact_formulieren` | Moderate | Medium | Wekelijks |
| `aanmeldingen` | Seasonal | Low-Medium | Maandelijks |

---

## 🔍 Monitoring & Health Checks

### Dagelijkse Checks

```sql
-- Tabel groottes (top 5)
SELECT 
    tablename,
    pg_size_pretty(pg_total_relation_size('public.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size('public.'||tablename) DESC
LIMIT 5;

-- Failed emails (laatste 24 uur)
SELECT COUNT(*) 
FROM verzonden_emails 
WHERE status = 'failed' 
  AND verzonden_op > NOW() - INTERVAL '24 hours';

-- Unprocessed incoming emails
SELECT COUNT(*) 
FROM incoming_emails 
WHERE is_processed = FALSE;
```

### Wekelijkse Checks

```sql
-- Dead tuples (bloat)
SELECT
    schemaname,
    tablename,
    n_dead_tup,
    ROUND(100 * n_dead_tup / NULLIF(n_live_tup + n_dead_tup, 0), 2) AS dead_ratio
FROM pg_stat_user_tables
WHERE n_dead_tup > 1000
ORDER BY n_dead_tup DESC
LIMIT 10;

-- Slow queries (requires pg_stat_statements)
SELECT 
    substring(query, 1, 100) as query,
    calls,
    ROUND(mean_exec_time::numeric, 2) as avg_ms
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;
```

---

## 🛠️ Troubleshooting

### Problem: Migratie Niet Toegepast

```bash
# Check migratie status
docker logs dkl-email-service --tail 100 | grep V1_47

# Herstart applicatie
docker-compose -f docker-compose.dev.yml restart app

# Check database migraties tabel
docker exec dkl-postgres psql -U postgres -d dklemailservice -c "SELECT versie, naam FROM migraties ORDER BY toegepast DESC LIMIT 5;"
```

### Problem: Slow Queries

```sql
-- Identify slow queries
SELECT 
    query,
    calls,
    total_exec_time,
    mean_exec_time
FROM pg_stat_statements
WHERE mean_exec_time > 1000  -- > 1 second
ORDER BY mean_exec_time DESC;

-- Check missing indexes
SELECT 
    schemaname,
    tablename,
    seq_scan,
    idx_scan,
    seq_scan - idx_scan AS too_much_seq,
    pg_size_pretty(pg_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_stat_user_tables
WHERE seq_scan - idx_scan > 0
  AND pg_relation_size(schemaname||'.'||tablename) > 1000000
ORDER BY too_much_seq DESC;
```

### Problem: High Disk Usage

```bash
# Check disk usage
docker exec dkl-postgres df -h /var/lib/postgresql/data

# Identify large tables
docker exec dkl-postgres psql -U postgres -d dklemailservice -c "
SELECT 
    tablename,
    pg_size_pretty(pg_total_relation_size('public.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size('public.'||tablename) DESC
LIMIT 10;"

# Run cleanup if needed
docker exec -i dkl-postgres psql -U postgres -d dklemailservice < database/scripts/data_cleanup.sql
```

---

## 🔐 Backup & Restore

### Backup Strategieën

#### Daily Automated Backup (Cron Job)

```bash
#!/bin/bash
# Save as: /root/backup_dkl_db.sh
# Cron: 0 2 * * * /root/backup_dkl_db.sh

BACKUP_DIR="/backups/postgresql"
DATE=$(date +%Y%m%d_%H%M%S)
RETENTION_DAYS=30

# Create backup
docker exec dkl-postgres pg_dump -U postgres dklemailservice | gzip > "$BACKUP_DIR/dkl_db_$DATE.sql.gz"

# Remove old backups
find "$BACKUP_DIR" -name "dkl_db_*.sql.gz" -mtime +$RETENTION_DAYS -delete

echo "Backup completed: dkl_db_$DATE.sql.gz"
```

#### Manual Backup

```bash
# Full database
docker exec dkl-postgres pg_dump -U postgres dklemailservice > backup_$(date +%Y%m%d).sql

# Compressed
docker exec dkl-postgres pg_dump -U postgres dklemailservice | gzip > backup_$(date +%Y%m%d).sql.gz

# Schema only
docker exec dkl-postgres pg_dump -U postgres --schema-only dklemailservice > schema_backup.sql
```

#### Restore

```bash
# Stop applicatie
docker-compose -f docker-compose.dev.yml stop app

# Restore
docker exec -i dkl-postgres psql -U postgres -d dklemailservice < backup_20251030.sql

# Of van compressed backup
gunzip -c backup_20251030.sql.gz | docker exec -i dkl-postgres psql -U postgres -d dklemailservice

# Start applicatie
docker-compose -f docker-compose.dev.yml start app
```

---

## 📈 Performance Benchmarking

### Baseline Metrics (Voor Optimalisaties)

Test deze queries voor en na optimalisaties:

```sql
-- Query 1: Dashboard contact forms
EXPLAIN ANALYZE
SELECT * FROM contact_formulieren 
WHERE status = 'nieuw' AND beantwoord = FALSE 
ORDER BY created_at DESC 
LIMIT 20;

-- Query 2: Recent emails per status
EXPLAIN ANALYZE
SELECT status, COUNT(*) 
FROM verzonden_emails 
WHERE verzonden_op > NOW() - INTERVAL '7 days'
GROUP BY status;

-- Query 3: User permissions lookup
EXPLAIN ANALYZE
SELECT * FROM user_permissions 
WHERE user_id = 'some-uuid-here';

-- Query 4: Chat messages in channel
EXPLAIN ANALYZE
SELECT * FROM chat_messages 
WHERE channel_id = 'some-uuid-here' 
ORDER BY created_at DESC 
LIMIT 50;
```

### Expected Results After V1_47

| Query | Before | After | Improvement |
|-------|--------|-------|-------------|
| Dashboard forms | ~50ms | ~5ms | 90% faster |
| Email stats | ~100ms | ~20ms | 80% faster |
| Permission lookup | ~30ms | ~3ms | 90% faster |
| Chat messages | ~40ms | ~5ms | 87% faster |

---

## 🎓 Best Practices

### DO's ✅

- ✅ Run VACUUM ANALYZE wekelijks
- ✅ Monitor tabel groottes dagelijks
- ✅ Backup daily (automated)
- ✅ Test backups maandelijks
- ✅ Review slow queries wekelijks
- ✅ Update statistieken na bulk inserts
- ✅ Archive oude data maandelijks

### DON'Ts ❌

- ❌ DROP indexes zonder analyse
- ❌ Run cleanup zonder backup
- ❌ Modify schema tijdens peak hours
- ❌ Ignore dead tuple warnings
- ❌ Skip migration testing
- ❌ Forget to ANALYZE after big changes

---

## 📞 Support & Contact

Voor vragen over database optimalisaties:

1. Check [DATABASE_ANALYSIS.md](../docs/DATABASE_ANALYSIS.md) voor details
2. Review [DATABASE_QUICK_REFERENCE.md](../docs/DATABASE_QUICK_REFERENCE.md) voor commando's
3. Consult [POSTGRESQL_CONFIGURATION.md](../docs/POSTGRESQL_CONFIGURATION.md) voor configuratie

---

**Database Versie**: PostgreSQL 15 Alpine  
**Laatst Bijgewerkt**: 30 oktober 2025  
**Optimalisatie Versie**: V1_47  
**Status**: Production Ready ✅