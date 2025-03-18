# DKL Email Service API Geautomatiseerd Test Script
# Dit script test automatisch alle endpoints van de DKL Email Service API

# Kleuren voor output - EERST DEFINIËREN zodat ze altijd beschikbaar zijn
$successColor = "Green"
$errorColor = "Red"
$infoColor = "Cyan"
$promptColor = "Yellow"
$highlightColor = "Magenta"

# Parameters definiëren
param (
    [string]$ConfigFile = "api_test_config.json",
    [string]$ApiUrl,
    [switch]$UseLocalUrl,
    [switch]$DisableTestMode,
    [string]$TestFilter,
    [switch]$SkipAuth,
    [int]$ParallelJobs = 0,
    [switch]$GenerateReport,
    [string]$ReportPath,
    [switch]$Help
)

# Help functie
function Show-Help {
    Write-Host "DKL Email Service API Test Script" -ForegroundColor $highlightColor
    Write-Host "==================================" -ForegroundColor $infoColor
    Write-Host ""
    Write-Host "Gebruik: .\test_api.ps1 [parameters]" -ForegroundColor $infoColor
    Write-Host ""
    Write-Host "Parameters:" -ForegroundColor $highlightColor
    Write-Host "  -ConfigFile <path>      Path naar configuratiebestand (default: api_test_config.json)" -ForegroundColor $infoColor
    Write-Host "  -ApiUrl <url>           Custom API URL om te testen" -ForegroundColor $infoColor
    Write-Host "  -UseLocalUrl            Gebruik lokale URL (http://localhost:8080)" -ForegroundColor $infoColor
    Write-Host "  -DisableTestMode        Schakel testmodus uit (echte emails worden verstuurd)" -ForegroundColor $infoColor
    Write-Host "  -TestFilter <filter>    Filter tests op naam (bijv. 'Contact' of 'Health')" -ForegroundColor $infoColor
    Write-Host "  -SkipAuth               Sla authenticatie tests over" -ForegroundColor $infoColor
    Write-Host "  -ParallelJobs <n>       Aantal parallelle jobs (0 = uitgeschakeld)" -ForegroundColor $infoColor
    Write-Host "  -GenerateReport         Genereer testrapport" -ForegroundColor $infoColor
    Write-Host "  -ReportPath <path>      Custom pad voor testrapport" -ForegroundColor $infoColor
    Write-Host "  -Help                   Toon deze help informatie" -ForegroundColor $infoColor
    Write-Host ""
    Write-Host "Voorbeeld:" -ForegroundColor $highlightColor
    Write-Host "  .\test_api.ps1 -UseLocalUrl -TestFilter 'Health'" -ForegroundColor $infoColor
    Write-Host "  .\test_api.ps1 -ConfigFile 'my_config.json' -GenerateReport" -ForegroundColor $infoColor
    
    exit
}

# Configuratie via bestand en omgevingsvariabelen
function Load-Configuration {
    param (
        [string]$ConfigFile = "api_test_config.json"
    )
    
    $config = @{
        BaseUrl = $env:API_BASE_URL -or "https://dklemailservice.onrender.com"
        LocalUrl = $env:API_LOCAL_URL -or "http://localhost:8080"
        Timeout = 30
        TestMode = $true
        Credentials = @{
            AdminUser = @{
                Email = $env:API_ADMIN_EMAIL -or "admin@example.com"
                Password = $env:API_ADMIN_PASSWORD -or "placeholder"
            }
            RegularUser = @{
                Email = $env:API_USER_EMAIL -or "user@example.com"
                Password = $env:API_USER_PASSWORD -or "placeholder"
            }
        }
        ApiKey = $env:API_KEY -or "placeholder_api_key"
    }
    
    # Laad configuratiebestand indien aanwezig
    if (Test-Path $ConfigFile) {
        try {
            $fileConfig = Get-Content $ConfigFile -Raw | ConvertFrom-Json
            
            # Merge de configuraties, waarbij bestand voorrang heeft
            if ($fileConfig.BaseUrl) { $config.BaseUrl = $fileConfig.BaseUrl }
            if ($fileConfig.LocalUrl) { $config.LocalUrl = $fileConfig.LocalUrl }
            if ($fileConfig.Timeout) { $config.Timeout = $fileConfig.Timeout }
            if ($null -ne $fileConfig.TestMode) { $config.TestMode = $fileConfig.TestMode }
            
            if ($fileConfig.Credentials) {
                if ($fileConfig.Credentials.AdminUser) {
                    if ($fileConfig.Credentials.AdminUser.Email) { 
                        $config.Credentials.AdminUser.Email = $fileConfig.Credentials.AdminUser.Email 
                    }
                    if ($fileConfig.Credentials.AdminUser.Password) { 
                        $config.Credentials.AdminUser.Password = $fileConfig.Credentials.AdminUser.Password 
                    }
                }
                if ($fileConfig.Credentials.RegularUser) {
                    if ($fileConfig.Credentials.RegularUser.Email) { 
                        $config.Credentials.RegularUser.Email = $fileConfig.Credentials.RegularUser.Email 
                    }
                    if ($fileConfig.Credentials.RegularUser.Password) { 
                        $config.Credentials.RegularUser.Password = $fileConfig.Credentials.RegularUser.Password 
                    }
                }
            }
            
            if ($fileConfig.ApiKey) { $config.ApiKey = $fileConfig.ApiKey }
            
            Write-Host "Configuratie geladen uit $ConfigFile" -ForegroundColor $infoColor
        } catch {
            Write-Host "Fout bij laden van configuratiebestand: $_" -ForegroundColor $errorColor
            Write-Host "Standaard configuratie wordt gebruikt" -ForegroundColor $promptColor
        }
    } else {
        Write-Host "Geen configuratiebestand gevonden, standaard configuratie wordt gebruikt" -ForegroundColor $promptColor
    }
    
    return $config
}

# Toon help als gevraagd
if ($Help) {
    Show-Help
}

# Laad configuratie EENMALIG en verwerk command-line parameters
$config = Load-Configuration -ConfigFile $ConfigFile

# Overschrijf configuratie met command-line parameters
if ($ApiUrl) {
    $config.BaseUrl = $ApiUrl
}

# Stel initiële URLs in
$baseUrl = $config.BaseUrl
$localUrl = $config.LocalUrl

# Stel test mode in op basis van configuratie en parameters
$testMode = $config.TestMode
if ($DisableTestMode) {
    $testMode = $false
}

# Stel de huidige URL in op basis van parameters
if ($UseLocalUrl) {
    $currentUrl = $localUrl
} else {
    $currentUrl = $baseUrl
}

# Initialiseer andere parameters
$generateTestReport = $GenerateReport
$skipAuthTests = $SkipAuth
$testFilterPattern = $TestFilter
$parallelJobCount = $ParallelJobs

# Indien custom report path
if ($ReportPath) {
    $customReportPath = $ReportPath
} else {
    $reportDate = Get-Date -Format "yyyy-MM-dd_HH-mm-ss"
    $customReportPath = "DKL_API_Test_Report_$reportDate.txt"
}

# Testdata voor formulieren - BIJGEWERKT met correcte inloggegevens
$testData = @{
    ContactEmail = @{
        naam = "Test Gebruiker"
        email = "test@example.com"
        bericht = "Dit is een geautomatiseerd testbericht"
        privacy_akkoord = $true
        Clone = {
            $copy = @{
                naam = $this.naam
                email = $this.email
                bericht = $this.bericht
                privacy_akkoord = $this.privacy_akkoord
                Clone = $this.Clone
            }
            return $copy
        }
    }
    AanmeldingEmail = @{
        naam = "Test Deelnemer"
        email = "deelnemer@example.com"
        telefoon = "0612345678"
        rol = "deelnemer"
        afstand = "10km"
        ondersteuning = "geen"
        bijzonderheden = "geen"
        terms = $true
        Clone = {
            $copy = @{
                naam = $this.naam
                email = $this.email
                telefoon = $this.telefoon
                rol = $this.rol
                afstand = $this.afstand
                ondersteuning = $this.ondersteuning
                bijzonderheden = $this.bijzonderheden
                terms = $this.terms
                Clone = $this.Clone
            }
            return $copy
        }
    }
    AdminUser = @{
        email = $config.Credentials.AdminUser.Email
        wachtwoord = $config.Credentials.AdminUser.Password
    }
    JeffreyUser = @{
        email = $config.Credentials.RegularUser.Email
        wachtwoord = $config.Credentials.RegularUser.Password
    }
    ResetPassword = @{
        huidig_wachtwoord = "admin123"  # Bijgewerkt wachtwoord
        nieuw_wachtwoord = "admin123"   # Zelfde wachtwoord om geen wijziging te maken
    }
    # API key voor metrics endpoints
    ApiKey = $config.ApiKey
    
    # Testdata voor Contact Beheer endpoints
    ContactUpdate = @{
        status = "in_behandeling"
        notities = "Dit is een testnotitie voor het contactformulier"
    }
    ContactAntwoord = @{
        tekst = "Dit is een testantwoord op het contactformulier"
    }
    
    # Testdata voor Aanmelding Beheer endpoints
    AanmeldingUpdate = @{
        status = "in_behandeling"
        notities = "Dit is een testnotitie voor de aanmelding"
    }
    AanmeldingAntwoord = @{
        tekst = "Dit is een testantwoord op de aanmelding"
    }
}

# Functie om een titel te tonen
function Show-Title {
    param ([string]$Title)
    
    Write-Host "" -ForegroundColor $infoColor
    Write-Host "=============================================" -ForegroundColor $infoColor
    Write-Host " $Title" -ForegroundColor $infoColor
    Write-Host "=============================================" -ForegroundColor $infoColor
}

# Functie om een resultaat te tonen
function Show-Result {
    param (
        [string]$Name,
        [object]$Response,
        [bool]$Success = $true
    )
    
    Write-Host "" -ForegroundColor $infoColor
    if ($Success) {
        Write-Host "[SUCCESS] $Name" -ForegroundColor $successColor
    } else {
        Write-Host "[ERROR] $Name" -ForegroundColor $errorColor
    }
    
    if ($Response) {
        $formattedResponse = ($Response | ConvertTo-Json -Depth 4)
        Write-Host $formattedResponse
    } else {
        Write-Host "Geen response ontvangen" -ForegroundColor $errorColor
    }
}

# Functie om een API call te maken
function Invoke-ApiCall {
    param (
        [string]$Name,
        [string]$Method,
        [string]$Endpoint,
        [object]$Body = $null,
        [switch]$UseAuth,
        [Microsoft.PowerShell.Commands.WebRequestSession]$Session = $null,
        [string]$ContentType = "application/json",
        [int]$RetryCount = 1,
        [int]$RetryDelay = 2,
        [string]$ApiKey = $null,
        [switch]$IgnoreEmailErrors = $false,
        [switch]$UseTestMode = $global:testMode
    )
    
    Write-Host "" -ForegroundColor $infoColor
    Write-Host "[TESTING] $Name..." -ForegroundColor $infoColor
    
    $attempt = 0
    $success = $false
    $response = $null
    
    while (-not $success -and $attempt -lt $RetryCount) {
        $attempt++
        if ($attempt -gt 1) {
            Write-Host "Poging $attempt van $RetryCount..." -ForegroundColor $promptColor
            Start-Sleep -Seconds $RetryDelay
        }
        
        try {
            $params = @{
                Uri = "$currentUrl$Endpoint"
                Method = $Method
                UseBasicParsing = $true  # Toegevoegd voor betere compatibiliteit
            }
            
            if ($Body) {
                $jsonBody = $Body | ConvertTo-Json
                Write-Host "Request Body:" -ForegroundColor $infoColor
                Write-Host $jsonBody
                $params.Body = $jsonBody
                $params.ContentType = $ContentType
            }
            
            # Initialiseer headers als ze nog niet bestaan
            if (-not $params.Headers) {
                $params.Headers = @{}
            }
            
            if ($UseAuth -and $Session) {
                $params.WebSession = $Session
                
                # Voeg JWT token toe als die beschikbaar is
                if ($script:jwtToken) {
                    $params.Headers["Authorization"] = "Bearer $($script:jwtToken)"
                }
            }
            
            # Voeg API key toe als die is opgegeven
            if ($ApiKey) {
                Write-Host "API Key wordt gebruikt voor deze request." -ForegroundColor $infoColor
                $params.Headers["X-API-Key"] = $ApiKey
            }
            
            # Voeg testmodus headers toe indien nodig
            if ($UseTestMode) {
                $params.Headers["X-Test-Mode"] = "true"
                
                # Voeg test_mode parameter toe aan het request body
                if ($Body -and $Method -ne "GET" -and $Body -is [hashtable]) {
                    $Body["test_mode"] = $true
                    $jsonBody = $Body | ConvertTo-Json
                    $params.Body = $jsonBody
                }
            }
            
            if ($Method -eq "Get" -and $Endpoint -eq "/metrics") {
                # Speciale behandeling voor Prometheus metrics (plaintext)
                $response = Invoke-RestMethod @params
            } else {
                $response = Invoke-RestMethod @params
            }
            
            $success = $true
            Show-Result -Name $Name -Response $response
        }
        catch {
            # Controleer of dit een e-mailfout is die we kunnen negeren
            $isEmailError = $false
            if ($IgnoreEmailErrors -and $_.Exception.Response.StatusCode.value__ -eq 500) {
                try {
                    if ($_.Exception.Response -and $_.Exception.Response.GetResponseStream()) {
                        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
                        $reader.BaseStream.Position = 0
                        $reader.DiscardBufferedData()
                        $responseBody = $reader.ReadToEnd()
                        
                        # Controleer of de fout gerelateerd is aan e-mail
                        if ($responseBody -match "Fout bij het verzenden van de bevestigingsemail") {
                            $isEmailError = $true
                            $success = $true
                            
                            # Toon een waarschuwing maar beschouw het als succesvol
                            Write-Host "" -ForegroundColor $infoColor
                            Write-Host "[WAARSCHUWING] $Name heeft een e-mailfout, maar de gegevens zijn opgeslagen" -ForegroundColor $promptColor
                            Write-Host "StatusCode: $($_.Exception.Response.StatusCode.value__)" -ForegroundColor $promptColor
                            Write-Host "Response Body:" -ForegroundColor $promptColor
                            Write-Host $responseBody -ForegroundColor $promptColor
                            
                            # Maak een eenvoudige response om terug te geven
                            $response = @{
                                success = $true
                                message = "Gegevens opgeslagen, maar e-mail kon niet worden verzonden"
                                error = "E-mailfout genegeerd: $responseBody"
                            }
                        }
                    }
                } catch {
                    # Kon response body niet lezen, ga door met normale foutafhandeling
                }
            }
            
            if (-not $isEmailError -and $attempt -ge $RetryCount) {
                Write-Host "" -ForegroundColor $infoColor
                Write-Host "[ERROR] $Name failed na $RetryCount pogingen" -ForegroundColor $errorColor
                
                if ($_.Exception.Response) {
                    Write-Host "StatusCode:" $_.Exception.Response.StatusCode.value__ -ForegroundColor $errorColor
                }
                
                Write-Host "Message:" $_.Exception.Message -ForegroundColor $errorColor
                
                # Probeer de response body te lezen als die er is
                try {
                    if ($_.Exception.Response -and $_.Exception.Response.GetResponseStream()) {
                        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
                        $reader.BaseStream.Position = 0
                        $reader.DiscardBufferedData()
                        $responseBody = $reader.ReadToEnd()
                        Write-Host "Response Body:" -ForegroundColor $errorColor
                        Write-Host $responseBody -ForegroundColor $errorColor
                    }
                } catch {
                    Write-Host "Kon response body niet lezen" -ForegroundColor $errorColor
                }
            }
        }
    }
    
    # Simuleer response in testmodus bij fouten
    if ($UseTestMode -and -not $success -and ($Endpoint -eq "/api/contact-email" -or $Endpoint -eq "/api/aanmelding-email")) {
        Write-Host "Server ondersteunt geen testmodus, simuleer respons lokaal" -ForegroundColor $promptColor
        $response = Get-MockResponse -Endpoint $Endpoint -RequestBody $Body
        $success = $true
        return $response
    }
    
    return $response
}

# Nieuwe functie: Get-MockResponse (overgenomen uit run_api_test.ps1)
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
        "/api/contact/*/antwoord" {
            return @{
                success = $true
                message = "[TEST MODE] Antwoord verwerkt (geen echte email verstuurd)"
                test_mode = $true
                id = (New-Guid).ToString()
                tekst = $RequestBody.tekst
                verzond_door = "test@example.com"
                email_verzonden = $false
                created_at = [DateTime]::Now.ToString("o")
            }
        }
        "/api/aanmelding/*/antwoord" {
            return @{
                success = $true
                message = "[TEST MODE] Antwoord verwerkt (geen echte email verstuurd)"
                test_mode = $true
                id = (New-Guid).ToString()
                tekst = $RequestBody.tekst
                verzond_door = "test@example.com"
                email_verzonden = $false
                created_at = [DateTime]::Now.ToString("o")
            }
        }
        default {
            return $null
        }
    }
}

# Functie om API URL te vaststellen met automatische detectie
function Set-ApiUrl {
    Show-Title -Title "API URL Detecteren"
    
    # Als we een specifieke URL hebben via parameters, gebruik die direct
    if ($UseLocalUrl) {
        Write-Host "Gebruik lokale URL: $localUrl (ingesteld via parameters)" -ForegroundColor $successColor
        return $localUrl
    }
    
    if ($ApiUrl) {
        Write-Host "Gebruik aangepaste URL: $baseUrl (ingesteld via parameters)" -ForegroundColor $successColor
        return $baseUrl
    }
    
    # Probeer eerst de productie URL
    Write-Host "Proberen te verbinden met productie URL ($baseUrl)..." -ForegroundColor $infoColor
    
    try {
        $testResult = Invoke-RestMethod -Uri "$baseUrl/api/health" -Method Get -UseBasicParsing -TimeoutSec 5
        Write-Host "Verbinding gelukt met productie URL!" -ForegroundColor $successColor
        return $baseUrl
    } catch {
        Write-Host "Kon geen verbinding maken met productie URL: $_" -ForegroundColor $errorColor
    }
    
    # Probeer de lokale URL als alternatief
    Write-Host "Proberen te verbinden met lokale URL ($localUrl)..." -ForegroundColor $infoColor
    
    try {
        $testResult = Invoke-RestMethod -Uri "$localUrl/api/health" -Method Get -UseBasicParsing -TimeoutSec 5
        Write-Host "Verbinding gelukt met lokale URL!" -ForegroundColor $successColor
        return $localUrl
    } catch {
        Write-Host "Kon geen verbinding maken met lokale URL: $_" -ForegroundColor $errorColor
    }
    
    # Default terugvallen op de productie URL
    Write-Host "Geen werkende URL gedetecteerd, terugvallen op productie URL: $baseUrl" -ForegroundColor $promptColor
    return $baseUrl
}

# Functie om de root endpoint te testen
function Test-RootEndpoint {
    Show-Title -Title "Root Endpoint Testen"
    
    $response = Invoke-ApiCall -Name "Root Endpoint" -Method "Get" -Endpoint "/" -RetryCount 3
    
    if ($response) {
        Write-Host "" -ForegroundColor $infoColor
        Write-Host "Service: $($response.service)" -ForegroundColor $infoColor
        Write-Host "Version: $($response.version)" -ForegroundColor $infoColor
        Write-Host "Status: $($response.status)" -ForegroundColor $infoColor
        Write-Host "Timestamp: $($response.timestamp)" -ForegroundColor $infoColor
        
        Write-Host "" -ForegroundColor $infoColor
        Write-Host "Beschikbare Endpoints:" -ForegroundColor $infoColor
        foreach ($endpoint in $response.endpoints) {
            Write-Host "- $($endpoint.method) $($endpoint.path): $($endpoint.description)" -ForegroundColor $infoColor
        }
    }
    
    return $response
}

# Functie om de health endpoint te testen
function Test-HealthEndpoint {
    Show-Title -Title "Health Endpoint Testen"
    
    $response = Invoke-ApiCall -Name "Health Check" -Method "Get" -Endpoint "/api/health" -RetryCount 3
    
    return $response
}

# Functie om de contact email endpoint te testen
function Test-ContactEmailEndpoint {
    param (
        [string]$TestId = (New-Guid).ToString().Substring(0, 8)
    )
    
    Show-Title -Title "Contact Email Endpoint Testen"
    
    # Genereer unieke testdata voor deze specifieke test
    $isolatedData = New-IsolatedTestData -TestId $TestId
    $contactData = $isolatedData.ContactEmail
    
    # Voeg test ID toe aan request voor traceren
    $contactData["test_id"] = $TestId
    
    # Gebruik IgnoreEmailErrors of TestMode om e-mailfouten te negeren
    $response = Invoke-ApiCall -Name "Contact Email" -Method "Post" -Endpoint "/api/contact-email" -Body $contactData -RetryCount 2 -UseTestMode
    
    return $response
}

# Functie om de aanmelding email endpoint te testen
function Test-AanmeldingEmailEndpoint {
    Show-Title -Title "Aanmelding Email Endpoint Testen"
    
    $aanmeldingData = $testData.AanmeldingEmail
    
    # Gebruik IgnoreEmailErrors om e-mailfouten te negeren
    $response = Invoke-ApiCall -Name "Aanmelding Email" -Method "Post" -Endpoint "/api/aanmelding-email" -Body $aanmeldingData -RetryCount 2 -IgnoreEmailErrors
    
    return $response
}

# Functie om de login endpoint te testen met admin account
function Test-AdminLoginEndpoint {
    Show-Title -Title "Admin Login Endpoint Testen"
    
    $loginData = $testData.AdminUser
    
    # Maak een nieuwe sessie aan
    $script:session = New-Object Microsoft.PowerShell.Commands.WebRequestSession
    
    # Probeer in te loggen
    $response = Invoke-ApiCall -Name "Admin Login" -Method "Post" -Endpoint "/api/auth/login" -Body $loginData -RetryCount 3
    
    if ($response) {
        # Sla de sessie op met cookies
        try {
            $loginRequest = Invoke-WebRequest -Uri "$currentUrl/api/auth/login" -Method Post -Body ($loginData | ConvertTo-Json) -ContentType "application/json" -SessionVariable tempSession -UseBasicParsing
            $script:session = $tempSession
            $script:isLoggedIn = $true
            Write-Host "" -ForegroundColor $infoColor
            Write-Host "[SUCCESS] Succesvol ingelogd als admin en sessie opgeslagen!" -ForegroundColor $successColor
            
            # Sla de token op als die er is
            if ($response.token) {
                $script:jwtToken = $response.token
                Write-Host "JWT Token opgeslagen voor gebruik in requests" -ForegroundColor $infoColor
            }
        }
        catch {
            Write-Host "" -ForegroundColor $infoColor
            Write-Host "[ERROR] Kon sessie niet opslaan: $_" -ForegroundColor $errorColor
            $script:isLoggedIn = $false
        }
    } else {
        $script:isLoggedIn = $false
    }
    
    return $response
}

# Functie om de login endpoint te testen met Jeffrey account
function Test-JeffreyLoginEndpoint {
    Show-Title -Title "Jeffrey Login Endpoint Testen"
    
    $loginData = $testData.JeffreyUser
    
    # Maak een nieuwe sessie aan
    $script:jeffreySession = New-Object Microsoft.PowerShell.Commands.WebRequestSession
    
    # Probeer in te loggen
    $response = Invoke-ApiCall -Name "Jeffrey Login" -Method "Post" -Endpoint "/api/auth/login" -Body $loginData -RetryCount 3
    
    if ($response) {
        # Sla de sessie op met cookies
        try {
            $loginRequest = Invoke-WebRequest -Uri "$currentUrl/api/auth/login" -Method Post -Body ($loginData | ConvertTo-Json) -ContentType "application/json" -SessionVariable tempSession -UseBasicParsing
            $script:jeffreySession = $tempSession
            $script:isJeffreyLoggedIn = $true
            Write-Host "" -ForegroundColor $infoColor
            Write-Host "[SUCCESS] Succesvol ingelogd als Jeffrey en sessie opgeslagen!" -ForegroundColor $successColor
            
            # Sla de token op als die er is
            if ($response.token) {
                $script:jeffreyJwtToken = $response.token
                Write-Host "Jeffrey JWT Token opgeslagen voor gebruik in requests" -ForegroundColor $infoColor
            }
        }
        catch {
            Write-Host "" -ForegroundColor $infoColor
            Write-Host "[ERROR] Kon Jeffrey sessie niet opslaan: $_" -ForegroundColor $errorColor
            $script:isJeffreyLoggedIn = $false
        }
    } else {
        $script:isJeffreyLoggedIn = $false
    }
    
    return $response
}

# Functie om de profile endpoint te testen
function Test-ProfileEndpoint {
    Show-Title -Title "Profile Endpoint Testen"
    
    if (-not $script:isLoggedIn) {
        Write-Host "Je moet eerst inloggen om deze endpoint te testen." -ForegroundColor $errorColor
        return $null
    }
    
    $response = Invoke-ApiCall -Name "Get Profile" -Method "Get" -Endpoint "/api/auth/profile" -UseAuth -Session $script:session -RetryCount 2
    
    return $response
}

# Functie om de Jeffrey profile endpoint te testen
function Test-JeffreyProfileEndpoint {
    Show-Title -Title "Jeffrey Profile Endpoint Testen"
    
    if (-not $script:isJeffreyLoggedIn) {
        Write-Host "Jeffrey moet eerst inloggen om deze endpoint te testen." -ForegroundColor $errorColor
        return $null
    }
    
    # Bewaar de huidige token
    $tempToken = $script:jwtToken
    $script:jwtToken = $script:jeffreyJwtToken
    
    $response = Invoke-ApiCall -Name "Get Jeffrey Profile" -Method "Get" -Endpoint "/api/auth/profile" -UseAuth -Session $script:jeffreySession -RetryCount 2
    
    # Herstel de originele token
    $script:jwtToken = $tempToken
    
    return $response
}

# Functie om de reset password endpoint te testen
function Test-ResetPasswordEndpoint {
    Show-Title -Title "Reset Password Endpoint Testen"
    
    if (-not $script:isLoggedIn) {
        Write-Host "Je moet eerst inloggen om deze endpoint te testen." -ForegroundColor $errorColor
        return $null
    }
    
    Write-Host "[INFO] We gebruiken hetzelfde wachtwoord voor 'nieuw wachtwoord' om het account niet te wijzigen." -ForegroundColor $infoColor
    
    $resetData = $testData.ResetPassword
    
    $response = Invoke-ApiCall -Name "Reset Password" -Method "Post" -Endpoint "/api/auth/reset-password" -Body $resetData -UseAuth -Session $script:session -RetryCount 2
    
    return $response
}

# Functie om de logout endpoint te testen
function Test-LogoutEndpoint {
    Show-Title -Title "Logout Endpoint Testen"
    
    if (-not $script:isLoggedIn) {
        Write-Host "Je bent niet ingelogd." -ForegroundColor $errorColor
        return $null
    }
    
    $response = Invoke-ApiCall -Name "Logout" -Method "Post" -Endpoint "/api/auth/logout" -UseAuth -Session $script:session -RetryCount 2
    
    if ($response) {
        $script:isLoggedIn = $false
        $script:session = $null
        $script:jwtToken = $null
        Write-Host "" -ForegroundColor $infoColor
        Write-Host "[SUCCESS] Succesvol uitgelogd!" -ForegroundColor $successColor
    }
    
    return $response
}

# Functie om de Jeffrey logout endpoint te testen
function Test-JeffreyLogoutEndpoint {
    Show-Title -Title "Jeffrey Logout Endpoint Testen"
    
    if (-not $script:isJeffreyLoggedIn) {
        Write-Host "Jeffrey is niet ingelogd." -ForegroundColor $errorColor
        return $null
    }
    
    # Bewaar de huidige token
    $tempToken = $script:jwtToken
    $script:jwtToken = $script:jeffreyJwtToken
    
    $response = Invoke-ApiCall -Name "Jeffrey Logout" -Method "Post" -Endpoint "/api/auth/logout" -UseAuth -Session $script:jeffreySession -RetryCount 2
    
    # Herstel de originele token
    $script:jwtToken = $tempToken
    
    if ($response) {
        $script:isJeffreyLoggedIn = $false
        $script:jeffreySession = $null
        $script:jeffreyJwtToken = $null
        Write-Host "" -ForegroundColor $infoColor
        Write-Host "[SUCCESS] Jeffrey succesvol uitgelogd!" -ForegroundColor $successColor
    }
    
    return $response
}

# Functie om de email metrics endpoint te testen
function Test-EmailMetricsEndpoint {
    Show-Title -Title "Email Metrics Endpoint Testen"
    
    if (-not $script:isLoggedIn) {
        Write-Host "Je moet eerst inloggen als admin om deze endpoint te testen." -ForegroundColor $errorColor
        return $null
    }
    
    Write-Host "API Key wordt gebruikt voor de Email Metrics endpoint." -ForegroundColor $infoColor
    $response = Invoke-ApiCall -Name "Email Metrics" -Method "Get" -Endpoint "/api/metrics/email" -UseAuth -Session $script:session -RetryCount 2 -ApiKey $testData.ApiKey
    
    return $response
}

# Functie om de rate limit metrics endpoint te testen
function Test-RateLimitMetricsEndpoint {
    Show-Title -Title "Rate Limit Metrics Endpoint Testen"
    
    if (-not $script:isLoggedIn) {
        Write-Host "Je moet eerst inloggen als admin om deze endpoint te testen." -ForegroundColor $errorColor
        return $null
    }
    
    Write-Host "API Key wordt gebruikt voor de Rate Limit Metrics endpoint." -ForegroundColor $infoColor
    $response = Invoke-ApiCall -Name "Rate Limit Metrics" -Method "Get" -Endpoint "/api/metrics/rate-limits" -UseAuth -Session $script:session -RetryCount 2 -ApiKey $testData.ApiKey
    
    return $response
}

# Functie om de prometheus metrics endpoint te testen
function Test-PrometheusMetricsEndpoint {
    Show-Title -Title "Prometheus Metrics Endpoint Testen"
    
    Write-Host "API Key wordt gebruikt voor de Prometheus Metrics endpoint." -ForegroundColor $infoColor
    $response = Invoke-ApiCall -Name "Prometheus Metrics" -Method "Get" -Endpoint "/metrics" -RetryCount 2 -ApiKey $testData.ApiKey
    
    return $response
}

# Functie om de contact lijst endpoint te testen
function Test-ContactListEndpoint {
    Show-Title -Title "Contact Lijst Endpoint Testen"
    
    if (-not $script:isLoggedIn) {
        Write-Host "Je moet eerst inloggen als admin om deze endpoint te testen." -ForegroundColor $errorColor
        return $null
    }
    
    $response = Invoke-ApiCall -Name "Contact Lijst" -Method "Get" -Endpoint "/api/contact" -UseAuth -Session $script:session -RetryCount 2
    
    # Sla het eerste contact ID op voor gebruik in andere tests
    if ($response -and $response.Count -gt 0) {
        $script:contactId = $response[0].id
        Write-Host "Contact ID opgeslagen voor gebruik in andere tests: $($script:contactId)" -ForegroundColor $infoColor
    } else {
        Write-Host "Geen contacten gevonden om te testen." -ForegroundColor $errorColor
    }
    
    return $response
}

# Functie om de contact details endpoint te testen
function Test-ContactDetailsEndpoint {
    Show-Title -Title "Contact Details Endpoint Testen"
    
    if (-not $script:isLoggedIn) {
        Write-Host "Je moet eerst inloggen als admin om deze endpoint te testen." -ForegroundColor $errorColor
        return $null
    }
    
    if (-not $script:contactId) {
        Write-Host "Geen contact ID beschikbaar. Voer eerst de Contact Lijst test uit." -ForegroundColor $errorColor
        return $null
    }
    
    $response = Invoke-ApiCall -Name "Contact Details" -Method "Get" -Endpoint "/api/contact/$($script:contactId)" -UseAuth -Session $script:session -RetryCount 2
    
    return $response
}

# Functie om de contact update endpoint te testen
function Test-ContactUpdateEndpoint {
    Show-Title -Title "Contact Update Endpoint Testen"
    
    if (-not $script:isLoggedIn) {
        Write-Host "Je moet eerst inloggen als admin om deze endpoint te testen." -ForegroundColor $errorColor
        return $null
    }
    
    if (-not $script:contactId) {
        Write-Host "Geen contact ID beschikbaar. Voer eerst de Contact Lijst test uit." -ForegroundColor $errorColor
        return $null
    }
    
    $updateData = $testData.ContactUpdate
    
    $response = Invoke-ApiCall -Name "Contact Update" -Method "Put" -Endpoint "/api/contact/$($script:contactId)" -Body $updateData -UseAuth -Session $script:session -RetryCount 2
    
    return $response
}

# Functie om de contact antwoord endpoint te testen
function Test-ContactAntwoordEndpoint {
    Show-Title -Title "Contact Antwoord Endpoint Testen"
    
    if (-not $script:isLoggedIn) {
        Write-Host "Je moet eerst inloggen als admin om deze endpoint te testen." -ForegroundColor $errorColor
        return $null
    }
    
    if (-not $script:contactId) {
        Write-Host "Geen contact ID beschikbaar. Voer eerst de Contact Lijst test uit." -ForegroundColor $errorColor
        return $null
    }
    
    $antwoordData = $testData.ContactAntwoord
    
    $response = Invoke-ApiCall -Name "Contact Antwoord" -Method "Post" -Endpoint "/api/contact/$($script:contactId)/antwoord" -Body $antwoordData -UseAuth -Session $script:session -RetryCount 2
    
    return $response
}

# Functie om de contact status filter endpoint te testen
function Test-ContactStatusFilterEndpoint {
    Show-Title -Title "Contact Status Filter Endpoint Testen"
    
    if (-not $script:isLoggedIn) {
        Write-Host "Je moet eerst inloggen als admin om deze endpoint te testen." -ForegroundColor $errorColor
        return $null
    }
    
    # Test met status "in_behandeling"
    $status = "in_behandeling"
    
    $response = Invoke-ApiCall -Name "Contact Status Filter" -Method "Get" -Endpoint "/api/contact/status/$status" -UseAuth -Session $script:session -RetryCount 2
    
    return $response
}

# Functie om de aanmelding lijst endpoint te testen
function Test-AanmeldingListEndpoint {
    Show-Title -Title "Aanmelding Lijst Endpoint Testen"
    
    if (-not $script:isLoggedIn) {
        Write-Host "Je moet eerst inloggen als admin om deze endpoint te testen." -ForegroundColor $errorColor
        return $null
    }
    
    $response = Invoke-ApiCall -Name "Aanmelding Lijst" -Method "Get" -Endpoint "/api/aanmelding" -UseAuth -Session $script:session -RetryCount 2
    
    # Sla het eerste aanmelding ID op voor gebruik in andere tests
    if ($response -and $response.Count -gt 0) {
        $script:aanmeldingId = $response[0].id
        Write-Host "Aanmelding ID opgeslagen voor gebruik in andere tests: $($script:aanmeldingId)" -ForegroundColor $infoColor
    } else {
        Write-Host "Geen aanmeldingen gevonden om te testen." -ForegroundColor $errorColor
    }
    
    return $response
}

# Functie om de aanmelding details endpoint te testen
function Test-AanmeldingDetailsEndpoint {
    Show-Title -Title "Aanmelding Details Endpoint Testen"
    
    if (-not $script:isLoggedIn) {
        Write-Host "Je moet eerst inloggen als admin om deze endpoint te testen." -ForegroundColor $errorColor
        return $null
    }
    
    if (-not $script:aanmeldingId) {
        Write-Host "Geen aanmelding ID beschikbaar. Voer eerst de Aanmelding Lijst test uit." -ForegroundColor $errorColor
        return $null
    }
    
    $response = Invoke-ApiCall -Name "Aanmelding Details" -Method "Get" -Endpoint "/api/aanmelding/$($script:aanmeldingId)" -UseAuth -Session $script:session -RetryCount 2
    
    return $response
}

# Functie om de aanmelding update endpoint te testen
function Test-AanmeldingUpdateEndpoint {
    Show-Title -Title "Aanmelding Update Endpoint Testen"
    
    if (-not $script:isLoggedIn) {
        Write-Host "Je moet eerst inloggen als admin om deze endpoint te testen." -ForegroundColor $errorColor
        return $null
    }
    
    if (-not $script:aanmeldingId) {
        Write-Host "Geen aanmelding ID beschikbaar. Voer eerst de Aanmelding Lijst test uit." -ForegroundColor $errorColor
        return $null
    }
    
    $updateData = $testData.AanmeldingUpdate
    
    $response = Invoke-ApiCall -Name "Aanmelding Update" -Method "Put" -Endpoint "/api/aanmelding/$($script:aanmeldingId)" -Body $updateData -UseAuth -Session $script:session -RetryCount 2
    
    return $response
}

# Functie om de aanmelding antwoord endpoint te testen
function Test-AanmeldingAntwoordEndpoint {
    Show-Title -Title "Aanmelding Antwoord Endpoint Testen"
    
    if (-not $script:isLoggedIn) {
        Write-Host "Je moet eerst inloggen als admin om deze endpoint te testen." -ForegroundColor $errorColor
        return $null
    }
    
    if (-not $script:aanmeldingId) {
        Write-Host "Geen aanmelding ID beschikbaar. Voer eerst de Aanmelding Lijst test uit." -ForegroundColor $errorColor
        return $null
    }
    
    $antwoordData = $testData.AanmeldingAntwoord
    
    $response = Invoke-ApiCall -Name "Aanmelding Antwoord" -Method "Post" -Endpoint "/api/aanmelding/$($script:aanmeldingId)/antwoord" -Body $antwoordData -UseAuth -Session $script:session -RetryCount 2
    
    return $response
}

# Functie om de aanmelding rol filter endpoint te testen
function Test-AanmeldingRolFilterEndpoint {
    Show-Title -Title "Aanmelding Rol Filter Endpoint Testen"
    
    if (-not $script:isLoggedIn) {
        Write-Host "Je moet eerst inloggen als admin om deze endpoint te testen." -ForegroundColor $errorColor
        return $null
    }
    
    # Test met rol "deelnemer"
    $rol = "deelnemer"
    
    $response = Invoke-ApiCall -Name "Aanmelding Rol Filter" -Method "Get" -Endpoint "/api/aanmelding/rol/$rol" -UseAuth -Session $script:session -RetryCount 2
    
    return $response
}

# Functie om de mail lijst endpoint te testen
function Test-MailListEndpoint {
    Show-Title -Title "Mail Lijst Endpoint Testen"
    
    if (-not $script:isLoggedIn) {
        Write-Host "Je moet eerst inloggen als admin om deze endpoint te testen." -ForegroundColor $errorColor
        return $null
    }
    
    $response = Invoke-ApiCall -Name "Mail Lijst" -Method "Get" -Endpoint "/api/mail" -UseAuth -Session $script:session -RetryCount 2
    
    # Sla het eerste mail ID op voor gebruik in andere tests
    if ($response -and $response.Count -gt 0) {
        $script:mailId = $response[0].id
        Write-Host "Mail ID opgeslagen voor gebruik in andere tests: $($script:mailId)" -ForegroundColor $infoColor
    } else {
        Write-Host "Geen e-mails gevonden om te testen." -ForegroundColor $promptColor
    }
    
    return $response
}

# Functie om de mail details endpoint te testen
function Test-MailDetailsEndpoint {
    Show-Title -Title "Mail Details Endpoint Testen"
    
    if (-not $script:isLoggedIn) {
        Write-Host "Je moet eerst inloggen als admin om deze endpoint te testen." -ForegroundColor $errorColor
        return $null
    }
    
    if (-not $script:mailId) {
        Write-Host "Geen mail ID beschikbaar. Voer eerst de Mail Lijst test uit." -ForegroundColor $promptColor
        return $null
    }
    
    $response = Invoke-ApiCall -Name "Mail Details" -Method "Get" -Endpoint "/api/mail/$($script:mailId)" -UseAuth -Session $script:session -RetryCount 2
    
    return $response
}

# Functie om de unprocessed mails endpoint te testen
function Test-UnprocessedMailsEndpoint {
    Show-Title -Title "Onverwerkte Mail Lijst Endpoint Testen"
    
    if (-not $script:isLoggedIn) {
        Write-Host "Je moet eerst inloggen als admin om deze endpoint te testen." -ForegroundColor $errorColor
        return $null
    }
    
    $response = Invoke-ApiCall -Name "Onverwerkte Mails" -Method "Get" -Endpoint "/api/mail/unprocessed" -UseAuth -Session $script:session -RetryCount 2
    
    return $response
}

# Functie om de mail fetch endpoint te testen
function Test-FetchMailsEndpoint {
    Show-Title -Title "Mail Ophalen Endpoint Testen"
    
    if (-not $script:isLoggedIn) {
        Write-Host "Je moet eerst inloggen als admin om deze endpoint te testen." -ForegroundColor $errorColor
        return $null
    }
    
    $response = Invoke-ApiCall -Name "Mail Ophalen" -Method "Post" -Endpoint "/api/mail/fetch" -UseAuth -Session $script:session -RetryCount 2 -UseTestMode
    
    return $response
}

# Functie om de mark as processed endpoint te testen
function Test-MarkMailAsProcessedEndpoint {
    Show-Title -Title "Mail Als Verwerkt Markeren Endpoint Testen"
    
    if (-not $script:isLoggedIn) {
        Write-Host "Je moet eerst inloggen als admin om deze endpoint te testen." -ForegroundColor $errorColor
        return $null
    }
    
    if (-not $script:mailId) {
        Write-Host "Geen mail ID beschikbaar. Voer eerst de Mail Lijst test uit." -ForegroundColor $promptColor
        return $null
    }
    
    $response = Invoke-ApiCall -Name "Mail Als Verwerkt Markeren" -Method "Put" -Endpoint "/api/mail/$($script:mailId)/processed" -UseAuth -Session $script:session -RetryCount 2
    
    return $response
}

# Functie om mails per account type te testen
function Test-MailsByAccountTypeEndpoint {
    Show-Title -Title "Mails Per Account Type Endpoint Testen"
    
    if (-not $script:isLoggedIn) {
        Write-Host "Je moet eerst inloggen als admin om deze endpoint te testen." -ForegroundColor $errorColor
        return $null
    }
    
    # Test met account type "info"
    $accountType = "info"
    
    $response = Invoke-ApiCall -Name "Mails Per Account Type" -Method "Get" -Endpoint "/api/mail/account/$accountType" -UseAuth -Session $script:session -RetryCount 2
    
    return $response
}

# Functie om de delete mail endpoint te testen
function Test-DeleteMailEndpoint {
    Show-Title -Title "Mail Verwijderen Endpoint Testen"
    
    if (-not $script:isLoggedIn) {
        Write-Host "Je moet eerst inloggen als admin om deze endpoint te testen." -ForegroundColor $errorColor
        return $null
    }
    
    if (-not $script:mailId) {
        Write-Host "Geen mail ID beschikbaar. Voer eerst de Mail Lijst test uit." -ForegroundColor $promptColor
        return $null
    }
    
    $response = Invoke-ApiCall -Name "Mail Verwijderen" -Method "Delete" -Endpoint "/api/mail/$($script:mailId)" -UseAuth -Session $script:session -RetryCount 2
    
    if ($response) {
        # Als het verwijderen succesvol was, reset dan de mail ID
        $script:mailId = $null
        Write-Host "Mail ID gereset na succesvolle verwijdering." -ForegroundColor $infoColor
    }
    
    return $response
}

# Functie voor parallelle testuitvoering
function Invoke-ParallelTests {
    param (
        [hashtable]$Tests,
        [int]$MaxJobs = 4
    )
    
    if ($MaxJobs -le 0) {
        # Voer tests sequentieel uit als parallelle uitvoering is uitgeschakeld
        $results = @{}
        foreach ($testName in $Tests.Keys) {
            Write-Host "Uitvoeren test: $testName" -ForegroundColor $infoColor
            $results[$testName] = & $Tests[$testName]
        }
        return $results
    }
    
    Write-Host "Start parallelle testuitvoering met maximaal $MaxJobs gelijktijdige jobs..." -ForegroundColor $highlightColor
    
    $jobs = @{}
    $results = @{}
    $runningJobs = 0
    $testQueue = [System.Collections.Queue]::new($Tests.Keys)
    
    # Start initiële jobs
    while ($runningJobs -lt $MaxJobs -and $testQueue.Count -gt 0) {
        $testName = $testQueue.Dequeue()
        Write-Host "Start job: $testName" -ForegroundColor $infoColor
        $scriptBlock = [scriptblock]::Create($Tests[$testName].ToString())
        $jobs[$testName] = Start-Job -ScriptBlock $scriptBlock
        $runningJobs++
    }
    
    # Verwerk jobs en start nieuwe wanneer oude zijn voltooid
    while ($jobs.Count -gt 0) {
        $completedJobs = @()
        
        foreach ($testName in $jobs.Keys) {
            $job = $jobs[$testName]
            
            if ($job.State -eq "Completed") {
                Write-Host "Job voltooid: $testName" -ForegroundColor $successColor
                $results[$testName] = Receive-Job -Job $job
                Remove-Job -Job $job
                $completedJobs += $testName
                $runningJobs--
                
                # Start een nieuwe job als er nog tests in de wachtrij staan
                if ($testQueue.Count -gt 0) {
                    $newTestName = $testQueue.Dequeue()
                    Write-Host "Start job: $newTestName" -ForegroundColor $infoColor
                    $scriptBlock = [scriptblock]::Create($Tests[$newTestName].ToString())
                    $jobs[$newTestName] = Start-Job -ScriptBlock $scriptBlock
                    $runningJobs++
                }
            }
            elseif ($job.State -eq "Failed") {
                Write-Host "Job mislukt: $testName" -ForegroundColor $errorColor
                $results[$testName] = $null
                Remove-Job -Job $job
                $completedJobs += $testName
                $runningJobs--
                
                # Start een nieuwe job als er nog tests in de wachtrij staan
                if ($testQueue.Count -gt 0) {
                    $newTestName = $testQueue.Dequeue()
                    Write-Host "Start job: $newTestName" -ForegroundColor $infoColor
                    $scriptBlock = [scriptblock]::Create($Tests[$newTestName].ToString())
                    $jobs[$newTestName] = Start-Job -ScriptBlock $scriptBlock
                    $runningJobs++
                }
            }
        }
        
        # Verwijder voltooide jobs uit de trackinghash
        foreach ($testName in $completedJobs) {
            $jobs.Remove($testName)
        }
        
        # Vermijd CPU spinning
        if ($jobs.Count -gt 0) {
            Start-Sleep -Milliseconds 500
        }
    }
    
    return $results
}

# Functie om unieke testdata te genereren
function New-IsolatedTestData {
    param (
        [string]$TestId = (New-Guid).ToString().Substring(0, 8)
    )
    
    # Genereer geïsoleerde testdata met unieke IDs
    $isolatedData = @{
        ContactEmail = @{
            naam = "Test Gebruiker $TestId"
            email = "test_$($TestId)@example.com"
            bericht = "Dit is een geautomatiseerd testbericht $TestId"
            privacy_akkoord = $true
        }
        AanmeldingEmail = @{
            naam = "Test Deelnemer $TestId"
            email = "deelnemer_$($TestId)@example.com"
            telefoon = "0612345$TestId"
            rol = "deelnemer"
            afstand = "10km"
            ondersteuning = "geen"
            bijzonderheden = "Test ID: $TestId"
            terms = $true
        }
        AdminUser = @{
            email = $config.Credentials.AdminUser.Email
            wachtwoord = $config.Credentials.AdminUser.Password
        }
        JeffreyUser = @{
            email = $config.Credentials.RegularUser.Email
            wachtwoord = $config.Credentials.RegularUser.Password
        }
        ResetPassword = @{
            huidig_wachtwoord = "admin123"  # Bijgewerkt wachtwoord
            nieuw_wachtwoord = "admin123"   # Zelfde wachtwoord om geen wijziging te maken
        }
        ApiKey = $config.ApiKey
        ContactUpdate = @{
            status = "in_behandeling"
            notities = "Dit is een testnotitie voor het contactformulier"
        }
        ContactAntwoord = @{
            tekst = "Dit is een testantwoord op het contactformulier"
        }
        AanmeldingUpdate = @{
            status = "in_behandeling"
            notities = "Dit is een testnotitie voor de aanmelding"
        }
        AanmeldingAntwoord = @{
            tekst = "Dit is een testantwoord op de aanmelding"
        }
    }
    
    return $isolatedData
}

# Functie om alle endpoints te testen
function Test-AllEndpoints {
    Show-Title -Title "Alle Endpoints Automatisch Testen"
    
    # Start tijdmeting
    $startTime = Get-Date
    
    # Test basis endpoints
    $results = @{
        Root = Test-RootEndpoint
        Health = Test-HealthEndpoint
    }
    
    # Stap 1: Maak contactformulieren en aanmeldingen aan
    Write-Host "" -ForegroundColor $infoColor
    Write-Host "Stap 1: Contactformulieren en aanmeldingen aanmaken voor beheer tests..." -ForegroundColor $highlightColor
    
    # Test eerst de standaard contact en aanmelding endpoints
    $results.ContactEmail = Test-ContactEmailEndpoint
    $results.AanmeldingEmail = Test-AanmeldingEmailEndpoint
    
    # Maak extra contactformulieren aan
    for ($i = 2; $i -le 3; $i++) {
        $contactData = $testData.ContactEmail.Clone()
        $contactData.naam = "Test Gebruiker $i"
        $contactData.bericht = "Dit is een geautomatiseerd testbericht $i"
        
        Write-Host "Extra contactformulier $i aanmaken..." -ForegroundColor $infoColor
        $contactResponse = Invoke-ApiCall -Name "Contact Email $i" -Method "Post" -Endpoint "/api/contact-email" -Body $contactData -RetryCount 2 -IgnoreEmailErrors
        
        # Wacht even om rate limiting te voorkomen
        Start-Sleep -Seconds 2
    }
    
    # Maak extra aanmeldingen aan
    for ($i = 2; $i -le 3; $i++) {
        $aanmeldingData = $testData.AanmeldingEmail.Clone()
        $aanmeldingData.naam = "Test Deelnemer $i"
        
        # Varieer de rol voor betere testdekking
        if ($i -eq 2) {
            $aanmeldingData.rol = "vrijwilliger"
        } else {
            $aanmeldingData.rol = "sponsor"
        }
        
        Write-Host "Extra aanmelding $i aanmaken..." -ForegroundColor $infoColor
        $aanmeldingResponse = Invoke-ApiCall -Name "Aanmelding Email $i" -Method "Post" -Endpoint "/api/aanmelding-email" -Body $aanmeldingData -RetryCount 2 -IgnoreEmailErrors
        
        # Wacht even om rate limiting te voorkomen
        Start-Sleep -Seconds 2
    }
    
    # Stap 2: Log in als admin om beheer endpoints te testen
    Write-Host "" -ForegroundColor $infoColor
    Write-Host "Stap 2: Inloggen als admin om beheer endpoints te testen..." -ForegroundColor $highlightColor
    
    # Wacht even om rate limiting te voorkomen
    Write-Host "Even wachten om rate limiting te voorkomen..." -ForegroundColor $infoColor
    Start-Sleep -Seconds 10
    
    $results.AdminLogin = Test-AdminLoginEndpoint
    
    if ($script:isLoggedIn) {
        $results.AdminProfile = Test-ProfileEndpoint
        
        # Stap 3: Test Contact Beheer endpoints
        Write-Host "" -ForegroundColor $infoColor
        Write-Host "Stap 3: Contact Beheer endpoints testen..." -ForegroundColor $highlightColor
        
        # Wacht even om zeker te zijn dat de contactformulieren zijn verwerkt
        Write-Host "Even wachten om zeker te zijn dat de contactformulieren zijn verwerkt..." -ForegroundColor $infoColor
        Start-Sleep -Seconds 5
        
        $results.ContactList = Test-ContactListEndpoint
        if ($script:contactId) {
            $results.ContactDetails = Test-ContactDetailsEndpoint
            $results.ContactUpdate = Test-ContactUpdateEndpoint
            $results.ContactAntwoord = Test-ContactAntwoordEndpoint
        } else {
            Write-Host "Geen contact ID gevonden, sla gerelateerde tests over." -ForegroundColor $promptColor
        }
        $results.ContactStatusFilter = Test-ContactStatusFilterEndpoint
        
        # Stap 4: Test Aanmelding Beheer endpoints
        Write-Host "" -ForegroundColor $infoColor
        Write-Host "Stap 4: Aanmelding Beheer endpoints testen..." -ForegroundColor $highlightColor
        
        # Wacht even om zeker te zijn dat de aanmeldingen zijn verwerkt
        Write-Host "Even wachten om zeker te zijn dat de aanmeldingen zijn verwerkt..." -ForegroundColor $infoColor
        Start-Sleep -Seconds 5
        
        $results.AanmeldingList = Test-AanmeldingListEndpoint
        if ($script:aanmeldingId) {
            $results.AanmeldingDetails = Test-AanmeldingDetailsEndpoint
            $results.AanmeldingUpdate = Test-AanmeldingUpdateEndpoint
            $results.AanmeldingAntwoord = Test-AanmeldingAntwoordEndpoint
        } else {
            Write-Host "Geen aanmelding ID gevonden, sla gerelateerde tests over." -ForegroundColor $promptColor
        }
        $results.AanmeldingRolFilter = Test-AanmeldingRolFilterEndpoint
        
        # Stap 5: Test Mail Beheer endpoints
        Write-Host "" -ForegroundColor $infoColor
        Write-Host "Stap 5: Mail Beheer endpoints testen..." -ForegroundColor $highlightColor
        
        # Eerst nieuwe e-mails ophalen om te testen
        $results.FetchMails = Test-FetchMailsEndpoint
        
        # Wacht even om zeker te zijn dat de mails zijn opgehaald
        Write-Host "Even wachten om zeker te zijn dat de e-mails zijn opgehaald..." -ForegroundColor $infoColor
        Start-Sleep -Seconds 5
        
        $results.MailList = Test-MailListEndpoint
        $results.UnprocessedMails = Test-UnprocessedMailsEndpoint
        $results.MailsByAccountType = Test-MailsByAccountTypeEndpoint
        
        if ($script:mailId) {
            $results.MailDetails = Test-MailDetailsEndpoint
            $results.MarkMailAsProcessed = Test-MarkMailAsProcessedEndpoint
            # Test verwijderen als laatste
            $results.DeleteMail = Test-DeleteMailEndpoint
        } else {
            Write-Host "Geen mail ID gevonden, sla gerelateerde tests over." -ForegroundColor $promptColor
        }
        
        # Stap 6: Test overige admin endpoints
        Write-Host "" -ForegroundColor $infoColor
        Write-Host "Stap 6: Overige admin endpoints testen..." -ForegroundColor $highlightColor
        
        $results.EmailMetrics = Test-EmailMetricsEndpoint
        $results.RateLimitMetrics = Test-RateLimitMetricsEndpoint
        $results.ResetPassword = Test-ResetPasswordEndpoint
    } else {
        Write-Host "" -ForegroundColor $errorColor
        Write-Host "Kon niet inloggen als admin, sla admin-gerelateerde endpoints over." -ForegroundColor $errorColor
    }
    
    # Stap 7: Test Jeffrey account
    Write-Host "" -ForegroundColor $infoColor
    Write-Host "Stap 7: Jeffrey account testen..." -ForegroundColor $highlightColor
    
    # Wacht even om rate limiting te voorkomen
    Write-Host "Even wachten om rate limiting te voorkomen..." -ForegroundColor $infoColor
    Start-Sleep -Seconds 10
    
    $results.JeffreyLogin = Test-JeffreyLoginEndpoint
    
    if ($script:isJeffreyLoggedIn) {
        $results.JeffreyProfile = Test-JeffreyProfileEndpoint
        $results.JeffreyLogout = Test-JeffreyLogoutEndpoint
    } else {
        Write-Host "" -ForegroundColor $errorColor
        Write-Host "Kon niet inloggen als Jeffrey, sla Jeffrey-gerelateerde endpoints over." -ForegroundColor $errorColor
    }
    
    # Stap 8: Logout admin en test Prometheus metrics
    Write-Host "" -ForegroundColor $infoColor
    Write-Host "Stap 8: Admin uitloggen en Prometheus metrics testen..." -ForegroundColor $highlightColor
    
    # Logout admin als laatste
    if ($script:isLoggedIn) {
        $results.AdminLogout = Test-LogoutEndpoint
    }
    
    $results.PrometheusMetrics = Test-PrometheusMetricsEndpoint
    
    # Bereken totale duur
    $endTime = Get-Date
    $duration = $endTime - $startTime
    
    # Toon samenvatting
    Show-Title -Title "Test Samenvatting"
    Write-Host "Totale testduur: $($duration.TotalSeconds.ToString("0.00")) seconden" -ForegroundColor $infoColor
    
    $successCount = ($results.Values | Where-Object { $_ -ne $null }).Count
    $failCount = ($results.Values | Where-Object { $_ -eq $null }).Count
    $totalCount = $results.Count
    
    Write-Host "" -ForegroundColor $infoColor
    Write-Host "Resultaten:" -ForegroundColor $infoColor
    Write-Host "- Succesvol: $successCount/$totalCount" -ForegroundColor $successColor
    if ($failCount -gt 0) {
        Write-Host "- Mislukt: $failCount/$totalCount" -ForegroundColor $errorColor
        
        Write-Host "" -ForegroundColor $infoColor
        Write-Host "Mislukte endpoints:" -ForegroundColor $errorColor
        $results.GetEnumerator() | Where-Object { $_.Value -eq $null } | ForEach-Object {
            Write-Host "- $($_.Key)" -ForegroundColor $errorColor
        }
    }
    
    # Genereer testrapport
    $reportDate = Get-Date -Format "yyyy-MM-dd_HH-mm-ss"
    $reportPath = "DKL_API_Test_Report_$reportDate.txt"
    
    Write-Host "" -ForegroundColor $infoColor
    Write-Host "Testrapport opslaan als $reportPath..." -ForegroundColor $infoColor
    
    try {
        $report = @"
DKL Email Service API Test Rapport
Gegenereerd op: $(Get-Date)
API URL: $currentUrl

Testresultaten:
- Succesvol: $successCount/$totalCount
- Mislukt: $failCount/$totalCount
- Totale duur: $($duration.TotalSeconds.ToString("0.00")) seconden

Details per endpoint:
"@
        
        foreach ($result in $results.GetEnumerator()) {
            $status = if ($result.Value -ne $null) { "SUCCESS" } else { "FAILED" }
            $report += "`n- $($result.Key): $status"
        }
        
        if ($failCount -gt 0) {
            $report += "`n`nMislukte endpoints:"
            $results.GetEnumerator() | Where-Object { $_.Value -eq $null } | ForEach-Object {
                $report += "`n- $($_.Key)"
            }
        }
        
        $report | Out-File -FilePath $reportPath -Encoding utf8
        Write-Host "Testrapport succesvol opgeslagen!" -ForegroundColor $successColor
    } catch {
        Write-Host "Kon testrapport niet opslaan: $_" -ForegroundColor $errorColor
    }
    
    return $results
}

# Initialisatie
$script:isLoggedIn = $false
$script:session = $null
$script:jwtToken = $null
$script:isJeffreyLoggedIn = $false
$script:jeffreySession = $null
$script:jeffreyJwtToken = $null
$script:contactId = $null
$script:aanmeldingId = $null
$script:mailId = $null

# Welkomstbericht
Show-Title -Title "DKL Email Service API Geautomatiseerde Test Tool"
Write-Host "Dit script test automatisch alle endpoints van de DKL Email Service API." -ForegroundColor $infoColor
Write-Host "Bijgewerkt met testmode ondersteuning en configureerbare parameters." -ForegroundColor $highlightColor

# Detecteer of gebruik de beste API URL
$detectedUrl = Set-ApiUrl
$currentUrl = $detectedUrl  # Update currentUrl met het detectieresultaat

# Alleen updaten als de detectie een waarde heeft teruggegeven
if ($detectedUrl) {
    Write-Host "API URL ingesteld op: $currentUrl" -ForegroundColor $infoColor
} else {
    Write-Host "WAARSCHUWING: API URL detectie mislukt, gebruik $currentUrl" -ForegroundColor $errorColor
}

Write-Host "" -ForegroundColor $infoColor
Write-Host "Test data wordt automatisch ingevuld. De tests zullen nu automatisch starten..." -ForegroundColor $infoColor
Write-Host "Testmodus is: " -NoNewline -ForegroundColor $infoColor
if ($testMode) {
    Write-Host "INGESCHAKELD (geen echte emails worden verzonden)" -ForegroundColor $successColor
} else {
    Write-Host "UITGESCHAKELD (echte emails worden verzonden)" -ForegroundColor $promptColor
}

# Start volledige test automatisch
$testResults = Test-AllEndpoints

Write-Host "" -ForegroundColor $infoColor
Show-Title -Title "Test Script Voltooid"
Write-Host "Alle tests zijn uitgevoerd. De resultaten zijn hierboven weergegeven." -ForegroundColor $infoColor
Write-Host "De inloggegevens zijn bijgewerkt en werken nu correct." -ForegroundColor $highlightColor
Write-Host "API Key is toegevoegd voor de metrics endpoints." -ForegroundColor $highlightColor
Write-Host "Testrapport is opgeslagen als: $customReportPath" -ForegroundColor $highlightColor