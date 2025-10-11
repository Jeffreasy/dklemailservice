-- Migratie: V1_31I__migrate_under_construction.sql
-- Beschrijving: Migrate under_construction data from Supabase
-- Versie: 1.31.0

-- Create under_construction table
CREATE TABLE IF NOT EXISTS under_construction (
    id SERIAL PRIMARY KEY,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    title TEXT,
    message TEXT,
    footer_text TEXT,
    logo_url TEXT,
    expected_date TIMESTAMP WITH TIME ZONE,
    social_links JSONB,
    progress_percentage INTEGER,
    contact_email TEXT,
    newsletter_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Insert under_construction data
INSERT INTO under_construction (id, is_active, title, message, footer_text, logo_url, expected_date, social_links, progress_percentage, contact_email, newsletter_enabled, created_at, updated_at) VALUES
(1, false, 'Website in onderhoud', 'We stomen ons klaar voor De Koninklijke Loop 2026, op dit moment is de website helaas niet bereikbaar', 'Bedankt voor uw geduld!', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png', '2026-01-31 18:00:00+00', '[{"url": "https://twitter.com/koninklijkeloop", "platform": "Twitter"}, {"url": "https://instagram.com/koninklijkeloop", "platform": "Instagram"}, {"url": "https://www.youtube.com/@DeKoninklijkeLoop", "platform": "YouTube"}]', 85, 'info@koninklijkeloop.nl', false, '2025-09-26 17:37:22.197854+00', '2025-10-09 21:00:29.391392+00')
ON CONFLICT (id) DO NOTHING;