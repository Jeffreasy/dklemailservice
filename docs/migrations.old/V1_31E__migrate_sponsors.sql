-- Migratie: V1_31E__migrate_sponsors.sql
-- Beschrijving: Migrate sponsors data from Supabase
-- Versie: 1.31.0

-- Create sponsors table
CREATE TABLE IF NOT EXISTS sponsors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    logo_url TEXT,
    website_url TEXT,
    order_number INTEGER,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    visible BOOLEAN NOT NULL DEFAULT TRUE
);

-- Insert sponsors data
INSERT INTO sponsors (id, name, description, logo_url, website_url, order_number, is_active, created_at, updated_at, visible) VALUES
('484576a1-2a60-4201-b582-e1f3754ab12a', '3x3 Anders', '3x3 Anders is een zorgbemiddelingsbureau gespecialiseerd in het matchen van zorgaanbieders met gekwalificeerde zorgprofessionals.', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1733166671/3x3anderslogo_itwm3g.webp', 'https://3x3anders.nl/', 4, true, '2024-11-29 09:49:35+00', '2025-04-30 13:07:58.478643+00', true),
('6408a640-cca1-4aaa-b845-04888f62ccec', 'Sterk In Vloeren', 'De website van Sterk In Vloeren biedt een uitgebreid assortiment aan vloeren, waaronder laminaat, PVC-vloeren en tapijt. Ze benadrukken heldere afspraken en hanteren all-in prijzen.', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1733166669/SterkinVloerenLOGO_zrdofb.webp', 'https://sterkinvloeren.nl/', 1, true, '2024-11-29 09:49:35+00', '2025-04-30 13:07:58.478643+00', true),
('6acd1b1e-8fed-4c8b-89cb-85eee9053536', 'Beeldpakker', 'Johan Groot Jebbink, een fotograaf met meer dan tien jaar ervaring, gespecialiseerd in portretfotografie. Actief in Ermelo en internationaal.', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1733166670/BeeldpakkerLogo_wijjmq.webp', 'https://beeldpakker.nl/', 2, true, '2024-11-29 09:49:35+00', '2025-04-30 13:07:58.478643+00', true),
('88cec1c1-3f08-4a6f-9623-c71605fe35b5', 'Bas Visual Story Telling', 'BAS Visual Storytelling, heeft een passie voor content en verhalen. Mijn hobby is uitgegroeid tot een eigen onderneming in het vastleggen van verhalen. Bij BAS Visual Storytelling laten we verhalen niet verstoffen op de plank, maar brengen ze tot leven! Waar ik ga of sta, mijn camera''s gaan met mij mee, leg de mooiste beelden haarscherp vast en breng jouw verhaal tot leven. Dus vertel eens, ''wat is jouw verhaal?''

', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1746017513/krqjbwwerv9hs6hyrhcy.png', 'https://basvisualstorytelling.nl/', 5, true, '2025-04-30 12:51:53.997413+00', '2025-05-01 09:30:16.725422+00', true),
('caa59f1f-65b4-442f-84d2-22cc52212dea', 'Mojo Dojo', 'Mojo Dojo is een veelzijdige studio in Rotterdam die diensten aanbiedt voor creatieve producties, waaronder muziekopnames, podcasts en livestreams.', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1733166669/LogoLayout_1_iphclc.webp', 'https://mojodojo.studio/', 3, true, '2024-11-29 09:49:35+00', '2025-04-30 13:07:58.478643+00', true)
ON CONFLICT (id) DO NOTHING;