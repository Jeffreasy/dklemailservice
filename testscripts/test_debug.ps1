# Vereenvoudigde versie om te debuggen
param(
    [switch]$UseLocalUrl,
    [switch]$DisableTestMode
)

# Kleuren voor output
$successColor = "Green"
$errorColor = "Red" 
$infoColor = "Cyan"
$promptColor = "Yellow"
$warningColor = "DarkYellow"

# Configuratie
$baseUrl = "https://dklemailservice.onrender.com"
$localUrl = "http://localhost:8080"
$testMode = -not $DisableTestMode

# API URL instellen
if ($UseLocalUrl) {
    $currentUrl = $localUrl
} else {
    $currentUrl = $baseUrl
}

Write-Host "Test script gestart" -ForegroundColor $infoColor
Write-Host "URL: $currentUrl" -ForegroundColor $infoColor
Write-Host "Test mode: $testMode" -ForegroundColor $infoColor

# Test de root endpoint
try {
    Write-Host "Root endpoint testen..." -ForegroundColor $infoColor
    $response = Invoke-RestMethod -Uri "$currentUrl/" -Method Get -UseBasicParsing
    
    Write-Host "Root endpoint succesvol" -ForegroundColor $successColor
    Write-Host "Service: $($response.service)" -ForegroundColor $infoColor
    Write-Host "Version: $($response.version)" -ForegroundColor $infoColor
} catch {
    Write-Host "Root endpoint test mislukt: $_" -ForegroundColor $errorColor
}

Write-Host "Test script voltooid" -ForegroundColor $successColor 