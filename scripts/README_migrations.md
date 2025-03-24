# Database Migratie Instructies

Dit document bevat instructies voor het uitvoeren van database migraties om het schema te synchroniseren met de Go models.

## SQL Script Methode (Aanbevolen)

1. Open het bestand `scripts/schema_sync.sql` in een tekstverwerker
2. Kopieer de inhoud van dit bestand
3. Voer dit SQL script uit met een van de volgende methoden:
   
   ### Met pgAdmin (Grafische interface)
   
   1. Open pgAdmin en maak verbinding met de database
   2. Klik met de rechtermuisknop op de database en selecteer "Query Tool"
   3. Plak het SQL script en klik op "Execute"
   
   ### Met een andere PostgreSQL client
   
   Gebruik je favoriete PostgreSQL client om het script uit te voeren.
   
   ### Met psql (Command line)
   
   Als je psql hebt geïnstalleerd, kun je het script als volgt uitvoeren:
   
   ```bash
   psql -h dpg-cnbvvnf109ks73f9rvr0-a.oregon-postgres.render.com -U dklemailservice_user -d dklemailservice -f scripts/schema_sync.sql
   ```

## Database Verbindingsgegevens

- **Host**: dpg-cnbvvnf109ks73f9rvr0-a.oregon-postgres.render.com
- **Poort**: 5432
- **Gebruiker**: dklemailservice_user
- **Database**: dklemailservice
- **SSL Mode**: Require

## Migraties Samenvatting

De volgende migraties worden uitgevoerd:

1. **V1_8__sync_contact_formulieren**:
   - Verwijdert oude kolommen (onderwerp, ip_adres)
   - Voegt nieuwe kolommen toe voor email verwerking en privacy
   - Voegt velden toe voor behandeling en antwoorden

2. **V1_9__sync_aanmeldingen**:
   - Verwijdert oude kolommen (evenement, extra_info, ip_adres)
   - Voegt kolommen toe voor deelnemersrol en hardloopeigenschappen
   - Voegt behandelingsvelden toe

3. **V1_10__fix_antwoord_tables**:
   - Werkt de structuur van contact_antwoorden en aanmelding_antwoorden bij
   - Voegt indices toe voor betere prestaties
   - Controleert en repareert foreign key relaties

## Veiligheidsinformatie

Alle SQL instructies gebruiken `IF NOT EXISTS` en `DROP COLUMN IF EXISTS` om te voorkomen dat er fouten optreden als een kolom al bestaat of niet bestaat. Dit script kan veilig meerdere keren worden uitgevoerd.

## Na de Migratie

Na het uitvoeren van de migraties moet je de applicatieserver opnieuw opstarten om ervoor te zorgen dat de nieuwe schema-structuur wordt gebruikt door de applicatie. 