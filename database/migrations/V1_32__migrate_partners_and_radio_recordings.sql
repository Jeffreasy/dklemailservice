-- V1_32__migrate_partners_and_radio_recordings.sql
-- Migrate partners and radio_recordings data from Supabase

-- Create partners table
CREATE TABLE IF NOT EXISTS partners (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    logo TEXT,
    website TEXT,
    tier TEXT,
    since DATE,
    visible BOOLEAN NOT NULL DEFAULT TRUE,
    order_number INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create radio_recordings table
CREATE TABLE IF NOT EXISTS radio_recordings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    description TEXT,
    date TEXT,
    audio_url TEXT,
    thumbnail_url TEXT,
    visible BOOLEAN NOT NULL DEFAULT TRUE,
    order_number INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Insert partners data (idempotent)
INSERT INTO "public"."partners" ("id", "name", "description", "logo", "website", "tier", "since", "visible", "order_number", "created_at", "updated_at") VALUES
('26f11d04-c0da-4755-b2e1-5fcb5e9887d4', 'Accress', 'Accres beheert in Apeldoorn ruim 60 locaties, waaronder sporthallen, wijkcentra, zwembaden, kinderboerderijen en een stadspark.

Zij helpen ond bij het halen van ons doel!', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1744388421/accres_logo_ochsmg.jpg', 'https://www.accres.nl/', 'bronze', '2025-04-01', 'true', '5', '2025-04-11 16:21:40+00', '2025-04-11 16:45:10.227784+00'),
('510d8e4d-6a7d-4ab3-b311-17b32df7f01f', 'Apeldoorn', 'Apeldoorn ondersteunt ons in ons doel.', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1734895194/nw1qxouzupzsshckzkab.png', 'https://www.apeldoorn.nl/', 'bronze', '2024-12-22', 'true', '0', '2024-12-22 19:19:55.414207+00', '2025-02-26 21:10:59.040907+00'),
('5e8b8390-e637-4c6d-80d8-25ff4155d5a9', 'Sheeren Loo', 'Samen met bewoners van SheerenLoo wordt deze loop georganiseerd.', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1734894570/mtvucaouruenat2cllsi.png', 'https://www.sheerenloo.nl/', 'silver', '2024-12-22', 'true', '0', '2024-12-22 19:09:31.099813+00', '2025-01-07 14:11:39.051238+00'),
('9c311d0f-4db7-4da9-9d87-a836e48078cd', 'Liliane Fonds', 'Samen maken we ons sterk', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1734893709/qsygajx2tdxxbqbfyurr.png', 'https://www.lilianefonds.nl/', 'bronze', '2024-12-22', 'true', '0', '2024-12-22 18:55:09.823207+00', '2025-01-08 19:43:58.464461+00'),
('eda49448-3db9-4db5-9d01-a7331879f5b5', 'De Grote Kerk', 'De grote kerk ondersteund ons al vanaf het begin. Hier is het allemaal begonnen.', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1734895146/ri4vclttn4nn2wh53wj0.jpg', 'https://www.grotekerkapeldoorn.nl/', 'gold', '2024-12-22', 'true', '0', '2024-12-22 19:19:07.18416+00', '2025-04-11 16:46:10.191766+00')
ON CONFLICT (id) DO NOTHING;

-- Insert radio_recordings data (idempotent)
INSERT INTO "public"."radio_recordings" ("id", "title", "description", "date", "audio_url", "thumbnail_url", "visible", "order_number", "created_at", "updated_at") VALUES
('a6e73425-6af2-4bbd-84f8-67b16d195f99', 'De koninklijke Loop 2025 uitzending!', 'Luister naar het live radioverslag van De Koninklijke Loop 2025, uitgezonden op RTV Apeldoorn, met interviews met de organisatie!

', '14 mei 2025', 'https://res.cloudinary.com/dgfuv7wif/video/upload/v1747733438/DKLRTV2025_dpdydc.wav', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png', 'true', '1', '2025-05-20 09:32:01+00', '2025-05-20 09:33:27.540519+00'),
('c2d84ca9-97a9-4990-9ee4-1fe1718a8c5b', 'Radioverslag Koninklijke Loop 2024 (RTV Apeldoorn)', 'Luister naar het live radioverslag van De Koninklijke Loop 2024, uitgezonden op RTV Apeldoorn, met interviews en sfeerimpressies.', '15 mei 2024', 'https://res.cloudinary.com/dgfuv7wif/video/upload/v1714042357/matinee_1_nbm0ph.wav', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png', 'true', '2', '2025-04-08 19:50:35.125804+00', '2025-05-20 09:32:30.549188+00')
ON CONFLICT (id) DO NOTHING;