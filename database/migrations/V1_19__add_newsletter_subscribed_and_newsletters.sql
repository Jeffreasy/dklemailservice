-- Add newsletter_subscribed to gebruikers
ALTER TABLE IF EXISTS gebruikers
    ADD COLUMN IF NOT EXISTS newsletter_subscribed boolean DEFAULT false;

-- Create newsletters table
CREATE TABLE IF NOT EXISTS newsletters (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    subject text NOT NULL,
    content text NOT NULL,
    sent_at timestamp with time zone,
    batch_id text,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_newsletters_sent_at ON newsletters (sent_at);

