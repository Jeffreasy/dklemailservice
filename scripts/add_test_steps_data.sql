-- Script om test data toe te voegen voor steps/leaderboard testing
-- DKL Email Service - SIMPLIFIED VERSION

BEGIN;

-- 1. Voeg test event toe en bewaar het ID
WITH new_event AS (
    INSERT INTO events (name, description, start_time, end_time, is_active, status)
    VALUES (
        'Test Wandel Event 2025',
        'Test event voor steps en leaderboard testing',
        '2025-01-01 10:00:00+00',
        '2025-12-31 18:00:00+00',
        true,
        'active'
    )
    ON CONFLICT DO NOTHING
    RETURNING id
)
SELECT id FROM new_event;

-- 2. Voeg test participants toe (met vaste UUIDs voor makkelijke referentie)
INSERT INTO participants (id, naam, email, telefoon, terms, test_mode)
VALUES 
    ('11111111-1111-1111-1111-111111111111', 'Jan de Tester', 'jan.tester@example.com', '0612345678', true, true),
    ('22222222-2222-2222-2222-222222222222', 'Marie Runner', 'marie.runner@example.com', '0687654321', true, true),
    ('33333333-3333-3333-3333-333333333333', 'Peter Wandelaar', 'peter.wandelaar@example.com', '0698765432', true, true),
    ('44444444-4444-4444-4444-444444444444', 'Sophia Sprint', 'sophia.sprint@example.com', '0676543210', true, true),
    ('55555555-5555-5555-5555-555555555555', 'Lucas Marathon', 'lucas.marathon@example.com', '0654321098', true, true)
ON CONFLICT (id) DO NOTHING;

-- 3. Voeg event registrations toe met steps (waar de ECHTE steps data staat!)
INSERT INTO event_registrations (
    event_id, 
    participant_id,
    distance_route,
    status,
    steps,
    terms,
    test_mode
)
SELECT 
    (SELECT id FROM events WHERE name = 'Test Wandel Event 2025' ORDER BY created_at DESC LIMIT 1),
    participant_id,
    distance_route,
    'confirmed',
    steps,
    true,
    true
FROM (VALUES
    ('11111111-1111-1111-1111-111111111111'::uuid, '10 KM', 15000),
    ('22222222-2222-2222-2222-222222222222'::uuid, '15 KM', 23500),
    ('33333333-3333-3333-3333-333333333333'::uuid, '6 KM', 8900),
    ('44444444-4444-4444-4444-444444444444'::uuid, '20 KM', 31250),
    ('55555555-5555-5555-5555-555555555555'::uuid, '15 KM', 19800)
) AS data(participant_id, distance_route, steps)
ON CONFLICT DO NOTHING;

COMMIT;

-- Verificatie queries
SELECT '=== Total Test Events ===' as info;
SELECT COUNT(*) as count FROM events WHERE name LIKE 'Test%';

SELECT '=== Total Test Participants ===' as info;
SELECT COUNT(*) as count FROM participants WHERE test_mode = true;

SELECT '=== Total Test Registrations ===' as info;
SELECT COUNT(*) as count FROM event_registrations WHERE test_mode = true;

SELECT '=== Total Steps (from event_registrations) ===' as info;
SELECT COALESCE(SUM(steps), 0) as total_steps FROM event_registrations WHERE test_mode = true;

SELECT '=== Leaderboard Preview (Top 5) ===' as info;
SELECT 
    ROW_NUMBER() OVER (ORDER BY er.steps DESC) as rank,
    p.naam,
    er.steps,
    er.distance_route as route
FROM event_registrations er
JOIN participants p ON p.id = er.participant_id
WHERE er.test_mode = true AND er.steps > 0
ORDER BY er.steps DESC
LIMIT 5;