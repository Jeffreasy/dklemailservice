-- V33: Create auto_responses table
-- Description: Creates table for managing email auto-response configurations

-- Create auto_responses table
CREATE TABLE IF NOT EXISTS auto_responses (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,
    subject VARCHAR(255),
    message TEXT,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index on email for faster lookups
CREATE INDEX IF NOT EXISTS idx_auto_responses_email ON auto_responses(email);

-- Create index on is_active for filtering active responses
CREATE INDEX IF NOT EXISTS idx_auto_responses_active ON auto_responses(is_active);

-- Add updated_at trigger
CREATE OR REPLACE FUNCTION update_auto_responses_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER auto_responses_updated_at
    BEFORE UPDATE ON auto_responses
    FOR EACH ROW
    EXECUTE FUNCTION update_auto_responses_updated_at();