-- V1_33__add_partner_radio_permissions.sql
-- Add permissions for partners and radio recordings management

-- Insert permissions for partners
INSERT INTO permissions (resource, action, description, is_system_permission) VALUES
('partner', 'read', 'Partners bekijken', true),
('partner', 'write', 'Partners aanmaken/bewerken', true),
('partner', 'delete', 'Partners verwijderen', true),

-- Insert permissions for radio recordings
('radio_recording', 'read', 'Radio opnames bekijken', true),
('radio_recording', 'write', 'Radio opnames aanmaken/bewerken', true),
('radio_recording', 'delete', 'Radio opnames verwijderen', true)
ON CONFLICT (resource, action) DO NOTHING;

-- Assign partner and radio_recording permissions to admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin' AND r.is_system_role = true
  AND p.resource IN ('partner', 'radio_recording')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign read permissions to staff role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'staff' AND r.is_system_role = true
  AND p.resource IN ('partner', 'radio_recording')
  AND p.action = 'read'
ON CONFLICT (role_id, permission_id) DO NOTHING;