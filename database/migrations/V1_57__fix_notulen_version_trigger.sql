-- V1_57__fix_notulen_version_trigger.sql
-- Fix the create_notulen_version trigger to use created_by instead of updated_by

-- Drop the existing trigger and function
DROP TRIGGER IF EXISTS notulen_version_trigger ON notulen;
DROP FUNCTION IF EXISTS create_notulen_version() CASCADE;

-- Recreate the function with the fix
CREATE OR REPLACE FUNCTION create_notulen_version()
RETURNS TRIGGER AS $$
BEGIN
    -- Only create version if this is an update (not insert)
    IF TG_OP = 'UPDATE' THEN
        -- Insert current version into history table
        INSERT INTO notulen_versies (
            notulen_id, versie, titel, vergadering_datum, locatie,
            voorzitter, notulist, aanwezigen, afwezigen,
            agenda_items, besluiten, actiepunten, notities, status,
            gewijzigd_door, wijziging_reden
        ) VALUES (
            OLD.id, OLD.versie, OLD.titel, OLD.vergadering_datum, OLD.locatie,
            OLD.voorzitter, OLD.notulist, OLD.aanwezigen, OLD.afwezigen,
            OLD.agenda_items, OLD.besluiten, OLD.actiepunten, OLD.notities, OLD.status,
            NEW.created_by, 'Automatic version snapshot'  -- Changed from NEW.updated_by to NEW.created_by
        );

        -- Increment version number
        NEW.versie = OLD.versie + 1;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Recreate the trigger
CREATE OR REPLACE TRIGGER notulen_version_trigger
BEFORE UPDATE ON notulen
FOR EACH ROW
EXECUTE FUNCTION create_notulen_version();