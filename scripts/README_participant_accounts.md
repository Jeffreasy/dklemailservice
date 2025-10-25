# Deelnemers Accounts Aanmaken

Deze scripts zorgen ervoor dat alle deelnemers uit de `aanmeldingen` tabel een gebruikersaccount krijgen om in te loggen in de step tracker app.

## Stappen

### 1. Database Migratie Uitvoeren

Voer eerst de migratie uit die de `gebruiker_id` kolom toevoegt aan de aanmeldingen tabel:

```bash
# De migratie wordt automatisch uitgevoerd bij de volgende deployment
# Of handmatig uitvoeren:
psql -d your_database -f database/migrations/V1_33__add_gebruiker_id_to_aanmeldingen.sql
```

### 2. Gebruikersaccounts Aanmaken

**Optie A: SQL Script (Aanbevolen voor productie)**

Voer het SQL script direct uit op de database:

```bash
psql -d your_database -f scripts/create_participant_accounts.sql
```

Of via database console (Render, Supabase, etc.):
1. Open de database SQL editor
2. Kopieer de inhoud van `scripts/create_participant_accounts.sql`
3. Voer het script uit
4. Bekijk de verificatie output

**Optie B: Go Script (Lokaal)**

```bash
cd scripts
go run create_participant_accounts.go
```

**Beide scripts:**
- ✅ Maken nieuwe gebruikersaccounts aan voor alle deelnemers zonder account
- ✅ Gebruiken het standaard wachtwoord: **DKL2025!**
- ✅ Linken aanmeldingen aan gebruikersaccounts via email matching
- ✅ Slaan duplicaten over (als een email al bestaat)
- ✅ Geven een samenvatting van de acties

### 3. Resultaat

Na het uitvoeren zie je een output zoals:

```
Verbonden met database
Gevonden 54 aanmeldingen zonder gebruikersaccount

✓ Gebruiker aangemaakt voor: Theun (diesbosje@hotmail.com)
✓ Gebruiker aangemaakt voor: Han van Doornik (LaanvanGS.26@sheerenloo.nl)
→ Aanmelding gelinkt aan bestaande gebruiker: Anneke van de Glind (Klaskehiddes@gmail.com)
...

=== Samenvatting ===
Nieuwe gebruikers aangemaakt: 45
Aanmeldingen gelinkt: 54
Overgeslagen (fouten): 0

Standaard wachtwoord voor nieuwe accounts: DKL2025!
Gebruikers moeten hun wachtwoord wijzigen via de app!
```

## Wat is er aangepast?

### 1. Database Schema
- ✅ Nieuwe kolom `gebruiker_id` in `aanmeldingen` tabel
- ✅ Foreign key relatie naar `gebruikers` tabel
- ✅ Index voor snellere queries

### 2. Models
- ✅ `models.Aanmelding` heeft nu `GebruikerID` veld
- ✅ Link tussen aanmelding en gebruiker

### 3. Services
- ✅ Nieuwe methode `GetParticipantDashboardByUserID()` in `StepsService`
- ✅ Dashboard data ophalen via gebruiker ID

### 4. Handlers
- ✅ `GetParticipantDashboard()` gebruikt nu gebruiker ID uit context
- ✅ Fallback naar ID parameter voor admin toegang

## Gebruikers Informeren

Stuur alle deelnemers een email met:

**Onderwerp:** Inloggegevens De Koninklijke Loop Step Tracker

**Bericht:**
```
Beste deelnemer,

Je kunt nu inloggen in de De Koninklijke Loop step tracker app!

Email: [hun email adres]
Tijdelijk wachtwoord: DKL2025!

Wijzig je wachtwoord na eerste inlog via je profiel.

Veel succes met stappen verzamelen!

Team De Koninklijke Loop
```

## Wachtwoord Reset

Deelnemers kunnen hun wachtwoord wijzigen via:
- POST /api/auth/reset-password (als ze ingelogd zijn)

Body:
```json
{
  "huidig_wachtwoord": "DKL2025!",
  "nieuw_wachtwoord": "hun_nieuwe_wachtwoord"
}
```

## Troubleshooting

### "Database connection failed"
- Check of de `.env` file correct is ingevuld
- Verifieer database credentials

### "Kon gebruiker niet aanmaken"
- Mogelijk duplicate email (normaal, wordt overgeslagen)
- Check database logs voor details

### "Deelnemer niet gevonden" error
- Run het script opnieuw om ontbrekende links te maken
- Verifieer of de gebruiker_id correct is gezet

## Verificatie

Controleer of alles correct is:

```sql
-- Tel gebruikers met rol 'deelnemer'
SELECT COUNT(*) FROM gebruikers WHERE rol = 'deelnemer';

-- Tel aanmeldingen met gebruiker_id
SELECT COUNT(*) FROM aanmeldingen WHERE gebruiker_id IS NOT NULL;

-- Bekijk aanmeldingen zonder gebruiker_id
SELECT id, naam, email FROM aanmeldingen WHERE gebruiker_id IS NULL;
```

## Security Note

⚠️ **Belangrijk:**
- Het standaard wachtwoord is **tijdelijk**
- Gebruikers **moeten** hun wachtwoord wijzigen
- Overweeg geforceerde wachtwoord wijziging bij eerste login te implementeren