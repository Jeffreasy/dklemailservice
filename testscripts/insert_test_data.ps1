# Insert Test Data script
# Dit script voert de insert_testdata.sql statements uit op de database

param (
    [switch]$Force
)

# Database configuratie - standaardwaarden uit README gebruiken
$DB_HOST = $env:DB_HOST -or "dpg-cva4c01c1ekc738q6q0g-a.oregon-postgres.render.com"
$DB_PORT = $env:DB_PORT -or "5432"
$DB_NAME = $env:DB_NAME -or "dekoninklijkeloopdatabase"
$DB_USER = $env:DB_USER -or "dekoninklijkeloopdatabase_user"
$DB_PASSWORD = $env:DB_PASSWORD -or "I4QP3JwyCcEbn8tGl6k3ErEvjUZ9V5rB"
$DB_SSL_MODE = $env:DB_SSL_MODE -or "require"

Write-Host "=============================================" -ForegroundColor Cyan
Write-Host " DKL Email Service - Test Data Invoegen" -ForegroundColor Cyan
Write-Host "=============================================" -ForegroundColor Cyan
Write-Host "Database: $DB_HOST / $DB_NAME"

if (-not $Force) {
    Write-Host ""
    Write-Host "WAARSCHUWING: Dit script voegt testgegevens in de database in." -ForegroundColor Yellow
    Write-Host "Dit kan bestaande gegevens beïnvloeden of duplicaten veroorzaken." -ForegroundColor Yellow
    Write-Host ""
    $confirm = Read-Host "Weet je zeker dat je door wilt gaan? (j/n)"
    
    if ($confirm -ne "j") {
        Write-Host "Operatie geannuleerd." -ForegroundColor Red
        exit
    }
}

# Controleer of het SQL bestand bestaat
$sqlFile = "..\insert_testdata.sql"
if (-not (Test-Path $sqlFile)) {
    Write-Host "SQL bestand niet gevonden: $sqlFile" -ForegroundColor Red
    exit
}

# PSQL uitvoeren met het SQL bestand
try {
    # Maak een connection string
    $env:PGPASSWORD = $DB_PASSWORD
    
    Write-Host "SQL statements worden uitgevoerd..." -ForegroundColor Cyan
    
    # Voer het commando uit
    # Gebruik PowerShell om het command uit te voeren, zodat we de uitvoer kunnen opvangen
    $command = "psql -h $DB_HOST -p $DB_PORT -d $DB_NAME -U $DB_USER -f $sqlFile"
    
    # Controleer of psql beschikbaar is
    try {
        $psqlVersion = (Invoke-Expression "psql --version") 2>&1
        Write-Host "PSQL versie: $psqlVersion" -ForegroundColor Green
    }
    catch {
        Write-Host "PSQL niet gevonden. Zorg ervoor dat PostgreSQL is geïnstalleerd en in het PATH staat." -ForegroundColor Red
        Write-Host "Je kunt het SQL bestand handmatig uitvoeren met een andere PostgreSQL client."
        
        # Toon de inhoud van het SQL bestand
        Write-Host ""
        Write-Host "SQL bestand inhoud:" -ForegroundColor Yellow
        Get-Content $sqlFile | ForEach-Object { Write-Host $_ }
        exit
    }
    
    $result = Invoke-Expression $command
    
    # Controleer het resultaat
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Testgegevens zijn succesvol ingevoegd!" -ForegroundColor Green
        
        # Toon aantal ingevoegde records
        Write-Host ""
        Write-Host "Ingevoegde gegevens:" -ForegroundColor Cyan
        Write-Host "- 2 contactformulieren" -ForegroundColor White
        Write-Host "- 11 aanmeldingen" -ForegroundColor White
    }
    else {
        Write-Host "Er is een fout opgetreden bij het uitvoeren van de SQL statements." -ForegroundColor Red
        Write-Host $result
    }
}
catch {
    Write-Host "Fout bij het uitvoeren van de SQL:" -ForegroundColor Red
    Write-Host $_.Exception.Message
}
finally {
    # Verwijder de password environment variable voor veiligheid
    $env:PGPASSWORD = ""
} 