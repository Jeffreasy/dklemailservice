-- Create access tokens table for server-side JWT token storage and invalidation
-- Migration: V44__create_access_tokens_table.sql
-- This enables token rotation and immediate invalidation of access tokens

CREATE TABLE IF NOT EXISTS access_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id VARCHAR(255) NOT NULL, -- UUID of the user (gebruiker or participant)
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMP NULL,
    is_revoked BOOLEAN DEFAULT FALSE
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_access_tokens_token ON access_tokens(token);
CREATE INDEX IF NOT EXISTS idx_access_tokens_owner_id ON access_tokens(owner_id);
CREATE INDEX IF NOT EXISTS idx_access_tokens_expires_at ON access_tokens(expires_at);
CREATE INDEX IF NOT EXISTS idx_access_tokens_is_revoked ON access_tokens(is_revoked);

-- Add composite index for owner lookup
CREATE INDEX IF NOT EXISTS idx_access_tokens_owner ON access_tokens(owner_id, is_revoked);

-- Add comment to table
COMMENT ON TABLE access_tokens IS 'Stores access tokens for server-side validation and invalidation';
COMMENT ON COLUMN access_tokens.owner_id IS 'UUID of the token owner (references gebruiker.id or participant.id)';
COMMENT ON COLUMN access_tokens.token IS 'Base64 encoded JWT access token';
COMMENT ON COLUMN access_tokens.expires_at IS 'Token expiration timestamp (20 minutes from creation)';
COMMENT ON COLUMN access_tokens.is_revoked IS 'Whether the token has been revoked/invalidated';
COMMENT ON COLUMN access_tokens.revoked_at IS 'Timestamp when the token was revoked';