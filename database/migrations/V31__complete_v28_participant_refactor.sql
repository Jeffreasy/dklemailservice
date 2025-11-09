-- ============================================================================
-- V31: Complete V28 Participant Refactor
-- ============================================================================
-- Purpose: Fix V28 migration failure and complete the participant/event
--          data separation according to the architecture design
-- ============================================================================

-- ----------------------------------------------------------------------------
-- PART 1: ADD MISSING COLUMNS TO EVENT_REGISTRATIONS
-- ----------------------------------------------------------------------------

ALTER TABLE event_registrations
    ADD COLUMN IF NOT EXISTS steps INTEGER DEFAULT 0,
    ADD COLUMN IF NOT EXISTS ondersteuning TEXT,
    ADD COLUMN IF NOT EXISTS test_mode BOOLEAN DEFAULT false,
    ADD COLUMN IF NOT EXISTS email_verzonden BOOLEAN DEFAULT false,
    ADD COLUMN IF NOT EXISTS email_verzonden_op TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS behandeld_door TEXT,
    ADD COLUMN IF NOT EXISTS behandeld_op TIMESTAMPTZ;

-- Rename current_steps to steps if it exists
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'event_registrations' AND column_name = 'current_steps'
    ) AND NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'event_registrations' AND column_name = 'steps'
    ) THEN
        ALTER TABLE event_registrations RENAME COLUMN current_steps TO steps;
        RAISE NOTICE 'Renamed current_steps to steps';
    ELSIF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'event_registrations' AND column_name = 'current_steps'
    ) AND EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'event_registrations' AND column_name = 'steps'
    ) THEN
        -- Both exist, migrate data and drop current_steps
        UPDATE event_registrations 
        SET steps = COALESCE(current_steps, steps, 0)
        WHERE steps = 0 AND current_steps > 0;
        
        ALTER TABLE event_registrations DROP COLUMN current_steps;
        RAISE NOTICE 'Merged current_steps into steps and dropped current_steps';
    END IF;
END $$;

COMMENT ON COLUMN event_registrations.steps IS 'Number of steps taken during event (migrated from participants table in V28/V31)';
COMMENT ON COLUMN event_registrations.ondersteuning IS 'Support/assistance needs (migrated from participants table in V28/V31)';

-- ----------------------------------------------------------------------------
-- PART 2: MIGRATE REMAINING DATA
-- ----------------------------------------------------------------------------

-- Only migrate if we have an active event and there's data to migrate
DO $$
DECLARE
    v_active_event_id UUID;
    v_migrated_count INTEGER := 0;
BEGIN
    -- Find active event
    SELECT get_active_event() INTO v_active_event_id;
    
    IF v_active_event_id IS NULL THEN
        RAISE WARNING '[V31] No active event found. Cannot migrate participant data.';
        RAISE NOTICE '[V31] You can manually create event_registrations or set an active event.';
        RETURN;
    END IF;

    -- Migrate data for participants that don't have event_registrations yet
    INSERT INTO event_registrations (
        event_id,
        participant_id,
        registered_at,
        tracking_status,
        steps,
        ondersteuning,
        bijzonderheden,
        terms,
        notities,
        status,
        distance_route,
        participant_role_name,
        test_mode,
        email_verzonden,
        email_verzonden_op,
        behandeld_door,
        behandeld_op,
        antwoorden_count
    )
    SELECT
        v_active_event_id,
        p.id,
        p.created_at,
        'registered',
        COALESCE(p.steps, 0),
        p.ondersteuning,
        p.bijzonderheden,
        COALESCE(p.terms, false),
        p.notities,
        COALESCE(p.status, 'registered'),
        p.distance_route,
        p.participant_role_name,
        COALESCE(p.test_mode, false),
        COALESCE(p.email_verzonden, false),
        p.email_verzonden_op,
        p.behandeld_door,
        p.behandeld_op,
        (SELECT COUNT(*) FROM participant_antwoorden WHERE participant_id = p.id)
    FROM participants p
    WHERE NOT EXISTS (
        SELECT 1 FROM event_registrations er
        WHERE er.participant_id = p.id AND er.event_id = v_active_event_id
    );

    GET DIAGNOSTICS v_migrated_count = ROW_COUNT;
    RAISE NOTICE '[V31] Migrated % participants to event_registrations', v_migrated_count;
END $$;

-- ----------------------------------------------------------------------------
-- PART 3: REMOVE OLD COLUMNS FROM PARTICIPANTS (Optional - commented out)
-- ----------------------------------------------------------------------------

-- UNCOMMENT THESE LINES AFTER VERIFYING DATA MIGRATION IS SUCCESSFUL:
-- 
-- Optionally remove old columns from participants table
-- These are kept for now for backward compatibility
-- 
-- ALTER TABLE participants 
--     DROP COLUMN IF EXISTS afstand,
--     DROP COLUMN IF EXISTS rol,
--     DROP COLUMN IF EXISTS ondersteuning,
--     DROP COLUMN IF EXISTS bijzonderheden,
--     DROP COLUMN IF EXISTS steps,
--     DROP COLUMN IF EXISTS status,
--     DROP COLUMN IF EXISTS email_verzonden,
--     DROP COLUMN IF EXISTS email_verzonden_op,
--     DROP COLUMN IF EXISTS behandeld_door,
--     DROP COLUMN IF EXISTS behandeld_op,
--     DROP COLUMN IF EXISTS notities,
--     DROP COLUMN IF NOT EXISTS antwoorden_count;
-- 
-- RAISE NOTICE '[V31] Removed old columns from participants table';

-- ----------------------------------------------------------------------------
-- PART 4: ADD MISSING INDEXES
-- ----------------------------------------------------------------------------

CREATE INDEX IF NOT EXISTS idx_event_registrations_participant_role 
    ON event_registrations(participant_role_name);

CREATE INDEX IF NOT EXISTS idx_event_registrations_distance 
    ON event_registrations(distance_route);

CREATE INDEX IF NOT EXISTS idx_event_registrations_steps 
    ON event_registrations(steps) 
    WHERE steps > 0;

-- ----------------------------------------------------------------------------
-- COMPLETION LOG
-- ----------------------------------------------------------------------------

DO $$
BEGIN
    RAISE NOTICE '✅ [V31] Participant refactor completion successful!';
    RAISE NOTICE '=== WHAT WAS DONE ===';
    RAISE NOTICE '1. Added missing columns to event_registrations (steps, ondersteuning, etc.)';
    RAISE NOTICE '2. Migrated data from participants to event_registrations';
    RAISE NOTICE '3. Added performance indexes';
    RAISE NOTICE '';
    RAISE NOTICE '=== NEXT STEPS ===';
    RAISE NOTICE '1. Verify data migration: SELECT COUNT(*) FROM event_registrations;';
    RAISE NOTICE '2. Test API endpoints with new schema';
    RAISE NOTICE '3. After verification, uncomment PART 3 to remove old columns';
    RAISE NOTICE '4. Update Go models to match new schema';
    RAISE NOTICE '';
    RAISE NOTICE '⚠️  Old columns on participants are KEPT for backward compatibility';
    RAISE NOTICE '⚠️  Remove them after thorough testing by uncommenting PART 3';
END $$;