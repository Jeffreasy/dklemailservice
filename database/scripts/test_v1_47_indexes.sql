-- Test Script: Verify V1_47 Migration Success
-- Usage: psql "$DATABASE_URL" -f database/scripts/test_v1_47_indexes.sql

\echo '============================================'
\echo 'V1_47 MIGRATION VERIFICATION TEST'
\echo '============================================'
\echo ''

-- Check if V1_47 migration was applied
\echo '1. Checking V1_47 migration status...'
SELECT 
    versie, 
    naam, 
    toegepast::text 
FROM migraties 
WHERE versie = '1.47.0';

\echo ''
\echo '2. Checking new indexes...'
\echo ''

-- Check critical FK indexes
\echo 'FK Indexes:'
SELECT 
    CASE 
        WHEN EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_gebruikers_role_id') 
        THEN '✅ idx_gebruikers_role_id'
        ELSE '❌ idx_gebruikers_role_id - MISSING'
    END as status
UNION ALL
SELECT 
    CASE 
        WHEN EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_aanmeldingen_gebruiker_id') 
        THEN '✅ idx_aanmeldingen_gebruiker_id'
        ELSE '❌ idx_aanmeldingen_gebruiker_id - MISSING'
    END
UNION ALL
SELECT 
    CASE 
        WHEN EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_verzonden_emails_contact_id') 
        THEN '✅ idx_verzonden_emails_contact_id'
        ELSE '❌ idx_verzonden_emails_contact_id - MISSING'
    END
UNION ALL
SELECT 
    CASE 
        WHEN EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_verzonden_emails_aanmelding_id') 
        THEN '✅ idx_verzonden_emails_aanmelding_id'
        ELSE '❌ idx_verzonden_emails_aanmelding_id - MISSING'
    END
UNION ALL
SELECT 
    CASE 
        WHEN EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_contact_antwoorden_contact_id') 
        THEN '✅ idx_contact_antwoorden_contact_id'
        ELSE '❌ idx_contact_antwoorden_contact_id - MISSING'
    END;

\echo ''
\echo 'Compound Indexes:'
SELECT 
    CASE 
        WHEN EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_contact_formulieren_status_created') 
        THEN '✅ idx_contact_formulieren_status_created'
        ELSE '❌ idx_contact_formulieren_status_created - MISSING'
    END as status
UNION ALL
SELECT 
    CASE 
        WHEN EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_aanmeldingen_status_created') 
        THEN '✅ idx_aanmeldingen_status_created'
        ELSE '❌ idx_aanmeldingen_status_created - MISSING'
    END;

\echo ''
\echo 'Full-Text Search Indexes:'
SELECT 
    CASE 
        WHEN EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_contact_formulieren_fts') 
        THEN '✅ idx_contact_formulieren_fts'
        ELSE '❌ idx_contact_formulieren_fts - MISSING'
    END as status
UNION ALL
SELECT 
    CASE 
        WHEN EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_aanmeldingen_fts') 
        THEN '✅ idx_aanmeldingen_fts'
        ELSE '❌ idx_aanmeldingen_fts - MISSING'
    END
UNION ALL
SELECT 
    CASE 
        WHEN EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_chat_messages_fts') 
        THEN '✅ idx_chat_messages_fts'
        ELSE '❌ idx_chat_messages_fts - MISSING'
    END;

\echo ''
\echo '3. Index count per table:'
SELECT 
    tablename, 
    COUNT(*) as index_count 
FROM pg_indexes 
WHERE schemaname = 'public' 
GROUP BY tablename 
HAVING COUNT(*) > 2
ORDER BY index_count DESC
LIMIT 10;

\echo ''
\echo '4. Top 5 largest tables:'
SELECT 
    tablename,
    pg_size_pretty(pg_total_relation_size('public.'||tablename)) AS size,
    pg_total_relation_size('public.'||tablename) as bytes
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY bytes DESC
LIMIT 5;

\echo ''
\echo '5. Recent migrations:'
SELECT 
    versie,
    naam,
    toegepast::text
FROM migraties
ORDER BY toegepast DESC
LIMIT 5;

\echo ''
\echo '============================================'
\echo 'TEST COMPLETED'
\echo '============================================'
\echo ''
\echo 'If V1_47 is found and all indexes exist:'
\echo 'Run: psql "$DATABASE_URL" -c "ANALYZE;"'
\echo ''
\echo 'Then test performance with queries from DATABASE_QUICK_REFERENCE.md'