-- Migratie: V1_43__replace_title_sections_with_content.sql
-- Beschrijving: Replace title_sections table with title_section_content table
-- Versie: 1.43.0

-- Drop the old title_sections table
DROP TABLE IF EXISTS title_sections;

-- Create the new title_section_content table
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

-- Insert title_section_content data from the provided SQL
INSERT INTO title_section_content (id, event_title, event_subtitle, image_url, image_alt, detail_1_title, detail_1_description, detail_2_title, detail_2_description, detail_3_title, detail_3_description, created_at, updated_at, participant_count) VALUES
('1', 'De Koninklijke Loop (DKL) 2025', 'Op de koninklijke weg in Apeldoorn kunnen mensen met een beperking samen wandelen tijdens dit unieke, rolstoelvriendelijke sponsorloop (DKL), samen met hun verwanten, vrijwilligers of begeleiders.', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1760112848/Wij_gaan_17_mei_lopen_voor_hen_3_zllxno_zoqd7z.webp', 'Promotiebanner De Koninklijke Loop (DKL) 2025: Wij gaan 17 mei lopen voor hen', '17 mei 2025', 'Starttijden variëren per afstand. Zie programma.', 'Voor iedereen', 'wandelaars met of zonder beperking (rolstoelvriendelijk).', 'Lopen voor een goed doel', 'Steun het goede doel via dit unieke wandelevenement.', '2025-04-16 01:31:29.48241+00', '2025-10-10 16:21:36.786249+00', 69)
ON CONFLICT (id) DO NOTHING;

-- Update permissions to use title_section_content instead of title_sections
UPDATE permissions SET resource = 'title_section_content' WHERE resource = 'title_sections';