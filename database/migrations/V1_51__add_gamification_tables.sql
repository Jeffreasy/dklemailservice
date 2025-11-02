-- Migration: V1_51__add_gamification_tables.sql
-- Description: Voegt tabellen toe voor badges, achievements en leaderboard functionaliteit
-- Date: 2025-01-02
--
-- IDEMPOTENT: Deze migratie kan veilig meerdere keren worden uitgevoerd

-- =====================================================
-- 1. BADGES TABLE
-- =====================================================
-- Badges die admin kan aanmaken en beheren
CREATE TABLE IF NOT EXISTS badges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    icon_url VARCHAR(500),
    criteria JSONB NOT NULL DEFAULT '{}',
    points INTEGER NOT NULL DEFAULT 0 CHECK (points >= 0),
    is_active BOOLEAN NOT NULL DEFAULT true,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index voor actieve badges query
CREATE INDEX IF NOT EXISTS idx_badges_active ON badges(is_active, display_order);
CREATE INDEX IF NOT EXISTS idx_badges_points ON badges(points);

-- Trigger voor updated_at
DROP TRIGGER IF EXISTS update_badges_updated_at ON badges;
CREATE TRIGGER update_badges_updated_at
    BEFORE UPDATE ON badges
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments (alleen als tabel bestaat)
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'badges') THEN
        COMMENT ON TABLE badges IS 'Badges die deelnemers kunnen verdienen (admin-beheerd)';
        COMMENT ON COLUMN badges.criteria IS 'JSON criteria voor badge verdienen, bijv: {"min_steps": 10000, "min_days": 7}';
        COMMENT ON COLUMN badges.points IS 'Punten die deelnemer krijgt bij verdienen van badge';
        COMMENT ON COLUMN badges.display_order IS 'Volgorde waarin badges worden getoond';
    END IF;
END $$;

-- =====================================================
-- 2. PARTICIPANT ACHIEVEMENTS TABLE
-- =====================================================
-- Verdiende badges per deelnemer
CREATE TABLE IF NOT EXISTS participant_achievements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    participant_id UUID NOT NULL,
    badge_id UUID NOT NULL,
    earned_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Voeg constraints toe (alleen als ze nog niet bestaan)
DO $$
BEGIN
    -- Foreign key naar aanmeldingen
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'fk_participant' 
        AND table_name = 'participant_achievements'
    ) THEN
        ALTER TABLE participant_achievements
        ADD CONSTRAINT fk_participant 
        FOREIGN KEY (participant_id) REFERENCES aanmeldingen(id) ON DELETE CASCADE;
    END IF;
    
    -- Foreign key naar badges
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'fk_badge' 
        AND table_name = 'participant_achievements'
    ) THEN
        ALTER TABLE participant_achievements
        ADD CONSTRAINT fk_badge 
        FOREIGN KEY (badge_id) REFERENCES badges(id) ON DELETE CASCADE;
    END IF;
    
    -- Unique constraint
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'unique_participant_badge' 
        AND table_name = 'participant_achievements'
    ) THEN
        ALTER TABLE participant_achievements
        ADD CONSTRAINT unique_participant_badge UNIQUE(participant_id, badge_id);
    END IF;
END $$;

-- Indexes voor snelle queries
CREATE INDEX IF NOT EXISTS idx_participant_achievements_participant ON participant_achievements(participant_id);
CREATE INDEX IF NOT EXISTS idx_participant_achievements_badge ON participant_achievements(badge_id);
CREATE INDEX IF NOT EXISTS idx_participant_achievements_earned_at ON participant_achievements(earned_at DESC);

-- Comments
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'participant_achievements') THEN
        COMMENT ON TABLE participant_achievements IS 'Verdiende badges per deelnemer';
    END IF;
END $$;

-- =====================================================
-- 3. LEADERBOARD VIEW
-- =====================================================
-- View voor leaderboard met ranking
CREATE OR REPLACE VIEW leaderboard_view AS
SELECT 
    a.id,
    a.naam,
    a.afstand as route,
    a.steps,
    COALESCE(SUM(b.points), 0) as achievement_points,
    a.steps + COALESCE(SUM(b.points), 0) as total_score,
    RANK() OVER (ORDER BY (a.steps + COALESCE(SUM(b.points), 0)) DESC) as rank,
    COUNT(pa.id) as badge_count,
    a.created_at as joined_at
FROM aanmeldingen a
LEFT JOIN participant_achievements pa ON a.id = pa.participant_id
LEFT JOIN badges b ON pa.badge_id = b.id AND b.is_active = true
GROUP BY a.id, a.naam, a.afstand, a.steps, a.created_at
ORDER BY total_score DESC, a.steps DESC;

COMMENT ON VIEW leaderboard_view IS 'Leaderboard met ranking gebaseerd op steps + achievement points';

-- =====================================================
-- 4. SEED DATA - DEFAULT BADGES
-- =====================================================
-- Voeg enkele standaard badges toe (ON CONFLICT DO NOTHING zorgt voor idempotency)
INSERT INTO badges (name, description, icon_url, criteria, points, display_order) VALUES
('First Steps', 'Je eerste 1000 stappen gezet', '/icons/badges/first-steps.svg', '{"min_steps": 1000}', 10, 1),
('5K Champion', '5000 stappen bereikt', '/icons/badges/5k-champion.svg', '{"min_steps": 5000}', 50, 2),
('10K Master', '10000 stappen bereikt', '/icons/badges/10k-master.svg', '{"min_steps": 10000}', 100, 3),
('Marathon Walker', '42195 stappen bereikt (marathon)', '/icons/badges/marathon.svg', '{"min_steps": 42195}', 500, 4),
('Early Bird', 'Binnen eerste 50 deelnemers ingeschreven', '/icons/badges/early-bird.svg', '{"early_participant": true}', 25, 5),
('Consistent Walker', '7 dagen op rij stappen gelogd', '/icons/badges/consistent.svg', '{"consecutive_days": 7}', 75, 6),
('Team Player', 'Deel van een team', '/icons/badges/team-player.svg', '{"has_team": true}', 20, 7),
('Distance Hero', 'Langste route gekozen (20 KM)', '/icons/badges/distance-hero.svg', '{"route": "20 KM"}', 150, 8)
ON CONFLICT (name) DO NOTHING;

-- =====================================================
-- 5. HELPER FUNCTIONS
-- =====================================================
-- Functie om participant total score te berekenen
CREATE OR REPLACE FUNCTION get_participant_score(p_participant_id UUID)
RETURNS INTEGER AS $$
DECLARE
    v_steps INTEGER;
    v_achievement_points INTEGER;
BEGIN
    -- Haal steps op
    SELECT steps INTO v_steps
    FROM aanmeldingen
    WHERE id = p_participant_id;
    
    -- Haal achievement points op
    SELECT COALESCE(SUM(b.points), 0) INTO v_achievement_points
    FROM participant_achievements pa
    JOIN badges b ON pa.badge_id = b.id
    WHERE pa.participant_id = p_participant_id
    AND b.is_active = true;
    
    RETURN COALESCE(v_steps, 0) + COALESCE(v_achievement_points, 0);
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_participant_score IS 'Berekent totale score voor deelnemer (steps + achievement points)';

-- Functie om te checken of participant badge heeft verdiend
CREATE OR REPLACE FUNCTION check_badge_eligibility(
    p_participant_id UUID,
    p_badge_id UUID
) RETURNS BOOLEAN AS $$
DECLARE
    v_criteria JSONB;
    v_steps INTEGER;
    v_route VARCHAR(50);
    v_participant_count INTEGER;
BEGIN
    -- Haal badge criteria op
    SELECT criteria INTO v_criteria
    FROM badges
    WHERE id = p_badge_id AND is_active = true;
    
    IF v_criteria IS NULL THEN
        RETURN false;
    END IF;
    
    -- Haal participant data op
    SELECT steps, afstand INTO v_steps, v_route
    FROM aanmeldingen
    WHERE id = p_participant_id;
    
    -- Check min_steps criterium
    IF v_criteria ? 'min_steps' THEN
        IF v_steps < (v_criteria->>'min_steps')::INTEGER THEN
            RETURN false;
        END IF;
    END IF;
    
    -- Check route criterium
    IF v_criteria ? 'route' THEN
        IF v_route != v_criteria->>'route' THEN
            RETURN false;
        END IF;
    END IF;
    
    -- Check early_participant criterium
    IF v_criteria ? 'early_participant' AND (v_criteria->>'early_participant')::BOOLEAN THEN
        SELECT COUNT(*) INTO v_participant_count
        FROM aanmeldingen
        WHERE created_at <= (SELECT created_at FROM aanmeldingen WHERE id = p_participant_id);
        
        IF v_participant_count > 50 THEN
            RETURN false;
        END IF;
    END IF;
    
    RETURN true;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION check_badge_eligibility IS 'Controleert of participant badge heeft verdiend op basis van criteria';

-- =====================================================
-- 6. PERMISSIONS
-- =====================================================
-- Voeg permissions toe voor gamification
INSERT INTO permissions (resource, action, description) VALUES
('badges', 'read', 'Badges kunnen bekijken'),
('badges', 'write', 'Badges kunnen aanmaken en bewerken'),
('achievements', 'read', 'Achievements kunnen bekijken'),
('achievements', 'write', 'Achievements kunnen toekennen'),
('leaderboard', 'read', 'Leaderboard kunnen bekijken')
ON CONFLICT (resource, action) DO NOTHING;

-- Ken permissions toe aan rollen
-- Admin krijgt volledige toegang
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'admin'
AND p.resource IN ('badges', 'achievements', 'leaderboard')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Staff krijgt read toegang tot badges en achievements, write voor achievements toekennen
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'staff'
AND (
    (p.resource = 'badges' AND p.action = 'read')
    OR (p.resource = 'achievements' AND p.action IN ('read', 'write'))
    OR (p.resource = 'leaderboard' AND p.action = 'read')
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Deelnemer krijgt read toegang tot leaderboard en eigen achievements
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'deelnemer'
AND p.resource IN ('leaderboard', 'achievements')
AND p.action = 'read'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- =====================================================
-- MIGRATION COMPLETE
-- =====================================================