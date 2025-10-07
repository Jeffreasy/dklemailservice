-- Migratie: V1_20__create_rbac_tables.sql
-- Beschrijving: Create Role-Based Access Control (RBAC) tables for flexible permission management
-- Versie: 1.20.0

-- Roles table - defines available roles in the system
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    is_system_role BOOLEAN NOT NULL DEFAULT FALSE, -- System roles cannot be deleted
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES gebruikers(id), -- Who created this role
    UNIQUE(name)
);

-- Permissions table - defines granular permissions (resource + action)
CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource VARCHAR(100) NOT NULL, -- e.g., 'contact', 'aanmelding', 'user', 'newsletter'
    action VARCHAR(50) NOT NULL,    -- e.g., 'create', 'read', 'update', 'delete', 'manage'
    description TEXT,
    is_system_permission BOOLEAN NOT NULL DEFAULT FALSE, -- System permissions cannot be deleted
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(resource, action)
);

-- Role-Permission relationship table
CREATE TABLE IF NOT EXISTS role_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    assigned_by UUID REFERENCES gebruikers(id), -- Who assigned this permission
    UNIQUE(role_id, permission_id)
);

-- User-Role relationship table (for multiple roles per user if needed)
CREATE TABLE IF NOT EXISTS user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES gebruikers(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    assigned_by UUID REFERENCES gebruikers(id), -- Who assigned this role
    expires_at TIMESTAMP WITH TIME ZONE, -- Optional expiration
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE(user_id, role_id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_roles_name ON roles(name);
CREATE INDEX IF NOT EXISTS idx_permissions_resource_action ON permissions(resource, action);
CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON role_permissions(permission_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON user_roles(role_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_active ON user_roles(is_active) WHERE is_active = true;

-- Add role_id column to gebruikers table for backward compatibility
-- This allows gradual migration from string-based roles to UUID-based roles
ALTER TABLE gebruikers ADD COLUMN IF NOT EXISTS role_id UUID REFERENCES roles(id);

-- Create a view for easy querying of user permissions
CREATE OR REPLACE VIEW user_permissions AS
SELECT
    ur.user_id,
    u.email,
    r.name as role_name,
    p.resource,
    p.action,
    rp.assigned_at as permission_assigned_at,
    ur.assigned_at as role_assigned_at
FROM user_roles ur
JOIN roles r ON ur.role_id = r.id
JOIN role_permissions rp ON r.id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
JOIN gebruikers u ON ur.user_id = u.id
WHERE ur.is_active = true
ORDER BY ur.user_id, r.name, p.resource, p.action;

-- Registreer de migratie
INSERT INTO migraties (versie, naam, toegepast)
VALUES ('1.20.0', 'Create RBAC tables for flexible permission management', CURRENT_TIMESTAMP)
ON CONFLICT (versie) DO NOTHING;