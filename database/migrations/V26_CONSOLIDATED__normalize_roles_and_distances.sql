-- ============================================================================
-- V26 CONSOLIDATED: Normalize Participant Roles and Distances
-- ============================================================================
-- Consolidates V26_01 through V26_14
-- Purpose: Create lookup tables for participant roles and distances,
--          migrate data from string fields to normalized foreign keys
-- ============================================================================

-- ----------------------------------------------------------------------------
-- PART 1: PARTICIPANT ROLES NORMALIZATION
-- ----------------------------------------------------------------------------

-- V26_01: Create participant_roles lookup table
CREATE TABLE IF NOT EXISTS participant_roles (
    name TEXT PRIMARY KEY NOT NULL,
    description TEXT
);

-- V26_02: Add table comment
COMMENT ON TABLE participant_roles IS 'Bron van waarheid voor deelnemer-rollen (Deelnemer, Vrijwilliger, etc) (V26)';

-- V26_03: Seed participant roles
INSERT INTO participant_roles (name, description) VALUES
('Deelnemer', 'Een standaard deelnemer aan het evenement.'),
('Begeleider', 'Een begeleider van een of meerdere deelnemers.'),
('Vrijwilliger', 'Een vrijwilliger die helpt bij het evenement.'),
('Sponsor', 'Een sponsor of partner (indien deze zich kunnen aanmelden).')
ON CONFLICT (name) DO NOTHING;

-- V26_04: Add new FK column to aanmeldingen (will become participants in V28)
ALTER TABLE aanmeldingen
    ADD COLUMN IF NOT EXISTS participant_role_name TEXT;

-- V26_05: Migrate role data from old 'rol' column to new column
UPDATE aanmeldingen
SET participant_role_name = 
    CASE
        WHEN LOWER(TRIM(rol)) = 'deelnemer' THEN 'Deelnemer'
        WHEN LOWER(TRIM(rol)) = 'begeleider' THEN 'Begeleider'
        WHEN LOWER(TRIM(rol)) = 'vrijwilliger' THEN 'Vrijwilliger'
        WHEN LOWER(TRIM(rol)) = 'sponsor' THEN 'Sponsor'
        ELSE NULL
    END
WHERE participant_role_name IS NULL AND rol IS NOT NULL;

-- V26_06: Add foreign key constraint
ALTER TABLE aanmeldingen
    DROP CONSTRAINT IF EXISTS fk_aanmeldingen_participant_role,
    ADD CONSTRAINT fk_aanmeldingen_participant_role
    FOREIGN KEY (participant_role_name) REFERENCES participant_roles(name)
    ON UPDATE CASCADE ON DELETE SET NULL;

-- ----------------------------------------------------------------------------
-- PART 2: DISTANCES NORMALIZATION
-- ----------------------------------------------------------------------------

-- V26_07: Rename route_funds table to distances (with idempotency)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'route_funds')
       AND NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'distances') THEN
        ALTER TABLE route_funds RENAME TO distances;
    ELSIF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'route_funds')
       AND EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'distances') THEN
        DROP VIEW IF EXISTS event_participants_view CASCADE;
        DROP TABLE distances CASCADE;
        ALTER TABLE route_funds RENAME TO distances;
    ELSIF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'route_funds')
       AND EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'distances') THEN
        RAISE NOTICE 'Migration V26_07: route_funds already renamed to distances';
    ELSE
        RAISE NOTICE 'Migration V26_07: Neither route_funds nor distances table exists';
    END IF;
END $$;

-- V26_08: Rename 'amount' column to 'fund_amount'
ALTER TABLE distances RENAME COLUMN amount TO fund_amount;

-- V26_09: Make 'route' the primary key (drop old constraints first)
ALTER TABLE distances
    DROP CONSTRAINT IF EXISTS route_funds_pkey,
    DROP CONSTRAINT IF EXISTS route_funds_route_key,
    ADD PRIMARY KEY (route);

-- V26_10: Add table comment
COMMENT ON TABLE distances IS 'Bron van waarheid voor alle afstanden en hun fondsenwerving (V17, V26)';

-- V26_11: Add new FK column to aanmeldingen
ALTER TABLE aanmeldingen
    ADD COLUMN IF NOT EXISTS distance_route TEXT;

-- V26_12: Migrate distance data from old 'afstand' column
UPDATE aanmeldingen
SET distance_route = UPPER(TRIM(afstand))
WHERE distance_route IS NULL AND afstand IS NOT NULL AND TRIM(afstand) != '';

-- V26_13: Add foreign key constraint
ALTER TABLE aanmeldingen
    DROP CONSTRAINT IF EXISTS fk_aanmeldingen_distance,
    ADD CONSTRAINT fk_aanmeldingen_distance
    FOREIGN KEY (distance_route) REFERENCES distances(route)
    ON UPDATE CASCADE ON DELETE SET NULL;

-- V26_14: Log completion and instructions
DO $$
BEGIN
    RAISE NOTICE '[V26 CONSOLIDATED] Normalisatie van rollen en afstanden is voltooid.';
    RAISE NOTICE '=== BREAKING CHANGE ===';
    RAISE NOTICE '1. (Go Code) Pas je `Aanmelding` GORM-model aan.';
    RAISE NOTICE '   - Vervang `Rol string` door `ParticipantRoleName string`';
    RAISE NOTICE '   - Vervang `Afstand string` door `DistanceRoute string`';
    RAISE NOTICE '2. Old columns (rol, afstand) kunnen nu veilig worden verwijderd als niet meer nodig.';
END $$;