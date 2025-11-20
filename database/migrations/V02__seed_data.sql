-- GECONSOLIDEERDE V2 - SEED DATA
-- Dit is de inhoud van 002_seed_data.sql,
-- maar zonder de 'INSERT INTO migraties' die bij het oude systeem hoorde.

-- Controleer of er al een admin gebruiker is
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM gebruikers WHERE email = 'admin@dekoninklijkeloop.nl') THEN
        -- Maak admin gebruiker aan (wachtwoord: admin)
        -- NOTE: rol column wordt later toegevoegd en gemigreerd naar RBAC systeem
        INSERT INTO gebruikers (naam, email, wachtwoord_hash, is_actief, created_at, updated_at)
        VALUES (
            'Admin',
            'admin@dekoninklijkeloop.nl',
            '$2a$10$5Yse5i2BJV.bwTzbmywa9e/3G.XxzQPayGPlTsut/nBrZr05pKMCK',
            TRUE,
            CURRENT_TIMESTAMP,
            CURRENT_TIMESTAMP
        );
    END IF;
END $$;

-- Haal de admin gebruiker ID op
DO $$
DECLARE
    admin_id UUID;
BEGIN
    -- Wacht tot de gebruiker daadwerkelijk is aangemaakt (voor het geval de vorige DO-blok nog bezig is)
    -- In de praktijk zal dit zelden nodig zijn, maar het is veiliger.
    SELECT id INTO admin_id FROM gebruikers WHERE email = 'admin@dekoninklijkeloop.nl';

    -- Als admin_id nog steeds NULL is (om wat voor reden dan ook), stop dan.
    IF admin_id IS NULL THEN
        RAISE NOTICE 'Admin gebruiker niet gevonden, kan email templates niet seeden.';
        RETURN;
    END IF;

    -- Maak standaard email templates aan als ze nog niet bestaan
    IF NOT EXISTS (SELECT 1 FROM email_templates WHERE naam = 'contact_admin_email') THEN
        INSERT INTO email_templates (naam, onderwerp, inhoud, beschrijving, is_actief, created_at, updated_at, created_by)
        VALUES (
            'contact_admin_email',
            'Nieuw contactformulier',
            '<p>Er is een nieuw contactformulier ingevuld door {{.Contact.Naam}}.</p><p>Email: {{.Contact.Email}}</p><p>Bericht: {{.Contact.Bericht}}</p>',
            'Email die naar de admin wordt gestuurd bij een nieuw contactformulier',
            TRUE,
            CURRENT_TIMESTAMP,
            CURRENT_TIMESTAMP,
            admin_id
        );
    END IF;

    IF NOT EXISTS (SELECT 1 FROM email_templates WHERE naam = 'contact_email') THEN
        INSERT INTO email_templates (naam, onderwerp, inhoud, beschrijving, is_actief, created_at, updated_at, created_by)
        VALUES (
            'contact_email',
            'Bedankt voor je bericht',
            '<p>Beste {{.Contact.Naam}},</p><p>Bedankt voor je bericht. We nemen zo snel mogelijk contact met je op.</p>',
            'Bevestigingsemail die naar de gebruiker wordt gestuurd bij een contactformulier',
            TRUE,
            CURRENT_TIMESTAMP,
            CURRENT_TIMESTAMP,
            admin_id
        );
    END IF;

    IF NOT EXISTS (SELECT 1 FROM email_templates WHERE naam = 'aanmelding_admin_email') THEN
        INSERT INTO email_templates (naam, onderwerp, inhoud, beschrijving, is_actief, created_at, updated_at, created_by)
        VALUES (
            'aanmelding_admin_email',
            'Nieuwe aanmelding ontvangen',
            '<p>Er is een nieuwe aanmelding ontvangen van {{.Aanmelding.Naam}}.</p><p>Email: {{.Aanmelding.Email}}</p>',
            'Email die naar de admin wordt gestuurd bij een nieuwe aanmelding',
            TRUE,
            CURRENT_TIMESTAMP,
            CURRENT_TIMESTAMP,
            admin_id
        );
    END IF;

    IF NOT EXISTS (SELECT 1 FROM email_templates WHERE naam = 'aanmelding_email') THEN
        INSERT INTO email_templates (naam, onderwerp, inhoud, beschrijving, is_actief, created_at, updated_at, created_by)
        VALUES (
            'aanmelding_email',
            'Bedankt voor je aanmelding',
            '<p>Beste {{.Aanmelding.Naam}},</p><p>Bedankt voor je aanmelding. We hebben je aanmelding ontvangen en zullen deze zo snel mogelijk verwerken.</p>',
            'Bevestigingsemail die naar de gebruiker wordt gestuurd bij een aanmelding',
            TRUE,
            CURRENT_TIMESTAMP,
            CURRENT_TIMESTAMP,
            admin_id
        );
    END IF;
END $$;