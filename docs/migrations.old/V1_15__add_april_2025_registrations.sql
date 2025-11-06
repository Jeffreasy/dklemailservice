-- V1_15__add_april_2025_registrations.sql
-- Beschrijving: Voegt nieuwe aanmeldingen toe van april 2025
-- Versie: 1.15

DO $$
BEGIN
    -- Controleer of de eerste nieuwe aanmelding al bestaat
    IF NOT EXISTS (SELECT 1 FROM aanmeldingen WHERE id = '61a3f823-82f6-4107-a9a9-b6a18e6b12c3') THEN
        INSERT INTO "public"."aanmeldingen"
        ("id", "created_at", "updated_at", "naam", "email", "telefoon", "rol", "afstand", "ondersteuning", "bijzonderheden", "terms", "email_verzonden", "email_verzonden_op", "test_mode", "status")
        VALUES
        ('61a3f823-82f6-4107-a9a9-b6a18e6b12c3', '2025-04-14 20:03:47.361483+00', '2025-04-15 13:23:43.293548+00', 'Henk Rekers ', 'h.rekers59@kpnmail.nl', null, 'Deelnemer', '6 KM', 'Nee', '', 'true', 'true', '2025-04-15 13:23:40.051+00', FALSE, 'verwerkt'),
        ('8a5c2c59-ebca-411e-9471-a38be47e1192', '2025-04-14 20:00:47.479484+00', '2025-04-15 13:23:42.628968+00', 'Hilde Rekers ', 'h.rekers59@kpnmail.nl', null, 'Deelnemer', '6 KM', 'Nee', '', 'true', 'true', '2025-04-15 13:23:39.203+00', FALSE, 'verwerkt'),
        ('3c420d85-5f78-4645-88dd-c9e90eb8b6f7', '2025-04-12 10:21:37.470747+00', '2025-04-15 15:38:46.288552+00', 'Henk Rekers ', 'h.rekers1959@kpnmail.nl', null, 'Deelnemer', '6 KM', 'Nee', 'email mislukt', 'true', 'false', '2025-04-12 16:52:21.796+00', FALSE, 'nieuw'),
        ('721d297d-d670-46b1-a581-bcc095565bbd', '2025-04-12 10:20:25.399025+00', '2025-04-15 15:38:52.374917+00', 'Hilde Rekers ', 'h.rekers1959@kpnmail.nl', null, 'Deelnemer', '6 KM', 'Nee', 'email mislukt', 'true', 'false', '2025-04-12 16:52:20.444+00', FALSE, 'nieuw'),
        ('02f9605b-8b26-4461-9ca7-ed7ebbd12311', '2025-04-12 07:09:19.884142+00', '2025-04-12 16:52:21.459029+00', 'Theun ', 'diesbosje@hotmail.com', null, 'Deelnemer', '6 KM', 'Nee', '', 'true', 'true', '2025-04-12 16:52:19.056+00', FALSE, 'verwerkt'),
        ('282e6ec6-2c97-4b46-a992-273c826c1f91', '2025-04-12 07:08:26.565956+00', '2025-04-12 16:52:20.449671+00', 'Albert ', 'diesbosje@hotmail.com', null, 'Deelnemer', '6 KM', 'Nee', '', 'true', 'true', '2025-04-12 16:52:18.027+00', FALSE, 'verwerkt'),
        ('d10065f9-229d-4a1d-9178-324bdffa57c0', '2025-04-12 07:06:59.872906+00', '2025-04-12 16:52:19.826576+00', 'Diesmer ', 'diesbosje@hotmail.com', '0613429612', 'Begeleider', '6 KM', 'Nee', '', 'true', 'true', '2025-04-12 16:52:17.255+00', FALSE, 'verwerkt'),
        ('613770dd-4733-4b54-964d-c3753cfebfd7', '2025-04-06 18:48:14.993723+00', '2025-04-10 11:15:45.930907+00', 'Sylvia Dijkstra', 'sylvia.dijkstra@sheerenloo.nl', '0683081728 ', 'Begeleider', '2.5 KM', 'Nee', '', 'true', 'true', '2025-04-10 11:15:46.282+00', FALSE, 'verwerkt'),
        ('2aaffc78-4dca-4b06-9a6a-3799cbfbe67b', '2025-03-31 11:57:23.536871+00', '2025-04-01 17:11:03.413137+00', 'Noa hiddes', 'Klaskehiddes@gmail.com', null, 'Deelnemer', '2.5 KM', 'Nee', '', 'true', 'true', '2025-04-01 17:11:04.313+00', FALSE, 'verwerkt'),
        ('200f862d-4d80-4ccd-8f06-bf795be727fb', '2025-03-31 11:56:37.582458+00', '2025-04-01 17:11:04.334788+00', 'Anneke van de Glind ', 'Klaskehiddes@gmail.com', null, 'Deelnemer', '2.5 KM', 'Nee', '', 'true', 'true', '2025-04-01 17:11:05.257+00', FALSE, 'verwerkt');

        -- Logbericht voor monitoring
        RAISE NOTICE 'Nieuwe aanmeldingen toegevoegd (10 records)';
    ELSE
        RAISE NOTICE 'Aanmeldingen vanaf april 2025 bestaan al, geen nieuwe records toegevoegd';
    END IF;
END;
$$;

-- Registreer de migratie
INSERT INTO migraties (versie, naam, toegepast)
VALUES ('1.15', 'Toevoegen van aanmeldingen april 2025', CURRENT_TIMESTAMP)
ON CONFLICT (versie) DO NOTHING; 