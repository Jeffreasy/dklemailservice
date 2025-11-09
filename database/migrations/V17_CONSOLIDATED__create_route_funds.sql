-- ============================================================================
-- V17 CONSOLIDATED: Create Route Funds Table
-- ============================================================================
-- Consolidates V17_01, V17_02, V17_03
-- Purpose: Create route_funds table for fundraising per distance
-- Note: This table will be renamed to 'distances' in V26
-- ============================================================================

-- Create route_funds table
CREATE TABLE IF NOT EXISTS route_funds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    route VARCHAR(50) NOT NULL UNIQUE,
    amount INTEGER NOT NULL CHECK (amount >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index
CREATE INDEX IF NOT EXISTS idx_route_funds_route ON route_funds(route);

-- Seed routes
INSERT INTO route_funds (route, amount) VALUES
    ('2.5 KM', 25),
    ('6 KM', 50),
    ('10 KM', 75),
    ('15 KM', 100),
    ('20 KM', 125)
ON CONFLICT (route) DO NOTHING;

COMMENT ON TABLE route_funds IS 'Fondsenwerving doelen per route (V17, wordt distances in V26)';