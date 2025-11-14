-- Check current user and participant data for user 19ff57b7-bb7a-4658-b3c4-28df672cdee1

-- 1. Check if user exists and has participant
SELECT 
    g.id as user_id,
    g.email as user_email,
    p.id as participant_id,
    p.naam,
    p.email as participant_email,
    p.account_type,
    p.has_app_access
FROM gebruikers g
LEFT JOIN participants p ON p.gebruiker_id = g.id
WHERE g.id = '19ff57b7-bb7a-4658-b3c4-28df672cdee1';

-- 2. Check if there's an event_registration
SELECT 
    er.id,
    er.participant_id,
    er.event_id,
    er.status,
    er.steps,
    er.distance_route
FROM event_registrations er
JOIN participants p ON er.participant_id = p.id
WHERE p.gebruiker_id = '19ff57b7-bb7a-4658-b3c4-28df672cdee1';

-- 3. Check if there's an active event
SELECT id, name, start_time, end_time, is_active
FROM events
WHERE is_active = true
ORDER BY created_at DESC
LIMIT 1;

-- 4. If user has participant but no event_registration, create one
-- This assumes there's an active event
DO $$
DECLARE
    v_participant_id UUID;
    v_active_event_id UUID;
    v_gebruiker_id UUID := '19ff57b7-bb7a-4658-b3c4-28df672cdee1';
BEGIN
    -- Get participant_id for this user
    SELECT p.id INTO v_participant_id
    FROM participants p
    WHERE p.gebruiker_id = v_gebruiker_id;
    
    -- Get active event
    SELECT id INTO v_active_event_id
    FROM events
    WHERE is_active = true
    ORDER BY created_at DESC
    LIMIT 1;
    
    -- Only create if participant exists, event exists, and no registration exists yet
    IF v_participant_id IS NOT NULL AND v_active_event_id IS NOT NULL THEN
        IF NOT EXISTS (
            SELECT 1 FROM event_registrations 
            WHERE participant_id = v_participant_id 
            AND event_id = v_active_event_id
        ) THEN
            INSERT INTO event_registrations (
                participant_id,
                event_id,
                status,
                steps,
                distance_route,
                participant_role_name,
                registered_at
            ) VALUES (
                v_participant_id,
                v_active_event_id,
                'bevestigd',  -- Default confirmed status
                0,            -- Start with 0 steps
                '5km',        -- Default route
                'deelnemer',  -- Default role
                CURRENT_TIMESTAMP
            );
            
            RAISE NOTICE 'Event registration created for user %', v_gebruiker_id;
        ELSE
            RAISE NOTICE 'Event registration already exists for user %', v_gebruiker_id;
        END IF;
    ELSE
        IF v_participant_id IS NULL THEN
            RAISE WARNING 'No participant found for user %', v_gebruiker_id;
        END IF;
        IF v_active_event_id IS NULL THEN
            RAISE WARNING 'No active event found';
        END IF;
    END IF;
END $$;

-- 5. Verify the fix
SELECT 
    'AFTER FIX' as status,
    er.id as registration_id,
    p.naam,
    p.email,
    er.status,
    er.steps,
    er.distance_route,
    e.name as event_name
FROM event_registrations er
JOIN participants p ON er.participant_id = p.id
JOIN events e ON er.event_id = e.id
WHERE p.gebruiker_id = '19ff57b7-bb7a-4658-b3c4-28df672cdee1';