-- V47: Create sessions table for session management
-- This migration creates the sessions table to track user sessions across devices

-- Create sessions table (idempotent - safe to run multiple times)
DO $$
BEGIN
    -- Create sessions table if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'sessions') THEN
        CREATE TABLE sessions (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            owner_id UUID NOT NULL, -- Can reference either gebruikers.id or participants.id
            owner_type VARCHAR(20) NOT NULL DEFAULT 'gebruiker' CHECK (owner_type IN ('gebruiker', 'participant')),
            access_token VARCHAR(500) UNIQUE, -- Link to current access token
            device_info JSONB, -- Device fingerprinting data
            ip_address INET, -- Client IP address
            user_agent TEXT, -- Browser/client user agent
            location_info JSONB, -- Optional location data
            is_active BOOLEAN NOT NULL DEFAULT true,
            is_current BOOLEAN NOT NULL DEFAULT false, -- Marks the current session for this user
            expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
            last_activity TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
            created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
            updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
        );

        -- Create indexes for performance
        CREATE INDEX idx_sessions_owner_id ON sessions(owner_id);
        CREATE INDEX idx_sessions_owner_type ON sessions(owner_type);
        CREATE INDEX idx_sessions_access_token ON sessions(access_token);
        CREATE INDEX idx_sessions_is_active ON sessions(is_active);
        CREATE INDEX idx_sessions_is_current ON sessions(is_current);
        CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
        CREATE INDEX idx_sessions_last_activity ON sessions(last_activity);

        -- Create partial index for active sessions
        CREATE INDEX idx_sessions_active_recent ON sessions(owner_id, last_activity DESC)
        WHERE is_active = true;

        -- Add comments
        COMMENT ON TABLE sessions IS 'Tracks user sessions for multi-device authentication and session management';
        COMMENT ON COLUMN sessions.owner_id IS 'ID of the user (gebruiker) or participant who owns this session';
        COMMENT ON COLUMN sessions.owner_type IS 'Type of owner: gebruiker or participant';
        COMMENT ON COLUMN sessions.access_token IS 'Current access token for this session';
        COMMENT ON COLUMN sessions.device_info IS 'Device fingerprinting information (OS, browser, etc.)';
        COMMENT ON COLUMN sessions.ip_address IS 'Client IP address when session was created';
        COMMENT ON COLUMN sessions.user_agent IS 'Browser/client user agent string';
        COMMENT ON COLUMN sessions.location_info IS 'Optional location data derived from IP';
        COMMENT ON COLUMN sessions.is_active IS 'Whether this session is still active';
        COMMENT ON COLUMN sessions.is_current IS 'Whether this is the current session for the user';
        COMMENT ON COLUMN sessions.expires_at IS 'When this session expires';
        COMMENT ON COLUMN sessions.last_activity IS 'Last activity timestamp for this session';
    END IF;
END $$;

-- Add foreign key constraint (optional - allows NULL for flexibility)
-- ALTER TABLE sessions ADD CONSTRAINT fk_sessions_owner_gebruiker
--     FOREIGN KEY (owner_id) REFERENCES gebruikers(id) ON DELETE CASCADE
--     WHERE owner_type = 'gebruiker';

-- Add comments
COMMENT ON TABLE sessions IS 'Tracks user sessions for multi-device authentication and session management';
COMMENT ON COLUMN sessions.owner_id IS 'ID of the user (gebruiker) or participant who owns this session';
COMMENT ON COLUMN sessions.owner_type IS 'Type of owner: gebruiker or participant';
COMMENT ON COLUMN sessions.access_token IS 'Current access token for this session';
COMMENT ON COLUMN sessions.device_info IS 'Device fingerprinting information (OS, browser, etc.)';
COMMENT ON COLUMN sessions.ip_address IS 'Client IP address when session was created';
COMMENT ON COLUMN sessions.user_agent IS 'Browser/client user agent string';
COMMENT ON COLUMN sessions.location_info IS 'Optional location data derived from IP';
COMMENT ON COLUMN sessions.is_active IS 'Whether this session is still active';
COMMENT ON COLUMN sessions.is_current IS 'Whether this is the current session for the user';
COMMENT ON COLUMN sessions.expires_at IS 'When this session expires';
COMMENT ON COLUMN sessions.last_activity IS 'Last activity timestamp for this session';