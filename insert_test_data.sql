-- Contactformulieren toevoegen
INSERT INTO contact_formulieren (
    id, created_at, updated_at, naam, email, bericht, 
    email_verzonden, privacy_akkoord, status
) VALUES 
    ('4a2eee9a-1a28-4534-87f4-83a70cd5161b', '2025-03-15T00:55:59.272Z', '2025-03-15T00:55:59.272Z', 'Test Contact 1', 'test1@example.com', 'Dit is een testbericht voor contactformulier 1.', true, true, 'nieuw'),
    ('7f4db580-7383-4d30-984e-cd9c8fcf3d24', '2025-03-15T00:55:59.272Z', '2025-03-15T00:55:59.272Z', 'Test Contact 2', 'test2@example.com', 'Dit is een testbericht voor contactformulier 2.', true, true, 'in_behandeling'),
    ('e0151d3f-605c-4622-8700-29d2a49f8c47', '2025-03-15T00:55:59.272Z', '2025-03-15T00:55:59.272Z', 'Test Contact 3', 'test3@example.com', 'Dit is een testbericht voor contactformulier 3.', true, true, 'beantwoord');

-- Antwoord toevoegen aan het laatste contactformulier
INSERT INTO contact_antwoorden (
    id, contact_id, tekst, verzonden_op, verzonden_door, email_verzonden
) VALUES 
    ('e1fd7df6-c6f7-4310-b117-94157251e58f', 'e0151d3f-605c-4622-8700-29d2a49f8c47', 'Dit is een testantwoord op contactformulier 3.', '2025-03-15T00:55:59.272Z', 'admin@dekoninklijkeloop.nl', true);

-- Update contactformulier met antwoord
UPDATE contact_formulieren 
SET beantwoord = true, 
    antwoord_tekst = 'Dit is een testantwoord op contactformulier 3.', 
    antwoord_datum = '2025-03-15T00:55:59.272Z', 
    antwoord_door = 'admin@dekoninklijkeloop.nl' 
WHERE id = 'e0151d3f-605c-4622-8700-29d2a49f8c47';

-- Aanmeldingen toevoegen
INSERT INTO aanmeldingen (
    id, created_at, updated_at, naam, email, telefoon, 
    rol, afstand, ondersteuning, bijzonderheden, terms, 
    email_verzonden, status
) VALUES 
    ('b97e4783-997c-4757-a740-88c3f55e435b', '2025-03-15T00:55:59.272Z', '2025-03-15T00:55:59.272Z', 'Test Aanmelding 1', 'aanmelding1@example.com', '0612345678', 'deelnemer', '10km', 'geen', 'Geen bijzonderheden', true, true, 'nieuw'),
    ('6259b809-f954-4cad-9dc3-93983909afb1', '2025-03-15T00:55:59.272Z', '2025-03-15T00:55:59.272Z', 'Test Aanmelding 2', 'aanmelding2@example.com', '0687654321', 'vrijwilliger', '', 'geen', 'Geen bijzonderheden', true, true, 'in_behandeling'),
    ('300d6125-9a97-4111-9dfc-069d633c2edc', '2025-03-15T00:55:59.272Z', '2025-03-15T00:55:59.272Z', 'Test Aanmelding 3', 'aanmelding3@example.com', '0611223344', 'sponsor', '', 'geen', 'Geen bijzonderheden', true, true, 'beantwoord');

-- Antwoord toevoegen aan het laatste aanmelding
INSERT INTO aanmelding_antwoorden (
    id, aanmelding_id, tekst, verzonden_op, verzonden_door, email_verzonden
) VALUES 
    ('a3cd2987-54b8-444e-ad23-cecdf1ef224b', '300d6125-9a97-4111-9dfc-069d633c2edc', 'Dit is een testantwoord op aanmelding 3.', '2025-03-15T00:55:59.272Z', 'admin@dekoninklijkeloop.nl', true);
