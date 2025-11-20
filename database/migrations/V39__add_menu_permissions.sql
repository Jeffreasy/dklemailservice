-- V39: Menu permissies toevoegen voor frontend zichtbaarheid controle
-- Deze permissies bepalen welke menu-items zichtbaar zijn voor verschillende gebruikersrollen

INSERT INTO permissions (resource, action, description, is_system_permission) VALUES
-- Dashboard menu items
('menu', 'dashboard', 'Dashboard menu-item zichtbaar', true),
('menu', 'analytics', 'Analytics menu-item zichtbaar', true),

-- Content management menu items
('menu', 'media', 'Media beheer menu-item zichtbaar', true),
('menu', 'albums', 'Albums menu-item zichtbaar', true),
('menu', 'videos', 'Videos menu-item zichtbaar', true),
('menu', 'photos', 'Foto''s menu-item zichtbaar', true),

-- Communication menu items
('menu', 'emails', 'Emails menu-item zichtbaar', true),
('menu', 'contacts', 'Contact formulieren menu-item zichtbaar', true),
('menu', 'newsletters', 'Nieuwsbrieven menu-item zichtbaar', true),

-- User management menu items
('menu', 'users', 'Gebruikers menu-item zichtbaar', true),
('menu', 'participants', 'Deelnemers menu-item zichtbaar', true),

-- Event management menu items
('menu', 'events', 'Evenementen menu-item zichtbaar', true),
('menu', 'registrations', 'Aanmeldingen menu-item zichtbaar', true),

-- Administrative menu items
('menu', 'admin', 'Admin paneel menu-item zichtbaar', true),
('menu', 'settings', 'Instellingen menu-item zichtbaar', true),
('menu', 'reports', 'Rapporten menu-item zichtbaar', true),

-- Community features menu items
('menu', 'chat', 'Chat menu-item zichtbaar', true),
('menu', 'leaderboard', 'Leaderboard menu-item zichtbaar', true),
('menu', 'achievements', 'Achievements menu-item zichtbaar', true),

-- Special menu items
('menu', 'under_construction', 'Under construction pagina zichtbaar', true),
('menu', 'program_schedule', 'Programma schema menu-item zichtbaar', true),
('menu', 'radio_recording', 'Radio opnames menu-item zichtbaar', true)
ON CONFLICT (resource, action) DO NOTHING;

-- Wijs alle menu permissies toe aan de 'admin' rol
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin' AND r.is_system_role = true
  AND p.resource = 'menu'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Wijs alle menu permissies toe aan de 'staff' rol (behalve admin)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'staff' AND r.is_system_role = true
  AND p.resource = 'menu'
  AND p.action != 'admin'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Log de voltooiing
DO $$
BEGIN
    RAISE NOTICE '[V39] Menu permissies toegevoegd voor frontend zichtbaarheid controle';
END $$;