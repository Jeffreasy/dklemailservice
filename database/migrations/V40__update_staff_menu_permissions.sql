-- V40: Update staff role menu permissions - give access to all menus except admin
-- This migration updates the staff role to have access to all menu items except the admin panel

-- Remove existing staff menu permissions
DELETE FROM role_permissions
WHERE role_id IN (
    SELECT id FROM roles WHERE name = 'staff' AND is_system_role = true
)
AND permission_id IN (
    SELECT id FROM permissions WHERE resource = 'menu'
);

-- Add all menu permissions to staff role (except admin)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'staff' AND r.is_system_role = true
  AND p.resource = 'menu'
  AND p.action != 'admin'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Log the completion
DO $$
BEGIN
    RAISE NOTICE '[V40] Staff role updated with all menu permissions except admin';
END $$;