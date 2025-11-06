-- GECONSOLIDEERDE V26 - NORMALISEER ROLLEN EN AFSTANDEN
--
-- DOEL: Vervangt "magic strings" in de 'aanmeldingen' tabel door
--       robuuste Foreign Key relaties.
--
-- WAAROM:
-- 1. Data Integriteit: Voorkomt typfouten (bv. '6 KM' vs '6km').
-- 2. Beheerbaarheid: Je beheert afstanden en rollen op één centrale plek.
-- 3. Efficiëntie: Database kan efficiënter werken met ID's dan met strings.
--
-- IMPACT: BREAKING CHANGE!
-- Je MOET je GORM-modellen in de Go-code aanpassen.
-- 'aanmeldingen.rol' (string) en 'aanmeldingen.afstand' (string)
-- worden vervangen door 'participant_role_name' (FK) en 'distance_route' (FK).
--

-- =====================================================
-- 1. NORMALISEER DE DEELNEMER-ROLLEN (aanmeldingen.rol)
-- =====================================================

-- Stap 1.1: Maak de "lookup-tabel" voor rollen.
CREATE TABLE IF NOT EXISTS participant_roles (
    name TEXT PRIMARY KEY NOT NULL,
    description TEXT
);
COMMENT ON TABLE participant_roles IS 'Bron van waarheid voor deelnemer-rollen (Deelnemer, Vrijwilliger, etc) (V26)';

-- Stap 1.2: Vul de tabel met de rollen die je al gebruikt.
-- (Deze zijn afkomstig uit je RBAC seeding in V7)
INSERT INTO participant_roles (name, description) VALUES
('Deelnemer', 'Een standaard deelnemer aan het evenement.'),
('Begeleider', 'Een begeleider van een of meerdere deelnemers.'),
('Vrijwilliger', 'Een vrijwilliger die helpt bij het evenement.'),
('Sponsor', 'Een sponsor of partner (indien deze zich kunnen aanmelden).')
ON CONFLICT (name) DO NOTHING;

-- Stap 1.3: Voeg de nieuwe Foreign Key kolom toe aan 'aanmeldingen'.
ALTER TABLE aanmeldingen
    ADD COLUMN IF NOT EXISTS participant_role_name TEXT;

-- Stap 1.4: Migreer de oude string-data naar de nieuwe kolom.
-- Dit normaliseert ook data (bv. alles met 'hoofdletter').
UPDATE aanmeldingen
SET participant_role_name = 
    CASE
        WHEN LOWER(TRIM(rol)) = 'deelnemer' THEN 'Deelnemer'
        WHEN LOWER(TRIM(rol)) = 'begeleider' THEN 'Begeleider'
        WHEN LOWER(TRIM(rol)) = 'vrijwilliger' THEN 'Vrijwilliger'
        WHEN LOWER(TRIM(rol)) = 'sponsor' THEN 'Sponsor'
        ELSE NULL -- Of zet een default, bv. 'Deelnemer'
    END
WHERE participant_role_name IS NULL AND rol IS NOT NULL;

-- Stap 1.5: Voeg de Foreign Key constraint toe.
-- Dit faalt als je 'rol' waarden had die NIET in Stap 1.2 zijn ingevoegd.
ALTER TABLE aanmeldingen
    DROP CONSTRAINT IF EXISTS fk_aanmeldingen_participant_role,
    ADD CONSTRAINT fk_aanmeldingen_participant_role
    FOREIGN KEY (participant_role_name) REFERENCES participant_roles(name)
    ON UPDATE CASCADE ON DELETE SET NULL;

-- =====================================================
-- 2. NORMALISEER DE AFSTANDEN (aanmeldingen.afstand)
-- =====================================================

-- Stap 2.1: Hergebruik en hernoem de 'route_funds' (V17) tabel.
-- Dit is al je "lookup-tabel"!
ALTER TABLE IF EXISTS route_funds RENAME TO distances;

-- Stap 2.2: Pas de kolomnaam aan voor de duidelijkheid.
ALTER TABLE distances RENAME COLUMN IF EXISTS amount TO fund_amount;

-- Stap 2.3: Zorg dat de 'route' (de naam, bv '6 KM') de Primary Key is.
-- Dit maakt de migratie makkelijker.
ALTER TABLE distances
    DROP CONSTRAINT IF EXISTS route_funds_pkey,
    DROP CONSTRAINT IF EXISTS route_funds_route_key,
    ADD PRIMARY KEY (route);
COMMENT ON TABLE distances IS 'Bron van waarheid voor alle afstanden en hun fondsenwerving (V17, V26)';

-- Stap 2.4: Voeg de nieuwe Foreign Key kolom toe aan 'aanmeldingen'.
ALTER TABLE aanmeldingen
    ADD COLUMN IF NOT EXISTS distance_route TEXT;

-- Stap 2.5: Migreer de oude 'afstand' data.
-- We gebruiken UPPER() en TRIM() om data op te schonen (bv. ' 6 km ' -> '6 KM')
-- aangenomen dat je routes in 'distances' ook in hoofdletters staan.
UPDATE aanmeldingen
SET distance_route = UPPER(TRIM(afstand))
WHERE distance_route IS NULL AND afstand IS NOT NULL AND TRIM(afstand) != '';

-- Stap 2.6: Voeg de Foreign Key constraint toe.
-- Dit faalt als 'aanmeldingen' afstanden bevat die NIET in de 'distances' tabel staan.
ALTER TABLE aanmeldingen
    DROP CONSTRAINT IF EXISTS fk_aanmeldingen_distance,
    ADD CONSTRAINT fk_aanmeldingen_distance
    FOREIGN KEY (distance_route) REFERENCES distances(route)
    ON UPDATE CASCADE ON DELETE SET NULL; -- SET NULL als een afstand wordt verwijderd

-- =====================================================
-- 3. OPRAKELEN (Optioneel, pas uitvoeren na Go-code update)
-- =====================================================

/*
-- Nadat je GORM-modellen zijn aangepast en alles werkt:
ALTER TABLE aanmeldingen DROP COLUMN IF EXISTS rol;
ALTER TABLE aanmeldingen DROP COLUMN IF EXISTS afstand;
*/

-- =====================================================
-- 4. LOG INSTRUCTIES
-- =====================================================
DO $$
BEGIN
    RAISE NOTICE '[V26] Normalisatie van rollen en afstanden is voorbereid.';
    RAISE NOTICE '=== BREAKING CHANGE ===';
    RAISE NOTICE '1. (Go Code) Pas je `Aanmelding` GORM-model aan.';
    RAISE NOTICE '   - Vervang `Rol string` door `ParticipantRoleName string` (of `ParticipantRole ParticipantRole `gorm:"foreignKey:ParticipantRoleName"`)';
    RAISE NOTICE '   - Vervang `Afstand string` door `DistanceRoute string` (of `Distance Distance `gorm:"foreignKey:DistanceRoute"`)';
    RAISE NOTICE '2. (Optioneel) Voer Stap 3 (DROP COLUMNs) uit om de database definitief op te schonen.';
END $$;