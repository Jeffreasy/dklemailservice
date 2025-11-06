-- Migratie: V1_41__migrate_title_sections.sql
-- Beschrijving: Migrate title_sections data from Supabase
-- Versie: 1.41.0

-- Create title_sections table
CREATE TABLE IF NOT EXISTS title_sections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    subtitle TEXT,
    cta_text TEXT,
    image_url TEXT,
    event_details JSONB,
    styling JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Insert title_sections data
INSERT INTO title_sections (id, title, subtitle, cta_text, image_url, event_details, created_at, updated_at, styling) VALUES
('8b52223c-f86e-4329-8cb5-dd7cfbe4aae0', 'De Koninklijke Loop 2025', 'op de koninklijke weg kunnen mensen met een beperking samen wandelen met hun verwanten, vrijwilligers of begeleiders', 'inschrijv', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1734323119/oo09xs7ahavvplwomj30.jpg', '[{"icon": "Wanneer?", "title": "17 mei 2025", "description": "Start om nog onbekend"}, {"icon": "Voor wie?", "title": "Voor iedereen", "description": "Alle leeftijden welkom"}, {"icon": "Waarvoor?", "title": "Lopen voor een goed doel", "description": "Voor verschillende categorieën"}]', '2024-12-22 19:30:53.366424+00', '2024-12-22 19:30:53.366424+00', '{"title": {"color": "#000000", "fontSize": "32px", "textAlign": "center", "fontWeight": "bold", "lineHeight": "1.5", "letterSpacing": "0px"}, "cta_text": {"color": "#000000", "fontSize": "24px", "textAlign": "center", "fontWeight": "semibold", "lineHeight": "1.5", "letterSpacing": "0px"}, "subtitle": {"color": "#000000", "fontSize": "24px", "textAlign": "center", "fontWeight": "normal", "lineHeight": "1.5", "letterSpacing": "0px"}}')
ON CONFLICT (id) DO NOTHING;