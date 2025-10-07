-- Migratie: V1_21__seed_rbac_data.sql
-- Beschrijving: Seed initial RBAC data based on current system roles and permissions
-- Versie: 1.21.0

-- Insert system roles (based on current hardcoded roles)
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

-- Insert system permissions (based on current middleware usage)
INSERT INTO permissions (resource, action, description, is_system_permission) VALUES
-- Contact management permissions
('contact', 'read', 'Contactformulieren bekijken', true),
('contact', 'write', 'Contactformulieren bewerken (status, notities, antwoorden)', true),
('contact', 'delete', 'Contactformulieren verwijderen', true),

-- Registration management permissions
('aanmelding', 'read', 'Aanmeldingen bekijken', true),
('aanmelding', 'write', 'Aanmeldingen bewerken (status, notities, antwoorden)', true),
('aanmelding', 'delete', 'Aanmeldingen verwijderen', true),

-- Newsletter permissions
('newsletter', 'read', 'Nieuwsbrieven bekijken', true),
('newsletter', 'write', 'Nieuwsbrieven aanmaken/bewerken', true),
('newsletter', 'send', 'Nieuwsbrieven verzenden', true),
('newsletter', 'delete', 'Nieuwsbrieven verwijderen', true),

-- Email management permissions
('email', 'read', 'Inkomende emails bekijken', true),
('email', 'write', 'Emails bewerken (markeren als verwerkt)', true),
('email', 'delete', 'Emails verwijderen', true),
('email', 'fetch', 'Nieuwe emails ophalen', true),

-- Admin email permissions
('admin_email', 'send', 'Emails verzenden namens admin', true),

-- User management permissions
('user', 'read', 'Gebruikers bekijken', true),
('user', 'write', 'Gebruikers aanmaken/bewerken', true),
('user', 'delete', 'Gebruikers verwijderen', true),
('user', 'manage_roles', 'Gebruikersrollen beheren', true),

-- Chat permissions
('chat', 'read', 'Chat kanalen en berichten bekijken', true),
('chat', 'write', 'Berichten verzenden', true),
('chat', 'manage_channel', 'Kanalen aanmaken/beheren', true),
('chat', 'moderate', 'Berichten modereren (bewerken/verwijderen)', true),

-- Notification permissions
('notification', 'read', 'Notificaties bekijken', true),
('notification', 'write', 'Notificaties aanmaken', true),
('notification', 'delete', 'Notificaties verwijderen', true),

-- System permissions
('system', 'admin', 'Volledige systeemtoegang', true)
ON CONFLICT (resource, action) DO NOTHING;

-- Assign permissions to roles based on current middleware usage

-- Admin role gets all permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin' AND r.is_system_role = true
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Staff role gets read permissions and some management
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'staff' AND r.is_system_role = true
  AND p.resource IN ('user', 'contact', 'aanmelding', 'newsletter', 'email', 'chat', 'notification')
  AND p.action = 'read'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Chat owner gets full chat permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'owner' AND r.is_system_role = true
  AND p.resource = 'chat'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Chat admin gets most chat permissions except channel creation
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'chat_admin' AND r.is_system_role = true
  AND p.resource = 'chat'
  AND p.action IN ('read', 'write', 'moderate')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Chat member gets basic chat permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'member' AND r.is_system_role = true
  AND p.resource = 'chat'
  AND p.action IN ('read', 'write')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Regular user gets basic permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'user' AND r.is_system_role = true
  AND p.resource = 'chat'
  AND p.action IN ('read', 'write')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Event roles (deelnemer, begeleider, vrijwilliger) get no special permissions
-- They are just labels for registration categorization

-- Registreer de migratie
INSERT INTO migraties (versie, naam, toegepast)
VALUES ('1.21.0', 'Seed initial RBAC data based on current system', CURRENT_TIMESTAMP)
ON CONFLICT (versie) DO NOTHING;