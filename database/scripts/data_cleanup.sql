-- Data Cleanup Script
-- Run this monthly to remove old/expired data
-- IMPORTANT: Always create a backup before running cleanup!
-- Usage: docker exec dkl-postgres psql -U postgres -d dklemailservice -f /path/to/data_cleanup.sql

-- ============================================
-- SAFETY CHECK: CREATE BACKUP FIRST!
-- ============================================
\echo 'WARNING: This script will permanently delete data!'
\echo 'Make sure you have created a backup before continuing.'
\echo 'Press Ctrl+C to cancel, or press Enter to continue...'
\prompt 'Continue? (yes/no): ' confirmation

-- ============================================
-- SECTION 1: EXPIRED REFRESH TOKENS
-- ============================================
\echo ''
\echo '=== Cleaning up expired refresh tokens ==='

-- Count before cleanup
SELECT COUNT(*) AS expired_tokens_count
FROM refresh_tokens
WHERE expires_at < NOW() - INTERVAL '30 days';

-- Delete expired tokens (older than 30 days past expiration)
DELETE FROM refresh_tokens
WHERE expires_at < NOW() - INTERVAL '30 days';

\echo 'Expired refresh tokens deleted.'

-- ============================================
-- SECTION 2: OLD VERZONDEN EMAILS (ARCHIVE)
-- ============================================
\echo ''
\echo '=== Archiving old sent emails ==='

-- Count emails older than 1 year
SELECT COUNT(*) AS old_emails_count
FROM verzonden_emails
WHERE verzonden_op < NOW() - INTERVAL '1 year';

-- Option 1: Create archive table (recommended)
CREATE TABLE IF NOT EXISTS verzonden_emails_archive (
    LIKE verzonden_emails INCLUDING ALL
);

-- Move to archive instead of deleting
WITH moved AS (
    DELETE FROM verzonden_emails
    WHERE verzonden_op < NOW() - INTERVAL '1 year'
    RETURNING *
)
INSERT INTO verzonden_emails_archive
SELECT * FROM moved;

\echo 'Old emails moved to archive table.'

-- ============================================
-- SECTION 3: PROCESSED INCOMING EMAILS
-- ============================================
\echo ''
\echo '=== Cleaning up old processed incoming emails ==='

-- Count processed emails older than 6 months
SELECT COUNT(*) AS processed_old_count
FROM incoming_emails
WHERE is_processed = TRUE
  AND processed_at < NOW() - INTERVAL '6 months';

-- Archive or delete old processed emails
CREATE TABLE IF NOT EXISTS incoming_emails_archive (
    LIKE incoming_emails INCLUDING ALL
);

WITH moved AS (
    DELETE FROM incoming_emails
    WHERE is_processed = TRUE
      AND processed_at < NOW() - INTERVAL '6 months'
    RETURNING *
)
INSERT INTO incoming_emails_archive
SELECT * FROM moved;

\echo 'Old processed emails moved to archive.'

-- ============================================
-- SECTION 4: SOFT DELETE CLEANUP
-- ============================================
\echo ''
\echo '=== Cleaning up soft-deleted uploaded images ==='

-- Count soft-deleted images older than 3 months
SELECT COUNT(*) AS deleted_images_count
FROM uploaded_images
WHERE deleted_at IS NOT NULL
  AND deleted_at < NOW() - INTERVAL '3 months';

-- Permanently delete soft-deleted images
DELETE FROM uploaded_images
WHERE deleted_at IS NOT NULL
  AND deleted_at < NOW() - INTERVAL '3 months';

\echo 'Old soft-deleted images permanently removed.'

-- ============================================
-- SECTION 5: OLD CHAT MESSAGES (OPTIONAL)
-- ============================================
\echo ''
\echo '=== Archiving old chat messages (optional) ==='

-- Count messages older than 2 years
SELECT COUNT(*) AS old_messages_count
FROM chat_messages
WHERE created_at < NOW() - INTERVAL '2 years';

-- Uncomment to archive old messages
-- CREATE TABLE IF NOT EXISTS chat_messages_archive (
--     LIKE chat_messages INCLUDING ALL
-- );

-- WITH moved AS (
--     DELETE FROM chat_messages
--     WHERE created_at < NOW() - INTERVAL '2 years'
--     RETURNING *
-- )
-- INSERT INTO chat_messages_archive
-- SELECT * FROM moved;

\echo 'Chat message archiving skipped (uncomment to enable).'

-- ============================================
-- SECTION 6: CLEANUP STATISTICS
-- ============================================
\echo ''
\echo '=== CLEANUP STATISTICS ==='

-- Updated table sizes
SELECT
    'verzonden_emails' AS table_name,
    pg_size_pretty(pg_total_relation_size('verzonden_emails')) AS size,
    (SELECT COUNT(*) FROM verzonden_emails) AS row_count
UNION ALL
SELECT
    'verzonden_emails_archive',
    pg_size_pretty(pg_total_relation_size('verzonden_emails_archive')),
    (SELECT COUNT(*) FROM verzonden_emails_archive)
UNION ALL
SELECT
    'incoming_emails',
    pg_size_pretty(pg_total_relation_size('incoming_emails')),
    (SELECT COUNT(*) FROM incoming_emails)
UNION ALL
SELECT
    'incoming_emails_archive',
    pg_size_pretty(pg_total_relation_size('incoming_emails_archive')),
    (SELECT COUNT(*) FROM incoming_emails_archive)
UNION ALL
SELECT
    'refresh_tokens',
    pg_size_pretty(pg_total_relation_size('refresh_tokens')),
    (SELECT COUNT(*) FROM refresh_tokens)
UNION ALL
SELECT
    'uploaded_images',
    pg_size_pretty(pg_total_relation_size('uploaded_images')),
    (SELECT COUNT(*) FROM uploaded_images);

-- ============================================
-- SECTION 7: VACUUM AFTER CLEANUP
-- ============================================
\echo ''
\echo '=== Running VACUUM to reclaim space ==='

VACUUM ANALYZE verzonden_emails;
VACUUM ANALYZE incoming_emails;
VACUUM ANALYZE refresh_tokens;
VACUUM ANALYZE uploaded_images;

\echo ''
\echo 'Data cleanup completed successfully!'
\echo 'Remember to run a full VACUUM ANALYZE for optimal performance.'