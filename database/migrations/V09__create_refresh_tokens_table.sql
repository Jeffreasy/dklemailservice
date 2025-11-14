-- GECONSOLIDEERDE V09 - REFRESH TOKENS (V34 Updated)
-- Dit is de logica van V1_28.
-- FIX: Omgezet naar TIMESTAMPTZ voor GORM-compatibiliteit.
-- V34: Changed user_id to owner_id (no FK constraint) to support both gebruikers and participants

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL,  -- V34: Was user_id, now owner_id (can be gebruiker or participant)
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMPTZ,
    is_revoked BOOLEAN DEFAULT FALSE
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_owner_id ON refresh_tokens(owner_id);  -- V34: Was user_id
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(owner_id);  -- Legacy index name for compatibility
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_is_revoked ON refresh_tokens(is_revoked);

-- Cleanup index for expired/revoked tokens
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_cleanup ON refresh_tokens(expires_at) WHERE is_revoked = FALSE;

-- Comment on table
COMMENT ON TABLE refresh_tokens IS 'Stores refresh tokens for JWT authentication with 7-day expiry';
COMMENT ON COLUMN refresh_tokens.owner_id IS 'Can reference either gebruikers.id (admin/staff) or participants.id (deelnemers). No FK constraint to support both types.';
COMMENT ON COLUMN refresh_tokens.token IS 'Base64 encoded random token (32 bytes)';
COMMENT ON COLUMN refresh_tokens.expires_at IS 'Token expiration timestamp (7 days from creation)';
COMMENT ON COLUMN refresh_tokens.is_revoked IS 'Whether the token has been revoked (for token rotation)';