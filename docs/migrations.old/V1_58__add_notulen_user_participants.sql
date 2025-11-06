-- Migration: V1_58__add_notulen_user_participants.sql
-- Description: Add UUID participant columns to notulen tables for proper user linking
-- Date: 2026-01-04
--
-- IDEMPOTENT: Deze migratie kan veilig meerdere keren worden uitgevoerd

-- =====================================================
-- 1. ADD UUID PARTICIPANT COLUMNS TO NOTULEN TABLE
-- =====================================================

-- Add columns for registered user participants (UUID arrays)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'notulen' AND column_name = 'aanwezigen_gebruikers') THEN
        ALTER TABLE notulen ADD COLUMN aanwezigen_gebruikers UUID[];
        COMMENT ON COLUMN notulen.aanwezigen_gebruikers IS 'UUIDs van geregistreerde gebruikers die aanwezig waren';
    END IF;

    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'notulen' AND column_name = 'afwezigen_gebruikers') THEN
        ALTER TABLE notulen ADD COLUMN afwezigen_gebruikers UUID[];
        COMMENT ON COLUMN notulen.afwezigen_gebruikers IS 'UUIDs van geregistreerde gebruikers die afwezig waren';
    END IF;
END $$;

-- Add columns for guest participants (text arrays) - keep existing text arrays for backwards compatibility
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'notulen' AND column_name = 'aanwezigen_gasten') THEN
        ALTER TABLE notulen ADD COLUMN aanwezigen_gasten TEXT[];
        COMMENT ON COLUMN notulen.aanwezigen_gasten IS 'Namen van niet-geregistreerde gasten die aanwezig waren';
    END IF;

    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'notulen' AND column_name = 'afwezigen_gasten') THEN
        ALTER TABLE notulen ADD COLUMN afwezigen_gasten TEXT[];
        COMMENT ON COLUMN notulen.afwezigen_gasten IS 'Namen van niet-geregistreerde gasten die afwezig waren';
    END IF;
END $$;

-- =====================================================
-- 2. ADD UUID PARTICIPANT COLUMNS TO NOTULEN_VERSIES TABLE
-- =====================================================

-- Add same columns to versions table
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'notulen_versies' AND column_name = 'aanwezigen_gebruikers') THEN
        ALTER TABLE notulen_versies ADD COLUMN aanwezigen_gebruikers UUID[];
        COMMENT ON COLUMN notulen_versies.aanwezigen_gebruikers IS 'UUIDs van geregistreerde gebruikers die aanwezig waren (versie snapshot)';
    END IF;

    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'notulen_versies' AND column_name = 'afwezigen_gebruikers') THEN
        ALTER TABLE notulen_versies ADD COLUMN afwezigen_gebruikers UUID[];
        COMMENT ON COLUMN notulen_versies.afwezigen_gebruikers IS 'UUIDs van geregistreerde gebruikers die afwezig waren (versie snapshot)';
    END IF;

    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'notulen_versies' AND column_name = 'aanwezigen_gasten') THEN
        ALTER TABLE notulen_versies ADD COLUMN aanwezigen_gasten TEXT[];
        COMMENT ON COLUMN notulen_versies.aanwezigen_gasten IS 'Namen van niet-geregistreerde gasten die aanwezig waren (versie snapshot)';
    END IF;

    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'notulen_versies' AND column_name = 'afwezigen_gasten') THEN
        ALTER TABLE notulen_versies ADD COLUMN afwezigen_gasten TEXT[];
        COMMENT ON COLUMN notulen_versies.afwezigen_gasten IS 'Namen van niet-geregistreerde gasten die afwezig waren (versie snapshot)';
    END IF;
END $$;

-- =====================================================
-- 3. ADD INDEXES FOR PERFORMANCE
-- =====================================================

-- Indexes for efficient user participant lookups
CREATE INDEX IF NOT EXISTS idx_notulen_aanwezigen_gebruikers ON notulen USING GIN(aanwezigen_gebruikers);
CREATE INDEX IF NOT EXISTS idx_notulen_afwezigen_gebruikers ON notulen USING GIN(afwezigen_gebruikers);
CREATE INDEX IF NOT EXISTS idx_notulen_versies_aanwezigen_gebruikers ON notulen_versies USING GIN(aanwezigen_gebruikers);
CREATE INDEX IF NOT EXISTS idx_notulen_versies_afwezigen_gebruikers ON notulen_versies USING GIN(afwezigen_gebruikers);

-- =====================================================
-- 4. HELPER FUNCTIONS FOR USER NAME RESOLUTION
-- =====================================================

-- Function to get user names from UUID array
CREATE OR REPLACE FUNCTION get_user_names_from_uuids(user_uuids UUID[])
RETURNS TEXT[] AS $$
DECLARE
    user_names TEXT[];
BEGIN
    IF user_uuids IS NULL OR array_length(user_uuids, 1) = 0 THEN
        RETURN ARRAY[]::TEXT[];
    END IF;

    SELECT array_agg(naam ORDER BY naam)
    INTO user_names
    FROM gebruikers
    WHERE id = ANY(user_uuids)
    AND is_actief = true;

    RETURN COALESCE(user_names, ARRAY[]::TEXT[]);
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_user_names_from_uuids IS 'Converteert array van user UUIDs naar array van user namen';

-- Function to get user UUIDs from names (best effort matching)
CREATE OR REPLACE FUNCTION get_user_uuids_from_names(user_names TEXT[])
RETURNS UUID[] AS $$
DECLARE
    user_uuids UUID[];
BEGIN
    IF user_names IS NULL OR array_length(user_names, 1) = 0 THEN
        RETURN ARRAY[]::UUID[];
    END IF;

    SELECT array_agg(id ORDER BY naam)
    INTO user_uuids
    FROM gebruikers
    WHERE naam = ANY(user_names)
    AND is_actief = true;

    RETURN COALESCE(user_uuids, ARRAY[]::UUID[]);
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_user_uuids_from_names IS 'Converteert array van user namen naar array van user UUIDs (best effort)';

-- =====================================================
-- 5. VIEWS FOR BACKWARDS COMPATIBILITY AND EASY ACCESS
-- =====================================================

-- Drop existing views to avoid column mismatch errors
DROP VIEW IF EXISTS notulen_with_participants CASCADE;
DROP VIEW IF EXISTS notulen_versies_with_participants CASCADE;

-- View that combines user UUIDs and guest names into single arrays (for API compatibility)
CREATE VIEW notulen_with_participants AS
SELECT
    n.*,
    -- Combined aanwezigen (users + guests)
    CASE
        WHEN n.aanwezigen_gebruikers IS NOT NULL AND n.aanwezigen_gasten IS NOT NULL THEN
            get_user_names_from_uuids(n.aanwezigen_gebruikers) || n.aanwezigen_gasten
        WHEN n.aanwezigen_gebruikers IS NOT NULL THEN
            get_user_names_from_uuids(n.aanwezigen_gebruikers)
        WHEN n.aanwezigen_gasten IS NOT NULL THEN
            n.aanwezigen_gasten
        ELSE
            ARRAY[]::TEXT[]
    END as aanwezigen_combined,
    -- Combined afwezigen (users + guests)
    CASE
        WHEN n.afwezigen_gebruikers IS NOT NULL AND n.afwezigen_gasten IS NOT NULL THEN
            get_user_names_from_uuids(n.afwezigen_gebruikers) || n.afwezigen_gasten
        WHEN n.afwezigen_gebruikers IS NOT NULL THEN
            get_user_names_from_uuids(n.afwezigen_gebruikers)
        WHEN n.afwezigen_gasten IS NOT NULL THEN
            n.afwezigen_gasten
        ELSE
            ARRAY[]::TEXT[]
    END as afwezigen_combined
FROM notulen n;

COMMENT ON VIEW notulen_with_participants IS 'Notulen view met gecombineerde participant lijsten (voor API backwards compatibility)';

-- Similar view for versions
CREATE VIEW notulen_versies_with_participants AS
SELECT
    nv.*,
    -- Combined aanwezigen (users + guests)
    CASE
        WHEN nv.aanwezigen_gebruikers IS NOT NULL AND nv.aanwezigen_gasten IS NOT NULL THEN
            get_user_names_from_uuids(nv.aanwezigen_gebruikers) || nv.aanwezigen_gasten
        WHEN nv.aanwezigen_gebruikers IS NOT NULL THEN
            get_user_names_from_uuids(nv.aanwezigen_gebruikers)
        WHEN nv.aanwezigen_gasten IS NOT NULL THEN
            nv.aanwezigen_gasten
        ELSE
            ARRAY[]::TEXT[]
    END as aanwezigen_combined,
    -- Combined afwezigen (users + guests)
    CASE
        WHEN nv.afwezigen_gebruikers IS NOT NULL AND nv.afwezigen_gasten IS NOT NULL THEN
            get_user_names_from_uuids(nv.afwezigen_gebruikers) || nv.afwezigen_gasten
        WHEN nv.afwezigen_gebruikers IS NOT NULL THEN
            get_user_names_from_uuids(nv.afwezigen_gebruikers)
        WHEN nv.afwezigen_gasten IS NOT NULL THEN
            nv.afwezigen_gasten
        ELSE
            ARRAY[]::TEXT[]
    END as afwezigen_combined
FROM notulen_versies nv;

COMMENT ON VIEW notulen_versies_with_participants IS 'Notulen versies view met gecombineerde participant lijsten';

-- =====================================================
-- MIGRATION NOTES
-- =====================================================
/*
NIEUWE PARTICIPANT STRUCTUUR:

1. VIER KOLOMMEN PER TABEL:
   - aanwezigen_gebruikers: UUID[] - Geregistreerde gebruikers die aanwezig waren
   - afwezigen_gebruikers: UUID[] - Geregistreerde gebruikers die afwezig waren
   - aanwezigen_gasten: TEXT[] - Niet-geregistreerde gasten die aanwezig waren
   - afwezigen_gasten: TEXT[] - Niet-geregistreerde gasten die afwezig waren

2. BACKWARDS COMPATIBILITY:
   - Oude kolommen (aanwezigen, afwezigen) blijven bestaan
   - Views combineren nieuwe kolommen voor API compatibility
   - Applicatie kan geleidelijk migreren naar nieuwe structuur

3. VOORDELEN:
   - Proper user linking met UUIDs
   - Guest support blijft mogelijk
   - Data integrity via foreign keys (optioneel)
   - Betere performance met GIN indexes
   - Version management toont echte namen ipv UUIDs

4. IMPLEMENTATIE STAPPEN:
   - Update models om nieuwe velden toe te voegen
   - Update services om UUIDs te verwerken
   - Update handlers om UUID arrays te accepteren
   - Update frontend om user selectie mogelijk te maken
   - Migrate bestaande data waar mogelijk

5. API CHANGES:
   - Create/Update requests kunnen nu zowel user UUIDs als guest names bevatten
   - Response bevat zowel UUIDs als resolved names
   - Version endpoints gebruiken resolved names ipv raw UUIDs
*/

-- =====================================================
-- MIGRATION COMPLETE
-- =====================================================