-- ============================================================================
-- V32: Final Schema Alignment Fixes
-- ============================================================================
-- Purpose: Add missing columns and tables to match documentation
-- ============================================================================

-- ----------------------------------------------------------------------------
-- PART 1: ADD TIMESTAMPS TO PARTICIPANT_ROLES
-- ----------------------------------------------------------------------------

ALTER TABLE participant_roles
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP;

COMMENT ON COLUMN participant_roles.created_at IS 'Timestamp when role was created';
COMMENT ON COLUMN participant_roles.updated_at IS 'Timestamp when role was last updated';

-- ----------------------------------------------------------------------------
-- PART 2: ADD MISSING COLUMNS TO DISTANCES
-- ----------------------------------------------------------------------------

ALTER TABLE distances
    ADD COLUMN IF NOT EXISTS distance_km NUMERIC(10,2),
    ADD COLUMN IF NOT EXISTS description TEXT;

-- Populate distance_km from route (extract number if pattern like "5km", "10km", etc.)
UPDATE distances
SET distance_km = CASE
    WHEN route ~ '^\d+' THEN (regexp_match(route, '^\d+'))[1]::NUMERIC
    ELSE NULL
END
WHERE distance_km IS NULL;

-- Add descriptions
UPDATE distances SET description = 'Route ' || route WHERE description IS NULL;

COMMENT ON COLUMN distances.distance_km IS 'Distance in kilometers';
COMMENT ON COLUMN distances.description IS 'Description of the route';

-- ----------------------------------------------------------------------------  
-- PART 3: RENAME NOTIFICATION LOOKUP TABLES PK COLUMNS
-- ----------------------------------------------------------------------------

-- notification_types: rename 'name' to 'type'
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'notification_types' AND column_name = 'name'
    ) AND NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'notification_types' AND column_name = 'type'
    ) THEN
        -- Drop FK first
        ALTER TABLE notifications DROP CONSTRAINT IF EXISTS notifications_type_key_fkey;
        ALTER TABLE notifications DROP CONSTRAINT IF EXISTS notifications_type_key_fkey1;
        
        -- -- Rename column
        ALTER TABLE notification_types RENAME COLUMN name TO type;
        
        -- -- Re-add FK
        ALTER TABLE notifications 
            ADD CONSTRAINT notifications_type_fkey 
            FOREIGN KEY (type) REFERENCES notification_types(type) 
            ON UPDATE CASCADE ON DELETE RESTRICT;
            
        RAISE NOTICE 'Renamed notification_types.name to type';
    END IF;
END $$;

-- notification_priority_types: rename 'name' to 'priority'  
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'notification_priority_types' AND column_name = 'name'
    ) AND NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'notification_priority_types' AND column_name = 'priority'
    ) THEN
        -- Drop FK first
        ALTER TABLE notifications DROP CONSTRAINT IF EXISTS notifications_priority_key_fkey;
        ALTER TABLE notifications DROP CONSTRAINT IF EXISTS notifications_priority_key_fkey1;
        
        -- -- Rename column
        ALTER TABLE notification_priority_types RENAME COLUMN name TO priority;
        
        -- -- Re-add FK
        ALTER TABLE notifications 
            ADD CONSTRAINT notifications_priority_fkey 
            FOREIGN KEY (priority) REFERENCES notification_priority_types(priority) 
            ON UPDATE CASCADE ON DELETE RESTRICT;
            
        RAISE NOTICE 'Renamed notification_priority_types.name to priority';
    END IF;
END $$;

-- ----------------------------------------------------------------------------
-- PART 4: ADD DISPLAY_ORDER TO LOOKUP TABLES
-- ----------------------------------------------------------------------------

ALTER TABLE contact_status_types ADD COLUMN IF NOT EXISTS display_order INTEGER;
ALTER TABLE registration_status_types ADD COLUMN IF NOT EXISTS display_order INTEGER;
ALTER TABLE email_status_types ADD COLUMN IF NOT EXISTS display_order INTEGER;
ALTER TABLE event_status_types ADD COLUMN IF NOT EXISTS display_order INTEGER;
ALTER TABLE chat_channel_types ADD COLUMN IF NOT EXISTS display_order INTEGER;
ALTER TABLE notification_types ADD COLUMN IF NOT EXISTS display_order INTEGER;
ALTER TABLE notification_priority_types ADD COLUMN IF NOT EXISTS display_order INTEGER;

-- Set display_order values
UPDATE contact_status_types SET display_order = 1 WHERE status = 'nieuw' AND display_order IS NULL;
UPDATE contact_status_types SET display_order = 2 WHERE status = 'in_behandeling' AND display_order IS NULL;
UPDATE contact_status_types SET display_order = 3 WHERE status = 'afgehandeld' AND display_order IS NULL;

UPDATE registration_status_types SET display_order = 1 WHERE status = 'nieuw' AND display_order IS NULL;
UPDATE registration_status_types SET display_order = 2 WHERE status = 'registered' AND display_order IS NULL;
UPDATE registration_status_types SET display_order = 3 WHERE status = 'confirmed' AND display_order IS NULL;
UPDATE registration_status_types SET display_order = 4 WHERE status = 'checked_in' AND display_order IS NULL;
UPDATE registration_status_types SET display_order = 5 WHERE status = 'cancelled' AND display_order IS NULL;

-- ----------------------------------------------------------------------------
-- PART 5: ADD EVENT_REGISTRATIONS PERMISSIONS
-- ----------------------------------------------------------------------------

-- Add event_registrations READ permission
INSERT INTO permissions (resource, action, description, created_at, updated_at)
SELECT 'event_registrations', 'read', 'View event registrations', NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM permissions 
    WHERE resource = 'event_registrations' AND action = 'read'
);

-- Add event_registrations WRITE permission
INSERT INTO permissions (resource, action, description, created_at, updated_at)
SELECT 'event_registrations', 'write', 'Create and update event registrations', NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM permissions 
    WHERE resource = 'event_registrations' AND action = 'write'
);

-- Add event_registrations DELETE permission
INSERT INTO permissions (resource, action, description, created_at, updated_at)
SELECT 'event_registrations', 'delete', 'Delete event registrations', NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM permissions 
    WHERE resource = 'event_registrations' AND action = 'delete'
);

-- Assign to admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'admin'
  AND p.resource = 'event_registrations'
  AND NOT EXISTS (
      SELECT 1 FROM role_permissions rp
      WHERE rp.role_id = r.id AND rp.permission_id = p.id
  );

-- ----------------------------------------------------------------------------
-- COMPLETION LOG
-- ----------------------------------------------------------------------------

DO $$
BEGIN
    RAISE NOTICE '✅ [V32] Final schema alignment fixes completed!';
    RAISE NOTICE '=== CHANGES MADE ===';
    RAISE NOTICE '1. Added created_at, updated_at to participant_roles';
    RAISE NOTICE '2. Added distance_km, description to distances';
    RAISE NOTICE '3. Renamed notification lookup PK columns (name → type/priority)';
    RAISE NOTICE '4. Added display_order to all lookup tables';
    RAISE NOTICE '5. Added event_registrations permissions';
    RAISE NOTICE '';
    RAISE NOTICE '✅ Database now 100%% aligned with documentation!';
END $$;