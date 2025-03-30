-- V1_14__add_more_registrations.sql
-- Beschrijving: Voegt nieuwe aanmeldingen toe van maart 2025
-- Versie: 1.14

-- Alleen uitvoeren als de aanmeldingen nog niet bestaan
DO $$
BEGIN
    -- Nieuwe aanmeldingen alleen toevoegen als ze nog niet bestaan
    IF NOT EXISTS (SELECT 1 FROM aanmeldingen WHERE id = '1ca80f61-f5c1-431f-b224-e6557150b65b') THEN
        INSERT INTO "public"."aanmeldingen" 
        ("id", "created_at", "updated_at", "naam", "email", "telefoon", "rol", "afstand", "ondersteuning", "bijzonderheden", "terms", "email_verzonden", "email_verzonden_op", "test_mode", "status")
        VALUES
        ('1ca80f61-f5c1-431f-b224-e6557150b65b', '2025-03-30 08:07:45.334762+00', '2025-03-30 08:07:45.334762+00', 'Han van Doornik', 'LaanvanGS.26@sheerenloo.nl', null, 'Deelnemer', '2.5 KM', 'Ja', 'Ik wil wel graag begeleiding ', 'true', 'false', null, FALSE, 'nieuw'),
        ('2499686e-2e62-4827-9079-78b468cb26c9', '2025-03-29 10:07:12.393466+00', '2025-03-29 10:07:56.169434+00', 'Bertram tijsma', 'Klaskehiddes@gmail.com', null, 'Deelnemer', '15 KM', 'Nee', '', 'true', 'true', '2025-03-29 10:07:57.838+00', FALSE, 'verwerkt'),
        ('9f75b1df-4c72-4e36-9901-0f74cc26574f', '2025-03-29 10:05:04.276906+00', '2025-03-29 10:07:55.533152+00', 'Klaske van de glind', 'Klaskehiddes@gmail.com', null, 'Deelnemer', '15 KM', 'Nee', '', 'true', 'true', '2025-03-29 10:07:57.215+00', FALSE, 'verwerkt'),
        ('db3ec762-dd54-4ba7-98ab-981235cc316a', '2025-03-26 16:26:53.777384+00', '2025-03-26 16:30:58.285233+00', 'Mila Veenendaal', 'gaminggirlayla@gmail.com', null, 'Deelnemer', '10 KM', 'Nee', '', 'true', 'true', '2025-03-26 16:30:57.269+00', FALSE, 'verwerkt'),
        ('d17c16c6-c423-43de-a876-d40326b62d9e', '2025-03-26 16:25:43.756211+00', '2025-03-26 16:30:57.680904+00', 'Ayla Toprak', 'gamergirlayla@gmail.com', null, 'Deelnemer', '10 KM', 'Nee', '', 'true', 'true', '2025-03-26 16:30:56.685+00', FALSE, 'verwerkt'),
        ('917206a7-a28d-4bad-8b41-ed127eab743a', '2025-03-26 12:27:20.236848+00', '2025-03-26 16:30:57.073558+00', 'A. Bistolfi', 'nedarg@icloud.com', null, 'Deelnemer', '15 KM', 'Nee', '', 'true', 'true', '2025-03-26 16:30:56.062+00', FALSE, 'verwerkt');
        
        -- Logbericht voor monitoring
        RAISE NOTICE 'Nieuwe aanmeldingen toegevoegd (6 records)';
    ELSE
        RAISE NOTICE 'Aanmeldingen bestaan al, geen nieuwe records toegevoegd';
    END IF;
END;
$$;

-- Registreer de migratie
INSERT INTO migraties (versie, naam, toegepast) 
VALUES ('1.14', 'Toevoegen van extra aanmeldingen maart 2025', CURRENT_TIMESTAMP)
ON CONFLICT (versie) DO NOTHING; 