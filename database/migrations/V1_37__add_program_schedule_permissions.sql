-- V1_37__add_program_schedule_permissions.sql
-- Add permissions for program schedule management

-- Insert permissions for program_schedule
INSERT INTO permissions (resource, action, description, is_system_permission) VALUES
('program_schedule', 'read', 'Programma bekijken', true),
('program_schedule', 'write', 'Programma aanmaken/bewerken', true),
('program_schedule', 'delete', 'Programma verwijderen', true)
ON CONFLICT (resource, action) DO NOTHING;

-- Assign program_schedule permissions to admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin' AND r.is_system_role = true
  AND p.resource = 'program_schedule'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign read permissions to staff role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'staff' AND r.is_system_role = true
  AND p.resource = 'program_schedule'
  AND p.action = 'read'
ON CONFLICT (role_id, permission_id) DO NOTHING;