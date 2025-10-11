-- Migratie: V1_31B__migrate_albums.sql
-- Beschrijving: Migrate albums data from Supabase
-- Versie: 1.31.0

-- Create albums table
CREATE TABLE IF NOT EXISTS albums (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    description TEXT,
    cover_photo_id UUID REFERENCES photos(id),
    visible BOOLEAN NOT NULL DEFAULT TRUE,
    order_number INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Insert albums data
INSERT INTO albums (id, title, description, cover_photo_id, visible, order_number, created_at, updated_at) VALUES
('72831c18-4c6c-4dc8-9a2a-ace696e8996b', 'Voorbereidingen ', 'Voor werk', 'dbb68d91-3f6f-46c3-b751-a13f34269b79', true, 3, '2025-02-14 14:35:32.369+00', '2025-02-14 14:35:32.37+00'),
('ce8df963-f118-4296-9f5c-e33308dc7bfa', 'DKL-2024', 'DKL 2024', '8a4d5c20-ea73-4336-8c6b-af7a197ef7c2', true, 2, '2024-12-26 15:17:54.268+00', '2025-02-17 13:32:50.793+00'),
('d51cff45-b958-4370-a983-51e650ffa43e', 'DKL 2025', 'De koninklijke Loop 2025!', '08362e92-340a-432a-b306-153ad27ee686', true, 1, '2025-05-17 20:10:00+00', '2025-05-17 20:10:06.643082+00')
ON CONFLICT (id) DO NOTHING;