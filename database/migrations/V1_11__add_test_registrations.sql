-- V1_11__add_test_registrations.sql
-- Voegt testgegevens toe voor aanmeldingen

-- Alleen uitvoeren als de aanmeldingen nog niet bestaan
DO $$
BEGIN
    -- Aanmeldingen toevoegen als ze niet bestaan
    IF NOT EXISTS (SELECT 1 FROM aanmeldingen WHERE id = '3e62d5d3-070d-47b1-a1ef-30665f982789') THEN
        INSERT INTO "public"."aanmeldingen" ("id", "created_at", "updated_at", "naam", "email", "telefoon", "rol", "afstand", "ondersteuning", "bijzonderheden", "terms", "email_verzonden", "email_verzonden_op") 
        VALUES ('3e62d5d3-070d-47b1-a1ef-30665f982789', '2025-03-23 17:06:26.132297+00', '2025-03-23 17:06:26.132297+00', 'TGTest', 'laventejeffrey@gmail.com', '06123456789', 'Begeleider', '15 KM', 'Anders', 'Telegram Test bericht - officiele weg', 'true', 'false', null);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM aanmeldingen WHERE id = '51205069-a20a-4231-bcb9-cb2d6fd042c8') THEN
        INSERT INTO "public"."aanmeldingen" ("id", "created_at", "updated_at", "naam", "email", "telefoon", "rol", "afstand", "ondersteuning", "bijzonderheden", "terms", "email_verzonden", "email_verzonden_op") 
        VALUES ('51205069-a20a-4231-bcb9-cb2d6fd042c8', '2025-03-23 14:08:53.328628+00', '2025-03-23 14:08:53.328628+00', 'Manuela van Zwam', 'rik.van-harxen@sheerenloo.nl', null, 'Deelnemer', '2.5 KM', 'Ja', 'Vaste begeleider die meeloopt', 'true', 'false', null);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM aanmeldingen WHERE id = '275490c0-1021-4bf4-9005-7df9884b0fe6') THEN
        INSERT INTO "public"."aanmeldingen" ("id", "created_at", "updated_at", "naam", "email", "telefoon", "rol", "afstand", "ondersteuning", "bijzonderheden", "terms", "email_verzonden", "email_verzonden_op") 
        VALUES ('275490c0-1021-4bf4-9005-7df9884b0fe6', '2025-03-22 16:43:03.19496+00', '2025-03-22 16:43:03.19496+00', 'Bas heijenk ', 'basheijenk96@gmail.com', null, 'Deelnemer', '2.5 KM', 'Nee', '', 'true', 'false', null);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM aanmeldingen WHERE id = 'b5e67c64-6bfa-46d4-ae40-98b093c8b720') THEN
        INSERT INTO "public"."aanmeldingen" ("id", "created_at", "updated_at", "naam", "email", "telefoon", "rol", "afstand", "ondersteuning", "bijzonderheden", "terms", "email_verzonden", "email_verzonden_op") 
        VALUES ('b5e67c64-6bfa-46d4-ae40-98b093c8b720', '2025-03-17 19:55:07.379647+00', '2025-03-21 06:47:40.512037+00', 'Salih', 'topraks@gmail.com', null, 'Deelnemer', '2.5 KM', 'Nee', '', 'true', 'true', '2025-03-21 06:47:38.608+00');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM aanmeldingen WHERE id = 'b2fd3412-8368-409f-8029-b2cdd581ade1') THEN
        INSERT INTO "public"."aanmeldingen" ("id", "created_at", "updated_at", "naam", "email", "telefoon", "rol", "afstand", "ondersteuning", "bijzonderheden", "terms", "email_verzonden", "email_verzonden_op") 
        VALUES ('b2fd3412-8368-409f-8029-b2cdd581ade1', '2025-03-11 11:28:05.483427+00', '2025-03-21 06:47:41.806005+00', 'Manuela van zwam', 'benjaminlaan.64a@sheerenloo.nl', null, 'Deelnemer', '2.5 KM', 'Nee', '', 'true', 'true', '2025-03-21 06:47:39.928+00');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM aanmeldingen WHERE id = '391f63c5-f034-466e-8a1f-ba9d06ed1192') THEN
        INSERT INTO "public"."aanmeldingen" ("id", "created_at", "updated_at", "naam", "email", "telefoon", "rol", "afstand", "ondersteuning", "bijzonderheden", "terms", "email_verzonden", "email_verzonden_op") 
        VALUES ('391f63c5-f034-466e-8a1f-ba9d06ed1192', '2025-03-09 16:52:09.564437+00', '2025-03-21 06:47:42.987352+00', 'Joyce Thielen', 'Joyce.thielen@sheerenloo.nl', '', 'Begeleider', '6 KM', 'Nee', '', 'true', 'true', '2025-03-21 06:47:41.071+00');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM aanmeldingen WHERE id = '391f2579-d7cb-4ef3-afbe-14dc4115c519') THEN
        INSERT INTO "public"."aanmeldingen" ("id", "created_at", "updated_at", "naam", "email", "telefoon", "rol", "afstand", "ondersteuning", "bijzonderheden", "terms", "email_verzonden", "email_verzonden_op") 
        VALUES ('391f2579-d7cb-4ef3-afbe-14dc4115c519', '2025-03-08 10:28:37.053379+00', '2025-03-08 10:28:38.042991+00', 'Dick van Norden', 'Enckerkamp.27@sheerenloo.nl', null, 'Deelnemer', '6 KM', 'Nee', '', 'true', 'true', '2025-03-08 10:28:42.2+00');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM aanmeldingen WHERE id = '90e477cc-89d1-4524-8edc-b697be8c504d') THEN
        INSERT INTO "public"."aanmeldingen" ("id", "created_at", "updated_at", "naam", "email", "telefoon", "rol", "afstand", "ondersteuning", "bijzonderheden", "terms", "email_verzonden", "email_verzonden_op") 
        VALUES ('90e477cc-89d1-4524-8edc-b697be8c504d', '2025-03-08 10:27:31.078013+00', '2025-03-08 10:27:32.642032+00', 'Angelo van Ingen', 'Enckerkamp.27@sheerenloo.nl', null, 'Deelnemer', '6 KM', 'Nee', '', 'true', 'true', '2025-03-08 10:27:36.787+00');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM aanmeldingen WHERE id = '26ea058b-2608-49d0-862a-611e98d7dc61') THEN
        INSERT INTO "public"."aanmeldingen" ("id", "created_at", "updated_at", "naam", "email", "telefoon", "rol", "afstand", "ondersteuning", "bijzonderheden", "terms", "email_verzonden", "email_verzonden_op") 
        VALUES ('26ea058b-2608-49d0-862a-611e98d7dc61', '2025-02-20 11:55:30.818333+00', '2025-02-20 11:55:32.209292+00', 'Janny van de Wall', 'mjvdwal@hotmail.com', null, 'Deelnemer', '10 KM', 'Nee', '', 'true', 'true', '2025-02-20 11:55:33.048+00');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM aanmeldingen WHERE id = 'f4fc2312-ec8a-4dfc-90b5-a8da317618e6') THEN
        INSERT INTO "public"."aanmeldingen" ("id", "created_at", "updated_at", "naam", "email", "telefoon", "rol", "afstand", "ondersteuning", "bijzonderheden", "terms", "email_verzonden", "email_verzonden_op") 
        VALUES ('f4fc2312-ec8a-4dfc-90b5-a8da317618e6', '2025-01-28 21:58:37.55756+00', '2025-01-28 21:58:38.527642+00', 'Martin van der Wal', 'mjvdwal@hotmail.com', '', 'Deelnemer', '10 KM', 'Ja', 'loopt samen met Dirk-Jan mee als vrijwilliger', 'true', 'true', '2025-01-28 21:58:39.243+00');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM aanmeldingen WHERE id = '4bfe814b-e0b8-4e60-9f46-fe38852d9ecb') THEN
        INSERT INTO "public"."aanmeldingen" ("id", "created_at", "updated_at", "naam", "email", "telefoon", "rol", "afstand", "ondersteuning", "bijzonderheden", "terms", "email_verzonden", "email_verzonden_op") 
        VALUES ('4bfe814b-e0b8-4e60-9f46-fe38852d9ecb', '2025-01-28 21:56:10.099964+00', '2025-01-28 21:56:11.309572+00', 'Dirk-Jan Hempe', 'mjvdwal@hotmail.com', '', 'Deelnemer', '10 KM', 'Nee', '', 'true', 'true', '2025-01-28 21:56:12.004+00');
    END IF;
END;
$$; 