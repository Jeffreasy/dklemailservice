-- GECONSOLIDEERDE V7 - RBAC SEEDING
-- Dit bestand combineert de logica van V1_21, V1_23, V1_24, V1_25, en V1_26.
-- Het seed alle rollen, permissies, en rol-permissie koppelingen in één keer.

-- Stap 1: Insert system roles (van V1_21)
INSERT INTO roles (name, description, is_system_role) VALUES
('admin', 'Volledige beheerder met toegang tot alle functies', true),
('staff', 'Ondersteunend personeel met beperkte beheerrechten', true),
('user', 'Standaard gebruiker', true),
('owner', 'Chat kanaal eigenaar', true),
('chat_admin', 'Chat kanaal beheerder', true),
('member', 'Chat kanaal lid', true),
('deelnemer', 'Evenement deelnemer', true),
('begeleider', 'Evenement begeleider', true),
('vrijwilliger', 'Evenement vrijwilliger', true)
ON CONFLICT (name) DO NOTHING;

-- Stap 2: Insert system permissions (gecombineerd van V1_21, V1_23, V1_24)
INSERT INTO permissions (resource, action, description, is_system_permission) VALUES
-- Van V1_21
('contact', 'read', 'Contactformulieren bekijken', true),
('contact', 'write', 'Contactformulieren bewerken (status, notities, antwoorden)', true),
('contact', 'delete', 'Contactformulieren verwijderen', true),
('aanmelding', 'read', 'Aanmeldingen bekijken', true),
('aanmelding', 'write', 'Aanmeldingen bewerken (status, notities, antwoorden)', true),
('aanmelding', 'delete', 'Aanmeldingen verwijderen', true),
('newsletter', 'read', 'Nieuwsbrieven bekijken', true),
('newsletter', 'write', 'Nieuwsbrieven aanmaken/bewerken', true),
('newsletter', 'send', 'Nieuwsbrieven verzenden', true),
('newsletter', 'delete', 'Nieuwsbrieven verwijderen', true),
('email', 'read', 'Inkomende emails bekijken', true),
('email', 'write', 'Emails bewerken (markeren als verwerkt)', true),
('email', 'delete', 'Emails verwijderen', true),
('email', 'fetch', 'Nieuwe emails ophalen', true),
('admin_email', 'send', 'Emails verzenden namens admin', true),
('user', 'read', 'Gebruikers bekijken', true),
('user', 'write', 'Gebruikers aanmaken/bewerken', true),
('user', 'delete', 'Gebruikers verwijderen', true),
('user', 'manage_roles', 'Gebruikersrollen beheren', true),
('chat', 'read', 'Chat kanalen en berichten bekijken', true),
('chat', 'write', 'Berichten verzenden', true),
('chat', 'manage_channel', 'Kanalen aanmaken/beheren', true),
('chat', 'moderate', 'Berichten modereren (bewerken/verwijderen)', true),
('notification', 'read', 'Notificaties bekijken', true),
('notification', 'write', 'Notificaties aanmaken', true),
('notification', 'delete', 'Notificaties verwijderen', true),
('system', 'admin', 'Volledige systeemtoegang', true),
-- Van V1_23
('admin', 'access', 'Volledige admin toegang', true),
('staff', 'access', 'Toegang tot staff functies', true),
-- Van V1_24
('photo', 'read', 'Foto''s bekijken', true),
('photo', 'write', 'Foto''s uploaden/bewerken', true),
('photo', 'delete', 'Foto''s verwijderen', true),
('album', 'read', 'Albums bekijken', true),
('album', 'write', 'Albums aanmaken/bewerken', true),
('album', 'delete', 'Albums verwijderen', true),
('partner', 'read', 'Partners bekijken', true),
('partner', 'write', 'Partners aanmaken/bewerken', true),
('partner', 'delete', 'Partners verwijderen', true),
('sponsor', 'read', 'Sponsors bekijken', true),
('sponsor', 'write', 'Sponsors aanmaken/bewerken', true),
('sponsor', 'delete', 'Sponsors verwijderen', true),
('video', 'read', 'Video''s bekijken', true),
('video', 'write', 'Video''s uploaden/bewerken', true),
('video', 'delete', 'Video''s verwijderen', true)
ON CONFLICT (resource, action) DO NOTHING;

-- Stap 3: Wijs permissies toe aan rollen

-- Admin role gets all permissions (gecombineerd van V1_21, V1_25)
-- Dit pakt nu ALLE permissies die in Stap 2 zijn aangemaakt.
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin' AND r.is_system_role = true
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Staff role gets specifieke permissies (gecombineerd van V1_21, V1_23, V1_26)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'staff' AND r.is_system_role = true
  AND (
    -- Van V1_21
    (p.resource IN ('user', 'contact', 'aanmelding', 'newsletter', 'email', 'chat', 'notification') AND p.action = 'read')
    -- Van V1_23
    OR (p.resource = 'staff' AND p.action = 'access')
    -- Van V1_26
    OR (p.resource IN ('photo', 'album', 'partner', 'sponsor', 'video') AND p.action = 'read')
  )
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Chat owner gets full chat permissions (van V1_21)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'owner' AND r.is_system_role = true
  AND p.resource = 'chat'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Chat admin gets most chat permissions (van V1_21)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'chat_admin' AND r.is_system_role = true
  AND p.resource = 'chat'
  AND p.action IN ('read', 'write', 'moderate')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Chat member gets basic chat permissions (van V1_21)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'member' AND r.is_system_role = true
  AND p.resource = 'chat'
  AND p.action IN ('read', 'write')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Regular user gets basic permissions (van V1_21)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'user' AND r.is_system_role = true
  AND p.resource = 'chat'
  AND p.action IN ('read', 'write')
ON CONFLICT (role_id, permission_id) DO NOTHING;