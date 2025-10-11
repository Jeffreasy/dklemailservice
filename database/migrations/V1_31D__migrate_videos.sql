-- Migratie: V1_31D__migrate_videos.sql
-- Beschrijving: Migrate videos data from Supabase
-- Versie: 1.31.0

-- Create videos table
CREATE TABLE IF NOT EXISTS videos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    video_id TEXT NOT NULL,
    url TEXT NOT NULL,
    title TEXT,
    description TEXT,
    thumbnail_url TEXT,
    visible BOOLEAN NOT NULL DEFAULT TRUE,
    order_number INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Insert videos data
INSERT INTO videos (id, video_id, url, title, description, thumbnail_url, visible, order_number, created_at, updated_at) VALUES
('14ee164b-50e3-4f59-a3e9-3a54312af9cd', 'q9ngqu', 'https://streamable.com/e/q9ngqu', 'De koninklijkeloop!', 'Preview!!!', null, true, 1, '2025-03-28 21:45:53+00', '2025-04-21 10:49:03.744+00'),
('18d951d2-f5d1-4b6e-95af-a72ac5ff18ff', 'x8zj4k', 'https://streamable.com/e/x8zj4k', 'Promotie De Koninklijke Loop: Flyers verspreiden', 'Bekijk hoe vrijwilligers flyers uitdelen om mensen uit te nodigen voor het DKL wandelevenement.', null, true, 4, '2025-03-03 20:25:54.251108+00', '2025-03-03 20:25:54.251108+00'),
('87502f84-91db-419f-9766-4071ade3e94f', 'tt6k80', 'https://streamable.com/e/0o2qf9', 'Highlights Koninklijke Loop 2024 (Wandelevenement Apeldoorn)', 'Herbeleef de mooiste momenten en de sfeer van De Koninklijke Loop 2024 in deze highlight video.', null, true, 2, '2024-12-22 20:23:06.654419+00', '2024-12-26 23:25:42.699+00'),
('99bbe55b-32ef-46ab-bb59-860ce92f1d58', 'cvfrpi', 'https://streamable.com/e/cvfrpi', 'De spannende start van de Koninklijke Loop 2024', 'Bekijk de start van de deelnemers aan de sponsorloop De Koninklijke Loop 2024.', null, true, 3, '2024-12-22 20:23:06.654419+00', '2024-12-26 23:25:43.979+00'),
('ac987839-edc1-468f-9080-064f894b3e5d', 'tt6k80', 'https://streamable.com/e/tt6k80', 'Koninklijke Loop 2024 - Hoofdevenement', 'Een sfeerimpressie van het hoofdevenement van de Koninklijke Loop 2024, met deelnemers, vrijwilligers en muziek.', null, true, 5, '2024-12-22 20:23:06.654419+00', '2024-12-26 23:25:41.341+00');