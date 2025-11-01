# Reprocess Existing Emails
# Dit script roept het nieuwe reprocess endpoint aan om bestaande emails opnieuw te decoderen

param(
    [string]$Environment = "production",
    [string]$Token = ""
)

$ErrorActionPreference = "Stop"

function Get-BaseUrl {
    param([string]$Env)
    
    if ($Env -eq "docker") {
        return "http://localhost:8082"
    } else {
        return "https://dklemailservice.onrender.com"
    }
}

$baseUrl = Get-BaseUrl -Env $Environment

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "   Email Reprocessing Tool" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Environment: $Environment" -ForegroundColor White
Write-Host "Base URL: $baseUrl" -ForegroundColor Blue
Write-Host "========================================`n" -ForegroundColor Cyan

# Get token if not provided
if (-not $Token) {
    Write-Host "Please provide your admin JWT token:" -ForegroundColor Yellow
    $Token = Read-Host "Token"
    
    if (-not $Token) {
        Write-Host "Error: Token is required" -ForegroundColor Red
        exit 1
    }
}

# First, check health
Write-Host "Checking service health..." -ForegroundColor Yellow
try {
    $health = Invoke-RestMethod -Uri "$baseUrl/api/health" -Method Get -ErrorAction Stop
    Write-Host "Service Status: $($health.status)" -ForegroundColor Green
} catch {
    Write-Host "Error: Cannot connect to service" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    exit 1
}

# Show warning
Write-Host "`n⚠️  WARNING ⚠️" -ForegroundColor Yellow
Write-Host "This will reprocess ALL emails in the database with improved decoding." -ForegroundColor Yellow
Write-Host "This includes:" -ForegroundColor White
Write-Host "  - Quoted-printable decoding (=92, =85, etc.)" -ForegroundColor Gray
Write-Host "  - Windows-1252 to UTF-8 conversion" -ForegroundColor Gray
Write-Host "  - MIME boundary removal" -ForegroundColor Gray
Write-Host "  - Multipart parsing improvements" -ForegroundColor Gray
Write-Host "`nDo you want to continue? (Y/N): " -NoNewline -ForegroundColor Yellow
$confirm = Read-Host

if ($confirm -ne "Y" -and $confirm -ne "y") {
    Write-Host "Operation cancelled" -ForegroundColor Yellow
    exit 0
}

# Call reprocess endpoint
Write-Host "`nStarting email reprocessing..." -ForegroundColor Cyan

try {
    $headers = @{
        'Authorization' = "Bearer $Token"
        'Content-Type' = 'application/json'
    }
    
    $response = Invoke-RestMethod -Uri "$baseUrl/api/admin/mail/reprocess" `
        -Method Post `
        -Headers $headers `
        -ErrorAction Stop
    
    Write-Host "`n========================================" -ForegroundColor Green
    Write-Host "   Reprocessing Complete!" -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Green
    Write-Host "Processed: $($response.processed) emails" -ForegroundColor White
    Write-Host "Failed: $($response.failed) emails" -ForegroundColor $(if($response.failed -gt 0){"Yellow"}else{"Green"})
    Write-Host "Message: $($response.message)" -ForegroundColor Gray
    Write-Host "========================================`n" -ForegroundColor Green
    
    if ($response.processed -gt 0) {
        Write-Host "✓ Emails have been reprocessed with improved decoding" -ForegroundColor Green
        Write-Host "✓ Quoted-printable artifacts (=92, =85) are now fixed" -ForegroundColor Green
        Write-Host "✓ MIME boundaries are removed" -ForegroundColor Green
        Write-Host "✓ Windows-1252 characters converted to UTF-8" -ForegroundColor Green
    }
    
} catch {
    Write-Host "`nError during reprocessing:" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    
    if ($_.Exception.Response.StatusCode -eq 'Unauthorized') {
        Write-Host "`nYour token is invalid or expired." -ForegroundColor Yellow
        Write-Host "Please login again to get a new token." -ForegroundColor Yellow
    } elseif ($_.Exception.Response.StatusCode -eq 'Forbidden') {
        Write-Host "`nYou don't have permission to reprocess emails." -ForegroundColor Yellow
        Write-Host "This operation requires 'admin_email:send' permission." -ForegroundColor Yellow
    }
    
    exit 1
}

Write-Host "`nReprocessing complete! Check your emails in the admin panel.`n" -ForegroundColor Green