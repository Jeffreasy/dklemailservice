-- Migratie: V1_31H__migrate_social_links.sql
-- Beschrijving: Migrate social_links data from Supabase
-- Versie: 1.31.0

-- Create social_links table
CREATE TABLE IF NOT EXISTS social_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    platform TEXT NOT NULL,
    url TEXT NOT NULL,
    bg_color_class TEXT,
    icon_color_class TEXT,
    order_number INTEGER,
    visible BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Insert social_links data
INSERT INTO social_links (id, platform, url, bg_color_class, icon_color_class, order_number, visible, created_at, updated_at) VALUES
('1de3ed0c-9bf3-4924-8d76-72045cc1c0ec', 'instagram', 'https://www.instagram.com/koninklijkeloop', null, null, 2, true, '2024-12-22 19:49:18.04807+00', '2024-12-22 19:49:18.04807+00'),
('29988672-0d98-4083-9908-18f6bbe34f3f', 'linkedin', 'https://www.linkedin.com/company/koninklijkeloop', null, null, 4, true, '2024-12-22 19:49:18.04807+00', '2024-12-22 19:49:18.04807+00'),
('dc917a65-6bb4-45cc-a1dc-890b1cdf1f5b', 'facebook', 'https://www.facebook.com/koninklijkeloop', null, null, 1, true, '2024-12-22 19:49:18.04807+00', '2024-12-22 19:49:18.04807+00'),
('f713d59e-e8a8-40d2-bb11-fc64f89b9ae9', 'youtube', 'https://www.youtube.com/@koninklijkeloop', null, null, 3, true, '2024-12-22 19:49:18.04807+00', '2024-12-22 19:49:18.04807+00');