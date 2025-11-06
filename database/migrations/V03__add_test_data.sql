-- GECONSOLIDEERDE V3 - TEST DATA
-- Dit bestand combineert de logica van V1_11, V1_12, V1_13, en V1_14.
-- Alle INSERTs zijn gestandaardiseerd naar de efficiënte 'ON CONFLICT' syntax.

-- Toevoegen van alle test-aanmeldingen (van V1_11, V1_13, V1_14)
INSERT INTO "public"."aanmeldingen" 
(
    "id", 
    "naam", 
    "email", 
    "telefoon", 
    "status", 
    "created_at", 
    "updated_at", 
    "rol", 
    "afstand", 
    "ondersteuning", 
    "bijzonderheden", 
    "terms", 
    "email_verzonden", 
    "email_verzonden_op", 
    "test_mode"
)
VALUES 
-- Data van V1_11 (kolommen 'status' en 'test_mode' toegevoegd)
('3e62d5d3-070d-47b1-a1ef-30665f982789', 'TGTest', 'laventejeffrey@gmail.com', '06123456789', 'nieuw', '2025-03-23 17:06:26.132297+00', '2025-03-23 17:06:26.132297+00', 'Begeleider', '15 KM', 'Anders', 'Telegram Test bericht - officiele weg', 'true', 'false', null, 'false'),
('51205069-a20a-4231-bcb9-cb2d6fd042c8', 'Manuela van Zwam', 'rik.van-harxen@sheerenloo.nl', null, 'nieuw', '2025-03-23 14:08:53.328628+00', '2025-03-23 14:08:53.328628+00', 'Deelnemer', '2.5 KM', 'Ja', 'Vaste begeleider die meeloopt', 'true', 'false', null, 'false'),
('275490c0-1021-4bf4-9005-7df9884b0fe6', 'Bas heijenk ', 'basheijenk96@gmail.com', null, 'nieuw', '2025-03-22 16:43:03.19496+00', '2025-03-22 16:43:03.19496+00', 'Deelnemer', '2.5 KM', 'Nee', '', 'true', 'false', null, 'false'),
('b5e67c64-6bfa-46d4-ae40-98b093c8b720', 'Salih', 'topraks@gmail.com', null, 'verwerkt', '2025-03-17 19:55:07.379647+00', '2025-03-21 06:47:40.512037+00', 'Deelnemer', '2.5 KM', 'Nee', '', 'true', 'true', '2025-03-21 06:47:38.608+00', 'false'),
('b2fd3412-8368-409f-8029-b2cdd581ade1', 'Manuela van zwam', 'benjaminlaan.64a@sheerenloo.nl', null, 'verwerkt', '2025-03-11 11:28:05.483427+00', '2025-03-21 06:47:41.806005+00', 'Deelnemer', '2.5 KM', 'Nee', '', 'true', 'true', '2025-03-21 06:47:39.928+00', 'false'),
('391f63c5-f034-466e-8a1f-ba9d06ed1192', 'Joyce Thielen', 'Joyce.thielen@sheerenloo.nl', '', 'verwerkt', '2025-03-09 16:52:09.564437+00', '2025-03-21 06:47:42.987352+00', 'Begeleider', '6 KM', 'Nee', '', 'true', 'true', '2025-03-21 06:47:41.071+00', 'false'),
('391f2579-d7cb-4ef3-afbe-14dc4115c519', 'Dick van Norden', 'Enckerkamp.27@sheerenloo.nl', null, 'verwerkt', '2025-03-08 10:28:37.053379+00', '2025-03-08 10:28:38.042991+00', 'Deelnemer', '6 KM', 'Nee', '', 'true', 'true', '2025-03-08 10:28:42.2+00', 'false'),
('90e477cc-89d1-4524-8edc-b697be8c504d', 'Angelo van Ingen', 'Enckerkamp.27@sheerenloo.nl', null, 'verwerkt', '2025-03-08 10:27:31.078013+00', '2025-03-08 10:27:32.642032+00', 'Deelnemer', '6 KM', 'Nee', '', 'true', 'true', '2025-03-08 10:27:36.787+00', 'false'),
('26ea058b-2608-49d0-862a-611e98d7dc61', 'Janny van de Wall', 'mjvdwal@hotmail.com', null, 'verwerkt', '2025-02-20 11:55:30.818333+00', '2025-02-20 11:55:32.209292+00', 'Deelnemer', '10 KM', 'Nee', '', 'true', 'true', '2025-02-20 11:55:33.048+00', 'false'),
('f4fc2312-ec8a-4dfc-90b5-a8da317618e6', 'Martin van der Wal', 'mjvdwal@hotmail.com', '', 'verwerkt', '2025-01-28 21:58:37.55756+00', '2025-01-28 21:58:38.527642+00', 'Deelnemer', '10 KM', 'Ja', 'loopt samen met Dirk-Jan mee als vrijwilliger', 'true', 'true', '2025-01-28 21:58:39.243+00', 'false'),
('4bfe814b-e0b8-4e60-9f46-fe38852d9ecb', 'Dirk-Jan Hempe', 'mjvdwal@hotmail.com', '', 'verwerkt', '2025-01-28 21:56:10.099964+00', '2025-01-28 21:56:11.309572+00', 'Deelnemer', '10 KM', 'Nee', '', 'true', 'true', '2025-01-28 21:56:12.004+00', 'false'),

-- Data van V1_13
('9f464844-8c93-4190-90f0-e74765c7f09a', 'Karin de Jong', 'karin.de.jong82@outlook.com', NULL, 'nieuw', '2025-03-24 18:47:05.053651', '2025-03-24 18:47:05.053651', 'Deelnemer', '6 KM', 'Nee', '', TRUE, FALSE, NULL, FALSE),
('51855fec-eab9-494a-9321-c40d22da4ffc', 'Mirjam Kerkvliet', 'mirjam.kerkvliet@gmail.com', NULL, 'nieuw', '2025-03-24 17:27:25.191642', '2025-03-24 17:27:25.191642', 'Deelnemer', '15 KM', 'Nee', '', TRUE, FALSE, NULL, FALSE),
('47775742-8950-4b94-9dd1-571ff4902688', 'Arno Kerkvliet', 'arno.kerkvliet@gmail.com', NULL, 'nieuw', '2025-03-24 17:26:12.1724', '2025-03-24 17:26:12.1724', 'Deelnemer', '15 KM', 'Nee', '', TRUE, FALSE, NULL, FALSE),
('d92ed75c-c275-47a4-88a9-ff7a4106f8ee', 'Jean-paul Hup', 'molenkamp.19@sheerenloo.nl', NULL, 'nieuw', '2025-03-24 09:17:59.726501', '2025-03-24 09:17:59.726501', 'Deelnemer', '6 KM', 'Nee', '', TRUE, FALSE, NULL, FALSE),
('ecb8332b-ea39-4611-9f58-64921226f2a6', 'Annerieke Mandemaker-Timmer', 'annerieketimmer@hotmail.com', '06 17 37 28 40 ', 'nieuw', '2025-03-24 09:16:42.111808', '2025-03-24 09:16:42.111808', 'Begeleider', '6 KM', 'Nee', '', TRUE, FALSE, NULL, FALSE),

-- Data van V1_14
('1ca80f61-f5c1-431f-b224-e6557150b65b', 'Han van Doornik', 'LaanvanGS.26@sheerenloo.nl', null, 'nieuw', '2025-03-30 08:07:45.334762+00', '2025-03-30 08:07:45.334762+00', 'Deelnemer', '2.5 KM', 'Ja', 'Ik wil wel graag begeleiding ', 'true', 'false', null, FALSE),
('2499686e-2e62-4827-9079-78b468cb26c9', 'Bertram tijsma', 'Klaskehiddes@gmail.com', null, 'verwerkt', '2025-03-29 10:07:12.393466+00', '2025-03-29 10:07:56.169434+00', 'Deelnemer', '15 KM', 'Nee', '', 'true', 'true', '2025-03-29 10:07:57.838+00', FALSE),
('9f75b1df-4c72-4e36-9901-0f74cc26574f', 'Klaske van de glind', 'Klaskehiddes@gmail.com', null, 'verwerkt', '2025-03-29 10:05:04.276906+00', '2025-03-29 10:07:55.533152+00', 'Deelnemer', '15 KM', 'Nee', '', 'true', 'true', '2025-03-29 10:07:57.215+00', FALSE),
('db3ec762-dd54-4ba7-98ab-981235cc316a', 'Mila Veenendaal', 'gaminggirlayla@gmail.com', null, 'verwerkt', '2025-03-26 16:26:53.777384+00', '2025-03-26 16:30:58.285233+00', 'Deelnemer', '10 KM', 'Nee', '', 'true', 'true', '2025-03-26 16:30:57.269+00', FALSE),
('d17c16c6-c423-43de-a876-d40326b62d9e', 'Ayla Toprak', 'gamergirlayla@gmail.com', null, 'verwerkt', '2025-03-26 16:25:43.756211+00', '2025-03-26 16:30:57.680904+00', 'Deelnemer', '10 KM', 'Nee', '', 'true', 'true', '2025-03-26 16:30:56.685+00', FALSE),
('917206a7-a28d-4bad-8b41-ed127eab743a', 'A. Bistolfi', 'nedarg@icloud.com', null, 'verwerkt', '2025-03-26 12:27:20.236848+00', '2025-03-26 16:30:57.073558+00', 'Deelnemer', '15 KM', 'Nee', '', 'true', 'true', '2025-03-26 16:30:56.062+00', FALSE)
ON CONFLICT (id) DO NOTHING;

-- Toevoegen van test-contactformulieren (van V1_12)
INSERT INTO "public"."contact_formulieren" 
(
    "id", 
    "naam", 
    "email", 
    "bericht", 
    "status", 
    "created_at", 
    "updated_at", 
    "email_verzonden", 
    "email_verzonden_op", 
    "privacy_akkoord", 
    "behandeld_door", 
    "behandeld_op", 
    "notities",
    "test_mode"
) 
VALUES 
('6ce2e9a3-59fd-4430-aa5c-66df48fbd695', 'je geheime liefde', 'de.konining@willem.alexander.nl', 'Gedeelte doneren doet het niet.', 'afgehandeld', '2025-01-28 00:03:05.830054+00', '2025-02-05 02:23:23.029501+00', 'true', '2025-01-28 00:03:08.041+00', 'true', 'marieke@dekoninklijkeloop.nl', '2025-02-05 02:23:22.97+00', null, 'false'),
('78910428-f760-485d-ae57-653db478ca35', 'Bas heijenk ', 'basheijenk96@gmail.com', 'Hallo ik heb me op gegeven maar ik kan helaas niet sorry ', 'nieuw', '2025-03-23 07:15:20.8851+00', '2025-03-23 07:15:20.8851+00', 'false', null, 'true', null, null, null, 'false')
ON CONFLICT (id) DO NOTHING;