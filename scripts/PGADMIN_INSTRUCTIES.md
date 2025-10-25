# PgAdmin Instructies - Deelnemers Accounts Aanmaken

## Stap 1: Open pgAdmin

1. Start pgAdmin 4
2. Connect met je database (dekoninklijkeloopdatabase)
3. Klik op de database naam in de linkerbalk

## Stap 2: Open Query Tool

1. Klik rechts op je database naam
2. Selecteer: **Query Tool** (of druk F5)
3. Een nieuw SQL editor venster opent

## Stap 3: Kopieer SQL Script

1. Open het bestand: `scripts/create_participant_accounts.sql`
2. Selecteer ALLES (Ctrl+A)
3. Kopieer (Ctrl+C)

## Stap 4: Plak in Query Tool

1. Ga terug naar pgAdmin Query Tool
2. Plak de SQL code (Ctrl+V)

## Stap 5: Voer Script Uit

1. Klik op de **Execute/Run** knop (▶️ play icon)
   - Of druk F5
   - Of ga naar Query → Execute/Refresh
2. Wacht tot het script klaar is (enkele seconden)

## Stap 6: Bekijk Resultaten

Je ziet meerdere output tabs onderaan:

### Tab 1: INSERT Resultaat
```
INSERT 0 XX  (XX nieuwe gebruikersaccounts aangemaakt)
```

### Tab 2: UPDATE Resultaat  
```
UPDATE XX    (XX aanmeldingen gelinkt)
```

### Tab 3: RBAC Rollen
```
INSERT 0 XX  (XX user_roles entries aangemaakt)
```

### Tab 4-8: Verificatie Queries
- Aantal deelnemers met account
- Aantal aanmeldingen gelinkt
- Lijst van alle nieuwe accounts
- RBAC rollen overzicht
- Eventuele problemen

## Verwachte Output Voorbeeld

```
Data Output - Tab 1
╔════════════════════════════════╗
║  INSERT 0 54                   ║
╚════════════════════════════════╝

Data Output - Tab 2
╔════════════════════════════════╗
║  UPDATE 54                     ║
╚════════════════════════════════╝

Data Output - Tab 3
╔════════════════════════════════╗
║  INSERT 0 54                   ║
╚════════════════════════════════╝

Data Output - Tab 4 (Verificatie)
╔═══════════════════════════╤════════╗
║ category                  │ count  ║
╠═══════════════════════════╪════════╣
║ Deelnemers met account    │ 54     ║
║ Aanmeldingen gelinkt      │ 54     ║
║ Aanmeldingen NIET gelinkt │ 0      ║
╚═══════════════════════════╧════════╝

Data Output - Tab 5 (Gebruikers Lijst)
╔════════╤═══════════╤════════════════════════════╤══════════════╗
║ id     │ naam      │ email                      │ rol          ║
╠════════╪═══════════╪════════════════════════════╪══════════════╣
║ uuid   │ Theun     │ diesbosje@hotmail.com      │ deelnemer    ║
║ uuid   │ Diesmer   │ diesbosje@hotmail.com      │ begeleider   ║
║ ...    │ ...       │ ...                        │ ...          ║
╚════════╧═══════════╧════════════════════════════╧══════════════╝
```

## Troubleshooting

### Error: "duplicate key value violates unique constraint"
- **Oorzaak:** Email bestaat al in gebruikers tabel
- **Oplossing:** Dit is normaal! Het script slaat deze over met `NOT EXISTS`
- **Actie:** Geen, het script werkt correct

### Error: "relation does not exist"
- **Oorzaak:** Database migratie V1_33 is niet uitgevoerd
- **Oplossing:** Wacht tot deployment succesvol is afgerond
- **Check:** Kijk of `aanmeldingen.gebruiker_id` kolom bestaat

### Error: "column does not exist"
- **Oorzaak:** Verkeerde database geselecteerd
- **Oplossing:** Zorg dat je verbonden bent met: `dekoninklijkeloopdatabase`

### Geen output zichtbaar
- **Oplossing:** Scroll naar beneden in het Query Tool venster
- Klik op de verschillende tabs (Data Output, Messages, etc.)

## Na Uitvoeren

### Test een Login:

Open een nieuwe Query tab en test:

```sql
-- Check of een specifieke gebruiker bestaat
SELECT 
    id, naam, email, rol, is_actief
FROM gebruikers 
WHERE email = 'diesbosje@hotmail.com';

-- Check RBAC rol koppeling
SELECT 
    g.naam, g.email, g.rol as gebruiker_rol,
    r.name as rbac_rol, r.description
FROM gebruikers g
JOIN user_roles ur ON ur.user_id = g.id
JOIN roles r ON r.id = ur.role_id
WHERE g.email = 'diesbosje@hotmail.com'
AND ur.is_active = true;
```

### Test Aanmelding Link:

```sql
-- Check of aanmelding gelinkt is aan gebruiker
SELECT 
    a.naam as aanmelding_naam,
    a.email as aanmelding_email,
    a.rol as aanmelding_rol,
    g.naam as gebruiker_naam,
    g.email as gebruiker_email,
    g.rol as gebruiker_rol
FROM aanmeldingen a
LEFT JOIN gebruikers g ON g.id = a.gebruiker_id
WHERE a.email = 'diesbosje@hotmail.com';
```

## Stap 7: Verifieer in Applicatie

Nu kunnen deelnemers inloggen:

**Test Login:**
```
URL: https://api.dekoninklijkeloop.nl/api/auth/login
Method: POST
Body: {
  "email": "diesbosje@hotmail.com",
  "wachtwoord": "DKL2025!"
}
```

**Verwacht Response:**
```json
{
  "success": true,
  "token": "eyJhbGc...",
  "user": {
    "id": "uuid",
    "email": "diesbosje@hotmail.com",
    "naam": "Diesmer",
    "rol": "begeleider"
  }
}
```

## Belangrijk! 🔐

- **Standaard wachtwoord:** DKL2025!
- Stuur alle deelnemers een email met hun inloggegevens
- Vraag ze om hun wachtwoord te wijzigen na eerste login
- Bewaar het script voor toekomstige deelnemers

## Hulp Nodig?

Als je errors krijgt of twijfelt:
1. Maak een screenshot van de error
2. Check welke database je hebt geselecteerd
3. Verifieer dat migratie V1_33 is uitgevoerd