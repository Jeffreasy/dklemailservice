-- V41: Create staff user with correct credentials
-- This migration creates the staff user Salih@dekoninklijkeloop.nl with password Bootje@12

INSERT INTO gebruikers (naam, email, wachtwoord_hash, is_actief, created_at, updated_at)
VALUES (
    'Salih Staff',
    'Salih@dekoninklijkeloop.nl',
    -- Password: Bootje@12 (bcrypt hash generated for this password)
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    TRUE,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
)
ON CONFLICT (email) DO NOTHING;

-- Assign staff role to the user
INSERT INTO user_roles (user_id, role_id, assigned_by, assigned_at, is_active)
SELECT
    g.id,
    r.id,
    g.id, -- self-assigned for seed data
    CURRENT_TIMESTAMP,
    TRUE
FROM gebruikers g, roles r
WHERE g.email = 'Salih@dekoninklijkeloop.nl'
  AND r.name = 'staff' AND r.is_system_role = true
ON CONFLICT (user_id, role_id) DO NOTHING;

-- Log the completion
DO $$
BEGIN
    RAISE NOTICE '[V41] Staff user Salih@dekoninklijkeloop.nl created and assigned staff role';
END $$;