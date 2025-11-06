```sql
-- Gecombineerde Migratie: Geïntegreerde RBAC seeding, migratie, permissies en nieuwe tabellen
-- Beschrijving: Combinatie van migraties V1_21 t/m V1_30 voor een enkelvoudig script met RBAC setup, seeds en toevoegingen
-- Versie: 1.2.0 (geconsolideerd)

-- Zorg ervoor dat de pgcrypto extensie beschikbaar is voor gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Maak refresh_tokens tabel aan (uit V1_28)
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES gebruikers(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMP,
    is_revoked BOOLEAN DEFAULT FALSE
);

-- Indices voor refresh_tokens (uit V1_28)
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_is_revoked ON refresh_tokens(is_revoked);

-- Commentaren voor refresh_tokens (uit V1_28)
COMMENT ON TABLE refresh_tokens IS 'Stores refresh tokens for JWT authentication with 7-day expiry';
COMMENT ON COLUMN refresh_tokens.token IS 'Base64 encoded random token (32 bytes)';
COMMENT ON COLUMN refresh_tokens.expires_at IS 'Token expiration timestamp (7 days from creation)';
COMMENT ON COLUMN refresh_tokens.is_revoked IS 'Whether the token has been revoked (for token rotation)';

-- Alter chat_messages voor thumbnail_url (uit V1_29)
ALTER TABLE IF EXISTS chat_messages
ADD COLUMN IF NOT EXISTS thumbnail_url TEXT;

-- Maak uploaded_images tabel aan (uit V1_30)
CREATE TABLE IF NOT EXISTS uploaded_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES gebruikers(id) ON DELETE CASCADE,
    public_id TEXT NOT NULL UNIQUE,
    url TEXT NOT NULL,
    secure_url TEXT NOT NULL,
    filename TEXT NOT NULL,
    size BIGINT NOT NULL,
    mime_type TEXT NOT NULL,
    width INTEGER,
    height INTEGER,
    folder TEXT NOT NULL,
    thumbnail_url TEXT,
    deleted_at TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indices voor uploaded_images (uit V1_30)
CREATE INDEX IF NOT EXISTS idx_uploaded_images_user_id ON uploaded_images(user_id);
CREATE INDEX IF NOT EXISTS idx_uploaded_images_public_id ON uploaded_images(public_id);
CREATE INDEX IF NOT EXISTS idx_uploaded_images_folder ON uploaded_images(folder);
CREATE INDEX IF NOT EXISTS idx_uploaded_images_deleted_at ON uploaded_images(deleted_at);
CREATE INDEX IF NOT EXISTS idx_uploaded_images_created_at ON uploaded_images(created_at DESC);

-- RBAC seeding: Inserts voor roles (uit V1_21)
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

-- RBAC seeding: Inserts voor permissions (gecombineerd uit V1_21, V1_23, V1_24)
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
    ('system', 'admin', 'Volledige systeemtoegang', true),
    -- Admin and staff access
    ('admin', 'access', 'Volledige admin toegang', true),
    ('staff', 'access', 'Toegang tot staff functies', true),
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

-- Assign permissions to roles (gecombineerd uit V1_21, V1_23, V1_25, V1_26)
-- Admin gets all permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin' AND r.is_system_role = true
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Staff gets read permissions and staff access
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'staff' AND r.is_system_role = true
  AND ((p.resource IN ('user', 'contact', 'aanmelding', 'newsletter', 'email', 'chat', 'notification', 'photo', 'album', 'partner', 'sponsor', 'video')
        AND p.action = 'read')
       OR (p.resource = 'staff' AND p.action = 'access'))
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Chat owner gets full chat permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'owner' AND r.is_system_role = true
  AND p.resource = 'chat'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Chat admin gets most chat permissions
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

-- Regular user gets basic chat permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'user' AND r.is_system_role = true
  AND p.resource = 'chat'
  AND p.action IN ('read', 'write')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign admin role to admin user (uit V1_22_assign_admin_role)
INSERT INTO user_roles (user_id, role_id, assigned_at, is_active)
SELECT u.id, r.id, CURRENT_TIMESTAMP, true
FROM gebruikers u
CROSS JOIN roles r
WHERE u.email = 'admin@dekoninklijkeloop.nl'
  AND r.name = 'admin'
  AND r.is_system_role = true
ON CONFLICT (user_id, role_id) DO NOTHING;

-- Assign staff role to jeffrey (uit V1_27)
INSERT INTO user_roles (user_id, role_id, assigned_at, is_active)
SELECT u.id, r.id, CURRENT_TIMESTAMP, true
FROM gebruikers u
CROSS JOIN roles r
WHERE u.email = 'jeffrey@dekoninklijkeloop.nl'
  AND r.name = 'staff'
  AND r.is_system_role = true
ON CONFLICT (user_id, role_id) DO NOTHING;

-- Migrate legacy roles to RBAC (uit V1_22_migrate_legacy_roles_to_rbac)
-- Voeg user_roles toe voor bestaande gebruikers met legacy roles
INSERT INTO user_roles (user_id, role_id, is_active, assigned_at)
SELECT 
    g.id as user_id,
    r.id as role_id,
    true as is_active,
    COALESCE(g.created_at, NOW()) as assigned_at
FROM gebruikers g
JOIN roles r ON LOWER(r.name) = LOWER(g.rol)
WHERE g.rol IS NOT NULL 
  AND g.rol != ''
ON CONFLICT (user_id, role_id) DO NOTHING;

-- Voeg standaard 'user' role toe voor gebruikers zonder specifieke rol
INSERT INTO user_roles (user_id, role_id, is_active, assigned_at)
SELECT 
    g.id as user_id,
    r.id as role_id,
    true as is_active,
    COALESCE(g.created_at, NOW()) as assigned_at
FROM gebruikers g
CROSS JOIN roles r
WHERE r.name = 'user'
  AND r.is_system_role = true
  AND (g.rol IS NULL OR g.rol = '' OR g.rol = 'gebruiker')
ON CONFLICT (user_id, role_id) DO NOTHING;

-- Create view voor migratie status (uit V1_22_migrate_legacy_roles_to_rbac)
CREATE OR REPLACE VIEW v_user_role_migration_status AS
SELECT 
    g.id as user_id,
    g.email,
    g.naam,
    g.rol as legacy_role,
    COALESCE(
        STRING_AGG(r.name, ', ' ORDER BY r.name), 
        'GEEN RBAC ROL'
    ) as rbac_roles,
    COUNT(ur.id) as rbac_role_count,
    CASE 
        WHEN COUNT(ur.id) = 0 THEN 'MISSING RBAC'
        WHEN g.rol IS NULL OR g.rol = '' THEN 'NO LEGACY'
        WHEN EXISTS (
            SELECT 1 FROM user_roles ur2
            JOIN roles r2 ON ur2.role_id = r2.id
            WHERE ur2.user_id = g.id 
            AND LOWER(r2.name) = LOWER(g.rol)
            AND ur2.is_active = true
        ) THEN 'MIGRATED'
        ELSE 'MISMATCH'
    END as migration_status
FROM gebruikers g
LEFT JOIN user_roles ur ON g.id = ur.user_id AND ur.is_active = true
LEFT JOIN roles r ON ur.role_id = r.id
GROUP BY g.id, g.email, g.naam, g.rol
ORDER BY 
    CASE 
        WHEN COUNT(ur.id) = 0 THEN 1
        ELSE 2
    END,
    g.email;

-- Log migratie resultaten (uit V1_22_migrate_legacy_roles_to_rbac)
DO $$
DECLARE
    migrated_count INTEGER;
    total_users INTEGER;
    users_without_rbac INTEGER;
    status_record RECORD;
    problem_record RECORD;
    problem_count INTEGER := 0;
BEGIN
    -- Tel totaal aantal gebruikers
    SELECT COUNT(*) INTO total_users FROM gebruikers;
    
    -- Tel gebruikers met RBAC roles
    SELECT COUNT(DISTINCT user_id) INTO migrated_count FROM user_roles WHERE is_active = true;
    
    -- Tel gebruikers zonder RBAC roles
    users_without_rbac := total_users - migrated_count;
    
    -- Log resultaten
    RAISE NOTICE 'Legacy to RBAC Migration Results:';
    RAISE NOTICE '  Total users: %', total_users;
    RAISE NOTICE '  Users with RBAC roles: %', migrated_count;
    RAISE NOTICE '  Users without RBAC roles: %', users_without_rbac;
    
    -- Waarschuwing als er gebruikers zonder RBAC rollen zijn
    IF users_without_rbac > 0 THEN
        RAISE WARNING 'Er zijn % gebruikers zonder RBAC rollen! Controleer de rol mapping.', users_without_rbac;
    END IF;

    -- Toon status overzicht
    RAISE NOTICE '=== Migration Status per Category ===';
    FOR status_record IN 
        SELECT 
            migration_status,
            COUNT(*) as count
        FROM v_user_role_migration_status
        GROUP BY migration_status
        ORDER BY 
            CASE migration_status
                WHEN 'MIGRATED' THEN 1
                WHEN 'NO LEGACY' THEN 2
                WHEN 'MISMATCH' THEN 3
                WHEN 'MISSING RBAC' THEN 4
            END
    LOOP
        RAISE NOTICE '  %: % users', status_record.migration_status, status_record.count;
    END LOOP;

    -- Toon problematische gevallen
    RAISE NOTICE '=== Problematische Gebruikers (max 10) ===';
    FOR problem_record IN 
        SELECT user_id, email, naam, legacy_role, rbac_roles, migration_status
        FROM v_user_role_migration_status
        WHERE migration_status IN ('MISMATCH', 'MISSING RBAC')
        LIMIT 10
    LOOP
        problem_count := problem_count + 1;
        RAISE NOTICE '  % - % (%) | Legacy: % | RBAC: % | Status: %', 
            problem_count,
            problem_record.email,
            problem_record.naam,
            COALESCE(problem_record.legacy_role, 'NULL'),
            problem_record.rbac_roles,
            problem_record.migration_status;
    END LOOP;
    
    IF problem_count = 0 THEN
        RAISE NOTICE '✓ Geen problematische gebruikers gevonden!';
    END IF;

    -- Toon instructies voor vervolgstappen
    RAISE NOTICE '';
    RAISE NOTICE '=== NEXT STEPS ===';
    RAISE NOTICE '1. Review migration status: SELECT * FROM v_user_role_migration_status;';
    RAISE NOTICE '2. Fix any mismatches or missing RBAC roles manually';
    RAISE NOTICE '3. Test RBAC permissions thoroughly';
    RAISE NOTICE '4. Update JWT generation to use RBAC roles (code change required)';
    RAISE NOTICE '5. After thorough testing, consider deprecating gebruikers.rol field';
    RAISE NOTICE '';
    RAISE NOTICE '⚠️  BELANGRIJK: gebruikers.rol field is NIET verwijderd voor backward compatibility';
    RAISE NOTICE '⚠️  Beide systemen werken nu naast elkaar';
END $$;

-- Registreer de gecombineerde migratie
INSERT INTO migraties (versie, naam, toegepast) 
VALUES ('1.2.0', 'Geïntegreerde RBAC seeding, migratie, permissies en nieuwe tabellen (refresh_tokens, thumbnails, uploaded_images)', CURRENT_TIMESTAMP)
ON CONFLICT (versie) DO NOTHING;
```