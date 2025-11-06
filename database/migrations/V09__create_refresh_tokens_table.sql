-- GECONSOLIDEERDE V09 - REFRESH TOKENS
-- Dit is de logica van V1_28.
-- FIX: Omgezet naar TIMESTAMPTZ voor GORM-compatibiliteit.

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES gebruikers(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMPTZ,
    is_revoked BOOLEAN DEFAULT FALSE
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_is_revoked ON refresh_tokens(is_revoked);

-- Comment on table
COMMENT ON TABLE refresh_tokens IS 'Stores refresh tokens for JWT authentication with 7-day expiry';
COMMENT ON COLUMN refresh_tokens.token IS 'Base64 encoded random token (32 bytes)';
COMMENT ON COLUMN refresh_tokens.expires_at IS 'Token expiration timestamp (7 days from creation)';
COMMENT ON COLUMN refresh_tokens.is_revoked IS 'Whether the token has been revoked (for token rotation)';