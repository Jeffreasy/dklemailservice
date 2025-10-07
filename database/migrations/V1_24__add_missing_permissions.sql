-- Migratie: V1_24__add_missing_permissions.sql
-- Beschrijving: Add missing permissions for Photos, Albums, Partners, Sponsors, Videos
-- Versie: 1.24.0

-- Voeg ontbrekende permissions toe voor Photos, Albums, Partners, Sponsors, Videos
INSERT INTO permissions (resource, action, description, is_system_permission) VALUES
-- Photos
('photo', 'read', 'Foto''s bekijken', true),
('photo', 'write', 'Foto''s uploaden/bewerken', true),
('photo', 'delete', 'Foto''s verwijderen', true),

-- Albums
('album', 'read', 'Albums bekijken', true),
('album', 'write', 'Albums aanmaken/bewerken', true),
('album', 'delete', 'Albums verwijderen', true),

-- Partners
('partner', 'read', 'Partners bekijken', true),
('partner', 'write', 'Partners aanmaken/bewerken', true),
('partner', 'delete', 'Partners verwijderen', true),

-- Sponsors
('sponsor', 'read', 'Sponsors bekijken', true),
('sponsor', 'write', 'Sponsors aanmaken/bewerken', true),
('sponsor', 'delete', 'Sponsors verwijderen', true),

-- Videos
('video', 'read', 'Video''s bekijken', true),
('video', 'write', 'Video''s uploaden/bewerken', true),
('video', 'delete', 'Video''s verwijderen', true)
ON CONFLICT (resource, action) DO NOTHING;

-- Registreer de migratie
INSERT INTO migraties (versie, naam, toegepast)
VALUES ('1.24.0', 'Add missing permissions for Photos, Albums, Partners, Sponsors, Videos', CURRENT_TIMESTAMP)
ON CONFLICT (versie) DO NOTHING;