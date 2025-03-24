#!/bin/bash
# Script om de nieuwe synchronisatie migraties uit te voeren

# Variabelen voor database connectie
DB_HOST=${DB_HOST:-"dpg-cnbvvnf109ks73f9rvr0-a.oregon-postgres.render.com"}
DB_PORT=${DB_PORT:-"5432"}
DB_USER=${DB_USER:-"dklemailservice_user"}
DB_NAME=${DB_NAME:-"dklemailservice"}

# Wachtwoord uit omgevingsvariabele of vraag erom
if [ -z "$DB_PASSWORD" ]; then
  echo "Voer het database wachtwoord in:"
  read -s DB_PASSWORD
fi

echo "Start uitvoeren van synchronisatie migraties..."

# Synchroniseren van contact_formulieren
echo "Synchroniseren van contact_formulieren tabel..."
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f database/migrations/V1_8__sync_contact_formulieren.sql

# Synchroniseren van aanmeldingen
echo "Synchroniseren van aanmeldingen tabel..."
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f database/migrations/V1_9__sync_aanmeldingen.sql

# Repareren van antwoord tabellen
echo "Repareren van antwoord tabellen..."
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f database/migrations/V1_10__fix_antwoord_tables.sql

echo "Synchronisatie migraties succesvol uitgevoerd!"

# Toevoegen van migratie records in de migraties tabel
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
INSERT INTO migraties (versie, naam, toegepast) 
VALUES 
  ('1.8', 'Synchronisatie van contact_formulieren', CURRENT_TIMESTAMP),
  ('1.9', 'Synchronisatie van aanmeldingen', CURRENT_TIMESTAMP),
  ('1.10', 'Reparatie van antwoord tabellen', CURRENT_TIMESTAMP)
ON CONFLICT (versie) DO NOTHING;
"

echo "Migratie records toegevoegd aan de migraties tabel."
echo "Database schema is nu gesynchroniseerd met de Go models!" 