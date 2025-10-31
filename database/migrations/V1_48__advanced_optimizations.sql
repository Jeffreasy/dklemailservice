-- Migratie: V1_48__advanced_optimizations.sql
-- Beschrijving: Advanced optimizations - triggers, constraints, cleanup duplicates
-- Versie: 1.48.0
-- Datum: 2025-10-31
-- Impact: MEDIUM - Schema cleanup and automated timestamp management

-- ============================================
-- SECTION 1: CLEANUP DUPLICATE CONSTRAINTS
-- ============================================
-- Remove duplicate foreign key constraints (GORM creates both)

-- contact_antwoorden: Remove duplicate FK (keep the newer one)
ALTER TABLE contact_antwoorden DROP CONSTRAINT IF EXISTS fk_contact_antwoorden_contact_id CASCADE;
COMMENT ON CONSTRAINT contact_antwoorden_contact_id_fkey ON contact_antwoorden IS 'FK to contact_formulieren';

-- aanmelding_antwoorden: Remove duplicate FK
ALTER TABLE aanmelding_antwoorden DROP CONSTRAINT IF EXISTS fk_aanmelding_antwoorden_aanmelding_id CASCADE;
COMMENT ON CONSTRAINT aanmelding_antwoorden_aanmelding_id_fkey ON aanmelding_antwoorden IS 'FK to aanmeldingen';

-- ============================================
-- SECTION 2: AUTO UPDATE_AT TRIGGERS
-- ============================================
-- Automatically update updated_at timestamp on row changes

-- Create generic trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION update_updated_at_column() IS 'Generic trigger function to update updated_at timestamp';

-- Apply to all tables with updated_at column
-- gebruikers
DROP TRIGGER IF EXISTS trigger_gebruikers_updated_at ON gebruikers;
CREATE TRIGGER trigger_gebruikers_updated_at
    BEFORE UPDATE ON gebruikers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- contact_formulieren
DROP TRIGGER IF EXISTS trigger_contact_formulieren_updated_at ON contact_formulieren;
CREATE TRIGGER trigger_contact_formulieren_updated_at
    BEFORE UPDATE ON contact_formulieren
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- contact_antwoorden
DROP TRIGGER IF EXISTS trigger_contact_antwoorden_updated_at ON contact_antwoorden;
CREATE TRIGGER trigger_contact_antwoorden_updated_at
    BEFORE UPDATE ON contact_antwoorden
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- aanmeldingen
DROP TRIGGER IF EXISTS trigger_aanmeldingen_updated_at ON aanmeldingen;
CREATE TRIGGER trigger_aanmeldingen_updated_at
    BEFORE UPDATE ON aanmeldingen
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- aanmelding_antwoorden
DROP TRIGGER IF EXISTS trigger_aanmelding_antwoorden_updated_at ON aanmelding_antwoorden;
CREATE TRIGGER trigger_aanmelding_antwoorden_updated_at
    BEFORE UPDATE ON aanmelding_antwoorden
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- email_templates
DROP TRIGGER IF EXISTS trigger_email_templates_updated_at ON email_templates;
CREATE TRIGGER trigger_email_templates_updated_at
    BEFORE UPDATE ON email_templates
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- verzonden_emails
DROP TRIGGER IF EXISTS trigger_verzonden_emails_updated_at ON verzonden_emails;
CREATE TRIGGER trigger_verzonden_emails_updated_at
    BEFORE UPDATE ON verzonden_emails
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- incoming_emails
DROP TRIGGER IF EXISTS trigger_incoming_emails_updated_at ON incoming_emails;
CREATE TRIGGER trigger_incoming_emails_updated_at
    BEFORE UPDATE ON incoming_emails
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Chat tables
DROP TRIGGER IF EXISTS trigger_chat_channels_updated_at ON chat_channels;
CREATE TRIGGER trigger_chat_channels_updated_at
    BEFORE UPDATE ON chat_channels
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS trigger_chat_messages_updated_at ON chat_messages;
CREATE TRIGGER trigger_chat_messages_updated_at
    BEFORE UPDATE ON chat_messages
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS trigger_chat_user_presence_updated_at ON chat_user_presence;
CREATE TRIGGER trigger_chat_user_presence_updated_at
    BEFORE UPDATE ON chat_user_presence
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Content tables
DROP TRIGGER IF EXISTS trigger_newsletters_updated_at ON newsletters;
CREATE TRIGGER trigger_newsletters_updated_at
    BEFORE UPDATE ON newsletters
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS trigger_uploaded_images_updated_at ON uploaded_images;
CREATE TRIGGER trigger_uploaded_images_updated_at
    BEFORE UPDATE ON uploaded_images
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS trigger_photos_updated_at ON photos;
CREATE TRIGGER trigger_photos_updated_at
    BEFORE UPDATE ON photos
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS trigger_albums_updated_at ON albums;
CREATE TRIGGER trigger_albums_updated_at
    BEFORE UPDATE ON albums
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS trigger_videos_updated_at ON videos;
CREATE TRIGGER trigger_videos_updated_at
    BEFORE UPDATE ON videos
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS trigger_sponsors_updated_at ON sponsors;
CREATE TRIGGER trigger_sponsors_updated_at
    BEFORE UPDATE ON sponsors
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- SECTION 3: DATA VALIDATION CONSTRAINTS
-- ============================================
-- Add constraints for data integrity

-- Email validation (basic regex check)
ALTER TABLE gebruikers DROP CONSTRAINT IF EXISTS gebruikers_email_check;
ALTER TABLE gebruikers ADD CONSTRAINT gebruikers_email_check 
    CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');

ALTER TABLE contact_formulieren DROP CONSTRAINT IF EXISTS contact_formulieren_email_check;
ALTER TABLE contact_formulieren ADD CONSTRAINT contact_formulieren_email_check 
    CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');

ALTER TABLE aanmeldingen DROP CONSTRAINT IF EXISTS aanmeldingen_email_check;
ALTER TABLE aanmeldingen ADD CONSTRAINT aanmeldingen_email_check 
    CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');

-- Steps must be non-negative
ALTER TABLE aanmeldingen DROP CONSTRAINT IF EXISTS aanmeldingen_steps_check;
ALTER TABLE aanmeldingen ADD CONSTRAINT aanmeldingen_steps_check 
    CHECK (steps >= 0);

-- Status validation for contact_formulieren
ALTER TABLE contact_formulieren DROP CONSTRAINT IF EXISTS contact_formulieren_status_check;
ALTER TABLE contact_formulieren ADD CONSTRAINT contact_formulieren_status_check 
    CHECK (status IN ('nieuw', 'in_behandeling', 'beantwoord', 'gesloten'));

-- Status validation for aanmeldingen
ALTER TABLE aanmeldingen DROP CONSTRAINT IF EXISTS aanmeldingen_status_check;
ALTER TABLE aanmeldingen ADD CONSTRAINT aanmeldingen_status_check 
    CHECK (status IN ('nieuw', 'bevestigd', 'geannuleerd', 'voltooid'));

-- Email status consistency
ALTER TABLE contact_formulieren DROP CONSTRAINT IF EXISTS contact_formulieren_email_consistency;
ALTER TABLE contact_formulieren ADD CONSTRAINT contact_formulieren_email_consistency 
    CHECK (
        (email_verzonden = FALSE AND email_verzonden_op IS NULL) OR
        (email_verzonden = TRUE AND email_verzonden_op IS NOT NULL)
    );

ALTER TABLE aanmeldingen DROP CONSTRAINT IF EXISTS aanmeldingen_email_consistency;
ALTER TABLE aanmeldingen ADD CONSTRAINT aanmeldingen_email_consistency 
    CHECK (
        (email_verzonden = FALSE AND email_verzonden_op IS NULL) OR
        (email_verzonden = TRUE AND email_verzonden_op IS NOT NULL)
    );

-- ============================================
-- SECTION 4: MISSING UPDATED_AT DEFAULT
-- ============================================
-- Ensure all updated_at columns have proper defaults

-- This was already done in schema, but ensuring consistency
ALTER TABLE gebruikers ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE contact_formulieren ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE contact_antwoorden ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE aanmeldingen ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE aanmelding_antwoorden ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE email_templates ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE verzonden_emails ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE incoming_emails ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP;

-- ============================================
-- SECTION 5: PERFORMANCE HINTS FOR QUERY PLANNER
-- ============================================
-- Update statistics targets for frequently queried columns

-- Increase statistics for email columns (used in lookups)
ALTER TABLE gebruikers ALTER COLUMN email SET STATISTICS 1000;
ALTER TABLE contact_formulieren ALTER COLUMN email SET STATISTICS 1000;
ALTER TABLE aanmeldingen ALTER COLUMN email SET STATISTICS 1000;

-- Increase statistics for status columns (used in filtering)
ALTER TABLE contact_formulieren ALTER COLUMN status SET STATISTICS 500;
ALTER TABLE aanmeldingen ALTER COLUMN status SET STATISTICS 500;
ALTER TABLE verzonden_emails ALTER COLUMN status SET STATISTICS 500;

-- Increase statistics for foreign keys
ALTER TABLE verzonden_emails ALTER COLUMN contact_id SET STATISTICS 500;
ALTER TABLE verzonden_emails ALTER COLUMN aanmelding_id SET STATISTICS 500;
ALTER TABLE contact_antwoorden ALTER COLUMN contact_id SET STATISTICS 500;
ALTER TABLE aanmelding_antwoorden ALTER COLUMN aanmelding_id SET STATISTICS 500;

-- ============================================
-- SECTION 6: DENORMALIZATION OPPORTUNITIES
-- ============================================
-- Add computed columns for frequently accessed aggregates

-- Add cached count of antwoorden to contact_formulieren
ALTER TABLE contact_formulieren ADD COLUMN IF NOT EXISTS antwoorden_count INTEGER DEFAULT 0;
COMMENT ON COLUMN contact_formulieren.antwoorden_count IS 'Cached count of responses - updated via trigger';

-- Create trigger to maintain count
CREATE OR REPLACE FUNCTION update_contact_antwoorden_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE contact_formulieren 
        SET antwoorden_count = antwoorden_count + 1 
        WHERE id = NEW.contact_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE contact_formulieren 
        SET antwoorden_count = GREATEST(0, antwoorden_count - 1) 
        WHERE id = OLD.contact_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_contact_antwoorden_count ON contact_antwoorden;
CREATE TRIGGER trigger_contact_antwoorden_count
    AFTER INSERT OR DELETE ON contact_antwoorden
    FOR EACH ROW
    EXECUTE FUNCTION update_contact_antwoorden_count();

-- Initialize counts
UPDATE contact_formulieren cf
SET antwoorden_count = (
    SELECT COUNT(*) 
    FROM contact_antwoorden ca 
    WHERE ca.contact_id = cf.id
);

-- Same for aanmeldingen
ALTER TABLE aanmeldingen ADD COLUMN IF NOT EXISTS antwoorden_count INTEGER DEFAULT 0;
COMMENT ON COLUMN aanmeldingen.antwoorden_count IS 'Cached count of responses - updated via trigger';

CREATE OR REPLACE FUNCTION update_aanmelding_antwoorden_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE aanmeldingen 
        SET antwoorden_count = antwoorden_count + 1 
        WHERE id = NEW.aanmelding_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE aanmeldingen 
        SET antwoorden_count = GREATEST(0, antwoorden_count - 1) 
        WHERE id = OLD.aanmelding_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_aanmelding_antwoorden_count ON aanmelding_antwoorden;
CREATE TRIGGER trigger_aanmelding_antwoorden_count
    AFTER INSERT OR DELETE ON aanmelding_antwoorden
    FOR EACH ROW
    EXECUTE FUNCTION update_aanmelding_antwoorden_count();

-- Initialize counts
UPDATE aanmeldingen a
SET antwoorden_count = (
    SELECT COUNT(*) 
    FROM aanmelding_antwoorden aa 
    WHERE aa.aanmelding_id = a.id
);

-- ============================================
-- SECTION 7: MATERIALIZED VIEWS
-- ============================================
-- Create materialized views for frequently accessed aggregates

-- Dashboard statistics view
CREATE MATERIALIZED VIEW IF NOT EXISTS dashboard_stats AS
SELECT
    'contact_formulieren' as entity,
    status,
    beantwoord,
    COUNT(*) as count,
    MAX(created_at) as last_created
FROM contact_formulieren
GROUP BY status, beantwoord
UNION ALL
SELECT
    'aanmeldingen' as entity,
    status,
    NULL as beantwoord,
    COUNT(*) as count,
    MAX(created_at) as last_created
FROM aanmeldingen
GROUP BY status
UNION ALL
SELECT
    'verzonden_emails' as entity,
    status,
    NULL as beantwoord,
    COUNT(*) as count,
    MAX(verzonden_op) as last_created
FROM verzonden_emails
WHERE verzonden_op > NOW() - INTERVAL '30 days'
GROUP BY status;

-- Create index on materialized view
CREATE UNIQUE INDEX IF NOT EXISTS idx_dashboard_stats_entity_status 
ON dashboard_stats(entity, status);

COMMENT ON MATERIALIZED VIEW dashboard_stats IS 'Cached dashboard statistics - refresh hourly or on demand';

-- Refresh function
CREATE OR REPLACE FUNCTION refresh_dashboard_stats()
RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY dashboard_stats;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION refresh_dashboard_stats() IS 'Refresh dashboard statistics view (concurrent safe)';

-- ============================================
-- SECTION 8: DATA QUALITY CONSTRAINTS
-- ============================================

-- Naam cannot be empty string
ALTER TABLE gebruikers DROP CONSTRAINT IF EXISTS gebruikers_naam_not_empty;
ALTER TABLE gebruikers ADD CONSTRAINT gebruikers_naam_not_empty 
    CHECK (LENGTH(TRIM(naam)) > 0);

ALTER TABLE contact_formulieren DROP CONSTRAINT IF EXISTS contact_formulieren_naam_not_empty;
ALTER TABLE contact_formulieren ADD CONSTRAINT contact_formulieren_naam_not_empty 
    CHECK (LENGTH(TRIM(naam)) > 0);

ALTER TABLE aanmeldingen DROP CONSTRAINT IF EXISTS aanmeldingen_naam_not_empty;
ALTER TABLE aanmeldingen ADD CONSTRAINT aanmeldingen_naam_not_empty 
    CHECK (LENGTH(TRIM(naam)) > 0);

-- Bericht cannot be empty
ALTER TABLE contact_formulieren DROP CONSTRAINT IF EXISTS contact_formulieren_bericht_not_empty;
ALTER TABLE contact_formulieren ADD CONSTRAINT contact_formulieren_bericht_not_empty 
    CHECK (LENGTH(TRIM(bericht)) > 0);

-- Telefoon format (optional - basic check)
ALTER TABLE aanmeldingen DROP CONSTRAINT IF EXISTS aanmeldingen_telefoon_format;
ALTER TABLE aanmeldingen ADD CONSTRAINT aanmeldingen_telefoon_format 
    CHECK (telefoon IS NULL OR LENGTH(TRIM(telefoon)) >= 10);

-- ============================================
-- SECTION 9: PERFORMANCE INDEXES ON COMPUTED COLUMNS
-- ============================================

-- Index on antwoorden_count for filtering
CREATE INDEX IF NOT EXISTS idx_contact_formulieren_antwoorden_count 
ON contact_formulieren(antwoorden_count) 
WHERE antwoorden_count > 0;

CREATE INDEX IF NOT EXISTS idx_aanmeldingen_antwoorden_count 
ON aanmeldingen(antwoorden_count) 
WHERE antwoorden_count > 0;

-- ============================================
-- SECTION 10: MISSING INDEXES ON STATUS TRANSITIONS
-- ============================================

-- Track status changes over time
CREATE INDEX IF NOT EXISTS idx_contact_formulieren_behandeld 
ON contact_formulieren(behandeld_op DESC NULLS LAST, status);

CREATE INDEX IF NOT EXISTS idx_aanmeldingen_behandeld 
ON aanmeldingen(behandeld_op DESC NULLS LAST, status);

-- ============================================
-- SECTION 11: COMPOSITE UNIQUE CONSTRAINTS
-- ============================================
-- Prevent duplicate submissions (same email within timeframe)

-- Note: This is intentionally commented out as it might be too restrictive
-- Uncomment if you want to prevent rapid duplicate submissions

-- CREATE UNIQUE INDEX IF NOT EXISTS idx_contact_formulieren_email_recent
-- ON contact_formulieren(email, DATE(created_at))
-- WHERE created_at > NOW() - INTERVAL '1 day';

-- CREATE UNIQUE INDEX IF NOT EXISTS idx_aanmeldingen_email_recent
-- ON aanmeldingen(email, DATE(created_at))
-- WHERE created_at > NOW() - INTERVAL '1 day';

-- ============================================
-- REGISTER MIGRATION
-- ============================================

INSERT INTO migraties (versie, naam, toegepast) 
VALUES ('1.48.0', 'Advanced optimizations: triggers, constraints, materialized views, data quality', CURRENT_TIMESTAMP)
ON CONFLICT (versie) DO NOTHING;

-- ============================================
-- POST-MIGRATION INSTRUCTIONS
-- ============================================
-- After running this migration:
-- 1. ANALYZE; -- Update query planner statistics
-- 2. SELECT refresh_dashboard_stats(); -- Initialize materialized view
-- 3. Set up hourly refresh job (if using cron):
--    0 * * * * psql $DATABASE_URL -c "SELECT refresh_dashboard_stats();"
-- 4. Monitor trigger performance on high-traffic tables
-- 5. Review constraint violations in application logs

-- Expected improvements:
-- - Automatic updated_at management (no application code changes needed)
-- - Better data quality (email validation, non-empty names)
-- - Fast dashboard stats via materialized view (100x faster for aggregates)
-- - Denormalized counts avoid expensive COUNT queries
-- - Cleanup of duplicate constraints reduces overhead