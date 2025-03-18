# Script om mail endpoint logs te bekijken
$token = ""
$apiUrl = "https://dklemailservice.onrender.com"

# Eerst inloggen om token te krijgen
Write-Host "Inloggen om token te krijgen..."
$loginBody = @{
    email = "admin@dekoninklijkeloop.nl"
    wachtwoord = "admin123"
} | ConvertTo-Json

try {
    $loginResponse = Invoke-RestMethod -Uri "$apiUrl/api/auth/login" -Method Post -Body $loginBody -ContentType "application/json"
    $token = $loginResponse.token
    Write-Host "Succesvol ingelogd. Token ontvangen."
} catch {
    Write-Host "Inloggen mislukt: $_"
    exit
}

# Nu de mail endpoint testen
Write-Host "`nFetching e-mails via /api/mail/fetch endpoint..."
try {
    $headers = @{
        Authorization = "Bearer $token"
    }
    
    $fetchResponse = Invoke-RestMethod -Uri "$apiUrl/api/mail/fetch" -Method Post -Headers $headers -Body "{}" -ContentType "application/json"
    Write-Host "Response:" -ForegroundColor Green
    $fetchResponse | ConvertTo-Json -Depth 10
} catch {
    Write-Host "Fetch emails fout: $_" -ForegroundColor Red
    try {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $reader.BaseStream.Position = 0
        $reader.DiscardBufferedData()
        $responseBody = $reader.ReadToEnd()
        Write-Host "Response Body: $responseBody" -ForegroundColor Red
    } catch {}
}

# Nu de lijst endpoint testen
Write-Host "`nLijst van e-mails ophalen via /api/mail endpoint..."
try {
    $headers = @{
        Authorization = "Bearer $token"
    }
    
    $listResponse = Invoke-RestMethod -Uri "$apiUrl/api/mail" -Method Get -Headers $headers
    Write-Host "Response:" -ForegroundColor Green
    $listResponse | ConvertTo-Json -Depth 10
} catch {
    Write-Host "List emails fout: $_" -ForegroundColor Red
    try {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $reader.BaseStream.Position = 0
        $reader.DiscardBufferedData()
        $responseBody = $reader.ReadToEnd()
        Write-Host "Response Body: $responseBody" -ForegroundColor Red
    } catch {}
} 