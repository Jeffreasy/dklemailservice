-- V1_12__add_test_contact_forms.sql
-- Voegt testgegevens toe voor contactformulieren

-- Alleen uitvoeren als de contactformulieren nog niet bestaan
DO $$
BEGIN
    -- Contactformulieren toevoegen als ze niet bestaan
    IF NOT EXISTS (SELECT 1 FROM contact_formulieren WHERE id = '6ce2e9a3-59fd-4430-aa5c-66df48fbd695') THEN
        INSERT INTO "public"."contact_formulieren" ("id", "created_at", "updated_at", "naam", "email", "bericht", "email_verzonden", "email_verzonden_op", "privacy_akkoord", "status", "behandeld_door", "behandeld_op", "notities") 
        VALUES ('6ce2e9a3-59fd-4430-aa5c-66df48fbd695', '2025-01-28 00:03:05.830054+00', '2025-02-05 02:23:23.029501+00', 'je geheime liefde', 'de.konining@willem.alexander.nl', 'Gedeelte doneren doet het niet.', 'true', '2025-01-28 00:03:08.041+00', 'true', 'afgehandeld', 'marieke@dekoninklijkeloop.nl', '2025-02-05 02:23:22.97+00', null);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM contact_formulieren WHERE id = '78910428-f760-485d-ae57-653db478ca35') THEN
        INSERT INTO "public"."contact_formulieren" ("id", "created_at", "updated_at", "naam", "email", "bericht", "email_verzonden", "email_verzonden_op", "privacy_akkoord", "status", "behandeld_door", "behandeld_op", "notities") 
        VALUES ('78910428-f760-485d-ae57-653db478ca35', '2025-03-23 07:15:20.8851+00', '2025-03-23 07:15:20.8851+00', 'Bas heijenk ', 'basheijenk96@gmail.com', 'Hallo ik heb me op gegeven maar ik kan helaas niet sorry ', 'false', null, 'true', 'nieuw', null, null, null);
    END IF;
END;
$$; 