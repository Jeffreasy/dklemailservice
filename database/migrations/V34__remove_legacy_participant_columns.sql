-- V34: Remove legacy columns from participants table
-- These columns were moved to event_registrations in V28
-- but the columns themselves were never dropped

-- STEP 1: Drop legacy indexes first (these prevent column drops)
DROP INDEX IF EXISTS idx_aanmeldingen_rol CASCADE;
DROP INDEX IF EXISTS idx_aanmeldingen_afstand CASCADE;
DROP INDEX IF EXISTS idx_aanmeldingen_status CASCADE;

-- STEP 2: Remove legacy event-specific columns that now belong in event_registrations
-- Using CASCADE to force drop any remaining dependencies
ALTER TABLE participants 
DROP COLUMN IF EXISTS rol CASCADE,
DROP COLUMN IF EXISTS afstand CASCADE,
DROP COLUMN IF EXISTS ondersteuning CASCADE,
DROP COLUMN IF EXISTS bijzonderheden CASCADE,
DROP COLUMN IF EXISTS email_verzonden CASCADE,
DROP COLUMN IF EXISTS email_verzonden_op CASCADE,
DROP COLUMN IF EXISTS behandeld_door CASCADE,
DROP COLUMN IF EXISTS behandeld_op CASCADE,
DROP COLUMN IF EXISTS notities CASCADE,
DROP COLUMN IF EXISTS steps CASCADE,
DROP COLUMN IF EXISTS antwoorden_count CASCADE,
DROP COLUMN IF EXISTS participant_role_name CASCADE,
DROP COLUMN IF EXISTS distance_route CASCADE,
DROP COLUMN IF EXISTS status CASCADE;

-- These fields now exist ONLY in event_registrations table
-- participants table should only contain person information + account type info

COMMENT ON TABLE participants IS 'V34: Cleaned up - removed legacy event-specific columns. Person data only + V30 account type fields.';