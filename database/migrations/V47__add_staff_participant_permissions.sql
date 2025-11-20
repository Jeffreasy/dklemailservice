-- V47: Add participant permissions to staff role
-- This migration gives staff members access to view all participant registrations

-- Add participant:view_all_registrations permission to staff role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'staff' AND r.is_system_role = true
  AND p.resource = 'participant' AND p.action = 'view_all_registrations'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Also add some other useful participant permissions for staff
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'staff' AND r.is_system_role = true
  AND p.resource = 'participant' AND p.action IN ('read', 'write')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Update admin password to "DKL2026!"
UPDATE gebruikers
SET wachtwoord_hash = '$2a$10$lHPwGgMrPQimDJVbOQakduAWYAhTM.ohG8V3OSaSZFbKLlIulhKT.'  -- bcrypt hash for "DKL2026!"
WHERE email = 'admin@dekoninklijkeloop.nl';

-- Log the changes
DO $$
DECLARE
    staff_role_id UUID;
    admin_user_id UUID;
BEGIN
    -- Get staff role ID
    SELECT id INTO staff_role_id FROM roles WHERE name = 'staff' AND is_system_role = true;

    -- Get admin user ID
    SELECT id INTO admin_user_id FROM gebruikers WHERE email = 'admin@dekoninklijkeloop.nl';

    RAISE NOTICE 'V47: Added participant permissions to staff role (ID: %)', staff_role_id;
    RAISE NOTICE 'V47: Updated password for admin user (ID: %)', admin_user_id;
END $$;