-- GECONSOLIDEERDE V14 - OVERIGE CMS PERMISSIES
-- Dit bestand combineert de logica van V1_33 (radio), V1_37, V1_38, V1_39, en V1_40.
-- De permissies voor partner, album, video, en sponsor waren al toegevoegd in V7.

-- Stap 1: Voeg de nieuwe permissies toe
INSERT INTO permissions (resource, action, description, is_system_permission) VALUES
-- Van V1_33
('radio_recording', 'read', 'Radio opnames bekijken', true),
('radio_recording', 'write', 'Radio opnames aanmaken/bewerken', true),
('radio_recording', 'delete', 'Radio opnames verwijderen', true),
-- Van V1_37
('program_schedule', 'read', 'Programma bekijken', true),
('program_schedule', 'write', 'Programma aanmaken/bewerken', true),
('program_schedule', 'delete', 'Programma verwijderen', true),
-- Van V1_38
('social_embed', 'read', 'Social embeds bekijken', true),
('social_embed', 'write', 'Social embeds aanmaken/bewerken', true),
('social_embed', 'delete', 'Social embeds verwijderen', true),
-- Van V1_39
('social_link', 'read', 'Social links bekijken', true),
('social_link', 'write', 'Social links aanmaken/bewerken', true),
('social_link', 'delete', 'Social links verwijderen', true),
-- Van V1_40
('under_construction', 'read', 'Under construction bekijken', true),
('under_construction', 'write', 'Under construction aanmaken/bewerken', true),
('under_construction', 'delete', 'Under construction verwijderen', true)
ON CONFLICT (resource, action) DO NOTHING;

-- Stap 2: Wijs alle nieuwe permissies toe aan 'admin'
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin' AND r.is_system_role = true
  AND p.resource IN (
    'radio_recording', 
    'program_schedule', 
    'social_embed', 
    'social_link', 
    'under_construction'
  )
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Stap 3: Wijs 'read' permissies toe aan 'staff'
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'staff' AND r.is_system_role = true
  AND p.resource IN (
    'radio_recording', 
    'program_schedule', 
    'social_embed', 
    'social_link', 
    'under_construction'
  )
  AND p.action = 'read'
ON CONFLICT (role_id, permission_id) DO NOTHING;