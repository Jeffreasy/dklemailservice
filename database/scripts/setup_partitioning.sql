-- Partitioning Setup Script
-- This script sets up time-based partitioning for large tables
-- IMPORTANT: This requires taking tables offline during migration!
-- Usage: docker exec dkl-postgres psql -U postgres -d dklemailservice -f /path/to/setup_partitioning.sql

-- ============================================
-- PREREQUISITES
-- ============================================
-- 1. Full database backup created
-- 2. Maintenance window scheduled
-- 3. Application stopped (no active connections)

\echo '============================================'
\echo 'TABLE PARTITIONING SETUP'
\echo '============================================'
\echo ''
\echo 'This script will partition large tables for better performance.'
\echo 'IMPORTANT: This operation requires downtime!'
\echo ''
\prompt 'Continue? (yes/no): ' confirmation

-- ============================================
-- SECTION 1: VERZONDEN_EMAILS PARTITIONING
-- ============================================
\echo ''
\echo '=== Setting up partitioning for verzonden_emails ==='

-- Step 1: Rename existing table
ALTER TABLE verzonden_emails RENAME TO verzonden_emails_old;

-- Step 2: Create new partitioned table
CREATE TABLE verzonden_emails (
    id UUID DEFAULT gen_random_uuid(),
    ontvanger VARCHAR(255) NOT NULL,
    onderwerp VARCHAR(255) NOT NULL,
    inhoud TEXT NOT NULL,
    verzonden_op TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) NOT NULL DEFAULT 'verzonden',
    fout_bericht TEXT,
    contact_id UUID REFERENCES contact_formulieren(id),
    aanmelding_id UUID REFERENCES aanmeldingen(id),
    template_id UUID REFERENCES email_templates(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id, verzonden_op)
) PARTITION BY RANGE (verzonden_op);

-- Step 3: Create partitions for current and future months
-- Past partitions (adjust dates as needed)
CREATE TABLE verzonden_emails_2024_01 PARTITION OF verzonden_emails
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');

CREATE TABLE verzonden_emails_2024_02 PARTITION OF verzonden_emails
    FOR VALUES FROM ('2024-02-01') TO ('2024-03-01');

CREATE TABLE verzonden_emails_2024_03 PARTITION OF verzonden_emails
    FOR VALUES FROM ('2024-03-01') TO ('2024-04-01');

CREATE TABLE verzonden_emails_2024_04 PARTITION OF verzonden_emails
    FOR VALUES FROM ('2024-04-01') TO ('2024-05-01');

CREATE TABLE verzonden_emails_2024_05 PARTITION OF verzonden_emails
    FOR VALUES FROM ('2024-05-01') TO ('2024-06-01');

CREATE TABLE verzonden_emails_2024_06 PARTITION OF verzonden_emails
    FOR VALUES FROM ('2024-06-01') TO ('2024-07-01');

CREATE TABLE verzonden_emails_2024_07 PARTITION OF verzonden_emails
    FOR VALUES FROM ('2024-07-01') TO ('2024-08-01');

CREATE TABLE verzonden_emails_2024_08 PARTITION OF verzonden_emails
    FOR VALUES FROM ('2024-08-01') TO ('2024-09-01');

CREATE TABLE verzonden_emails_2024_09 PARTITION OF verzonden_emails
    FOR VALUES FROM ('2024-09-01') TO ('2024-10-01');

CREATE TABLE verzonden_emails_2024_10 PARTITION OF verzonden_emails
    FOR VALUES FROM ('2024-10-01') TO ('2024-11-01');

CREATE TABLE verzonden_emails_2024_11 PARTITION OF verzonden_emails
    FOR VALUES FROM ('2024-11-01') TO ('2024-12-01');

CREATE TABLE verzonden_emails_2024_12 PARTITION OF verzonden_emails
    FOR VALUES FROM ('2024-12-01') TO ('2025-01-01');

-- Current year 2025
CREATE TABLE verzonden_emails_2025_01 PARTITION OF verzonden_emails
    FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');

CREATE TABLE verzonden_emails_2025_02 PARTITION OF verzonden_emails
    FOR VALUES FROM ('2025-02-01') TO ('2025-03-01');

CREATE TABLE verzonden_emails_2025_03 PARTITION OF verzonden_emails
    FOR VALUES FROM ('2025-03-01') TO ('2025-04-01');

CREATE TABLE verzonden_emails_2025_04 PARTITION OF verzonden_emails
    FOR VALUES FROM ('2025-04-01') TO ('2025-05-01');

CREATE TABLE verzonden_emails_2025_05 PARTITION OF verzonden_emails
    FOR VALUES FROM ('2025-05-01') TO ('2025-06-01');

CREATE TABLE verzonden_emails_2025_06 PARTITION OF verzonden_emails
    FOR VALUES FROM ('2025-06-01') TO ('2025-07-01');

CREATE TABLE verzonden_emails_2025_07 PARTITION OF verzonden_emails
    FOR VALUES FROM ('2025-07-01') TO ('2025-08-01');

CREATE TABLE verzonden_emails_2025_08 PARTITION OF verzonden_emails
    FOR VALUES FROM ('2025-08-01') TO ('2025-09-01');

CREATE TABLE verzonden_emails_2025_09 PARTITION OF verzonden_emails
    FOR VALUES FROM ('2025-09-01') TO ('2025-10-01');

CREATE TABLE verzonden_emails_2025_10 PARTITION OF verzonden_emails
    FOR VALUES FROM ('2025-10-01') TO ('2025-11-01');

CREATE TABLE verzonden_emails_2025_11 PARTITION OF verzonden_emails
    FOR VALUES FROM ('2025-11-01') TO ('2025-12-01');

CREATE TABLE verzonden_emails_2025_12 PARTITION OF verzonden_emails
    FOR VALUES FROM ('2025-12-01') TO ('2026-01-01');

-- Future months (2026)
CREATE TABLE verzonden_emails_2026_01 PARTITION OF verzonden_emails
    FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');

CREATE TABLE verzonden_emails_2026_02 PARTITION OF verzonden_emails
    FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');

CREATE TABLE verzonden_emails_2026_03 PARTITION OF verzonden_emails
    FOR VALUES FROM ('2026-03-01') TO ('2026-04-01');

-- Step 4: Recreate indexes on partitioned table
CREATE INDEX idx_verzonden_emails_contact_id ON verzonden_emails(contact_id);
CREATE INDEX idx_verzonden_emails_aanmelding_id ON verzonden_emails(aanmelding_id);
CREATE INDEX idx_verzonden_emails_template_id ON verzonden_emails(template_id);
CREATE INDEX idx_verzonden_emails_status ON verzonden_emails(status);
CREATE INDEX idx_verzonden_emails_ontvanger ON verzonden_emails(ontvanger);
CREATE INDEX idx_verzonden_emails_status_tijd ON verzonden_emails(status, verzonden_op DESC);

-- Step 5: Migrate data from old table
INSERT INTO verzonden_emails 
SELECT * FROM verzonden_emails_old;

-- Step 6: Verify data migration
SELECT 
    'Old table' AS source,
    COUNT(*) AS row_count
FROM verzonden_emails_old
UNION ALL
SELECT 
    'New partitioned table',
    COUNT(*)
FROM verzonden_emails;

\echo 'verzonden_emails partitioning completed.'

-- ============================================
-- SECTION 2: CHAT_MESSAGES PARTITIONING
-- ============================================
\echo ''
\echo '=== Setting up partitioning for chat_messages ==='

-- Step 1: Rename existing table
ALTER TABLE chat_messages RENAME TO chat_messages_old;

-- Step 2: Create new partitioned table
CREATE TABLE chat_messages (
    id UUID DEFAULT gen_random_uuid(),
    channel_id UUID REFERENCES chat_channels(id) ON DELETE CASCADE,
    user_id UUID,
    content TEXT,
    message_type TEXT DEFAULT 'text' CHECK (message_type IN ('text', 'image', 'file', 'system')),
    file_url TEXT,
    file_name TEXT,
    file_size INTEGER,
    thumbnail_url TEXT,
    reply_to_id UUID,
    edited_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

-- Create partitions for 2024-2026
CREATE TABLE chat_messages_2024_q1 PARTITION OF chat_messages
    FOR VALUES FROM ('2024-01-01') TO ('2024-04-01');

CREATE TABLE chat_messages_2024_q2 PARTITION OF chat_messages
    FOR VALUES FROM ('2024-04-01') TO ('2024-07-01');

CREATE TABLE chat_messages_2024_q3 PARTITION OF chat_messages
    FOR VALUES FROM ('2024-07-01') TO ('2024-10-01');

CREATE TABLE chat_messages_2024_q4 PARTITION OF chat_messages
    FOR VALUES FROM ('2024-10-01') TO ('2025-01-01');

CREATE TABLE chat_messages_2025_q1 PARTITION OF chat_messages
    FOR VALUES FROM ('2025-01-01') TO ('2025-04-01');

CREATE TABLE chat_messages_2025_q2 PARTITION OF chat_messages
    FOR VALUES FROM ('2025-04-01') TO ('2025-07-01');

CREATE TABLE chat_messages_2025_q3 PARTITION OF chat_messages
    FOR VALUES FROM ('2025-07-01') TO ('2025-10-01');

CREATE TABLE chat_messages_2025_q4 PARTITION OF chat_messages
    FOR VALUES FROM ('2025-10-01') TO ('2026-01-01');

CREATE TABLE chat_messages_2026_q1 PARTITION OF chat_messages
    FOR VALUES FROM ('2026-01-01') TO ('2026-04-01');

CREATE TABLE chat_messages_2026_q2 PARTITION OF chat_messages
    FOR VALUES FROM ('2026-04-01') TO ('2026-07-01');

-- Step 3: Recreate indexes
CREATE INDEX idx_chat_messages_channel_id_created_at ON chat_messages(channel_id, created_at DESC);
CREATE INDEX idx_chat_messages_user_id ON chat_messages(user_id);
CREATE INDEX idx_chat_messages_reply_to ON chat_messages(reply_to_id, created_at DESC) 
WHERE reply_to_id IS NOT NULL;
CREATE INDEX idx_chat_messages_files ON chat_messages(channel_id, created_at DESC) 
WHERE message_type IN ('image', 'file');
CREATE INDEX idx_chat_messages_fts ON chat_messages 
USING gin(to_tsvector('dutch', COALESCE(content, '')));

-- Step 4: Migrate data
INSERT INTO chat_messages 
SELECT * FROM chat_messages_old;

-- Step 5: Verify migration
SELECT 
    'Old table' AS source,
    COUNT(*) AS row_count
FROM chat_messages_old
UNION ALL
SELECT 
    'New partitioned table',
    COUNT(*)
FROM chat_messages;

\echo 'chat_messages partitioning completed.'

-- ============================================
-- SECTION 3: CREATE PARTITION MAINTENANCE FUNCTION
-- ============================================
\echo ''
\echo '=== Creating automatic partition creation function ==='

-- Function to create next month's partition
CREATE OR REPLACE FUNCTION create_next_month_partition()
RETURNS void AS $$
DECLARE
    next_month DATE;
    partition_name TEXT;
    start_date DATE;
    end_date DATE;
BEGIN
    -- Calculate next month
    next_month := date_trunc('month', CURRENT_DATE + INTERVAL '1 month');
    start_date := next_month;
    end_date := next_month + INTERVAL '1 month';
    
    -- Create partition for verzonden_emails
    partition_name := 'verzonden_emails_' || to_char(next_month, 'YYYY_MM');
    EXECUTE format(
        'CREATE TABLE IF NOT EXISTS %I PARTITION OF verzonden_emails 
         FOR VALUES FROM (%L) TO (%L)',
        partition_name, start_date, end_date
    );
    RAISE NOTICE 'Created partition: %', partition_name;
    
    -- Create partition for chat_messages (quarterly)
    IF EXTRACT(MONTH FROM next_month) IN (1, 4, 7, 10) THEN
        partition_name := 'chat_messages_' || to_char(next_month, 'YYYY') || '_q' || 
                         EXTRACT(QUARTER FROM next_month)::TEXT;
        EXECUTE format(
            'CREATE TABLE IF NOT EXISTS %I PARTITION OF chat_messages 
             FOR VALUES FROM (%L) TO (%L)',
            partition_name, start_date, start_date + INTERVAL '3 months'
        );
        RAISE NOTICE 'Created partition: %', partition_name;
    END IF;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION create_next_month_partition() IS 
'Creates partitions for the next month. Run this monthly via cron job.';

\echo 'Partition maintenance function created.'

-- ============================================
-- SECTION 4: CLEANUP
-- ============================================
\echo ''
\echo '=== Final steps ==='

-- Drop old tables after verification
-- UNCOMMENT AFTER VERIFYING DATA MIGRATION:
-- DROP TABLE verzonden_emails_old;
-- DROP TABLE chat_messages_old;

\echo ''
\echo '============================================'
\echo 'Partitioning setup completed!'
\echo '============================================'
\echo ''
\echo 'IMPORTANT NEXT STEPS:'
\echo '1. Verify data integrity in partitioned tables'
\echo '2. Run ANALYZE on all partitions'
\echo '3. Test application functionality'
\echo '4. Drop old tables after verification (see cleanup section)'
\echo '5. Set up monthly cron job to run create_next_month_partition()'
\echo ''
\echo 'Example cron job:'
\echo '0 0 1 * * docker exec dkl-postgres psql -U postgres -d dklemailservice -c "SELECT create_next_month_partition();"'