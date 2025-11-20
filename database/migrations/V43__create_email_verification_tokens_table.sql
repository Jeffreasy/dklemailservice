-- Create email verification tokens table for email verification functionality
-- Migration: V43__create_email_verification_tokens_table.sql

CREATE TABLE IF NOT EXISTS email_verification_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    is_used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    used_at TIMESTAMP NULL,
    user_id VARCHAR(255) NOT NULL, -- UUID of the user (gebruiker or participant)
    user_type VARCHAR(50) NOT NULL -- 'gebruiker' or 'participant'
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_email_verification_tokens_token ON email_verification_tokens(token);
CREATE INDEX IF NOT EXISTS idx_email_verification_tokens_email ON email_verification_tokens(email);
CREATE INDEX IF NOT EXISTS idx_email_verification_tokens_expires_at ON email_verification_tokens(expires_at);
CREATE INDEX IF NOT EXISTS idx_email_verification_tokens_is_used ON email_verification_tokens(is_used);
CREATE INDEX IF NOT EXISTS idx_email_verification_tokens_user_id ON email_verification_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_email_verification_tokens_user_type ON email_verification_tokens(user_type);

-- Add composite index for user lookup
CREATE INDEX IF NOT EXISTS idx_email_verification_tokens_user ON email_verification_tokens(user_id, user_type);

-- Add comment to table
COMMENT ON TABLE email_verification_tokens IS 'Stores email verification tokens for account activation';
COMMENT ON COLUMN email_verification_tokens.email IS 'Email address of the user to be verified';
COMMENT ON COLUMN email_verification_tokens.token IS 'Base64 encoded verification token sent via email';
COMMENT ON COLUMN email_verification_tokens.expires_at IS 'Token expiration timestamp (24 hours from creation)';
COMMENT ON COLUMN email_verification_tokens.is_used IS 'Whether the token has been used for email verification';
COMMENT ON COLUMN email_verification_tokens.used_at IS 'Timestamp when the token was used';
COMMENT ON COLUMN email_verification_tokens.user_id IS 'UUID of the user (references gebruiker.id or participant.id)';
COMMENT ON COLUMN email_verification_tokens.user_type IS 'Type of user: gebruiker or participant';