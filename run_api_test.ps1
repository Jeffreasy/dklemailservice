# PowerShell script om het test_api.ps1 script uit te voeren

Write-Host "DKL Email Service - API Test Runner" -ForegroundColor Cyan
Write-Host "======================================" -ForegroundColor Cyan
Write-Host ""

# Controleer of het test_api.ps1 script bestaat
if (-not (Test-Path -Path "test_api.ps1")) {
    Write-Host "Het test_api.ps1 script is niet gevonden in de huidige directory." -ForegroundColor Red
    Write-Host "Zorg ervoor dat je dit script uitvoert vanuit dezelfde directory als test_api.ps1." -ForegroundColor Yellow
    exit 1
}

Write-Host "Het test_api.ps1 script is gevonden." -ForegroundColor Green
Write-Host ""
Write-Host "Dit script zal het test_api.ps1 script uitvoeren om alle API endpoints te testen." -ForegroundColor Cyan
Write-Host "De volgende inloggegevens worden gebruikt:" -ForegroundColor Cyan
Write-Host "1. Admin account: admin@dekoninklijkeloop.nl met wachtwoord admin123" -ForegroundColor Yellow
Write-Host "2. Jeffrey account: jeffrey@dekoninklijkeloop.nl met wachtwoord DKL2025!" -ForegroundColor Yellow
Write-Host ""
Write-Host "Het script zal de volgende endpoints testen:" -ForegroundColor Cyan
Write-Host "- Root endpoint (/)" -ForegroundColor Magenta
Write-Host "- Health endpoint (/api/health)" -ForegroundColor Magenta
Write-Host "- Contact email endpoint (/api/contact-email)" -ForegroundColor Magenta
Write-Host "- Aanmelding email endpoint (/api/aanmelding-email)" -ForegroundColor Magenta
Write-Host "- Admin login endpoint (/api/auth/login)" -ForegroundColor Magenta
Write-Host "- Admin profile endpoint (/api/auth/profile)" -ForegroundColor Magenta
Write-Host "- Email metrics endpoint (/api/metrics/email)" -ForegroundColor Magenta
Write-Host "- Rate limit metrics endpoint (/api/metrics/rate-limits)" -ForegroundColor Magenta
Write-Host "- Reset password endpoint (/api/auth/reset-password)" -ForegroundColor Magenta
Write-Host "- Jeffrey login endpoint (/api/auth/login)" -ForegroundColor Magenta
Write-Host "- Jeffrey profile endpoint (/api/auth/profile)" -ForegroundColor Magenta
Write-Host "- Jeffrey logout endpoint (/api/auth/logout)" -ForegroundColor Magenta
Write-Host "- Admin logout endpoint (/api/auth/logout)" -ForegroundColor Magenta
Write-Host "- Prometheus metrics endpoint (/metrics)" -ForegroundColor Magenta
Write-Host ""

# Vraag om bevestiging
$confirmation = Read-Host "Wil je doorgaan met het uitvoeren van het test script? (j/n)"
if ($confirmation -ne "j") {
    Write-Host "Script geannuleerd." -ForegroundColor Yellow
    exit 0
}

Write-Host ""
Write-Host "Het test_api.ps1 script wordt nu uitgevoerd..." -ForegroundColor Cyan
Write-Host "======================================" -ForegroundColor Cyan
Write-Host ""

# Voer het test_api.ps1 script uit
try {
    & .\test_api.ps1
    if ($LASTEXITCODE -ne 0) {
        Write-Host ""
        Write-Host "Het test_api.ps1 script is beëindigd met foutcode: $LASTEXITCODE" -ForegroundColor Red
    }
} catch {
    Write-Host ""
    Write-Host "Er is een fout opgetreden bij het uitvoeren van het test_api.ps1 script:" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "Het test_api.ps1 script is voltooid." -ForegroundColor Green
Write-Host "Druk op een toets om af te sluiten..." -ForegroundColor Cyan
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown") 