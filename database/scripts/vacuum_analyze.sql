-- Vacuum Analyze Script
-- Run this weekly for optimal database performance
-- Usage: docker exec dkl-postgres psql -U postgres -d dklemailservice -f /path/to/vacuum_analyze.sql

-- ============================================
-- SECTION 1: FULL VACUUM ANALYZE
-- ============================================
-- This updates table statistics and reclaims space

\echo 'Starting VACUUM ANALYZE on all tables...'

-- High-traffic tables (prioritize these)
\echo 'Vacuuming verzonden_emails...'
VACUUM ANALYZE verzonden_emails;

\echo 'Vacuuming chat_messages...'
VACUUM ANALYZE chat_messages;

\echo 'Vacuuming contact_formulieren...'
VACUUM ANALYZE contact_formulieren;

\echo 'Vacuuming aanmeldingen...'
VACUUM ANALYZE aanmeldingen;

\echo 'Vacuuming incoming_emails...'
VACUUM ANALYZE incoming_emails;

\echo 'Vacuuming gebruikers...'
VACUUM ANALYZE gebruikers;

\echo 'Vacuuming refresh_tokens...'
VACUUM ANALYZE refresh_tokens;

\echo 'Vacuuming chat_channel_participants...'
VACUUM ANALYZE chat_channel_participants;

-- All other tables
\echo 'Vacuuming remaining tables...'
VACUUM ANALYZE;

\echo 'VACUUM ANALYZE completed successfully!'

-- ============================================
-- SECTION 2: STATISTICS REPORT
-- ============================================

\echo ''
\echo '=== DATABASE STATISTICS REPORT ==='
\echo ''

-- Table sizes
\echo '=== TOP 10 LARGEST TABLES ==='
SELECT 
    tablename,
    pg_size_pretty(pg_total_relation_size('public.'||tablename)) AS total_size,
    pg_size_pretty(pg_relation_size('public.'||tablename)) AS table_size,
    pg_size_pretty(pg_total_relation_size('public.'||tablename) - pg_relation_size('public.'||tablename)) AS index_size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size('public.'||tablename) DESC
LIMIT 10;

\echo ''
\echo '=== DEAD TUPLES (BLOAT) ==='
SELECT
    schemaname,
    tablename,
    n_live_tup AS live_rows,
    n_dead_tup AS dead_rows,
    ROUND(100 * n_dead_tup / NULLIF(n_live_tup + n_dead_tup, 0), 2) AS dead_ratio
FROM pg_stat_user_tables
WHERE n_dead_tup > 100
ORDER BY n_dead_tup DESC
LIMIT 10;

\echo ''
\echo '=== INDEX USAGE ==='
SELECT
    schemaname,
    tablename,
    indexname,
    idx_scan AS scans,
    pg_size_pretty(pg_relation_size(indexrelid)) AS size
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
ORDER BY idx_scan DESC
LIMIT 15;

\echo ''
\echo '=== POTENTIALLY UNUSED INDEXES ==='
SELECT
    schemaname,
    tablename,
    indexname,
    idx_scan AS scans,
    pg_size_pretty(pg_relation_size(indexrelid)) AS size
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
    AND idx_scan < 10
    AND indexrelname NOT LIKE '%_pkey'
ORDER BY pg_relation_size(indexrelid) DESC
LIMIT 10;

\echo ''
\echo 'Report completed!'