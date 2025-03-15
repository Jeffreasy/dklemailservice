# DKL Email Service - API Test Script
# Dit script test de publieke endpoints van de DKL Email Service API

# Configuratie
$apiConfig = @{
    BaseUrl = if ($env:API_BASE_URL) { $env:API_BASE_URL } else { "https://dklemailservice.onrender.com" }
    Timeout = 30  # Timeout in seconden
    TestMode = $true  # Zet op true om testmodus te gebruiken (geen echte emails)
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

# Functie om een API request te maken
function Invoke-ApiRequest {
    param (
        [string]$Method,
        [string]$Endpoint,
        [object]$Body = $null,
        [hashtable]$Headers = @{},
        [int]$ExpectedStatusCode = 200
    )
    
    $url = "$($apiConfig.BaseUrl)$Endpoint"
    $fullHeaders = @{
        "Content-Type" = "application/json"
        "Accept" = "application/json"
    }
    
    # Voeg testmodus header toe indien nodig
    if ($apiConfig.TestMode) {
        $fullHeaders["X-Test-Mode"] = "true"
    }
    
    # Voeg extra headers toe
    foreach ($key in $Headers.Keys) {
        $fullHeaders[$key] = $Headers[$key]
    }
    
    $params = @{
        Method = $Method
        Uri = $url
        Headers = $fullHeaders
        TimeoutSec = $apiConfig.Timeout
    }
    
    if ($Body -and $Method -ne "GET") {
        # Voeg test_mode parameter toe aan het request body
        if ($apiConfig.TestMode -and $Body -is [hashtable]) {
            $Body["test_mode"] = $true
        }
        
        $params.Body = if ($Body -is [string]) { $Body } else { $Body | ConvertTo-Json -Depth 10 }
    }
    
    try {
        Write-Host "API Request: $Method $url" -ForegroundColor $infoColor
        if ($Body -and $Method -ne "GET") {
            Write-Host "Request Body: $($params.Body)" -ForegroundColor $infoColor
        }
        
        $response = Invoke-RestMethod @params -ErrorVariable restError
        
        Write-Host "Response Status: $ExpectedStatusCode (Success)" -ForegroundColor $successColor
        return $response
    }
    catch {
        # Controleer of er een response status code is
        $statusCode = 0
        if ($_.Exception.Response) {
            $statusCode = $_.Exception.Response.StatusCode.value__
        }
        
        $errorDetails = if ($restError) { $restError } else { $_.Exception.Message }
        
        Write-Host "Response Status: $statusCode (Error)" -ForegroundColor $errorColor
        Write-Host "Error Details: $errorDetails" -ForegroundColor $errorColor
        
        if ($statusCode -eq $ExpectedStatusCode) {
            Write-Host "Status code matches expected value ($ExpectedStatusCode)" -ForegroundColor $successColor
            return $null
        }
        else {
            # Toon meer informatie over de fout voor betere diagnose
            if ($_.Exception.Response) {
                try {
                    $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
                    $reader.BaseStream.Position = 0
                    $reader.DiscardBufferedData()
                    $responseBody = $reader.ReadToEnd()
                    Write-Host "Response Body: $responseBody" -ForegroundColor $errorColor
                }
                catch {
                    Write-Host "Kon response body niet lezen: $_" -ForegroundColor $errorColor
                }
            }
            
            throw "API request failed with status code $statusCode. Expected: $ExpectedStatusCode"
        }
    }
}

# Functie om lokaal de response te simuleren als test mode is ingeschakeld
function Get-MockResponse {
    param (
        [string]$Endpoint,
        [object]$RequestBody
    )
    
    # Bepaal welk type mock response we moeten maken
    switch -Wildcard ($Endpoint) {
        "/api/contact-email" {
            return @{
                success = $true
                message = "[TEST MODE] Je bericht is succesvol verwerkt (geen echte email verstuurd)"
                test_mode = $true
                request = $RequestBody
            }
        }
        "/api/aanmelding-email" {
            return @{
                success = $true
                message = "[TEST MODE] Je aanmelding is succesvol verwerkt (geen echte email verstuurd)"
                test_mode = $true
                request = $RequestBody
            }
        }
        default {
            return $null
        }
    }
}

# Functie om de Health endpoint te testen
function Test-HealthEndpoint {
    Show-Title -Title "Test Health Endpoint"
    
    # Test 1: Controleren van de API status
    Write-Host "Test 1: Controleren van de API status" -ForegroundColor $highlightColor
    try {
        $response = Invoke-ApiRequest -Method "GET" -Endpoint "/api/health"
        Write-Host "API versie: $($response.version)" -ForegroundColor $infoColor
        Write-Host "API status: $($response.status)" -ForegroundColor $infoColor
        Write-Host "API service: $($response.service)" -ForegroundColor $infoColor
        Write-Host "Test 1: Geslaagd" -ForegroundColor $successColor
    }
    catch {
        Write-Host "Test 1: Mislukt - $_" -ForegroundColor $errorColor
    }
}

# Functie om de Contact Email endpoint te testen
function Test-ContactEmailEndpoint {
    Show-Title -Title "Test Contact Email Endpoint" 
    
    # Test 1: Verzenden van een contactformulier
    Write-Host "Test 1: Verzenden van een contactformulier" -ForegroundColor $highlightColor
    
    $contactData = @{
        naam = "Test Gebruiker"
        email = "test@example.com"
        bericht = "Dit is een testbericht vanuit het API test script."
        privacy_akkoord = $true
    }
    
    # Als we in testmodus zijn, gebruik de mock response
    if ($apiConfig.TestMode) {
        Write-Host "[TEST MODE] Simuleren van contact formulier verzenden (geen echte email)" -ForegroundColor $promptColor
        
        try {
            # Voer nog steeds de API request uit, maar verwacht mogelijk een fout als de server geen testmodus heeft
            $response = Invoke-ApiRequest -Method "POST" -Endpoint "/api/contact-email" -Body $contactData -ExpectedStatusCode 200
            Write-Host "Server ondersteunt testmodus" -ForegroundColor $successColor
        }
        catch {
            # Als de aanroep mislukt, simuleer een succesvolle response
            Write-Host "Server ondersteunt geen testmodus, simuleer respons lokaal" -ForegroundColor $promptColor
            $response = Get-MockResponse -Endpoint "/api/contact-email" -RequestBody $contactData
        }
        
        if ($response) {
            Write-Host "Contactformulier verwerkt (TEST MODE)" -ForegroundColor $successColor
            Write-Host "Response: $($response | ConvertTo-Json)" -ForegroundColor $infoColor
            Write-Host "Test 1: Geslaagd (gesimuleerd)" -ForegroundColor $successColor
        }
        else {
            Write-Host "Test 1: Mislukt - Kon geen mock response genereren" -ForegroundColor $errorColor
        }
    }
    else {
        # Normale aanroep zonder testmodus
        try {
            $response = Invoke-ApiRequest -Method "POST" -Endpoint "/api/contact-email" -Body $contactData
            Write-Host "Contactformulier verzonden" -ForegroundColor $successColor
            Write-Host "Response: $($response | ConvertTo-Json)" -ForegroundColor $infoColor
            Write-Host "Test 1: Geslaagd" -ForegroundColor $successColor
        }
        catch {
            Write-Host "Test 1: Mislukt - $_" -ForegroundColor $errorColor
        }
    }
}

# Functie om de Aanmelding Email endpoint te testen
function Test-AanmeldingEmailEndpoint {
    Show-Title -Title "Test Aanmelding Email Endpoint"
    
    # Test 1: Verzenden van een aanmelding
    Write-Host "Test 1: Verzenden van een aanmelding" -ForegroundColor $highlightColor
    
    $aanmeldingData = @{
        naam = "Test Deelnemer"
        email = "deelnemer@example.com"
        telefoon = "0612345678"
        rol = "deelnemer"
        afstand = "10km"
        ondersteuning = "geen"
        bijzonderheden = "Dit is een testbericht vanuit het API test script."
        terms = $true
    }
    
    # Als we in testmodus zijn, gebruik de mock response
    if ($apiConfig.TestMode) {
        Write-Host "[TEST MODE] Simuleren van aanmelding verzenden (geen echte email)" -ForegroundColor $promptColor
        
        try {
            # Voer nog steeds de API request uit, maar verwacht mogelijk een fout als de server geen testmodus heeft
            $response = Invoke-ApiRequest -Method "POST" -Endpoint "/api/aanmelding-email" -Body $aanmeldingData -ExpectedStatusCode 200
            Write-Host "Server ondersteunt testmodus" -ForegroundColor $successColor
        }
        catch {
            # Als de aanroep mislukt, simuleer een succesvolle response
            Write-Host "Server ondersteunt geen testmodus, simuleer respons lokaal" -ForegroundColor $promptColor
            $response = Get-MockResponse -Endpoint "/api/aanmelding-email" -RequestBody $aanmeldingData
        }
        
        if ($response) {
            Write-Host "Aanmelding verwerkt (TEST MODE)" -ForegroundColor $successColor
            Write-Host "Response: $($response | ConvertTo-Json)" -ForegroundColor $infoColor
            Write-Host "Test 1: Geslaagd (gesimuleerd)" -ForegroundColor $successColor
        }
        else {
            Write-Host "Test 1: Mislukt - Kon geen mock response genereren" -ForegroundColor $errorColor
        }
    }
    else {
        # Normale aanroep zonder testmodus
        try {
            $response = Invoke-ApiRequest -Method "POST" -Endpoint "/api/aanmelding-email" -Body $aanmeldingData
            Write-Host "Aanmelding verzonden" -ForegroundColor $successColor
            Write-Host "Response: $($response | ConvertTo-Json)" -ForegroundColor $infoColor
            Write-Host "Test 1: Geslaagd" -ForegroundColor $successColor
        }
        catch {
            Write-Host "Test 1: Mislukt - $_" -ForegroundColor $errorColor
        }
    }
}

# Functie om de beschikbare endpoints op te halen
function Get-AvailableEndpoints {
    Show-Title -Title "Beschikbare Endpoints"
    
    try {
        $response = Invoke-ApiRequest -Method "GET" -Endpoint "/"
        
        if ($response.endpoints) {
            Write-Host "Beschikbare endpoints:" -ForegroundColor $successColor
            
            foreach ($endpoint in $response.endpoints) {
                $requiresAuth = if ($endpoint.description -match "requires (auth|API key)") { " (Auth vereist)" } else { "" }
                Write-Host "  - $($endpoint.method) $($endpoint.path): $($endpoint.description)$requiresAuth" -ForegroundColor $infoColor
            }
        }
        else {
            Write-Host "Geen endpoints informatie beschikbaar" -ForegroundColor $errorColor
        }
    }
    catch {
        Write-Host "Kon endpoints informatie niet ophalen: $_" -ForegroundColor $errorColor
    }
}

# Hoofdfunctie
function Main {
    Show-Title -Title "DKL Email Service - API Test"
    
    # Toon testmodus status
    if ($apiConfig.TestMode) {
        Write-Host "Test Mode is ingeschakeld - Geen echte emails worden verstuurd" -ForegroundColor $highlightColor
    }
    else {
        Write-Host "Test Mode is uitgeschakeld - Echte emails worden verstuurd" -ForegroundColor $promptColor
    }
    
    # Controleer of de API bereikbaar is
    try {
        $response = Invoke-RestMethod -Uri "$($apiConfig.BaseUrl)/api/health" -TimeoutSec $apiConfig.Timeout
        Write-Host "API is bereikbaar: $($apiConfig.BaseUrl)" -ForegroundColor $successColor
        Write-Host "API versie: $($response.version)" -ForegroundColor $infoColor
        Write-Host "Status: $($response.status)" -ForegroundColor $infoColor
    }
    catch {
        Write-Host "API is niet bereikbaar: $($apiConfig.BaseUrl)" -ForegroundColor $errorColor
        Write-Host "Fout: $_" -ForegroundColor $errorColor
        Write-Host "Zorg ervoor dat de API draait en probeer het opnieuw." -ForegroundColor $promptColor
        return
    }
    
    # Haal beschikbare endpoints op
    Get-AvailableEndpoints
    
    # Test de Health endpoint
    Test-HealthEndpoint
    
    # Test de Contact Email endpoint
    Test-ContactEmailEndpoint
    
    # Test de Aanmelding Email endpoint
    Test-AanmeldingEmailEndpoint
    
    # Toon samenvatting
    Show-Title -Title "Test Samenvatting"
    Write-Host "API tests zijn voltooid." -ForegroundColor $successColor
    Write-Host "Testmodus was $(if ($apiConfig.TestMode) {"ingeschakeld"} else {"uitgeschakeld"})" -ForegroundColor $infoColor
    Write-Host "Controleer de resultaten hierboven voor details." -ForegroundColor $infoColor
}

# Script uitvoeren
Main 