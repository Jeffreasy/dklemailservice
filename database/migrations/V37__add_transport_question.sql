-- V37: Add transport question to event registrations
-- Date: 2025-11-17
-- Purpose: Add "Heb je vervoer?" boolean question to event registration form

-- Add transport question column to event_registrations table
ALTER TABLE event_registrations
ADD COLUMN IF NOT EXISTS heeft_vervoer BOOLEAN;

-- Add index for performance on transport question queries
CREATE INDEX IF NOT EXISTS idx_event_registrations_heeft_vervoer
ON event_registrations(heeft_vervoer);

-- Add comment for documentation
COMMENT ON COLUMN event_registrations.heeft_vervoer
IS 'V37: Geeft aan of participant eigen vervoer heeft (Ja/Nee vraag bij registratie)';

-- Verification query (optional - can be removed after testing)
DO $$
DECLARE
    column_exists BOOLEAN;
BEGIN
    SELECT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'event_registrations'
        AND column_name = 'heeft_vervoer'
    ) INTO column_exists;

    IF column_exists THEN
        RAISE NOTICE 'V37: heeft_vervoer column successfully added to event_registrations';
    ELSE
        RAISE EXCEPTION 'V37: Failed to add heeft_vervoer column';
    END IF;
END $$;