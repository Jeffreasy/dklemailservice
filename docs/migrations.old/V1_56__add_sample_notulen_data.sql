-- V1_56__add_sample_notulen_data.sql
-- Add sample notulen data for testing and demonstration

-- Insert the new notulen data
DO $$
DECLARE
    admin_user_id UUID;
BEGIN
    -- Get the admin user ID from the gebruikers table
    SELECT id INTO admin_user_id FROM gebruikers WHERE email = 'admin@dekoninklijkeloop.nl' LIMIT 1;

    -- Only proceed if we found an admin user
    IF admin_user_id IS NOT NULL THEN
        IF NOT EXISTS (SELECT 1 FROM notulen WHERE titel = 'Notulen 30 oktober 2025') THEN
            INSERT INTO notulen (
                titel,
                vergadering_datum,
                locatie,
                voorzitter,
                aanwezigen,
                afwezigen,
                agenda_items,
                besluiten,
                actiepunten,
                notities,
                status,
                created_by
            ) VALUES (
                'Notulen 30 oktober 2025',
                '2025-10-30',
                'Teams',
                'Salih',
                ARRAY['Salih', 'Angelique', 'Jeffrey', 'Marieke', 'Lida', 'Ginelli'],
                ARRAY[]::TEXT[],
                '{"items": [{"title": "Afspraken voor DKL 2026", "details": "Datum, Doel, Streefdoel, Organisatie teams, Vaste meetings, Partners, Website & online"}, {"title": "Ideeën voor DKL26", "details": "Social media, Mascotte, Merch, Betrokkenheid doelgroep"}]}'::JSONB,
                '{"besluiten": [{"besluit": "Datum: zaterdag 16 mei 2026"}, {"besluit": "Doel: Only Friends[](https://onlyfriends.nl/), stichting in Apeldoorn, goed doel omdat het lokaal is en mooi aansluit bij samenwerking met Alessandro Bistolfi."}, {"besluit": "Streefdoel aantal deelnemers: 100"}, {"besluit": "Verdeling van de organisatie in miniteams", "teams": {"Algemene organisatie": ["Salih", "Angelique", "Jeffrey", "Marieke"], "Sponsoring/subsidie aanvragen": ["evt Marieke"], "PR & communicatie": ["Lida", "Ginelli", "Jeffrey", "Marieke"], "Doelgroep betrokkenheid": [], "Partners en samenwerkingen": ["Salih"]}}, {"besluit": "Vaste meetings: Eenmaal in de 6 weken op de maandag in Teams"}, {"besluit": "Partners: Accres (beheer sportfaciliteiten en evenementen voor de gemeente Apeldoorn), Alessandro Bistolfi van Nedarg handbikes[](https://www.nedarg.com)"}, {"besluit": "Website & online: Stappenteller ontwikkeld waaraan donatie optie gekoppeld, iedere deelnemer krijgt account voor prestaties en gezamenlijke teller, opties via site of app"}]}'::JSONB,
                '{"acties": [{"actie": "Salih gaat een meeting cyclus op Teams aanmaken.", "verantwoordelijke": "Salih"}, {"actie": "Salih neemt contact op met Alessandro.", "verantwoordelijke": "Salih"}, {"actie": "Lida kijkt naar scholen in Harderwijk.", "verantwoordelijke": "Lida"}, {"actie": "Salih kijkt naar optie bij Nijmegen: https://www.maartenschool.nl/home.", "verantwoordelijke": "Salih"}, {"actie": "Jeffrey deelt eerste versie van de app.", "verantwoordelijke": "Jeffrey"}, {"actie": "Angelique benadert collega''s via intranet.", "verantwoordelijke": "Angelique"}, {"actie": "Lida en Ginelli gaan aan de slag met social media content (vlogs, live, intro videos).", "verantwoordelijke": ["Lida", "Ginelli"]}, {"actie": "Brainstorm over eigen mascotte (met doelgroep bedenken of zelf ontwikkelen).", "verantwoordelijke": "Team"}, {"actie": "Angelique kijkt naar merch opties (draagtasjes, bekers, later T-shirts) met Wesleys en collega grafisch ontwerper.", "verantwoordelijke": "Angelique"}, {"actie": "Salih informeert bij Kathelijn over betrokkenheid doelgroep met cliëntenraad (pas vanaf januari).", "verantwoordelijke": "Salih"}]}'::JSONB,
                'Pas vanaf januari doelgroep betrekken, anders te vroeg.',
                'draft',
                admin_user_id
            );
        END IF;
    END IF;
END $$;