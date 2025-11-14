-- V30: DUAAL REGISTRATIESYSTEEM MET VOLLEDIGE RBAC INTEGRATIE
-- Datum: 2025-11-10
-- KRITIEKE INTEGRATIE: Koppelt V30 participant systeem volledig aan RBAC

-- ==============================================================================
-- STAP 1: PARTICIPANTS TABEL UITBREIDINGEN (V30 Core)
-- ==============================================================================

-- Voeg nieuwe kolommen toe voor duaal systeem
ALTER TABLE participants ADD COLUMN IF NOT EXISTS account_type TEXT DEFAULT 'temporary';
ALTER TABLE participants ADD COLUMN IF NOT EXISTS registration_year INTEGER;
ALTER TABLE participants ADD COLUMN IF NOT EXISTS wachtwoord_hash TEXT;
ALTER TABLE participants ADD COLUMN IF NOT EXISTS has_app_access BOOLEAN DEFAULT FALSE;
ALTER TABLE participants ADD COLUMN IF NOT EXISTS upgraded_to_gebruiker_id UUID;
ALTER TABLE participants ADD COLUMN IF NOT EXISTS upgraded_at TIMESTAMPTZ;

-- Update bestaande NULL waarden (maak compatible met nieuwe systeem)
UPDATE participants SET account_type = 'temporary' WHERE account_type IS NULL;
UPDATE participants SET has_app_access = FALSE WHERE has_app_access IS NULL;
UPDATE participants SET registration_year = 2026 WHERE registration_year IS NULL AND account_type = 'temporary';

-- Maak kolommen NOT NULL waar nodig
ALTER TABLE participants ALTER COLUMN account_type SET NOT NULL;
ALTER TABLE participants ALTER COLUMN has_app_access SET NOT NULL;

-- Voeg constraints toe (gebruik DO block voor IF NOT EXISTS logica)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'check_account_type' AND conrelid = 'participants'::regclass
    ) THEN
        ALTER TABLE participants ADD CONSTRAINT check_account_type
            CHECK (account_type IN ('full', 'temporary'));
    END IF;
END $$;

-- Foreign key voor upgrade tracking
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_participants_upgraded_to_gebruiker' AND conrelid = 'participants'::regclass
    ) THEN
        ALTER TABLE participants ADD CONSTRAINT fk_participants_upgraded_to_gebruiker
            FOREIGN KEY (upgraded_to_gebruiker_id) REFERENCES gebruikers(id) ON DELETE SET NULL;
    END IF;
END $$;

-- Performance indexes
CREATE INDEX IF NOT EXISTS idx_participants_account_type ON participants(account_type);
CREATE INDEX IF NOT EXISTS idx_participants_registration_year ON participants(registration_year);
CREATE INDEX IF NOT EXISTS idx_participants_has_app_access ON participants(has_app_access);
CREATE INDEX IF NOT EXISTS idx_participants_upgraded_to_gebruiker_id ON participants(upgraded_to_gebruiker_id);

-- Unieke constraint: 1 temporary registratie per email per jaar
-- Eerst: Verwijder duplicates (behoud oudste registratie per email/jaar)
DELETE FROM participants p1
WHERE account_type = 'temporary'
  AND EXISTS (
    SELECT 1 FROM participants p2
    WHERE p2.account_type = 'temporary'
      AND p2.email = p1.email
      AND p2.registration_year = p1.registration_year
      AND p2.created_at < p1.created_at
  );

-- Dan: Maak unique index
CREATE UNIQUE INDEX IF NOT EXISTS idx_participants_temp_year_unique
ON participants(email, registration_year)
WHERE account_type = 'temporary';

-- ==============================================================================
-- STAP 2: PARTICIPANT UPGRADES AUDIT TABEL
-- ==============================================================================

CREATE TABLE IF NOT EXISTS participant_upgrades (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    participant_id UUID NOT NULL REFERENCES participants(id) ON DELETE CASCADE,
    gebruiker_id UUID NOT NULL REFERENCES gebruikers(id) ON DELETE CASCADE,
    upgraded_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    upgraded_by UUID REFERENCES gebruikers(id), -- Admin die upgrade triggerde (NULL = self-upgrade)
    notes TEXT,
    
    -- Metadata
    old_account_type TEXT DEFAULT 'temporary',
    new_account_type TEXT DEFAULT 'full',
    
    CONSTRAINT check_upgrade_types CHECK (
        old_account_type = 'temporary' AND new_account_type = 'full'
    )
);

CREATE INDEX IF NOT EXISTS idx_participant_upgrades_participant_id ON participant_upgrades(participant_id);
CREATE INDEX IF NOT EXISTS idx_participant_upgrades_gebruiker_id ON participant_upgrades(gebruiker_id);
CREATE INDEX IF NOT EXISTS idx_participant_upgrades_upgraded_at ON participant_upgrades(upgraded_at);

COMMENT ON TABLE participant_upgrades IS 'V30: Audit trail voor account upgrades van temporary naar full';

-- ==============================================================================
-- STAP 3: RBAC PERMISSIONS VOOR PARTICIPANTS (NIEUW)
-- ==============================================================================

-- Voeg participant-specifieke permissions toe aan het systeem
INSERT INTO permissions (resource, action, description, is_system_permission) VALUES
-- Basis participant permissions
('participant', 'read_own', 'Eigen participant gegevens bekijken', true),
('participant', 'write_own', 'Eigen participant gegevens wijzigen', true),
('participant', 'register_event', 'Registreren voor events', true),
('participant', 'view_registrations', 'Eigen event registraties bekijken', true),
('participant', 'cancel_registration', 'Eigen registratie annuleren', true),

-- App toegang permissions
('app', 'access', 'Toegang tot DKL Step App', true),
('app', 'login', 'Inloggen in DKL Step App', true),

-- Steps & Gamification permissions
('steps', 'track', 'Stappen bijhouden en synchroniseren', true),
('steps', 'view_own', 'Eigen stappen geschiedenis bekijken', true),
('achievements', 'view', 'Achievements bekijken', true),
('achievements', 'earn', 'Achievements verdienen', true),
('badges', 'view', 'Badges bekijken', true),
('badges', 'earn', 'Badges verdienen', true),
('leaderboard', 'view', 'Leaderboards bekijken', true),
('leaderboard', 'participate', 'Deelnemen aan leaderboard', true),

-- Community permissions
('community', 'view', 'Community features bekijken', true),
('community', 'participate', 'Deelnemen aan community activiteiten', true),

-- Admin-only participant permissions
('participant', 'read', 'Alle participants bekijken (admin)', true),
('participant', 'write', 'Participants bewerken (admin)', true),
('participant', 'delete', 'Participants verwijderen (admin)', true),
('participant', 'manage_upgrades', 'Account upgrades beheren (admin)', true),
('participant', 'view_all_registrations', 'Alle registraties bekijken (admin)', true)
ON CONFLICT (resource, action) DO NOTHING;

-- ==============================================================================
-- STAP 4: NIEUWE RBAC ROLLEN VOOR PARTICIPANT SYSTEEM
-- ==============================================================================

-- Voeg participant-specifieke rollen toe
INSERT INTO roles (name, description, is_system_role) VALUES
-- Basis participant rol (voor alle full account users)
('participant_user', 'Participant met full account en app toegang', true),

-- Rol-specifieke rollen (optionele uitbreidingen)
('participant_guide', 'Begeleider met extra rechten', true),
('participant_volunteer', 'Vrijwilliger met extra rechten', true)
ON CONFLICT (name) DO NOTHING;

-- ==============================================================================
-- STAP 5: WIJ PERMISSIONS TOE AAN PARTICIPANT ROLLEN
-- ==============================================================================

-- participant_user rol krijgt basis permissions (ALLE full account users)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'participant_user' 
  AND r.is_system_role = true
  AND (
    -- Eigen data permissions
    (p.resource = 'participant' AND p.action IN ('read_own', 'write_own', 'register_event', 'view_registrations', 'cancel_registration'))
    
    -- App toegang (KRITIEK voor full accounts)
    OR (p.resource = 'app' AND p.action IN ('access', 'login'))
    
    -- Steps & gamification (kern features)
    OR (p.resource = 'steps' AND p.action IN ('track', 'view_own'))
    OR (p.resource = 'achievements' AND p.action IN ('view', 'earn'))
    OR (p.resource = 'badges' AND p.action IN ('view', 'earn'))
    OR (p.resource = 'leaderboard' AND p.action IN ('view', 'participate'))
    
    -- Community features
    OR (p.resource = 'community' AND p.action IN ('view', 'participate'))
  )
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- participant_guide rol krijgt extra permissions (begeleiders)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'participant_guide'
  AND r.is_system_role = true
  AND (
    -- Alle participant_user permissions (inherit base permissions)
    p.id IN (
        SELECT rp.permission_id 
        FROM role_permissions rp 
        JOIN roles r2 ON rp.role_id = r2.id 
        WHERE r2.name = 'participant_user'
    )
    -- Plus extra: kan groepen beheren, anderen ondersteunen
    OR (p.resource = 'community' AND p.action = 'moderate')
  )
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- participant_volunteer rol krijgt extra permissions (vrijwilligers)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'participant_volunteer'
  AND r.is_system_role = true
  AND (
    -- Alle participant_user permissions
    p.id IN (
        SELECT rp.permission_id 
        FROM role_permissions rp 
        JOIN roles r2 ON rp.role_id = r2.id 
        WHERE r2.name = 'participant_user'
    )
    -- Plus extra: event support rechten
    OR (p.resource = 'event' AND p.action = 'support')
  )
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Admin rol krijgt ALLE participant permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin'
  AND r.is_system_role = true
  AND p.resource IN ('participant', 'app', 'steps', 'achievements', 'badges', 'leaderboard', 'community')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- ==============================================================================
-- STAP 6: AUTOMATISCHE ROL TOEWIJZING FUNCTIE
-- ==============================================================================

-- Functie om automatisch participant_user rol toe te wijzen aan nieuwe gebruikers
CREATE OR REPLACE FUNCTION assign_participant_user_role()
RETURNS TRIGGER AS $$
DECLARE
    participant_user_role_id UUID;
BEGIN
    -- Haal participant_user rol ID op
    SELECT id INTO participant_user_role_id
    FROM roles
    WHERE name = 'participant_user' AND is_system_role = true
    LIMIT 1;
    
    -- Als rol bestaat, wijs toe aan nieuwe gebruiker
    IF participant_user_role_id IS NOT NULL THEN
        INSERT INTO user_roles (user_id, role_id, assigned_at, is_active)
        VALUES (NEW.id, participant_user_role_id, CURRENT_TIMESTAMP, true)
        ON CONFLICT (user_id, role_id) DO NOTHING;
        
        RAISE NOTICE 'Participant_user rol toegewezen aan gebruiker %', NEW.id;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger om participant_user rol automatisch toe te wijzen bij nieuwe gebruiker
DROP TRIGGER IF EXISTS trigger_assign_participant_user_role ON gebruikers;
CREATE TRIGGER trigger_assign_participant_user_role
    AFTER INSERT ON gebruikers
    FOR EACH ROW
    EXECUTE FUNCTION assign_participant_user_role();

COMMENT ON FUNCTION assign_participant_user_role() IS 'V30: Wijst automatisch participant_user rol toe aan nieuwe gebruikers';

-- ==============================================================================
-- STAP 7: ROL-SPECIFIEKE RBAC KOPPELING FUNCTIE
-- ==============================================================================

-- Functie om participant event rol te koppelen aan RBAC rol
CREATE OR REPLACE FUNCTION sync_participant_role_to_rbac()
RETURNS TRIGGER AS $$
DECLARE
    rbac_role_name TEXT;
    rbac_role_id UUID;
    participant_gebruiker_id UUID;
BEGIN
    -- Alleen voor full accounts met gebruiker_id
    IF NEW.account_type = 'full' AND NEW.gebruiker_id IS NOT NULL THEN
        -- Haal participant rol op van meest recente event registratie
        SELECT er.participant_role_name INTO rbac_role_name
        FROM event_registrations er
        WHERE er.participant_id = NEW.id
        ORDER BY er.registered_at DESC
        LIMIT 1;
        
        -- Map participant rol naar RBAC rol
        IF rbac_role_name IS NOT NULL THEN
            CASE 
                WHEN LOWER(rbac_role_name) = 'begeleider' THEN
                    rbac_role_name := 'participant_guide';
                WHEN LOWER(rbac_role_name) = 'vrijwilliger' THEN
                    rbac_role_name := 'participant_volunteer';
                ELSE
                    -- Deelnemer of onbekend krijgt geen extra rol (alleen participant_user)
                    rbac_role_name := NULL;
            END CASE;
            
            -- Wijs rol toe als er een mapping is
            IF rbac_role_name IS NOT NULL THEN
                SELECT id INTO rbac_role_id
                FROM roles
                WHERE name = rbac_role_name AND is_system_role = true;
                
                IF rbac_role_id IS NOT NULL THEN
                    INSERT INTO user_roles (user_id, role_id, assigned_at, is_active)
                    VALUES (NEW.gebruiker_id, rbac_role_id, CURRENT_TIMESTAMP, true)
                    ON CONFLICT (user_id, role_id) DO NOTHING;
                    
                    RAISE NOTICE 'Rol % toegewezen aan gebruiker %', rbac_role_name, NEW.gebruiker_id;
                END IF;
            END IF;
        END IF;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger om participant rol te syncen met RBAC bij update
DROP TRIGGER IF EXISTS trigger_sync_participant_role_to_rbac ON participants;
CREATE TRIGGER trigger_sync_participant_role_to_rbac
    AFTER INSERT OR UPDATE OF account_type, gebruiker_id ON participants
    FOR EACH ROW
    EXECUTE FUNCTION sync_participant_role_to_rbac();

COMMENT ON FUNCTION sync_participant_role_to_rbac() IS 'V30: Synchroniseert participant event rol naar RBAC systeem rol';

-- ==============================================================================
-- STAP 8: MIGREER BESTAANDE FULL ACCOUNT PARTICIPANTS
-- ==============================================================================

-- Voor participants die al een gebruiker_id hebben maar nog geen account_type
-- (dit zijn waarschijnlijk oude full accounts)
UPDATE participants 
SET 
    account_type = 'full',
    has_app_access = TRUE,
    wachtwoord_hash = (
        SELECT wachtwoord_hash 
        FROM gebruikers 
        WHERE gebruikers.id = participants.gebruiker_id
    )
WHERE gebruiker_id IS NOT NULL 
  AND account_type = 'temporary'; -- Update alleen als nog niet geüpdatet

-- Wijs participant_user rol toe aan bestaande gebruikers die al een participant zijn
INSERT INTO user_roles (user_id, role_id, assigned_at, is_active)
SELECT 
    g.id,
    r.id,
    CURRENT_TIMESTAMP,
    true
FROM gebruikers g
JOIN participants p ON p.gebruiker_id = g.id
CROSS JOIN roles r
WHERE r.name = 'participant_user' 
  AND r.is_system_role = true
  AND p.account_type = 'full'
  AND p.has_app_access = true
  -- Voorkom duplicates
  AND NOT EXISTS (
      SELECT 1 FROM user_roles ur 
      WHERE ur.user_id = g.id AND ur.role_id = r.id
  )
ON CONFLICT (user_id, role_id) DO NOTHING;

-- ==============================================================================
-- STAP 9: VIEWS VOOR PARTICIPANT PERMISSIONS
-- ==============================================================================

-- View voor makkelijke opzoeking van participant permissions via gebruiker_id
CREATE OR REPLACE VIEW participant_user_permissions AS
SELECT
    p.id as participant_id,
    p.email as participant_email,
    p.account_type,
    p.has_app_access,
    p.gebruiker_id,
    g.email as gebruiker_email,
    r.name as role_name,
    perm.resource,
    perm.action,
    perm.description
FROM participants p
LEFT JOIN gebruikers g ON p.gebruiker_id = g.id
LEFT JOIN user_roles ur ON g.id = ur.user_id AND ur.is_active = true
LEFT JOIN roles r ON ur.role_id = r.id
LEFT JOIN role_permissions rp ON r.id = rp.role_id
LEFT JOIN permissions perm ON rp.permission_id = perm.id
WHERE p.account_type = 'full'
  AND p.has_app_access = true
ORDER BY p.id, r.name, perm.resource, perm.action;

COMMENT ON VIEW participant_user_permissions IS 'V30: Overzicht van alle permissions voor full account participants';

-- ==============================================================================
-- STAP 10: HELPER FUNCTIES VOOR PERMISSION CHECKS
-- ==============================================================================

-- Functie om te checken of een participant specifieke permission heeft
CREATE OR REPLACE FUNCTION participant_has_permission(
    p_participant_id UUID,
    p_resource TEXT,
    p_action TEXT
) RETURNS BOOLEAN AS $$
DECLARE
    has_perm BOOLEAN;
BEGIN
    -- Check via gebruiker_id → user_roles → role_permissions → permissions
    SELECT EXISTS(
        SELECT 1
        FROM participants p
        JOIN gebruikers g ON p.gebruiker_id = g.id
        JOIN user_roles ur ON g.id = ur.user_id
        JOIN role_permissions rp ON ur.role_id = rp.role_id
        JOIN permissions perm ON rp.permission_id = perm.id
        WHERE p.id = p_participant_id
          AND p.account_type = 'full'
          AND p.has_app_access = true
          AND ur.is_active = true
          AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
          AND perm.resource = p_resource
          AND perm.action = p_action
    ) INTO has_perm;
    
    RETURN has_perm;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION participant_has_permission(UUID, TEXT, TEXT) IS 'V30: Check of participant specifieke permission heeft';

-- Functie om te checken of participant app toegang heeft
CREATE OR REPLACE FUNCTION participant_can_access_app(p_participant_id UUID)
RETURNS BOOLEAN AS $$
DECLARE
    can_access BOOLEAN;
BEGIN
    SELECT 
        p.account_type = 'full' 
        AND p.has_app_access = true 
        AND p.gebruiker_id IS NOT NULL
        AND participant_has_permission(p.id, 'app', 'access')
    INTO can_access
    FROM participants p
    WHERE p.id = p_participant_id;
    
    RETURN COALESCE(can_access, false);
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION participant_can_access_app(UUID) IS 'V30: Comprehensive check voor app toegang (account type + flag + permission)';

-- ==============================================================================
-- STAP 11: STATISTIEKEN VIEWS
-- ==============================================================================

-- View voor account type statistieken
CREATE OR REPLACE VIEW participant_account_stats AS
SELECT 
    registration_year,
    account_type,
    COUNT(*) as total_count,
    COUNT(*) FILTER (WHERE has_app_access = true) as with_app_access,
    COUNT(*) FILTER (WHERE gebruiker_id IS NOT NULL) as with_gebruiker,
    COUNT(*) FILTER (WHERE upgraded_at IS NOT NULL) as upgraded_count,
    MIN(created_at) as first_registration,
    MAX(created_at) as last_registration
FROM participants
WHERE registration_year IS NOT NULL
GROUP BY registration_year, account_type
ORDER BY registration_year DESC, account_type;

COMMENT ON VIEW participant_account_stats IS 'V30: Statistieken over account types per jaar';

-- ==============================================================================
-- STAP 12: DATA INTEGRITEIT CHECKS
-- ==============================================================================

-- Constraint: Full accounts MOETEN wachtwoord hebben, gebruiker_id wordt later toegevoegd tijdens registratie
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'check_full_account_requirements' AND conrelid = 'participants'::regclass
    ) THEN
        ALTER TABLE participants ADD CONSTRAINT check_full_account_requirements
            CHECK (
                (account_type = 'temporary') OR
                (account_type = 'full' AND wachtwoord_hash IS NOT NULL)
            );
    END IF;
END $$;

-- Constraint: Temporary accounts MOGEN GEEN gebruiker_id hebben (voorkom verwarring)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'check_temporary_account_no_gebruiker' AND conrelid = 'participants'::regclass
    ) THEN
        ALTER TABLE participants ADD CONSTRAINT check_temporary_account_no_gebruiker
            CHECK (
                (account_type = 'full') OR
                (account_type = 'temporary' AND gebruiker_id IS NULL)
            );
    END IF;
END $$;

-- Constraint: Temporary accounts moeten registration_year hebben
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'check_temporary_account_year' AND conrelid = 'participants'::regclass
    ) THEN
        ALTER TABLE participants ADD CONSTRAINT check_temporary_account_year
            CHECK (
                (account_type = 'full') OR
                (account_type = 'temporary' AND registration_year IS NOT NULL)
            );
    END IF;
END $$;

-- ==============================================================================
-- STAP 13: AUDIT LOGGING SETUP
-- ==============================================================================

-- Tabel voor participant RBAC audit events
CREATE TABLE IF NOT EXISTS participant_rbac_audit (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    participant_id UUID NOT NULL REFERENCES participants(id) ON DELETE CASCADE,
    gebruiker_id UUID REFERENCES gebruikers(id) ON DELETE SET NULL,
    event_type TEXT NOT NULL, -- 'role_assigned', 'role_revoked', 'permission_grant', 'app_access_granted'
    role_name TEXT,
    permission_name TEXT,
    performed_by UUID REFERENCES gebruikers(id),
    performed_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    details JSONB,
    
    CONSTRAINT check_event_type CHECK (
        event_type IN ('role_assigned', 'role_revoked', 'permission_grant', 'app_access_granted', 'account_upgraded')
    )
);

CREATE INDEX IF NOT EXISTS idx_participant_rbac_audit_participant ON participant_rbac_audit(participant_id);
CREATE INDEX IF NOT EXISTS idx_participant_rbac_audit_gebruiker ON participant_rbac_audit(gebruiker_id);
CREATE INDEX IF NOT EXISTS idx_participant_rbac_audit_event_type ON participant_rbac_audit(event_type);
CREATE INDEX IF NOT EXISTS idx_participant_rbac_audit_performed_at ON participant_rbac_audit(performed_at);

COMMENT ON TABLE participant_rbac_audit IS 'V30: Audit trail voor alle RBAC-gerelateerde participant acties';

-- ==============================================================================
-- STAP 14: VERIFICATIE QUERIES
-- ==============================================================================

-- Verificatie: Toon account type distributie
DO $$
DECLARE
    full_count INTEGER;
    temp_count INTEGER;
    with_access INTEGER;
BEGIN
    SELECT COUNT(*) INTO full_count FROM participants WHERE account_type = 'full';
    SELECT COUNT(*) INTO temp_count FROM participants WHERE account_type = 'temporary';
    SELECT COUNT(*) INTO with_access FROM participants WHERE has_app_access = true;
    
    RAISE NOTICE '=== V30 RBAC INTEGRATIE VERIFICATIE ===';
    RAISE NOTICE 'Full accounts: %', full_count;
    RAISE NOTICE 'Temporary accounts: %', temp_count;
    RAISE NOTICE 'Met app access: %', with_access;
    
    -- Verificatie: Alle full accounts hebben gebruiker_id
    IF EXISTS (SELECT 1 FROM participants WHERE account_type = 'full' AND gebruiker_id IS NULL) THEN
        RAISE WARNING 'WAARSCHUWING: Er zijn full accounts zonder gebruiker_id!';
    ELSE
        RAISE NOTICE '✓ Alle full accounts hebben gebruiker_id';
    END IF;
    
    -- Verificatie: Participant_user rol bestaat en heeft permissions
    IF EXISTS (SELECT 1 FROM roles WHERE name = 'participant_user') THEN
        RAISE NOTICE '✓ participant_user rol bestaat';
        
        SELECT COUNT(*) INTO full_count 
        FROM role_permissions rp
        JOIN roles r ON rp.role_id = r.id
        WHERE r.name = 'participant_user';
        
        RAISE NOTICE '  - Heeft % permissions', full_count;
    ELSE
        RAISE WARNING 'WAARSCHUWING: participant_user rol bestaat niet!';
    END IF;
END $$;

-- ==============================================================================
-- MIGRATIE VOLTOOID
-- ==============================================================================

-- Commentaren voor documentatie
COMMENT ON COLUMN participants.account_type IS 'V30: Account type - full (met app) of temporary (alleen event)';
COMMENT ON COLUMN participants.registration_year IS 'V30: Voor temporary accounts - jaar van registratie';
COMMENT ON COLUMN participants.wachtwoord_hash IS 'V30: Bcrypt hash van wachtwoord (alleen full accounts)';
COMMENT ON COLUMN participants.has_app_access IS 'V30: Expliciete flag voor app toegang';
COMMENT ON COLUMN participants.upgraded_to_gebruiker_id IS 'V30: Tracking - naar welke gebruiker geüpgraded';
COMMENT ON COLUMN participants.upgraded_at IS 'V30: Tijdstip van upgrade naar full account';