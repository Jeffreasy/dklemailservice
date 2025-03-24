# PowerShell script om de nieuwe synchronisatie migraties uit te voeren

# Variabelen voor database connectie
$DB_HOST = if ($env:DB_HOST) { $env:DB_HOST } else { "dpg-cnbvvnf109ks73f9rvr0-a.oregon-postgres.render.com" }
$DB_PORT = if ($env:DB_PORT) { $env:DB_PORT } else { "5432" }
$DB_USER = if ($env:DB_USER) { $env:DB_USER } else { "dklemailservice_user" }
$DB_NAME = if ($env:DB_NAME) { $env:DB_NAME } else { "dklemailservice" }

# Wachtwoord uit omgevingsvariabele of vraag erom
if (-not $env:DB_PASSWORD) {
  $SecurePassword = Read-Host -Prompt "Voer het database wachtwoord in" -AsSecureString
  $BSTR = [System.Runtime.InteropServices.Marshal]::SecureStringToBSTR($SecurePassword)
  $DB_PASSWORD = [System.Runtime.InteropServices.Marshal]::PtrToStringAuto($BSTR)
} else {
  $DB_PASSWORD = $env:DB_PASSWORD
}

Write-Host "Start uitvoeren van synchronisatie migraties..." -ForegroundColor Cyan

# Controleer of psql beschikbaar is
try {
  $psqlVersion = & psql --version
  Write-Host "PostgreSQL client beschikbaar: $psqlVersion" -ForegroundColor Green
} catch {
  Write-Host "PostgreSQL client (psql) niet gevonden. Installeer PostgreSQL command-line tools." -ForegroundColor Red
  exit 1
}

# Synchroniseren van contact_formulieren
Write-Host "Synchroniseren van contact_formulieren tabel..." -ForegroundColor Cyan
$env:PGPASSWORD = $DB_PASSWORD
& psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f database/migrations/V1_8__sync_contact_formulieren.sql

# Synchroniseren van aanmeldingen
Write-Host "Synchroniseren van aanmeldingen tabel..." -ForegroundColor Cyan
& psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f database/migrations/V1_9__sync_aanmeldingen.sql

# Repareren van antwoord tabellen
Write-Host "Repareren van antwoord tabellen..." -ForegroundColor Cyan
& psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f database/migrations/V1_10__fix_antwoord_tables.sql

Write-Host "Synchronisatie migraties succesvol uitgevoerd!" -ForegroundColor Green

# Toevoegen van migratie records in de migraties tabel
$migratieQuery = @"
INSERT INTO migraties (versie, naam, toegepast) 
VALUES 
  ('1.8', 'Synchronisatie van contact_formulieren', CURRENT_TIMESTAMP),
  ('1.9', 'Synchronisatie van aanmeldingen', CURRENT_TIMESTAMP),
  ('1.10', 'Reparatie van antwoord tabellen', CURRENT_TIMESTAMP)
ON CONFLICT (versie) DO NOTHING;
"@

& psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c $migratieQuery

Write-Host "Migratie records toegevoegd aan de migraties tabel." -ForegroundColor Green
Write-Host "Database schema is nu gesynchroniseerd met de Go models!" -ForegroundColor Green

# Maak het wachtwoord leeg
$env:PGPASSWORD = $null 