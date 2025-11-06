-- GECONSOLIDEERDE V23 - EVENTS & TRACKING
-- Logica van V1_53.
-- OPTIMALISATIE:
-- 1. 'CREATE OR REPLACE FUNCTION update_updated_at_column' is verwijderd (bestaat al uit V19).
-- 2. 'ADD CONSTRAINT unique_resource_action' is verwijderd (bestaat al uit V6).
-- 3. 'ADD CONSTRAINT unique_role_permission' is verwijderd (bestaat al uit V6).

-- =====================================================
-- 1. EVENTS TABLE
-- =====================================================
CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE,
    status VARCHAR(50) NOT NULL DEFAULT 'upcoming',
    geofences JSONB NOT NULL DEFAULT '[]',
    event_config JSONB DEFAULT '{}',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES gebruikers(id),
    CONSTRAINT chk_event_status CHECK (status IN ('upcoming', 'active', 'completed', 'cancelled'))
);

CREATE INDEX IF NOT EXISTS idx_events_status ON events(status);
CREATE INDEX IF NOT EXISTS idx_events_active ON events(is_active);
CREATE INDEX IF NOT EXISTS idx_events_start_time ON events(start_time);
CREATE INDEX IF NOT EXISTS idx_events_geofences ON events USING gin(geofences);

-- Trigger voor updated_at (de functie is al aangemaakt in V19)
DROP TRIGGER IF EXISTS update_events_updated_at ON events;
CREATE TRIGGER update_events_updated_at
    BEFORE UPDATE ON events
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'events') THEN
        COMMENT ON TABLE events IS 'Events zoals De Koninklijke Loop met tracking configuratie';
        COMMENT ON COLUMN events.start_time IS 'Wanneer het event begint (UTC)';
        COMMENT ON COLUMN events.geofences IS 'Array van geofence objecten: [{type: "start|checkpoint|finish", lat: number, long: number, radius: number}]';
        COMMENT ON COLUMN events.event_config IS 'Extra configuratie zoals distance requirements, minStepsInterval, etc.';
        COMMENT ON COLUMN events.status IS 'upcoming, active, completed, cancelled';
    END IF;
END $$;

-- =====================================================
-- 2. EVENT PARTICIPANTS TABLE
-- =====================================================
CREATE TABLE IF NOT EXISTS event_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    participant_id UUID NOT NULL REFERENCES aanmeldingen(id) ON DELETE CASCADE,
    registered_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    check_in_time TIMESTAMP WITH TIME ZONE,
    start_time TIMESTAMP WITH TIME ZONE,
    finish_time TIMESTAMP WITH TIME ZONE,
    tracking_status VARCHAR(50) DEFAULT 'registered',
    last_location_update TIMESTAMP WITH TIME ZONE,
    total_distance DECIMAL(10,2) DEFAULT 0,
    current_steps INT DEFAULT 0,
    CONSTRAINT chk_tracking_status CHECK (tracking_status IN ('registered', 'checked_in', 'started', 'in_progress', 'finished', 'dnf'))
);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'event_participants_event_id_participant_id_key'
        AND conrelid = 'event_participants'::regclass
    ) THEN
        ALTER TABLE event_participants ADD CONSTRAINT event_participants_event_id_participant_id_key UNIQUE (event_id, participant_id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_event_participants_event ON event_participants(event_id);
CREATE INDEX IF NOT EXISTS idx_event_participants_participant ON event_participants(participant_id);
CREATE INDEX IF NOT EXISTS idx_event_participants_status ON event_participants(tracking_status);

DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'event_participants') THEN
        COMMENT ON TABLE event_participants IS 'Koppeling tussen events en participants met tracking status';
        COMMENT ON COLUMN event_participants.tracking_status IS 'registered, checked_in, started, in_progress, finished, dnf';
        COMMENT ON COLUMN event_participants.total_distance IS 'Totale afstand in kilometers';
    END IF;
END $$;

-- =====================================================
-- 3. SEED DATA - DEFAULT EVENT
-- =====================================================
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM events WHERE name = 'De Koninklijke Loop 2025') THEN
        INSERT INTO events (
        name,
        description,
        start_time,
        end_time,
        status,
        geofences,
        event_config,
        is_active
        ) VALUES (
        'De Koninklijke Loop 2025',
        'De jaarlijkse Koninklijke Loop hardloopevenement',
        '2025-05-16 09:00:00+00',
        '2025-05-16 16:00:00+00',
        'upcoming',
        '[
            {
            "type": "start",
            "lat": 52.0907,
            "long": 5.1214,
            "radius": 50,
            "name": "Start Locatie"
            },
            {
            "type": "checkpoint",
            "lat": 52.0950,
            "long": 5.1300,
            "radius": 30,
            "name": "Checkpoint 5KM"
            },
            {
            "type": "checkpoint",
            "lat": 52.1000,
            "long": 5.1400,
            "radius": 30,
            "name": "Checkpoint 10KM"
            },
            {
            "type": "finish",
            "lat": 52.0907,
            "long": 5.1214,
            "radius": 50,
            "name": "Finish Locatie"
            }
        ]'::jsonb,
        '{
            "minStepsInterval": 10,
            "requireGeofenceCheckin": true,
            "distanceThreshold": 100,
            "accuracyLevel": "balanced"
        }'::jsonb,
        true
        );
    END IF;
END $$;

-- =====================================================
-- 4. PERMISSIONS
-- =====================================================
INSERT INTO permissions (resource, action, description) VALUES
('events', 'read', 'Events kunnen bekijken'),
('events', 'write', 'Events kunnen aanmaken en bewerken'),
('events', 'delete', 'Events kunnen verwijderen'),
('event_tracking', 'read', 'Event tracking data kunnen bekijken'),
('event_tracking', 'write', 'Event tracking data kunnen bijwerken')
ON CONFLICT (resource, action) DO NOTHING;

-- Admin krijgt volledige toegang
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'admin'
AND p.resource IN ('events', 'event_tracking')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Staff krijgt read toegang
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'staff'
AND p.resource IN ('events', 'event_tracking')
AND p.action = 'read'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Deelnemers krijgen read toegang tot events en write voor hun eigen tracking
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'deelnemer'
AND (
(p.resource = 'events' AND p.action = 'read')
OR (p.resource = 'event_tracking' AND p.action IN ('read', 'write'))
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- =====================================================
-- 5. HELPER FUNCTIONS
-- =====================================================
CREATE OR REPLACE FUNCTION get_active_event()
RETURNS UUID AS $$ 
DECLARE
    v_event_id UUID;
BEGIN
    SELECT id INTO v_event_id
    FROM events
    WHERE is_active = true
    AND status IN ('upcoming', 'active')
    ORDER BY start_time ASC
    LIMIT 1;
    RETURN v_event_id;
END; 
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_active_event IS 'Retourneert ID van het eerstvolgende actieve event';

CREATE OR REPLACE FUNCTION is_in_geofence(
    p_lat DECIMAL,
    p_long DECIMAL,
    p_geofence JSONB
) RETURNS BOOLEAN AS $$ 
DECLARE
    v_fence_lat DECIMAL;
    v_fence_long DECIMAL;
    v_radius DECIMAL;
    v_distance DECIMAL;
BEGIN
    v_fence_lat := (p_geofence->>'lat')::DECIMAL;
    v_fence_long := (p_geofence->>'long')::DECIMAL;
    v_radius := (p_geofence->>'radius')::DECIMAL;
    
    -- Haversine formula (vereenvoudigd)
    v_distance := 111320 * SQRT(
        POW(v_fence_lat - p_lat, 2) +
        POW((v_fence_long - p_long) * COS(RADIANS(p_lat)), 2)
    );
    RETURN v_distance <= v_radius;
END; 
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION is_in_geofence IS 'Controleert of coördinaten binnen een geofence radius vallen';

-- =====================================================
-- 6. VIEWS
-- =====================================================
DROP VIEW IF EXISTS event_participants_view;
CREATE VIEW event_participants_view AS
SELECT
    ep.id as event_participant_id,
    ep.event_id,
    e.name as event_name,
    e.start_time as event_start_time,
    e.status as event_status,
    ep.participant_id,
    a.naam as participant_name,
    a.email as participant_email,
    a.afstand as participant_route,
    a.steps as participant_total_steps,
    ep.tracking_status,
    ep.registered_at,
    ep.check_in_time,
    ep.start_time as participant_start_time,
    ep.finish_time,
    ep.total_distance,
    ep.current_steps,
    ep.last_location_update,
    CASE
        WHEN ep.finish_time IS NOT NULL AND ep.start_time IS NOT NULL THEN
            EXTRACT(EPOCH FROM (ep.finish_time - ep.start_time))/60
        ELSE NULL
    END as duration_minutes
FROM event_participants ep
JOIN events e ON ep.event_id = e.id
JOIN aanmeldingen a ON ep.participant_id = a.id
ORDER BY ep.registered_at DESC;

COMMENT ON VIEW event_participants_view IS 'Overzicht van event participants met details';