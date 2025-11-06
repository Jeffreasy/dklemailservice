-- V1_59__fix_notulen_version_complete.sql
-- Fix: Include UUID participant fields in version snapshots
-- This migration ensures all participant data (UUIDs and guest names) is properly versioned

-- Drop the existing trigger and function
DROP TRIGGER IF EXISTS notulen_version_trigger ON notulen;
DROP FUNCTION IF EXISTS create_notulen_version() CASCADE;

-- Recreate the function with complete field coverage
CREATE OR REPLACE FUNCTION create_notulen_version()
RETURNS TRIGGER AS $$
BEGIN
    -- Only create version if this is an update (not insert)
    IF TG_OP = 'UPDATE' THEN
        -- Insert current version into history table with ALL fields
        INSERT INTO notulen_versies (
            notulen_id, versie, titel, vergadering_datum, locatie,
            voorzitter, notulist, 
            aanwezigen, afwezigen,  -- Legacy fields
            aanwezigen_gebruikers, afwezigen_gebruikers,  -- NEW: UUID arrays
            aanwezigen_gasten, afwezigen_gasten,          -- NEW: Guest text arrays
            agenda_items, besluiten, actiepunten, notities, status,
            gewijzigd_door, wijziging_reden
        ) VALUES (
            OLD.id, OLD.versie, OLD.titel, OLD.vergadering_datum, OLD.locatie,
            OLD.voorzitter, OLD.notulist, 
            OLD.aanwezigen, OLD.afwezigen,  -- Legacy fields
            OLD.aanwezigen_gebruikers, OLD.afwezigen_gebruikers,  -- NEW: UUID arrays
            OLD.aanwezigen_gasten, OLD.afwezigen_gasten,          -- NEW: Guest text arrays
            OLD.agenda_items, OLD.besluiten, OLD.actiepunten, OLD.notities, OLD.status,
            COALESCE(NEW.updated_by, NEW.created_by), 'Automatic version snapshot'
        );

        -- Increment version number
        NEW.versie = OLD.versie + 1;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Recreate the trigger
CREATE TRIGGER notulen_version_trigger
BEFORE UPDATE ON notulen
FOR EACH ROW
EXECUTE FUNCTION create_notulen_version();

-- Add comment for documentation
COMMENT ON FUNCTION create_notulen_version IS 'Creates version snapshot with complete participant data (UUIDs + guests)';