-- ============================================================================
-- V29 CONSOLIDATED: Update Permissions for Participants
-- ============================================================================
-- Consolidates V29 (single file)
-- Purpose: Update permission resource name from 'aanmelding' to 'participant'
--          to align with V28 table renaming
-- ============================================================================

-- START FIX: Verwijder eerst de oude 'aanmelding' permissies
-- die al als 'participant' permissies bestaan om duplicates te voorkomen.
DELETE FROM permissions p
WHERE p.resource = 'aanmelding'
  AND EXISTS (
    SELECT 1 
    FROM permissions p_new
    WHERE p_new.resource = 'participant' 
      AND p_new.action = p.action
  );

-- Update de resterende permissies
UPDATE permissions
SET resource = 'participant'
WHERE resource = 'aanmelding';
-- EINDE FIX

-- Log completion
DO $$
BEGIN
    RAISE NOTICE '[V29] Permissie resource ''aanmelding'' hernoemd naar ''participant''';
    RAISE NOTICE 'Dit matcht met de V28 tabel hernoeming (aanmeldingen → participants)';
END $$;