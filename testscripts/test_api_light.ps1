# DKL Email Service API Test Script
# Efficiënte tester voor alle endpoints met testmode ondersteuning

# Parameters moeten aan het begin van het script staan voor correcte PowerShell functionaliteit
param(
    [switch]$UseLocalUrl,
    [switch]$DisableTestMode,
    [string]$ApiUrl,
    [switch]$IncludeSecuredEndpoints,
    [string]$Username,
    [string]$Password,
    [string]$ApiKey,
    [switch]$SaveResults,
    [switch]$DetailedHealth,
    [switch]$SkipEmailTests,
    [switch]$TestMetrics,
    [switch]$TestMailEndpoints,
    [string]$OutputFile = "DKL_API_Test_Report_$(Get-Date -Format 'yyyy-MM-dd_HH-mm-ss').txt"
)

# Kleuren voor output
$successColor = "Green"
$errorColor = "Red"
$infoColor = "Cyan"
$promptColor = "Yellow"
$highlightColor = "Magenta"
$warningColor = "DarkYellow"

# Configuratie
$baseUrl = "https://dklemailservice.onrender.com"
$localUrl = "http://localhost:8080"
$testMode = -not $DisableTestMode
$authToken = $null
$totalTests = 0
$successfulTests = 0
$apiVersion = "Unknown"
$apiStatus = "Unknown"
$healthDetails = @{}
$defaultApiKey = "dkl_metrics_api_key_2025"

# API URL instellen
if ($UseLocalUrl) {
    $currentUrl = $localUrl
} elseif ($ApiUrl) {
    $currentUrl = $ApiUrl
} else {
    $currentUrl = $baseUrl
}

# Testdata voor formulieren
$testData = @{
    ContactEmail = @{
        naam = "Test Gebruiker"
        email = "test@example.com"
        bericht = "Dit is een geautomatiseerd testbericht"
        privacy_akkoord = $true
    }
    AanmeldingEmail = @{
        naam = "Test Deelnemer"
        email = "deelnemer@example.com"
        telefoon = "0612345678"
        rol = "deelnemer"
        afstand = "10km"
        ondersteuning = "geen"
        bijzonderheden = "Dit is een testbericht vanuit het API test script."
        terms = $true
    }
    Auth = @{
        email = if ($Username) { $Username } else { "admin@dekoninklijkeloop.nl" }
        wachtwoord = if ($Password) { $Password } else { "admin123" }
    }
}

# Start transcript als resultaten moeten worden opgeslagen
if ($SaveResults) {
    Start-Transcript -Path $OutputFile -Force
}

# Functie om een titel te tonen
function Show-Title {
    param ([string]$Title)
    
    Write-Host "" -ForegroundColor $infoColor
    Write-Host "=============================================" -ForegroundColor $infoColor
    Write-Host " $Title" -ForegroundColor $infoColor
    Write-Host "=============================================" -ForegroundColor $infoColor
}

# Functie om een API call te maken
function Invoke-ApiCall {
    param (
        [string]$Name,
        [string]$Method,
        [string]$Endpoint,
        [object]$Body = $null,
        [int]$RetryCount = 1,
        [int]$RetryDelay = 2,
        [switch]$UseTestMode = $script:testMode,
        [switch]$RequiresAuth = $false,
        [switch]$RequiresApiKey = $false,
        [string]$ApiKey = "",
        [switch]$SuppressOutput = $false
    )
    
    $script:totalTests++
    
    if (-not $SuppressOutput) {
        Write-Host ""
        Write-Host "Test: ${Name}" -ForegroundColor $infoColor
    }
    
    # Voeg test_mode toe aan body als nodig
    if ($UseTestMode -and $Body -and $Method -ne "GET") {
        if (-not $SuppressOutput) {
            Write-Host "[TEST MODE] Simuleren van $Name (geen echte email)" -ForegroundColor $promptColor
        }
        $Body["test_mode"] = $true
    }
    
    if (-not $SuppressOutput) {
        Write-Host "API Request: $Method $currentUrl$Endpoint" -ForegroundColor $infoColor
        if ($Body) {
            $jsonBody = $Body | ConvertTo-Json -Depth 5
            Write-Host "Request Body: $jsonBody" -ForegroundColor $infoColor
        }
    }
    
    $attempt = 0
    $success = $false
    $response = $null
    
    # Als authenticatie nodig is maar we hebben geen token, probeer in te loggen
    if ($RequiresAuth -and -not $script:authToken -and $IncludeSecuredEndpoints) {
        if (-not $SuppressOutput) {
            Write-Host "Authenticatie vereist. Probeert in te loggen..." -ForegroundColor $promptColor
        }
        $authResponse = Invoke-ApiCall -Name "Login" -Method "Post" -Endpoint "/api/auth/login" -Body $testData.Auth -SuppressOutput:$true
        if ($authResponse -and $authResponse.token) {
            $script:authToken = $authResponse.token
            if (-not $SuppressOutput) {
                Write-Host "Succesvol ingelogd" -ForegroundColor $successColor
            }
        } else {
            if (-not $SuppressOutput) {
                Write-Host "Inloggen mislukt. Kan beveiligde endpoints niet testen." -ForegroundColor $errorColor
            }
            return $null
        }
    }
    
    # Overslaan van beveiligde endpoints als ze niet zijn inbegrepen
    if ($RequiresAuth -and -not $IncludeSecuredEndpoints) {
        if (-not $SuppressOutput) {
            Write-Host "Beveiligde endpoint test overgeslagen (gebruik -IncludeSecuredEndpoints om te testen)" -ForegroundColor $promptColor
        }
        return $null
    }
    
    # API key bepalen
    $effectiveApiKey = ""
    if ($RequiresApiKey) {
        if ($ApiKey) {
            $effectiveApiKey = $ApiKey
        } elseif ($script:ApiKey) {
            $effectiveApiKey = $script:ApiKey
        } else {
            $effectiveApiKey = $script:defaultApiKey
        }
        
        if (-not $SuppressOutput) {
            Write-Host "API Key vereist voor deze endpoint" -ForegroundColor $promptColor
        }
    }
    
    while (-not $success -and $attempt -lt $RetryCount) {
        $attempt++
        if ($attempt -gt 1 -and -not $SuppressOutput) {
            Write-Host "Poging $attempt van $RetryCount..." -ForegroundColor $promptColor
            Start-Sleep -Seconds $RetryDelay
        }
        
        try {
            $params = @{
                Uri = "$currentUrl$Endpoint"
                Method = $Method
                UseBasicParsing = $true
            }
            
            # Headers voorbereiden
            $headers = @{}
            
            # Voeg testmodus header toe indien nodig
            if ($UseTestMode) {
                $headers["X-Test-Mode"] = "true"
            }
            
            # Voeg auth token toe indien nodig
            if ($RequiresAuth -and $script:authToken) {
                $headers["Authorization"] = "Bearer $script:authToken"
            }
            
            # Voeg API key toe indien nodig
            if ($RequiresApiKey -and $effectiveApiKey) {
                $headers["X-API-Key"] = $effectiveApiKey
            }
            
            # Headers toevoegen als er zijn
            if ($headers.Count -gt 0) {
                $params.Headers = $headers
            }
            
            # Body toevoegen indien aanwezig
            if ($Body) {
                $params.Body = $Body | ConvertTo-Json -Depth 5
                $params.ContentType = "application/json"
            }
            
            $response = Invoke-RestMethod @params
            $success = $true
            
            if (-not $SuppressOutput) {
                Write-Host "Response Status: 200 (Success)" -ForegroundColor $successColor
                if ($response) {
                    if ($response.version) {
                        Write-Host "API versie: $($response.version)" -ForegroundColor $infoColor
                        $script:apiVersion = $response.version
                    }
                    if ($response.status) {
                        Write-Host "API status: $($response.status)" -ForegroundColor $infoColor
                        $script:apiStatus = $response.status
                    }
                    if ($response.service) {
                        Write-Host "API service: $($response.service)" -ForegroundColor $infoColor
                    }
                    
                    # Sla health details op voor later gebruik
                    if ($Endpoint -eq "/api/health" -and $DetailedHealth) {
                        $script:healthDetails = $response
                    }
                    
                    if ($UseTestMode -and ($response.test_mode -or $response.testMode)) {
                        Write-Host "Server ondersteunt testmodus" -ForegroundColor $successColor
                        if ($Name -match "Contact Email") {
                            Write-Host "Contactformulier verwerkt (TEST MODE)" -ForegroundColor $successColor
                        } elseif ($Name -match "Aanmelding") {
                            Write-Host "Aanmelding verwerkt (TEST MODE)" -ForegroundColor $successColor
                        }
                    }
                    
                    if ($Endpoint -eq "/" -and $response.endpoints) {
                        Write-Host "Beschikbare endpoints:" -ForegroundColor $infoColor
                        foreach ($endpoint in $response.endpoints) {
                            $authText = if ($endpoint.auth_required) { " (Auth vereist)" } else { "" }
                            Write-Host "  - $($endpoint.method) $($endpoint.path): $($endpoint.description)$authText" -ForegroundColor $infoColor
                        }
                    } elseif ($Method -ne "GET" -or $Endpoint -eq "/api/auth/profile" -or $Endpoint -eq "/api/metrics/email") {
                        $responseJson = $response | ConvertTo-Json -Depth 5
                        if ($responseJson.Length -gt 1000) {
                            Write-Host "Response: $(($responseJson).Substring(0, 997))..." -ForegroundColor $infoColor
                        } else {
                            Write-Host "Response: $responseJson" -ForegroundColor $infoColor
                        }
                    }
                }
                
                Write-Host "Test ${Name}: Geslaagd" -ForegroundColor $successColor
                if ($UseTestMode -and ($response.test_mode -or $response.testMode)) {
                    Write-Host "Test ${Name}: Geslaagd (gesimuleerd)" -ForegroundColor $successColor
                }
            }
            
            $script:successfulTests++
        }
        catch {
            if ($attempt -ge $RetryCount -and -not $SuppressOutput) {
                $statusCode = if ($_.Exception.Response) { $_.Exception.Response.StatusCode.value__ } else { "Unknown" }
                Write-Host "Response Status: $statusCode (Error)" -ForegroundColor $errorColor
                Write-Host "Error: $($_.Exception.Message)" -ForegroundColor $errorColor
                
                try {
                    if ($_.Exception.Response -and $_.Exception.Response.GetResponseStream()) {
                        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
                        $reader.BaseStream.Position = 0
                        $reader.DiscardBufferedData()
                        $responseBody = $reader.ReadToEnd()
                        Write-Host "Response Body: $responseBody" -ForegroundColor $errorColor
                    }
                } catch {
                    Write-Host "Kon response body niet lezen" -ForegroundColor $errorColor
                }
                
                Write-Host "Test ${Name}: Mislukt" -ForegroundColor $errorColor
            }
        }
    }
    
    # Lokale mock response als testmodus is ingeschakeld maar server ondersteunt het niet
    if ($UseTestMode -and -not $success -and -not $SuppressOutput -and ($Endpoint -eq "/api/contact-email" -or $Endpoint -eq "/api/aanmelding-email")) {
        Write-Host "Server ondersteunt geen testmodus, simuleer respons lokaal" -ForegroundColor $promptColor
        $mockMessage = if ($Endpoint -eq "/api/contact-email") {
            "[TEST MODE] Je bericht is succesvol verwerkt (geen echte email verstuurd)"
        } else {
            "[TEST MODE] Je aanmelding is succesvol verwerkt (geen echte email verstuurd)"
        }
        
        $response = @{
            success = $true
            message = $mockMessage
            test_mode = $true
            request = $Body
        }
        
        Write-Host "Response: $($response | ConvertTo-Json -Depth 5)" -ForegroundColor $infoColor
        Write-Host "Test ${Name}: Geslaagd (lokaal gesimuleerd)" -ForegroundColor $successColor
        $success = $true
        $script:successfulTests++
    }
    
    return $response
}

# Test de root endpoint
function Test-RootEndpoint {
    Show-Title -Title "Beschikbare Endpoints"
    return Invoke-ApiCall -Name "Root Endpoint" -Method "GET" -Endpoint "/" -RetryCount 3
}

# Test de health endpoint
function Test-HealthEndpoint {
    Show-Title -Title "Health Endpoint"
    return Invoke-ApiCall -Name "Health Check" -Method "GET" -Endpoint "/api/health" -RetryCount 3
}

# Test de publieke email endpoints
function Test-PublicEmailEndpoints {
    if ($SkipEmailTests) {
        Write-Host "Email tests overgeslagen (-SkipEmailTests parameter)" -ForegroundColor $promptColor
        return $null
    }
    
    Show-Title -Title "Publieke Email Endpoints"
    
    # Contact email endpoint
    $contactResult = Invoke-ApiCall -Name "Contact Email" -Method "POST" -Endpoint "/api/contact-email" -Body $testData.ContactEmail -RetryCount 2
    
    # Aanmelding email endpoint
    $aanmeldingResult = Invoke-ApiCall -Name "Aanmelding Email" -Method "POST" -Endpoint "/api/aanmelding-email" -Body $testData.AanmeldingEmail -RetryCount 2
    
    return @{
        ContactResult = $contactResult
        AanmeldingResult = $aanmeldingResult
    }
}

# Test de metrics endpoints
function Test-MetricsEndpoints {
    if (-not $TestMetrics) {
        return $null
    }
    
    Show-Title -Title "Metrics Endpoints"
    
    # Email metrics endpoint
    $emailMetricsResult = Invoke-ApiCall -Name "Email Metrics" -Method "GET" -Endpoint "/api/metrics/email" -RequiresApiKey -ApiKey $ApiKey -RetryCount 2
    
    # Rate limit metrics endpoint
    $rateLimitMetricsResult = Invoke-ApiCall -Name "Rate Limit Metrics" -Method "GET" -Endpoint "/api/metrics/rate-limits" -RequiresApiKey -ApiKey $ApiKey -RetryCount 2
    
    # Prometheus metrics endpoint
    $prometheusMetricsResult = Invoke-ApiCall -Name "Prometheus Metrics" -Method "GET" -Endpoint "/metrics" -RequiresApiKey -ApiKey $ApiKey -RetryCount 2
    
    return @{
        EmailMetricsResult = $emailMetricsResult
        RateLimitMetricsResult = $rateLimitMetricsResult
        PrometheusMetricsResult = $prometheusMetricsResult
    }
}

# Test de beveiligde endpoints als optie is ingeschakeld
function Test-SecuredEndpoints {
    if (-not $IncludeSecuredEndpoints) {
        return $null
    }
    
    Show-Title -Title "Beveiligde Endpoints (Auth Vereist)"
    
    # Login endpoint (al getest in Invoke-ApiCall als nodig)
    $loginResult = Invoke-ApiCall -Name "Login" -Method "POST" -Endpoint "/api/auth/login" -Body $testData.Auth -RetryCount 2
    
    # Als login is gelukt, test andere endpoints
    if ($loginResult -and $loginResult.token) {
        # Profile endpoint
        $profileResult = Invoke-ApiCall -Name "User Profile" -Method "GET" -Endpoint "/api/auth/profile" -RequiresAuth -RetryCount 2
        
        # Contact beheer endpoints
        $contactsResult = Invoke-ApiCall -Name "Contactformulieren Lijst" -Method "GET" -Endpoint "/api/contact" -RequiresAuth -RetryCount 2
        
        # Aanmelding beheer endpoints
        $aanmeldingenResult = Invoke-ApiCall -Name "Aanmeldingen Lijst" -Method "GET" -Endpoint "/api/aanmelding" -RequiresAuth -RetryCount 2
        
        return @{
            ProfileResult = $profileResult
            ContactsResult = $contactsResult
            AanmeldingenResult = $aanmeldingenResult
        }
    }
    
    return $null
}

# Test de mail gerelateerde endpoints
function Test-MailEndpoints {
    if (-not $TestMailEndpoints -or -not $IncludeSecuredEndpoints) {
        if (-not $TestMailEndpoints) {
            Write-Host "Mail endpoints tests overgeslagen (-TestMailEndpoints parameter niet opgegeven)" -ForegroundColor $promptColor
        } elseif (-not $IncludeSecuredEndpoints) {
            Write-Host "Mail endpoints tests overgeslagen (vereist -IncludeSecuredEndpoints)" -ForegroundColor $promptColor
        }
        return $null
    }
    
    Show-Title -Title "Mail Endpoints (Auth Vereist)"
    
    # Test het ophalen van e-mails (GET /api/mail)
    $mailListResult = Invoke-ApiCall -Name "Inkomende e-mails overzicht" -Method "GET" -Endpoint "/api/mail" -RequiresAuth -RetryCount 2
    
    # Test het ophalen van specifieke e-mail (GET /api/mail/1)
    # Alleen als we een lijst van e-mails hebben kunnen ophalen
    $mailItemResult = $null
    if ($mailListResult -and $mailListResult.emails -and $mailListResult.emails.Count -gt 0) {
        $firstEmailId = $mailListResult.emails[0].id
        $mailItemResult = Invoke-ApiCall -Name "Specifieke e-mail ophalen" -Method "GET" -Endpoint "/api/mail/$firstEmailId" -RequiresAuth -RetryCount 2
    }
    
    # Test het actief ophalen van nieuwe e-mails (POST /api/mail/fetch)
    $fetchEmailsResult = Invoke-ApiCall -Name "E-mails ophalen van server" -Method "POST" -Endpoint "/api/mail/fetch" -RequiresAuth -RetryCount 2 -Body @{}
    
    return @{
        MailListResult = $mailListResult
        MailItemResult = $mailItemResult
        FetchEmailsResult = $fetchEmailsResult
    }
}

# Toont gedetailleerde analyse van health check
function Show-HealthAnalysis {
    param ([object]$HealthData)
    
    if (-not $HealthData -or -not $DetailedHealth) {
        return
    }
    
    Show-Title -Title "Gedetailleerde Health Analyse"
    
    # Algemene informatie
    Write-Host "API Versie: $($HealthData.version)" -ForegroundColor $infoColor
    Write-Host "Status: $($HealthData.status)" -ForegroundColor $(if ($HealthData.status -eq "healthy") { $successColor } elseif ($HealthData.status -eq "degraded") { $warningColor } else { $errorColor })
    Write-Host "Omgeving: $($HealthData.environment)" -ForegroundColor $infoColor
    Write-Host "Uptime: $($HealthData.uptime)" -ForegroundColor $infoColor
    Write-Host "Tijdstempel: $($HealthData.timestamp)" -ForegroundColor $infoColor
    
    # System resources
    if ($HealthData.system) {
        Write-Host "`nSystem Resources:" -ForegroundColor $infoColor
        Write-Host "  Go versie: $($HealthData.system.go_version)" -ForegroundColor $infoColor
        Write-Host "  CPU cores: $($HealthData.system.num_cpu)" -ForegroundColor $infoColor
        Write-Host "  Goroutines: $($HealthData.system.num_goroutines)" -ForegroundColor $infoColor
    }
    
    # Memory info
    if ($HealthData.memory) {
        Write-Host "`nGeheugen Gebruik:" -ForegroundColor $infoColor
        $allocMB = [math]::Round($HealthData.memory.heap_alloc / 1024 / 1024, 2)
        Write-Host "  Heap Geheugen: $allocMB MB" -ForegroundColor $infoColor
        Write-Host "  GC Cycles: $($HealthData.memory.num_gc)" -ForegroundColor $infoColor
    }
    
    # Component checks
    if ($HealthData.checks) {
        Write-Host "`nComponent Status:" -ForegroundColor $highlightColor
        
        # SMTP check
        if ($HealthData.checks.smtp) {
            $smtpStatus = $HealthData.checks.smtp
            $defaultSmtp = if ($smtpStatus.default) { "✅ Actief" } else { "❌ Inactief" }
            $regSmtp = if ($smtpStatus.registration) { "✅ Actief" } else { "❌ Inactief" }
            Write-Host "  SMTP Services:" -ForegroundColor $infoColor
            Write-Host "    - Default SMTP: $defaultSmtp" -ForegroundColor $(if ($smtpStatus.default) { $successColor } else { $errorColor })
            Write-Host "    - Registration SMTP: $regSmtp" -ForegroundColor $(if ($smtpStatus.registration) { $successColor } else { $errorColor })
        }
        
        # Rate limiter
        if ($HealthData.checks.rate_limiter) {
            $rateLimiter = $HealthData.checks.rate_limiter
            $rlStatus = if ($rateLimiter.status) { "✅ Actief" } else { "❌ Inactief" }
            Write-Host "  Rate Limiter: $rlStatus" -ForegroundColor $(if ($rateLimiter.status) { $successColor } else { $errorColor })
            
            if ($rateLimiter.limits) {
                Write-Host "    Geconfigureerde limieten:" -ForegroundColor $infoColor
                foreach ($limitKey in $rateLimiter.limits.PSObject.Properties.Name) {
                    $limit = $rateLimiter.limits.$limitKey
                    $scopeText = if ($limit.per_ip) { "Per IP" } else { "Globaal" }
                    $windowSeconds = [math]::Round($limit.window / 1000000000, 0)
                    Write-Host "    - $limitKey ($scopeText): $($limit.count) requests per $windowSeconds seconden" -ForegroundColor $infoColor
                }
            }
        }
        
        # Templates
        if ($HealthData.checks.templates) {
            $templates = $HealthData.checks.templates
            $templatesStatus = if ($templates.status) { "✅ Actief" } else { "❌ Inactief" }
            Write-Host "  Templates: $templatesStatus" -ForegroundColor $(if ($templates.status) { $successColor } else { $errorColor })
            
            if ($templates.available -and $templates.available.Count -gt 0) {
                Write-Host "    Beschikbare templates:" -ForegroundColor $infoColor
                foreach ($template in $templates.available) {
                    Write-Host "    - $template" -ForegroundColor $infoColor
                }
            }
        }
    }
    
    # Diagnose van "degraded" status
    if ($HealthData.status -eq "degraded") {
        Write-Host "`nDiagnose van 'degraded' status:" -ForegroundColor $warningColor
        $foundIssues = $false
        
        if ($HealthData.checks.smtp -and (-not $HealthData.checks.smtp.default -or -not $HealthData.checks.smtp.registration)) {
            Write-Host "  ⚠️ SMTP service heeft problemen - e-mail verzending kan mislukken" -ForegroundColor $warningColor
            $foundIssues = $true
        }
        
        if ($HealthData.checks.rate_limiter -and -not $HealthData.checks.rate_limiter.status) {
            Write-Host "  ⚠️ Rate limiter is inactief - dit kan een configuratie probleem zijn" -ForegroundColor $warningColor
            Write-Host "     De rate limiter wordt normaal gesproken gedisabled in ontwikkelomgevingen" -ForegroundColor $infoColor
            $foundIssues = $true
        }
        
        if ($HealthData.checks.templates -and -not $HealthData.checks.templates.status) {
            Write-Host "  ⚠️ Template systeem heeft problemen - e-mail templates kunnen niet worden geladen" -ForegroundColor $warningColor
            $foundIssues = $true
        }
        
        if (-not $foundIssues) {
            Write-Host "  ⚠️ Geen specifieke problemen gedetecteerd in health check data" -ForegroundColor $warningColor
            Write-Host "     Er kan een issue zijn met een component die niet in de health check is opgenomen" -ForegroundColor $infoColor
        }
    }
}

# Hoofdprogramma
function Start-ApiTest {
    Show-Title -Title "DKL Email Service - API Test"
    
    $testModeText = if ($testMode) {
        "ingeschakeld - Geen echte emails worden verstuurd"
    } else {
        "uitgeschakeld - Echte emails worden verstuurd"
    }
    Write-Host "Test Mode is $testModeText" -ForegroundColor $promptColor
    Write-Host "API is bereikbaar op: $currentUrl" -ForegroundColor $infoColor
    if ($DetailedHealth) {
        Write-Host "Gedetailleerde health analyse: ingeschakeld" -ForegroundColor $infoColor
    }
    if ($TestMetrics) {
        $apiKeyText = if ($ApiKey) { "aangepast" } else { "standaard" }
        Write-Host "Metrics testen: ingeschakeld (API key: $apiKeyText)" -ForegroundColor $infoColor
    }
    if ($TestMailEndpoints) {
        Write-Host "Mail endpoints testen: ingeschakeld" -ForegroundColor $infoColor
    }
    
    # Test beschikbare endpoints via root endpoint
    $rootResult = Test-RootEndpoint
    
    # Test health endpoint
    $healthResult = Test-HealthEndpoint
    
    # Test publieke email endpoints
    $emailResults = Test-PublicEmailEndpoints
    
    # Test metrics endpoints
    $metricsResults = Test-MetricsEndpoints
    
    # Test beveiligde endpoints indien ingeschakeld
    $securedResults = Test-SecuredEndpoints
    
    # Test mail endpoints indien ingeschakeld
    $mailResults = Test-MailEndpoints
    
    # Toon health check analyse indien ingeschakeld
    if ($DetailedHealth -and $healthResult) {
        Show-HealthAnalysis -HealthData $healthResult
    }
    
    # Toon samenvatting
    Show-Title -Title "Test Samenvatting"
    
    $successRate = if ($totalTests -gt 0) { ($successfulTests / $totalTests) * 100 } else { 0 }
    Write-Host "Tests uitgevoerd: $totalTests" -ForegroundColor $infoColor
    Write-Host "Tests geslaagd: $successfulTests" -ForegroundColor $successColor
    Write-Host "Succes percentage: $([math]::Round($successRate, 2))%" -ForegroundColor $(if ($successRate -ge 90) { $successColor } elseif ($successRate -ge 75) { $promptColor } else { $errorColor })
    
    $testModusText = if ($testMode) {
        "ingeschakeld"
    } else {
        "uitgeschakeld"
    }
    Write-Host "Testmodus was $testModusText" -ForegroundColor $infoColor
    
    # API status rapporteren
    if ($healthResult) {
        $apiStatusColor = if ($healthResult.status -eq "healthy") { $successColor } elseif ($healthResult.status -eq "degraded") { $warningColor } else { $errorColor }
        Write-Host "API Status: $($healthResult.status)" -ForegroundColor $apiStatusColor
        
        if ($healthResult.status -eq "degraded" -and -not $DetailedHealth) {
            Write-Host "Voer het script uit met -DetailedHealth parameter voor meer informatie over de oorzaak" -ForegroundColor $promptColor
        }
    }
    
    # Test mode werkend?
    if ($testMode -and $emailResults) {
        $testModeSupported = ($emailResults.ContactResult -and $emailResults.ContactResult.test_mode) -or 
                             ($emailResults.AanmeldingResult -and $emailResults.AanmeldingResult.test_mode)
        
        if ($testModeSupported) {
            Write-Host "Test Mode: ✅ Volledig ondersteund door de API" -ForegroundColor $successColor
        } else {
            Write-Host "Test Mode: ⚠️ Wordt lokaal gesimuleerd (niet ondersteund door API)" -ForegroundColor $warningColor
        }
    }
    
    # Beveiligde endpoints toegang
    if ($IncludeSecuredEndpoints) {
        if ($securedResults) {
            Write-Host "Beveiligde Endpoints: ✅ Authenticatie succesvol" -ForegroundColor $successColor
        } else {
            Write-Host "Beveiligde Endpoints: ❌ Authenticatie mislukt" -ForegroundColor $errorColor
            Write-Host "  Controleer gebruikersnaam en wachtwoord parameters" -ForegroundColor $promptColor
            Write-Host "  Standaard credentials: admin@dekoninklijkeloop.nl / admin123" -ForegroundColor $promptColor
            Write-Host "  Alternatieve gebruikers: jeffrey@dekoninklijkeloop.nl" -ForegroundColor $promptColor
            Write-Host "  Voorbeeld: -IncludeSecuredEndpoints -Username admin@dekoninklijkeloop.nl -Password admin123" -ForegroundColor $promptColor
        }
    }
    
    # Mail endpoints toegang
    if ($TestMailEndpoints -and $IncludeSecuredEndpoints) {
        if ($mailResults -and ($mailResults.MailListResult -or $mailResults.FetchEmailsResult)) {
            Write-Host "Mail Endpoints: ✅ Toegang verkregen" -ForegroundColor $successColor
            
            if ($mailResults.FetchEmailsResult -and $mailResults.FetchEmailsResult.count -gt 0) {
                Write-Host "  Er zijn $($mailResults.FetchEmailsResult.count) nieuwe e-mails opgehaald" -ForegroundColor $successColor
            } elseif ($mailResults.FetchEmailsResult) {
                Write-Host "  Er zijn geen nieuwe e-mails opgehaald (0)" -ForegroundColor $infoColor
            }
            
            if ($mailResults.MailListResult -and $mailResults.MailListResult.emails -and $mailResults.MailListResult.emails.Count -gt 0) {
                Write-Host "  Er zijn $($mailResults.MailListResult.emails.Count) e-mails in de database" -ForegroundColor $infoColor
            }
        } else {
            Write-Host "Mail Endpoints: ❌ Toegang mislukt of endpoints niet beschikbaar" -ForegroundColor $errorColor
        }
    }
    
    # Metrics toegang
    if ($TestMetrics) {
        if ($metricsResults -and ($metricsResults.EmailMetricsResult -or $metricsResults.RateLimitMetricsResult -or $metricsResults.PrometheusMetricsResult)) {
            Write-Host "Metrics Endpoints: ✅ Toegang verkregen" -ForegroundColor $successColor
        } else {
            Write-Host "Metrics Endpoints: ❌ Toegang geweigerd" -ForegroundColor $errorColor
            Write-Host "  Controleer API key parameter" -ForegroundColor $promptColor
            Write-Host "  Standaard API key: $defaultApiKey" -ForegroundColor $promptColor
        }
    }
    
    Write-Host "Controleer de resultaten hierboven voor details." -ForegroundColor $infoColor
}

# Start de test
Start-ApiTest

# Stop transcript als resultaten worden opgeslagen
if ($SaveResults) {
    Stop-Transcript
    Write-Host ""
    Write-Host "Resultaten opgeslagen in: $OutputFile" -ForegroundColor $infoColor
}