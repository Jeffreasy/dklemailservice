-- GECONSOLIDEERDE V17 - ROUTE FUNDS TABEL
-- Logica van V1_46.
-- De trigger-functie is verwijderd; V19 zal een generieke trigger toevoegen.

-- Create route_funds table for configurable fund allocation per route
CREATE TABLE IF NOT EXISTS route_funds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    route VARCHAR(50) NOT NULL UNIQUE,
    amount INTEGER NOT NULL CHECK (amount >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on route for fast lookups (only if it doesn't exist)
CREATE INDEX IF NOT EXISTS idx_route_funds_route ON route_funds(route);

-- Insert default values (only if they don't exist)
INSERT INTO route_funds (route, amount)
VALUES
    ('6 KM', 50),
    ('10 KM', 75),
    ('15 KM', 100),
    ('20 KM', 125)
ON CONFLICT (route) DO NOTHING;