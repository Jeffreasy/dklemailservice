-- Create route_funds table for configurable fund allocation per route
CREATE TABLE route_funds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    route VARCHAR(50) NOT NULL UNIQUE,
    amount INTEGER NOT NULL CHECK (amount >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on route for fast lookups
CREATE INDEX idx_route_funds_route ON route_funds(route);

-- Insert default values (can be modified by admin)
INSERT INTO route_funds (route, amount) VALUES
    ('6 KM', 50),
    ('10 KM', 75),
    ('15 KM', 100),
    ('20 KM', 125);

-- Create trigger for updated_at
CREATE OR REPLACE FUNCTION update_route_funds_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_route_funds_updated_at
    BEFORE UPDATE ON route_funds
    FOR EACH ROW
    EXECUTE FUNCTION update_route_funds_updated_at();