# DKL Email Service Mail API Test Script
# Dit script test specifiek de mail endpoints van de DKL Email Service API

# Kleuren voor output
$successColor = "Green"
$errorColor = "Red"
$infoColor = "Cyan"
$promptColor = "Yellow"
$highlightColor = "Magenta"

# Configuratie
$baseUrl = "https://dklemailservice.onrender.com"
$localUrl = "http://localhost:8080"
$currentUrl = $localUrl  # Default lokaal testen

# Admin gebruiker voor authenticatie
$adminCredentials = @{
    email = "admin@example.com"
    wachtwoord = "admin123"
}

# Globale variabelen
$global:isLoggedIn = $false
$global:session = $null
$global:jwtToken = $null
$global:mailId = $null

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
        [Microsoft.PowerShell.Commands.WebRequestSession]$Session = $null
    )
    
    Write-Host "" -ForegroundColor $infoColor
    Write-Host "[TESTING] $Name..." -ForegroundColor $infoColor
    
    try {
        $params = @{
            Uri = "$currentUrl$Endpoint"
            Method = $Method
            UseBasicParsing = $true
        }
        
        if ($Body) {
            $jsonBody = $Body | ConvertTo-Json
            Write-Host "Request Body:" -ForegroundColor $infoColor
            Write-Host $jsonBody
            $params.Body = $jsonBody
            $params.ContentType = "application/json"
        }
        
        # Initialiseer headers als ze nog niet bestaan
        if (-not $params.Headers) {
            $params.Headers = @{}
        }
        
        if ($UseAuth -and $global:jwtToken) {
            $params.Headers["Authorization"] = "Bearer $($global:jwtToken)"
        }
        
        $response = Invoke-RestMethod @params
        
        Show-Result -Name $Name -Response $response
        return $response
    }
    catch {
        Write-Host "" -ForegroundColor $infoColor
        Write-Host "[ERROR] $Name failed" -ForegroundColor $errorColor
        
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
        
        return $null
    }
}

# Functie om in te loggen als admin
function Test-AdminLogin {
    Show-Title -Title "Admin Login"
    
    $response = Invoke-ApiCall -Name "Admin Login" -Method "Post" -Endpoint "/api/auth/login" -Body $adminCredentials
    
    if ($response -and $response.token) {
        $global:jwtToken = $response.token
        $global:isLoggedIn = $true
        Write-Host "JWT Token opgeslagen voor gebruik in requests" -ForegroundColor $infoColor
    } else {
        $global:isLoggedIn = $false
    }
    
    return $response
}

# Functie om de mail lijst endpoint te testen
function Test-MailList {
    Show-Title -Title "Mail Lijst Endpoint Testen"
    
    if (-not $global:isLoggedIn) {
        Write-Host "Je moet eerst inloggen als admin om deze endpoint te testen." -ForegroundColor $errorColor
        return $null
    }
    
    $response = Invoke-ApiCall -Name "Mail Lijst" -Method "Get" -Endpoint "/api/mail" -UseAuth
    
    # Sla het eerste mail ID op voor gebruik in andere tests
    if ($response -and $response.Count -gt 0) {
        $global:mailId = $response[0].id
        Write-Host "Mail ID opgeslagen voor gebruik in andere tests: $($global:mailId)" -ForegroundColor $infoColor
    } else {
        Write-Host "Geen e-mails gevonden om te testen." -ForegroundColor $promptColor
    }
    
    return $response
}

# Functie om de mail details endpoint te testen
function Test-MailDetails {
    Show-Title -Title "Mail Details Endpoint Testen"
    
    if (-not $global:isLoggedIn) {
        Write-Host "Je moet eerst inloggen als admin om deze endpoint te testen." -ForegroundColor $errorColor
        return $null
    }
    
    if (-not $global:mailId) {
        Write-Host "Geen mail ID beschikbaar. Voer eerst de Mail Lijst test uit." -ForegroundColor $promptColor
        return $null
    }
    
    $response = Invoke-ApiCall -Name "Mail Details" -Method "Get" -Endpoint "/api/mail/$($global:mailId)" -UseAuth
    
    return $response
}

# Functie om de unprocessed mails endpoint te testen
function Test-UnprocessedMails {
    Show-Title -Title "Onverwerkte Mail Lijst Endpoint Testen"
    
    if (-not $global:isLoggedIn) {
        Write-Host "Je moet eerst inloggen als admin om deze endpoint te testen." -ForegroundColor $errorColor
        return $null
    }
    
    $response = Invoke-ApiCall -Name "Onverwerkte Mails" -Method "Get" -Endpoint "/api/mail/unprocessed" -UseAuth
    
    return $response
}

# Functie om de mail fetch endpoint te testen
function Test-FetchMails {
    Show-Title -Title "Mail Ophalen Endpoint Testen"
    
    if (-not $global:isLoggedIn) {
        Write-Host "Je moet eerst inloggen als admin om deze endpoint te testen." -ForegroundColor $errorColor
        return $null
    }
    
    $response = Invoke-ApiCall -Name "Mail Ophalen" -Method "Post" -Endpoint "/api/mail/fetch" -UseAuth
    
    return $response
}

# Functie om de mark as processed endpoint te testen
function Test-MarkMailAsProcessed {
    Show-Title -Title "Mail Als Verwerkt Markeren Endpoint Testen"
    
    if (-not $global:isLoggedIn) {
        Write-Host "Je moet eerst inloggen als admin om deze endpoint te testen." -ForegroundColor $errorColor
        return $null
    }
    
    if (-not $global:mailId) {
        Write-Host "Geen mail ID beschikbaar. Voer eerst de Mail Lijst test uit." -ForegroundColor $promptColor
        return $null
    }
    
    $response = Invoke-ApiCall -Name "Mail Als Verwerkt Markeren" -Method "Put" -Endpoint "/api/mail/$($global:mailId)/processed" -UseAuth
    
    return $response
}

# Functie om mails per account type te testen
function Test-MailsByAccountType {
    Show-Title -Title "Mails Per Account Type Endpoint Testen"
    
    if (-not $global:isLoggedIn) {
        Write-Host "Je moet eerst inloggen als admin om deze endpoint te testen." -ForegroundColor $errorColor
        return $null
    }
    
    # Test met account type "info"
    $accountType = "info"
    
    $response = Invoke-ApiCall -Name "Mails Per Account Type" -Method "Get" -Endpoint "/api/mail/account/$accountType" -UseAuth
    
    return $response
}

# Functie om de delete mail endpoint te testen
function Test-DeleteMail {
    Show-Title -Title "Mail Verwijderen Endpoint Testen"
    
    if (-not $global:isLoggedIn) {
        Write-Host "Je moet eerst inloggen als admin om deze endpoint te testen." -ForegroundColor $errorColor
        return $null
    }
    
    if (-not $global:mailId) {
        Write-Host "Geen mail ID beschikbaar. Voer eerst de Mail Lijst test uit." -ForegroundColor $promptColor
        return $null
    }
    
    $response = Invoke-ApiCall -Name "Mail Verwijderen" -Method "Delete" -Endpoint "/api/mail/$($global:mailId)" -UseAuth
    
    if ($response) {
        # Als het verwijderen succesvol was, reset dan de mail ID
        $global:mailId = $null
        Write-Host "Mail ID gereset na succesvolle verwijdering." -ForegroundColor $infoColor
    }
    
    return $response
}

# Functie om uit te loggen
function Test-Logout {
    Show-Title -Title "Logout"
    
    if (-not $global:isLoggedIn) {
        Write-Host "Je bent niet ingelogd." -ForegroundColor $errorColor
        return $null
    }
    
    $response = Invoke-ApiCall -Name "Logout" -Method "Post" -Endpoint "/api/auth/logout" -UseAuth
    
    if ($response) {
        $global:isLoggedIn = $false
        $global:jwtToken = $null
        Write-Host "Succesvol uitgelogd!" -ForegroundColor $successColor
    }
    
    return $response
}

# Hoofdfunctie om alle mail endpoints te testen
function Test-AllMailEndpoints {
    Show-Title -Title "Alle Mail API Endpoints Testen"
    
    # Stap 1: Login
    $loginResult = Test-AdminLogin
    
    if (-not $global:isLoggedIn) {
        Write-Host "Login mislukt. Tests worden gestopt." -ForegroundColor $errorColor
        return
    }
    
    # Stap 2: Haal nieuwe e-mails op
    $fetchResult = Test-FetchMails
    
    # Wacht even om zeker te zijn dat de mails zijn opgehaald
    Write-Host "Even wachten om zeker te zijn dat de e-mails zijn opgehaald..." -ForegroundColor $infoColor
    Start-Sleep -Seconds 5
    
    # Stap 3: Test alle mail endpoints
    $mailListResult = Test-MailList
    $unprocessedMailsResult = Test-UnprocessedMails
    $mailsByAccountTypeResult = Test-MailsByAccountType
    
    if ($global:mailId) {
        $mailDetailsResult = Test-MailDetails
        $markAsProcessedResult = Test-MarkMailAsProcessed
        # Test verwijderen als laatste
        $deleteMailResult = Test-DeleteMail
    } else {
        Write-Host "Geen mail ID gevonden, sla gerelateerde tests over." -ForegroundColor $promptColor
    }
    
    # Stap 4: Logout
    $logoutResult = Test-Logout
    
    # Toon samenvatting
    Show-Title -Title "Test Samenvatting"
    Write-Host "Alle mail API tests zijn uitgevoerd!" -ForegroundColor $successColor
}

# Start alle tests
Test-AllMailEndpoints 