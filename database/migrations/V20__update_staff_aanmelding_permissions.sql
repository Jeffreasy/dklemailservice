-- GECONSOLIDEERDE V20 - UPDATE STAFF AANMELDING PERMISSIES
-- Logica van V1_50.
-- De permissies 'aanmelding' (read, write, delete) en de 'admin' toewijzing
-- zijn al uitgevoerd in V7.
-- Dit script voegt de 'write' permissie toe aan de 'staff' rol.

-- Wijs read en write permissies toe aan staff role (geen delete)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'staff' AND r.is_system_role = true
  AND p.resource = 'aanmelding'
  AND p.action IN ('read', 'write')
ON CONFLICT (role_id, permission_id) DO NOTHING;