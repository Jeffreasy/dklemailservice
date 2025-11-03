-- Script: create_test_event_dronten.sql
-- Description: Maakt test event aan voor GPS/geofencing testing in Dronten
-- Location: Spiegelstraat 6, Dronten (52.5185, 5.7220)
-- Date: 2025-11-03
--
-- USAGE:
-- docker exec -i dkl-postgres psql -U dekoninklijkeloopdatabase_user -d dekoninklijkeloopdatabase < database/scripts/create_test_event_dronten.sql

-- =====================================================
-- 1. DELETE OLD TEST EVENTS (CLEANUP)
-- =====================================================
-- Verwijder oude test events om duplicates te voorkomen
DELETE FROM events 
WHERE name LIKE '%Test%' 
OR (name = 'De Koninklijke Loop 2025' AND description LIKE '%Test%');

-- =====================================================
-- 2. CREATE ACTIVE TEST EVENT
-- =====================================================
-- Event is ACTIEF (start in verleden, end in toekomst)
-- Locatie: Dronten, Spiegelstraat 6
INSERT INTO events (
    name,
    description,
    start_time,
    end_time,
    status,
    is_active,
    geofences,
    event_config
) VALUES (
    'De Koninklijke Loop 2025 - GPS Test',
    'Test event voor GPS geofencing in Dronten',
    '2025-01-01 00:00:00+00',  -- Start: 1 jan (in verleden = altijd actief)
    '2025-12-31 23:59:59+00',  -- End: 31 dec (hele jaar actief)
    'active',                   -- Status: ACTIEF
    true,                       -- Is active: JA
    '[
        {
            "type": "start",
            "lat": 52.5185,
            "long": 5.7220,
            "radius": 500,
            "name": "Start - Dronten Spiegelstraat"
        },
        {
            "type": "checkpoint",
            "lat": 52.5228,
            "long": 5.7306,
            "radius": 300,
            "name": "Checkpoint Noord"
        },
        {
            "type": "checkpoint",
            "lat": 52.5142,
            "long": 5.7134,
            "radius": 300,
            "name": "Checkpoint Zuid"
        },
        {
            "type": "finish",
            "lat": 52.5185,
            "long": 5.7220,
            "radius": 500,
            "name": "Finish - Dronten Spiegelstraat"
        }
    ]'::jsonb,
    '{
        "minStepsInterval": 10,
        "requireGeofenceCheckin": true,
        "distanceThreshold": 100,
        "accuracyLevel": "balanced"
    }'::jsonb
)
ON CONFLICT DO NOTHING;

-- =====================================================
-- 3. VERIFICATION
-- =====================================================
-- Check dat event correct is aangemaakt
DO $$
DECLARE
    v_event_id UUID;
    v_geofence_count INT;
BEGIN
    -- Haal active event op
    SELECT id INTO v_event_id FROM get_active_event();
    
    IF v_event_id IS NULL THEN
        RAISE NOTICE '❌ FOUT: Geen actief event gevonden!';
    ELSE
        -- Check geofences
        SELECT jsonb_array_length(geofences) INTO v_geofence_count
        FROM events WHERE id = v_event_id;
        
        RAISE NOTICE '✅ SUCCESS: Test event aangemaakt!';
        RAISE NOTICE 'Event ID: %', v_event_id;
        RAISE NOTICE 'Geofences: % stuks', v_geofence_count;
        RAISE NOTICE '';
        RAISE NOTICE 'Test met:';
        RAISE NOTICE 'curl http://localhost:8082/api/events/active';
        RAISE NOTICE 'OF (production):';
        RAISE NOTICE 'curl https://dklemailservice.onrender.com/api/events/active';
    END IF;
END $$;

-- =====================================================
-- 4. DISPLAY EVENT DATA
-- =====================================================
SELECT 
    id,
    name,
    status,
    is_active,
    start_time,
    end_time,
    jsonb_array_length(geofences) as geofence_count,
    event_config->>'accuracyLevel' as accuracy
FROM events
WHERE is_active = true
ORDER BY created_at DESC
LIMIT 1;

-- =====================================================
-- SCRIPT COMPLETE
-- =====================================================