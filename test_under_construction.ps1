# Test Under Construction API - Docker & Production
# Dit script test de under construction/maintenance mode status

param(
    [string]$Environment = "production"
)

$ErrorActionPreference = "Continue"

function Test-UnderConstructionEndpoint {
    param(
        [string]$BaseUrl,
        [string]$EnvName
    )
    
    Write-Host "`n========================================" -ForegroundColor Cyan
    Write-Host "Testing Under Construction in $EnvName" -ForegroundColor Cyan
    Write-Host "Base URL: $BaseUrl" -ForegroundColor Cyan
    Write-Host "========================================`n" -ForegroundColor Cyan
    
    # Test 1: Health Check
    Write-Host "1. Health Check..." -ForegroundColor Yellow
    try {
        $health = Invoke-RestMethod -Uri "$BaseUrl/api/health" -Method Get -ErrorAction Stop
        Write-Host "SUCCESS: Health check successful" -ForegroundColor Green
        Write-Host "  Status: $($health.status)" -ForegroundColor Gray
    } catch {
        Write-Host "FAILED: Health check failed: $($_.Exception.Message)" -ForegroundColor Red
        return
    }
    
    # Test 2: Get Active Under Construction
    Write-Host "`n2. Get Active Under Construction Status..." -ForegroundColor Yellow
    try {
        $uc = Invoke-RestMethod -Uri "$BaseUrl/api/under-construction/active" -Method Get -ErrorAction Stop
        Write-Host "SUCCESS: Maintenance mode is ACTIVE" -ForegroundColor Yellow
        Write-Host "`n  Maintenance Mode Details:" -ForegroundColor Cyan
        Write-Host "    ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor DarkGray
        Write-Host "    Title: $($uc.title)" -ForegroundColor White
        Write-Host "    Message: $($uc.message)" -ForegroundColor Gray
        if ($uc.footer_text) {
            Write-Host "    Footer: $($uc.footer_text)" -ForegroundColor Gray
        }
        if ($uc.logo_url) {
            Write-Host "    Logo URL: $($uc.logo_url)" -ForegroundColor Blue
        }
        if ($uc.expected_date) {
            Write-Host "    Expected Date: $($uc.expected_date)" -ForegroundColor Green
        }
        if ($uc.progress_percentage) {
            Write-Host "    Progress: $($uc.progress_percentage)%" -ForegroundColor Cyan
        }
        if ($uc.contact_email) {
            Write-Host "    Contact: $($uc.contact_email)" -ForegroundColor Gray
        }
        Write-Host "    Newsletter: $(if($uc.newsletter_enabled){'Enabled'}else{'Disabled'})" -ForegroundColor Gray
        
        if ($uc.social_links) {
            Write-Host "`n    Social Links:" -ForegroundColor Cyan
            foreach ($link in $uc.social_links) {
                Write-Host "      - $($link.platform): $($link.url)" -ForegroundColor DarkCyan
            }
        }
        Write-Host "    ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor DarkGray
        
    } catch {
        if ($_.Exception.Response.StatusCode -eq 404) {
            Write-Host "SUCCESS: No maintenance mode active (site is operational)" -ForegroundColor Green
            Write-Host "  The website is available for normal use" -ForegroundColor Gray
        } else {
            Write-Host "FAILED: $($_.Exception.Message)" -ForegroundColor Red
            if ($_.ErrorDetails) {
                Write-Host "  Details: $($_.ErrorDetails.Message)" -ForegroundColor Red
            }
        }
    }
    
    Write-Host "`n========================================" -ForegroundColor Cyan
    Write-Host "Completed testing $EnvName" -ForegroundColor Cyan
    Write-Host "========================================`n" -ForegroundColor Cyan
}

# Main execution
switch ($Environment.ToLower()) {
    "docker" {
        Test-UnderConstructionEndpoint -BaseUrl "http://localhost:8082" -EnvName "Docker (Dev)"
    }
    "production" {
        Test-UnderConstructionEndpoint -BaseUrl "https://dklemailservice.onrender.com" -EnvName "Production"
    }
    "both" {
        Test-UnderConstructionEndpoint -BaseUrl "http://localhost:8082" -EnvName "Docker (Dev)"
        Test-UnderConstructionEndpoint -BaseUrl "https://dklemailservice.onrender.com" -EnvName "Production"
    }
    default {
        Write-Host "Invalid environment. Use: docker, production, or both" -ForegroundColor Red
        Write-Host "Example: .\test_under_construction.ps1 -Environment production" -ForegroundColor Yellow
    }
}

Write-Host "`nTest completed!" -ForegroundColor Green