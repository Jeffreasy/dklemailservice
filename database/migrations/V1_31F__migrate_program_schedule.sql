-- Migratie: V1_31F__migrate_program_schedule.sql
-- Beschrijving: Migrate program_schedule data from Supabase
-- Versie: 1.31.0

-- Create program_schedule table
CREATE TABLE IF NOT EXISTS program_schedule (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    time TEXT NOT NULL,
    event_description TEXT NOT NULL,
    category TEXT,
    icon_name TEXT,
    order_number INTEGER,
    visible BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8)
);

-- Insert program_schedule data
INSERT INTO program_schedule (id, time, event_description, category, icon_name, order_number, visible, created_at, updated_at, latitude, longitude) VALUES
('075095c7-925d-411e-bebc-a7fc96a3000a', '12:00u', 'Aanvang Deelnemers 10km bij het coördinatiepunt', 'Aanvang', 'aanvang', 50, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', null, null),
('1302ac11-1b7a-4738-a2dc-535990a09e69', '10:15u', 'Aanvang Deelnemers 15km bij het coördinatiepunt', 'Aanvang', 'aanvang', 10, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', null, null),
('1572e571-5f90-4cd9-b930-173d31df0124', '11:05u', 'Deelnemers 15km aanwezig startpunt (Kootwijk)', 'Aanwezig', 'aanwezig', 30, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', '52.18474064170016', '5.77074939643035'),
('16f07558-27cc-4e70-8d2f-4093d5e47009', '15:35u', 'START 2,5KM, Hervatting 6km, 10km en 15km', 'Start', 'start', 190, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', '52.22044761505996', '5.9288957497254176'),
('2eb4673c-6bde-465a-b4e8-27425bc32d54', '15:15u', 'Verwachte aankomst 15, 10, 6 km lopers bij rustpunt (Berg & Bos - 15 min pauze)', 'Rustpunt', 'rustpunt', 180, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', null, null),
('425c461a-cee7-480c-a399-7d469cc7efbe', '12:50u', 'Deelnemers 10km aanwezig bij het startpunt (Halte Assel)', 'Aanwezig', 'aanwezig', 80, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', '52.20071361748442', '5.836023236413426'),
('460625b7-ec90-446d-a895-2ddaefb98335', '14:15u', 'START 6KM, Hervatting 10km en 15km', 'Start', 'start', 140, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', '52.219164379355604', '5.872559210794147'),
('4f201162-2524-4b03-b0ea-0fb5ccaf29c3', '15:00u', 'Vertrek deelnemers 2,5 km met de pendelbussen naar het startpunt 2,5km', 'Vertrek', 'vertrek', 160, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', null, null),
('5af9c5b7-9b73-44aa-9153-61e2277b9233', '17:00u – 18:00u', 'INHULDIGINGSFEEST', 'Feest', 'feest', 220, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', null, null),
('652b1d86-b062-4b63-bbf0-5294c71979d0', '12:45u', 'Verwachte aankomst 15 km lopers bij rustpunt (Halte Assel - 15 min pauze)', 'Rustpunt', 'rustpunt', 70, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', null, null),
('65dd7eb0-4de3-4d54-99b7-2f278a08ece4', '12:30u', 'Vertrek deelnemers 10km met de pendelbussen naar het startpunt 10km', 'Vertrek', 'vertrek', 60, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', null, null),
('6730216e-9a4c-4df5-aaca-119da8595eef', '10:45u', 'Vertrek pendelbussen naar startpunt 15km', 'Vertrek', 'vertrek', 20, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', null, null),
('7ae60214-dc2d-4b23-a4d1-0993fbb56e46', '14:00u', 'Deelnemers 6km aanwezig bij het startpunt (Hoog Soeren)', 'Aanwezig', 'aanwezig', 130, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', '52.219164379355604', '5.872559210794147'),
('91666bd3-32fd-4054-ab31-4d9f7532c0ce', '14:30u', 'Aanvang Deelnemers 2,5km bij het coördinatiepunt', 'Aanvang', 'aanvang', 150, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', null, null),
('9564d363-17d5-45b4-b868-45d165a82c72', '16:10u – 16:30u', 'FINISH', 'Finish', 'finish', 210, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', null, null),
('9ff53726-0e9a-48de-b4f9-35cfa7666756', '15:55u', 'Aankomst bij De Naald / START INHULDIGINGSLOOP', 'Aankomst', 'aankomst', 200, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', null, null),
('b9c046ea-de08-48e5-995e-9d4a98d76b6e', '14:00u', 'Verwachte aankomst 15, 10 km lopers bij rustpunt (Hoog Soeren - 15 min pauze)', 'Rustpunt', 'rustpunt', 120, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', null, null),
('bc30a0ea-0f22-443a-8124-0cd52e10a2b3', '11:15u', 'START 15KM', 'Start', 'start', 40, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', '52.18474064170016', '5.77074939643035'),
('cf1514db-2e06-45ab-891b-76736eb308d3', '13:00u', 'START 10KM, Hervatting 15km', 'Start', 'start', 90, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', '52.20071361748442', '5.836023236413426'),
('de0cd4c1-2fe0-4f13-9b94-00c03ce38527', '15:05u', 'Deelnemers 2,5km aanwezig bij het startpunt (Berg & Bos)', 'Aanwezig', 'aanwezig', 170, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', '52.22044761505996', '5.9288957497254176'),
('fd80912e-f112-4bee-a2dd-570c4cf88c89', '13:15u', 'Aanvang Deelnemers 6km bij het coördinatiepunt', 'Aanvang', 'aanvang', 100, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', null, null),
('fec9b64c-4cdd-4fc1-a507-9c222fdeb958', '13:45u', 'Vertrek deelnemers 6 km met de pendelbussen naar het startpunt 6km', 'Vertrek', 'vertrek', 110, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', null, null)
ON CONFLICT (id) DO NOTHING;