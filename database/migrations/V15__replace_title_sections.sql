-- GECONSOLIDEERDE V15 - VERVANG TITLE SECTIONS
-- Dit script is de logica van V1_43.
-- Het vervangt de (nooit gebruikte) 'title_sections' tabel door 'title_section_content'.
-- De 'ALTER TABLE ... ADD COLUMN steps' is hieruit gehaald en naar V16 verplaatst.

-- Drop de oude tabel die in V1_41 was aangemaakt
DROP TABLE IF EXISTS title_sections;

-- Create de nieuwe title_section_content tabel
CREATE TABLE IF NOT EXISTS title_section_content (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_title TEXT NOT NULL,
    event_subtitle TEXT,
    image_url TEXT,
    image_alt TEXT,
    detail_1_title TEXT,
    detail_1_description TEXT,
    detail_2_title TEXT,
    detail_2_description TEXT,
    detail_3_title TEXT,
    detail_3_description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    participant_count INTEGER DEFAULT 0
);

-- Voeg de content in
INSERT INTO title_section_content (id, event_title, event_subtitle, image_url, image_alt, detail_1_title, detail_1_description, detail_2_title, detail_2_description, detail_3_title, detail_3_description, created_at, updated_at, participant_count) VALUES
('550e8400-e29b-41d4-a716-446655440001', 'De Koninklijke Loop (DKL) 2025', 'Op de koninklijke weg in Apeldoorn kunnen mensen met een beperking samen wandelen tijdens dit unieke, rolstoelvriendelijke sponsorloop (DKL), samen met hun verwanten, vrijwilligers of begeleiders.', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1760112848/Wij_gaan_17_mei_lopen_voor_hen_3_zllxno_zoqd7z.webp', 'Promotiebanner De Koninklijke Loop (DKL) 2025: Wij gaan 17 mei lopen voor hen', '17 mei 2025', 'Starttijden variëren per afstand. Zie programma.', 'Voor iedereen', 'wandelaars met of zonder beperking (rolstoelvriendelijk).', 'Lopen voor een goed doel', 'Steun het goede doel via dit unieke wandelevenement.', '2025-04-16 01:31:29.48241+00', '2025-10-10 16:21:36.786249+00', 69)
ON CONFLICT (id) DO NOTHING;

-- Update de permissies van 'title_section' (V1_42) naar de nieuwe tabelnaam
UPDATE permissions SET resource = 'title_section_content' WHERE resource = 'title_section';