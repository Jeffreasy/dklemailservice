INSERT INTO "public"."aanmeldingen" ("id", "created_at", "updated_at", "naam", "email", "telefoon", "rol", "afstand", "ondersteuning", "bijzonderheden", "terms", "email_verzonden", "email_verzonden_op")
SELECT * FROM (VALUES
    ('9f464844-8c93-4190-90f0-e74765c7f09a', '2025-03-24 18:47:05.053651+00', '2025-03-24 18:47:05.053651+00', 'Karin de Jong', 'karin.de.jong82@outlook.com', null, 'Deelnemer', '6 KM', 'Nee', '', 'true', 'false', null),
    ('51855fec-eab9-494a-9321-c40d22da4ffc', '2025-03-24 17:27:25.191642+00', '2025-03-24 17:27:25.191642+00', 'Mirjam Kerkvliet', 'mirjam.kerkvliet@gmail.com', null, 'Deelnemer', '15 KM', 'Nee', '', 'true', 'false', null),
    ('47775742-8950-4b94-9dd1-571ff4902688', '2025-03-24 17:26:12.1724+00', '2025-03-24 17:26:12.1724+00', 'Arno Kerkvliet', 'arno.kerkvliet@gmail.com', null, 'Deelnemer', '15 KM', 'Nee', '', 'true', 'false', null),
    ('d92ed75c-c275-47a4-88a9-ff7a4106f8ee', '2025-03-24 09:17:59.726501+00', '2025-03-24 09:17:59.726501+00', 'Jean-paul Hup', 'molenkamp.19@sheerenloo.nl', null, 'Deelnemer', '6 KM', 'Nee', '', 'true', 'false', null),
    ('ecb8332b-ea39-4611-9f58-64921226f2a6', '2025-03-24 09:16:42.111808+00', '2025-03-24 09:16:42.111808+00', 'Annerieke Mandemaker-Timmer', 'annerieketimmer@hotmail.com', '06 17 37 28 40 ', 'Begeleider', '6 KM', 'Nee', '', 'true', 'false', null)
) AS new_registrations
WHERE NOT EXISTS (
    SELECT 1 FROM "public"."aanmeldingen" WHERE id IN (
        '9f464844-8c93-4190-90f0-e74765c7f09a',
        '51855fec-eab9-494a-9321-c40d22da4ffc',
        '47775742-8950-4b94-9dd1-571ff4902688',
        'd92ed75c-c275-47a4-88a9-ff7a4106f8ee',
        'ecb8332b-ea39-4611-9f58-64921226f2a6'
    )
); 