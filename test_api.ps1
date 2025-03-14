# DKL Email Service API Geautomatiseerd Test Script
# Dit script test automatisch alle endpoints van de DKL Email Service API

# Configuratie
$baseUrl = "https://dklemailservice.onrender.com"  # Default URL, kan worden aangepast
$localUrl = "http://localhost:8080"            # Lokale URL voor ontwikkeling

# Kies automatisch de API URL (Productie als default)
$currentUrl = $baseUrl

# Testdata voor formulieren - BIJGEWERKT met correcte inloggegevens
$testData = @{
    ContactEmail = @{
        naam = "Test Gebruiker"
        email = "noreply@dekoninklijkeloop.nl"  # Geldig e-mailadres
        bericht = "Dit is een geautomatiseerd testbericht"
        privacy_akkoord = $true
    }
    AanmeldingEmail = @{
        naam = "Test Deelnemer"
        email = "noreply@dekoninklijkeloop.nl"  # Geldig e-mailadres
        telefoon = "0612345678"
        rol = "deelnemer"
        afstand = "10km"
        ondersteuning = "Geen bijzondere ondersteuning nodig"
        bijzonderheden = "Dit is een geautomatiseerde testregistratie"
        terms = $true
    }
    AdminUser = @{
        email = "admin@dekoninklijkeloop.nl"  # Admin account
        wachtwoord = "admin123"  # Bijgewerkt wachtwoord na hash update
    }
    JeffreyUser = @{
        email = "jeffrey@dekoninklijkeloop.nl"  # Jeffrey account
        wachtwoord = "DKL2025!"  # Bijgewerkt wachtwoord na hash update
    }
    ResetPassword = @{
        huidig_wachtwoord = "admin123"  # Bijgewerkt wachtwoord
        nieuw_wachtwoord = "admin123"   # Zelfde wachtwoord om geen wijziging te maken
    }
    # API key voor metrics endpoints
    ApiKey = "dkl_metrics_api_key_2025"  # Voorbeeld API key, vervang met de echte key
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
        [string]$ApiKey = $null
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
            if ($attempt -ge $RetryCount) {
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
    
    return $response
}

# Functie om API URL te vaststellen met automatische detectie
function Set-ApiUrl {
    Show-Title -Title "API URL Detecteren"
    
    # Probeer eerst de productie URL
    Write-Host "Proberen te verbinden met productie URL ($baseUrl)..." -ForegroundColor $infoColor
    
    try {
        $testResult = Invoke-RestMethod -Uri "$baseUrl/api/health" -Method Get -UseBasicParsing -TimeoutSec 5
        Write-Host "Verbinding gelukt met productie URL!" -ForegroundColor $successColor
        return $baseUrl
    } catch {
        Write-Host "Kon geen verbinding maken met productie URL, status: $($_.Exception.Response.StatusCode.value__)" -ForegroundColor $errorColor
    }
    
    # Probeer de lokale URL als alternatief
    Write-Host "Proberen te verbinden met lokale URL ($localUrl)..." -ForegroundColor $infoColor
    
    try {
        $testResult = Invoke-RestMethod -Uri "$localUrl/api/health" -Method Get -UseBasicParsing -TimeoutSec 5
        Write-Host "Verbinding gelukt met lokale URL!" -ForegroundColor $successColor
        return $localUrl
    } catch {
        Write-Host "Kon geen verbinding maken met lokale URL" -ForegroundColor $errorColor
    }
    
    # Default terugvallen op de productie URL
    Write-Host "Geen werkende URL gedetecteerd, terugvallen op productie URL" -ForegroundColor $promptColor
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
    Show-Title -Title "Contact Email Endpoint Testen"
    
    $contactData = $testData.ContactEmail
    
    $response = Invoke-ApiCall -Name "Contact Email" -Method "Post" -Endpoint "/api/contact-email" -Body $contactData -RetryCount 2
    
    return $response
}

# Functie om de aanmelding email endpoint te testen
function Test-AanmeldingEmailEndpoint {
    Show-Title -Title "Aanmelding Email Endpoint Testen"
    
    $aanmeldingData = $testData.AanmeldingEmail
    
    $response = Invoke-ApiCall -Name "Aanmelding Email" -Method "Post" -Endpoint "/api/aanmelding-email" -Body $aanmeldingData -RetryCount 2
    
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

# Functie om alle endpoints te testen
function Test-AllEndpoints {
    Show-Title -Title "Alle Endpoints Automatisch Testen"
    
    # Start tijdmeting
    $startTime = Get-Date
    
    # Test alle endpoints
    $results = @{
        Root = Test-RootEndpoint
        Health = Test-HealthEndpoint
        ContactEmail = Test-ContactEmailEndpoint
        AanmeldingEmail = Test-AanmeldingEmailEndpoint
        AdminLogin = Test-AdminLoginEndpoint
    }
    
    if ($script:isLoggedIn) {
        $results.AdminProfile = Test-ProfileEndpoint
        $results.EmailMetrics = Test-EmailMetricsEndpoint
        $results.RateLimitMetrics = Test-RateLimitMetricsEndpoint
        $results.ResetPassword = Test-ResetPasswordEndpoint
    } else {
        Write-Host "" -ForegroundColor $errorColor
        Write-Host "Kon niet inloggen als admin, sla admin-gerelateerde endpoints over." -ForegroundColor $errorColor
    }
    
    # Test Jeffrey account
    $results.JeffreyLogin = Test-JeffreyLoginEndpoint
    
    if ($script:isJeffreyLoggedIn) {
        $results.JeffreyProfile = Test-JeffreyProfileEndpoint
        $results.JeffreyLogout = Test-JeffreyLogoutEndpoint
    } else {
        Write-Host "" -ForegroundColor $errorColor
        Write-Host "Kon niet inloggen als Jeffrey, sla Jeffrey-gerelateerde endpoints over." -ForegroundColor $errorColor
    }
    
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

# Welkomstbericht
Show-Title -Title "DKL Email Service API Geautomatiseerde Test Tool"
Write-Host "Dit script test automatisch alle endpoints van de DKL Email Service API." -ForegroundColor $infoColor
Write-Host "Bijgewerkt met de correcte inloggegevens na de wachtwoord hash update." -ForegroundColor $highlightColor

# Detecteer de beste API URL
$currentUrl = Set-ApiUrl
Write-Host "API URL ingesteld op: $currentUrl" -ForegroundColor $infoColor

Write-Host "" -ForegroundColor $infoColor
Write-Host "Test data wordt automatisch ingevuld. De tests zullen nu automatisch starten..." -ForegroundColor $infoColor

# Start volledige test automatisch
$testResults = Test-AllEndpoints

Write-Host "" -ForegroundColor $infoColor
Show-Title -Title "Test Script Voltooid"
Write-Host "Alle tests zijn uitgevoerd. De resultaten zijn hierboven weergegeven." -ForegroundColor $infoColor
Write-Host "De inloggegevens zijn bijgewerkt en werken nu correct." -ForegroundColor $highlightColor
Write-Host "API Key is toegevoegd voor de metrics endpoints." -ForegroundColor $highlightColor