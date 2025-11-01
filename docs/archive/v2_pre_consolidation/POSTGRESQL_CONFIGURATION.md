# PostgreSQL Configuration Optimization

Optimalisatie instellingen voor PostgreSQL 15 in Docker omgeving.

---

## 📋 Huidige Configuratie

```yaml
# docker-compose.dev.yml
postgres:
  image: postgres:15-alpine
  environment:
    POSTGRES_USER: postgres
    POSTGRES_PASSWORD: postgres
    POSTGRES_DB: dklemailservice
```

---

## 🎯 Aanbevolen Configuratie Wijzigingen

### 1. Memory Settings

Voeg toe aan `docker-compose.dev.yml`:

```yaml
postgres:
  image: postgres:15-alpine
  environment:
    # ... existing env vars ...
  command: >
    postgres
    -c shared_buffers=256MB
    -c effective_cache_size=1GB
    -c maintenance_work_mem=128MB
    -c checkpoint_completion_target=0.9
    -c wal_buffers=16MB
    -c default_statistics_target=100
    -c random_page_cost=1.1
    -c effective_io_concurrency=200
    -c work_mem=4MB
    -c min_wal_size=1GB
    -c max_wal_size=4GB
```

### 2. Custom postgresql.conf (Aanbevolen voor Productie)

Maak `database/config/postgresql.conf`:

```conf
# PostgreSQL 15 Configuration for DKL Email Service
# Optimized for moderate load with 4GB RAM

# ============================================
# CONNECTIONS AND AUTHENTICATION
# ============================================
max_connections = 100
superuser_reserved_connections = 3

# ============================================
# MEMORY SETTINGS
# ============================================
shared_buffers = 1GB              # 25% of RAM
effective_cache_size = 3GB        # 75% of RAM
work_mem = 4MB                    # Per operation
maintenance_work_mem = 256MB      # For VACUUM, CREATE INDEX
wal_buffers = 16MB
temp_buffers = 8MB

# ============================================
# QUERY PLANNER
# ============================================
random_page_cost = 1.1            # SSD optimized (default 4.0 for HDD)
effective_io_concurrency = 200    # SSD parallel I/O
default_statistics_target = 100   # More stats for better plans

# ============================================
# WRITE AHEAD LOG (WAL)
# ============================================
wal_level = replica               # Enable replication support
min_wal_size = 1GB
max_wal_size = 4GB
checkpoint_completion_target = 0.9
checkpoint_timeout = 10min
archive_mode = off                # Enable for PITR backups

# ============================================
# REPLICATION (For Future Scaling)
# ============================================
max_wal_senders = 3
max_replication_slots = 3
hot_standby = on

# ============================================
# LOGGING
# ============================================
logging_collector = on
log_directory = 'log'
log_filename = 'postgresql-%Y-%m-%d_%H%M%S.log'
log_rotation_age = 1d
log_rotation_size = 100MB
log_min_duration_statement = 1000  # Log queries > 1 second
log_checkpoints = on
log_connections = on
log_disconnections = on
log_duration = off
log_line_prefix = '%t [%p]: [%l-1] user=%u,db=%d,app=%a,client=%h '
log_statement = 'ddl'              # Log DDL statements
log_temp_files = 0                 # Log all temp files

# ============================================
# AUTOVACUUM
# ============================================
autovacuum = on
autovacuum_max_workers = 3
autovacuum_naptime = 1min
autovacuum_vacuum_threshold = 50
autovacuum_analyze_threshold = 50
autovacuum_vacuum_scale_factor = 0.1
autovacuum_analyze_scale_factor = 0.05
autovacuum_vacuum_cost_delay = 10ms
autovacuum_vacuum_cost_limit = 200

# ============================================
# QUERY PERFORMANCE
# ============================================
# Enable query statistics extension
shared_preload_libraries = 'pg_stat_statements'
pg_stat_statements.max = 10000
pg_stat_statements.track = all

# ============================================
# PERFORMANCE MONITORING
# ============================================
track_activities = on
track_counts = on
track_io_timing = on
track_functions = all

# ============================================
# CLIENT CONNECTION DEFAULTS
# ============================================
timezone = 'Europe/Amsterdam'
lc_messages = 'en_US.UTF-8'
lc_monetary = 'en_US.UTF-8'
lc_numeric = 'en_US.UTF-8'
lc_time = 'en_US.UTF-8'
default_text_search_config = 'pg_catalog.dutch'
```

### Mount deze configuratie in Docker:

```yaml
postgres:
  image: postgres:15-alpine
  volumes:
    - postgres_data:/var/lib/postgresql/data
    - ./database/config/postgresql.conf:/etc/postgresql/postgresql.conf:ro
  command: postgres -c config_file=/etc/postgresql/postgresql.conf
```

---

## 🔧 Performance Tuning Calculator

Gebruik [PGTune](https://pgtune.leopard.in.ua/) voor aangepaste configuraties:

**Input:**
- DB Version: 15
- OS Type: Linux
- DB Type: Web application
- Total Memory: 4GB (pas aan naar jouw server)
- CPUs: 2 (pas aan naar jouw server)
- Connections: 100
- Data Storage: SSD

---

## 📊 Monitoring Queries

### Check Current Configuration

```sql
-- Show all non-default settings
SELECT name, setting, unit, source 
FROM pg_settings 
WHERE source != 'default' 
ORDER BY name;

-- Check shared_buffers
SHOW shared_buffers;

-- Check effective_cache_size
SHOW effective_cache_size;

-- Check work_mem
SHOW work_mem;
```

### Enable pg_stat_statements

```sql
-- Add to postgresql.conf or docker command
shared_preload_libraries = 'pg_stat_statements'

-- Restart PostgreSQL
-- Then create extension:
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

-- View slow queries
SELECT 
    substring(query, 1, 100) as query_snippet,
    calls,
    total_exec_time,
    mean_exec_time,
    max_exec_time
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 20;
```

---

## 🚀 Connection Pooling (Application Level)

De Go applicatie gebruikt database/sql pool:

```go
// Aanbevolen settings in config/database.go
db.SetMaxOpenConns(25)           // Max connections
db.SetMaxIdleConns(5)            // Idle connections
db.SetConnMaxLifetime(5 * time.Minute)
db.SetConnMaxIdleTime(5 * time.Minute)
```

### PgBouncer (Optioneel - Voor Hoge Load)

Voeg toe aan `docker-compose.dev.yml`:

```yaml
pgbouncer:
  image: pgbouncer/pgbouncer:latest
  environment:
    DATABASES_HOST: postgres
    DATABASES_PORT: 5432
    DATABASES_USER: postgres
    DATABASES_PASSWORD: postgres
    DATABASES_DBNAME: dklemailservice
    PGBOUNCER_POOL_MODE: transaction
    PGBOUNCER_MAX_CLIENT_CONN: 1000
    PGBOUNCER_DEFAULT_POOL_SIZE: 20
  ports:
    - "6432:6432"
  depends_on:
    - postgres
```

Update applicatie connection string:
```
DB_HOST=pgbouncer
DB_PORT=6432
```

---

## 🔍 Performance Verification

Na het toepassen van configuratie wijzigingen:

### 1. Check Settings
```bash
docker exec dkl-postgres psql -U postgres -d dklemailservice -c "SHOW ALL;"
```

### 2. Restart with New Config
```bash
docker-compose -f docker-compose.dev.yml restart postgres
docker logs dkl-postgres --tail 50
```

### 3. Verify Performance
```sql
-- Check buffer hit ratio (should be > 99%)
SELECT 
    sum(heap_blks_hit) / (sum(heap_blks_hit) + sum(heap_blks_read)) * 100 
    AS cache_hit_ratio
FROM pg_statio_user_tables;

-- Check index usage
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes
ORDER BY idx_scan DESC;
```

---

## 🎯 Optimization Checklist

### Immediate Actions
- [ ] Update `docker-compose.dev.yml` with memory settings
- [ ] Enable `pg_stat_statements` extension
- [ ] Set appropriate connection limits
- [ ] Configure autovacuum aggressively

### Short-term (1 week)
- [ ] Monitor slow queries via `pg_stat_statements`
- [ ] Create custom `postgresql.conf`
- [ ] Implement connection pooling in app
- [ ] Set up query performance monitoring

### Long-term (1 month)
- [ ] Consider PgBouncer for high load
- [ ] Implement read replicas if needed
- [ ] Set up automated backups
- [ ] Review and tune based on actual usage patterns

---

## 📚 Resources

- [PostgreSQL 15 Documentation](https://www.postgresql.org/docs/15/)
- [PGTune Configuration Generator](https://pgtune.leopard.in.ua/)
- [Postgres Performance Optimization](https://wiki.postgresql.org/wiki/Performance_Optimization)
- [Understanding PostgreSQL Configuration](https://postgresqlco.nf/doc/en/param/)

---

**Laatst bijgewerkt**: 30 oktober 2025  
**Versie**: 1.0  
**Review**: Na toepassing en 1 week monitoring