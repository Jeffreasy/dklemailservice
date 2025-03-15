# DKL Email Service - Database Test Data Insertion Script
# Dit script voegt testgegevens toe aan de database voor het testen van de Contact en Aanmelding Beheer endpoints

# Configuratie
$dbConfig = @{
    Host     = if ($env:DB_HOST) { $env:DB_HOST } else { "dpg-cva4c01c1ekc738q6q0g-a" }
    Port     = if ($env:DB_PORT) { $env:DB_PORT } else { "5432" }
    Database = if ($env:DB_NAME) { $env:DB_NAME } else { "dekoninklijkeloopdatabase" }
    Username = if ($env:DB_USER) { $env:DB_USER } else { "dekoninklijkeloopdatabase_user" }
    Password = if ($env:DB_PASSWORD) { $env:DB_PASSWORD } else { "I4QP3JwyCcEbn8tGl6k3ErEvjUZ9V5rB" }
    SSLMode  = if ($env:DB_SSL_MODE) { $env:DB_SSL_MODE } else { "require" }
}

# Kleuren voor output
$successColor = "Green"
$errorColor = "Red"
$infoColor = "Cyan"
$promptColor = "Yellow"
$highlightColor = "Magenta"

# Functie om een titel te tonen
function Show-Title {
    param ([string]$Title)
    
    Write-Host "" -ForegroundColor $infoColor
    Write-Host "=============================================" -ForegroundColor $infoColor
    Write-Host " $Title" -ForegroundColor $infoColor
    Write-Host "=============================================" -ForegroundColor $infoColor
}

# Functie om een UUID te genereren
function New-UUID {
    return [guid]::NewGuid().ToString()
}

# Functie om de huidige tijd in ISO 8601 formaat te krijgen
function Get-ISOTime {
    return (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ss.fffZ")
}

# Functie om testgegevens direct in te voegen met SQL queries
function Insert-TestData {
    Show-Title -Title "Testgegevens invoegen met SQL queries"
    
    # Maak een tijdelijk SQL bestand
    $sqlFile = "insert_test_data.sql"
    
    # Genereer UUIDs en timestamps
    $contactIds = @()
    $aanmeldingIds = @()
    $now = Get-ISOTime
    
    for ($i = 1; $i -le 3; $i++) {
        $contactIds += New-UUID
        $aanmeldingIds += New-UUID
    }
    
    $antwoordContactId = New-UUID
    $antwoordAanmeldingId = New-UUID
    
    # Maak SQL queries
    $sql = @"
-- Contactformulieren toevoegen
INSERT INTO contact_formulieren (
    id, created_at, updated_at, naam, email, bericht, 
    email_verzonden, privacy_akkoord, status
) VALUES 
    ('$($contactIds[0])', '$now', '$now', 'Test Contact 1', 'test1@example.com', 'Dit is een testbericht voor contactformulier 1.', true, true, 'nieuw'),
    ('$($contactIds[1])', '$now', '$now', 'Test Contact 2', 'test2@example.com', 'Dit is een testbericht voor contactformulier 2.', true, true, 'in_behandeling'),
    ('$($contactIds[2])', '$now', '$now', 'Test Contact 3', 'test3@example.com', 'Dit is een testbericht voor contactformulier 3.', true, true, 'beantwoord');

-- Antwoord toevoegen aan het laatste contactformulier
INSERT INTO contact_antwoorden (
    id, contact_id, tekst, verzonden_op, verzonden_door, email_verzonden
) VALUES 
    ('$antwoordContactId', '$($contactIds[2])', 'Dit is een testantwoord op contactformulier 3.', '$now', 'admin@dekoninklijkeloop.nl', true);

-- Update contactformulier met antwoord
UPDATE contact_formulieren 
SET beantwoord = true, 
    antwoord_tekst = 'Dit is een testantwoord op contactformulier 3.', 
    antwoord_datum = '$now', 
    antwoord_door = 'admin@dekoninklijkeloop.nl' 
WHERE id = '$($contactIds[2])';

-- Aanmeldingen toevoegen
INSERT INTO aanmeldingen (
    id, created_at, updated_at, naam, email, telefoon, 
    rol, afstand, ondersteuning, bijzonderheden, terms, 
    email_verzonden, status
) VALUES 
    ('$($aanmeldingIds[0])', '$now', '$now', 'Test Aanmelding 1', 'aanmelding1@example.com', '0612345678', 'deelnemer', '10km', 'geen', 'Geen bijzonderheden', true, true, 'nieuw'),
    ('$($aanmeldingIds[1])', '$now', '$now', 'Test Aanmelding 2', 'aanmelding2@example.com', '0687654321', 'vrijwilliger', '', 'geen', 'Geen bijzonderheden', true, true, 'in_behandeling'),
    ('$($aanmeldingIds[2])', '$now', '$now', 'Test Aanmelding 3', 'aanmelding3@example.com', '0611223344', 'sponsor', '', 'geen', 'Geen bijzonderheden', true, true, 'beantwoord');

-- Antwoord toevoegen aan het laatste aanmelding
INSERT INTO aanmelding_antwoorden (
    id, aanmelding_id, tekst, verzonden_op, verzonden_door, email_verzonden
) VALUES 
    ('$antwoordAanmeldingId', '$($aanmeldingIds[2])', 'Dit is een testantwoord op aanmelding 3.', '$now', 'admin@dekoninklijkeloop.nl', true);
"@
    
    # Schrijf SQL naar bestand
    $sql | Out-File -FilePath $sqlFile -Encoding utf8
    
    # Toon instructies voor het uitvoeren van de SQL
    Write-Host "SQL bestand aangemaakt: $sqlFile" -ForegroundColor $successColor
    Write-Host "" -ForegroundColor $infoColor
    Write-Host "Om de testgegevens in te voegen, voer de volgende commando's uit:" -ForegroundColor $highlightColor
    Write-Host "" -ForegroundColor $infoColor
    Write-Host "Met psql (als je PostgreSQL client tools hebt):" -ForegroundColor $infoColor
    Write-Host "psql -h $($dbConfig.Host) -p $($dbConfig.Port) -d $($dbConfig.Database) -U $($dbConfig.Username) -f $sqlFile" -ForegroundColor $promptColor
    Write-Host "" -ForegroundColor $infoColor
    Write-Host "Of gebruik een PostgreSQL beheertools zoals pgAdmin om de SQL uit te voeren." -ForegroundColor $infoColor
    Write-Host "" -ForegroundColor $infoColor
    Write-Host "Nadat je de testgegevens hebt ingevoegd, kun je de API tests uitvoeren met het run_api_test.ps1 script." -ForegroundColor $highlightColor
    
    # Toon samenvatting van de ingevoegde gegevens
    Show-Title -Title "Samenvatting van testgegevens"
    Write-Host "Contactformulieren:" -ForegroundColor $successColor
    for ($i = 0; $i -lt $contactIds.Count; $i++) {
        $status = if ($i -eq 0) { "nieuw" } elseif ($i -eq 1) { "in_behandeling" } else { "beantwoord" }
        Write-Host "  - ID: $($contactIds[$i]) (Status: $status)" -ForegroundColor $infoColor
    }
    
    Write-Host "" -ForegroundColor $infoColor
    Write-Host "Aanmeldingen:" -ForegroundColor $successColor
    $rollen = @("deelnemer", "vrijwilliger", "sponsor")
    for ($i = 0; $i -lt $aanmeldingIds.Count; $i++) {
        $status = if ($i -eq 0) { "nieuw" } elseif ($i -eq 1) { "in_behandeling" } else { "beantwoord" }
        Write-Host "  - ID: $($aanmeldingIds[$i]) (Rol: $($rollen[$i]), Status: $status)" -ForegroundColor $infoColor
    }
    
    return @{
        ContactIds = $contactIds
        AanmeldingIds = $aanmeldingIds
        SqlFile = $sqlFile
    }
}

# Hoofdfunctie
function Main {
    Show-Title -Title "DKL Email Service - Database Test Data Insertion"
    
    # Voeg testgegevens in
    $result = Insert-TestData
    
    # Toon instructies voor het uitvoeren van de API tests
    Write-Host "" -ForegroundColor $infoColor
    Write-Host "Na het invoegen van de testgegevens kun je de API tests uitvoeren met:" -ForegroundColor $highlightColor
    Write-Host ".\run_api_test.ps1" -ForegroundColor $promptColor
}

# Script uitvoeren
Main 