-- ============================================================================
-- V28 CONSOLIDATED: Rename Tables - Participant Refactor
-- ============================================================================
-- Consolidates V28_01 through V28_11
-- Purpose: Rename tables from Dutch to English and refactor participant/event
--          relationship to separate person data from event registration data
-- ============================================================================

-- ----------------------------------------------------------------------------
-- PART 1: RENAME MAIN TABLES
-- ----------------------------------------------------------------------------

-- V28_01: Rename 'aanmeldingen' to 'participants'
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'aanmeldingen')
       AND NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'participants') THEN
        ALTER TABLE aanmeldingen RENAME TO participants;
        RAISE NOTICE 'Renamed aanmeldingen to participants';
    ELSIF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'aanmeldingen')
       AND EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'participants') THEN
        DROP VIEW IF EXISTS participant_view CASCADE;
        DROP TABLE participants CASCADE;
        ALTER TABLE aanmeldingen RENAME TO participants;
        RAISE NOTICE 'Dropped existing participants and renamed aanmeldingen';
    ELSIF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'aanmeldingen')
       AND EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'participants') THEN
        RAISE NOTICE 'Migration V28_01: aanmeldingen already renamed to participants';
    ELSE
        RAISE NOTICE 'Migration V28_01: Neither aanmeldingen nor participants table exists';
    END IF;
END $$;

-- V28_02: Rename 'event_participants' to 'event_registrations'
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'event_participants')
       AND NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'event_registrations') THEN
        ALTER TABLE event_participants RENAME TO event_registrations;
        RAISE NOTICE 'Renamed event_participants to event_registrations';
    ELSIF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'event_participants')
       AND EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'event_registrations') THEN
        DROP VIEW IF EXISTS event_participants_view CASCADE;
        DROP TABLE event_registrations CASCADE;
        ALTER TABLE event_participants RENAME TO event_registrations;
        RAISE NOTICE 'Dropped existing event_registrations and renamed event_participants';
    ELSIF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'event_participants')
       AND EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'event_registrations') THEN
        RAISE NOTICE 'Migration V28_02: event_participants already renamed to event_registrations';
    ELSE
        RAISE NOTICE 'Migration V28_02: Neither event_participants nor event_registrations table exists';
    END IF;
END $$;

-- V28_03: Drop old event_participants view if it exists
DROP VIEW IF EXISTS event_participants_view CASCADE;

-- V28_04: Rename user_participant view
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.views WHERE table_name = 'user_participant_view') THEN
        DROP VIEW IF EXISTS participant_user_view CASCADE;
        ALTER VIEW user_participant_view RENAME TO participant_user_view;
        RAISE NOTICE 'Renamed user_participant_view to participant_user_view';
    END IF;
END $$;

-- ----------------------------------------------------------------------------
-- PART 2: ADD EVENT-SPECIFIC COLUMNS TO EVENT_REGISTRATIONS
-- ----------------------------------------------------------------------------

-- V28_05: Add event-specific columns (moved from participants table)
ALTER TABLE event_registrations
    ADD COLUMN IF NOT EXISTS participant_role_name TEXT REFERENCES participant_roles(name) ON UPDATE CASCADE,
    ADD COLUMN IF NOT EXISTS distance_route TEXT REFERENCES distances(route) ON UPDATE CASCADE,
    ADD COLUMN IF NOT EXISTS status TEXT REFERENCES registration_status_types(status) ON UPDATE CASCADE,
    ADD COLUMN IF NOT EXISTS bijzonderheden TEXT,
    ADD COLUMN IF NOT EXISTS terms BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS notities TEXT,
    ADD COLUMN IF NOT EXISTS antwoorden_count INTEGER DEFAULT 0;

-- V28_06: Rename 'stappen' column to 'steps' in event_registrations
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'event_registrations' AND column_name = 'stappen'
    ) THEN
        ALTER TABLE event_registrations RENAME COLUMN stappen TO steps;
        RAISE NOTICE 'Renamed stappen to steps in event_registrations';
    END IF;
END $$;

-- ----------------------------------------------------------------------------
-- PART 3: MIGRATE DATA FROM PARTICIPANTS TO EVENT_REGISTRATIONS
-- ----------------------------------------------------------------------------

-- V28_07: Migrate participant data to event_registrations for active event
DO $$
DECLARE
    v_active_event_id UUID;
BEGIN
    -- Find the active event
    SELECT get_active_event() INTO v_active_event_id;
    
    IF v_active_event_id IS NULL THEN
        RAISE WARNING '[V28] Geen actief event gevonden. Data migratie overgeslagen.';
        RETURN;
    END IF;

    RAISE NOTICE '[V28] Actief event ID: %. Migreren van data...', v_active_event_id;

    -- Copy data from participants to event_registrations
    INSERT INTO event_registrations (
        event_id, 
        participant_id, 
        registered_at, 
        tracking_status,
        -- steps, status,
        distance_route,
        participant_role_name,
        bijzonderheden,
        terms,
        notities,
        antwoorden_count
    )
    SELECT
        v_active_event_id,           -- Link to active event
        p.id,                         -- Participant ID
        p.created_at,                 -- Original registration date
        'registered',                 -- Default tracking status
        -- COALESCE(p.steps, 0),         -- COALESCE(p.status, 'registered'), -- Status (from old table)
        p.distance_route,             -- Distance (from V26)
        p.participant_role_name,      -- Role (from V26)
        p.bijzonderheden,
        COALESCE(p.terms, false),
        p.notities,
        COALESCE(p.antwoorden_count, 0)
    FROM participants p
    WHERE NOT EXISTS (
        SELECT 1 FROM event_registrations er
        WHERE er.participant_id = p.id AND er.event_id = v_active_event_id
    );

    RAISE NOTICE '[V28] Data migratie voltooid.';
END $$;

-- ----------------------------------------------------------------------------
-- PART 4: RENAME ANTWOORDEN TABLES
-- ----------------------------------------------------------------------------

-- V28_08: Rename 'aanmelding_antwoorden' to 'participant_antwoorden'
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'aanmelding_antwoorden')
       AND NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'participant_antwoorden') THEN
        ALTER TABLE aanmelding_antwoorden RENAME TO participant_antwoorden;
        RAISE NOTICE 'Renamed aanmelding_antwoorden to participant_antwoorden';
    ELSIF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'aanmelding_antwoorden')
       AND EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'participant_antwoorden') THEN
        DROP VIEW IF EXISTS participant_antwoorden_view CASCADE;
        DROP TABLE participant_antwoorden CASCADE;
        ALTER TABLE aanmelding_antwoorden RENAME TO participant_antwoorden;
        RAISE NOTICE 'Dropped existing participant_antwoorden and renamed aanmelding_antwoorden';
    ELSIF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'aanmelding_antwoorden')
       AND EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'participant_antwoorden') THEN
        RAISE NOTICE 'Migration V28_08: aanmelding_antwoorden already renamed';
    ELSE
        RAISE NOTICE 'Migration V28_08: Neither table exists';
    END IF;
END $$;

-- V28_09: Fix foreign key in participant_antwoorden
DO $$
BEGIN
    -- Drop old FK if exists
    ALTER TABLE participant_antwoorden 
        DROP CONSTRAINT IF EXISTS aanmelding_antwoorden_aanmelding_id_fkey;
    
    -- Rename column if needed
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'participant_antwoorden' AND column_name = 'aanmelding_id'
    ) THEN
        ALTER TABLE participant_antwoorden RENAME COLUMN aanmelding_id TO participant_id;
    END IF;
    
    -- Add new FK
    ALTER TABLE participant_antwoorden
        DROP CONSTRAINT IF EXISTS fk_participant_antwoorden_participant,
        ADD CONSTRAINT fk_participant_antwoorden_participant
        FOREIGN KEY (participant_id) REFERENCES participants(id)
        ON UPDATE CASCADE ON DELETE CASCADE;
    
    RAISE NOTICE 'Fixed FK in participant_antwoorden';
END $$;

-- V28_10: Add comment to participant_antwoorden table
COMMENT ON TABLE participant_antwoorden IS 'Antwoorden van staff op participant vragen (V1, V28)';

-- ----------------------------------------------------------------------------
-- COMPLETION LOG
-- ----------------------------------------------------------------------------

-- V28_11: Log completion
DO $$
BEGIN
    RAISE NOTICE '[V28 CONSOLIDATED] Tabel hernoeming en participant refactor voltooid.';
    RAISE NOTICE '=== BREAKING CHANGES ===';
    RAISE NOTICE '1. Tabellen hernoemd:';
    RAISE NOTICE '   - aanmeldingen → participants';
    RAISE NOTICE '   - event_participants → event_registrations';
    RAISE NOTICE '   - aanmelding_antwoorden → participant_antwoorden';
    RAISE NOTICE '2. Data separation:';
    RAISE NOTICE '   - participants: Alleen persoonsgegevens (naam, email, telefoon)';
    RAISE NOTICE '   - event_registrations: Event-specifieke data (status, rol, afstand, steps)';
    RAISE NOTICE '3. (Go Code) Update alle GORM TableName() methods en references.';
END $$;